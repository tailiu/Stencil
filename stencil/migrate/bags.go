package migrate

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/helper"
	"stencil/reference_resolution"
	"strings"

	"github.com/gookit/color"
)

func (self *MigrationWorkerV2) ConstructBagNode(bagMember, bagMemberID, bagRowID string) (*DependencyNode, error) {

	bagTag, err := self.SrcAppConfig.GetTagByMember(bagMember)
	if err != nil {
		self.Logger.Fatal(fmt.Sprintf("@MigrateBags: UNABLE TO GET BAG TAG BY MEMBER : %s | %s", bagMember, err))
		return nil, err
	}

	if ql := self.GetTagQLForBag(*bagTag); len(ql) > 0 {
		where := fmt.Sprintf(" WHERE %s.member = %s AND %s.id = %s", bagMember, bagMemberID, bagMember, bagRowID)
		sql := ql + where
		if data, err := db.DataCall1(self.logTxn.DBconn, sql); err == nil && len(data) > 0 {
			bagData := make(map[string]interface{})
			if err := json.Unmarshal(data["json_data"].([]byte), &bagData); err != nil {
				fmt.Println("@ConstructBagNode >> ", data)
				self.Logger.Fatal(fmt.Sprintf("@ConstructBagNode: UNABLE TO CONVERT BAG TO MAP: %s", err))
				return nil, err
			}
			var bagPKs []int64
			if err := json.Unmarshal(data["pks_json"].([]byte), &bagPKs); err != nil {
				fmt.Println("@ConstructBagNode >> ", data)
				self.Logger.Fatal(fmt.Sprintf("@ConstructBagNode: UNABLE TO CONVERT pks_json TO []int64: %s", err))
				return nil, err
			}
			bagNode := &DependencyNode{Tag: *bagTag, SQL: sql, Data: bagData, PKs: bagPKs}
			return bagNode, nil
		} else {
			if err == nil {
				err = errors.New("no data returned for root node, doesn't exist?")
			} else {
				fmt.Println("@ConstructBagNode > DataCall1 | ", err)
			}
			fmt.Println(sql)
			self.Logger.Fatal("@ConstructBagNode: ", err)
			return nil, err
		}
	} else {
		self.Logger.Fatal("@ConstructBagNode > GetTagQLForBag : Failed to get tag ql | ", ql)
	}

	self.Logger.Fatal("@ConstructBagNode: End")

	return nil, nil
}

func (self *MigrationWorkerV2) GetRowsFromIDTable(app, member, id interface{}, getFrom bool) ([]IDRow, error) {

	var idRows []IDRow
	var err error
	var idRowsDB []map[string]interface{}

	if !getFrom {
		idRowsDB, err = db.GetRowsFromIDTableByTo(self.logTxn.DBconn, app, member, id)
	} else {
		idRowsDB, err = db.GetRowsFromIDTableByFrom(self.logTxn.DBconn, app, member, id)
	}

	if err != nil {
		self.Logger.Fatalf("Unable to get bags | %s | getFrom:%v uid:%v app:%v member:%v id:%v txnID:%v", err, getFrom, self.uid, app, member, id, self.logTxn.Txn_id)
		return nil, err
	}

	for _, idRowDB := range idRowsDB {
		fromAppID := fmt.Sprint(idRowDB["from_app"])
		fromAppName, err := db.GetAppNameByAppID(self.logTxn.DBconn, fromAppID)
		if err != nil {
			self.Logger.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get name | ", fromAppID, err)
			return nil, err
		}

		toAppID := fmt.Sprint(idRowDB["to_app"])
		toAppName, err := db.GetAppNameByAppID(self.logTxn.DBconn, toAppID)
		if err != nil {
			self.Logger.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get name | ", toAppID, err)
			return nil, err
		}

		fromMemberID := fmt.Sprint(idRowDB["from_member"])
		fromMember, err := db.TableName(self.logTxn.DBconn, fromMemberID, fromAppID)
		if err != nil {
			self.Logger.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get fromMember | ", fromMemberID, fromAppID, err)
			return nil, err
		}

		toMemberID := fmt.Sprint(idRowDB["to_member"])
		toMember, err := db.TableName(self.logTxn.DBconn, toMemberID, toAppID)
		if err != nil {
			self.Logger.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get toMember | ", toMemberID, toAppID, err)
			return nil, err
		}

		idRows = append(idRows, IDRow{
			FromAppName:  fromAppName,
			FromAppID:    fromAppID,
			FromMemberID: fromMemberID,
			FromMember:   fromMember,
			FromID:       idRowDB["from_id"],
			ToAppName:    toAppName,
			ToAppID:      toAppID,
			ToMember:     toMember,
			ToMemberID:   toMemberID,
			ToID:         idRowDB["to_id"]})
	}
	return idRows, nil
}

func (self *MigrationWorkerV2) MergeBagDataWithMappedData(mappedData *MappedData, node *DependencyNode, toTable config.ToTable) error {

	toTableData := make(map[string]interface{})

	prevUIDs := reference_resolution.GetPrevUserIDs(self.logTxn.DBconn, self.SrcAppConfig.AppID, self.uid)
	if prevUIDs == nil {
		prevUIDs = make(map[string]string)
	}
	prevUIDs[self.SrcAppConfig.AppID] = self.uid

	for fromTable := range mappedData.srcTables {
		if fromTableID, err := db.TableID(self.logTxn.DBconn, fromTable, self.SrcAppConfig.AppID); err == nil {
			if fromID, ok := node.Data[fromTable+".id"]; ok {
				visitedRows := make(map[string]bool)
				if err := self.FetchDataFromBags(visitedRows, toTableData, prevUIDs, self.SrcAppConfig.AppID, fromTableID, fromID, toTable.TableID, toTable.Table, toTable.Table); err != nil {
					self.Logger.Fatal("@MigrateNode > FetchDataFromBags | ", err)
				}
			} else {
				self.Logger.Fatal("@MigrateNode > FetchDataFromBags > id doesn't exist in table ", fromTable, err)
			}
		} else {
			self.Logger.Fatal("@MigrateNode > FetchDataFromBags > TableID, fromTable: error in getting table id for member! ", fromTable, err)
		}
	}

	if len(toTableData) > 0 {

		mappedCols := strings.Split(mappedData.cols, ",")
		for col, val := range toTableData {
			if !helper.Contains(mappedCols, col) {
				mappedData.cols += "," + col
				mappedData.ivals = append(mappedData.ivals, val)
				mappedData.vals += fmt.Sprintf(",$%d", len(mappedData.ivals))
			}
		}
		mappedData.Trim(",")
		self.Logger.Trace("@MigrateNode > FetchDataFromBags > Data merged for: ", toTable.Table)
	}

	return nil
}

func (self *MigrationWorkerV2) FetchDataFromBags(visitedRows map[string]bool, toTableData map[string]interface{}, prevUIDs map[string]string, app, member string, id interface{}, dstMemberID, dstMemberName, toTableName string) error {

	currentRow := fmt.Sprintf("%s:%s:%s", app, member, id)

	if _, ok := visitedRows[currentRow]; ok {
		return nil
	} else {
		visitedRows[currentRow] = true
	}

	idRows, err := self.GetRowsFromIDTable(app, member, id, false)

	if err != nil {
		self.Logger.Fatal("@FetchDataFromBags > GetRowsFromIDTable, Unable to get IDRows | ", app, member, id, false, err)
		return err
	} else if len(idRows) < 1 {
		// self.Logger.Trace("@FetchDataFromBags > GetRowsFromIDTable | No IDRows found | ", app, member, id, false)
		return nil
	} else {
		self.Logger.Trace("@FetchDataFromBags > GetRowsFromIDTable | ", idRows)
	}

	for _, idRow := range idRows {

		bagRow, err := db.GetBagByAppMemberIDV2(self.logTxn.DBconn, prevUIDs[idRow.FromAppID], idRow.FromAppID, idRow.FromMemberID, idRow.FromID, self.logTxn.Txn_id)
		if err != nil {
			self.Logger.Fatal("@FetchDataFromBags > GetBagByAppMemberIDV2, Unable to get bags | ", prevUIDs[idRow.FromAppID], idRow.FromAppID, idRow.FromMemberID, idRow.FromID, self.logTxn.Txn_id, err)
			return err
		}
		if bagRow != nil {

			bagData := make(map[string]interface{})
			if err := json.Unmarshal(bagRow["data"].([]byte), &bagData); err != nil {
				self.Logger.Debug(bagRow["data"])
				self.Logger.Debug(bagRow)
				self.Logger.Fatal("@FetchDataFromBags: UNABLE TO CONVERT BAG TO MAP ", bagRow, err)
				return err
			}

			if bagMappedApp, mapping, found := self.FetchMappingsForBag(idRow.FromAppName, idRow.FromAppID, self.DstAppConfig.AppName, self.DstAppConfig.AppID, idRow.FromMember, dstMemberName); found {
				self.Logger.Trace("@FetchDataFromBags > FetchMappingsForBag, Mappings found for | ", idRow.FromAppName, idRow.FromAppID, self.DstAppConfig.AppName, self.DstAppConfig.AppID, idRow.FromMember, dstMemberName)
				for _, toTable := range mapping.ToTables {
					if !strings.EqualFold(toTable.Table, toTableName) {
						continue
					}
					for toAttr, fromAttr := range toTable.Mapping {
						cleanedFromAttr := fromAttr
						if _, ok := toTableData[toAttr]; !ok {
							if fromAttr[0:1] == "$" {
								if inputVal, err := bagMappedApp.GetInput(fromAttr); err == nil {
									toTableData[toAttr] = inputVal
								} else {
									self.Logger.Debugf("@FetchDataFromBags | fromAttr [%s]", fromAttr)
									self.Logger.Fatal(err)
								}
							} else if bagVal, _, decodedFromAttr, _, found, err := self.DecodeMappingValue(fromAttr, bagData, true); err == nil {
								if found && bagVal != nil {
									toTableData[toAttr] = bagVal
								}
								cleanedFromAttr = decodedFromAttr
								self.Logger.Tracef("@FetchDataFromBags > DecodeMappingValue | Added | toTable: [%s], fromAttr: [%s], toAttr: [%s], BagVal: [%v], Found: [%v]", toTable.Table, fromAttr, toAttr, bagVal, found)
							} else {
								self.Logger.Debug(bagData)
								self.Logger.Fatalf("Unable to decode mapped val | fromAttr: [%s], toAttr: [%s]", fromAttr, toAttr)
							}
						} else {
							self.Logger.Tracef("@FetchDataFromBags > DecodeMappingValue | Exists | toTable: [%s], fromAttr: [%s], toAttr: [%s], BagVal: [%v]", toTable.Table, fromAttr, toAttr, toTableData[toAttr])
						}
						self.Logger.Tracef("@FetchDataFromBags > %s | [%s] | [%v]", color.FgMagenta.Render("Deleting Attr From Bag"), color.FgMagenta.Render(cleanedFromAttr), color.FgMagenta.Render(bagData))
						delete(bagData, cleanedFromAttr)
					}
				}

				if self.IsNodeDataEmpty(bagData) {
					log.Println(fmt.Sprintf("%s | PK: %v", color.FgLightRed.Render("Deleting BAG"), fmt.Sprint(bagRow["pk"])))
					if err := db.DeleteBagV2(self.tx.StencilTx, fmt.Sprint(bagRow["pk"])); err != nil {
						self.Logger.Fatal("@FetchDataFromBags > DeleteBagV2, Unable to delete bag | ", bagRow["pk"])
						return err
					} else {
						log.Println(fmt.Sprintf("%s | PK: %v", color.FgLightRed.Render("Deleted BAG"), bagData["pk"]))
					}
				} else {
					log.Println(fmt.Sprintf("%s | PK: %v", color.FgYellow.Render("BAG NOT EMPTY"), bagData))
					if jsonData, err := json.Marshal(bagData); err == nil {
						if err := db.UpdateBag(self.tx.StencilTx, fmt.Sprint(bagRow["pk"]), self.logTxn.Txn_id, jsonData); err != nil {
							self.Logger.Fatal("@FetchDataFromBags: UNABLE TO UPDATE BAG ", bagRow, err)
							return err
						}
					} else {
						self.Logger.Fatal("@FetchDataFromBags > len(bagData) != 0, Unable to marshall bag | ", bagData)
						return err
					}
				}
			} else {
				self.Logger.Warnf("No mappings found from [%s:%s] to [%s:%s]", idRow.FromAppName, idRow.FromMember, self.DstAppConfig.AppName, dstMemberName)
			}
		} else {
			log.Println("@FetchDataFromBags > GetBagByAppMemberIDV2, No bags found for | ", prevUIDs[idRow.FromAppID], idRow.FromAppID, idRow.FromMember, idRow.FromID, self.logTxn.Txn_id)
		}

		self.Logger.Tracef("@FetchDataFromBags > FetchDataFromBags: %s |\nFromAppID: %s FromMemberID: %s FromID: %v dstMemberID: %s dstMemberName: %s \ntoTableData: %v\nprevUIDs: %v\n ",
			color.FgLightMagenta.Render("Recursive Traversal"), idRow.FromAppID, idRow.FromMemberID, idRow.FromID, dstMemberID, dstMemberName,
			toTableData, prevUIDs)
		if err := self.FetchDataFromBags(visitedRows, toTableData, prevUIDs, idRow.FromAppID, idRow.FromMemberID, idRow.FromID, dstMemberID, dstMemberName, toTableName); err != nil {
			self.Logger.Fatal("@FetchDataFromBags > FetchDataFromBags: Error while recursing | ", toTableData, idRow.FromAppID, idRow.FromMember, idRow.FromID)
			return err
		}
	}
	return nil
}

func (self *MigrationWorkerV2) SendMemberToBag(node *DependencyNode, member, ownerID string) error {
	if self.mtype != DELETION && self.mtype != BAGS {
		return nil
	}
	if memberID, err := db.TableID(self.logTxn.DBconn, member, self.SrcAppConfig.AppID); err != nil {
		self.Logger.Fatal("@SendMemberToBag > TableID: error in getting table id for member! ", member, err)
		return err
	} else {
		bagData := self.GetMemberDataFromNode(member, node.Data)
		if len(bagData) > 0 {
			if srcID, ok := node.Data[member+".id"]; ok && srcID != nil {
				if jsonData, err := json.Marshal(bagData); err == nil {
					if err := db.CreateNewBag(self.tx.StencilTx, self.SrcAppConfig.AppID, memberID, srcID, ownerID, fmt.Sprint(self.logTxn.Txn_id), jsonData); err != nil {
						fmt.Println(self.SrcAppConfig.AppName, member, srcID, ownerID)
						fmt.Println(bagData)
						self.Logger.Fatal("@SendMemberToBag: error in creating bag! ", err)
						return err
					}
					if self.mtype == BAGS {

					} else {
						if derr := db.ReallyDeleteRowFromAppDB(self.tx.SrcTx, member, srcID); derr != nil {
							fmt.Println("@SendMemberToBag > DeleteRowFromAppDB")
							self.Logger.Fatal(derr)
							return derr
						}
					}
				} else {
					fmt.Println(bagData)
					self.Logger.Fatal("@SendMemberToBag: unable to convert bag data to JSON ", err)
					return err
				}
			} else {
				// fmt.Println(node.Data)
				// fmt.Println(node.SQL)
				log.Println("@SendMemberToBag: '", member, "' doesn't contain id! ", srcID)
				// return err
			}
		}
	}

	return nil
}

func (self *MigrationWorkerV2) SendNodeToBagWithOwnerID(node *DependencyNode, ownerID string) error {
	if self.mtype != DELETION && self.mtype != BAGS {
		return nil
	}
	for _, member := range node.Tag.Members {
		if err := self.SendMemberToBag(node, member, ownerID); err != nil {
			fmt.Println(node)
			self.Logger.Fatal("@SendNodeToBagWithOwnerID > SendMemberToBag: ownerID error! ")
			return err
		} else {
			log.Println(fmt.Sprintf("%s { %s :  %s } | Owner ID: %v ", color.FgYellow.Render("BAG"), node.Tag.Name, member, ownerID))
		}
	}
	return nil
}

func (self *MigrationWorkerV2) SendNodeToBag(node *DependencyNode) error {
	if self.mtype != DELETION {
		return nil
	}
	if ownerID, _ := self.GetNodeOwner(node); len(ownerID) > 0 {
		if err := self.SendNodeToBagWithOwnerID(node, ownerID); err != nil {
			return err
		}
	} else {
		fmt.Println(node)
		self.Logger.Fatal("@SendNodeToBag > GetNodeOwner: ownerID error! ")
	}

	return nil
}

func (bagWorker *MigrationWorkerV2) DeleteBag(bagNode *DependencyNode) error {
	for _, pk := range bagNode.PKs {
		bagPK := fmt.Sprint(pk)
		if err := db.DeleteBagV2(bagWorker.tx.StencilTx, bagPK); err != nil {
			fmt.Println(bagNode)
			bagWorker.Logger.Fatal(fmt.Sprintf("@DeleteBag: UNABLE TO DELETE BAG : %s | %s ", bagPK, err))
			return err
		} else {
			log.Println(fmt.Sprintf("%s { Tag: %s } | PK: %s", color.FgLightRed.Render("Deleted BAG"), bagNode.Tag.Name, bagPK))
		}
	}
	return nil
}

func (self *MigrationWorkerV2) UpdateProcessedBagIDs(bagNode *DependencyNode, processedBags map[string]bool) {
	for _, pk := range bagNode.PKs {
		pkStr := fmt.Sprint(pk)
		processedBags[pkStr] = true
	}
}

func (self *MigrationWorkerV2) MigrateBags(threadID int, isBlade ...bool) error {

	prevIDs := reference_resolution.GetPrevUserIDs(self.logTxn.DBconn, self.SrcAppConfig.AppID, self.uid)
	if prevIDs == nil {
		prevIDs = make(map[string]string)
	}

	prevIDs[self.SrcAppConfig.AppID] = self.uid

	for bagAppID, userID := range prevIDs {

		// if bagAppID != "2" {
		// 	continue
		// }

		fmt.Println("########################################################################")
		log.Println(fmt.Sprintf("Starting Bags for User: %s App: %s", userID, bagAppID))

		bags, err := db.GetBagsV2(self.logTxn.DBconn, bagAppID, userID, self.logTxn.Txn_id)

		if err != nil {
			self.Logger.Fatal(fmt.Sprintf("UNABLE TO FETCH BAGS FOR USER: %s | %s", userID, err))
			return err
		}

		bagWorker := CreateBagWorkerV2(userID, bagAppID, self.DstAppConfig.AppID, self.logTxn, BAGS, threadID, isBlade...)

		log.Println(fmt.Sprintf("Bag Worker Created | %s -> %s ", bagWorker.SrcAppConfig.AppName, bagWorker.DstAppConfig.AppName))

		processedBags := make(map[string]bool)

		for _, bag := range bags {

			bagPK := fmt.Sprint(bag["pk"])

			if _, ok := processedBags[bagPK]; ok {
				continue
			}

			bagRowID := fmt.Sprint(bag["id"])
			bagUserID := fmt.Sprint(bag["user_id"])
			bagMemberID := fmt.Sprint(bag["member"])
			bagMemberName, err := db.TableName(bagWorker.logTxn.DBconn, bagMemberID, bagAppID)

			if err != nil {
				bagWorker.Logger.Fatal("@MigrateBags > Table Name: ", err)
			}

			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			log.Println(fmt.Sprintf("Current Bag: { %s } | PK: %s, App: %s ", color.FgLightYellow.Render(bagMemberName), bagPK, bagAppID))

			if bagNode, err := bagWorker.ConstructBagNode(bagMemberName, bagMemberID, bagRowID); err == nil {

				if err := bagWorker.InitTransactions(); err != nil {
					return err
				}

				if migrated, err := bagWorker.HandleMigration(bagNode); err != nil {
					log.Println(fmt.Sprintf("%s { %s } | PK: %s | Owner: %s | %s ", color.FgGreen.Render("BAG NOT MIGRATED"), bagNode.Tag.Name, bagPK, bagUserID, err))
					if err := bagWorker.RollbackTransactions(); err != nil {
						bagWorker.Logger.Fatal(fmt.Sprintf("UNABLE TO ROLLBACK bag { %s } | Owner: %s | %s | PK: %s", bagNode.Tag.Name, bagUserID, err, bagPK))
					} else {
						log.Println(fmt.Sprintf("ROLLBACK bag { %s } | PK: %s ", bagNode.Tag.Name, bagPK))
					}
				} else {
					if migrated {
						log.Println(fmt.Sprintf("%s { %s } | PK: %s", color.FgLightGreen.Render("BAG MIGRATED"), bagNode.Tag.Name, bagPK))

						if !bagWorker.IsNodeDataEmpty(bagNode.Data) {
							bagWorker.Logger.Tracef("LeftOver Bag Data =>\n%v\n", bagNode.Data)
							log.Println(fmt.Sprintf("HANDLING LEFTOVER DATA { %s } | PK: %s | Owner: %s", bagNode.Tag.Name, bagPK, bagUserID))
							if err := bagWorker.SendNodeToBagWithOwnerID(bagNode, bagUserID); err != nil {
								bagWorker.Logger.Debug(bagNode)
								bagWorker.Logger.Fatal(fmt.Sprintf("@MigrateBags > SendNodeToBagWithOwnerID | { %s } | PK: %s", bagNode.Tag.Name, bagPK))
							} else {
								log.Println(fmt.Sprintf("%s { %s } | PK: %s | Owner: %s", color.FgYellow.Render("BAGGED LEFTOVER DATA"), bagNode.Tag.Name, bagPK, bagUserID))
							}
						}
					} else {
						log.Println(fmt.Sprintf("%s { %s } | PK: %s", color.FgGreen.Render("BAG NOT MIGRATED; NO ERR"), bagNode.Tag.Name, bagPK))
					}

					if err := bagWorker.CommitTransactions(); err != nil {
						bagWorker.Logger.Fatal(fmt.Sprintf("UNABLE TO COMMIT bag { %s } | %s ", bagNode.Tag.Name, err))
					} else {
						log.Println(fmt.Sprintf("COMMITTED bag { %s } | PK: %s", bagNode.Tag.Name, bagPK))
					}

					bagWorker.UpdateProcessedBagIDs(bagNode, processedBags)
				}

			} else {
				bagWorker.Logger.Fatal(fmt.Sprintf("UNABLE TO CONSTRUCT bag node | bagMemberName:%s  bagMemberID:%s  bagRowID:%s | %s", bagMemberName, bagMemberID, bagRowID, err))
			}
		}

		bagWorker.CloseDBConns()
	}

	return nil
}

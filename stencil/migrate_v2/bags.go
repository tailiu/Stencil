package migrate_v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	config "stencil/config/v2"
	"stencil/db"
	"stencil/helper"
	"stencil/reference_resolution"
	"strings"

	"github.com/gookit/color"
)

func (self *MigrationWorker) ConstructBagNode(bagMember, bagMemberID, bagRowID string) (*DependencyNode, error) {

	bagTag, err := self.SrcAppConfig.GetTagByMember(bagMember)
	if err != nil {
		self.Logger.Fatal(fmt.Sprintf("@ConstructBagNode: UNABLE TO GET BAG TAG BY MEMBER : %s | %s", bagMember, err))
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

func (self *MigrationWorker) GetRowsFromIDTable(app, member string, id interface{}, getFrom bool) ([]IDRow, error) {

	var idRows []IDRow
	var err error
	var idRowsDB []map[string]interface{}

	if !getFrom {
		idRowsDB, err = db.GetRowsFromIDTableByTo(self.logTxn.DBconn, app, member, helper.GetInt64(id))
	} else {
		idRowsDB, err = db.GetRowsFromIDTableByFrom(self.logTxn.DBconn, app, member, helper.GetInt64(id))
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
			FromID:       idRowDB["from_id"].(int64),
			ToAppName:    toAppName,
			ToAppID:      toAppID,
			ToMember:     toMember,
			ToMemberID:   toMemberID,
			ToID:         idRowDB["to_id"].(int64)})
	}
	return idRows, nil
}

func (self *MigrationWorker) MergeBagDataWithMappedData(mappedData *MappedData, node *DependencyNode, toTable config.ToTable) error {

	toTableData := make(map[string]ValueWithReference)

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
				self.Logger.Debug(node.Data)
				self.Logger.Fatal("@MigrateNode > FetchDataFromBags > id doesn't exist in table ", fromTable)
			}
		} else {
			self.Logger.Fatal("@MigrateNode > FetchDataFromBags > TableID, fromTable: error in getting table id for member! ", fromTable, err)
		}
	}

	if len(toTableData) > 0 {

		mappedCols := strings.Split(mappedData.cols, ",")
		for col, valWithRef := range toTableData {
			if !helper.Contains(mappedCols, col) {
				mappedData.cols += "," + col
				mappedData.ivals = append(mappedData.ivals, valWithRef.value)
				mappedData.vals += fmt.Sprintf(",$%d", len(mappedData.ivals))
				self.Logger.Tracef("@MigrateNode > FetchDataFromBags > Attr merged for: '%s.%s' = '%v'", toTable.Table, col, valWithRef.value)
				if valWithRef.ref != nil {
					mappedData.refs = append(mappedData.refs, *valWithRef.ref)
					self.Logger.Tracef("@MigrateNode > FetchDataFromBags > Ref merged for: %s.%s\n", toTable.Table, col)
				}

			}
		}
		mappedData.Trim(",")
		self.Logger.Tracef("@MigrateNode > FetchDataFromBags > Data merged for: %s\nData | %v", toTable.Table, toTableData)
	}

	return nil
}

func (self *MigrationWorker) FetchDataFromBags(visitedRows map[string]bool, toTableData map[string]ValueWithReference, prevUIDs map[string]string, app, member string, id interface{}, dstMemberID, dstMemberName, toTableName string) error {

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

			self.Logger.Tracef("@FetchDataFromBags > Processing Bag | ID: %v | PK: %v | App: %v | Member: %v", bagRow["id"], bagRow["pk"], bagRow["app"], bagRow["member"])

			bagData := make(map[string]interface{})
			if err := json.Unmarshal(bagRow["data"].([]byte), &bagData); err != nil {
				self.Logger.Debug(bagRow["data"])
				self.Logger.Debug(bagRow)
				self.Logger.Fatal("@FetchDataFromBags: UNABLE TO CONVERT BAG TO MAP ", bagRow, err)
				return err
			}
			fmt.Printf("bag data | %v\n", bagData)

			if bagMappedApp, mapping, found := self.FetchMappingsForBag(idRow.FromAppName, idRow.FromAppID, self.DstAppConfig.AppName, self.DstAppConfig.AppID, idRow.FromMember, dstMemberName); found {
				// self.Logger.Trace("@FetchDataFromBags > FetchMappingsForBag, Mappings found for | ", idRow.FromAppName, idRow.FromAppID, self.DstAppConfig.AppName, self.DstAppConfig.AppID, idRow.FromMember, dstMemberName)
				bagAppConfig, err := config.CreateAppConfig(idRow.FromAppName, idRow.FromAppID)
				if err != nil {
					log.Fatal("@FetchDataFromBags: ", err)
				}
				for _, toTable := range mapping.ToTables {
					if !strings.EqualFold(toTable.Table, toTableName) {
						continue
					}
					for toAttr, fromAttr := range toTable.Mapping {

						if strings.EqualFold("id", toAttr) {
							continue
						}

						var valueForNode interface{}
						var refForNode *MappingRef
						cleanedFromAttr := fromAttr

						if strings.Contains(fromAttr, "#FETCH") {
							continue
						} else if fromAttr[0:1] == "$" {
							if inputVal, err := bagMappedApp.GetInput(fromAttr); err == nil {
								valueForNode = inputVal
							} else {
								self.Logger.Debugf("@FetchDataFromBags | fromAttr [%s]", fromAttr)
								self.Logger.Fatal(err)
							}
						} else if bagVal, _, decodedFromAttr, ref, found, err := self.DecodeMappingValue(fromAttr, bagData, true); err == nil {
							cleanedFromAttr = decodedFromAttr
							if found && bagVal != nil {
								valueForNode = bagVal
								if bagAppConfig.AppID == self.DstAppConfig.AppID {
									fromAttrTokens := strings.Split(cleanedFromAttr, ".")
									if bagTag, err := bagAppConfig.GetTagByMember(fromAttrTokens[0]); err == nil {

										if depRefs, err := CreateReferencesViaDependencies(bagAppConfig, *bagTag, bagData, cleanedFromAttr); err != nil {
											log.Println(bagData)
											self.Logger.Fatal("@FetchDataFromBags > CreateReferencesViaDependencies: ", err)
										} else if len(depRefs) > 0 && ref == nil {
											ref = &depRefs[0]
										} else {
											log.Println("@FetchDataFromBags > CreateReferencesViaDependencies > No reference created: ", cleanedFromAttr)
										}

										if bagTag.Name != "root" {
											if ownRefs, err := CreateReferencesViaOwnerships(bagAppConfig, *bagTag, bagData, cleanedFromAttr); err != nil {
												log.Println(bagData)
												self.Logger.Fatal("@FetchDataFromBags > CreateReferencesViaOwnerships: ", err)
											} else if len(ownRefs) > 0 && ref == nil {
												ref = &ownRefs[0]
											} else {
												log.Println("@FetchDataFromBags > CreateReferencesViaOwnerships > No reference created: ", cleanedFromAttr)
											}
										}
									}
								}
								if ref != nil {
									ref.appID = fmt.Sprint(bagRow["app"])
									ref.mergedFromBag = true
									refForNode = ref
								}
							}
						} else {
							self.Logger.Debug(bagData)
							self.Logger.Fatalf("Unable to decode mapped val | fromAttr: [%s], toAttr: [%s]", fromAttr, toAttr)
						}

						if _, ok := toTableData[toAttr]; !ok {
							if valueForNode != nil {
								toTableData[toAttr] = ValueWithReference{value: valueForNode, ref: refForNode}
								// self.Logger.Tracef("@FetchDataFromBags > DecodeMappingValue | Added New | toTable: [%s], fromAttr: [%s], cleanedFromAttr: [%s], toAttr: [%s], valueForNode: [%v], Found: [%v]", toTable.Table, fromAttr, cleanedFromAttr, toAttr, valueForNode, found)
							} else {
								// self.Logger.Tracef("@FetchDataFromBags > DecodeMappingValue | Decoded but not Added | toTable: [%s], fromAttr: [%s], cleanedFromAttr: [%s], toAttr: [%s], valueForNode: [%v], Found: [%v]", toTable.Table, fromAttr, cleanedFromAttr, toAttr, valueForNode, found)
							}
						} else {
							// self.Logger.Tracef("@FetchDataFromBags > DecodeMappingValue | Exists | toTable: [%s], cleanedFromAttr: [%s], fromAttr: [%s], toAttr: [%s], BagVal: [%v]", toTable.Table, cleanedFromAttr, fromAttr, toAttr, toTableData[toAttr])
						}

						if strings.Contains(cleanedFromAttr, ".id") {
							// self.Logger.Tracef("@FetchDataFromBags > %s | [cleanedFromAttr:%s] | BagData:[%v]", color.FgLightCyan.Render("Not Deleting Attr From Bag"), color.FgBlue.Render(cleanedFromAttr), color.FgBlue.Render(bagData))
						} else {
							// self.Logger.Tracef("@FetchDataFromBags > %s | [%s] | BagData:[%v]", color.FgCyan.Render("Deleting Attr From Bag"), color.FgBlue.Render(cleanedFromAttr), color.FgBlue.Render(bagData))
							delete(bagData, cleanedFromAttr)
						}
					}
				}

				bagAppConfig.CloseDBConns()

				if self.IsNodeDataEmpty(bagData) {
					if err := db.DeleteBagV2(self.tx.StencilTx, fmt.Sprint(bagRow["pk"])); err != nil {
						self.Logger.Fatal("@FetchDataFromBags > DeleteBagV2, Unable to delete bag | ", bagRow["pk"])
						return err
					} else {
						log.Println(fmt.Sprintf("%s | PK: %v", color.FgLightRed.Render("Deleted BAG"), bagRow["pk"]))
					}
				} else {
					log.Println(fmt.Sprintf("%s | %v", color.FgYellow.Render("BAG NOT EMPTY"), bagData))
					if jsonData, err := json.Marshal(bagData); err == nil {
						if err := db.UpdateBag(self.tx.StencilTx, fmt.Sprint(bagRow["pk"]), self.logTxn.Txn_id, jsonData); err != nil {
							self.Logger.Fatal("@FetchDataFromBags: UNABLE TO UPDATE BAG ", bagRow, err)
							return err
						} else {
							log.Println(fmt.Sprintf("%s | PK: %v", color.FgLightYellow.Render("Updated BAG"), bagRow["pk"]))
						}
					} else {
						self.Logger.Fatal("@FetchDataFromBags > len(bagData) != 0, Unable to marshall bag | ", bagData)
						return err
					}
				}
			} else {
				// self.Logger.Warnf("No mappings found from [%s:%s] to [%s:%s]", idRow.FromAppName, idRow.FromMember, self.DstAppConfig.AppName, dstMemberName)
			}
		} else {
			// log.Println("@FetchDataFromBags > GetBagByAppMemberIDV2, No bags found for | ", prevUIDs[idRow.FromAppID], idRow.FromAppID, idRow.FromMember, idRow.FromID, self.logTxn.Txn_id)
		}

		// self.Logger.Tracef("@FetchDataFromBags > FetchDataFromBags: %s |\nFromAppID: %s FromMemberID: %s FromID: %v dstMemberID: %s dstMemberName: %s \ntoTableData: %v\nprevUIDs: %v ",
		// color.FgLightMagenta.Render("Recursive Traversal"), idRow.FromAppID, idRow.FromMemberID, idRow.FromID, dstMemberID, dstMemberName,
		// toTableData, prevUIDs)
		if err := self.FetchDataFromBags(visitedRows, toTableData, prevUIDs, idRow.FromAppID, idRow.FromMemberID, idRow.FromID, dstMemberID, dstMemberName, toTableName); err != nil {
			self.Logger.Fatal("@FetchDataFromBags > FetchDataFromBags: Error while recursing | ", toTableData, idRow.FromAppID, idRow.FromMember, idRow.FromID)
			return err
		}
	}
	return nil
}

func (self *MigrationWorker) SendMemberToBag(node *DependencyNode, member, ownerID string) error {
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
					} else {
						self.Logger.Infof("Bag Created for Member '%s' with ID '%v' \nData | %v", member, srcID, bagData)
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
				self.Logger.Warn("@SendMemberToBag: '", member, "' doesn't contain id! ", srcID)
				// return err
			}
		}
	}

	return nil
}

func (self *MigrationWorker) SendNodeToBagWithOwnerID(node *DependencyNode, ownerID string) error {
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

func (self *MigrationWorker) SendNodeToBag(node *DependencyNode) error {
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

func (bagWorker *MigrationWorker) DeleteBag(bagNode *DependencyNode) error {
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

func (self *MigrationWorker) UpdateProcessedBagIDs(bagNode *DependencyNode, processedBags map[string]bool) {
	for _, pk := range bagNode.PKs {
		pkStr := fmt.Sprint(pk)
		processedBags[pkStr] = true
	}
}

func (self *MigrationWorker) StartBagsMigration(userID, bagAppID string, threadID int, isBlade ...bool) error {
	color.Magenta.Println("########################################################################")
	color.LightMagenta.Printf("Starting Bags for User: '%s' | App: '%s' \n", userID, bagAppID)
	color.Magenta.Println("########################################################################")

	if bags, err := db.GetBagsV2(self.logTxn.DBconn, bagAppID, userID, self.logTxn.Txn_id); err != nil {
		self.Logger.Fatal(fmt.Sprintf("UNABLE TO FETCH BAGS FOR USER: %s | %s", userID, err))
		return err
	} else if len(bags) > 0 {
		log.Println("Bags fetched:  ", len(bags))

		bagWorker := self.mThread.CreateBagWorker(userID, bagAppID, self.DstAppConfig.AppID, threadID)

		log.Println(fmt.Sprintf("Bag Worker Created | %s -> %s ", bagWorker.SrcAppConfig.AppName, bagWorker.DstAppConfig.AppName))

		processedBags := make(map[string]bool)

		for _, bag := range bags {

			if strings.EqualFold(bagAppID, "4") {
				// break
			}

			bagPK := fmt.Sprint(bag["pk"])

			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			self.Logger.Tracef("@StartBagsMigration > Processing Bag | ID: %v | PK: %v | %s", bag["id"], bag["pk"], bagPK)

			if _, ok := processedBags[bagPK]; ok {
				self.Logger.Tracef("@StartBagsMigration > Bag Already Processed | ID: %v | PK: %v | %s", bag["id"], bag["pk"], bagPK)
				continue
			}

			bagRowID := fmt.Sprint(bag["id"])
			bagUserID := fmt.Sprint(bag["user_id"])
			bagMemberID := fmt.Sprint(bag["member"])
			bagMemberName, err := db.TableName(bagWorker.logTxn.DBconn, bagMemberID, bagAppID)

			if err != nil {
				bagWorker.Logger.Fatal("@StartBagsMigration > Table Name: ", err)
			}

			log.Println(fmt.Sprintf("Current Bag: { %s } | PK: %s, App: %s ", color.FgLightCyan.Render(bagMemberName), bagPK, bagAppID))

			if bagNode, err := bagWorker.ConstructBagNode(bagMemberName, bagMemberID, bagRowID); err == nil {

				fmt.Printf("bag data | %v\n", bagNode.Data)

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
								bagWorker.Logger.Fatal(fmt.Sprintf("@StartBagsMigration > SendNodeToBagWithOwnerID | { %s } | PK: %s", bagNode.Tag.Name, bagPK))
							} else {
								log.Println(fmt.Sprintf("%s { %s } | PK: %s | Owner: %s", color.FgLightYellow.Render("BAGGED LEFTOVER DATA"), bagNode.Tag.Name, bagPK, bagUserID))
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
	} else {
		self.Logger.Info("No Bags found!")
	}
	return nil
}

func (self *MigrationWorker) MigrateBags(threadID int, isBlade ...bool) error {

	bagAppID, userID := self.SrcAppConfig.AppID, self.uid

	for {
		if err := self.StartBagsMigration(userID, bagAppID, threadID, isBlade...); err != nil {
			self.Logger.Fatal(err)
		}
		self.Logger.Infof("BagWorker finished | App: '%s', UID: '%s' \n", bagAppID, userID)
		var prevIDErr error
		if bagAppID, userID, prevIDErr = self.GetUserIDAppIDFromPreviousMigration(bagAppID, userID); prevIDErr != nil {
			self.Logger.Fatal(prevIDErr)
		} else if len(bagAppID) <= 0 && len(userID) <= 0 {
			self.Logger.Info("Bag migration finished!")
			break
		}
	}

	return nil
}

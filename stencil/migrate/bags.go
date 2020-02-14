package migrate

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/reference_resolution"
	"strings"
)

func (self *MigrationWorkerV2) ConstructBagNode(bagMember, bagMemberID, bagRowID string) (*DependencyNode, error) {

	bagTag, err := self.SrcAppConfig.GetTagByMember(bagMember)
	if err != nil {
		log.Fatal(fmt.Sprintf("@MigrateBags: UNABLE TO GET BAG TAG BY MEMBER : %s | %s", bagMember, err))
		return nil, err
	}

	if ql := self.GetTagQLForBag(*bagTag); len(ql) > 0 {
		where := fmt.Sprintf(" WHERE %s.member = %s AND %s.id = %s", bagMember, bagMemberID, bagMember, bagRowID)
		sql := ql + where
		if data, err := db.DataCall1(self.logTxn.DBconn, sql); err == nil && len(data) > 0 {
			bagData := make(map[string]interface{})
			if err := json.Unmarshal(data["json_data"].([]byte), &bagData); err != nil {
				fmt.Println("@ConstructBagNode >> ", data)
				log.Fatal(fmt.Sprintf("@ConstructBagNode: UNABLE TO CONVERT BAG TO MAP: %s", err))
				return nil, err
			}
			bagNode := &DependencyNode{Tag: *bagTag, SQL: sql, Data: bagData, IDs: data["ids"].([]uint8)}
			return bagNode, nil
		} else {
			if err == nil {
				err = errors.New("no data returned for root node, doesn't exist?")
			} else {
				fmt.Println("@ConstructBagNode > DataCall1 | ", err)
			}
			fmt.Println(sql)
			log.Fatal("@ConstructBagNode: ", err)
			return nil, err
		}
	} else {
		log.Fatal("@ConstructBagNode > GetTagQLForBag : Failed to get tag ql | ", ql)
	}

	log.Fatal("@ConstructBagNode: End")

	return nil, nil
}

func (self *MigrationWorkerV2) GetRowsFromIDTable(app, member, id string, getFrom bool) ([]IDRow, error) {
	var idRows []IDRow
	var err error
	var idRowsDB []map[string]interface{}
	if !getFrom {
		idRowsDB, err = db.GetRowsFromIDTableByTo(self.logTxn.DBconn, app, member, id)
	} else {
		idRowsDB, err = db.GetRowsFromIDTableByFrom(self.logTxn.DBconn, app, member, id)
	}

	if err != nil {
		log.Fatal("@GetRowsFromIDTable > db.GetRowsFromIDTable, Unable to get bags | ", getFrom, self.uid, app, member, id, self.logTxn.Txn_id, err)
		return nil, err
	}
	for _, idRowDB := range idRowsDB {
		fromAppID := fmt.Sprint(idRowDB["from_app"])
		fromAppName, err := db.GetAppNameByAppID(self.logTxn.DBconn, fromAppID)
		if err != nil {
			log.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get name | ", fromAppID, err)
			return nil, err
		}

		toAppID := fmt.Sprint(idRowDB["to_app"])
		toAppName, err := db.GetAppNameByAppID(self.logTxn.DBconn, toAppID)
		if err != nil {
			log.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get name | ", toAppID, err)
			return nil, err
		}

		fromMemberID := fmt.Sprint(idRowDB["from_member"])
		fromMember, err := db.TableName(self.logTxn.DBconn, fromMemberID, fromAppID)
		if err != nil {
			log.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get fromMember | ", fromMemberID, fromAppID, err)
			return nil, err
		}

		toMemberID := fmt.Sprint(idRowDB["to_member"])
		toMember, err := db.TableName(self.logTxn.DBconn, toMemberID, toAppID)
		if err != nil {
			log.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get toMember | ", toMemberID, toAppID, err)
			return nil, err
		}

		idRows = append(idRows, IDRow{
			FromAppName:  fromAppName,
			FromAppID:    fromAppID,
			FromMemberID: fromMemberID,
			FromMember:   fromMember,
			FromID:       fmt.Sprint(idRowDB["from_id"]),
			ToAppName:    toAppName,
			ToAppID:      toAppID,
			ToMember:     toMember,
			ToMemberID:   toMemberID,
			ToID:         fmt.Sprint(idRowDB["to_id"])})
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
				if err := self.FetchDataFromBags(visitedRows, toTableData, prevUIDs, self.SrcAppConfig.AppID, fromTableID, fmt.Sprint(fromID), toTable.TableID, toTable.Table); err != nil {
					log.Fatal("@MigrateNode > FetchDataFromBags | ", err)
				}
			} else {
				log.Fatal("@MigrateNode > FetchDataFromBags > id doesn't exist in table ", fromTable, err)
			}
		} else {
			log.Fatal("@MigrateNode > FetchDataFromBags > TableID, fromTable: error in getting table id for member! ", fromTable, err)
		}
	}

	if len(toTableData) > 0 {

		for col, val := range toTableData {
			if !strings.Contains(mappedData.cols, col) {
				mappedData.cols += "," + col
				mappedData.ivals = append(mappedData.ivals, val)
				mappedData.vals += fmt.Sprintf(",$%d", len(mappedData.ivals))
			}
		}
		mappedData.Trim(",")
		log.Println("@MigrateNode > FetchDataFromBags > Data merged for: ", toTable.Table)
	}

	return nil
}

func (self *MigrationWorkerV2) FetchDataFromBags(visitedRows map[string]bool, toTableData map[string]interface{}, prevUIDs map[string]string, app, member, id, dstMemberID, dstMemberName string) error {

	currentRow := fmt.Sprintf("%s:%s:%s", app, member, id)

	if _, ok := visitedRows[currentRow]; ok {
		return nil
	} else {
		visitedRows[currentRow] = true
	}

	idRows, err := self.GetRowsFromIDTable(app, member, id, false)

	if err != nil {
		log.Fatal("@FetchDataFromBags > GetRowsFromIDTable, Unable to get IDRows | ", app, member, id, false, err)
		return err
	} else {
		log.Println("@FetchDataFromBags > GetRowsFromIDTable | ", idRows)
	}

	for _, idRow := range idRows {

		bagRow, err := db.GetBagByAppMemberIDV2(self.logTxn.DBconn, prevUIDs[idRow.FromAppID], idRow.FromAppID, idRow.FromMemberID, idRow.FromID, self.logTxn.Txn_id)
		if err != nil {
			log.Fatal("@FetchDataFromBags > GetBagByAppMemberIDV2, Unable to get bags | ", prevUIDs[idRow.FromAppID], idRow.FromAppID, idRow.FromMemberID, idRow.FromID, self.logTxn.Txn_id, err)
			return err
		}
		if bagRow != nil {

			bagData := make(map[string]interface{})
			if err := json.Unmarshal(bagRow["data"].([]byte), &bagData); err != nil {
				fmt.Println(bagRow["data"])
				fmt.Println(bagRow)
				log.Fatal("@FetchDataFromBags: UNABLE TO CONVERT BAG TO MAP ", bagRow, err)
				return err
			}

			if mapping, found := self.FetchMappingsForBag(idRow.FromAppName, idRow.FromAppID, self.DstAppConfig.AppName, self.DstAppConfig.AppID, idRow.FromMember, dstMemberName); found {
				log.Println("@FetchDataFromBags > FetchMappingsForBag, Mappings found for | ", idRow.FromAppName, idRow.FromAppID, self.DstAppConfig.AppName, self.DstAppConfig.AppID, idRow.FromMember, dstMemberName)
				for _, toTable := range mapping.ToTables {
					for fromAttr, toAttr := range toTable.Mapping {
						if _, ok := toTableData[fromAttr]; !ok {
							if bagVal, exists := bagData[toAttr]; exists {
								toTableData[fromAttr] = bagVal
							}
						}
						delete(bagData, toAttr)
					}
				}

				if len(bagData) == 0 {
					if err := db.DeleteBagV2(self.tx.StencilTx, fmt.Sprint(bagRow["pk"])); err != nil {
						log.Fatal("@FetchDataFromBags > DeleteBagV2, Unable to delete bag | ", bagRow["pk"])
						return err
					}
				} else {
					if jsonData, err := json.Marshal(bagData); err == nil {
						if err := db.UpdateBag(self.tx.StencilTx, fmt.Sprint(bagRow["pk"]), self.logTxn.Txn_id, jsonData); err != nil {
							log.Fatal("@FetchDataFromBags: UNABLE TO UPDATE BAG ", bagRow, err)
							return err
						}
					} else {
						log.Fatal("@FetchDataFromBags > len(bagData) != 0, Unable to marshall bag | ", bagData)
						return err
					}
				}
			} else {
				log.Println("@FetchDataFromBags > FetchMappingsForBag, No mappings found for | ", idRow.FromAppName, idRow.FromAppID, self.DstAppConfig.AppName, self.DstAppConfig.AppID, idRow.FromMember, dstMemberName)
			}
		} else {
			log.Println("@FetchDataFromBags > GetBagByAppMemberIDV2, No bags found for | ", prevUIDs[idRow.FromAppID], idRow.FromAppID, idRow.FromMember, idRow.FromID, self.logTxn.Txn_id)
		}

		log.Println("@FetchDataFromBags > FetchDataFromBags: Recursive Traversal | ", toTableData, prevUIDs, idRow.FromAppID, idRow.FromMemberID, idRow.FromID, dstMemberID, dstMemberName)
		if err := self.FetchDataFromBags(visitedRows, toTableData, prevUIDs, idRow.FromAppID, idRow.FromMemberID, idRow.FromID, dstMemberID, dstMemberName); err != nil {
			log.Fatal("@FetchDataFromBags > FetchDataFromBags: Error while recursing | ", toTableData, idRow.FromAppID, idRow.FromMember, idRow.FromID)
			return err
		}
	}
	return nil
}

func (self *MigrationWorkerV2) SendMemberToBag(node *DependencyNode, member, ownerID string, fromNode bool) error {
	if !strings.EqualFold(self.mtype, DELETION) {
		return nil
	}
	if memberID, err := db.TableID(self.logTxn.DBconn, member, self.SrcAppConfig.AppID); err != nil {
		log.Fatal("@SendMemberToBag > TableID: error in getting table id for member! ", member, err)
		return err
	} else {
		bagData := make(map[string]interface{})
		for col, val := range node.Data {
			colTokens := strings.Split(col, ".")
			colMember := colTokens[0]
			// colAttr := colTokens[1]
			if strings.Contains(colMember, member) {
				bagData[col] = val
			}
		}
		if len(bagData) > 0 {
			if id, ok := node.Data[member+".id"]; ok && id != nil {
				srcID := fmt.Sprint(id)
				if jsonData, err := json.Marshal(bagData); err == nil {
					if err := db.CreateNewBag(self.tx.StencilTx, self.SrcAppConfig.AppID, memberID, srcID, ownerID, fmt.Sprint(self.logTxn.Txn_id), jsonData); err != nil {
						fmt.Println(self.SrcAppConfig.AppName, member, srcID, ownerID)
						fmt.Println(bagData)
						log.Fatal("@SendMemberToBag: error in creating bag! ", err)
						return err
					}
					// if serr := db.SaveForEvaluation(self.logTxn.DBconn, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, memberID, "n/a", srcID, "n/a", "*", "n/a", fmt.Sprint(self.logTxn.Txn_id)); serr != nil {
					// 	log.Fatal("@SendMemberToBag > SaveForEvaluation =>", serr)
					// }
					if derr := db.ReallyDeleteRowFromAppDB(self.tx.SrcTx, member, fmt.Sprint(id)); derr != nil {
						fmt.Println("@SendMemberToBag > DeleteRowFromAppDB")
						log.Fatal(derr)
						return derr
					}
				} else {
					fmt.Println(bagData)
					log.Fatal("@SendMemberToBag: unable to convert bag data to JSON ", err)
					return err
				}
				if !fromNode {
					if err := self.AddInnerReferences(node, member); err != nil {
						fmt.Println(node.Tag.Members)
						fmt.Println(node.Tag.InnerDependencies)
						fmt.Println(node.Data)
						log.Fatal("@SendMemberToBag > AddInnerReferences: Adding Inner References failed ", err)
						return err
					}
				}
			} else {
				// fmt.Println(node.Data)
				// fmt.Println(node.SQL)
				log.Println("@SendMemberToBag: '", member, "' doesn't contain id! ", id)
				// return err
			}
		}
	}

	return nil
}

func (self *MigrationWorkerV2) SendNodeToBagWithOwnerID(node *DependencyNode, ownerID string) error {
	if !strings.EqualFold(self.mtype, DELETION) {
		return nil
	}
	for _, member := range node.Tag.Members {
		if err := self.SendMemberToBag(node, member, ownerID, true); err != nil {
			fmt.Println(node)
			log.Fatal("@SendNodeToBagWithOwnerID > SendMemberToBag: ownerID error! ")
			return err
		}
	}
	if err := self.AddInnerReferences(node, ""); err != nil {
		fmt.Println(node.Tag.Members)
		fmt.Println(node.Tag.InnerDependencies)
		fmt.Println(node.Data)
		log.Fatal("@SendNodeToBagWithOwnerID > AddInnerReferences: Adding Inner References failed ", err)
		return err
	}
	return nil
}

func (self *MigrationWorkerV2) SendNodeToBag(node *DependencyNode) error {
	if !strings.EqualFold(self.mtype, DELETION) {
		return nil
	}
	if ownerID, _ := self.GetNodeOwner(node); len(ownerID) > 0 {
		if err := self.SendNodeToBagWithOwnerID(node, ownerID); err != nil {
			return err
		}
	} else {
		fmt.Println(node)
		log.Fatal("@SendNodeToBag > GetNodeOwner: ownerID error! ")
	}

	return nil
}

func (self *MigrationWorkerV2) MigrateBags(threadID int, isBlade ...bool) error {

	prevIDs := reference_resolution.GetPrevUserIDs(self.logTxn.DBconn, self.SrcAppConfig.AppID, self.uid)
	prevIDs[self.SrcAppConfig.AppID] = self.uid

	for bagAppID, userID := range prevIDs {

		log.Println(fmt.Sprintf("x%2dx Starting Bags for User: %s App: %s", threadID, userID, bagAppID))

		bags, err := db.GetBagsV2(self.logTxn.DBconn, bagAppID, userID, self.logTxn.Txn_id)

		if err != nil {
			log.Fatal(fmt.Sprintf("x%2dx UNABLE TO FETCH BAGS FOR USER: %s | %s", threadID, userID, err))
			return err
		}

		bagWorker := CreateBagWorkerV2(userID, bagAppID, self.DstAppConfig.AppID, self.logTxn, BAGS, threadID, isBlade...)

		log.Println(fmt.Sprintf("x%2dx Bag Worker Created | %s -> %s ", threadID, bagWorker.SrcAppConfig.AppName, bagWorker.DstAppConfig.AppName))

		for _, bag := range bags {

			bagRowID := fmt.Sprint(bag["id"])
			bagMemberID := fmt.Sprint(bag["member"])
			bagMemberName, err := db.TableName(bagWorker.logTxn.DBconn, bagMemberID, bagAppID)

			if err != nil {
				log.Fatal("@MigrateBags > Table Name: ", err)
			}

			bagID := fmt.Sprint(bag["pk"])

			log.Println(fmt.Sprintf("~%2d~ Current    Bag: { %s } | ID: %s, App: %s ", threadID, bagMemberName, bagID, bagAppID))

			if bagNode, err := bagWorker.ConstructBagNode(bagMemberName, bagMemberID, bagRowID); err == nil {

				if err := bagWorker.InitTransactions(); err != nil {
					return err
				}

				toCommit := true
				if err := bagWorker.HandleMigration(bagNode, true); err != nil {
					toCommit = false
					// fmt.Println(bag)
					log.Println(fmt.Sprintf("x%2dx UNABLE TO MIGRATE BAG { %s } | ID: %s | %s ", threadID, bagNode.Tag.Name, bagID, err))
				} else {
					log.Println(fmt.Sprintf("x%2dx MIGRATED bag { %s } | ID: %s", threadID, bagNode.Tag.Name, bagID))

					if bagWorker.IsNodeDataEmpty(bagNode) {
						if err := db.DeleteBagV2(bagWorker.tx.StencilTx, bagID); err != nil {
							toCommit = false
							fmt.Println(bag)
							log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO DELETE BAG : %s | %s ", threadID, bagID, err))
							// return err
						} else {
							log.Println(fmt.Sprintf("x%2dx DELETED bag { %s } ", threadID, bagNode.Tag.Name))
						}
					} else {
						if jsonData, err := json.Marshal(bagNode.Data); err == nil {
							if err := db.UpdateBag(bagWorker.tx.StencilTx, bagID, bagWorker.logTxn.Txn_id, jsonData); err != nil {
								toCommit = false
								fmt.Println(bag)
								log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO UPDATE BAG : %s | %s ", threadID, bagID, err))
								// return err
							} else {
								log.Println(fmt.Sprintf("x%2dx UPDATED bag { %s } ", threadID, bagNode.Tag.Name))
							}
						} else {
							toCommit = false
							fmt.Println(bagNode.Data)
							log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO MARSHALL BAG DATA : %s | %s ", threadID, bagID, err))
							// return err
						}
					}
				}
				if toCommit {
					if err := bagWorker.CommitTransactions(); err != nil {
						log.Fatal(fmt.Sprintf("x%2dx UNABLE TO COMMIT bag { %s } | %s ", threadID, bagNode.Tag.Name, err))
						// return err
					} else {
						log.Println(fmt.Sprintf("x%2dx COMMITTED bag { %s } ", threadID, bagNode.Tag.Name))
					}
				} else {
					log.Println(fmt.Sprintf("x%2dx ROLLBACK bag { %s } | ID: %s ", threadID, bagNode.Tag.Name, bagID))
					bagWorker.tx.DstTx.Rollback()
					bagWorker.tx.SrcTx.Rollback()
					bagWorker.tx.StencilTx.Rollback()
				}

			} else {
				log.Fatal(fmt.Sprintf("x%2dx UNABLE TO CONSTRUCT bag node | bagMemberName:%s  bagMemberID:%s  bagRowID:%s | %s", threadID, bagMemberName, bagMemberID, bagRowID, err))
			}
		}

		bagWorker.CloseDBConns()
	}

	return nil
}
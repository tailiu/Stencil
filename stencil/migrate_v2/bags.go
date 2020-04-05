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

func (mWorker *MigrationWorker) ConstructBagNode(bagMember, bagMemberID, bagRowID string) (*DependencyNode, error) {

	bagTag, err := mWorker.SrcAppConfig.GetTagByMember(bagMember)
	if err != nil {
		mWorker.Logger.Fatal(fmt.Sprintf("@ConstructBagNode: UNABLE TO GET BAG TAG BY MEMBER : %s | %s", bagMember, err))
		return nil, err
	}

	if ql := mWorker.GetTagQLForBag(*bagTag); len(ql) > 0 {
		where := fmt.Sprintf(" WHERE %s.member = %s AND %s.id = %s", bagMember, bagMemberID, bagMember, bagRowID)
		sql := ql + where
		if data, err := db.DataCall1(mWorker.logTxn.DBconn, sql); err == nil && len(data) > 0 {
			bagData := make(map[string]interface{})
			if err := json.Unmarshal(data["json_data"].([]byte), &bagData); err != nil {
				fmt.Println("@ConstructBagNode >> ", data)
				mWorker.Logger.Fatal(fmt.Sprintf("@ConstructBagNode: UNABLE TO CONVERT BAG TO MAP: %s", err))
				return nil, err
			}
			var bagPKs []int64
			if err := json.Unmarshal(data["pks_json"].([]byte), &bagPKs); err != nil {
				fmt.Println("@ConstructBagNode >> ", data)
				mWorker.Logger.Fatal(fmt.Sprintf("@ConstructBagNode: UNABLE TO CONVERT pks_json TO []int64: %s", err))
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
			mWorker.Logger.Fatal("@ConstructBagNode: ", err)
			return nil, err
		}
	} else {
		mWorker.Logger.Fatal("@ConstructBagNode > GetTagQLForBag : Failed to get tag ql | ", ql)
	}

	mWorker.Logger.Fatal("@ConstructBagNode: End")

	return nil, nil
}

func (mWorker *MigrationWorker) GetRowsFromIDTable(app, member string, id interface{}, getFrom bool) ([]AttrRow, error) {

	var AttrRows []AttrRow
	var err error
	var AttrRowsDB []map[string]interface{}

	if !getFrom {
		AttrRowsDB, err = db.GetRowsFromIDTableByTo(mWorker.logTxn.DBconn, app, member, helper.GetInt64(id))
	} else {
		AttrRowsDB, err = db.GetRowsFromIDTableByFrom(mWorker.logTxn.DBconn, app, member, helper.GetInt64(id))
	}

	if err != nil {
		mWorker.Logger.Fatalf("Unable to get bags | %s | getFrom:%v uid:%v app:%v member:%v id:%v txnID:%v", err, getFrom, mWorker.uid, app, member, id, mWorker.logTxn.Txn_id)
		return nil, err
	}

	for _, AttrRowDB := range AttrRowsDB {
		fromAppID := fmt.Sprint(AttrRowDB["from_app"])
		fromAppName, err := db.GetAppNameByAppID(mWorker.logTxn.DBconn, fromAppID)
		if err != nil {
			mWorker.Logger.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get name | ", fromAppID, err)
			return nil, err
		}

		toAppID := fmt.Sprint(AttrRowDB["to_app"])
		toAppName, err := db.GetAppNameByAppID(mWorker.logTxn.DBconn, toAppID)
		if err != nil {
			mWorker.Logger.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get name | ", toAppID, err)
			return nil, err
		}

		fromMemberID := fmt.Sprint(AttrRowDB["from_member"])
		fromMember, err := db.TableName(mWorker.logTxn.DBconn, fromMemberID, fromAppID)
		if err != nil {
			mWorker.Logger.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get fromMember | ", fromMemberID, fromAppID, err)
			return nil, err
		}

		toMemberID := fmt.Sprint(AttrRowDB["to_member"])
		toMember, err := db.TableName(mWorker.logTxn.DBconn, toMemberID, toAppID)
		if err != nil {
			mWorker.Logger.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get toMember | ", toMemberID, toAppID, err)
			return nil, err
		}

		AttrRows = append(AttrRows, AttrRow{
			FromAppName:  fromAppName,
			FromAppID:    fromAppID,
			FromMemberID: fromMemberID,
			FromMember:   fromMember,
			FromID:       AttrRowDB["from_id"].(int64),
			ToAppName:    toAppName,
			ToAppID:      toAppID,
			ToMember:     toMember,
			ToMemberID:   toMemberID,
			ToID:         AttrRowDB["to_id"].(int64)})
	}
	return AttrRows, nil
}

func (mWorker *MigrationWorker) MergeBagDataWithMappedData(mappedData *MappedMemberData, node *DependencyNode, toTable config.ToTable) error {

	toTableData := make(map[string]ValueWithReference)

	prevUIDs := reference_resolution.GetPrevUserIDs(mWorker.logTxn.DBconn, mWorker.SrcAppConfig.AppID, mWorker.uid)
	if prevUIDs == nil {
		prevUIDs = make(map[string]string)
	}
	prevUIDs[mWorker.SrcAppConfig.AppID] = mWorker.uid

	for _, fromTable := range mappedData.SrcTables() {
		if fromID, ok := node.Data[fromTable.Name+".id"]; ok {
			visitedRows := make(map[string]bool)
			if err := mWorker.FetchDataFromBags(visitedRows, toTableData, prevUIDs, mWorker.SrcAppConfig.AppID, fromTable.ID, fromID, toTable.TableID, toTable.Table, toTable.Table); err != nil {
				mWorker.Logger.Fatal("@MigrateNode > FetchDataFromBags | ", err)
			}
		} else {
			mWorker.Logger.Debug(node.Data)
			mWorker.Logger.Fatal("@MigrateNode > FetchDataFromBags > id doesn't exist in table ", fromTable)
		}
	}

	if len(toTableData) > 0 {

		mappedCols := mappedData.ToCols()
		for col, valWithRef := range toTableData {
			if !helper.Contains(mappedCols, col) {
				mmv := MappedMemberValue{
					Value:  valWithRef.value,
					DBConn: mWorker.logTxn.DBconn,
				}
				if valWithRef.ref != nil {
					mmv.Ref = valWithRef.ref
					mWorker.Logger.Tracef("@MigrateNode > FetchDataFromBags > Ref merged for: %s.%s\n", toTable.Table, col)
				}
				mappedData.Data[col] = mmv
				mWorker.Logger.Tracef("@MigrateNode > FetchDataFromBags > Attr merged for: '%s.%s' = '%v'", toTable.Table, col, valWithRef.value)
			}
		}
		mWorker.Logger.Tracef("@MigrateNode > FetchDataFromBags > Data merged for: %s\nData | %v", toTable.Table, toTableData)
	}

	return nil
}

func (mWorker *MigrationWorker) FetchDataFromBags(visitedRows map[string]bool, toTableData map[string]ValueWithReference, prevUIDs map[string]string, app, member string, id interface{}, dstMemberID, dstMemberName, toTableName string) error {

	currentRow := fmt.Sprintf("%s:%s:%s", app, member, id)

	if _, ok := visitedRows[currentRow]; ok {
		return nil
	} else {
		visitedRows[currentRow] = true
	}

	AttrRows, err := mWorker.GetRowsFromIDTable(app, member, id, false)

	if err != nil {
		mWorker.Logger.Fatal("@FetchDataFromBags > GetRowsFromIDTable, Unable to get AttrRows | ", app, member, id, false, err)
		return err
	} else if len(AttrRows) < 1 {
		// mWorker.Logger.Trace("@FetchDataFromBags > GetRowsFromIDTable | No AttrRows found | ", app, member, id, false)
		return nil
	} else {
		mWorker.Logger.Trace("@FetchDataFromBags > GetRowsFromIDTable | ", AttrRows)
	}

	for _, AttrRow := range AttrRows {

		bagRow, err := db.GetBagByAppMemberIDV2(mWorker.logTxn.DBconn, prevUIDs[AttrRow.FromAppID], AttrRow.FromAppID, AttrRow.FromMemberID, AttrRow.FromID, mWorker.logTxn.Txn_id)
		if err != nil {
			mWorker.Logger.Fatal("@FetchDataFromBags > GetBagByAppMemberIDV2, Unable to get bags | ", prevUIDs[AttrRow.FromAppID], AttrRow.FromAppID, AttrRow.FromMemberID, AttrRow.FromID, mWorker.logTxn.Txn_id, err)
			return err
		}
		if bagRow != nil {

			mWorker.Logger.Tracef("@FetchDataFromBags > Processing Bag | ID: %v | PK: %v | App: %v | Member: %v", bagRow["id"], bagRow["pk"], bagRow["app"], bagRow["member"])

			bagData := make(DataMap)
			if err := json.Unmarshal(bagRow["data"].([]byte), &bagData); err != nil {
				mWorker.Logger.Debug(bagRow["data"])
				mWorker.Logger.Debug(bagRow)
				mWorker.Logger.Fatal("@FetchDataFromBags: UNABLE TO CONVERT BAG TO MAP ", bagRow, err)
				return err
			}
			fmt.Printf("bag data | %v\n", bagData)

			if bagMappedApp, mapping, found := mWorker.FetchMappingsForBag(AttrRow.FromAppName, AttrRow.FromAppID, mWorker.DstAppConfig.AppName, mWorker.DstAppConfig.AppID, AttrRow.FromMember, dstMemberName); found {
				// mWorker.Logger.Trace("@FetchDataFromBags > FetchMappingsForBag, Mappings found for | ", AttrRow.FromAppName, AttrRow.FromAppID, mWorker.DstAppConfig.AppName, mWorker.DstAppConfig.AppID, AttrRow.FromMember, dstMemberName)
				bagAppConfig, err := config.CreateAppConfig(AttrRow.FromAppName, AttrRow.FromAppID)
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
								mWorker.Logger.Debugf("@FetchDataFromBags | fromAttr [%s]", fromAttr)
								mWorker.Logger.Fatal(err)
							}
						} else if mmv, err := mWorker.ResolveMappedStatement(fromAttr, bagData); err == nil {
							cleanedFromAttr = mmv.GetMemberAttr()
							if mmv.Value != nil {
								valueForNode = mmv.Value
								if bagAppConfig.AppID == mWorker.DstAppConfig.AppID {
									if bagTag, err := bagAppConfig.GetTagByMember(mmv.FromMember); err == nil {
										if depRefs, err := CreateReferencesViaDependencies(bagAppConfig, *bagTag, bagData, cleanedFromAttr); err != nil {
											log.Println(bagData)
											mWorker.Logger.Fatal("@FetchDataFromBags > CreateReferencesViaDependencies: ", err)
										} else if len(depRefs) > 0 && mmv.Ref == nil {
											mmv.Ref = &depRefs[0]
										} else {
											log.Println("@FetchDataFromBags > CreateReferencesViaDependencies > No reference created: ", cleanedFromAttr)
										}

										if bagTag.Name != "root" {
											if ownRefs, err := CreateReferencesViaOwnerships(bagAppConfig, *bagTag, bagData, cleanedFromAttr); err != nil {
												log.Println(bagData)
												mWorker.Logger.Fatal("@FetchDataFromBags > CreateReferencesViaOwnerships: ", err)
											} else if len(ownRefs) > 0 && mmv.Ref == nil {
												mmv.Ref = &ownRefs[0]
											} else {
												log.Println("@FetchDataFromBags > CreateReferencesViaOwnerships > No reference created: ", cleanedFromAttr)
											}
										}
									}
								}
								if mmv.Ref != nil {
									mmv.Ref.appID = fmt.Sprint(bagRow["app"])
									mmv.Ref.mergedFromBag = true
									refForNode = mmv.Ref
								}
							}
						} else {
							mWorker.Logger.Debug(bagData)
							mWorker.Logger.Fatalf("Unable to decode mapped val | fromAttr: [%s], toAttr: [%s]", fromAttr, toAttr)
						}

						if _, ok := toTableData[toAttr]; !ok {
							if valueForNode != nil {
								toTableData[toAttr] = ValueWithReference{value: valueForNode, ref: refForNode}
								// mWorker.Logger.Tracef("@FetchDataFromBags > DecodeMappingValue | Added New | toTable: [%s], fromAttr: [%s], cleanedFromAttr: [%s], toAttr: [%s], valueForNode: [%v], Found: [%v]", toTable.Table, fromAttr, cleanedFromAttr, toAttr, valueForNode, found)
							} else {
								// mWorker.Logger.Tracef("@FetchDataFromBags > DecodeMappingValue | Decoded but not Added | toTable: [%s], fromAttr: [%s], cleanedFromAttr: [%s], toAttr: [%s], valueForNode: [%v], Found: [%v]", toTable.Table, fromAttr, cleanedFromAttr, toAttr, valueForNode, found)
							}
						} else {
							// mWorker.Logger.Tracef("@FetchDataFromBags > DecodeMappingValue | Exists | toTable: [%s], cleanedFromAttr: [%s], fromAttr: [%s], toAttr: [%s], BagVal: [%v]", toTable.Table, cleanedFromAttr, fromAttr, toAttr, toTableData[toAttr])
						}

						if strings.Contains(cleanedFromAttr, ".id") {
							// mWorker.Logger.Tracef("@FetchDataFromBags > %s | [cleanedFromAttr:%s] | BagData:[%v]", color.FgLightCyan.Render("Not Deleting Attr From Bag"), color.FgBlue.Render(cleanedFromAttr), color.FgBlue.Render(bagData))
						} else {
							// mWorker.Logger.Tracef("@FetchDataFromBags > %s | [%s] | BagData:[%v]", color.FgCyan.Render("Deleting Attr From Bag"), color.FgBlue.Render(cleanedFromAttr), color.FgBlue.Render(bagData))
							delete(bagData, cleanedFromAttr)
						}
					}
				}

				bagAppConfig.CloseDBConns()

				if bagData.IsEmptyExcept() {
					if err := db.DeleteBagV2(mWorker.tx.StencilTx, fmt.Sprint(bagRow["pk"])); err != nil {
						mWorker.Logger.Fatal("@FetchDataFromBags > DeleteBagV2, Unable to delete bag | ", bagRow["pk"])
						return err
					} else {
						log.Println(fmt.Sprintf("%s | PK: %v", color.FgLightRed.Render("Deleted BAG"), bagRow["pk"]))
					}
				} else {
					log.Println(fmt.Sprintf("%s | %v", color.FgYellow.Render("BAG NOT EMPTY"), bagData))
					if jsonData, err := json.Marshal(bagData); err == nil {
						if err := db.UpdateBag(mWorker.tx.StencilTx, fmt.Sprint(bagRow["pk"]), mWorker.logTxn.Txn_id, jsonData); err != nil {
							mWorker.Logger.Fatal("@FetchDataFromBags: UNABLE TO UPDATE BAG ", bagRow, err)
							return err
						} else {
							log.Println(fmt.Sprintf("%s | PK: %v", color.FgLightYellow.Render("Updated BAG"), bagRow["pk"]))
						}
					} else {
						mWorker.Logger.Fatal("@FetchDataFromBags > len(bagData) != 0, Unable to marshall bag | ", bagData)
						return err
					}
				}
			} else {
				// mWorker.Logger.Warnf("No mappings found from [%s:%s] to [%s:%s]", AttrRow.FromAppName, AttrRow.FromMember, mWorker.DstAppConfig.AppName, dstMemberName)
			}
		} else {
			// log.Println("@FetchDataFromBags > GetBagByAppMemberIDV2, No bags found for | ", prevUIDs[AttrRow.FromAppID], AttrRow.FromAppID, AttrRow.FromMember, AttrRow.FromID, mWorker.logTxn.Txn_id)
		}

		// mWorker.Logger.Tracef("@FetchDataFromBags > FetchDataFromBags: %s |\nFromAppID: %s FromMemberID: %s FromID: %v dstMemberID: %s dstMemberName: %s \ntoTableData: %v\nprevUIDs: %v ",
		// color.FgLightMagenta.Render("Recursive Traversal"), AttrRow.FromAppID, AttrRow.FromMemberID, AttrRow.FromID, dstMemberID, dstMemberName,
		// toTableData, prevUIDs)
		if err := mWorker.FetchDataFromBags(visitedRows, toTableData, prevUIDs, AttrRow.FromAppID, AttrRow.FromMemberID, AttrRow.FromID, dstMemberID, dstMemberName, toTableName); err != nil {
			mWorker.Logger.Fatal("@FetchDataFromBags > FetchDataFromBags: Error while recursing | ", toTableData, AttrRow.FromAppID, AttrRow.FromMember, AttrRow.FromID)
			return err
		}
	}
	return nil
}

func (mWorker *MigrationWorker) SendMemberToBag(node *DependencyNode, member, ownerID string) error {
	if mWorker.mtype != DELETION && mWorker.mtype != BAGS {
		return nil
	}
	if memberID, err := db.TableID(mWorker.logTxn.DBconn, member, mWorker.SrcAppConfig.AppID); err != nil {
		mWorker.Logger.Fatal("@SendMemberToBag > TableID: error in getting table id for member! ", member, err)
		return err
	} else {
		bagData := mWorker.GetMemberDataFromNode(member, node.Data)
		if len(bagData) > 0 {
			if srcID, ok := node.Data[member+".id"]; ok && srcID != nil {
				if jsonData, err := json.Marshal(bagData); err == nil {
					if err := db.CreateNewBag(mWorker.tx.StencilTx, mWorker.SrcAppConfig.AppID, memberID, srcID, ownerID, fmt.Sprint(mWorker.logTxn.Txn_id), jsonData); err != nil {
						fmt.Println(mWorker.SrcAppConfig.AppName, member, srcID, ownerID)
						fmt.Println(bagData)
						mWorker.Logger.Fatal("@SendMemberToBag: error in creating bag! ", err)
						return err
					} else {
						mWorker.Logger.Infof("Bag Created for Member '%s' with ID '%v' \nData | %v", member, srcID, bagData)
					}
					if mWorker.mtype == BAGS {

					} else {
						if derr := db.ReallyDeleteRowFromAppDB(mWorker.tx.SrcTx, member, srcID); derr != nil {
							fmt.Println("@SendMemberToBag > DeleteRowFromAppDB")
							mWorker.Logger.Fatal(derr)
							return derr
						}
					}
				} else {
					fmt.Println(bagData)
					mWorker.Logger.Fatal("@SendMemberToBag: unable to convert bag data to JSON ", err)
					return err
				}
			} else {
				// fmt.Println(node.Data)
				// fmt.Println(node.SQL)
				mWorker.Logger.Warn("@SendMemberToBag: '", member, "' doesn't contain id! ", srcID)
				// return err
			}
		}
	}

	return nil
}

func (mWorker *MigrationWorker) SendNodeToBagWithOwnerID(node *DependencyNode, ownerID string) error {
	if mWorker.mtype != DELETION && mWorker.mtype != BAGS {
		return nil
	}
	for _, member := range node.Tag.Members {
		if err := mWorker.SendMemberToBag(node, member, ownerID); err != nil {
			fmt.Println(node)
			mWorker.Logger.Fatal("@SendNodeToBagWithOwnerID > SendMemberToBag: ownerID error! ")
			return err
		}
		log.Printf("%s { %s :  %s } | Owner ID: %v \n", color.FgYellow.Render("BAG"), node.Tag.Name, member, ownerID)
	}
	return nil
}

func (mWorker *MigrationWorker) SendNodeToBag(node *DependencyNode) error {
	if mWorker.mtype != DELETION {
		return nil
	}
	if ownerID, _ := mWorker.GetNodeOwner(node); len(ownerID) > 0 {
		if err := mWorker.SendNodeToBagWithOwnerID(node, ownerID); err != nil {
			return err
		}
	} else {
		fmt.Println(node)
		mWorker.Logger.Fatal("@SendNodeToBag > GetNodeOwner: ownerID error! ")
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

func (mWorker *MigrationWorker) UpdateProcessedBagIDs(bagNode *DependencyNode, processedBags map[string]bool) {
	for _, pk := range bagNode.PKs {
		pkStr := fmt.Sprint(pk)
		processedBags[pkStr] = true
	}
}

func (mWorker *MigrationWorker) StartBagsMigration(userID, bagAppID string, threadID int, isBlade ...bool) error {
	color.Magenta.Println("########################################################################")
	color.LightMagenta.Printf("Starting Bags for User: '%s' | App: '%s' \n", userID, bagAppID)
	color.Magenta.Println("########################################################################")

	if bags, err := db.GetBagsV2(mWorker.logTxn.DBconn, bagAppID, userID, mWorker.logTxn.Txn_id); err != nil {
		mWorker.Logger.Fatal(fmt.Sprintf("UNABLE TO FETCH BAGS FOR USER: %s | %s", userID, err))
		return err
	} else if len(bags) > 0 {
		log.Println("Bags fetched:  ", len(bags))

		bagWorker := mWorker.mThread.CreateBagWorker(userID, bagAppID, mWorker.DstAppConfig.AppID, threadID)

		log.Println(fmt.Sprintf("Bag Worker Created | %s -> %s ", bagWorker.SrcAppConfig.AppName, bagWorker.DstAppConfig.AppName))

		processedBags := make(map[string]bool)

		for _, bag := range bags {

			if strings.EqualFold(bagAppID, "4") {
				// break
			}

			bagPK := fmt.Sprint(bag["pk"])

			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			mWorker.Logger.Tracef("@StartBagsMigration > Processing Bag | ID: %v | PK: %v | %s", bag["id"], bag["pk"], bagPK)

			if _, ok := processedBags[bagPK]; ok {
				mWorker.Logger.Tracef("@StartBagsMigration > Bag Already Processed | ID: %v | PK: %v | %s", bag["id"], bag["pk"], bagPK)
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

						if !bagNode.IsEmptyExcept() {
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
		mWorker.Logger.Info("No Bags found!")
	}
	return nil
}

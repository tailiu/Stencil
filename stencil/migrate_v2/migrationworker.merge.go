package migrate_v2

import (
	"encoding/json"
	"fmt"
	"log"
	"stencil/db"
	"stencil/reference_resolution"
	"strings"

	"github.com/gookit/color"
)

// func (mWorker *MigrationWorker) MergeMemberDataWithBagsData(mmd *MappedMemberData, node *DependencyNode) error {

// 	currUID := mWorker.uid
// 	currApp := &App{
// 		Name: mWorker.SrcAppConfig.AppName,
// 		ID:   mWorker.SrcAppConfig.AppID,
// 	}

// 	for {
// 		if prevApp, prevUID, err := mWorker.GetUserIDAppIDFromPreviousMigration(currApp.ID, currUID); err != nil {
// 			mWorker.Logger.Fatal(err)
// 		} else if len(prevApp.ID) <= 0 && len(prevUID) <= 0 {
// 			break
// 		} else {
// 			for _, currMember := range mmd.SrcTables() {
// 				if err := mWorker.FetchDataFromBags(mmd, *prevApp, prevUID, currMember, *currApp); err != nil {
// 					mWorker.Logger.Fatal("@MergeMemberDataWithBagsData | ", err)
// 				}
// 			}
// 			currApp, currUID = prevApp, prevUID
// 		}
// 	}

// 	return nil
// }

// func (mWorker *MigrationWorker) FetchDataFromBags(mmd *MappedMemberData, prevApp App, prevUID string, currMember Member, currApp App) error {

// 	var currID string

// 	if mmv, ok := mmd.Data[currMember.Name+".id"]; ok {
// 		currID = fmt.Sprint(mmv.Value)
// 	} else {
// 		mWorker.Logger.Fatal("@FetchDataFromBags > id doesn't exist in table ", currMember)
// 	}

// 	attrRows, err := mWorker.GetRowsFromAttrTable(currApp.ID, currMember.ID, currID, false)
// 	if err != nil {
// 		mWorker.Logger.Fatal("@FetchDataFromBags > GetRowsFromIDTable, Unable to get AttrRows | ", currApp, currMember, currID, false, err)
// 		return err
// 	} else if len(attrRows) < 1 {
// 		return nil
// 	} else {
// 		mWorker.Logger.Trace("@FetchDataFromBags > GetRowsFromIDTable | ", attrRows)
// 	}

// }

func (mWorker *MigrationWorker) MergeBagDataWithMappedData(mmd *MappedMemberData, node *DependencyNode) error {

	prevUIDs := reference_resolution.GetPrevUserIDs(mWorker.logTxn.DBconn, mWorker.SrcAppConfig.AppID, mWorker.uid)
	if prevUIDs == nil {
		prevUIDs = make(map[string]string)
	}
	prevUIDs[mWorker.SrcAppConfig.AppID] = mWorker.uid

	for _, fromTable := range mmd.SrcTables() {
		if fromID, ok := node.Data[fromTable.Name+".id"]; ok {
			visitedRows := make(map[string]bool)
			if err := mWorker.FetchDataFromBags(visitedRows, prevUIDs, mmd, mWorker.SrcAppConfig.AppID, fromTable.ID, fromID, mmd.ToMemberID, mmd.ToMember); err != nil {
				mWorker.Logger.Fatal("@MigrateNode > FetchDataFromBags | ", err)
			}
		} else {
			mWorker.Logger.Debug(node.Data)
			mWorker.Logger.Fatal("@MigrateNode > FetchDataFromBags > id doesn't exist in table ", fromTable)
		}
	}

	// if len(toTableData) > 0 {

	// 	mappedCols := mappedData.ToCols()
	// 	for col, valWithRef := range toTableData {
	// 		if !helper.Contains(mappedCols, col) {
	// 			mmv := MappedMemberValue{
	// 				Value:  valWithRef.value,
	// 				DBConn: mWorker.logTxn.DBconn,
	// 			}
	// 			if valWithRef.ref != nil {
	// 				mmv.Ref = valWithRef.ref
	// 				mWorker.Logger.Tracef("@MigrateNode > FetchDataFromBags > Ref merged for: %s.%s\n", mappedData.ToMember, col)
	// 			}
	// 			mappedData.Data[col] = mmv
	// 			mWorker.Logger.Tracef("@MigrateNode > FetchDataFromBags > Attr merged for: '%s.%s' = '%v'", mappedData.ToMember, col, valWithRef.value)
	// 		}
	// 	}
	// 	mWorker.Logger.Tracef("@MigrateNode > FetchDataFromBags > Data merged for: %s\nData | %v", mappedData.ToMember, toTableData)
	// }

	return nil
}

func (mWorker *MigrationWorker) FetchDataFromBags(visitedRows map[string]bool, prevUIDs map[string]string, mmd *MappedMemberData, app, member string, id interface{}, dstMemberID, dstMemberName string) error {

	currentRow := fmt.Sprintf("%s:%s:%s", app, member, id)
	if _, ok := visitedRows[currentRow]; ok {
		return nil
	}
	visitedRows[currentRow] = true

	AttrRows, err := mWorker.GetRowsFromAttrTable(app, member, id, false)

	if err != nil {
		mWorker.Logger.Fatal("@FetchDataFromBags > GetRowsFromIDTable, Unable to get AttrRows | ", app, member, id, false, err)
		return err
	} else if len(AttrRows) < 1 {
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

			if mapping, found := mWorker.FetchMappingsForBag(AttrRow.FromAppName, AttrRow.FromAppID, mWorker.DstAppConfig.AppName, mWorker.DstAppConfig.AppID, AttrRow.FromMember, dstMemberName); found {

				for _, toTable := range mapping.ToTables {
					if !strings.EqualFold(toTable.Table, mmd.ToMember) {
						continue
					}
					for toAttr, fromAttr := range toTable.Mapping {

						if strings.EqualFold("id", toAttr) || strings.Contains(fromAttr, "#FETCH") || fromAttr[0:1] == "$" {
							continue
						} else if _, ok := mmd.Data[toAttr]; ok {
							continue
						}

						if mmv, err := mWorker.ResolveMappedStatement(fromAttr, bagData); err == nil {
							if mmv.Value != nil {
								if mmv.Ref != nil {
									mmv.Ref.appID = fmt.Sprint(bagRow["app"])
									mmv.Ref.mergedFromBag = true
									mWorker.Logger.Tracef("@FetchDataFromBags > REF merged for: %s.%s | %s.%s \n", toTable.Table, toAttr, AttrRow.FromAppName, mmv.FromAttr)
								}
								mmd.Data[toAttr] = *mmv
								delete(bagData, mmv.FromAttr)
								mWorker.Logger.Tracef("@FetchDataFromBags > ATTR merged for: %s.%s\n", toTable.Table, toAttr)
							}
						} else {
							mWorker.Logger.Debug(bagData)
							mWorker.Logger.Fatalf("Unable to ResolveMappedStatement | fromAttr: [%s], toAttr: [%s]", fromAttr, toAttr)
						}
					}
				}

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
			}
		}
		if err := mWorker.FetchDataFromBags(visitedRows, prevUIDs, mmd, AttrRow.FromAppID, AttrRow.FromMemberID, AttrRow.FromID, dstMemberID, dstMemberName); err != nil {
			mWorker.Logger.Fatal("@FetchDataFromBags > FetchDataFromBags: Error while recursing | ", AttrRow.FromAppID, AttrRow.FromMember, AttrRow.FromID)
			return err
		}
	}
	return nil
}

// func (mWorker *MigrationWorker) FetchDataFromBags(visitedRows map[string]bool, toTableData map[string]ValueWithReference, prevUIDs map[string]string, app, member string, id interface{}, dstMemberID, dstMemberName, toTableName string) error {

// 	currentRow := fmt.Sprintf("%s:%s:%s", app, member, id)

// 	if _, ok := visitedRows[currentRow]; ok {
// 		return nil
// 	} else {
// 		visitedRows[currentRow] = true
// 	}

// 	AttrRows, err := mWorker.GetRowsFromAttrTable(app, member, id, false)

// 	if err != nil {
// 		mWorker.Logger.Fatal("@FetchDataFromBags > GetRowsFromIDTable, Unable to get AttrRows | ", app, member, id, false, err)
// 		return err
// 	} else if len(AttrRows) < 1 {
// 		return nil
// 	} else {
// 		mWorker.Logger.Trace("@FetchDataFromBags > GetRowsFromIDTable | ", AttrRows)
// 	}

// 	for _, AttrRow := range AttrRows {

// 		bagRow, err := db.GetBagByAppMemberIDV2(mWorker.logTxn.DBconn, prevUIDs[AttrRow.FromAppID], AttrRow.FromAppID, AttrRow.FromMemberID, AttrRow.FromID, mWorker.logTxn.Txn_id)
// 		if err != nil {
// 			mWorker.Logger.Fatal("@FetchDataFromBags > GetBagByAppMemberIDV2, Unable to get bags | ", prevUIDs[AttrRow.FromAppID], AttrRow.FromAppID, AttrRow.FromMemberID, AttrRow.FromID, mWorker.logTxn.Txn_id, err)
// 			return err
// 		}
// 		if bagRow != nil {

// 			mWorker.Logger.Tracef("@FetchDataFromBags > Processing Bag | ID: %v | PK: %v | App: %v | Member: %v", bagRow["id"], bagRow["pk"], bagRow["app"], bagRow["member"])

// 			bagData := make(DataMap)
// 			if err := json.Unmarshal(bagRow["data"].([]byte), &bagData); err != nil {
// 				mWorker.Logger.Debug(bagRow["data"])
// 				mWorker.Logger.Debug(bagRow)
// 				mWorker.Logger.Fatal("@FetchDataFromBags: UNABLE TO CONVERT BAG TO MAP ", bagRow, err)
// 				return err
// 			}
// 			fmt.Printf("bag data | %v\n", bagData)

// 			if bagMappedApp, mapping, found := mWorker.FetchMappingsForBag(AttrRow.FromAppName, AttrRow.FromAppID, mWorker.DstAppConfig.AppName, mWorker.DstAppConfig.AppID, AttrRow.FromMember, dstMemberName); found {
// 				bagAppConfig, err := config.CreateAppConfig(AttrRow.FromAppName, AttrRow.FromAppID)
// 				if err != nil {
// 					log.Fatal("@FetchDataFromBags: ", err)
// 				}
// 				for _, toTable := range mapping.ToTables {
// 					if !strings.EqualFold(toTable.Table, toTableName) {
// 						continue
// 					}
// 					for toAttr, fromAttr := range toTable.Mapping {

// 						if strings.EqualFold("id", toAttr) {
// 							continue
// 						}

// 						var valueForNode interface{}
// 						var refForNode *MappingRef
// 						cleanedFromAttr := fromAttr

// 						if strings.Contains(fromAttr, "#FETCH") {
// 							continue
// 						} else if fromAttr[0:1] == "$" {
// 							if inputVal, err := bagMappedApp.GetInput(fromAttr); err == nil {
// 								valueForNode = inputVal
// 							} else {
// 								mWorker.Logger.Debugf("@FetchDataFromBags | fromAttr [%s]", fromAttr)
// 								mWorker.Logger.Fatal(err)
// 							}
// 						} else if mmv, err := mWorker.ResolveMappedStatement(fromAttr, bagData); err == nil {
// 							cleanedFromAttr = mmv.GetMemberAttr()
// 							if mmv.Value != nil {
// 								valueForNode = mmv.Value
// 								if bagAppConfig.AppID == mWorker.DstAppConfig.AppID {
// 									if bagTag, err := bagAppConfig.GetTagByMember(mmv.FromMember); err == nil {
// 										if depRefs, err := mWorker.CreateReferencesViaDependencies(bagAppConfig, *bagTag, bagData, cleanedFromAttr); err != nil {
// 											log.Println(bagData)
// 											mWorker.Logger.Fatal("@FetchDataFromBags > CreateReferencesViaDependencies: ", err)
// 										} else if len(depRefs) > 0 && mmv.Ref == nil {
// 											mmv.Ref = &depRefs[0]
// 										} else {
// 											log.Println("@FetchDataFromBags > CreateReferencesViaDependencies > No reference created: ", cleanedFromAttr)
// 										}

// 										if bagTag.Name != "root" {
// 											if ownRefs, err := mWorker.CreateReferencesViaOwnerships(bagAppConfig, *bagTag, bagData, cleanedFromAttr); err != nil {
// 												log.Println(bagData)
// 												mWorker.Logger.Fatal("@FetchDataFromBags > CreateReferencesViaOwnerships: ", err)
// 											} else if len(ownRefs) > 0 && mmv.Ref == nil {
// 												mmv.Ref = &ownRefs[0]
// 											} else {
// 												log.Println("@FetchDataFromBags > CreateReferencesViaOwnerships > No reference created: ", cleanedFromAttr)
// 											}
// 										}
// 									}
// 								}
// 								if mmv.Ref != nil {
// 									mmv.Ref.appID = fmt.Sprint(bagRow["app"])
// 									mmv.Ref.mergedFromBag = true
// 									refForNode = mmv.Ref
// 								}
// 							}
// 						} else {
// 							mWorker.Logger.Debug(bagData)
// 							mWorker.Logger.Fatalf("Unable to decode mapped val | fromAttr: [%s], toAttr: [%s]", fromAttr, toAttr)
// 						}

// 						if _, ok := toTableData[toAttr]; !ok {
// 							if valueForNode != nil {
// 								toTableData[toAttr] = ValueWithReference{value: valueForNode, ref: refForNode}
// 							}
// 						}

// 						if !strings.Contains(cleanedFromAttr, ".id") {
// 							delete(bagData, cleanedFromAttr)
// 						}
// 					}
// 				}

// 				bagAppConfig.CloseDBConns()

// 				if bagData.IsEmptyExcept() {
// 					if err := db.DeleteBagV2(mWorker.tx.StencilTx, fmt.Sprint(bagRow["pk"])); err != nil {
// 						mWorker.Logger.Fatal("@FetchDataFromBags > DeleteBagV2, Unable to delete bag | ", bagRow["pk"])
// 						return err
// 					} else {
// 						log.Println(fmt.Sprintf("%s | PK: %v", color.FgLightRed.Render("Deleted BAG"), bagRow["pk"]))
// 					}
// 				} else {
// 					log.Println(fmt.Sprintf("%s | %v", color.FgYellow.Render("BAG NOT EMPTY"), bagData))
// 					if jsonData, err := json.Marshal(bagData); err == nil {
// 						if err := db.UpdateBag(mWorker.tx.StencilTx, fmt.Sprint(bagRow["pk"]), mWorker.logTxn.Txn_id, jsonData); err != nil {
// 							mWorker.Logger.Fatal("@FetchDataFromBags: UNABLE TO UPDATE BAG ", bagRow, err)
// 							return err
// 						} else {
// 							log.Println(fmt.Sprintf("%s | PK: %v", color.FgLightYellow.Render("Updated BAG"), bagRow["pk"]))
// 						}
// 					} else {
// 						mWorker.Logger.Fatal("@FetchDataFromBags > len(bagData) != 0, Unable to marshall bag | ", bagData)
// 						return err
// 					}
// 				}
// 			}
// 		}
// 		if err := mWorker.FetchDataFromBags(visitedRows, toTableData, prevUIDs, AttrRow.FromAppID, AttrRow.FromMemberID, AttrRow.FromID, dstMemberID, dstMemberName, toTableName); err != nil {
// 			mWorker.Logger.Fatal("@FetchDataFromBags > FetchDataFromBags: Error while recursing | ", toTableData, AttrRow.FromAppID, AttrRow.FromMember, AttrRow.FromID)
// 			return err
// 		}
// 	}
// 	return nil
// }

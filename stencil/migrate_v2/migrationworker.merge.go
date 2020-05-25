package migrate_v2

import (
	"fmt"
	"strings"
)

func (mWorker *MigrationWorker) MergeBagDataWithMappedData(mmd *MappedMemberData, node *DependencyNode, bagManager *BagManager) error {

	for _, fromTable := range mmd.SrcTables() {
		if fromID, ok := node.Data[fromTable.Name+".id"]; ok {
			bagManager.ClearVisitedRows()
			if err := mWorker.FetchDataFromBags(bagManager, mmd, mWorker.SrcAppConfig.AppID, fromTable.ID, fromID, mmd.ToMemberID, mmd.ToMember); err != nil {
				mWorker.Logger.Fatal("@MigrateNode > FetchDataFromBags | ", err)
			}
		} else {
			mWorker.Logger.Debug(node.Data)
			mWorker.Logger.Fatal("@MigrateNode > FetchDataFromBags > id doesn't exist in table ", fromTable)
		}
	}

	return nil
}

func (mWorker *MigrationWorker) FetchDataFromBags(bagManager *BagManager, mmd *MappedMemberData, app, member string, id interface{}, dstMemberID, dstMemberName string) error {

	if bagManager.IsRowVisited(app, member, fmt.Sprint(id)) {
		return nil
	}

	if attrRows, err := mWorker.GetRowsFromAttrTable(app, member, id, false); err != nil {
		mWorker.Logger.Fatal("@FetchDataFromBags > GetRowsFromIDTable, Unable to get attrRows | ", app, member, id, false, err)
		return err
	} else {
		mWorker.Logger.Trace("Fetched attrRows | ", attrRows)
		for _, attrRow := range attrRows {
			mWorker.Logger.Trace("Current attrRow | ", attrRow)
			if bag := bagManager.GetBagsFromDB(attrRow.FromAppID, attrRow.FromMemberID, attrRow.FromID, mWorker.logTxn.Txn_id); bag != nil {
				mWorker.Logger.Tracef("Processing Bag | ID: %v | PK: %v | App: %v | Member: %v\nBag Data | %v\n", bag.ID, bag.PK, bag.AppID, bag.MemberID, bag.Data)
				if mapping, found := mWorker.FetchMappingsForBag(attrRow.FromAppName, attrRow.FromAppID, mWorker.DstAppConfig.AppName, mWorker.DstAppConfig.AppID, attrRow.FromMember, dstMemberName); found {
					mWorker.Logger.Tracef("Mapping found | %s(%s) : %s -> %s(%s) : %s \n", attrRow.FromAppName, attrRow.FromAppID, attrRow.FromMember, mWorker.DstAppConfig.AppName, mWorker.DstAppConfig.AppID, dstMemberName)
					mWorker.visitedBags.MarkAsVisited(bag)
					for _, toTable := range mapping.ToTables {
						if !strings.EqualFold(toTable.Table, mmd.ToMember) {
							continue
						}
						for toAttr, mappedStmt := range toTable.Mapping {
							if strings.EqualFold("id", toAttr) || strings.Contains(mappedStmt, "#FETCH") || mappedStmt[0:1] == "$" {
								continue
							}
							if mmv, err := mWorker.ResolveMappedStatement(mappedStmt, bag.Data, bag.AppID); err == nil {
								if mmv != nil && mmv.Value != nil {
									bag.AddAttrtoRemove(mmv.GetMemberAttr())
									if existingMMV, ok := mmd.Data[toAttr]; ok {
										if existingMMV.Value != nil {
											mWorker.Logger.Tracef("ATTR exists in node: %s.%s | '%s' \n", toTable.Table, toAttr, mmv.GetMemberAttr())
											continue
										}
									}
									if mmv.Ref != nil {
										mmv.Ref.appID = bag.AppID
										mmv.Ref.mergedFromBag = true
										mWorker.Logger.Tracef("REF merged for: %s.%s | %s.%s \n", toTable.Table, toAttr, attrRow.FromAppName, mmv.GetMemberAttr())
									}
									mmd.Data[toAttr] = *mmv
									mWorker.Logger.Tracef("ATTR merged for: %s.%s | '%s' \n", toTable.Table, toAttr, mmv.GetMemberAttr())
								}
							} else {
								mWorker.Logger.Debug(err)
								mWorker.Logger.Debug(bag.Data)
								mWorker.Logger.Fatalf("Unable to ResolveMappedStatement | mappedStmt: [%s], toAttr: [%s]", mappedStmt, toAttr)
							}
						}
					}
				} else {
					mWorker.Logger.Tracef("Mapping not found | %s(%s) : %s -> %s(%s) : %s \n", attrRow.FromAppName, attrRow.FromAppID, attrRow.FromMember, mWorker.DstAppConfig.AppName, mWorker.DstAppConfig.AppID, dstMemberName)
				}
			}
			if err := mWorker.FetchDataFromBags(bagManager, mmd, attrRow.FromAppID, attrRow.FromMemberID, attrRow.FromID, dstMemberID, dstMemberName); err != nil {
				mWorker.Logger.Fatal("@FetchDataFromBags > FetchDataFromBags: Error while recursing | ", attrRow.FromAppID, attrRow.FromMember, attrRow.FromID)
				return err
			}
		}
	}
	return nil
}

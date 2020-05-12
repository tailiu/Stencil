package migrate_v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"stencil/db"
	"stencil/helper"

	"github.com/gookit/color"
)

func (mWorker *MigrationWorker) ConstructBagNode(bagMember, bagMemberID, bagRowID string) (*DependencyNode, error) {

	bagTag, err := mWorker.SrcAppConfig.GetTagByMember(bagMember)
	if err != nil {
		mWorker.Logger.Fatal(fmt.Sprintf("@ConstructBagNode: UNABLE TO GET BAG TAG BY MEMBER : %s | %s", bagMember, err))
		return nil, err
	}

	if ql := mWorker.GetTagQLForBag(*bagTag); len(ql) > 0 {
		where := fmt.Sprintf(" WHERE \"%s\".member = %s AND \"%s\".id = %s", bagMember, bagMemberID, bagMember, bagRowID)
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

func (mWorker *MigrationWorker) GetRowsFromAttrTable(app, member string, id interface{}, getFrom bool) ([]AttrRow, error) {

	var AttrRows []AttrRow
	var err error
	var AttrRowsDB []map[string]interface{}

	if !getFrom {
		AttrRowsDB, err = db.GetRowsFromAttrTableByTo(mWorker.logTxn.DBconn, app, member, helper.GetInt64(id))
	} else {
		AttrRowsDB, err = db.GetRowsFromAttrTableByFrom(mWorker.logTxn.DBconn, app, member, helper.GetInt64(id))
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

func (bagWorker *MigrationWorker) CreateBagStruct(bagInfo DataMap) (Bag, error) {

	var bag Bag

	bag.ID = fmt.Sprint(bagInfo["id"])
	bag.PK = fmt.Sprint(bagInfo["pk"])
	bag.UID = fmt.Sprint(bagInfo["user_id"])
	bag.MemberID = fmt.Sprint(bagInfo["member"])
	bag.AppID = fmt.Sprint(bagInfo["app"])

	if member, err := db.TableName(bagWorker.logTxn.DBconn, bag.MemberID, bag.AppID); err != nil {
		bagWorker.Logger.Fatal("@CreateBagStruct > Table Name: ", err)
		return bag, err
	} else {
		bag.Member = member
	}

	if node, err := bagWorker.ConstructBagNode(bag.Member, bag.MemberID, bag.ID); err != nil {
		bagWorker.Logger.Fatal("@CreateBagStruct > ConstructBagNode: ", err)
		return bag, err
	} else {
		bag.Node = node
	}

	return bag, nil
}

func (mWorker *MigrationWorker) CheckRawBag(node *DependencyNode) (bool, error) {
	for _, table := range node.Tag.Members {
		if tableID, err := db.TableID(mWorker.logTxn.DBconn, table, mWorker.SrcAppConfig.AppID); err == nil {
			if id, ok := node.Data[table+".id"]; ok {
				if idRows, err := mWorker.GetRowsFromAttrTable(mWorker.SrcAppConfig.AppID, tableID, id, true); err == nil {
					if len(idRows) == 0 {
						return true, nil
					}
				} else {
					mWorker.Logger.Debug(node.Data)
					mWorker.Logger.Fatal("@CheckRawBag > GetRowsFromIDTable > ", mWorker.SrcAppConfig.AppID, tableID, id, err)
				}
			} else {
				mWorker.Logger.Debug(node.Data)
				mWorker.Logger.Warn("@CheckRawBag > id doesn't exist in table ", table+".id")
			}
		} else {
			mWorker.Logger.Fatal("@CheckRawBag > TableID, fromTable: error in getting table id for member! ", table, err)
		}
	}
	return false, nil
}

func (bagWorker *MigrationWorker) HandleBagDeletion(node *DependencyNode) error {

	if bagWorker.mtype != BAGS {
		bagWorker.Logger.Fatal("Shouldn't be here. It's Bag Deletion Method, being called by non-bag deletion migration.")
		return nil
	}

	if !node.IsEmptyExcept() {
		log.Printf("%s { %s } | PKs: %v \n Data | %v", color.FgLightYellow.Render("HANDLING LEFTOVER DATA"), node.Tag.Name, node.PKs, node.Data)
		if err := bagWorker.SendNodeToBagWithOwnerID(node, bagWorker.uid); err != nil {
			bagWorker.Logger.Debug("Data | ", node.Data)
			bagWorker.Logger.Fatalf("@HandleBagDeletion > SendNodeToBagWithOwnerID | { %s } | PK: %v", node.Tag.Name, node.PKs)
			return err
		} else {
			log.Printf("%s { %s } | PK: %v  \n", color.FgLightYellow.Render("BAGGED LEFTOVER DATA"), node.Tag.Name, node.PKs)
		}
	}

	if err := bagWorker.DeleteBag(node); err != nil {
		bagWorker.Logger.Fatal("@HandleBagDeletion > DeleteBag:", err)
		return err
	}

	return nil
}

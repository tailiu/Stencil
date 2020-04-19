package migrate_v2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"stencil/db"
	"stencil/helper"
	"stencil/reference_resolution"
	"strings"

	"github.com/gookit/color"
)

func (bm *BagManager) Init(dbConn *sql.DB, srcAppID, srcUID string) {
	bm.Bags = make(map[string]*DBBag)
	bm.VisitedRows = make(map[string]bool)
	bm.DBConn = dbConn
	bm.PrevUIDs = reference_resolution.GetPrevUserIDs(bm.DBConn, srcAppID, srcUID)
	if bm.PrevUIDs == nil {
		bm.PrevUIDs = make(map[string]string)
	}
	bm.PrevUIDs[srcAppID] = srcUID
}

func (bm *BagManager) ClearVisitedRows() {
	bm.VisitedRows = make(map[string]bool)
}

func (bm *BagManager) IsRowVisited(app, member, id string) bool {
	currentRow := fmt.Sprintf("%s:%s:%s", app, member, id)
	if _, ok := bm.VisitedRows[currentRow]; ok {
		return true
	}
	bm.VisitedRows[currentRow] = true
	return false
}

func (bm *BagManager) GetBagsFromDB(fromAppID, fromMemberID string, fromID int64, txnID int) *DBBag {
	if bagRow, err := db.GetBagByAppMemberIDV2(bm.DBConn, bm.PrevUIDs[fromAppID], fromAppID, fromMemberID, fromID, txnID); err != nil {
		log.Fatal("@GetBagsFromDB > GetBagByAppMemberIDV2, Unable to get bags | ", bm.PrevUIDs[fromAppID], fromAppID, fromMemberID, fromID, err)
	} else {
		if bag, ok := bm.Bags[fmt.Sprint(bagRow["pk"])]; ok {
			return bag
		} else {
			bag = &DBBag{
				PK:       fmt.Sprint(bagRow["pk"]),
				ID:       fmt.Sprint(bagRow["id"]),
				UID:      fmt.Sprint(bagRow["user_id"]),
				AppID:    fmt.Sprint(bagRow["app"]),
				MemberID: fmt.Sprint(bagRow["member"]),
				TxnID:    fmt.Sprint(bagRow["migration_id"]),
			}

			bagData := make(DataMap)
			if err := json.Unmarshal(bagRow["data"].([]byte), &bagData); err != nil {
				fmt.Println(bag)
				fmt.Println(bagRow["data"])
				log.Fatal("@BM.GetBagsFromDB: UNABLE TO CONVERT BAG TO MAP | ", err)
			}
			bag.Data = bagData
			bm.Bags[bag.PK] = bag
			return bm.Bags[bag.PK]
		}
	}
	return nil
}

func (bm *BagManager) UpdateBags(stenctilTx *sql.Tx, txnID int) {
	for _, bag := range bm.Bags {
		if bag.RemoveAttrs() {
			if bag.Data.IsEmptyExcept() {
				if err := db.DeleteBagV2(stenctilTx, bag.PK); err != nil {
					log.Fatalf("@UpdateBags > DeleteBagV2, Unable to delete bag | %s | %s", bag.PK, err)
				} else {
					log.Println(fmt.Sprintf("%s | PK: %v", color.FgLightRed.Render("Deleted BAG"), bag.PK))
				}
			} else {
				log.Println(fmt.Sprintf("%s | %v", color.FgYellow.Render("BAG NOT EMPTY"), bag.Data))
				if jsonData, err := json.Marshal(bag.Data); err == nil {
					if err := db.UpdateBag(stenctilTx, bag.PK, txnID, jsonData); err != nil {
						log.Fatalf("@UpdateBags: UNABLE TO UPDATE BAG | %s | %s", bag.PK, err)
					} else {
						log.Println(fmt.Sprintf("%s | PK: %v", color.FgLightYellow.Render("Updated BAG"), bag.PK))
						fmt.Println("Updated Bag Data | ", bag.Data)
					}
				} else {
					log.Fatal("@UpdateBags > len(bag.Data) != 0, Unable to marshall bag | ", bag.Data)
				}
			}
		}
	}
}

func (bag *DBBag) RemoveAttrs() bool {
	if len(bag.AttrsToRemove) == 0 {
		return false
	}
	for _, attr := range bag.AttrsToRemove {
		log.Println("Deleting attr from bag data: ", attr)
		delete(bag.Data, attr)
	}
	return true
}

func (bag *DBBag) AddAttrtoRemove(attr string) {
	if !strings.Contains(attr, ".id") {
		if !helper.Contains(bag.AttrsToRemove, attr) {
			bag.AttrsToRemove = append(bag.AttrsToRemove, attr)
		}
	}
}

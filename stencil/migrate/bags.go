package migrate

import (
	"encoding/json"
	"fmt"
	"log"
	"stencil/db"
	"stencil/reference_resolution"
)

func ConstructBagNode() {

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

			srcMemberID := fmt.Sprint(bag["member"])
			srcMemberName, err := db.TableName(bagWorker.logTxn.DBconn, srcMemberID, bagAppID)

			if err != nil {
				log.Fatal("@MigrateBags > Table Name: ", err)
			}

			bagID := fmt.Sprint(bag["pk"])

			log.Println(fmt.Sprintf("~%2d~ Current    Bag: { %s } | ID: %s, App: %s ", threadID, srcMemberName, bagID, bagAppID))

			bagData := make(map[string]interface{})

			if err := json.Unmarshal(bag["data"].([]byte), &bagData); err != nil {
				fmt.Println("BAG >> ", bag)
				log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO CONVERT BAG TO MAP: %s | %s", threadID, bagWorker.uid, err))
				return err
			}

			bagTag, err := bagWorker.SrcAppConfig.GetTagByMember(srcMemberName)
			if err != nil {
				log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO GET BAG TAG BY MEMBER : %s | %s", threadID, srcMemberName, err))
				return err
			}

			if err := bagWorker.InitTransactions(); err != nil {
				return err
			}

			toCommit := true

			bagNode := DependencyNode{Tag: *bagTag, Data: bagData}
			if err := bagWorker.HandleMigration(&bagNode, true); err != nil {
				toCommit = false
				// fmt.Println(bag)
				log.Println(fmt.Sprintf("x%2dx UNABLE TO MIGRATE BAG { %s } | ID: %s | %s ", threadID, bagTag.Name, bagID, err))
			} else {
				log.Println(fmt.Sprintf("x%2dx MIGRATED bag { %s } | ID: %s", threadID, bagTag.Name, bagID))

				if bagWorker.IsNodeDataEmpty(&bagNode) {
					if err := db.DeleteBagV2(bagWorker.tx.StencilTx, bagID); err != nil {
						toCommit = false
						fmt.Println(bag)
						log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO DELETE BAG : %s | %s ", threadID, bagID, err))
						// return err
					} else {
						log.Println(fmt.Sprintf("x%2dx DELETED bag { %s } ", threadID, bagTag.Name))
					}
				} else {
					if jsonData, err := json.Marshal(bagNode.Data); err == nil {
						if err := db.UpdateBag(bagWorker.tx.StencilTx, bagID, bagWorker.logTxn.Txn_id, jsonData); err != nil {
							toCommit = false
							fmt.Println(bag)
							log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO UPDATE BAG : %s | %s ", threadID, bagID, err))
							// return err
						} else {
							log.Println(fmt.Sprintf("x%2dx UPDATED bag { %s } ", threadID, bagTag.Name))
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
					log.Fatal(fmt.Sprintf("x%2dx UNABLE TO COMMIT bag { %s } | %s ", threadID, bagTag.Name, err))
					// return err
				} else {
					log.Println(fmt.Sprintf("x%2dx COMMITTED bag { %s } ", threadID, bagTag.Name))
				}
			} else {
				log.Println(fmt.Sprintf("x%2dx ROLLBACK bag { %s } | ID: %s ", threadID, bagTag.Name, bagID))
				bagWorker.tx.DstTx.Rollback()
				bagWorker.tx.SrcTx.Rollback()
				bagWorker.tx.StencilTx.Rollback()
			}
		}

		bagWorker.CloseDBConns()
	}

	return nil
}

package migrate_v2

import (
	"fmt"
	"log"
	"stencil/db"
	"strings"

	"github.com/gookit/color"
)

func (mWorker *MigrationWorker) CallMigration(node *DependencyNode, threadID int) error {

	if ownerID, isRoot := mWorker.GetNodeOwner(node); isRoot && len(ownerID) > 0 {
		log.Println(fmt.Sprintf("OWNED   node { %s } | root [%s] : owner [%s]", node.Tag.Name, mWorker.uid, ownerID))
		if err := mWorker.InitTransactions(); err != nil {
			return err
		} else {
			defer mWorker.RollbackTransactions()
		}

		if !strings.EqualFold(node.Tag.Name, "root") {
			log.Println(fmt.Sprintf("CHECKING NEXT NODES { %s }", node.Tag.Name))
			if err := mWorker.CheckNextNode(node); err != nil {
				return err
			}
		}

		log.Println(fmt.Sprintf("HANDLING MIGRATION { %s }", node.Tag.Name))

		if migrated, err := mWorker.HandleMigration(node); err == nil {
			if err := mWorker.HandleNodeDeletion(node, true); err != nil {
				mWorker.Logger.Fatal(err)
			}
			if migrated {
				log.Println(fmt.Sprintf("%s  node { %s } ", color.FgLightGreen.Render("Migrated"), node.Tag.Name))
			} else {
				log.Println(fmt.Sprintf("%s  node { %s } ", color.FgGreen.Render("Not Migrated / No Err"), node.Tag.Name))
			}
		} else {
			if strings.EqualFold(err.Error(), "3") {
				log.Println(fmt.Sprintf("UNMAPPED  node { %s } ", node.Tag.Name))
			} else if strings.EqualFold(err.Error(), "2") {
				log.Println(fmt.Sprintf("%s  node { %s } ", color.FgLightYellow.Render("BAGGED"), node.Tag.Name))
			} else {
				log.Println(fmt.Sprintf("FAILED    node { %s } ", node.Tag.Name))
				if strings.EqualFold(err.Error(), "0") {
					log.Println(err)
					return err
				}
				return err
			}
		}

		if err := mWorker.CommitTransactions(); err != nil {
			return err
		} else {
			log.Println(fmt.Sprintf("COMMITTED node { %s } ", node.Tag.Name))
		}
	} else {
		log.Println(fmt.Sprintf("VISITED  node { %s } | root [%s] : owner [%s] | root: %v", node.Tag.Name, mWorker.uid, ownerID, isRoot))
		mWorker.visitedNodes.MarkAsVisited(node)
	}
	fmt.Println("------------------------------------------------------------------------")
	return nil
}

func (mWorker *MigrationWorker) CallMigrationX(node *DependencyNode, threadID int) error {
	if ownerID, isRoot := mWorker.GetNodeOwner(node); isRoot && len(ownerID) > 0 {
		if err := mWorker.InitTransactions(); err != nil {
			return err
		} else {
			defer mWorker.tx.SrcTx.Rollback()
			defer mWorker.tx.DstTx.Rollback()
			defer mWorker.tx.StencilTx.Rollback()
		}
		if migrated, err := mWorker.HandleMigration(node); err == nil {
			if migrated {
				log.Println(fmt.Sprintf("%s  node { %s } ", color.FgLightGreen.Render("Migrated"), node.Tag.Name))
			} else {
				log.Println(fmt.Sprintf("%s  node { %s } ", color.FgGreen.Render("Not Migrated / No Err"), node.Tag.Name))
			}
		} else {
			log.Println(fmt.Sprintf("RCVD ERR  node { %s } ", node.Tag.Name), err)

			if strings.EqualFold(err.Error(), "3") {
				log.Println(fmt.Sprintf("IGNORED   node { %s } ", node.Tag.Name))
			} else if strings.EqualFold(err.Error(), "2") {
				log.Println(fmt.Sprintf("BAGGED?   node { %s } ", node.Tag.Name))
			} else {
				log.Println(fmt.Sprintf("FAILED    node { %s } ", node.Tag.Name), err)
				if strings.EqualFold(err.Error(), "0") {
					log.Println(err)
					return err
				}
				if strings.Contains(err.Error(), "deadlock") {
					return err
				}
			}
		}
		if err := mWorker.CommitTransactions(); err != nil {
			mWorker.Logger.Fatal(fmt.Sprintf("UNABEL to COMMIT node { %s } ", node.Tag.Name))
			return err
		} else {
			log.Println(fmt.Sprintf("COMMITTED node { %s } ", node.Tag.Name))
		}
	} else {
		log.Println(fmt.Sprintf("VISITED  node { %s } | root [%s] : owner [%s]", node.Tag.Name, mWorker.uid, ownerID))
	}
	mWorker.visitedNodes.MarkAsVisited(node)
	return nil
}

func (mWorker *MigrationWorker) CallBagsMigration(userID, bagAppID string, threadID int) error {

	if bagRows, err := db.GetBagsV2(mWorker.logTxn.DBconn, bagAppID, userID, mWorker.logTxn.Txn_id); err != nil {
		mWorker.Logger.Fatalf("UNABLE TO FETCH BAGS FOR USER: %s | %s", userID, err)
		return err
	} else if len(bagRows) > 0 {
		fmt.Println("\n\n========================================================================")
		mWorker.Logger.Infof("%s | User: '%s' | App: '%s' | Bags Count: '%v' \n", color.LightMagenta.Render("Starting Bag Migration"), userID, bagAppID, len(bagRows))

		bagWorker := mWorker.mThread.CreateBagWorker(userID, bagAppID, mWorker.DstAppConfig.AppID, mWorker)
		defer bagWorker.CloseDBConns()

		for _, bagRow := range bagRows {

			dbBag := &DBBag{
				PK:       fmt.Sprint(bagRow["pk"]),
				ID:       fmt.Sprint(bagRow["id"]),
				UID:      fmt.Sprint(bagRow["user_id"]),
				AppID:    fmt.Sprint(bagRow["app"]),
				MemberID: fmt.Sprint(bagRow["member"]),
				TxnID:    fmt.Sprint(bagRow["migration_id"]),
			}

			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			if bagWorker.visitedBags.PKExists(dbBag.PK) {
				mWorker.Logger.Infof("@CallBagsMigration: Bag Already Processed | ID: %v | PK: %v \n", dbBag.ID, dbBag.PK)
				continue
			} else if bagWorker.visitedBags.IsVisited(dbBag) {
				mWorker.Logger.Infof("@CallBagsMigration: Bag already visited | ID: %v | PK: %v | App: %v | Member: %v\nBag Data | %v\n", dbBag.ID, dbBag.PK, dbBag.AppID, dbBag.MemberID, dbBag.Data)
				continue
			}

			if bagStruct, err := bagWorker.CreateBagStruct(bagRow); err == nil {

				bagWorker.visitedBags.UpdatePKs(bagStruct.Node)

				bagWorker.Logger.Infof("Processing Bag | %s | ID: %v | PK: %v \n Data | %v \n", bagStruct.Node.Tag.Name, bagStruct.ID, bagStruct.PK, bagStruct.Node.Data)

				if err := bagWorker.InitTransactions(); err != nil {
					return err
				}

				if migrated, err := bagWorker.HandleMigration(bagStruct.Node); err != nil {
					log.Printf("%s { %s } | ROLLBACK | PK: %s | ID: %v | Owner: %s | %s \n", color.FgGreen.Render("BAG NOT MIGRATED"), bagStruct.Node.Tag.Name, bagStruct.PK, bagStruct.ID, bagStruct.UID, color.FgYellow.Render(err))
					bagWorker.RollbackTransactions()
					continue
				} else if migrated {
					if err := bagWorker.HandleBagDeletion(bagStruct.Node); err != nil {
						bagWorker.Logger.Fatal(err)
					}
					log.Printf("%s { %s } | PK: %s \n", color.FgLightGreen.Render("BAG MIGRATED"), bagStruct.Node.Tag.Name, bagStruct.PK)
				} else {
					log.Printf("%s { %s } | PK: %s | ID: %s \n", color.FgGreen.Render("BAG NOT MIGRATED; NO ERR"), bagStruct.Node.Tag.Name, bagStruct.PK, bagStruct.ID)
				}

				if err := bagWorker.CommitTransactions(); err != nil {
					bagWorker.Logger.Fatal(fmt.Sprintf("UNABLE TO COMMIT bag { %s } | %s ", bagStruct.Node.Tag.Name, err))
					return err
				}
				log.Println(fmt.Sprintf("COMMITTED bag { %s } | PK: %s", bagStruct.Node.Tag.Name, bagStruct.PK))
			} else {
				bagWorker.Logger.Fatal(fmt.Sprintf("UNABLE TO CREATE bag struct | bagMemberID:%s  bagRowID:%s | %s", bagRow["member"], bagRow["id"], err))
			}
		}

	} else {
		mWorker.Logger.Infof("%s | User: '%s' | App: '%s' \n", color.LightMagenta.Render("No Bags Found"), userID, bagAppID)
	}
	return nil
}

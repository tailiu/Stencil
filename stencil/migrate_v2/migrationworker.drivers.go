package migrate_v2

import (
	"fmt"
	"log"
	"strings"

	"github.com/gookit/color"
)

func (self *MigrationWorker) DeletionMigration(node *DependencyNode, threadID int) error {

	rootTagName := "root"

	if nodeIDAttr, err := node.Tag.ResolveTagAttr("id"); err != nil {
		self.Logger.Fatalf("node.Tag.ResolveTagAttr(id): %s", err)
	} else {

		if strings.EqualFold(node.Tag.Name, rootTagName) {
			log.Println(fmt.Sprintf("Current   Node { %s } | ID: %v ", color.FgLightCyan.Render(node.Tag.Name), node.Data[nodeIDAttr]))
			if err := self.CallMigration(node, threadID); err != nil {
				return err
			}
		}

		for {
			if adjNode, err := self.GetAdjNode(node, threadID); err != nil {
				return err
			} else {
				if adjNode == nil {
					break
				}
				if adjNodeIDAttr, err := adjNode.Tag.ResolveTagAttr("id"); err != nil {
					self.Logger.Fatalf("adjNode.Tag.ResolveTagAttr(id): %s", err)
				} else {
					log.Println(fmt.Sprintf("Current   Node { %s } | ID: %v ", color.FgLightCyan.Render(node.Tag.Name), node.Data[nodeIDAttr]))
					log.Println(fmt.Sprintf("Adjacent  Node { %s } | ID: %v ", color.FgLightCyan.Render(adjNode.Tag.Name), adjNode.Data[adjNodeIDAttr]))
					if err := self.DeletionMigration(adjNode, threadID); err != nil {
						self.Logger.Fatal(fmt.Sprintf("ERROR! NODE { %s } | ID: %v, ADJ NODE : { %s } | ID: %v | err: [ %s ]", node.Tag.Name, node.Data[nodeIDAttr], adjNode.Tag.Name, adjNode.Data[adjNodeIDAttr], err))
						return err
					}
				}
			}
		}

		log.Println(fmt.Sprintf("PROCESS Node { %s } ", color.FgLightCyan.Render(node.Tag.Name)))

		if strings.EqualFold(node.Tag.Name, rootTagName) {
			return self.DeleteRoot(threadID)
		} else {
			if err := self.CallMigration(node, threadID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (self *MigrationWorker) ConsistentMigration(threadID int) error {

	if err := self.CallMigrationX(self.Root, threadID); err != nil {
		return err
	}

	for {
		if node, err := self.GetOwnedNode(threadID); err == nil {
			if node == nil {
				return nil
			}
			nodeIDAttr, _ := node.Tag.ResolveTagAttr("id")
			log.Println(fmt.Sprintf("~%2d~ | Current   Node: { %s } ID: %v", threadID, node.Tag.Name, node.Data[nodeIDAttr]))
			if err := self.CallMigrationX(node, threadID); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	log.Println(self.mtype, " MIGRATION DONE!")

	return nil
}

func (self *MigrationWorker) IndependentMigration(threadID int) error {

	return self.ConsistentMigration(threadID)
}

func (self *MigrationWorker) NaiveMigration(threadID int) error {

	if err := self.CallMigrationX(self.Root, threadID); err != nil {
		return err
	}

	for {
		if node, err := self.GetOwnedNode(threadID); err == nil {
			if node == nil {
				break
			}
			nodeIDAttr, _ := node.Tag.ResolveTagAttr("id")
			log.Println(fmt.Sprintf("~%2d~ | Current   Node: { %s } ID: %v", threadID, node.Tag.Name, node.Data[nodeIDAttr]))
			if err := self.CallMigrationX(node, threadID); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if err := self.DeleteRoot(threadID); err != nil {
		log.Println(fmt.Sprintf("~%2d~ | Root not deleted!", threadID))
		log.Fatal(err)
	}

	log.Println("NAIVE MIGRATION DONE!")

	return nil
}

func (mWorker *MigrationWorker) BagsMigration(threadID int) error {

	bagUID := mWorker.uid
	bagApp := &App{
		Name: mWorker.SrcAppConfig.AppName,
		ID:   mWorker.SrcAppConfig.AppID,
	}

	for {
		if err := mWorker.CallBagsMigration(bagUID, bagApp.ID, threadID); err != nil {
			mWorker.Logger.Fatal(err)
		}
		mWorker.Logger.Infof("BagWorker finished | Thread # %v | App: '%s', UID: '%s' \n", threadID, bagApp.ID, bagUID)
		var prevIDErr error
		if bagApp, bagUID, prevIDErr = mWorker.GetUserIDAppIDFromPreviousMigration(bagApp.ID, bagUID); prevIDErr != nil {
			mWorker.Logger.Fatal(prevIDErr)
		} else if bagApp == nil && len(bagUID) <= 0 {
			mWorker.Logger.Info("Bag Migration Finished!")
			break
		}
	}

	return nil
}

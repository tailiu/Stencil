package migrate_v1

import (
	"fmt"
	"log"
)

func (self *MigrationWorkerV2) NaiveMigration(threadID int) error {

	if err := self.CallMigrationX(self.root, threadID); err != nil {
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

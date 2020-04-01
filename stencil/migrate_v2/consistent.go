package migrate_v2

import (
	"fmt"
	"log"
)

func (self *MigrationWorkerV2) ConsistentMigration(threadID int) error {

	if err := self.CallMigrationX(self.root, threadID); err != nil {
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

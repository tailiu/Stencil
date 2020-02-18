package migrate

import (
	"fmt"
	"log"
	"strings"

	"github.com/gookit/color"
)

func (self *MigrationWorkerV2) DeletionMigration(node *DependencyNode, threadID int) error {

	nodeIDAttr, _ := node.Tag.ResolveTagAttr("id")

	if strings.EqualFold(node.Tag.Name, "root") {
		log.Println(fmt.Sprintf("Current   Node { %s } | ID: %v ", color.FgLightYellow.Render(node.Tag.Name), node.Data[nodeIDAttr]))
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
			adjNodeIDAttr, _ := adjNode.Tag.ResolveTagAttr("id")
			log.Println(fmt.Sprintf("Current   Node { %s } | ID: %v ", color.FgLightYellow.Render(node.Tag.Name), node.Data[nodeIDAttr]))
			log.Println(fmt.Sprintf("Adjacent  Node { %s } | ID: %v ", color.FgLightYellow.Render(adjNode.Tag.Name), adjNode.Data[adjNodeIDAttr]))
			if err := self.DeletionMigration(adjNode, threadID); err != nil {
				self.Logger.Fatal(fmt.Sprintf("ERROR! NODE { %s } | ID: %v, ADJ_NODE : { %s } | ID: %v | err: [ %s ]", node.Tag.Name, node.Data[nodeIDAttr], adjNode.Tag.Name, adjNode.Data[adjNodeIDAttr], err))
				return err
			}
		}
	}

	log.Println(fmt.Sprintf("PROCESS Node { %s } ", color.FgLightYellow.Render(node.Tag.Name)))

	if strings.EqualFold(node.Tag.Name, "root") {
		return self.DeleteRoot(threadID)
	} else {
		if err := self.CallMigration(node, threadID); err != nil {
			return err
		}
	}

	return nil
}

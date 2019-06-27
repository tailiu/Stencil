package migrate

import (
	"strings"
)

func (self InvalidList) Exists(node DependencyNode) bool {
	for _, iNode := range self.Nodes {
		if strings.EqualFold(node.Tag.Name, iNode.Tag.Name) {
			if idAttr, err := node.Tag.ResolveTagAttr("id"); err == nil {
				if iNodeVal, iNodeErr := iNode.GetValueForKey(idAttr); iNodeErr == nil {
					if nodeVal, nodeErr := node.GetValueForKey(idAttr); nodeErr == nil {
						if iNodeVal == nodeVal {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func (self *InvalidList) Add(node *DependencyNode) {
	iNode := node
	// if idAttr, err := iNode.Tag.ResolveTagAttr("id"); err == nil {
	// 	for key := range iNode.Data {
	// 		if !strings.EqualFold(key, idAttr) {
	// 			delete(iNode.Data, key)
	// 		}
	// 	}
	// }
	self.Nodes = append(self.Nodes, iNode)
}

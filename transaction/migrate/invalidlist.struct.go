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

func (self *InvalidList) Add(node DependencyNode) {
	node.SQL = ""
	if idAttr, err := node.Tag.ResolveTagAttr("id"); err == nil {
		for key := range node.Data {
			if !strings.EqualFold(key, idAttr) {
				delete(node.Data, key)
			}
		}
	}
	self.Nodes = append(self.Nodes, node)
}
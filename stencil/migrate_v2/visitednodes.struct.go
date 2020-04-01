package migrate_v2

import (
	"fmt"
	"log"
	config "stencil/config/v2"
	"strings"
)

// Init : Initializes nodes map
func (vNodes *VisitedNodes) Init() {
	vNodes.Nodes = make(map[string]map[string]bool)
}

// GetIDsByTag : Fetch ids of visited nodes by tag member
func (vNodes VisitedNodes) GetIDsByTag(tagMember string) (map[string]bool, bool) {
	res, ok := vNodes.Nodes[tagMember]
	return res, ok
}

// CheckIfIDExistsByTag : Check if the passed ID has been visited
func (vNodes VisitedNodes) CheckIfIDExistsByTag(tagMember, id string) bool {
	if _, ok := vNodes.Nodes[tagMember]; ok {
		_, found := vNodes.Nodes[tagMember]
		return found
	}
	return false
}

// ExcludeVisited : Return a string of IDs to be excluded while fetching new nodes
func (vNodes VisitedNodes) ExcludeVisited(tag config.Tag) string {
	visited := ""
	for _, tagMember := range tag.Members {
		if memberIDs, ok := vNodes.Nodes[tagMember]; ok {
			pks := ""
			for pk := range memberIDs {
				if len(pk) > 0 {
					pks += pk + ","
				}
			}
			if pks != "" {
				pks = strings.Trim(pks, ",")
				visited += fmt.Sprintf(" AND %s.id NOT IN (%s) ", tagMember, pks)
			}

		}
	}
	return visited
}

// IsVisited : Check if the passed node has already been visited
func (vNodes VisitedNodes) IsVisited(node *DependencyNode) bool {

	for _, tagMember := range node.Tag.Members {
		if _, ok := vNodes.Nodes[tagMember]; !ok {
			continue
		}
		idCol := fmt.Sprintf("%s.id", tagMember)
		if _, ok := node.Data[idCol]; ok {
			srcID := fmt.Sprint(node.Data[idCol])
			if _, ok := vNodes.Nodes[tagMember][srcID]; ok {
				return true
			}
		} else {
			log.Println("In: IsVisited | node.Data =>", node.Data)
			log.Fatal(idCol, "NOT PRESENT IN NODE DATA")
		}
	}
	return false
}

// MarkAsVisited : Mark the passed node as visited
func (vNodes *VisitedNodes) MarkAsVisited(node *DependencyNode) {
	for _, tagMember := range node.Tag.Members {
		idCol := fmt.Sprintf("%s.id", tagMember)
		if nodeVal, ok := node.Data[idCol]; ok {
			if nodeVal == nil {
				continue
			}
			if _, ok := vNodes.Nodes[tagMember]; !ok {
				vNodes.Nodes[tagMember] = make(map[string]bool)
			}
			srcID := fmt.Sprint(node.Data[idCol])
			vNodes.Nodes[tagMember][srcID] = true
		} else {
			log.Println("In: MarkAsVisited | node.Data =>", node.Data)
			log.Fatal(idCol, "NOT PRESENT IN NODE DATA")
		}
	}
}

package migrate_v2

import (
	"database/sql"
	"fmt"
	"log"
	"stencil/db"
)

func (vBags *VisitedBags) Init(dbConn *sql.DB) {
	vBags.Bags = make(map[string]map[string]map[string]bool)
	vBags.DBConn = dbConn
	vBags.InitPKs()
}

func (vBags *VisitedBags) InitPKs() {
	vBags.BagPKs = make(map[string]bool)
}

func (vBags *VisitedBags) AddPK(id string) {
	vBags.BagPKs[id] = true
}

func (vBags *VisitedBags) PKExists(id string) bool {
	_, ok := vBags.BagPKs[id]
	return ok
}

func (vBags *VisitedBags) UpdatePKs(bagNode *DependencyNode) {
	for _, pk := range bagNode.PKs {
		vBags.AddPK(fmt.Sprint(pk))
	}
}

func (vBags *VisitedBags) IsVisited(bag *DBBag) bool {
	if _, ok := vBags.Bags[bag.AppID]; !ok {
		return false
	}
	if _, ok := vBags.Bags[bag.AppID][bag.MemberID]; !ok {
		return false
	}
	if _, ok := vBags.Bags[bag.AppID][bag.MemberID][bag.ID]; !ok {
		return false
	}
	return true
}

func (vBags *VisitedBags) MarkAsVisited(bag *DBBag) {

	if _, ok := vBags.Bags[bag.AppID]; !ok {
		vBags.Bags[bag.AppID] = make(map[string]map[string]bool)
	}
	if _, ok := vBags.Bags[bag.AppID][bag.MemberID]; !ok {
		vBags.Bags[bag.AppID][bag.MemberID] = make(map[string]bool)
	}
	vBags.Bags[bag.AppID][bag.MemberID][bag.ID] = true
}

func (vBags *VisitedBags) AddNode(node *DependencyNode, appID string) {
	if _, ok := vBags.Bags[appID]; !ok {
		vBags.Bags[appID] = make(map[string]map[string]bool)
	}
	for _, tagMember := range node.Tag.Members {
		idCol := fmt.Sprintf("%s.id", tagMember)
		if memberID, err := db.TableID(vBags.DBConn, tagMember, appID); err == nil {
			if nodeVal, ok := node.Data[idCol]; ok {
				if nodeVal == nil {
					continue
				}
				if _, ok := vBags.Bags[memberID]; !ok {
					vBags.Bags[appID][memberID] = make(map[string]bool)
				}
				srcID := fmt.Sprint(node.Data[idCol])
				vBags.Bags[appID][memberID][srcID] = true
			} else {
				log.Println("@vBags.AddNode | node.Data =>", node.Data)
				log.Fatal(idCol, "NOT PRESENT IN NODE DATA")
			}
		} else {
			fmt.Println(tagMember, appID)
			log.Fatal("@vBags.AddNode: ", err)
		}
	}
}

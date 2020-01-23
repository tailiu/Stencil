package counter

import (
	"errors"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/migrate"
	"strings"
)

func CreateCounter(appName, appID string) Counter {
	AppConfig, err := config.CreateAppConfig(appName, appID)
	if err != nil {
		log.Fatal(err)
	}
	AppConfig.QR.Migration = true
	counter := Counter{
		AppConfig:     AppConfig,
		AppDBConn:     db.GetDBConn("diaspora_count"),
		StencilDBConn: db.GetDBConn(db.STENCIL_DB),
		visitedNodes:  make(map[string]map[string]bool),
		NodeCount:     0,
		EdgeCount:     0}

	return counter
}

func (self *Counter) FetchUserNode(uid string) (*migrate.DependencyNode, error) {

	if root, err := self.AppConfig.GetTag("root"); err == nil {
		rootTable, rootCol := self.AppConfig.GetItemsFromKey(root, "root_id")
		where := fmt.Sprintf("%s.%s = '%s'", rootTable, rootCol, uid)
		ql := self.GetTagQL(root)
		sql := fmt.Sprintf("%s WHERE %s ", ql, where)
		sql += root.ResolveRestrictions()
		if data, err := db.DataCall1(self.AppDBConn, sql); err == nil && len(data) > 0 {
			return &migrate.DependencyNode{Tag: root, SQL: sql, Data: data}, nil
		} else {
			if err == nil {
				err = errors.New("no data returned for root node, doesn't exist?")
			}
			return nil, err
		}
	} else {
		log.Fatal("Can't fetch root tag:", err)
		return nil, err
	}
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	return nil, errors.New("nothing returned?")
}

func (self *Counter) GetAdjNode(node *migrate.DependencyNode) (*migrate.DependencyNode, error) {
	if strings.EqualFold(node.Tag.Name, "root") {
		return self.GetOwnedNode(node)
	}
	return self.GetDependentNode(node)
}

func (self *Counter) GetDependentNode(node *migrate.DependencyNode) (*migrate.DependencyNode, error) {

	for _, dep := range self.AppConfig.ShuffleDependencies(self.AppConfig.GetSubDependencies(node.Tag.Name)) {
		if child, err := self.AppConfig.GetTag(dep.Tag); err == nil {
			where := node.ResolveDependencyConditions(self.AppConfig, dep, child)
			ql := self.GetTagQL(child)
			sql := fmt.Sprintf("%s WHERE %s ", ql, where)
			sql += child.ResolveRestrictions()
			sql += " ORDER BY random()"
			if data, err := db.DataCall1(self.AppDBConn, sql); err == nil {
				if len(data) > 0 {
					log.Println(fmt.Sprintf("FETCHING  tag for dependency { %s > %s } ", node.Tag.Name, dep.Tag))
					return &migrate.DependencyNode{Tag: child, SQL: sql, Data: data}, nil
				}
			} else {
				fmt.Println(err)
				log.Fatal(sql)
				return nil, err
			}
		} else {
			log.Fatal("Unable to fetch tag for: ", dep.Tag)
		}
	}
	return nil, nil
}

func (self *Counter) GetOwnedNode(root *migrate.DependencyNode) (*migrate.DependencyNode, error) {

	for _, own := range self.AppConfig.GetShuffledOwnerships() {
		if child, err := self.AppConfig.GetTag(own.Tag); err == nil {
			where := root.ResolveOwnershipConditions(own, child)
			ql := self.GetTagQL(child)
			sql := fmt.Sprintf("%s WHERE %s ", ql, where)
			sql += child.ResolveRestrictions()
			sql += " ORDER BY random() "
			// log.Fatal(sql)
			if data, err := db.DataCall1(self.AppDBConn, sql); err == nil {
				if len(data) > 0 {
					log.Println(fmt.Sprintf("FETCHING  tag  for ownership { %s } ", own.Tag))
					return &migrate.DependencyNode{Tag: child, SQL: sql, Data: data}, nil
				}
			} else {
				fmt.Println(err)
				log.Fatal(sql)
				return nil, err
			}
		}
	}
	return nil, nil
}

func (self *Counter) GetTagQL(tag config.Tag) string {

	sql := "SELECT %s FROM %s "

	if len(tag.InnerDependencies) > 0 {
		cols := ""
		joinMap := tag.CreateInDepMap()
		seenMap := make(map[string]bool)
		joinStr := ""
		for fromTable, toTablesMap := range joinMap {
			if _, ok := seenMap[fromTable]; !ok {
				joinStr += fromTable
				_, colStr := db.GetColumnsForTable(self.AppDBConn, fromTable)
				cols += colStr + ","
			}
			for toTable, conditions := range toTablesMap {
				if conditions != nil {
					conditions = append(conditions, joinMap[toTable][fromTable]...)
					if joinMap[toTable][fromTable] != nil {
						joinMap[toTable][fromTable] = nil
					}
					joinStr += fmt.Sprintf(" JOIN %s ON %s ", toTable, strings.Join(conditions, " AND "))
					_, colStr := db.GetColumnsForTable(self.AppDBConn, toTable)
					cols += colStr + ","
					seenMap[toTable] = true
				}
			}
			seenMap[fromTable] = true
		}
		sql = fmt.Sprintf(sql, strings.Trim(cols, ","), joinStr)
	} else {
		table := tag.Members["member1"]
		_, cols := db.GetColumnsForTable(self.AppConfig.DBConn, table)
		sql = fmt.Sprintf(sql, cols, table)
	}
	return sql
}

func (self *Counter) RunCounter() error {

	offset := 0

	for {
		if person_id, err := db.GetNextUserFromAppDB("diaspora", "people", "id", offset); err == nil {
			if len(person_id) < 1 {
				break
			}
			log.Println("Current User:", person_id)
			if personNode, err := self.FetchUserNode(person_id); err == nil {
				if err := self.Traverse(personNode); err == nil {
					offset += 1
				} else {
					fmt.Println("User - Offset: ", person_id, offset)
					log.Fatal("Error while traversing: ", err)
				}
			} else {
				log.Fatal("User Node Not Created: ", err)
			}
		} else {
			fmt.Println("User offset: ", offset)
			log.Fatal("Crashed while running counter: ", err)
		}
	}
	fmt.Println("Counter Finished!")
	fmt.Println("Offset: ", offset)
	fmt.Println("Nodes: ", self.NodeCount)
	fmt.Println("Edges: ", self.EdgeCount)
	return nil
}

func (self *Counter) GetAllPreviousNodes(node *migrate.DependencyNode) ([]*migrate.DependencyNode, error) {
	var nodes []*migrate.DependencyNode
	for _, dep := range self.AppConfig.GetParentDependencies(node.Tag.Name) {
		for _, pdep := range dep.DependsOn {
			if parent, err := self.AppConfig.GetTag(pdep.Tag); err == nil {
				where := node.ResolveParentDependencyConditions(pdep.Conditions, parent)
				ql := self.GetTagQL(parent)
				sql := fmt.Sprintf("%s WHERE %s ", ql, where)
				sql += parent.ResolveRestrictions()
				if data, err := db.DataCall(self.AppDBConn, sql); err == nil {
					for _, datum := range data {
						newNode := new(migrate.DependencyNode)
						newNode.Tag = parent
						newNode.SQL = sql
						newNode.Data = datum
						nodes = append(nodes, newNode)
					}
				} else {
					fmt.Println(sql)
					log.Fatal("@GetAllPreviousNodes: Error while DataCall: ", err)
					return nil, err
				}
			} else {
				log.Fatal("@GetAllPreviousNodes: Tag doesn't exist? ", pdep.Tag)
			}
		}
	}
	return nodes, nil
}

func (self *Counter) DeleteNode(node *migrate.DependencyNode) error {
	tx, err := self.AppDBConn.Begin()
	if err != nil {
		log.Fatal("Error creating Source DB Transaction! ", err)
		return err
	}
	defer tx.Rollback()
	for _, tagMember := range node.Tag.Members {
		idCol := fmt.Sprintf("%s.id", tagMember)
		if _, ok := node.Data[idCol]; ok {
			srcID := fmt.Sprint(node.Data[idCol])
			if derr := db.ReallyDeleteRowFromAppDB(tx, tagMember, srcID); derr != nil {
				fmt.Println("@ERROR_DeleteRowFromAppDB", derr)
				fmt.Println("@QARGS:", tagMember, srcID)
				return derr
			}
		} else {
			log.Println("node.Data =>", node.Data)
			log.Fatal(idCol, "NOT PRESENT IN NODE DATA")
		}
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(node.Data)
		log.Fatal("Unable to commit deletion for node ", node.Tag.Name)
	}
	return nil
}

func (self *Counter) Traverse(node *migrate.DependencyNode) error {
	nodeIDAttr, _ := node.Tag.ResolveTagAttr("id")
	for {
		if adjNode, err := self.GetAdjNode(node); err != nil {
			return err
		} else if adjNode == nil {
			break
		} else {
			adjNodeIDAttr, _ := adjNode.Tag.ResolveTagAttr("id")
			log.Println(fmt.Sprintf("Current   Node: { %s } | ID: %v ", node.Tag.Name, node.Data[nodeIDAttr]))
			log.Println(fmt.Sprintf("Adjacent  Node: { %s } | ID: %v ", adjNode.Tag.Name, adjNode.Data[adjNodeIDAttr]))
			if err := self.Traverse(adjNode); err != nil {
				log.Fatal(fmt.Sprintf("ERROR! NODE : { %s } | ID: %v, ADJ_NODE : { %s } | ID: %v | err: [ %s ]", node.Tag.Name, node.Data[nodeIDAttr], adjNode.Tag.Name, adjNode.Data[adjNodeIDAttr], err))
				return err
			}
		}
	}

	self.NodeCount += 1
	if previousNodes, err := self.GetAllPreviousNodes(node); err == nil {
		for _, previousNode := range previousNodes {
			self.EdgeCount += 1
		}
	} else {
		log.Fatal("Error while getting previous nodes for the leaf!")
	}
	if err := self.DeleteNode(node); err != nil {
		fmt.Println(node.Data)
		log.Fatal("Error while deleting node ", node.Tag.Name, node.Data[nodeIDAttr])
	}
	fmt.Println("CURRENT NODES:", self.NodeCount)
	fmt.Println("CURRENT EDGES:", self.EdgeCount)
	return nil
}

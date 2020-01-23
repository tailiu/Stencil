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
		AppDBConn:     db.GetDBConn(appName),
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
		sql += self.ResolveRestrictions(root)
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

func (self *Counter) GetDependentNode(node *migrate.DependencyNode) (*migrate.DependencyNode, error) {

	for _, dep := range self.AppConfig.ShuffleDependencies(self.AppConfig.GetSubDependencies(node.Tag.Name)) {
		if child, err := self.AppConfig.GetTag(dep.Tag); err == nil {
			log.Println(fmt.Sprintf("FETCHING  tag for dependency { %s > %s } ", node.Tag.Name, dep.Tag))
			where := self.ResolveDependencyConditions(node, dep, child)
			ql := self.GetTagQL(child)
			sql := fmt.Sprintf("%s WHERE %s ", ql, where)
			sql += self.ResolveRestrictions(child)
			sql += " ORDER BY random()"
			// log.Fatal(sql)
			if data, err := db.DataCall1(self.AppDBConn, sql); err == nil {
				if len(data) > 0 {
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

func (self *Counter) ResolveRestrictions(tag config.Tag) string {
	restrictions := ""
	for _, restriction := range tag.Restrictions {
		if restrictionAttr, err := tag.ResolveTagAttr(restriction["col"]); err == nil {
			restrictions += fmt.Sprintf(" AND %s = '%s' ", restrictionAttr, restriction["val"])
		}

	}
	return restrictions
}

func (self *Counter) ResolveDependencyConditions(node *migrate.DependencyNode, dep config.Dependency, tag config.Tag) string {

	where := ""
	for _, depOn := range dep.DependsOn {
		if depOnTag, err := self.AppConfig.GetTag(depOn.Tag); err == nil {
			if strings.EqualFold(depOnTag.Name, node.Tag.Name) {
				for _, condition := range depOn.Conditions {
					conditionStr := ""
					tagAttr, err := tag.ResolveTagAttr(condition.TagAttr)
					if err != nil {
						log.Println(err, tag.Name, condition.TagAttr)
						break
					}
					depOnAttr, err := depOnTag.ResolveTagAttr(condition.DependsOnAttr)
					if err != nil {
						log.Println(err, depOnTag.Name, condition.DependsOnAttr)
						break
					}
					if _, ok := node.Data[depOnAttr]; ok {
						if conditionStr != "" || where != "" {
							conditionStr += " AND "
						}
						conditionStr += fmt.Sprintf("%s = '%v'", tagAttr, node.Data[depOnAttr])
					} else {
						fmt.Println(depOnTag)
						log.Fatal("ResolveDependencyConditions:", depOnAttr, " doesn't exist in ", depOnTag.Name)
					}
					if len(condition.Restrictions) > 0 {
						restrictions := ""
						for _, restriction := range condition.Restrictions {
							if restrictionAttr, err := tag.ResolveTagAttr(restriction["col"]); err == nil {
								if restrictions != "" {
									restrictions += " OR "
								}
								restrictions += fmt.Sprintf(" %s = '%s' ", restrictionAttr, restriction["val"])
							}

						}
						if restrictions == "" {
							log.Fatal(condition.Restrictions)
						}
						conditionStr += fmt.Sprintf(" AND (%s) ", restrictions)
					}
					where += conditionStr
				}
			}
		}
	}
	return where
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

func (self *Counter) Traverse(node *migrate.DependencyNode) error {
	for {
		if adjNode, err := self.GetDependentNode(node); err != nil {
			return err
		} else {
			if adjNode == nil {
				break
			}
			nodeIDAttr, _ := node.Tag.ResolveTagAttr("id")
			adjNodeIDAttr, _ := adjNode.Tag.ResolveTagAttr("id")
			log.Println(fmt.Sprintf("Current   Node: { %s } | ID: %v ", node.Tag.Name, node.Data[nodeIDAttr]))
			log.Println(fmt.Sprintf("Adjacent  Node: { %s } | ID: %v ", adjNode.Tag.Name, adjNode.Data[adjNodeIDAttr]))
			if err := self.Traverse(adjNode); err != nil {
				log.Fatal(fmt.Sprintf("ERROR! NODE : { %s } | ID: %v, ADJ_NODE : { %s } | ID: %v | err: [ %s ]", node.Tag.Name, node.Data[nodeIDAttr], adjNode.Tag.Name, adjNode.Data[adjNodeIDAttr], err))
				return err
			}
		}
	}
	return nil
}

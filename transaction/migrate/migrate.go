package migrate

import (
	"fmt"
	"log"
	"strings"
	"transaction/config"
	"transaction/db"
	"transaction/helper"
)

func addUserToApplication(uid string, dstApp config.AppConfig) {

}

func removeUserFromApplication(uid string, srcApp config.AppConfig) {

}

func checkUserInApp(uid string, dstApp config.AppConfig) bool {
	return true
}

func UpdateMigrationState(uid string, srcApp, dstApp config.AppConfig) {

}

func GetRoot(appConfig config.AppConfig, uid string) *DependencyNode {
	tagName := "root"
	if root, err := appConfig.GetTag(tagName); err == nil {
		sql := "SELECT %s FROM %s WHERE %s "
		rootTable, rootCol := appConfig.GetItemsFromKey(root, "root_id")
		where := fmt.Sprintf("%s.%s = $1", rootTable, rootCol)
		if len(root.InnerDependencies) > 0 {
			cols := ""
			joinMap := appConfig.CreateInDepMap(root)
			seenMap := make(map[string]bool)
			joinStr := ""
			for fromTable, toTablesMap := range joinMap {
				if _, ok := seenMap[fromTable]; !ok {
					joinStr += fromTable
					_, colStr := db.GetColumnsForTable(appConfig.AppName, fromTable)
					cols += colStr + ","
				}
				for toTable, conditions := range toTablesMap {
					if conditions != nil {
						conditions = append(conditions, joinMap[toTable][fromTable]...)
						if joinMap[toTable][fromTable] != nil {
							joinMap[toTable][fromTable] = nil
						}
						joinStr += " JOIN " + toTable + " ON " + strings.Join(conditions, " AND ")
						_, colStr := db.GetColumnsForTable(appConfig.AppName, toTable)
						cols += colStr + ","
						seenMap[toTable] = true
					}
				}
				seenMap[fromTable] = true
			}
			sql = fmt.Sprintf(sql, strings.Trim(cols, ","), joinStr, where)
		} else {
			table := root.Members["member1"]
			_, cols := db.GetColumnsForTable(appConfig.AppName, table)
			sql = fmt.Sprintf(sql, cols, table, where)
		}
		rootNode := new(DependencyNode)
		rootNode.Tag = "root"
		rootNode.SQL = sql
		// fmt.Println(sql)
		rootNode.Data = db.DataCall(appConfig.AppName, sql, uid)
		return rootNode
	}
	return nil
}

// Handle restrictions tag in depends on conditions
func ResolveDependencyConditions(node *DependencyNode, appConfig config.AppConfig, dep config.Dependency) string {
	where := ""
	if tag, err := appConfig.GetTag(dep.Tag); err == nil {
		for _, depOn := range dep.DependsOn {
			if depOnTag, err := appConfig.GetTag(depOn.Tag); err == nil {
				if strings.EqualFold(depOnTag.Name, node.Tag) {
					// fmt.Println(tag.Name, depOnTag.Name)
					for _, condition := range depOn.Conditions {
						conditionStr := ""
						tagAttr, err := appConfig.ResolveTagAttr(tag.Name, condition.TagAttr)
						if err != nil {
							log.Println(err, tag.Name, condition.TagAttr)
							break
						}
						depOnAttr, err := appConfig.ResolveTagAttr(depOnTag.Name, condition.DependsOnAttr)
						if err != nil {
							log.Println(err, depOnTag.Name, condition.DependsOnAttr)
							break
						}
						// fmt.Print(tagAttr, "==", depOnAttr, " | ")
						for _, datum := range node.Data {
							if _, ok := datum[depOnAttr]; ok {
								// fmt.Println(depOnAttr, datum[depOnAttr])
								if conditionStr != "" || where != "" {
									conditionStr += " AND "
								}
								conditionStr += fmt.Sprintf("%s = '%v'", tagAttr, datum[depOnAttr])
							} else {
								fmt.Println(depOnAttr, "doesn't exist in ", depOnTag.Name)
							}
						}
						if len(condition.Restrictions) > 0 {
							restrictions := ""
							for _, restriction := range condition.Restrictions {
								if restrictions != "" {
									restrictions += " OR "
								}
								if restrictionAttr, err := appConfig.ResolveTagAttr(tag.Name, restriction["col"]); err == nil {
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
	}
	return where
}

func GetAdjNode(node *DependencyNode, appConfig config.AppConfig, uid string) *DependencyNode {

	for _, dep := range helper.ShuffleDependencies(appConfig.GetSubDependencies(node.Tag)) {
		if where := ResolveDependencyConditions(node, appConfig, dep); where != "" {
			limit := " LIMIT 1 "
			orderby := " ORDER BY random() "
			if child, err := appConfig.GetTag(dep.Tag); err == nil {
				sql := "SELECT %s FROM %s WHERE %s %s %s"
				if len(child.Restrictions) > 0 {
					restrictions := ""
					for _, restriction := range child.Restrictions {
						if restrictions != "" {
							restrictions += " OR "
						}
						if restrictionAttr, err := appConfig.ResolveTagAttr(child.Name, restriction["col"]); err == nil {
							restrictions += fmt.Sprintf(" %s = '%s' ", restrictionAttr, restriction["val"])
						}

					}
					where += fmt.Sprintf(" AND (%s) ", restrictions)
				}
				if len(child.InnerDependencies) > 0 {
					cols := ""
					joinMap := appConfig.CreateInDepMap(child)
					seenMap := make(map[string]bool)
					joinStr := ""
					for fromTable, toTablesMap := range joinMap {
						if _, ok := seenMap[fromTable]; !ok {
							joinStr += fromTable
							_, colStr := db.GetColumnsForTable(appConfig.AppName, fromTable)
							cols += colStr + ","
						}
						for toTable, conditions := range toTablesMap {
							if conditions != nil {
								conditions = append(conditions, joinMap[toTable][fromTable]...)
								if joinMap[toTable][fromTable] != nil {
									joinMap[toTable][fromTable] = nil
								}
								// joinStr += " JOIN " + toTable + " ON " + strings.Join(conditions, " AND ")
								joinStr += fmt.Sprintf(" JOIN %s ON %s ", toTable, strings.Join(conditions, " AND "))
								_, colStr := db.GetColumnsForTable(appConfig.AppName, toTable)
								cols += colStr + ","
								seenMap[toTable] = true
							}
						}
						seenMap[fromTable] = true
					}
					sql = fmt.Sprintf(sql, strings.Trim(cols, ","), joinStr, where, orderby, limit)
				} else {
					table := child.Members["member1"]
					_, cols := db.GetColumnsForTable(appConfig.AppName, table)
					sql = fmt.Sprintf(sql, cols, table, where, orderby, limit)
				}
				if nodeData := db.DataCall(appConfig.AppName, sql); len(nodeData) > 0 {
					newNode := new(DependencyNode)
					newNode.Tag = dep.Tag
					newNode.SQL = sql
					newNode.Data = nodeData
					// fmt.Println(sql)
					return newNode
				}
			}
		}
	}
	return nil
}

func MigrateNode(node *DependencyNode, srcApp, dstApp config.AppConfig) {
	if mappings := config.GetSchemaMappingsFor(srcApp.AppName, dstApp.AppName); mappings == nil {
		log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcApp.AppName, dstApp.AppName))
	} else {
		if nodeTag, err := srcApp.GetTag(node.Tag); err == nil {
			for _, appMapping := range mappings.Mappings {
				if tagMembers, err := srcApp.GetTagMembers(nodeTag.Name); err == nil {
					if mappedTables := helper.IntersectString(tagMembers, appMapping.FromTables); len(mappedTables) > 0 {
						if len(tagMembers) == len(appMapping.FromTables) {
							fmt.Println("Fully Mapped:", mappedTables)
							for _, toTable := range appMapping.ToTables {
								if len(toTable.Conditions) > 0 {
									breakCondition := false
									fmt.Println("toTable.Conditions", toTable.Conditions)
									for conditionKey, conditionVal := range toTable.Conditions {
										if nodeVal := node.GetValueForKey(conditionKey); nodeVal != nil {
											if !strings.EqualFold(*nodeVal, conditionVal) {
												breakCondition = true
												fmt.Println(*nodeVal, "!=", conditionVal)
											} else {
												fmt.Println(*nodeVal, "==", conditionVal)
											}
										}
									}
									if breakCondition {
										continue // Move on to the next mapping.
									}
								}
								cols, vals := "", ""
								for toAttr, fromAttr := range toTable.Mapping {
									if val := node.GetValueForKey(fromAttr); val != nil {
										vals += fmt.Sprintf("'%s',", *val)
										cols += fmt.Sprintf("%s,", toAttr)
									}
								}
								if cols != "" && vals != "" {
									cols := strings.Trim(cols, ",")
									vals := strings.Trim(vals, ",")
									isql := fmt.Sprintf("INSERT INTO %s (%s) VALUEs (%s);", toTable.Table, cols, vals)
									fmt.Println("** Insert Query:", isql)
								} else {
									fmt.Println("## Insert Query Error:", cols, vals)
								}

							}
						} else {
							fmt.Println("Partially Mapped:", mappedTables)
						}
					}
				}
			}
		}
		log.Fatal(fmt.Sprintf("Mappings found from [%s] to [%s].", srcApp.AppName, dstApp.AppName))
	}
}

func MigrateProcess(uid string, srcApp, dstApp config.AppConfig, node *DependencyNode) {

	// try:

	if node.Tag == "root" && !checkUserInApp(uid, dstApp) {
		addUserToApplication(uid, dstApp)
	}

	for child := GetAdjNode(node, srcApp, uid); child != nil; child = GetAdjNode(node, srcApp, uid) {
		fmt.Println("------------------------------------------------------------------------")
		log.Println("Current Node:", node.Tag)
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		log.Println("Child Node:", child.Tag)
		fmt.Println("------------------------------------------------------------------------")
		MigrateProcess(uid, srcApp, dstApp, child)
	}
	// acquirePredicateLock(*node)
	// for child := GetAdjNode(node, srcApp, uid); child != nil; child = GetAdjNode(node, srcApp, uid) {
	// 	MigrateProcess(uid, srcApp, dstApp, child)
	// }
	MigrateNode(node, srcApp, dstApp) // Log before migrating
	// releasePredicateLock(*node)

	// catch NodeNotFound:

	// t.releaseAllLocks()
	// if node.Tag == "root" {
	// 	MigrateProcess(uid, srcApp, dstApp, GetRoot(srcApp, uid))
	// } else {
	// 	if checkUserInApp(uid, srcApp) {
	// 		removeUserFromApplication(uid, srcApp)
	// 	}
	// 	UpdateMigrationState(uid, srcApp, dstApp)
	// 	log.Println("Congratulations, this migration worker has finished it's job!")
	// }
}

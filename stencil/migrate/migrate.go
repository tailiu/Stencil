package migrate

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/display"
	"stencil/helper"
	"stencil/transaction"
	"strings"
)

var USEREXISTSINAPP = false

func ResetUserExistsInApp() {
	USEREXISTSINAPP = false
}

func remove(s []config.Tag, i int) []config.Tag {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func addUserToApplication(node *DependencyNode, srcApp, dstApp config.AppConfig, log_txn *transaction.Log_txn) bool {
	if mappings := config.GetSchemaMappingsFor(srcApp.AppName, dstApp.AppName); mappings == nil {
		log.Fatal(fmt.Sprintf("addUserToApplication: Can't find mappings from [%s] to [%s].", srcApp.AppName, dstApp.AppName))
	} else {
		tagMembers := node.Tag.GetTagMembers()
		for _, appMapping := range mappings.Mappings {
			// GenerateAndInsert(mappings, appMapping.ToTables, node)
			if mappedTables := helper.IntersectString(tagMembers, appMapping.FromTables); len(mappedTables) > 0 {
				if len(tagMembers) == len(appMapping.FromTables) {
					errs := GenerateAndInsert(mappings, dstApp, appMapping.ToTables, node, log_txn)
					fmt.Println(errs)
					USEREXISTSINAPP = true
					return true
				}
			}
		}
	}
	return false
}

func removeUserFromApplication(uid string, srcApp config.AppConfig) {

}

func checkUserInApp(uid string, dstApp config.AppConfig) bool {
	return USEREXISTSINAPP
}

func UpdateMigrationState(uid string, srcApp, dstApp config.AppConfig) {

}

func GetRoot(appConfig config.AppConfig, uid string) (*DependencyNode, error) {
	tagName := "root"
	if root, err := appConfig.GetTag(tagName); err == nil {
		sql := "SELECT %s FROM %s WHERE %s "
		rootTable, rootCol := appConfig.GetItemsFromKey(root, "root_id")
		where := fmt.Sprintf("%s.%s = $1", rootTable, rootCol)
		if len(root.InnerDependencies) > 0 {
			cols := ""
			joinMap := root.CreateInDepMap()
			seenMap := make(map[string]bool)
			joinStr := ""
			for fromTable, toTablesMap := range joinMap {
				if _, ok := seenMap[fromTable]; !ok {
					joinStr += fromTable
					_, colStr := db.GetColumnsForTable(appConfig.DBConn, fromTable)
					cols += colStr + ","
				}
				for toTable, conditions := range toTablesMap {
					if conditions != nil {
						conditions = append(conditions, joinMap[toTable][fromTable]...)
						if joinMap[toTable][fromTable] != nil {
							joinMap[toTable][fromTable] = nil
						}
						joinStr += " JOIN " + toTable + " ON " + strings.Join(conditions, " AND ")
						_, colStr := db.GetColumnsForTable(appConfig.DBConn, toTable)
						cols += colStr + ","
						seenMap[toTable] = true
					}
				}
				seenMap[fromTable] = true
			}
			sql = fmt.Sprintf(sql, strings.Trim(cols, ","), joinStr, where)
		} else {
			table := root.Members["member1"]
			_, cols := db.GetColumnsForTable(appConfig.DBConn, table)
			sql = fmt.Sprintf(sql, cols, table, where)
		}
		rootNode := new(DependencyNode)
		rootNode.Tag = root
		rootNode.SQL = sql
		// fmt.Println(sql)
		rootNode.Data, err = db.DataCall1(appConfig.DBConn, sql, uid)
		return rootNode, err
	} else {
		return nil, errors.New("Can't fetch tag:" + tagName)
	}
}

func ResolveDependencyConditions(node *DependencyNode, appConfig config.AppConfig, dep config.Dependency) string {
	where := ""
	if tag, err := appConfig.GetTag(dep.Tag); err == nil {
		for _, depOn := range dep.DependsOn {
			if depOnTag, err := appConfig.GetTag(depOn.Tag); err == nil {
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
								if restrictions != "" {
									restrictions += " OR "
								}
								if restrictionAttr, err := tag.ResolveTagAttr(restriction["col"]); err == nil {
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

func GetAdjNode(node *DependencyNode, appConfig config.AppConfig, uid string, wList *WaitingList, invalidList *InvalidList) (*DependencyNode, error) {

	for _, dep := range config.ShuffleDependencies(appConfig.GetSubDependencies(node.Tag.Name)) {
		if where := ResolveDependencyConditions(node, appConfig, dep); where != "" {
			orderby := " ORDER BY random() "
			if child, err := appConfig.GetTag(dep.Tag); err == nil {
				sql := "SELECT %s FROM %s WHERE %s %s "
				if len(child.Restrictions) > 0 {
					restrictions := ""
					for _, restriction := range child.Restrictions {
						if restrictions != "" {
							restrictions += " OR "
						}
						if restrictionAttr, err := child.ResolveTagAttr(restriction["col"]); err == nil {
							restrictions += fmt.Sprintf(" %s = '%s' ", restrictionAttr, restriction["val"])
						}

					}
					where += fmt.Sprintf(" AND (%s) ", restrictions)
				}
				if len(child.InnerDependencies) > 0 {
					cols := ""
					joinMap := child.CreateInDepMap()
					seenMap := make(map[string]bool)
					joinStr := ""
					for fromTable, toTablesMap := range joinMap {
						if _, ok := seenMap[fromTable]; !ok {
							joinStr += fromTable
							_, colStr := db.GetColumnsForTable(appConfig.DBConn, fromTable)
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
								_, colStr := db.GetColumnsForTable(appConfig.DBConn, toTable)
								cols += colStr + ","
								seenMap[toTable] = true
							}
						}
						seenMap[fromTable] = true
					}
					sql = fmt.Sprintf(sql, strings.Trim(cols, ","), joinStr, where, orderby)
				} else {
					table := child.Members["member1"]
					_, cols := db.GetColumnsForTable(appConfig.DBConn, table)
					sql = fmt.Sprintf(sql, cols, table, where, orderby)
				}
				if nodeData, err := db.DataCall1(appConfig.DBConn, sql); err == nil {
					if len(nodeData) > 0 {
						newNode := new(DependencyNode)
						newNode.Tag = child
						newNode.SQL = sql
						newNode.Data = nodeData
						if !wList.IsAlreadyWaiting(*newNode) && !invalidList.Exists(*newNode) {
							return newNode, nil
						}
					}
				} else {
					return nil, err
				}
			}
		}
	}
	return nil, nil
}

func GenerateAndInsert(mappings *config.MappedApp, dstApp config.AppConfig, toTables []config.ToTable, node *DependencyNode, log_txn *transaction.Log_txn) []error {
	// var isqls []string
	var errs []error
	for _, toTable := range toTables {
		if len(toTable.Conditions) > 0 {
			breakCondition := false
			for conditionKey, conditionVal := range toTable.Conditions {
				if nodeVal, err := node.GetValueForKey(conditionKey); err == nil {
					if !strings.EqualFold(nodeVal, conditionVal) {
						breakCondition = true
						log.Println(nodeVal, "!=", conditionVal)
					} else {
						// fmt.Println(*nodeVal, "==", conditionVal)
					}
				} else {
					breakCondition = true
					log.Println("Condition Key", conditionKey, "doesn't exist!")
				}
			}
			if breakCondition {
				continue // Move on to the next mapping.
			}
		}
		undoAction := new(transaction.UndoAction)
		cols, vals := "", ""
		orgCols, orgColsLeft := "", ""
		var ivals []interface{}
		for toAttr, fromAttr := range toTable.Mapping {
			// if val, err := node.GetValueForKey(fromAttr); err == nil {
			if val, ok := node.Data[fromAttr]; ok {
				// vals += fmt.Sprintf("'%v',", val)
				ivals = append(ivals, val)
				vals += fmt.Sprintf("$%d,", len(ivals))
				cols += fmt.Sprintf("%s,", toAttr)
				orgCols += fmt.Sprintf("%s,", strings.Split(fromAttr, ".")[1])
				undoAction.AddData(fromAttr, val)
				undoAction.AddOrgTable(strings.Split(fromAttr, ".")[0])
			} else if strings.Contains(fromAttr, "$") {
				if inputVal, err := mappings.GetInput(fromAttr); err == nil {
					// vals += fmt.Sprintf("'%s',", inputVal)
					ivals = append(ivals, inputVal)
					vals += fmt.Sprintf("$%d,", len(ivals))
					cols += fmt.Sprintf("%s,", toAttr)
					orgCols += fmt.Sprintf("%s,", fromAttr)
				}
			} else if strings.Contains(fromAttr, "#") {
				// Resolve Mapping Method
			} else {
				orgColsLeft += fmt.Sprintf("%s,", strings.Split(fromAttr, ".")[1])
			}
		}
		if cols != "" && vals != "" {
			orgCols := strings.Trim(orgCols, ",")
			cols := strings.Trim(cols, ",")
			vals := strings.Trim(vals, ",")
			isql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ", toTable.Table, cols, vals)
			undoAction.AddDstTable(toTable.Table)
			undoActionSerialized, _ := json.Marshal(undoAction)
			if id, err := db.Insert(dstApp.DBConn, isql, ivals...); err == nil {
				transaction.LogChange(string(undoActionSerialized), log_txn)
				displayFlag := false
				if strings.EqualFold(node.Tag.Name, "root") {
					displayFlag = true
				}
				if err := display.GenDisplayFlag(log_txn.DBconn, dstApp.AppName, toTable.Table, id, displayFlag, log_txn.Txn_id); err != nil {
					log.Println("## DISPLAY ERROR!", err)
					errs = append(errs, err)
				}
				for _, fromTable := range undoAction.OrgTables {
					srcID := "0"
					if _, ok := node.Data[fmt.Sprintf("%s.id", fromTable)]; ok {
						srcID = fmt.Sprint(node.Data[fmt.Sprintf("%s.id", fromTable)])
					}
					if serr := db.SaveForEvaluation(log_txn.DBconn, "diaspora", dstApp.AppName, fromTable, toTable.Table, srcID, fmt.Sprint(id), orgCols, cols, fmt.Sprint(log_txn.Txn_id)); serr != nil {
						log.Fatal(serr)
					}
				}

			} else {
				fmt.Println("\n@ERROR")
				fmt.Println("@SQL:", isql)
				fmt.Println("@ARGS:", ivals)
				fmt.Println(err)
				db.LogError(log_txn.DBconn, isql, fmt.Sprint(ivals), fmt.Sprint(log_txn.Txn_id), dstApp.AppName, err.Error())
				errs = append(errs, err)
			}
		} else {
			errs = append(errs, errors.New("Insert Query Error"))
			log.Fatal("## Insert Query Error:", cols, vals)
		}
	}
	return errs
}

func GetMappedDataForTable(mappings *config.MappedApp, toTable config.ToTable, node *DependencyNode) (string, string, string, []interface{}) {
	cols, vals, orgCols, orgColsLeft := "", "", "", ""
	var ivals []interface{}
	for toAttr, fromAttr := range toTable.Mapping {
		// if val, err := node.GetValueForKey(fromAttr); err == nil {
		if val, ok := node.Data[fromAttr]; ok {
			// vals += fmt.Sprintf("'%v',", val)
			ivals = append(ivals, val)
			vals += fmt.Sprintf("$%d,", len(ivals))
			cols += fmt.Sprintf("%s,", toAttr)
			orgCols += fmt.Sprintf("%s,", strings.Split(fromAttr, ".")[1])
		} else if strings.Contains(fromAttr, "$") {
			if inputVal, err := mappings.GetInput(fromAttr); err == nil {
				// vals += fmt.Sprintf("'%s',", inputVal)
				ivals = append(ivals, inputVal)
				vals += fmt.Sprintf("$%d,", len(ivals))
				cols += fmt.Sprintf("%s,", toAttr)
				orgCols += fmt.Sprintf("%s,", fromAttr)
			}
		} else if strings.Contains(fromAttr, "#") {
			// Resolve Mapping Method
		} else {
			orgColsLeft += fmt.Sprintf("%s,", strings.Split(fromAttr, ".")[1])
		}
	}
	orgCols = strings.Trim(orgCols, ",")
	cols = strings.Trim(cols, ",")
	vals = strings.Trim(vals, ",")
	return cols, vals, orgCols, ivals
}

func PostProcessInsert(id int, dstApp config.AppConfig, toTable config.ToTable, log_txn *transaction.Log_txn, cols, orgCols string, node *DependencyNode, dbConn *sql.DB) {
	undoAction := new(transaction.UndoAction)
	undoAction.AddDstTable(toTable.Table)
	undoActionSerialized, _ := json.Marshal(undoAction)
	transaction.LogChange(string(undoActionSerialized), log_txn)
	displayFlag := false
	if strings.EqualFold(node.Tag.Name, "root") {
		displayFlag = true
	}
	if err := display.GenDisplayFlag(dbConn, dstApp.AppName, toTable.Table, id, displayFlag, log_txn.Txn_id); err != nil {
		log.Println("## DISPLAY ERROR!", err)
	}
	for _, fromTable := range undoAction.OrgTables {
		srcID := "0"
		if _, ok := node.Data[fmt.Sprintf("%s.id", fromTable)]; ok {
			srcID = fmt.Sprint(node.Data[fmt.Sprintf("%s.id", fromTable)])
		}
		if serr := db.SaveForEvaluation(dbConn, "diaspora", dstApp.AppName, fromTable, toTable.Table, srcID, fmt.Sprint(id), orgCols, cols, fmt.Sprint(log_txn.Txn_id)); serr != nil {
			log.Fatal(serr)
		}
	}
}

func MigrateNode(node *DependencyNode, srcApp, dstApp config.AppConfig, wList *WaitingList, invalidList *InvalidList, log_txn *transaction.Log_txn) {
	if mappings := config.GetSchemaMappingsFor(srcApp.AppName, dstApp.AppName); mappings == nil {
		log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcApp.AppName, dstApp.AppName))
	} else {
		mappingFound := false
		for _, appMapping := range mappings.Mappings {
			tagMembers := node.Tag.GetTagMembers()
			if mappedTables := helper.IntersectString(tagMembers, appMapping.FromTables); len(mappedTables) > 0 {
				mappingFound = true
				if len(tagMembers) == len(appMapping.FromTables) {
					GenerateAndInsert(mappings, dstApp, appMapping.ToTables, node, log_txn)
					invalidList.Add(node)
				} else {
					if waitingNode, err := wList.UpdateIfBeingLookedFor(node); err == nil {
						if waitingNode.IsComplete() {
							tempCombinedDataDependencyNode := waitingNode.GenDependencyDataNode()
							GenerateAndInsert(mappings, dstApp, appMapping.ToTables, &tempCombinedDataDependencyNode, log_txn)
							invalidList.Add(node)
						} else {
							// fmt.Println("-->> IS NOT COMPLETE!")
						}
					} else {
						adjTags := srcApp.GetTagsByTables(appMapping.FromTables)
						if err := wList.AddNewToWaitingList(node, adjTags, srcApp); err != nil {
							fmt.Println("!! ERROR !!", err)
						}
					}
				}
				break
			}
		}
		if !mappingFound {
			invalidList.Add(node)
			for _, tagMember := range node.Tag.Members {
				if _, ok := node.Data[fmt.Sprintf("%s.id", tagMember)]; ok {
					srcID := fmt.Sprint(node.Data[fmt.Sprintf("%s.id", tagMember)])
					if serr := db.SaveForEvaluation(log_txn.DBconn, "diaspora", dstApp.AppName, tagMember, "n/a", srcID, "n/a", "*", "n/a", fmt.Sprint(log_txn.Txn_id)); serr != nil {
						log.Fatal(serr)
					}
				}
			}
			// fmt.Println(node)
		}
		// log.Fatal(fmt.Sprintf("Mappings found from [%s] to [%s].", srcApp.AppName, dstApp.AppName))
	}
}

func MigrateProcess(uid string, srcApp, dstApp config.AppConfig, node *DependencyNode, wList *WaitingList, invalidList *InvalidList, log_txn *transaction.Log_txn) error {

	// try:

	if strings.EqualFold(node.Tag.Name, "root") && !checkUserInApp(uid, dstApp) {
		log.Println("++ Adding User from ", srcApp.AppName, " to ", dstApp.AppName)
		addUserToApplication(node, srcApp, dstApp, log_txn)
	}

	for child, err := GetAdjNode(node, srcApp, uid, wList, invalidList); child != nil; child, err = GetAdjNode(node, srcApp, uid, wList, invalidList) {
		if err != nil {
			return err
		}
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		nodeIDAttr, _ := node.Tag.ResolveTagAttr("id")
		childIDAttr, _ := child.Tag.ResolveTagAttr("id")
		log.Println("-- Currrent Node:", node.Tag.Name, "ID:", node.Data[nodeIDAttr])
		log.Println("-- Adjacent Node:", child.Tag.Name, "ID:", child.Data[childIDAttr])
		if mErr := MigrateProcess(uid, srcApp, dstApp, child, wList, invalidList, log_txn); mErr != nil {
			return mErr
		}
	}
	// acquirePredicateLock(*node)
	// for child := GetAdjNode(node, srcApp, uid); child != nil; child = GetAdjNode(node, srcApp, uid) {
	// 	MigrateProcess(uid, srcApp, dstApp, child)
	// }

	if !strings.EqualFold(node.Tag.Name, "root") {
		log.Println("++ Began migrating node *", node.Tag.Name, "from", srcApp.AppName, "to", dstApp.AppName)
		MigrateNode(node, srcApp, dstApp, wList, invalidList, log_txn) // Log before migrating
		log.Println("++ Finished migrating node *", node.Tag.Name, "from", srcApp.AppName, "to", dstApp.AppName)
	}
	fmt.Println("------------------------------------------------------------------------")
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
	return nil
}

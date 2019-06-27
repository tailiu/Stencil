package migrate

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/display"
	"stencil/helper"
	m2 "stencil/migrate"
	"stencil/qr"
	"stencil/transaction"
	"strings"
)

func RemoveUserFromApp(uid, app_id string, log_txn *transaction.Log_txn) bool {
	sql := "DELETE FROM user_table WHERE user_id = $1 AND app_id = $2"
	if err := db.Delete(log_txn.DBconn, sql, uid, app_id); err == nil {
		return true
	}
	return false
}

func checkUserInApp(uid, app_id string, log_txn *transaction.Log_txn) bool {
	sql := "SELECT user_id FROM user_table WHERE user_id = $1 AND app_id = $2"
	res := db.DataCall1(log_txn.DBconn, sql, uid, app_id)
	if len(res) > 0 {
		return true
	}
	return false
}

func addUserToApp(uid, app_id string, log_txn *transaction.Log_txn) bool {
	query := "INSERT INTO user_table (user_id, app_id) VALUES ($1, $2)"
	log_txn.DBconn.QueryRow(query, uid, app_id)
	return true
}

func UpdateMigrationState(uid string, srcApp, dstApp config.AppConfig) {

}

func GetRoot(appConfig config.AppConfig, uid string, log_txn *transaction.Log_txn) *m2.DependencyNode {
	tagName := "root"
	if root, err := appConfig.GetTag(tagName); err == nil {
		var sql string
		qs := qr.CreateQS(appConfig.QR)
		rootTable, rootCol := appConfig.GetItemsFromKey(root, "root_id")
		if len(root.InnerDependencies) > 0 {
			joinMap := root.CreateInDepMap()
			seenMap := make(map[string]bool)
			for fromTable, toTablesMap := range joinMap {
				if _, ok := seenMap[fromTable]; !ok {
					qs.FromSimple(fromTable)
					qs.ColSimple(fromTable + ".*")
					qs.ColPK(fromTable)
				}
				for toTable, conditions := range toTablesMap {
					if conditions != nil {
						conditions = append(conditions, joinMap[toTable][fromTable]...)
						if joinMap[toTable][fromTable] != nil {
							joinMap[toTable][fromTable] = nil
						}
						qs.FromJoinList(toTable, conditions)
						qs.ColSimple(toTable + ".*")
						qs.ColPK(toTable)
						seenMap[toTable] = true
					}
				}
				seenMap[fromTable] = true
			}
		} else {
			table := root.Members["member1"]
			qs.FromSimple(table)
			qs.ColSimple(rootTable + ".*")
			qs.ColPK(rootTable)
		}
		qs.WhereSimpleVal(rootTable+"."+rootCol, "=", uid)
		sql = qs.GenSQL()
		rootNode := new(m2.DependencyNode)
		rootNode.Tag = root
		rootNode.SQL = sql
		rootNode.Data = db.DataCall1(log_txn.DBconn, sql)
		return rootNode
	} else {
		log.Fatal("Can't fetch tag:", tagName)
	}
	return nil
}

func ResolveDependencyConditions(node *m2.DependencyNode, appConfig config.AppConfig, dep config.Dependency, tag config.Tag, qs *qr.QS) string {

	where := qr.CreateQS(appConfig.QR)
	where.TableAliases = qs.TableAliases
	for _, depOn := range dep.DependsOn {
		if depOnTag, err := appConfig.GetTag(depOn.Tag); err == nil {
			if strings.EqualFold(depOnTag.Name, node.Tag.Name) {
				for _, condition := range depOn.Conditions {
					conditionStr := qr.CreateQS(appConfig.QR)
					conditionStr.TableAliases = qs.TableAliases
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
						conditionStr.WhereOperatorInterface("AND", tagAttr, "=", node.Data[depOnAttr])
					} else {
						fmt.Println(depOnTag)
						log.Fatal("ResolveDependencyConditions:", depOnAttr, " doesn't exist in ", depOnTag.Name)
					}
					if len(condition.Restrictions) > 0 {
						restrictions := qr.CreateQS(appConfig.QR)
						restrictions.TableAliases = qs.TableAliases
						for _, restriction := range condition.Restrictions {
							if restrictionAttr, err := tag.ResolveTagAttr(restriction["col"]); err == nil {
								restrictions.WhereOperatorInterface("OR", restrictionAttr, "=", restriction["val"])
							}

						}
						if restrictions.Where == "" {
							log.Fatal(condition.Restrictions)
						}
						// log.Fatal("restrictions.Where", restrictions.Where)
						conditionStr.WhereString("AND", restrictions.Where)
					}
					// log.Fatal("conditionStr.Where", conditionStr.Where)
					where.WhereString("AND", conditionStr.Where)
				}
			}
		}
	}
	// log.Fatal("where.Where", where.Where)
	return where.Where
}

func GetAdjNode(node *m2.DependencyNode, appConfig config.AppConfig, uid string, wList *m2.WaitingList, invalidList *m2.InvalidList, log_txn *transaction.Log_txn) *m2.DependencyNode {

	for _, dep := range config.ShuffleDependencies(appConfig.GetSubDependencies(node.Tag.Name)) {
		// for _, dep := range appConfig.GetSubDependencies(node.Tag.Name) {
		// if where := ResolveDependencyConditions(node, appConfig, dep); where != "" {
		if child, err := appConfig.GetTag(dep.Tag); err == nil {
			var sql string
			qs := qr.CreateQS(appConfig.QR)
			if len(child.InnerDependencies) > 0 {
				joinMap := child.CreateInDepMap()
				seenMap := make(map[string]bool)
				for fromTable, toTablesMap := range joinMap {
					if _, ok := seenMap[fromTable]; !ok {
						qs.FromSimple(fromTable)
						qs.ColSimple(fromTable + ".*")
						qs.ColPK(fromTable)
					}
					for toTable, conditions := range toTablesMap {
						if conditions != nil {
							conditions = append(conditions, joinMap[toTable][fromTable]...)
							if joinMap[toTable][fromTable] != nil {
								joinMap[toTable][fromTable] = nil
							}
							qs.FromJoinList(toTable, conditions)
							qs.ColSimple(toTable + ".*")
							qs.ColPK(toTable)
							seenMap[toTable] = true
						}
					}
					seenMap[fromTable] = true
				}
			} else {
				table := child.Members["member1"]
				qs = qr.CreateQS(appConfig.QR)
				qs.FromSimple(table)
				qs.ColPK(table)
				qs.ColSimple(table + ".*")

				// qs.WhereSimpleVal(table+"."+rootCol, "=", uid)
			}
			if len(child.Restrictions) > 0 {
				restrictions := qr.CreateQS(appConfig.QR)
				restrictions.TableAliases = qs.TableAliases
				for _, restriction := range child.Restrictions {
					if restrictionAttr, err := child.ResolveTagAttr(restriction["col"]); err == nil {
						restrictions.WhereOperatorInterface("OR", restrictionAttr, "=", restriction["val"])
					}

				}
				// log.Fatal("restrictions2.Where", restrictions.Where)
				qs.WhereString("AND", restrictions.Where)
			}
			where := ResolveDependencyConditions(node, appConfig, dep, child, qs)
			qs.WhereString("AND", where)
			qs.OrderBy("random()")
			sql = qs.GenSQL()
			// fmt.Println("** PHY SQL:", sql)
			if nodeData := db.DataCall1(log_txn.DBconn, sql); len(nodeData) > 0 {
				newNode := new(m2.DependencyNode)
				newNode.Tag = child
				newNode.SQL = sql
				newNode.Data = nodeData
				if !wList.IsAlreadyWaiting(*newNode) && !invalidList.Exists(*newNode) {
					return newNode
				}
			}
		}
		// }
	}
	return nil
}

func GenerateAndInsert(mappings *config.MappedApp, dstApp config.AppConfig, toTables []config.ToTable, node *m2.DependencyNode, log_txn *transaction.Log_txn) []error {
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

func CheckMappingConditions(toTable config.ToTable, node *m2.DependencyNode) bool {
	breakCondition := false
	if len(toTable.Conditions) > 0 {
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
				fmt.Println("node data:", node.Data)
				fmt.Println("node sql:", node.SQL)
				log.Fatal("stop here and check")
			}
		}
	}
	return breakCondition
}

func GetMappedDataForTable(mappings *config.MappedApp, toTable config.ToTable, node *m2.DependencyNode) (string, string, string, []interface{}) {
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

func PostProcessInsert(id int, dstApp config.AppConfig, toTable config.ToTable, log_txn *transaction.Log_txn, cols, orgCols string, node *m2.DependencyNode) {
	undoAction := new(transaction.UndoAction)
	undoAction.AddDstTable(toTable.Table)
	undoActionSerialized, _ := json.Marshal(undoAction)
	transaction.LogChange(string(undoActionSerialized), log_txn)
	displayFlag := false
	if strings.EqualFold(node.Tag.Name, "root") {
		displayFlag = true
	}
	if err := display.GenDisplayFlag(log_txn.DBconn, dstApp.AppName, toTable.Table, id, displayFlag, log_txn.Txn_id); err != nil {
		log.Println("## DISPLAY ERROR!", err)
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
}

func UpdateRowDesc(mappings *config.MappedApp, dstApp config.AppConfig, toTables []config.ToTable, node *m2.DependencyNode, log_txn *transaction.Log_txn) bool {

	if tx, err := log_txn.DBconn.Begin(); err != nil {
		log.Fatal("Can't create UpdateRowDesc transaction!")
	} else {
		// var errs []error
		var updated []string

		for _, toTable := range toTables {
			if CheckMappingConditions(toTable, node) {
				continue
			}
			for col, val := range node.Data {
				if strings.Contains(col, "pk.") {
					pk := fmt.Sprint(val)
					if !helper.Contains(updated, pk) {
						if err := db.SetAppID(tx, pk, dstApp.AppID); err == nil {
							updated = append(updated, pk)
						} else {
							tx.Rollback()
							fmt.Println("\n@ERROR:", err)
							// db.LogError(log_txn.DBconn, usql, "", fmt.Sprint(log_txn.Txn_id), dstApp.AppName, err.Error())
							// errs = append(errs, err)
							return false
						}
					}
				}
			}
		}
		tx.Rollback()
		// tx.Commit()
	}
	return true
}

func MigrateNode(node *m2.DependencyNode, srcApp, dstApp config.AppConfig, wList *m2.WaitingList, invalidList *m2.InvalidList, log_txn *transaction.Log_txn) bool {
	if mappings := config.GetSchemaMappingsFor(srcApp.AppName, dstApp.AppName); mappings == nil {
		log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcApp.AppName, dstApp.AppName))
	} else {
		mappingFound := false
		for _, appMapping := range mappings.Mappings {
			tagMembers := node.Tag.GetTagMembers()
			if mappedTables := helper.IntersectString(tagMembers, appMapping.FromTables); len(mappedTables) > 0 {
				mappingFound = true
				if len(tagMembers) == len(appMapping.FromTables) {
					invalidList.Add(node)
					return UpdateRowDesc(mappings, dstApp, appMapping.ToTables, node, log_txn)
				} else {
					log.Println("!! Node [", node.Tag.Name, "] needs to wait?")
					if waitingNode, err := wList.UpdateIfBeingLookedFor(*node); err == nil {
						log.Println("!! Node [", node.Tag.Name, "] updated an existing waiting node!")
						if waitingNode.IsComplete() {
							log.Println("!! Node [", node.Tag.Name, "] completed a waiting node!")
							tempCombinedDataDependencyNode := waitingNode.GenDependencyDataNode()
							invalidList.Add(node)
							return UpdateRowDesc(mappings, dstApp, appMapping.ToTables, &tempCombinedDataDependencyNode, log_txn)
						} else {
							log.Println("!! Node [", node.Tag.Name, "] added to an incomplete waiting node!")
							return false
						}
					} else {
						adjTags := srcApp.GetTagsByTables(appMapping.FromTables)
						if err := wList.AddNewToWaitingList(*node, adjTags, srcApp); err != nil {
							fmt.Println("!! ERROR WHILE TRYING TO ADD TO WAITING LIST !!", err)
							return false
						} else {
							log.Println("!! Node [", node.Tag.Name, "] added to a new waiting node!")
							return true
						}
					}
				}
				// break
				// return true
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
			log.Println("!! Couldn't find mappings for the node [", node.Tag.Name, "]")
			return false
		}
	}
	fmt.Println("~~~~~~~~~~~~ why here?")
	return false
}

func MigrateProcess(uid string, srcApp, dstApp config.AppConfig, node *m2.DependencyNode, wList *m2.WaitingList, invalidList *m2.InvalidList, log_txn *transaction.Log_txn) {

	if strings.EqualFold(node.Tag.Name, "root") && !checkUserInApp(uid, dstApp.AppID, log_txn) {
		log.Println("++ Adding User from ", srcApp.AppName, " to ", dstApp.AppName)
		addUserToApp(uid, dstApp.AppID, log_txn)
	}

	for child := GetAdjNode(node, srcApp, uid, wList, invalidList, log_txn); child != nil; child = GetAdjNode(node, srcApp, uid, wList, invalidList, log_txn) {
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		nodeIDAttr, _ := node.Tag.ResolveTagAttr("id")
		childIDAttr, _ := child.Tag.ResolveTagAttr("id")
		log.Println("~~ Current   Node: {", node.Tag.Name, "} ID:", node.Data[nodeIDAttr])
		log.Println("~~ Adjacent  Node: {", child.Tag.Name, "} ID:", child.Data[childIDAttr])
		MigrateProcess(uid, srcApp, dstApp, child, wList, invalidList, log_txn)
	}

	log.Println(fmt.Sprintf("## Migrating node { %s } From [%s] to [%s]", node.Tag.Name, srcApp.AppName, dstApp.AppName))

	if success := MigrateNode(node, srcApp, dstApp, wList, invalidList, log_txn); success {
		log.Println(fmt.Sprintf("xx Finished  node { %s } From [%s] to [%s]", node.Tag.Name, srcApp.AppName, dstApp.AppName))
	} else {
		log.Println(fmt.Sprintf("xx FAILED    node { %s } From [%s] to [%s]", node.Tag.Name, srcApp.AppName, dstApp.AppName))
	}

	fmt.Println("------------------------------------------------------------------------")
}

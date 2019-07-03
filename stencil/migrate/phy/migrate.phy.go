package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/helper"
	m2 "stencil/migrate"
	"stencil/qr"
	"stencil/transaction"
	"strings"
)

func RemoveUserFromApp(uid, app_id string, dbConn *sql.DB) bool {
	sql := "DELETE FROM user_table WHERE user_id = $1 AND app_id = $2"
	if err := db.Delete(dbConn, sql, uid, app_id); err == nil {
		return true
	}
	return false
}

func CheckUserInApp(uid, app_id string, dbConn *sql.DB) bool {
	sql := "SELECT user_id FROM user_table WHERE user_id = $1 AND app_id = $2"
	if res, err := db.DataCall1(dbConn, sql, uid, app_id); err == nil {
		if len(res) > 0 {
			return true
		}
	} else {
		log.Fatal(err)
	}
	return false
}

func AddUserToApp(uid, app_id string, dbConn *sql.DB) bool {
	query := "INSERT INTO user_table (user_id, app_id) VALUES ($1, $2)"
	dbConn.Exec(query, uid, app_id)
	return true
}

func UpdateMigrationState(uid string, srcApp, dstApp config.AppConfig) {

}

func GetRoot(appConfig config.AppConfig, uid string, dbConn *sql.DB) (*m2.DependencyNode, error) {
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
		qs.WhereMFlag(qr.EXISTS, "0", appConfig.AppID)
		sql = qs.GenSQL()
		rootNode := new(m2.DependencyNode)
		rootNode.Tag = root
		rootNode.SQL = sql
		if rootNode.Data, err = db.DataCall1(dbConn, sql); err == nil && len(rootNode.Data) > 0 {
			return rootNode, nil
		} else {
			return nil, err
		}
	} else {
		log.Fatal("Can't fetch tag:", tagName)
	}
	return nil, nil
}

func GetAdjNode(node *m2.DependencyNode, appConfig config.AppConfig, uid string, wList *m2.WaitingList, dbConn *sql.DB, unmappedTags *m2.UnmappedTags) (*m2.DependencyNode, error) {

	for _, dep := range config.ShuffleDependencies(appConfig.GetSubDependencies(node.Tag.Name)) {
		// for _, dep := range appConfig.GetSubDependencies(node.Tag.Name) {
		// if where := ResolveDependencyConditions(node, appConfig, dep); where != "" {
		if unmappedTags.Exists(dep.Tag) {
			continue
		}
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
			}
			if len(child.Restrictions) > 0 {
				restrictions := qr.CreateQS(appConfig.QR)
				restrictions.TableAliases = qs.TableAliases
				for _, restriction := range child.Restrictions {
					if restrictionAttr, err := child.ResolveTagAttr(restriction["col"]); err == nil {
						restrictions.WhereOperatorInterface("OR", restrictionAttr, "=", restriction["val"])
					}

				}
				qs.WhereString("AND", restrictions.Where)
			}
			where := ResolveDependencyConditions(node, appConfig, dep, child, qs)
			qs.WhereString("AND", where)
			qs.WhereMFlag(qr.EXISTS, "0", appConfig.AppID)
			qs.OrderBy("random()")
			sql = qs.GenSQL()
			if nodeData, err := db.DataCall1(dbConn, sql); err == nil {
				if len(nodeData) > 0 {
					newNode := new(m2.DependencyNode)
					newNode.Tag = child
					newNode.SQL = sql
					newNode.Data = nodeData
					if !wList.IsAlreadyWaiting(*newNode) {
						return newNode, nil
					}
				}
			} else {
				return nil, err
			}
		}
	}
	return nil, nil
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

func UpdateRowDesc(mappings *config.MappedApp, dstApp config.AppConfig, toTables []config.ToTable, node *m2.DependencyNode, dbConn *sql.DB) error {

	if tx, err := dbConn.Begin(); err != nil {
		log.Println("Can't create UpdateRowDesc transaction!")
		return errors.New("0")
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
							if err := db.SetMFlag(tx, pk, "1"); err == nil {
								updated = append(updated, pk)
							} else {
								tx.Rollback()
								fmt.Println("\n@ERROR_SET_MFLAG:", err)
								return err
							}
						} else {
							tx.Rollback()
							fmt.Println("\n@ERROR_SET_APPID:", err)
							return err
						}
					}
				}
			}
		}
		// tx.Rollback()
		tx.Commit()
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

func MigrateNode(node *m2.DependencyNode, srcApp, dstApp config.AppConfig, wList *m2.WaitingList, log_txn *transaction.Log_txn, dbConn *sql.DB, unmappedTags *m2.UnmappedTags) error {
	if mappings := config.GetSchemaMappingsFor(srcApp.AppName, dstApp.AppName); mappings == nil {
		log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcApp.AppName, dstApp.AppName))
	} else {
		mappingFound := false
		for _, appMapping := range mappings.Mappings {
			tagMembers := node.Tag.GetTagMembers()
			if mappedTables := helper.IntersectString(tagMembers, appMapping.FromTables); len(mappedTables) > 0 {
				mappingFound = true
				if len(tagMembers) == len(appMapping.FromTables) {
					return UpdateRowDesc(mappings, dstApp, appMapping.ToTables, node, dbConn)
				} else {
					log.Println("!! Node [", node.Tag.Name, "] needs to wait?")
					if waitingNode, err := wList.UpdateIfBeingLookedFor(*node); err == nil {
						log.Println("!! Node [", node.Tag.Name, "] updated an existing waiting node!")
						if waitingNode.IsComplete() {
							log.Println("!! Node [", node.Tag.Name, "] completed a waiting node!")
							tempCombinedDataDependencyNode := waitingNode.GenDependencyDataNode()
							return UpdateRowDesc(mappings, dstApp, appMapping.ToTables, &tempCombinedDataDependencyNode, dbConn)
						} else {
							log.Println("!! Node [", node.Tag.Name, "] added to an incomplete waiting node!")
							return nil
						}
					} else {
						adjTags := srcApp.GetTagsByTables(appMapping.FromTables)
						return wList.AddNewToWaitingList(*node, adjTags, srcApp)
						// if err := wList.AddNewToWaitingList(*node, adjTags, srcApp); err != nil {
						// 	fmt.Println("!! ERROR WHILE TRYING TO ADD TO WAITING LIST !!", err)
						// 	return err
						// } else {
						// 	log.Println("!! Node [", node.Tag.Name, "] added to a new waiting node!")
						// 	return true
						// }
					}
				}
				// break
				// return true
			}
		}
		if !mappingFound {
			// set m flag?
			for _, tagMember := range node.Tag.Members {
				if _, ok := node.Data[fmt.Sprintf("%s.id", tagMember)]; ok {
					srcID := fmt.Sprint(node.Data[fmt.Sprintf("%s.id", tagMember)])
					if serr := db.SaveForEvaluation(dbConn, "diaspora", dstApp.AppName, tagMember, "n/a", srcID, "n/a", "*", "n/a", fmt.Sprint(log_txn.Txn_id)); serr != nil {
						log.Fatal(serr)
					}
				}
			}
			log.Println("!! Couldn't find mappings for the tag [", node.Tag.Name, "]")
			unmappedTags.Add(node.Tag.Name)
			return errors.New("1")
		}
	}
	fmt.Println("~~~~~~~~~~~~ why here?")
	return nil
}

func MigrateProcess(uid string, srcApp, dstApp config.AppConfig, node *m2.DependencyNode, wList *m2.WaitingList, log_txn *transaction.Log_txn, dbConn *sql.DB, unmappedTags *m2.UnmappedTags, thread_id int) error {

	if unmappedTags.Exists(node.Tag.Name) {
		return nil
	}

	if strings.EqualFold(node.Tag.Name, "root") && !CheckUserInApp(uid, dstApp.AppID, dbConn) {
		log.Println("++ Adding User from ", srcApp.AppName, " to ", dstApp.AppName)
		AddUserToApp(uid, dstApp.AppID, dbConn)
	}

	for child, err := GetAdjNode(node, srcApp, uid, wList, dbConn, unmappedTags); child != nil; child, err = GetAdjNode(node, srcApp, uid, wList, dbConn, unmappedTags) {
		if err != nil {
			return err
		}
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		nodeIDAttr, _ := node.Tag.ResolveTagAttr("id")
		childIDAttr, _ := child.Tag.ResolveTagAttr("id")
		log.Println(fmt.Sprintf("~%d~ Current   Node: { %s } ID: %v", thread_id, node.Tag.Name, node.Data[nodeIDAttr]))
		log.Println(fmt.Sprintf("~%d~ Adjacent  Node: {%s} ID: %v", thread_id, child.Tag.Name, child.Data[childIDAttr]))
		if err := MigrateProcess(uid, srcApp, dstApp, child, wList, log_txn, dbConn, unmappedTags, thread_id); err != nil {
			return err
		}
	}

	log.Println(fmt.Sprintf("#%d# Migrating node { %s } From [%s] to [%s]", thread_id, node.Tag.Name, srcApp.AppName, dstApp.AppName))

	if err := MigrateNode(node, srcApp, dstApp, wList, log_txn, dbConn, unmappedTags); err == nil {
		log.Println(fmt.Sprintf("x%dx Finished  node { %s } From [%s] to [%s]", thread_id, node.Tag.Name, srcApp.AppName, dstApp.AppName))
	} else {
		if strings.EqualFold(err.Error(), "0") {
			log.Println(err)
			return err
		}
		log.Println(fmt.Sprintf("x%dx FAILED    node { %s } From [%s] to [%s]", thread_id, node.Tag.Name, srcApp.AppName, dstApp.AppName))
	}

	fmt.Println("------------------------------------------------------------------------")

	return nil
}

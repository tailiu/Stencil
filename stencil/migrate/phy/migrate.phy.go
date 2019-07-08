package migrate

import (
	"database/sql"
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
	"strconv"
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

func ResolveDependencyConditions(node *m2.DependencyNode, appConfig config.AppConfig, dep config.Dependency, tag config.Tag, qs *qr.QS) {

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
	if where.Where != "" {
		qs.WhereString("AND", where.Where)
	}
}

func GetDataNodeQS(tag config.Tag, appConfig config.AppConfig) *qr.QS {
	qs := qr.CreateQS(appConfig.QR)
	if len(tag.InnerDependencies) > 0 {
		joinMap := tag.CreateInDepMap()
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
		table := tag.Members["member1"]
		qs = qr.CreateQS(appConfig.QR)
		qs.FromSimple(table)
		qs.ColPK(table)
		qs.ColSimple(table + ".*")
	}
	if len(tag.Restrictions) > 0 {
		restrictions := qr.CreateQS(appConfig.QR)
		restrictions.TableAliases = qs.TableAliases
		for _, restriction := range tag.Restrictions {
			if restrictionAttr, err := tag.ResolveTagAttr(restriction["col"]); err == nil {
				restrictions.WhereOperatorInterface("OR", restrictionAttr, "=", restriction["val"])
			}

		}
		qs.WhereString("AND", restrictions.Where)
	}
	return qs
}

func GetRoot(appConfig config.AppConfig, uid string, dbConn *sql.DB) (*m2.DependencyNode, error) {
	tagName := "root"
	if root, err := appConfig.GetTag(tagName); err == nil {
		qs := GetDataNodeQS(root, appConfig)
		rootTable, rootCol := appConfig.GetItemsFromKey(root, "root_id")
		qs.WhereSimpleVal(rootTable+"."+rootCol, "=", uid)
		qs.WhereMFlag(qr.EXISTS, "0", appConfig.AppID)
		sql := qs.GenSQL()
		if data, err := db.DataCall1(dbConn, sql); err == nil && len(data) > 0 {
			rootNode := new(m2.DependencyNode)
			rootNode.Tag = root
			rootNode.SQL = sql
			rootNode.Data = data
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
		if unmappedTags.Exists(dep.Tag) {
			continue
		}
		if child, err := appConfig.GetTag(dep.Tag); err == nil {
			qs := GetDataNodeQS(child, appConfig)
			ResolveDependencyConditions(node, appConfig, dep, child, qs)
			qs.WhereMFlag(qr.EXISTS, "0", appConfig.AppID)
			qs.OrderBy("random()")
			sql := qs.GenSQL()
			if data, err := db.DataCall1(dbConn, sql); err == nil {
				if len(data) > 0 {
					newNode := new(m2.DependencyNode)
					newNode.Tag = child
					newNode.SQL = sql
					newNode.Data = data
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

func UpdateRowDesc(mappings *config.MappedApp, srcApp, dstApp config.AppConfig, toTables []config.ToTable, node *m2.DependencyNode, dbConn *sql.DB, log_txn *transaction.Log_txn) error {

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
					pk := strconv.FormatInt(val.(int64), 10)
					if val != nil && !helper.Contains(updated, pk) {
						if err := db.MUpdate(tx, pk, "1", dstApp.AppID); err == nil {
							updated = append(updated, pk)
							if err := display.GenDisplayFlag(log_txn.DBconn, dstApp.AppName, toTable.Table, pk, false, log_txn.Txn_id); err != nil {
								log.Fatal("## DISPLAY ERROR!", err)
							}
						} else {
							tx.Rollback()
							fmt.Println("\n@ERROR_MUpdate:", err)
							log.Fatal("pk:", pk, "appid", dstApp.AppID)
							return err
						}
						// if err := db.SetAppID(tx, pk, dstApp.AppID); err == nil {
						// 	if err := db.SetMFlag(tx, pk, "1"); err == nil {
						// 		updated = append(updated, pk)
						// 	} else {
						// 		tx.Rollback()
						// 		fmt.Println("\n@ERROR_SET_MFLAG:", err)
						// 		return err
						// 	}
						// } else {
						// 	tx.Rollback()
						// 	fmt.Println("\n@ERROR_SET_APPID:", err)
						// 	return err
						// }
					}
				}
			}
		}
		// tx.Rollback()
		if undoActionJSON, err := transaction.GenUndoActionJSON(updated, srcApp.AppID, dstApp.AppID); err == nil {
			if log_err := transaction.LogChange(undoActionJSON, log_txn); log_err != nil {
				log.Fatal("UpdateRowDesc: unable to LogChange", log_err)
			}
		} else {
			log.Fatal("UpdateRowDesc: unable to GenUndoActionJSON", err)
		}
		tx.Commit()
	}
	return nil
}

func HandleWaitingList(mappings *config.MappedApp, appMapping config.Mapping, tagMembers []string, node *m2.DependencyNode, srcApp, dstApp config.AppConfig, wList *m2.WaitingList) (*m2.DependencyNode, error) {

	log.Println("!! Node [", node.Tag.Name, "] needs to wait?")
	log.Println("tagMembers:", tagMembers, "appMapping.FromTables", appMapping.FromTables)
	if waitingNode, err := wList.UpdateIfBeingLookedFor(node); err == nil {
		log.Println("!! Node [", node.Tag.Name, "] updated an EXISITNG waiting node!")
		if waitingNode.IsComplete() {
			log.Println("!! Node [", node.Tag.Name, "] COMPLETED a waiting node!")
			tempCombinedDataDependencyNode := waitingNode.GenDependencyDataNode()
			return tempCombinedDataDependencyNode, nil
		}
		return nil, errors.New("1")
	}
	adjTags := srcApp.GetTagsByTables(appMapping.FromTables)
	if err := wList.AddNewToWaitingList(node, adjTags, srcApp); err == nil {
		log.Println("!! Node [", node.Tag.Name, "] added to a NEW waiting node!")
		return nil, errors.New("1")
	} else {
		log.Println("!! Node [", node.Tag.Name, "] ", err)
		return nil, err
	}
}

func HandleUnmappedTags(node *m2.DependencyNode, dbConn *sql.DB, unmappedTags *m2.UnmappedTags, dstApp config.AppConfig, log_txn *transaction.Log_txn) error {
	log.Println("!! Couldn't find mappings for the tag [", node.Tag.Name, "]")
	unmappedTags.Add(node.Tag.Name)
	// save for evaluation
	for _, tagMember := range node.Tag.Members {
		if _, ok := node.Data[fmt.Sprintf("%s.id", tagMember)]; ok {
			srcID := fmt.Sprint(node.Data[fmt.Sprintf("%s.id", tagMember)])
			if serr := db.SaveForEvaluation(dbConn, "diaspora", dstApp.AppName, tagMember, "n/a", srcID, "n/a", "*", "n/a", fmt.Sprint(log_txn.Txn_id)); serr != nil {
				log.Fatal(serr)
			}
		}
	}
	return errors.New("2")
}

func HandleUnmappedNode(node *m2.DependencyNode, dbConn *sql.DB, log_txn *transaction.Log_txn) error {
	if tx, err := dbConn.Begin(); err != nil {
		log.Println("Can't create UpdateRowDesc transaction!")
		return errors.New("0")
	} else {
		var updated []string
		for col, val := range node.Data {
			if strings.Contains(col, "pk.") {
				pk := fmt.Sprint(val)
				if val != nil && !helper.Contains(updated, pk) {
					if err := db.SetMFlag(tx, pk, "1"); err == nil {
						updated = append(updated, pk)
					} else {
						tx.Rollback()
						fmt.Println("\n@ERROR_MUpdate:", err)
						log.Fatal("pk:", pk, "appid: unmapped")
						return err
					}
				}
			}
		}
		if undoActionJSON, err := transaction.GenUndoActionJSON(updated, "0", "0"); err == nil {
			if log_err := transaction.LogChange(undoActionJSON, log_txn); log_err != nil {
				log.Fatal("HandleUnmappedNode: unable to LogChange", log_err)
			}
		} else {
			log.Fatal("HandleUnmappedNode: unable to GenUndoActionJSON", err)
		}
		tx.Commit()
		return nil
	}
}

func MigrateNode(node *m2.DependencyNode, srcApp, dstApp config.AppConfig, wList *m2.WaitingList, log_txn *transaction.Log_txn, dbConn *sql.DB, unmappedTags *m2.UnmappedTags) error {
	mappings := config.GetSchemaMappingsFor(srcApp.AppName, dstApp.AppName)
	if mappings == nil {
		log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcApp.AppName, dstApp.AppName))
	}
	for _, appMapping := range mappings.Mappings {
		tagMembers := node.Tag.GetTagMembers()
		if mappedTables := helper.IntersectString(tagMembers, appMapping.FromTables); len(mappedTables) > 0 {
			if helper.Sublist(tagMembers, appMapping.FromTables) {
				return UpdateRowDesc(mappings, srcApp, dstApp, appMapping.ToTables, node, dbConn, log_txn)
			}
			if wNode, err := HandleWaitingList(mappings, appMapping, tagMembers, node, srcApp, dstApp, wList); wNode != nil && err == nil {
				return UpdateRowDesc(mappings, srcApp, dstApp, appMapping.ToTables, wNode, dbConn, log_txn)
			} else {
				return err
			}
		}
	}
	return HandleUnmappedNode(node, dbConn, log_txn)
	// return HandleUnmappedTags(node, dbConn, unmappedTags, dstApp, log_txn)
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
		log.Println(fmt.Sprintf("~%d~ Adjacent  Node: { %s } ID: %v", thread_id, child.Tag.Name, child.Data[childIDAttr]))
		if err := MigrateProcess(uid, srcApp, dstApp, child, wList, log_txn, dbConn, unmappedTags, thread_id); err != nil {
			return err
		}
	}

	log.Println(fmt.Sprintf("#%d# Process   node { %s } From [%s] to [%s]", thread_id, node.Tag.Name, srcApp.AppName, dstApp.AppName))
	if err := MigrateNode(node, srcApp, dstApp, wList, log_txn, dbConn, unmappedTags); err == nil {
		log.Println(fmt.Sprintf("x%dx MIGRATED  node { %s } From [%s] to [%s]", thread_id, node.Tag.Name, srcApp.AppName, dstApp.AppName))
	} else {
		if strings.EqualFold(err.Error(), "2") {
			log.Println(fmt.Sprintf("x%dx IGNORED   node { %s } From [%s] to [%s]", thread_id, node.Tag.Name, srcApp.AppName, dstApp.AppName))
		} else {
			log.Println(fmt.Sprintf("x%dx FAILED    node { %s } From [%s] to [%s]", thread_id, node.Tag.Name, srcApp.AppName, dstApp.AppName))
			if strings.EqualFold(err.Error(), "0") {
				log.Println(err)
				return err
			}
		}
	}

	fmt.Println("------------------------------------------------------------------------")

	return nil
}

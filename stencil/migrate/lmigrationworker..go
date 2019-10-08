package migrate

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/app_display"
	"stencil/helper"
	"stencil/transaction"
	"strings"
)

func CreateLMigrationWorkerWithAppsConfig(uid, srcApp, srcAppID, dstApp, dstAppID string, logTxn *transaction.Log_txn, mtype string, srcAppConfig, dstAppConfig config.AppConfig) LMigrationWorker {
	mappings := config.GetSchemaMappingsFor(srcApp, dstApp)
	if mappings == nil {
		log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcApp, dstApp))
	}
	dstAppConfig.QR.Migration = true
	srcAppConfig.QR.Migration = true
	mWorker := LMigrationWorker{
		uid:          uid,
		SrcAppConfig: srcAppConfig,
		DstAppConfig: dstAppConfig,
		mappings:     mappings,
		wList:        WaitingList{},
		unmappedTags: CreateUnmappedTags(),
		SrcDBConn:    db.GetDBConn(srcApp),
		DstDBConn:    db.GetDBConn2(dstApp),
		logTxn:       &transaction.Log_txn{DBconn: logTxn.DBconn, Txn_id: logTxn.Txn_id},
		mtype:        mtype,
		visitedNodes: make(map[string]bool)}
	if err := mWorker.FetchRoot(); err != nil {
		log.Fatal(err)
	}
	return mWorker
}

func CreateLMigrationWorker(uid, srcApp, srcAppID, dstApp, dstAppID string, logTxn *transaction.Log_txn, mtype string, mappings *config.MappedApp) LMigrationWorker {
	srcAppConfig, err := config.CreateAppConfig(srcApp, srcAppID)
	if err != nil {
		log.Fatal(err)
	}
	dstAppConfig, err := config.CreateAppConfig(dstApp, dstAppID)
	if err != nil {
		log.Fatal(err)
	}
	dstAppConfig.QR.Migration = true
	srcAppConfig.QR.Migration = true
	mWorker := LMigrationWorker{
		uid:          uid,
		SrcAppConfig: srcAppConfig,
		DstAppConfig: dstAppConfig,
		mappings:     mappings,
		wList:        WaitingList{},
		unmappedTags: CreateUnmappedTags(),
		SrcDBConn:    db.GetDBConn(srcApp),
		DstDBConn:    db.GetDBConn2(dstApp),
		logTxn:       &transaction.Log_txn{DBconn: logTxn.DBconn, Txn_id: logTxn.Txn_id},
		mtype:        mtype,
		visitedNodes: make(map[string]bool)}
	if err := mWorker.FetchRoot(); err != nil {
		log.Fatal(err)
	}
	return mWorker
}

func (self *LMigrationWorker) RenewDBConn() {
	if self.SrcDBConn != nil {
		self.SrcDBConn.Close()
	}
	if self.DstDBConn != nil {
		self.DstDBConn.Close()
	}
	self.SrcDBConn = db.GetDBConn(self.SrcAppConfig.AppName)
	self.DstDBConn = db.GetDBConn2(self.DstAppConfig.AppName)
}

func (self *LMigrationWorker) Finish() {
	self.SrcDBConn.Close()
	self.DstDBConn.Close()
}

func (self *LMigrationWorker) GetRoot() *DependencyNode {
	return self.root
}

func (self *LMigrationWorker) MType() string {
	return self.mtype
}

func (self *LMigrationWorker) UserID() string {
	return self.uid
}

func (self *LMigrationWorker) MigrationID() int {
	return self.logTxn.Txn_id
}

func (self *LMigrationWorker) ResolveDependencyConditions(node *DependencyNode, dep config.Dependency, tag config.Tag) string {

	where := ""
	for _, depOn := range dep.DependsOn {
		if depOnTag, err := self.SrcAppConfig.GetTag(depOn.Tag); err == nil {
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

func (self *LMigrationWorker) ResolveOwnershipConditions(own config.Ownership, tag config.Tag) string {

	where := ""
	for _, condition := range own.Conditions {
		conditionStr := ""
		tagAttr, err := tag.ResolveTagAttr(condition.TagAttr)
		if err != nil {
			fmt.Println("data1", self.root.Data)
			log.Fatal(err, tag.Name, condition.TagAttr)
			break
		}
		depOnAttr, err := self.root.Tag.ResolveTagAttr(condition.DependsOnAttr)
		if err != nil {
			fmt.Println("data2", self.root.Data)
			log.Fatal(err, tag.Name, condition.DependsOnAttr)
			break
		}
		if _, ok := self.root.Data[depOnAttr]; ok {
			if conditionStr != "" || where != "" {
				conditionStr += " AND "
			}
			conditionStr += fmt.Sprintf("%s = '%v'", tagAttr, self.root.Data[depOnAttr])
		} else {
			fmt.Println("data3", self.root.Data)
			log.Fatal("ResolveOwnershipConditions:", depOnAttr, " doesn't exist in ", tag.Name)
		}
		where += conditionStr
	}
	return where
}

func (self *LMigrationWorker) ResolveRestrictions(tag config.Tag) string {
	restrictions := ""
	if len(tag.Restrictions) > 0 {
		restrictions := ""
		for _, restriction := range tag.Restrictions {
			if restrictionAttr, err := tag.ResolveTagAttr(restriction["col"]); err == nil {
				if restrictions != "" {
					restrictions += " AND "
				}
				restrictions += fmt.Sprintf("%s = '%s'", restrictionAttr, restriction["val"])
			}

		}
	}
	return restrictions
}

func (self *LMigrationWorker) GetTagQL(tag config.Tag) (string, string) {

	sql := "SELECT %s FROM %s "
	where_mad := ""

	if len(tag.InnerDependencies) > 0 {
		cols := ""
		joinMap := tag.CreateInDepMap()
		seenMap := make(map[string]bool)
		joinStr := ""
		for fromTable, toTablesMap := range joinMap {
			if _, ok := seenMap[fromTable]; !ok {
				joinStr += fromTable
				_, colStr := db.GetColumnsForTable(self.SrcDBConn, fromTable)
				cols += colStr + ","
				if where_mad != "" {
					where_mad += " AND "
				}
				where_mad += fmt.Sprintf("%s.mark_as_delete = false", fromTable)
			}
			for toTable, conditions := range toTablesMap {
				if conditions != nil {
					conditions = append(conditions, joinMap[toTable][fromTable]...)
					if joinMap[toTable][fromTable] != nil {
						joinMap[toTable][fromTable] = nil
					}
					joinStr += fmt.Sprintf(" JOIN %s ON %s ", toTable, strings.Join(conditions, " AND "))
					_, colStr := db.GetColumnsForTable(self.SrcDBConn, toTable)
					cols += colStr + ","
					seenMap[toTable] = true
					if where_mad != "" {
						where_mad += " AND "
					}
					where_mad += fmt.Sprintf("%s.mark_as_delete = false", toTable)
				}
			}
			seenMap[fromTable] = true
		}
		sql = fmt.Sprintf(sql, strings.Trim(cols, ","), joinStr)
	} else {
		table := tag.Members["member1"]
		_, cols := db.GetColumnsForTable(self.SrcAppConfig.DBConn, table)
		sql = fmt.Sprintf(sql, cols, table)
		where_mad = fmt.Sprintf("%s.mark_as_delete = false", table)
	}
	return sql, where_mad
}

func (self *LMigrationWorker) FetchRoot() error {
	tagName := "root"
	if root, err := self.SrcAppConfig.GetTag(tagName); err == nil {
		rootTable, rootCol := self.SrcAppConfig.GetItemsFromKey(root, "root_id")
		where := fmt.Sprintf("%s.%s = '%s'", rootTable, rootCol, self.uid)
		ql, wmad := self.GetTagQL(root)
		sql := fmt.Sprintf("%s WHERE %s AND %s", ql, where, wmad)
		if restrictions := self.ResolveRestrictions(root); restrictions != "" {
			sql += restrictions
		}
		// log.Fatal(sql)
		if data, err := db.DataCall1(self.SrcDBConn, sql); err == nil && len(data) > 0 {
			rootNode := new(DependencyNode)
			rootNode.Tag = root
			rootNode.SQL = sql
			rootNode.Data = data
			self.root = rootNode
			return nil
		} else {
			if err == nil {
				err = errors.New("no data returned for root node, doesn't exist?")
			}
			return err
		}
	} else {
		log.Fatal("Can't fetch root tag:", err)
		return err
	}
}

func (self *LMigrationWorker) GetAdjNode(node *DependencyNode, threadID int) (*DependencyNode, error) {

	for _, dep := range self.SrcAppConfig.ShuffleDependencies(self.SrcAppConfig.GetSubDependencies(node.Tag.Name)) {
		if child, err := self.SrcAppConfig.GetTag(dep.Tag); err == nil {
			log.Println(fmt.Sprintf("x%2dx | FETCHING  tag  { %s } ", threadID, dep.Tag))
			where := self.ResolveDependencyConditions(node, dep, child)
			ql, wmad := self.GetTagQL(child)
			sql := fmt.Sprintf("%s WHERE %s AND %s", ql, where, wmad)
			if restrictions := self.ResolveRestrictions(child); restrictions != "" {
				sql += restrictions
			}
			sql += " ORDER BY random()"
			// log.Fatal(sql)
			if data, err := db.DataCall(self.SrcDBConn, sql); err == nil {
				for _, datum := range data {
					newNode := new(DependencyNode)
					newNode.Tag = child
					newNode.SQL = sql
					newNode.Data = datum
					if !self.wList.IsAlreadyWaiting(*newNode) && !self.IsVisited(newNode){
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

func (self *LMigrationWorker) GetOwnedNodes(threadID, nodelimit int) ([]*DependencyNode, error) {

	for _, own := range self.SrcAppConfig.GetOrderedOwnerships() {
		log.Println(fmt.Sprintf("x%2dx |         FETCHING  tag  { %s } ", threadID, own.Tag))
		if self.unmappedTags.Exists(own.Tag) {
			log.Println(fmt.Sprintf("x%2dx |         UNMAPPED  tag  { %s } ", threadID, own.Tag))
			continue
		}
		if child, err := self.SrcAppConfig.GetTag(own.Tag); err == nil {
			where := self.ResolveOwnershipConditions(own, child)
			ql, wmad := self.GetTagQL(child)
			sql := fmt.Sprintf("%s WHERE %s AND %s", ql, where, wmad)
			if restrictions := self.ResolveRestrictions(child); restrictions != "" {
				sql += restrictions
			}
			sql += " ORDER BY random() LIMIT " + fmt.Sprint(nodelimit)
			// log.Fatal(sql)
			if result, err := db.DataCall(self.SrcDBConn, sql); err == nil {
				var nodes []*DependencyNode
				for _, data := range result {
					if len(data) > 0 {
						newNode := new(DependencyNode)
						newNode.Tag = child
						newNode.SQL = sql
						newNode.Data = data
						if !self.wList.IsAlreadyWaiting(*newNode) {
							nodes = append(nodes, newNode)
						}
					}
				}
				if len(nodes) > 0 {
					return nodes, nil
				}
			} else {
				return nil, err
			}
		}
	}
	return nil, nil
}

func (self *LMigrationWorker) PushData(dtable config.ToTable, pk, orgCols, cols string, undoAction *transaction.UndoAction, node *DependencyNode) error {

	undoActionSerialized, _ := json.Marshal(undoAction)
	transaction.LogChange(string(undoActionSerialized), self.logTxn)
	if err := app_display.GenDisplayFlag(self.logTxn.DBconn, self.DstAppConfig.AppID, dtable.TableID, pk, fmt.Sprint(self.logTxn.Txn_id)); err != nil {
		log.Fatal("## DISPLAY ERROR!", err)
		return errors.New("0")
	}

	for _, fromTable := range undoAction.OrgTables {
		if _, ok := node.Data[fmt.Sprintf("%s.id", fromTable)]; ok {
			srcID := fmt.Sprint(node.Data[fmt.Sprintf("%s.id", fromTable)])
			if serr := db.SaveForLEvaluation(self.logTxn.DBconn, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTable, dtable.Table, srcID, pk, orgCols, cols, fmt.Sprint(self.logTxn.Txn_id)); serr != nil {
				log.Fatal(serr)
				return errors.New("0")
			}
		}

	}
	return nil
}

func (self *LMigrationWorker) CheckMappingConditions(toTable config.ToTable, node *DependencyNode) bool {
	breakCondition := false
	if len(toTable.Conditions) > 0 {
		for conditionKey, conditionVal := range toTable.Conditions {
			if nodeVal, ok := node.Data[conditionKey]; ok {
				if conditionVal[:1] == "#" {
					// fmt.Println("VerifyMappingConditions: conditionVal[:1] == #")
					// fmt.Println(conditionKey, conditionVal, nodeVal)
					// fmt.Scanln()
					switch conditionVal {
						case "#NULL": {
							if nodeVal != nil {
								// log.Println(nodeVal, "!=", conditionVal)
								// fmt.Println(conditionKey, conditionVal, nodeVal)
								// log.Fatal("@VerifyMappingConditions: return false, from case #NULL:")
								return false
							}
						}
						case "#NOTNULL": {
							if nodeVal == nil {
								// log.Println(nodeVal, "!=", conditionVal)
								// fmt.Println(conditionKey, conditionVal, nodeVal)
								// log.Fatal("@VerifyMappingConditions: return false, from case #NOTNULL:")
								return false
							}		
						}
						default: {
							fmt.Println(toTable.Table, conditionKey, conditionVal)
							log.Fatal("@CheckMappingConditions: Case not found:" + conditionVal)
						}
					}
				} else if conditionVal[:1] == "$" {
					// fmt.Println("VerifyMappingConditions: conditionVal[:1] == $")
					// fmt.Println(conditionKey, conditionVal, nodeVal)
					// fmt.Scanln()
					if inputVal, err := self.mappings.GetInput(conditionVal); err == nil {
						if !strings.EqualFold(fmt.Sprint(nodeVal), inputVal) {
							log.Println(nodeVal, "!=", inputVal)
							fmt.Println(conditionKey, conditionVal, inputVal, nodeVal)
							log.Fatal("@CheckMappingConditions: return false, from conditionVal[:1] == $")
							return false
						}
					}else {
						fmt.Println("node data:", node.Data)
						fmt.Println(conditionKey, conditionVal)
						log.Fatal("@CheckMappingConditions: input doesn't exist?", err)
					}
				} else if !strings.EqualFold(fmt.Sprint(nodeVal), conditionVal) {
					breakCondition = true
					// log.Println(conditionKey, conditionVal, "!=", nodeVal)
					return true
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

func (self *LMigrationWorker) IsNodeOwnedByRoot(node *DependencyNode) bool {

	if ownership := self.SrcAppConfig.GetOwnership(node.Tag.Name, self.root.Tag.Name); ownership != nil {
		for _, condition := range ownership.Conditions{	
			tagAttr, err := node.Tag.ResolveTagAttr(condition.TagAttr)
			if err != nil {
				log.Fatal("Resolving TagAttr in IsNodeOwnedByRoot", err, node.Tag.Name, condition.TagAttr)
				break
			}
			depOnAttr, err := self.root.Tag.ResolveTagAttr(condition.DependsOnAttr)
			if err != nil {
				log.Fatal("Resolving depOnAttr in IsNodeOwnedByRoot",err, node.Tag.Name, condition.DependsOnAttr)
				break
			}
			if nodeVal, err := node.GetValueForKey(tagAttr); err == nil {
				if rootVal, err := self.root.GetValueForKey(depOnAttr); err == nil {
					if !strings.EqualFold(nodeVal, rootVal) {
						return false
					} 
					// else {
					// 	fmt.Println(nodeVal, "==", rootVal)
					// }
				}else{
					fmt.Println("Ownership Condition Key in Root Data:", depOnAttr, "doesn't exist!")
					fmt.Println("root data:", self.root.Data)
					log.Fatal("stop here and check ownership conditions wrt root")	
				}
			} else {
				fmt.Println("Ownership Condition Key", tagAttr, "doesn't exist!")
				fmt.Println("node data:", node.Data)
				fmt.Println("node sql:", node.SQL)
				log.Fatal("stop here and check ownership conditions")
			}
		}
	}else{
		// log.Fatal("Ownership not found in IsNodeOwnedByRoot:", node.Tag.Name)
	}
	return true
}

func (self *LMigrationWorker) GetMappedData(toTable config.ToTable, node *DependencyNode) (string, string, []interface{}, string, string, *transaction.UndoAction) {
	undoAction := new(transaction.UndoAction)
	cols, vals := "", ""
	orgCols, orgColsLeft := "", ""
	var ivals []interface{}
	for toAttr, fromAttr := range toTable.Mapping {
		if val, ok := node.Data[fromAttr]; ok {
			ivals = append(ivals, val)
			vals += fmt.Sprintf("$%d,", len(ivals))
			cols += fmt.Sprintf("%s,", toAttr)
			orgCols += fmt.Sprintf("%s,", strings.Split(fromAttr, ".")[1])
			undoAction.AddData(fromAttr, val)
			undoAction.AddOrgTable(strings.Split(fromAttr, ".")[0])
		} else if strings.Contains(fromAttr, "$") {
			if inputVal, err := self.mappings.GetInput(fromAttr); err == nil {
				ivals = append(ivals, inputVal)
				vals += fmt.Sprintf("$%d,", len(ivals))
				cols += fmt.Sprintf("%s,", toAttr)
				orgCols += fmt.Sprintf("%s,", fromAttr)
			}
		} else if strings.Contains(fromAttr, "#") {
			assignedTabCol := strings.Trim(fromAttr, "#ASSIGN()")
			if strings.Contains(fromAttr, "#ASSIGN"){
				// fmt.Println(fromAttr)
				if nodeVal, ok := node.Data[assignedTabCol]; ok {
					// fmt.Println("found attr: ", assignedTabCol, " ", nodeVal)
					ivals = append(ivals, nodeVal)
					vals += fmt.Sprintf("$%d,", len(ivals))
					cols += fmt.Sprintf("%s,", toAttr)
					orgCols += fmt.Sprintf("%s,", assignedTabCol)
				}
				// fmt.Println("IN: ", strings.Trim(cols, ","), strings.Trim(vals, ","), ivals, strings.Trim(orgCols, ","), orgColsLeft, undoAction)
				// log.Fatal("found attr: ", assignedTabCol)
			}else{
				switch fromAttr {
					case "#GUID": {
						ivals = append(ivals, uuid.New())
						vals += fmt.Sprintf("$%d,", len(ivals))
						cols += fmt.Sprintf("%s,", toAttr)
						orgCols += fmt.Sprintf("%s,", assignedTabCol)
					}
					case "#RANDINT": {
						ivals = append(ivals, self.SrcAppConfig.QR.NewRowId())
						vals += fmt.Sprintf("$%d,", len(ivals))
						cols += fmt.Sprintf("%s,", toAttr)
						orgCols += fmt.Sprintf("%s,", assignedTabCol)
					}
					default: {
						fmt.Println(toTable.Table, toAttr, fromAttr)
						log.Fatal("@GetMappedData: Case not found:" + fromAttr)
					}
				}
			}
			// log.Fatal(fromAttr)
		} else {
			orgColsLeft += fmt.Sprintf("%s,", strings.Split(fromAttr, ".")[1])
		}
	}
	// fmt.Println(strings.Trim(cols, ","), strings.Trim(vals, ","), ivals, strings.Trim(orgCols, ","), orgColsLeft, undoAction)
	return strings.Trim(cols, ","), strings.Trim(vals, ","), ivals, strings.Trim(orgCols, ","), orgColsLeft, undoAction
}

func (self *LMigrationWorker) MarkRowAsDeleted(node *DependencyNode, tx *sql.Tx) error {
	for _, tagMember := range node.Tag.Members {
		idCol := fmt.Sprintf("%s.id", tagMember)
		if _, ok := node.Data[idCol]; ok {
			srcID := fmt.Sprint(node.Data[idCol])
			if derr := db.DeleteRowFromAppDB(tx, tagMember, srcID); derr != nil {
				fmt.Println("@ERROR_DeleteRowFromAppDB", derr)
				fmt.Println("@QARGS:", tagMember, srcID)
				// log.Fatal(derr)
				return derr
			}
			if derr := db.UpdateLEvaluation(self.logTxn.DBconn, tagMember, srcID, self.logTxn.Txn_id); derr != nil {
				fmt.Println("@ERROR_UpdateLEvaluation", derr)
				fmt.Println("@QARGS:", tagMember, srcID, self.logTxn.Txn_id)
				// log.Fatal(derr)
				return derr
			}
		} else {
			log.Println("node.Data =>", node.Data)
			log.Fatal(idCol, "NOT PRESENT IN NODE DATA")
		}
	}
	return nil
}

func (self *LMigrationWorker) HandleMigration(toTables []config.ToTable, node *DependencyNode) error {

	srctx, err := self.SrcDBConn.Begin()
	if err != nil {
		log.Println("Can't create SrcDBConn transaction!")
		return errors.New("0")
	}
	defer srctx.Rollback()

	dsttx, err := self.DstDBConn.Begin()
	if err != nil {
		log.Println("Can't create DstDBConn transaction!")
		return errors.New("0")
	}
	defer dsttx.Rollback()

	for _, toTable := range toTables {
		if self.CheckMappingConditions(toTable, node) {
			continue
		}
		if cols, placeholders, ivals, orgCols, _, undoAction := self.GetMappedData(toTable, node); len(cols) > 0 && len(placeholders) > 0 && len(ivals) > 0 {
			undoAction.AddDstTable(toTable.Table)
			// if strings.Contains(toTable.Table, "status_"){
			// 	fmt.Println(toTable.Table, cols, placeholders, ivals)
			// 	log.Fatal("--------------")
			// }
			if id, err := db.InsertRowIntoAppDB(dsttx, toTable.Table, cols, placeholders, ivals...); err == nil {
				if err := self.PushData(toTable, fmt.Sprint(id), orgCols, cols, undoAction, node); err != nil {
					fmt.Println("@ERROR_PushData")
					fmt.Println("@Params:", toTable.Table, fmt.Sprint(id), orgCols, cols, undoAction, node)
					log.Fatal(err)
					return err
				}
			} else {
				if !strings.Contains(err.Error(), "duplicate key value") {
					fmt.Println("@ERROR_Insert")
					fmt.Println("@QARGS:", toTable.Table, cols, placeholders, ivals)
					log.Fatal(err)
					return err
				} else {
					// log.Println("@Already_Exists in:", toTable.Table, node.Data)
					if err := self.MarkRowAsDeleted(node, srctx); err == nil {
						srctx.Commit()
						dsttx.Commit()
					}
					return errors.New("3")
					// break
				}
			}
		} else {
			log.Fatal("@ERROR_GetMappedData:", cols, placeholders, ivals, orgCols, undoAction)
		}
	}

	// if self.mtype == DELETION {
	if err := self.MarkRowAsDeleted(node, srctx); err == nil {
		srctx.Commit()
		dsttx.Commit()
	}
	// }

	return nil
}

func (self *LMigrationWorker) HandleWaitingList(appMapping config.Mapping, tagMembers []string, node *DependencyNode) (*DependencyNode, error) {

	srctx, err := self.SrcDBConn.Begin()
	if err != nil {
		log.Println("Can't create HandleWaitingList transaction!")
		return nil, errors.New("0")
	}
	if err := self.MarkRowAsDeleted(node, srctx); err != nil {
		fmt.Println("@ERROR_HandleWaitingList")
		fmt.Println("@Node:", node)
		log.Fatal(err)
	}
	srctx.Commit()

	log.Println("!! Node [", node.Tag.Name, "] needs to wait?")
	log.Println("tagMembers:", tagMembers, "appMapping.FromTables", appMapping.FromTables)
	if waitingNode, err := self.wList.UpdateIfBeingLookedFor(node); err == nil {
		log.Println("!! Node [", node.Tag.Name, "] updated an EXISITNG waiting node!")
		if waitingNode.IsComplete() {
			log.Println("!! Node [", node.Tag.Name, "] COMPLETED a waiting node!")
			return waitingNode.GenDependencyDataNode(), nil
		}
		return nil, errors.New("1")
	}
	adjTags := self.SrcAppConfig.GetTagsByTables(appMapping.FromTables)
	if err := self.wList.AddNewToWaitingList(node, adjTags, self.SrcAppConfig); err == nil {
		log.Println("!! Node [", node.Tag.Name, "] added to a NEW waiting node!")
		return nil, errors.New("1")
	} else {
		log.Println("!! Node [", node.Tag.Name, "] ", err)
		return nil, err
	}
}

func (self *LMigrationWorker) HandleUnmappedTags(node *DependencyNode) error {
	log.Println("!! Couldn't find mappings for the tag [", node.Tag.Name, "]")
	self.unmappedTags.Add(node.Tag.Name)
	// save for evaluation
	for _, tagMember := range node.Tag.Members {
		if _, ok := node.Data[fmt.Sprintf("%s.id", tagMember)]; ok {
			srcID := fmt.Sprint(node.Data[fmt.Sprintf("%s.id", tagMember)])
			if serr := db.SaveForEvaluation(self.logTxn.DBconn, self.SrcAppConfig.AppName, self.DstAppConfig.AppName, tagMember, "n/a", srcID, "n/a", "*", "n/a", fmt.Sprint(self.logTxn.Txn_id)); serr != nil {
				log.Fatal(serr)
			}
		}
	}
	return errors.New("2")
}

func (self *LMigrationWorker) HandleUnmappedNode(node *DependencyNode) error {
	if !strings.EqualFold(self.mtype, DELETION) {
		return errors.New("2")
	}
	if tx, err := self.SrcDBConn.Begin(); err != nil {
		log.Println("Can't create HandleUnmappedNode transaction!")
		return errors.New("0")
	} else {
		var updated []string
		for _, tagMember := range node.Tag.Members {
			idCol := fmt.Sprintf("%s.id", tagMember)
			if _, ok := node.Data[idCol]; ok {
				srcID := fmt.Sprint(node.Data[idCol])
				if derr := db.DeleteRowFromAppDB(tx, tagMember, srcID); derr != nil {
					fmt.Println("@ERROR_Delete")
					fmt.Println("@SQL:", tagMember, srcID)
					fmt.Println("@ARGS:", tagMember, srcID)
					// log.Fatal(derr)
					return derr
				}
				if serr := db.SaveForEvaluation(self.logTxn.DBconn, self.SrcAppConfig.AppName, self.DstAppConfig.AppName, tagMember, "n/a", srcID, "n/a", "*", "n/a", fmt.Sprint(self.logTxn.Txn_id)); serr != nil {
					log.Fatal("@ERROR_HandleUnmappedNode_SaveForEvaluation =>", serr)
				}
				updated = append(updated, srcID)
			} else {
				log.Println("node.Data =>", node.Data)
				log.Fatal(idCol, "NOT PRESENT IN NODE DATA")
			}
		}
		if undoActionJSON, err := transaction.GenUndoActionJSON(updated, "0", "0"); err == nil {
			if log_err := transaction.LogChange(undoActionJSON, self.logTxn); log_err != nil {
				log.Fatal("@ERROR_HandleUnmappedNode_LogChange =>", log_err)
			}
		} else {
			log.Fatal("@ERROR_HandleUnmappedNode_GenUndoActionJSON =>", err)
		}
		tx.Commit()
		return errors.New("2")
	}
}

func (self *LMigrationWorker) FetchMappingsForNode(node *DependencyNode) (config.Mapping, bool) {
	var combinedMapping config.Mapping
	tagMembers := node.Tag.GetTagMembers()
	mappingFound := false
	for _, mapping := range self.mappings.Mappings {	
		if mappedTables := helper.IntersectString(tagMembers, mapping.FromTables); len(mappedTables) > 0 {
			combinedMapping.FromTables = append(combinedMapping.FromTables, mapping.FromTables...)
			combinedMapping.ToTables = append(combinedMapping.ToTables, mapping.ToTables...)
			mappingFound = true
		}
	}
	return combinedMapping, mappingFound
}

func (self *LMigrationWorker) _MigrateNode(node *DependencyNode, isBag bool) error {

	for _, appMapping := range self.mappings.Mappings {
		tagMembers := node.Tag.GetTagMembers()
		if mappedTables := helper.IntersectString(tagMembers, appMapping.FromTables); len(mappedTables) > 0 {
			if helper.Sublist(tagMembers, appMapping.FromTables) {
				return self.HandleMigration(appMapping.ToTables, node)
			}
			if wNode, err := self.HandleWaitingList(appMapping, tagMembers, node); wNode != nil && err == nil {
				return self.HandleMigration(appMapping.ToTables, wNode)
			} else {
				return err
			}
		}
	}
	if isBag {
		return fmt.Errorf("no mapping found for bag: %s", node.Tag.Name)
	}
	if !strings.EqualFold(self.mtype, DELETION) {
		self.unmappedTags.Add(node.Tag.Name)
		return fmt.Errorf("no mapping found for node: %s", node.Tag.Name)
	}
	return self.HandleUnmappedNode(node)
}

func (self *LMigrationWorker) MigrateNode(node *DependencyNode, isBag bool) error {
	
	if mapping, found := self.FetchMappingsForNode(node); found {
		tagMembers := node.Tag.GetTagMembers()
		if helper.Sublist(tagMembers, mapping.FromTables) { // other mappings HANDLE!
			return self.HandleMigration(mapping.ToTables, node)
		}
		if wNode, err := self.HandleWaitingList(mapping, tagMembers, node); wNode != nil && err == nil {
			return self.HandleMigration(mapping.ToTables, wNode)
		} else {
			return err
		}
	} else {
		if isBag || !strings.EqualFold(self.mtype, DELETION) {
			self.unmappedTags.Add(node.Tag.Name)
			return fmt.Errorf("no mapping found for node: %s", node.Tag.Name)
		}
		return self.HandleUnmappedNode(node)
	}
}

func (self *LMigrationWorker) HandleLeftOverWaitingNodes() {

	for _, waitingNode := range self.wList.Nodes {
		for _, containedNode := range waitingNode.ContainedNodes {
			self.HandleUnmappedNode(containedNode)
		}
	}
}

func (self *LMigrationWorker) IsVisited(node *DependencyNode) bool{
	for _, tagMember := range node.Tag.Members {
		idCol := fmt.Sprintf("%s.id", tagMember)
		if _, ok := node.Data[idCol]; ok {
			srcID := fmt.Sprint(node.Data[idCol])
			if _, ok := self.visitedNodes[srcID]; ok {
				return true
			}
		} else {
			log.Println("In: IsVisited | node.Data =>", node.Data)
			log.Fatal(idCol, "NOT PRESENT IN NODE DATA")
		}
	}
	return false
}

func (self *LMigrationWorker) MarkAsVisited(node *DependencyNode) {
	for _, tagMember := range node.Tag.Members {
		idCol := fmt.Sprintf("%s.id", tagMember)
		if _, ok := node.Data[idCol]; ok {
			srcID := fmt.Sprint(node.Data[idCol])
			self.visitedNodes[srcID] = true
		} else {
			log.Println("In: MarkAsVisited | node.Data =>", node.Data)
			log.Fatal(idCol, "NOT PRESENT IN NODE DATA")
		}
	}
}

func (self *LMigrationWorker) DeletionMigration(node *DependencyNode, threadID int) error {

	// if strings.EqualFold(node.Tag.Name, "root") && !db.CheckUserInApp(self.uid, self.DstAppConfig.AppID, self.DstDBConn) {
	// 	log.Println("++ Adding User from ", self.SrcAppConfig.AppName, " to ", self.DstAppConfig.AppName)
	// 	db.AddUserToApp(self.uid, self.DstAppConfig.AppID, self.SrcDBConn)
	// }

	for child, err := self.GetAdjNode(node, threadID); child != nil; child, err = self.GetAdjNode(node, threadID) {
		if err != nil {
			return err
		}
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		nodeIDAttr, _ := node.Tag.ResolveTagAttr("id")
		childIDAttr, _ := child.Tag.ResolveTagAttr("id")
		log.Println(fmt.Sprintf("~%d~ Current   Node: { %s } ID: %v", threadID, node.Tag.Name, node.Data[nodeIDAttr]))
		log.Println(fmt.Sprintf("~%d~ Adjacent  Node: { %s } ID: %v", threadID, child.Tag.Name, child.Data[childIDAttr]))
		if err := self.DeletionMigration(child, threadID); err != nil {
			return err
		}
	}

	log.Println(fmt.Sprintf("#%d# Process   node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))

	if self.IsNodeOwnedByRoot(node){
		if err := self.MigrateNode(node, false); err == nil {
			log.Println(fmt.Sprintf("x%2dx MIGRATED  node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
		} else {
			if strings.EqualFold(err.Error(), "2") {
				log.Println(fmt.Sprintf("x%2dx IGNORED   node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
			} else {
				log.Println(fmt.Sprintf("x%2dx FAILED    node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
				if strings.EqualFold(err.Error(), "0") {
					log.Println(err)
					return err
				}
				if strings.Contains(err.Error(), "deadlock"){
					return err
				}
			}
		}
	}else{
		log.Println(fmt.Sprintf("x%2dx UN-OWNED  node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
	}

	self.MarkAsVisited(node)

	fmt.Println("------------------------------------------------------------------------")

	return nil
}

func (self *LMigrationWorker) SecondPhase(threadID int) error {

	nodelimit := 1
	for nodes, err := self.GetOwnedNodes(threadID, nodelimit); err != nil || nodes != nil; nodes, err = self.GetOwnedNodes(threadID, nodelimit) {
		if err != nil {
			return err
		}
		for _, node := range nodes {
			if err := self.DeletionMigration(node, threadID); err != nil {
				return err
			}
		}
	}
	return nil
}

// func (self *LMigrationWorker) FinishMigration(mtype string, number_of_threads int) bool {
// 	return db.FinishMigration(self.logTxn.DBconn, self.logTxn.Txn_id)
// }

func (self *LMigrationWorker) ConsistentMigration(threadID int) error {

	nodelimit := 1
	for nodes, err := self.GetOwnedNodes(threadID, nodelimit); err != nil || nodes != nil; nodes, err = self.GetOwnedNodes(threadID, nodelimit) {
		if err != nil {
			return err
		}
		totalNodes := len(nodes)
		existingNodesCount := 0 // consecutive 10 already exist? break loop!
		for nodeNum, node := range nodes {
			nodeIDAttr, _ := node.Tag.ResolveTagAttr("id")
			log.Println(fmt.Sprintf("~%2d~ | %d/%d | Current   Node: { %s } ID: %v", threadID, nodeNum, totalNodes, node.Tag.Name, node.Data[nodeIDAttr]))
			if err := self.MigrateNode(node, false); err == nil {
				existingNodesCount = 0
				log.Println(fmt.Sprintf("x%2dx | %d/%d | MIGRATED  node { %s } From [%s] to [%s]", threadID, nodeNum, totalNodes, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
			} else {
				log.Println(fmt.Sprintf("x%2dx | %d/%d | RCVD ERR  node { %s } From [%s] to [%s]", threadID, nodeNum, totalNodes, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName), err)
				if self.unmappedTags.Exists(node.Tag.Name) {
					log.Println(fmt.Sprintf("x%2dx | %d/%d | BREAKLOOP node { %s } From [%s] to [%s]", threadID, nodeNum, totalNodes, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName), err)
					break
				}
				if strings.EqualFold(err.Error(), "2") {
					log.Println(fmt.Sprintf("x%2dx | %d/%d | IGNORED   node { %s } From [%s] to [%s]", threadID, nodeNum, totalNodes, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
				} else if strings.EqualFold(err.Error(), "3") {
					existingNodesCount++
					log.Println(fmt.Sprintf("x%2dx | %d/%d | EXISTS    node { %s } From [%s] to [%s]", threadID, nodeNum, totalNodes, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
					if existingNodesCount > 10 {
						break
					}
				} else {
					log.Println(fmt.Sprintf("x%2dx | %d/%d | FAILED    node { %s } From [%s] to [%s]", threadID, nodeNum, totalNodes, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName), err)
					if strings.EqualFold(err.Error(), "0") {
						log.Println(err)
						return err
					}
					if strings.Contains(err.Error(), "deadlock"){
						return err
					}
				}
			}
		}
	}
	if err := self.MigrateNode(self.root, false); err == nil {
		log.Println(fmt.Sprintf("x%2dx | MIGRATED  node { %s } From [%s] to [%s]", threadID, self.root.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
	}else{
		log.Println(fmt.Sprintf("x%2dx | MIGRATED? node { %s } From [%s] to [%s]", threadID, self.root.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
	}
	return nil
}

func (self *LMigrationWorker) IndependentMigration(threadID int) error {

	return self.ConsistentMigration(threadID)
}

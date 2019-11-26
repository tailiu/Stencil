package migrate

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"stencil/app_display"
	"stencil/config"
	"stencil/db"
	"stencil/helper"
	"stencil/transaction"
	"strings"

	"github.com/google/uuid"
)

func CreateMigrationWorkerV2WithAppsConfig(uid, srcApp, srcAppID, dstApp, dstAppID string, logTxn *transaction.Log_txn, mtype string, srcAppConfig, dstAppConfig config.AppConfig) MigrationWorkerV2 {
	mappings := config.GetSchemaMappingsFor(srcApp, dstApp)
	if mappings == nil {
		log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcApp, dstApp))
	}
	dstAppConfig.QR.Migration = true
	srcAppConfig.QR.Migration = true
	mWorker := MigrationWorkerV2{
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
		FTPClient:    GetFTPClient(),
		visitedNodes: make(map[string]map[string]bool)}
	if err := mWorker.FetchRoot(); err != nil {
		log.Fatal(err)
	}
	return mWorker
}

func CreateMigrationWorkerV2(uid, srcApp, srcAppID, dstApp, dstAppID string, logTxn *transaction.Log_txn, mtype string, mappings *config.MappedApp) MigrationWorkerV2 {
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
	mWorker := MigrationWorkerV2{
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
		FTPClient:    GetFTPClient(),
		visitedNodes: make(map[string]map[string]bool)}
	if err := mWorker.FetchRoot(); err != nil {
		log.Fatal(err)
	}
	return mWorker
}

func (self *MigrationWorkerV2) RenewDBConn() {
	if self.SrcDBConn != nil {
		self.SrcDBConn.Close()
	}
	if self.DstDBConn != nil {
		self.DstDBConn.Close()
	}
	self.SrcDBConn = db.GetDBConn(self.SrcAppConfig.AppName)
	self.DstDBConn = db.GetDBConn2(self.DstAppConfig.AppName)
}

func (self *MigrationWorkerV2) Finish() {
	self.SrcDBConn.Close()
	self.DstDBConn.Close()
}

func (self *MigrationWorkerV2) GetRoot() *DependencyNode {
	return self.root
}

func (self *MigrationWorkerV2) MType() string {
	return self.mtype
}

func (self *MigrationWorkerV2) UserID() string {
	return self.uid
}

func (self *MigrationWorkerV2) MigrationID() int {
	return self.logTxn.Txn_id
}

func (self *MigrationWorkerV2) ResolveParentDependencyConditions(node *DependencyNode, dconditions []config.DCondition, parentTag config.Tag) string {

	conditionStr := ""
	for _, condition := range dconditions {
		tagAttr, err := node.Tag.ResolveTagAttr(condition.TagAttr)
		if err != nil {
			log.Println(err, node.Tag.Name, condition.TagAttr)
			log.Fatal("@ResolveParentDependencyConditions: tagAttr in condition doesn't exist? ", condition.TagAttr)
			break
		}
		if len(condition.Restrictions) > 0 {
			restricted := false
			for _, restriction := range condition.Restrictions {
				if restrictionAttr, err := node.Tag.ResolveTagAttr(restriction["col"]); err == nil {
					if val, ok := node.Data[restrictionAttr]; ok {
						if strings.EqualFold(fmt.Sprint(val), restriction["val"]) {
							restricted = true
						}
					} else {
						fmt.Println(node.Data)
						log.Fatal("@ResolveParentDependencyConditions:", tagAttr, " doesn't exist in node data? ", node.Tag.Name)
					}
				} else {
					log.Fatal("@ResolveParentDependencyConditions: Col in restrictions doesn't exist? ", restriction["col"])
					break
				}
			}
			if restricted {
				return ""
			}
		}
		depOnAttr, err := parentTag.ResolveTagAttr(condition.DependsOnAttr)
		if err != nil {
			log.Println(err, parentTag.Name, condition.DependsOnAttr)
			log.Fatal("@ResolveParentDependencyConditions: depOnAttr in condition doesn't exist? ", condition.DependsOnAttr)
			break
		}
		if val, ok := node.Data[tagAttr]; ok {
			if conditionStr != "" {
				conditionStr += " AND "
			}
			conditionStr += fmt.Sprintf("%s = '%v'", depOnAttr, val)
		} else {
			fmt.Println(node.Data)
			log.Fatal("ResolveDependencyConditions:", tagAttr, " doesn't exist in node data? ", node.Tag.Name)
		}
	}
	return conditionStr
}

func (self *MigrationWorkerV2) ResolveDependencyConditions(node *DependencyNode, dep config.Dependency, tag config.Tag) string {

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

func (self *MigrationWorkerV2) ResolveOwnershipConditions(own config.Ownership, tag config.Tag) string {

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

func (self *MigrationWorkerV2) ExcludeVisited(tag config.Tag) string {
	visited := ""
	for _, tagMember := range tag.Members {
		if memberIDs, ok := self.visitedNodes[tagMember]; ok {
			pks := ""
			for pk := range memberIDs {
				pks += pk + ","
			}
			pks = strings.Trim(pks, ",")
			visited += fmt.Sprintf(" AND %s.id NOT IN (%s) ", tagMember, pks)
		}
	}
	return visited
}

func (self *MigrationWorkerV2) ResolveRestrictions(tag config.Tag) string {
	restrictions := ""
	for _, restriction := range tag.Restrictions {
		if restrictionAttr, err := tag.ResolveTagAttr(restriction["col"]); err == nil {
			restrictions += fmt.Sprintf(" AND %s = '%s' ", restrictionAttr, restriction["val"])
		}

	}
	return restrictions
}

func (self *MigrationWorkerV2) GetTagQL(tag config.Tag) string {

	sql := "SELECT %s FROM %s "

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
				}
			}
			seenMap[fromTable] = true
		}
		sql = fmt.Sprintf(sql, strings.Trim(cols, ","), joinStr)
	} else {
		table := tag.Members["member1"]
		_, cols := db.GetColumnsForTable(self.SrcAppConfig.DBConn, table)
		sql = fmt.Sprintf(sql, cols, table)
	}
	return sql
}

func (self *MigrationWorkerV2) FetchRoot() error {
	tagName := "root"
	if root, err := self.SrcAppConfig.GetTag(tagName); err == nil {
		rootTable, rootCol := self.SrcAppConfig.GetItemsFromKey(root, "root_id")
		where := fmt.Sprintf("%s.%s = '%s'", rootTable, rootCol, self.uid)
		ql := self.GetTagQL(root)
		sql := fmt.Sprintf("%s WHERE %s ", ql, where)
		sql += self.ResolveRestrictions(root)
		if data, err := db.DataCall1(self.SrcDBConn, sql); err == nil && len(data) > 0 {
			self.root = &DependencyNode{Tag: root, SQL: sql, Data: data}
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
	return nil
}

func (self *MigrationWorkerV2) GetAllNextNodes(node *DependencyNode) ([]*DependencyNode, error) {
	var nodes []*DependencyNode
	for _, dep := range self.SrcAppConfig.GetSubDependencies(node.Tag.Name) {
		if child, err := self.SrcAppConfig.GetTag(dep.Tag); err == nil {
			where := self.ResolveDependencyConditions(node, dep, child)
			ql := self.GetTagQL(child)
			sql := fmt.Sprintf("%s WHERE %s ", ql, where)
			sql += self.ResolveRestrictions(child)
			log.Fatal("@GetAllNextNodes | ", sql)
			if data, err := db.DataCall(self.SrcDBConn, sql); err == nil {
				for _, datum := range data {
					newNode := new(DependencyNode)
					newNode.Tag = child
					newNode.SQL = sql
					newNode.Data = datum
					nodes = append(nodes, newNode)
				}
			} else {
				log.Fatal("@GetAllNextNodes: Error while DataCall: ", err)
				return nil, err
			}
		} else {
			log.Fatal("@GetAllNextNodes: Tag doesn't exist? ", dep.Tag)
		}
	}
	return nodes, nil
}

func (self *MigrationWorkerV2) GetAllPreviousNodes(node *DependencyNode) ([]*DependencyNode, error) {
	var nodes []*DependencyNode
	for _, dep := range self.SrcAppConfig.GetParentDependencies(node.Tag.Name) {
		for _, pdep := range dep.DependsOn {
			if parent, err := self.SrcAppConfig.GetTag(pdep.Tag); err == nil {
				where := self.ResolveParentDependencyConditions(node, pdep.Conditions, parent)
				ql := self.GetTagQL(parent)
				sql := fmt.Sprintf("%s WHERE %s ", ql, where)
				sql += self.ResolveRestrictions(parent)
				// fmt.Println(node.SQL)
				// log.Fatal("@GetAllPreviousNodes | ", sql)
				if data, err := db.DataCall(self.SrcDBConn, sql); err == nil {
					for _, datum := range data {
						newNode := new(DependencyNode)
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

func (self *MigrationWorkerV2) GetAdjNode(node *DependencyNode, threadID int) (*DependencyNode, error) {
	if strings.EqualFold(node.Tag.Name, "root") {
		return self.GetOwnedNode(threadID)
	}
	return self.GetDependentNode(node, threadID)
}

func (self *MigrationWorkerV2) GetDependentNode(node *DependencyNode, threadID int) (*DependencyNode, error) {

	for _, dep := range self.SrcAppConfig.ShuffleDependencies(self.SrcAppConfig.GetSubDependencies(node.Tag.Name)) {
		if child, err := self.SrcAppConfig.GetTag(dep.Tag); err == nil {
			log.Println(fmt.Sprintf("x%2dx | FETCHING  tag  { %s } ", threadID, dep.Tag))
			where := self.ResolveDependencyConditions(node, dep, child)
			ql := self.GetTagQL(child)
			sql := fmt.Sprintf("%s WHERE %s ", ql, where)
			sql += self.ResolveRestrictions(child)
			sql += self.ExcludeVisited(child)
			sql += " ORDER BY random()"

			if data, err := db.DataCall1(self.SrcDBConn, sql); err == nil {
				if len(data) > 0 {
					newNode := DependencyNode{Tag: child, SQL: sql, Data: data}
					if !self.wList.IsAlreadyWaiting(newNode) && !self.IsVisited(&newNode) {
						return &newNode, nil
					}
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

func (self *MigrationWorkerV2) GetOwnedNode(threadID int) (*DependencyNode, error) {

	for _, own := range self.SrcAppConfig.GetShuffledOwnerships() {
		log.Println(fmt.Sprintf("x%2dx |         FETCHING  tag  { %s } ", threadID, own.Tag))
		// if self.unmappedTags.Exists(own.Tag) {
		// 	log.Println(fmt.Sprintf("x%2dx |         UNMAPPED  tag  { %s } ", threadID, own.Tag))
		// 	continue
		// }
		if child, err := self.SrcAppConfig.GetTag(own.Tag); err == nil {
			where := self.ResolveOwnershipConditions(own, child)
			ql := self.GetTagQL(child)
			sql := fmt.Sprintf("%s WHERE %s ", ql, where)
			sql += self.ResolveRestrictions(child)
			sql += " ORDER BY random() "
			// log.Fatal(sql)
			if data, err := db.DataCall1(self.SrcDBConn, sql); err == nil {
				if len(data) > 0 {
					newNode := DependencyNode{Tag: child, SQL: sql, Data: data}
					if !self.wList.IsAlreadyWaiting(newNode) {
						return &newNode, nil
					}
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

func (self *MigrationWorkerV2) PushData(tx *sql.Tx, dtable config.ToTable, pk, orgCols, cols string, undoAction *transaction.UndoAction, node *DependencyNode) error {

	undoActionSerialized, _ := json.Marshal(undoAction)
	transaction.LogChange(string(undoActionSerialized), self.logTxn)
	if err := app_display.GenDisplayFlag(self.logTxn.DBconn, self.DstAppConfig.AppID, dtable.TableID, pk, fmt.Sprint(self.logTxn.Txn_id)); err != nil {
		log.Fatal("## DISPLAY ERROR!", err)
		return errors.New("0")
	}

	for _, fromTable := range undoAction.OrgTables {
		if _, ok := node.Data[fmt.Sprintf("%s.id", fromTable)]; ok {
			srcID := fmt.Sprint(node.Data[fmt.Sprintf("%s.id", fromTable)])
			if fromTableID, err := db.TableID(self.logTxn.DBconn, fromTable, self.SrcAppConfig.AppID); err == nil {
				if err := db.InsertIntoIdentityTable(tx, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTableID, dtable.TableID, srcID, pk, fmt.Sprint(self.logTxn.Txn_id)); err != nil {
					log.Println("@PushData:db.InsertIntoIdentityTable: ", self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTableID, dtable.TableID, srcID, pk, fmt.Sprint(self.logTxn.Txn_id))
					log.Fatal(err)
					return errors.New("0")
				}
				if serr := db.SaveForLEvaluation(tx, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTable, dtable.TableID, srcID, pk, orgCols, cols, fmt.Sprint(self.logTxn.Txn_id)); serr != nil {
					log.Println("@PushData:db.SaveForLEvaluation: ", self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTable, dtable.TableID, srcID, pk, orgCols, cols, fmt.Sprint(self.logTxn.Txn_id))
					log.Fatal(serr)
					return errors.New("0")
				}
			} else {
				log.Println("@PushData:db.TableID: ", fromTable, self.SrcAppConfig.AppID)
				log.Fatal(err)
			}
		}
	}
	return nil
}

func (self *MigrationWorkerV2) CheckMappingConditions(toTable config.ToTable, node *DependencyNode) bool {
	breakCondition := false
	if len(toTable.Conditions) > 0 {
		for conditionKey, conditionVal := range toTable.Conditions {
			if nodeVal, ok := node.Data[conditionKey]; ok {
				if conditionVal[:1] == "#" {
					// fmt.Println("VerifyMappingConditions: conditionVal[:1] == #")
					// fmt.Println(conditionKey, conditionVal, nodeVal)
					// fmt.Scanln()
					switch conditionVal {
					case "#NULL":
						{
							if nodeVal != nil {
								// log.Println(nodeVal, "!=", conditionVal)
								// fmt.Println(conditionKey, conditionVal, nodeVal)
								// log.Fatal("@VerifyMappingConditions: return false, from case #NULL:")
								return false
							}
						}
					case "#NOTNULL":
						{
							if nodeVal == nil {
								// log.Println(nodeVal, "!=", conditionVal)
								// fmt.Println(conditionKey, conditionVal, nodeVal)
								// log.Fatal("@VerifyMappingConditions: return false, from case #NOTNULL:")
								return false
							}
						}
					default:
						{
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
					} else {
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
				log.Fatal("@CheckMappingConditions: stop here and check")
			}
		}
	}
	return breakCondition
}

func (self *MigrationWorkerV2) GetNodeOwner(node *DependencyNode) (string, bool) {

	if ownership := self.SrcAppConfig.GetOwnership(node.Tag.Name, self.root.Tag.Name); ownership != nil {
		for _, condition := range ownership.Conditions {
			tagAttr, err := node.Tag.ResolveTagAttr(condition.TagAttr)
			if err != nil {
				log.Fatal("Resolving TagAttr in GetNodeOwner", err, node.Tag.Name, condition.TagAttr)
				break
			}
			depOnAttr, err := self.root.Tag.ResolveTagAttr(condition.DependsOnAttr)
			if err != nil {
				log.Fatal("Resolving depOnAttr in GetNodeOwner", err, node.Tag.Name, condition.DependsOnAttr)
				break
			}
			if nodeVal, err := node.GetValueForKey(tagAttr); err == nil {
				if rootVal, err := self.root.GetValueForKey(depOnAttr); err == nil {
					if !strings.EqualFold(nodeVal, rootVal) {
						return nodeVal, true
					} else {
						return nodeVal, false
					}
				} else {
					fmt.Println("Ownership Condition Key in Root Data:", depOnAttr, "doesn't exist!")
					fmt.Println("root data:", self.root.Data)
					log.Fatal("@GetNodeOwner: stop here and check ownership conditions wrt root")
				}
			} else {
				fmt.Println("Ownership Condition Key", tagAttr, "doesn't exist!")
				fmt.Println("node data:", node.Data)
				fmt.Println("node sql:", node.SQL)
				log.Fatal("@GetNodeOwner: stop here and check ownership conditions")
			}
		}
	} else {
		log.Fatal("Ownership not found in GetNodeOwner:", node.Tag.Name)
	}
	return "", false
}

func (self *MappedData) UpdateData(col, orgCol, fromTable string, ival interface{}) {
	self.ivals = append(self.ivals, ival)
	self.vals += fmt.Sprintf("$%d,", len(self.ivals))
	self.cols += fmt.Sprintf("%s,", col)
	self.orgCols += fmt.Sprintf("%s,", orgCol)
	if fromTable != "" {
		self.srcTables[fromTable] = true
	}
}

func (self *MappedData) Trim(chars string) {
	self.vals = strings.Trim(self.vals, chars)
	self.cols = strings.Trim(self.cols, chars)
	self.orgCols = strings.Trim(self.orgCols, chars)
}

func (self *MigrationWorkerV2) FetchFromMapping(node *DependencyNode, toAttr, assignedTabCol string, data *MappedData) error {
	args := strings.Split(assignedTabCol, ",")
	for i, arg := range args {
		args[i] = strings.Trim(arg, "()")
	}
	if nodeVal, ok := node.Data[args[2]]; ok {
		targetTabCol := strings.Split(args[0], ".")
		comparisonTabCol := strings.Split(args[1], ".")
		if res, err := db.FetchForMapping(self.SrcAppConfig.DBConn, targetTabCol[0], targetTabCol[1], comparisonTabCol[1], fmt.Sprint(nodeVal)); err != nil {
			fmt.Println(targetTabCol[0], targetTabCol[1], comparisonTabCol[1], fmt.Sprint(nodeVal))
			log.Fatal("@GetMappedData: FetchForMapping | ", err)
			return err
		} else {
			data.UpdateData(toAttr, assignedTabCol, targetTabCol[0], res[targetTabCol[1]])
			if len(args) > 3 {
				toMemberTokens := strings.Split(args[3], ".")
				data.refs = append(data.refs, MappingRef{
					fromID:     fmt.Sprint(res[targetTabCol[1]]),
					fromMember: targetTabCol[0],
					fromAttr:   targetTabCol[1],
					toID:       fmt.Sprint(res[targetTabCol[1]]),
					toMember:   toMemberTokens[0],
					toAttr:     toMemberTokens[1]})
			}
		}
	} else {
		fmt.Println(node.Tag.Name, node.Data)
		log.Fatal("@GetMappedData: unable to fetch ", args[2])
		return errors.New("Unable to fetch data from node")
	}
	return nil
}

func (self *MigrationWorkerV2) GetMappedData(toTable config.ToTable, node *DependencyNode) (MappedData, error) {

	data := MappedData{
		cols:        "",
		vals:        "",
		orgCols:     "",
		orgColsLeft: "",
		srcTables:   make(map[string]bool),
		undoAction:  new(transaction.UndoAction)}

	for toAttr, fromAttr := range toTable.Mapping {
		if val, ok := node.Data[fromAttr]; ok {
			fromTokens := strings.Split(fromAttr, ".")
			data.UpdateData(toAttr, fromTokens[1], fromTokens[0], val)
			data.undoAction.AddData(fromAttr, val)
			data.undoAction.AddOrgTable(fromTokens[0])
		} else if strings.Contains(fromAttr, "$") {
			if inputVal, err := self.mappings.GetInput(fromAttr); err == nil {
				data.UpdateData(toAttr, fromAttr, "", inputVal)
			}
		} else if strings.Contains(fromAttr, "#") {
			assignedTabCol := strings.Trim(fromAttr, "(#ASSIGN#FETCH#REF)")
			if strings.Contains(fromAttr, "#ASSIGN") {
				if nodeVal, ok := node.Data[assignedTabCol]; ok {
					assignedTabColTokens := strings.Split(assignedTabCol, ".")
					data.UpdateData(toAttr, assignedTabColTokens[1], assignedTabColTokens[0], nodeVal)
				}
			} else if strings.Contains(fromAttr, "#REF") {
				if strings.Contains(fromAttr, "#FETCH") {
					self.FetchFromMapping(node, toAttr, assignedTabCol, &data)
				} else {
					args := strings.Split(assignedTabCol, ",")
					if nodeVal, ok := node.Data[args[0]]; ok {
						data.UpdateData(toAttr, assignedTabCol, "", nodeVal)
					}
					var toID, fromID string

					if val, ok := node.Data[args[0]]; ok {
						toID = fmt.Sprint(val)
					} else {
						fmt.Println(args[0], " | ", args)
						fmt.Println(node.Data)
						log.Fatal("@GetMappedData > #REF > toID: Unable to find ref value in node data")
						return data, errors.New("Unable to find ref value in node data")
					}

					firstMemberTokens := strings.Split(args[0], ".")
					secondMemberTokens := strings.Split(args[1], ".")

					if val, ok := node.Data[firstMemberTokens[0]+".id"]; ok {
						fromID = fmt.Sprint(val)
					} else {
						fmt.Println(args[0], " | ", args)
						fmt.Println(node.Data)
						log.Fatal("@GetMappedData > #REF > fromID: Unable to find ref value in node data")
						return data, errors.New("Unable to find ref value in node data")
					}

					data.refs = append(data.refs, MappingRef{fromID: fromID, fromMember: firstMemberTokens[0], fromAttr: fromAttr, toID: toID, toAttr: secondMemberTokens[1], toMember: secondMemberTokens[0]})
				}
			} else if strings.Contains(fromAttr, "#FETCH") {
				// #FETCH(targetSrcTable.targetSrcCol, targetSrcTable.srcColToCompare, currentSrcTable.currentSrcColForComparison)
				// # Do we need to create an identity entry for row referenced in fetch?
				if err := self.FetchFromMapping(node, toAttr, assignedTabCol, &data); err != nil {
					fmt.Println(node.Data)
					fmt.Println(toAttr, assignedTabCol)
					log.Fatal("@GetMappedData > #FETCH > FetchFromMapping: Unable to fetch")
					return data, err
				}
				// fmt.Println(strings.Trim(cols, ","), strings.Trim(vals, ","), ivals, strings.Trim(orgCols, ","), orgColsLeft)
				// log.Fatal("check")
			} else {
				switch fromAttr {
				case "#GUID":
					{
						data.UpdateData(toAttr, assignedTabCol, "", uuid.New())
					}
				case "#RANDINT":
					{
						data.UpdateData(toAttr, assignedTabCol, "", self.SrcAppConfig.QR.NewRowId())
					}
				default:
					{
						fmt.Println(toTable.Table, toAttr, fromAttr)
						log.Fatal("@GetMappedData: Case not found:" + fromAttr)
					}
				}
			}
			// log.Fatal(fromAttr)
		} else {
			data.orgColsLeft += fmt.Sprintf("%s,", strings.Split(fromAttr, ".")[1])
		}
	}
	// fmt.Println(strings.Trim(cols, ","), strings.Trim(vals, ","), ivals, strings.Trim(orgCols, ","), orgColsLeft, undoAction)
	data.Trim(",")
	return data, nil
}

func (self *MigrationWorkerV2) DeleteRow(node *DependencyNode) error {
	for _, tagMember := range node.Tag.Members {
		idCol := fmt.Sprintf("%s.id", tagMember)
		if _, ok := node.Data[idCol]; ok {
			srcID := fmt.Sprint(node.Data[idCol])
			if derr := db.ReallyDeleteRowFromAppDB(self.tx.SrcTx, tagMember, srcID); derr != nil {
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

func (self *MigrationWorkerV2) TransferMedia(filePath string) error {

	file, err := os.Open(filePath)
	if err != nil {
		log.Println(fmt.Sprintf("Can't open the file at: %s | ", filePath), err)
		return err
	}

	fpTokens := strings.Split(filePath, "/")
	fileName := fpTokens[len(fpTokens)-1]
	fsName := "/" + fileName

	log.Println(fmt.Sprintf("Transferring file [%s] with name [%s] to [%s]...", filePath, fileName, fsName))
	if err := self.FTPClient.Stor(fsName, file); err != nil {
		log.Println("File Transfer Failed: ", err)
		return err
	}

	return nil
}

func (self *MigrationWorkerV2) HandleUnmappedMembersOfNode(mapping config.Mapping, node *DependencyNode) error {

	if self.mtype != DELETION {
		return nil
	}
	for _, nodeMember := range node.Tag.GetTagMembers() {
		if !helper.Contains(mapping.FromTables, nodeMember) {
			if err := self.SendMemberToBag(node, nodeMember, self.uid, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func (self *MigrationWorkerV2) MigrateNode(mapping config.Mapping, node *DependencyNode) error {

	// Handle partial bags

	for _, toTable := range mapping.ToTables {
		if self.CheckMappingConditions(toTable, node) {
			continue
		}
		if mappedData, _ := self.GetMappedData(toTable, node); len(mappedData.cols) > 0 && len(mappedData.vals) > 0 && len(mappedData.ivals) > 0 {
			mappedData.undoAction.AddDstTable(toTable.Table)
			// if strings.Contains(toTable.Table, "status_"){
			// 	fmt.Println(toTable.Table, cols, placeholders, ivals)
			// 	log.Fatal("--------------")
			// }
			if id, err := db.InsertRowIntoAppDB(self.tx.DstTx, toTable.Table, mappedData.cols, mappedData.vals, mappedData.ivals...); err == nil {
				if toTableID, err := db.TableID(self.logTxn.DBconn, toTable.Table, self.DstAppConfig.AppID); err != nil {
					for fromTable := range mappedData.srcTables {
						if fromTableID, err := db.TableID(self.logTxn.DBconn, fromTable, self.SrcAppConfig.AppID); err != nil {
							fromID := fmt.Sprint(node.Data[fromTable+".id"])
							if err := db.InsertIntoIdentityTable(self.tx.StencilTx, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTableID, toTableID, fromID, fmt.Sprint(id), fmt.Sprint(self.logTxn.Txn_id)); err != nil {
								fmt.Println("@ERROR_PushData")
								fmt.Println("@Params:", toTable.Table, fmt.Sprint(id), mappedData.orgCols, mappedData.cols, mappedData.undoAction, node)
								log.Fatal(err)
								return err
							}
						}
					}
					if err := self.PushData(self.tx.StencilTx, toTable, fmt.Sprint(id), mappedData.orgCols, mappedData.cols, mappedData.undoAction, node); err != nil {
						fmt.Println("@ERROR_PushData")
						fmt.Println("@Params:", toTable.Table, fmt.Sprint(id), mappedData.orgCols, mappedData.cols, mappedData.undoAction, node)
						log.Fatal(err)
						return err
					}
					if len(toTable.Media) > 0 {
						if filePathCol, ok := toTable.Media["path"]; ok {
							if filePath, ok := node.Data[filePathCol]; ok {
								if err := self.TransferMedia(fmt.Sprint(filePath)); err != nil {
									log.Fatal("@MigrateNode: ", err)
								}
							}
						} else {
							log.Fatal("@MigrateNode > toTable.Media: Path not found in map!")
						}
					}
				}
			} else {
				log.Fatal("@MigrateNode:", err)
				return err
			}

			if err := self.AddMappedReferences(mappedData.refs); err != nil {
				log.Fatal("@MigrateNode > AddMappedReferences: ", err)
				return err
			}
		} else {
			log.Fatal("@MigrateNode > GetMappedData:", mappedData)
		}
	}

	if self.mtype == DELETION {
		if err := self.HandleUnmappedMembersOfNode(mapping, node); err != nil {
			return err
		}
		if err := self.DeleteRow(node); err != nil {
			return err
		}
	}

	return nil
}

func (self *MigrationWorkerV2) HandleWaitingList(appMapping config.Mapping, tagMembers []string, node *DependencyNode) (*DependencyNode, error) {

	srctx, err := self.SrcDBConn.Begin()
	if err != nil {
		log.Println("Can't create HandleWaitingList transaction!")
		return nil, errors.New("0")
	}
	if err := self.DeleteRow(node); err != nil {
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

func (self *MigrationWorkerV2) HandleUnmappedTags(node *DependencyNode) error {
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

func (self *MigrationWorkerV2) HandleUnmappedNode(node *DependencyNode) error {
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

func (self *MigrationWorkerV2) FetchMappingsForNode(node *DependencyNode) (config.Mapping, bool) {
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

func (self *MigrationWorkerV2) SendMemberToBag(node *DependencyNode, member, ownerID string, fromNode bool) error {
	if memberID, err := db.TableID(self.logTxn.DBconn, member, self.SrcAppConfig.AppID); err != nil {
		log.Fatal("@SendMemberToBag > TableID: error in getting table id for member! ", member, err)
		return err
	} else {
		var bagData map[string]interface{}
		for col, val := range node.Data {
			colTokens := strings.Split(col, ".")
			colMember := colTokens[0]
			colAttr := colTokens[1]
			if strings.Contains(colMember, member) {
				bagData[colAttr] = val
			}
		}
		if len(bagData) > 0 {
			if id, ok := node.Data[member+".id"]; ok {
				if jsonData, err := json.Marshal(bagData); err == nil {
					if err := db.CreateNewBag(self.tx.StencilTx, self.SrcAppConfig.AppID, memberID, fmt.Sprint(id), ownerID, fmt.Sprint(self.logTxn.Txn_id), jsonData); err != nil {
						log.Fatal("@SendMemberToBag: error in creating bag! ", err)
						return err
					}
				} else {
					fmt.Println(bagData)
					log.Fatal("@SendMemberToBag: unable to convert bag data to JSON ", err)
					return err
				}
			} else {
				fmt.Println(node.Data)
				log.Fatal("@SendMemberToBag: member doesn't contain id! ", member)
				return err
			}
			if !fromNode {
				if err := self.AddInnerReferences(node, member); err != nil {
					fmt.Println(node.Tag.Members)
					fmt.Println(node.Tag.InnerDependencies)
					fmt.Println(node.Data)
					log.Fatal("@SendMemberToBag > AddInnerReferences: Adding Inner References failed ", err)
					return err
				}
			}
		}
	}

	return nil
}

func (self *MigrationWorkerV2) SendNodeToBag(node *DependencyNode) error {
	if ownerID, _ := self.GetNodeOwner(node); len(ownerID) > 0 {
		for _, member := range node.Tag.Members {
			if err := self.SendMemberToBag(node, member, ownerID, true); err != nil {
				fmt.Println(node)
				log.Fatal("@SendNodeToBag > SendMemberToBag: ownerID error! ")
			}
		}
		if err := self.AddInnerReferences(node, ""); err != nil {
			fmt.Println(node.Tag.Members)
			fmt.Println(node.Tag.InnerDependencies)
			fmt.Println(node.Data)
			log.Fatal("@SendNodeToBag > AddInnerReferences: Adding Inner References failed ", err)
			return err
		}
	} else {
		fmt.Println(node)
		log.Fatal("@SendNodeToBag > GetNodeOwner: ownerID error! ")
	}

	return nil
}

func (self *MigrationWorkerV2) HandleMigration(node *DependencyNode, isBag bool) error {

	if mapping, found := self.FetchMappingsForNode(node); found {
		tagMembers := node.Tag.GetTagMembers()
		if helper.Sublist(tagMembers, mapping.FromTables) { // other mappings HANDLE!
			return self.MigrateNode(mapping, node)
		}
		if wNode, err := self.HandleWaitingList(mapping, tagMembers, node); wNode != nil && err == nil {
			return self.MigrateNode(mapping, wNode)
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

func (self *MigrationWorkerV2) HandleLeftOverWaitingNodes() {

	for _, waitingNode := range self.wList.Nodes {
		for _, containedNode := range waitingNode.ContainedNodes {
			self.HandleUnmappedNode(containedNode)
		}
	}
}

func (self *MigrationWorkerV2) IsVisited(node *DependencyNode) bool {

	for _, tagMember := range node.Tag.Members {
		if _, ok := self.visitedNodes[tagMember]; !ok {
			continue
		}
		idCol := fmt.Sprintf("%s.id", tagMember)
		if _, ok := node.Data[idCol]; ok {
			srcID := fmt.Sprint(node.Data[idCol])
			if _, ok := self.visitedNodes[tagMember][srcID]; ok {
				return true
			}
		} else {
			log.Println("In: IsVisited | node.Data =>", node.Data)
			log.Fatal(idCol, "NOT PRESENT IN NODE DATA")
		}
	}
	return false
}

func (self *MigrationWorkerV2) MarkAsVisited(node *DependencyNode) {
	for _, tagMember := range node.Tag.Members {
		idCol := fmt.Sprintf("%s.id", tagMember)
		if _, ok := node.Data[idCol]; ok {
			if _, ok := self.visitedNodes[tagMember]; !ok {
				self.visitedNodes[tagMember] = make(map[string]bool)
			}
			srcID := fmt.Sprint(node.Data[idCol])
			self.visitedNodes[tagMember][srcID] = true
		} else {
			log.Println("In: MarkAsVisited | node.Data =>", node.Data)
			log.Fatal(idCol, "NOT PRESENT IN NODE DATA")
		}
	}
}

func (self *MigrationWorkerV2) CheckNextNode(node *DependencyNode) error {
	if nextNodes, err := self.GetAllNextNodes(node); err == nil {
		for _, nextNode := range nextNodes {
			self.AddToReferences(nextNode, node)
			if precedingNodes, err := self.GetAllPreviousNodes(node); err != nil {
				return err
			} else if len(precedingNodes) <= 1 {
				if err := self.CheckNextNode(nextNode); err != nil {
					return err
				}
				if err := self.SendNodeToBag(nextNode); err != nil {
					return err
				}
			}
		}
		return nil
	} else {
		return err
	}
}

func (self *MigrationWorkerV2) AddMappedReferences(refs []MappingRef) error {

	for _, ref := range refs {

		dependeeMemberID, err := db.TableID(self.logTxn.DBconn, ref.fromMember, self.SrcAppConfig.AppID)
		if err != nil {
			log.Fatal("@AddMappedReferences: Unable to resolve id for dependeeMember ", ref.fromMember)
			return err
		}

		depOnMemberID, err := db.TableID(self.logTxn.DBconn, ref.toMember, self.SrcAppConfig.AppID)
		if err != nil {
			log.Fatal("@AddMappedReferences: Unable to resolve id for depOnMember ", ref.toMember)
			return err
		}

		if err := db.CreateNewReference(self.tx.StencilTx, self.SrcAppConfig.AppID, dependeeMemberID, ref.fromID, depOnMemberID, ref.toID, fmt.Sprint(self.logTxn.Txn_id), ref.fromAttr, ref.toAttr); err != nil {
			fmt.Println("#Args: ", self.SrcAppConfig.AppID, dependeeMemberID, ref.fromID, depOnMemberID, ref.toID, fmt.Sprint(self.logTxn.Txn_id), ref.fromAttr, ref.toAttr)
			log.Fatal("@AddMappedReferences: Unable to CreateNewReference: ", err)
			return err
		}
	}

	return nil
}

func (self *MigrationWorkerV2) AddInnerReferences(node *DependencyNode, member string) error {

	for _, innerDependency := range node.Tag.InnerDependencies {
		for dependee, dependsOn := range innerDependency {

			depTokens := strings.Split(dependee, ".")
			dependeeMember := node.Tag.Members[depTokens[0]]
			dependeeMemberID, err := db.TableID(self.logTxn.DBconn, dependeeMember, self.SrcAppConfig.AppID)
			if err != nil {
				log.Fatal("@AddInnerReferences: Unable to resolve id for dependeeMember ", dependeeMember)
			}

			depOnTokens := strings.Split(dependsOn, ".")
			depOnMember := node.Tag.Members[depOnTokens[0]]
			depOnMemberID, err := db.TableID(self.logTxn.DBconn, depOnMember, self.SrcAppConfig.AppID)
			if err != nil {
				log.Fatal("@AddInnerReferences: Unable to resolve id for depOnMember ", depOnMember)
			}

			if member != "" {
				if !strings.EqualFold(dependeeMember, member) && !strings.EqualFold(depOnMember, member) {
					continue
				}
			}

			var fromID, toID string

			if val, ok := node.Data[dependeeMember+".id"]; ok {
				fromID = fmt.Sprint(val)
			} else {
				fmt.Println(node.Data)
				log.Fatal("@AddInnerReferences:", dependeeMember+".id", " doesn't exist in node data? ", node.Tag.Name)
			}

			if val, ok := node.Data[depOnMember+".id"]; ok {
				toID = fmt.Sprint(val)
			} else {
				fmt.Println(node.Data)
				log.Fatal("@AddInnerReferences:", depOnMember+".id", " doesn't exist in node data? ", node.Tag.Name)
			}

			if err := db.CreateNewReference(self.tx.StencilTx, self.SrcAppConfig.AppID, dependeeMemberID, fromID, depOnMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), depTokens[1], depOnTokens[1]); err != nil {
				fmt.Println("#Args: ", self.SrcAppConfig.AppID, dependeeMemberID, fromID, depOnMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), depTokens[1], depOnTokens[1])
				log.Fatal("@AddInnerReferences: Unable to CreateNewReference: ", err)
				return err
			}
		}
	}

	return nil
}

func (self *MigrationWorkerV2) AddToReferences(currentNode *DependencyNode, referencedNode *DependencyNode) error {

	if dep, err := self.SrcAppConfig.CheckDependency(currentNode.Tag.Name, referencedNode.Tag.Name); err != nil {
		fmt.Println(err)
		log.Fatal("@AddToReferences: CheckDependency can't find dependency!")
	} else {
		for _, condition := range dep.Conditions {
			tagAttr, err := currentNode.Tag.ResolveTagAttr(condition.TagAttr)
			if err != nil {
				log.Println(err, currentNode.Tag.Name, condition.TagAttr)
				log.Fatal("@AddToReferences: tagAttr in condition doesn't exist? ", condition.TagAttr)
				break
			}
			tagAttrTokens := strings.Split(tagAttr, ".")
			fromMember := tagAttrTokens[0]
			fromReference := tagAttrTokens[1]
			fromMemberID, err := db.TableID(self.logTxn.DBconn, fromMember, self.SrcAppConfig.AppID)
			if err != nil {
				log.Fatal("@AddToReferences: Unable to resolve id for fromMember ", fromMember)
			}

			depOnAttr, err := referencedNode.Tag.ResolveTagAttr(condition.DependsOnAttr)
			if err != nil {
				log.Println(err, referencedNode.Tag.Name, condition.DependsOnAttr)
				log.Fatal("@AddToReferences: depOnAttr in condition doesn't exist? ", condition.DependsOnAttr)
				break
			}
			depOnAttrTokens := strings.Split(depOnAttr, ".")
			toMember := depOnAttrTokens[0]
			toReference := depOnAttrTokens[1]
			toMemberID, err := db.TableID(self.logTxn.DBconn, toMember, self.SrcAppConfig.AppID)
			if err != nil {
				log.Fatal("@AddToReferences: Unable to resolve id for toMember ", toMember)
			}

			var fromID, toID string

			if val, ok := currentNode.Data[fromMember+".id"]; ok {
				fromID = fmt.Sprint(val)
			} else {
				fmt.Println(currentNode.Data)
				log.Fatal("@AddToReferences:", fromMember+".id", " doesn't exist in node data? ", currentNode.Tag.Name)
			}

			if val, ok := referencedNode.Data[toMember+".id"]; ok {
				toID = fmt.Sprint(val)
			} else {
				fmt.Println(referencedNode.Data)
				log.Fatal("@AddToReferences:", toMember+".id", " doesn't exist in node data? ", referencedNode.Tag.Name)
			}

			if err := db.CreateNewReference(self.tx.StencilTx, self.SrcAppConfig.AppID, fromMemberID, fromID, toMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), fromReference, toReference); err != nil {
				fmt.Println("#Args: ", self.SrcAppConfig.AppID, fromMemberID, fromID, toMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), fromReference, toReference)
				log.Fatal("@AddToReferences: Unable to CreateNewReference: ", err)
				return err
			}
		}
	}
	return nil
}

func (self *MigrationWorkerV2) InitTransactions() error {
	var err error
	self.tx.SrcTx, err = self.SrcDBConn.Begin()
	if err != nil {
		log.Fatal("Error creating Source DB Transaction! ", err)
		return err
	}
	self.tx.DstTx, err = self.DstDBConn.Begin()
	if err != nil {
		log.Fatal("Error creating Dst DB Transaction! ", err)
		return err
	}
	self.tx.StencilTx, err = self.logTxn.DBconn.Begin()
	if err != nil {
		log.Fatal("Error creating Stencil DB Transaction! ", err)
		return err
	}
	return nil
}

func (self *MigrationWorkerV2) CommitTransactions() error {
	log.Fatal("@CommitTransactions: About to Commit!")
	if err := self.tx.SrcTx.Commit(); err != nil {
		log.Fatal("Error committing Source DB Transaction! ", err)
		return err
	}
	if err := self.tx.DstTx.Commit(); err != nil {
		log.Fatal("Error committing Destination DB Transaction! ", err)
		return err
	}
	if err := self.tx.StencilTx.Commit(); err != nil {
		log.Fatal("Error committing Stencil DB Transaction! ", err)
		return err
	}
	return nil
}

func (self *MigrationWorkerV2) DeletionMigration(node *DependencyNode, threadID int) error {

	for {
		if adjNode, err := self.GetAdjNode(node, threadID); err != nil {
			return err
		} else {
			if adjNode == nil {
				break
			}
			fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
			log.Println(fmt.Sprintf("~%2d~ Current   Node: { %s } ", threadID, node.Tag.Name))
			log.Println(fmt.Sprintf("~%2d~ Adjacent  Node: { %s } ", threadID, adjNode.Tag.Name))
			if err := self.DeletionMigration(adjNode, threadID); err != nil {
				log.Fatal(fmt.Sprintf("~%2d~ ERROR! NODE : { %s }, ADJ_NODE : { %s } | err: [ %s ]", threadID, node.Tag.Name, adjNode.Tag.Name, err))
				return err
			}
		}
	}

	log.Println(fmt.Sprintf("#%2d# Process   Node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))

	if ownerID, isRoot := self.GetNodeOwner(node); isRoot && len(ownerID) > 0 {
		if err := self.InitTransactions(); err != nil {
			return err
		} else {
			defer self.tx.SrcTx.Rollback()
			defer self.tx.DstTx.Rollback()
			defer self.tx.StencilTx.Rollback()
		}

		if err := self.CheckNextNode(node); err != nil {
			return err
		}

		if previousNodes, err := self.GetAllPreviousNodes(node); err == nil {
			for _, previousNode := range previousNodes {
				self.AddToReferences(node, previousNode)
			}
		} else {
			return err
		}

		if err := self.HandleMigration(node, false); err == nil {
			log.Println(fmt.Sprintf("x%2dx MIGRATED  node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
		} else {
			if strings.EqualFold(err.Error(), "2") {
				log.Println(fmt.Sprintf("x%2dx UNMAPPED  node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
			} else {
				log.Println(fmt.Sprintf("x%2dx FAILED    node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
				if strings.EqualFold(err.Error(), "0") {
					log.Println(err)
					return err
				}
				return err
				// if strings.Contains(err.Error(), "deadlock") {
				// 	return err
				// }
			}
		}
		if err := self.CommitTransactions(); err != nil {
			return err
		} else {
			log.Println(fmt.Sprintf("x%2dx COMMITTED node { %s } ", threadID, node.Tag.Name))
		}
	} else {
		log.Println(fmt.Sprintf("x%2dx VISITED   node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
		self.MarkAsVisited(node)
	}

	fmt.Println("------------------------------------------------------------------------")

	return nil
}

func (self *MigrationWorkerV2) ConsistentMigration(threadID int) error {

	for {
		if node, err := self.GetOwnedNode(threadID); err == nil {
			if node == nil {
				return nil
			}
			nodeIDAttr, _ := node.Tag.ResolveTagAttr("id")
			log.Println(fmt.Sprintf("~%2d~ | Current   Node: { %s } ID: %v", threadID, node.Tag.Name, node.Data[nodeIDAttr]))
			if err := self.HandleMigration(node, false); err == nil {
				log.Println(fmt.Sprintf("x%2dx | MIGRATED  node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
			} else {
				log.Println(fmt.Sprintf("x%2dx | RCVD ERR  node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName), err)
				if self.unmappedTags.Exists(node.Tag.Name) {
					log.Println(fmt.Sprintf("x%2dx | BREAKLOOP node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName), err)
					break
				}
				if strings.EqualFold(err.Error(), "2") {
					log.Println(fmt.Sprintf("x%2dx | IGNORED   node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
				} else if strings.EqualFold(err.Error(), "3") {
					log.Println(fmt.Sprintf("x%2dx | EXISTS    node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
				} else {
					log.Println(fmt.Sprintf("x%2dx | FAILED    node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName), err)
					if strings.EqualFold(err.Error(), "0") {
						log.Println(err)
						return err
					}
					if strings.Contains(err.Error(), "deadlock") {
						return err
					}
				}
			}
		} else {
			return err
		}
	}

	if err := self.HandleMigration(self.root, false); err == nil {
		log.Println(fmt.Sprintf("x%2dx | MIGRATED  node { %s } From [%s] to [%s]", threadID, self.root.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
	} else {
		log.Println(fmt.Sprintf("x%2dx | MIGRATED? node { %s } From [%s] to [%s]", threadID, self.root.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
	}
	return nil
}

func (self *MigrationWorkerV2) IndependentMigration(threadID int) error {

	return self.ConsistentMigration(threadID)
}

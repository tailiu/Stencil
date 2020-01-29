package migrate

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"stencil/SA1_display"
	"stencil/config"
	"stencil/db"
	"stencil/helper"
	"stencil/reference_resolution"
	"stencil/transaction"
	"strings"

	"github.com/google/uuid"
)

func CreateMigrationWorkerV2WithAppsConfig(uid string, logTxn *transaction.Log_txn, mtype string, srcAppConfig, dstAppConfig config.AppConfig, threadID int) MigrationWorkerV2 {
	mappings := config.GetSchemaMappingsFor(srcAppConfig.AppName, dstAppConfig.AppName)
	if mappings == nil {
		log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcAppConfig.AppName, dstAppConfig.AppName))
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
		SrcDBConn:    db.GetDBConn(srcAppConfig.AppName),
		DstDBConn:    db.GetDBConn(dstAppConfig.AppName),
		logTxn:       &transaction.Log_txn{DBconn: logTxn.DBconn, Txn_id: logTxn.Txn_id},
		mtype:        mtype,
		visitedNodes: make(map[string]map[string]bool)}
	if err := mWorker.FetchRoot(threadID); err != nil {
		log.Fatal(err)
	}
	mWorker.FTPClient = GetFTPClient()
	return mWorker
}

func CreateMigrationWorkerV2(uid, srcApp, srcAppID, dstApp, dstAppID string, logTxn *transaction.Log_txn, mtype string, mappings *config.MappedApp, threadID int) MigrationWorkerV2 {
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
		DstDBConn:    db.GetDBConn(dstApp),
		logTxn:       &transaction.Log_txn{DBconn: logTxn.DBconn, Txn_id: logTxn.Txn_id},
		mtype:        mtype,
		visitedNodes: make(map[string]map[string]bool)}

	if err := mWorker.FetchRoot(threadID); err != nil {
		log.Fatal(err)
	}
	mWorker.FTPClient = GetFTPClient()
	log.Println("Worker Created for thread: ", threadID)
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
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

func (node *DependencyNode) ResolveParentDependencyConditions(dconditions []config.DCondition, parentTag config.Tag) string {

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

func (node *DependencyNode) ResolveDependencyConditions(SrcAppConfig config.AppConfig, dep config.Dependency, tag config.Tag) string {

	where := ""
	for _, depOn := range dep.DependsOn {
		if depOnTag, err := SrcAppConfig.GetTag(depOn.Tag); err == nil {
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

func (root *DependencyNode) ResolveOwnershipConditions(own config.Ownership, tag config.Tag) string {

	where := ""
	for _, condition := range own.Conditions {
		conditionStr := ""
		tagAttr, err := tag.ResolveTagAttr(condition.TagAttr)
		if err != nil {
			fmt.Println("data1", root.Data)
			log.Fatal(err, tag.Name, condition.TagAttr)
			break
		}
		depOnAttr, err := root.Tag.ResolveTagAttr(condition.DependsOnAttr)
		if err != nil {
			fmt.Println("data2", root.Data)
			log.Fatal(err, tag.Name, condition.DependsOnAttr)
			break
		}
		if _, ok := root.Data[depOnAttr]; ok {
			if conditionStr != "" || where != "" {
				conditionStr += " AND "
			}
			conditionStr += fmt.Sprintf("%s = '%v'", tagAttr, root.Data[depOnAttr])
		} else {
			fmt.Println("data3", root.Data)
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
					joinStr += fmt.Sprintf(" FULL JOIN %s ON %s ", toTable, strings.Join(conditions, " AND "))
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

func (self *MigrationWorkerV2) FetchRoot(threadID int) error {
	tagName := "root"
	if root, err := self.SrcAppConfig.GetTag(tagName); err == nil {
		rootTable, rootCol := self.SrcAppConfig.GetItemsFromKey(root, "root_id")
		where := fmt.Sprintf("%s.%s = '%s'", rootTable, rootCol, self.uid)
		ql := self.GetTagQL(root)
		sql := fmt.Sprintf("%s WHERE %s ", ql, where)
		sql += root.ResolveRestrictions()
		if data, err := db.DataCall1(self.SrcDBConn, sql); err == nil && len(data) > 0 {
			self.root = &DependencyNode{Tag: root, SQL: sql, Data: data}
		} else {
			if err == nil {
				err = errors.New("no data returned for root node, doesn't exist?")
			} else {
				fmt.Println("@FetchRoot > DataCall1 | ", err)
			}
			fmt.Println(sql)
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
			where := node.ResolveDependencyConditions(self.SrcAppConfig, dep, child)
			ql := self.GetTagQL(child)
			sql := fmt.Sprintf("%s WHERE %s ", ql, where)
			sql += child.ResolveRestrictions()
			// log.Println("@GetAllNextNodes | ", sql)
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
	// if len(self.SrcAppConfig.GetSubDependencies(node.Tag.Name)) > 0 {
	// 	log.Println("@GetAllNextNodes:", len(nodes))
	// 	log.Fatal(nodes)
	// }
	return nodes, nil
}

func (self *MigrationWorkerV2) GetAllPreviousNodes(node *DependencyNode) ([]*DependencyNode, error) {
	var nodes []*DependencyNode
	for _, dep := range self.SrcAppConfig.GetParentDependencies(node.Tag.Name) {
		for _, pdep := range dep.DependsOn {
			if parent, err := self.SrcAppConfig.GetTag(pdep.Tag); err == nil {
				where := node.ResolveParentDependencyConditions(pdep.Conditions, parent)
				ql := self.GetTagQL(parent)
				sql := fmt.Sprintf("%s WHERE %s ", ql, where)
				sql += parent.ResolveRestrictions()
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
			log.Println(fmt.Sprintf("x%2dx | FETCHING  tag for dependency { %s > %s } ", threadID, node.Tag.Name, dep.Tag))
			where := node.ResolveDependencyConditions(self.SrcAppConfig, dep, child)
			ql := self.GetTagQL(child)
			sql := fmt.Sprintf("%s WHERE %s ", ql, where)
			sql += child.ResolveRestrictions()
			sql += self.ExcludeVisited(child)
			sql += " ORDER BY random()"
			// log.Fatal(sql)
			if data, err := db.DataCall1(self.SrcDBConn, sql); err == nil {
				if len(data) > 0 {
					newNode := DependencyNode{Tag: child, SQL: sql, Data: data}
					if !self.wList.IsAlreadyWaiting(newNode) && !self.IsVisited(&newNode) {
						return &newNode, nil
					}
				}
			} else {
				fmt.Println("@GetDependentNode > DataCall1 | ", err)
				log.Fatal(sql)
				return nil, err
			}
		}
	}
	return nil, nil
}

func (self *MigrationWorkerV2) GetOwnedNode(threadID int) (*DependencyNode, error) {

	for _, own := range self.SrcAppConfig.GetShuffledOwnerships() {
		log.Println(fmt.Sprintf("x%2dx | FETCHING  tag  for ownership { %s } ", threadID, own.Tag))
		// if self.unmappedTags.Exists(own.Tag) {
		// 	log.Println(fmt.Sprintf("x%2dx |         UNMAPPED  tag  { %s } ", threadID, own.Tag))
		// 	continue
		// }
		if child, err := self.SrcAppConfig.GetTag(own.Tag); err == nil {
			where := self.root.ResolveOwnershipConditions(own, child)
			ql := self.GetTagQL(child)
			sql := fmt.Sprintf("%s WHERE %s ", ql, where)
			sql += child.ResolveRestrictions()
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
				fmt.Println("@GetOwnedNode > DataCall1 | ", err)
				log.Fatal(sql)
				return nil, err
			}
		}
	}
	return nil, nil
}

func (self *MigrationWorkerV2) PushData(tx *sql.Tx, dtable config.ToTable, pk string, mappedData MappedData, node *DependencyNode) error {

	undoActionSerialized, _ := json.Marshal(mappedData.undoAction)
	transaction.LogChange(string(undoActionSerialized), self.logTxn)
	if err := SA1_display.GenDisplayFlag(self.logTxn.DBconn, self.DstAppConfig.AppID, dtable.TableID, pk, fmt.Sprint(self.logTxn.Txn_id)); err != nil {
		fmt.Println(self.DstAppConfig.AppID, dtable.TableID, pk, fmt.Sprint(self.logTxn.Txn_id))
		log.Fatal("## DISPLAY ERROR!", err)
		return errors.New("0")
	}

	for fromTable, fromCols := range mappedData.srcTables {
		if _, ok := node.Data[fmt.Sprintf("%s.id", fromTable)]; ok {
			srcID := fmt.Sprint(node.Data[fmt.Sprintf("%s.id", fromTable)])
			if fromTableID, err := db.TableID(self.logTxn.DBconn, fromTable, self.SrcAppConfig.AppID); err == nil {
				// if err := db.InsertIntoIdentityTable(tx, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTableID, dtable.TableID, srcID, pk, fmt.Sprint(self.logTxn.Txn_id)); err != nil {
				// 	log.Println("@PushData:db.InsertIntoIdentityTable: ", self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTableID, dtable.TableID, srcID, pk, fmt.Sprint(self.logTxn.Txn_id))
				// 	log.Fatal(err)
				// 	return errors.New("0")
				// }
				if serr := db.SaveForLEvaluation(tx, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTableID, dtable.TableID, srcID, pk, strings.Join(fromCols, ","), mappedData.cols, fmt.Sprint(self.logTxn.Txn_id)); serr != nil {
					log.Println("@PushData:db.SaveForLEvaluation: ", self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTableID, dtable.TableID, srcID, pk, strings.Join(fromCols, ","), mappedData.cols, fmt.Sprint(self.logTxn.Txn_id))
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

	if strings.EqualFold(node.Tag.Name, "root") {
		return self.uid, true
	}

	if ownership := self.SrcAppConfig.GetOwnership(node.Tag.Name, self.root.Tag.Name); ownership != nil {
		for _, condition := range ownership.Conditions {
			tagAttr, err := node.Tag.ResolveTagAttr(condition.TagAttr)
			if err != nil {
				log.Fatal("@GetNodeOwner: Resolving TagAttr", err, node.Tag.Name, condition.TagAttr)
				break
			}
			depOnAttr, err := self.root.Tag.ResolveTagAttr(condition.DependsOnAttr)
			if err != nil {
				log.Fatal("@GetNodeOwner: Resolving depOnAttr", err, node.Tag.Name, condition.DependsOnAttr)
				break
			}
			if nodeVal, err := node.GetValueForKey(tagAttr); err == nil {
				if rootVal, err := self.root.GetValueForKey(depOnAttr); err == nil {
					if !strings.EqualFold(nodeVal, rootVal) {
						// fmt.Println(fmt.Sprintf("root:%s:%s; user:%s:%s", depOnAttr, rootVal, tagAttr, nodeVal))
						return nodeVal, false
					} else {
						return nodeVal, true
					}
				} else {
					fmt.Println("@GetNodeOwner: Ownership Condition Key in Root Data:", depOnAttr, "doesn't exist!")
					fmt.Println("@GetNodeOwner: root data:", self.root.Data)
					log.Fatal("@GetNodeOwner: stop here and check ownership conditions wrt root")
				}
			} else {
				fmt.Println("@GetNodeOwner: Ownership Condition Key", tagAttr, "doesn't exist!")
				fmt.Println("@GetNodeOwner: node data:", node.Data)
				fmt.Println("@GetNodeOwner: node sql:", node.SQL)
				log.Fatal("@GetNodeOwner: stop here and check ownership conditions")
			}
		}
	} else {
		log.Fatal("@GetNodeOwner: Ownership not found:", node.Tag.Name)
	}
	return "", false
}

func (self *MappedData) UpdateData(col, orgCol, fromTable string, ival interface{}) {
	self.ivals = append(self.ivals, ival)
	self.vals += fmt.Sprintf("$%d,", len(self.ivals))
	self.cols += fmt.Sprintf("%s,", col)
	self.orgCols += fmt.Sprintf("%s,", orgCol)
	if fromTable != "" {
		if _, ok := self.srcTables[fromTable]; !ok {
			self.srcTables[fromTable] = []string{strings.Split(orgCol, ".")[1]}
		} else {
			self.srcTables[fromTable] = append(self.srcTables[fromTable], strings.Split(orgCol, ".")[1])
		}
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
			log.Fatal("@FetchFromMapping: FetchForMapping | ", err)
			return err
		} else {
			data.UpdateData(toAttr, assignedTabCol, targetTabCol[0], res[targetTabCol[1]])
			node.Data[args[0]] = res[targetTabCol[1]]
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
		log.Fatal("@FetchFromMapping: unable to fetch ", args[2])
		return errors.New("Unable to fetch data from node")
	}
	return nil
}

func (self *MigrationWorkerV2) RemoveMappedDataFromNodeData(mappedData MappedData, node *DependencyNode) {
	for _, col := range strings.Split(mappedData.orgCols, ",") {
		for key := range node.Data {
			if strings.Contains(key, col) && !strings.Contains(key, ".id") {
				delete(node.Data, key)
			}
		}
	}
}

func (self *MigrationWorkerV2) IsNodeDataEmpty(node *DependencyNode) bool {
	for key := range node.Data {
		if !strings.Contains(key, ".id") {
			return false
		}
	}
	return true
}

func (self *MigrationWorkerV2) GetMappedData(toTable config.ToTable, node *DependencyNode) (MappedData, error) {

	data := MappedData{
		cols:        "",
		vals:        "",
		orgCols:     "",
		orgColsLeft: "",
		srcTables:   make(map[string][]string),
		undoAction:  new(transaction.UndoAction)}

	for toAttr, fromAttr := range toTable.Mapping {
		if strings.EqualFold("id", toAttr) {
			continue
		}
		if val, ok := node.Data[fromAttr]; ok {
			fromTokens := strings.Split(fromAttr, ".")
			data.UpdateData(toAttr, fromAttr, fromTokens[0], val)
			data.undoAction.AddData(fromAttr, val)
			data.undoAction.AddOrgTable(fromTokens[0])
		} else if strings.Contains(fromAttr, "$") {
			if inputVal, err := self.mappings.GetInput(fromAttr); err == nil {
				data.UpdateData(toAttr, fromAttr, "", inputVal)
			}
		} else if strings.Contains(fromAttr, "#") {
			assignedTabCol := strings.Trim(fromAttr, "(#ASSIGN#FETCH#REF)")
			if strings.Contains(fromAttr, "#REF") {
				if strings.Contains(fromAttr, "#FETCH") {
					self.FetchFromMapping(node, toAttr, assignedTabCol, &data)
				} else if strings.Contains(fromAttr, "#ASSIGN") {
					assignedTabColTokens := strings.Split(assignedTabCol, ",")
					referredTabCol := assignedTabColTokens[1]
					assignedTabCol = strings.Trim(assignedTabColTokens[0], "()")
					if nodeVal, ok := node.Data[assignedTabCol]; ok {
						assignedTabColTokens := strings.Split(assignedTabCol, ".")
						referredTabColTokens := strings.Split(referredTabCol, ".")
						data.UpdateData(toAttr, assignedTabCol, assignedTabColTokens[0], nodeVal)
						var fromID string
						if val, ok := node.Data[assignedTabColTokens[0]+".id"]; ok {
							fromID = fmt.Sprint(val)
						} else {
							fmt.Println(assignedTabColTokens[0], " | ", assignedTabColTokens)
							fmt.Println(node.Data)
							log.Fatal("@GetMappedData > #REF #ASSIGN> fromID: Unable to find ref value in node data")
							return data, errors.New("Unable to find ref value in node data")
						}
						data.refs = append(data.refs, MappingRef{fromID: fromID, fromMember: assignedTabColTokens[0], fromAttr: assignedTabColTokens[1], toID: fmt.Sprint(nodeVal), toAttr: referredTabColTokens[1], toMember: referredTabColTokens[0]})
						// fmt.Println(data.refs)
					}
					// log.Fatal("FOUND #REF#ASSIGN MAPPING")
				} else {
					args := strings.Split(assignedTabCol, ",")
					if nodeVal, ok := node.Data[args[0]]; ok {
						argsTokens := strings.Split(args[0], ".")
						data.UpdateData(toAttr, assignedTabCol, argsTokens[0], nodeVal)
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

					data.refs = append(data.refs, MappingRef{fromID: fromID, fromMember: firstMemberTokens[0], fromAttr: firstMemberTokens[1], toID: toID, toAttr: secondMemberTokens[1], toMember: secondMemberTokens[0]})
					// fmt.Println(toTable.Table, toAttr, fromAttr, data.refs[len(data.refs)-1])
				}
			} else if strings.Contains(fromAttr, "#ASSIGN") {
				if nodeVal, ok := node.Data[assignedTabCol]; ok {
					assignedTabColTokens := strings.Split(assignedTabCol, ".")
					data.UpdateData(toAttr, assignedTabCol, assignedTabColTokens[0], nodeVal)
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
			data.orgColsLeft += fmt.Sprintf("%s,", fromAttr)
		}
	}
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
			if tagMemberID, err := db.TableID(self.logTxn.DBconn, tagMember, self.SrcAppConfig.AppID); err == nil {
				if derr := db.UpdateLEvaluation(self.logTxn.DBconn, tagMemberID, srcID, self.logTxn.Txn_id); derr != nil {
					fmt.Println("@ERROR_UpdateLEvaluation", derr)
					fmt.Println("@QARGS:", tagMember, srcID, self.logTxn.Txn_id)
					log.Fatal(derr)
					return derr
				}
			} else {
				log.Fatal("@DeleteRow>TableID: ", err)
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

func (self *MigrationWorkerV2) GetRowsFromIDTable(app, member, id string, getFrom bool) ([]IDRow, error) {
	var idRows []IDRow
	var err error
	var idRowsDB []map[string]interface{}
	if !getFrom {
		idRowsDB, err = db.GetRowsFromIDTableByTo(self.logTxn.DBconn, app, member, id)
	} else {
		idRowsDB, err = db.GetRowsFromIDTableByFrom(self.logTxn.DBconn, app, member, id)
	}

	if err != nil {
		log.Fatal("@GetRowsFromIDTable > db.GetRowsFromIDTable, Unable to get bags | ", getFrom, self.uid, app, member, id, self.logTxn.Txn_id, err)
		return nil, err
	}
	for _, idRowDB := range idRowsDB {
		fromAppID := fmt.Sprint(idRowDB["from_app"])
		fromAppName, err := db.GetAppNameByAppID(self.logTxn.DBconn, fromAppID)
		if err != nil {
			log.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get name | ", fromAppID, err)
			return nil, err
		}

		toAppID := fmt.Sprint(idRowDB["to_app"])
		toAppName, err := db.GetAppNameByAppID(self.logTxn.DBconn, toAppID)
		if err != nil {
			log.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get name | ", toAppID, err)
			return nil, err
		}

		idRows = append(idRows, IDRow{
			FromAppName: fromAppName,
			FromAppID:   fromAppID,
			FromMember:  fmt.Sprint(idRowDB["from_member"]),
			FromID:      fmt.Sprint(idRowDB["from_id"]),
			ToAppName:   toAppName,
			ToAppID:     toAppID,
			ToMember:    fmt.Sprint(idRowDB["to_member"]),
			ToID:        fmt.Sprint(idRowDB["to_id"])})
	}
	return idRows, nil
}

func (self *MigrationWorkerV2) FetchDataFromBags(nodeData map[string]interface{}, app, member, id, dstMember string) error {
	idRows, err := self.GetRowsFromIDTable(app, member, id, false)
	if err != nil {
		log.Fatal("@FetchDataFromBags > GetRowsFromIDTable, Unable to get IDRows | ", app, member, id, false, err)
		return err
	}
	for _, idRow := range idRows {
		bagRow, err := db.GetBagByAppMemberIDV2(self.logTxn.DBconn, self.uid, idRow.FromAppID, idRow.FromMember, idRow.FromID, self.logTxn.Txn_id)
		if err != nil {
			log.Fatal("@FetchDataFromBags > GetBagByAppMemberIDV2, Unable to get bags | ", self.uid, idRow.FromAppID, idRow.FromMember, idRow.FromID, self.logTxn.Txn_id, err)
			return err
		}
		if bagRow != nil {
			bagData := make(map[string]interface{})
			if err := json.Unmarshal([]byte(fmt.Sprint(bagRow["data"])), &bagData); err != nil {
				fmt.Println(bagRow["data"])
				fmt.Println(bagRow)
				log.Fatal("@FetchDataFromBags: UNABLE TO CONVERT BAG TO MAP ", bagRow, err)
				return err
			}
			if mapping, found := self.FetchMappingsForBag(idRow.FromAppName, idRow.FromMember, self.DstAppConfig.AppName, dstMember); found {
				for _, toTable := range mapping.ToTables {
					for fromAttr, toAttr := range toTable.Mapping {
						if bagVal, ok := bagRow[fromAttr]; ok {
							if _, exists := bagData[toAttr]; !exists {
								bagData[toAttr] = bagVal
							}
							delete(bagData, fromAttr)

						}
					}
				}

				if len(bagData) == 0 {
					if err := db.DeleteBagV2(self.tx.StencilTx, fmt.Sprint(bagRow["pk"])); err != nil {
						log.Fatal("@FetchDataFromBags > DeleteBagV2, Unable to delete bag | ", bagRow["pk"])
						return err
					}
				} else {
					if jsonData, err := json.Marshal(bagData); err == nil {
						if err := db.UpdateBag(self.tx.StencilTx, fmt.Sprint(bagRow["pk"]), self.logTxn.Txn_id, jsonData); err != nil {
							log.Fatal("@FetchDataFromBags: UNABLE TO UPDATE BAG ", bagRow, err)
							return err
						} else {
							log.Fatal("@FetchDataFromBags > UpdateBag | ", fmt.Sprint(bagRow["pk"]), self.logTxn.Txn_id, jsonData)
							return err
						}
					} else {
						log.Fatal("@FetchDataFromBags > len(bagData) != 0, Unable to marshall bag | ", bagData)
						return err
					}
				}
			}
		}
		if err := self.FetchDataFromBags(nodeData, idRow.FromAppID, idRow.FromMember, idRow.FromID, dstMember); err != nil {
			log.Fatal("@FetchDataFromBags > FetchDataFromBags: Error while recursing | ", nodeData, idRow.FromAppID, idRow.FromMember, idRow.FromID)
			return err
		}
	}
	return nil
}

func (self *MigrationWorkerV2) DeleteNode(mapping config.Mapping, node *DependencyNode) error {

	if self.mtype == DELETION {

		if !self.IsNodeDataEmpty(node) {
			if err := self.SendNodeToBagWithOwnerID(node, self.uid); err != nil {
				fmt.Println(node.Tag.Name)
				fmt.Println(node.Data)
				log.Fatal("@DeleteNode > SendNodeToBagWithOwnerID:", err)
				return err
			}
		} else {
			fmt.Println(node.Data)
			log.Fatal("@DeleteNode > IsNodeDataEmpty: true")
		}

		if err := self.HandleUnmappedMembersOfNode(mapping, node); err != nil {
			fmt.Println(node.Tag.Name)
			fmt.Println(mapping, node)
			log.Fatal("@DeleteNode > HandleUnmappedMembersOfNode:", err)
			return err
		}

		if err := self.DeleteRow(node); err != nil {
			fmt.Println(node.Tag.Name)
			fmt.Println(node)
			log.Fatal("@DeleteNode > DeleteRow:", err)
			return err
		}
	}

	return nil
}

func (self *MigrationWorkerV2) DeleteRoot(threadID int) error {
	if err := self.InitTransactions(); err != nil {
		log.Fatal("@DeleteRoot > InitTransactions", err)
		return err
	} else {
		defer self.tx.SrcTx.Rollback()
		defer self.tx.DstTx.Rollback()
		defer self.tx.StencilTx.Rollback()
	}
	if mapping, found := self.FetchMappingsForNode(self.root); found {
		if err := self.DeleteNode(mapping, self.root); err != nil {
			log.Fatal("@DeleteRoot:", err)
			return err
		}
	} else {
		fmt.Println(self.root)
		log.Fatal("@DeleteRoot: Can't find mappings for root | ", mapping, found)
	}
	if err := self.CommitTransactions(); err != nil {
		return err
	} else {
		log.Println(fmt.Sprintf("x%2dx COMMITTED node { %s } ", threadID, self.root.Tag.Name))
	}
	return nil
}

func (self *MigrationWorkerV2) MigrateNode(mapping config.Mapping, node *DependencyNode, isBag bool) error {

	// fetchDataFromBags, id table recursion?
	var allMappedData []MappedData
	for _, toTable := range mapping.ToTables {
		if self.CheckMappingConditions(toTable, node) {
			continue
		}
		if mappedData, _ := self.GetMappedData(toTable, node); len(mappedData.cols) > 0 && len(mappedData.vals) > 0 && len(mappedData.ivals) > 0 {
			mappedData.undoAction.AddDstTable(toTable.Table)
			if id, err := db.InsertRowIntoAppDB(self.tx.DstTx, toTable.Table, mappedData.cols, mappedData.vals, mappedData.ivals...); err == nil {
				if toTableID, err := db.TableID(self.logTxn.DBconn, toTable.Table, self.DstAppConfig.AppID); err == nil {
					for fromTable := range mappedData.srcTables {
						if fromTableID, err := db.TableID(self.logTxn.DBconn, fromTable, self.SrcAppConfig.AppID); err == nil {
							if val, ok := node.Data[fromTable+".id"]; ok {
								fromID := fmt.Sprint(val)
								if err := db.InsertIntoIdentityTable(self.tx.StencilTx, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTableID, toTableID, fromID, fmt.Sprint(id), fmt.Sprint(self.logTxn.Txn_id)); err != nil {
									fmt.Println("@MigrateNode: InsertIntoIdentityTable")
									fmt.Println("@Params:", toTable.Table, fmt.Sprint(id), mappedData.orgCols, mappedData.cols, mappedData.undoAction, node)
									log.Fatal(err)
									return err
								}
							} else {
								fmt.Println(node.Data)
								log.Fatal("@MigrateNode: InsertIntoIdentityTable | " + fromTable + ".id doesn't exist")
							}
						} else {
							log.Fatal("@MigrateNode > TableID, fromTable: error in getting table id for member! ", toTable.Table, err)
							return err
						}
					}

					if err := self.PushData(self.tx.StencilTx, toTable, fmt.Sprint(id), mappedData, node); err != nil {
						fmt.Println("@MigrateNode")
						fmt.Println("@Params:", toTable.Table, fmt.Sprint(id), mappedData.orgCols, mappedData.cols, mappedData.undoAction, node)
						log.Fatal(err)
						return err
					}

					if len(toTable.Media) > 0 {
						if filePathCol, ok := toTable.Media["path"]; ok {
							if filePath, ok := node.Data[filePathCol]; ok {
								if err := self.TransferMedia(fmt.Sprint(filePath)); err != nil {
									log.Fatal("@MigrateNode > TransferMedia: ", err)
								}
							}
						} else {
							log.Fatal("@MigrateNode > toTable.Media: Path not found in map!")
						}
					}
				} else {
					log.Fatal("@MigrateNode > TableID, toTable: error in getting table id for member! ", toTable.Table, err)
					return err
				}
				allMappedData = append(allMappedData, mappedData)
			} else {
				fmt.Println(toTable.Table, mappedData.cols, mappedData.vals, mappedData.ivals)
				log.Fatal("@MigrateNode > InsertRowIntoAppDB:", err)
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

	for _, mappedData := range allMappedData {
		self.RemoveMappedDataFromNodeData(mappedData, node)
	}

	// log.Fatal("Check here!")

	if !isBag && !strings.EqualFold(node.Tag.Name, "root") {
		if err := self.DeleteNode(mapping, node); err != nil {
			log.Fatal("@MigrateNode > DeleteNode:", err)
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

	for _, tagMember := range node.Tag.Members {
		if _, ok := node.Data[fmt.Sprintf("%s.id", tagMember)]; ok {
			srcID := fmt.Sprint(node.Data[fmt.Sprintf("%s.id", tagMember)])
			if tagMemberID, err := db.TableID(self.logTxn.DBconn, tagMember, self.SrcAppConfig.AppID); err == nil {
				if serr := db.SaveForEvaluation(self.logTxn.DBconn, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, tagMemberID, "n/a", srcID, "n/a", "*", "n/a", fmt.Sprint(self.logTxn.Txn_id)); serr != nil {
					log.Fatal(serr)
				}
			} else {
				log.Fatal("@HandleUnmappedTags > Table id:", err)
			}
		}
	}
	return errors.New("2")
}

func (self *MigrationWorkerV2) HandleUnmappedNode(node *DependencyNode) error {
	if !strings.EqualFold(self.mtype, DELETION) {
		return errors.New("3")
	} else {
		if err := self.SendNodeToBag(node); err != nil {
			return err
		} else {
			return errors.New("2")
		}
	}
}

func (self *MigrationWorkerV2) FetchMappingsForBag(srcApp, dstApp, srcMember, dstMember string) (config.Mapping, bool) {
	var combinedMapping config.Mapping
	appMappings := config.GetSchemaMappingsFor(srcApp, dstApp)
	mappingFound := false
	for _, mapping := range appMappings.Mappings {
		if mappedTables := helper.IntersectString([]string{srcMember}, mapping.FromTables); len(mappedTables) > 0 {
			combinedMapping.FromTables = append(combinedMapping.FromTables, mapping.FromTables...)
			combinedMapping.ToTables = append(combinedMapping.ToTables, mapping.ToTables...)
			mappingFound = true
		}
	}
	fmt.Println("srcMember, dstMember", srcMember, dstMember)
	fmt.Println(combinedMapping)
	log.Fatal("Check mappings for bag.")

	return combinedMapping, mappingFound
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
		bagData := make(map[string]interface{})
		for col, val := range node.Data {
			colTokens := strings.Split(col, ".")
			colMember := colTokens[0]
			// colAttr := colTokens[1]
			if strings.Contains(colMember, member) {
				bagData[col] = val
			}
		}
		if len(bagData) > 0 {
			if id, ok := node.Data[member+".id"]; ok {
				srcID := fmt.Sprint(id)
				if jsonData, err := json.Marshal(bagData); err == nil {
					if err := db.CreateNewBag(self.tx.StencilTx, self.SrcAppConfig.AppID, memberID, srcID, ownerID, fmt.Sprint(self.logTxn.Txn_id), jsonData); err != nil {
						log.Fatal("@SendMemberToBag: error in creating bag! ", err)
						return err
					}
					// if serr := db.SaveForEvaluation(self.logTxn.DBconn, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, memberID, "n/a", srcID, "n/a", "*", "n/a", fmt.Sprint(self.logTxn.Txn_id)); serr != nil {
					// 	log.Fatal("@SendMemberToBag > SaveForEvaluation =>", serr)
					// }
					if derr := db.ReallyDeleteRowFromAppDB(self.tx.SrcTx, member, fmt.Sprint(id)); derr != nil {
						fmt.Println("@SendMemberToBag > DeleteRowFromAppDB")
						log.Fatal(derr)
						return derr
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

func (self *MigrationWorkerV2) SendNodeToBagWithOwnerID(node *DependencyNode, ownerID string) error {
	for _, member := range node.Tag.Members {
		if err := self.SendMemberToBag(node, member, ownerID, true); err != nil {
			fmt.Println(node)
			log.Fatal("@SendNodeToBagWithOwnerID > SendMemberToBag: ownerID error! ")
			return err
		}
	}
	if err := self.AddInnerReferences(node, ""); err != nil {
		fmt.Println(node.Tag.Members)
		fmt.Println(node.Tag.InnerDependencies)
		fmt.Println(node.Data)
		log.Fatal("@SendNodeToBagWithOwnerID > AddInnerReferences: Adding Inner References failed ", err)
		return err
	}
	return nil
}

func (self *MigrationWorkerV2) SendNodeToBag(node *DependencyNode) error {
	if ownerID, _ := self.GetNodeOwner(node); len(ownerID) > 0 {
		if err := self.SendNodeToBagWithOwnerID(node, ownerID); err != nil {
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
			return self.MigrateNode(mapping, node, isBag)
		}
		if wNode, err := self.HandleWaitingList(mapping, tagMembers, node); wNode != nil && err == nil {
			return self.MigrateNode(mapping, wNode, isBag)
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
	return nil
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
	return nil
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

func (self *MigrationWorkerV2) MigrateBags(threadID int) error {

	prevIDs := reference_resolution.GetPrevUserIDs(self.SrcAppConfig.AppID, self.uid)
	prevIDs = append(prevIDs, []string{self.SrcAppConfig.AppID, self.uid})

	for _, prevID := range prevIDs {

		appID, userID := prevID[0], prevID[1]

		bags, err := db.GetBagsV2(self.logTxn.DBconn, appID, userID, self.logTxn.Txn_id)
		if err != nil {
			log.Fatal(fmt.Sprintf("x%2dx UNABLE TO FETCH BAGS FOR USER: %s | %s", threadID, self.uid, err))
			return err
		}
		for _, bag := range bags {
			bagAppID := fmt.Sprint(bag["app"])
			srcMember := fmt.Sprint(bag["member"])
			srcMemberName, err := db.TableName(self.logTxn.DBconn, srcMember, bagAppID)
			if err != nil {
				log.Fatal("@MigrateBags > Table Name: ", err)
			}
			bagID := fmt.Sprint(bag["id"])
			log.Println(fmt.Sprintf("~%2d~ Current    Bag: { %s } | ID: %s, App: %s ", threadID, srcMemberName, bagID, bagAppID))
			bagData := make(map[string]interface{})
			if err := json.Unmarshal(bag["data"].([]byte), &bagData); err != nil {
				fmt.Println("BAG >> ", bag)
				log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO CONVERT BAG TO MAP: %s | %s", threadID, self.uid, err))
				return err
			}

			// bagAppName, err := db.GetAppNameByAppID(self.logTxn.DBconn, bagAppID)
			// if err != nil {
			// 	log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO GET BAG APP NAME BY ID : %s | %s", threadID, bagAppID, err))
			// 	return err
			// }
			// bagAppConfig := self.DstAppConfig
			// if !strings.EqualFold(bagAppID, self.DstAppConfig.AppID) {
			// 	bagAppConfig, err = config.CreateAppConfig(bagAppName, bagAppID)
			// 	if err != nil {
			// 		log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO CREATE BAG APP CONFIG : %s, %s | %s", threadID, bagAppID, bagAppName, err))
			// 		log.Fatal(err)
			// 	}
			// }
			// bagWorker := CreateMigrationWorkerV2WithAppsConfig(self.uid, self.logTxn, self.mtype, self.SrcAppConfig, bagAppConfig, threadID)
			// bagTag, err := bagAppConfig.GetTagByMember(srcMember)
			// if err != nil {
			// 	log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO GET BAG TAG BY MEMBER : %s | %s", threadID, srcMember, err))
			// 	return err
			// }
			// bagNode := DependencyNode{Tag: *bagTag, Data: bagData}
			// if err := bagWorker.HandleMigration(&bagNode, true); err != nil {
			// 	fmt.Println(bag)
			// 	log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO MIGRATE BAG : %s | %s ", threadID, bagID, err))
			// 	return err
			// }
			// if self.IsNodeDataEmpty(&bagNode) {
			// 	if err := db.DeleteBagV2(self.tx.StencilTx, bagID); err != nil {
			// 		fmt.Println(bag)
			// 		log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO DELETE BAG : %s | %s ", threadID, bagID, err))
			// 		return err
			// 	}
			// } else {
			// 	if jsonData, err := json.Marshal(bagNode.Data); err == nil {
			// 		if err := db.UpdateBag(self.tx.StencilTx, bagID, self.logTxn.Txn_id, jsonData); err != nil {
			// 			fmt.Println(bag)
			// 			log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO UPDATE BAG : %s | %s ", threadID, bagID, err))
			// 			return err
			// 		}
			// 	} else {
			// 		fmt.Println(bagNode.Data)
			// 		log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO MARSHALL BAG DATA : %s | %s ", threadID, bagID, err))
			// 		return err
			// 	}
			// }
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
	// log.Fatal("@CommitTransactions: About to Commit!")
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

func (self *MigrationWorkerV2) CallMigration(node *DependencyNode, threadID int) error {

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
			if strings.EqualFold(err.Error(), "3") {
				log.Println(fmt.Sprintf("x%2dx UNMAPPED  node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
			} else if strings.EqualFold(err.Error(), "2") {
				log.Println(fmt.Sprintf("x%2dx Sent2Bag  node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
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
		log.Println(fmt.Sprintf("x%2dx VISITED   node { %s } | root [%s] : owner [%s]", threadID, node.Tag.Name, self.uid, ownerID))
		self.MarkAsVisited(node)
	}
	return nil
}

func (self *MigrationWorkerV2) DeletionMigration(node *DependencyNode, threadID int) error {

	if strings.EqualFold(node.Tag.Name, "root") {
		if err := self.CallMigration(node, threadID); err != nil {
			return err
		}
	}

	for {
		if adjNode, err := self.GetAdjNode(node, threadID); err != nil {
			return err
		} else {
			if adjNode == nil {
				break
			}
			nodeIDAttr, _ := node.Tag.ResolveTagAttr("id")
			adjNodeIDAttr, _ := adjNode.Tag.ResolveTagAttr("id")
			log.Println(fmt.Sprintf("~%2d~ Current   Node: { %s } | ID: %v ", threadID, node.Tag.Name, node.Data[nodeIDAttr]))
			log.Println(fmt.Sprintf("~%2d~ Adjacent  Node: { %s } | ID: %v ", threadID, adjNode.Tag.Name, adjNode.Data[adjNodeIDAttr]))
			if err := self.DeletionMigration(adjNode, threadID); err != nil {
				log.Fatal(fmt.Sprintf("~%2d~ ERROR! NODE : { %s } | ID: %v, ADJ_NODE : { %s } | ID: %v | err: [ %s ]", threadID, node.Tag.Name, node.Data[nodeIDAttr], adjNode.Tag.Name, adjNode.Data[adjNodeIDAttr], err))
				return err
			}
		}
	}

	log.Println(fmt.Sprintf("#%2d# Process   Node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))

	if strings.EqualFold(node.Tag.Name, "root") {
		return self.DeleteRoot(threadID)
	} else {
		if err := self.CallMigration(node, threadID); err != nil {
			return err
		}
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
				if strings.EqualFold(err.Error(), "3") {
					log.Println(fmt.Sprintf("x%2dx | IGNORED   node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
				} else if strings.EqualFold(err.Error(), "2") {
					log.Println(fmt.Sprintf("x%2dx | BAGGED?   node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
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

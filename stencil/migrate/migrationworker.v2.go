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

func CreateBagWorkerV2(uid, srcAppID, dstAppID string, logTxn *transaction.Log_txn, mtype string, threadID int, isBlade ...bool) MigrationWorkerV2 {

	srcApp, err := db.GetAppNameByAppID(logTxn.DBconn, srcAppID)
	if err != nil {
		log.Fatal(err)
	}
	dstApp, err := db.GetAppNameByAppID(logTxn.DBconn, dstAppID)
	if err != nil {
		log.Fatal(err)
	}

	srcAppConfig, err := config.CreateAppConfig(srcApp, srcAppID)
	if err != nil {
		log.Fatal(err)
	}
	dstAppConfig, err := config.CreateAppConfig(dstApp, dstAppID, isBlade...)
	if err != nil {
		log.Fatal(err)
	}

	var mappings *config.MappedApp

	if srcAppID == dstAppID {
		mappings = config.GetSelfSchemaMappings(logTxn.DBconn, srcAppID, srcApp)
		// log.Fatal(mappings)
	} else {
		mappings = config.GetSchemaMappingsFor(srcAppConfig.AppName, dstAppConfig.AppName)
		if mappings == nil {
			log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcAppConfig.AppName, dstAppConfig.AppName))
		}
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
		DstDBConn:    db.GetDBConn(dstAppConfig.AppName, isBlade...),
		logTxn:       &transaction.Log_txn{DBconn: logTxn.DBconn, Txn_id: logTxn.Txn_id},
		mtype:        mtype,
		visitedNodes: make(map[string]map[string]bool)}
	if err := mWorker.FetchRoot(threadID); err != nil {
		log.Fatal(err)
	}
	mWorker.FTPClient = GetFTPClient()
	log.Println("Bag Worker Created for thread: ", threadID)
	fmt.Println("************************************************************************")
	return mWorker
}

func CreateMigrationWorkerV2(uid, srcApp, srcAppID, dstApp, dstAppID string, logTxn *transaction.Log_txn, mtype string, mappings *config.MappedApp, threadID int, isBlade ...bool) MigrationWorkerV2 {
	srcAppConfig, err := config.CreateAppConfig(srcApp, srcAppID)
	if err != nil {
		log.Fatal(err)
	}
	dstAppConfig, err := config.CreateAppConfig(dstApp, dstAppID, isBlade...)
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
		DstDBConn:    db.GetDBConn(dstApp, isBlade...),
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

func (self *MigrationWorkerV2) CloseDBConns() {

	self.SrcDBConn.Close()
	self.DstDBConn.Close()
	self.SrcAppConfig.CloseDBConns()
	self.DstAppConfig.CloseDBConns()
}

func (self *MigrationWorkerV2) RenewDBConn(isBlade ...bool) {
	self.CloseDBConns()
	self.logTxn.DBconn.Close()
	self.logTxn.DBconn = db.GetDBConn(db.STENCIL_DB)
	self.SrcDBConn = db.GetDBConn(self.SrcAppConfig.AppName)
	self.DstDBConn = db.GetDBConn(self.DstAppConfig.AppName, isBlade...)
	self.SrcAppConfig.DBConn = db.GetDBConn(self.SrcAppConfig.AppName)
	self.SrcAppConfig.DBConn = db.GetDBConn(self.DstAppConfig.AppName, isBlade...)
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

func (node *DependencyNode) ResolveParentDependencyConditions(dconditions []config.DCondition, parentTag config.Tag) (string, error) {

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
				return "", errors.New("Returning empty from restricted. Why?")
			}
		}
		depOnAttr, err := parentTag.ResolveTagAttr(condition.DependsOnAttr)
		if err != nil {
			log.Println(err, parentTag.Name, condition.DependsOnAttr)
			log.Fatal("@ResolveParentDependencyConditions: depOnAttr in condition doesn't exist? ", condition.DependsOnAttr)
			break
		}
		if val, ok := node.Data[tagAttr]; ok {
			if val == nil {
				return "", errors.New(fmt.Sprintf("trying to assign %s = %s, value is nil in node %s ", tagAttr, depOnAttr, node.Tag.Name))
			}
			if conditionStr != "" {
				conditionStr += " AND "
			}
			conditionStr += fmt.Sprintf("%s = '%v'", depOnAttr, val)
		} else {
			fmt.Println(node.Data)
			log.Fatal("ResolveParentDependencyConditions:", tagAttr, " doesn't exist in node data? ", node.Tag.Name)
		}
	}
	return conditionStr, nil
}

func (node *DependencyNode) ResolveDependencyConditions(SrcAppConfig config.AppConfig, dep config.Dependency, tag config.Tag) (string, error) {

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
					if nodeVal, ok := node.Data[depOnAttr]; ok {
						if nodeVal == nil {
							return "", errors.New(fmt.Sprintf("trying to assign %s = %s, value is nil in node %s ", tagAttr, depOnAttr, node.Tag.Name))
						}
						if conditionStr != "" || where != "" {
							conditionStr += " AND "
						}
						conditionStr += fmt.Sprintf("%s = '%v'", tagAttr, nodeVal)
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
	return where, nil
}

func (root *DependencyNode) ResolveOwnershipConditions(own config.Ownership, tag config.Tag) (string, error) {

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
		if nodeVal, ok := root.Data[depOnAttr]; ok {
			if nodeVal == nil {
				return "", errors.New(fmt.Sprintf("trying to assign %s = %s, value is nil in node %s ", tagAttr, depOnAttr, root.Tag.Name))
			}
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
	return where, nil
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
			if where, err := node.ResolveDependencyConditions(self.SrcAppConfig, dep, child); err == nil {
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
				log.Println("@GetAllNextNodes > ResolveDependencyConditions | ", err)
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
				if where, err := node.ResolveParentDependencyConditions(pdep.Conditions, parent); err == nil {
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
			if where, err := node.ResolveDependencyConditions(self.SrcAppConfig, dep, child); err == nil {
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
			} else {
				log.Println("@GetDependentNode > ResolveDependencyConditions | ", err)
				// log.Fatal(err)
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
			if where, err := self.root.ResolveOwnershipConditions(own, child); err == nil {
				ql := self.GetTagQL(child)
				sql := fmt.Sprintf("%s WHERE %s ", ql, where)
				sql += child.ResolveRestrictions()
				sql += self.ExcludeVisited(child)
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
			} else {

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
			srcID := node.Data[fmt.Sprintf("%s.id", fromTable)]
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

func (self *MigrationWorkerV2) ValidateMappingConditions(toTable config.ToTable, node *DependencyNode) bool {
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
							log.Fatal("@ValidateMappingConditions: Case not found:" + conditionVal)
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
							log.Fatal("@ValidateMappingConditions: return false, from conditionVal[:1] == $")
							return false
						}
					} else {
						fmt.Println("node data:", node.Data)
						fmt.Println(conditionKey, conditionVal)
						log.Fatal("@ValidateMappingConditions: input doesn't exist?", err)
					}
				} else if !strings.EqualFold(fmt.Sprint(nodeVal), conditionVal) {
					// log.Println(conditionKey, conditionVal, "!=", nodeVal)
					return false
				} else {
					// fmt.Println(*nodeVal, "==", conditionVal)
				}
			} else {
				log.Println("Condition Key", conditionKey, "doesn't exist!")
				fmt.Println("node data:", node.Data)
				fmt.Println("node sql:", node.SQL)
				log.Fatal("@ValidateMappingConditions: stop here and check")
				return false
			}
		}
	}
	return true
}

func (self *MigrationWorkerV2) ValidateMappedTableData(toTable config.ToTable, mappedData MappedData) bool {
	for mappedCol, srcMappedCol := range toTable.Mapping {
		if strings.Contains(srcMappedCol, "$") {
			continue
		}
		for i, mCol := range strings.Split(mappedData.cols, ",") {
			if strings.EqualFold(mappedCol, mCol) {
				if mappedData.ivals[i] != nil {
					return true
				}
			}
		}
	}
	return false
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
	if ival == nil {
		return
	}
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
		} else if len(res) > 0 {
			// fmt.Println("FETCHED DATA ", res)
			data.UpdateData(toAttr, args[0], targetTabCol[0], res[targetTabCol[1]])
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
		} else {
			log.Println("@FetchFromMapping: FetchForMapping | Returned data is nil! Previous node already migrated?", res, targetTabCol[0], targetTabCol[1], comparisonTabCol[1], fmt.Sprint(nodeVal))
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
					if nodeVal, ok := node.Data[assignedTabCol]; ok && nodeVal != nil {
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
						data.UpdateData(toAttr, args[0], argsTokens[0], nodeVal)
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
	data.undoAction.AddDstTable(toTable.Table)
	return data, nil
}

func (self *MigrationWorkerV2) DeleteRow(node *DependencyNode) error {
	for _, tagMember := range node.Tag.Members {
		idCol := fmt.Sprintf("%s.id", tagMember)
		if nodeVal, ok := node.Data[idCol]; ok && nodeVal != nil {
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
			// log.Println("node.Data =>", node.Data)
			log.Println("@DeleteRow:  '", idCol, "' not present or is null in node data!", nodeVal)
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
	for {
		if err := self.FTPClient.Stor(fsName, file); err != nil {
			log.Println("File Transfer Failed: ", err)
			self.FTPClient.Quit()
			self.FTPClient = GetFTPClient()
			continue
			// return err
		}
		break
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

		fromMemberID := fmt.Sprint(idRowDB["from_member"])
		fromMember, err := db.TableName(self.logTxn.DBconn, fromMemberID, fromAppID)
		if err != nil {
			log.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get fromMember | ", fromMemberID, fromAppID, err)
			return nil, err
		}

		toMemberID := fmt.Sprint(idRowDB["to_member"])
		toMember, err := db.TableName(self.logTxn.DBconn, toMemberID, toAppID)
		if err != nil {
			log.Fatal("@GetRowsFromIDTable > db.GetAppNameByAppID, Unable to get toMember | ", toMemberID, toAppID, err)
			return nil, err
		}

		idRows = append(idRows, IDRow{
			FromAppName:  fromAppName,
			FromAppID:    fromAppID,
			FromMemberID: fromMemberID,
			FromMember:   fromMember,
			FromID:       fmt.Sprint(idRowDB["from_id"]),
			ToAppName:    toAppName,
			ToAppID:      toAppID,
			ToMember:     toMember,
			ToMemberID:   toMemberID,
			ToID:         fmt.Sprint(idRowDB["to_id"])})
	}
	return idRows, nil
}

func (self *MigrationWorkerV2) MergeBagDataWithMappedData(mappedData *MappedData, node *DependencyNode, toTable config.ToTable) error {

	toTableData := make(map[string]interface{})

	prevUIDs := reference_resolution.GetPrevUserIDs(self.SrcAppConfig.AppID, self.uid)
	prevUIDs[self.SrcAppConfig.AppID] = self.uid

	for fromTable := range mappedData.srcTables {
		if fromTableID, err := db.TableID(self.logTxn.DBconn, fromTable, self.SrcAppConfig.AppID); err == nil {
			if fromID, ok := node.Data[fromTable+".id"]; ok {
				if err := self.FetchDataFromBags(toTableData, prevUIDs, self.SrcAppConfig.AppID, fromTableID, fmt.Sprint(fromID), toTable.TableID, toTable.Table); err != nil {
					log.Fatal("@MigrateNode > FetchDataFromBags | ", err)
				}
			} else {
				log.Fatal("@MigrateNode > FetchDataFromBags > id doesn't exist in table ", fromTable, err)
			}
		} else {
			log.Fatal("@MigrateNode > FetchDataFromBags > TableID, fromTable: error in getting table id for member! ", fromTable, err)
		}
	}

	if len(toTableData) > 0 {

		for col, val := range toTableData {
			if !strings.Contains(mappedData.cols, col) {
				mappedData.cols += "," + col
				mappedData.ivals = append(mappedData.ivals, val)
				mappedData.vals += fmt.Sprintf(",$%d", len(mappedData.ivals))
			}
		}
		mappedData.Trim(",")
		log.Println("@MigrateNode > FetchDataFromBags > Data merged for: ", toTable.Table)
	}

	return nil
}

func (self *MigrationWorkerV2) FetchDataFromBags(toTableData map[string]interface{}, prevUIDs map[string]string, app, member, id, dstMemberID, dstMemberName string) error {

	idRows, err := self.GetRowsFromIDTable(app, member, id, false)

	if err != nil {
		log.Fatal("@FetchDataFromBags > GetRowsFromIDTable, Unable to get IDRows | ", app, member, id, false, err)
		return err
	}
	for _, idRow := range idRows {

		bagRow, err := db.GetBagByAppMemberIDV2(self.logTxn.DBconn, prevUIDs[idRow.FromAppID], idRow.FromAppID, idRow.FromMemberID, idRow.FromID, self.logTxn.Txn_id)
		if err != nil {
			log.Fatal("@FetchDataFromBags > GetBagByAppMemberIDV2, Unable to get bags | ", prevUIDs[idRow.FromAppID], idRow.FromAppID, idRow.FromMemberID, idRow.FromID, self.logTxn.Txn_id, err)
			return err
		}
		if bagRow != nil {

			bagData := make(map[string]interface{})
			if err := json.Unmarshal(bagRow["data"].([]byte), &bagData); err != nil {
				fmt.Println(bagRow["data"])
				fmt.Println(bagRow)
				log.Fatal("@FetchDataFromBags: UNABLE TO CONVERT BAG TO MAP ", bagRow, err)
				return err
			}

			if mapping, found := self.FetchMappingsForBag(idRow.FromAppName, idRow.FromAppID, self.DstAppConfig.AppName, self.DstAppConfig.AppID, idRow.FromMember, dstMemberName); found {

				for _, toTable := range mapping.ToTables {
					for fromAttr, toAttr := range toTable.Mapping {
						if _, ok := toTableData[fromAttr]; !ok {
							if bagVal, exists := bagData[toAttr]; exists {
								toTableData[fromAttr] = bagVal
							}
						}
						delete(bagData, toAttr)
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
						}
					} else {
						log.Fatal("@FetchDataFromBags > len(bagData) != 0, Unable to marshall bag | ", bagData)
						return err
					}
				}
			}
		} else {
			log.Fatal("@FetchDataFromBags > GetBagByAppMemberIDV2, No bags found for | ", prevUIDs[idRow.FromAppID], idRow.FromAppID, idRow.FromMember, idRow.FromID, self.logTxn.Txn_id)
		}

		if err := self.FetchDataFromBags(toTableData, prevUIDs, idRow.FromAppID, idRow.FromMemberID, idRow.FromID, dstMemberID, dstMemberName); err != nil {
			log.Fatal("@FetchDataFromBags > FetchDataFromBags: Error while recursing | ", toTableData, idRow.FromAppID, idRow.FromMember, idRow.FromID)
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
		if self.mtype == NAIVE {
			if err := self.DeleteRow(self.root); err != nil {
				log.Fatal("@DeleteRoot:", err)
				return err
			}
		} else if self.mtype == DELETION {
			if err := self.DeleteNode(mapping, self.root); err != nil {
				log.Fatal("@DeleteRoot:", err)
				return err
			}
		} else {
			log.Fatal("ATTEMPTED DELETION IN DISALLOWED MIGRATION TYPE!")
		}
	} else {
		fmt.Println(self.root)
		log.Fatal("@DeleteRoot: Can't find mappings for root | ", mapping, found)
	}
	if err := self.CommitTransactions(); err != nil {
		return err
	} else {
		log.Println(fmt.Sprintf("x%2dx DELETED ROOT ", threadID))
	}
	return nil
}

func (self *MigrationWorkerV2) MigrateNode(mapping config.Mapping, node *DependencyNode, isBag bool) error {

	// fetchDataFromBags, id table recursion?

	var allMappedData []MappedData
	for _, toTable := range mapping.ToTables {
		if !self.ValidateMappingConditions(toTable, node) {
			continue
		}
		if mappedData, _ := self.GetMappedData(toTable, node); len(mappedData.cols) > 0 && len(mappedData.vals) > 0 && len(mappedData.ivals) > 0 {
			if !self.ValidateMappedTableData(toTable, mappedData) {
				continue
			}

			if self.mtype == DELETION {
				// fmt.Println("Unmerged Mapped Data: ", mappedData)
				if err := self.MergeBagDataWithMappedData(&mappedData, node, toTable); err != nil {
					log.Fatal("@MigrateNode > MergeDataFromBagsWithMappedData | ", err)
				}
				// log.Fatal("Merged Mapped Data: ", mappedData)
			}

			if id, err := db.InsertRowIntoAppDB(self.tx.DstTx, toTable.Table, mappedData.cols, mappedData.vals, mappedData.ivals...); err == nil {
				for fromTable := range mappedData.srcTables {
					if fromTableID, err := db.TableID(self.logTxn.DBconn, fromTable, self.SrcAppConfig.AppID); err == nil {
						if fromID, ok := node.Data[fromTable+".id"]; ok {
							// fromID := fmt.Sprint(val.(int))
							if err := db.InsertIntoIdentityTable(self.tx.StencilTx, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTableID, toTable.TableID, fromID, fmt.Sprint(id), fmt.Sprint(self.logTxn.Txn_id)); err != nil {
								fmt.Println("@MigrateNode: InsertIntoIdentityTable")
								fmt.Println("@Args: ", self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTableID, toTable.TableID, fromID, fmt.Sprint(id), fmt.Sprint(self.logTxn.Txn_id))
								fmt.Println("@Params:", toTable.Table, fmt.Sprint(id), mappedData.orgCols, mappedData.cols, mappedData.undoAction, node)
								log.Fatal(err)
								return err
							}
						} else {
							fmt.Println(node.Data)
							log.Fatal("@MigrateNode: InsertIntoIdentityTable | " + fromTable + ".id doesn't exist")
						}
					} else {
						log.Fatal("@MigrateNode > TableID, fromTable: error in getting table id for member! ", fromTable, err)
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
				allMappedData = append(allMappedData, mappedData)
			} else {
				fmt.Println("@Args: ", toTable.Table, mappedData.cols, mappedData.vals, mappedData.ivals, mappedData.srcTables)
				fmt.Println("@NODE: ", node.Tag.Name, node.Data)
				log.Fatalln("@MigrateNode > InsertRowIntoAppDB: ", err)
				return err
			}

			if err := self.AddMappedReferences(mappedData.refs); err != nil {
				log.Println(mappedData.refs)
				log.Fatal("@MigrateNode > AddMappedReferences: ", err)
				return err
			}
		} else {
			// fmt.Println("cols:", mappedData.cols)
			// fmt.Println("vals:", mappedData.vals)
			// fmt.Println("ivals:", mappedData.ivals)
			// fmt.Println("toTable.Table:", toTable.Table)
			// fmt.Println("toTable.Mapping:", toTable.Mapping)
			// fmt.Println("node.Data:", node.Data)
			// log.Println("@MigrateNode > GetMappedData > If Conditions failed | ", node.Tag.Name, " -> ", toTable.Table)
			// fmt.Println(node.Tag.Name, " -> ", toTable.Table)
			// time.Sleep(time.Second * 5)
			continue
		}
	}

	for _, mappedData := range allMappedData {
		self.RemoveMappedDataFromNodeData(mappedData, node)
	}

	// log.Fatal("Check here!")

	if !isBag && !strings.EqualFold(node.Tag.Name, "root") {
		if self.mtype == DELETION {
			if err := self.DeleteNode(mapping, node); err != nil {
				log.Fatal("@MigrateNode > DeleteNode:", err)
				return err
			}
		} else if self.mtype == NAIVE {
			if err := self.DeleteRow(node); err != nil {
				log.Fatal("@MigrateNode > DeleteNode:", err)
				return err
			} else {
				log.Println(fmt.Sprintf("DELETED node { %s }", node.Tag.Name))
			}
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

func (self *MigrationWorkerV2) FetchMappingsForBag(srcApp, srcAppID, dstApp, dstAppID, srcMember, dstMember string) (config.Mapping, bool) {

	var combinedMapping config.Mapping
	var appMappings config.MappedApp
	if srcApp == dstApp {
		appMappings = *config.GetSelfSchemaMappings(self.logTxn.DBconn, srcAppID, srcApp)
	} else {
		appMappings = *config.GetSchemaMappingsFor(srcApp, dstApp)
	}
	mappingFound := false
	for _, mapping := range appMappings.Mappings {
		if mappedTables := helper.IntersectString([]string{srcMember}, mapping.FromTables); len(mappedTables) > 0 {
			for _, toTableMapping := range mapping.ToTables {
				if strings.EqualFold(dstMember, toTableMapping.Table) {
					combinedMapping.FromTables = append(combinedMapping.FromTables, mapping.FromTables...)
					combinedMapping.ToTables = append(combinedMapping.ToTables, mapping.ToTables...)
					mappingFound = true
				}
			}

		}
	}
	// fmt.Println(">>>>>>>>", srcApp, srcAppID, dstApp, dstAppID, srcMember, dstMember, " | Mappings | ", combinedMapping)
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
	if !strings.EqualFold(self.mtype, DELETION) {
		return nil
	}
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
			if id, ok := node.Data[member+".id"]; ok && id != nil {
				srcID := fmt.Sprint(id)
				if jsonData, err := json.Marshal(bagData); err == nil {
					if err := db.CreateNewBag(self.tx.StencilTx, self.SrcAppConfig.AppID, memberID, srcID, ownerID, fmt.Sprint(self.logTxn.Txn_id), jsonData); err != nil {
						fmt.Println(self.SrcAppConfig.AppName, member, srcID, ownerID)
						fmt.Println(bagData)
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
				if !fromNode {
					if err := self.AddInnerReferences(node, member); err != nil {
						fmt.Println(node.Tag.Members)
						fmt.Println(node.Tag.InnerDependencies)
						fmt.Println(node.Data)
						log.Fatal("@SendMemberToBag > AddInnerReferences: Adding Inner References failed ", err)
						return err
					}
				}
			} else {
				// fmt.Println(node.Data)
				// fmt.Println(node.SQL)
				log.Println("@SendMemberToBag: '", member, "' doesn't contain id! ", id)
				// return err
			}
		}
	}

	return nil
}

func (self *MigrationWorkerV2) SendNodeToBagWithOwnerID(node *DependencyNode, ownerID string) error {
	if !strings.EqualFold(self.mtype, DELETION) {
		return nil
	}
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
	if !strings.EqualFold(self.mtype, DELETION) {
		return nil
	}
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
		// log.Println(fmt.Sprintf("NEXT NODES FETCHED { %s } | nodes [%d]", node.Tag.Name, len(nextNodes)))
		if len(nextNodes) > 0 {
			for _, nextNode := range nextNodes {
				// log.Println(fmt.Sprintf("CURRENT NEXT NODE { %s > %s } %d/%d", node.Tag.Name, nextNode.Tag.Name, i, len(nextNodes)))
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
			// log.Println(fmt.Sprintf("NEXT NODES RETURNING %s", node.Tag.Name))
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
			fmt.Println(refs)
			fmt.Println("#Args: ", self.SrcAppConfig.AppID, ref.fromMember, dependeeMemberID, ref.fromID, ref.toMember, depOnMemberID, ref.toID, fmt.Sprint(self.logTxn.Txn_id), ref.fromAttr, ref.toAttr)
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

func (self *MigrationWorkerV2) MigrateBags(threadID int, isBlade ...bool) error {

	prevIDs := reference_resolution.GetPrevUserIDs(self.SrcAppConfig.AppID, self.uid)
	// prevIDs = append(prevIDs, []string{self.SrcAppConfig.AppID, self.uid})
	prevIDs[self.SrcAppConfig.AppID] = self.uid
	log.Fatal(prevIDs)
	for bagAppID, userID := range prevIDs {

		log.Println(fmt.Sprintf("x%2dx Starting Bags for User: %s App: %s", threadID, userID, bagAppID))

		bags, err := db.GetBagsV2(self.logTxn.DBconn, bagAppID, userID, self.logTxn.Txn_id)

		if err != nil {
			log.Fatal(fmt.Sprintf("x%2dx UNABLE TO FETCH BAGS FOR USER: %s | %s", threadID, self.uid, err))
			return err
		}

		bagWorker := CreateBagWorkerV2(self.uid, bagAppID, self.DstAppConfig.AppID, self.logTxn, BAGS, threadID, isBlade...)

		// log.Fatal(fmt.Sprintf("x%2dx Bag Worker Created | %s -> %s ", threadID, bagWorker.SrcAppConfig.AppName, bagWorker.DstAppConfig.AppName))

		for _, bag := range bags {

			srcMemberID := fmt.Sprint(bag["member"])
			srcMemberName, err := db.TableName(bagWorker.logTxn.DBconn, srcMemberID, bagAppID)

			if err != nil {
				log.Fatal("@MigrateBags > Table Name: ", err)
			}

			bagID := fmt.Sprint(bag["pk"])

			log.Println(fmt.Sprintf("~%2d~ Current    Bag: { %s } | ID: %s, App: %s ", threadID, srcMemberName, bagID, bagAppID))

			bagData := make(map[string]interface{})

			if err := json.Unmarshal(bag["data"].([]byte), &bagData); err != nil {
				fmt.Println("BAG >> ", bag)
				log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO CONVERT BAG TO MAP: %s | %s", threadID, bagWorker.uid, err))
				return err
			}

			bagTag, err := bagWorker.SrcAppConfig.GetTagByMember(srcMemberName)
			if err != nil {
				log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO GET BAG TAG BY MEMBER : %s | %s", threadID, srcMemberName, err))
				return err
			}

			if err := bagWorker.InitTransactions(); err != nil {
				return err
			}

			toCommit := true

			bagNode := DependencyNode{Tag: *bagTag, Data: bagData}
			if err := bagWorker.HandleMigration(&bagNode, true); err != nil {
				toCommit = false
				// fmt.Println(bag)
				log.Println(fmt.Sprintf("x%2dx UNABLE TO MIGRATE BAG { %s } | ID: %s | %s ", threadID, bagTag.Name, bagID, err))
			} else {
				log.Println(fmt.Sprintf("x%2dx MIGRATED bag { %s } | ID: %s", threadID, bagTag.Name, bagID))

				if bagWorker.IsNodeDataEmpty(&bagNode) {
					if err := db.DeleteBagV2(bagWorker.tx.StencilTx, bagID); err != nil {
						toCommit = false
						fmt.Println(bag)
						log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO DELETE BAG : %s | %s ", threadID, bagID, err))
						// return err
					} else {
						log.Println(fmt.Sprintf("x%2dx DELETED bag { %s } ", threadID, bagTag.Name))
					}
				} else {
					if jsonData, err := json.Marshal(bagNode.Data); err == nil {
						if err := db.UpdateBag(bagWorker.tx.StencilTx, bagID, bagWorker.logTxn.Txn_id, jsonData); err != nil {
							toCommit = false
							fmt.Println(bag)
							log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO UPDATE BAG : %s | %s ", threadID, bagID, err))
							// return err
						} else {
							log.Println(fmt.Sprintf("x%2dx UPDATED bag { %s } ", threadID, bagTag.Name))
						}
					} else {
						toCommit = false
						fmt.Println(bagNode.Data)
						log.Fatal(fmt.Sprintf("x%2dx @MigrateBags: UNABLE TO MARSHALL BAG DATA : %s | %s ", threadID, bagID, err))
						// return err
					}
				}
			}

			if toCommit {
				if err := bagWorker.CommitTransactions(); err != nil {
					log.Fatal(fmt.Sprintf("x%2dx UNABLE TO COMMIT bag { %s } | %s ", threadID, bagTag.Name, err))
					// return err
				} else {
					log.Println(fmt.Sprintf("x%2dx COMMITTED bag { %s } ", threadID, bagTag.Name))
				}
			} else {
				log.Println(fmt.Sprintf("x%2dx ROLLBACK bag { %s } | ID: %s ", threadID, bagTag.Name, bagID))
				bagWorker.tx.DstTx.Rollback()
				bagWorker.tx.SrcTx.Rollback()
				bagWorker.tx.StencilTx.Rollback()
			}
		}

		bagWorker.CloseDBConns()
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
		log.Println(fmt.Sprintf("x%2dx | OWNED   node { %s } | root [%s] : owner [%s]", threadID, node.Tag.Name, self.uid, ownerID))
		if err := self.InitTransactions(); err != nil {
			return err
		} else {
			defer self.tx.SrcTx.Rollback()
			defer self.tx.DstTx.Rollback()
			defer self.tx.StencilTx.Rollback()
		}

		log.Println(fmt.Sprintf("x%2dx | CHECKING NEXT NODES { %s }", threadID, node.Tag.Name))

		if err := self.CheckNextNode(node); err != nil {
			return err
		}

		log.Println(fmt.Sprintf("x%2dx | CHECKING PREVIOUS NODES { %s }", threadID, node.Tag.Name))

		if previousNodes, err := self.GetAllPreviousNodes(node); err == nil {
			for _, previousNode := range previousNodes {
				self.AddToReferences(node, previousNode)
			}
		} else {
			return err
		}

		log.Println(fmt.Sprintf("x%2dx | HANDLING MIGRATION { %s }", threadID, node.Tag.Name))

		if err := self.HandleMigration(node, false); err == nil {
			log.Println(fmt.Sprintf("x%2dx MIGRATED  node { %s } ", threadID, node.Tag.Name))
		} else {
			if strings.EqualFold(err.Error(), "3") {
				log.Println(fmt.Sprintf("x%2dx UNMAPPED  node { %s } ", threadID, node.Tag.Name))
			} else if strings.EqualFold(err.Error(), "2") {
				log.Println(fmt.Sprintf("x%2dx Sent2Bag  node { %s } ", threadID, node.Tag.Name))
			} else {
				log.Println(fmt.Sprintf("x%2dx FAILED    node { %s } ", threadID, node.Tag.Name))
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
	fmt.Println("------------------------------------------------------------------------")
	return nil
}

func (self *MigrationWorkerV2) DeletionMigration(node *DependencyNode, threadID int) error {

	if strings.EqualFold(node.Tag.Name, "root") {
		log.Println(fmt.Sprintf("~%2d~ MIGRATING ROOT {%s}", threadID, node.Tag.Name))
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
			log.Println(fmt.Sprintf("~%2d~ Current   Node { %s } | ID: %v ", threadID, node.Tag.Name, node.Data[nodeIDAttr]))
			log.Println(fmt.Sprintf("~%2d~ Adjacent  Node { %s } | ID: %v ", threadID, adjNode.Tag.Name, adjNode.Data[adjNodeIDAttr]))
			if err := self.DeletionMigration(adjNode, threadID); err != nil {
				log.Fatal(fmt.Sprintf("~%2d~ ERROR! NODE { %s } | ID: %v, ADJ_NODE : { %s } | ID: %v | err: [ %s ]", threadID, node.Tag.Name, node.Data[nodeIDAttr], adjNode.Tag.Name, adjNode.Data[adjNodeIDAttr], err))
				return err
			}
		}
	}

	log.Println(fmt.Sprintf("#%2d# | PROCESS Node { %s } ", threadID, node.Tag.Name))

	if strings.EqualFold(node.Tag.Name, "root") {
		return self.DeleteRoot(threadID)
	} else {
		if err := self.CallMigration(node, threadID); err != nil {
			return err
		}
	}

	return nil
}

func (self *MigrationWorkerV2) CallMigrationX(node *DependencyNode, threadID int) error {
	if ownerID, isRoot := self.GetNodeOwner(node); isRoot && len(ownerID) > 0 {
		if err := self.InitTransactions(); err != nil {
			return err
		} else {
			defer self.tx.SrcTx.Rollback()
			defer self.tx.DstTx.Rollback()
			defer self.tx.StencilTx.Rollback()
		}
		if err := self.HandleMigration(node, false); err == nil {
			log.Println(fmt.Sprintf("x%2dx | MIGRATED  node { %s } ", threadID, node.Tag.Name))
		} else {
			log.Println(fmt.Sprintf("x%2dx | RCVD ERR  node { %s } ", threadID, node.Tag.Name), err)
			// if self.unmappedTags.Exists(node.Tag.Name) {
			// 	log.Println(fmt.Sprintf("x%2dx | BREAKLOOP node { %s } ", threadID, node.Tag.Name), err)
			// 	continue
			// }
			if strings.EqualFold(err.Error(), "3") {
				log.Println(fmt.Sprintf("x%2dx | IGNORED   node { %s } ", threadID, node.Tag.Name))
			} else if strings.EqualFold(err.Error(), "2") {
				log.Println(fmt.Sprintf("x%2dx | BAGGED?   node { %s } ", threadID, node.Tag.Name))
			} else {
				log.Println(fmt.Sprintf("x%2dx | FAILED    node { %s } ", threadID, node.Tag.Name), err)
				if strings.EqualFold(err.Error(), "0") {
					log.Println(err)
					return err
				}
				if strings.Contains(err.Error(), "deadlock") {
					return err
				}
			}
		}
		if err := self.CommitTransactions(); err != nil {
			log.Fatal(fmt.Sprintf("x%2dx | UNABEL to COMMIT node { %s } ", threadID, node.Tag.Name))
			return err
		} else {
			log.Println(fmt.Sprintf("x%2dx COMMITTED node { %s } ", threadID, node.Tag.Name))
		}
	} else {
		log.Println(fmt.Sprintf("x%2dx VISITED   node { %s } | root [%s] : owner [%s]", threadID, node.Tag.Name, self.uid, ownerID))
	}
	self.MarkAsVisited(node)
	return nil
}

func (self *MigrationWorkerV2) NaiveMigration(threadID int) error {

	if err := self.CallMigrationX(self.root, threadID); err != nil {
		return err
	}

	for {
		if node, err := self.GetOwnedNode(threadID); err == nil {
			if node == nil {
				break
			}
			nodeIDAttr, _ := node.Tag.ResolveTagAttr("id")
			log.Println(fmt.Sprintf("~%2d~ | Current   Node: { %s } ID: %v", threadID, node.Tag.Name, node.Data[nodeIDAttr]))
			if err := self.CallMigrationX(node, threadID); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if err := self.DeleteRoot(threadID); err != nil {
		log.Println(fmt.Sprintf("~%2d~ | Root not deleted!", threadID))
		log.Fatal(err)
	}

	log.Println("NAIVE MIGRATION DONE!")

	return nil
}

func (self *MigrationWorkerV2) ConsistentMigration(threadID int) error {

	if err := self.CallMigrationX(self.root, threadID); err != nil {
		return err
	}

	for {
		if node, err := self.GetOwnedNode(threadID); err == nil {
			if node == nil {
				return nil
			}
			nodeIDAttr, _ := node.Tag.ResolveTagAttr("id")
			log.Println(fmt.Sprintf("~%2d~ | Current   Node: { %s } ID: %v", threadID, node.Tag.Name, node.Data[nodeIDAttr]))
			if err := self.CallMigrationX(node, threadID); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	log.Println(self.mtype, " MIGRATION DONE!")

	return nil
}

func (self *MigrationWorkerV2) IndependentMigration(threadID int) error {

	return self.ConsistentMigration(threadID)
}

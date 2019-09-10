package migrate

import (
	"errors"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/display"
	"stencil/helper"
	"stencil/qr"
	"stencil/transaction"
	"strconv"
	"github.com/google/uuid"
	"strings"
	"database/sql"
)

func CreateMigrationWorker(uid, srcApp, srcAppID, dstApp, dstAppID string, logTxn *transaction.Log_txn, mtype string) MigrationWorker {
	srcAppConfig, err := config.CreateAppConfig(srcApp, srcAppID)
	if err != nil {
		log.Fatal(err)
	}
	dstAppConfig, err := config.CreateAppConfig(dstApp, dstAppID)
	if err != nil {
		log.Fatal(err)
	}
	mappings := config.GetSchemaMappingsFor(srcApp, dstApp)
	if mappings == nil {
		log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcApp, dstApp))
	}
	dstAppConfig.QR.Migration = true
	srcAppConfig.QR.Migration = true
	mWorker := MigrationWorker{
		uid:          uid,
		SrcAppConfig: srcAppConfig,
		DstAppConfig: dstAppConfig,
		mappings:     mappings,
		wList:        WaitingList{},
		unmappedTags: CreateUnmappedTags(),
		DBConn:       db.GetDBConn(db.STENCIL_DB),
		logTxn:       logTxn,
		mtype:        mtype,
		visitedNodes: make(map[string]bool)}
	if err := mWorker.FetchRoot(); err != nil {
		log.Fatal(err)
	}
	// log.Fatal(mWorker.root.Data)
	return mWorker
}

func (self *MigrationWorker) GetUserBags() ([]map[string]interface{}, error) {
	return db.GetUserBags(self.DBConn, self.uid, self.SrcAppConfig.AppID)
}

func (self *MigrationWorker) RenewDBConn() {
	if self.DBConn != nil {
		self.DBConn.Close()
	}
	self.DBConn = db.GetDBConn(db.STENCIL_DB)
}

func (self *MigrationWorker) Finish() {
	self.DBConn.Close()
}

func (self *MigrationWorker) GetRoot() *DependencyNode {
	return self.root
}

func (self *MigrationWorker) MType() string {
	return self.mtype
}

func (self *MigrationWorker) UserID() string {
	return self.uid
}

func (self *MigrationWorker) MigrationID() int {
	return self.logTxn.Txn_id
}

func (self *MigrationWorker) ResolveDependencyConditions(node *DependencyNode, dep config.Dependency, tag config.Tag, qs *qr.QS) {

	where := qr.CreateQS(self.SrcAppConfig.QR)
	where.TableAliases = qs.TableAliases
	for _, depOn := range dep.DependsOn {
		if depOnTag, err := self.SrcAppConfig.GetTag(depOn.Tag); err == nil {
			if strings.EqualFold(depOnTag.Name, node.Tag.Name) {
				for _, condition := range depOn.Conditions {
					conditionStr := qr.CreateQS(self.SrcAppConfig.QR)
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
						conditionStr.AdditionalWhereWithValue("AND", tagAttr, "=", fmt.Sprint(node.Data[depOnAttr]))
					} else {
						fmt.Println(depOnTag)
						log.Fatal("ResolveDependencyConditions:", depOnAttr, " doesn't exist in ", depOnTag.Name)
					}
					if len(condition.Restrictions) > 0 {
						restrictions := qr.CreateQS(self.SrcAppConfig.QR)
						restrictions.TableAliases = qs.TableAliases
						for _, restriction := range condition.Restrictions {
							if restrictionAttr, err := tag.ResolveTagAttr(restriction["col"]); err == nil {
								restrictions.AdditionalWhereWithValue("OR", restrictionAttr, "=", restriction["val"])
							}

						}
						if restrictions.Where == "" {
							log.Fatal(condition.Restrictions)
						}
						// log.Fatal("restrictions.Where", restrictions.Where)
						conditionStr.AddWhereAsString("AND", restrictions.Where)
					}
					// log.Fatal("conditionStr.Where", conditionStr.Where)
					where.AddWhereAsString("AND", conditionStr.Where)
				}
			}
		}
	}
	// log.Fatal("where.Where", where.Where)
	if where.Where != "" {
		qs.AddWhereAsString("AND", where.Where)
	}
}

func (self *MigrationWorker) ResolveOwnershipConditions(own config.Ownership, tag config.Tag, qs *qr.QS) {

	where := qr.CreateQS(self.SrcAppConfig.QR)
	where.TableAliases = qs.TableAliases
	for _, condition := range own.Conditions {
		conditionStr := qr.CreateQS(self.SrcAppConfig.QR)
		conditionStr.TableAliases = qs.TableAliases
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
			conditionStr.AdditionalWhereWithValue("AND", tagAttr, "=", fmt.Sprint(self.root.Data[depOnAttr]))
		} else {
			fmt.Println("data3", self.root.Data)
			log.Fatal("ResolveOwnershipConditions:", depOnAttr, " doesn't exist in ", tag.Name)
		}
		where.AddWhereAsString("AND", conditionStr.Where)
	}
	if where.Where != "" {
		qs.AddWhereAsString("AND", where.Where)
	}
}

func (self *MigrationWorker) FetchRoot() error {
	tagName := "root"
	if root, err := self.SrcAppConfig.GetTag(tagName); err == nil {
		qs := self.SrcAppConfig.GetTagQS(root)
		rootTable, rootCol := self.SrcAppConfig.GetItemsFromKey(root, "root_id")
		qs.AddWhereWithValue(rootTable+"."+rootCol, "=", self.uid)
		sql := qs.GenSQL()
		// log.Fatal(sql)
		if data, err := db.DataCall1(self.DBConn, sql); err == nil && len(data) > 0 {
			rootNode := new(DependencyNode)
			rootNode.Tag = root
			rootNode.SQL = sql
			rootNode.Data = data
			self.root = rootNode
			// log.Fatal(self.root.Data)
			return nil
		} else {
			if err == nil {
				err = errors.New("no data returned for root node, doesn't exist?")
			}			
			// log.Fatal(sql)
			return err
		}
	} else {
		log.Fatal("Can't fetch root tag:", err)
		return err
	}
}

func (self *MigrationWorker) GetAdjNode(node *DependencyNode, threadID int) (*DependencyNode, error) {

	for _, dep := range self.SrcAppConfig.ShuffleDependencies(self.SrcAppConfig.GetSubDependencies(node.Tag.Name)) {
		if child, err := self.SrcAppConfig.GetTag(dep.Tag); err == nil {
			log.Println(fmt.Sprintf("x%dx | FETCHING  tag  { %s } ", threadID, dep.Tag))
			qs := self.SrcAppConfig.GetTagQS(child)
			self.ResolveDependencyConditions(node, dep, child, qs)
			qs.OrderByFunction("random()")
			qs.ExcludeRowIDs(strings.Join(self.VisitedPKs(), ","))
			sql := qs.GenSQL()
			// log.Fatal(sql)
			if data, err := db.DataCall1(self.DBConn, sql); err == nil {
				if len(data) > 0 {
					newNode := new(DependencyNode)
					newNode.Tag = child
					newNode.SQL = sql
					newNode.Data = data
					if !self.wList.IsAlreadyWaiting(*newNode){
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

func (self *MigrationWorker) GetBagNodes(tagName, bagpks string) ([]*DependencyNode, error) {

	if tag, err := self.SrcAppConfig.GetTag(tagName); err == nil {
		qs := self.SrcAppConfig.GetTagQS(tag)
		qs.RowIDs(bagpks)
		sql := qs.GenSQL()
		log.Fatal(sql)
		if data, err := db.DataCall(self.DBConn, sql); err == nil && len(data) > 0 {
			var bagNodes []*DependencyNode
			for _, datum := range data {
				bagNode := new(DependencyNode)
				bagNode.Tag = tag
				bagNode.SQL = sql
				bagNode.Data = datum
				bagNodes = append(bagNodes, bagNode)
			}
			return bagNodes, nil
		} else {
			log.Println("sql", sql)
			log.Fatal("Problem getting BagNode data:", err, data)
			return nil, err
		}
	} else {
		log.Fatal("Can't fetch bag tag:", tagName, err)
		return nil, err
	}
}

func (self *MigrationWorker) GetOwnedNodes(threadID, limit int) ([]*DependencyNode, error) {

	for _, own := range self.SrcAppConfig.GetShuffledOwnerships() {
		log.Println(fmt.Sprintf("x%dx | FETCHING  tag  { %s } ", threadID, own.Tag))
		if self.unmappedTags.Exists(own.Tag) {
			log.Println(fmt.Sprintf("x%dx | UNMAPPED  tag  { %s } ", threadID, own.Tag))
			continue
		}
		if child, err := self.SrcAppConfig.GetTag(own.Tag); err == nil {
			qs := self.SrcAppConfig.GetTagQS(child)
			self.ResolveOwnershipConditions(own, child, qs)
			qs.OrderByFunction("random()")
			qs.LimitResult(fmt.Sprint(limit))
			sql := qs.GenSQL()
			if result, err := db.DataCall(self.DBConn, sql); err == nil {
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

func (self *MigrationWorker) PushData(stable, dtable, pk string) error {
	if err := display.GenDisplayFlag(self.logTxn.DBconn, self.DstAppConfig.AppName, dtable, pk, false, self.logTxn.Txn_id); err != nil {
		log.Println("## DISPLAY ERROR!", err)
		return errors.New("0")
	}
	if err := db.SaveForEvaluation(self.DBConn, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, stable, dtable, pk, pk, "-", "-", fmt.Sprint(self.logTxn.Txn_id)); err != nil {
		log.Println("## SaveForEvaluation ERROR!", err)
		return errors.New("0")
	}
	return nil
}

func (self *MigrationWorker) VerifyMappingConditions(toTable config.ToTable, node *DependencyNode) bool {

	if len(toTable.Conditions) > 0 {
		for conditionKey, conditionVal := range toTable.Conditions {
			if nodeVal, ok := node.Data[conditionKey]; ok {
				if conditionVal[:1] == "#" {
					fmt.Println("VerifyMappingConditions: conditionVal[:1] == #")
					fmt.Println(conditionKey, conditionVal, nodeVal)
					fmt.Scanln()
					switch conditionVal {
						case "#NULL": {
							if nodeVal != nil {
								log.Println(nodeVal, "!=", conditionVal)
								fmt.Println(conditionKey, conditionVal, nodeVal)
								log.Fatal("@VerifyMappingConditions: return false, from case #NULL:")
								return false
							}
						}
						case "#NOTNULL": {
							if nodeVal == nil {
								log.Println(nodeVal, "!=", conditionVal)
								fmt.Println(conditionKey, conditionVal, nodeVal)
								log.Fatal("@VerifyMappingConditions: return false, from case #NOTNULL:")
								return false
							}		
						}
						default: {
							fmt.Println(toTable.Table, conditionKey, conditionVal)
							log.Fatal("@VerifyMappingConditions: Case not found:" + conditionVal)
						}
					}
				} else if conditionVal[:1] == "$" {
					fmt.Println("VerifyMappingConditions: conditionVal[:1] == $")
					fmt.Println(conditionKey, conditionVal, nodeVal)
					fmt.Scanln()
					if inputVal, err := self.mappings.GetInput(conditionVal); err == nil {
						if !strings.EqualFold(fmt.Sprint(nodeVal), inputVal) {
							log.Println(nodeVal, "!=", inputVal)
							fmt.Println(conditionKey, conditionVal, inputVal, nodeVal)
							log.Fatal("@VerifyMappingConditions: return false, from conditionVal[:1] == $")
							return false
						}
					}else {
						fmt.Println("node data:", node.Data)
						fmt.Println(conditionKey, conditionVal)
						log.Fatal("@VerifyMappingConditions: input doesn't exist?")
					}
				} else {
					if !strings.EqualFold(fmt.Sprint(nodeVal), conditionVal) {
						log.Println(nodeVal, "!=", conditionVal)
						return false
					}
				}
			} else {
				fmt.Println("node data:", node.Data)
				log.Fatal("@VerifyMappingConditions: failed node.Data["+conditionKey+"]")
				return false
			}
		}
	}
	return true
}

func (self *MigrationWorker) CreateMissingData(toTable config.ToTable, node *DependencyNode) map[string]string {
	
	newRows := make(map[string]string)
	for toCol, mappedTabCol := range toTable.Mapping {
		if mappedTabCol[:1] == "$" {
			if inputVal, err := self.mappings.GetInput(mappedTabCol); err == nil {
				newRows[toCol] = inputVal
			} else {
				fmt.Println(toTable.Table, toCol, mappedTabCol)
				log.Fatal("@CreateMissingData: input doesn't exist?")
			}
		} else if mappedTabCol[:1] == "#" {
			if strings.Contains(mappedTabCol, "#ASSIGN"){
				assignedTabCol := strings.Trim(mappedTabCol, "#ASSIGN()")
				if nodeVal, ok := node.Data[assignedTabCol]; ok {
					newRows[toCol] = fmt.Sprint(nodeVal)
				}
			}else{
				switch mappedTabCol {
					case "#GUID": {
						newRows[toCol] = fmt.Sprint(uuid.New())
					}
					case "#RANDINT": {
						newRows[toCol] = fmt.Sprint(self.SrcAppConfig.QR.NewRowId())
					}
					default: {
						fmt.Println(toTable.Table, toCol, mappedTabCol)
						log.Fatal("@CreateMissingData: Case not found:" + mappedTabCol)
					}
				}
			}
		}
	}
	return newRows
}

func (self *MigrationWorker) InsertMissingData(tx *sql.Tx, table, rowid string, data map[string]string) error {
	var cols []string
	var vals []interface{}
	if pk, err := strconv.Atoi(rowid); err == nil {
		for col, val := range data {
			cols, vals = append(cols, col), append(vals, val)
		}
		qi := qr.CreateQI(table, cols, vals, qr.QTInsert)
		qis := self.DstAppConfig.QR.ResolveInsertWithoutRowDesc(qi, int32(pk))
		for _, qi := range qis {
			query, args := qi.GenSQL()
			fmt.Println(query, args)
			if _, err := tx.Exec(query, args...); err != nil {
				fmt.Println("Some error:", err)
				fmt.Println(query, args)
				fmt.Println(qi)
				log.Fatal(err)
				return err
			}
		}
		fmt.Println("InsertMissingData")
		fmt.Scanln()
		return nil
	}else{
		fmt.Println("PK", rowid)
		log.Fatal(err)
		return err
	}
}

func (self *MigrationWorker) IsNodeOwnedByRoot(node *DependencyNode) bool {

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

func (self *MigrationWorker) UpdatePhyRowIDsOfSourceTables(tx *sql.Tx, mapping config.Mapping, node *DependencyNode) error {
	
	if allRowIDs := strings.Split(node.Data["rowids"].(string), ","); len(allRowIDs) > 1 {
		pk     := allRowIDs[0]
		rowids := strings.Join(allRowIDs[1:], ",")
		for _, table := range mapping.FromTables {
			phyTab := self.SrcAppConfig.QR.GetPhyMappingForLogicalTable(table)
			for ptab := range phyTab {
				if err := db.PKReplace(tx, pk, rowids, ptab); err != nil {
					fmt.Println(err)
					fmt.Println(pk, rowids, table, ptab)
					log.Fatal("@UpdatePhyRowIDsOfSourceTables")
					return err
				}
			}
		}
		if err := db.DeleteFromRowDescByRowIDAndAppID(tx, rowids, self.SrcAppConfig.AppID); err != nil{
			fmt.Println(err)
			fmt.Println(pk, rowids)
			log.Fatal("@UpdatePhyRowIDsOfSourceTables: DeleteFromRowDescByRowIDAndAppID")
			return err
		}
		if err := db.PKReplaceRowDesc(tx, pk, rowids); err != nil{
			fmt.Println(err)
			fmt.Println(pk, rowids)
			log.Fatal("@UpdatePhyRowIDsOfSourceTables: PKReplaceRowDesc")
			return err
		}
		for col, val := range node.Data {
			if strings.Contains(col, "pk.") && val != nil {
				node.Data[col] = pk
			}
		}
		node.Data["rowids"] = pk
		// log.Fatal("@UpdatePhyRowIDsOfSourceTables: success")
	}
	return nil
}

func (self *MigrationWorker) HandleMappedMembersOfNode(tx *sql.Tx, mapping config.Mapping, node *DependencyNode) ([]string, error) {
	var updatedPKs []string
	for _, toTable := range mapping.ToTables {
		if self.VerifyMappingConditions(toTable, node) {
			dst_rowid := ""		
			for _, fromTable := range toTable.FromTables() {
				src_rowids := strings.Split(fmt.Sprint(node.Data[fromTable+".rowids_str"]), ",")
				for _, src_rowid := range src_rowids {
					if dst_rowid == "" {
						dst_rowid = src_rowid
					}
					fmt.Println(self.DstAppConfig.AppID, dst_rowid, src_rowid, toTable.Table, self.logTxn.Txn_id)
					cow := "false"
					if self.mtype == INDEPENDENT {
						cow = "true"
					}
					if err := db.InsertIntoMigrationTable(tx, self.DstAppConfig.AppID, dst_rowid, src_rowid, cow, toTable.Table, "1", fmt.Sprint(self.logTxn.Txn_id)); err != nil {
						return nil, err
					}
					if err := db.SaveForEvaluation(self.DBConn, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTable, toTable.Table, src_rowid, dst_rowid, "-", "-", fmt.Sprint(self.logTxn.Txn_id)); err != nil {
						log.Println("## SaveForEvaluation ERROR!", err)
						return nil, errors.New("0")
					}
					updatedPKs = append(updatedPKs, src_rowid)
				}
			}
			if newRow := self.CreateMissingData(toTable, node); len(newRow) > 0 {
				fmt.Println("CreateMissingData:", toTable.Table)
				fmt.Println(newRow)
				fmt.Scanln()
				dst_rowid := fmt.Sprint(node.Data["rowids"])
				self.InsertMissingData(tx, toTable.Table, dst_rowid, newRow)
			}
		}
	}
	return updatedPKs, nil
}

func (self *MigrationWorker) HandleUnmappedMembersOfNode(tx *sql.Tx, mapping config.Mapping, node *DependencyNode) error {

	for _, nodeMember := range node.Tag.GetTagMembers() {
		if !helper.Contains(mapping.FromTables, nodeMember) {
			dst_rowid := ""		
			for _, fromTable := range mapping.FromTables {
				src_rowids := strings.Split(fmt.Sprint(node.Data[fromTable+".rowids_str"]), ",")
				for _, src_rowid := range src_rowids {
					if dst_rowid == "" {
						dst_rowid = src_rowid
					}
					if err := db.NewBag(tx, dst_rowid, src_rowid, self.uid, nodeMember, self.SrcAppConfig.AppID, self.logTxn.Txn_id); err != nil {
						fmt.Println("Args: ", dst_rowid, src_rowid, self.uid, nodeMember, self.SrcAppConfig.AppID, self.logTxn.Txn_id)
						log.Fatal("HandleUnmappedMembersOfNode :: NewBag :", err)
						return err
					}
				}
			}
		}
		
	}
	return nil
}

func (self *MigrationWorker) MigrateNode(mapping config.Mapping, node *DependencyNode) error {

	if tx, err := self.DBConn.Begin(); err != nil {
		log.Println("Can't create MigrateNode transaction!")
		return errors.New("0")
	} else {
		defer tx.Rollback()
		if updatedPKs, err := self.HandleMappedMembersOfNode(tx, mapping, node); err == nil {
			if len(updatedPKs) > 0 {
				if err := self.HandleUnmappedMembersOfNode(tx, mapping, node); err != nil {
					return err
				}
				if undoActionJSON, err := transaction.GenUndoActionJSON(updatedPKs, self.SrcAppConfig.AppID, self.DstAppConfig.AppID); err == nil {
					if log_err := transaction.LogChange(undoActionJSON, self.logTxn); log_err != nil {
						log.Println("MigrateNode: unable to LogChange", log_err)
						return errors.New("0")
					}
				} else {
					log.Fatal("MigrateNode: unable to GenUndoActionJSON", err)
				}
				log.Fatal("killed at commit!")
				// tx.Commit()
			} else {
				return self.HandleUnmappedNode(node)
			}
		}else {
			return err
		}
	}
	return nil
}

func (self *MigrationWorker) HandleWaitingList(appMapping config.Mapping, tagMembers []string, node *DependencyNode) (*DependencyNode, error) {

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

func (self *MigrationWorker) HandleUnmappedTags(node *DependencyNode) error {
	log.Println("!! Couldn't find mappings for the tag [", node.Tag.Name, "]")
	self.unmappedTags.Add(node.Tag.Name)
	// save for evaluation
	for _, tagMember := range node.Tag.Members {
		if _, ok := node.Data[fmt.Sprintf("%s.id", tagMember)]; ok {
			srcID := fmt.Sprint(node.Data[fmt.Sprintf("%s.id", tagMember)])
			if serr := db.SaveForEvaluation(self.DBConn, "diaspora", self.DstAppConfig.AppName, tagMember, "n/a", srcID, "n/a", "*", "n/a", fmt.Sprint(self.logTxn.Txn_id)); serr != nil {
				log.Fatal(serr)
			}
		}
	}
	return errors.New("2")
}

func (self *MigrationWorker) HandleUnmappedNode(node *DependencyNode) error {
	if !strings.EqualFold(self.mtype, DELETION) {
		return errors.New("2")
	}
	if tx, err := self.DBConn.Begin(); err != nil {
		log.Println("Can't create HandleUnmappedNode transaction!")
		return errors.New("0")
	} else {
		var updated []string
		for _, nodeMember := range node.Tag.GetTagMembers() {
			dst_rowid := ""		
			src_rowids := strings.Split(fmt.Sprint(node.Data[nodeMember+".rowids_str"]), ",")
			for _, src_rowid := range src_rowids {
				if dst_rowid == "" {
					dst_rowid = src_rowid
				}
				if err := db.NewBag(tx, dst_rowid, src_rowid, self.uid, nodeMember, self.SrcAppConfig.AppID, self.logTxn.Txn_id); err != nil {
					fmt.Println("Args: ", dst_rowid, src_rowid, self.uid, nodeMember, self.SrcAppConfig.AppID, self.logTxn.Txn_id)
					log.Fatal("HandleUnmappedMembersOfNode :: NewBag :", err)
					return err
				}
				updated = append(updated, src_rowid)
			}
		}
		if undoActionJSON, err := transaction.GenUndoActionJSON(updated, "0", "0"); err == nil {
			if log_err := transaction.LogChange(undoActionJSON, self.logTxn); log_err != nil {
				log.Fatal("HandleUnmappedNode: unable to LogChange", log_err)
			}
		} else {
			log.Fatal("HandleUnmappedNode: unable to GenUndoActionJSON", err)
		}
		tx.Commit()
		return errors.New("2")
	}
}

func (self *MigrationWorker) HandleMigration(node *DependencyNode, isBag bool) error {

	for _, mapping := range self.mappings.Mappings {
		tagMembers := node.Tag.GetTagMembers()
		if mappedTables := helper.IntersectString(tagMembers, mapping.FromTables); len(mappedTables) > 0 {
			if helper.Sublist(tagMembers, mapping.FromTables) { // other mappings HANDLE!
				return self.MigrateNode(mapping, node)
			}
			if wNode, err := self.HandleWaitingList(mapping, tagMembers, node); wNode != nil && err == nil {
				return self.MigrateNode(mapping, wNode)
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

func (self *MigrationWorker) HandleLeftOverWaitingNodes() {

	for _, waitingNode := range self.wList.Nodes {
		for _, containedNode := range waitingNode.ContainedNodes {
			self.HandleUnmappedNode(containedNode)
		}
	}
}

func (self *MigrationWorker) IsVisited(node *DependencyNode) bool {
	rowids := strings.Split(node.Data["rowids"].(string), ",")
	for _, pk := range rowids {
		if _, ok := self.visitedNodes[pk]; ok {
			return true
		}
	}
	return false
}

func (self *MigrationWorker) MarkAsVisited(node *DependencyNode) {
	rowids := strings.Split(node.Data["rowids"].(string), ",")
	for _, pk := range rowids {
		self.visitedNodes[pk] = true
	}
}

func (self *MigrationWorker) VisitedPKs() []string {
	var pks []string
	for pk := range self.visitedNodes {
		pks = append(pks, pk)
	}
	return pks
}

func (self *MigrationWorker) DeletionMigration(node *DependencyNode, threadID int) error {

	if strings.EqualFold(node.Tag.Name, "root") && !db.CheckUserInApp(self.uid, self.DstAppConfig.AppID, self.DBConn) {
		log.Println("++ Adding User from ", self.SrcAppConfig.AppName, " to ", self.DstAppConfig.AppName)
		db.AddUserToApp(self.uid, self.DstAppConfig.AppID, self.DBConn)
	}

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
		if err := self.HandleMigration(node, false); err == nil {
			log.Println(fmt.Sprintf("x%dx MIGRATED  node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
		} else {
			if strings.EqualFold(err.Error(), "2") {
				log.Println(fmt.Sprintf("x%dx BAGGED    node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
			} else {
				log.Println(fmt.Sprintf("x%dx FAILED    node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
				if strings.EqualFold(err.Error(), "0") {
					log.Println(err)
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

func (self *MigrationWorker) SecondPhase(threadID int) error {

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

func (self *MigrationWorker) RegisterMigration(mtype string, number_of_threads int) bool {
	return db.RegisterMigration(self.uid, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, mtype, self.logTxn.Txn_id, number_of_threads, self.DBConn, false)
	// db.DeleteExistingMigrationRegistrations(self.uid, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, self.DBConn)
	// if !db.CheckMigrationRegistration(self.uid, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, self.DBConn) {
	// 	return db.RegisterMigration(self.uid, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, mtype, self.logTxn.Txn_id, number_of_threads, self.DBConn)
	// } else {
	// 	log.Println("Migration Already Registered!")
	// 	return true
	// }
}

func (self *MigrationWorker) FinishMigration(mtype string, number_of_threads int) bool {
	return db.FinishMigration(self.logTxn.DBconn, self.logTxn.Txn_id)
}

func (self *MigrationWorker) MigrateProcessBags(bag map[string]interface{}) error {

	if bagNodes, err := self.GetBagNodes(fmt.Sprint(bag["tag"]), fmt.Sprint(bag["rowids"])); err != nil {
		log.Fatal(err)
		return nil
	} else {
		for _, bagNode := range bagNodes {
			if err := self.HandleMigration(bagNode, true); err == nil {
				if err := db.DeleteBagsByRowIDS(self.DBConn, fmt.Sprint(bag["rowids"])); err != nil {
					log.Println(err)
					return err
				}
				return nil
			} else if strings.EqualFold(err.Error(), "0") {
				log.Println(err)
				return err
			}
		}
		return nil
	}
}

func (self *MigrationWorker) ConsistentMigration(threadID int) error {
	
	nodelimit := 100
	for nodes, err := self.GetOwnedNodes(threadID, nodelimit); err != nil || nodes != nil; nodes, err = self.GetOwnedNodes(threadID, nodelimit) {
		if err != nil {
			return err
		}
		for _, node := range nodes {
			nodeIDAttr, _ := node.Tag.ResolveTagAttr("id")
			log.Println(fmt.Sprintf("~%d~ | Current   Node: { %s } ID: %v", threadID, node.Tag.Name, node.Data[nodeIDAttr]))
			if err := self.HandleMigration(node, false); err == nil {
				log.Println(fmt.Sprintf("x%dx | MIGRATED  node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
			} else {
				log.Println(fmt.Sprintf("x%dx | RCVD ERR  node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName), err)
				if self.unmappedTags.Exists(node.Tag.Name) {
					log.Println(fmt.Sprintf("x%dx | BREAKLOOP node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName), err)
					break
				}
				if strings.EqualFold(err.Error(), "2") {
					log.Println(fmt.Sprintf("x%dx | IGNORED   node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
				} else {
					log.Println(fmt.Sprintf("x%dx | FAILED    node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
					if strings.EqualFold(err.Error(), "0") {
						log.Println(err)
						return err
					}
				}
			}
		}
	}
	return nil
}

func (self *MigrationWorker) IndependentMigration(threadID int) error {

	return self.ConsistentMigration(threadID)
}

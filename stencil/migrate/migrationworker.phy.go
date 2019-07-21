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
	"strings"
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
	mWorker := MigrationWorker{
		uid:          uid,
		SrcAppConfig: srcAppConfig,
		DstAppConfig: dstAppConfig,
		mappings:     mappings,
		wList:        WaitingList{},
		unmappedTags: CreateUnmappedTags(),
		DBConn:       db.GetDBConn("stencil"),
		logTxn:       logTxn,
		mtype:        mtype}
	if err := mWorker.FetchRoot(); err != nil {
		log.Fatal(err)
	}
	return mWorker
}

func (self *MigrationWorker) GetUserBags() ([]map[string]interface{}, error) {
	return db.GetUserBags(self.DBConn, self.uid, self.SrcAppConfig.AppID)
}

func (self *MigrationWorker) RenewDBConn() {
	if self.DBConn != nil {
		self.DBConn.Close()
	}
	self.DBConn = db.GetDBConn("stencil")
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
						conditionStr.WhereOperatorInterface("AND", tagAttr, "=", node.Data[depOnAttr])
					} else {
						fmt.Println(depOnTag)
						log.Fatal("ResolveDependencyConditions:", depOnAttr, " doesn't exist in ", depOnTag.Name)
					}
					if len(condition.Restrictions) > 0 {
						restrictions := qr.CreateQS(self.SrcAppConfig.QR)
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

func (self *MigrationWorker) FetchRoot() error {
	tagName := "root"
	if root, err := self.SrcAppConfig.GetTag(tagName); err == nil {
		qs := self.SrcAppConfig.GetTagQS(root)
		rootTable, rootCol := self.SrcAppConfig.GetItemsFromKey(root, "root_id")
		qs.WhereSimpleVal(rootTable+"."+rootCol, "=", self.uid)
		qs.WhereMFlag(qr.EXISTS, "0", self.SrcAppConfig.AppID)
		sql := qs.GenSQL()
		if data, err := db.DataCall1(self.DBConn, sql); err == nil && len(data) > 0 {
			rootNode := new(DependencyNode)
			rootNode.Tag = root
			rootNode.SQL = sql
			rootNode.Data = data
			self.root = rootNode
			return nil
		} else {
			log.Fatal("Problem getting RootNode data:", err, data)
			return err
		}
	} else {
		log.Fatal("Can't fetch root tag:", err)
		return err
	}
}

func (self *MigrationWorker) GetAdjNode(node *DependencyNode) (*DependencyNode, error) {

	for _, dep := range config.ShuffleDependencies(self.SrcAppConfig.GetSubDependencies(node.Tag.Name)) {
		if child, err := self.SrcAppConfig.GetTag(dep.Tag); err == nil {
			qs := self.SrcAppConfig.GetTagQS(child)
			self.ResolveDependencyConditions(node, dep, child, qs)
			qs.WhereMFlag(qr.EXISTS, "0", self.SrcAppConfig.AppID)
			qs.OrderBy("random()")
			sql := qs.GenSQL()
			if data, err := db.DataCall1(self.DBConn, sql); err == nil {
				if len(data) > 0 {
					newNode := new(DependencyNode)
					newNode.Tag = child
					newNode.SQL = sql
					newNode.Data = data
					if !self.wList.IsAlreadyWaiting(*newNode) {
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
		sql := qs.GenSQLWith(bagpks)
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

func (self *MigrationWorker) CheckMappingConditions(toTable config.ToTable, node *DependencyNode) bool {
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

func (self *MigrationWorker) UpdateRowDesc(toTables []config.ToTable, node *DependencyNode) error {

	if tx, err := self.DBConn.Begin(); err != nil {
		log.Println("Can't create UpdateRowDesc transaction!")
		return errors.New("0")
	} else {
		// var errs []error
		var updated []string

		for _, toTable := range toTables {
			if self.CheckMappingConditions(toTable, node) {
				continue
			}
			for col, val := range node.Data {
				if strings.Contains(col, "pk.") && val != nil {
					pk := strconv.FormatInt(val.(int64), 10)
					if val != nil && !helper.Contains(updated, pk) {
						if err := db.MUpdate(tx, pk, "1", self.DstAppConfig.AppID); err == nil {
							updated = append(updated, pk)
							if err := display.GenDisplayFlag(self.logTxn.DBconn, self.DstAppConfig.AppName, toTable.Table, pk, false, self.logTxn.Txn_id); err != nil {
								log.Fatal("## DISPLAY ERROR!", err)
							}
						} else {
							tx.Rollback()
							fmt.Println("\n@ERROR_MUpdate:", err)
							log.Fatal("pk:", pk, "appid", self.DstAppConfig.AppID)
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
		if undoActionJSON, err := transaction.GenUndoActionJSON(updated, self.SrcAppConfig.AppID, self.DstAppConfig.AppID); err == nil {
			if log_err := transaction.LogChange(undoActionJSON, self.logTxn); log_err != nil {
				log.Fatal("UpdateRowDesc: unable to LogChange", log_err)
			}
		} else {
			log.Fatal("UpdateRowDesc: unable to GenUndoActionJSON", err)
		}
		tx.Commit()
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
			tempCombinedDataDependencyNode := waitingNode.GenDependencyDataNode()
			return tempCombinedDataDependencyNode, nil
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
	if tx, err := self.DBConn.Begin(); err != nil {
		log.Println("Can't create UpdateRowDesc transaction!")
		return errors.New("0")
	} else {
		var updated []string
		for col, val := range node.Data {
			if strings.Contains(col, "pk.") {
				pk := fmt.Sprint(val)
				if val != nil && !helper.Contains(updated, pk) {
					if err := db.SetMFlag(tx, pk, "1"); err == nil {
						if bag_err := db.NewBag(tx, pk, self.uid, node.Tag.Name, self.logTxn.Txn_id); bag_err == nil {
							updated = append(updated, pk)
						} else {
							tx.Rollback()
							// log.Fatal("\n@ERROR_New_Bag:", bag_err)
							return err
						}
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

func (self *MigrationWorker) MigrateNode(node *DependencyNode, isBag bool) error {

	for _, appMapping := range self.mappings.Mappings {
		tagMembers := node.Tag.GetTagMembers()
		if mappedTables := helper.IntersectString(tagMembers, appMapping.FromTables); len(mappedTables) > 0 {
			if helper.Sublist(tagMembers, appMapping.FromTables) {
				return self.UpdateRowDesc(appMapping.ToTables, node)
			}
			if wNode, err := self.HandleWaitingList(appMapping, tagMembers, node); wNode != nil && err == nil {
				return self.UpdateRowDesc(appMapping.ToTables, wNode)
			} else {
				return err
			}
		}
	}
	if isBag {
		return fmt.Errorf("no mapping found for bag: %s", node.Tag.Name)
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

func (self *MigrationWorker) DeletionMigration(node *DependencyNode, threadID int) error {

	if strings.EqualFold(node.Tag.Name, "root") && !db.CheckUserInApp(self.uid, self.DstAppConfig.AppID, self.DBConn) {
		log.Println("++ Adding User from ", self.SrcAppConfig.AppName, " to ", self.DstAppConfig.AppName)
		db.AddUserToApp(self.uid, self.DstAppConfig.AppID, self.DBConn)
	}

	for child, err := self.GetAdjNode(node); child != nil; child, err = self.GetAdjNode(node) {
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
	if err := self.MigrateNode(node, false); err == nil {
		log.Println(fmt.Sprintf("x%dx MIGRATED  node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
	} else {
		if strings.EqualFold(err.Error(), "2") {
			log.Println(fmt.Sprintf("x%dx IGNORED   node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
		} else {
			log.Println(fmt.Sprintf("x%dx FAILED    node { %s } From [%s] to [%s]", threadID, node.Tag.Name, self.SrcAppConfig.AppName, self.DstAppConfig.AppName))
			if strings.EqualFold(err.Error(), "0") {
				log.Println(err)
				return err
			}
		}
	}

	fmt.Println("------------------------------------------------------------------------")

	return nil
}

func (self *MigrationWorker) RegisterMigration(mtype string) bool {
	db.DeleteExistingMigrationRegistrations(self.uid, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, self.DBConn)
	if !db.CheckMigrationRegistration(self.uid, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, self.DBConn) {
		return db.RegisterMigration(self.uid, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, mtype, self.logTxn.Txn_id, self.DBConn)
	} else {
		log.Println("Migration Already Registered!")
		return true
	}
}

func (self *MigrationWorker) MigrateProcessBags(bag map[string]interface{}) error {

	if bagNodes, err := self.GetBagNodes(fmt.Sprint(bag["tag"]), fmt.Sprint(bag["rowids"])); err != nil {
		log.Fatal(err)
		return nil
	} else {
		for _, bagNode := range bagNodes {
			if err := self.MigrateNode(bagNode, true); err == nil {
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

func (self *MigrationWorker) ConsistentMigration(node *DependencyNode, threadID int) error {

	return nil
}

func (self *MigrationWorker) IndependentMigration(node *DependencyNode, threadID int) error {

	return nil
}

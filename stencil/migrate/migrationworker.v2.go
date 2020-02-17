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
	"stencil/transaction"
	"strings"

	"github.com/google/uuid"
	"github.com/gookit/color"
)

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
				err = errors.New("Can't fetch Root node. Check if it exists. UID: " + self.uid)
			} else {
				fmt.Println("@FetchRoot > DataCall1 | ", err)
			}
			return err
		}
	} else {
		self.Logger.Fatal("Can't fetch root tag:", err)
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
					self.Logger.Fatal("@GetAllNextNodes: Error while DataCall: ", err)
					return nil, err
				}
			} else {
				log.Println("@GetAllNextNodes > ResolveDependencyConditions | ", err)
			}
		} else {
			self.Logger.Fatal("@GetAllNextNodes: Tag doesn't exist? ", dep.Tag)
		}
	}
	// if len(self.SrcAppConfig.GetSubDependencies(node.Tag.Name)) > 0 {
	// 	log.Println("@GetAllNextNodes:", len(nodes))
	// 	self.Logger.Fatal(nodes)
	// }
	return nodes, nil
}

func (self *MigrationWorkerV2) GetAllPreviousNodes(node *DependencyNode) ([]*DependencyNode, error) {
	var nodes []*DependencyNode

	if node.Tag.Name != "root" {
		if ownership := self.SrcAppConfig.GetOwnership(node.Tag.Name, "root"); ownership != nil {
			if where, err := node.ResolveParentOwnershipConditions(ownership, self.root.Tag); err == nil {
				ql := self.GetTagQL(self.root.Tag)
				sql := fmt.Sprintf("%s WHERE %s ", ql, where)
				sql += self.root.Tag.ResolveRestrictions()
				if data, err := db.DataCall(self.SrcDBConn, sql); err == nil {
					for _, datum := range data {
						newNode := new(DependencyNode)
						newNode.Tag = self.root.Tag
						newNode.SQL = sql
						newNode.Data = datum
						nodes = append(nodes, newNode)
					}
				} else {
					fmt.Println(sql)
					self.Logger.Fatal("@GetAllPreviousNodes: Error while DataCall: ", err)
					return nil, err
				}
			} else {
				log.Println("@GetAllPreviousNodes > ResolveParentOwnershipConditions: ", err)
			}
		} else {
			self.Logger.Fatal("@GetAllPreviousNodes: Ownership doesn't exist? ", node.Tag.Name, "root")
		}
	}

	// for _, dep := range self.SrcAppConfig.GetParentDependencies(node.Tag.Name) {
	// 	for _, pdep := range dep.DependsOn {
	// 		if parent, err := self.SrcAppConfig.GetTag(pdep.Tag); err == nil {
	// 			if where, err := node.ResolveParentDependencyConditions(pdep.Conditions, parent); err == nil {
	// 				ql := self.GetTagQL(parent)
	// 				sql := fmt.Sprintf("%s WHERE %s ", ql, where)
	// 				sql += parent.ResolveRestrictions()
	// 				// fmt.Println(node.SQL)
	// 				// self.Logger.Fatal("@GetAllPreviousNodes | ", sql)
	// 				if data, err := db.DataCall(self.SrcDBConn, sql); err == nil {
	// 					for _, datum := range data {
	// 						newNode := new(DependencyNode)
	// 						newNode.Tag = parent
	// 						newNode.SQL = sql
	// 						newNode.Data = datum
	// 						nodes = append(nodes, newNode)
	// 					}
	// 				} else {
	// 					fmt.Println(sql)
	// 					self.Logger.Fatal("@GetAllPreviousNodes: Error while DataCall: ", err)
	// 					return nil, err
	// 				}
	// 			} else {
	// 				log.Println("@GetAllPreviousNodes > ResolveParentDependencyConditions: ", err)
	// 			}
	// 		} else {
	// 			self.Logger.Fatal("@GetAllPreviousNodes: Tag doesn't exist? ", pdep.Tag)
	// 		}
	// 	}
	// }
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
				// self.Logger.Fatal(sql)
				if data, err := db.DataCall1(self.SrcDBConn, sql); err == nil {
					if len(data) > 0 {
						newNode := DependencyNode{Tag: child, SQL: sql, Data: data}
						if !self.wList.IsAlreadyWaiting(newNode) && !self.IsVisited(&newNode) {
							return &newNode, nil
						}
					}
				} else {
					fmt.Println("@GetDependentNode > DataCall1 | ", err)
					self.Logger.Fatal(sql)
					return nil, err
				}
			} else {
				log.Println("@GetDependentNode > ResolveDependencyConditions | ", err)
				// self.Logger.Fatal(err)
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
				// self.Logger.Fatal(sql)
				if data, err := db.DataCall1(self.SrcDBConn, sql); err == nil {
					if len(data) > 0 {
						newNode := DependencyNode{Tag: child, SQL: sql, Data: data}
						if !self.wList.IsAlreadyWaiting(newNode) {
							return &newNode, nil
						}
					}
				} else {
					fmt.Println("@GetOwnedNode > DataCall1 | ", err)
					self.Logger.Fatal(sql)
					return nil, err
				}
			} else {
				log.Println("@GetOwnedNode > ResolveOwnershipConditions | ", err)
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
		self.Logger.Fatal("## DISPLAY ERROR!", err)
		return errors.New("0")
	}

	for fromTable, fromCols := range mappedData.srcTables {
		if _, ok := node.Data[fmt.Sprintf("%s.id", fromTable)]; ok {
			srcID := node.Data[fmt.Sprintf("%s.id", fromTable)]
			if fromTableID, err := db.TableID(self.logTxn.DBconn, fromTable, self.SrcAppConfig.AppID); err == nil {
				// if err := db.InsertIntoIdentityTable(tx, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTableID, dtable.TableID, srcID, pk, fmt.Sprint(self.logTxn.Txn_id)); err != nil {
				// 	log.Println("@PushData:db.InsertIntoIdentityTable: ", self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTableID, dtable.TableID, srcID, pk, fmt.Sprint(self.logTxn.Txn_id))
				// 	self.Logger.Fatal(err)
				// 	return errors.New("0")
				// }
				if serr := db.SaveForLEvaluation(tx, self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTableID, dtable.TableID, srcID, pk, strings.Join(fromCols, ","), mappedData.cols, fmt.Sprint(self.logTxn.Txn_id)); serr != nil {
					log.Println("@PushData:db.SaveForLEvaluation: ", self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTableID, dtable.TableID, srcID, pk, strings.Join(fromCols, ","), mappedData.cols, fmt.Sprint(self.logTxn.Txn_id))
					self.Logger.Fatal(serr)
					return errors.New("0")
				}
			} else {
				log.Println("@PushData:db.TableID: ", fromTable, self.SrcAppConfig.AppID)
				self.Logger.Fatal(err)
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
								// self.Logger.Fatal("@VerifyMappingConditions: return false, from case #NULL:")
								return false
							}
						}
					case "#NOTNULL":
						{
							if nodeVal == nil {
								// log.Println(nodeVal, "!=", conditionVal)
								// fmt.Println(conditionKey, conditionVal, nodeVal)
								// self.Logger.Fatal("@VerifyMappingConditions: return false, from case #NOTNULL:")
								return false
							}
						}
					default:
						{
							fmt.Println(toTable.Table, conditionKey, conditionVal)
							self.Logger.Fatal("@ValidateMappingConditions: Case not found:" + conditionVal)
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
							self.Logger.Fatal("@ValidateMappingConditions: return false, from conditionVal[:1] == $")
							return false
						}
					} else {
						fmt.Println("node data:", node.Data)
						fmt.Println(conditionKey, conditionVal)
						self.Logger.Fatal("@ValidateMappingConditions: input doesn't exist?", err)
					}
				} else if !strings.EqualFold(fmt.Sprint(nodeVal), conditionVal) {
					// log.Println(conditionKey, conditionVal, "!=", nodeVal)
					return false
				} else {
					// fmt.Println(*nodeVal, "==", conditionVal)
				}
			} else {
				// log.Println("Condition Key", conditionKey, "doesn't exist!")
				// fmt.Println("node data:", node.Data)
				// fmt.Println("node sql:", node.SQL)
				// self.Logger.Fatal("@ValidateMappingConditions: stop here and check")
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
				self.Logger.Fatal("@GetNodeOwner: Resolving TagAttr", err, node.Tag.Name, condition.TagAttr)
				break
			}
			depOnAttr, err := self.root.Tag.ResolveTagAttr(condition.DependsOnAttr)
			if err != nil {
				self.Logger.Fatal("@GetNodeOwner: Resolving depOnAttr", err, node.Tag.Name, condition.DependsOnAttr)
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
					self.Logger.Fatal("@GetNodeOwner: stop here and check ownership conditions wrt root")
				}
			} else {
				fmt.Println("@GetNodeOwner: Ownership Condition Key", tagAttr, "doesn't exist!")
				fmt.Println("@GetNodeOwner: node data:", node.Data)
				fmt.Println("@GetNodeOwner: node sql:", node.SQL)
				self.Logger.Fatal("@GetNodeOwner: stop here and check ownership conditions")
			}
		}
	} else {
		self.Logger.Fatal("@GetNodeOwner: Ownership not found:", node.Tag.Name)
	}
	return "", false
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
			self.Logger.Fatal("@FetchFromMapping: FetchForMapping | ", err)
			return err
		} else if len(res) > 0 {
			data.UpdateData(toAttr, args[0], targetTabCol[0], res[targetTabCol[1]])
			node.Data[args[0]] = res[targetTabCol[1]]
			if len(args) > 3 {
				toMemberTokens := strings.Split(args[3], ".")
				data.UpdateRefs(res[targetTabCol[1]], targetTabCol[0], targetTabCol[1], res[targetTabCol[1]], toMemberTokens[0], toMemberTokens[1])
			}
		} else {
			log.Println("@FetchFromMapping: FetchForMapping | Returned data is nil! Previous node already migrated?", res, targetTabCol[0], targetTabCol[1], comparisonTabCol[1], fmt.Sprint(nodeVal))
		}
	} else {
		fmt.Println(node.Tag.Name, node.Data)
		self.Logger.Fatal("@FetchFromMapping: unable to fetch ", args[2])
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

func (self *MigrationWorkerV2) IsNodeDataEmpty(data map[string]interface{}) bool {
	for key := range data {
		if !(strings.Contains(key, ".id") || strings.Contains(key, ".display_flag")) {
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
						var fromID interface{}
						if val, ok := node.Data[assignedTabColTokens[0]+".id"]; ok {
							fromID = val
						} else {
							fmt.Println(assignedTabColTokens[0], " | ", assignedTabColTokens)
							fmt.Println(node.Data)
							self.Logger.Fatal("@GetMappedData > #REF #ASSIGN> fromID: Unable to find ref value in node data")
							return data, errors.New("Unable to find ref value in node data")
						}
						data.UpdateRefs(fromID, assignedTabColTokens[0], assignedTabColTokens[1], nodeVal, referredTabColTokens[0], referredTabColTokens[1])
					}
				} else {
					args := strings.Split(assignedTabCol, ",")
					if nodeVal, ok := node.Data[args[0]]; ok {
						argsTokens := strings.Split(args[0], ".")
						data.UpdateData(toAttr, args[0], argsTokens[0], nodeVal)
					}
					if self.mtype != BAGS {
						if toID, fromID, err := data.GetIDs(args[0], node.Data); err == nil {
							secondMemberTokens := strings.Split(args[1], ".")
							firstMemberTokens := strings.Split(args[0], ".")
							data.UpdateRefs(fromID, firstMemberTokens[0], firstMemberTokens[1], toID, secondMemberTokens[0], secondMemberTokens[1])
						} else {
							fmt.Printf("args[0]: '%v' \n", args[0])
							fmt.Printf("toID: '%v' | fromID: '%v' \n", toID, fromID)
							self.Logger.Fatal("@GetMappedData > GetIDs | ", err)
							return data, err
						}
					}
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
					self.Logger.Fatal("@GetMappedData > #FETCH > FetchFromMapping: Unable to fetch")
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
						self.Logger.Fatal("@GetMappedData: Case not found:" + fromAttr)
					}
				}
			}
			// self.Logger.Fatal(fromAttr)
		} else {
			data.orgColsLeft += fmt.Sprintf("%s,", fromAttr)
		}
	}
	data.undoAction.AddDstTable(toTable.Table)
	data.Trim(", ")
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
				// self.Logger.Fatal(derr)
				return derr
			}
			if tagMemberID, err := db.TableID(self.logTxn.DBconn, tagMember, self.SrcAppConfig.AppID); err == nil {
				if derr := db.UpdateLEvaluation(self.logTxn.DBconn, tagMemberID, srcID, self.logTxn.Txn_id); derr != nil {
					fmt.Println("@ERROR_UpdateLEvaluation", derr)
					fmt.Println("@QARGS:", tagMember, srcID, self.logTxn.Txn_id)
					self.Logger.Fatal(derr)
					return derr
				}
			} else {
				self.Logger.Fatal("@DeleteRow>TableID: ", err)
			}

		} else {
			// log.Println("node.Data =>", node.Data)
			self.Logger.Info("@DeleteRow:  '", idCol, "' not present or is null in node data!", nodeVal)
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
			if err := self.SendMemberToBag(node, nodeMember, self.uid); err != nil {
				return err
			}
		}
	}
	return nil
}

func (self *MigrationWorkerV2) DeleteNode(mapping config.Mapping, node *DependencyNode) error {

	if self.mtype == DELETION {

		if !self.IsNodeDataEmpty(node.Data) {
			if err := self.SendNodeToBagWithOwnerID(node, self.uid); err != nil {
				fmt.Println(node.Tag.Name)
				fmt.Println(node.Data)
				self.Logger.Fatal("@DeleteNode > SendNodeToBagWithOwnerID:", err)
				return err
			}
		}

		if err := self.DeleteRow(node); err != nil {
			fmt.Println(node.Tag.Name)
			fmt.Println(node)
			self.Logger.Fatal("@DeleteNode > DeleteRow:", err)
			return err
		} else {
			log.Println(fmt.Sprintf("%s node { %s }", color.FgRed.Render("Deleted"), node.Tag.Name))
		}
	}

	return nil
}

func (self *MigrationWorkerV2) DeleteRoot(threadID int) error {

	if err := self.InitTransactions(); err != nil {
		self.Logger.Fatal("@DeleteRoot > InitTransactions", err)
		return err
	} else {
		defer self.tx.SrcTx.Rollback()
		defer self.tx.DstTx.Rollback()
		defer self.tx.StencilTx.Rollback()
	}
	if mapping, found := self.FetchMappingsForNode(self.root); found {
		if self.mtype == NAIVE {
			if err := self.DeleteRow(self.root); err != nil {
				self.Logger.Fatal("@DeleteRoot:", err)
				return err
			}
		} else if self.mtype == DELETION {
			if err := self.DeleteNode(mapping, self.root); err != nil {
				self.Logger.Fatal("@DeleteRoot:", err)
				return err
			}
		} else {
			self.Logger.Fatal("ATTEMPTED DELETION IN DISALLOWED MIGRATION TYPE!")
		}
	} else {
		fmt.Println(self.root)
		self.Logger.Fatal("@DeleteRoot: Can't find mappings for root | ", mapping, found)
	}
	if err := self.CommitTransactions(); err != nil {
		self.Logger.Fatal("@DeleteRoot: ERROR COMMITING TRANSACTIONS! ")
		return err
	}
	return nil
}

func (self *MigrationWorkerV2) MigrateNode(mapping config.Mapping, node *DependencyNode) (bool, error) {

	migrated := false

	var allMappedData []MappedData

	for _, toTable := range mapping.ToTables {

		if !self.ValidateMappingConditions(toTable, node) {
			continue
		}
		if mappedData, mappedDataErr := self.GetMappedData(toTable, node); mappedDataErr != nil {
			self.Logger.Info(node.Data)
			self.Logger.Info(mappedData)
			self.Logger.Fatal("@MigrateNode > GetMappedData Error | ", mappedDataErr)
		} else if len(mappedData.cols) > 0 && len(mappedData.vals) > 0 && len(mappedData.ivals) > 0 {
			if !self.ValidateMappedTableData(toTable, mappedData) {
				self.Logger.Warn("@MigrateNode > ValidateMappedTableData: Validation failed!")
				continue
			}

			if self.mtype == DELETION || self.mtype == BAGS {
				self.Logger.Infof("Before Merging Data | %s\n%v\n---", toTable.Table, mappedData)
				if err := self.MergeBagDataWithMappedData(&mappedData, node, toTable); err != nil {
					self.Logger.Fatal("@MigrateNode > MergeDataFromBagsWithMappedData | ", err)
				}
				self.Logger.Infof("After Merging Data | %s\n%v\n---", toTable.Table, mappedData)
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
								self.Logger.Fatal(err)
								return migrated, err
							}
						} else {
							fmt.Println(node.Data)
							self.Logger.Fatal("@MigrateNode: InsertIntoIdentityTable | " + fromTable + ".id doesn't exist")
						}
					} else {
						self.Logger.Fatal("@MigrateNode > TableID, fromTable: error in getting table id for member! ", fromTable, err)
						return migrated, err
					}
				}

				if err := self.PushData(self.tx.StencilTx, toTable, fmt.Sprint(id), mappedData, node); err != nil {
					self.Logger.Warn("@MigrateNode")
					self.Logger.Warn("@Params:", toTable.Table, fmt.Sprint(id), mappedData.orgCols, mappedData.cols, mappedData.undoAction, node)
					self.Logger.Fatal(err)
					return migrated, err
				}

				if len(toTable.Media) > 0 {
					if filePathCol, ok := toTable.Media["path"]; ok {
						if filePath, ok := node.Data[filePathCol]; ok {
							if err := self.TransferMedia(fmt.Sprint(filePath)); err != nil {
								self.Logger.Fatal("@MigrateNode > TransferMedia: ", err)
							}
						}
					} else {
						self.Logger.Fatal("@MigrateNode > toTable.Media: Path not found in map!")
					}
				}
				allMappedData = append(allMappedData, mappedData)
			} else {
				self.Logger.Warn(fmt.Sprintf("@Args: [toTable: %s], [cols: %s], [vals: %s], [ivals: %v], [srcTables: %s]", toTable.Table, mappedData.cols, mappedData.vals, mappedData.ivals, mappedData.srcTables))
				self.Logger.Warn("@NODE: ", node.Tag.Name, "|", node.Data)
				self.Logger.Error("@MigrateNode > InsertRowIntoAppDB: ", err)
				if self.mtype == DELETION {
					self.Logger.Fatal("@MigrateNode > InsertRowIntoAppDB: STOP HERE!")
				}
				return migrated, err
			}

			if self.mtype != BAGS {
				if err := self.AddMappedReferences(mappedData.refs); err != nil {
					log.Println(mappedData.refs)
					self.Logger.Fatal("@MigrateNode > AddMappedReferences: ", err)
					return migrated, err
				}
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

	if len(allMappedData) > 0 {
		migrated = true
		self.Logger.Info("Migrated Data:\n", allMappedData)
		for _, mappedData := range allMappedData {
			self.RemoveMappedDataFromNodeData(mappedData, node)
		}
		if self.mtype == BAGS {
			if err := self.DeleteBag(node); err != nil {
				self.Logger.Fatal("@MigrateNode > DeleteBag:", err)
				return false, err
			}
		}
	}

	if !strings.EqualFold(node.Tag.Name, "root") {
		switch self.mtype {
		case DELETION:
			{
				if err := self.DeleteNode(mapping, node); err != nil {
					self.Logger.Fatal("@MigrateNode > DeleteNode:", err)
					return false, err
				}
			}
		case NAIVE:
			{
				if err := self.DeleteRow(node); err != nil {
					self.Logger.Fatal("@MigrateNode > DeleteRow:", err)
					return false, err
				} else {
					log.Println(fmt.Sprintf("%s node { %s }", color.FgRed.Render("Deleted"), node.Tag.Name))
				}
			}
		}
	}

	return migrated, nil
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
		self.Logger.Fatal(err)
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
					self.Logger.Fatal(serr)
				}
			} else {
				self.Logger.Fatal("@HandleUnmappedTags > Table id:", err)
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

func (self *MigrationWorkerV2) HandleMigration(node *DependencyNode) (bool, error) {

	if mapping, found := self.FetchMappingsForNode(node); found {
		tagMembers := node.Tag.GetTagMembers()
		if helper.Sublist(tagMembers, mapping.FromTables) {
			return self.MigrateNode(mapping, node)
		}
		if wNode, err := self.HandleWaitingList(mapping, tagMembers, node); wNode != nil && err == nil {
			return self.MigrateNode(mapping, wNode)
		} else {
			return false, err
		}
	} else {
		if strings.EqualFold(self.mtype, BAGS) || !strings.EqualFold(self.mtype, DELETION) {
			self.unmappedTags.Add(node.Tag.Name)
			return false, fmt.Errorf("no mapping found for node: %s", node.Tag.Name)
		}
		return false, self.HandleUnmappedNode(node)
	}
}

func (self *MigrationWorkerV2) HandleLeftOverWaitingNodes() {

	for _, waitingNode := range self.wList.Nodes {
		for _, containedNode := range waitingNode.ContainedNodes {
			self.HandleUnmappedNode(containedNode)
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

		if migrated, err := self.HandleMigration(node); err == nil {
			if migrated {
				log.Println(fmt.Sprintf("x%2dx %s  node { %s } ", threadID, color.FgLightGreen.Render("Migrated"), node.Tag.Name))
			} else {
				log.Println(fmt.Sprintf("x%2dx %s  node { %s } ", threadID, color.FgGreen.Render("Not Migrated / No Err"), node.Tag.Name))
			}
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
		log.Println(fmt.Sprintf("x%2dx %s  node { %s } | root [%s] : owner [%s]", threadID, color.FgBlue.Render("Visited"), node.Tag.Name, self.uid, ownerID))
		self.MarkAsVisited(node)
	}
	fmt.Println("------------------------------------------------------------------------")
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
		if migrated, err := self.HandleMigration(node); err == nil {
			if migrated {
				log.Println(fmt.Sprintf("x%2dx %s  node { %s } ", threadID, color.FgLightGreen.Render("Migrated"), node.Tag.Name))
			} else {
				log.Println(fmt.Sprintf("x%2dx %s  node { %s } ", threadID, color.FgGreen.Render("Not Migrated / No Err"), node.Tag.Name))
			}
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
			self.Logger.Fatal(fmt.Sprintf("x%2dx | UNABEL to COMMIT node { %s } ", threadID, node.Tag.Name))
			return err
		} else {
			log.Println(fmt.Sprintf("x%2dx COMMITTED node { %s } ", threadID, node.Tag.Name))
		}
	} else {
		log.Println(fmt.Sprintf("x%2dx %s  node { %s } | root [%s] : owner [%s]", threadID, color.FgBlue.Render("Visited"), node.Tag.Name, self.uid, ownerID))
	}
	self.MarkAsVisited(node)
	return nil
}

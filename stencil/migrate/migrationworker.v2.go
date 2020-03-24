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
		where := fmt.Sprintf("\"%s\".\"%s\" = '%s'", rootTable, rootCol, self.uid)
		ql := self.GetTagQL(root)
		sql := fmt.Sprintf("%s WHERE %s ", ql, where)
		sql += root.ResolveRestrictions()
		// self.Logger.Debug(sql)
		if data, err := db.DataCall1(self.SrcDBConn, sql); err == nil {
			if len(data) > 0 {
				self.root = &DependencyNode{Tag: root, SQL: sql, Data: data}
			} else {
				self.Logger.Trace(sql)
				return errors.New("Can't fetch Root node. Check if it exists. UID: " + self.uid)
			}
		} else {
			self.Logger.Debug(sql)
			return err
		}
	} else {
		self.Logger.Fatalf("Can't fetch root tag '%s' | App => %s, %s | err: %v", tagName, self.SrcAppConfig.AppID, self.SrcAppConfig.AppName, err)
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

	for _, dep := range self.SrcAppConfig.GetParentDependencies(node.Tag.Name) {
		for _, pdep := range dep.DependsOn {
			if parent, err := self.SrcAppConfig.GetTag(pdep.Tag); err == nil {
				if where, err := node.ResolveParentDependencyConditions(pdep.Conditions, parent); err == nil {
					ql := self.GetTagQL(parent)
					sql := fmt.Sprintf("%s WHERE %s ", ql, where)
					sql += parent.ResolveRestrictions()
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
						self.Logger.Fatal("@GetAllPreviousNodes: Error while DataCall: ", err)
						return nil, err
					}
				} else {
					// log.Println("@GetAllPreviousNodes > ResolveParentDependencyConditions: ", err)
				}
			} else {
				self.Logger.Fatal("@GetAllPreviousNodes: Tag doesn't exist? ", pdep.Tag)
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
			log.Println(fmt.Sprintf("FETCHING  tag for dependency { %s > %s } ", node.Tag.Name, dep.Tag))
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
		log.Println(fmt.Sprintf("FETCHING  tag  for ownership { %s } ", own.Tag))
		// if self.unmappedTags.Exists(own.Tag) {
		// 	log.Println(fmt.Sprintf("        UNMAPPED  tag  { %s } ", own.Tag))
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
	if err := SA1_display.GenDisplayFlagTx(self.tx.StencilTx, self.DstAppConfig.AppID, dtable.TableID, pk, fmt.Sprint(self.logTxn.Txn_id)); err != nil {
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
			// self.Logger.Debugf("Checking Condition | conditionKey [%s] conditionVal [%s]", conditionKey, conditionVal)
			if nodeVal, ok := node.Data[conditionKey]; ok {
				// self.Logger.Debugf("Checking Condition | nodeVal [%v]", nodeVal)
				if conditionVal[:1] == "#" {
					// fmt.Println("VerifyMappingConditions: conditionVal[:1] == #")
					// fmt.Println(conditionKey, conditionVal, nodeVal)
					// fmt.Scanln()
					switch conditionVal {
					case "#NULL":
						{
							if nodeVal != nil {
								// log.Println("Case #NULL | ", nodeVal, "!=", conditionVal)
								// fmt.Println(conditionKey, conditionVal, nodeVal)
								// self.Logger.Fatal("@VerifyMappingConditions: return false, from case #NULL:")
								return false
							} else {
								// log.Println("Case #NULL | ", nodeVal, "==", conditionVal)
							}
						}
					case "#NOTNULL":
						{
							if nodeVal == nil {
								// log.Println("Case #NOTNULL | ", nodeVal, "!=", conditionVal)
								// fmt.Println(conditionKey, conditionVal, nodeVal)
								// self.Logger.Fatal("@VerifyMappingConditions: return false, from case #NOTNULL:")
								return false
							} else {
								// log.Println("Case #NOTNULL | ", nodeVal, "==", conditionVal)
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
				self.Logger.Warnf("Checking Condition | nodeVal doesn't exist | [%s]", conditionKey)
				return false
			}
		}
	} else {
		// self.Logger.Debugf("No mapping conditions exist for table: %s", toTable.Table)
	}

	return true
}

func (self *MigrationWorkerV2) ValidateMappedTableData(toTable config.ToTable, mappedData MappedData) bool {
	for mappedCol, srcMappedCol := range toTable.Mapping {
		if srcMappedCol[0:1] == "$" || mappedCol == "id" {
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
		self.Logger.Debug(self.SrcAppConfig.Ownerships)
		self.Logger.Fatal("@GetNodeOwner: Ownership not found:", node.Tag.Name)
	}
	return "", false
}

func (self *MigrationWorkerV2) FetchFromMapping(nodeData map[string]interface{}, fromAttr string) (interface{}, string, string, *MappingRef, error) {

	var mappedVal interface{}
	var ref *MappingRef

	fromTable, cleanedFromAttr := "", ""

	args := strings.Split(fromAttr, ",")
	cleanedFromAttr = args[0]
	// fmt.Println(color.FgLightRed.Render("#############################################################################################################"))
	// self.Logger.Debugf("\n#FETCH: fromAttr: [%s] | cleanedFromAttr: [%s]", fromAttr, cleanedFromAttr)

	if nodeVal, ok := nodeData[args[2]]; ok {
		targetTabCol := strings.Split(args[0], ".")
		comparisonTabCol := strings.Split(args[1], ".")
		if res, err := db.FetchForMapping(self.SrcAppConfig.DBConn, targetTabCol[0], targetTabCol[1], comparisonTabCol[1], fmt.Sprint(nodeVal)); err != nil {
			self.Logger.Debug(targetTabCol[0], targetTabCol[1], comparisonTabCol[1], fmt.Sprint(nodeVal))
			self.Logger.Fatal("@FetchFromMapping: FetchForMapping | ", err)
		} else if len(res) > 0 {
			mappedVal = res[targetTabCol[1]]
			fromTable = targetTabCol[0]
			nodeData[args[0]] = res[targetTabCol[1]]
			if len(args) > 3 {
				toMemberTokens := strings.Split(args[3], ".")
				ref = &MappingRef{
					appID:      fmt.Sprint(self.SrcAppConfig.AppID),
					fromID:     res[targetTabCol[1]].(int64),
					fromMember: fmt.Sprint(targetTabCol[0]),
					fromAttr:   fmt.Sprint(targetTabCol[1]),
					toID:       res[targetTabCol[1]].(int64),
					toMember:   fmt.Sprint(toMemberTokens[0]),
					toAttr:     fmt.Sprint(toMemberTokens[1])}
			}
		} else {
			err = errors.New(fmt.Sprintf("@FetchFromMapping: FetchForMapping | Returned data is nil! Previous node already migrated? Args: [%s, %s, %s, %s]", targetTabCol[0], targetTabCol[1], comparisonTabCol[1], fmt.Sprint(nodeVal)))
			return mappedVal, fromTable, cleanedFromAttr, ref, nil
		}
	} else {
		fmt.Println(nodeData)
		self.Logger.Fatal("@FetchFromMapping: unable to #FETCH ", args[2])
	}
	// self.Logger.Debugf("\n#FETCH EXIT: fromAttr: [%s] | cleanedFromAttr: [%s] | FromTable: [%s], val: [%v]", fromAttr, cleanedFromAttr, fromTable, mappedVal)
	// fmt.Println(color.FgLightRed.Render("#############################################################################################################"))
	return mappedVal, fromTable, cleanedFromAttr, ref, nil
}

func (self *MigrationWorkerV2) RemoveMappedDataFromNodeData(mappedData MappedData, node *DependencyNode) {
	fmt.Println("Deleting Cols | ", mappedData.orgCols)
	for _, col := range strings.Split(mappedData.orgCols, ",") {
		for key := range node.Data {
			if strings.EqualFold(key, col) && !strings.Contains(key, ".id") {
				fmt.Printf("Deleting Key | Key: '%v' | Col: '%v'\n", key, col)
				delete(node.Data, key)
			}
		}
	}
}

func (self *MigrationWorkerV2) IsNodeDataEmpty(data map[string]interface{}) bool {
	for key, val := range data {
		if !(strings.Contains(key, ".id") || strings.Contains(key, ".display_flag")) && val != nil {
			return false
		}
	}
	return true
}

func (self *MigrationWorkerV2) DecodeMappingValue(fromAttr string, nodeData map[string]interface{}, args ...bool) (interface{}, string, string, *MappingRef, bool, error) {

	isBag := false
	// rawBag := false

	// if len(args) > 1 {
	// 	rawBag = args[1]
	// }

	if len(args) > 0 {
		isBag = args[0]
	}

	if self.mtype == BAGS {
		isBag = true
	}

	// self.Logger.Tracef("@DecodeMappingValue | fromAttr [%s] | isBag: [%v] | args: [%v] | data: %v", fromAttr, isBag, args, nodeData)

	var mappedVal interface{}
	var ref *MappingRef

	fromTable := ""
	found := true
	cleanedFromAttr := fromAttr

	switch fromAttr[0:1] {
	case "$":
		{
			if inputVal, err := self.mappings.GetInput(fromAttr); err == nil {
				mappedVal = inputVal
			} else {
				self.Logger.Debugf("@DecodeMappingValue | fromAttr [%s] | isBag: [%v] | data: %v", fromAttr, isBag, nodeData)
				self.Logger.Fatal(err)
			}
		}
	case "#":
		{
			cleanedFromAttr = self.CleanMappingAttr(fromAttr)
			// cleanedFromAttr = strings.Trim(cleanedFromAttr, "#ASSIGNFETCHREF")
			// self.Logger.Debug(color.FgLightYellow.Render(fmt.Sprintf("FromAttr: %s | Cleaned: %s", fromAttr, cleanedFromAttr)))
			if strings.Contains(fromAttr, "#REF") {
				if strings.Contains(fromAttr, "#FETCH") {
					var err error
					if !isBag {
						if mappedVal, fromTable, cleanedFromAttr, ref, err = self.FetchFromMapping(nodeData, cleanedFromAttr); err != nil {
							self.Logger.Fatal("@DecodeMappingValue > FetchFromMapping: ", cleanedFromAttr, err)
						}
					} else {
						found = false
					}
				} else if strings.Contains(fromAttr, "#ASSIGN") {
					cleanedFromAttrTokens := strings.Split(cleanedFromAttr, ",")
					referredTabCol := cleanedFromAttrTokens[1]
					cleanedFromAttr = strings.Trim(cleanedFromAttrTokens[0], "()")
					if nodeVal, ok := nodeData[cleanedFromAttr]; ok && nodeVal != nil {
						cleanedFromAttrTokens := strings.Split(cleanedFromAttr, ".")
						referredTabColTokens := strings.Split(referredTabCol, ".")
						fromTable = cleanedFromAttrTokens[0]
						mappedVal = nodeVal

						// if !isBag || rawBag {
						var fromID interface{}
						if val, ok := nodeData[cleanedFromAttrTokens[0]+".id"]; ok {
							fromID = val
						} else {
							fmt.Println(cleanedFromAttrTokens[0], " | ", cleanedFromAttrTokens)
							fmt.Println(nodeData)
							self.Logger.Fatal("@DecodeMappingValue > #REF > #ASSIGN > fromID: Unable to find ref value in node data | ", cleanedFromAttrTokens[0])
						}
						ref = &MappingRef{
							appID:      fmt.Sprint(self.SrcAppConfig.AppID),
							fromID:     helper.GetInt64(fromID),
							fromMember: fmt.Sprint(cleanedFromAttrTokens[0]),
							fromAttr:   fmt.Sprint(cleanedFromAttrTokens[1]),
							toID:       helper.GetInt64(nodeVal),
							toMember:   fmt.Sprint(referredTabColTokens[0]),
							toAttr:     fmt.Sprint(referredTabColTokens[1]),
						}
						// }
					} else {
						self.Logger.Debugf("fromAttr: [%s], cleanedFromAttr: [%s], nodeData: %v", fromAttr, cleanedFromAttr, nodeData)
						if isBag {
							self.Logger.Debugf("Unable to DecodeMappingValue | value found = [%v]", ok)
						} else {
							self.Logger.Fatalf("Unable to DecodeMappingValue | value found = [%v]", ok)
						}
					}
				} else {
					args := strings.Split(cleanedFromAttr, ",")
					cleanedFromAttr = args[0]
					if nodeVal, ok := nodeData[args[0]]; ok {
						argsTokens := strings.Split(args[0], ".")
						mappedVal = nodeVal
						fromTable = argsTokens[0]
					}
					hardRef := false
					if strings.Contains(fromAttr, "#REFHARD") {
						hardRef = true
					}
					// if !isBag || rawBag {
					if toID, fromID, err := GetIDsFromNodeData(args[0], args[1], nodeData, hardRef); err == nil {
						secondMemberTokens := strings.Split(args[1], ".")
						firstMemberTokens := strings.Split(args[0], ".")
						ref = &MappingRef{
							appID:      fmt.Sprint(self.SrcAppConfig.AppID),
							fromID:     fromID,
							fromMember: fmt.Sprint(firstMemberTokens[0]),
							fromAttr:   fmt.Sprint(firstMemberTokens[1]),
							toID:       toID,
							toMember:   fmt.Sprint(secondMemberTokens[0]),
							toAttr:     fmt.Sprint(secondMemberTokens[1]),
						}
					} else {
						self.Logger.Debugf("fromAttr: '%v' \n", fromAttr)
						self.Logger.Debugf("args[0]: '%v' \n", args[0])
						self.Logger.Debugf("toID: '%v' | fromID: '%v' \n", toID, fromID)
						fmt.Println(nodeData)
						// if !rawBag && !isBag {
						// 	self.Logger.Fatal("@DecodeMappingValue > GetIDs | ", err)
						// } else {
						self.Logger.Warn("@DecodeMappingValue > GetIDs | ", err)
						// }
					}
					// }
				}
			} else if strings.Contains(fromAttr, "#ASSIGN") {
				if nodeVal, ok := nodeData[cleanedFromAttr]; ok {
					cleanedFromAttrTokens := strings.Split(cleanedFromAttr, ".")
					mappedVal = nodeVal
					fromTable = cleanedFromAttrTokens[0]
				} else {
					self.Logger.Fatalf("@DecodeMappingValue > #ASSIGN > Can't find assigned attr in data | cleanedFromAttr:[%s]", cleanedFromAttr)
				}
			} else if strings.Contains(fromAttr, "#FETCH") {
				if !isBag {
					var err error
					if mappedVal, fromTable, cleanedFromAttr, ref, err = self.FetchFromMapping(nodeData, cleanedFromAttr); err != nil {
						self.Logger.Debug(nodeData)
						self.Logger.Debug(cleanedFromAttr)
						self.Logger.Fatal("@DecodeMappingValue > #FETCH > FetchFromMapping: Unable to fetch | ", err)
					}
				} else {
					found = false
				}
			} else {
				switch fromAttr {
				case "#GUID":
					{
						mappedVal = uuid.New()
					}
				case "#RANDINT":
					{
						mappedVal = self.SrcAppConfig.QR.NewRowId()
					}
				default:
					{
						self.Logger.Fatal("@DecodeMappingValue: Case not found:" + fromAttr)
					}
				}
			}
		}
	default:
		{
			if val, ok := nodeData[fromAttr]; ok {
				mappedVal = val
				fromTable = strings.Split(fromAttr, ".")[0]
			} else {
				found = false
			}
		}
	}

	// if isBag {
	// self.Logger.Debugf("@DecodeMappingValue | mappedVal: [%v], fromTable: [%s], cleanedFromAttr: [%s], fromAttr: [%s], found: [%v], isBag: [%v]", mappedVal, fromTable, cleanedFromAttr, fromAttr, found, isBag)
	// }

	return mappedVal, fromTable, cleanedFromAttr, ref, found, nil
}

func (self *MigrationWorkerV2) GetMappedData(toTable config.ToTable, node *DependencyNode, isBag, rawBag bool) (MappedData, error) {

	data := MappedData{
		cols:        "",
		vals:        "",
		orgCols:     "",
		orgColsLeft: "",
		srcTables:   make(map[string][]string),
		undoAction:  new(transaction.UndoAction)}

	newRowId := db.GetNewRowIDForTable(self.DstDBConn, toTable.Table)
	data.UpdateData("id", "", "", newRowId)
	// color.Red.Printf("Getting mapped data | %v | %v\n", toTable.Table, toTable.Mapping)
	for toAttr, fromAttr := range toTable.Mapping {
		// color.Red.Printf("Getting mapped data | toAttr : %v , fromAttr : %v  \n", toAttr, fromAttr)
		if strings.EqualFold("id", toAttr) {
			// fmt.Println("toAttr is id  ")
			// if self.mtype != BAGS && strings.Contains(fromAttr, "#REF") {
			if strings.Contains(fromAttr, "#REF") {
				// fmt.Println("fromAttr contains #REF  ")
				assignedTabCol := self.CleanMappingAttr(fromAttr)
				args := strings.Split(assignedTabCol, ",")
				hardRef := false
				if strings.Contains(fromAttr, "#REFHARD") {
					// fmt.Println("fromAttr contains #REFHARD  ")
					hardRef = true
				}
				// color.Red.Printf("Creating reference | args[0] : %v , args[1] : %v , node.Data : %v , hardRef : %v  \n", args[0], args[1], node.Data, hardRef)
				if toID, fromID, err := GetIDsFromNodeData(args[0], args[1], node.Data, hardRef); err == nil {
					secondMemberTokens := strings.Split(args[1], ".")
					firstMemberTokens := strings.Split(args[0], ".")
					data.UpdateRefs(self.SrcAppConfig.AppID, fromID, firstMemberTokens[0], firstMemberTokens[1], toID, secondMemberTokens[0], secondMemberTokens[1])
				} else {
					fmt.Printf("args[0]: '%v' \n", args[0])
					fmt.Printf("toID: '%v' | fromID: '%v' \n", toID, fromID)
					fmt.Printf("data: [%v] \n", node.Data)
					self.Logger.Fatal("@GetMappedData > id > GetIDs | ", err)
					return data, err
				}
			} else {
				// fmt.Println("fromAttr doesn't contain #REF  ")
			}
		} else if mappedValue, fromTable, cleanedFromAttr, ref, found, err := self.DecodeMappingValue(fromAttr, node.Data, isBag, rawBag); err == nil {
			if found {
				if mappedValue != nil {
					data.UpdateData(toAttr, cleanedFromAttr, fromTable, mappedValue)
					if ref != nil {
						data.refs = append(data.refs, *ref)
					}
					if len(fromTable) > 0 {
						data.undoAction.AddOrgTable(fromTable)
						data.undoAction.AddData(cleanedFromAttr, mappedValue)
					}
				} else {

				}
			} else {
				data.orgColsLeft += fmt.Sprintf("%s,", cleanedFromAttr)
			}
		} else {
			self.Logger.Fatalf("@DecodeMappingValue | fromAttr: %s | err: %s | Data: %v", fromAttr, err, node.Data)
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
			self.Logger.Tracef("@DeleteRow: '%v' not present or is null in node data! | %v", idCol, nodeVal)
		}
	}
	return nil
}

func (self *MigrationWorkerV2) TransferMedia(filePath string) error {

	color.Red.Println("File transfer is turned off!")
	return nil

	file, err := os.Open(filePath)
	if err != nil {
		log.Println(fmt.Sprintf("Can't open the file at: %s | ", filePath), err)
		return err
	}

	fpTokens := strings.Split(filePath, "/")
	fileName := fpTokens[len(fpTokens)-1]
	fsName := "/" + fileName

	log.Println(color.FgLightWhite.Render(fmt.Sprintf("Transferring file [%s] with name [%s] to [%s]...", filePath, fileName, fsName)))
	for {
		if err := self.FTPClient.Stor(fsName, file); err != nil {
			self.Logger.Errorf("File Transfer Failed: %v", err)
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

func (self *MigrationWorkerV2) CheckRawBag(node *DependencyNode) (bool, error) {
	for _, table := range node.Tag.Members {
		if tableID, err := db.TableID(self.logTxn.DBconn, table, self.SrcAppConfig.AppID); err == nil {
			if id, ok := node.Data[table+".id"]; ok {
				if idRows, err := self.GetRowsFromIDTable(self.SrcAppConfig.AppID, tableID, id, true); err == nil {
					if len(idRows) == 0 {
						return true, nil
					}
				} else {
					self.Logger.Debug(node.Data)
					self.Logger.Fatal("@CheckRawBag > GetRowsFromIDTable > ", self.SrcAppConfig.AppID, tableID, id, err)
				}
			} else {
				self.Logger.Debug(node.Data)
				self.Logger.Warn("@CheckRawBag > id doesn't exist in table ", table+".id")
			}
		} else {
			self.Logger.Fatal("@CheckRawBag > TableID, fromTable: error in getting table id for member! ", table, err)
		}
	}
	return false, nil
}

func (self *MigrationWorkerV2) MigrateNode(mapping config.Mapping, node *DependencyNode) (bool, error) {

	migrated, rawBag, isBag := false, false, false

	if self.mtype == BAGS {
		isBag = true
		if res, err := self.CheckRawBag(node); err == nil {
			rawBag = res
			if rawBag {
				self.Logger.Info("{{{{{ RAW BAG }}}}}")
			} else {
				self.Logger.Info("{{{{{ NOT RAW BAG }}}}}")
			}
		} else {
			self.Logger.Fatal("@MigrateNode > CheckRawBag > ", err)
		}
	}

	var allMappedData []MappedData

	for _, toTable := range mapping.ToTables {

		if !self.ValidateMappingConditions(toTable, node) {
			self.Logger.Infof("toTable: %s | ValidateMappingConditions | Mapping Conditions Not Validated", toTable.Table)
			continue
		} else {
			// self.Logger.Infof("toTable: %s | ValidateMappingConditions | Mapping Conditions Validated", toTable.Table)
		}
		fmt.Println(".........................................")
		if mappedData, mappedDataErr := self.GetMappedData(toTable, node, isBag, rawBag); mappedDataErr != nil {
			self.Logger.Debug(node.Data)
			self.Logger.Debug(mappedData)
			self.Logger.Fatal("@MigrateNode > GetMappedData Error | ", mappedDataErr)
		} else if len(mappedData.cols) > 0 && len(mappedData.vals) > 0 && len(mappedData.ivals) > 0 {
			if !self.ValidateMappedTableData(toTable, mappedData) {
				self.Logger.Tracef("toTable: %s | mappedData: %v", toTable.Table, mappedData)
				self.Logger.Warn("@MigrateNode > ValidateMappedTableData: All Nulls?")
				continue
			}

			if self.mtype == DELETION || self.mtype == BAGS {
				// self.Logger.Tracef("Before Merging Data | %s\n%v | %v\n---", toTable.Table, mappedData.cols, mappedData.ivals)
				if err := self.MergeBagDataWithMappedData(&mappedData, node, toTable); err != nil {
					self.Logger.Fatal("@MigrateNode > MergeDataFromBagsWithMappedData | ", err)
				}
				// self.Logger.Tracef("After Merging Data | %s\n%v | %v\n---", toTable.Table, mappedData.cols, mappedData.ivals)
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
							} else {
								color.LightBlue.Printf("New IDRow | FromApp: %s, DstApp: %s, FromTable: %s, ToTable: %s, FromID: %v, toID: %s, MigrationID: %s\n", self.SrcAppConfig.AppID, self.DstAppConfig.AppID, fromTableID, toTable.TableID, fromID, fmt.Sprint(id), fmt.Sprint(self.logTxn.Txn_id))
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
					self.Logger.Debug("@Params:", toTable.Table, fmt.Sprint(id), mappedData.orgCols, mappedData.cols, mappedData.undoAction, node)
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
				self.Logger.Infof("Inserted into '%s' with ID '%v' \ncols | %s\nvals | %v", toTable.Table, id, mappedData.cols, mappedData.ivals)
				allMappedData = append(allMappedData, mappedData)
			} else {
				self.Logger.Debugf("@Args | [toTable: %s], [cols: %s], [vals: %s], [ivals: %v], [srcTables: %s], [srcCols: %s]", toTable.Table, mappedData.cols, mappedData.vals, mappedData.ivals, mappedData.srcTables, mappedData.orgCols)
				self.Logger.Debugf("@NODE: %s | Data: %v", node.Tag.Name, node.Data)

				if self.mtype == DELETION {
					self.Logger.Fatal("@MigrateNode > InsertRowIntoAppDB: ", err)
				} else if self.mtype == BAGS {
					self.Logger.Error("@MigrateNode > InsertRowIntoAppDB: ", err)
				}
				return migrated, err
			}

			if self.mtype != BAGS {
				if err := self.AddMappedReferences(mappedData.refs); err != nil {
					log.Println(mappedData.refs)
					self.Logger.Fatal("@MigrateNode > AddMappedReferences: ", err)
					return migrated, err
				}
			} else if self.mtype == BAGS || rawBag {
				if self.SrcAppConfig.AppID == self.DstAppConfig.AppID {
					if err := self.AddInnerReferences(node); err != nil {
						log.Println(node)
						self.Logger.Fatal("@MigrateNode > AddInnerReferences: ", err)
						return migrated, err
					}
					if err := self.AddToReferencesViaDependencies(node); err != nil {
						log.Println(node)
						self.Logger.Fatal("@MigrateNode > AddToReferencesViaDependencies: ", err)
						return migrated, err
					}
				} else {
					if err := self.AddMappedReferencesIfNotExist(mappedData.refs); err != nil {
						log.Println(mappedData.refs)
						self.Logger.Fatal("@MigrateNode > AddMappedReferencesIfNotExist: ", err)
						return migrated, err
					}
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
		// self.Logger.Tracef("Migrated Data:\n", allMappedData)
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
				// self.AddToReferences(nextNode, node)
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
		log.Println(fmt.Sprintf("OWNED   node { %s } | root [%s] : owner [%s]", node.Tag.Name, self.uid, ownerID))
		if err := self.InitTransactions(); err != nil {
			return err
		} else {
			defer self.tx.SrcTx.Rollback()
			defer self.tx.DstTx.Rollback()
			defer self.tx.StencilTx.Rollback()
		}

		log.Println(fmt.Sprintf("CHECKING NEXT NODES { %s }", node.Tag.Name))

		if err := self.CheckNextNode(node); err != nil {
			return err
		}

		// log.Println(fmt.Sprintf("CHECKING PREVIOUS NODES { %s }", node.Tag.Name))

		// if previousNodes, err := self.GetAllPreviousNodes(node); err == nil {
		// for _, previousNode := range previousNodes {
		// self.AddToReferences(node, previousNode)
		// }
		// } else {
		// return err
		// }

		log.Println(fmt.Sprintf("HANDLING MIGRATION { %s }", node.Tag.Name))

		if migrated, err := self.HandleMigration(node); err == nil {
			if migrated {
				log.Println(fmt.Sprintf("%s  node { %s } ", color.FgLightGreen.Render("Migrated"), node.Tag.Name))
			} else {
				log.Println(fmt.Sprintf("%s  node { %s } ", color.FgGreen.Render("Not Migrated / No Err"), node.Tag.Name))
			}
		} else {
			if strings.EqualFold(err.Error(), "3") {
				log.Println(fmt.Sprintf("UNMAPPED  node { %s } ", node.Tag.Name))
			} else if strings.EqualFold(err.Error(), "2") {
				log.Println(fmt.Sprintf("Sent2Bag  node { %s } ", node.Tag.Name))
			} else {
				log.Println(fmt.Sprintf("FAILED    node { %s } ", node.Tag.Name))
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
			log.Println(fmt.Sprintf("COMMITTED node { %s } ", node.Tag.Name))
		}
	} else {
		log.Println(fmt.Sprintf("VISITED  node { %s } | root [%s] : owner [%s]", node.Tag.Name, self.uid, ownerID))
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
				log.Println(fmt.Sprintf("%s  node { %s } ", color.FgLightGreen.Render("Migrated"), node.Tag.Name))
			} else {
				log.Println(fmt.Sprintf("%s  node { %s } ", color.FgGreen.Render("Not Migrated / No Err"), node.Tag.Name))
			}
		} else {
			log.Println(fmt.Sprintf("RCVD ERR  node { %s } ", node.Tag.Name), err)
			// if self.unmappedTags.Exists(node.Tag.Name) {
			// 	log.Println(fmt.Sprintf("BREAKLOOP node { %s } ", node.Tag.Name), err)
			// 	continue
			// }
			if strings.EqualFold(err.Error(), "3") {
				log.Println(fmt.Sprintf("IGNORED   node { %s } ", node.Tag.Name))
			} else if strings.EqualFold(err.Error(), "2") {
				log.Println(fmt.Sprintf("BAGGED?   node { %s } ", node.Tag.Name))
			} else {
				log.Println(fmt.Sprintf("FAILED    node { %s } ", node.Tag.Name), err)
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
			self.Logger.Fatal(fmt.Sprintf("UNABEL to COMMIT node { %s } ", node.Tag.Name))
			return err
		} else {
			log.Println(fmt.Sprintf("COMMITTED node { %s } ", node.Tag.Name))
		}
	} else {
		log.Println(fmt.Sprintf("VISITED  node { %s } | root [%s] : owner [%s]", node.Tag.Name, self.uid, ownerID))
	}
	self.MarkAsVisited(node)
	return nil
}

package migrate_v2

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"stencil/SA1_display"
	config "stencil/config/v2"
	"stencil/db"
	"stencil/helper"
	"stencil/transaction"
	"strings"

	"github.com/google/uuid"
	"github.com/gookit/color"
)

func (mWorker *MigrationWorker) FetchRoot(threadID int) error {
	tagName := "root"
	if root, err := mWorker.SrcAppConfig.GetTag(tagName); err == nil {
		rootTable, rootCol := mWorker.SrcAppConfig.GetItemsFromKey(root, "root_id")
		where := fmt.Sprintf("\"%s\".\"%s\" = '%s'", rootTable, rootCol, mWorker.uid)
		ql := mWorker.GetTagQL(root)
		sql := fmt.Sprintf("%s WHERE %s ", ql, where)
		sql += root.ResolveRestrictions()
		// mWorker.Logger.Debug(sql)
		if data, err := db.DataCall1(mWorker.SrcAppConfig.DBConn, sql); err == nil {
			if len(data) > 0 {
				mWorker.Root = &DependencyNode{Tag: root, SQL: sql, Data: data}
			} else {
				mWorker.Logger.Trace(sql)
				return errors.New("Can't fetch Root node. Check if it exists. UID: " + mWorker.uid)
			}
		} else {
			mWorker.Logger.Debug(sql)
			return err
		}
	} else {
		mWorker.Logger.Fatalf("Can't fetch root tag '%s' | App => %s, %s | err: %v", tagName, mWorker.SrcAppConfig.AppID, mWorker.SrcAppConfig.AppName, err)
		return err
	}
	return nil
}

func (mWorker *MigrationWorker) GetAllNextNodes(node *DependencyNode) ([]*DependencyNode, error) {
	var nodes []*DependencyNode
	for _, dep := range mWorker.SrcAppConfig.GetSubDependencies(node.Tag.Name) {
		if child, err := mWorker.SrcAppConfig.GetTag(dep.Tag); err == nil {
			if where, err := node.ResolveDependencyConditions(mWorker.SrcAppConfig, dep, child); err == nil {
				ql := mWorker.GetTagQL(child)
				sql := fmt.Sprintf("%s WHERE %s ", ql, where)
				sql += child.ResolveRestrictions()
				// log.Println("@GetAllNextNodes | ", sql)
				if data, err := db.DataCall(mWorker.SrcAppConfig.DBConn, sql); err == nil {
					for _, datum := range data {
						newNode := new(DependencyNode)
						newNode.Tag = child
						newNode.SQL = sql
						newNode.Data = datum
						nodes = append(nodes, newNode)
					}
				} else {
					mWorker.Logger.Fatal("@GetAllNextNodes: Error while DataCall: ", err)
					return nil, err
				}
			} else {
				log.Println("@GetAllNextNodes > ResolveDependencyConditions | ", err)
			}
		} else {
			mWorker.Logger.Fatal("@GetAllNextNodes: Tag doesn't exist? ", dep.Tag)
		}
	}
	// if len(mWorker.SrcAppConfig.GetSubDependencies(node.Tag.Name)) > 0 {
	// 	log.Println("@GetAllNextNodes:", len(nodes))
	// 	mWorker.Logger.Fatal(nodes)
	// }
	return nodes, nil
}

func (mWorker *MigrationWorker) GetAllPreviousNodes(node *DependencyNode) ([]*DependencyNode, error) {
	var nodes []*DependencyNode

	if node.Tag.Name != "root" {
		if ownership := mWorker.SrcAppConfig.GetOwnership(node.Tag.Name, "root"); ownership != nil {
			if where, err := node.ResolveParentOwnershipConditions(ownership, mWorker.Root.Tag); err == nil {
				ql := mWorker.GetTagQL(mWorker.Root.Tag)
				sql := fmt.Sprintf("%s WHERE %s ", ql, where)
				sql += mWorker.Root.Tag.ResolveRestrictions()
				if data, err := db.DataCall(mWorker.SrcAppConfig.DBConn, sql); err == nil {
					for _, datum := range data {
						newNode := new(DependencyNode)
						newNode.Tag = mWorker.Root.Tag
						newNode.SQL = sql
						newNode.Data = datum
						nodes = append(nodes, newNode)
					}
				} else {
					fmt.Println(sql)
					mWorker.Logger.Fatal("@GetAllPreviousNodes: Error while DataCall: ", err)
					return nil, err
				}
			} else {
				log.Println("@GetAllPreviousNodes > ResolveParentOwnershipConditions: ", err)
			}
		} else {
			mWorker.Logger.Fatal("@GetAllPreviousNodes: Ownership doesn't exist? ", node.Tag.Name, "root")
		}
	}

	for _, dep := range mWorker.SrcAppConfig.GetParentDependencies(node.Tag.Name) {
		for _, pdep := range dep.DependsOn {
			if parent, err := mWorker.SrcAppConfig.GetTag(pdep.Tag); err == nil {
				if where, err := node.ResolveParentDependencyConditions(pdep.Conditions, parent); err == nil {
					ql := mWorker.GetTagQL(parent)
					sql := fmt.Sprintf("%s WHERE %s ", ql, where)
					sql += parent.ResolveRestrictions()
					if data, err := db.DataCall(mWorker.SrcAppConfig.DBConn, sql); err == nil {
						for _, datum := range data {
							newNode := new(DependencyNode)
							newNode.Tag = parent
							newNode.SQL = sql
							newNode.Data = datum
							nodes = append(nodes, newNode)
						}
					} else {
						fmt.Println(sql)
						mWorker.Logger.Fatal("@GetAllPreviousNodes: Error while DataCall: ", err)
						return nil, err
					}
				} else {
					// log.Println("@GetAllPreviousNodes > ResolveParentDependencyConditions: ", err)
				}
			} else {
				mWorker.Logger.Fatal("@GetAllPreviousNodes: Tag doesn't exist? ", pdep.Tag)
			}
		}
	}

	return nodes, nil
}

func (mWorker *MigrationWorker) GetAdjNode(node *DependencyNode, threadID int) (*DependencyNode, error) {
	if strings.EqualFold(node.Tag.Name, "root") {
		return mWorker.GetOwnedNode(threadID)
	}
	return mWorker.GetDependentNode(node, threadID)
}

func (mWorker *MigrationWorker) GetDependentNode(node *DependencyNode, threadID int) (*DependencyNode, error) {

	for _, dep := range mWorker.SrcAppConfig.ShuffleDependencies(mWorker.SrcAppConfig.GetSubDependencies(node.Tag.Name)) {
		if child, err := mWorker.SrcAppConfig.GetTag(dep.Tag); err == nil {
			log.Println(fmt.Sprintf("FETCHING  tag for dependency { %s > %s } ", node.Tag.Name, dep.Tag))
			if where, err := node.ResolveDependencyConditions(mWorker.SrcAppConfig, dep, child); err == nil {
				ql := mWorker.GetTagQL(child)
				sql := fmt.Sprintf("%s WHERE %s ", ql, where)
				sql += child.ResolveRestrictions()
				sql += mWorker.visitedNodes.ExcludeVisited(child)
				sql += " ORDER BY random()"
				// mWorker.Logger.Fatal(sql)
				if data, err := db.DataCall1(mWorker.SrcAppConfig.DBConn, sql); err == nil {
					if len(data) > 0 {
						newNode := DependencyNode{Tag: child, SQL: sql, Data: data}
						// if !mWorker.wList.IsAlreadyWaiting(newNode) &&
						if !mWorker.visitedNodes.IsVisited(&newNode) {
							return &newNode, nil
						}
					}
				} else {
					fmt.Println("@GetDependentNode > DataCall1 | ", err)
					mWorker.Logger.Fatal(sql)
					return nil, err
				}
			} else {
				log.Println("@GetDependentNode > ResolveDependencyConditions | ", err)
				// mWorker.Logger.Fatal(err)
			}
		}
	}
	return nil, nil
}

func (mWorker *MigrationWorker) GetOwnedNode(threadID int) (*DependencyNode, error) {

	for _, own := range mWorker.SrcAppConfig.GetShuffledOwnerships() {
		log.Println(fmt.Sprintf("FETCHING  tag  for ownership { %s } ", own.Tag))

		if child, err := mWorker.SrcAppConfig.GetTag(own.Tag); err == nil {
			if where, err := mWorker.Root.ResolveOwnershipConditions(own, child); err == nil {
				ql := mWorker.GetTagQL(child)
				sql := fmt.Sprintf("%s WHERE %s ", ql, where)
				sql += child.ResolveRestrictions()
				sql += mWorker.visitedNodes.ExcludeVisited(child)
				sql += " ORDER BY random() "
				// mWorker.Logger.Fatal(sql)
				if data, err := db.DataCall1(mWorker.SrcAppConfig.DBConn, sql); err == nil {
					if len(data) > 0 {
						newNode := DependencyNode{Tag: child, SQL: sql, Data: data}
						// if !mWorker.wList.IsAlreadyWaiting(newNode) {
						return &newNode, nil
						// }
					}
				} else {
					fmt.Println("@GetOwnedNode > DataCall1 | ", err)
					mWorker.Logger.Fatal(sql)
					return nil, err
				}
			} else {
				log.Println("@GetOwnedNode > ResolveOwnershipConditions | ", err)
			}
		}
	}
	return nil, nil
}

func (mWorker *MigrationWorker) PushData(tx *sql.Tx, dtable config.ToTable, pk string, mappedData MappedData, node *DependencyNode) error {

	undoActionSerialized, _ := json.Marshal(mappedData.undoAction)
	transaction.LogChange(string(undoActionSerialized), mWorker.logTxn)
	if err := SA1_display.GenDisplayFlagTx(mWorker.tx.StencilTx, mWorker.DstAppConfig.AppID, dtable.TableID, pk, fmt.Sprint(mWorker.logTxn.Txn_id)); err != nil {
		fmt.Println(mWorker.DstAppConfig.AppID, dtable.TableID, pk, fmt.Sprint(mWorker.logTxn.Txn_id))
		mWorker.Logger.Fatal("## DISPLAY ERROR!", err)
		return errors.New("0")
	}

	for fromTable, fromCols := range mappedData.srcTables {
		if _, ok := node.Data[fmt.Sprintf("%s.id", fromTable)]; ok {
			srcID := node.Data[fmt.Sprintf("%s.id", fromTable)]
			if fromTableID, err := db.TableID(mWorker.logTxn.DBconn, fromTable, mWorker.SrcAppConfig.AppID); err == nil {
				// if err := db.InsertIntoIdentityTable(tx, mWorker.SrcAppConfig.AppID, mWorker.DstAppConfig.AppID, fromTableID, dtable.TableID, srcID, pk, fmt.Sprint(mWorker.logTxn.Txn_id)); err != nil {
				// 	log.Println("@PushData:db.InsertIntoIdentityTable: ", mWorker.SrcAppConfig.AppID, mWorker.DstAppConfig.AppID, fromTableID, dtable.TableID, srcID, pk, fmt.Sprint(mWorker.logTxn.Txn_id))
				// 	mWorker.Logger.Fatal(err)
				// 	return errors.New("0")
				// }
				if serr := db.SaveForLEvaluation(tx, mWorker.SrcAppConfig.AppID, mWorker.DstAppConfig.AppID, fromTableID, dtable.TableID, srcID, pk, strings.Join(fromCols, ","), mappedData.cols, fmt.Sprint(mWorker.logTxn.Txn_id)); serr != nil {
					log.Println("@PushData:db.SaveForLEvaluation: ", mWorker.SrcAppConfig.AppID, mWorker.DstAppConfig.AppID, fromTableID, dtable.TableID, srcID, pk, strings.Join(fromCols, ","), mappedData.cols, fmt.Sprint(mWorker.logTxn.Txn_id))
					mWorker.Logger.Fatal(serr)
					return errors.New("0")
				}
			} else {
				log.Println("@PushData:db.TableID: ", fromTable, mWorker.SrcAppConfig.AppID)
				mWorker.Logger.Fatal(err)
			}
		}
	}
	return nil
}

func (mWorker *MigrationWorker) ValidateMappingConditions(toTable config.ToTable, node *DependencyNode) bool {
	if len(toTable.Conditions) > 0 {
		for conditionKey, conditionVal := range toTable.Conditions {
			// mWorker.Logger.Debugf("Checking Condition | conditionKey [%s] conditionVal [%s]", conditionKey, conditionVal)
			if nodeVal, ok := node.Data[conditionKey]; ok {
				// mWorker.Logger.Debugf("Checking Condition | nodeVal [%v]", nodeVal)
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
								// mWorker.Logger.Fatal("@VerifyMappingConditions: return false, from case #NULL:")
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
								// mWorker.Logger.Fatal("@VerifyMappingConditions: return false, from case #NOTNULL:")
								return false
							} else {
								// log.Println("Case #NOTNULL | ", nodeVal, "==", conditionVal)
							}
						}
					default:
						{
							fmt.Println(toTable.Table, conditionKey, conditionVal)
							mWorker.Logger.Fatal("@ValidateMappingConditions: Case not found:" + conditionVal)
						}
					}
				} else if conditionVal[:1] == "$" {
					// fmt.Println("VerifyMappingConditions: conditionVal[:1] == $")
					// fmt.Println(conditionKey, conditionVal, nodeVal)
					// fmt.Scanln()
					if inputVal, err := mWorker.mappings.GetInput(conditionVal); err == nil {
						if !strings.EqualFold(fmt.Sprint(nodeVal), inputVal) {
							log.Println(nodeVal, "!=", inputVal)
							fmt.Println(conditionKey, conditionVal, inputVal, nodeVal)
							mWorker.Logger.Fatal("@ValidateMappingConditions: return false, from conditionVal[:1] == $")
							return false
						}
					} else {
						fmt.Println("node data:", node.Data)
						fmt.Println(conditionKey, conditionVal)
						mWorker.Logger.Fatal("@ValidateMappingConditions: input doesn't exist?", err)
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
				// mWorker.Logger.Fatal("@ValidateMappingConditions: stop here and check")
				mWorker.Logger.Warnf("Checking Condition | nodeVal doesn't exist | [%s]", conditionKey)
				return false
			}
		}
	} else {
		// mWorker.Logger.Debugf("No mapping conditions exist for table: %s", toTable.Table)
	}

	return true
}

func (mWorker *MigrationWorker) ValidateMappedTableData(toTable config.ToTable, mappedData MappedData) bool {
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

func (mWorker *MigrationWorker) GetNodeOwner(node *DependencyNode) (string, bool) {

	if strings.EqualFold(node.Tag.Name, "root") {
		return mWorker.uid, true
	}

	if ownership := mWorker.SrcAppConfig.GetOwnership(node.Tag.Name, mWorker.Root.Tag.Name); ownership != nil {
		for _, condition := range ownership.Conditions {
			tagAttr, err := node.Tag.ResolveTagAttr(condition.TagAttr)
			if err != nil {
				mWorker.Logger.Fatal("@GetNodeOwner: Resolving TagAttr", err, node.Tag.Name, condition.TagAttr)
				break
			}
			depOnAttr, err := mWorker.Root.Tag.ResolveTagAttr(condition.DependsOnAttr)
			if err != nil {
				mWorker.Logger.Fatal("@GetNodeOwner: Resolving depOnAttr", err, node.Tag.Name, condition.DependsOnAttr)
				break
			}
			if nodeVal, err := node.GetValueForKey(tagAttr); err == nil {
				if rootVal, err := mWorker.Root.GetValueForKey(depOnAttr); err == nil {
					if !strings.EqualFold(nodeVal, rootVal) {
						// fmt.Println(fmt.Sprintf("root:%s:%s; user:%s:%s", depOnAttr, rootVal, tagAttr, nodeVal))
						return nodeVal, false
					} else {
						return nodeVal, true
					}
				} else {
					fmt.Println("@GetNodeOwner: Ownership Condition Key in Root Data:", depOnAttr, "doesn't exist!")
					fmt.Println("@GetNodeOwner: root data:", mWorker.Root.Data)
					mWorker.Logger.Fatal("@GetNodeOwner: stop here and check ownership conditions wrt root")
				}
			} else {
				fmt.Println("@GetNodeOwner: Ownership Condition Key", tagAttr, "doesn't exist!")
				fmt.Println("@GetNodeOwner: node data:", node.Data)
				fmt.Println("@GetNodeOwner: node sql:", node.SQL)
				mWorker.Logger.Fatal("@GetNodeOwner: stop here and check ownership conditions")
			}
		}
	} else {
		mWorker.Logger.Debug(mWorker.SrcAppConfig.Ownerships)
		mWorker.Logger.Fatal("@GetNodeOwner: Ownership not found:", node.Tag.Name)
	}
	return "", false
}

func (mWorker *MigrationWorker) FetchFromMapping(nodeData map[string]interface{}, fromAttr string) (interface{}, string, string, *MappingRef, error) {

	var mappedVal interface{}
	var ref *MappingRef

	fromTable, cleanedFromAttr := "", ""

	args := strings.Split(fromAttr, ",")
	cleanedFromAttr = args[0]
	// fmt.Println(color.FgLightRed.Render("#############################################################################################################"))
	// mWorker.Logger.Debugf("\n#FETCH: fromAttr: [%s] | cleanedFromAttr: [%s]", fromAttr, cleanedFromAttr)

	if nodeVal, ok := nodeData[args[2]]; ok {
		targetTabCol := strings.Split(args[0], ".")
		comparisonTabCol := strings.Split(args[1], ".")
		if res, err := db.FetchForMapping(mWorker.SrcAppConfig.DBConn, targetTabCol[0], targetTabCol[1], comparisonTabCol[1], fmt.Sprint(nodeVal)); err != nil {
			mWorker.Logger.Debug(targetTabCol[0], targetTabCol[1], comparisonTabCol[1], fmt.Sprint(nodeVal))
			mWorker.Logger.Fatal("@FetchFromMapping: FetchForMapping | ", err)
		} else if len(res) > 0 {
			mappedVal = res[targetTabCol[1]]
			fromTable = targetTabCol[0]
			nodeData[args[0]] = res[targetTabCol[1]]
			if len(args) > 3 {
				toMemberTokens := strings.Split(args[3], ".")
				ref = &MappingRef{
					appID:      fmt.Sprint(mWorker.SrcAppConfig.AppID),
					fromVal:    fmt.Sprint(res[targetTabCol[1]]),
					fromMember: fmt.Sprint(targetTabCol[0]),
					fromAttr:   fmt.Sprint(targetTabCol[1]),
					toVal:      fmt.Sprint(res[targetTabCol[1]]),
					toMember:   fmt.Sprint(toMemberTokens[0]),
					toAttr:     fmt.Sprint(toMemberTokens[1])}
			}
		} else {
			err = fmt.Errorf("@FetchFromMapping: FetchForMapping | Returned data is nil! Previous node already migrated? Args: [%s, %s, %s, %s]", targetTabCol[0], targetTabCol[1], comparisonTabCol[1], fmt.Sprint(nodeVal))
			return mappedVal, fromTable, cleanedFromAttr, ref, nil
		}
	} else {
		fmt.Println(nodeData)
		mWorker.Logger.Fatal("@FetchFromMapping: unable to #FETCH ", args[2])
	}
	// mWorker.Logger.Debugf("\n#FETCH EXIT: fromAttr: [%s] | cleanedFromAttr: [%s] | FromTable: [%s], val: [%v]", fromAttr, cleanedFromAttr, fromTable, mappedVal)
	// fmt.Println(color.FgLightRed.Render("#############################################################################################################"))
	return mappedVal, fromTable, cleanedFromAttr, ref, nil
}

func (mWorker *MigrationWorker) RemoveMappedDataFromNodeData(mappedData MappedData, node *DependencyNode) {
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

func (mWorker *MigrationWorker) IsNodeDataEmpty(data map[string]interface{}) bool {
	for key, val := range data {
		if !(strings.Contains(key, ".id") || strings.Contains(key, ".display_flag")) && val != nil {
			return false
		}
	}
	return true
}

func (mWorker *MigrationWorker) DecodeMappingValue(fromAttr string, nodeData map[string]interface{}, args ...bool) (interface{}, string, string, *MappingRef, bool, error) {

	isBag := false
	// rawBag := false

	// if len(args) > 1 {
	// 	rawBag = args[1]
	// }

	if len(args) > 0 {
		isBag = args[0]
	}

	if mWorker.mtype == BAGS {
		isBag = true
	}

	// mWorker.Logger.Tracef("@DecodeMappingValue | fromAttr [%s] | isBag: [%v] | args: [%v] | data: %v", fromAttr, isBag, args, nodeData)

	var mappedVal interface{}
	var ref *MappingRef

	fromTable := ""
	found := true
	cleanedFromAttr := fromAttr

	switch fromAttr[0:1] {
	case "$":
		{
			if inputVal, err := mWorker.mappings.GetInput(fromAttr); err == nil {
				mappedVal = inputVal
			} else {
				mWorker.Logger.Debugf("@DecodeMappingValue | fromAttr [%s] | isBag: [%v] | data: %v", fromAttr, isBag, nodeData)
				mWorker.Logger.Fatal(err)
			}
		}
	case "#":
		{
			cleanedFromAttr = mWorker.CleanMappingAttr(fromAttr)
			// cleanedFromAttr = strings.Trim(cleanedFromAttr, "#ASSIGNFETCHREF")
			// mWorker.Logger.Debug(color.FgLightYellow.Render(fmt.Sprintf("FromAttr: %s | Cleaned: %s", fromAttr, cleanedFromAttr)))
			if strings.Contains(fromAttr, "#REF") {
				if strings.Contains(fromAttr, "#FETCH") {
					var err error
					if !isBag {
						if mappedVal, fromTable, cleanedFromAttr, ref, err = mWorker.FetchFromMapping(nodeData, cleanedFromAttr); err != nil {
							mWorker.Logger.Fatal("@DecodeMappingValue > FetchFromMapping: ", cleanedFromAttr, err)
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
							mWorker.Logger.Fatal("@DecodeMappingValue > #REF > #ASSIGN > fromID: Unable to find ref value in node data | ", cleanedFromAttrTokens[0])
						}
						ref = &MappingRef{
							appID:      fmt.Sprint(mWorker.SrcAppConfig.AppID),
							fromVal:    fmt.Sprint(fromID),
							fromMember: fmt.Sprint(cleanedFromAttrTokens[0]),
							fromAttr:   fmt.Sprint(cleanedFromAttrTokens[1]),
							toVal:      fmt.Sprint(nodeVal),
							toMember:   fmt.Sprint(referredTabColTokens[0]),
							toAttr:     fmt.Sprint(referredTabColTokens[1]),
						}
						// }
					} else {
						mWorker.Logger.Debugf("fromAttr: [%s], cleanedFromAttr: [%s], nodeData: %v", fromAttr, cleanedFromAttr, nodeData)
						if isBag {
							mWorker.Logger.Debugf("Unable to DecodeMappingValue | value found = [%v]", ok)
						} else {
							mWorker.Logger.Fatalf("Unable to DecodeMappingValue | value found = [%v]", ok)
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
							appID:      fmt.Sprint(mWorker.SrcAppConfig.AppID),
							fromVal:    fmt.Sprint(fromID),
							fromMember: fmt.Sprint(firstMemberTokens[0]),
							fromAttr:   fmt.Sprint(firstMemberTokens[1]),
							toVal:      fmt.Sprint(toID),
							toMember:   fmt.Sprint(secondMemberTokens[0]),
							toAttr:     fmt.Sprint(secondMemberTokens[1]),
						}
					} else {
						mWorker.Logger.Debugf("fromAttr: '%v' \n", fromAttr)
						mWorker.Logger.Debugf("args[0]: '%v' \n", args[0])
						mWorker.Logger.Debugf("toID: '%v' | fromID: '%v' \n", toID, fromID)
						fmt.Println(nodeData)
						// if !rawBag && !isBag {
						// 	mWorker.Logger.Fatal("@DecodeMappingValue > GetIDs | ", err)
						// } else {
						mWorker.Logger.Warn("@DecodeMappingValue > GetIDs | ", err)
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
					mWorker.Logger.Fatalf("@DecodeMappingValue > #ASSIGN > Can't find assigned attr in data | cleanedFromAttr:[%s]", cleanedFromAttr)
				}
			} else if strings.Contains(fromAttr, "#FETCH") {
				if !isBag {
					var err error
					if mappedVal, fromTable, cleanedFromAttr, ref, err = mWorker.FetchFromMapping(nodeData, cleanedFromAttr); err != nil {
						mWorker.Logger.Debug(nodeData)
						mWorker.Logger.Debug(cleanedFromAttr)
						mWorker.Logger.Fatal("@DecodeMappingValue > #FETCH > FetchFromMapping: Unable to fetch | ", err)
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
						mappedVal = mWorker.SrcAppConfig.QR.NewRowId()
					}
				default:
					{
						mWorker.Logger.Fatal("@DecodeMappingValue: Case not found:" + fromAttr)
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
	// mWorker.Logger.Debugf("@DecodeMappingValue | mappedVal: [%v], fromTable: [%s], cleanedFromAttr: [%s], fromAttr: [%s], found: [%v], isBag: [%v]", mappedVal, fromTable, cleanedFromAttr, fromAttr, found, isBag)
	// }

	return mappedVal, fromTable, cleanedFromAttr, ref, found, nil
}

func (mWorker *MigrationWorker) GetMappedData(toTable config.ToTable, node *DependencyNode, isBag, rawBag bool) (MappedData, error) {

	data := MappedData{
		cols:        "",
		vals:        "",
		orgCols:     "",
		orgColsLeft: "",
		srcTables:   make(map[string][]string),
		undoAction:  new(transaction.UndoAction)}

	newRowId := db.GetNewRowIDForTable(mWorker.DstAppConfig.DBConn, toTable.Table)
	data.UpdateData("id", "", "", newRowId)
	// color.Red.Printf("Getting mapped data | %v | %v\n", toTable.Table, toTable.Mapping)
	for toAttr, fromAttr := range toTable.Mapping {
		// color.Red.Printf("Getting mapped data | toAttr : %v , fromAttr : %v  \n", toAttr, fromAttr)
		if strings.EqualFold("id", toAttr) {
			// fmt.Println("toAttr is id  ")
			// if mWorker.mtype != BAGS && strings.Contains(fromAttr, "#REF") {
			if strings.Contains(fromAttr, "#REF") {
				// fmt.Println("fromAttr contains #REF  ")
				assignedTabCol := mWorker.CleanMappingAttr(fromAttr)
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
					data.UpdateRefs(mWorker.SrcAppConfig.AppID, fromID, firstMemberTokens[0], firstMemberTokens[1], toID, secondMemberTokens[0], secondMemberTokens[1])
				} else {
					fmt.Printf("args[0]: '%v' \n", args[0])
					fmt.Printf("toID: '%v' | fromID: '%v' \n", toID, fromID)
					fmt.Printf("data: [%v] \n", node.Data)
					mWorker.Logger.Fatal("@GetMappedData > id > GetIDs | ", err)
					return data, err
				}
			} else {
				// fmt.Println("fromAttr doesn't contain #REF  ")
			}
		} else if mappedValue, fromTable, cleanedFromAttr, ref, found, err := mWorker.DecodeMappingValue(fromAttr, node.Data, isBag, rawBag); err == nil {
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
			mWorker.Logger.Fatalf("@DecodeMappingValue | fromAttr: %s | err: %s | Data: %v", fromAttr, err, node.Data)
		}
	}

	data.undoAction.AddDstTable(toTable.Table)
	data.Trim(", ")

	return data, nil
}

func (mWorker *MigrationWorker) DeleteRow(node *DependencyNode) error {
	for _, tagMember := range node.Tag.Members {
		idCol := fmt.Sprintf("%s.id", tagMember)
		if nodeVal, ok := node.Data[idCol]; ok && nodeVal != nil {
			srcID := fmt.Sprint(node.Data[idCol])
			if derr := db.ReallyDeleteRowFromAppDB(mWorker.tx.SrcTx, tagMember, srcID); derr != nil {
				fmt.Println("@ERROR_DeleteRowFromAppDB", derr)
				fmt.Println("@QARGS:", tagMember, srcID)
				// mWorker.Logger.Fatal(derr)
				return derr
			}
			if tagMemberID, err := db.TableID(mWorker.logTxn.DBconn, tagMember, mWorker.SrcAppConfig.AppID); err == nil {
				if derr := db.UpdateLEvaluation(mWorker.logTxn.DBconn, tagMemberID, srcID, mWorker.logTxn.Txn_id); derr != nil {
					fmt.Println("@ERROR_UpdateLEvaluation", derr)
					fmt.Println("@QARGS:", tagMember, srcID, mWorker.logTxn.Txn_id)
					mWorker.Logger.Fatal(derr)
					return derr
				}
			} else {
				mWorker.Logger.Fatal("@DeleteRow>TableID: ", err)
			}

		} else {
			// log.Println("node.Data =>", node.Data)
			mWorker.Logger.Tracef("@DeleteRow: '%v' not present or is null in node data! | %v", idCol, nodeVal)
		}
	}
	return nil
}

func (mWorker *MigrationWorker) TransferMedia(filePath string) error {

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
		if err := mWorker.FTPClient.Stor(fsName, file); err != nil {
			mWorker.Logger.Errorf("File Transfer Failed: %v", err)
			mWorker.FTPClient.Quit()
			mWorker.FTPClient = GetFTPClient()
			continue
			// return err
		}
		break
	}

	return nil
}

func (mWorker *MigrationWorker) HandleUnmappedMembersOfNode(mapping config.Mapping, node *DependencyNode) error {

	if mWorker.mtype != DELETION {
		return nil
	}
	for _, nodeMember := range node.Tag.GetTagMembers() {
		if !helper.Contains(mapping.FromTables, nodeMember) {
			if err := mWorker.SendMemberToBag(node, nodeMember, mWorker.uid); err != nil {
				return err
			}
		}
	}
	return nil
}

func (mWorker *MigrationWorker) DeleteNode(mapping config.Mapping, node *DependencyNode) error {

	if mWorker.mtype == DELETION {

		if !mWorker.IsNodeDataEmpty(node.Data) {
			if err := mWorker.SendNodeToBagWithOwnerID(node, mWorker.uid); err != nil {
				fmt.Println(node.Tag.Name)
				fmt.Println(node.Data)
				mWorker.Logger.Fatal("@DeleteNode > SendNodeToBagWithOwnerID:", err)
				return err
			}
		}

		if err := mWorker.DeleteRow(node); err != nil {
			fmt.Println(node.Tag.Name)
			fmt.Println(node)
			mWorker.Logger.Fatal("@DeleteNode > DeleteRow:", err)
			return err
		} else {
			log.Println(fmt.Sprintf("%s node { %s }", color.FgRed.Render("Deleted"), node.Tag.Name))
		}
	}

	return nil
}

func (mWorker *MigrationWorker) DeleteRoot(threadID int) error {

	if err := mWorker.InitTransactions(); err != nil {
		mWorker.Logger.Fatal("@DeleteRoot > InitTransactions", err)
		return err
	} else {
		defer mWorker.tx.SrcTx.Rollback()
		defer mWorker.tx.DstTx.Rollback()
		defer mWorker.tx.StencilTx.Rollback()
	}
	if mapping, found := mWorker.FetchMappingsForNode(mWorker.Root); found {
		if mWorker.mtype == NAIVE {
			if err := mWorker.DeleteRow(mWorker.Root); err != nil {
				mWorker.Logger.Fatal("@DeleteRoot:", err)
				return err
			}
		} else if mWorker.mtype == DELETION {
			if err := mWorker.DeleteNode(mapping, mWorker.Root); err != nil {
				mWorker.Logger.Fatal("@DeleteRoot:", err)
				return err
			}
		} else {
			mWorker.Logger.Fatal("ATTEMPTED DELETION IN DISALLOWED MIGRATION TYPE!")
		}
	} else {
		fmt.Println(mWorker.Root)
		mWorker.Logger.Fatal("@DeleteRoot: Can't find mappings for root | ", mapping, found)
	}
	if err := mWorker.CommitTransactions(); err != nil {
		mWorker.Logger.Fatal("@DeleteRoot: ERROR COMMITING TRANSACTIONS! ")
		return err
	}
	return nil
}

func (mWorker *MigrationWorker) CheckRawBag(node *DependencyNode) (bool, error) {
	for _, table := range node.Tag.Members {
		if tableID, err := db.TableID(mWorker.logTxn.DBconn, table, mWorker.SrcAppConfig.AppID); err == nil {
			if id, ok := node.Data[table+".id"]; ok {
				if idRows, err := mWorker.GetRowsFromIDTable(mWorker.SrcAppConfig.AppID, tableID, id, true); err == nil {
					if len(idRows) == 0 {
						return true, nil
					}
				} else {
					mWorker.Logger.Debug(node.Data)
					mWorker.Logger.Fatal("@CheckRawBag > GetRowsFromIDTable > ", mWorker.SrcAppConfig.AppID, tableID, id, err)
				}
			} else {
				mWorker.Logger.Debug(node.Data)
				mWorker.Logger.Warn("@CheckRawBag > id doesn't exist in table ", table+".id")
			}
		} else {
			mWorker.Logger.Fatal("@CheckRawBag > TableID, fromTable: error in getting table id for member! ", table, err)
		}
	}
	return false, nil
}

func (mWorker *MigrationWorker) MigrateNode(mapping config.Mapping, node *DependencyNode) (bool, error) {

	migrated, rawBag, isBag := false, false, false

	if mWorker.mtype == BAGS {
		isBag = true
		if res, err := mWorker.CheckRawBag(node); err == nil {
			rawBag = res
			if rawBag {
				mWorker.Logger.Info("{{{{{ RAW BAG }}}}}")
			} else {
				mWorker.Logger.Info("{{{{{ NOT RAW BAG }}}}}")
			}
		} else {
			mWorker.Logger.Fatal("@MigrateNode > CheckRawBag > ", err)
		}
	}

	var allMappedData []MappedData

	for _, toTable := range mapping.ToTables {

		if !mWorker.ValidateMappingConditions(toTable, node) {
			mWorker.Logger.Infof("toTable: %s | ValidateMappingConditions | Mapping Conditions Not Validated", toTable.Table)
			continue
		} else {
			// mWorker.Logger.Infof("toTable: %s | ValidateMappingConditions | Mapping Conditions Validated", toTable.Table)
		}
		fmt.Println(".........................................")
		if mappedData, mappedDataErr := mWorker.GetMappedData(toTable, node, isBag, rawBag); mappedDataErr != nil {
			mWorker.Logger.Debug(node.Data)
			mWorker.Logger.Debug(mappedData)
			mWorker.Logger.Fatal("@MigrateNode > GetMappedData Error | ", mappedDataErr)
		} else if len(mappedData.cols) > 0 && len(mappedData.vals) > 0 && len(mappedData.ivals) > 0 {
			if !mWorker.ValidateMappedTableData(toTable, mappedData) {
				mWorker.Logger.Tracef("toTable: %s | mappedData: %v", toTable.Table, mappedData)
				mWorker.Logger.Warn("@MigrateNode > ValidateMappedTableData: All Nulls?")
				continue
			}

			if mWorker.mtype == DELETION || mWorker.mtype == BAGS {
				// mWorker.Logger.Tracef("Before Merging Data | %s\n%v | %v\n---", toTable.Table, mappedData.cols, mappedData.ivals)
				if err := mWorker.MergeBagDataWithMappedData(&mappedData, node, toTable); err != nil {
					mWorker.Logger.Fatal("@MigrateNode > MergeDataFromBagsWithMappedData | ", err)
				}
				// mWorker.Logger.Tracef("After Merging Data | %s\n%v | %v\n---", toTable.Table, mappedData.cols, mappedData.ivals)
			}

			if id, err := db.InsertRowIntoAppDB(mWorker.tx.DstTx, toTable.Table, mappedData.cols, mappedData.vals, mappedData.ivals...); err == nil {
				for fromTable := range mappedData.srcTables {
					if fromTableID, err := db.TableID(mWorker.logTxn.DBconn, fromTable, mWorker.SrcAppConfig.AppID); err == nil {
						if fromID, ok := node.Data[fromTable+".id"]; ok {
							// fromID := fmt.Sprint(val.(int))
							if err := db.InsertIntoIdentityTable(mWorker.tx.StencilTx, mWorker.SrcAppConfig.AppID, mWorker.DstAppConfig.AppID, fromTableID, toTable.TableID, fromID, fmt.Sprint(id), fmt.Sprint(mWorker.logTxn.Txn_id)); err != nil {
								fmt.Println("@MigrateNode: InsertIntoIdentityTable")
								fmt.Println("@Args: ", mWorker.SrcAppConfig.AppID, mWorker.DstAppConfig.AppID, fromTableID, toTable.TableID, fromID, fmt.Sprint(id), fmt.Sprint(mWorker.logTxn.Txn_id))
								fmt.Println("@Params:", toTable.Table, fmt.Sprint(id), mappedData.orgCols, mappedData.cols, mappedData.undoAction, node)
								mWorker.Logger.Fatal(err)
								return migrated, err
							} else {
								color.LightBlue.Printf("New IDRow | FromApp: %s, DstApp: %s, FromTable: %s, ToTable: %s, FromID: %v, toID: %s, MigrationID: %s\n", mWorker.SrcAppConfig.AppID, mWorker.DstAppConfig.AppID, fromTableID, toTable.TableID, fromID, fmt.Sprint(id), fmt.Sprint(mWorker.logTxn.Txn_id))
							}
						} else {
							fmt.Println(node.Data)
							mWorker.Logger.Fatal("@MigrateNode: InsertIntoIdentityTable | " + fromTable + ".id doesn't exist")
						}
					} else {
						mWorker.Logger.Fatal("@MigrateNode > TableID, fromTable: error in getting table id for member! ", fromTable, err)
						return migrated, err
					}
				}

				if err := mWorker.PushData(mWorker.tx.StencilTx, toTable, fmt.Sprint(id), mappedData, node); err != nil {
					mWorker.Logger.Debug("@Params:", toTable.Table, fmt.Sprint(id), mappedData.orgCols, mappedData.cols, mappedData.undoAction, node)
					mWorker.Logger.Fatal(err)
					return migrated, err
				}

				if len(toTable.Media) > 0 {
					if filePathCol, ok := toTable.Media["path"]; ok {
						if filePath, ok := node.Data[filePathCol]; ok {
							if err := mWorker.TransferMedia(fmt.Sprint(filePath)); err != nil {
								mWorker.Logger.Fatal("@MigrateNode > TransferMedia: ", err)
							}
						}
					} else {
						mWorker.Logger.Fatal("@MigrateNode > toTable.Media: Path not found in map!")
					}
				}
				mWorker.Logger.Infof("Inserted into '%s' with ID '%v' \ncols | %s\nvals | %v", toTable.Table, id, mappedData.cols, mappedData.ivals)
				allMappedData = append(allMappedData, mappedData)
			} else {
				mWorker.Logger.Debugf("@Args | [toTable: %s], [cols: %s], [vals: %s], [ivals: %v], [srcTables: %s], [srcCols: %s]", toTable.Table, mappedData.cols, mappedData.vals, mappedData.ivals, mappedData.srcTables, mappedData.orgCols)
				mWorker.Logger.Debugf("@NODE: %s | Data: %v", node.Tag.Name, node.Data)

				if mWorker.mtype == DELETION {
					mWorker.Logger.Fatal("@MigrateNode > InsertRowIntoAppDB: ", err)
				} else if mWorker.mtype == BAGS {
					mWorker.Logger.Error("@MigrateNode > InsertRowIntoAppDB: ", err)
				}
				return migrated, err
			}

			if mWorker.mtype != BAGS {
				if err := mWorker.AddMappedReferences(mappedData.refs); err != nil {
					log.Println(mappedData.refs)
					mWorker.Logger.Fatal("@MigrateNode > AddMappedReferences: ", err)
					return migrated, err
				}
			} else if mWorker.mtype == BAGS || rawBag {
				if mWorker.SrcAppConfig.AppID == mWorker.DstAppConfig.AppID {
					if inDepRefs, err := CreateInnerDependencyReferences(mWorker.SrcAppConfig, node.Tag, node.Data, ""); err != nil {
						log.Println(node)
						mWorker.Logger.Fatal("@MigrateNode > CreateInnerDependencyReferences: ", err)
						return migrated, err
					} else if len(inDepRefs) > 0 {
						mappedData.refs = append(mappedData.refs, inDepRefs...)
					}

					if depRefs, err := CreateReferencesViaDependencies(mWorker.SrcAppConfig, node.Tag, node.Data, ""); err != nil {
						log.Println(node)
						mWorker.Logger.Fatal("@MigrateNode > CreateReferencesViaDependencies: ", err)
						return migrated, err
					} else if len(depRefs) > 0 {
						mappedData.refs = append(mappedData.refs, depRefs...)
					}

					if node.Tag.Name != "root" {
						if ownRefs, err := CreateReferencesViaOwnerships(mWorker.SrcAppConfig, node.Tag, node.Data, ""); err != nil {
							log.Println(node)
							mWorker.Logger.Fatal("@MigrateNode > CreateReferencesViaOwnerships: ", err)
							return migrated, err
						} else if len(ownRefs) > 0 {
							mappedData.refs = append(mappedData.refs, ownRefs...)
						}
					}
				}
				if err := mWorker.AddMappedReferencesIfNotExist(mappedData.refs); err != nil {
					log.Println(mappedData.refs)
					mWorker.Logger.Fatal("@MigrateNode > AddMappedReferencesIfNotExist: ", err)
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
		// mWorker.Logger.Tracef("Migrated Data:\n", allMappedData)
		for _, mappedData := range allMappedData {
			mWorker.RemoveMappedDataFromNodeData(mappedData, node)
		}
		if mWorker.mtype == BAGS {
			if err := mWorker.DeleteBag(node); err != nil {
				mWorker.Logger.Fatal("@MigrateNode > DeleteBag:", err)
				return false, err
			}
		}
	}

	if !strings.EqualFold(node.Tag.Name, "root") {
		switch mWorker.mtype {
		case DELETION:
			{
				if err := mWorker.DeleteNode(mapping, node); err != nil {
					mWorker.Logger.Fatal("@MigrateNode > DeleteNode:", err)
					return false, err
				}
			}
		case NAIVE:
			{
				if err := mWorker.DeleteRow(node); err != nil {
					mWorker.Logger.Fatal("@MigrateNode > DeleteRow:", err)
					return false, err
				} else {
					log.Println(fmt.Sprintf("%s node { %s }", color.FgRed.Render("Deleted"), node.Tag.Name))
				}
			}
		}
	}

	return migrated, nil
}

func (mWorker *MigrationWorker) HandleUnmappedNode(node *DependencyNode) error {
	if !strings.EqualFold(mWorker.mtype, DELETION) {
		return errors.New("3")
	} else {
		if err := mWorker.SendNodeToBag(node); err != nil {
			return err
		} else {
			return errors.New("2")
		}
	}
}

func (mWorker *MigrationWorker) HandleMigration(node *DependencyNode) (bool, error) {

	if mapping, found := mWorker.FetchMappingsForNode(node); found {
		tagMembers := node.Tag.GetTagMembers()
		if helper.Sublist(tagMembers, mapping.FromTables) {
			return mWorker.MigrateNode(mapping, node)
		} else {
			mWorker.Logger.Fatal("Waiting List Case!")
			return false, errors.New("Waiting List Case")
		}
	} else {
		if strings.EqualFold(mWorker.mtype, BAGS) || !strings.EqualFold(mWorker.mtype, DELETION) {
			return false, fmt.Errorf("no mapping found for node: %s", node.Tag.Name)
		}
		return false, mWorker.HandleUnmappedNode(node)
	}
}

func (mWorker *MigrationWorker) CheckNextNode(node *DependencyNode) error {

	if nextNodes, err := mWorker.GetAllNextNodes(node); err == nil {
		// log.Println(fmt.Sprintf("NEXT NODES FETCHED { %s } | nodes [%d]", node.Tag.Name, len(nextNodes)))
		if len(nextNodes) > 0 {
			for _, nextNode := range nextNodes {
				// log.Println(fmt.Sprintf("CURRENT NEXT NODE { %s > %s } %d/%d", node.Tag.Name, nextNode.Tag.Name, i, len(nextNodes)))
				// mWorker.AddToReferences(nextNode, node)
				if precedingNodes, err := mWorker.GetAllPreviousNodes(node); err != nil {
					return err
				} else if len(precedingNodes) <= 1 {
					if err := mWorker.CheckNextNode(nextNode); err != nil {
						return err
					}
					if err := mWorker.SendNodeToBag(nextNode); err != nil {
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

func (mWorker *MigrationWorker) CallMigration(node *DependencyNode, threadID int) error {

	if ownerID, isRoot := mWorker.GetNodeOwner(node); isRoot && len(ownerID) > 0 {
		log.Println(fmt.Sprintf("OWNED   node { %s } | root [%s] : owner [%s]", node.Tag.Name, mWorker.uid, ownerID))
		if err := mWorker.InitTransactions(); err != nil {
			return err
		} else {
			defer mWorker.tx.SrcTx.Rollback()
			defer mWorker.tx.DstTx.Rollback()
			defer mWorker.tx.StencilTx.Rollback()
		}

		log.Println(fmt.Sprintf("CHECKING NEXT NODES { %s }", node.Tag.Name))

		if err := mWorker.CheckNextNode(node); err != nil {
			return err
		}

		// log.Println(fmt.Sprintf("CHECKING PREVIOUS NODES { %s }", node.Tag.Name))

		// if previousNodes, err := mWorker.GetAllPreviousNodes(node); err == nil {
		// for _, previousNode := range previousNodes {
		// mWorker.AddToReferences(node, previousNode)
		// }
		// } else {
		// return err
		// }

		log.Println(fmt.Sprintf("HANDLING MIGRATION { %s }", node.Tag.Name))

		if migrated, err := mWorker.HandleMigration(node); err == nil {
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

		if err := mWorker.CommitTransactions(); err != nil {
			return err
		} else {
			log.Println(fmt.Sprintf("COMMITTED node { %s } ", node.Tag.Name))
		}
	} else {
		log.Println(fmt.Sprintf("VISITED  node { %s } | root [%s] : owner [%s]", node.Tag.Name, mWorker.uid, ownerID))
		mWorker.visitedNodes.MarkAsVisited(node)
	}
	fmt.Println("------------------------------------------------------------------------")
	return nil
}

func (mWorker *MigrationWorker) CallMigrationX(node *DependencyNode, threadID int) error {
	if ownerID, isRoot := mWorker.GetNodeOwner(node); isRoot && len(ownerID) > 0 {
		if err := mWorker.InitTransactions(); err != nil {
			return err
		} else {
			defer mWorker.tx.SrcTx.Rollback()
			defer mWorker.tx.DstTx.Rollback()
			defer mWorker.tx.StencilTx.Rollback()
		}
		if migrated, err := mWorker.HandleMigration(node); err == nil {
			if migrated {
				log.Println(fmt.Sprintf("%s  node { %s } ", color.FgLightGreen.Render("Migrated"), node.Tag.Name))
			} else {
				log.Println(fmt.Sprintf("%s  node { %s } ", color.FgGreen.Render("Not Migrated / No Err"), node.Tag.Name))
			}
		} else {
			log.Println(fmt.Sprintf("RCVD ERR  node { %s } ", node.Tag.Name), err)

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
		if err := mWorker.CommitTransactions(); err != nil {
			mWorker.Logger.Fatal(fmt.Sprintf("UNABEL to COMMIT node { %s } ", node.Tag.Name))
			return err
		} else {
			log.Println(fmt.Sprintf("COMMITTED node { %s } ", node.Tag.Name))
		}
	} else {
		log.Println(fmt.Sprintf("VISITED  node { %s } | root [%s] : owner [%s]", node.Tag.Name, mWorker.uid, ownerID))
	}
	mWorker.visitedNodes.MarkAsVisited(node)
	return nil
}

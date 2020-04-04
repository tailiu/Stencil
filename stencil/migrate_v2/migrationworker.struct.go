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
	"strconv"
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

func (mWorker *MigrationWorker) ValidateMappingConditions(conditions map[string]string, data DataMap) bool {

	for conditionKey, conditionVal := range conditions {
		// mWorker.Logger.Debugf("Checking Condition | conditionKey [%s] conditionVal [%s]", conditionKey, conditionVal)
		if nodeVal, ok := data[conditionKey]; ok {
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
						fmt.Println(conditionKey, conditionVal)
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
					fmt.Println("node data:", data)
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
			// fmt.Println("node sql:", node.SQL)
			// mWorker.Logger.Fatal("@ValidateMappingConditions: stop here and check")
			mWorker.Logger.Warnf("Checking Condition | nodeVal doesn't exist | [%s]", conditionKey)
			return false
		}
	}

	return true
}

func (mWorker *MigrationWorker) ResolveMappingMethodRef(cleanedMappedStmt string, data DataMap) (*MappedMemberValue, error) {

	args := strings.Split(cleanedMappedStmt, ",")
	fromAttr = args[0]

	if nodeVal, ok := data[fromAttr]; ok {
		mmv := &MappedMemberValue{
			AppID:    mWorker.SrcAppConfig.AppID,
			IsMethod: true,
			Value:    nodeVal,
			DBConn:   mWorker.logTxn.DBconn,
		}
		mmv.StoreMemberAndAttr(fromAttr)
		return mmv, nil
	}
	return nil, nil
}

func (mWorker *MigrationWorker) ResolveMappingMethodAssign(cleanedMappedStmt string, data DataMap) (*MappedMemberValue, error) {
	if nodeVal, ok := data[cleanedMappedStmt]; ok {
		mmv := &MappedMemberValue{
			AppID:    mWorker.SrcAppConfig.AppID,
			IsMethod: true,
			Value:    nodeVal,
			DBConn:   mWorker.logTxn.DBconn,
		}
		mmv.StoreMemberAndAttr(cleanedMappedStmt)
		return mmv, nil
	} else {
		fmt.Println(data)
		mWorker.Logger.Fatalf("@ResolveMappingMethodAssign > Can't find assigned attr in data | cleanedMappedStmt:[%s]", cleanedMappedStmt)
		return nil, nil
	}
}

func (mWorker *MigrationWorker) ResolveMappingMethodFetch(cleanedMappedStmt string, data DataMap) (*MappedMemberValue, error) {
	args := strings.Split(cleanedMappedStmt, ",")

	attrToFetch := args[0]
	attrToCompare := args[1]
	attrToCompareWith := args[2]

	if nodeVal, ok := data[attrToCompareWith]; ok {
		targetAttrTokens := strings.Split(attrToFetch, ".")
		comparisonAttrTokens := strings.Split(attrToCompare, ".")
		if res, err := db.FetchForMapping(mWorker.SrcAppConfig.DBConn, targetAttrTokens[0], targetAttrTokens[1], comparisonAttrTokens[1], fmt.Sprint(nodeVal)); err != nil {
			mWorker.Logger.Debug(argetAttrTokens[0], targetAttrTokens[1], comparisonAttrTokens[1], nodeVal)
			mWorker.Logger.Fatal("@ResolveMappingMethodFetch: FetchForMapping | ", err)
			return nil, err
		} else if len(res) > 0 {
			data[attrToFetch] = res[targetAttrTokens[1]]
			mmv := &MappedMemberValue{
				AppID:    mWorker.SrcAppConfig.AppID,
				IsMethod: true,
				Value:    data[attrToFetch],
				DBConn:   mWorker.logTxn.DBconn,
			}
			mmv.StoreMemberAndAttr(attrToFetch)
			return mmv, nil
		} else {
			return nil, fmt.Errorf("ResolveMappingMethodFetch | Returned data is nil. Previous node already migrated? | Args: '%s', '%s', '%s', '%s'", targetTabCol[0], targetTabCol[1], comparisonTabCol[1], fmt.Sprint(nodeVal))
		}
	} else {
		err := fmt.Errorf("@ResolveMappingMethodFetch: unable to #FETCH '%s'", args[2])
		mWorker.Logger.Debug(data)
		mWorker.Logger.Fatal(err)
		return nil, err
	}
}

func (mWorker *MigrationWorker) ResolveMappedStatement(mappedStmt string, data DataMap) (*MappedMemberValue, error) {

	// first character in the mappedStmt identifies the statment type.
	// "$" and "#" identify special cases.
	switch mappedStmt[0:1] {
	case "$":
		{
			if inputVal, err := mWorker.mappings.GetInput(mappedStmt); err == nil {
				mmv := &MappedMemberValue{
					IsInput: true,
					Value:   inputVal,
					DBConn:  mWorker.logTxn.DBconn,
				}
				return mmv, nil
			}
			mWorker.Logger.Debugf("@ResolveMappedStatement | mappedStmt [%s] | data: %v", mappedStmt, data)
			mWorker.Logger.Fatal(err)
			return nil, err
		}
	case "#":
		{
			cleanedMappedStmt = mWorker.CleanMappingAttr(mappedStmt)
			if strings.Contains(mappedStmt, "#REF") {

				var mmv *MappedMemberValue
				var err error

				if strings.Contains(fromAttr, "#FETCH") {
					mmv, err = mWorker.ResolveMappingMethodFetch(cleanedMappedStmt, data)
					if mmv != nil && err == nil {
						tokens := strings.Split(cleanedMappedStmt, ",")
						if err := mmv.CreateReference(tokens[0], tokens[3], mappedStmt, data); err != nil {
							mWorker.Logger.Debug(Err)
						}
					}
				} else if strings.Contains(fromAttr, "#ASSIGN") {
					mmv, err = mWorker.ResolveMappingMethodAssign(cleanedMappedStmt, data)
					if mmv != nil && err == nil {
						tokens := strings.Split(cleanedMappedStmt, ",")
						if err := mmv.CreateReference(tokens[0], tokens[2], mappedStmt, data); err != nil {
							mWorker.Logger.Debug(Err)
						}
					}
				} else {
					mmv, err = mWorker.ResolveMappingMethodRef(cleanedMappedStmt, data)
					if mmv != nil && err == nil {
						tokens := strings.Split(cleanedMappedStmt, ",")
						if err := mmv.CreateReference(tokens[0], tokens[1], mappedStmt, data); err != nil {
							mWorker.Logger.Debug(Err)
						}
					}
				}

				if err != nil {
					mWorker.Logger.Debug(mappedStmt, cleanedMappedStmt)
					mWorker.Logger.Debug(data)
					mWorker.Logger.Fatal(err)
				}

				if mmv == nil {
					err = fmt.Errorf("Value fetched is nil for mappedStmt: '%s', '%s'", mappedStmt, cleanedMappedStmt)
					mWorker.Logger.Debug(data)
					mWorker.Logger.Debug(errMsg)
				}

				return mmv, err

			} else if strings.Contains(mappedStmt, "#ASSIGN") {
				return mWorker.ResolveMappingMethodAssign(cleanedMappedStmt, data)
			} else if strings.Contains(mappedStmt, "#FETCH") {
				return mWorker.ResolveMappingMethodFetch(cleanedMappedStmt, data)
			} else if strings.Contains(mappedStmt, "#GUID") {
				mmv := &MappedMemberValue{
					IsExpression: true,
					Value:        uuid.New(),
					DBConn:       mWorker.logTxn.DBconn,
				}
				return mmv, nil
			} else if strings.Contains(mappedStmt, "#RANDINT") {
				mmv := &MappedMemberValue{
					IsExpression: true,
					Value:        db.NewRandInt(),
					DBConn:       mWorker.logTxn.DBconn,
				}
				return mmv, nil
			} else {
				errMsg := "Unidentified Mapping Method in mappedStmt: " + mappedStmt
				mWorker.Logger.Fatal(errMsg)
				return nil, errors.New(errMsg)
			}
		}
	default:
		{
			if val, ok := data[mappedStmt]; ok {
				mmv := &MappedMemberValue{
					AppID:  mWorker.SrcAppConfig.AppID,
					Value:  val,
					DBConn: mWorker.logTxn.DBconn,
				}
				mmv.StoreMemberAndAttr(mappedStmt)
				return mmv, nil
			}
			return nil, nil
		}
	}
}

func (mWorker *MigrationWorker) GetMappedMemberData(toTable config.ToTable, node *DependencyNode) MappedMemberData {

	newRowID := db.GetNewRowIDForTable(mWorker.DstDBConn, toTable.Table)
	mappedMemberData := MappedMemberData{
		ToID:   newRowID,
		AppID:  mWorker.SrcAppConfig.AppID,
		Data:   make(map[string]MappedMemberValue),
		DBConn: mWorker.logTxn.DBconn,
	}
	mappedMemberData.SetMember(toTable.Table)

	for toMemberAttr, mappedStmt := range toTable.Mapping {
		if mmv, err := mWorker.ResolveMappedStatement(mappedStmt, node.Data); err == nil && mmv != nil {
			mmv.ToID = newRowID
			if strings.EqualFold(toMemberAttr, "id") {
				mmv.Value = mmv.ToID
			}
			mappedMemberData.Data[toMemberAttr] = *mmv
		} else if err != nil {
			mWorker.Logger.Fatal(err)
		} else {
			mWorker.Logger.Tracef("%s.%s not resolved!", toTable.Table, toMemberAttr)
		}
	}

	return mappedMemberData
}

func (mWorker *MigrationWorker) DeleteNode(node *DependencyNode) error {
	for _, tagMember := range node.Tag.Members {
		idCol := fmt.Sprintf("%s.id", tagMember)
		if nodeVal, ok := node.Data[idCol]; ok && nodeVal != nil {
			srcID := fmt.Sprint(node.Data[idCol])
			if derr := db.ReallyDeleteRowFromAppDB(mWorker.tx.SrcTx, tagMember, srcID); derr != nil {
				fmt.Println("@ERROR_DeleteRowFromAppDB", derr)
				fmt.Println("@QARGS:", tagMember, srcID)
				mWorker.Logger.Fatal(derr)
				return derr
			}
		} else {
			mWorker.Logger.Tracef("@DeleteRow: '%v' not present or is null in node data! | %v", idCol, nodeVal)
		}
	}
	log.Printf("%s node { %s } \n", color.FgRed.Render("Deleted"), node.Tag.Name)
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

func (mWorker *MigrationWorker) DeleteRoot(threadID int) error {

	if err := mWorker.InitTransactions(); err != nil {
		mWorker.Logger.Fatal("@DeleteRoot > InitTransactions", err)
		return err
	}
	defer mWorker.RollbackTransactions()

	if err := mWorker.DeleteNode(mWorker.Root); err != nil {
		mWorker.Logger.Fatal("@DeleteRoot:", err)
		return err
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

func (mWorker *MigrationWorker) HandleBagsMigration(mapping config.Mapping, node *DependencyNode) (bool, error) {

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
			mWorker.Logger.Infof("toTable: %s | Mapping Conditions NOT Validated\n", toTable.Table)
			continue
		} else {
			mWorker.Logger.Infof("toTable: %s | Mapping Conditions Validated\n", toTable.Table)
		}

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

func (mWorker *MigrationWorker) MigrateMemberData(mmd *MappedMemberData, node *DependencyNode) error {

	colStr, phStr, valList := mmd.GetQueryArgs()

	if id, err := db.InsertRowIntoAppDB(mWorker.tx.DstTx, mmd.Member, colStr, phStr, valList); err == nil {
		mWorker.Logger.Infof("Inserted into '%s' with ID '%v', ToID: '%v' \ncols | %s\nvals | %v", mmd.ToMember, id, mmd.ToID, colStr, valList)
		return nil
	} else {
		mWorker.Logger.Infof("Failed to insert new row into '%s' | err: %s \n", mmd.Member, err)
		mWorker.Logger.Infof("Args | [%s], [%s], [%v]  \n", colStr, phStr, valList)
		return err
	}
}

func (mWorker *MigrationWorker) CreateAttributeRows(mmd *MappedMemberData) error {

	for toAttr, mmv := range mmd.Data {
		var fromValue interface{}
		if strings.EqualFold(toAttr, "id") {
			fromValue = mmv.FromID
		} else if mmv.Ref != nil {
			fromValue = mmv.Value
		} else {
			continue
		}
		if err := db.InsertIntoAttrTable(mWorker.tx.StencilTx, mmv.AppID, mWorker.DstAppConfig.AppID, mmv.FromMemberID, mmd.ToMemberID, mmv.FromID, mmv.ToID, mmv.FromAttr, toAttr, fromValue, mmv.Value, fmt.Sprint(mWorker.logTxn.Txn_id)); err != nil {
			mWorker.Logger.Debug("Args | FromApp: %s, DstApp: %s, FromTable: %s, ToTable: %s, FromID: %v, toID: %s, FromAttr: %s, ToAttr: %s, fromVal: %v, toVal: %v \n", mmv.AppID, mWorker.DstAppConfig.AppID, mmv.FromMemberID, mmd.ToMemberID, mmv.FromID, mmv.ToID, mmv.FromAttr, toAttr, fromValue, mmv.Value)
			mWorker.Logger.Fatal("Unable to insert into attrTable")
			return err
		} else {
			color.LightBlue.Printf("New AttrRow | FromApp: %s, DstApp: %s, FromTable: %s, ToTable: %s, FromID: %v, toID: %s, FromAttr: %s, ToAttr: %s, fromVal: %v, toVal: %v \n", mmv.AppID, mWorker.DstAppConfig.AppID, mmv.FromMemberID, mmd.ToMemberID, mmv.FromID, mmv.ToID, mmv.FromAttr, toAttr, fromValue, mmv.Value)
		}
	}
	return nil
}

func (mWorker *MigrationWorker) CreateReferenceRows(mmd *MappedMemberData) error {

	for toAttr, mmv := range mmd.Data {
		if mmv.Ref == nil {
			continue
		}
		if err := db.CreateNewReferenceV2(mWorker.tx.StencilTx, mmv.Ref.appID, mmv.Ref.fromMember, mmv.Ref.fromID, mmv.Ref.fromVal, mmv.Ref.toMember, mmv.Ref.toVal, mWorker.logTxn.Txn_id, mmv.Ref.fromAttr, mmv.Ref.toAttr); err != nil {
			mWorker.Logger.Debugf("App: %s, FromMember: %s, FromAttr: %s, FromID: %v, FromVal: %s, ToMember: %s, ToAttr: %s, ToVal: %s\n", mmv.Ref.appID, mmv.Ref.fromMember, mmv.Ref.fromAttr, mmv.Ref.fromID, mmv.Ref.fromVal, mmv.Ref.toMember, mmv.Ref.toAttr, mmv.Ref.toVal)
			mWorker.Logger.Debug(mmv.Ref)
			mWorker.Logger.Fatal(err)
		} else {
			color.LightBlue.Printf("App: %s, FromMember: %s, FromAttr: %s, FromID: %v, FromVal: %s, ToMember: %s, ToAttr: %s, ToVal: %s", mmv.Ref.appID, mmv.Ref.fromMember, mmv.Ref.fromAttr, mmv.Ref.fromID, mmv.Ref.fromVal, mmv.Ref.toMember, mmv.Ref.toAttr, mmv.Ref.toVal)
		}
	}
}

func (mWorker *MigrationWorker) HandleNodeDeletion(node *DependencyNode, rootCheck bool) error {

	if mWorker.mtype == INDEPENDENT || mWorker.mtype == CONSISTENT {
		return nil
	}

	if rootCheck && strings.EqualFold(node.Tag.Name, "root") {
		return nil
	}

	if mWorker.mtype == DELETION {
		if !node.IsEmptyExcept() {
			if err := mWorker.SendNodeToBagWithOwnerID(node, mWorker.uid); err != nil {
				fmt.Println(node.Tag.Name)
				fmt.Println(node.Data)
				mWorker.Logger.Fatal("@DeleteNode > SendNodeToBagWithOwnerID:", err)
				return err
			}
			log.Printf("%s { %s } | Owner ID: %v \n", color.FgLightYellow.Render("BAG"), node.Tag.Name, mWorker.uid)
		}
	}

	if err := mWorker.DeleteNode(node); err != nil {
		mWorker.Logger.Fatal("@MigrateNode > DeleteNode:", err)
		return err
	}

	return nil
}

func (mWorker *MigrationWorker) HandleMigration(node *DependencyNode) (bool, error) {

	if mapping, found := mWorker.FetchMappingsForNode(node); found {
		tagMembers := node.Tag.GetTagMembers()
		if helper.Sublist(tagMembers, mapping.FromTables) {

			var mappedMemberData []MappedMemberData

			for _, toTable := range mapping.ToTables {

				// Merge data from bags
				mappedMemberDatum := mWorker.GetMappedMemberData(toTable, node)

				if !mappedMemberDatum.ValidateMappingConditions(toTable.Conditions, node.Data) {
					mWorker.Logger.Info("Mapping conditions not validated!")
					mWorker.Logger.Debug("mappedMemberDatum | ", mappedMemberDatum)
					mWorker.Logger.Debug("toTable.Conditions | ", toTable.Conditions)
					continue
				}

				if !mappedMemberDatum.ValidateMappedData() {
					mWorker.Logger.Info("Mapped data not validated!")
					mWorker.Logger.Debug("mappedMemberDatum | ", mappedMemberDatum)
					continue
				}

				if err := mWorker.MigrateMemberData(mappedMemberDatum, node); err != nil {
					mWorker.Logger.Fatal(err)
				}

				if err := mWorker.CreateAttributeRows(mappedMemberDatum); err != nil {
					mWorker.Logger.Fatal(err)
				}

				// Check all different reference creation scenarios
				if err := mWorker.CreateReferenceRows(mappedMemberDatum); err != nil {
					mWorker.Logger.Fatal(err)
				}

				mappedMemberData = append(mappedMemberData, mappedMemberDatum)
			}

			node.DeleteMappedDataFromNode(mappedMemberData)

			if err := mWorker.HandleNodeDeletion(node, true); err != nil {
				mWorker.Logger.Fatal(err)
			}

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

func (mWorker *MigrationWorker) defunct__HandleMigration(node *DependencyNode) (bool, error) {

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

func (mWorker *MigrationWorker) CloseDBConns() {

	mWorker.SrcAppConfig.CloseDBConns()
	mWorker.DstAppConfig.CloseDBConns()
}

func (mWorker *MigrationWorker) RenewDBConn(isBlade ...bool) {
	mWorker.CloseDBConns()
	mWorker.logTxn.DBconn.Close()
	mWorker.logTxn.DBconn = db.GetDBConn(db.STENCIL_DB)
	mWorker.SrcAppConfig.DBConn = db.GetDBConn(mWorker.SrcAppConfig.AppName)
	mWorker.SrcAppConfig.DBConn = db.GetDBConn(mWorker.DstAppConfig.AppName, isBlade...)
}

func (mWorker *MigrationWorker) UserID() string {
	return mWorker.uid
}

func (mWorker *MigrationWorker) MigrationID() int {
	return mWorker.logTxn.Txn_id
}

func (mWorker *MigrationWorker) GetMemberDataFromNode(member string, nodeData DataMap) DataMap {
	memberData := make(DataMap)
	for col, val := range nodeData {
		colTokens := strings.Split(col, ".")
		colMember := colTokens[0]
		// colAttr := colTokens[1]
		if !strings.Contains(col, ".display_flag") && strings.Contains(colMember, member) && val != nil {
			memberData[col] = val
		}
	}
	return memberData
}

func (mWorker *MigrationWorker) GetTagQL(tag config.Tag) string {

	sql := "SELECT %s FROM %s "

	if len(tag.InnerDependencies) > 0 {
		cols := ""
		joinMap := tag.CreateInDepMap()
		seenMap := make(map[string]bool)
		joinStr := ""

		for fromTable, toTablesMap := range joinMap {
			if _, ok := seenMap[fromTable]; !ok {
				if len(joinStr) > 0 {
					joinStr += fmt.Sprintf(" FULL JOIN ")
				}
				joinStr += fmt.Sprintf("\"%s\"", fromTable)
				_, colStr := db.GetColumnsForTable(mWorker.SrcAppConfig.DBConn, fromTable)
				cols += colStr + ","
			}
			for toTable, conditions := range toTablesMap {
				if conditions != nil {
					conditions = append(conditions, joinMap[toTable][fromTable]...)
					if joinMap[toTable][fromTable] != nil {
						joinMap[toTable][fromTable] = nil
					}
					if _, ok := seenMap[toTable]; !ok {
						joinStr += fmt.Sprintf(" FULL JOIN \"%s\" ", toTable)
					}
					joinStr += fmt.Sprintf("  ON %s ", strings.Join(conditions, " AND "))
					_, colStr := db.GetColumnsForTable(mWorker.SrcAppConfig.DBConn, toTable)
					cols += colStr + ","
					seenMap[toTable] = true
				}
			}
			seenMap[fromTable] = true
		}
		sql = fmt.Sprintf(sql, strings.Trim(cols, ","), joinStr)
	} else {
		table := tag.Members["member1"]
		_, cols := db.GetColumnsForTable(mWorker.SrcAppConfig.DBConn, table)
		sql = fmt.Sprintf(sql, cols, table)
	}
	return sql
}

func (mWorker *MigrationWorker) GetTagQLForBag(tag config.Tag) string {

	if tableIDs, err := tag.MemberIDs(mWorker.logTxn.DBconn, mWorker.SrcAppConfig.AppID); err != nil {
		log.Fatal("@GetTagQLForBag: ", err)
	} else {

		sql := "SELECT array_to_json(array_remove(array[%s], NULL)) as pks_json, %s as json_data FROM %s "

		if len(tag.InnerDependencies) > 0 {
			idCols, cols := "", ""
			joinMap := tag.CreateInDepMap(true)
			// log.Fatalln(joinMap)
			seenMap := make(map[string]bool)
			joinStr := ""
			for fromTable, toTablesMap := range joinMap {
				// log.Print(fromTable, toTablesMap)
				if _, ok := seenMap[fromTable]; !ok {
					if len(joinStr) > 0 {
						joinStr += fmt.Sprintf(" FULL JOIN ")
					}
					joinStr += fmt.Sprintf("data_bags %s", fromTable)
					idCols += fmt.Sprintf("%s.pk,", fromTable)
					cols += fmt.Sprintf(" coalesce(%s.\"data\"::jsonb, '{}'::jsonb)  ||", fromTable)
				}
				for toTable, conditions := range toTablesMap {
					if conditions != nil {
						conditions = append(conditions, joinMap[toTable][fromTable]...)
						if joinMap[toTable][fromTable] != nil {
							joinMap[toTable][fromTable] = nil
						}
						if _, ok := seenMap[toTable]; !ok {
							joinStr += fmt.Sprintf(" FULL JOIN data_bags %s ", toTable)
						}
						joinStr += fmt.Sprintf(" ON %s.member = %s AND %s.member = %s AND %s ", fromTable, tableIDs[fromTable], toTable, tableIDs[toTable], strings.Join(conditions, " AND "))
						cols += fmt.Sprintf(" coalesce(%s.\"data\"::jsonb, '{}'::jsonb)  ||", toTable)
						idCols += fmt.Sprintf("%s.pk,", toTable)
						seenMap[toTable] = true
					}
				}
				seenMap[fromTable] = true
			}
			sql = fmt.Sprintf(sql, strings.Trim(idCols, ","), strings.Trim(cols, ",|"), joinStr)
		} else {
			table := tag.Members["member1"]
			joinStr := fmt.Sprintf("data_bags %s", table)
			idCols := fmt.Sprintf("%s.pk", table)
			cols := fmt.Sprintf(" coalesce(%s.\"data\"::jsonb, '{}'::jsonb)  ", table)
			sql = fmt.Sprintf(sql, idCols, cols, joinStr)
		}

		return sql
	}
	return ""
}

func (mWorker *MigrationWorker) InitTransactions() error {
	var err error
	mWorker.tx.SrcTx, err = mWorker.SrcAppConfig.DBConn.Begin()
	if err != nil {
		log.Fatal("Error creating Source DB Transaction! ", err)
		return err
	}
	mWorker.tx.DstTx, err = mWorker.DstAppConfig.DBConn.Begin()
	if err != nil {
		log.Fatal("Error creating Dst DB Transaction! ", err)
		return err
	}
	mWorker.tx.StencilTx, err = mWorker.logTxn.DBconn.Begin()
	if err != nil {
		log.Fatal("Error creating Stencil DB Transaction! ", err)
		return err
	}
	return nil
}

func (mWorker *MigrationWorker) CommitTransactions() error {
	// log.Fatal("@CommitTransactions: About to Commit!")
	if err := mWorker.tx.SrcTx.Commit(); err != nil {
		log.Fatal("Error committing Source DB Transaction! ", err)
		return err
	}
	if err := mWorker.tx.DstTx.Commit(); err != nil {
		log.Fatal("Error committing Destination DB Transaction! ", err)
		return err
	}
	if err := mWorker.tx.StencilTx.Commit(); err != nil {
		log.Fatal("Error committing Stencil DB Transaction! ", err)
		return err
	}
	return nil
}

func (mWorker *MigrationWorker) RollbackTransactions() error {
	if err := mWorker.tx.SrcTx.Rollback(); err != nil {
		log.Fatal("Error rolling back Source DB Transaction! ", err)
		return err
	}
	if err := mWorker.tx.DstTx.Rollback(); err != nil {
		log.Fatal("Error rolling back Dst DB Transaction! ", err)
		return err
	}
	if err := mWorker.tx.StencilTx.Rollback(); err != nil {
		log.Fatal("Error rolling back Stencil DB Transaction! ", err)
		return err
	}
	mWorker.Logger.Warn("Transactions Rolled Back!")
	return nil
}

func (mWorker *MigrationWorker) FetchMappingsForBag(srcApp, srcAppID, dstApp, dstAppID, srcMember, dstMember string) (config.MappedApp, config.Mapping, bool) {

	var combinedMapping config.Mapping
	var appMappings config.MappedApp
	if srcApp == dstApp {
		appMappings = *config.GetSelfSchemaMappings(mWorker.logTxn.DBconn, srcAppID, srcApp)
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
	return appMappings, combinedMapping, mappingFound
}

func (mWorker *MigrationWorker) CleanMappingAttr(attr string) string {
	cleanedAttr := strings.ReplaceAll(attr, "(", "")
	cleanedAttr = strings.ReplaceAll(cleanedAttr, ")", "")
	cleanedAttr = strings.ReplaceAll(cleanedAttr, "#ASSIGN", "")
	cleanedAttr = strings.ReplaceAll(cleanedAttr, "#FETCH", "")
	cleanedAttr = strings.ReplaceAll(cleanedAttr, "#REFHARD", "")
	cleanedAttr = strings.ReplaceAll(cleanedAttr, "#REF", "")
	return cleanedAttr
}

func (mWorker *MigrationWorker) FetchMappedAttribute(srcApp, srcAppID, dstApp, dstAppID, srcMember, dstMember, dstAttr string) (string, bool) {

	var appMappings config.MappedApp
	if srcApp == dstApp {
		appMappings = *config.GetSelfSchemaMappings(mWorker.logTxn.DBconn, srcAppID, srcApp)
	} else {
		appMappings = *config.GetSchemaMappingsFor(srcApp, dstApp)
	}
	for _, mapping := range appMappings.Mappings {
		if mappedTables := helper.IntersectString([]string{srcMember}, mapping.FromTables); len(mappedTables) > 0 {
			for _, toTableMapping := range mapping.ToTables {
				if strings.EqualFold(dstMember, toTableMapping.Table) {
					for toAttr, fromAttr := range toTableMapping.Mapping {
						if toAttr == dstAttr {
							cleanedAttr := mWorker.CleanMappingAttr(fromAttr)
							cleanedAttrTokens := strings.Split(cleanedAttr, ",")
							cleanedAttrTabCol := strings.Split(cleanedAttrTokens[0], ".")
							return cleanedAttrTabCol[1], true
						}
					}
				}
			}
		}
	}
	return "", false
}

func (mWorker *MigrationWorker) FetchMappingsForNode(node *DependencyNode) (config.Mapping, bool) {
	var combinedMapping config.Mapping
	tagMembers := node.Tag.GetTagMembers()
	mappingFound := false
	for _, mapping := range mWorker.mappings.Mappings {
		if mappedTables := helper.IntersectString(tagMembers, mapping.FromTables); len(mappedTables) > 0 {
			combinedMapping.FromTables = append(combinedMapping.FromTables, mapping.FromTables...)
			combinedMapping.ToTables = append(combinedMapping.ToTables, mapping.ToTables...)
			mappingFound = true
		}
	}
	return combinedMapping, mappingFound
}

func (mWorker *MigrationWorker) GetUserIDAppIDFromPreviousMigration(currentAppID, currentUID string) (string, string, error) {

	currentRootMemberID := db.GetAppRootMemberID(mWorker.logTxn.DBconn, currentAppID)

	currentUIDInt, err := strconv.ParseInt(currentUID, 10, 64)
	if err != nil {
		panic(err)
	}

	fmt.Printf("@GetUserIDAppIDFromPreviousMigration | Getting previous migration | App: '%v', UID: '%v', rootMemberID: '%v' \n", currentAppID, currentUIDInt, currentRootMemberID)

	if IDRows, err := mWorker.GetRowsFromIDTable(currentAppID, currentRootMemberID, currentUIDInt, false); err == nil {
		fmt.Println(IDRows)
		if len(IDRows) > 0 {
			for _, IDRow := range IDRows {
				prevRootMemberID := db.GetAppRootMemberID(mWorker.logTxn.DBconn, IDRow.FromAppID)
				if strings.EqualFold(IDRow.FromMemberID, prevRootMemberID) {
					fmt.Printf("@GetUserIDAppIDFromPreviousMigration | Previous migration found | App: '%v', UID: '%v', rootMemberID: '%v' \n", IDRow.FromAppID, IDRow.FromID, IDRow.FromMemberID)
					return IDRow.FromAppID, fmt.Sprint(IDRow.FromID), nil
				}
			}
		}
		fmt.Printf("@GetUserIDAppIDFromPreviousMigration | No previous migration found | App: '%v', UID: '%v', rootMemberID: '%v' \n", currentAppID, currentUIDInt, currentRootMemberID)
		return "", "", nil
	} else {
		log.Fatalf("@GetUserIDAppIDFromPreviousMigration | App: '%s', UID: '%v', rootMemberID: '%s' | err => %v \n", currentAppID, currentUIDInt, currentRootMemberID, err)
		return "", "", fmt.Errorf("no previous migration user and app id found for => currentAppID: %s, currentUID: %v", currentAppID, currentUIDInt)
	}
}

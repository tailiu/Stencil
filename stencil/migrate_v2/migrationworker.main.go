package migrate_v2

import (
	"errors"
	"fmt"
	"log"
	"os"
	"stencil/SA1_display"
	config "stencil/config/v2"
	"stencil/db"
	"stencil/helper"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gookit/color"
)

func (mWorker *MigrationWorker) CreateBagManager() *BagManager {
	var bagManager BagManager
	bagManager.Init(mWorker.logTxn.DBconn, mWorker.SrcAppConfig.AppID, mWorker.uid)
	return &bagManager
}

func (mWorker *MigrationWorker) FetchUserNode(rootKey, rootVal string) (*DependencyNode, error) {

	if root, err := mWorker.SrcAppConfig.GetTag("root"); err == nil {
		rootTable, rootCol := mWorker.SrcAppConfig.GetItemsFromKey(root, rootKey)
		where := fmt.Sprintf("\"%s\".\"%s\" = '%s'", rootTable, rootCol, rootVal)
		ql := mWorker.GetTagQL(root)
		sql := fmt.Sprintf("%s WHERE %s ", ql, where)
		sql += root.ResolveRestrictions()
		// mWorker.Logger.Debug(sql)
		if data, err := db.DataCall1(mWorker.SrcAppConfig.DBConn, sql); err == nil {
			if len(data) > 0 {
				return &DependencyNode{Tag: root, SQL: sql, Data: data}, nil
			} else {
				mWorker.Logger.Trace(sql)
				return nil, fmt.Errorf("Can't fetch Root node. rootKey: '%s' | rootVal: '%s' | rootTable: '%s' | rootCol: '%s'", rootKey, rootVal, rootTable, rootCol)
			}
		} else {
			mWorker.Logger.Debug(sql)
			return nil, err
		}
	} else {
		mWorker.Logger.Fatalf("Can't fetch root tag | App => %s, %s | err: %v", mWorker.SrcAppConfig.AppID, mWorker.SrcAppConfig.AppName, err)
		return nil, err
	}
}

func (mWorker *MigrationWorker) GetAllNextNodes(node *DependencyNode) ([]*DependencyNode, error) {
	var nodes []*DependencyNode

	for _, dep := range mWorker.SrcAppConfig.ShuffleDependencies(mWorker.SrcAppConfig.GetSubDependencies(node.Tag.Name)) {
		if child, err := mWorker.SrcAppConfig.GetTag(dep.Tag); err == nil {
			if where, err := node.ResolveDependencyConditions(mWorker.SrcAppConfig, dep, child); err == nil {
				ql := mWorker.GetTagQL(child)
				sql := fmt.Sprintf("%s WHERE %s ", ql, where)
				sql += child.ResolveRestrictions()
				// sql += mWorker.visitedNodes.ExcludeVisited(child)
				if res, err := db.DataCall(mWorker.SrcAppConfig.DBConn, sql); err == nil {
					for _, data := range res {
						newNode := DependencyNode{Tag: child, SQL: sql, Data: data}
						nodes = append(nodes, &newNode)
					}
				} else {
					fmt.Println("@GetAllNextNodes > GetDependentNode > DataCall1 | ", err)
					mWorker.Logger.Fatal(sql)
					return nil, err
				}
			} else {
				log.Println("@GetAllNextNodes > GetDependentNode > ResolveDependencyConditions | ", err)
			}
		}
	}

	return nodes, nil
}

func (mWorker *MigrationWorker) GetAllPreviousNodes(node *DependencyNode) ([]*DependencyNode, error) {
	var nodes []*DependencyNode

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
							// newNode.DeleteNulls()
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
						// newNode.DeleteNulls()
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

func (mWorker *MigrationWorker) PushData(mappedData MappedMemberData, node *DependencyNode) error {

	if err := SA1_display.GenDisplayFlagTx(mWorker.tx.StencilTx, mWorker.DstAppConfig.AppID, mappedData.ToMemberID, mappedData.ToID, fmt.Sprint(mWorker.logTxn.Txn_id)); err != nil {
		fmt.Println(mWorker.DstAppConfig.AppID, mappedData.ToMemberID, mappedData.ToID, fmt.Sprint(mWorker.logTxn.Txn_id))
		mWorker.Logger.Fatal("## DISPLAY ERROR!", err)
		return errors.New("0")
	}

	for _, fromTable := range mappedData.SrcTables() {
		if srcID, ok := node.Data[fmt.Sprintf("%s.id", fromTable.Name)]; ok {
			fromCols := strings.Join(mappedData.FromCols(fromTable.Name), ",")
			toCols := strings.Join(mappedData.ToCols(), ",")
			if serr := db.SaveForLEvaluation(mWorker.tx.StencilTx, mWorker.SrcAppConfig.AppID, mWorker.DstAppConfig.AppID, fromTable.ID, mappedData.ToMemberID, srcID, mappedData.ToID, fromCols, toCols, fmt.Sprint(mWorker.logTxn.Txn_id)); serr != nil {
				log.Println("@PushData:db.SaveForLEvaluation: ", mWorker.SrcAppConfig.AppID, mWorker.DstAppConfig.AppID, fromTable.ID, mappedData.ToMemberID, srcID, mappedData.ToID, fromCols, toCols, fmt.Sprint(mWorker.logTxn.Txn_id))
				mWorker.Logger.Fatal(serr)
				return errors.New("0")
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
				mWorker.Logger.Fatal("@GetNodeOwner > ResolveTagAttr: Resolving TagAttr |", err, node.Tag.Name, condition.TagAttr)
			}
			depOnAttr, err := mWorker.Root.Tag.ResolveTagAttr(condition.DependsOnAttr)
			if err != nil {
				mWorker.Logger.Fatal("@GetNodeOwner > ResolveTagAttr: Resolving depOnAttr | ", err, condition.DependsOnAttr)
				break
			}
			rootAttr, err := mWorker.Root.Tag.ResolveTagAttr("root_id")
			if err != nil {
				mWorker.Logger.Fatal("@GetNodeOwner > ResolveTagAttr: Resolving rootAttr | ", err)
			}
			if nodeVal, err := node.GetValueForKey(tagAttr); err == nil {
				if strings.EqualFold(rootAttr, depOnAttr) {
					if rootVal, err := mWorker.Root.GetValueForKey(depOnAttr); err == nil {
						isRoot := false
						if strings.EqualFold(nodeVal, rootVal) {
							isRoot = true
						}
						return nodeVal, isRoot
					} else {
						fmt.Println("@GetNodeOwner: Ownership Condition Key in Root Data:", depOnAttr, "doesn't exist!")
						fmt.Println("@GetNodeOwner: root data:", mWorker.Root.Data)
						mWorker.Logger.Fatal("@GetNodeOwner: stop here and check ownership conditions wrt root")
					}
				} else {
					if userNode, err := mWorker.FetchUserNode(condition.DependsOnAttr, nodeVal); err != nil {
						mWorker.Logger.Fatal(err)
					} else if userNode == nil {
						mWorker.Logger.Fatalf("User node not found | nodeVal: '%s' | dependsOnAttr: '%s'  ", nodeVal, condition.DependsOnAttr)
					} else {
						if rootVal, ok := userNode.Data[rootAttr]; ok {
							isRoot := false
							if strings.EqualFold(mWorker.uid, fmt.Sprint(rootVal)) {
								isRoot = true
							}
							return fmt.Sprint(rootVal), isRoot
						}
						fmt.Println(userNode.Data)
						mWorker.Logger.Fatalf("@GetNodeOwner: Can't find rootVal in userNode.Data[%s]", rootAttr)
					}
				}
			} else {
				fmt.Println("@GetNodeOwner: Ownership Condition Key", tagAttr, "doesn't exist!")
				fmt.Println("@GetNodeOwner: node data:", node.Data)
				fmt.Println("@GetNodeOwner: node sql:", node.SQL)
				mWorker.Logger.Fatal("@GetNodeOwner: stop here and check ownership conditions")
			}
		}
		mWorker.Logger.Fatal("@GetNodeOwner: Owner not found")
		return "", false
	} else {
		mWorker.Logger.Debug(mWorker.SrcAppConfig.Ownerships)
		mWorker.Logger.Fatal("@GetNodeOwner: Ownership not found:", node.Tag.Name)
		return "", false
	}
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

func (mWorker *MigrationWorker) ResolveMappingMethodRef(cleanedMappedStmt string, data DataMap, appID string) (*MappedMemberValue, error) {

	args := strings.Split(cleanedMappedStmt, ",")
	fromAttr := args[0]

	if nodeVal, ok := data[fromAttr]; ok {
		mmv := &MappedMemberValue{
			AppID:    appID,
			IsMethod: true,
			Value:    nodeVal,
			DBConn:   mWorker.logTxn.DBconn,
		}
		mmv.StoreMemberAndAttr(fromAttr)
		if err := mmv.SetFromID(data); err != nil {
			return nil, err
		}
		return mmv, nil
	}
	return nil, nil
}

func (mWorker *MigrationWorker) ResolveMappingMethodAssign(cleanedMappedStmt string, data DataMap, appID string) (*MappedMemberValue, error) {
	if nodeVal, ok := data[cleanedMappedStmt]; ok {
		mmv := &MappedMemberValue{
			AppID:    appID,
			IsMethod: true,
			Value:    nodeVal,
			DBConn:   mWorker.logTxn.DBconn,
		}
		mmv.StoreMemberAndAttr(cleanedMappedStmt)
		if err := mmv.SetFromID(data); err != nil {
			return nil, err
		}
		return mmv, nil
	} else {
		fmt.Println(data)
		mWorker.Logger.Infof("@ResolveMappingMethodAssign > Can't find assigned attr in data | cleanedMappedStmt:[%s]", cleanedMappedStmt)
		return nil, nil
	}
}

func (mWorker *MigrationWorker) ResolveMappingMethodFetch(cleanedMappedStmt string, data DataMap, appID string) (*MappedMemberValue, error) {
	args := strings.Split(cleanedMappedStmt, ",")

	attrToFetch := args[0]
	attrToCompare := args[1]
	attrToCompareWith := args[2]

	if nodeVal, ok := data[attrToCompareWith]; ok {
		targetAttrTokens := strings.Split(attrToFetch, ".")
		comparisonAttrTokens := strings.Split(attrToCompare, ".")
		if res, err := db.FetchForMapping(mWorker.SrcAppConfig.DBConn, targetAttrTokens[0], targetAttrTokens[1], comparisonAttrTokens[1], fmt.Sprint(nodeVal)); err != nil {

			mWorker.Logger.Debug(targetAttrTokens[0], targetAttrTokens[1], comparisonAttrTokens[1], nodeVal)
			mWorker.Logger.Fatal("@ResolveMappingMethodFetch: FetchForMapping | ", err)
			return nil, err
		} else if len(res) > 0 {
			data[attrToFetch] = res[targetAttrTokens[1]]
			mmv := &MappedMemberValue{
				AppID:    appID,
				IsMethod: true,
				Value:    data[attrToFetch],
				DBConn:   mWorker.logTxn.DBconn,
			}
			mmv.StoreMemberAndAttr(attrToFetch)
			if err := mmv.SetFromID(data); err != nil {
				return nil, err
			}
			return mmv, nil
		} else {
			return nil, fmt.Errorf("ResolveMappingMethodFetch | Returned data is nil. Previous node already migrated? | Args: '%s', '%s', '%s', '%s'", targetAttrTokens[0], targetAttrTokens[1], comparisonAttrTokens[1], fmt.Sprint(nodeVal))
		}
	} else {
		err := fmt.Errorf("@ResolveMappingMethodFetch: unable to #FETCH '%s'", args[2])
		mWorker.Logger.Debug(data)
		mWorker.Logger.Fatal(err)
		return nil, err
	}
}

func (mWorker *MigrationWorker) ResolveMappingMethodExpression(value interface{}) (*MappedMemberValue, error) {
	mmv := &MappedMemberValue{
		IsExpression: true,
		Value:        value,
		DBConn:       mWorker.logTxn.DBconn,
	}
	return mmv, nil
}

func (mWorker *MigrationWorker) ResolveMappedStatement(mappedStmt string, data DataMap, appID string) (*MappedMemberValue, error) {

	// first character in the mappedStmt identifies the statment type.
	// "$" and "#" identify special cases.
	// mWorker.Logger.Trace("=> Resolving: ", mappedStmt)
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
			} else {

				mWorker.Logger.Debugf("@ResolveMappedStatement | mappedStmt [%s] | data: %v", mappedStmt, data)
				mWorker.Logger.Fatal(err)
				return nil, err
			}
		}
	case "#":
		{
			cleanedMappedStmt := mWorker.CleanMappingAttr(mappedStmt)
			if strings.Contains(mappedStmt, "#REF") {

				var mmv *MappedMemberValue
				var err error
				var fromAttr, toAttr string

				tokens := strings.Split(cleanedMappedStmt, ",")

				if strings.Contains(mappedStmt, "#FETCH") {
					mmv, err = mWorker.ResolveMappingMethodFetch(cleanedMappedStmt, data, appID)
					fromAttr, toAttr = tokens[0], tokens[3]
				} else if strings.Contains(mappedStmt, "#ASSIGN") {
					mmv, err = mWorker.ResolveMappingMethodAssign(tokens[0], data, appID)
					fromAttr, toAttr = tokens[0], tokens[1]
				} else {
					mmv, err = mWorker.ResolveMappingMethodRef(cleanedMappedStmt, data, appID)
					fromAttr, toAttr = tokens[0], tokens[1]
				}

				if mmv != nil && err == nil {
					if err := mmv.CreateReference(fromAttr, toAttr, mappedStmt, data); err != nil {
						mWorker.Logger.Debug(err)
					}
				} else if err != nil {
					mWorker.Logger.Debug(mappedStmt, cleanedMappedStmt)
					mWorker.Logger.Debug(data)
					mWorker.Logger.Warn(err)
				} else {
					mWorker.Logger.Warnf("Value fetched is nil for mappedStmt: '%s', '%s'", mappedStmt, cleanedMappedStmt)
				}

				return mmv, err

			} else if strings.Contains(mappedStmt, "#ASSIGN") {
				return mWorker.ResolveMappingMethodAssign(cleanedMappedStmt, data, appID)
			} else if strings.Contains(mappedStmt, "#FETCH") {
				return mWorker.ResolveMappingMethodFetch(cleanedMappedStmt, data, appID)
			} else if strings.Contains(mappedStmt, "#GUID") {
				return mWorker.ResolveMappingMethodExpression(uuid.New())
			} else if strings.Contains(mappedStmt, "#RANDINT") {
				return mWorker.ResolveMappingMethodExpression(db.NewRandInt())
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
					AppID:  appID,
					Value:  val,
					DBConn: mWorker.logTxn.DBconn,
				}
				mmv.StoreMemberAndAttr(mappedStmt)
				if err := mmv.SetFromID(data); err != nil {
					return nil, err
				}
				// mWorker.Logger.Tracef("=> Resolved | From Member : %s(%s)\n", mmv.FromMember, mmv.FromMemberID)
				return mmv, nil
			}
			return nil, nil
		}
	}
}

func (mWorker *MigrationWorker) GetMappedMemberData(toTable config.ToTable, node *DependencyNode) MappedMemberData {

	newRowID := db.GetNewRowIDForTable(mWorker.DstAppConfig.DBConn, toTable.Table)
	mappedMemberData := MappedMemberData{
		ToID:    newRowID,
		ToAppID: mWorker.DstAppConfig.AppID,
		Data:    make(map[string]MappedMemberValue),
		DBConn:  mWorker.logTxn.DBconn,
	}
	mappedMemberData.SetMember(toTable.Table)

	for toMemberAttr, mappedStmt := range toTable.Mapping {
		if mmv, err := mWorker.ResolveMappedStatement(mappedStmt, node.Data, mWorker.SrcAppConfig.AppID); err == nil && mmv != nil {
			mmv.ToID = newRowID
			if strings.EqualFold(toMemberAttr, "id") {
				mmv.Value = mmv.ToID
			}
			mappedMemberData.Data[toMemberAttr] = *mmv
		} else if err != nil {
			mWorker.Logger.Debugf("%s.%s not resolved | %s", toTable.Table, toMemberAttr, err)
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

	if !mWorker.FTPFlag {
		color.LightRed.Println("***  File transfer is turned off  ***")
		return nil
	}

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

	if mWorker.DeleteRootFlag {
		color.LightRed.Println("***  Root deletion is turned off  ***")
		return nil
	}

	if err := mWorker.InitTransactions(); err != nil {
		mWorker.Logger.Fatal("@DeleteRoot > InitTransactions", err)
		return err
	}
	defer mWorker.RollbackTransactions()

	log.Println(fmt.Sprintf("CHECKING NEXT NODES { %s }", mWorker.Root.Tag.Name))
	if err := mWorker.CheckNextNode(mWorker.Root); err != nil {
		return err
	}

	if err := mWorker.HandleNodeDeletion(mWorker.Root, false); err != nil {
		mWorker.Logger.Fatal("@DeleteRoot:", err)
		return err
	}

	if err := mWorker.CommitTransactions(); err != nil {
		mWorker.Logger.Fatal("@DeleteRoot: ERROR COMMITING TRANSACTIONS! ")
		return err
	}
	return nil
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

func (mWorker *MigrationWorker) MigrateMemberData(mmd MappedMemberData, node *DependencyNode) error {

	colStr, phStr, valList := mmd.GetQueryArgs()

	if id, err := db.InsertRowIntoAppDB(mWorker.tx.DstTx, mmd.ToMember, colStr, phStr, valList...); err == nil {
		mWorker.Logger.Infof("Inserted into '%s' with ID '%v', ToID: '%v' \ncols | %s\nvals | %v\n", mmd.ToMember, id, mmd.ToID, colStr, valList)
		return nil
	} else {
		mWorker.Logger.Warnf("Failed to insert new row into '%s' | err: %s \ncols | %s \nplaceholders | %s \nvals | %v\n ", mmd.ToMember, err, colStr, phStr, valList)
		return err
	}
}

func (mWorker *MigrationWorker) CreateAttributeRowFromMMV(mmd MappedMemberData, mmv MappedMemberValue, toAttr string) error {
	if toAttrID, err := db.AttrID(mWorker.logTxn.DBconn, mmd.ToMemberID, toAttr); err != nil {
		mWorker.Logger.Fatal(err)
	} else {
		if err := db.InsertIntoAttrTable(mWorker.tx.StencilTx, mmv.AppID, mWorker.DstAppConfig.AppID, mmv.FromMemberID, mmd.ToMemberID, mmv.FromID, mmv.ToID, mmv.FromAttrID, toAttrID, helper.ConvertScientificNotationToString(mmv.Value), helper.ConvertScientificNotationToString(mmv.Value), fmt.Sprint(mWorker.logTxn.Txn_id)); err != nil {
			mWorker.Logger.Debugf("Args |\nFromApp: %s, DstApp: %s, FromTable: %s, ToTable: %s, FromID: %v, toID: %s, FromAttr: %s, ToAttr: %s, fromVal: %v, toVal: %v \n", mmv.AppID, mWorker.DstAppConfig.AppID, mmv.FromMemberID, mmd.ToMemberID, mmv.FromID, mmv.ToID, mmv.FromAttrID, toAttrID, helper.ConvertScientificNotationToString(mmv.Value), helper.ConvertScientificNotationToString(mmv.Value))
			mWorker.Logger.Fatal(err)
			return err
		} else {
			color.LightBlue.Printf("New AttrRow | FromApp: %s, DstApp: %s, FromTable: %s, ToTable: %s, FromID: %v, toID: %s, FromAttr: %s, ToAttr: %s, fromVal: %v, toVal: %v \n", mmv.AppID, mWorker.DstAppConfig.AppID, mmv.FromMemberID, mmd.ToMemberID, mmv.FromID, mmv.ToID, mmv.FromAttrID, toAttrID, mmv.Value, mmv.Value)
		}
	}
	return nil
}

func (mWorker *MigrationWorker) CreateAttributeRows(mmd MappedMemberData) error {

	for toAttr, mmv := range mmd.Data {

		if mmv.Ref == nil {
			continue
		}

		if err := mWorker.CreateAttributeRowFromMMV(mmd, mmv, toAttr); err != nil {
			mWorker.Logger.Fatal(err)
		}

	}
	return nil
}

func (mWorker *MigrationWorker) CreateIDAttributeRows(mmd MappedMemberData, data DataMap) error {

	for srcAppID, srcMembers := range mmd.GetSourceAppsAndTables() {
		if srcAppID != mWorker.SrcAppConfig.AppID {
			continue
		}
		for _, srcMember := range srcMembers {
			var srcID int64
			if val, ok := data[srcMember+".id"]; ok {
				if val != nil {
					srcID = helper.GetInt64(val)
				}
			} else {
				fmt.Println(data)
				mWorker.Logger.Fatal("Unable to get id for member: ", srcMember)
			}

			if srcMemberID, err := db.TableID(mmd.DBConn, srcMember, srcAppID); err != nil {
				fmt.Println(srcMember, srcAppID)
				log.Fatal("@CreateIDAttributeRows.TableID: ", err)
			} else {
				if srcAttrID, err := db.AttrID(mWorker.logTxn.DBconn, srcMemberID, "id"); err != nil {
					mWorker.Logger.Fatal("@CreateIDAttributeRows.SrcAttrID: ", err)
				} else {
					if toAttrID, err := db.AttrID(mWorker.logTxn.DBconn, mmd.ToMemberID, "id"); err != nil {
						mWorker.Logger.Fatal("@CreateIDAttributeRows.ToAttrID: ", err)
					} else {
						if err := db.InsertIntoAttrTable(mWorker.tx.StencilTx, srcAppID, mWorker.DstAppConfig.AppID, srcMemberID, mmd.ToMemberID, srcID, mmd.ToID, srcAttrID, toAttrID, srcID, mmd.ToID, fmt.Sprint(mWorker.logTxn.Txn_id)); err != nil {
							mWorker.Logger.Debugf("Args |\nFromApp: %s, DstApp: %s, FromTable: %s, ToTable: %s, FromID: %v, toID: %s, FromAttr: %s, ToAttr: %s, fromVal: %v, toVal: %v \n", srcAppID, mWorker.DstAppConfig.AppID, srcMemberID, mmd.ToMemberID, srcID, mmd.ToID, srcAttrID, toAttrID, srcID, mmd.ToID)
							mWorker.Logger.Fatal(err)
							return err
						} else {
							color.LightBlue.Printf("New IDRow | FromApp: %s, DstApp: %s, FromTable: %s, ToTable: %s, FromID: %v, toID: %s, FromAttr: %s, ToAttr: %s, fromVal: %v, toVal: %v \n", srcAppID, mWorker.DstAppConfig.AppID, srcMemberID, mmd.ToMemberID, srcID, mmd.ToID, srcAttrID, toAttrID, srcID, mmd.ToID)
						}
					}
				}
			}
		}
	}
	return nil
}

func (mWorker *MigrationWorker) CreateSelfReferenceRows(mmd MappedMemberData) error {

	for srcAppID, srcMembers := range mmd.GetSourceAppsAndTables() {
		if srcAppID == mWorker.DstAppConfig.AppID {
			mWorker.Logger.Tracef("Creating Self References | App: %v | Members: %v \n", srcAppID, srcMembers)
			var visitedMembers []string
			for _, srcMember := range srcMembers {
				if !helper.Contains(visitedMembers, srcMember) {
					if bagTag, err := mWorker.DstAppConfig.GetTagByMember(srcMember); err == nil {
						visitedMembers = append(visitedMembers, bagTag.GetTagMembers()...)
						if bagRefs, err := mmd.CreateSelfReferences(mWorker.DstAppConfig, *bagTag, mmd.GetSrcDataMap()); err == nil {
							for _, ref := range bagRefs {
								if mmvToAttr, refMMV := mmd.FindMMV(ref.appID, fmt.Sprint(ref.fromID), fmt.Sprint(ref.fromMemberID), ref.fromAttrID); refMMV != nil {
									// color.Info.Println(fmt.Sprintf(">>>>>>>>>>> MMV FOUND! app: %s, fromID: %s, fromMemberID: %s, fromAttrID: %s", ref.appID, fmt.Sprint(ref.fromID), fmt.Sprint(ref.fromMemberID), ref.fromAttrID))
									if err := mWorker.CreateAttributeRowFromMMV(mmd, *refMMV, mmvToAttr); err != nil {
										mWorker.Logger.Fatal(err)
									}
								} else {
									fmt.Println(mmd.GetSrcDataMap())
									mWorker.Logger.Warnf("MMV NOT FOUND! app: '%s', fromID: '%s', fromMemberID: '%s', fromAttrID: '%s'", ref.appID, fmt.Sprint(ref.fromID), fmt.Sprint(ref.fromMemberID), ref.fromAttrID)
								}
								if ref.fromVal == STENCIL_NULL || ref.toVal == STENCIL_NULL {
									color.Yellow.Printf("Self-Ref Exists | App: %s, FromMember: %s, FromAttr: %s, FromVal: %s, FromID: %v, ToMember: %s, ToAttr: %s, ToVal: %s\n", ref.appID, ref.fromMemberID, ref.fromAttrID, ref.fromVal, ref.fromID, ref.toMemberID, ref.toAttrID, ref.toVal)
									continue
								}
								if err := db.CreateNewReferenceV2(mWorker.tx.StencilTx, ref.appID, ref.fromMemberID, ref.fromVal, ref.fromID, ref.toMemberID, ref.toVal, fmt.Sprint(mWorker.logTxn.Txn_id), ref.fromAttrID, ref.toAttrID); err != nil {
									mWorker.Logger.Debugf("App: %s, FromMember: %s, FromAttr: %s, FromVal: %s, ToMember: %s, ToAttr: %s, ToVal: %s\n", ref.appID, ref.fromMemberID, ref.fromAttrID, ref.fromVal, ref.toMemberID, ref.toAttrID, ref.toVal)
									mWorker.Logger.Debug(ref)
									mWorker.Logger.Fatal(err)
									return err
								} else {
									color.Magenta.Printf("New Self-Ref | App: %s, FromMember: %s, FromAttr: %s, FromVal: %s, ToMember: %s, ToAttr: %s, ToVal: %s\n", ref.appID, ref.fromMemberID, ref.fromAttrID, ref.fromVal, ref.toMemberID, ref.toAttrID, ref.toVal)
								}
							}
						} else {
							mWorker.Logger.Fatal(err)
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func (mWorker *MigrationWorker) CreateReferenceRows(mmd MappedMemberData) error {

	for _, mmv := range mmd.Data {

		if mmv.Ref == nil {
			continue
		}

		if mmv.Ref.fromVal == STENCIL_NULL || mmv.Ref.toVal == STENCIL_NULL {
			color.Yellow.Printf("Ref Exists | App: %s, FromMember: %s, FromAttr: %s, FromVal: %s, FromID: %v, ToMember: %s, ToAttr: %s, ToVal: %s\n", mmv.Ref.appID, mmv.Ref.fromMemberID, mmv.Ref.fromAttrID, mmv.Ref.fromVal, mmv.Ref.fromID, mmv.Ref.toMemberID, mmv.Ref.toAttrID, mmv.Ref.toVal)
			continue
		}
		if err := db.CreateNewReferenceV2(mWorker.tx.StencilTx, mmv.Ref.appID, mmv.Ref.fromMemberID, mmv.Ref.fromVal, mmv.Ref.fromID, mmv.Ref.toMemberID, mmv.Ref.toVal, fmt.Sprint(mWorker.logTxn.Txn_id), mmv.Ref.fromAttrID, mmv.Ref.toAttrID); err != nil {
			mWorker.Logger.Debugf("App: %s, FromMember: %s, FromAttr: %s, FromVal: %s, ToMember: %s, ToAttr: %s, ToVal: %s\n", mmv.Ref.appID, mmv.Ref.fromMemberID, mmv.Ref.fromAttrID, mmv.Ref.fromVal, mmv.Ref.toMemberID, mmv.Ref.toAttrID, mmv.Ref.toVal)
			mWorker.Logger.Debug(mmv.Ref)
			mWorker.Logger.Fatal(err)
			return err
		} else {
			color.Magenta.Printf("New Ref | App: %s, FromMember: %s, FromAttr: %s, FromVal: %s, ToMember: %s, ToAttr: %s, ToVal: %s\n", mmv.Ref.appID, mmv.Ref.fromMemberID, mmv.Ref.fromAttrID, mmv.Ref.fromVal, mmv.Ref.toMemberID, mmv.Ref.toAttrID, mmv.Ref.toVal)
		}

	}

	return nil
}

func (mWorker *MigrationWorker) HandleNodeMedia(toTable config.ToTable, data DataMap) error {
	if len(toTable.Media) > 0 {
		if filePathCol, ok := toTable.Media["path"]; ok {
			if filePath, ok := data[filePathCol]; ok {
				if err := mWorker.TransferMedia(fmt.Sprint(filePath)); err != nil {
					mWorker.Logger.Fatal("@MigrateNode > TransferMedia: ", err)
					return err
				}
			}
		} else {
			err := errors.New("@HandleNodeMedia > toTable.Media: Path not found in map")
			mWorker.Logger.Fatal(err)
			return err
		}
	}
	return nil
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
			// log.Printf("%s { %s } | Owner ID: %v \n", color.FgLightYellow.Render("BAG"), node.Tag.Name, mWorker.uid)
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
			var migrated bool

			bagManager := mWorker.CreateBagManager()

			for _, toTable := range mapping.ToTables {

				mappedMemberDatum := mWorker.GetMappedMemberData(toTable, node)

				if !mWorker.ValidateMappingConditions(toTable.Conditions, node.Data) {
					mWorker.Logger.Infof("%s: Mapping conditions not validated!\n", toTable.Table)
					continue
				} else {
					mWorker.Logger.Infof("%s: Mapping conditions validated!\n", toTable.Table)
				}

				if !mappedMemberDatum.ValidateMappedData() {
					mWorker.Logger.Infof("%s: Mapped data not validated!\n", toTable.Table)
					mWorker.Logger.Debug("mappedMemberDatum | ", mappedMemberDatum)
					continue
				} else {
					mWorker.Logger.Infof("%s: Mapped data validated!\n", toTable.Table)
				}

				mWorker.Logger.Infof("Cols before merging: %v\n", mappedMemberDatum.ToCols())
				if err := mWorker.MergeBagDataWithMappedData(&mappedMemberDatum, node, bagManager); err != nil {
					mWorker.Logger.Fatal(err)
				} else {
					mWorker.Logger.Infof("Cols after merging: %v\n", mappedMemberDatum.ToCols())
				}

				if err := mWorker.MigrateMemberData(mappedMemberDatum, node); err != nil {
					if mWorker.mtype == BAGS {
						return false, err
					}
					mWorker.Logger.Fatal(err)
				}

				if err := mWorker.HandleNodeMedia(toTable, node.Data); err != nil {
					mWorker.Logger.Fatal(err)
				}

				if err := mWorker.CreateIDAttributeRows(mappedMemberDatum, node.Data); err != nil {
					mWorker.Logger.Fatal(err)
				}

				if err := mWorker.CreateAttributeRows(mappedMemberDatum); err != nil {
					mWorker.Logger.Fatal(err)
				}

				if err := mWorker.CreateReferenceRows(mappedMemberDatum); err != nil {
					mWorker.Logger.Fatal(err)
				}

				if err := mWorker.CreateSelfReferenceRows(mappedMemberDatum); err != nil {
					mWorker.Logger.Fatal(err)
				}

				if err := mWorker.PushData(mappedMemberDatum, node); err != nil {
					mWorker.Logger.Fatal(err)
				}

				mappedMemberData = append(mappedMemberData, mappedMemberDatum)
				migrated = true
			}

			if migrated {
				bagManager.UpdateBags(mWorker.tx.StencilTx, mWorker.logTxn.Txn_id)
				node.DeleteMappedDataFromNode(mappedMemberData)
			}

			return migrated, nil
		} else {
			mWorker.Logger.Fatal("Waiting List Case!")
			return false, errors.New("Waiting List Case")
		}
	} else {
		if !strings.EqualFold(mWorker.mtype, DELETION) {
			return false, fmt.Errorf("no mapping found for node: %s", node.Tag.Name)
		}
		return false, mWorker.HandleUnmappedNode(node)
	}
}

func (mWorker *MigrationWorker) CheckNextNode(node *DependencyNode) error {

	if nextNodes, err := mWorker.GetAllNextNodes(node); err == nil {
		log.Println(fmt.Sprintf("NEXT NODES FETCHED { %s } : %d", node.Tag.Name, len(nextNodes)))
		for _, nextNode := range nextNodes {
			// log.Println(fmt.Sprintf("CURRENT NEXT NODE { %s > %s } %d/%d", node.Tag.Name, nextNode.Tag.Name, i, len(nextNodes)))
			// mWorker.AddToReferences(nextNode, node)
			if precedingNodes, err := mWorker.GetAllPreviousNodes(nextNode); err != nil {
				return err
			} else if len(precedingNodes) <= 1 {
				// log.Printf("precedingNodes <= 1 { %s } \n", node.Tag.Name)
				if err := mWorker.CheckNextNode(nextNode); err != nil {
					return err
				}
				log.Printf("SEND NEXT NODE TO BAG { %s } \n", nextNode.Tag.Name)
				if err := mWorker.SendNodeToBag(nextNode); err != nil {
					return err
				}
			} else {
				// log.Printf("precedingNodes > 1 { %s } \n", node.Tag.Name)
			}
		}
		return nil
	} else {
		return err
	}
}

func (mWorker *MigrationWorker) CloseDBConns() {

	mWorker.SrcAppConfig.CloseDBConns()
	mWorker.DstAppConfig.CloseDBConns()
	mWorker.logTxn.DBconn.Close()
	color.Info.Println("v2.mWorker.CloseDBConns: Closing DB connection for stencil in logTxn")
}

func (mWorker *MigrationWorker) RenewDBConn(isBlade ...bool) {
	mWorker.CloseDBConns()
	mWorker.logTxn.DBconn.Close()
	color.Info.Println("v2.mWorker.RenewDBConn: Closing DB connection for stencil in logTxn")
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
					joinStr += fmt.Sprintf("data_bags \"%s\"", fromTable)
					idCols += fmt.Sprintf("\"%s\".pk,", fromTable)
					cols += fmt.Sprintf(" coalesce(\"%s\".\"data\"::jsonb, '{}'::jsonb)  ||", fromTable)
				}
				for toTable, conditions := range toTablesMap {
					if conditions != nil {
						conditions = append(conditions, joinMap[toTable][fromTable]...)
						if joinMap[toTable][fromTable] != nil {
							joinMap[toTable][fromTable] = nil
						}
						if _, ok := seenMap[toTable]; !ok {
							joinStr += fmt.Sprintf(" FULL JOIN data_bags \"%s\" ", toTable)
						}
						joinStr += fmt.Sprintf(" ON \"%s\".member = %s AND \"%s\".member = %s AND %s ", fromTable, tableIDs[fromTable], toTable, tableIDs[toTable], strings.Join(conditions, " AND "))
						cols += fmt.Sprintf(" coalesce(\"%s\".\"data\"::jsonb, '{}'::jsonb)  ||", toTable)
						idCols += fmt.Sprintf("\"%s\".pk,", toTable)
						seenMap[toTable] = true
					}
				}
				seenMap[fromTable] = true
			}
			sql = fmt.Sprintf(sql, strings.Trim(idCols, ","), strings.Trim(cols, ",|"), joinStr)
		} else {
			table := tag.Members["member1"]
			joinStr := fmt.Sprintf("data_bags \"%s\"", table)
			idCols := fmt.Sprintf("\"%s\".pk", table)
			cols := fmt.Sprintf(" coalesce(\"%s\".\"data\"::jsonb, '{}'::jsonb)  ", table)
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
	mWorker.tx.SrcTx.Rollback()
	mWorker.tx.DstTx.Rollback()
	mWorker.tx.StencilTx.Rollback()
	return nil
}

func (mWorker *MigrationWorker) FetchMappingsForBag(srcApp, srcAppID, dstApp, dstAppID, srcMember, dstMember string) (config.Mapping, bool) {

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

	return combinedMapping, mappingFound
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

func (mWorker *MigrationWorker) GetUserIDAppIDFromPreviousMigration(currentAppID, currentUID string) (*App, string, error) {

	currentRootMemberID := db.GetAppRootMemberID(mWorker.logTxn.DBconn, currentAppID)

	currentUIDInt, err := strconv.ParseInt(currentUID, 10, 64)
	if err != nil {
		panic(err)
	}

	mWorker.Logger.Infof("Getting previous migration | App: '%v', UID: '%v', rootMemberID: '%v' \n", currentAppID, currentUIDInt, currentRootMemberID)

	if IDRows, err := mWorker.GetRowsFromAttrTable(currentAppID, currentRootMemberID, currentUIDInt, false); err == nil {
		// mWorker.Logger.Info("Fetched AttrRows: ", len(IDRows))
		for _, IDRow := range IDRows {
			// fmt.Println(IDRow)
			prevRootMemberID := db.GetAppRootMemberID(mWorker.logTxn.DBconn, IDRow.FromAppID)
			if strings.EqualFold(IDRow.FromMemberID, prevRootMemberID) {
				mWorker.Logger.Infof("Previous migration found | App: '%v', UID: '%v', rootMemberID: '%v' \n", IDRow.FromAppID, IDRow.FromID, IDRow.FromMemberID)
				if appName, err := db.GetAppNameByAppID(mWorker.logTxn.DBconn, IDRow.FromAppID); err != nil {
					mWorker.Logger.Fatal(err)
				} else {
					return &App{Name: appName, ID: IDRow.FromAppID}, fmt.Sprint(IDRow.FromID), nil
				}

			}
		}

		mWorker.Logger.Infof("No previous migration found | App: '%v', UID: '%v', rootMemberID: '%v' \n", currentAppID, currentUIDInt, currentRootMemberID)

		return nil, "", nil
	} else {
		mWorker.Logger.Fatalf("@GetUserIDAppIDFromPreviousMigration | App: '%s', UID: '%v', rootMemberID: '%s' | err => %v \n", currentAppID, currentUIDInt, currentRootMemberID, err)
		return nil, "", fmt.Errorf("no previous migration user and app id found for => currentAppID: %s, currentUID: %v", currentAppID, currentUIDInt)
	}
}

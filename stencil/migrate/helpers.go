package migrate

import (
	"fmt"
	"log"
	"os"
	"stencil/config"
	"stencil/db"
	"stencil/helper"
	"stencil/transaction"
	"strconv"
	"strings"

	logg "github.com/withmandala/go-log"
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
		visitedNodes: make(map[string]map[string]bool),
		Logger:       logg.New(os.Stderr)}
	// if err := mWorker.FetchRoot(threadID); err != nil {
	// 	log.Fatal(err)
	// }
	mWorker.FTPClient = GetFTPClient()
	mWorker.Logger.WithTimestamp()
	mWorker.Logger.WithColor()
	mWorker.Logger.WithDebug()
	mWorker.Logger.Infof("Bag Worker Created for thread: %v", threadID)
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
		visitedNodes: make(map[string]map[string]bool),
		Logger:       logg.New(os.Stderr)}

	if err := mWorker.FetchRoot(threadID); err != nil {
		mWorker.Logger.Fatal(err)
	}
	mWorker.FTPClient = GetFTPClient()
	mWorker.Logger.WithTimestamp()
	mWorker.Logger.WithColor()
	mWorker.Logger.WithDebug()
	mWorker.Logger.Infof("Worker Created for thread: %v", threadID)

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

func (self *MigrationWorkerV2) ExcludeVisited(tag config.Tag) string {
	visited := ""
	for _, tagMember := range tag.Members {
		if memberIDs, ok := self.visitedNodes[tagMember]; ok {
			pks := ""
			for pk := range memberIDs {
				if len(pk) > 0 {
					pks += pk + ","
				}
			}
			if pks != "" {
				pks = strings.Trim(pks, ",")
				visited += fmt.Sprintf(" AND %s.id NOT IN (%s) ", tagMember, pks)
			}

		}
	}
	return visited
}

func (self *MigrationWorkerV2) GetMemberDataFromNode(member string, nodeData map[string]interface{}) map[string]interface{} {
	memberData := make(map[string]interface{})
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

func (self *MigrationWorkerV2) GetTagQL(tag config.Tag) string {

	sql := "SELECT %s FROM %s "

	if len(tag.InnerDependencies) > 0 {
		cols := ""
		joinMap := tag.CreateInDepMap()
		seenMap := make(map[string]bool)
		joinStr := ""
		for fromTable, toTablesMap := range joinMap {
			if _, ok := seenMap[fromTable]; !ok {
				joinStr += fmt.Sprintf("\"%s\"", fromTable)
				_, colStr := db.GetColumnsForTable(self.SrcDBConn, fromTable)
				cols += colStr + ","
			}
			for toTable, conditions := range toTablesMap {
				if conditions != nil {
					conditions = append(conditions, joinMap[toTable][fromTable]...)
					if joinMap[toTable][fromTable] != nil {
						joinMap[toTable][fromTable] = nil
					}
					joinStr += fmt.Sprintf(" FULL JOIN \"%s\" ON %s ", toTable, strings.Join(conditions, " AND "))
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

func (self *MigrationWorkerV2) GetTagQLForBag(tag config.Tag) string {

	if tableIDs, err := tag.MemberIDs(self.logTxn.DBconn, self.SrcAppConfig.AppID); err != nil {
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
						joinStr += fmt.Sprintf(" FULL JOIN data_bags %s ON %s.member = %s AND %s.member = %s AND %s ", toTable, fromTable, tableIDs[fromTable], toTable, tableIDs[toTable], strings.Join(conditions, " AND "))
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

func (self *MigrationWorkerV2) RollbackTransactions() error {
	// log.Fatal("@CommitTransactions: About to Commit!")
	if err := self.tx.SrcTx.Rollback(); err != nil {
		log.Fatal("Error rolling back Source DB Transaction! ", err)
		return err
	}
	if err := self.tx.DstTx.Rollback(); err != nil {
		log.Fatal("Error rolling back Destination DB Transaction! ", err)
		return err
	}
	if err := self.tx.StencilTx.Rollback(); err != nil {
		log.Fatal("Error rolling back Stencil DB Transaction! ", err)
		return err
	}
	return nil
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
		if nodeVal, ok := node.Data[idCol]; ok {
			if nodeVal == nil {
				continue
			}
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

func (self *MigrationWorkerV2) FetchMappingsForBag(srcApp, srcAppID, dstApp, dstAppID, srcMember, dstMember string) (config.MappedApp, config.Mapping, bool) {

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
	return appMappings, combinedMapping, mappingFound
}

func (self *MigrationWorkerV2) CleanMappingAttr(attr string) string {
	cleanedAttr := strings.ReplaceAll(attr, "(", "")
	cleanedAttr = strings.ReplaceAll(cleanedAttr, ")", "")
	cleanedAttr = strings.ReplaceAll(cleanedAttr, "#ASSIGN", "")
	cleanedAttr = strings.ReplaceAll(cleanedAttr, "#FETCH", "")
	cleanedAttr = strings.ReplaceAll(cleanedAttr, "#REFHARD", "")
	cleanedAttr = strings.ReplaceAll(cleanedAttr, "#REF", "")
	return cleanedAttr
}

func (self *MigrationWorkerV2) FetchMappedAttribute(srcApp, srcAppID, dstApp, dstAppID, srcMember, dstMember, dstAttr string) (string, bool) {

	var appMappings config.MappedApp
	if srcApp == dstApp {
		appMappings = *config.GetSelfSchemaMappings(self.logTxn.DBconn, srcAppID, srcApp)
	} else {
		appMappings = *config.GetSchemaMappingsFor(srcApp, dstApp)
	}
	for _, mapping := range appMappings.Mappings {
		if mappedTables := helper.IntersectString([]string{srcMember}, mapping.FromTables); len(mappedTables) > 0 {
			for _, toTableMapping := range mapping.ToTables {
				if strings.EqualFold(dstMember, toTableMapping.Table) {
					for toAttr, fromAttr := range toTableMapping.Mapping {
						if toAttr == dstAttr {
							cleanedAttr := self.CleanMappingAttr(fromAttr)
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

func (self *MigrationWorkerV2) GetUserIDAppIDFromPreviousMigration(currentAppID, currentUID string) (string, string, error) {

	rootMemberID := db.GetAppRootMemberID(self.logTxn.DBconn, currentAppID)

	currentUIDInt, err := strconv.ParseInt(currentUID, 10, 64)
	if err != nil {
		panic(err)
	}

	fmt.Printf("@GetUserIDAppIDFromPreviousMigration | Getting previous migration | App: '%v', UID: '%v', rootMemberID: '%v' \n", currentAppID, currentUIDInt, rootMemberID)

	if IDRows, err := self.GetRowsFromIDTable(currentAppID, rootMemberID, currentUIDInt, false); err == nil {
		fmt.Println(IDRows)
		if len(IDRows) > 0 {
			for _, IDRow := range IDRows {
				prevRootMemberID := db.GetAppRootMemberID(self.logTxn.DBconn, IDRow.FromAppID)
				if IDRow.FromMemberID == prevRootMemberID {
					fmt.Printf("@GetUserIDAppIDFromPreviousMigration | Previous migration found | App: '%v', UID: '%v', rootMemberID: '%v' \n", IDRows[0].FromAppID, IDRows[0].FromID, IDRows[0].FromMemberID)
					return IDRows[0].FromAppID, fmt.Sprint(IDRows[0].FromID), nil
				}
			}
		}
		fmt.Printf("@GetUserIDAppIDFromPreviousMigration | No previous migration found | App: '%v', UID: '%v', rootMemberID: '%v' \n", currentAppID, currentUIDInt, rootMemberID)
		return "", "", nil
	} else {
		log.Fatalf("@GetUserIDAppIDFromPreviousMigration | App: '%s', UID: '%v', rootMemberID: '%s' | err => %v \n", currentAppID, currentUIDInt, rootMemberID, err)
		return "", "", fmt.Errorf("no previous migration user and app id found for => currentAppID: %s, currentUID: %v", currentAppID, currentUIDInt)
	}
}

package migrate_v2

import (
	"fmt"
	"log"
	config "stencil/config/v2"
	"stencil/db"
	"stencil/helper"
	"strconv"
	"strings"
)

func (self *MigrationWorker) CloseDBConns() {

	self.SrcAppConfig.CloseDBConns()
	self.DstAppConfig.CloseDBConns()
}

func (self *MigrationWorker) RenewDBConn(isBlade ...bool) {
	self.CloseDBConns()
	self.logTxn.DBconn.Close()
	self.logTxn.DBconn = db.GetDBConn(db.STENCIL_DB)
	self.SrcAppConfig.DBConn = db.GetDBConn(self.SrcAppConfig.AppName)
	self.SrcAppConfig.DBConn = db.GetDBConn(self.DstAppConfig.AppName, isBlade...)
}

func (self *MigrationWorker) UserID() string {
	return self.uid
}

func (self *MigrationWorker) MigrationID() int {
	return self.logTxn.Txn_id
}

func (self *MigrationWorker) GetMemberDataFromNode(member string, nodeData map[string]interface{}) map[string]interface{} {
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

func (self *MigrationWorker) GetTagQL(tag config.Tag) string {

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
				_, colStr := db.GetColumnsForTable(self.SrcAppConfig.DBConn, fromTable)
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
					_, colStr := db.GetColumnsForTable(self.SrcAppConfig.DBConn, toTable)
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

func (self *MigrationWorker) GetTagQLForBag(tag config.Tag) string {

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

func (self *MigrationWorker) InitTransactions() error {
	var err error
	self.tx.SrcTx, err = self.SrcAppConfig.DBConn.Begin()
	if err != nil {
		log.Fatal("Error creating Source DB Transaction! ", err)
		return err
	}
	self.tx.DstTx, err = self.DstAppConfig.DBConn.Begin()
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

func (self *MigrationWorker) CommitTransactions() error {
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

func (self *MigrationWorker) RollbackTransactions() error {
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

func (self *MigrationWorker) FetchMappingsForBag(srcApp, srcAppID, dstApp, dstAppID, srcMember, dstMember string) (config.MappedApp, config.Mapping, bool) {

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

func (self *MigrationWorker) CleanMappingAttr(attr string) string {
	cleanedAttr := strings.ReplaceAll(attr, "(", "")
	cleanedAttr = strings.ReplaceAll(cleanedAttr, ")", "")
	cleanedAttr = strings.ReplaceAll(cleanedAttr, "#ASSIGN", "")
	cleanedAttr = strings.ReplaceAll(cleanedAttr, "#FETCH", "")
	cleanedAttr = strings.ReplaceAll(cleanedAttr, "#REFHARD", "")
	cleanedAttr = strings.ReplaceAll(cleanedAttr, "#REF", "")
	return cleanedAttr
}

func (self *MigrationWorker) FetchMappedAttribute(srcApp, srcAppID, dstApp, dstAppID, srcMember, dstMember, dstAttr string) (string, bool) {

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

func (self *MigrationWorker) FetchMappingsForNode(node *DependencyNode) (config.Mapping, bool) {
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

func (self *MigrationWorker) GetUserIDAppIDFromPreviousMigration(currentAppID, currentUID string) (string, string, error) {

	currentRootMemberID := db.GetAppRootMemberID(self.logTxn.DBconn, currentAppID)

	currentUIDInt, err := strconv.ParseInt(currentUID, 10, 64)
	if err != nil {
		panic(err)
	}

	fmt.Printf("@GetUserIDAppIDFromPreviousMigration | Getting previous migration | App: '%v', UID: '%v', rootMemberID: '%v' \n", currentAppID, currentUIDInt, currentRootMemberID)

	if IDRows, err := self.GetRowsFromIDTable(currentAppID, currentRootMemberID, currentUIDInt, false); err == nil {
		fmt.Println(IDRows)
		if len(IDRows) > 0 {
			for _, IDRow := range IDRows {
				prevRootMemberID := db.GetAppRootMemberID(self.logTxn.DBconn, IDRow.FromAppID)
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

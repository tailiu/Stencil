package SA1_display

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"stencil/common_funcs"
	"stencil/config"
	"stencil/db"
	"stencil/reference_resolution_v2"
	"stencil/schema_mappings"
	"strconv"
	"strings"

	"github.com/gookit/color"
)

func CreateDisplayConfig(migrationID int, resolveReference, useBladeServerAsDst,
	displayInFirstPhase, markAsDelete bool) *display {

	var display display

	var srcAppConfig srcAppConfig

	var dstAppConfig dstAppConfig

	stencilDBConn := db.GetDBConn("stencil")

	srcAppID, dstAppID, srcUserID :=
		common_funcs.GetSrcDstAppIDsUserIDByMigrationID(stencilDBConn, migrationID)

	srcAppName := common_funcs.GetAppNameByAppID(stencilDBConn, srcAppID)
	dstAppName := common_funcs.GetAppNameByAppID(stencilDBConn, dstAppID)

	allMappings, err1 := config.LoadSchemaMappings()
	if err1 != nil {
		log.Fatal(err1)
	}

	mappingsFromSrcToDst, err2 :=
		schema_mappings.GetToAppMappings(allMappings, srcAppName, dstAppName)

	if err2 != nil {
		log.Fatal(err2)
	}

	allApps := schema_mappings.GetAllApps(allMappings)

	mappingsFromOtherAppsToDst := make(map[string]*config.MappedApp)

	for _, app := range allApps {

		if app != dstAppName {

			toDstMapping, err3 := schema_mappings.GetToAppMappings(allMappings, app, dstAppName)
			if err3 != nil {
				log.Fatal(err3)
			}

			mappingsFromOtherAppsToDst[app] = toDstMapping
		}

	}

	dstDAG, err4 := common_funcs.LoadDAG(dstAppName)
	if err4 != nil {
		log.Fatal(err4)
	}

	var dstDBConn *sql.DB

	if useBladeServerAsDst {
		dstDBConn = db.GetDBConn(dstAppName, true)
	} else {
		dstDBConn = db.GetDBConn(dstAppName)
	}

	// Note that when the display thread is initializing, dstUserID could be nil
	// because the data has not been migrated yet
	dstRootMember, dstRootAttr, dstUserID := getDstRootMemberAttrID(
		stencilDBConn, dstAppID, migrationID, dstDAG)

	dstAppTableIDNamePairs := make(map[string]string)
	dstAppTableNameIDPairs := make(map[string]string)

	dstRes := common_funcs.GetTableIDNamePairsInApp(stencilDBConn, dstAppID)

	for _, dstRes1 := range dstRes {

		dstAppTableIDNamePairs[fmt.Sprint(dstRes1["pk"])] =
			fmt.Sprint(dstRes1["table_name"])

		dstAppTableNameIDPairs[fmt.Sprint(dstRes1["table_name"])] =
			fmt.Sprint(dstRes1["pk"])
	}

	srcAppTableNameIDPairs := make(map[string]string)

	srcRes := common_funcs.GetTableIDNamePairsInApp(stencilDBConn, srcAppID)

	for _, srcRes1 := range srcRes {
		srcAppTableNameIDPairs[fmt.Sprint(srcRes1["table_name"])] = fmt.Sprint(srcRes1["pk"])
	}

	appIDNamePairs := common_funcs.GetAppIDNamePairs(stencilDBConn)
	tableIDNamePairs := common_funcs.GetTableIDNamePairs(stencilDBConn)
	attrIDNamePairs := GetAttrIDNamePairs(stencilDBConn)
	dstAppColNameIDPairs := getAttrNameIDPairsInApp(stencilDBConn, dstAppID)

	appTableNameTableIDPairs := getAppTableNameTableIDPairs(stencilDBConn, appIDNamePairs)

	refResolution := reference_resolution_v2.InitializeReferenceResolution(
		migrationID, dstAppID, dstAppName, dstDBConn, stencilDBConn,
		dstAppTableNameIDPairs, appIDNamePairs, tableIDNamePairs,
		attrIDNamePairs, dstAppColNameIDPairs, allMappings, dstDAG,
	)

	srcAppConfig.appID = srcAppID
	srcAppConfig.appName = srcAppName
	srcAppConfig.userID = srcUserID
	srcAppConfig.tableNameIDPairs = srcAppTableNameIDPairs

	dstAppConfig.appID = dstAppID
	dstAppConfig.appName = dstAppName
	dstAppConfig.tableNameIDPairs = dstAppTableNameIDPairs
	dstAppConfig.colNameIDPairs = dstAppColNameIDPairs
	dstAppConfig.rootTable = dstRootMember
	dstAppConfig.rootAttr = dstRootAttr
	dstAppConfig.userID = dstUserID
	dstAppConfig.dag = dstDAG
	dstAppConfig.DBConn = dstDBConn
	dstAppConfig.ownershipDisplaySettingsSatisfied = false

	display.stencilDBConn = stencilDBConn
	display.appIDNamePairs = appIDNamePairs
	display.tableIDNamePairs = tableIDNamePairs
	display.attrIDNamePairs = attrIDNamePairs
	display.migrationID = migrationID
	display.resolveReference = resolveReference
	display.srcAppConfig = &srcAppConfig
	display.dstAppConfig = &dstAppConfig
	display.rr = refResolution
	display.displayInFirstPhase = displayInFirstPhase
	display.markAsDelete = markAsDelete
	display.mappingsFromSrcToDst = mappingsFromSrcToDst
	display.mappingsFromOtherAppsToDst = mappingsFromOtherAppsToDst
	display.appTableNameTableIDPairs = appTableNameTableIDPairs

	return &display

}

func closeDBConn(conn *sql.DB) {

	err := conn.Close()
	if err != nil {
		log.Fatal(err)
	}

}

func (display *display) closeDBConns() {

	log.Println("Close db connections in the SA1 display thread")

	closeDBConn(display.stencilDBConn)
	closeDBConn(display.dstAppConfig.DBConn)

}

func getDstRootMemberAttrID(stencilDBConn *sql.DB,
	appID string, migrationID int, dstDAG *common_funcs.DAG) (string, string, string) {

	// log.Println(*dstDAG)

	dstRootMember, dstRootAttr, err2 := dstDAG.GetRootMemberAttr()
	if err2 != nil {
		log.Fatal(err2)
	}

	tableID := getTableIDByTableName(stencilDBConn, appID, dstRootMember)

	// Since the in current settings, there is only one row and the root attribute is always id,
	// we only do in the following way. Note that this is not a generic way.
	query := fmt.Sprintf(
		`SELECT id FROM display_flags WHERE app_id = %s 
		and table_id = %s and migration_id = %d`,
		appID, tableID, migrationID)

	// log.Println(query)

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if data["id"] == nil {
		return dstRootMember, dstRootAttr, ""
	} else {
		return dstRootMember, dstRootAttr, fmt.Sprint(data["id"])
	}

}

func (display *display) GetUndisplayedMigratedData() []*HintStruct {

	var displayHints []*HintStruct

	query := fmt.Sprintf(`SELECT table_id, id FROM display_flags
		 WHERE app_id = %s and migration_id = %d and display_flag = true`,
		display.dstAppConfig.appID, display.migrationID)

	data := db.GetAllColsOfRows(display.stencilDBConn, query)
	// fmt.Println(data)

	for _, data1 := range data {

		displayHints = append(displayHints,
			TransformDisplayFlagDataToHint(display, data1))

	}
	// fmt.Println(displayHints)

	return displayHints

}

func (display *display) CheckMigrationComplete() bool {

	query := fmt.Sprintf(
		`SELECT 1 FROM txn_logs 
		WHERE action_id = %d and action_type='COMMIT' LIMIT 1`,
		display.migrationID)

	data := db.GetAllColsOfRows(display.stencilDBConn, query)

	if len(data) == 0 {

		return false

	} else {

		return true

	}

}

func CheckMigrationComplete1(stencilDBConn *sql.DB,
	migrationID int) bool {

	query := fmt.Sprintf(
		`SELECT 1 FROM txn_logs 
		WHERE action_id = %d and action_type='COMMIT' LIMIT 1`,
		migrationID)

	data := db.GetAllColsOfRows(stencilDBConn, query)

	if len(data) == 0 {

		return false

	} else {

		return true

	}

}

func getAppNameIDPairs(stencilDBConn *sql.DB) map[string]string {

	query := fmt.Sprintf(`SELECT pk, app_name FROM apps`)

	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	res := make(map[string]string)

	for _, data1 := range data {
		res[fmt.Sprint(data1["app_name"])] = fmt.Sprint(data1["pk"])
	}

	return res

}

func getAppTableNameTableIDPairs(stencilDBConn *sql.DB,
	appIDNamePairs map[string]string) map[string]string {

	query := fmt.Sprintf(`SELECT pk, app_id, table_name FROM app_tables`)

	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	res := make(map[string]string)

	for _, data1 := range data {

		tableID := fmt.Sprint(data1["pk"])
		appID := fmt.Sprint(data1["app_id"])
		tableName := fmt.Sprint(data1["table_name"])

		appName := appIDNamePairs[appID]

		res[appName+":"+tableName] = tableID

	}

	return res

}

func (display *display) isDataMigratedAndAlreadyDisplayed(dataHint *HintStruct) bool {

	query := fmt.Sprintf(
		`SELECT * FROM display_flags WHERE
		app_id = %s and table_id = %s and id = %d and display_flag = false`,
		display.dstAppConfig.appID,
		dataHint.TableID,
		dataHint.KeyVal["id"],
	)

	data := db.GetAllColsOfRows(display.stencilDBConn, query)

	if len(data) == 0 {
		return true
	} else {
		return false
	}

}

func (display *display) isDataNotMigratedAndAlreadyDisplayed(dataHint *HintStruct) bool {

	query := fmt.Sprintf(
		`SELECT * FROM display_flags
		WHERE app_id = %s and table_id = %s and id = %d`,
		display.dstAppConfig.appID,
		dataHint.TableID,
		dataHint.KeyVal["id"],
	)

	data := db.GetAllColsOfRows(display.stencilDBConn, query)

	if len(data) == 0 {
		return true
	} else {
		return false
	}

}

func (display *display) getAttributesToSetAsSTENCILNULLs(dataHint *HintStruct) []string {

	table := dataHint.Table
	tableID := dataHint.TableID
	id := fmt.Sprint(dataHint.Data["id"])
	
	attrsToBeUpdated := display.getAllAttributesToBeUpdated(table)
	
	log.Println("attributes to be updated:", attrsToBeUpdated)

	// Even though references have been updated in checking dependencies and ownership,
	// update reference the last time before displaying the data
	for _, attrToBeUpdated := range attrsToBeUpdated {
		colID := display.dstAppConfig.colNameIDPairs[table + ":" + attrToBeUpdated]
		attr := reference_resolution_v2.CreateAttribute(display.dstAppConfig.appID, tableID, colID, id, id)
		display.rr.ResolveReference(attr)
	}

	updatedAttrs := display.rr.GetUpdatedAttributes(tableID, id)

	log.Println("attributes has already been updated:", updatedAttrs)

	var attrsToBeSetToNULLs []string

	for _, attr := range attrsToBeUpdated {
		if _, ok := updatedAttrs[attr]; !ok {
			attrsToBeSetToNULLs = append(attrsToBeSetToNULLs, attr)
		}
	}

	return attrsToBeSetToNULLs
}

// Before displaying a piece of data, the unresolved references should be set to be NULLs
// Setting those references to be NULLs does not influence how data behaves in the application,
// but in the next migration it may influence how Stencil migrates data
// For example,
// If the reply field of  a status which was a comment in Diaspora is set to be NULLs because
// the corresponding post is not migrated, then in the next migration due the reply field is NULL,
// it may be migrated as a post. We might use a special field indicating that this field has been
// set to NULLs before and both applications and Stencil should be aware of that
// to solve this problem, but since setting unresolved references to be NULLs is enough
// in our current testing applications, we just use this simple method.
func (display *display) Display(dataHints []*HintStruct) error {

	var queries1, queries2 []string

	var query1, query2, query3 string

	for _, dataHint := range dataHints {

		// we try to display data in the UNIT of a node
		// Even if in the case where some data belonging to other migration should not be
		// displayed, that data is displayed as a unit of data
		// e.g., some data in a tweet is migrated in another migration which is still in the first phase,
		// and according to the rule, it can only be displayed in the first phase
		// if the parent tweet is there, but the parent tweet is not there for now.
		// other data in the tweet is going to be displayed in the migration in the second phase.
		// In this case, all the data will be displayed
		if !display.isNodeInCurrentMigration(dataHint) {
			log.Println(`This data is not migrated by the currently checked user,
				but this display thread will also display this data`)
		}

		log.Println("Set unresolved attributes to be STENCIL_NULLs before displaying the data")

		attrsToBeSetToNULLs := display.getAttributesToSetAsSTENCILNULLs(dataHint)

		query1 = fmt.Sprintf(`UPDATE "%s" SET display_flag = false`, dataHint.Table)

		for _, attr := range attrsToBeSetToNULLs {
			query1 += ", " + attr + " = '-1'"
		}

		query1Where := fmt.Sprintf(" WHERE id = %d;", dataHint.KeyVal["id"])

		query1 += query1Where

		query2 = fmt.Sprintf(
			`UPDATE Display_flags SET 
			display_flag = false, updated_at = now() 
			WHERE app_id = %s and table_id = %s and id = %d;`,
			display.dstAppConfig.appID, dataHint.TableID, dataHint.KeyVal["id"])

		query3 = fmt.Sprintf(
			`UPDATE evaluation SET
			displayed_at = now() WHERE migration_id = '%d' and
			src_app = '%s' and dst_app = '%s' and
			dst_table = '%s' and dst_id = '%d'`,
			display.migrationID, display.srcAppConfig.appID,
			display.dstAppConfig.appID, dataHint.TableID, dataHint.KeyVal["id"])

		log.Println("**************************************")
		log.Println("Update the application table:")
		log.Println(query1)
		log.Println("Attributes need to be set as NULL:")
		log.Println(attrsToBeSetToNULLs)
		log.Println("Update the display_flags table:")
		log.Println(query2)
		log.Println("Update the evaluation table:")
		log.Println(query3)
		log.Println("**************************************")

		queries1 = append(queries1, query1)
		queries2 = append(queries2, query2, query3)

	}

	if err1 := db.TxnExecute(display.dstAppConfig.DBConn, queries1); err1 != nil {
		return err1
	} else {
		if err2 := db.TxnExecute(display.stencilDBConn, queries2); err2 != nil {
			return err2
		} else {
			return nil
		}
	}

}

func (display *display) CheckDisplay(dataHint *HintStruct) bool {

	query := fmt.Sprintf(
		`SELECT display_flag from "%s" where id = %d`,
		dataHint.Table, dataHint.KeyVal["id"],
	)

	data1, err := db.DataCall1(display.dstAppConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data1)

	return !data1["display_flag"].(bool)

}

func (display *display) getAllAttributesToBeUpdated(table string) []string {
	
	attrsToBeUpdatedBasedOnMappings := display.getAllMappedAttributesContainingREFInMappingsFromAllApps(table)
	
	attrsToBeUpdatedBasedOnDag := display.dstAppConfig.dag.GetAllAttrsDepsOnBasedOnDag(table)
	
	combinedAttrs := append(attrsToBeUpdatedBasedOnMappings, attrsToBeUpdatedBasedOnDag...)

	return common_funcs.RemoveDuplicateElementsInSlice(combinedAttrs)
	
}

func (display *display) getAllMappedAttributesContainingREFInMappingsFromAllApps(table string) []string {

	var attrs []string

	for _, mapping := range display.mappingsFromOtherAppsToDst {

		attrsInApp := schema_mappings.GetAllMappedAttributesContainingREFInMappings(
			mapping, table,
		)

		for _, attrInApp := range attrsInApp {
			if !common_funcs.ExistsInSlice(attrs, attrInApp) {
				attrs = append(attrs, attrInApp)
			}
		}

	}

	return attrs
}

func getAttrNameByAttrID(stencilDBConn *sql.DB, attrID string) string {
	//Need to change
	query := fmt.Sprintf("select column_name from app_schemas where pk = %s", attrID)

	// log.Println(query)

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprint(data["app_name"])
}

func GetAttrIDNamePairs(stencilDBConn *sql.DB) map[string]string {

	attrIDNamePairs := make(map[string]string)

	query := fmt.Sprintf("select column_name, pk from app_schemas")

	// log.Println(query)

	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	for _, data1 := range data {
		attrIDNamePairs[fmt.Sprint(data1["pk"])] = fmt.Sprint(data1["column_name"])
	}

	return attrIDNamePairs
}

func getAttrNameIDPairsInApp(stencilDBConn *sql.DB, appID string) map[string]string {

	attrNameIDPairs := make(map[string]string)

	query := fmt.Sprintf(
		`SELECT t.table_name, s.column_name, s.pk FROM app_schemas as s JOIN app_tables as t ON
		s.table_id = t.pk WHERE t.app_id = %s`, appID,
	)

	// log.Println(query)

	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	for _, data1 := range data {
		tableNameColumnName := fmt.Sprint(data1["table_name"]) + ":" + fmt.Sprint(data1["column_name"])
		columnID := fmt.Sprint(data1["pk"])
		attrNameIDPairs[tableNameColumnName] = columnID
	}

	return attrNameIDPairs
}

func getTableIDByTableName(stencilDBConn *sql.DB, appID, tableName string) string {

	query := fmt.Sprintf("select pk from app_tables where app_id = %s and table_name = '%s'",
		appID, tableName)

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprint(data["pk"])

}

func ConvertMapToJSONString(data map[string]interface{}) string {

	convertedData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return string(convertedData)

}

func (display *display) chechPutIntoDataBag(secondRound bool, dataHints []*HintStruct) error {

	if secondRound {
		err9 := display.putIntoDataBag(dataHints)
		if err9 != nil {
			log.Fatal(err9)
		}
		return common_funcs.NoDataInNodeCanBeDisplayed
	} else {
		return common_funcs.NoDataInNodeCanBeDisplayed
	}
}

func (display *display) setDstUserIDIfNotSet() {

	if display.dstAppConfig.userID != "" {
		return
	}

	// Since the in current settings, there is only one row and the root attribute is always id,
	// we only do in the following way. Note that this is not a generic way.
	query := fmt.Sprintf(
		`SELECT id FROM display_flags WHERE app_id = %s 
		and table_id = %s and migration_id = %d`,
		display.dstAppConfig.appID,
		display.dstAppConfig.tableNameIDPairs[display.dstAppConfig.rootTable],
		display.migrationID)

	log.Println(query)

	data, err := db.DataCall1(display.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if data["id"] == nil {
		log.Fatal("This is the second phase, so the user id should not be nil!")
	} else {
		display.dstAppConfig.userID = fmt.Sprint(data["id"])
	}

}

func addTableNameToDataAttibutes(data map[string]interface{}, tableName string) map[string]interface{} {

	procData := make(map[string]interface{})

	for attr, val := range data {
		procAttr := tableName + "." + attr
		procData[procAttr] = val
	}

	return procData

}

// When putting data to dag bags, it does not matter whether we set unresolved references
// to NULLs or not, so we don't set those as NULLs.
func (display *display) putIntoDataBag(dataHints []*HintStruct) error {

	display.setDstUserIDIfNotSet()

	var queries1, queries2, queries3 []string

	var q1, q2, q3 string

	log.Println("Going to put data into data bags")

	for _, dataHint := range dataHints {

		// we try to avoid putting data into data bags if it is not in the
		// currently checked user's migration
		if !display.isNodeInCurrentMigration(dataHint) {
			log.Println(`This data is not migrated by the currently checked user,
				so it will not be put into data bags by this display thread`)
			continue
		}

		// dataHint.Data could be nil, which means there is no data,
		// if a thread crashes before executing queries3
		// and after executing queries1 and queries2, or data is deleted by services.
		// In both cases, there is no need to execute queries1 and queries2 again.
		if dataHint.Data != nil {

			attrsToBeSetToNULLs := display.getAttributesToSetAsSTENCILNULLs(dataHint)

			var STENCIL_NULL interface{}

			STENCIL_NULL = "-1"			

			for attr := range dataHint.Data {
				if common_funcs.ExistsInSlice(attrsToBeSetToNULLs, attr) {
					dataHint.Data[attr] = STENCIL_NULL
				}
			}

			procData := addTableNameToDataAttibutes(dataHint.Data, dataHint.Table)

			q1 = fmt.Sprintf(
				`INSERT INTO data_bags 
				(app, member, id, data, user_id, migration_id) VALUES 
				(%s, %s, %d, '%s', %s, %d)`,
				display.dstAppConfig.appID,
				dataHint.TableID,
				dataHint.KeyVal["id"],
				ConvertMapToJSONString(procData),
				display.dstAppConfig.userID,
				display.migrationID,
			)

			if !display.markAsDelete {
				q2 = fmt.Sprintf(
					`DELETE FROM "%s" WHERE id = %d`,
					dataHint.Table, dataHint.KeyVal["id"],
				)
			} else {
				q2 = fmt.Sprintf(
					`UPDATE "%s" SET mark_as_delete = true WHERE id = %d`,
					dataHint.Table, dataHint.KeyVal["id"],
				)
			}

		}

		q3 = fmt.Sprintf(
			`UPDATE display_flags SET 
			display_flag = false, updated_at = now() 
			WHERE app_id = %s and table_id = %s and id = %d;`,
			display.dstAppConfig.appID,
			dataHint.TableID, dataHint.KeyVal["id"],
		)

		log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
		log.Println("INSERT INTO data_bags:", q1)
		if !display.markAsDelete {
			log.Println("DELETE FROM the application:", q2)
		} else {
			log.Println("MARK AS DELETE in the application:", q2)
		}
		log.Println("UPDATE display_flags:", q3)
		log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")

		queries1 = append(queries1, q1)
		queries2 = append(queries2, q2)
		queries3 = append(queries3, q3)

	}

	// Since queries1 and queries3 need to be executed in Stencil, while queries2
	// need to be executed in the application database.
	// The sequence of executing the three queries ensure that
	// there is no anomaly.
	if err := db.TxnExecute(display.stencilDBConn, queries1); err != nil {
		return err
	} else {
		if err1 := db.TxnExecute(display.dstAppConfig.DBConn, queries2); err1 != nil {
			return err1
		} else {
			if err2 := db.TxnExecute(display.stencilDBConn, queries3); err2 != nil {
				return err2
			} else {
				return nil
			}
		}
	}
}

func (display *display) checkDisplayConditionsInNode(
	dataInNode []*HintStruct) ([]*HintStruct, []*HintStruct) {

	var displayedData, notDisplayedData []*HintStruct

	for _, oneDataInNode := range dataInNode {

		displayed := display.CheckDisplay(oneDataInNode)

		if !displayed {

			notDisplayedData = append(notDisplayedData, oneDataInNode)

		} else {

			displayedData = append(displayedData, oneDataInNode)

		}
	}

	return displayedData, notDisplayedData

}

func doesArgAttributeContainID(data string) bool {

	tmp := strings.Split(data, ".")

	if tmp[1] == "id" {
		return true
	} else {
		return false
	}
}

func getTableInArg(data string) string {

	tmp := strings.Split(data, ".")

	return tmp[0]
}

func (display *display) isNodeInCurrentMigration(oneMigratedData *HintStruct) bool {

	query := fmt.Sprintf(`select migration_id from display_flags 
		where app_id = %s and table_id = %s and id = %d and display_flag = true`,
		display.dstAppConfig.appID, oneMigratedData.TableID,
		oneMigratedData.KeyVal["id"])

	data, err := db.DataCall1(display.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if fmt.Sprint(data["migration_id"]) == strconv.Itoa(display.migrationID) {
		return true
	} else {
		return false
	}

}

func (display *display) getIDChanges(hint *HintStruct) string {

	// log.Println("ok")

	log.Println("Get ID changes:")

	query := fmt.Sprintf(
		`SELECT new_id FROM id_changes WHERE 
		app_id = %s and table_id = %s and old_id = %d`,
		display.dstAppConfig.appID,
		hint.TableID,
		hint.KeyVal["id"],
	)

	log.Println(query)

	data, err := db.DataCall1(display.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(data)

	if len(data) == 0 {
		return ""
	} else {
		return fmt.Sprint(data["new_id"])
	}

}

func (display *display) refreshCachedDataHints(hints []*HintStruct) {

	var err2, err3 error

	for i := range hints {

		// log.Println("=====")
		// log.Println(hints[i])
		// log.Println("=====")

		// hintID := strconv.Itoa(hints[i].KeyVal["id"])
		// hintDataID := fmt.Sprint(hints[i].Data["id"])

		// // Data id could change and the cached hint id could become stale and
		// // different from the got data id
		// // for example, in profile, user.id could change.
		// // There are two cases:
		// // 1. hint.id is old but data id is new
		// // so we use data id to update hint.id
		// // (this can only happen in the first phase since some attributes
		// // are not resolved because other data has not come)
		// // 2. hint.id is old and data id is old, then this does not cause problems
		// // display settings should be set to prevent this data from being displayed
		// // and this data should wait other data this data depends on to come
		// if hintID != hintDataID {

		// 	intHintDataID, err1 := strconv.Atoi(hintDataID)
		// 	if err1 != nil {
		// 		log.Fatal(err1)
		// 	}

		// 	hints[i].KeyVal["id"] = intHintDataID
		// }

		hints[i].Data, err2 = display.getOneRowBasedOnHint(hints[i])
		if err2 != nil {

			log.Println(err2)

			// log.Println(hints[i])

			newID := display.getIDChanges(hints[i])

			log.Println("new id:", newID)

			if newID == "" {
				// Note that for now this case is not considered
				panic("Since there is no application service, this data shoud not be deleted")

			} else {

				newHint := CreateHint(hints[i].Table, hints[i].TableID, newID)

				newHint.Data, err3 = display.getOneRowBasedOnHint(newHint)
				if err3 != nil {
					log.Fatal(err3)
				}

				hints[i] = newHint

				log.Println(hints[i])

			}

		}

	}

}

func getMigrationIDs(stencilDBConn *sql.DB, uid, srcAppID, dstAppID, migrationType string) []int {

	var mType string
	var migrationIDs []int

	switch migrationType {
	case "i":
		mType = "0"
	case "d":
		mType = "3"
	case "n":
		mType = "5"
	default:
		log.Fatal("Cannot find a corresponding migration type")
	}

	query := fmt.Sprintf(
		`select migration_id from migration_registration 
		where user_id = %s and src_app = %s and dst_app = %s 
		and migration_type = %s`,
		uid, srcAppID, dstAppID, mType)

	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	for _, data1 := range data {

		migrationID, ok := data1["migration_id"].(int64)
		if !ok {
			log.Fatal("Transform an interface type migrationID to an int64 fails")
		}

		migrationIDs = append(migrationIDs, int(migrationID))
	}

	return migrationIDs

}

func (display *display) logDisplayStartTime() {

	query := fmt.Sprintf(`
		INSERT INTO display_registration (start_time, migration_id)
		VALUES (now(), %d)`,
		display.migrationID,
	)

	err1 := db.TxnExecute1(display.stencilDBConn, query)
	if err1 != nil {
		log.Fatal(err1)
	}

}

func (display *display) logDisplayEndTime() {

	query := fmt.Sprintf(`
		UPDATE display_registration SET end_time = now()
		WHERE migration_id = %d`,
		display.migrationID,
	)

	err1 := db.TxnExecute1(display.stencilDBConn, query)
	if err1 != nil {
		log.Fatal(err1)
	}

}

func AddMarkAsDeleteToAllTables(dbConn *sql.DB) {

	query1 := `SELECT tablename FROM pg_catalog.pg_tables WHERE 
		schemaname != 'pg_catalog' AND schemaname != 'information_schema';`

	data := db.GetAllColsOfRows(dbConn, query1)

	// log.Println(data)

	for _, data1 := range data {

		query2 := fmt.Sprintf(`ALTER TABLE %s ADD mark_as_delete BOOLEAN DEFAULT FALSE;`,
			data1["tablename"])

		log.Println(query2)

		if _, err1 := dbConn.Exec(query2); err1 != nil {
			log.Fatal(err1)
		}

	}

}

func CreateIDChangesTable(dbConn *sql.DB) {

	query1 := `CREATE TABLE id_changes (
		app_id int8 NOT NULL,
		table_id int8 NOT NULL,
		old_id int8 NOT NULL,
		new_id int8 NOT NULL,
		migration_id int8 NOT NULL
	);`

	query2 := `CREATE INDEX ON id_changes (app_id, table_id, old_id);`

	queries := []string{
		query1, query2,
	}

	err1 := db.TxnExecute(dbConn, queries)
	if err1 != nil {
		log.Fatal(err1)
	}
}

func (display *display) getFirstArgsInREFByToTableToAttrInAllFromApps(
	toTable, toAttr string) map[string][]string {

	firstArgsFromApps := make(map[string][]string)

	for fromApp, mapping := range display.mappingsFromOtherAppsToDst {

		firstArgsFromApp := schema_mappings.GetFirstArgsInREFByToTableToAttr(mapping, toTable, toAttr)

		if len(firstArgsFromApp) != 0 {
			firstArgsFromApps[fromApp] = firstArgsFromApp
		}

	}

	return firstArgsFromApps

}

func (display *display) logUnresolvedRefAndData(tableName, tableID, id, unResolvedAttr string) {

	green := color.FgGreen.Render
	log.Println(green("Unresolved attribute is:"), green(tableName + ":" + unResolvedAttr))

	hint := CreateHint(tableName, tableID, id)
	data, err := display.getOneRowBasedOnHint(hint)

	if err != nil {
		log.Println(green("The data has been deleted by application services"))
	} else {
		log.Println(green("The data (node member) containing the unresolved attribute is:"))
		log.Println(green(transformMapToString(data)))
	}

}

func transformMapToString(data map[string]interface{}) string {

	procData := "| "

	for k, v := range data {
		procData += k + ": " + fmt.Sprint(v) + " | "
	}

	return procData

}

func (display *display) needToResolveReference(table, attr string) bool {
	
	// log.Println("Checking the need to resolve reference:", table, attr)

	for _, mapping := range display.mappingsFromOtherAppsToDst {

		if exists, err := schema_mappings.REFExists(mapping, table, attr); err != nil {
	
			// This can happen when there is no mapping
			// For example: 
			// When migrating from Diaspora to Mastodon:
			// there is no mapping to stream_entries.activity_id.
			// log.Println(err)
	
		} else {
			if exists {
				return true
			}
		}
	}

	if display.dstAppConfig.dag.IfDependsOnBasedOnDag(table, attr) {
		return true
	}

	return false
}
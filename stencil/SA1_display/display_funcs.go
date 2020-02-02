package SA1_display

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/reference_resolution"
	"stencil/schema_mappings"
	"strconv"
	"strings"
)

func CreateDisplayConfig(migrationID int,
	resolveReference, newDB, displayInFirstPhase bool) *displayConfig {

	var displayConfig displayConfig

	var srcAppConfig srcAppConfig

	var dstAppConfig dstAppConfig

	stencilDBConn := db.GetDBConn("stencil")

	srcAppID, dstAppID, srcUserID := 
		getSrcDstAppIDsUserIDByMigrationID(stencilDBConn, migrationID)

	srcAppName := getAppNameByAppID(stencilDBConn, srcAppID)
	dstAppName := getAppNameByAppID(stencilDBConn, dstAppID)

	allMappings, err1 := config.LoadSchemaMappings()
	if err1 != nil {
		log.Fatal(err1)
	}

	mappingsFromSrcToDst, err2 := 
		schema_mappings.GetToAppMappings(allMappings, srcAppName, dstAppName)
	
	if err2 != nil {
		log.Fatal(err2)
	}

	dstDAG, err4 := loadDAG(dstAppName)
	if err4 != nil {
		log.Fatal(err4)
	}

	var dstDBConn *sql.DB

	if newDB {
		dstDBConn = db.GetDBConn(dstAppName)
	} else {
		dstDBConn = db.GetDBConn(dstAppName, true)
	}

	// Note that when the display thread is initializing, dstUserID could be nil
	// because the data has not been migrated yet
	dstRootMember, dstRootAttr, dstUserID := getDstRootMemberAttrID(
		stencilDBConn, dstAppID, migrationID, dstDAG)

	dstAppTableIDNamePairs := make(map[string]string)
	dstAppTableNameIDPairs := make(map[string]string)

	dstRes := getTableIDNamePairsInApp(stencilDBConn, dstAppID)

	for _, dstRes1 := range dstRes {

		dstAppTableIDNamePairs[fmt.Sprint(dstRes1["pk"])] = 
			fmt.Sprint(dstRes1["table_name"])

		dstAppTableNameIDPairs[fmt.Sprint(dstRes1["table_name"])] = 
			fmt.Sprint(dstRes1["pk"])
	}

	srcAppTableNameIDPairs := make(map[string]string)

	srcRes := getTableIDNamePairsInApp(stencilDBConn, srcAppID)

	for _, srcRes1 := range srcRes {

		srcAppTableNameIDPairs[fmt.Sprint(srcRes1["table_name"])] = 
			fmt.Sprint(srcRes1["pk"])

	}

	appIDNamePairs := GetAppIDNamePairs(stencilDBConn)
	tableIDNamePairs := GetTableIDNamePairs(stencilDBConn)

	refResolutionConfig := reference_resolution.InitializeReferenceResolution(
		migrationID, dstAppID, dstAppName, dstDBConn, stencilDBConn,
		dstAppTableNameIDPairs, appIDNamePairs, tableIDNamePairs,
		allMappings, mappingsFromSrcToDst)

	srcAppConfig.appID = srcAppID
	srcAppConfig.appName = srcAppName
	srcAppConfig.userID = srcUserID
	srcAppConfig.tableNameIDPairs = srcAppTableNameIDPairs

	dstAppConfig.appID = dstAppID
	dstAppConfig.appName = dstAppName
	dstAppConfig.tableNameIDPairs = dstAppTableNameIDPairs
	dstAppConfig.rootTable = dstRootMember
	dstAppConfig.rootAttr = dstRootAttr
	dstAppConfig.userID = dstUserID
	dstAppConfig.dag = dstDAG
	dstAppConfig.DBConn = dstDBConn
	dstAppConfig.ownershipDisplaySettingsSatisfied = false

	displayConfig.stencilDBConn = stencilDBConn
	displayConfig.appIDNamePairs = appIDNamePairs
	displayConfig.tableIDNamePairs = tableIDNamePairs
	displayConfig.attrIDNamePairs = GetAttrIDNamePairs(stencilDBConn)
	displayConfig.migrationID = migrationID
	displayConfig.resolveReference = resolveReference
	displayConfig.srcAppConfig = &srcAppConfig
	displayConfig.dstAppConfig = &dstAppConfig
	displayConfig.refResolutionConfig = refResolutionConfig
	displayConfig.mappingsFromSrcToDst = mappingsFromSrcToDst
	displayConfig.displayInFirstPhase = displayInFirstPhase

	return &displayConfig

}

func closeDBConn(conn *sql.DB) {

	err := conn.Close()
	if err != nil {
		log.Fatal(err)
	}

}

func closeDBConns(displayConfig *displayConfig) {

	log.Println("Close db connections in the display thread")

	closeDBConn(displayConfig.stencilDBConn)
	closeDBConn(displayConfig.dstAppConfig.DBConn)

}

func getDstRootMemberAttrID(stencilDBConn *sql.DB,
	appID string, migrationID int, dstDAG *DAG) (string, string, string) {

	// log.Println(*dstDAG)

	dstRootMember, dstRootAttr, err2 := getRootMemberAttr(dstDAG)
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

func oldCreateDisplayConfig(migrationID int, resolveReference, newDB bool) *config.DisplayConfig {

	var displayConfig config.DisplayConfig

	stencilDBConn := db.GetDBConn(config.StencilDBName)

	srcAppID, dstAppID, userID := getSrcDstAppIDsUserIDByMigrationID(stencilDBConn, migrationID)

	dstAppName := getAppNameByAppID(stencilDBConn, dstAppID)
	srcAppName := getAppNameByAppID(stencilDBConn, srcAppID)

	// app_id := db.GetAppIDByAppName(stencilDBConn, app)

	appConfig, err := config.CreateAppConfigDisplay(dstAppName, dstAppID, stencilDBConn, newDB)
	if err != nil {
		log.Fatal(err)
	}

	allMappings, err1 := config.LoadSchemaMappings()
	if err1 != nil {
		log.Fatal(err1)
	}

	mappingsToDst, err2 := schema_mappings.GetToAppMappings(allMappings, srcAppName, dstAppName)
	if err2 != nil {
		log.Fatal(err2)
	}

	displayConfig.ResolveReference = resolveReference
	displayConfig.AllMappings = allMappings
	displayConfig.MappingsToDst = mappingsToDst
	displayConfig.SrcAppID = srcAppID
	displayConfig.SrcAppName = srcAppName
	displayConfig.AppConfig = &appConfig
	displayConfig.AttrIDNamePairs = GetAttrIDNamePairs(stencilDBConn)
	displayConfig.AppIDNamePairs = GetAppIDNamePairs(stencilDBConn)
	displayConfig.TableIDNamePairs = GetTableIDNamePairs(stencilDBConn)
	displayConfig.StencilDBConn = stencilDBConn
	displayConfig.MigrationID = migrationID
	displayConfig.UserID = userID

	return &displayConfig

}

func oldInitialize(app string) (*sql.DB, *config.AppConfig) {

	stencilDBConn := db.GetDBConn(config.StencilDBName)

	app_id := db.GetAppIDByAppName(stencilDBConn, app)

	appConfig, err := config.CreateAppConfigDisplay(app, app_id, stencilDBConn, false)
	if err != nil {
		log.Fatal(err)
	}

	return stencilDBConn, &appConfig

}

func GetUndisplayedMigratedData(displayConfig *displayConfig) []*HintStruct {

	var displayHints []*HintStruct

	query := fmt.Sprintf(`SELECT table_id, id FROM display_flags
		 WHERE app_id = %s and migration_id = %d and display_flag = true`,
		displayConfig.dstAppConfig.appID, displayConfig.migrationID)

	data := db.GetAllColsOfRows(displayConfig.stencilDBConn, query)
	// fmt.Println(data)

	for _, data1 := range data {

		displayHints = append(displayHints, TransformDisplayFlagDataToHint(displayConfig, data1))

	}
	// fmt.Println(displayHints)

	return displayHints

}

func CheckMigrationComplete(displayConfig *displayConfig) bool {

	query := fmt.Sprintf(
		`SELECT 1 FROM txn_logs 
		WHERE action_id = %d and action_type='COMMIT' LIMIT 1`,
		displayConfig.migrationID)

	data := db.GetAllColsOfRows(displayConfig.stencilDBConn, query)

	if len(data) == 0 {

		return false

	} else {

		return true

	}

}

func CheckMigrationComplete1(stencilDBConn *sql.DB, migrationID int) bool {

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
func Display(displayConfig *displayConfig, dataHints []*HintStruct) error {

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
		if !isNodeInCurrentMigration(displayConfig, dataHint) {

			log.Println(`This data is not migrated by the currently checked user,
				but this display thread will also display this data`)

		}

		log.Println("Check attributes to be set as NULLs before displaying the data")

		ID := dataHint.TransformHintToIdenity(displayConfig)

		// Even though references have been updated in checking dependencies and ownership,
		// update reference the last time before displaying the data
		reference_resolution.ResolveReference(displayConfig.refResolutionConfig, ID)

		updatedAttrs := reference_resolution.GetUpdatedAttributes(
			displayConfig.refResolutionConfig,
			ID,
		)

		attrsToBeUpdated := schema_mappings.GetAllMappedAttributesContainingREFInMappings(
			displayConfig.mappingsFromSrcToDst,
			dataHint.Table)

		var attrsToBeSetToNULLs []string

		for attr := range attrsToBeUpdated {

			if _, ok := updatedAttrs[attr]; !ok {
				attrsToBeSetToNULLs = append(attrsToBeSetToNULLs, attr)
			}

		}

		// query1 = fmt.Sprintf("UPDATE %s SET display_flag = false WHERE id = %d;",
		// 	dataHint.Table, dataHint.KeyVal["id"])

		query1 = fmt.Sprintf("UPDATE %s SET display_flag = false", dataHint.Table)

		for _, attr := range attrsToBeSetToNULLs {

			query1 += ", " + attr + " = NULL"

		}

		query1Where := fmt.Sprintf(" WHERE id = %d;", dataHint.KeyVal["id"])

		query1 += query1Where

		query2 = fmt.Sprintf(`UPDATE Display_flags SET 
			display_flag = false, updated_at = now() 
			WHERE app_id = %s and table_id = %s and id = %d;`,
			displayConfig.dstAppConfig.appID, dataHint.TableID, dataHint.KeyVal["id"])
		
		query3 = fmt.Sprintf(`UPDATE evaluation SET
			displayed_at = now() WHERE migration_id = '%d' and
			src_app = '%s' and dst_app = '%s' and
			dst_table = '%s' and dst_id = '%d'`,
			displayConfig.migrationID, displayConfig.srcAppConfig.appID,
			displayConfig.dstAppConfig.appID, dataHint.TableID, dataHint.KeyVal["id"])

		log.Println("**************************************")
		log.Println(query1)
		log.Println(query2)
		log.Println(query3)
		log.Println("**************************************")

		queries1 = append(queries1, query1)

		queries2 = append(queries2, query2, query3)
		
	}

	if err1 := db.TxnExecute(displayConfig.dstAppConfig.DBConn, queries1); err1 != nil {

		return err1

	} else {

		if err2 := db.TxnExecute(displayConfig.stencilDBConn, queries2); err2 != nil {

			return err2

		} else {

			return nil

		}
	}

}

func CheckDisplay(displayConfig *displayConfig, dataHint *HintStruct) bool {

	query := fmt.Sprintf("SELECT display_flag from %s where id = %d",
		dataHint.Table, dataHint.KeyVal["id"])

	data1, err := db.DataCall1(displayConfig.dstAppConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data1)

	return !data1["display_flag"].(bool)

}

func getAppNameByAppID(stencilDBConn *sql.DB, appID string) string {

	query := fmt.Sprintf("select app_name from apps where pk = %s", appID)

	// log.Println(query)

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprint(data["app_name"])

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

func GetAppIDNamePairs(stencilDBConn *sql.DB) map[string]string {

	appIDNamePairs := make(map[string]string)

	query := fmt.Sprintf("select app_name, pk from apps")

	// log.Println(query)

	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	for _, data1 := range data {
		appIDNamePairs[fmt.Sprint(data1["pk"])] = fmt.Sprint(data1["app_name"])
	}

	return appIDNamePairs

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

func getTableIDNamePairsInApp(stencilDBConn *sql.DB, app_id string) []map[string]interface{} {

	query := fmt.Sprintf("select pk, table_name from app_tables where app_id = %s", app_id)

	result, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func GetTableIDNamePairs(stencilDBConn *sql.DB) map[string]string {

	tableIDNamePairs := make(map[string]string)

	query := fmt.Sprintf("select pk, table_name from app_tables;")

	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	for _, data1 := range data {
		tableIDNamePairs[fmt.Sprint(data1["pk"])] = fmt.Sprint(data1["table_name"])
	}

	return tableIDNamePairs

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

func getSrcDstAppIDsUserIDByMigrationID(stencilDBConn *sql.DB,
	migrationID int) (string, string, string) {

	query := fmt.Sprintf("select src_app, dst_app, user_id from migration_registration")

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprint(data["src_app"]), fmt.Sprint(data["dst_app"]), fmt.Sprint(data["user_id"])

}

func ConvertMapToJSONString(data map[string]interface{}) string {

	convertedData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return string(convertedData)

}

func chechPutIntoDataBag(displayConfig *displayConfig,
	secondRound bool, dataHints []*HintStruct) error {

	if secondRound {

		err9 := putIntoDataBag(displayConfig, dataHints)
		if err9 != nil {
			log.Fatal(err9)
		}

		return NoNodeCanBeDisplayed

	} else {

		return NoNodeCanBeDisplayed
	}
}

func setDstUserIDIfNotSet(displayConfig *displayConfig) {

	if displayConfig.dstAppConfig.userID != "" {
		return
	}

	// Since the in current settings, there is only one row and the root attribute is always id,
	// we only do in the following way. Note that this is not a generic way.
	query := fmt.Sprintf(
		`SELECT id FROM display_flags WHERE app_id = %s 
		and table_id = %s and migration_id = %d`,
		displayConfig.dstAppConfig.appID,
		displayConfig.dstAppConfig.tableNameIDPairs[displayConfig.dstAppConfig.rootTable],
		displayConfig.migrationID)

	log.Println(query)

	data, err := db.DataCall1(displayConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if data["id"] == nil {
		log.Fatal("This is the second phase, so the user id should not be nil!")
	} else {
		displayConfig.dstAppConfig.userID = fmt.Sprint(data["id"])
	}

}

// When putting data to dag bags, it does not matter whether we set unresolved references
// to NULLs or not, so we don't set those as NULLs.
func putIntoDataBag(displayConfig *displayConfig, dataHints []*HintStruct) error {

	setDstUserIDIfNotSet(displayConfig)

	var queries1, queries2, queries3 []string

	var q1, q2, q3 string

	for _, dataHint := range dataHints {

		// we try to avoid putting data into data bags if it is not in the
		// currently checked user's migration
		if !isNodeInCurrentMigration(displayConfig, dataHint) {

			log.Println(`This data is not migrated by the currently checked user,
				so it will not be put into data bags by this display thread`)

			continue
		}

		// dataHint.Data could be nil, which means there is no data,
		// if a thread crashes before executing queries3
		// and after executing queries1 and queries2, or data is deleted by services.
		// In both cases, there is no need to execute queries1 and queries2 again.
		if dataHint.Data != nil {

			q1 = fmt.Sprintf(`INSERT INTO data_bags 
				(app, member, id, data, user_id, migration_id) VALUES 
				(%s, %s, %d, '%s', %s, %d)`,
				displayConfig.dstAppConfig.appID,
				dataHint.TableID,
				dataHint.KeyVal["id"],
				ConvertMapToJSONString(dataHint.Data),
				displayConfig.dstAppConfig.userID,
				displayConfig.migrationID)

			q2 = fmt.Sprintf("DELETE FROM %s WHERE id = %d",
				dataHint.Table, dataHint.KeyVal["id"])

		}

		q3 = fmt.Sprintf(`UPDATE display_flags SET 
			display_flag = false, updated_at = now() 
			WHERE app_id = %s and table_id = %s and id = %d;`,
			displayConfig.dstAppConfig.appID, dataHint.TableID, dataHint.KeyVal["id"])

		log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
		log.Println("INSERT INTO data_bags:", q1)
		log.Println("DELETE FROM the application:", q2)
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
	if err := db.TxnExecute(displayConfig.stencilDBConn, queries1); err != nil {

		return err

	} else {

		if err1 := db.TxnExecute(displayConfig.dstAppConfig.DBConn, queries2); err1 != nil {

			return err1

		} else {

			if err2 := db.TxnExecute(displayConfig.stencilDBConn, queries3); err2 != nil {

				return err2

			} else {

				return nil
			}
		}
	}
}

func checkDisplayConditionsInNode(displayConfig *displayConfig,
	dataInNode []*HintStruct) ([]*HintStruct, []*HintStruct) {

	var displayedData, notDisplayedData []*HintStruct

	for _, oneDataInNode := range dataInNode {

		displayed := CheckDisplay(displayConfig, oneDataInNode)

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

func isNodeInCurrentMigration(displayConfig *displayConfig,
	oneMigratedData *HintStruct) bool {

	query := fmt.Sprintf(`select migration_id from display_flags 
		where app_id = %s and table_id = %s and id = %d and display_flag = true`,
		displayConfig.dstAppConfig.appID, oneMigratedData.TableID,
		oneMigratedData.KeyVal["id"])

	data, err := db.DataCall1(displayConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if fmt.Sprint(data["migration_id"]) == strconv.Itoa(displayConfig.migrationID) {
		return true
	} else {
		return false
	}

}

func refreshCachedDataHints(displayConfig *displayConfig,
	hints []*HintStruct) {

	var err2 error

	for i, hint1 := range hints {

		hints[i].Data, err2 = getOneRowBasedOnHint(displayConfig, hint1)
		if err2 != nil {
			log.Fatal(err2)
		}

	}

}

func getMigrationIDs(stencilDBConn *sql.DB,
	uid, srcAppID, dstAppID, migrationType string) []int {

	var mType string
	var migrationIDs []int

	switch migrationType {
	case "d":
		mType = "3"
	case "n":
		mType = "5"
	default:
		log.Fatal("Cannot find a corresponding migration type")
	}

	query := fmt.Sprintf(`select migration_id from migration_registration 
		where user_id = %s and src_app = %s and dst_app = %s and migration_type = %s`,
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

func logDisplayStartTime(displayConfig *displayConfig) {

	query := fmt.Sprintf(`
		INSERT INTO display_registration (start_time, migration_id)
		VALUES (now(), %d)`,
		displayConfig.migrationID,
	)

	err1 := db.TxnExecute1(displayConfig.stencilDBConn, query); 
	if err1 != nil {
		log.Fatal(err1)
	}

}

func logDisplayEndTime(displayConfig *displayConfig) {

	query := fmt.Sprintf(`
		UPDATE display_registration SET end_time = now()
		WHERE migration_id = %d`,
		displayConfig.migrationID,
	)

	err1 := db.TxnExecute1(displayConfig.stencilDBConn, query); 
	if err1 != nil {
		log.Fatal(err1)
	}

}
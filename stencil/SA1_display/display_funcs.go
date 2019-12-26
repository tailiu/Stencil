package SA1_display

import (
	"database/sql"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/schema_mappings"
	"stencil/reference_resolution"
	"encoding/json"
)

func CreateDisplayConfig(migrationID int, resolveReference, newDB bool) *displayConfig {

	var displayConfig displayConfig

	var srcAppConfig srcAppConfig

	var dstAppConfig dstAppConfig

	stencilDBConn := db.GetDBConn(config.StencilDBName)

	srcAppID, dstAppID, srcUserID := getSrcDstAppIDsUserIDByMigrationID(stencilDBConn, migrationID)

	srcAppName := getAppNameByAppID(stencilDBConn, srcAppID)
	dstAppName := getAppNameByAppID(stencilDBConn, dstAppID)

	allMappings, err1 := config.LoadSchemaMappings()
	if err1 != nil {
		log.Fatal(err1)
	}

	mappingsToDst, err2 := schema_mappings.GetToAppMappings(allMappings, srcAppName, dstAppName)
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
		dstDBConn = db.GetDBConn2(dstAppName)
	}

	dstAppTableIDNamePair := make(map[string]string)

	res := getTableIDNamePairsInApp(stencilDBConn, dstAppID)

	for _, tableIDNamePair := range res {

		dstAppTableIDNamePair[fmt.Sprint(res["pk"])] = fmt.Sprint(res["table_name"])

	}

	srcAppConfig.appID = srcAppID
	srcAppConfig.appName = srcAppName
	srcAppConfig.userID = srcUserID

	dstAppConfig.appID = dstAppID
	dstAppConfig.appName = dstAppName
	dstAppConfig.tableNameIDPairs = dstAppTableIDNamePair
	dstAppConfig.userID = getDstUserID(stencilDBConn, dstAppID, migrationID, dstDAG)
	dstAppConfig.dag = dstDAG
	dstAppConfig.DBConn = dstDBConn

	displayConfig.stencilDBConn = stencilDBConn
	displayConfig.appIDNamePairs = GetAppIDNamePairs(stencilDBConn)
	displayConfig.tableIDNamePairs = GetTableIDNamePairs(stencilDBConn)
	displayConfig.attrIDNamePairs = GetAttrIDNamePairs(stencilDBConn)
	displayConfig.migrationID = migrationID
	displayConfig.allMappings = allMappings
	displayConfig.mappingsToDst = mappingsToDst
	displayConfig.resolveReference = resolveReference
	displayConfig.srcAppConfig = &srcAppConfig
	displayConfig.dstAppConfig = &dstAppConfig

	return &displayConfig

}

func getDstUserID(stencilDBConn *sql.DB, appID string, migrationID int, dstDAG *DAG) string {

	// log.Println(*dstDAG)

	dstRootMember, _, err2 := getRootMemberAttr(dstDAG)
	if err2 != nil {
		log.Fatal(err2)
	}

	tableID := getTableIDByTableName(stencilDBConn, appID, dstRootMember)

	// Since the in current settings, there is only one row and the root attribute is always id,
	// we only do in the following way. Note that this is not a generic way.
	query := fmt.Sprintf(`SELECT id FROM display_flags WHERE app_id = %s 
		and table_id = %s and migration_id = %d`, appID, tableID, migrationID)

	// log.Println(query)

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprint(data["id"])

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

func GetUndisplayedMigratedData(displayConfig *config.DisplayConfig) []*HintStruct {

	var displayHints []*HintStruct

	query := fmt.Sprintf(`SELECT table_id, id FROM display_flags
		 WHERE app_id = %s and migration_id = %d and display_flag = true`, 
		displayConfig.AppConfig.AppID, displayConfig.MigrationID)
	
	data := db.GetAllColsOfRows(displayConfig.StencilDBConn, query)
	// fmt.Println(data)

	for _, data1 := range data {

		displayHints = append(displayHints, TransformDisplayFlagDataToHint(displayConfig, data1))

	}
	// fmt.Println(displayHints)
	
	return displayHints

}

func CheckMigrationComplete(displayConfig *config.DisplayConfig) bool {
	
	query := fmt.Sprintf("SELECT 1 FROM txn_logs WHERE action_id = %d and action_type='COMMIT' LIMIT 1", 
		displayConfig.MigrationID)
	
	data := db.GetAllColsOfRows(displayConfig.StencilDBConn, query)
	
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
func Display(displayConfig *config.DisplayConfig, dataHints []*HintStruct) error {

	var queries1, queries2 []string

	var query1, query2 string

	for _, dataHint := range dataHints {
		
		ID := dataHint.TransformHintToIdenity(displayConfig)

		myUpdatedAttrs, _ := reference_resolution.ResolveReference(displayConfig, ID)

		attrsToBeUpdated := schema_mappings.GetAllMappedAttributesContainingREFInMappings(
			displayConfig.MappingsToDst,
			dataHint.Table)

		var attrsToBeSetToNULLs []string

		for attr := range attrsToBeUpdated {
			
			if _, ok := myUpdatedAttrs[attr]; !ok {
				attrsToBeSetToNULLs = append(attrsToBeSetToNULLs, attr)
			}

		}

		// query1 = fmt.Sprintf("UPDATE %s SET display_flag = false WHERE id = %d;",
		// 	dataHint.Table, dataHint.KeyVal["id"])
		
		query1 = fmt.Sprintf("UPDATE %s SET display_flag = false", dataHint.Table)

		for _, attr := range attrsToBeSetToNULLs {

			query1 += ", "  + attr + " = NULL"

		}

		query1Where := fmt.Sprintf(" WHERE id = %d;", dataHint.KeyVal["id"])

		query1 += query1Where


		query2 = fmt.Sprintf(`UPDATE Display_flags SET 
			display_flag = false, updated_at = now() 
			WHERE app_id = %s and table_id = %s and id = %d;`,
			displayConfig.AppConfig.AppID, dataHint.TableID, dataHint.KeyVal["id"])
		
		log.Println("**************************************")
		log.Println(query1)
		log.Println(query2)
		log.Println("**************************************")
		
		queries1 = append(queries1, query1)

		queries2 = append(queries2, query2)
		
	}

	if err1 := db.TxnExecute(displayConfig.AppConfig.DBConn, queries1); err1 != nil {

		return err1

	} else {

		if err2 := db.TxnExecute(displayConfig.StencilDBConn, queries2); err2 != nil {
			
			return err2
		
		} else {

			return nil
		}
	}

}

func CheckDisplay(displayConfig *config.DisplayConfig, dataHint *HintStruct) bool {

	query := fmt.Sprintf("SELECT display_flag from %s where id = %d",
		dataHint.Table, dataHint.KeyVal["id"])
	
	data1, err := db.DataCall1(displayConfig.AppConfig.DBConn, query)
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


// When putting data to dag bags, it does not matter whether we set unresolved references
// to NULLs or not, so we don't set those as NULLs. 
func PutIntoDataBag(displayConfig *config.DisplayConfig, dataHints []*HintStruct) error {
	
	var queries1, queries2, queries3 []string

	var q1, q2, q3 string

	for _, dataHint := range dataHints {
		
		// dataHint.Data could be nil, which means there is no data,
		// if a thread crashes before executing queries3
		// and after executing queries1 and queries2, or data is deleted by services.
		// In both cases, there is no need to execute queries1 and queries2 again.
		if dataHint.Data != nil {

			q1 = fmt.Sprintf(`INSERT INTO data_bags 
				(app, member, id, data, user_id, migration_id) VALUES 
				(%s, %s, %d, '%s', %s, %d)`, 
				displayConfig.AppConfig.AppID,
				dataHint.TableID,
				dataHint.KeyVal["id"],
				ConvertMapToJSONString(dataHint.Data),
				displayConfig.UserID,
				displayConfig.MigrationID)

			q2 = fmt.Sprintf("DELETE FROM %s WHERE id = %d", 
				dataHint.Table, dataHint.KeyVal["id"])
			
		}
		
		q3 = fmt.Sprintf(`UPDATE display_flags SET 
			display_flag = false, updated_at = now() 
			WHERE app_id = %s and table_id = %s and id = %d;`,
			displayConfig.AppConfig.AppID, dataHint.TableID, dataHint.KeyVal["id"])

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
	if err := db.TxnExecute(displayConfig.StencilDBConn, queries1); err != nil {

		return err

	} else {

		if err1 := db.TxnExecute(displayConfig.AppConfig.DBConn, queries2); err1 != nil {

			return err1

		} else {

			if err2 := db.TxnExecute(displayConfig.StencilDBConn, queries3); err2 != nil {
				
				return err2

			} else {

				return nil				
			}
		}
	}
}

func checkDisplayConditionsInNode(displayConfig *config.DisplayConfig, 
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

func isNodeMigratingUserRootNode(displayConfig *config.DisplayConfig, 
	dataInNode []*HintStruct) (bool, error) {

	
	return true, nil
}
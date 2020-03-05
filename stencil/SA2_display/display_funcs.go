package SA2_display

import (
	"database/sql"
	"fmt"
	"log"
	"stencil/common_funcs"
	"stencil/schema_mappings"
	"stencil/config"
	"stencil/db"
	"stencil/qr"
	"time"
	"math/rand"
	"math"
)

func CreateDisplayConfig(migrationID int, 
	displayInFirstPhase bool) *displayConfig {

	var displayConfig displayConfig

	var srcAppConfig srcAppConfig

	var dstAppConfig dstAppConfig

	stencilDBConn := db.GetDBConn("stencil")

	srcAppID, dstAppID, userID := 
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

	dstDAG, err4 := common_funcs.LoadDAG(dstAppName)
	if err4 != nil {
		log.Fatal(err4)
	}

	dstDBConn := db.GetDBConn(dstAppName)

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

		srcAppTableNameIDPairs[fmt.Sprint(srcRes1["table_name"])] = 
			fmt.Sprint(srcRes1["pk"])

	}

	appIDNamePairs := common_funcs.GetAppIDNamePairs(stencilDBConn)
	tableIDNamePairs := common_funcs.GetTableIDNamePairs(stencilDBConn)

	srcAppConfig.appID = srcAppID
	srcAppConfig.appName = srcAppName
	srcAppConfig.tableNameIDPairs = srcAppTableNameIDPairs

	dstAppConfig.appID = dstAppID
	dstAppConfig.appName = dstAppName
	dstAppConfig.tableNameIDPairs = dstAppTableNameIDPairs
	dstAppConfig.dag = dstDAG
	dstAppConfig.DBConn = dstDBConn
	dstAppConfig.ownershipDisplaySettingsSatisfied = false
	dstAppConfig.qr = qr.NewQR(dstAppName, dstAppID)

	displayConfig.stencilDBConn = stencilDBConn
	displayConfig.appIDNamePairs = appIDNamePairs
	displayConfig.tableIDNamePairs = tableIDNamePairs
	displayConfig.migrationID = migrationID
	displayConfig.srcAppConfig = &srcAppConfig
	displayConfig.dstAppConfig = &dstAppConfig
	displayConfig.mappingsFromSrcToDst = mappingsFromSrcToDst
	displayConfig.displayInFirstPhase = displayInFirstPhase
	displayConfig.userID = userID

	return &displayConfig

}

func closeDBConn(conn *sql.DB) {

	err := conn.Close()
	if err != nil {
		log.Fatal(err)
	}

}

func closeDBConns(displayConfig *displayConfig) {

	log.Println("Close db connections in the SA2 display thread")

	closeDBConn(displayConfig.stencilDBConn)
	closeDBConn(displayConfig.dstAppConfig.DBConn)

}

func RandomNonnegativeInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(math.MaxInt32)
}

func getUserIDByMigrationID(stencilDBConn *sql.DB, 
	migrationID int) string {

	query := fmt.Sprintf(
		`select user_id from migration_registration
		where migration_id = %d`,
		migrationID,
	)

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprint(data["user_id"])
}

func getDstAppIDUserIDByMigrationID(stencilDBConn *sql.DB,
	migrationID int) (string, string) {

	query := fmt.Sprintf(
		`SELECT dst_app, user_id FROM migration_registration 
		WHERE migration_id = %d`,
		migrationID,
	)

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	dstApp := fmt.Sprint(data["dst_app"])
	userID := fmt.Sprint(data["user_id"])

	return dstApp, userID

}

func GetUndisplayedMigratedData(displayConfig *displayConfig) []*HintStruct {
	
	var displayHints []*HintStruct
	
	// This is important that table id / group id should also be used to get results in the new design
	// For example, in the one-to-multiple mapping, the same row id has different group ids / table ids
	// Those rows could be displayed differently
	query := fmt.Sprintf(
		`SELECT table_id, array_agg(row_id) as row_ids 
		FROM migration_table where mflag = 1 and app_id = %s and migration_id = %d 
		group by group_id, table_id;`,
		displayConfig.dstAppConfig.appID,
		displayConfig.migrationID,
	)
	
	data := db.GetAllColsOfRows(displayConfig.stencilDBConn, query)
	// log.Println(data)

	// If we don't use physical schema, both table_name and id are necessary to identify a piece of migratd data.
	// Actually, in our physical schema, row_id itself is enough to identify a piece of migrated data.
	// We use table_name to optimize performance
	for _, data1 := range data {
		displayHints = append(displayHints, TransformRowToHint(displayConfig, data1))
	}

	// log.Println(displayHints)
	
	return displayHints
}

func CheckMigrationComplete(displayConfig *displayConfig) bool {

	query := fmt.Sprintf(
		`SELECT 1 FROM txn_logs WHERE action_id = %d 
		and action_type='COMMIT' LIMIT 1`,
		displayConfig.migrationID,
	)
	
	data := db.GetAllColsOfRows(displayConfig.stencilDBConn, query)
	
	if len(data) == 0 {
		
		return false
	
	} else {
		
		return true
	}
}

func logDisplayStartTime(displayConfig *displayConfig) {

	query := fmt.Sprintf(`
		INSERT INTO display_registration (start_time, migration_id)
		VALUES (now(), %d)`,
		displayConfig.migrationID,
	)

	err1 := db.TxnExecute1(displayConfig.stencilDBConn, query)
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

	err1 := db.TxnExecute1(displayConfig.stencilDBConn, query)
	if err1 != nil {
		log.Fatal(err1)
	}

}

func CheckDisplay(displayConfig *displayConfig, data *HintStruct) bool {

	// Here for one group, we only need to check 
	// one row_id to see whether the group is displayed or not
	// It should be noted that table_id / group_id should also be considered
	query := fmt.Sprintf(
		`SELECT mflag FROM migration_table 
		WHERE row_id = %d and app_id = %s and table_id = %s`, 
		data.RowIDs[0], displayConfig.dstAppConfig.appID, data.TableID,
	)
	
	// log.Println("==========")
	// log.Println(query)
	// log.Println(data)
	// log.Println("==========")
	
	data1, err := db.DataCall1(displayConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(data1)

	if fmt.Sprint(data1["mflag"]) == "1" {
		return false
	} else {
		return true
	}

}

func Display(displayConfig *displayConfig, dataInNode []*HintStruct) error {
	
	var queries []string

	appID := displayConfig.dstAppConfig.appID

	for _, dataHint := range dataInNode {

		// This is an optimization to prevent possible path conflict
		// We only need to test one rowID in a data hint
		if CheckDisplay(displayConfig, dataHint) {

			return PathConflictsWhenDisplayingData

		}

		data := dataHint.GetAllRowIDs(displayConfig)

		for _, data1 := range data {
			
			rowID := fmt.Sprint(data1["row_id"])

			// It should be noted that table_id / group_id should also be considered
			query := fmt.Sprintf(
				`UPDATE migration_table SET mflag = 0, updated_at = now() 
				WHERE row_id = %s and app_id = %s and table_id = %s`, 
				rowID, appID, dataHint.TableID,
			)
			
			query1 := fmt.Sprintf(
				`UPDATE evaluation SET displayed_at = now()
				WHERE migration_id = '%d' and dst_app = '%s' 
				and dst_table = '%s' and dst_id = '%s'`,
				displayConfig.migrationID, appID,
				dataHint.TableID, rowID,
			)

			queries = append(queries, query, query1)

		}
	}
	
	// if deletionHoldEnable {
		
	// 	var dhQueries []string
		
	// 	dhQueries, dhStack = AddToDeletionHoldStack(dhStack, dataHints, threadID)
	// 	queries = append(queries, dhQueries...)

	// }

	log.Println("**************************************")
	log.Println("Display Data:")
	for seq, q1 := range queries {
		log.Println("Query", seq + 1)
		log.Println(q1)
	}
	log.Println("**************************************")

	return db.TxnExecute(displayConfig.stencilDBConn, queries)

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

func chechPutIntoDataBag(displayConfig *displayConfig,
	secondRound bool, dataHints []*HintStruct) error {

	if secondRound {

		err9 := putIntoDataBag(displayConfig, dataHints)
		if err9 != nil {
			log.Println(err9)
		}

		return common_funcs.NoNodeCanBeDisplayed

	} else {

		return common_funcs.NoNodeCanBeDisplayed
	}
}

func alreadyInBag(displayConfig *displayConfig, data *HintStruct) bool {
	
	// Here for one group, we only need to check one to see whether the group is displayed or not
	// It should be noted that table_id / group_id should also be considered
	query := fmt.Sprintf(
		`SELECT bag FROM migration_table 
		WHERE row_id = %d and app_id = %s and table_id = %s`, 
		data.RowIDs[0], 
		displayConfig.dstAppConfig.appID,
		data.TableID,
	)

	// log.Println(query)
	
	data1, err := db.DataCall1(displayConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	
	// log.Println(data1)
	
	return data1["bag"].(bool)

}

func putIntoDataBag(displayConfig *displayConfig, dataHints []*HintStruct) error {
	
	var queries []string

	for _, dataHint := range dataHints {
		
		// Similar to displaying data, this is an optimization to prevent possible path conflict
		// We only need to test one rowID in a data hint
		if alreadyInBag(displayConfig, dataHint) {

			// log.Println("Found that there is a path conflict!! When putting data in a databag")
			
			// return errors.New("Path conflict")

			return PathConflictsWhenPuttingInBags
		
		}

		rowIDs := dataHint.GetAllRowIDs(displayConfig)
		
		for _, rowID := range rowIDs {

			// It should be noted that table_id / group_id should also be considered
			query := fmt.Sprintf(
				`UPDATE migration_table SET 
				user_id = %s, bag = true, mark_as_delete = true, 
				mflag = 0, updated_at = now() 
				WHERE row_id = %s and app_id = %s and table_id = %s`,
				displayConfig.userID, fmt.Sprint(rowID["row_id"]), 
				displayConfig.dstAppConfig.appID, 
				dataHint.TableID,
			)
			
			queries = append(queries, query)
		}
	}

	log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	log.Println("Put Data into Data Bags:")
	for seq, q1 := range queries {
		log.Println("Query", seq + 1)
		log.Println(q1)
	}
	log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")

	return db.TxnExecute(displayConfig.stencilDBConn, queries)
	
}

func getMigrationIDs(stencilDBConn *sql.DB,
	uid, srcAppID, dstAppID, migrationType string) []int {

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
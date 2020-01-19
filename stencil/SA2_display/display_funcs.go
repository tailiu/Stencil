package SA2_display

import (
	"database/sql"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"strconv"
	"time"
	"errors"
	"math/rand"
	"math"
)

const StencilDBName = "stencil"

func RandomNonnegativeInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(math.MaxInt32)
}

func getUserIDByMigrationID(stencilDBConn *sql.DB, migrationID int) string {
	query := fmt.Sprintf("select user_id from migration_registration where migration_id = %d", migrationID)

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprint(data["user_id"])
}

func Initialize(migrationID int, app string) (*sql.DB, *config.AppConfig, int, string) {
	stencilDBConn := db.GetDBConn(StencilDBName)

	app_id := db.GetAppIDByAppName(stencilDBConn, app)

	appConfig, err := config.CreateAppConfigDisplay(app, app_id, stencilDBConn, true)
	if err != nil {
		log.Fatal(err)
	}

	threadID := RandomNonnegativeInt()

	userID := getUserIDByMigrationID(stencilDBConn, migrationID)

	return stencilDBConn, &appConfig, threadID, userID
}

func GetUndisplayedMigratedData(stencilDBConn *sql.DB, migrationID int, appConfig *config.AppConfig) []HintStruct {
	var displayHints []HintStruct

	appID, _ := strconv.Atoi(appConfig.AppID)
	
	// This is important that table id / group id should also be used to get results in the new design
	// For example, in the one-to-multiple mapping, the same row id has different group ids / table ids
	// Those rows could be displayed differently
	query := fmt.Sprintf(
		"SELECT table_id, array_agg(row_id) as row_ids FROM migration_table where mflag = 1 and app_id = %d and migration_id = %d group by group_id, table_id;",
		appID, migrationID)
	
	data := db.GetAllColsOfRows(stencilDBConn, query)
	// log.Println(data)

	// If we don't use physical schema, both table_name and id are necessary to identify a piece of migratd data.
	// Actually, in our physical schema, row_id itself is enough to identify a piece of migrated data.
	// We use table_name to optimize performance
	for _, data1 := range data {
		displayHints = append(displayHints, TransformRowToHint(appConfig, data1))
	}

	// log.Println(displayHints)
	
	return displayHints
}

func CheckMigrationComplete(stencilDBConn *sql.DB, migrationID int) bool {

	query := fmt.Sprintf("SELECT 1 FROM txn_logs WHERE action_id = %d and action_type='COMMIT' LIMIT 1", migrationID)
	
	data := db.GetAllColsOfRows(stencilDBConn, query)
	
	if len(data) == 0 {
		
		return false
	
	} else {
		
		return true
	}
}

func CheckDisplay(stencilDBConn *sql.DB, appID string, data HintStruct) int64 {
	appID1, err1 := strconv.Atoi(appID)
	if err1 != nil {
		log.Fatal(err1)
	}

	// Here for one group, we only need to check one row_id to see whether the group is displayed or not
	// It should be noted that table_id / group_id should also be considered
	query := fmt.Sprintf("SELECT mflag FROM migration_table WHERE row_id = %d and app_id = %d and table_id = %s", data.RowIDs[0], appID1, data.TableID)
	
	// log.Println("==========")
	// log.Println(query)
	// log.Println(data)
	// log.Println("==========")
	
	data1, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(data1)
	return data1["mflag"].(int64)
}

func Display(stencilDBConn *sql.DB, appID string, dataHints []HintStruct, deletionHoldEnable bool, dhStack [][]int, threadID int) (error, [][]int) {
	var queries []string

	for _, dataHint := range dataHints {

		// This is an optimization to prevent possible path conflict
		// We only need to test one rowID in a data hint
		if CheckDisplay(stencilDBConn, appID, dataHint) == 0 {

			log.Println("Found that there is a path conflict!! When displaying data")
			return errors.New("Path conflict"), dhStack

		}

		rowIDs := dataHint.GetAllRowIDs(stencilDBConn, appID)

		for _, rowID := range rowIDs {
			
			// It should be noted that table_id / group_id should also be considered
			query := fmt.Sprintf("UPDATE migration_table SET mflag = 0, updated_at = now() WHERE row_id = %s and app_id = %s and table_id = %s", 
				fmt.Sprint(rowID["row_id"]), appID, dataHint.TableID)
			
				log.Println("**************************************")
			log.Println(query)
			log.Println("**************************************")
			
			queries = append(queries, query)

		}
	}
	if deletionHoldEnable {
		
		var dhQueries []string
		
		dhQueries, dhStack = AddToDeletionHoldStack(dhStack, dataHints, threadID)
		queries = append(queries, dhQueries...)

	}

	return db.TxnExecute(stencilDBConn, queries), dhStack
}

func alreadyInBag(stencilDBConn *sql.DB, appID string, data HintStruct) bool {

	appID1, err1 := strconv.Atoi(appID)
	if err1 != nil {
		log.Fatal(err1)
	}
	
	// Here for one group, we only need to check one to see whether the group is displayed or not
	// It should be noted that table_id / group_id should also be considered
	query := fmt.Sprintf("SELECT bag FROM migration_table WHERE row_id = %d and app_id = %d and table_id = %s", 
		data.RowIDs[0], appID1, data.TableID)
	// log.Println(query)
	
	data1, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	
	// log.Println(data1)
	
	return data1["bag"].(bool)

}

func PutIntoDataBag(stencilDBConn *sql.DB, 
	appID string, dataHints []HintStruct, userID string) error {
	
	var queries []string

	for _, dataHint := range dataHints {
		
		// Similar to displaying data, this is an optimization to prevent possible path conflict
		// We only need to test one rowID in a data hint
		if alreadyInBag(stencilDBConn, appID, dataHint) {

			log.Println("Found that there is a path conflict!! When putting data in a databag")
			
			return errors.New("Path conflict")
		
		}

		rowIDs := dataHint.GetAllRowIDs(stencilDBConn, appID)
		
		for _, rowID := range rowIDs {

			// It should be noted that table_id / group_id should also be considered
			query := fmt.Sprintf(`UPDATE migration_table SET 
				user_id = %s, bag = true, mark_as_delete = true, mflag = 0, updated_at = now() 
				WHERE row_id = %s and app_id = %s and table_id = %s`,
				userID, fmt.Sprint(rowID["row_id"]), appID, dataHint.TableID)
			
			log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			log.Println(query)
			log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			
			queries = append(queries, query)
		}
	}

	return db.TxnExecute(stencilDBConn, queries)
}

func GetTableNameByTableID(stencilDBConn *sql.DB, tableID string) string {

	iTableID, err := strconv.Atoi(tableID)
	if err != nil {
		log.Fatal(err)
	}

	query := fmt.Sprintf("select table_name from app_tables where pk = %d", iTableID)
	
	data1, err1 := db.DataCall1(stencilDBConn, query)
	if err1 != nil {
		log.Fatal(err1)
	}

	return data1["table_name"].(string)

}

func GetTableIDByTableName(stencilDBConn *sql.DB, tableName, appID string) string {
	appID1, err := strconv.Atoi(appID)
	if err != nil {
		log.Fatal(err)
	}

	query := fmt.Sprintf("select pk from app_tables where app_id = %d and table_name = '%s'", 
		appID1, tableName)

	// log.Println(query)
	
	data1, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	
	// log.Println(data1)
	
	return strconv.FormatInt(data1["pk"].(int64), 10)
}

func CheckAndGetTableNameAndID(stencilDBConn *sql.DB, data *HintStruct, appID string) {
	tableName := data.TableName
	
	tableID := data.TableID
	
	if tableName == "" &&  tableID != "" {
		data.TableName = GetTableNameByTableID(stencilDBConn, tableID)
	} 
	
	if tableName != "" &&  tableID == "" {
		data.TableID = GetTableIDByTableName(stencilDBConn, tableName, appID)
	}

	// log.Println(data.TableID)
	// log.Println(data.TableName)
}
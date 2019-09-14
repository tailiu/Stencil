package display

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
	"strings"
)

const StencilDBName = "stencil"

func RandomNonnegativeInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(math.MaxInt32)
}

func Initialize(app string) (*sql.DB, *config.AppConfig, int) {
	stencilDBConn := db.GetDBConn(StencilDBName)

	app_id := db.GetAppIDByAppName(stencilDBConn, app)

	appConfig, err := config.CreateAppConfigDisplay(app, app_id)
	if err != nil {
		log.Fatal(err)
	}

	threadID := RandomNonnegativeInt()

	return stencilDBConn, &appConfig, threadID
}

func GetUndisplayedMigratedData(stencilDBConn *sql.DB, migrationID int, appConfig *config.AppConfig) []HintStruct {
	var displayHints []HintStruct

	appID, _ := strconv.Atoi(appConfig.AppID)
	query := fmt.Sprintf(
		"SELECT table_id, array_agg(row_id) as row_ids FROM migration_table where mflag = 1 and app_id = %d and migration_id = %d group by group_id, table_id;",
		appID, migrationID)
	
	data := db.GetAllColsOfRows(stencilDBConn, query)
	// log.Println(data)

	// If we don't use physical schema, both table_name and id are necessary to identify a piece of migratd data.
	// Actually, in our physical schema, row_id itself is enough to identify a piece of migrated data.
	// We use table_name to optimize performance
	for _, data1 := range data {
		var rowIDs []int
		s := data1["row_ids"][1:len(data1["row_ids"]) - 1]
		s1 := strings.Split(s, ",")
		for _, strRowID := range s1 {
			rowID, err1 := strconv.Atoi(strRowID)
			if err1 != nil {
				log.Fatal(err1)
			} 
			rowIDs = append(rowIDs, rowID)
		}

		hint := HintStruct{}
		// hint.Table = GetTableNameByTableID(stencilDBConn, data1["table_id"])
		hint.TableID = data1["table_id"]
		hint.RowIDs = rowIDs
		displayHints = append(displayHints, hint)
	}
	log.Println(displayHints)
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
	// Here for one group, we only need to check one to see whether the group is displayed or not
	query := fmt.Sprintf("SELECT mflag FROM migration_table WHERE row_id = %d and app_id = %d", data.RowIDs[0], appID1)
	// log.Println(query)
	data1, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(data1)
	return data1["mflag"].(int64)
}

func Display(stencilDBConn *sql.DB, appID string, dataHints []HintStruct, deletionHoldEnable bool, dhStack [][]int, threadID int) (error, [][]int) {
	var queries []string
	
	appID1, err1 := strconv.Atoi(appID)
	if err1 != nil {
		log.Fatal(err1)
	}

	for _, dataHint := range dataHints {
		t := time.Now().Format(time.RFC3339)
		for _, rowID := range dataHint.RowIDs {

			// This is an optimization to prevent possible path conflict
			if CheckDisplay(stencilDBConn, appID, dataHint) == 0 {
				log.Println("There is a path conflict!!")
				return errors.New("Path conflict"), dhStack
			}
			
			query := fmt.Sprintf("UPDATE migration_table SET mflag = 0, updated_at = '%s' WHERE row_id = %d and app_id = %d", t, rowID, appID1)
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

func GetTableIDByTableName(stencilDBConn *sql.DB, tableName, appID string) int {
	appID1, err := strconv.Atoi(appID)
	if err != nil {
		log.Fatal(err)
	}
	query := fmt.Sprintf("select pk from app_tables where app_id = %d and table_name = '%s'", appID1, tableName)
	data1, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(data1)
	return int(data1["pk"].(int64))
}
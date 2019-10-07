package app_display

import (
	"database/sql"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"strconv"
	"time"
)

const StencilDBName = "stencil"

func Initialize(app string) (*sql.DB, config.AppConfig) {
	stencilDBConn := db.GetDBConn(StencilDBName)

	app_id := db.GetAppIDByAppName(stencilDBConn, app)

	appConfig, err := config.CreateAppConfigDisplay(app, app_id, false)
	if err != nil {
		log.Fatal(err)
	}

	return stencilDBConn, appConfig
}

func getTableIDNamePairsInApp(stencilDBConn *sql.DB, appConfig config.AppConfig) {
	query := fmt.Sprintf("select pk, table_name from app_tables where app_id = %s", appConfig.AppID)

	
}

func GetUndisplayedMigratedData(stencilDBConn *sql.DB, appConfig config.AppConfig, migrationID int) []HintStruct {
	var displayHints []HintStruct
	query := fmt.Sprintf("SELECT table_id, id FROM display_flags WHERE app_id = %s and migration_id = %d and display_flag = true", appConfig.AppID, migrationID)
	data := db.GetAllColsOfRows(stencilDBConn, query)
	// fmt.Println(data)
	for _, oneData := range data {
		hint := HintStruct{}
		table := oneData["table_name"]
		intVal, err := strconv.Atoi(oneData["id"])
		if err != nil {
			log.Fatal(err)
		}
		keyVal := map[string]int{
			"id": intVal,
		}
		hint.Table = table
		hint.KeyVal = keyVal
		displayHints = append(displayHints, hint)
	}
	// fmt.Println(displayHints)
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

func Display(stencilDBConn *sql.DB, appConfig config.AppConfig, dataHints []HintStruct, pks map[string]string) error {
	var queries []string

	for _, dataHint := range dataHints {
		table := dataHint.Table
		query := fmt.Sprintf("UPDATE Display_flags SET display_flag = false, updated_at = now() WHERE app_id = %s and table_name = '%s' and id = %d;",
			appConfig.AppID, table, dataHint.KeyVal["id"])
		log.Println("**************************************")
		log.Println(query)
		log.Println("**************************************")
		queries = append(queries, query)
	}

	return db.TxnExecute(stencilDBConn, queries)
}

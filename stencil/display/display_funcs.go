package display

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

func Initialize(app, app_id string) (*sql.DB, *sql.DB, config.AppConfig, map[string]string) {
	appConfig, err := config.CreateAppConfig(app, app_id)
	if err != nil {
		log.Fatal(err)
	}

	stencilDBConn := db.GetDBConn(StencilDBName)
	appDBConn := db.GetDBConn(app)

	pks := make(map[string]string)
	tables := db.GetTablesOfDB(appDBConn, app)
	for _, table := range tables {
		pk, err := db.GetPrimaryKeyOfTable(appDBConn, table)
		if err != nil {
			fmt.Println(err)
		}
		pks[table] = pk
	}

	return stencilDBConn, appDBConn, appConfig, pks
}

func GetUndisplayedMigratedData(stencilDBConn *sql.DB, app string, migrationID int, pks map[string]string) []HintStruct {
	var displayHints []HintStruct
	query := fmt.Sprintf("SELECT * FROM display_flags WHERE app = '%s' and migration_id = %d and display_flag = false", app, migrationID)
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
			pks[table]: intVal,
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

func Display(stencilDBConn *sql.DB, app string, dataHints []HintStruct, pks map[string]string) error {
	var queries []string

	for _, dataHint := range dataHints {
		table := dataHint.Table
		t := time.Now().Format(time.RFC3339)
		query := fmt.Sprintf("UPDATE Display_flags SET display_flag = true, updated_at = '%s' WHERE app = '%s' and table_name = '%s' and id = %d;",
			t, app, table, dataHint.KeyVal[pks[table]])
		log.Println("**************************************")
		log.Println(query)
		log.Println("**************************************")
		queries = append(queries, query)
	}

	return db.TxnExecute(stencilDBConn, queries)
}

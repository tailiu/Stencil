package display

import (
	"log"
	"fmt"
	"database/sql"
	"transaction/db"
	"transaction/config"
	"strconv"
)

const StencilDBName = "stencil"

func Initialize(app string) (*sql.DB, *sql.DB, config.AppConfig, map[string]string) {
	appConfig, err := config.CreateAppConfig(app)
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
		keyVal := map[string]int {
			pks[table]:	intVal,
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
		query := fmt.Sprintf("UPDATE Display_flags SET display_flag = true WHERE app = '%s' and table = '%s' and id = %d;",
							app, table, dataHint.KeyVal[pks[table]])
		queries = append(queries, query)
	}

	return db.TxnExecute(stencilDBConn, queries)
}
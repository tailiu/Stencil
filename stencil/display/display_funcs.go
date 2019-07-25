package display

import (
	"database/sql"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"strconv"
	"time"
	"stencil/qr"
)

const StencilDBName = "stencil"

func Initialize(app string) (*sql.DB, *config.AppConfig, map[string]string) {
	stencilDBConn := db.GetDBConn(StencilDBName)

	app_id := GetAppIDByAppName(stencilDBConn, app)

	appConfig, err := config.CreateAppConfig(app, app_id)
	if err != nil {
		log.Fatal(err)
	}

	pks := make(map[string]string)
	tables := db.GetTablesOfDB(appConfig.DBConn, app)
	for _, table := range tables {
		pk, err := db.GetPrimaryKeyOfTable(appConfig.DBConn, table)
		if err != nil {
			fmt.Println(err)
		}
		pks[table] = pk
	}

	return stencilDBConn, &appConfig, pks
}

func GetUndisplayedMigratedData(stencilDBConn *sql.DB, app string, migrationID int, pks map[string]string, appConfig *config.AppConfig) []HintStruct {
	var displayHints []HintStruct
	query := fmt.Sprintf(
		"SELECT d.table_name, d.id FROM row_desc AS r JOIN display_flags AS d on r.rowid = d.id where app = '%s' and migration_id = %d and mflag = 1;",
		app, migrationID)
	// query := fmt.Sprintf("SELECT * FROM display_flags WHERE app = '%s' and migration_id = %d and display_flag = false", app, migrationID)
	data := db.GetAllColsOfRows(stencilDBConn, query)
	fmt.Println(data)

	// If we don't use physical schema, both table_name and id are necessary to identify a piece of migratd data.
	// Actually, in our physical schema, row_id itself is enough to identify a piece of migrated data.
	// We use table_name to optimize performance
	for _, oneData := range data {
		hint := HintStruct{}
		table := oneData["table_name"]

		//********
		log.Println("**************")
		log.Println(oneData["id"])
		qs := qr.CreateQS(appConfig.QR)
		qs.FromSimple(table)
		qs.ColSimple(table + ".*")
		queryRow := qs.GenSQLWith(oneData["id"])
		row, err := db.DataCall1(stencilDBConn, queryRow)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(row)
		log.Println("**************")
		//********

		intVal, err1 := strconv.Atoi(oneData["id"])
		if err1 != nil {
			log.Fatal(err1)
		}
		keyVal := map[string]int{
			pks[table]: intVal,
		}
		// hint.Data = row
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

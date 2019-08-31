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

func Initialize(app string) (*sql.DB, *config.AppConfig, map[string]string) {
	stencilDBConn := db.GetDBConn(StencilDBName)

	app_id := db.GetAppIDByAppName(stencilDBConn, app)

	appConfig, err := config.CreateAppConfig(app, app_id)
	if err != nil {
		log.Fatal(err)
	}

	pks := make(map[string]string)
	// tables := db.GetTablesOfDB(appConfig.DBConn, app)
	// for _, table := range tables {
	// 	pk, err := db.GetPrimaryKeyOfTable(appConfig.DBConn, table)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	pks[table] = pk
	// }

	return stencilDBConn, &appConfig, pks
}

func GetUndisplayedMigratedData(stencilDBConn *sql.DB, app string, migrationID int, appConfig *config.AppConfig) []HintStruct {
	var displayHints []HintStruct

	query := fmt.Sprintf(
		"SELECT d.table_name, d.id FROM row_desc AS r JOIN display_flags AS d on r.rowid = d.id where app = '%s' and migration_id = %d and mflag = 1;",
		app, migrationID)
	// query := fmt.Sprintf("SELECT * FROM display_flags WHERE app = '%s' and migration_id = %d and display_flag = false", app, migrationID)
	data := db.GetAllColsOfRows(stencilDBConn, query)
	// log.Println(data)

	// If we don't use physical schema, both table_name and id are necessary to identify a piece of migratd data.
	// Actually, in our physical schema, row_id itself is enough to identify a piece of migrated data.
	// We use table_name to optimize performance
	for _, data1 := range data {
		table := data1["table_name"]

		hint := HintStruct{}
		hint.Table = table
		// log.Println(GetData1FromPhysicalSchemaByRowID(stencilDBConn, appConfig.QR, table + ".*", table, data1["id"]))
		hint.RowID = data1["id"]

		displayHints = append(displayHints, hint)
		// log.Println(hint)
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
	rowID, err := strconv.Atoi(data.RowID)
	if err != nil {
		log.Fatal(err)
	}
	appID1, err1 := strconv.Atoi(appID)
	if err1 != nil {
		log.Fatal(err1)
	}

	query := fmt.Sprintf("SELECT mflag FROM row_desc WHERE rowid = %d and app_id = %d", rowID, appID1)
	log.Println(query)
	data1, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(data1)
	return data1["mflag"].(int64)
}

func Display(stencilDBConn *sql.DB, appID string, dataHints []HintStruct) error {
	var queries []string
	
	appID1, err1 := strconv.Atoi(appID)
	if err1 != nil {
		log.Fatal(err1)
	}
	for _, dataHint := range dataHints {
		rowID, err := strconv.Atoi(dataHint.RowID)
		if err != nil {
			log.Fatal(err)
		}
		t := time.Now().Format(time.RFC3339)
		query := fmt.Sprintf("UPDATE row_desc SET mflag = 0, updated_at = '%s' WHERE rowid = %d and app_id = %d", t, rowID, appID1)
		log.Println("**************************************")
		log.Println(query)
		log.Println("**************************************")
		queries = append(queries, query)

	}

	return db.TxnExecute(stencilDBConn, queries)
}

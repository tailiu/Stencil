package app_display

import (
	"database/sql"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"strconv"
	// "time"
)

const StencilDBName = "stencil"

func Initialize(app string) (*sql.DB, config.AppConfig) {
	stencilDBConn := db.GetDBConn(StencilDBName)

	app_id := db.GetAppIDByAppName(stencilDBConn, app)

	appConfig, err := config.CreateAppConfigDisplay(app, app_id, stencilDBConn, false)
	if err != nil {
		log.Fatal(err)
	}

	return stencilDBConn, appConfig
}

func GetUndisplayedMigratedData(stencilDBConn *sql.DB, appConfig config.AppConfig, migrationID int) []HintStruct {
	var displayHints []HintStruct
	query := fmt.Sprintf("SELECT table_id, id FROM display_flags WHERE app_id = %s and migration_id = %d and display_flag = true", appConfig.AppID, migrationID)
	data := db.GetAllColsOfRows(stencilDBConn, query)
	// fmt.Println(data)

	for _, data1 := range data {
		displayHints = append(displayHints, app_display.TransformDisplayFlagDataToHint(data1))
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

func Display(stencilDBConn *sql.DB, appConfig config.AppConfig, dataHints []HintStruct) error {
	var queries1 []string
	var queries2 []string

	for _, dataHint := range dataHints {
		query1 := fmt.Sprintf("UPDATE %s SET display_flag = false WHERE id = %d;",
			dataHint.Table, dataHint.KeyVal["id"])
		query2 := fmt.Sprintf("UPDATE Display_flags SET display_flag = false, updated_at = now() WHERE app_id = %s and table_id = %s and id = %d;",
			appConfig.AppID, dataHint.TableID, dataHint.KeyVal["id"])
		log.Println("**************************************")
		log.Println(query1)
		log.Println(query2)
		log.Println("**************************************")
		queries1 = append(queries1, query1)
		queries2 = append(queries2, query2)
	}

	if err := db.TxnExecute(appConfig.DBConn, queries1); err != nil {
		return err
	} else {
		if err := db.TxnExecute(stencilDBConn, queries2); err != nil {
			return err
		} else {
			return nil
		}
	}
}

func CheckDisplay(stencilDBConn *sql.DB, appConfig config.AppConfig, dataHint HintStruct) bool {
	query := fmt.Sprintf("SELECT display_flag from %s where id = %d",
		dataHint.Table, dataHint.KeyVal["id"])
	
	data1, err := db.DataCall1(appConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(data1)
	return !data1["display_flag"].(bool)
}

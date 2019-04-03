package display

import (
	// "log"
	"fmt"
	"database/sql"
	"transaction/db"
)

func GetMigratedData(stencilDBConn *sql.DB, app string, migrationID int) []HintStruct {
	var displayHints []HintStruct
	query := fmt.Sprintf("SELECT * FROM display_flags WHERE app = '%s' and migration_id = %d", app, migrationID)
	data := db.GetAllColsOfRows(stencilDBConn, query)
	fmt.Println(data)

	return displayHints
}

func CheckMigrationComplete(stencilDBConn *sql.DB, migrationID int) bool {
	query := fmt.Sprintf("SELECT 1 FROM txn_log WHERE action_id = %d and action_type='COMMIT' LIMIT 1", migrationID)
	data := db.GetAllColsOfRows(stencilDBConn, query)
	if len(data) == 0 {
		return false
	} else {
		return true
	}
}
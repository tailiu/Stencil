package evaluation

import (
	"time"
	"log"
	"fmt"
	"database/sql"
	"stencil/db"
)

func getStartTime(stencilDBConn *sql.DB, migrationID int64) time.Time {
	query := fmt.Sprintf("select start_time from migration_registration where migration_id = %d;", 
		migrationID)

	startTime, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return startTime["start_time"].(time.Time)
}

func getEndTime(stencilDBConn *sql.DB, migrationID string) {


}

func GetMigrationTime(evalConfig *EvalConfig, migrationID int64) {
	startTime := getStartTime(evalConfig.StencilDBConn, migrationID)
	log.Println(startTime)
}
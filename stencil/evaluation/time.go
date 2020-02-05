package evaluation

import (
	"time"
	"fmt"
	"log"
	"database/sql"
	"stencil/db"
	"stencil/transaction"
)

func getMigrationStartTime(stencilDBConn *sql.DB, 
	migrationID int) time.Time {

	log_txn := new(transaction.Log_txn)
	log_txn.DBconn = stencilDBConn
	log_txn.Txn_id = migrationID
	
	if startTime := log_txn.GetCreatedAt("BEGIN_TRANSACTION"); 
		len(startTime) == 1 {
		
		return startTime[0]
	} else {

		panic("Should never happen here!")
	}
}

func GetMigrationTime(stencilDBConn *sql.DB, 
	migrationID int) time.Duration {
	
	return getMigrationEndTime(stencilDBConn, migrationID).
		Sub(getMigrationStartTime(stencilDBConn, migrationID))

}

func GetDisplayTime(stencilDBConn *sql.DB, 
	migrationID string) time.Duration {
	
	query := fmt.Sprintf(
		`SELECT start_time, end_time FROM display_registration 
		WHERE migration_id = %s`,
		migrationID,
	)

	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if len(data) != 1 {
		log.Fatal("Get more than one row in the display_registration")
	}

	endTime := data[0]["end_time"].(time.Time)
	startTime := data[0]["start_time"].(time.Time)
	displayTime := endTime.Sub(startTime)

	return displayTime

}
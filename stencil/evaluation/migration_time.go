package evaluation

import (
	"time"
	"database/sql"
	"stencil/transaction"
)

func getMigrationStartTime(stencilDBConn *sql.DB, migrationID int) time.Time {
	log_txn := new(transaction.Log_txn)
	log_txn.DBconn = stencilDBConn
	log_txn.Txn_id = migrationID
	if startTime := log_txn.GetCreatedAt("BEGIN_TRANSACTION"); len(startTime) == 1 {
		return startTime[0]
	} else {
		panic("Should never happen here!")
	}
}

func GetMigrationTime(stencilDBConn *sql.DB, migrationID int) time.Duration {
	return getMigrationEndTime(stencilDBConn, migrationID).Sub(getMigrationStartTime(stencilDBConn, migrationID))
}
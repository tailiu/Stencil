package transaction_log

import (
	"database/sql"
	"math/rand"
    "time"
	"transaction/db"
	"log"
	"fmt"
)

const stencilDBName = "stencil"

type Log_txn struct {
	DBconn      *sql.DB
	Txn_id		int
}

func randomNonnegativeInt() int {
	rand.Seed(time.Now().UnixNano())
    return rand.Intn(2147483647)
}

func Begin_transaction() *Log_txn {
	txn_id := randomNonnegativeInt()

	stencilDB := db.GetDBConn(stencilDBName)
	op := fmt.Sprintf("INSERT INTO txn_log (action_id, action_type) VALUES (%d, 'BEGIN_TRANSACTION');", txn_id)
	if _, err := stencilDB.Exec(op); err != nil {
        log.Fatal(err)
	}
	
	return &Log_txn{DBconn: stencilDB, Txn_id: txn_id}
}

func End_transaction(log_txn *Log_txn) {
	op := fmt.Sprintf("INSERT INTO txn_log (action_id, action_type) VALUES (%d, 'END_TRANSACTION');", log_txn.Txn_id)
	if _, err := log_txn.DBconn.Exec(op); err != nil {
        log.Fatal(err)
    }
}

func Log_change(srcAppID, tgtAppID, table, row_id string, log_txn *Log_txn) {
	undo_action := fmt.Sprintf("%s %s %s %s", tgtAppID, srcAppID, table, row_id)
	op := fmt.Sprintf("INSERT INTO txn_log (action_id, action_type, undo_action) VALUES (%d, 'CHANGE', '%s');", log_txn.Txn_id, undo_action)
	if _, err := log_txn.DBconn.Exec(op); err != nil {
        log.Fatal(err)
	}
}
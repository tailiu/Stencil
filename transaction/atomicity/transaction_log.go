package atomicity

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"
	"transaction/db"
)

const StencilDBName = "stencil"

type Log_txn struct {
	DBconn *sql.DB
	Txn_id int
}

func randomNonnegativeInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2147483647)
}

func CreateTxnLogTable() {
	stencilDB := db.GetDBConn(StencilDBName)
	op := `CREATE TABLE txn_log (
			id SERIAL PRIMARY KEY, 
			action_id INT NOT NULL, 
			action_type string NOT NULL CHECK (action_type IN ('COMMIT','ABORT','ABORTED', 'CHANGE', 'BEGIN_TRANSACTION')),
			undo_action string, 
			INDEX action_id_index (action_id),
			display_hint string)`
	if _, err := stencilDB.Exec(op); err != nil {
		log.Fatal(err)
	}
}

func BeginTransaction() *Log_txn {
	txn_id := randomNonnegativeInt()

	stencilDB := db.GetDBConn(StencilDBName)
	op := fmt.Sprintf("INSERT INTO txn_log (action_id, action_type) VALUES (%d, 'BEGIN_TRANSACTION');", txn_id)
	if _, err := stencilDB.Exec(op); err != nil {
		log.Fatal(err)
	}

	return &Log_txn{DBconn: stencilDB, Txn_id: txn_id}
}

func LogChange(undo_action string, display_hint []byte, log_txn *Log_txn) {
	op := fmt.Sprintf("INSERT INTO txn_log (action_id, action_type, undo_action, display_hint) VALUES (%d, 'CHANGE', '%s', '%s');", 
					log_txn.Txn_id, undo_action, display_hint)
	if _, err := log_txn.DBconn.Exec(op); err != nil {
		log.Fatal(err)
	}
}

func LogOutcome(log_txn *Log_txn, outcome string) {
	op := fmt.Sprintf("INSERT INTO txn_log (action_id, action_type) VALUES (%d, '%s');", log_txn.Txn_id, outcome)
	if _, err := log_txn.DBconn.Exec(op); err != nil {
		log.Fatal(err)
	}
}

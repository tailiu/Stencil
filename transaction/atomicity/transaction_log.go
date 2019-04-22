package atomicity

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
	"transaction/db"
)

type Log_txn struct {
	DBconn *sql.DB
	Txn_id int
}

type UndoAction struct {
	OrgTables []string
	DstTables []string
	Data      map[string]interface{}
}

func (self *UndoAction) AddOrgTable(newTable string) []string {
	for _, orgTable := range self.OrgTables {
		if strings.EqualFold(orgTable, newTable) {
			return self.OrgTables
		}
	}
	self.OrgTables = append(self.OrgTables, newTable)
	return self.OrgTables
}

func (self *UndoAction) AddDstTable(newTable string) []string {
	for _, dstTable := range self.DstTables {
		if strings.EqualFold(dstTable, newTable) {
			return self.DstTables
		}
	}
	self.DstTables = append(self.DstTables, newTable)
	return self.DstTables
}

func (self *UndoAction) AddData(key string, val interface{}) map[string]interface{} {
	if len(self.Data) <= 0 {
		self.Data = make(map[string]interface{})
	}
	self.Data[key] = val
	return self.Data
}

func randomNonnegativeInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2147483647)
}

func CreateTxnLogTable() {
	stencilDB := db.GetDBConn(db.STENCIL_DB)
	op := `CREATE TABLE txn_logs (
			id SERIAL PRIMARY KEY, 
			action_id INT NOT NULL, 
			action_type varchar NOT NULL CHECK (action_type IN ('COMMIT','ABORT','ABORTED', 'CHANGE', 'BEGIN_TRANSACTION')),
			undo_action varchar, 
			created_at TIMESTAMP NOT NULL);`
	if _, err := stencilDB.Exec(op); err != nil {
		log.Fatal(err)
	}
}

func BeginTransaction() (*Log_txn, error) {
	txn_id := randomNonnegativeInt()
	t := time.Now().Format(time.RFC3339)

	stencilDB := db.GetDBConn(db.STENCIL_DB)
	op := fmt.Sprintf("INSERT INTO txn_logs (action_id, action_type, created_at) VALUES (%d, 'BEGIN_TRANSACTION', '%s');",
		txn_id, t)
	if _, err := stencilDB.Exec(op); err != nil {
		return nil, err
	}

	return &Log_txn{DBconn: stencilDB, Txn_id: txn_id}, nil
}

func LogChange(undo_action string, log_txn *Log_txn) error {
	t := time.Now().Format(time.RFC3339)
	op := fmt.Sprintf("INSERT INTO txn_logs (action_id, action_type, undo_action, created_at) VALUES (%d, 'CHANGE', '%s', '%s');",
		log_txn.Txn_id, undo_action, t)
	if _, err := log_txn.DBconn.Exec(op); err != nil {
		return err
	}
	return nil
}

func LogOutcome(log_txn *Log_txn, outcome string) error {
	t := time.Now().Format(time.RFC3339)
	op := fmt.Sprintf("INSERT INTO txn_logs (action_id, action_type, created_at) VALUES (%d, '%s', '%s');",
		log_txn.Txn_id, outcome, t)
	if _, err := log_txn.DBconn.Exec(op); err != nil {
		return err
	}
	return nil
}

func CloseDBConn(log_txn *Log_txn) {
	log_txn.DBconn.Close()
}

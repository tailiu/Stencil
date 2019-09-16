package transaction

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"stencil/db"
	"stencil/helper"
	"strings"
	"time"
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

func GenUndoActionJSON(pks []string, srcAppID, dstAppID string) (string, error) {
	undoAction := make(map[string]string)
	undoAction["srcApp"] = srcAppID
	undoAction["dstApp"] = dstAppID
	undoAction["rows"] = strings.Join(helper.DistinctString(pks), ",")
	if undoActionSerialized, err := json.Marshal(undoAction); err == nil {
		return string(undoActionSerialized), nil
	} else {
		return "", err
	}
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

	stencilDB := db.GetDBConn(db.STENCIL_DB)
	op := fmt.Sprintf("INSERT INTO txn_logs (action_id, action_type, created_at) VALUES (%d, 'BEGIN_TRANSACTION', now());",
		txn_id)
	if _, err := stencilDB.Exec(op); err != nil {
		return nil, err
	}

	return &Log_txn{DBconn: stencilDB, Txn_id: txn_id}, nil
}

func LogChange(undo_action string, log_txn *Log_txn) error {
	op := fmt.Sprintf("INSERT INTO txn_logs (action_id, action_type, undo_action, created_at) VALUES (%d, 'CHANGE', '%s', now());",
		log_txn.Txn_id, undo_action)
	if _, err := log_txn.DBconn.Exec(op); err != nil {
		return err
	}
	return nil
}

func LogOutcome(log_txn *Log_txn, outcome string) error {
	op := fmt.Sprintf("INSERT INTO txn_logs (action_id, action_type, created_at) VALUES (%d, '%s', now());",
		log_txn.Txn_id, outcome)
	if _, err := log_txn.DBconn.Exec(op); err != nil {
		return err
	}
	return nil
}

func CloseDBConn(log_txn *Log_txn) {
	db.CloseDBConn(db.STENCIL_DB)
}

func (self *Log_txn) GetCreatedAt(action_type string) []time.Time {
	var result []time.Time

	op := fmt.Sprintf("SELECT created_at FROM txn_logs WHERE action_id = %d and action_type = '%s';",
		self.Txn_id, action_type)
	// fmt.Println(op)
	data, err := db.DataCall(self.DBconn, op)
	if err != nil {
		log.Fatal(err)
	}

	for _, data1 := range data {
		result = append(result, data1["created_at"].(time.Time))
	} 
	// fmt.Println(result)
	return result
}

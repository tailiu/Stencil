package atomicity

import (
	"transaction/db"
	"transaction/qr"
	"fmt"
	"log"
	"database/sql"
	"strings"
)

func rollbackOneRow(QR *qr.QR, undo_action sql.NullString) {
	parameters := strings.Fields(undo_action.String)

	if QR.AppID == "" {
		QR.SetAppID(parameters[0])
	}

	updQ := QR.PhyUpdateAppIDByRowID(parameters[1], parameters[2], []string{parameters[3]})
	fmt.Println(updQ)

	QR.MigrateOneLogicalRow(updQ)
}

func RollbackMigration(txn_id int) {
	stencilDB := db.GetDBConn(stencilDBName)

	getLogRecords := fmt.Sprintf("SELECT action_type, undo_action FROM txn_log WHERE action_id = %d ORDER BY PRIMARY KEY txn_log DESC", txn_id)
	rows, err := stencilDB.Query(getLogRecords)
	if err != nil {
        log.Fatal(err)
	}
	defer rows.Close()

	QR := qr.NewQR("", "stencil")

	for rows.Next() {
		var action_type string
		var undo_action sql.NullString
        if err := rows.Scan(&action_type, &undo_action); err != nil {
            log.Fatal(err)
        }
		// fmt.Printf("%s %s\n", action_type, undo_action)

		switch action_type {
			// case "COMMIT":
			// 	log.Fatal("Can't abort an already completed action.")
			case "ABORT", "ABORTED":
				log.Fatal("Can't abort an already aborted action.")
			case "CHANGE":
				rollbackOneRow(QR, undo_action)
			case "BEGIN_TRANSACTION":
				break
		}
	}

	// LogOutcome(&Log_txn{DBconn: stencilDB, Txn_id: txn_id}, "ABORTED")
}
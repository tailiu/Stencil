package main

import (
	"fmt"
	"log"
	"transaction/db"
	"transaction/qr"
)

func main() {

	stencil := "stencil"
	stencilDB := db.GetDBConn(stencil)

	appName := "diaspora"
	appDB := db.GetDBConn(appName)

	QR := qr.NewQRWithAppName(appName)
	tables := db.DataCall(appDB, "SHOW TABLES FROM "+appName)

	for _, tableRes := range tables {
		table := tableRes["table_name"].(string)
		sql := fmt.Sprintf("SELECT * FROM %s", table)
		for _, row := range db.DataCall(appDB, sql) {
			var cols []string
			var vals []interface{}
			for col, val := range row {
				cols = append(cols, col)
				vals = append(vals, val)
			}
			qi := qr.CreateQI(table, cols, vals, qr.QTInsert)
			if _, err := stencilDB.Begin(); err != nil {
				log.Fatal("ERROR! SOURCE TRANSACTION CAN'T BEGIN:", err)
			} else {
				for qnum, pq := range QR.ResolveInsert(qi) {
					fmt.Println(qnum, pq)
					// if _, err := tx.Exec(pq); err != nil {
					// 	panic("Can't Insert:", err)
					// }
					fmt.Println(qnum, table, " = Inserted!")
				}
				// tx.Commit()
			}
		}

	}

}

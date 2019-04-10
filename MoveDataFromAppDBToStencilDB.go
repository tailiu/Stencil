package main

import (
	"fmt"
	"strings"
	"transaction/db"
	"transaction/qr"

	escape "github.com/tj/go-pg-escape"
)

func main() {
	stencil := "stencil"
	appName := "diaspora"
	appDB := db.GetDBConn(appName)
	stencilDB := db.GetDBConn(stencil)
	QR := qr.NewQRWithAppName(appName)
	tables := db.DataCall(appDB, "SHOW TABLES FROM "+appName)

	for _, tableRes := range tables {
		table := tableRes["table_name"]

		sql := fmt.Sprintf("SELECT * FROM %s", table)

		for rownum, row := range db.DataCall(appName, sql) {
			cols := ""
			vals := ""
			for col, val := range row {
				cols += fmt.Sprintf("%s, ", col)
				vals += strings.TrimPrefix(fmt.Sprintf("%s, ", escape.Literal(val)), "E")
			}
			insql := escape.Escape("INSERT INTO %s (%s) VALUES (%s)", table, strings.Trim(cols, ", "), strings.Trim(vals, ", "))
			// fmt.Println(insql)

			// tx, err := stencilDB.Begin()
			// if err != nil {
			// 	fmt.Println(err)
			// 	panic("ERROR! SOURCE TRANSACTION CAN'T BEGIN")
			// }

			for qnum, pq := range QR.Resolve(insql) {
				fmt.Println(rownum, qnum, pq)

				// if _, err := tx.Exec(pq); err != nil {
				// 	fmt.Println(rownum, qnum, pq)
				// 	fmt.Println(err)
				// 	panic("Can't Insert")
				// } else {
				// 	fmt.Println(rownum, qnum, table, " = Inserted!")
				// }
			}
			// tx.Commit()
		}

	}

}

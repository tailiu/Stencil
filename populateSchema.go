package main

import (
	"fmt"
	"strings"
	"transaction/db"
	"transaction/qr"

	escape "github.com/tj/go-pg-escape"
)

func main() {

	stencilDB := db.GetDBConn("stencil")
	appDB := "app1"
	QR := qr.NewQR(appDB)
	// tables := []string{"customer", "history", "orderr", "new_order", "item", "stock", "order_line"}
	tables := []string{"stock"}

	for _, table := range tables {
		sql := fmt.Sprintf("SELECT * FROM %s", table)
		for rownum, row := range db.DataCall(appDB, sql) {
			cols := ""
			vals := ""
			for col, val := range row {
				cols += fmt.Sprintf("%s, ", col)
				vals += strings.TrimPrefix(fmt.Sprintf("%s, ", escape.Literal(val)), "E")
			}
			insql := escape.Escape("INSERT INTO %s (%s) VALUES (%s)", table, strings.Trim(cols, ", "), strings.Trim(vals, ", "))
			// fmt.Println(insql)

			tx, err := stencilDB.Begin()
			if err != nil {
				fmt.Println(err)
				panic("ERROR! SOURCE TRANSACTION CAN'T BEGIN")
			}
			for qnum, pq := range QR.Resolve(insql) {
				if _, err := tx.Exec(pq); err != nil {
					fmt.Println(rownum, qnum, pq)
					fmt.Println(err)
					panic("Can't Insert")
				} else {
					fmt.Println(rownum, qnum, table, " = Inserted!")
				}
			}
			tx.Commit()
		}
	}

}

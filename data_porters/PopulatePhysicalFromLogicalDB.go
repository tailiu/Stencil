package main

import (
	"fmt"
	"stencil/db"
	"stencil/qr"
	"log"
	"database/sql"
	// escape "github.com/tj/go-pg-escape"
)

func runTx(dbConn *sql.DB, QIs []*qr.QI) bool{
	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal("transaction can't even begin")
	}

	success := true
	
	for _, qi := range QIs {
		query, args := qi.GenSQL()
		fmt.Println(query)
		if _, err := tx.Exec(query, args...); err != nil {
			success = false
			fmt.Println("Some error:", err)
			break
		}
	}
	
	if success {
		tx.Commit()
	}else{
		tx.Rollback()
		fmt.Println("QIs :=v")
		fmt.Println(QIs)
		log.Fatal()
	}
	
	return success
}

func main() {
	appName, appID := "diaspora", "1"
	stencilDB := db.GetDBConn("stencil")
	appDB := db.GetDBConn(appName)
	QR := qr.NewQR(appName, appID)
	tables := db.GetTablesOfDB(appDB, appName)

	for _, table := range tables {
		q := fmt.Sprintf("SELECT * FROM %s", table)
		if ldata, err := db.DataCall(appDB, q); err != nil{
			fmt.Println(q)
			log.Fatal("Some problem with logical data query:", err)
		}else{
			for _, ldatum := range ldata {
				var cols []string
				var vals []interface{}
				for col, val := range ldatum {
					cols, vals = append(cols, col), append(vals, val)
				}
				qi := qr.CreateQI(table, cols, vals, qr.QTInsert)
				qis := QR.ResolveInsert(qi, QR.NewRowId())
				runTx(stencilDB, qis)
			}
		}
	}
}

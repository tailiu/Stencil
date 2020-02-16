package apis

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"stencil/db"
	"stencil/qr"
	"time"
)

func runTx(dbConn *sql.DB, QIs []*qr.QI) bool {
	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal("transaction can't even begin")
	}

	success := true

	for _, qi := range QIs {
		query, args := qi.GenSQL()
		// fmt.Println(query)
		if _, err := tx.Exec(query, args...); err != nil {
			success = false
			fmt.Println("Some error:", err)
			fmt.Println(query, args)
			fmt.Println(qi)
			break
		}
	}

	if success {
		tx.Commit()
	} else {
		tx.Rollback()
		fmt.Println("QIs :=v")
		fmt.Println(QIs)
		log.Fatal()
	}

	return success
}

func _transfer(QR *qr.QR, appDB, stencilDB *sql.DB, table string, limit, offset int64) {
	q := fmt.Sprintf("SELECT * FROM \"%s\" ORDER BY id LIMIT %d OFFSET %d", table, limit, offset)
	if ldata, err := db.DataCall(appDB, q); err != nil {
		fmt.Println(q)
		log.Fatal("Some problem with logical data query:", err)
	} else {
		for _, ldatum := range ldata {
			var cols []string
			var vals []interface{}
			for col, val := range ldatum {
				cols, vals = append(cols, col), append(vals, val)
			}
			qi := qr.CreateQI(table, cols, vals, qr.QTInsert)
			rowid := db.GetNewRowID(stencilDB)
			qis := QR.ResolveInsert(qi, rowid)
			runTx(stencilDB, qis)
		}
	}
}

func transfer(QR *qr.QR, appDB, stencilDB *sql.DB, table string, limit, offset int64) {

	log.Println("Populating ", table)
	// if totalRows, err := db.GetRowCount(appDB, table); err == nil {
	// for offset := int64(0); offset < totalRows; offset += limit {
	// log.Println(fmt.Sprintf(">> %s: %d - %d of %d | Remaining: %d", table, offset, offset+limit, totalRows, totalRows-offset))
	log.Println(fmt.Sprintf(">> %s: %d - %d", table, offset, offset+limit))
	_transfer(QR, appDB, stencilDB, table, limit, offset)
	// }
	// } else {
	// 	log.Fatal("Error while fetching total rows", err)
	// }
	log.Println("Done:", table)
}

func Port(appName, appID, table string, limit, offset int64) {
	rand.Seed(time.Now().UnixNano())
	transfer(qr.NewQR(appName, appID), db.GetDBConn(appName), db.GetDBConn(db.STENCIL_DB), table, limit, offset)
}

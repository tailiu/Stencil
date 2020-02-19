package apis

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"stencil/db"
	"stencil/qr"
	"time"

	"github.com/gookit/color"
)

func printQIs(QIs []*qr.QI) {
	for _, qi := range QIs {
		fmt.Println()
		fmt.Println(qi)
	}
}

func runBulkTx(dbConn *sql.DB, QIs [][]*qr.QI) bool {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal("transaction can't even begin")
	}

	success := true

	queries, args := qr.GenSQLBulk(QIs)

	if queriesLen, argsLen := len(queries), len(args); queriesLen > 0 && argsLen > 0 && queriesLen == argsLen {
		for i := 0; i < queriesLen; i++ {
			color.LightMagenta.Println(queries[i])
			fmt.Println(args[i])
			if _, err := tx.Exec(queries[i], args[i]...); err != nil {
				success = false
				color.Danger.Print("Execution Failed: ")
				log.Fatal(err)
			} else {
				color.Info.Println("Executed")
			}
		}
	} else {
		success = false
		color.Danger.Printf("Mismatched queries and args: %s | %s\n", queriesLen, argsLen)
	}

	if success {
		if err := tx.Commit(); err != nil {
			color.Danger.Print("Commit Failed: ")
			log.Fatal(err)
		} else {
			color.Success.Println("Committed!")
		}
	} else {
		tx.Rollback()
		fmt.Println(queries)
		fmt.Println()
		fmt.Println(args)
		fmt.Println()
		color.Danger.Print("Execution Halted! ")
		log.Fatal()
	}

	return success
}

func runTx(dbConn *sql.DB, QIs []*qr.QI) bool {
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
		printQIs(QIs)
		log.Fatal("Fatal: Not ported!")
	}

	return success
}

func transfer(QR *qr.QR, appDB, stencilDB *sql.DB, table string, limit, offset int64) {
	log.Println(fmt.Sprintf("Populating %s: %d - %d", color.FgYellow.Render(table), offset, offset+limit))
	q := fmt.Sprintf("SELECT * FROM \"%s\" ORDER BY id LIMIT %d OFFSET %d", table, limit, offset)
	if ldata, err := db.DataCall(appDB, q); err != nil {
		fmt.Println(q)
		log.Fatal("Some problem with logical data query:", err)
	} else {
		var groupedQIs [][]*qr.QI
		for _, ldatum := range ldata {
			var cols []string
			var vals []interface{}
			for col, val := range ldatum {
				cols, vals = append(cols, col), append(vals, val)
			}
			qi := qr.CreateQI(table, cols, vals, qr.QTInsert)
			rowid := db.GetNewRowID(stencilDB)
			groupedQIs = append(groupedQIs, QR.ResolveInsert(qi, rowid))
		}
		runBulkTx(stencilDB, groupedQIs)
	}
	color.Notice.Println("Done:", table)
}

func Port(appName, appID, table string, limit, offset int64, appDB, stencilDB *sql.DB) {
	rand.Seed(time.Now().UnixNano())
	transfer(qr.NewQRWithDBConn(appName, appID, stencilDB), appDB, stencilDB, table, limit, offset)
}

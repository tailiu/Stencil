package main

import (
	"fmt"
	"strings"
	"stencil/db"
	"stencil/qr"
	"log"
	"database/sql"
	"math/rand"
	"time"
	"sync"
	// escape "github.com/tj/go-pg-escape"
)

func FilterTablesFromList(tables []string, tablesToRemove[] string) []string{
	var filteredTables []string

	for _, table := range tables {
		remove := false
		for _, tableToRemove := range tablesToRemove {
			if strings.EqualFold(table, tableToRemove){
				remove = true
				break
			}
		}
		if !remove{
			filteredTables = append(filteredTables, table)
		}
	}

	return filteredTables
}

func runTx(dbConn *sql.DB, QIs []*qr.QI) bool{
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
	}else{
		tx.Rollback()
		fmt.Println("QIs :=v")
		fmt.Println(QIs)
		log.Fatal()
	}
	
	return success
}

func transfer(QR *qr.QR, appDB, stencilDB *sql.DB, table string, wg *sync.WaitGroup) {
	
	log.Println("Populating ",table)
	q := fmt.Sprintf("SELECT * FROM \"%s\"", table)
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
			rowid := db.GetNewRowID(stencilDB)
			qis := QR.ResolveInsert(qi, rowid)
			runTx(stencilDB, qis)
		}
	}
	log.Println("Done:", table)
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())

	appName, appID := "diaspora", "1"
	stencilDB := db.GetDBConn(db.STENCIL_DB)
	appDB := db.GetDBConn(appName)
	QR := qr.NewQR(appName, appID)
	tables := db.GetTablesOfDB(appDB, appName)
	// tables = FilterTablesFromList(tables, []string{"messages"})
	// tables := []string{"messages"}
	// log.Fatal(tables)
	for _, table := range tables {
		wg.Add(1)
		go transfer(QR, appDB, stencilDB, table, &wg)
	}
	wg.Wait()
}

package main

import (
	"fmt"
	"log"
	"stencil/db"
	"stencil/qr"
)

func main() {

	stencildb := "stencil"
	stencilDB := db.GetDBConn(stencildb)

	appName, appID := "diaspora", "1"
	appDB := db.GetDBConn(appName)

	QR := qr.NewQR(appName, appID)
	tables := db.DataCall(appDB, "select table_name from information_schema.tables where table_schema = 'public' and table_name = 'users'")
	orderby := "id"

	for _, tableRes := range tables {
		log.Println("Table:", tableRes)
		table := tableRes["table_name"].(string)
		res := db.DataCall1(appDB, fmt.Sprintf("SELECT COUNT(*) as num FROM \"%s\"", table))
		if val, ok := res["num"]; ok {
			rowcount := int(val.(int64))
			log.Println(rowcount)
			limit := 5000
			for i := 0; i < rowcount; i += limit {
				// fmt.Println("processing row", i)
				sql := fmt.Sprintf("SELECT * FROM \"%s\" ORDER BY %s ASC LIMIT %d OFFSET %d", table, orderby, limit, i)
				for count, row := range db.DataCall(appDB, sql) {
					// fmt.Println("processing rowcount", count)
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
						resolvedQI, pk := QR.ResolveInsert(qi)
						fmt.Print("@@ PQ PK => ", pk)
						for qnum, pqi := range resolvedQI {
							fmt.Print("@@ PQ => ")
							pq, args := pqi.GenSQL()
							fmt.Println(pq)
							if _, err := stencilDB.Exec(pq, args...); err != nil {
								log.Fatal("Can't Insert:", err)
							}
							log.Println(fmt.Sprintf("Row # %d | Inserted into * %s *;", qnum, table))
						}
						// tx.Commit()
					}
				}
			}
		}

	}
	// log.Fatal("stop here")

}

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
	tables := db.DataCall(appDB, "select table_name from information_schema.tables where table_schema = 'public' and table_name = 'users'")
	orderby := "id"
	for _, tableRes := range tables {
		log.Println("Table:", tableRes)
		table := tableRes["table_name"].(string)
		res := db.DataCall1(appDB, fmt.Sprintf("SELECT COUNT(*) as num FROM \"%s\"", table))
		if val, ok := res["num"]; ok {
			rowcount := int(val.(int64))
			// log.Fatal(rowcount)
			log.Println(rowcount)
			limit := 50000
			for i := 0; i < rowcount; i += limit {
				sql := fmt.Sprintf("SELECT * FROM \"%s\" ORDER BY %s ASC LIMIT %d OFFSET %d", table, orderby, limit, i)
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
						for qnum, pqi := range QR.ResolveInsert(qi) {
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

package display

import (
	"stencil/db"
	"database/sql"
	"stencil/qr"
	"log"
	// "fmt"
	"strings"
	"strconv"
)

func GetData1FromPhysicalSchema(stencilDBConn *sql.DB, QR *qr.QR, cols, from, col, op, val string) map[string]interface{}  {	
	qs := qr.CreateQS(QR)
	qs.FromSimple(from)
	qs.ColSimple(cols)
	qs.ColPK(from)
	qs.WhereSimpleVal(col, op, val)

	physicalQuery := qs.GenSQL()
	// log.Println(physicalQuery)

	result, err := db.DataCall1(stencilDBConn, physicalQuery)
	if err != nil {
		log.Fatal(err)
	}
	
	return result
}

func GetData1FromPhysicalSchemaByRowID(stencilDBConn *sql.DB, QR *qr.QR, cols, from, rowid string) map[string]interface{} {	
	qs := qr.CreateQS(QR)
	qs.FromSimple(from)
	qs.ColSimple(cols)
	physicalQuery := qs.GenSQLWith(rowid)

	result, err := db.DataCall1(stencilDBConn, physicalQuery)
	if err != nil {
		log.Fatal(err)
	}
	
	return result
}

func GetRowIDFromData(data map[string]interface{}) string {
	for key, val := range data {
		if strings.Contains(key, "pk.") && val != nil {
			return strconv.FormatInt(data[key].(int64), 10)
		}
	}

	return ""
}

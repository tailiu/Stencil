package display

import (
	"stencil/db"
	"database/sql"
	"stencil/qr"
	"log"
	"strings"
	"strconv"
)

func GetData1FromPhysicalSchema(stencilDBConn *sql.DB, QR *qr.QR, appID, cols, from, col, op, val string) map[string]interface{}  {	
	qs := qr.CreateQSold(QR)
	qs.FromSimple(from)
	qs.ColSimple(cols)
	qs.ColPK(from)
	qs.WhereSimpleVal(col, op, val)
	qs.WhereAppID(qr.EXISTS, appID)
	physicalQuery := qs.GenSQL()
	// log.Println(physicalQuery)

	result, err := db.DataCall1(stencilDBConn, physicalQuery)
	if err != nil {
		log.Fatal(err)
	}
	
	return result
}

func GetData1FromPhysicalSchemaByRowID(stencilDBConn *sql.DB, QR *qr.QR, appID, cols, from, rowid string, restrictions []map[string]string) map[string]interface{} {	
	qs := qr.CreateQSold(QR)
	qs.FromSimple(from)
	qs.ColSimple(cols)
	qs.ColPK(from)
	qs.WhereAppID(qr.EXISTS, appID)
	// qs.WhereSimpleVal(statuses.id, =, #numl)
	// qs.WhereOperatorVal(“AND”, “profiles.bio”, “=”, “student”)
	// qs.WhereOperatorVal(“OR”, “profiles.bio”, “=”, “student”)
	physicalQuery := qs.GenSQLWith(rowid)
	// log.Println(physicalQuery)

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

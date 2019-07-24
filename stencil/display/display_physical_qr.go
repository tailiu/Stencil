package display

import (
	"stencil/db"
	"database/sql"
	"stencil/qr"
	"log"
	"fmt"
)

func GetDataFromPhysicalSchema(stencilDBConn *sql.DB, QR *qr.QR, cols, from, col, op, val, limit string) []map[string]string {	
	qs := qr.CreateQS(QR)
	qs.FromSimple(from)
	qs.ColSimple(cols)
	qs.WhereSimpleVal(col, op, val)
	qs.LimitResult(limit)

	physicalQuery := qs.GenSQL()
	log.Println(physicalQuery)

	return db.GetAllColsOfRows(stencilDBConn, physicalQuery)
}

func GetAppIDByAppName(stencilDBConn *sql.DB, app string) string {
	query := fmt.Sprintf("SELECT pk from apps WHERE app_name = '%s'", app)
	res := db.GetAllColsOfRows(stencilDBConn, query)

	if res[0]["pk"] == "" {
		log.Fatal("AppID does not exist!")
	}

	return res[0]["pk"]
}


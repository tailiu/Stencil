package common_phy_funcs

import (
	"stencil/db"
	"stencil/config"
	"stencil/qr"
	"database/sql"
	"log"
)

func GetRowFromRowIDandTable(stencilDBConn *sql.DB, appConfig *config.AppConfig, row_id string, table string) {
	qs := qr.CreateQS(appConfig.QR)
	qs.FromSimple(table)
	qs.ColSimple(table + ".*")
	queryRow := qs.GenSQLWith(row_id)
	row, err := db.DataCall1(stencilDBConn, queryRow)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(row)
}
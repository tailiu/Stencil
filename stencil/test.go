package main

import (
	"fmt"
	"os"
	"stencil/config"
	"stencil/qr"
	"stencil/db"
	"strings"
)

func test(){
	stencilDB := db.GetDBConn(db.STENCIL_DB)
	appconfig, _ := config.CreateAppConfig("mastodon", "2")
	qs := qr.CreateQS(appconfig.QR)
	qs.FromSimple("statuses")
	qs.FromJoinList("accounts", []string{"statuses.account_id=accounts.id"})
	qs.ColSimple("statuses.*")
	qs.ColSimple("accounts.*")
	qs.LimitResult("10")
	sql := qs.GenSQL()
	fmt.Println(sql)
	fmt.Println("========================================================================================================================================================================")
	fmt.Println(db.DataCall(stencilDB, sql))
}

func test2(){
	stencilDB := db.GetDBConn(db.STENCIL_DB)
	appconfig, _ := config.CreateAppConfig("mastodon", "2")
	qs := qr.CreateQS(appconfig.QR)
	qs.FromSimple("accounts")
	qs.ColSimple("accounts.*")
	qs.WhereSimpleVal("accounts.id","=","1416")
	sql := qs.GenSQL()
	fmt.Println(sql)
	fmt.Println("========================================================================================================================================================================")
	fmt.Println(db.DataCall(stencilDB, sql))
}

func main() {
	test()
	return
	app := os.Args[1]
	table := os.Args[2]
	rowid := os.Args[3] //"653250685"
	with := os.Args[4]
	stencilDB := db.GetDBConn(db.STENCIL_DB)
	var appconfig config.AppConfig
	if strings.EqualFold(app, "diaspora") {
		appconfig, _ = config.CreateAppConfig("diaspora", "1")
	} else if strings.EqualFold(app, "mastodon") {
		appconfig, _ = config.CreateAppConfig("mastodon", "2")
	}

	qs := qr.CreateQS(appconfig.QR)
	qs.FromSimple(table)
	qs.ColSimple(table + ".*")

	sql := ""

	if strings.EqualFold(with, "simple") {
		sql = qs.GenSQL()
		
	} else if strings.EqualFold(with, "with") {
		sql = qs.GenSQLWith(rowid)
	}
	fmt.Println(sql)
	fmt.Println("========================================================================================================================================================================")
	fmt.Println(db.DataCall(stencilDB, sql))

	// return qs.GenSQL()

	// query := qs.GenSQL()

}

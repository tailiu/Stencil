package main

import (
	"fmt"
	"os"
	"stencil/config"
	"stencil/qr"
	"stencil/db"
	"strings"
)

func test3(){
	stencilDB := db.GetDBConn(db.STENCIL_DB)
	appconfig, _ := config.CreateAppConfig("mastodon", "2")
	qs := qr.CreateQS(appconfig.QR)
	qs.SelectColumns("accounts1.*, users.*")
	// qs.SelectColumns("accounts1.*, accounts2.*, users.*")
	qs.FromTable(map[string]string{"table":"accounts", "alias":"accounts1"})
	// qs.JoinTable(map[string]string{"table":"accounts", "alias":"accounts2", "condition1":"accounts1.id=accounts2.id", "condition2":"accounts1.id=accounts2.id"})
	qs.JoinTable(map[string]string{"table":"users", "condition1":"accounts1.id=users.account_id"})
	// qs.AddWhereWithValue("accounts1.id", "=", "1815")
	// qs.AddWhereWithColumn("accounts.id", "=", "accounts.id")
	// qs.AdditionalWhereWithColumn("AND", "accounts.id", "=", "accounts.id")
	sql := qs.GenSQL()
	fmt.Println(sql)
	if res, err := db.DataCall(stencilDB, sql); err == nil{
		for _, r := range res {
			fmt.Println("========================================================================================================================================================================")
			fmt.Println(r)
		}
		fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println("Total Rows:", len(res))
	}else{ 
		fmt.Println(err)
	}
}

func test4(){
	stencilDB := db.GetDBConn(db.STENCIL_DB)
	appconfig, _ := config.CreateAppConfig("diaspora", "1")
	qs := qr.CreateQS(appconfig.QR)
	qs.SelectColumns("people.*")
	qs.FromTable(map[string]string{"table":"people"})
	// qs.JoinTable(map[string]string{"table":"users", "condition1":"people.owner_id=users.id"})
	// qs.RowIDs("870909929")
	qs.AddWhereWithValue("people.id", "=", "23742")
	// qs.RowID("27299317,75744792,2041730142", []string{"accounts"})
	sql := qs.GenSQL()
	fmt.Println(sql)
	if res, err := db.DataCall(stencilDB, sql); err == nil{
		for _, r := range res {
			fmt.Println("========================================================================================================================================================================")
			fmt.Println(r)
		}
		fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println("Total Rows:", len(res))
	}else{ 
		fmt.Println(err)
	}
}

func test5() {
	stencilDB := db.GetDBConn(db.STENCIL_DB)
	appconfig, _ := config.CreateAppConfig("diaspora", "1")
	qs := qr.CreateQS(appconfig.QR)
	qs.SelectColumns("contacts.*")
	qs.FromTable(map[string]string{"table":"contacts"})
	qs.RowIDs("134238299")
	sql := qs.GenSQLSize()
	fmt.Println(sql)
	if res, err := db.DataCall(stencilDB, sql); err == nil{
		for _, r := range res {
			fmt.Println("========================================================================================================================================================================")
			fmt.Println(r)
		}
		fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println("Total Rows:", len(res))
	}else{ 
		fmt.Println(err)
	}
}

func main() {

	test5()
	return

	app := os.Args[1]
	table := os.Args[2]
	// rowid := os.Args[3] //"653250685"
	// with := os.Args[4]
	stencilDB := db.GetDBConn(db.STENCIL_DB)
	var appconfig config.AppConfig
	if strings.EqualFold(app, "diaspora") {
		appconfig, _ = config.CreateAppConfig("diaspora", "1")
	} else if strings.EqualFold(app, "mastodon") {
		appconfig, _ = config.CreateAppConfig("mastodon", "2")
	}

	qs := qr.CreateQS(appconfig.QR)
	qs.FromTable(map[string]string{"table":table})
	qs.SelectColumns(table + ".*")

	sql := qs.GenSQL()
	// sql = qs.GenSQLWith(rowid)
	fmt.Println(sql)
	fmt.Println("========================================================================================================================================================================")
	fmt.Println(db.DataCall(stencilDB, sql))

	// return qs.GenSQL()

	// query := qs.GenSQL()

}

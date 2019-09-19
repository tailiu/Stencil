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
	qs.SelectColumns("statuses.*")
	qs.FromTable(map[string]string{"table":"statuses"})
	// qs.JoinTable(map[string]string{"table":"users", "condition1":"accounts1.id=users.account_id"})
	qs.AddWhereWithValue("statuses.account_id", "=", "24434")
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
	qs.SelectColumns("people.*,users.*,profiles.*")
	qs.FromTable(map[string]string{"table":"people"})
	qs.JoinTable(map[string]string{"table":"users", "condition1":"people.owner_id=users.id"})
	qs.JoinTable(map[string]string{"table":"profiles", "condition1":"profiles.person_id=people.id"})
	qs.RowIDs("2042989316,1201370542,672206183")
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

func test6() {
	stencilDB := db.GetDBConn(db.STENCIL_DB)
	appconfig, _ := config.CreateAppConfig("diaspora", "1")
	qs := qr.CreateQS(appconfig.QR)
	qs.SelectColumns("notifications.*")
	qs.FromTable(map[string]string{"table":"notifications"})
	qs.RowIDs("2031573135,3044756")
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

func main() {

	test6()
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

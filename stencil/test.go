package main

import (
	"fmt"
	"os"
	"stencil/config"
	"stencil/qr"
	"strings"
)

func main() {
	app := os.Args[1]
	table := os.Args[2]
	rowid := os.Args[3] //"653250685"
	with := os.Args[4]

	var appconfig config.AppConfig
	if strings.EqualFold(app, "diaspora") {
		appconfig, _ = config.CreateAppConfig("diaspora", "1")
	} else if strings.EqualFold(app, "mastodon") {
		appconfig, _ = config.CreateAppConfig("mastodon", "2")
	}

	qs := qr.CreateQS(appconfig.QR)
	qs.FromSimple(table)
	qs.ColSimple(table + ".*")

	if strings.EqualFold(with, "simple") {
		fmt.Println(qs.GenSQL())
	} else if strings.EqualFold(with, "with") {
		fmt.Println(qs.GenSQLWith(rowid))
	}

	// return qs.GenSQL()

	// query := qs.GenSQL()

}

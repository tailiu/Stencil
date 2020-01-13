package main

import (
	"stencil/SA1_display"
	"stencil/db"
)

func main() {

	dbName := "mastodon"

	dbConn := db.GetDBConn2(dbName)

	SA1_display.AddDisplayFlagToAllTables(dbConn)

}
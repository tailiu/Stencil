package main

import (
	"stencil/SA1_display"
	"stencil/db"
)

func test1() {

	dbName := "mastodon"

	dbConn := db.GetDBConn(dbName, true)

	SA1_display.AddDisplayFlagToAllTables(dbConn)

}

func test2() {

	dbName := "mastodon"

	dbConn := db.GetDBConn(dbName, true)

	SA1_display.RemoveDisplayFlagInAllTables(dbConn)
}

func main() {

	test1()
	
	// test2()

}
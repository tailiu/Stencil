package main

import (
	"stencil/SA1_display"
	"stencil/db"
)

func test1() {

	dbName := "gnusocial_exp6"

	isBladeServer := false

	dbConn := db.GetDBConn(dbName, isBladeServer)

	SA1_display.AddDisplayFlagToAllTables(dbConn)

}

func test2() {

	dbName := "mastodon"

	dbConn := db.GetDBConn(dbName, true)

	SA1_display.RemoveDisplayFlagInAllTables(dbConn)
}

func test3() {

	dbName := "mastodon_exp3"

	dbConn := db.GetDBConn(dbName, true)

	SA1_display.AddMarkAsDeleteToAllTables(dbConn)
}
 
func main() {

	test1()
	
	// test2()

	// test3()
}
package main

import (
	"stencil/SA1_display"
	"stencil/evaluation"
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

	isBladeServer := true

	dbConn := db.GetDBConn(dbName, isBladeServer)

	SA1_display.RemoveDisplayFlagInAllTables(dbConn)
}

func test3() {

	dbName := "mastodon_exp3"

	isBladeServer := true

	dbConn := db.GetDBConn(dbName, isBladeServer)

	SA1_display.AddMarkAsDeleteToAllTables(dbConn)
}

func test4() {

	dbName := "diaspora_test2"

	isBladeServer := false

	dbConn := db.GetDBConn(dbName, isBladeServer)

	evaluation.AlterTableColumnsIntToInt8(dbConn)

}

func test5() {

	dbName := "gnusocial_exp6"

	isBladeServer := false

	dbConn := db.GetDBConn(dbName, isBladeServer)

	evaluation.AlterTableColumnsAddIDInt8IfNotExists(dbConn)

}

func main() {

	// test1()
	
	// test2()

	// test3()

	// test4()

	test5()

}
package main

import (
	"stencil/reference_resolution_v2"
	"stencil/SA1_display"
	"stencil/evaluation"
	"stencil/db"
)

func test1() {

	dbName := "gnusocial_exp6_3"

	isBladeServer := false

	dbConn := db.GetDBConn(dbName, isBladeServer)
	defer dbConn.Close()

	SA1_display.AddDisplayFlagToAllTables(dbConn)

}

func test2() {

	dbName := "mastodon"

	isBladeServer := true

	dbConn := db.GetDBConn(dbName, isBladeServer)
	defer dbConn.Close()

	SA1_display.RemoveDisplayFlagInAllTables(dbConn)
}

func test3() {

	dbName := "mastodon_100k_exp5"

	isBladeServer := true

	dbConn := db.GetDBConn(dbName, isBladeServer)
	defer dbConn.Close()

	SA1_display.AddMarkAsDeleteToAllTables(dbConn)
}

func test4() {

	dbName := "diaspora_test2"

	isBladeServer := false

	dbConn := db.GetDBConn(dbName, isBladeServer)
	defer dbConn.Close()

	evaluation.AlterTableColumnsIntToInt8(dbConn)

}

func test5() {

	dbName := "gnusocial_exp6"

	isBladeServer := false

	dbConn := db.GetDBConn(dbName, isBladeServer)
	defer dbConn.Close()

	evaluation.AlterTableColumnsAddIDInt8IfNotExists(dbConn)

}

func test6() {

	dbName := "gnusocial_exp6"

	col := "urlhash"

	isBladeServer := false

	dbConn := db.GetDBConn(dbName, isBladeServer)
	defer dbConn.Close()

	evaluation.GetTablesContainingCol(dbConn, col)

}

func test7() {

	dbName := "diaspora_100000_int8_template"

	isBladeServer := false

	dbConn := db.GetDBConn(dbName, isBladeServer)
	defer dbConn.Close()

	evaluation.AlterTableColumnsIntToInt8Concurrently(dbConn)

	
}

func test8() {

	dbName := "stencil_exp_template"

	isBladeServer := false

	dbConn := db.GetDBConn(dbName, isBladeServer)
	defer dbConn.Close()

	SA1_display.CreateIDChangesTable(dbConn)
	
}

func test9() {

	dbName := "diaspora_1m_exp6_0"

	isBladeServer := false

	dbConn := db.GetDBConn(dbName, isBladeServer)
	defer dbConn.Close()

	evaluation.DropForeignKeyConstraints(dbConn)
	
}

func test10() {

	dbName := "stencil_exp_sa2_1k_backup"

	isBladeServer := false

	dbConn := db.GetDBConn(dbName, isBladeServer)
	defer dbConn.Close()

	evaluation.CreateDagCounter(dbConn, "dag_counter_1K")
	
}

func test11() {

	dbName := "stencil_exp_sa2_100k_backup"

	isBladeServer := false

	dbConn := db.GetDBConn(dbName, isBladeServer)
	defer dbConn.Close()

	evaluation.CreateFourDagCounterTables(dbConn)

}

func test12() {

	dbName := "stencil_test"
	// dbName := "stencil_exp6_3"

	isBladeServer := false

	dbConn := db.GetDBConn(dbName, isBladeServer)
	defer dbConn.Close()

	reference_resolution_v2.CreateAttributeChangesTable(dbConn)

}

func test13() {

	// dbName := "stencil_test"
	dbName := "stencil_exp6_3"

	isBladeServer := false

	dbConn := db.GetDBConn(dbName, isBladeServer)
	defer dbConn.Close()

	reference_resolution_v2.CreateReferenceTableV2(dbConn)

}

func test14() {

	// dbName := "stencil_test"
	dbName := "stencil_exp6_3"

	isBladeServer := false

	dbConn := db.GetDBConn(dbName, isBladeServer)
	defer dbConn.Close()

	reference_resolution_v2.CreateResolvedReferencesTable(dbConn)

}

func main() {

	test1()
	
	// test2()

	// test3()

	// test4()

	// test5()

	// test6()

	// test7()

	// test8()

	// test9()

	// test10()

	// test11()

	// test12()
	
	// test13()

	// test14()
}
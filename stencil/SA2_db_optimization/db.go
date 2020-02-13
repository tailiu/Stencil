package SA2_db_optimization

import (
	"stencil/db"
	"strings"
	"log"
	"fmt"
)

func TruncateSA2Tables() {

	db.STENCIL_DB = "stencil"

	dbConn := db.GetDBConn(db.STENCIL_DB)

	query1 := `TRUNCATE migration_table`
	
	query2 := `SELECT tablename FROM pg_catalog.pg_tables WHERE 
	schemaname != 'pg_catalog' AND schemaname != 'information_schema';`

	query3 := "TRUNCATE "
	
	data := db.GetAllColsOfRows(dbConn, query2)

	for _, data1 := range data {

		tableName := data1["tablename"]

		if strings.Contains(tableName, "base_") {
			query3 += tableName + ", "
			continue
		}

		if strings.Contains(tableName, "supplementary_") &&
			tableName != "supplementary_tables" {
			query3 += tableName + ", "
			continue
		}

	}
	
	query3 = query3[:len(query3) - 2]

	log.Println(query1)
	log.Println(query3)

	queries := []string{query1, query3} 

	err := db.TxnExecute(dbConn, queries)
	if err != nil {
		log.Fatal(err)
	}

	err1 := dbConn.Close()
	if err1 != nil {
		log.Fatal(err1)
	}	

}

func GetTotalRowCountsOfDB() {

	dbName := "diaspora_1000000_template"

	dbConn := db.GetDBConn(dbName)

	query1 := `SELECT tablename FROM pg_catalog.pg_tables WHERE 
		schemaname != 'pg_catalog' AND schemaname != 'information_schema';`

	data := db.GetAllColsOfRows(dbConn, query1)
	
	// log.Println(data)

	var totalRows int64

	for _, data1 := range data {
		
		tableName := data1["tablename"]

		// references table will cause errors
		if tableName == "references" {
			continue
		}

		query2 := fmt.Sprintf(
			`select count(*) as num from %s`, 
			tableName,
		)

		// log.Println(query2)

		res, err := db.DataCall1(dbConn, query2)
		if err != nil {
			log.Fatal(err)
		}

		// log.Println(res)

		totalRows += res["num"].(int64)
		
	}

	log.Println("Total Rows:", totalRows)

}

func createPartionedMigrationTable() {
	
}
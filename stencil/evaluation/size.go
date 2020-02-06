package evaluation

import (
	"fmt"
	"log"
	"stencil/db"
	"strconv"
	// "time"
)

func getDanglingDataSizeOfMigration(evalConfig *EvalConfig, 
	migrationID string) (int64, int64) {

	var size1, size2 int64

	query1 := fmt.Sprintf(`
		SELECT pg_column_size(data), app FROM data_bags WHERE
		migration_id = %s`, migrationID)
	
	result1, err := db.DataCall(evalConfig.StencilDBConn, query1)
	if err != nil {
		log.Fatal(err)
	}

	for _, data1 := range result1 {

		appID := fmt.Sprint(data1["app"])

		if appID == "1" {
			size1 += data1["pg_column_size"].(int64)
		} else {
			size2 += data1["pg_column_size"].(int64)
		}
		
	}

	return size1, size2

}

func getAllMediaSize(evalConfig *EvalConfig) int64 {

	query := fmt.Sprintf(`SELECT id FROM photos`)
	
	result, err := db.DataCall(evalConfig.DiasporaDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	var size int64

	for _, data1 := range result {
		
		tmp := fmt.Sprint(data1["id"])
		id, err1 := strconv.Atoi(tmp)
		if err1 != nil {
			log.Fatal(err)
		}

		size += calculateMediaSize(
				evalConfig.DiasporaDBConn, 
				"photos",
				id,
				"1",
			)
	}

	return size

}

func getAllRowsSize(evalConfig *EvalConfig) int64 {

	query1 := `SELECT tablename FROM pg_catalog.pg_tables WHERE 
		schemaname != 'pg_catalog' AND schemaname != 'information_schema';`

	data := db.GetAllColsOfRows(evalConfig.DiasporaDBConn, query1)
	
	// log.Println(data)

	var totalSize int64

	for _, data1 := range data {
		
		tableName := data1["tablename"]

		// references table will cause errors
		if tableName == "references" {
			continue
		}

		// Subtract row header size 24 bytes for each row
		query2 := fmt.Sprintf(
			`select sum(pg_column_size(t) - 24) as size from %s as t`, 
			tableName,
		)

		log.Println(query2)

		res, err := db.DataCall(evalConfig.DiasporaDBConn, query2)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(res)

		// There could be no data in some tables like the block table
		if res[0]["size"] == nil {
			continue
		} else {
			totalSize += res[0]["size"].(int64)
		}
		
	}

	return totalSize

}
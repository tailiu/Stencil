package evaluation

import (
	"fmt"
	"log"
	"stencil/db"
	"database/sql"
	"strconv"
	"strings"
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

func getDanglingDataSizeOfApp(evalConfig *EvalConfig,
	appID string) int64 {

	query := fmt.Sprintf(`
		SELECT pg_column_size(data) FROM data_bags WHERE app = %s`, 
		appID,
	)

	// log.Println(query)

	result, err := db.DataCall(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	var size int64

	for _, data1 := range result {
		size += data1["pg_column_size"].(int64)
	}

	return size
}

func getAllMediaSize(dbConn *sql.DB, table, appID string) int64 {

	query := fmt.Sprintf(
		`SELECT id FROM %s`,
		table,
	)
	
	result, err := db.DataCall(dbConn, query)
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
				dbConn, 
				table,
				id,
				appID,
			)
	}

	return size

}

func getAllRowsSize(dbConn *sql.DB) int64 {

	query1 := `SELECT tablename FROM pg_catalog.pg_tables WHERE 
		schemaname != 'pg_catalog' AND schemaname != 'information_schema';`

	data := db.GetAllColsOfRows(dbConn, query1)
	
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

		res, err := db.DataCall(dbConn, query2)
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

func calculateMediaSize(AppDBConn *sql.DB, table string, 
	pKey int, AppID string) int64 {
	
	if AppID == "1" && table == "photos" {

		query := fmt.Sprintf(
			`select remote_photo_name from %s where id = %d`,
			table, pKey)
		
		res, err2 := db.DataCall1(AppDBConn, query)
		if err2 != nil {
			log.Fatal(err2)
		}

		return mediaSize[fmt.Sprint(res["remote_photo_name"])]

	} else if AppID == "2" && table == "media_attachments" {

		query := fmt.Sprintf(
			`select remote_url from %s where id = %d`,
			table, pKey)
		
		res, err2 := db.DataCall1(AppDBConn, query)
		if err2 != nil {
			log.Fatal(err2)
		}

		parts := strings.Split(fmt.Sprint(res["remote_url"]), "/")
		mediaName := parts[len(parts) - 1]
		return mediaSize[mediaName]

	} else if AppID == "3" && table == "tweets" {

		query := fmt.Sprintf(
			`select tweet_media from %s where id = %d`,
			table, pKey)
		
		res, err2 := db.DataCall1(AppDBConn, query)
		if err2 != nil {
			log.Fatal(err2)
		}

		parts := strings.Split(fmt.Sprint(res["tweet_media"]), "/")
		mediaName := parts[len(parts) - 1]
		return mediaSize[mediaName]

	} else if AppID == "4" && table == "file" {

		query := fmt.Sprintf(
			`select url from %s where id = %d`,
			table, pKey)
		
		res, err2 := db.DataCall1(AppDBConn, query)
		if err2 != nil {
			log.Fatal(err2)
		}

		parts := strings.Split(fmt.Sprint(res["url"]), "/")
		mediaName := parts[len(parts) - 1]
		return mediaSize[mediaName]
	
	} else {
		return 0
	}
}

func calculateRowSize(AppDBConn *sql.DB, 
	cols []string, table string, pKey int, 
	AppID string, checkMediaSize bool) int64 {

	selectQuery := "select"
	
	for i, col := range cols {
		selectQuery += " pg_column_size(" + col + ") "
		if i != len(cols) - 1 {
			selectQuery += " + "
		}
		if i == len(cols) - 1{
			selectQuery += " as cols_size "
		}
	}
	
	query := selectQuery + " from " + table + " where id = " + strconv.Itoa(pKey)
	// log.Println(table)
	// log.Println(query)
	
	row, err2 := db.DataCall1(AppDBConn, query)
	if err2 != nil {
		log.Fatal(err2)
	}
	// log.Println(row["cols_size"].(int64))
	// if table == "photos" {
	// 	fmt.Print(fmt.Sprint(pKey) + ":" + fmt.Sprint(calculateMediaSize(AppDBConn, table, pKey, AppID)) + ",")
	// }
	
	var mediaSize int64

	if checkMediaSize {
		mediaSize = calculateMediaSize(AppDBConn, table, pKey, AppID)
	}

	if row["cols_size"] == nil {

		return mediaSize
		
	} else {

		return row["cols_size"].(int64) + mediaSize
		
	}
	
}

func getDanglingObjectsOfMigration(evalConfig *EvalConfig,
	migrationID string) (int64, int64) {

	var size1, size2 int64

	query1 := fmt.Sprintf(`
		SELECT count(*) as num, app FROM data_bags WHERE
		migration_id = %s GROUP BY app`, migrationID)
	
	result1, err := db.DataCall(evalConfig.StencilDBConn, query1)
	if err != nil {
		log.Fatal(err)
	}

	for _, data1 := range result1 {

		appID := fmt.Sprint(data1["app"])
		num := data1["num"]

		if num == nil {
			continue
		}

		if appID == "1" {
			size1 = num.(int64)
		} else {
			size2 = num.(int64)
		}
		
	}

	return size1, size2

}

func getTotalRowCountsOfDB(dbConn *sql.DB) int64 {

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

	return totalRows

}

func getMediaCountsOfApp(dbConn *sql.DB, appName string) int64 {

	var mediaCounts int64

	if table, ok := appMediaTables[appName]; ok {

		query1 := fmt.Sprintf(`select count(*) as num from %s`, table)

		res1, err1 := db.DataCall1(dbConn, query1)
		if err1 != nil {
			log.Fatal(err1)
		}

		if res1["num"] != nil {
			mediaCounts = res1["num"].(int64)
		}

	}

	return mediaCounts

}

func getTotalObjsIncludingMediaOfApp(evalConfig *EvalConfig, 
	appName string) int64 {

	var totalObjs, rowCounts, mediaCounts int64

	switch appName {
	case "diaspora":
		rowCounts = getTotalRowCountsOfDB(evalConfig.DiasporaDBConn)
		mediaCounts = getMediaCountsOfApp(evalConfig.DiasporaDBConn, appName)
	case "mastodon":
		rowCounts = getTotalRowCountsOfDB(evalConfig.MastodonDBConn)
		mediaCounts = getMediaCountsOfApp(evalConfig.MastodonDBConn, appName)
	case "twitter":
		rowCounts = getTotalRowCountsOfDB(evalConfig.TwitterDBConn)
		mediaCounts = getMediaCountsOfApp(evalConfig.TwitterDBConn, appName)
	case "gnusocial":
		rowCounts = getTotalRowCountsOfDB(evalConfig.GnusocialDBConn)
		mediaCounts = getMediaCountsOfApp(evalConfig.GnusocialDBConn, appName)
	default:
		log.Fatal("Cannot find a connection for the app:", appName)
	}

	totalObjs = rowCounts + mediaCounts

	return totalObjs

}

func getTotalRowCountsOfTable(dbConn *sql.DB, tableName string) int64 {
	
	query := fmt.Sprintf(
		`select count(*) as num from %s`, 
		tableName,
	)

	// log.Println(query)

	res, err := db.DataCall1(dbConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if res["num"] == nil {
		return 0
	} else {
		return res["num"].(int64)
	} 

}

func getDanglingObjectsOfApp(evalConfig *EvalConfig, appID string) int64 {

	query := fmt.Sprintf(
		`select count(*) as num from data_bags where app = %s`, 
		appID,
	)

	// log.Println(query)

	res, err := db.DataCall1(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if res["num"] == nil {
		return 0
	} else {
		return res["num"].(int64)
	}

}

func getDanglingObjsIncludingMediaOfSystem(dbConn *sql.DB) int64 {

	query1 := fmt.Sprintf(`select count(*) as num from data_bags`)

	// log.Println(query)

	res1, err1 := db.DataCall1(dbConn, query1)
	if err1 != nil {
		log.Fatal(err1)
	}

	var objsInDB int64

	if res1["num"] != nil {
		objsInDB += res1["num"].(int64)
	}

	query2 := fmt.Sprintf(`select member from data_bags`)

	res2, err2 := db.DataCall(dbConn, query2)
	if err2 != nil {
		log.Fatal(err2)
	}

	for _, data := range res2 {

		table := fmt.Sprint(data["member"])

		if _, ok := mediaTables[table]; ok {
			objsInDB += 1
		}

	}

	return objsInDB

}

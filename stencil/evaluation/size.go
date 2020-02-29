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

func getTotalObjsIncludingMediaOfApp(dbConn *sql.DB, 
	appName string) int64 {
	
	rowCounts := getTotalRowCountsOfDB(dbConn)
	mediaCounts := getMediaCountsOfApp(dbConn, appName)

	totalObjs := rowCounts + mediaCounts

	return totalObjs

}

func getTotalObjsIncludingMediaOfAppInExp7V2(evalConfig *EvalConfig, 
	appName string, enableBags bool) int64 {

	var totalObjs int64

	var dDBConn, mDBConn, tDBConn, gDBConn *sql.DB

	if enableBags {
		dDBConn = evalConfig.DiasporaDBConn
		mDBConn = evalConfig.MastodonDBConn
		tDBConn = evalConfig.TwitterDBConn
		gDBConn = evalConfig.GnusocialDBConn
	} else {
		dDBConn = evalConfig.DiasporaDBConn1
		mDBConn = evalConfig.MastodonDBConn1
		tDBConn = evalConfig.TwitterDBConn1
		gDBConn = evalConfig.GnusocialDBConn1
	}

	switch appName {
	case "diaspora":
		totalObjs = getTotalRowCountsOfDB(dDBConn)
	case "mastodon":
		totalObjs = getTotalRowCountsOfDB(mDBConn)
	case "twitter":
		totalObjs = getTotalRowCountsOfDB(tDBConn)
	case "gnusocial":
		totalObjs = getTotalRowCountsOfDB(gDBConn)
	default:
		log.Fatal("Cannot find a connection for the app:", appName)
	}

	return totalObjs

}

func getTotalObjsIncludingMediaOfAppInExp7(evalConfig *EvalConfig, 
	appName string, enableBags bool) int64 {

	var totalObjs int64

	var dDBConn, mDBConn, tDBConn, gDBConn *sql.DB

	if enableBags {
		dDBConn = evalConfig.DiasporaDBConn
		mDBConn = evalConfig.MastodonDBConn
		tDBConn = evalConfig.TwitterDBConn
		gDBConn = evalConfig.GnusocialDBConn
	} else {
		dDBConn = evalConfig.DiasporaDBConn1
		mDBConn = evalConfig.MastodonDBConn1
		tDBConn = evalConfig.TwitterDBConn1
		gDBConn = evalConfig.GnusocialDBConn1
	}

	switch appName {
	case "diaspora":
		totalObjs = getTotalObjsIncludingMediaOfApp(dDBConn, appName)
	case "mastodon":
		totalObjs = getTotalObjsIncludingMediaOfApp(mDBConn, appName)
	case "twitter":
		totalObjs = getTotalObjsIncludingMediaOfApp(tDBConn, appName)
	case "gnusocial":
		totalObjs = getTotalObjsIncludingMediaOfApp(gDBConn, appName)
	default:
		log.Fatal("Cannot find a connection for the app:", appName)
	}

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

func getDanglingObjsIncludingMediaOfSystem(dbConn *sql.DB, 
	toApp string, totalMediaInMigrations int64, enableBags bool) int64 {

	query1 := fmt.Sprintf(`select count(*) as num from data_bags`)

	// log.Println(query)

	res1, err1 := db.DataCall1(dbConn, query1)
	if err1 != nil {
		log.Fatal(err1)
	}

	var objsInDB, mediaObjs, totalDanglingObjs int64

	if res1["num"] != nil {
		objsInDB = res1["num"].(int64)
	}

	// Media cannot be migrated to Twitter based on mappings 
	// media cannot be migrated if option databags is not enabled 
	if toApp == "twitter" || !enableBags {
		mediaObjs = totalMediaInMigrations
	}

	totalDanglingObjs = objsInDB + mediaObjs

	return totalDanglingObjs

}

func getDanglingObjsOfSystemV2(dbConn *sql.DB, 
	toApp string, enableBags bool, migrationID, 
	migratedUserID , toAppID, srcUserID, 
	fromAppID string) []map[string]interface{} {

	// Dangling objects put by display threads
	query1 := fmt.Sprintf(
		`select * from data_bags where 
		migration_id = %s and user_id = %s and app = %s`, 
		migrationID, migratedUserID, toAppID,
	)

	log.Println(query1)

	res1, err1 := db.DataCall(dbConn, query1)
	if err1 != nil {
		log.Fatal(err1)
	}

	query2 := fmt.Sprintf(
		`select * from data_bags where 
		migration_id = %s and user_id != %s and app = %s`, 
		migrationID, srcUserID, fromAppID,
	)

	log.Println(query2)

	res2, err2 := db.DataCall(dbConn, query2)
	if err1 != nil {
		log.Fatal(err2)
	}

	res1 = append(res1, res2...)

	return res1

}

func throughTwitter(migrationSeq []string) bool {

	for i, app := range migrationSeq {
		if app == "twitter" && i != len(migrationSeq) - 1{
			return true
		}
	}

	return false
}

func calculateDanglingAndTotalObjectsInExp7(
	evalConfig *EvalConfig, enableBags bool,
	totalMediaInMigrations, totalRemainingObjsInOriginalApp int64,
	toApp string, seqNum int, migrationSeq []string) map[string]int64 {

	var stencilDBConn *sql.DB

	if enableBags {
		stencilDBConn = evalConfig.StencilDBConn
	} else {
		stencilDBConn = evalConfig.StencilDBConn1
	}

	danglingObjs := getDanglingObjsIncludingMediaOfSystem(stencilDBConn,
		toApp, totalMediaInMigrations, enableBags)
	totalObjs := getTotalObjsIncludingMediaOfAppInExp7(evalConfig, 
		toApp, enableBags)

	seqLen := len(migrationSeq)

	// Only when the final application is Diaspora do we need to do this
	if seqNum == seqLen - 2 && toApp == "diaspora" {

		// If the option databags is not enabled and through *twitter*
		// then the total objs should not include migrated media
		if !enableBags && throughTwitter(migrationSeq) {
			totalObjs -= totalMediaInMigrations
		}

		totalObjs -= totalRemainingObjsInOriginalApp

	}

	objs := make(map[string]int64)
	objs["danglingObjs"] = danglingObjs
	objs["totalObjs"] = totalObjs

	return objs

}

func calculateDanglingAndTotalObjectsInExp7v2(
	evalConfig *EvalConfig, enableBags bool, 
	totalRemainingObjsInOriginalApp int64,
	toApp string, seqNum int, migrationSeq []string, migrationID, 
	migratedUserID, toAppID, srcUserID,
	fromAppID string) ([]map[string]interface{}, int64) {

	var stencilDBConn *sql.DB

	if enableBags {
		stencilDBConn = evalConfig.StencilDBConn
	} else {
		stencilDBConn = evalConfig.StencilDBConn1
	}

	danglingObjs := getDanglingObjsOfSystemV2(stencilDBConn,
		toApp, enableBags, migrationID, 
		migratedUserID, toAppID, srcUserID, fromAppID,
	)
	
	totalObjs := getTotalObjsIncludingMediaOfAppInExp7V2(evalConfig, 
		toApp, enableBags)

	seqLen := len(migrationSeq)

	// Only when the final application is Diaspora do we need to do this
	if seqNum == seqLen - 2 && toApp == "diaspora" {

		totalObjs -= totalRemainingObjsInOriginalApp

	}

	return danglingObjs, totalObjs

}

func removeMigratedDanglingObjsFromDataBags(
	evalConfig *EvalConfig, 
	totalDanglingObjs []map[string]interface{}) []map[string]interface{} {

	query1 := fmt.Sprintf(`SELECT pk, migration_id FROM data_bags`)
	
	res, err2 := db.DataCall(evalConfig.StencilDBConn, query1)
	if err2 != nil {
		log.Fatal(err2)
	}

	var deletedObjsIndex []int

	for i, obj1 := range totalDanglingObjs {

		pk2 := fmt.Sprint(obj1["pk"])
		migrationID2 := fmt.Sprint(obj1["migration_id"])

		foundObj := false
		migratedPartially := false

		for _, res1 := range res {

			pk1 := fmt.Sprint(res1["pk"])
			migrationID1 := fmt.Sprint(res1["migration_id"])
	
			if pk1 == pk2 {
				
				foundObj = true

				if migrationID2 != migrationID1 {

					log.Println("partially migrated!!!")	

					migratedPartially = true

				} 
				
				break
			}
		}

		if !foundObj || migratedPartially {
			deletedObjsIndex = append(deletedObjsIndex, i)
		}

	}

	log.Println("delete objs index length:", len(deletedObjsIndex))

	for m := len(deletedObjsIndex) - 1; m > -1; m-- {

		totalDanglingObjs = append(totalDanglingObjs[:m], totalDanglingObjs[m+1:]...)

	}

	return totalDanglingObjs

}
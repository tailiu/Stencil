package evaluation

import (
	"stencil/db"
	"stencil/config"
	"stencil/transaction"
	"stencil/common_funcs"
	"fmt"
	"database/sql"
	"log"
	"bufio"
	"os"
	"strconv"
	"time"
	"encoding/json"
	"strings"
	"math/rand"
	"sync"
)

func InitializeEvalConfig(isBladeServer ...bool) *EvalConfig {

	evalConfig := new(EvalConfig)
	
	if len(isBladeServer) == 1 {
		connectToDB(evalConfig, isBladeServer[0])
	} else {
		connectToDB(evalConfig)
	}

	evalConfig.MastodonAppID = db.GetAppIDByAppName(evalConfig.StencilDBConn, "mastodon")
	evalConfig.DiasporaAppID = db.GetAppIDByAppName(evalConfig.StencilDBConn, "diaspora")
	evalConfig.Dependencies = dependencies
	evalConfig.TableIDNamePairs = GetTableIDNamePairs(evalConfig.StencilDBConn)
	
	mastodonTableNameIDPairs := make(map[string]string)
	diasporaTableNameIDPairs := make(map[string]string)

	mastodonRes := getTableIDNamePairsInApp(evalConfig.StencilDBConn, evalConfig.MastodonAppID)

	for _, res1 := range mastodonRes {
		mastodonTableNameIDPairs[fmt.Sprint(res1["table_name"])] = fmt.Sprint(res1["pk"])
	}

	evalConfig.MastodonTableNameIDPairs = mastodonTableNameIDPairs

	diasporaRes := getTableIDNamePairsInApp(evalConfig.StencilDBConn, evalConfig.DiasporaAppID)

	for _, res1 := range diasporaRes {
		diasporaTableNameIDPairs[fmt.Sprint(res1["table_name"])] = fmt.Sprint(res1["pk"])
	}

   	evalConfig.DiasporaTableNameIDPairs = diasporaTableNameIDPairs
	
	evalConfig.AllAppNameIDs = common_funcs.GetAppIDNamePairs(evalConfig.StencilDBConn)

	attrNameIDPairsOfApps := make(map[string]map[string]string)
	for appID := range evalConfig.AllAppNameIDs {
		attrNameIDPairsInApp := common_funcs.GetAttrNameIDPairsInApp(evalConfig.StencilDBConn, appID)
		attrNameIDPairsOfApps[appID] = attrNameIDPairsInApp
	}
	evalConfig.AttrNameIDPairsOfApps = attrNameIDPairsOfApps

	// t := time.Now()
	evalConfig.SrcAnomaliesVsMigrationSizeFile, 
	evalConfig.DstAnomaliesVsMigrationSizeFile, 
	evalConfig.InterruptionDurationFile,
	evalConfig.MigrationRateFile,
	evalConfig.MigratedDataSizeFile,
	evalConfig.MigrationTimeFile,
	evalConfig.SrcDanglingDataInSystemFile,
	evalConfig.DstDanglingDataInSystemFile,
	evalConfig.DataDowntimeInStencilFile,
	evalConfig.DataDowntimeInNaiveFile,
	evalConfig.DataBags,
	evalConfig.MigratedDataSizeByDstFile,
	evalConfig.MigrationTimeByDstFile,
	evalConfig.MigratedDataSizeBySrcFile,
	evalConfig.MigrationTimeBySrcFile,
	evalConfig.DanglingDataFile,
	evalConfig.DanglingObjectsFile,
	evalConfig.Diaspora1KCounterFile,
	evalConfig.Diaspora10KCounterFile,
	evalConfig.Diaspora100KCounterFile,
	evalConfig.Diaspora1MCounterFile,
	evalConfig.DataDowntimeInPercentageInStencilFile,
	evalConfig.DataDowntimeInPercentageInNaiveFile = 
		"srcAnomaliesVsMigrationSize",
		"dstAnomaliesVsMigrationSize",
		"interruptionDuration",
		"migrationRate",
		"migratedDataSize",
		"migrationTime",
		"srcSystemDanglingData",
		"dstSystemDanglingData",
		"dataDowntimeInStencil",
		"dataDowntimeInNaive",
		"dataBags",
		"migratedDataSizeByDst",
		"migrationTimeByDst",
		"migratedDataSizeBySrc",
		"migrationTimeBySrc",
		"danglingData",
		"danglingObjects",
		"diaspora1KCounter",
		"diaspora10KCounter",
		"diaspora100KCounter",
		"diaspora1MCounter",
		"dataDowntimeInPercentageInStencil",
		"dataDowntimeInPercentageInNaive"

	return evalConfig
}

func getTableIDNamePairsInApp(stencilDBConn *sql.DB, 
	app_id string) []map[string]interface{} {

	query := fmt.Sprintf(
		`select pk, table_name from app_tables where app_id = %s`, 
		app_id)

	result, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	
	return result
}

func GetTableIDNamePairs(stencilDBConn *sql.DB) map[string]string {
	
	tableIDNamePairs := make(map[string]string)

	query := fmt.Sprintf("select pk, table_name from app_tables;")

	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	for _, data1 := range data {
		tableIDNamePairs[fmt.Sprint(data1["pk"])] = fmt.Sprint(data1["table_name"])
	}
	
	return tableIDNamePairs

}

func GetAllMigrationIDsOfAppWithConds(stencilDBConn *sql.DB, 
	appID string, extraConditions string) []map[string]interface{} {
	
	query := fmt.Sprintf(
		`select * from migration_registration where dst_app = '%s' %s;`, 
		appID, extraConditions)
	// log.Println(query)

	migrationIDs, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return migrationIDs
}

func GetAllMigrationIDs(evalConfig *EvalConfig) []string {

	query := fmt.Sprintf("select migration_id from migration_registration")
	// log.Println(query)

	data, err := db.DataCall(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	var migrationIDs []string
	for _, data1 := range data {
		migrationIDs = append(migrationIDs, fmt.Sprint(data1["migration_id"]))
	} 

	return migrationIDs

}

func GetAllMigrationIDsOrderByEndTime(evalConfig *EvalConfig) []string {

	query := fmt.Sprintf(
		"select migration_id from migration_registration order by end_time")
	// log.Println(query)

	data, err := db.DataCall(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	var migrationIDs []string
	for _, data1 := range data {
		migrationIDs = append(migrationIDs, fmt.Sprint(data1["migration_id"]))
	} 

	return migrationIDs

}

func getMigrationIDBySrcUserID(evalConfig *EvalConfig, 
	userID string) string {

	query := fmt.Sprintf(
		`SELECT migration_id FROM migration_registration 
		WHERE user_id = %s`, userID)
	
	result, err := db.DataCall(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if len(result) != 1 {
		log.Fatal("One user id", userID, "results in more than one migration ids")
	}

	migrationID := fmt.Sprint(result[0]["migration_id"])

	return migrationID

}

func getMigrationIDBySrcUserID1(dbConn *sql.DB, 
	userID string) string {

	query := fmt.Sprintf(
		`SELECT migration_id FROM migration_registration 
		WHERE user_id = %s`, userID)
	
	result, err := db.DataCall(dbConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if len(result) != 1 {
		log.Fatal("One user id", userID, "results in more than one migration ids")
	}

	migrationID := fmt.Sprint(result[0]["migration_id"])

	return migrationID

}

func getAllUserIDsInDiaspora(evalConfig *EvalConfig, 
	orderByUserIDs ...bool) []string {

	var query string

	if len(orderByUserIDs) == 1 && orderByUserIDs[0] {
		query = fmt.Sprintf(`SELECT id FROM people order by id`)
	} else {
		query = fmt.Sprintf(`SELECT id FROM people`)
	}
	
	result, err := db.DataCall(evalConfig.DiasporaDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	var userIDs []string

	for _, data1 := range result {
		userIDs = append(userIDs, fmt.Sprint(data1["id"]))
	}

	return userIDs
}

func GetMigrationData(evalConfig *EvalConfig) []map[string]interface{} {

	query := fmt.Sprintf(
		`select user_id, migration_id, migration_type 
		from migration_registration`)
	
	// log.Println(query)

	data, err := db.DataCall(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return data

}

func getMigrationIDBySrcUserIDMigrationType(dbConn *sql.DB, 
	userID, migrationType string) string {

	var mType string

	switch migrationType {
	case "i":
		mType = "0"
	case "d":
		mType = "3"
	case "n":
		mType = "5"
	default:
		log.Fatal("Cannot find a corresponding migration type")
	}

	query := fmt.Sprintf(
		`SELECT migration_id FROM migration_registration 
		WHERE user_id = %s and migration_type = %s`, 
		userID, mType)
	
	result, err := db.DataCall(dbConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if len(result) != 1 {
		log.Fatal("One user id", userID, "results in more than one migration ids")
	}

	migrationID := fmt.Sprint(result[0]["migration_id"])

	return migrationID

}

func getMigrationIDBySrcUserIDMigrationTypeFromToAppID(stencilDBConn *sql.DB,
	uid, srcAppID, dstAppID, migrationType string) string {

	var mType string

	switch migrationType {
	case "i":
		mType = "0"
	case "d":
		mType = "3"
	case "n":
		mType = "5"
	default:
		log.Fatal("Cannot find a corresponding migration type")
	}

	query := fmt.Sprintf(`select migration_id from migration_registration 
		where user_id = %s and src_app = %s and dst_app = %s and migration_type = %s`,
		uid, srcAppID, dstAppID, mType)

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprint(data["migration_id"])

}

func GetAllMigrationIDsAndTypesOfAppWithConds(stencilDBConn *sql.DB, appID string, 
	extraConditions string) []map[string]interface{} {
	
	query := fmt.Sprintf(
		`select migration_id, is_logical from migration_registration 
		where dst_app = '%s' %s;`, 
		appID, extraConditions)
	// log.Println(query)

	result, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func recreateDBByTemplate(dbConn *sql.DB, 
	dbName string, templateDB string) {

	query1 := fmt.Sprintf(
		"drop database %s", 
		dbName,
	)

	query2 := fmt.Sprintf(
		"create database %s template %s", 
		dbName, templateDB,
	)

	if err1 := db.TxnExecute1(dbConn, query1); err1 != nil {
		log.Fatal(err1)	
	} else {
		if err2 := db.TxnExecute1(dbConn, query2); err2 != nil {
			log.Fatal(err2)
		} else {
			return
		}
	}

}

func ConvertStringtoInt64(data string) int64 {

	res, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	
	return res
}

func ConvertFloat64ToString(data []float64) []string {
	
	var convertedData []string
	
	for _, data1 := range data {
		convertedData = append(convertedData, fmt.Sprintf("%f", data1))
	}
	
	return convertedData
}

func ConvertDurationToString(data []time.Duration) []string {

	var convertedData []string
	
	for _, data1 := range data {
		convertedData = append(convertedData, fmt.Sprintf("%f", data1.Seconds()))
	}
	
	return convertedData
}

func ConvertSingleDurationToString(data time.Duration) string {

	return fmt.Sprintf("%f", data.Seconds())

}

func ConvertMapToJSONString(data map[string]int) string {

	convertedData, err := json.Marshal(data)   
    if err != nil {
        fmt.Println(err.Error())
        log.Fatal()
    }
     
    return string(convertedData)
}

func ConvertMapStringToJSONString(data map[string]string) string {

	convertedData, err := json.Marshal(data)   
    if err != nil {
        fmt.Println(err.Error())
        log.Fatal()
    }
     
    return string(convertedData)
}

func ConvertMapIntToJSONString(data map[string]int) string {

	convertedData, err := json.Marshal(data)   
    if err != nil {
        fmt.Println(err.Error())
        log.Fatal()
    }
     
    return string(convertedData)
}

func ConvertMapInt64ToJSONString(data map[string]int64) string {

	convertedData, err := json.Marshal(data)   
    if err != nil {
        fmt.Println(err.Error())
        log.Fatal()
    }
     
    return string(convertedData)
}

func ConvertInt64ToString(data int64) string {
	return strconv.FormatInt(data, 10)
}

func ConvertInt64ArrToStringArr(data []int64) []string {

	var res []string
	
	for _, data1 := range data {
		res = append(res, ConvertInt64ToString(data1))
	}

	return res

}

func ReadStrLinesFromLog(fileName string, changeDefaultDir ...bool) []string {

	dir := logDir

	if len(changeDefaultDir) > 0 {
		if changeDefaultDir[0] {
			dir = logCounterDir
		}
	}

	file, err := os.Open(dir + fileName)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

	var data []string

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        data = append(data, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
	}
	
	return data

}

func WriteStrToLog(fileName string, data string, 
	changeDefaultDir ...bool) {

	dir := logDir

	if len(changeDefaultDir) > 0 {
		if changeDefaultDir[0] {
			dir = logCounterDir
		}
	}

	f, err := os.OpenFile(dir + fileName, 
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, data); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(f)
}

func WriteStrArrToLog(fileName string, data []string) {

	f, err := os.OpenFile(logDir + fileName, 
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for i, data1 := range data {
		if i != len(data) - 1 {
			data1 += ","
		}
		if _, err := fmt.Fprintf(f, data1); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Fprintln(f)
}

func transformTableKeyToNormalType(tableKey map[string]interface{}) (string, int) {
	src_table := tableKey["src_table"].(string)
	src_id_str := tableKey["src_id"].(string)
	src_id_int, err1 := strconv.Atoi(src_id_str)

	if err1 != nil {
		log.Fatal(err1)
	}

	return src_table, src_id_int
}

func transformTableKeyToNormalTypeInDstApp(tableKey map[string]interface{}) (string, int) {
	src_table := tableKey["dst_table"].(string)
	src_id_str := tableKey["dst_id"].(string)
	src_id_int, err1 := strconv.Atoi(src_id_str)

	if err1 != nil {
		log.Fatal(err1)
	}

	return src_table, src_id_int
}

func getLogicalRow(AppDBConn *sql.DB, table string, pKey int) map[string]interface{} {
	query := fmt.Sprintf("select * from %s where id = %d", table, pKey)
	// log.Println(query)
	row, err2 := db.DataCall1(AppDBConn, query)
	if err2 != nil {
		log.Fatal(err2)
	}
	return row
}

func getTableKeyInLogicalSchemaOfMigrationWithConditions(
	stencilDBConn *sql.DB, migrationID string, 
	side string, conditions string) []map[string]interface{} {

	query := fmt.Sprintf(`select %s_table, %s_id from evaluation 
		where migration_id = '%s' and %s;`, 
		side, side, migrationID, conditions)
	
	log.Println(query)
	
	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func getDependsOnTableKeys(evalConfig *EvalConfig, 
	app, table string) []string {
	
	return evalConfig.Dependencies[app][table]

}

func IncreaseMapValByMap(m1 map[string]int, m2 map[string]int) {
	for k, v := range m2 {
		if _, ok := m1[k]; ok {
			m1[k] += v
		} else {
			m1[k] = v
		}
	}
}

func IncreaseMapValByMapInt64(m1 map[string]int64, m2 map[string]int64) {
	for k, v := range m2 {
		if _, ok := m1[k]; ok {
			m1[k] += v
		} else {
			m1[k] = v
		}
	}
}

func increaseMapValOneByKey(m1 map[string]int, key string) {
	if _, ok := m1[key]; ok {
		m1[key] += 1
	} else {
		m1[key] = 1
	}
}

func getMigratedColsOfApp(stencilDBConn *sql.DB, 
	appID string, migration_id string) map[string][]string {
	
	mCols := make(map[string][]string)

	query := fmt.Sprintf(
		`select src_table, src_cols from evaluation
		 where src_app = '%s' and migration_id = '%s'`,
		appID, migration_id)

	tableCols, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	for _, tableCol := range tableCols {
		table := tableCol["src_table"].(string)
		cols := strings.Split(tableCol["src_cols"].(string), ",")
		if cols1, ok := mCols[table]; ok {
			var newCols []string
			for _, col := range cols {
				unique := true
				for _, col1 := range cols1 {
					if col == col1 {
						unique = false
						break
					}
				}
				if unique {
					newCols = append(newCols, col)
				}
			}
			mCols[table] = append(cols1, newCols...)
		} else {
			mCols[table] = cols
		}
	}

	return mCols
}

func getCountsSystem(dbConn *sql.DB, query string) int64 {
	data, err := db.DataCall(dbConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return data[0]["count"].(int64)
}

func getAllDisplayedData(evalConfig *EvalConfig, migrationID, appID string) []DisplayedData {

	query := fmt.Sprintf(
		`select table_id, array_agg(row_id) as row_ids from migration_table 
		where bag = false and app_id = %s and migration_id = %s and mark_as_delete = false 
		group by group_id, table_id;`,
		appID, migrationID)
	
	data := db.GetAllColsOfRows(evalConfig.StencilDBConn, query)

	var displayedData []DisplayedData

	for _, data1 := range data {

		var rowIDs []string

		s := data1["row_ids"][1:len(data1["row_ids"]) - 1]

		s1 := strings.Split(s, ",")
		
		for _, rowID := range s1 {
			rowIDs = append(rowIDs, rowID)
		}

		data2 := DisplayedData{}
		data2.TableID = data1["table_id"]
		data2.RowIDs = rowIDs
		
		displayedData = append(displayedData, data2)
	}

	// log.Println(displayedData)
	return displayedData
}

func getAppConfig(evalConfig *EvalConfig, app string) *config.AppConfig {

	app_id := db.GetAppIDByAppName(evalConfig.StencilDBConn, app)
	
	appConfig, err := config.CreateAppConfigDisplay(app, app_id, evalConfig.StencilDBConn, true)
	
	if err != nil {
		log.Fatal(err)
	}
	
	return &appConfig
}

func GetTableNameByTableID(evalConfig *EvalConfig, tableID string) string {
	
	query := fmt.Sprintf("select table_name from app_tables where pk = %s", tableID)
	
	data1, err1 := db.DataCall1(evalConfig.StencilDBConn, query)
	
	if err1 != nil {
		log.Fatal(err1)
	}
	
	return data1["table_name"].(string)
}

func getMigrationEndTime(stencilDBConn *sql.DB, migrationID int) time.Time {
	
	log_txn := new(transaction.Log_txn)
	
	log_txn.DBconn = stencilDBConn
	
	log_txn.Txn_id = migrationID
	
	if endTime := log_txn.GetCreatedAt("COMMIT"); len(endTime) == 1 {
		return endTime[0]
	} else {
		panic("Should never happen here!")
	}

}

func oldGetAllDataInDataBag(evalConfig *EvalConfig, 
	migrationID string, appConfig *config.AppConfig) []DataBagData {
	
	query := fmt.Sprintf(
		`select table_id, array_agg(row_id) as row_ids from migration_table 
		where bag = true and app_id = %s and migration_id = %s 
		group by group_id, table_id;`,
		appConfig.AppID, migrationID)
	
	data := db.GetAllColsOfRows(evalConfig.StencilDBConn, query)

	var dataBag []DataBagData
	
	for _, data1 := range data {
		
		var rowIDs []string
		
		s := data1["row_ids"][1:len(data1["row_ids"]) - 1]
		s1 := strings.Split(s, ",")
		
		for _, rowID := range s1 {
			rowIDs = append(rowIDs, rowID)
		}

		dataBagData := DataBagData{}
		dataBagData.TableID = data1["table_id"]
		dataBagData.RowIDs = rowIDs
		
		dataBag = append(dataBag, dataBagData)
	}

	log.Println(dataBag)
	
	return dataBag
}

func closeDBConn(conn *sql.DB) {

	err := conn.Close()
	if err != nil {
		log.Fatal(err)
	}

}

func closeDBConns(evalConfig *EvalConfig) {

	log.Println("Close db connections in the evaluation")

	closeDBConn(evalConfig.StencilDBConn)
	closeDBConn(evalConfig.StencilDBConn1)
	closeDBConn(evalConfig.StencilDBConn2)

	closeDBConn(evalConfig.MastodonDBConn)
	closeDBConn(evalConfig.MastodonDBConn1)
	closeDBConn(evalConfig.MastodonDBConn2)

	closeDBConn(evalConfig.DiasporaDBConn)
	closeDBConn(evalConfig.DiasporaDBConn1)

	closeDBConn(evalConfig.TwitterDBConn)
	closeDBConn(evalConfig.TwitterDBConn1)
	
	closeDBConn(evalConfig.GnusocialDBConn)
	closeDBConn(evalConfig.GnusocialDBConn1)

}

func procRes(res map[string]interface{}) map[string]string {

	procResult := make(map[string]string)

	for k, v := range res {
		procResult[k] = fmt.Sprint(v)
	}

	return procResult

}


func procRes1(res []map[string]interface{}) []map[string]string {

	var procResult []map[string]string

	for _, res1 := range res {

		procResult1 := make(map[string]string)

		for k, v := range res1 {
			procResult1[k] = fmt.Sprint(v)
		}

		procResult = append(procResult, procResult1)
	}

	return procResult

}

func getAllUserIDsSortByPhotosInDiaspora(evalConfig *EvalConfig) []map[string]string {

	query := fmt.Sprintf(`
		SELECT author_id, count(id) AS nums 
		FROM photos GROUP BY author_id 
		ORDER BY nums DESC
	`)

	data, err := db.DataCall(evalConfig.DiasporaDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	var res []map[string]string

	for _, data1 := range data {

		res = append(res, procRes(data1))
	}

	return res

}

func connectToDB(evalConfig *EvalConfig, isBladeServer ...bool) {

	bladeServer := true

	if len(isBladeServer) == 1 {
		bladeServer = isBladeServer[0]
	}

	evalConfig.StencilDBConn = db.GetDBConn(stencilDB)
	evalConfig.StencilDBConn1 = db.GetDBConn(stencilDB1)
	evalConfig.StencilDBConn2 = db.GetDBConn(stencilDB2)

	evalConfig.DiasporaDBConn = db.GetDBConn(diaspora)
	evalConfig.DiasporaDBConn1 = db.GetDBConn(diaspora1)

	evalConfig.MastodonDBConn = db.GetDBConn(mastodon, bladeServer)
	evalConfig.MastodonDBConn1 = db.GetDBConn(mastodon1, bladeServer)
	evalConfig.MastodonDBConn2 = db.GetDBConn(mastodon2, bladeServer)

	evalConfig.TwitterDBConn = db.GetDBConn(twitter, bladeServer)
	evalConfig.TwitterDBConn1 = db.GetDBConn(twitter1, bladeServer)

	evalConfig.GnusocialDBConn = db.GetDBConn(gnusocial, bladeServer)
	evalConfig.GnusocialDBConn1 = db.GetDBConn(gnusocial1, bladeServer)

}

func refreshEvalConfigDBConnections(evalConfig *EvalConfig, 
	isBladeServer ...bool) {

	closeDBConns(evalConfig)
	
	if len(isBladeServer) == 1 {
		connectToDB(evalConfig, isBladeServer[0])
	} else {
		connectToDB(evalConfig)
	}

}

func getDBConnByName(evalConfig *EvalConfig, dbName string) *sql.DB {

	var connection *sql.DB

	switch dbName {
	case stencilDB:
		connection = evalConfig.StencilDBConn
	case stencilDB1:
		connection = evalConfig.StencilDBConn1
	case stencilDB2:
		connection = evalConfig.StencilDBConn2
	case mastodon:
		connection = evalConfig.MastodonDBConn
	case mastodon1:
		connection = evalConfig.MastodonDBConn1
	case mastodon2:
		connection = evalConfig.MastodonDBConn2
	case diaspora:
		connection = evalConfig.DiasporaDBConn
	case twitter:
		connection = evalConfig.TwitterDBConn
	case gnusocial:
		connection = evalConfig.GnusocialDBConn
	case diaspora1:
		connection = evalConfig.DiasporaDBConn1
	case twitter1:
		connection = evalConfig.TwitterDBConn1
	case gnusocial1:
		connection = evalConfig.GnusocialDBConn1
	default:
		log.Fatal("Cannot find a connection by the provided connection name")
	}

	return connection

}

func getMigratedUserID(dbConn *sql.DB, 
	migrationID, dstAppID string) string {

	query1 := fmt.Sprintf(
		`SELECT root_member_id FROM app_root_member 
		WHERE app_id = %s`, 
		dstAppID,
	)

	data1, err := db.DataCall1(dbConn, query1)
	if err != nil {
		log.Fatal(err)
	}

	query2 := fmt.Sprintf(
		`SELECT id FROM display_flags WHERE migration_id = %s and 
		table_id = %s and app_id = %s`,
		migrationID, fmt.Sprint(data1["root_member_id"]), dstAppID,
	)

	data2, err := db.DataCall1(dbConn, query2)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprint(data2["id"])

}

func shuffleSlice(s []string) {
	
	rand.Seed(time.Now().UnixNano())
	
	rand.Shuffle(len(s), func(i, j int) { 
		s[i], s[j] = s[j], s[i] 
	})

}

func moveElementToStartOfSlice(s []string, 
	element string) []string {

	if len(s) == 0 || s[0] == element {
		return s
	} 
	
	if s[len(s)-1] == element {
		s = append([]string{element}, s[:len(s)-1]...)
		return s
	} 

	for i, value := range s {
		if value == element {
			s = append([]string{element}, append((s)[:i], (s)[i+1:]...)...)
			break
		}
	}

	return s

}

func getAllTablesOfDB(dbConn *sql.DB) []string {

	query1 := `SELECT tablename FROM pg_catalog.pg_tables WHERE 
		schemaname != 'pg_catalog' AND schemaname != 'information_schema';`

	data := db.GetAllColsOfRows(dbConn, query1)

	var tables []string 
	
	for _, data1 := range data {
		tables = append(tables, fmt.Sprint(data1["tablename"]))
	}
	
	return tables

}

func AlterTableColumnsIntToInt8(dbConn *sql.DB) {

	tables := getAllTablesOfDB(dbConn)

	for _, table := range tables {

		var columnsToBeUpdated []string

		query1 := fmt.Sprintf(
			`select column_name, data_type from information_schema.columns 
			where table_name = '%s'`, 
			table,
		)

		res1, err1 := db.DataCall(dbConn, query1)
		if err1 != nil {
			log.Fatal(err1)
		}

		// log.Println(res1)

		for _, data1 := range res1 {
			
			if data1["data_type"] == "integer" {
				columnsToBeUpdated = append(columnsToBeUpdated, 
					fmt.Sprint(data1["column_name"]))
			}
			
		}

		var queries []string

		for _, col := range columnsToBeUpdated {

			query2 := fmt.Sprintf(
				`ALTER TABLE %s ALTER COLUMN %s TYPE int8`,
				table, col,
			)
			queries = append(queries, query2)

		}

		err2 := db.TxnExecute(dbConn, queries)
		if err2 != nil {
			log.Fatal(err2)
		}		

		log.Println("Finish Modifying:", table)

	}

}

func alterATableColumnsIntToInt8(dbConn *sql.DB, table string, wg *sync.WaitGroup) {

	defer wg.Done()

	log.Println("Start Modifying:", table)

	var columnsToBeUpdated []string

	query1 := fmt.Sprintf(
		`select column_name, data_type from information_schema.columns 
		where table_name = '%s'`, 
		table,
	)

	res1, err1 := db.DataCall(dbConn, query1)
	if err1 != nil {
		log.Fatal(err1)
	}

	// log.Println(res1)

	for _, data1 := range res1 {
		
		if data1["data_type"] == "integer" {
			columnsToBeUpdated = append(columnsToBeUpdated, 
				fmt.Sprint(data1["column_name"]))
		}
		
	}

	var queries []string

	for _, col := range columnsToBeUpdated {

		query2 := fmt.Sprintf(
			`ALTER TABLE %s ALTER COLUMN %s TYPE int8`,
			table, col,
		)
		queries = append(queries, query2)

	}

	err2 := db.TxnExecute(dbConn, queries)
	if err2 != nil {
		log.Fatal(err2)
	}		

	log.Println("Finish Modifying:", table)

}

func AlterTableColumnsIntToInt8Concurrently(dbConn *sql.DB) {

	var wg sync.WaitGroup

	tables := getAllTablesOfDB(dbConn)

	wg.Add(len(tables))

	for _, table := range tables {

		go alterATableColumnsIntToInt8(dbConn, table, &wg)	

	}

	wg.Wait()

	log.Println("Finish Modifying All Tables")

}

func AlterTableColumnsAddIDInt8IfNotExists(dbConn *sql.DB) {

	tables := getAllTablesOfDB(dbConn)

	for _, table := range tables {

		query1 := fmt.Sprintf(
			`select column_name, data_type from information_schema.columns 
			where table_name = '%s'`, 
			table,
		)

		res1, err1 := db.DataCall(dbConn, query1)
		if err1 != nil {
			log.Fatal(err1)
		}

		// log.Println(res1)

		isIDmissing := true

		for _, data1 := range res1 {
			
			if fmt.Sprint(data1["column_name"]) == "id" {
				isIDmissing = false
			}
			
		}

		if isIDmissing {

			query2 := fmt.Sprintf(
				`ALTER TABLE %s ADD COLUMN id int8`,
				table,
			)
	
			err2 := db.TxnExecute1(dbConn, query2)
			if err2 != nil {
				log.Fatal(err2)
			}

		}

		log.Println("Finish Checking:", table)
	}
}

func GetTablesContainingCol(dbConn *sql.DB, col string) {

	tables := getAllTablesOfDB(dbConn)

	var tablesContainingCol []string

	for _, table := range tables {

		query1 := fmt.Sprintf(
			`select column_name, data_type from information_schema.columns 
			where table_name = '%s'`, 
			table,
		)

		res1, err1 := db.DataCall(dbConn, query1)
		if err1 != nil {
			log.Fatal(err1)
		}

		for _, data1 := range res1 {
			
			if fmt.Sprint(data1["column_name"]) == col {
				tablesContainingCol = append(tablesContainingCol, table)
			}
			
		}

	}

	log.Println(tablesContainingCol)

}

func DropForeignKeyConstraints(dbConn *sql.DB) {

	query1 := fmt.Sprintf(
		`select conrelid::regclass AS table_from, conname, pg_get_constraintdef(c.oid)
		from pg_constraint c join pg_namespace n ON n.oid = c.connamespace
		where contype in ('f') order by table_from;`,
	)
	
	fkConstraints := db.GetAllColsOfRows(dbConn, query1)

	for _, fkConstraint := range fkConstraints {

		query2 := fmt.Sprintf(
			`ALTER TABLE "%s" DROP CONSTRAINT %s;`,
			fkConstraint["table_from"],
			fkConstraint["conname"],
		)

		log.Println(query2)

		err2 := db.TxnExecute1(dbConn, query2)
		if err2 != nil {
			log.Fatal(err2)
		}
	}
}

func DropNotNullConstraints(dbConn *sql.DB) {

	tables := getAllTablesOfDB(dbConn)

	for _, table := range tables {

		if table == "user" || table == "references" ||
		   table == "signature_orders" || table == "chat_offline_messages" {
			continue
		}

		query := fmt.Sprintf(`SELECT dropNull('%s')`, table)
		log.Println(query)

		err := db.TxnExecute1(dbConn, query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func CreateDagCounter(dbConn *sql.DB, tableName string) {

	query := fmt.Sprintf(
		`CREATE TABLE %s (
			person_id varchar,
			edges int8,
			nodes int8)`,
		tableName,
	)

	err2 := db.TxnExecute1(dbConn, query)
	if err2 != nil {
		log.Fatal(err2)
	}

}

func CreateFourDagCounterTables(dbConn *sql.DB) {

	counterTables := []string {
		"dag_counter_1K", 
		"dag_counter_10K", 
		"dag_counter_100K",
		"dag_counter_1M",
	}

	for _, counterTable := range counterTables {
		CreateDagCounter(dbConn, counterTable)
	}

}

func getSrcUserIDByMigrationID(dbConn *sql.DB, migrationID string) string {

	query := fmt.Sprintf(
		`SELECT user_id FROM migration_registration WHERE migration_id = %s`,
		migrationID,
	)

	res, err2 := db.DataCall1(dbConn, query)
	if err2 != nil {
		log.Fatal(err2)
	}

	if res["user_id"] == nil {
		log.Fatal("Cannot find user id!")
	} 
		
	return fmt.Sprint(res["user_id"])

}

func logSeqMigsRes(logFile, userID string, appObjs int64) {

	objs := make(map[string]int64)
	objs["appObjs"] = appObjs
	objs["userID"] = ConvertStringtoInt64(userID)

	WriteStrToLog(
		logFile,
		ConvertMapInt64ToJSONString(objs),
	)

}

func (eval *EvalConfig) getRootMembersOfApps() map[string]string {

	query := `SELECT app_id, root_member_id from app_root_member`

	data, err := db.DataCall(eval.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	rootMembers := make(map[string]string)
	
	for _, data1 := range data {
		rootMembers[fmt.Sprint(data1["app_id"])] = fmt.Sprint(data1["root_member_id"])
	}

	return rootMembers
}

func (eval *EvalConfig) getNextUserID(stencilDBConn *sql.DB, migrationID string) string {

	appRootMembers := eval.getRootMembersOfApps()

	query := fmt.Sprintf(
		`SELECT user_id, src_app, dst_app FROM migration_registration
		WHERE migration_id = %s`, migrationID,
	)

	log.Println(query)

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(data)

	userID := fmt.Sprint(data["user_id"])
	srcApp := fmt.Sprint(data["src_app"])
	dstApp := fmt.Sprint(data["dst_app"])
	srcRootMemberID := appRootMembers[srcApp]
	dstRootMemberID := appRootMembers[dstApp]
	srcRootMemberName := eval.TableIDNamePairs[srcRootMemberID]
	dstRootMemberName := eval.TableIDNamePairs[dstRootMemberID]
	idInSrcApp := eval.AttrNameIDPairsOfApps[srcApp][srcRootMemberName + ":id"]
	idInDstApp := eval.AttrNameIDPairsOfApps[dstApp][dstRootMemberName + ":id"]

	query1 := fmt.Sprintf(
		`SELECT to_id FROM attribute_changes WHERE
		from_app = %s and from_member = %s and from_id = %s and from_attr = %s and from_val = '%s' 
		and to_app = %s and to_member = %s and to_attr = %s and migration_id = %s`,
		srcApp, srcRootMemberID, userID, idInSrcApp, userID,
		dstApp, dstRootMemberID, idInDstApp, migrationID, 
	)

	log.Println(query1)

	data1, err1 := db.DataCall1(stencilDBConn, query1)
	if err1 != nil {
		log.Fatal(err1)
	}

	if data1["to_id"] == nil {
		log.Fatal("Cannot get user ID in the destination application!")
	}

	return fmt.Sprint(data1["to_id"])

}

func (eval *EvalConfig) getRootUserIDsRandomly(app string, num int) []string {

	var query string
	var err error
	var data []map[string]interface{}

	switch app {
	case "diaspora":
		query = fmt.Sprintf("SELECT id from people ORDER BY random() LIMIT %d", num)
		data, err = db.DataCall(eval.DiasporaDBConn, query)
	case "mastodon":
		query = fmt.Sprintf("SELECT id from accounts ORDER BY random() LIMIT %d", num)
		data, err = db.DataCall(eval.MastodonDBConn, query)
	case "twitter":
		query = fmt.Sprintf("SELECT id from users ORDER BY random() LIMIT %d", num)
		data, err = db.DataCall(eval.TwitterDBConn, query)
	case "gnusocial":
		query = fmt.Sprintf("SELECT id from profile ORDER BY random() LIMIT %d", num)
		data, err = db.DataCall(eval.GnusocialDBConn, query)
	default:
		log.Fatal("Unrecognized application!")
	}

	if err != nil {
		log.Fatal(err)
	}

	var res []string
	for _, data1 := range data {
		res = append(res, fmt.Sprint(data1["id"]))
	}

	return res
}

func setDatabasesLogFilesForExp7(startApp, size string) (string, string) {

	var logFile, logFile1 string

	stencilDB = "stencil_exp7_0"
	diaspora = "diaspora_exp7_0"
	mastodon = "mastodon_exp7_0"
	twitter = "twitter_exp7_0"
	gnusocial = "gnusocial_exp7_0"

	stencilDB1 = "stencil_exp7_1"
	diaspora1 = "diaspora_exp7_1"
	mastodon1 = "mastodon_exp7_1"
	twitter1 = "twitter_exp7_1"
	gnusocial1 = "gnusocial_exp7_1"

	switch startApp {
	case "diaspora":
		diaspora = "diaspora_exp7_0_" + size
		diaspora1 = "diaspora_exp7_1_" + size
		logFile = "dataBagsEnabled_diaspora"
		logFile1 = "dataBagsNotEnabled_diaspora"
	case "mastodon":
		mastodon = "mastodon_exp7_0_" + size
		mastodon1 = "mastodon_exp7_1_" + size
		logFile = "dataBagsEnabled_mastodon"
		logFile1 = "dataBagsNotEnabled_mastodon"
	case "twitter":
		twitter = "twitter_exp7_0_" + size
		twitter1 = "twitter_exp7_1_" + size
		logFile = "dataBagsEnabled_twitter"
		logFile1 = "dataBagsNotEnabled_twitter"
	case "gnusocial":
		gnusocial = "gnusocial_exp7_0_" + size
		gnusocial1 = "gnusocial_exp7_1_" + size
		logFile = "dataBagsEnabled_gnusocial"
		logFile1 = "dataBagsNotEnabled_gnusocial"
	default:
		log.Fatal("Unrecognized application!")
	}
	
	return logFile, logFile1
}


package evaluation

import (
	"stencil/db"
	"stencil/config"
	"stencil/transaction"
	"fmt"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"
	"encoding/json"
	"strings"
)

func InitializeEvalConfig() *EvalConfig {

	evalConfig := new(EvalConfig)
	evalConfig.StencilDBConn = db.GetDBConn(stencilDB)
	evalConfig.MastodonDBConn = db.GetDBConn(mastodon, true)
	evalConfig.DiasporaDBConn = db.GetDBConn(diaspora)
	evalConfig.MastodonAppID = db.GetAppIDByAppName(evalConfig.StencilDBConn, mastodon)
	evalConfig.DiasporaAppID = db.GetAppIDByAppName(evalConfig.StencilDBConn, diaspora)
	evalConfig.Dependencies = dependencies
	evalConfig.TableIDNamePairs = GetTableIDNamePairs(evalConfig.StencilDBConn)
	
	mastodonTableNameIDPairs := make(map[string]string)
	diasporaTableNameIDPairs := make(map[string]string)

	mastodonRes := getTableIDNamePairsInApp(evalConfig.StencilDBConn,
		 evalConfig.MastodonAppID)

	for _, res1 := range mastodonRes {

		mastodonTableNameIDPairs[fmt.Sprint(res1["table_name"])] = 
			fmt.Sprint(res1["pk"])
	}

	evalConfig.MastodonTableNameIDPairs = mastodonTableNameIDPairs

	diasporaRes := getTableIDNamePairsInApp(evalConfig.StencilDBConn,
		evalConfig.DiasporaAppID)

	for _, res1 := range diasporaRes {

		diasporaTableNameIDPairs[fmt.Sprint(res1["table_name"])] = 
			fmt.Sprint(res1["pk"])
	}

   	evalConfig.DiasporaTableNameIDPairs = diasporaTableNameIDPairs

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
	evalConfig.DataBags = 
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
		"dataBags"

	return evalConfig
}

func getTableIDNamePairsInApp(stencilDBConn *sql.DB, app_id string) []map[string]interface{} {

	query := fmt.Sprintf("select pk, table_name from app_tables where app_id = %s", app_id)

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
	
	query := fmt.Sprintf("select * from migration_registration where dst_app = '%s' %s;", 
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

func GetMigrationData(evalConfig *EvalConfig) []map[string]interface{} {

	query := fmt.Sprintf("select user_id, migration_id from migration_registration")
	// log.Println(query)

	data, err := db.DataCall(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return data

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

func ConvertMapInt64ToJSONString(data map[string]int64) string {
	convertedData, err := json.Marshal(data)   
    if err != nil {
        fmt.Println(err.Error())
        log.Fatal()
    }
     
    return string(convertedData)
}

func WriteStrToLog(fileName string, data string) {

	f, err := os.OpenFile(logDir + fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, data); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(f)
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

func closeDBConns(evalConfig *EvalConfig) {

	evalConfig.StencilDBConn.Close()
	evalConfig.MastodonDBConn.Close()
	evalConfig.DiasporaDBConn.Close()

}
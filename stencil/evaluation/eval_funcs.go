package evaluation

import (
	"stencil/db"
	"fmt"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"
	"encoding/json"
	"strings"
)

const logDir = "./evaluation/logs/"

func GetAllMigrationIDsOfAppWithConds(stencilDBConn *sql.DB, appID string, extraConditions string) []map[string]interface{} {
	query := fmt.Sprintf("select * from migration_registration where dst_app = '%s' %s;", 
		appID, extraConditions)
	// log.Println(query)

	migrationIDs, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return migrationIDs
}

func GetAllMigrationIDsAndTypesOfAppWithConds(stencilDBConn *sql.DB, appID string, extraConditions string) []map[string]interface{} {
	query := fmt.Sprintf("select migration_id, is_logical from migration_registration where dst_app = '%s' %s;", 
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

func WriteStrArrToLog(fileName string, data []string) {
	f, err := os.OpenFile(logDir + fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func calculateRowSize(AppDBConn *sql.DB, cols []string, table string, pKey int) int64 {
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
	return row["cols_size"].(int64)
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

func getTableKeyInLogicalSchemaOfMigrationWithConditions(stencilDBConn *sql.DB, migrationID string, side string, conditions string) []map[string]interface{} {
	query := fmt.Sprintf("select %s_table, %s_id from evaluation where migration_id = '%s' and %s;", 
		side, side, migrationID, conditions)
	
	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func getDependsOnTableKeys(evalConfig *EvalConfig, app, table string) []string {
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

func getMigratedColsOfApp(stencilDBConn *sql.DB, appID string, migration_id string) map[string][]string {
	mCols := make(map[string][]string)

	query := fmt.Sprintf("select src_table, src_cols from evaluation where src_app = '%s' and migration_id = '%s'",
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
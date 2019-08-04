package evaluation

import (
	"stencil/db"
	"fmt"
	"database/sql"
	"log"
	"os"
	"strconv"
	// "strings"
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

func ConvertFloat64ToString(data []float64) []string {
	var convertedData []string
	for _, data1 := range data {
		convertedData = append(convertedData, fmt.Sprintf("%f", data1))
	}
	return convertedData
}

func WriteToLog(fileName string, data []string) {
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
	log.Println(table)
	log.Println(query)
	row, err2 := db.DataCall1(AppDBConn, query)
	if err2 != nil {
		log.Fatal(err2)
	}
	log.Println(row["cols_size"].(int64))
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

func getLogicalRow(AppDBConn *sql.DB, table string, pKey int) map[string]interface{} {
	query := fmt.Sprintf("select * from %s where id = %d", table, pKey)
	row, err2 := db.DataCall1(AppDBConn, query)
	if err2 != nil {
		log.Fatal(err2)
	}
	return row
}
package evaluation

import (
	"stencil/db"
	"fmt"
	"database/sql"
	"log"
	"strconv"
)

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

func GetDataInLogicalSchemaOfMigration(stencilDBConn *sql.DB, AppDBConn *sql.DB, migrationID string, side string) []map[string]interface{} {
	query := fmt.Sprintf("select %s_table, %s_id from evaluation where migration_id = '%s'", side, side, migrationID)
	
	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return data
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
	// log.Println(query)
	row, err2 := db.DataCall1(AppDBConn, query)
	if err2 != nil {
		log.Fatal(err2)
	}
	return row["cols_size"].(int64)
}

func getEntireRowSize(AppDBConn *sql.DB, data1 map[string]interface{}) int64 {
	table, pKey := transformTableKeyToNormalType(data1)
	query := fmt.Sprintf("select * from %s where id = %d", table, pKey)
	row, err2 := db.DataCall1(AppDBConn, query)
	if err2 != nil {
		log.Fatal(err2)
	}

	var keys []string
	for k, v := range row {	
		if v == nil {
			continue
		}
		keys = append(keys, k)
	}
	return calculateRowSize(AppDBConn, keys, table, pKey)
}

func GetTotalDataSize(stencilDBConn *sql.DB, AppDBConn *sql.DB, migrationID string) int64 {
	data := GetDataInLogicalSchemaOfMigration(stencilDBConn, AppDBConn, migrationID, "src")

	var totalDataSize int64
	for _, data1 := range data {
		totalDataSize += getEntireRowSize(AppDBConn, data1)
	}

	return totalDataSize
}

func getEntireUnmappedRowSize() {
	data := 
}

func GetLeftoverDataSize(stencilDBConn *sql.DB, AppDBConn *sql.DB, migrationID string) {
	getEntireUnmappedRowSize()
}
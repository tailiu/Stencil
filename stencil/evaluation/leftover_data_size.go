package evaluation

import (
	"stencil/db"
	"fmt"
	"database/sql"
	"log"
	"strings"
)


func getLeftoverDataInRowSize(AppDBConn *sql.DB, data1 map[string]interface{}) int64 {
	table, pKey := transformTableKeyToNormalType(data1)
	row := getLogicalRow(AppDBConn, table, pKey)

	var keys []string
	for k, v := range row {	
		if v == nil {
			continue
		}
		if strings.Contains(data1["src_cols"].(string), k) {
			continue
		}
		keys = append(keys, k)
	}
	log.Println(table)
	return calculateRowSize(AppDBConn, keys, table, pKey)
}

func getEntireRowInEvaluation(stencilDBConn *sql.DB, migrationID string) []map[string]interface{} {
	query := fmt.Sprintf("select * from evaluation where migration_id = '%s'", migrationID)
	
	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func getLeftoverDataInRowsSize(stencilDBConn *sql.DB, AppDBConn *sql.DB, migrationID string) {
	data := getEntireRowInEvaluation(stencilDBConn, migrationID)

	var leftoverDataInRowSize int64
	for _, data1 := range data {
		leftoverDataInRowSize += getLeftoverDataInRowSize(AppDBConn, data1)
	}
}

func getEntireUnmappedRowSize(stencilDBConn *sql.DB, migrationID string) {
	// query := fmt.Sprintf("select * from evaluation where migration_id = '%s' and ;", migrationID, ) 
}

func GetLeftoverDataSize(stencilDBConn *sql.DB, AppDBConn *sql.DB, migrationID string) {
	// getEntireUnmappedRowSize(stencilDBConn, migrationID)
	getLeftoverDataInRowsSize(stencilDBConn, AppDBConn, migrationID)
}
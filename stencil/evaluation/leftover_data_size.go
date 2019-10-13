package evaluation

import (
	"stencil/db"
	"fmt"
	"database/sql"
	"log"
	"strconv"
)

func getLeftoverDataInRowSize(AppDBConn *sql.DB, mCols map[string][]string, table string, pKey int, AppID string) int64 {
	row := getLogicalRow(AppDBConn, table, pKey)

	var keys []string
	for k, v := range row {
		if v == nil {
			continue
		}
		duplicate := false
		for _, col := range mCols[table] {
			if k == col {
				duplicate = true
				break
			}
		}
		if duplicate {
			continue
		} else {
			keys = append(keys, k)
		}
	}

	return calculateRowSize(AppDBConn, keys, table, pKey, AppID)
}

func getEntireRowInEvaluation(stencilDBConn *sql.DB, migrationID string) []map[string]interface{} {
	query := fmt.Sprintf("select * from evaluation where migration_id = '%s'", migrationID)
	
	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func getLeftoverDataInRowsSize(stencilDBConn *sql.DB, AppDBConn *sql.DB, AppID, migrationID string) int64 {
	data := getEntireRowInEvaluation(stencilDBConn, migrationID)
	mCols := getMigratedColsOfApp(stencilDBConn, AppID, migrationID)
	checkedRow := make(map[string]bool) 
	
	var leftoverDataInRowSize int64
	for _, data1 := range data {
		table, pKey := transformTableKeyToNormalType(data1)
		key := table + ":" + strconv.Itoa(pKey)
		if _, ok := checkedRow[key]; ok {
			continue
		} else {
			checkedRow[key] = true
			leftoverDataInRowSize += getLeftoverDataInRowSize(AppDBConn, mCols, table, pKey, AppID)
		}
	}
	return leftoverDataInRowSize
}

func getEntireUnmappedRowSize(stencilDBConn *sql.DB, AppDBConn *sql.DB, migrationID, AppID string) int64 {
	conditions := "dst_table = 'n/a'"
	data := getTableKeyInLogicalSchemaOfMigrationWithConditions(stencilDBConn, migrationID, "src", conditions)

	var entireUnmappedRowSize int64
	for _, data1 := range data {
		table, pKey := transformTableKeyToNormalType(data1)
		row := getLogicalRow(AppDBConn, table, pKey)

		var keys []string
		for k, v := range row {
			if v == nil {
				continue
			}
			keys = append(keys, k)
		}

		entireUnmappedRowSize += calculateRowSize(AppDBConn, keys, table, pKey, AppID)
	}

	return entireUnmappedRowSize
}

func GetLeftoverDataSize(stencilDBConn *sql.DB, AppDBConn *sql.DB, AppID, migrationID string) int64 {
	entireUnmappedRowSize := getEntireUnmappedRowSize(stencilDBConn, AppDBConn, migrationID, AppID)
	// log.Println(entireUnmappedRowSize)
	leftoverDataInRowSize := getLeftoverDataInRowsSize(stencilDBConn, AppDBConn, AppID, migrationID)
	// log.Println(leftoverDataInRowSize)
	
	return entireUnmappedRowSize + leftoverDataInRowSize
}
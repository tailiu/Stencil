package evaluation

import (
	"stencil/db"
	"fmt"
	"database/sql"
	"log"
	"strings"
	"strconv"
)

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

func getLeftoverDataInRowSize(AppDBConn *sql.DB, mCols map[string][]string, table string, pKey int) int64 {
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
			leftoverDataInRowSize += getLeftoverDataInRowSize(AppDBConn, mCols, table, pKey)
		}
	}
	return leftoverDataInRowSize
}

func getEntireUnmappedRowSize(stencilDBConn *sql.DB, migrationID string) {
	// query := fmt.Sprintf("select * from evaluation where migration_id = '%s' and ;", migrationID, ) 
}

func GetLeftoverDataSize(stencilDBConn *sql.DB, AppDBConn *sql.DB, AppID, migrationID string) {
	// getEntireUnmappedRowSize(stencilDBConn, migrationID)
	leftoverDataInRowSize := getLeftoverDataInRowsSize(stencilDBConn, AppDBConn, AppID, migrationID)
	console.log(leftoverDataInRowSize)
}
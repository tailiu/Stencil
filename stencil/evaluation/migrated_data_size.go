package evaluation

import (
	"database/sql"
	"strconv"
)

func getMigratedDataInRowSize(AppDBConn *sql.DB, data1 map[string]interface{}, mCols map[string][]string, table string, pKey int) int64 {
	row := getLogicalRow(AppDBConn, table, pKey)

	var keys []string
	for k, v := range row {	
		if v == nil {
			continue
		}
		for _, col := range mCols[table] {
			if k == col {
				keys = append(keys, k)
				break
			} 
		}
	}
	return calculateRowSize(AppDBConn, keys, table, pKey)
}

func GetMigratedDataSize(stencilDBConn *sql.DB, AppDBConn *sql.DB, AppID, migrationID string) int64 {
	conditions := "dst_table != 'n/a'"
	data := getTableKeyInLogicalSchemaOfMigrationWithConditions(stencilDBConn, migrationID, "src", conditions)
	mCols := getMigratedColsOfApp(stencilDBConn, AppID, migrationID)
	checkedRow := make(map[string]bool) 

	var migratedDataSize int64
	for _, data1 := range data {
		table, pKey := transformTableKeyToNormalType(data1)
		key := table + ":" + strconv.Itoa(pKey)
		if _, ok := checkedRow[key]; ok {
			continue
		} else {
			checkedRow[key] = true
			migratedDataSize += getMigratedDataInRowSize(AppDBConn, data1, mCols, table, pKey)
		}
	}

	return migratedDataSize
}
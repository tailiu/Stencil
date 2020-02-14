package SA2_db_optimization

import (
	"stencil/db"
	"database/sql"
	"strconv"
)

func existsInSlice(s []int, element int) bool {

	for _, v := range s {
		if element == v {
			return true
		}
	}

	return false

} 

func getAllTablesInDBs(dbConn *sql.DB) []map[string]string {

	query := `SELECT tablename FROM pg_catalog.pg_tables WHERE 
		schemaname != 'pg_catalog' AND schemaname != 'information_schema';`

	data := db.GetAllColsOfRows(dbConn, query)

	return data

}

func isSubPartitionTable(subPartitionTableIDs map[int]int, table string) bool {

	for _, subPartitionTableID := range subPartitionTableIDs {

		strID := strconv.Itoa(subPartitionTableID)

		subPartitionTable := "migration_table_" + strID

		if table == subPartitionTable {
			return true
		}

	}

	return false

}
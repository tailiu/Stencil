package SA2_db_populating

import (
	"stencil/db"
	"log"
	"fmt"
)

func DropPrimaryKeys(stencilDB string) {

	// db.STENCIL_DB = "stencil_exp_sa2_10k"
	db.STENCIL_DB = stencilDB

	dbConn := db.GetDBConn(db.STENCIL_DB)

	defer dbConn.Close()

	dropPrimaryKeysOfSA2TablesWithoutPartitions(dbConn)

}

func addPrimaryKeysToMigrationTable(stencilDB string) {

	db.STENCIL_DB = stencilDB

	dbConn := db.GetDBConn(db.STENCIL_DB)
	defer dbConn.Close()
	
	tables := getAllTablesInDB(dbConn)

	for _, t := range tables {
		
		table := t["tablename"]

		if isMigrationTable(table) {

			query := fmt.Sprintf(
				`ALTER TABLE %s ADD CONSTRAINT %s_pk 
				PRIMARY KEY (app_id, table_id, group_id, row_id, mark_as_delete);`,
				table, table,
			)

			log.Println(query)

			err := db.TxnExecute1(dbConn, query)
			if err != nil {
				log.Fatal(err)
			}

		}

	} 
}

func AddPrimaryKeys(stencilDB string) {

	addPrimaryKeysToMigrationTable(stencilDB)

	addPrimaryKeysToBaseSupTables(stencilDB)

}

func DeleteRowsByDuplicateColumnsInMigrationTable(stencilDB string) {

	db.STENCIL_DB = stencilDB

	uniqueCols := []string {
		"app_id", "table_id", "group_id", "row_id", "mark_as_delete",
	}
	
	dbConn := db.GetDBConn(db.STENCIL_DB)
	
	defer dbConn.Close()

	deleteRowsByDuplicateColumnsInMigrationTable(dbConn, uniqueCols)

}

func DeleteRowsByDuplicateColumnsInBaseSupTables(stencilDB string) {

	db.STENCIL_DB = stencilDB

	uniqueCols := []string {
		"pk",
	}
	
	dbConn := db.GetDBConn(db.STENCIL_DB)
	
	defer dbConn.Close()

	deleteRowsByDuplicateColumnsInBaseSupTables(dbConn, uniqueCols)

}

func DeleteDuplicateColumns(stencilDB string) {
	 
	DeleteRowsByDuplicateColumnsInMigrationTable(stencilDB)

	DeleteRowsByDuplicateColumnsInBaseSupTables(stencilDB)

}
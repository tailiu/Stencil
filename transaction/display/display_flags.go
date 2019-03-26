package display

import (
	"fmt"
	"log"
	"database/sql"
)

func CreateDisplayFlagTable(dbConn *sql.DB) {
	op := `CREATE TABLE display_flag (
			tableName string NOT NULL,
			id int NOT NULL,
			display bool default true, 
			INDEX id_index (id),
			INDEX table_index (tableName))`
	if _, err := dbConn.Exec(op); err != nil {
		log.Fatal(err)
	}
}

// id: primary key value
// table: tableName
func GenDisplayFlag(dbConn *sql.DB, id int, table string, display bool) {
	op := fmt.Sprintf("INSERT INTO display_flag (tableName, id, display) VALUES ('%s', %d, %t);",
						table, id, display)

	if _, err := dbConn.Exec(op); err != nil {
		log.Fatal(err)
	}
}

func UpdateDisplayFlag(dbConn *sql.DB, id int, table string, display bool) {
	op := fmt.Sprintf("UPDATE display_flag SET display = %t WHERE tableName = '%s' and id = %d",
						display, table, id)
	if _, err := dbConn.Exec(op); err != nil {
		log.Fatal(err)
	}
}
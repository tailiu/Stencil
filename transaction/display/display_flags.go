package display

import (
	"fmt"
	"errors"
	"log"
	"database/sql"
)

func CreateDisplayFlagsTable(dbConn *sql.DB) {
	op := `CREATE TABLE display_flags (
			tableName string NOT NULL,
			id int NOT NULL,
			display_flag bool default true, 
			INDEX id_index (id),
			INDEX table_index (tableName))`
	if _, err := dbConn.Exec(op); err != nil {
		log.Fatal(err)
	}
}

// id: primary key value
// table: tableName
func GenDisplayFlag(dbConn *sql.DB, id int, table string, display_flag bool) {
	op := fmt.Sprintf("INSERT INTO display_flags (tableName, id, display_flag) VALUES ('%s', %d, %t);",
						table, id, display_flag)

	if _, err := dbConn.Exec(op); err != nil {
		log.Fatal(err)
	}
}

func UpdateDisplayFlag(dbConn *sql.DB, id int, table string, display_flag bool) {
	op := fmt.Sprintf("UPDATE display_flags SET display_flag = %t WHERE tableName = '%s' and id = %d",
						display_flag, table, id)
	if _, err := dbConn.Exec(op); err != nil {
		log.Fatal(err)
	}
}

func CheckDisplayFlag(dbConn *sql.DB, id int, table string) (bool, error) {
	op := fmt.Sprintf("SELECT display_flag FROM display_flags WHERE tableName = '%s' and id = %d LIMIT 1",
						table, id)
	row, err := dbConn.Query(op)
	if err != nil {
		log.Fatal(err)
	}

	var display_flag bool
	find := false
	for row.Next() {
		if err := row.Scan(&display_flag); err != nil {
			log.Fatal(err)
		}
		find = true
	}

	if !find {
		return false, errors.New("Check Display Flag Error: Data Being Checked Does Not Exist!")
	} else {
		return display_flag, nil
	}
	
}
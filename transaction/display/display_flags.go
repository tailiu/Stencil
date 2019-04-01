package display

import (
	"fmt"
	"errors"
	"log"
	"database/sql"
)

func CreateDisplayFlagsTable(dbConn *sql.DB) {
	op := `CREATE TABLE display_flags (
			app string NOT NULL,
			tableName string NOT NULL,
			id int NOT NULL,
			display_flag bool default true, 
			INDEX app_index (app),
			INDEX id_index (id),
			INDEX table_index (tableName),
			INDEX display_flag_index (display_flag))`
	if _, err := dbConn.Exec(op); err != nil {
		log.Fatal(err)
	}
}

// app: application Name, id: primary key value, table: tableName
func GenDisplayFlag(dbConn *sql.DB, app, table string, id int, display_flag bool) error {
	op := fmt.Sprintf("INSERT INTO display_flags (app, tableName, id, display_flag) VALUES ('%s', '%s', %d, %t);",
						app, table, id, display_flag)

	if _, err := dbConn.Exec(op); err != nil {
		return err
	}
	return nil
}

func UpdateDisplayFlag(dbConn *sql.DB, app, table string, id int, display_flag bool) error {
	op := fmt.Sprintf("UPDATE display_flags SET display_flag = %t WHERE app = '%s' and tableName = '%s' and id = %d;",
						display_flag, app, table, id)
	if _, err := dbConn.Exec(op); err != nil {
		return err
	}
	return nil
}

func GetDisplayFlag(dbConn *sql.DB, app, table string, id int) (bool, error) {
	op := fmt.Sprintf("SELECT display_flag FROM display_flags WHERE app = '%s' and tableName = '%s' and id = %d LIMIT 1;",
						app, table, id)
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
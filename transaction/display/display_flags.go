package display

import (
	"fmt"
	"errors"
	"log"
	"database/sql"
	"time"
)

func CreateDisplayFlagsTable(dbConn *sql.DB) {
	op := `CREATE TABLE display_flags (
			app varchar NOT NULL,
			table_name varchar NOT NULL,
			id int NOT NULL,
			display_flag bool default true, 
			migration_id int,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL);`
	if _, err := dbConn.Exec(op); err != nil {
		log.Fatal(err)
	}
}

// app: application Name, id: primary key value, table: table_name
func GenDisplayFlag(dbConn *sql.DB, app, table string, id int, display_flag bool, migration_id ...int) error {
	var op string
	t := time.Now().Format(time.RFC3339)

	if len(migration_id) == 0 {
		op = fmt.Sprintf("INSERT INTO display_flags (app, table_name, id, display_flag, created_at, updated_at) VALUES ('%s', '%s', %d, %t, '%s', '%s');",
						app, table, id, display_flag, t, t)
	} else if len(migration_id) == 1 {
		op = fmt.Sprintf("INSERT INTO display_flags (app, table_name, id, display_flag, migration_id, created_at, updated_at) VALUES ('%s', '%s', %d, %t, %d, '%s', '%s');",
						app, table, id, display_flag, migration_id[0], t, t)
	} else {
		return errors.New("Argument Num Error: Please Input Only One Migration ID")
	}
	
	if _, err := dbConn.Exec(op); err != nil {
		return err
	}
	return nil
}

func UpdateDisplayFlag(dbConn *sql.DB, app, table string, id int, display_flag bool) error {
	t := time.Now().Format(time.RFC3339)
	op := fmt.Sprintf("UPDATE display_flags SET display_flag = %t, updated_at = '%s' WHERE app = '%s' and table_name = '%s' and id = %d;",
						display_flag, t, app, table, id)
	if _, err := dbConn.Exec(op); err != nil {
		return err
	}
	return nil
}

func GetDisplayFlag(dbConn *sql.DB, app, table string, id int) (bool, error) {
	op := fmt.Sprintf("SELECT display_flag FROM display_flags WHERE app = '%s' and table_name = '%s' and id = %d LIMIT 1;",
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
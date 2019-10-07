package app_display

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

func CreateDisplayFlagsTable(dbConn *sql.DB) {
	op := `CREATE TABLE display_flags (
			app_id int NOT NULL,
			table_id int NOT NULL,
			id int NOT NULL,
			migration_id int NOT NULL,
			display_flag bool NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL);`
	if _, err := dbConn.Exec(op); err != nil {
		log.Fatal(err)
	}
}

func GenDisplayFlag(dbConn *sql.DB, app_id, table_id, id, migration_id int) error {
	var op string
	op = fmt.Sprintf("INSERT INTO display_flags (app_id, table_id, id, migration_id, display_flag, created_at, updated_at) VALUES ( %d, %d, %d, %d, true, now(), now());",
		app_id, table_id, id, migration_id)
	if _, err := dbConn.Exec(op); err != nil {
		fmt.Println(op)
		return err
	}
	return nil
}

// func UpdateDisplayFlag(dbConn *sql.DB, app, table string, id int, display_flag bool) error {
// 	t := time.Now().Format(time.RFC3339)
// 	op := fmt.Sprintf("UPDATE display_flags SET display_flag = %t, updated_at = '%s' WHERE app = '%s' and table_name = '%s' and id = %d;",
// 		display_flag, t, app, table, id)
// 	fmt.Println("**************************************")
// 	fmt.Println(op)
// 	fmt.Println("**************************************")
// 	if _, err := dbConn.Exec(op); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func GetDisplayFlag(dbConn *sql.DB, app, table string, id int) (bool, error) {
// 	op := fmt.Sprintf("SELECT display_flag FROM display_flags WHERE app = '%s' and table_name = '%s' and id = %d LIMIT 1;",
// 		app, table, id)
// 	row, err := dbConn.Query(op)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var display_flag bool
// 	find := false
// 	for row.Next() {
// 		if err := row.Scan(&display_flag); err != nil {
// 			log.Fatal(err)
// 		}
// 		find = true
// 	}

// 	if !find {
// 		return false, errors.New("Check Display Flag Error: Data Being Checked Does Not Exist!")
// 	} else {
// 		return display_flag, nil
// 	}

// }

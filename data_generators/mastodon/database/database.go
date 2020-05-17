/*
 * Database Operations
 */

package database

import (
    "database/sql"
	"log"
	_ "github.com/lib/pq"
	"fmt"
)

func ConnectToDB(address string) *sql.DB {
	dbConn, err := sql.Open("postgres", address)
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
	}
	return dbConn
}

func BeginTx(dbConn *sql.DB) *sql.Tx {
	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	return tx
}

func Execute(tx *sql.Tx, queries []string) bool {
	haveErr := false
	for _, query := range queries {
		// fmt.Println(query)
		if _, err := tx.Exec(query); err != nil {
			fmt.Println(err)
			haveErr = true
			break
		}
	}
	return !haveErr
}

func CheckExists(tx *sql.Tx, query string) int {
	var exists int
	rows, err := tx.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&exists); err != nil {
			log.Fatal(err)
		}
	}
	return exists
}

func DataCall(dbConn *sql.DB, sql string, args ...interface{}) []map[string]string {

	var result []map[string]string

	rows, err := dbConn.Query(sql, args...)
	if err != nil {
		log.Fatal(err)
	}
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		data := make(map[string]string)
		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		for i, col := range cols {
			data[col] = columns[i]
		}
		result = append(result, data)
	}
	rows.Close()
	return result
}
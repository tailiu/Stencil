/*
 * Database Operations
 */

package database

import (
    "database/sql"
	"log"
    _ "github.com/lib/pq"
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

func Execute(tx *sql.Tx, queries []string) {
	for _, query := range queries {
		// fmt.Println(query)
		if _, err := tx.Exec(query); err != nil {
			log.Fatal(err)
		}	
	}
	tx.Commit()
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
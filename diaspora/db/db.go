package db

import (
	"database/sql"
	"diaspora/config"
	"fmt"
	"log"

	_ "github.com/lib/pq" // postgres driver
)

func GetDBConn(dbname string) *sql.DB {

	dbConnAddr := "postgresql://%s@%s:%s/%s?sslmode=disable"
	dbConn, err := sql.Open("postgres",
		fmt.Sprintf(dbConnAddr, config.DB_USER, config.DB_ADDR, config.DB_PORT, dbname))
	if err != nil {
		log.Println("Can't connect to DB:", dbname)
		log.Fatal(err)
	} else {
		log.Println("Connected to DB:", dbname)
	}
	return dbConn
}

func RunTxWQnArgsReturningId(tx *sql.Tx, query string, args ...interface{}) (int, error) {
	lastInsertId := -1
	err := tx.QueryRow(query, args...).Scan(&lastInsertId)
	if err != nil || lastInsertId == -1 {
		log.Println("# Can't insert!", err)
		tx.Rollback()
		log.Println(query, args)
		log.Println("Transaction rolled back!")
		return lastInsertId, err
	}
	return lastInsertId, err
}

func RunTxWQnArgs(tx *sql.Tx, query string, args ...interface{}) error {
	if _, err := tx.Exec(query, args...); err != nil {
		log.Println("# Can't execute!", err)
		tx.Rollback()
		log.Println(query, args)
		log.Println("Transaction rolled back!")
		return err
	}
	return nil
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
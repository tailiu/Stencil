package db

import (
	"database/sql"
	"diaspora/config"
	"fmt"
	"log"

	_ "github.com/lib/pq" // postgres driver
)

func GetDBConn(dbname string) *sql.DB {

	// dbConnAddr := "postgresql://%s@%s:%s/%s?sslmode=disable"
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", config.DB_ADDR, config.DB_PORT, config.DB_USER, config.DB_PASSWORD, dbname)

	dbConn, err := sql.Open("postgres", psqlInfo)
	// sql.Open("postgres",fmt.Sprintf(dbConnAddr, config.DB_USER, config.DB_ADDR, config.DB_PORT, dbname))
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

func DataCall1(dbConn *sql.DB, sql string, args ...interface{}) []map[string]string {

	var result []map[string]string

	rows, err := dbConn.Query(sql+" LIMIT 1", args...)

	if err != nil {
		log.Fatal(err)
	}
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	if rows.Next() {
		// log.Println("DataCall1: in rows.next()")
		data := make(map[string]string)
		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		// log.Println("DataCall1: scan pointers")
		rows.Scan(columnPointers...)
		// log.Println("DataCall1: pointers scanned")
		for i, col := range cols {
			data[col] = columns[i]
		}
		// log.Println("DataCall1: appending result")
		result = append(result, data)
	}
	// log.Println("DataCall1: rows.close")
	rows.Close()
	// log.Println("DataCall1: return")
	return result
}

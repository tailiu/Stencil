/*
 * DB Handler
 */

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq" // postgres driver
)

var dbConns map[string]*sql.DB

func GetDBConn(app string) *sql.DB {

	if dbConns == nil {
		dbConns = make(map[string]*sql.DB)
	}

	if _, ok := dbConns[app]; !ok {
		log.Println("Creating new db conn for:", app)
		dbConnAddr := "postgresql://root@10.224.45.154:26257/%s?sslmode=disable"
		dbConn, err := sql.Open("postgres", fmt.Sprintf(dbConnAddr, app))
		if err != nil {
			fmt.Println("error connecting to the db app:", app)
			log.Fatal(err)
		}
		dbConns[app] = dbConn
	}
	// log.Println("Returning dbconn for:", app)
	return dbConns[app]
}

func DataCall(app, sql string, args ...interface{}) []map[string]string {

	var result []map[string]string

	db := GetDBConn(app)

	rows, err := db.Query(sql, args...)
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

func DataCall1(app, sql string, args ...interface{}) (map[string]string, error) {

	data := make(map[string]string)

	db := GetDBConn(app)

	rows, err := db.Query(sql+" LIMIT 1", args...)
	if err != nil {
		log.Fatal(err)
	}
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	if rows.Next() {
		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		for i, col := range cols {
			data[col] = columns[i]
		}

		return data, nil
	}

	return data, errors.New("no result found for sql: " + sql)
}

func GetPK(app, table string) []string {

	var result []string

	db := GetDBConn(app)

	sql := "SHOW CONSTRAINTS FROM " + table

	rows, err := db.Query(sql)
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

		if data["Type"] == "PRIMARY KEY" {
			result = strings.Split(data["Column(s)"], ",")
			break
		}
	}
	rows.Close()
	return result
}

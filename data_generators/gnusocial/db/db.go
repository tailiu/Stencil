package db

import (
	"database/sql"
	"errors"
	"fmt"
	"gnusocial/config"
	"log"
	"math/rand"
	"strings"

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

func GetNewRowIDForTable(dbConn *sql.DB, table string) string {

	var rowid int32
	for {
		rowid = rand.Int31n(2147483647)
		q := fmt.Sprintf("SELECT id FROM \"%s\" WHERE id = %d", table, rowid)
		if v, err := DataCall1(dbConn, q); err != nil {
			fmt.Println(q)
			log.Println("@db.GetNewRowIDForTable: ", table)
			log.Fatal(err)
		} else if v == nil {
			break
		}
	}
	return fmt.Sprint(rowid)
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
		tx.Rollback()
		if !strings.Contains(err.Error(), "row_desc_pk") {
			log.Println("# Can't execute!", err)
			log.Println(query, args)
			log.Println("Transaction rolled back!")
		} else {
			return errors.New("# ROW_DESC ROWID EXISTS!")
		}
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

func DataCall1(db *sql.DB, SQL string, args ...interface{}) (map[string]interface{}, error) {

	// db := GetDBConn(app)
	// log.Println(SQL, args)
	if rows, err := db.Query(SQL+" LIMIT 1", args...); err != nil {
		// log.Println(SQL, args)
		// log.Println("## DB ERROR: ", err)
		// log.Fatal("check datacall1 in stencil.db")
		return nil, err
	} else {
		defer rows.Close()

		if colNames, err := rows.Columns(); err != nil {
			return nil, err
		} else {
			if rows.Next() {
				var data = make(map[string]interface{})
				cols := make([]interface{}, len(colNames))
				colPtrs := make([]interface{}, len(colNames))
				for i := 0; i < len(colNames); i++ {
					colPtrs[i] = &cols[i]
				}
				// for rows.Next() {
				err = rows.Scan(colPtrs...)
				if err != nil {
					log.Fatal(err)
				}
				for i, col := range cols {
					data[colNames[i]] = col
				}
				// Do something with the map
				// for key, val := range data {
				// 	fmt.Println("Key:", key, "Value Type:", reflect.TypeOf(val), fmt.Sprint(val))
				// }
				return data, nil
			} else {
				return nil, nil
			}
		}
	}
}

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
	"transaction/config"

	_ "github.com/lib/pq" // postgres driver
	"github.com/tj/go-pg-escape"
)

func GetDBConn(app string) *sql.DB {

	dbConnAddr := "postgresql://root@10.224.45.158:26257/%s?sslmode=disable"

	dbConn, err := sql.Open("postgres", fmt.Sprintf(dbConnAddr, app))
	if err != nil {
		fmt.Println("error connecting to the db app:", app)
		log.Fatal(err)
	}
	return dbConn
}

func MoveData(srcApp, tgtApp string, sql config.DataQuery, mappings config.Mapping, uid int) error {

	if appMapping, ok := mappings[tgtApp]; ok {

		if tableMapping, ok := appMapping[strings.ToLower(sql.Table)]; ok {

			srcDB := GetDBConn(srcApp)
			tgtDB := GetDBConn(tgtApp)

			rows, err := srcDB.Query(sql.SQL, uid)
			if err != nil {
				log.Fatal(err)
			}
			cols, err := rows.Columns()
			if err != nil {
				log.Fatal(err)
			}

			data := make(map[string]string)
			columns := make([]string, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i := range columns {
				columnPointers[i] = &columns[i]
			}

			tx, err := tgtDB.Begin()

			for rows.Next() {

				rows.Scan(columnPointers...)

				for i, col := range cols {
					data[col] = columns[i]
				}

				for tgtTable, tgtMap := range tableMapping {

					var cols, vals string
					for col1, col2 := range tgtMap {
						cols += col1 + ","
						vals += escape.Literal(data[col2]) + ","
						// vals += fmt.Sprintf("\"%s\",", data[col2])
					}
					cols = strings.TrimSuffix(cols, ",")
					vals = strings.TrimSuffix(vals, ",")
					insql := escape.Escape("INSERT INTO %s (%s) VALUES (%s)", tgtTable, cols, vals)
					fmt.Println(insql)

					stmt, err := tx.Prepare(insql)
					if err != nil {
						tx.Rollback()
						return err
					}
					defer stmt.Close()

					if _, err := stmt.Exec(insql); err != nil {
						tx.Rollback() // return an error too, we may want to wrap them
						return err
					}

					// if _, err = tgtDB.Exec(insql); err != nil {
					// 	fmt.Println(">>>>>>>>>>> Can't insert!")
					// 	panic(err)
					// }
				}
			}
			rows.Close()
			return tx.Commit()
		}
		return errors.New("mapping doesn't exist for table:" + sql.Table)
	}
	return errors.New("mapping doesn't exist for app:" + tgtApp)
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

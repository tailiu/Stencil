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

var dbConns map[string]*sql.DB

func GetDBConn(app string) *sql.DB {

	if dbConns == nil {
		dbConns = make(map[string]*sql.DB)
	}

	if _, ok := dbConns[app]; !ok {
		log.Println("Creating new db conn for:", app)
		dbConnAddr := "postgresql://root@10.224.45.158:26257/%s?sslmode=disable"
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

func MoveData(srcApp, tgtApp string, sql config.DataQuery, mappings config.Mapping, uid int) error {

	if appMapping, ok := mappings[tgtApp]; ok {

		if tableMapping, ok := appMapping[strings.ToLower(sql.Table)]; ok {

			srcDB := GetDBConn(srcApp)
			tgtDB := GetDBConn(tgtApp)

			for {

				row, err := DataCall1(srcApp, sql.SQL, uid)
				if err == nil {

					ttx, err := tgtDB.Begin()
					if err != nil {
						log.Println("ERROR! TARGET TRANSACTION CAN'T BEGIN")
						return err
					}

					stx, err := srcDB.Begin()
					if err != nil {
						log.Println("ERROR! SOURCE TRANSACTION CAN'T BEGIN")
						return err
					}

					defer ttx.Rollback()
					defer stx.Rollback()

					ucond := ""

					for col, val := range row {
						if !strings.EqualFold(col, "mark_delete") && val != "" {
							ucond += fmt.Sprintf(" %s = %s AND", col, escape.Literal(val))
						}
						// fmt.Println("col", col, "data", columns[i])
					}

					ucond = strings.TrimSuffix(ucond, "AND")
					usql := fmt.Sprintf("UPDATE %s SET mark_delete = 'true' WHERE %s", sql.Table, ucond)

					if _, err = stx.Exec(usql); err != nil {
						fmt.Println(">>>>>>>>>>> Can't update!")
						return err
					} else {
						fmt.Println("Updated!")
					}

					for tgtTable, tgtMap := range tableMapping {

						var cols, vals string
						for scol, tcol := range tgtMap {
							cols += tcol + ","
							vals += escape.Literal(row[scol]) + ","
						}
						cols = strings.TrimSuffix(cols, ",")
						vals = strings.TrimSuffix(vals, ",")
						insql := escape.Escape("INSERT INTO %s (%s) VALUES (%s)", tgtTable, cols, vals)

						if _, err = ttx.Exec(insql); err != nil {
							log.Println("# Can't insert!")
							return err
						} else {
							fmt.Println("Inserted!")
						}
					}

					stx.Commit()
					ttx.Commit()
				} else if err != nil {
					log.Println("# No more rows!")
					break
				}

			}
			return nil
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

func DataCall1(app, sql string, args ...interface{}) (map[string]string, error) {

	data := make(map[string]string)

	db := GetDBConn(app)

	rows, err := db.Query(sql, args...)
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

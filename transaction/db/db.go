/*
 * DB Handler
 */

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
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
		dbConnAddr := "postgresql://root@10.230.12.75:26257/%s?sslmode=disable"
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

func GetColumnsForTable(app, table string) ([]string, string) {
	var resultList []string
	resultStr := ""

	db := GetDBConn(app)

	rows, err := db.Query("SHOW COLUMNS FROM " + table)
	if err != nil {
		log.Fatal(err)
	}
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		for i, col := range cols {

			if strings.EqualFold(col, "column_name") {
				resultList = append(resultList, columns[i])
				// resultStr += fmt.Sprintf("IFNULL(%s.%s, 'NULL') AS \"%s.%s\",", table, columns[i], table, columns[i])
				resultStr += table + "." + columns[i] + " AS \"" + table + "." + columns[i] + "\","
			}

		}

	}
	rows.Close()
	return resultList, strings.Trim(resultStr, ",")
}

func GetRow(rows *sql.Rows) map[string]interface{} {
	var myMap = make(map[string]interface{})

	colNames, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
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
		myMap[colNames[i]] = col
	}
	// Do something with the map
	for key, val := range myMap {
		fmt.Println("Key:", key, "Value Type:", reflect.TypeOf(val))
	}
	// }
	return myMap
}

func DataCall(app, SQL string, args ...interface{}) []map[string]interface{} {

	var result []map[string]interface{}

	db := GetDBConn(app)

	if rows, err := db.Query(SQL, args...); err != nil {
		log.Println(SQL, args)
		log.Fatal(err)
	} else {

		if colNames, err := rows.Columns(); err != nil {
			log.Fatal(err)
		} else {

			for rows.Next() {
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
				result = append(result, data)
			}
			rows.Close()
		}
	}
	return result
}

func DataCall1(app, SQL string, args ...interface{}) map[string]interface{} {

	db := GetDBConn(app)

	if rows, err := db.Query(SQL+" LIMIT 1", args...); err != nil {
		log.Println(SQL, args)
		log.Fatal("## DB ERROR: ", err)
	} else {

		if colNames, err := rows.Columns(); err != nil {
			log.Fatal(err)
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
				return data
			}
			rows.Close()
		}
	}
	return make(map[string]interface{})
}

func _DataCall1(app, sql string, args ...interface{}) (map[string]string, error) {

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

	// defer rows.Close()

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

		rows.Close()
		return data, nil
	}

	rows.Close()
	return data, errors.New("no result found for sql: " + sql)
}

// func GetAppId(app_name string) (string, error) {
// 	sql := "SELECT row_id from apps WHERE app_name = $1"

// 	if result, err := DataCall1("stencil", sql, app_name); err == nil {
// 		return result["row_id"], nil
// 	}
// 	return "-1", errors.New("App Not Found: " + app_name)
// }

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

func getAllColsOfRows(dbConn *sql.DB, query string) []map[string]string {
	rows, err := dbConn.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	var allRows []map[string]string

	for rows.Next() {
		data := make(map[string]string)

		cols, err := rows.Columns()
		if err != nil {
			log.Fatal(err)
		}

		columns := make([]sql.NullString, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		err1 := rows.Scan(columnPointers...)
		if err1 != nil {
			log.Fatal(err1)
		}
		for i, colName := range cols {
			if columns[i].Valid {
				data[colName] = columns[i].String
			} else {
				data[colName] = "NULL"
			}
		}
		allRows = append(allRows, data)
	}
	// fmt.Println(query)
	// fmt.Println(data)
	return allRows
}

// NOTE: We assume that primary key is only one string!!!
func GetPrimaryKeyOfTable(dbConn *sql.DB, table string) (string, error) {
	query := fmt.Sprintf("SHOW CONSTRAINTS FROM %s;", table)
	constraints := getAllColsOfRows(dbConn, query)

	for _, constraint := range constraints {
		if constraint["constraint_type"] == "PRIMARY KEY" {
			details := constraint["details"]
			s1 := strings.Split(details, "(")[1]
			s2 := strings.Split(s1, ")")[0]
			s3 := strings.Split(s2, " ")[0]
			return s3, nil
		}
	}

	return "", fmt.Errorf("Get Primary Key Error: No Primary Key Found For Table %s", table)
}

func GetOneRowBasedOnHint(dbConn *sql.DB, app, depDataValue, depDataValueType, depDataKey, depDataTable string) (map[string]string, error) {
	var query string
	switch valueType := depDataValueType; valueType {
	case "int":
		value, err := strconv.Atoi(depDataValue)
		if err != nil {
			log.Fatal(err)
		}
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s = %d LIMIT 1;", depDataTable, depDataKey, value)
	default:
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s = '%s' LIMIT 1;", depDataTable, depDataKey, depDataValue)
	}

	data := getAllColsOfRows(dbConn, query)
	if len(data) == 0 {
		return nil, errors.New("Check Remaining Data Exists Error: Original Data Not Exists")
	} else {
		return data[0], nil
	}
}

func GetOneRowBasedOnDependency(dbConn *sql.DB, app string, val int, dep string) (map[string]string, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %d LIMIT 1;", strings.Split(dep, ".")[0], strings.Split(dep, ".")[1], val)
	// fmt.Println(query)
	data := getAllColsOfRows(dbConn, query)
	if len(data) == 0 {
		return nil, errors.New("Check Remaining Data Exists Error: Data Not Exists")
	} else {
		return data[0], nil
	}
}

// NOTE: This should be changed to get one row RANDOMLY!!
func GetOneRowInParentNodeRandomly(dbConn *sql.DB, depDataValue, depDataValueType, depDataKey, depDataTable string, conditions []string) (map[string]string, string, error) {
	query := fmt.Sprintf("SELECT %s.* FROM ", "t"+strconv.Itoa(len(conditions)))
	from := ""
	table := ""
	for i, condition := range conditions {
		table1 := strings.Split(condition, ":")[0]
		table2 := strings.Split(condition, ":")[1]
		t1 := strings.Split(table1, ".")[0]
		a1 := strings.Split(table1, ".")[1]
		t2 := strings.Split(table2, ".")[0]
		a2 := strings.Split(table2, ".")[1]
		seq1 := "t" + strconv.Itoa(i)
		seq2 := "t" + strconv.Itoa(i+1)
		if i == 0 {
			from += fmt.Sprintf("%s %s JOIN %s %s ON %s.%s = %s.%s ",
				t1, seq1, t2, seq2, seq1, a1, seq2, a2)
		} else {
			from += fmt.Sprintf("JOIN %s %s on %s.%s = %s.%s ",
				t2, seq2, seq1, a1, seq2, a2)
		}
		if i == len(conditions)-1 {
			where := ""
			if depDataValueType == "int" {
				val, err := strconv.Atoi(depDataValue)
				if err != nil {
					log.Fatal(err)
				}
				where = fmt.Sprintf("WHERE %s.%s = %d LIMIT 1;", "t0", depDataKey, val)
			} else {
				where = fmt.Sprintf("WHERE %s.%s = '%s' LIMIT 1;", "t0", depDataKey, depDataValue)
			}
			table = t2
			query += from + where
		}
	}
	fmt.Println(query)

	data := getAllColsOfRows(dbConn, query)
	if len(data) == 0 {
		return nil, "", errors.New("Error In Get Data: Fail To Get One Data This Data Depends On")
	} else {
		return data[0], table, nil
	}
}

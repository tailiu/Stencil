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
	"strings"
	"math/rand"
	_ "github.com/lib/pq" // postgres driver
)

// func GetDBConn(app string) *sql.DB {

// 	if dbConns == nil {
// 		dbConns = make(map[string]*sql.DB)
// 	}

// 	if _, ok := dbConns[app]; !ok {
// 		log.Println("Creating new db conn for:", app)
// 		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
// 			"password=%s dbname=%s sslmode=disable", DB_ADDR, DB_PORT, DB_USER, DB_PASSWORD, app)
// 		// dbConnAddr := "postgresql://root@10.230.12.75:26257/%s?sslmode=disable"
// 		// fmt.Println(psqlInfo)
// 		dbConn, err := sql.Open("postgres", psqlInfo)
// 		if err != nil {
// 			fmt.Println("error connecting to the db app:", app)
// 			log.Fatal(err)
// 		}
// 		dbConns[app] = dbConn
// 	}
// 	// log.Println("Returning dbconn for:", app)
// 	return dbConns[app]
// }

func GetDBConn(app string) *sql.DB {
	// log.Println("Creating new db conn for:", app)
	var host string
	if strings.Contains(app, "old") {
		host = DB_ADDR_old
		app = strings.Split(app, "old_")[1]
	} else {
		host = DB_ADDR
	}
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", host, DB_PORT, DB_USER, DB_PASSWORD, app)
	// dbConnAddr := "postgresql://root@10.230.12.75:26257/%s?sslmode=disable"
	// fmt.Println(psqlInfo)
	dbConn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("error connecting to the db app:", app)
		log.Fatal(err)
	}
	return dbConn
}

func GetDBConn2(app string) *sql.DB {
	// log.Println("Creating new db conn for:", app)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", DB_ADDR_old, DB_PORT, DB_USER, DB_PASSWORD, app)
	// dbConnAddr := "postgresql://root@10.230.12.75:26257/%s?sslmode=disable"
	// fmt.Println(psqlInfo)
	dbConn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("error connecting to the db app:", app)
		log.Fatal(err)
	}
	return dbConn
}

func CloseDBConn(app string) {
	if _, ok := dbConns[app]; ok {
		dbConns[app].Close()
		delete(dbConns, app)
		log.Println("DBConn closed for:", app)
	}
}

func Delete(db *sql.DB, SQL string, args ...interface{}) error {

	if _, err := db.Query(SQL, args...); err != nil {
		log.Println(SQL, args)
		log.Fatal("## DB ERROR: ", err)
		return err
	}
	return nil
}

func Insert(dbConn *sql.DB, query string, args ...interface{}) (int, error) {

	lastInsertId := -1
	err := dbConn.QueryRow(query+" RETURNING id; ", args...).Scan(&lastInsertId)
	if err != nil || lastInsertId == -1 {
		return lastInsertId, err
	}
	return lastInsertId, err
}

func UpdateTx(tx *sql.Tx, query string, args ...interface{}) error {

	_, err := tx.Exec(query, args...)
	return err
}

func InsertRowIntoAppDB(tx *sql.Tx, table, cols, placeholders string, args ...interface{}) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id;", table, cols, placeholders)
	lastInsertId := -1
	err := tx.QueryRow(query, args...).Scan(&lastInsertId)
	if err != nil || lastInsertId == -1 {
		return lastInsertId, err
	}
	return lastInsertId, err
}

func DeleteRowFromAppDB(tx *sql.Tx, table, id string) error {
	query := fmt.Sprintf("UPDATE %s SET mark_as_delete = $1 WHERE id = $2", table)
	if _, err := tx.Exec(query, true, id); err != nil {
		log.Println(query, "true", id)
		log.Fatal("## DB ERROR: ", err)
		return err
	}
	return nil
}

func NewBag(tx *sql.Tx, rowid, user_id, tagName string, migration_id int) error {
	query := "INSERT INTO data_bags (rowid, user_id, tag, migration_id) VALUES ($1, $2, $3, $4)"
	_, err := tx.Exec(query, rowid, user_id, tagName, migration_id)
	return err
}

func NewRow(tx *sql.Tx, rowid, app_id, mflag string, copy_on_write bool) error {
	query := "INSERT INTO row_desc (rowid, app_id, copy_on_write, mflag) VALUES ($1, $2, $3, $4)"
	_, err := tx.Exec(query, rowid, app_id, copy_on_write, mflag)
	return err
}

func GetUserBags(dbConn *sql.DB, user_id, app_id string) ([]map[string]interface{}, error) {
	query := "SELECT string_agg(data_bags.rowid::varchar, ',') as rowids, data_bags.tag as tag FROM data_bags JOIN row_desc ON data_bags.rowid = row_desc.rowid WHERE data_bags.user_id = $1 AND row_desc.mflag = 1 AND row_desc.app_id = $2 group by data_bags.tag"
	return DataCall(dbConn, query, user_id, app_id)
}

func DeleteBag(dbConn *sql.DB, bag_id string) error {
	query := "DELETE FROM data_bags WHERE pk = $1"
	return Delete(dbConn, query, bag_id)
}

func DeleteBagsByRowIDS(dbConn *sql.DB, rowids string) error {
	query := "DELETE FROM data_bags WHERE rowid IN ($1)"
	_, err := dbConn.Exec(query, rowids)
	return err
}

func SetAppID(tx *sql.Tx, pk, app_id string) error {

	q := "UPDATE row_desc SET app_id = $1 WHERE rowid = $2"
	_, err := tx.Exec(q, app_id, pk)
	return err
}

func SetMFlag(tx *sql.Tx, pk, flag string) error {

	q := "UPDATE row_desc SET mflag = $1 WHERE rowid IN ($2)"
	_, err := tx.Exec(q, flag, pk)
	return err
}

func MUpdate(tx *sql.Tx, pk, flag, app_id string) error {

	q := "UPDATE row_desc SET mflag = $1, app_id = $2 WHERE rowid IN ($3)"
	_, err := tx.Exec(q, flag, app_id, pk)
	return err
}

func BUpdate(dbConn *sql.DB, pk, flag, app_id string) error {

	q := "UPDATE row_desc SET mflag = $1, app_id = $2 WHERE rowid IN ($3)"
	_, err := dbConn.Exec(q, flag, app_id, pk)
	return err
}

func GetUnmigratedUsers() ([]map[string]interface{}, error) {
	dbConn := GetDBConn("stencil")
	sql := "SELECT user_id FROM user_table WHERE user_id NOT IN (SELECT DISTINCT user_id FROM migration_registration) ORDER BY user_id ASC"
	return DataCall(dbConn, sql)
}

func RemoveUserFromApp(uid, app_id string, dbConn *sql.DB) bool {
	sql := "DELETE FROM user_table WHERE user_id = $1 AND app_id = $2"
	if err := Delete(dbConn, sql, uid, app_id); err == nil {
		return true
	}
	return false
}

func CheckUserInApp(uid, app_id string, dbConn *sql.DB) bool {
	sql := "SELECT user_id FROM user_table WHERE user_id = $1 AND app_id = $2"
	if res, err := DataCall1(dbConn, sql, uid, app_id); err == nil {
		if len(res) > 0 {
			return true
		}
	} else {
		log.Fatal(err)
	}
	return false
}

func AddUserToApp(uid, app_id string, dbConn *sql.DB) bool {
	query := "INSERT INTO user_table (user_id, app_id) VALUES ($1, $2)"
	dbConn.Exec(query, uid, app_id)
	return true
}

func AddOwnedData(uid, row_id string, dbConn *sql.DB) bool {
	query := "INSERT INTO owned_data (user_id, row_id) VALUES ($1, $2)"
	if _, err := dbConn.Exec(query, uid, row_id); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return true
		}
		log.Fatal("Error in db", err)
	}
	return true
}

func TruncateOwnedData(dbConn *sql.DB) {
	query := "TRUNCATE TABLE owned_data"
	if _, err := dbConn.Exec(query); err != nil {
		log.Fatal("Error in db", err)
	}
}

func DeleteExistingMigrationRegistrations(uid, src_app, dst_app string, dbConn *sql.DB) bool {
	query := "DELETE FROM migration_registration WHERE user_id = $1 AND src_app = $2 AND dst_app = $3"
	if _, err := dbConn.Exec(query, uid, src_app, dst_app); err != nil {
		log.Fatal("DELETE Error in DeleteExistingMigrationRegistrations", err)
		return false
	}
	return true
}

func CheckMigrationRegistration(uid, src_app, dst_app string, dbConn *sql.DB) bool {
	sql := "SELECT migration_id FROM migration_registration WHERE user_id = $1 AND src_app = $2 AND dst_app = $3"
	if res, err := DataCall1(dbConn, sql, uid, src_app, dst_app); err == nil {
		fmt.Println(res)
		if len(res) > 0 {
			return true
		}
	} else {
		log.Fatal(err)
	}
	return false
}

func RegisterMigration(uid, src_app, dst_app, mtype string, migrationID, number_of_threads int, dbConn *sql.DB, logical bool) bool {
	query := "INSERT INTO migration_registration (migration_id, user_id, src_app, dst_app, migration_type, number_of_threads, is_logical) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	if _, err := dbConn.Exec(query, migrationID, uid, src_app, dst_app, mtype, number_of_threads, logical); err != nil {
		log.Fatal("Insert Error in RegisterMigration", err)
		return false
	}
	return true
}

func GetColumnsForTable(db *sql.DB, table string) ([]string, string) {
	var resultList []string
	resultStr := ""

	// db := GetDBConn(app)

	query := "select column_name, data_type, is_nullable, column_default, generation_expression from INFORMATION_SCHEMA.COLUMNS where table_name = $1;"
	rows, err := db.Query(query, table)
	// rows, err := db.Query("SHOW COLUMNS FROM " + table)
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

func DataCall(db *sql.DB, SQL string, args ...interface{}) ([]map[string]interface{}, error) {

	// db := GetDBConn(app)
	// log.Println(SQL, args)
	if rows, err := db.Query(SQL, args...); err != nil {
		log.Println(SQL, args)
		return nil, err
	} else {
		defer rows.Close()

		if colNames, err := rows.Columns(); err != nil {
			return nil, err
		} else {
			var result []map[string]interface{}

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
			return result, nil
		}
	}
}

func DataCall1(db *sql.DB, SQL string, args ...interface{}) (map[string]interface{}, error) {

	// db := GetDBConn(app)
	// log.Println(SQL, args)
	if rows, err := db.Query(SQL+" LIMIT 1", args...); err != nil {
		log.Println(SQL, args)
		log.Println("## DB ERROR: ", err)
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

func GetNewRowID(dbConn *sql.DB) int32 {
	
	var rowid int32
	for{
		rowid = rand.Int31n(2147483647)
		q := "SELECT rowid FROM row_desc WHERE rowid = $1"
		if v, err := DataCall1(dbConn, q, rowid); err == nil && v == nil {
			break
		}
	}
	return rowid
}

func GetAllColsOfRows(dbConn *sql.DB, query string) []map[string]string {
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
	query := fmt.Sprintf("SELECT c.column_name FROM information_schema.key_column_usage AS c LEFT JOIN information_schema.table_constraints AS t ON t.constraint_name = c.constraint_name WHERE t.table_name = '%s' AND t.constraint_type = 'PRIMARY KEY';", table)
	primaryKey := GetAllColsOfRows(dbConn, query)

	if len(primaryKey) == 0 {
		return "", fmt.Errorf("Get Primary Key Error: No Primary Key Found For Table %s", table)
	}

	if pk, ok := primaryKey[0]["column_name"]; ok {
		return pk, nil
	} else {
		return "", fmt.Errorf("Get Primary Key Error: No Primary Key Found For Table %s", table)
	}
}

func GetTablesOfDB(dbConn *sql.DB, app string) []string {
	query := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE';")
	tablesMapArr := GetAllColsOfRows(dbConn, query)

	var tables []string
	for _, element := range tablesMapArr {
		for _, table := range element {
			tables = append(tables, table)
		}
	}

	return tables
}

func TxnExecute(dbConn *sql.DB, queries []string) error {
	tx, err := dbConn.Begin()
	if err != nil {
		return err
	}

	for _, query := range queries {
		// fmt.Println(query)
		if _, err := tx.Exec(query); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil

}

func SaveForEvaluation(dbConn *sql.DB, srcApp, dstApp, srcTable, dstTable, srcID, dstID, srcCol, dstCol, migrationID string) error {
	query := "INSERT INTO evaluation (src_app, dst_app, src_table, dst_table, src_id, dst_id, src_cols, dst_cols, migration_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	_, err := Insert(dbConn, query, srcApp, dstApp, srcTable, dstTable, srcID, dstID, srcCol, dstCol, migrationID)
	return err
}

func SaveForLEvaluation(dbConn *sql.DB, srcApp, dstApp, srcTable, dstTable, srcID, dstID, srcCol, dstCol, migrationID string) error {
	query := "INSERT INTO evaluation (src_app, dst_app, src_table, dst_table, src_id, dst_id, src_cols, dst_cols, migration_id, added_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now())"
	_, err := Insert(dbConn, query, srcApp, dstApp, srcTable, dstTable, srcID, dstID, srcCol, dstCol, migrationID)
	return err
}

func UpdateLEvaluation(dbCOnn *sql.DB, srcTable, srcID string, migrationID int) error {
	query := "UPDATE evaluation SET deleted_at = now() WHERE migration_id = $1 AND src_table = $2 AND src_ID = $3"
	_, err := dbCOnn.Exec(query, migrationID, srcTable, srcID)
	return err
}

func LogError(dbConn *sql.DB, dbquery, args, migration_id, dst_app, qerr string) error {

	query := "INSERT INTO error_log (query, args, migration_id, dst_app, error) VALUES ($1, $2, $3, $4, $5)"
	_, err := Insert(dbConn, query, dbquery, args, migration_id, dst_app, qerr)
	return err
}

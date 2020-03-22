/*
 * DB Handler
 */

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"reflect"
	"strings"

	"github.com/gookit/color"
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

func GetDBConn(app string, isBlade ...bool) *sql.DB {

	var dbName string

	switch app {
	case "diaspora":
		{
			dbName = DIASPORA_DB
		}
	case "mastodon":
		{
			dbName = MASTODON_DB
		}
	case "twitter":
		{
			dbName = TWITTER_DB
		}
	case "gnusocial":
		{
			dbName = GNUSOCIAL_DB
		}
	case "stencil":
		{
			dbName = STENCIL_DB
		}
	default:
		{
			dbName = app
		}
	}

	dbAddr := DB_ADDR

	if len(isBlade) > 0 {
		if isBlade[0] {
			dbAddr = DB_ADDR_old
		}
	}

	color.Info.Println(fmt.Sprintf("Connecting to DB \"%s\" @ [%s] for App {%s} ...", dbName, dbAddr, app))
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=600", dbAddr, DB_PORT, DB_USER, DB_PASSWORD, dbName)

	dbConn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		color.Error.Println("error connecting to the db app:", app)
		// fmt.Println("error connecting to the db app:", app)
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

func GetRowCount(dbConn *sql.DB, table string) (int64, error) {
	sql := fmt.Sprintf("SELECT COUNT(*) as total_rows from %s", table)
	if result, err := DataCall1(dbConn, sql); err == nil {
		if val, ok := result["total_rows"]; ok {
			return val.(int64), nil
		}
		log.Fatal("db.GetRowCount: Can't find total row count for table: ", table)
		return -1, err
	} else {
		log.Fatal("db.GetRowCount:  ", err)
		return -1, err
	}
}

func GetAppIDByAppName(dbConn *sql.DB, app string) string {
	sql := fmt.Sprintf("SELECT pk from apps WHERE app_name = '%s'", app)
	if result, err := DataCall1(dbConn, sql); err == nil {
		if val, ok := result["pk"]; ok {
			return fmt.Sprint(val)
		}
		log.Fatal("db.GetAppIDByAppName: Can't find app id for app ", app)
	} else {
		log.Fatal("db.GetAppIDByAppName:  ", err)
	}
	return ""
}

func GetAppNameByAppID(dbConn *sql.DB, appID string) (string, error) {
	sql := "SELECT app_name FROM apps WHERE pk = $1"
	if result, err := DataCall1(dbConn, sql, appID); err == nil {
		if val, ok := result["app_name"]; ok {
			return fmt.Sprint(val), nil
		}
		log.Fatal("db.GetAppIDByAppName: Can't find app_name for id ", appID)
	} else {
		log.Fatal("db.GetAppIDByAppName:  ", err)
	}
	return "", fmt.Errorf("can't find app_name by id %s", appID)
}

func Delete(db *sql.DB, SQL string, args ...interface{}) error {

	if _, err := db.Query(SQL, args...); err != nil {
		log.Println(SQL, args)
		// log.Fatal("## DB ERROR: ", err)
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
	query := fmt.Sprintf("INSERT INTO \"%s\" (%s, display_flag) VALUES (%s, true) RETURNING id;", table, cols, placeholders)
	lastInsertId := -1
	err := tx.QueryRow(query, args...).Scan(&lastInsertId)
	if err != nil || lastInsertId == -1 {
		return lastInsertId, err
	}
	return lastInsertId, err
}

func InsertIntoIdentityTable(tx *sql.Tx, srcApp, dstApp, srcTable, dstTable, srcID, dstID, migrationID interface{}) error {
	query := "INSERT INTO identity_table (from_app, from_member, from_id, to_app, to_member, to_id, migration_id) VALUES ($1, $2, $3, $4, $5, $6, $7);"
	_, err := tx.Exec(query, srcApp, srcTable, srcID, dstApp, dstTable, dstID, migrationID)
	return err
}

func _DropAndRecreateDB(dbConn *sql.DB, dbname string) error {
	q := fmt.Sprintf("SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname in ('%s', 'diaspora_1000000') AND pid <> pg_backend_pid();", dbname)

	if _, err := dbConn.Exec(q); err != nil {
		fmt.Println(q)
		log.Fatal(err)
	}

	var cmd *exec.Cmd

	q = fmt.Sprintf("DROP DATABASE %s;", dbname)
	cmd = exec.Command(fmt.Sprintf("PGPASSWORD=123456 psql -h 10.230.12.86 -U cow -d stencil -c '%s'", q))

	if err := cmd.Run(); err != nil {
		fmt.Println(q)
		log.Fatal(err)
	}

	q = fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE diaspora_1000000 OWNER cow;", dbname)
	cmd = exec.Command(fmt.Sprintf("PGPASSWORD=123456 psql -h 10.230.12.86 -U cow -d stencil -c '%s'", q))

	if err := cmd.Run(); err != nil {
		fmt.Println(q)
		log.Fatal(err)
	}

	return nil
}

func DropAndRecreateDB(dbConn *sql.DB, dbname string) error {
	q := fmt.Sprintf("SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname in ('%s', 'diaspora_1000000') AND pid <> pg_backend_pid();", dbname)
	if _, err := dbConn.Exec(q); err != nil {
		fmt.Println(q)
		log.Fatal(err)
	}

	q = fmt.Sprintf("DROP DATABASE %s;", dbname)
	if _, err := dbConn.Exec(q); err != nil {
		fmt.Println(q)
		log.Fatal(err)
	}

	q = fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE diaspora_1000000 OWNER cow;", dbname)
	if _, err := dbConn.Exec(q); err != nil {
		fmt.Println(q)
		log.Fatal(err)
	}

	return nil
}

func InsertIntoDAGCounter(dbConn *sql.DB, person_id string, edges, nodes int) error {
	query := "INSERT INTO dag_counter (person_id, edges, nodes) VALUES ($1, $2, $3);"
	_, err := dbConn.Exec(query, person_id, edges, nodes)
	return err
}

func ReallyDeleteRowFromAppDB(tx *sql.Tx, table, id interface{}) error {
	query := fmt.Sprintf("DELETE FROM \"%s\" WHERE id = $1", table)
	if _, err := tx.Exec(query, id); err != nil {
		log.Println(query, id)
		log.Fatal("## DB ERROR: ", err)
		return err
	}
	return nil
}

func DeleteRowFromAppDB(tx *sql.Tx, table, id string) error {
	query := fmt.Sprintf("UPDATE \"%s\" SET mark_as_delete = $1 WHERE id = $2", table)
	if _, err := tx.Exec(query, true, id); err != nil {
		log.Println(query, "true", id)
		// log.Fatal("## DB ERROR: ", err)
		return err
	}
	return nil
}

func NewBag(tx *sql.Tx, pk, rowid, user_id, table, app string, migration_id int) error {
	query := "INSERT INTO data_bags (pk, rowid, user_id, \"table\", app, migration_id) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING"
	_, err := tx.Exec(query, pk, rowid, user_id, table, app, migration_id)
	return err
}

func UpdateBag(tx *sql.Tx, pk string, migration_id int, data []byte) error {
	query := "UPDATE data_bags SET data = $1, migration_id = $2 WHERE pk = $3;"
	_, err := tx.Exec(query, data, migration_id, pk)
	return err
}

func DeleteBagV2(tx *sql.Tx, pk string) error {
	query := "DELETE FROM data_bags WHERE pk = $1;"
	_, err := tx.Exec(query, pk)
	return err
}

func CreateNewBag(tx *sql.Tx, app, member, id, user_id, migration_id interface{}, data []byte) error {
	query := "INSERT INTO data_bags (app, member, id, data, user_id, migration_id) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING"
	_, err := tx.Exec(query, app, member, id, data, user_id, migration_id)
	return err
}

func GetRowsFromIDTableByTo(dbConn *sql.DB, app, member string, id int64) ([]map[string]interface{}, error) {

	query := "SELECT from_app, from_member, from_id, to_app, to_member, to_id, migration_id FROM identity_table WHERE to_app = $1 AND to_member = $2 AND to_id = $3;"
	return DataCall(dbConn, query, app, member, id)
}

func GetRowsFromIDTableByFrom(dbConn *sql.DB, app, member string, id int64) ([]map[string]interface{}, error) {
	query := "SELECT from_app, from_member, from_id, to_app, to_member, to_id, migration_id FROM identity_table WHERE from_app = $1 AND from_member = $2 AND from_id = $3;"
	return DataCall(dbConn, query, app, member, id)
}

func GetBagsV2(dbConn *sql.DB, app_id, user_id string, migration_id int) ([]map[string]interface{}, error) {
	query := "SELECT app, member, id, data, pk, user_id FROM data_bags WHERE user_id = $1 AND app = $2 AND migration_id != $3 ORDER BY pk DESC"
	return DataCall(dbConn, query, user_id, app_id, migration_id)
}

func GetBagByAppMemberIDV2(dbConn *sql.DB, user_id, app, member string, id int64, migration_id int) (map[string]interface{}, error) {
	query := "SELECT app, member, id, data, pk FROM data_bags WHERE user_id = $1 AND app = $2 AND member = $3 and id = $4 AND migration_id != $5"
	return DataCall1(dbConn, query, user_id, app, member, id, migration_id)
}

func CheckIfReferenceExists(dbConn *sql.DB, app, fromMember string, fromID int64, fromReference string) bool {

	q := fmt.Sprintf("SELECT * FROM reference_table WHERE app = '%s' AND from_member = '%s' AND from_id = '%v' AND from_reference = '%s'", app, fromMember, fromID, fromReference)
	if res, err := DataCall1(dbConn, q); err == nil {
		if len(res) > 0 {
			return true
		}
	} else {
		fmt.Println(q)
		log.Fatal(err)
	}
	return false
}

func CheckIfCompleteReferenceExists(dbConn *sql.DB, app, fromMember string, fromID int64, toMember string, toID int, fromReference, toReference string) bool {

	q := fmt.Sprintf("SELECT * FROM reference_table WHERE app = '%s' AND from_member = '%s' AND from_id = '%v' AND from_reference = '%s' AND to_member = '%s' AND to_id = '%v' AND to_reference = '%s'", app, fromMember, fromID, fromReference, toMember, toID, toReference)
	if res, err := DataCall1(dbConn, q); err == nil {
		if len(res) > 0 {
			return true
		}
	} else {
		fmt.Println(q)
		log.Fatal(err)
	}
	return false
}

func CreateNewReference(tx *sql.Tx, app, fromMember string, fromID int64, toMember string, toID int64, migration_id, fromReference, toReference string) error {

	// query := "INSERT INTO reference_table (app, from_member, from_id, from_reference, to_member, to_id, to_reference, migration_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING;"
	query := fmt.Sprintf("INSERT INTO reference_table (app, from_member, from_id, from_reference, to_member, to_id, to_reference, migration_id) VALUES ('%s', '%s', '%v', '%s', '%s', '%v', '%s', '%s') ON CONFLICT DO NOTHING;", app, fromMember, fromID, fromReference, toMember, toID, toReference, migration_id)

	_, err := tx.Exec(query)
	if err != nil {
		fmt.Println(query)
	}
	return err
}

func NewBagOld(tx *sql.Tx, rowid, user_id, tagName string, migration_id int) error {
	query := "INSERT INTO data_bags (rowid, user_id, tag, migration_id) VALUES ($1, $2, $3, $4)"
	_, err := tx.Exec(query, rowid, user_id, tagName, migration_id)
	return err
}

func NewRow(tx *sql.Tx, rowid, app_id, mflag string, copy_on_write bool) error {
	log.Fatal("Arrived at db.NewRow. Check why!")
	query := "INSERT INTO row_desc (rowid, app_id, copy_on_write, mflag) VALUES ($1, $2, $3, $4)"
	_, err := tx.Exec(query, rowid, app_id, copy_on_write, mflag)
	return err
}

func GetUserBags(dbConn *sql.DB, user_id, app_id string) ([]map[string]interface{}, error) {
	query := "SELECT string_agg(rowid::varchar, ',') as rowids, \"table\" FROM data_bags WHERE user_id = $1 AND app = $2 GROUP BY pk, \"table\""
	return DataCall(dbConn, query, user_id, app_id)
}

func GetBagAppAndTablesForUser(dbConn *sql.DB, user_id string) ([]map[string]interface{}, error) {
	query := "SELECT string_agg(\"table\"::varchar, ',') as tables, app FROM data_bags WHERE user_id = $1 AND app = $2 GROUP BY app"
	return DataCall(dbConn, query, user_id)
}

func GetAppsThatHaveBagsForUser(dbConn *sql.DB, user_id string) ([]map[string]interface{}, error) {
	query := "SELECT DISTINCT(app_id), app_name FROM migration_table join apps ON apps.pk = app_id WHERE user_id = $1 AND bag = true"
	return DataCall(dbConn, query, user_id)
}

func GetUserBagsByTables(dbConn *sql.DB, user_id, app_id, table string) ([]map[string]interface{}, error) {
	query := "SELECT string_agg(rowid::varchar, ',') as rowids FROM migration_table WHERE bag = true AND user_id = $1 AND app = $2 AND \"table_id\" = $3 GROUP BY group_id ORDER BY random()"
	return DataCall(dbConn, query, user_id, app_id, table)
}

func GetTablesForApp(dbConn *sql.DB, app_id string) ([]map[string]interface{}, error) {
	query := "SELECT pk as table_id FROM app_tables WHERE app_id = $1"
	return DataCall(dbConn, query, app_id)
}

func GetColumnsFromAppSchema(dbConn *sql.DB, table_id string) ([]map[string]interface{}, error) {
	query := "SELECT column_name FROM app_schemas WHERE table_id = $1"
	return DataCall(dbConn, query, table_id)
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

func FetchForMapping(dbConn *sql.DB, targetTable, targetCol, conditionCol, conditionVal string) (map[string]interface{}, error) {
	q := fmt.Sprintf("SELECT %s FROM %s WHERE %s = '%s'", targetCol, targetTable, conditionCol, conditionVal)
	// fmt.Println(q)
	return DataCall1(dbConn, q)
}

func SetAppID(tx *sql.Tx, pk, app_id string) error {
	log.Fatal("Arrived at db.SetAppID. Check why!")
	q := "UPDATE row_desc SET app_id = $1 WHERE rowid = $2"
	_, err := tx.Exec(q, app_id, pk)
	return err
}

func SetMFlag(tx *sql.Tx, pk, flag string) error {
	log.Fatal("Arrived at db.SetMFlag. Check why!")
	q := "UPDATE row_desc SET mflag = $1 WHERE rowid IN ($2)"
	_, err := tx.Exec(q, flag, pk)
	return err
}

func MUpdate(tx *sql.Tx, pk, flag, app_id string) error {
	log.Fatal("Arrived at db.MUpdate. Check why!")
	q := "UPDATE row_desc SET mflag = $1, app_id = $2 WHERE rowid IN ($3)"
	_, err := tx.Exec(q, flag, app_id, pk)
	return err
}

func BUpdate(dbConn *sql.DB, pk, flag, app_id string) error {
	log.Fatal("Arrived at db.BUpdate. Check why!")
	q := "UPDATE row_desc SET mflag = $1, app_id = $2 WHERE rowid IN ($3)"
	_, err := dbConn.Exec(q, flag, app_id, pk)
	return err
}

func PKReplaceRowDesc(tx *sql.Tx, newRowID, oldRowIDs string) error {
	log.Fatal("Arrived at db.PKReplaceRowDesc. Check why!")
	q := fmt.Sprintf("UPDATE row_desc SET rowid = $1 WHERE rowid IN (%s)", oldRowIDs)
	_, err := tx.Exec(q, newRowID)
	return err
}

func PKReplace(tx *sql.Tx, newPK, oldPK, table string) error {

	q := fmt.Sprintf("UPDATE %s SET pk = $1 WHERE pk IN (%s)", table, oldPK)
	_, err := tx.Exec(q, newPK)
	return err
}

func DeleteFromRowDescByRowID(tx *sql.Tx, rowid string) error {
	log.Fatal("Arrived at db.DeleteFromRowDescByRowID. Check why!")
	q := "DELETE FROM row_desc WHERE rowid = $1"
	_, err := tx.Exec(q, rowid)
	return err
}

func DeleteFromMigrationTable(tx *sql.Tx, rowid, table_id string) error {
	q := "DELETE FROM migration_table WHERE row_id = $1 AND table_id = $2"
	_, err := tx.Exec(q, rowid, table_id)
	return err
}

func MarkRowAsDeleted(tx *sql.Tx, rowid, table_id string) error {

	q := "UPDATE migration_table SET mark_as_delete = $1 WHERE row_id = $2 AND table_id = $3"
	_, err := tx.Exec(q, true, rowid, table_id)
	return err
}

func RemoveBag(tx *sql.Tx, rowid, table_id string) error {

	q := "UPDATE migration_table SET bag = $1, migration_id = NULL, mark_as_delete = $4 WHERE row_id = $2 AND table_id = $3"
	_, err := tx.Exec(q, false, rowid, table_id, true)
	return err
}

func RevertBag(tx *sql.Tx, rowid, table_id, migration_id string) error {

	q := "UPDATE migration_table SET bag = $1, mark_as_delete = $2, mflag = $3, migration_id = $4 WHERE row_id = $5 AND table_id = $6"
	_, err := tx.Exec(q, false, false, 1, migration_id, rowid, table_id)
	return err
}

func MarkRowAsBag(tx *sql.Tx, rowid, table_id, migration_id, user_id string) error {

	q := "UPDATE migration_table SET bag = true, mark_as_delete = true, migration_id = $1, user_id = $4 WHERE row_id = $2 AND table_id = $3"
	_, err := tx.Exec(q, migration_id, rowid, table_id, user_id)
	return err
}

func DeleteFromRowDescByRowIDAndAppID(tx *sql.Tx, rowids, appid string) error {
	log.Fatal("Arrived at db.DeleteFromRowDescByRowIDAndAppID. Check why!")
	q := fmt.Sprintf("DELETE FROM row_desc WHERE rowid in (%s) and app_id = $1", rowids)
	_, err := tx.Exec(q, appid)
	return err
}

func InsertIntoMigrationTable(tx *sql.Tx, dstApp, dstRow, orgRow, COW, dstTable, mflag, migration_id string) error {
	q := "INSERT INTO migration_table (app_id, group_id, row_id, copy_on_write, table_id, mflag, migration_id) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := tx.Exec(q, dstApp, dstRow, orgRow, COW, dstTable, mflag, migration_id)
	return err
}

func GetUnmigratedUsers() ([]map[string]interface{}, error) {
	dbConn := GetDBConn(STENCIL_DB)
	sql := "SELECT user_id FROM user_table WHERE user_id NOT IN (SELECT DISTINCT user_id FROM migration_registration) ORDER BY user_id ASC"
	return DataCall(dbConn, sql)
}

func GetAppRootMemberID(stencilDBConn *sql.DB, appID string) string {

	query := fmt.Sprintf(`SELECT root_member_id from app_root_member where app_id = %s`, appID)
	fmt.Printf("@db.GetAppRootMemberID | %s\n", query)
	data, err := DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprint(data["root_member_id"])
}

func TableID(dbConn *sql.DB, table, app string) (string, error) {
	sql := fmt.Sprintf("SELECT pk FROM app_tables WHERE app_id = '%s' and table_name = '%s'", app, table)
	if res, err := DataCall1(dbConn, sql); err == nil {
		if pk, ok := res["pk"]; ok {
			return fmt.Sprint(pk), nil
		} else {
			fmt.Println(fmt.Sprintf("@db.TableID | Args | table: %s| app: %s", table, app))
			return "", errors.New("Something bad with the returned result!")
		}
	} else {
		return "", err
	}
}

func TableName(dbConn *sql.DB, table, app string) (string, error) {
	sql := fmt.Sprintf("SELECT table_name FROM app_tables WHERE app_id = '%s' and pk = '%s'", app, table)
	if res, err := DataCall1(dbConn, sql); err == nil {
		if pk, ok := res["table_name"]; ok {
			return fmt.Sprint(pk), nil
		} else {
			fmt.Println(fmt.Sprintf("@db.TableName | Args | table: %s| app: %s", table, app))
			return "", errors.New("Something bad with the returned result!")
		}
	} else {
		return "", err
	}
}

func CheckPhyRowExists(ptab, rowid string, dbConn *sql.DB) bool {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE pk = $1", ptab)
	if res, err := DataCall1(dbConn, sql, rowid); err == nil {
		if len(res) > 0 {
			return true
		}
	} else {
		log.Fatal(err)
	}
	return false
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

func FinishMigration(dbConn *sql.DB, migrationID, size int) bool {
	query := "UPDATE migration_registration SET end_time = now(), msize = $2 WHERE migration_id = $1;"
	if _, err := dbConn.Exec(query, migrationID, size); err != nil {
		fmt.Println(query, migrationID, size)
		log.Fatal("Insert Error in FinishMigration", err)
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
				// resultStr += table + "." + columns[i] + " AS \"" + table + "." + columns[i] + "\","
				resultStr += fmt.Sprintf("\"%s\".\"%s\" AS \"%s.%s\",", table, columns[i], table, columns[i])
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

func DataCallIgnoreVisited(db *sql.DB, SQL string, visited []string, args ...interface{}) ([]map[string]interface{}, error) {

	// db := GetDBConn(app)
	// log.Println(SQL, args)
	if rows, err := db.Query(SQL, args...); err != nil {
		color.Red.Println("ERROR:", SQL, args, err)
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

func DataCall(db *sql.DB, SQL string, args ...interface{}) ([]map[string]interface{}, error) {

	// db := GetDBConn(app)
	// log.Println(SQL, args)
	if rows, err := db.Query(SQL, args...); err != nil {
		color.Danger.Printf("ERROR | %s\nQuery | %s\nArgs | %v\n", err, SQL, args)
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
	for {
		rowid = rand.Int31n(2147483647)
		q := fmt.Sprintf("SELECT row_id FROM migration_table WHERE row_id = %d", rowid)
		if v, err := DataCall1(dbConn, q); err != nil {
			fmt.Println(q)
			log.Fatal("@db.GetNewRowID: ", err)
		} else if v == nil {
			break
		}
	}
	return rowid
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

func TxnExecute1(dbConn *sql.DB, query string) error {
	tx, err := dbConn.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(query); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
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

func SaveForLEvaluation(tx *sql.Tx, srcApp, dstApp, srcTable, dstTable, srcID, dstID, srcCol, dstCol, migrationID interface{}) error {
	query := "INSERT INTO evaluation (src_app, dst_app, src_table, dst_table, src_id, dst_id, src_cols, dst_cols, migration_id, added_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now())"
	_, err := tx.Exec(query, srcApp, dstApp, srcTable, dstTable, srcID, dstID, srcCol, dstCol, migrationID)
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

func GetUserListFromAppDB(appName, userTable, userCol string) []string {
	dbConn := GetDBConn(appName)
	defer dbConn.Close()
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY random()", userCol, userTable)
	if res, err := DataCall(dbConn, query); err == nil {
		var users []string
		for _, row := range res {
			users = append(users, fmt.Sprint(row[userCol]))
		}
		return users
	} else {
		fmt.Println(query, userTable, userCol, appName)
		log.Fatal(err)
	}
	return nil
}

func GetNextUserFromAppDB(appName, userTable, userCol string, offset int) (string, error) {
	dbConn := GetDBConn(appName)
	defer dbConn.Close()
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY %s ASC OFFSET %d", userCol, userTable, userCol, offset)
	if res, err := DataCall1(dbConn, query); err == nil {
		return fmt.Sprint(res[userCol]), nil
	} else {
		fmt.Println(query, userTable, userCol, appName)
		log.Fatal(err)
		return "", err
	}
}

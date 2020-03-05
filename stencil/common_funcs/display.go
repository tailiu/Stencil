package common_funcs

import (
	"database/sql"
	"stencil/db"
	"fmt"
	"log"
)

func GetSrcDstAppIDsUserIDByMigrationID(stencilDBConn *sql.DB,
	migrationID int) (string, string, string) {

	query := fmt.Sprintf(
		`SELECT src_app, dst_app, user_id FROM migration_registration 
		WHERE migration_id = %d`,
		migrationID,
	)

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	srcApp := fmt.Sprint(data["src_app"]) 
	dstApp := fmt.Sprint(data["dst_app"])
	userID := fmt.Sprint(data["user_id"])

	return srcApp, dstApp, userID

}

func GetAppNameByAppID(stencilDBConn *sql.DB, appID string) string {

	query := fmt.Sprintf("select app_name from apps where pk = %s", appID)

	// log.Println(query)

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprint(data["app_name"])

}

func GetTableIDNamePairsInApp(stencilDBConn *sql.DB, 
	app_id string) []map[string]interface{} {

	query := fmt.Sprintf(
		`select pk, table_name from app_tables where app_id = %s`,
		app_id,
	)

	result, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func GetAppIDNamePairs(stencilDBConn *sql.DB) map[string]string {

	appIDNamePairs := make(map[string]string)

	query := fmt.Sprintf("select app_name, pk from apps")

	// log.Println(query)

	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	for _, data1 := range data {
		appIDNamePairs[fmt.Sprint(data1["pk"])] = fmt.Sprint(data1["app_name"])
	}

	return appIDNamePairs

}

func GetTableIDNamePairs(stencilDBConn *sql.DB) map[string]string {

	tableIDNamePairs := make(map[string]string)

	query := fmt.Sprintf("select pk, table_name from app_tables;")

	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	for _, data1 := range data {
		tableIDNamePairs[fmt.Sprint(data1["pk"])] = 
			fmt.Sprint(data1["table_name"])
	}

	return tableIDNamePairs

}
package db

import (
	"log"
	"database/sql"
	"fmt"
)

func ConvertAppTableToPhyTable(appName, tableName string) {

}

func GetAppIDByAppName(stencilDBConn *sql.DB, app string) string {
	query := fmt.Sprintf("SELECT pk from apps WHERE app_name = '%s'", app)
	res := GetAllColsOfRows(stencilDBConn, query)

	if res[0]["pk"] == "" {
		log.Fatal("AppID does not exist!")
	}

	return res[0]["pk"]
}

package db

import (
	// "log"
	// "database/sql"
	// "fmt"
)

func ConvertAppTableToPhyTable(appName, tableName string) {

}

// func GetAppIDByAppName(stencilDBConn *sql.DB, app string) string {
// 	if app == "mastodon" || app == "mastodon_old" {
// 		app = "mastodon"
// 	} else if (app == "diaspora" || app == "diaspora_old") {
// 		app = "diaspora"
// 	}
// 	query := fmt.Sprintf("SELECT pk from apps WHERE app_name = '%s'", app)
// 	res := GetAllColsOfRows(stencilDBConn, query)

// 	if res[0]["pk"] == "" {
// 		log.Fatal("AppID does not exist!")
// 	}

// 	return res[0]["pk"]
// }

package evaluation

import (
	"stencil/db"
	"fmt"
	"database/sql"
	"log"
	// "strings"
)

func GetAllMigrationIDsOfAppWithConds(stencilDBConn *sql.DB, appID string, extraConditions string) []map[string]interface{} {
	query := fmt.Sprintf("select * from migration_registration where dst_app = '%s' %s;", 
		appID, extraConditions)
	// log.Println(query)

	migrationIDs, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return migrationIDs
}

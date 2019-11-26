package reference_resolution

import (
	"stencil/db"
	"database/sql"
	"fmt"
)

func GetFromReferences(stencilDBConn *sql.DB, appConfig *config.AppConfig, IDRow map[string]interface{}) []map[string]interface{} {

	query := fmt.Sprintf("SELECT * FROM references WHERE to_app = %d and to_member = %d and to_id = %d and migration_id = %d",
		appConfig.AppID, member, id, migrationID)
	
	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)

	return data

}
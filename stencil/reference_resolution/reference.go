package reference_resolution

import (
	"stencil/db"
	"database/sql"
	"fmt"
	"log"
)

func getFromReferences(stencilDBConn *sql.DB, migrationID int, IDRow map[string]string) []map[string]interface{} {

	query := fmt.Sprintf("SELECT * FROM reference_table WHERE app = %s and from_member = %s and from_id = %s and migration_id = %d;",
		IDRow["from_app"], IDRow["from_member"], IDRow["from_id"], migrationID)
	
	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)

	return data

}

func updateReferences() {
	
}
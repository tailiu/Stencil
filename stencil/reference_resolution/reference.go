package reference_resolution

import (
	"stencil/db"
	"stencil/config"
	"fmt"
	"log"
)

func getFromReferences(displayConfig *config.DisplayConfig, IDRow map[string]string) []map[string]interface{} {

	query := fmt.Sprintf("SELECT * FROM reference_table WHERE app = %s and from_member = %s and from_id = %s and migration_id = %d;",
		IDRow["from_app"], IDRow["from_member"], IDRow["from_id"], displayConfig.MigrationID)
	
	data, err := db.DataCall(displayConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)

	return data

}

func updateReferences() {

}
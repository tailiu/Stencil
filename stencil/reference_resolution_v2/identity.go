package reference_resolution_v2

import (
	"fmt"
	"log"
	"stencil/db"
)

func getUpdateIDInDisplayFlagsQuery(refResolutionConfig *RefResolutionConfig, 
	table, IDToBeUpdated, id string) string {
	
	query := fmt.Sprintf(
		`UPDATE display_flags SET id = %s, updated_at = now() 
		WHERE app_id = %s and table_id = %s 
		and id = %s and migration_id = %d;`,
		id, refResolutionConfig.appID, 
		refResolutionConfig.appTableNameIDPairs[table],
		IDToBeUpdated, refResolutionConfig.migrationID,
	)

	return query
}

func getInsertIntoIDChangesTableQuery(refResolutionConfig *RefResolutionConfig,
	table, IDToBeUpdated, id string) string {

	query := fmt.Sprintf(
		`INSERT INTO id_changes (app_id, table_id, old_id, new_id, migration_id)
		VALUES (%s, %s, %s, %s, %d)`,
		refResolutionConfig.appID, 
		refResolutionConfig.appTableNameIDPairs[table],
		IDToBeUpdated,
		id,
		refResolutionConfig.migrationID,
	)

	return query
}

func getIDsOfDataToBeUpdated(refResolutionConfig *RefResolutionConfig,
	memberToBeUpdated, attrValToBeUpdated, attrToBeUpdated string) []string {
	
	query := fmt.Sprintf(
		"SELECT id FROM %s WHERE %s = '%s'",
		memberToBeUpdated, attrToBeUpdated, attrValToBeUpdated,
	)

	log.Println(query)

	data, err := db.DataCall(refResolutionConfig.appDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	var res []string
	for _, data1 := range data {
		res = append(res, fmt.Sprint(data1["id"])) 
	}

	return res

}
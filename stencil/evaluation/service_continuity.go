package evaluation

import (
	"fmt"
	"log"
	"stencil/db"
	"time"
)

func getDataDowntimeOfMigration(evalConfig *EvalConfig, migrationID string) []time.Duration {

	var downtime []time.Duration

	query := fmt.Sprintf(
		`SELECT created_at, displayed_at FROM evaluation
		WHERE migration_id = '%s' and dst_table != 'n/a' 
		and displayed_at is not null`, migrationID)

	log.Println(query)
	
	result, err := db.DataCall(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	for _, data1 := range result {

		downtime = append(
			downtime, 
			data1["displayed_at"].(time.Time).Sub(data1["created_at"].(time.Time)),
		)

	}

	return downtime

}
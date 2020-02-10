package evaluation

import (
	"fmt"
	"log"
	"stencil/db"
	"database/sql"
	"time"
)

func getDataDowntimeOfMigration(dbConn *sql.DB,
	migrationID string) []time.Duration {

	var downtime []time.Duration

	query := fmt.Sprintf(
		`SELECT created_at, displayed_at, dst_table, dst_id
		FROM evaluation WHERE migration_id = '%s'
		and displayed_at is not null`, 
		migrationID)

	// log.Println(query)
	
	result, err := db.DataCall(dbConn, query)
	if err != nil {
		log.Fatal(err)
	}

	checkedData := make(map[string]bool)

	for _, data1 := range result {

		dstTable := fmt.Sprint(data1["dst_table"])
		dstID := fmt.Sprint(data1["dst_id"])

		key := dstTable + ":" + dstID

		if _, ok := checkedData[key]; ok {
			// log.Println("Duplicate key:", key)
			continue
		}

		downtime = append(
			downtime, 
			data1["displayed_at"].(time.Time).Sub(data1["created_at"].(time.Time)),
		)

		checkedData[key] = true

	}

	return downtime

}
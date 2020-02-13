package evaluation

import (
	"strconv"
	"time"
	"fmt"
	"log"
	"database/sql"
	"stencil/db"
	"stencil/transaction"
)

func getMigrationStartTime(dbConn *sql.DB, 
	migrationID int) time.Time {

	log_txn := new(transaction.Log_txn)
	log_txn.DBconn = dbConn
	log_txn.Txn_id = migrationID
	
	if startTime := log_txn.GetCreatedAt("BEGIN_TRANSACTION"); 
		len(startTime) == 1 {
		
		return startTime[0]
	} else {

		panic("Should never happen here!")
	}
}

func GetMigrationTime(stencilDBConn *sql.DB, 
	migrationID int) time.Duration {
	
	return getMigrationEndTime(stencilDBConn, migrationID).
		Sub(getMigrationStartTime(stencilDBConn, migrationID))

}

func GetDisplayTime(stencilDBConn *sql.DB, 
	migrationID string) time.Duration {
	
	query := fmt.Sprintf(
		`SELECT start_time, end_time FROM display_registration 
		WHERE migration_id = %s`,
		migrationID,
	)

	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if len(data) != 1 {
		log.Fatal("Get more than one row in the display_registration")
	}

	endTime := data[0]["end_time"].(time.Time)
	startTime := data[0]["start_time"].(time.Time)
	displayTime := endTime.Sub(startTime)

	return displayTime

}

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

func getTotalTimeOfMigration(dbConn *sql.DB,
	migrationID string) time.Duration {

	query := fmt.Sprintf(
		`SELECT end_time FROM display_registration 
		WHERE migration_id = %s`,
		migrationID,
	)

	data, err := db.DataCall1(dbConn, query)
	if err != nil {
		log.Fatal(err)
	}

	migrationIDInt, err1 := strconv.Atoi(migrationID)
	if err1 != nil {
		log.Fatal(err1)
	}

	migrationStartTime := getMigrationStartTime(dbConn, migrationIDInt)

	totalTime := data["end_time"].(time.Time).Sub(migrationStartTime)

	return totalTime

}

func calculateTimeInPercentage(times []time.Duration, 
	totalTime time.Duration) []float64 {

	var timesInPercentage []float64
	
	for _, time := range times {
		timesInPercentage = append(timesInPercentage, time.Seconds() / totalTime.Seconds())
	}

	return timesInPercentage

}
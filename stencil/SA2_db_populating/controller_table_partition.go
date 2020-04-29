package SA2_db_populating

import (
	"stencil/db"
	"sync"
	"time"
	"log"
)

func PupulatingControllerWithCheckpointAndTruncateWithTablePartition() {

	var limit, startPoint, endPoint, checkpointFeq int64

	// ******************* Setting Parameters Start *******************
	
	dbID := "10"

	table := "aspect_visibilities"
	startPoint = 1000000
	endPoint = -1

	db.STENCIL_DB = "stencil_exp_sa2_" + dbID

	srcDB := db.STENCIL_DB
	dstDB := "stencil_exp_sa2_100k" 
	
	appName := "diaspora_100000_sa2_" + dbID
	appID := "1"

	checkpointFeq = 200000

	limit = 2500
	threadNum := 10

	isStencilOnBladeServer := false
	isAppOnBladeServer := false

	// ****************************** End ******************************

	log.Println("Start populating data from table:", table)
	log.Println("Start point:", startPoint)

	startTime := time.Now()

	stencilDBConn := db.GetDBConn(db.STENCIL_DB, isStencilOnBladeServer)
	defer stencilDBConn.Close()

	appDBConn := db.GetDBConn(appName, isAppOnBladeServer)
	defer appDBConn.Close()

	var dataSeqStart, dataSeqStep, dataSeqEnd, rowCount int64

	if endPoint == -1 {
		rowCount = getTotalRowCountOfTable(appDBConn, table)
	} else {
		rowCount = endPoint
	}
	
	tableNum, ok := tableNameRangeIndexMap[table]
	if !ok {
		log.Fatal("Error: Table name is incorrect!")
	}

	migrationTable := "migration_table_" + tableNum

	log.Println("End point:", rowCount)
	log.Println("Corresponding migration table:", migrationTable)

	dataSeqStart = startPoint

	var wg sync.WaitGroup

	for {
		dataSeqEnd = dataSeqStart + checkpointFeq 
		if dataSeqEnd > rowCount + checkpointFeq {
			log.Fatal("Error: Something is wrong here!")
		} else if dataSeqEnd == rowCount + checkpointFeq {
			break
		} else if dataSeqEnd > rowCount && dataSeqEnd < rowCount + checkpointFeq {
			dataSeqStep = (rowCount - dataSeqStart) / int64(threadNum)
			wg.Add(threadNum)
			startPopulatingThreads(
				stencilDBConn, appDBConn, threadNum, 
				dataSeqStart, dataSeqStep, rowCount, limit,
				appName, appID, table, &wg,
			)
			wg.Wait()
			break
		} else {
			dataSeqStep = (dataSeqEnd - dataSeqStart) / int64(threadNum)
			wg.Add(threadNum)
			startPopulatingThreads(
				stencilDBConn, appDBConn, threadNum, 
				dataSeqStart, dataSeqStep, dataSeqEnd, limit,
				appName, appID, table, &wg,
			)
			wg.Wait()
		}
		dataSeqStart = dataSeqEnd
		checkpointTruncateWithTablePartition(srcDB, dstDB, migrationTable)
	}

	checkpointTruncateWithTablePartition(srcDB, dstDB, migrationTable)

	endTime := time.Now()

	log.Println("Populating", table, "is done")
	log.Println("Time used:", endTime.Sub(startTime))

}
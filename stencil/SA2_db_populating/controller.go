package SA2_db_populating

import (
	"stencil/apis"
	"stencil/db"
	"sync"
	"database/sql"
	"time"
	"log"
)

// Note that [dataSeqStart, dataSeqEnd)
func PopulateRangeOfOneTable(stencilDBConn, appDBConn *sql.DB, 
	dataSeqStart, dataSeqEnd, limit int64, wg *sync.WaitGroup,
	appName, appID, table string) {

	log.Println("Thread working on range:", dataSeqStart, "-", dataSeqEnd)

	defer wg.Done()

	if dataSeqStart > dataSeqEnd {
		return 
	}

	offset := dataSeqStart

	for {
		if offset + limit > dataSeqEnd {
			// if offset (start) is equal to dataSeqEnd, there is no need to populate
			if offset != dataSeqEnd {
				apis.Port(
					appName, appID, table, 
					dataSeqEnd - offset, offset,
					appDBConn, stencilDBConn,
				)

			}
			break
		} else {
			apis.Port(
				appName, appID, table, 
				limit, offset,
				appDBConn, stencilDBConn,
			)	
		}
		offset += limit
	}

}

func PupulatingControllerForOneTable(fromApp string, toStencilDB string, tableName string, end int64) {

	var limit, startPoint, endPoint int64

	// ******************* Setting Parameters Start *******************
	
	// id := "12"

	table := tableName

	startPoint = 0
	endPoint = end

	// db.STENCIL_DB = "stencil_sa2_100k"
	db.STENCIL_DB = toStencilDB

	// appName := "diaspora_100000"
	appName := fromApp
	appID := "1"

	limit = 1000
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

	var wg sync.WaitGroup

	wg.Add(threadNum)

	var dataSeqStart, dataSeqStep, rowCount int64

	if endPoint == -1 {
		rowCount = getTotalRowCountOfTable(appDBConn, table)
	} else {
		rowCount = endPoint
	}

	log.Println("End point:", rowCount) 

	dataSeqStart = startPoint
	
	dataSeqStep = (rowCount - startPoint) / int64(threadNum)

	for i := 0; i < threadNum; i++ {
		if i != threadNum - 1 {
			go PopulateRangeOfOneTable(
				stencilDBConn, appDBConn, 
				dataSeqStart, dataSeqStart + dataSeqStep, limit, &wg,
				appName, appID, table,
			)
		} else {
			go PopulateRangeOfOneTable(
				stencilDBConn, appDBConn, 
				dataSeqStart, rowCount, limit, &wg,
				appName, appID, table,
			)
		}
		dataSeqStart += dataSeqStep
	}

	wg.Wait()

	endTime := time.Now()

	log.Println("Populating", table, "is done")
	log.Println("Time used:", endTime.Sub(startTime))

}

func PupulatingControllerForAllTables(fromApp string, toStencilDB string) {

	startTime := time.Now()
	
	appDBConn := db.GetDBConn(fromApp)
	defer appDBConn.Close()

	rowCounts := listRowCountsOfDB(appDBConn)

	for table, rowCount := range rowCounts {
		if rowCount != 0 {

			log.Println("==============================")
			log.Println("Start Populating Table:", table)
			log.Println("==============================")
			
			PupulatingControllerForOneTable(fromApp, toStencilDB, table, rowCount)
		
		} else {

			log.Println("Skip table:", table, "since its row count is 0")

		}
	}

	endTime := time.Now()

	log.Println("Populating DB is Done!")
	log.Println("Time used:", endTime.Sub(startTime))
}

func PupulatingControllerForAllTablesHandlingPKs(fromApp string, toStencilDB string) {

	DropPrimaryKeysOfSA2TablesWithoutPartitions(toStencilDB)

	PupulatingControllerForAllTables(fromApp, toStencilDB)

}

func startPopulatingThreads(stencilDBConn, appDBConn *sql.DB,
	threadNum int,  dataSeqStart, dataSeqStep, dataSeqEnd, limit int64,
	appName, appID, table string, wg *sync.WaitGroup) {

	for i := 0; i < threadNum; i++ {
				
		if i != threadNum - 1 {

			go PopulateRangeOfOneTable(
				stencilDBConn, appDBConn, 
				dataSeqStart, dataSeqStart + dataSeqStep, limit, wg,
				appName, appID, table,
			)

		} else {

			go PopulateRangeOfOneTable(
				stencilDBConn, appDBConn, 
				dataSeqStart, dataSeqEnd, limit, wg,
				appName, appID, table,
			)

		}

		dataSeqStart += dataSeqStep

	}

}

func PupulatingControllerWithCheckpointAndTruncate() {

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

		checkpointTruncate(srcDB, dstDB, migrationTable)

	}

	checkpointTruncate(srcDB, dstDB, migrationTable)

	endTime := time.Now()

	log.Println("Populating", table, "is done")
	log.Println("Time used:", endTime.Sub(startTime))

}
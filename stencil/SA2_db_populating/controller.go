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
func PopulateSA2Tables(stencilDBConn, appDBConn *sql.DB, 
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

// First population (stencil_exp_sa2):
// My machine: people(finished), users(finished), likes, comments
// VM: profiles(finished), notifications, posts
// Blade server: conversations(done), conversation_visibilities(done),
//				notification_actors, aspect_visibilities				

// Second population (stencil_exp_sa2_1):
// My machine: 
// Blade server: photos

// Third population (stencil_exp_sa2_2):
// Blade server: remaining comments
// My machine:

// Forth population (stencil_exp_sa2_4):
// Blade server: remaining likes
// My machine:

// Fifth population (stencil_exp_sa2_5)
// Blade server: remaining notifications
// My machine:

// sixth population (stencil_exp_sa2_6)
// Blade server: remaining notification_actors
// My machine:

// seventh population (stencil_exp_sa2_7)
// Blade server: conversations(done and deleted), 
// 				conversation_visibilities(done and deleted),
// 				posts
// My machine:

// eighth population (stencil_exp_sa2_8)
// Blade server: messages
// My machine:
func PupulatingController() {

	var limit, startPoint, endPoint int64

	// ******************* Setting Parameters Start *******************
	
	id := "12"

	table := "notifications"

	startPoint = 0
	endPoint = 200000

	db.STENCIL_DB = "stencil_exp_sa2_" + id
	
	appName := "diaspora_100000_sa2_" + id
	appID := "1"

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

			go PopulateSA2Tables(
				stencilDBConn, appDBConn, 
				dataSeqStart, dataSeqStart + dataSeqStep, limit, &wg,
				appName, appID, table,
			)

		} else {

			go PopulateSA2Tables(
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

func startPopulatingThreads(stencilDBConn, appDBConn *sql.DB,
	threadNum int,  dataSeqStart, dataSeqStep, dataSeqEnd, limit int64,
	appName, appID, table string, wg *sync.WaitGroup) {

	for i := 0; i < threadNum; i++ {
				
		if i != threadNum - 1 {

			go PopulateSA2Tables(
				stencilDBConn, appDBConn, 
				dataSeqStart, dataSeqStart + dataSeqStep, limit, wg,
				appName, appID, table,
			)

		} else {

			go PopulateSA2Tables(
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
	
	dbID := "12"

	table := "notifications"
	startPoint = 0
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
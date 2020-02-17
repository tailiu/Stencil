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

			apis.Port(
				appName, appID, table, 
				dataSeqEnd - offset, offset,
				appDBConn, stencilDBConn,
			)

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

// First population:
// My machine: people(finished), users(finished), likes, comments
// VM: profiles(finished), notifications, posts
// Blade server: notification_actors, aspect_visibilities
// Second population:
// My machine: users 
// Blade server: people
func PupulatingController() {

	var limit int64

	// ******************* Setting Parameters Start *******************
	
	db.STENCIL_DB = "stencil_exp_sa2_1"

	table := "users"

	appName := "diaspora_1000000_sa2_1"
	appID := "1"

	limit = 10000
	threadNum := 10

	// ****************************** End ******************************

	log.Println("Start populating data from table:", table)

	startTime := time.Now()

	stencilDBConn := db.GetDBConn(db.STENCIL_DB)
	defer stencilDBConn.Close()

	appDBConn := db.GetDBConn(appName)
	defer appDBConn.Close()

	var wg sync.WaitGroup

	wg.Add(threadNum)

	rowCount := getTotalRowCountOfTable(appDBConn, table)

	log.Println("Total row count:", rowCount)
	
	var dataSeqStart, dataSeqStep int64

	dataSeqStart = 0
	
	dataSeqStep = rowCount / int64(threadNum)

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
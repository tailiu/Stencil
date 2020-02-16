package SA2_db_populating

import (
	"stencil/apis"
	"stencil/db"
	"sync"
)

// Note that [dataSeqStart, dataSeqEnd)
func PopulateSA2Tables(dataSeqStart, dataSeqEnd, ) {

	var offset, limit int64 

	offset = 99
	limit = 2000	

	apis.Port(appName, appID, table, limit, offset)

}

// First population:
// My machine: people(finished), users(finished), likes, comments
// VM: profiles(finished), notifications, posts
// Blade server: notification_actors, aspect_visibilities
// Second population:
// VM: conversations(finished), conversation_visibilities(finished), 
// 		users(finished), photos 
// Blade server: people(finished), profiles(finished)
func PupulatingController() {

	db.STENCIL_DB = "stencil_exp_sa2_test"

	appName := "diaspora_1000000"
	appID := "1"

	threadNum := 50 

	table := "notifications"

	var wg sync.WaitGroup

	wg.Add(threadNum)

	rowCount := getTotalRowCountOfTable(appName, table)

	log.Println("Total row count:", rowCount)
	
	dataSeqStart := 0
	
	dataSeqStep := rowCount / threadNum

	for i := 0; i < threadNum; i++ {
		
		if i != threadNum - 1 {

			go PopulateSA2Tables(dataSeqStart, dataSeqStart + dataSeqStep)

		} else {

			go PopulateSA2Tables(dataSeqStart, rowCount)

		}

		dataSeqStart += dataSeqStep

	}

	wg.Wait()


}
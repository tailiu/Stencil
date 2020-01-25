package SA1_display

import (
	"stencil/config"
	"stencil/db"
	"log"
	"sync"
	"time"
)

const WAIT_FOR_MIGRATION_START_INTERVAL = 100 * time.Millisecond

func displayController(migrationID, threadNum int, wg *sync.WaitGroup) {

	// If the destination app database is not in the new server, newDB is false
	newDB := false

	// If the display controller needs to resolve references, resolveReference is true
	resolveReference := true

	dConfig := CreateDisplayConfig(migrationID, resolveReference, newDB)

	log.Println("Migration ID:", migrationID)

	log.Println("Total Display Thread(s):", threadNum)

	for i := 0; i < threadNum; i++ {

		log.Println("Start thread:", i + 1)

		go func(dConfig *displayConfig) {

			defer wg.Done()
			
			DisplayThread(dConfig)
		
		} (dConfig)

	}

	// log.Println("############### End Display Controller ###############")

}

func StartDisplay(uid, srcAppID, dstAppID, migrationType string, 
	threadNum int, wg *sync.WaitGroup) {

	var migrationIDs []int

	stencilDBConn := db.GetDBConn(config.StencilDBName)

	for {

		migrationIDs = getMigrationIDs(stencilDBConn, 
			uid, srcAppID, dstAppID, migrationType)
		
		// log.Println("*******")
		// log.Println(migrationIDs)

		if migrationNum := len(migrationIDs); migrationNum == 0 {
			time.Sleep(WAIT_FOR_MIGRATION_START_INTERVAL)
		} else if migrationNum == 1 {
			break
		} else {
			log.Fatal("There are more than one migration")
		}

	}

	displayController(migrationIDs[0], threadNum, wg)

}
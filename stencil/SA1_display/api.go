package SA1_display

import (
	// "stencil/config"
	"stencil/db"
	"log"
	"sync"
	"time"
)

const WAIT_FOR_MIGRATION_START_INTERVAL = 100 * time.Millisecond
const CHECK_MIGRATION_COMPLETE_INTERVAL = time.Second

func displayController(migrationID, threadNum int, 
	wg *sync.WaitGroup, displayInFirstPhase bool) {

	// If the destination app database is not in the new server, newDB is false
	newDB := false

	// If the display controller needs to resolve references, resolveReference is true
	resolveReference := true

	dConfig := CreateDisplayConfig(migrationID, resolveReference, 
		newDB, displayInFirstPhase)

	log.Println("Migration ID:", migrationID)

	log.Println("Total Display Thread(s):", threadNum)

	for i := 0; i < threadNum; i++ {

		log.Println("Start Display Thread:", i + 1)

		go func(dConfig *displayConfig) {

			defer wg.Done()
			
			DisplayThread(dConfig)
		
		} (dConfig)

	}

	// log.Println("############### End Display Controller ###############")

}

func waitGetMigrationID(uid, srcAppID, dstAppID, migrationType string) int {

	var migrationIDs []int

	stencilDBConn := db.GetDBConn("stencil")

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

	closeDBConn(stencilDBConn)

	return migrationIDs[0]

}

func waitForMigrationComplete(migrationID int, wg *sync.WaitGroup) {

	defer wg.Done()
	
	stencilDBConn := db.GetDBConn("stencil")

	for !CheckMigrationComplete1(stencilDBConn, migrationID) {
		time.Sleep(CHECK_MIGRATION_COMPLETE_INTERVAL)
	}

	closeDBConn(stencilDBConn)

}

func StartDisplay(uid, srcAppID, dstAppID, migrationType string, 
	threadNum int, wg *sync.WaitGroup, 
	enableDisplay, displayInFirstPhase bool) {

	migrationID := waitGetMigrationID(uid, srcAppID, dstAppID, migrationType)

	if enableDisplay {
		displayController(migrationID, threadNum, wg, displayInFirstPhase)
	} else {
		waitForMigrationComplete(migrationID, wg)
	}

}
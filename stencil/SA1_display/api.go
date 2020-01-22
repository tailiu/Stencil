package SA1_display

import (
	"log"
	"sync"
)

func displayController(migrationID, threadNum int) {

	var wg sync.WaitGroup

	wg.Add(threadNum)

	// If the destination app database is not in the new server, newDB is false
	newDB := false

	// If the display controller needs to resolve references, resolveReference is true
	resolveReference := true

	displayConfig := CreateDisplayConfig(migrationID, resolveReference, newDB, &wg)

	log.Println("############### Start Display Controller ###############")

	log.Println("Migration ID:", migrationID)

	log.Println("Total Display Thread(s):", threadNum)

	for i := 0; i < threadNum; i++ {

		log.Println("Start thread:", i + 1)

		go DisplayThread(displayConfig)

	}

	wg.Wait()

	log.Println("############### End Display Controller ###############")

}

func StartDisplay(uid, srcAppID, dstAppID, migrationType string, threadNum int) {

	migrationIDs := getMigrationIDs(uid, srcAppID, dstAppID, migrationType)

	if len(migrationIDs) != 1 {
		log.Fatal("There are more than one migration")
	}

	displayController(migrationIDs[0], threadNum)

}
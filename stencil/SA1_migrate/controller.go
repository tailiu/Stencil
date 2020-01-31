package SA1_migrate

import (
	"stencil/apis"
	"stencil/SA1_display"
	"sync"
	"log"
)

func Controller(uid, srcAppName, srcAppID, 
	dstAppName, dstAppID, migrationType string, threadNum int, 
	enableDisplay, displayInFirstPhase bool) {
	
	var wg sync.WaitGroup
	
	if enableDisplay {
		log.Println("############### Start Migration and Display Controller ###############")
	} else {
		log.Println("############### Start Migration Controller ###############")
	}
	

	// Instead of waiting for all display threads to finish,
	// we only need to wait for one display thread to finish
	wg.Add(1)

	go apis.StartMigration(uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType)

	go SA1_display.StartDisplay(uid, srcAppID, dstAppID, migrationType, 
		threadNum, &wg, enableDisplay, displayInFirstPhase)

	wg.Wait()

	if enableDisplay {
		log.Println("############### End Migration and Display Controller ###############")
	} else {
		log.Println("############### End Migration Controller ###############")
	}

}
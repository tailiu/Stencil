package SA2_migrate

import (
	"log"
	"stencil/SA2_display"
	"stencil/apis"
	"sync"
)

// By default, enableDisplay, displayInFirstPhase, markAsDelete
// are true, true, false 
func handleArgs(args []bool) (bool, bool, bool) {

	enableDisplay, displayInFirstPhase, enableBags := true, true, false
	
	for i, arg := range args {
		switch i {
		case 0:
			enableDisplay = arg
		case 1:
			displayInFirstPhase = arg
		case 2:
			enableBags = arg
		default:
			log.Fatal(`The input args of the migration and display controller 
				do not satisfy requirements!`)
		}
	}

	return enableDisplay, displayInFirstPhase, enableBags

}

func Controller(uid, srcAppName, srcAppID,
	dstAppName, dstAppID, migrationType string,
	threadNum int, args ...bool) {
	
	enableDisplay, displayInFirstPhase, enableBags := handleArgs(args)
	
	var wg sync.WaitGroup

	if enableDisplay {
		log.Println("############### Start SA2 Migration and Display Controller ###############")
	} else {
		log.Println("############### Start SA2 Migration Controller ###############")
	}

	// Instead of waiting for all display threads to finish,
	// we only need to wait for one display thread to finish
	wg.Add(1)

	go apis.StartMigration(uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, enableBags)

	go SA2_display.StartDisplay(
		uid, srcAppID, dstAppID, migrationType, threadNum, &wg, 
		enableDisplay, displayInFirstPhase,
	)

	wg.Wait()

	if enableDisplay {
		log.Println("############### End SA2 Migration and Display Controller ###############")
	} else {
		log.Println("############### End SA2 Migration Controller ###############")
	}

}
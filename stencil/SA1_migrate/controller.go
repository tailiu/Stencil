package SA1_migrate

import (
	"log"
	"stencil/SA1_display"
	"stencil/apis"
	"sync"
)

// By default, enableDisplay, displayInFirstPhase, markAsDelete
// are true, true, false 
func handleArgs(args []bool) (bool, bool, bool, bool, bool, bool, bool) {

	enableDisplay, displayInFirstPhase, markAsDelete, useBladeServerAsDst, enableBags, enableFTP, resolveRefs :=
		true, true, false, true, false, true, true
	
	for i, arg := range args {
		switch i {
		case 0:
			enableDisplay = arg
		case 1:
			displayInFirstPhase = arg
		case 2:
			markAsDelete = arg
		case 3:
			useBladeServerAsDst = arg
		case 4:
			enableBags = arg
		case 5:
			enableFTP = arg
		case 6:
			resolveRefs = arg
		default:
			log.Fatal(`The input args of the migration and display controller do not satisfy requirements!`)
		}
	}

	return enableDisplay, displayInFirstPhase, markAsDelete, useBladeServerAsDst, enableBags, enableFTP, resolveRefs

}

func Controller(uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType string, threadNum int, args ...bool) {
	
	enableDisplay, displayInFirstPhase, markAsDelete, useBladeServerAsDst, enableBags, enableFTP, resolveRefs := handleArgs(args)
	
	var wg sync.WaitGroup

	if enableDisplay {
		log.Println("############### Start SA1 Migration and Display Controller ###############")
	} else {
		log.Println("############### Start SA1 Migration Controller ###############")
	}

	// Instead of waiting for all display threads to finish,
	// we only need to wait for one display thread to finish
	wg.Add(1)

	go apis.StartMigration(
		uid, srcAppName, srcAppID, dstAppName, dstAppID, 
		migrationType, useBladeServerAsDst, enableBags, enableFTP,
	)

	go SA1_display.StartDisplay(
		uid, srcAppID, dstAppID, migrationType, threadNum, &wg, 
		enableDisplay, displayInFirstPhase, markAsDelete, useBladeServerAsDst, resolveRefs,
	)

	wg.Wait()

	if enableDisplay {
		log.Println("############### End SA1 Migration and Display Controller ###############")
	} else {
		log.Println("############### End SA1 Migration Controller ###############")
	}

}

package SA1_migrate

import (
	"log"
	"stencil/SA1_display"
	"stencil/apis"
	"sync"
)

// By default, enableDisplay, displayInFirstPhase, markAsDelete
// are true, true, false
func handleArgs(args []bool) (bool, bool, bool) {

	enableDisplay, displayInFirstPhase, markAsDelete :=
		true, true, false

	for i, arg := range args {
		switch i {
		case 0:
			enableDisplay = arg
		case 1:
			displayInFirstPhase = arg
		case 2:
			markAsDelete = arg
		default:
			log.Fatal(`The input args of the migration and display controller 
				do not satisfy requirements!`)
		}
	}

	return enableDisplay, displayInFirstPhase, markAsDelete

}

func Controller(uid, srcAppName, srcAppID,
	dstAppName, dstAppID, migrationType string,
	threadNum int, args ...bool) {

	enableDisplay, displayInFirstPhase, markAsDelete := handleArgs(args)

	useBladeServerAsDst := true

	var wg sync.WaitGroup

	if enableDisplay {
		log.Println("############### Start Migration and Display Controller ###############")
	} else {
		log.Println("############### Start Migration Controller ###############")
	}

	// Instead of waiting for all display threads to finish,
	// we only need to wait for one display thread to finish
	wg.Add(1)

	go apis.StartMigration(uid, srcAppName, srcAppID,
		dstAppName, dstAppID, migrationType, useBladeServerAsDst, false)

	go SA1_display.StartDisplay(uid, srcAppID, dstAppID, migrationType,
		threadNum, &wg, enableDisplay, displayInFirstPhase, markAsDelete)

	wg.Wait()

	if enableDisplay {
		log.Println("############### End Migration and Display Controller ###############")
	} else {
		log.Println("############### End Migration Controller ###############")
	}

}

func Controller2(uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType string, threadNum int, useBladeServerAsDst, enableDisplay, enableBags bool) {

	displayInFirstPhase, markAsDelete := false, false

	var wg sync.WaitGroup

	if enableDisplay {
		log.Println("############### Start Migration and Display Controller ###############")
	} else {
		log.Println("############### Start Migration Controller ###############")
	}

	// Instead of waiting for all display threads to finish,
	// we only need to wait for one display thread to finish
	wg.Add(1)

	go apis.StartMigration(uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, useBladeServerAsDst, enableBags)

	go SA1_display.StartDisplay(uid, srcAppID, dstAppID, migrationType, threadNum, &wg, enableDisplay, displayInFirstPhase, markAsDelete)

	wg.Wait()

	if enableDisplay {
		log.Println("############### End Migration and Display Controller ###############")
	} else {
		log.Println("############### End Migration Controller ###############")
	}

}

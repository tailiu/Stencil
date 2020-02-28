package SA1_display

import (
	// "stencil/config"
	"stencil/db"
	"log"
	"sync"
	"time"
)

const WAIT_FOR_MIGRATION_START_INTERVAL = 100 * time.Millisecond

const CHECK_MIGRATION_COMPLETE_INTERVAL1 = time.Second

const CHECK_MIGRATION_COMPLETE_INTERVAL2 = 500 * time.Millisecond

func displayController(migrationID, threadNum int, wg *sync.WaitGroup, 
	displayInFirstPhase, markAsDelete, useBladeServerAsDst bool) {
	
	deletionHoldEnable := false
	
	stencilDBConn := db.GetDBConn("stencil")
	defer stencilDBConn.Close()

	if !displayInFirstPhase {
		for !CheckMigrationComplete(stencilDBConn) {
			time.Sleep(CHECK_MIGRATION_COMPLETE_INTERVAL2)
		}
	}

	log.Println("Migration ID:", migrationID)

	log.Println("Total Display Thread(s):", threadNum)
	
	logDisplayStartTime(stencilDBConn, migrationID)

	for i := 0; i < threadNum; i++ {

		log.Println("Start Display Thread:", i + 1)

		go func(migrationID int, deletionHoldEnable bool) {

			defer wg.Done()
			
			DisplayThread(migrationID, deletionHoldEnable)
		
		} (migrationID, deletionHoldEnable)

	}

	// log.Println("############### End Display Controller ###############")

}

func waitToGetMigrationID(uid, srcAppID, dstAppID, migrationType string) int {

	var migrationIDs []int

	stencilDBConn := db.GetDBConn("stencil")

	for {

		migrationIDs = getMigrationIDs(stencilDBConn, 
			uid, srcAppID, dstAppID, migrationType)
		
		// log.Println("*******")
		// log.Println(migrationIDs)
		
		// log.Println("user ID:", uid)

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

	for !CheckMigrationComplete(stencilDBConn, migrationID) {
		time.Sleep(CHECK_MIGRATION_COMPLETE_INTERVAL1)
	}

	closeDBConn(stencilDBConn)

}

func handlArgs(args []bool) (bool, bool, bool, bool) {

	enableDisplay, displayInFirstPhase, markAsDelete, useBladeServerAsDst :=
		true, true, false, true
	
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
		default:
			log.Fatal(`The input args of the display controller 
				do not satisfy requirements!`)
		}
	}

	return enableDisplay, displayInFirstPhase, markAsDelete, useBladeServerAsDst

}

func StartDisplay(uid, srcAppID, dstAppID, 
	migrationType string, 
	threadNum int, wg *sync.WaitGroup, 
	args ...bool) {

	migrationID := waitToGetMigrationID(uid, srcAppID, dstAppID, migrationType)

	log.Println("Migration ID Found by Display:", migrationID)

	enableDisplay, displayInFirstPhase, markAsDelete, useBladeServerAsDst := handlArgs(args)

	if enableDisplay {
		displayController(
			migrationID, threadNum, wg, 
			displayInFirstPhase, markAsDelete,
			useBladeServerAsDst,
		)
	} else {
		waitForMigrationComplete(migrationID, wg)
	}

}
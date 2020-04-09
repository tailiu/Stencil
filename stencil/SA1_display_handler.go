package main

import (
	"fmt"
	"log"
	"stencil/SA1_display"
	"sync"
)

func test1() {

	threadNum := 1

	migrationID := 754595238

	// If the destination app database is not in the new server, newDB is false
	newDB := false

	// If the display controller needs to resolve references, resolveReference is true
	resolveReference := true

	displayInFirstPhase := true

	displayConfig := SA1_display.CreateDisplayConfig(migrationID,
		resolveReference, newDB, displayInFirstPhase)

	log.Println("Migration ID:", migrationID)

	for i := 0; i < threadNum; i++ {

		go displayConfig.DisplayThread()

	}

	for {

		fmt.Scanln()

	}

}

// func test2() {

// 	threadNum := 1

// 	migrationID := 800666262

// 	SA1_display.StartDisplay(migrationID, threadNum)

// }

func test3() {

	var wg sync.WaitGroup

	log.Println("############### Start Display Controller ###############")

	// Instead of waiting for all display threads to finish,
	// we only need to wait for one display thread to finish
	wg.Add(1)

	uid, srcAppID, dstAppID, migrationType, threadNum :=
		"24214", "1", "2", "d", 1

	enableDisplay, displayInFirstPhase := true, true

	// SA1_display.StartDisplay(uid, srcAppID, dstAppID, migrationType, threadNum, nil)

	go SA1_display.StartDisplay(uid, srcAppID, dstAppID, migrationType,
		threadNum, &wg, enableDisplay, displayInFirstPhase)

	wg.Wait()

	log.Println("############### End Display Controller ###############")

}

func main() {

	// test1()
	// test2()
	test3()
}

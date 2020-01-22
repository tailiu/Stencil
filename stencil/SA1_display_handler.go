package main

import (
	"stencil/SA1_display"
	"fmt"
	"log"
)

func test1() {

	threadNum := 1
	
	migrationID := 1370370281

	// If the destination app database is not in the new server, newDB is false
	newDB := false

	// If the display controller needs to resolve references, resolveReference is true
	resolveReference := true

	displayConfig := SA1_display.CreateDisplayConfig(migrationID, resolveReference, newDB, nil)

	log.Println("Migration ID:",migrationID)

	for i := 0; i < threadNum; i++ {

		go SA1_display.DisplayThread(displayConfig)

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

	uid, srcAppID, dstAppID, migrationType, threadNum := 
		"44773", "1", "2", "d", 1

	SA1_display.StartDisplay(uid, srcAppID, dstAppID, migrationType, threadNum)

}

func main() {
	
	// test2()
	test3()
}

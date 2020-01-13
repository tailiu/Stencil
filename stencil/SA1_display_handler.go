package main

import (
	"stencil/SA1_display"
	"fmt"
)

func main() {
	
	threadNum := 1
	
	migrationID := 700588914

	// If the destination app database is not in the new server, newDB is false
	newDB := false

	// If the display controller needs to resolve references, resolveReference is true
	resolveReference := true

	displayConfig := SA1_display.CreateDisplayConfig(migrationID, resolveReference, newDB)

	for i := 0; i < threadNum; i++ {

		go SA1_display.DisplayThread(displayConfig)

	}

	for {

		fmt.Scanln()

	}
}

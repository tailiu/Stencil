package main

import (
	"stencil/app_display"
	"stencil/app_display_algorithm"
	"fmt"
)

func main() {
	
	threadNum := 1
	
	migrationID := 955012936

	// If the destination app database is not in the new server, newDB is false
	newDB := false

	// If the display controller needs to resolve references, resolveReference is true
	resolveReference := true

	displayConfig := app_display.CreateDisplayConfig(migrationID, resolveReference, newDB)

	for i := 0; i < threadNum; i++ {

		go app_display_algorithm.DisplayThread(displayConfig)

	}

	for {

		fmt.Scanln()

	}
}

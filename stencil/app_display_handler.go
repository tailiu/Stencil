package main

import (
	"stencil/app_display"
	"stencil/app_display_algorithm"
	"fmt"
)

func main() {
	
	threadNum := 1
	
	dstApp := "mastodon"
	// dstApp := "diaspora"
	
	migrationID := 2124890507

	newDB := false

	displayConfig := app_display.CreateDisplayConfig(dstApp, migrationID, newDB)

	for i := 0; i < threadNum; i++ {

		go app_display_algorithm.DisplayThread(displayConfig)

	}

	for {

		fmt.Scanln()

	}
}

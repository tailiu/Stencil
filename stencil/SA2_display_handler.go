package main

import (
	"stencil/SA2_display"
	"fmt"
)

func main() {

	threadNum := 1
	
	dstApp := "mastodon"
	// dstApp := "diaspora"
	
	migrationID := 658943258

	deletionHoldEnable := false

	for i := 0; i < threadNum; i++ {
		go SA2_display.DisplayThread(dstApp, migrationID, deletionHoldEnable)
	}

	for {
		fmt.Scanln()
	}
}

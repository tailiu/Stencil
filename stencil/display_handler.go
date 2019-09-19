package main

import (
	"stencil/display_algorithm"
	"fmt"
)

func main() {
	threadNum := 1
	// dstApp := "mastodon"
	dstApp := "mastodon"
	migrationID := 1923587293

	deletionHoldEnable := false

	for i := 0; i < threadNum; i++ {
		go display_algorithm.DisplayThread(dstApp, migrationID, deletionHoldEnable)
	}

	for {
		fmt.Scanln()
	}
}

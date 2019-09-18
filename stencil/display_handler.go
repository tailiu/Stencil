package main

import (
	"stencil/display_algorithm"
	"fmt"
)

func main() {
	threadNum := 5
	dstApp := "mastodon"
	migrationID := 859988132

	deletionHoldEnable := false

	for i := 0; i < threadNum; i++ {
		go display_algorithm.DisplayThread(dstApp, migrationID, deletionHoldEnable)
	}

	for {
		fmt.Scanln()
	}
}

package main

import (
	"stencil/app_display_algorithm"
	"fmt"
)

func main() {
	threadNum := 1
	dstApp := "mastodon"
	// dstApp := "diaspora"
	migrationID := 228970122

	for i := 0; i < threadNum; i++ {
		go app_display_algorithm.DisplayThread(dstApp, migrationID)
	}

	for {
		fmt.Scanln()
	}
}

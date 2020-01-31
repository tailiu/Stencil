package main

import (
	"stencil/SA1_migrate"
)

// The diaspora database needs to be changed to diaspora_test for testing
func main() {
	
	uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
		"1019", "diaspora", "1", "mastodon", "2", "d", 1

	// If enableDisplay is set to be true, then display threads will be started
	// If displayInFirstPhase is set to be true, 
	// then display threads will check in the first phase
	enableDisplay, displayInFirstPhase := true, true
	
	SA1_migrate.Controller(uid, srcAppName, srcAppID, 
		dstAppName, dstAppID, migrationType, threadNum,
		enableDisplay, displayInFirstPhase,
	)

}
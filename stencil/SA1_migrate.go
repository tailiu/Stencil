package main

import (
	"stencil/SA1_migrate"
)

// The diaspora database needs to be changed to diaspora_test for testing
func main() {
<<<<<<< HEAD
	
	uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
=======

	uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum :=
>>>>>>> d30c473b8c03ce16f10ee415b4d45e13f6499e2c
		"162610", "diaspora", "1", "mastodon", "2", "d", 1

	// If enableDisplay is set to be true, then display threads will be started
	// If displayInFirstPhase is set to be true,
	// then display threads will check in the first phase
	enableDisplay, displayInFirstPhase := true, false

	SA1_migrate.Controller(uid, srcAppName, srcAppID,
		dstAppName, dstAppID, migrationType, threadNum,
		enableDisplay, displayInFirstPhase,
	)

}

package main

import (
	"stencil/SA1_migrate"
)

// The diaspora database needs to be changed to diaspora_test for testing
func main() {

	uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum :=
		"1001", "diaspora", "1", "mastodon", "2", "d", 1

	// If enableDisplay is set to be true, then display threads will be started
	// If displayInFirstPhase is set to be true,
	// then display threads will check in the first phase
	// If markAsDelete is set to be true,
	// the display threads will mark data to be put into data bags
	// as delete instead of deleting data
	enableDisplay, displayInFirstPhase, markAsDelete := true, true, false

	SA1_migrate.Controller(uid, srcAppName, srcAppID,
		dstAppName, dstAppID, migrationType, threadNum,
		enableDisplay, displayInFirstPhase, markAsDelete,
	)

}

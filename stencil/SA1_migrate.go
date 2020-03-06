package main

import (
	"stencil/SA1_migrate"
	"stencil/db"
)

// The diaspora database needs to be changed to diaspora_test for testing
func main() {

	db.STENCIL_DB = "stencil_test2"
	// db.DIASPORA_DB = "diaspora_test2"
	db.DIASPORA_DB = "diaspora_1000000_exp6"
	db.MASTODON_DB = "mastodon_test1"

	// uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum :=
	// 	"1017", "diaspora", "1", "mastodon", "2", "d", 1

	uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum :=
		"884", "diaspora", "1", "mastodon", "2", "i", 1

	// If enableDisplay is set to be true, then display threads will be started
	// If displayInFirstPhase is set to be true,
	// then display threads will check in the first phase
	// If markAsDelete is set to be true,
	// the display threads will mark data to be put into data bags
	// as delete instead of deleting data

	enableDisplay, displayInFirstPhase, markAsDelete, useBladeServerAsDst, enableBags :=
		true, true, false, false, true

	SA1_migrate.Controller(uid, srcAppName, srcAppID,
		dstAppName, dstAppID, migrationType, threadNum,
		enableDisplay, displayInFirstPhase, markAsDelete,
		useBladeServerAsDst, enableBags,
	)

}

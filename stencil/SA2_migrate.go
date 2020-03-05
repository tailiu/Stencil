package main

import (
	// "stencil/apis"
	"stencil/db"
	"stencil/SA2_migrate"
)

func main() {
	
	// db.STENCIL_DB = "stencil_exp_sa2_100k"

	// uid, srcApp, srcAppID, dstApp, dstAppID, migrationType, enableBags :=
	// 	"4061", "diaspora", "1", "mastodon", "2", "d", false

	// apis.StartMigrationSA2(uid, srcApp, srcAppID, dstApp, dstAppID, migrationType, enableBags)

	db.STENCIL_DB = "stencil_exp_sa2_1k_exp"

	uid, srcApp, srcAppID, dstApp, dstAppID, migrationType, threadNum :=
		"41", "diaspora", "1", "mastodon", "2", "d", 1

	enableDisplay, displayInFirstPhase, enableBags := true, true, true

	SA2_migrate.Controller(
		uid, srcApp, srcAppID, dstApp, dstAppID, 
		migrationType, threadNum, enableDisplay, 
		displayInFirstPhase, enableBags,
	)

}
package main

import (
	"stencil/apis"
	"stencil/db"
)

func main() {
	
	db.STENCIL_DB = "stencil_exp_sa2_1k_exp"

	uid, srcApp, srcAppID, dstApp, dstAppID, migrationType, enableBags :=
		"140", "diaspora", "1", "mastodon", "2", "d", false

	apis.StartMigrationSA2(uid, srcApp, srcAppID, dstApp, dstAppID, migrationType, enableBags)

}
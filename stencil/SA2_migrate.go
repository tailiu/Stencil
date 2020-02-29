package main

import (
	"stencil/apis"
	"stencil/db"
)

func main() {
	
	db.STENCIL_DB = "stencil_exp_sa2_100k"

	uid, srcApp, srcAppID, dstApp, dstAppID, migrationType, enableBags :=
		"4061", "diaspora", "1", "mastodon", "2", "d", false

	apis.StartMigrationSA2(uid, srcApp, srcAppID, dstApp, dstAppID, migrationType, enableBags)

}
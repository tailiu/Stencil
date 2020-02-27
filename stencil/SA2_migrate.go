package main

import (
	"stencil/apis"
	"stencil/db"
)

func main() {
	
	db.STENCIL_DB = "stencil_exp_sa2_10k_exp"
	db.DIASPORA_DB = "diaspora_test2"
	db.MASTODON_DB = "mastodon_test2"

	uid, srcApp, srcAppID, dstApp, dstAppID, migrationType, enableBags :=
		"44854", "diaspora", "1", "mastodon", "2", "d", false

	apis.StartMigrationSA2(uid, srcApp, srcAppID, dstApp, dstAppID, migrationType, enableBags)

}
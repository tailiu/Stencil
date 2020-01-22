package main

import (
	"stencil/apis"
	"stencil/SA1_display"
)

func main() {
	
	uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType := 
		"44781", "diaspora", "1", "mastodon", "2", "d"

	apis.StartMigration(uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType)

	SA1_display.StartDisplay()

}
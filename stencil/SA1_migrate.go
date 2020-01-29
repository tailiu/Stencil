package main

import (
	"stencil/SA1_migrate"
)

// The diaspora database needs to be changed to diaspora_test for testing
func main() {
	
	uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
		"10190", "diaspora", "1", "mastodon", "2", "d", 1

	SA1_migrate.Controller(uid, srcAppName, srcAppID, 
		dstAppName, dstAppID, migrationType, threadNum)

}
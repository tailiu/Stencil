package main

import (
	"stencil/SA1_migrate"
)

func main() {
	
	uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
		"447932", "diaspora", "1", "mastodon", "2", "d", 1

	SA1_migrate.Controller(uid, srcAppName, srcAppID, 
		dstAppName, dstAppID, migrationType, threadNum)

}
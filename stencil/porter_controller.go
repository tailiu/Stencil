package main

import (
	"stencil/apis"
	"stencil/db"
)

func main() {
	appName, appID, table := "diaspora", "1", "posts"
	apis.Port(appName, appID, table, 100, 200, db.GetDBConn(appName), db.GetDBConn(db.STENCIL_DB))
}

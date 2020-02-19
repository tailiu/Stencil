package main

import (
	"stencil/apis"
	"stencil/db"
)

func main() {
	for _, table := range []string{"photos", "posts", "likes", "comments", "conversations", "messages", "contacts", "notifications", "users", "people", "profiles"} {
		appName, appID, table := "diaspora", "1", table
		apis.Port(appName, appID, table, 1000, 200, db.GetDBConn(appName), db.GetDBConn(db.STENCIL_DB))
	}
}

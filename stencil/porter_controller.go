package main

import (
	"stencil/apis"
)

func main() {
	appName, appID, table := "diaspora", "1", "profiles"
	apis.Port(appName, appID, table, 100)
}

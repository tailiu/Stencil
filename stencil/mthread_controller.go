package main

import migrate "stencil/migrate_v2"

func main() {

	mtController := migrate.MigrationThreadController{
		UID:        "54123",
		MType:      "d",
		SrcAppInfo: migrate.App{Name: "diaspora", ID: 1},
		DstAppInfo: migrate.App{Name: "mastodon", ID: 2},
	}

	mtController.Init()
	mtController.Run()
	mtController.Stop()
}

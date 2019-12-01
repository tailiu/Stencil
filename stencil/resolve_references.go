package main

import (
	"stencil/reference_resolution"
	"stencil/app_display"
)

func main() { 
	migrationID := 908913181

	// StencilDBName := "stencil"
	// stencilDBConn := db.GetDBConn(StencilDBName)

	app := "mastodon"
	// app_id := db.GetAppIDByAppName(stencilDBConn, app)

	// appConfig, err := config.CreateAppConfigDisplay(app, app_id, stencilDBConn, true)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	displayConfig := app_display.CreateDisplayConfig(app, migrationID, true)

	var hint = app_display.HintStruct{
		Table:		"favourites",
		TableID:	"72",
		KeyVal:		map[string]int{"id":24},
	}

	// reference_resolution.ResolveReferenceByBackTraversal(stencilDBConn, &appConfig, migrationID, &hint)
	reference_resolution.ResolveReference(displayConfig, &hint)
}

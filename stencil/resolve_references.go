package main

import (
	// "stencil/reference_resolution"
	// "stencil/app_display"
	"log"
	"stencil/config"
)

func main() { 
	// migrationID := 908913181

	// // StencilDBName := "stencil"
	// // stencilDBConn := db.GetDBConn(StencilDBName)

	// app := "mastodon"
	// // app_id := db.GetAppIDByAppName(stencilDBConn, app)

	// // appConfig, err := config.CreateAppConfigDisplay(app, app_id, stencilDBConn, true)
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }

	// displayConfig := app_display.CreateDisplayConfig(app, migrationID, true)

	// var hint = app_display.HintStruct{
	// 	Table:		"favourites",
	// 	TableID:	"72",
	// 	KeyVal:		map[string]int{"id":24},
	// }

	// // reference_resolution.ResolveReference(displayConfig, &hint)

	schemaMappings, err := config.LoadSchemaMappings()
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(schemaMappings["mastodon"])
	fromApp, fromTable, fromAttr, toApp, toTable := "diaspora", "posts", "posts.id", "mastodon", "statuses"

	for _, mappings := range schemaMappings.AllMappings {
		if mappings.FromApp == fromApp {
			for _, app := range mappings.ToApps {
				if app.Name == toApp {
					for _, mapping := range app.Mappings {
						for _, fTable := range mapping.FromTables {
							if fTable == fromTable {
								for _, tTable := range mapping.ToTables {
									if tTable.Table == toTable {
										// log.Println(tTable)
										for tAttr, fAttr := range tTable.Mapping {
											if fAttr == fromAttr {
												log.Println(tAttr)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

}

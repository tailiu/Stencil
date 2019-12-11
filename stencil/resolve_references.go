package main

import (
	"stencil/reference_resolution"
	"stencil/app_display"
	"stencil/config"
	"log"
)

/**
 *
 * Diaspora -> Mastodon
 * 
 * Identity:
 * 	a like (id:12) in Diaspora likes table (id:8) -> a favourite (id:235893) in Mastodon favourite table (id:72)
 * 	a post (id:40) in Diaspora posts table (id:37) -> a status (id:21778) in Mastodon statuses table (id:92)
 *
 * Reference:
 *	Diaspora (1), likes (id:8), like (id:12), target_id -> Diaspora (1), posts (id:37), post (id:40), id
 *
 *
**/

func test1(displayConfig *config.DisplayConfig) {

	var hint = app_display.HintStruct{
		Table:		"favourites",
		TableID:	"72",
		KeyVal:		map[string]int{"id":235893},
	}

	updatedAttrs := reference_resolution.ResolveReference(displayConfig, &hint)

	log.Println(updatedAttrs)

}


/**
 *
 * Diaspora -> Mastodon
 * 
 * Identity:
 * 	a like (id:15) in Diaspora likes table (id:8) -> a favourite (id:235977) in Mastodon favourite table (id:72)
 * 	a post (id:50) in Diaspora posts table (id:37) -> a status (id:79738) in Mastodon statuses table (id:92)
 *
 * Reference:
 *	Diaspora (1), likes table (id:8), like (id:15), target_id -> Diaspora (1), posts table (id:37), post (id:50), id
 *
**/

func test2(displayConfig *config.DisplayConfig) {

	var hint = app_display.HintStruct{
		Table:		"statuses",
		TableID:	"92",
		KeyVal:		map[string]int{"id":79738},
	}

	updatedAttrs := reference_resolution.ResolveReference(displayConfig, &hint)

	log.Println(updatedAttrs)

}

func main() { 

	migrationID := 434969759

	newDB := false

	resolveReference := true

	displayConfig := app_display.CreateDisplayConfig(migrationID, resolveReference, newDB)
	
	// test1(displayConfig)

	test2(displayConfig)
	
}

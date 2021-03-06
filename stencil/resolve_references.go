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
 * 	a like (id:101) in Diaspora likes table (id:8) 
 * 		-> a favourite (id:235970) in Mastodon favourite table (id:72)
 * 	a post (id:201) in Diaspora posts table (id:37) 
 *		-> a status (id:79744) in Mastodon statuses table (id:92)
 *
 * Reference:
 *	Diaspora (1), likes (id:8), like (id:101), target_id 
 * 		-> Diaspora (1), posts (id:37), post (id:201), id
 *
 *
**/

func test1(display *config.DisplayConfig) {

	var hint = &app_display.HintStruct{
		Table:		"favourites",
		TableID:	"72",
		KeyVal:		map[string]int{"id":235970},
	}

	ID := hint.TransformHintToIdenity(display)

	myUpdatedAttrs, othersUpdatedAttrs := 
		reference_resolution.ResolveReference(display, ID)

	log.Println(myUpdatedAttrs, othersUpdatedAttrs)

}


/**
 *
 * Diaspora -> Mastodon
 * 
 * Identity:
 * 	a like (id:25) in Diaspora likes table (id:8) 
 *		-> a favourite (id:235959) in Mastodon favourite table (id:72)
 * 	a post (id:70) in Diaspora posts table (id:37) 
 *		-> a status (id:21783) in Mastodon statuses table (id:92)
 *
 * Reference:
 *	Diaspora (1), likes table (id:8), like (id:25), target_id 
 *		-> Diaspora (1), posts table (id:37), post (id:70), id
 *
**/

func test2(display *config.DisplayConfig) {

	var hint = &app_display.HintStruct{
		Table:		"statuses",
		TableID:	"92",
		KeyVal:		map[string]int{"id":21783},
	}

	ID := hint.TransformHintToIdenity(display)

	myUpdatedAttrs, othersUpdatedAttrs := 
		reference_resolution.ResolveReference(display, ID)

	log.Println(myUpdatedAttrs, othersUpdatedAttrs)

}

func test3() {

	

}

func main() { 

	migrationID := 434969759

	newDB := false
	
	referenceResolutionConfig := reference_resolution.InitializeReferenceResolution()

	test1(referenceResolutionConfig)

	// test2(display)
	
}

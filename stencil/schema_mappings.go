package main

import (
	"log"
	"stencil/config"
	"stencil/schema_mappings"
	// "stencil/SA1_display"
)


// func test1(display *DisplayConfig) {
	
// 	// fromApp, fromTable, fromAttr, toApp, toTable := 
// 		// "diaspora", "posts", "posts.id", "mastodon", "statuses"
// 	// fromApp, fromTable, fromAttr, toApp, toTable, ignoreREF := 
// 		// 	"diaspora", "comments", "comments.commentable_id", "mastodon", "statuses", false
// 	fromApp, fromTable, fromAttr, toApp, toTable, ignoreREF := 
// 		"diaspora", "posts", "posts.id", "mastodon", "status_stats", false
	
// 	attr, _ := schema_mappings.GetMappedAttributesFromSchemaMappings(display.AllMappings,
// 		fromApp, fromTable, fromAttr, toApp, toTable, ignoreREF)

// 	log.Println(attr)

// }

// func test2(display *DisplayConfig) {

// 	// toTable, toAttr := "accounts", "id"
// 	// toTable, toAttr := "users", "account_id"
// 	toTable, toAttr := "statuses", "in_reply_to_id"

// 	exists, err := schema_mappings.REFExists(display.MappingsToDst, toTable, toAttr)
// 	if err != nil {
// 		log.Println(err)
// 	} else {
// 		log.Println(exists)
// 	}

// }

// func test3(display *DisplayConfig) {

// 	toTable := "statuses"

// 	attrs := schema_mappings.GetAllMappedAttributesContainingREFInMappings(
// 		display.MappingsToDst, toTable)
	
// 	log.Println(attrs)

// }

// func test4(display *config.DisplayConfig) {

// 	fromApp, fromAttr, toApp, toTable := 
// 		"diaspora", "posts.id", "mastodon", "media_attachments"

// 	attrs, err := schema_mappings.GetMappedAttributesFromSchemaMappingsByFETCH(
// 		display.AllMappings, fromApp, fromAttr, toApp, toTable)
	
// 	if err != nil {

// 		log.Println(err)

// 	} else {

// 		log.Println(attrs)
// 	}
// }

func test5() {

	// pairwiseSchemaMappings, err := schema_mappings.LoadPairwiseSchemaMappings()
	// if err != nil {
	// 	log.Println(err)
	// }

	// log.Println(pairwiseSchemaMappings)

	_, err := schema_mappings.DeriveMappingsByPSM()
	if err != nil {
		log.Println(err)
	}

	// log.Println(schema_mappings)

}

func test6() {

	fromAttr := "#REF(#ASSIGN(messages.id),messages.id)"

	procAttr := schema_mappings.RemoveASSIGNAllRightParenthesesIfExists(fromAttr)

	log.Println(procAttr)

}

func test7() {

	allMappings, err1 := config.LoadSchemaMappings()
	if err1 != nil {
		log.Fatal(err1)
	}

	// app, attr, attrTable, attrToUpdate, attrToUpdateTable := "mastodon", "id", "accounts", "account_id", "statuses"

	app, attr, attrTable, attrToUpdate, attrToUpdateTable := "mastodon", "id", "statuses", "conversation_id", "statuses"

	log.Println(app, attr, attrTable, attrToUpdate, attrToUpdateTable)

	exists := schema_mappings.ReferenceExistsBasedOnMappings(
		allMappings, app, attr, attrTable, attrToUpdate, attrToUpdateTable,
	)

	log.Println(exists)
	
}

func main() {

	// migrationID := 955012936

	// // If the destination app database is not in the new server, newDB is false
	// newDB := false

	// resolveReference := true

	// display := SA1_display.CreateDisplayConfig(migrationID, resolveReference, newDB)

	// test1(display)

	// test2(display)

	// test3(display)

	// test4(display)

	// test5()

	// test6()

	test7()

}
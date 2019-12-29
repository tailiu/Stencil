package main

import (
	"log"
	"stencil/schema_mappings"
	// "stencil/SA1_display"
)


// func test1(displayConfig *DisplayConfig) {
	
// 	// fromApp, fromTable, fromAttr, toApp, toTable := 
// 		// "diaspora", "posts", "posts.id", "mastodon", "statuses"
// 	// fromApp, fromTable, fromAttr, toApp, toTable, ignoreREF := 
// 		// 	"diaspora", "comments", "comments.commentable_id", "mastodon", "statuses", false
// 	fromApp, fromTable, fromAttr, toApp, toTable, ignoreREF := 
// 		"diaspora", "posts", "posts.id", "mastodon", "status_stats", false
	
// 	attr, _ := schema_mappings.GetMappedAttributesFromSchemaMappings(displayConfig.AllMappings,
// 		fromApp, fromTable, fromAttr, toApp, toTable, ignoreREF)

// 	log.Println(attr)

// }

// func test2(displayConfig *DisplayConfig) {

// 	// toTable, toAttr := "accounts", "id"
// 	// toTable, toAttr := "users", "account_id"
// 	toTable, toAttr := "statuses", "in_reply_to_id"

// 	exists, err := schema_mappings.REFExists(displayConfig.MappingsToDst, toTable, toAttr)
// 	if err != nil {
// 		log.Println(err)
// 	} else {
// 		log.Println(exists)
// 	}

// }

// func test3(displayConfig *DisplayConfig) {

// 	toTable := "statuses"

// 	attrs := schema_mappings.GetAllMappedAttributesContainingREFInMappings(
// 		displayConfig.MappingsToDst, toTable)
	
// 	log.Println(attrs)

// }

// func test4(displayConfig *config.DisplayConfig) {

// 	fromApp, fromAttr, toApp, toTable := 
// 		"diaspora", "posts.id", "mastodon", "media_attachments"

// 	attrs, err := schema_mappings.GetMappedAttributesFromSchemaMappingsByFETCH(
// 		displayConfig.AllMappings, fromApp, fromAttr, toApp, toTable)
	
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

	schema_mappings, err := schema_mappings.DeriveMappingsByPSM()
	if err != nil {
		log.Println(err)
	}

	log.Println(schema_mappings)

}

func main() {

	// migrationID := 955012936

	// // If the destination app database is not in the new server, newDB is false
	// newDB := false

	// resolveReference := true

	// displayConfig := SA1_display.CreateDisplayConfig(migrationID, resolveReference, newDB)

	// test1(displayConfig)

	// test2(displayConfig)

	// test3(displayConfig)

	// test4(displayConfig)

	test5()

}
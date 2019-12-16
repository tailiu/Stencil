package main

import (
	"log"
	"stencil/schema_mappings"
	"stencil/app_display"
	"stencil/config"
)


func test1(displayConfig *config.DisplayConfig) {
	
	// fromApp, fromTable, fromAttr, toApp, toTable := 
		// "diaspora", "posts", "posts.id", "mastodon", "statuses"
	// fromApp, fromTable, fromAttr, toApp, toTable, ignoreREF := 
		// 	"diaspora", "comments", "comments.commentable_id", "mastodon", "statuses", false
	fromApp, fromTable, fromAttr, toApp, toTable, ignoreREF := 
		"diaspora", "posts", "posts.id", "mastodon", "status_stats", false
	
	attr, _ := schema_mappings.GetMappedAttributesFromSchemaMappings(
		fromApp, fromTable, fromAttr, toApp, toTable, ignoreREF)

	log.Println(attr)

}

func test2(displayConfig *config.DisplayConfig) {

	// toTable, toAttr := "accounts", "id"
	// toTable, toAttr := "users", "account_id"
	toTable, toAttr := "statuses", "in_reply_to_id"

	exists, err := schema_mappings.REFExists(displayConfig, toTable, toAttr)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(exists)
	}

}

func test3(displayConfig *config.DisplayConfig) {

	toTable := "statuses"

	attrs := schema_mappings.GetAllMappedAttributesContainingREFInMappings(
		displayConfig, toTable)
	
	log.Println(attrs)

}

func main() {

	migrationID := 955012936

	// If the destination app database is not in the new server, newDB is false
	newDB := false

	resolveReference := true

	displayConfig := app_display.CreateDisplayConfig(migrationID, resolveReference, newDB)

	// test1(displayConfig)

	// test2(displayConfig)

	test3(displayConfig)
	
}
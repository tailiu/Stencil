package main

import (
	"log"
	"stencil/schema_mappings"
)

func main() {

	// fromApp, fromTable, fromAttr, toApp, toTable := 
		// "diaspora", "posts", "posts.id", "mastodon", "statuses"
	fromApp, fromTable, fromAttr, toApp, toTable := 
		"diaspora", "comments", "comments.commentable_id", "mastodon", "statuses"

	attr, _ := schema_mappings.GetMappedAttributesFromSchemaMappings(
		fromApp, fromTable, fromAttr, toApp, toTable)

	log.Println(attr)

}
package schema_mapping

import (
	"stencil/config"
	// "database/sql"
	"log"
)

func GetMappedAttributeFromSchemaMappings(displayConfig *config.DisplayConfig, 
		fromApp, fromTable, fromAttr, toApp, toTable string) string {

	schemaMappings, err := config.LoadSchemaMappings()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(schemaMappings)

	return ""
	
}
package schema_mapping

import (
	"stencil/config"
	// "database/sql"
	"log"
)

func GetMappedAttributeFromSchemaMappings(displayConfig *config.DisplayConfig, 
		fromApp, fromTable, fromAttr, toApp, toTable string) string {
	// fromAppID := getAppNameByAppID(stencilDBConn, fromApp)

	// fromTableName := getTableNameByTableID(stencilDBConn, fromTable)

	// fromAttrName := getAttrNameByAttrID(stencilDBConn, fromAttr)

	// toTableName := getTableNameByTableID(stencilDBConn, toTable)

	schemaMappings, err := config.LoadSchemaMappings()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(schemaMappings)

	return ""
}
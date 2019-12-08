package schema_mappings

import (
	"stencil/config"
	// "database/sql"
	"log"
)

func GetMappedAttributeFromSchemaMappings(
		fromApp, fromTable, fromAttr, toApp, toTable string) (string, error) {

	schemaMappings, err := config.LoadSchemaMappings()
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(schemaMappings)
	log.Println(fromApp, fromTable, fromAttr, toApp, toTable)

	for _, mappings := range schemaMappings.AllMappings {
		// fromApp
		if mappings.FromApp == fromApp {
			for _, app := range mappings.ToApps {
				// toApp
				if app.Name == toApp {
					for _, mapping := range app.Mappings {
						for _, fTable := range mapping.FromTables {
							// fromTable
							if fTable == fromTable {
								for _, tTable := range mapping.ToTables {
									// toTable
									if tTable.Table == toTable {
										for tAttr, fAttr := range tTable.Mapping {
											// fromAttr
											if fAttr == fromAttr {
												return tAttr, nil
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

	return "", NoMappedAttrFound
	
}

func handleREF() {
	
}
package schema_mappings

import (
	"stencil/config"
	// "database/sql"
	"strings"
	"log"
)

func GetMappedAttributesFromSchemaMappings(
		fromApp, fromTable, fromAttr, toApp, toTable string, ignoreREF bool) ([]string, error) {

	schemaMappings, err := config.LoadSchemaMappings()
	if err != nil {
		log.Fatal(err)
	}

	var attributes []string

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
										log.Println(tTable)
										for tAttr, fAttr := range tTable.Mapping {
											// fromAttr

											// If not ignore #REF
											if !ignoreREF {
												// If there exists #REF
												if strings.Contains(fAttr, "#REF(") {
													fAttr = getFirstArgFromREF(fAttr)
												}
											}
																						
											if fAttr == fromAttr {
												attributes = append(attributes, tAttr)
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

	if len(attributes) == 0 {
		return nil, NoMappedAttrFound
	} else {
		return attributes, nil
	}
	
}

// Return the first argument of #REF
func getFirstArgFromREF(ref string) string {

	tmp := strings.Split(ref, "#REF(")
	
	return strings.Split(tmp[1], ",")[0]

}
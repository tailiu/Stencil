package schema_mappings

import (
	"stencil/config"
	// "database/sql"
	"strings"
	"log"
)

// Return the first argument of #REF
func getFirstArgFromREF(ref string) string {

	tmp := strings.Split(ref, "#REF(")
	
	return strings.Split(tmp[1], ",")[0]

}

func containREF(data string) bool {

	if strings.Contains(data, "#REF(") {
		
		return true

	} else {

		return false
	}
}

func GetToAppMappings(fromApp, toApp string) (config.MappedApp, error) {

	schemaMappings, err := config.LoadSchemaMappings()
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(schemaMappings)

	for _, mappings := range schemaMappings.AllMappings {
		// fromApp
		if mappings.FromApp == fromApp {
			for _, app := range mappings.ToApps {
				// toApp
				if app.Name == toApp {

					return app, nil

				}
			}
		}
	}

	return config.MappedApp{}, MappingsToAppNotFound
}

func GetMappedAttributesFromSchemaMappings(
		fromApp, fromTable, fromAttr, toApp, toTable string, ignoreREF bool) ([]string, error) {

	var attributes []string

	log.Println(fromApp, fromTable, fromAttr, toApp, toTable)

	toAppMappings, err := GetToAppMappings(fromApp, toApp)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, mapping := range toAppMappings.Mappings {
		for _, fTable := range mapping.FromTables {
			// fromTable
			if fTable == fromTable {
				for _, tTable := range mapping.ToTables {
					// toTable
					if tTable.Table == toTable {
						// log.Println(tTable)
						for tAttr, fAttr := range tTable.Mapping {
							// fromAttr

							// If not ignore #REF
							if !ignoreREF {

								// If there exists #REF
								if containREF(fAttr) {
									fAttr = getFirstArgFromREF(fAttr)

									if fAttr == fromAttr {
										attributes = append(attributes, tAttr)
									}
								}

							} else {

								// If there does not exist #REF
								if !containREF(fAttr) {
									
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

	if len(attributes) == 0 {
		return nil, NoMappedAttrFound
	} else {
		return attributes, nil
	}
}

func REFExists(displayConfig *config.DisplayConfig, toTable, toAttr string) (bool, error) {

	log.Println(toTable, toAttr)

	for _, mapping := range displayConfig.MappingsToDst.Mappings {
		// toTable
		for _, tTable := range mapping.ToTables {
			if tTable.Table == toTable {
				// toAttr
				if mappedAttr, ok := tTable.Mapping[toAttr]; ok {
					// log.Println(mappedAttr)

					// If there exists #REF
					if containREF(mappedAttr) {	

						return true, nil

					} else {

						return false, nil 
					}
				} else {

					return false, NoMappedAttrFound
				}
			}
		}
	}

	return false, NoMappedAttrFound

}
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

func GetToAppMappings(allMappings *config.SchemaMappings, 
	fromApp, toApp string) (*config.MappedApp, error) {

	// log.Println(allMappings)

	for _, mappings := range allMappings.AllMappings {

		// fromApp
		if mappings.FromApp == fromApp {
			for _, app := range mappings.ToApps {

				// toApp
				if app.Name == toApp {

					return &app, nil

				}
			}
		}
	}

	return &config.MappedApp{}, MappingsToAppNotFound
}

func GetMappedAttributesFromSchemaMappings(allMappings *config.SchemaMappings,
		fromApp, fromTable, fromAttr, toApp, toTable string, ignoreREF bool) ([]string, error) {

	var attributes []string

	log.Println(fromApp, fromTable, fromAttr, toApp, toTable)

	toAppMappings, err := GetToAppMappings(allMappings, fromApp, toApp)
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

// Here we need to check each from-and-to tables pair.
// For example, there are two #REFs in the mappings 
// from Diaspora.Posts to Mastodon.Statuses, while there are three #REFs in 
// the mappings from Diaspora.Comments to Mastodon.Statuses.
// REFExists will check all these possiblities
func REFExists(mappings *config.MappedApp, toTable, toAttr string) (bool, error) {

	// log.Println(toTable, toAttr)

	for _, mapping := range mappings.Mappings {
		
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
				} 

				// These are commented to take into account the cases mentioned above
				// else {
				// 	return false, NoMappedAttrFound
				// }
			}
		}
	}

	return false, NoMappedAttrFound

}

func GetAllMappedAttributesContainingREFInMappings(mappings *config.MappedApp,
	toTable string) map[string]bool {

	mappedAttrsWithREF := make(map[string]bool)
	
	for _, mapping := range mappings.Mappings {
		
		// toTable
		for _, tTable := range mapping.ToTables {
			if tTable.Table == toTable {

				// attributes
				for k, v := range tTable.Mapping {
					
					// contain #REF
					if containREF(v) {

						// add only if there does not exist 
						// to make sure mappedAttrsWithREF contains unique attrs
						if _, ok := mappedAttrsWithREF[k]; !ok {
							mappedAttrsWithREF[k] = true
						}
					}
				}
			}
		}
	}

	return mappedAttrsWithREF
}

func containFETCH(data string) bool {

	if strings.Contains(data, "#FETCH(") {
		
		return true

	} else {

		return false
	}
}

func getFirstArgFromFETCH(data string) string {
	
	tmp1 := strings.Split(data, "#REF(")

	tmp2 := strings.Split(tmp1[1], "#FETCH(")

	return strings.Split(tmp2[1], ",")[0]

}

func GetMappedAttributesFromSchemaMappingsByFETCH(allMappings *config.SchemaMappings,
	fromApp, fromAttr, toApp, toTable string) ([]string, error) {
	
	var attributes []string

	log.Println(fromApp, fromAttr, toApp, toTable)

	toAppMappings, err := GetToAppMappings(allMappings, fromApp, toApp)
	if err != nil {
		log.Fatal(err)
	}

	for _, mapping := range toAppMappings.Mappings {

		for _, tTable := range mapping.ToTables {
			
			// toTable
			if tTable.Table == toTable {
				// log.Println(tTable)
				for tAttr, fAttr := range tTable.Mapping {
					// fromAttr

					// If there exists #REF
					if containREF(fAttr) {

						if containFETCH(fAttr) {

							fAttr = getFirstArgFromFETCH(fAttr)

							if fAttr == fromAttr {
								attributes = append(attributes, tAttr)
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

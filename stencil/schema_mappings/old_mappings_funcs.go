package schema_mappings

import (
	"stencil/config"
	"strings"
	"log"
)

// Return the first argument of #REF
func oldGetFirstArgFromREF(ref string) string {

	tmp := strings.Split(ref, "#REF(")
	
	return strings.Split(tmp[1], ",")[0]

}

func OldGetMappedAttributesFromSchemaMappings(allMappings *config.SchemaMappings,
	fromApp, fromTable, fromAttr,
	toApp, toTable string, ignoreREF bool) ([]string, error) {

	var attributes []string

	// In the case: diaspora posts posts.id mastodon statuses
	// there are two ids in the result if we don't use uniqueAttrs
	// because posts are mapped to statuses in two different conditions: 
	// "posts.type": "StatusMessage" and "posts.type": "Reshare"
	uniqueAttrs := make(map[string]bool)

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

							fAttr = RemoveASSIGNAllRightParenthesesIfExists(fAttr)

							// If not ignore #REF
							if !ignoreREF {

								// If there exists #REF
								if containREForREFHARD(fAttr) {
									fAttr = getFirstArgFromREF(fAttr)

									if fAttr == fromAttr {
										uniqueAttrs[tAttr] = true
									}
								}

							} else {

								// If there does not exist #REF
								if !containREForREFHARD(fAttr) {

									if fAttr == fromAttr {
										uniqueAttrs[tAttr] = true
									}
								}
							}									
							
						}
					}
				}
			}
		}
	}

	for attr := range uniqueAttrs {
		attributes = append(attributes, attr)
	}

	if len(attributes) == 0 {

		return nil, NoMappedAttrFound

	} else {

		return attributes, nil
	}
}

func OldGetMappedAttributesFromSchemaMappingsByFETCH(allMappings *config.SchemaMappings,
	fromApp, fromAttr, toApp, toTable string) ([]string, error) {

	var attributes []string

	uniqueAttrs := make(map[string]bool)

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
					// #FETCH and #ASSIGN don't coexist, so we don't check #ASSIGN
					if containREForREFHARD(fAttr) {

						if containFETCH(fAttr) {

							fAttr = getFirstArgFromFETCH(fAttr)

							if fAttr == fromAttr {
								
								uniqueAttrs[tAttr] = true
							}

						}
						
					}
				}
			}
		}
	}

	for atrr := range  uniqueAttrs {
		attributes = append(attributes, atrr)
	}

	if len(attributes) == 0 {

		return nil, NoMappedAttrFound

	} else {

		return attributes, nil
	}
}
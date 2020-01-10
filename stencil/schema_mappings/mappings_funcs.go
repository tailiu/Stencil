package schema_mappings

import (
	"stencil/config"
	"strings"
	"log"
	"os"
	"io/ioutil"
	"encoding/json"
	combinations "github.com/mxschmitt/golang-combinations"
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


func writeMappingsToFile(pairwiseMappings *config.SchemaMappings) {

	bytes, err := json.MarshalIndent(pairwiseMappings, "", "	")

	if err != nil {
		log.Fatal(err)
	}
 
	err1 := ioutil.WriteFile(OUTPUTFILEPATH, bytes, 0644)

	if err1 != nil {
		log.Fatal(err1)
	}

}

// Get all unique applications in the pairwise schema mappings 
func getApplications(pairwiseMappings *config.SchemaMappings) []string {

	var apps []string

	for _, mapping := range pairwiseMappings.AllMappings {
		apps = append(apps, mapping.FromApp)
	}

	return apps

}

// Get all the permutations of an array
func permutations(arr []string) [][]string {

    var helper func([]string, int)
	
	res := [][]string{}

    helper = func(arr []string, n int) {

        if n == 1 {

            tmp := make([]string, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
			
        } else {
			
			for i := 0; i < n; i++ {
				
				helper(arr, n - 1)
		
				if n % 2 == 1 {

                    tmp := arr[i]
					arr[i] = arr[n - 1]
					arr[n - 1] = tmp

                } else {

                    tmp := arr[0]
					arr[0] = arr[n - 1]
					arr[n - 1] = tmp

                }
            }
        }
	}
	
	helper(arr, len(arr))
	
	return res

}

// find all possible mappings through different paths from a source app to a destination app
// This is equivalent to getting all permutations of an array
func getMappingsPaths(apps []string) [][]string {

	// // i is the index of the source app
	// for i := 0; i < len(apps); i++ {
		
	// 	srcApp := apps[i]

	// 	// j is the index of the destination app
	// 	for j := 0; j < len(apps); j++ {
			
	// 		// If i == j, this means the source and destination apps are the same
	// 		if i == j {
	// 			continue
	// 		}

	// 		dstApp := apps[j]

	// 	}
	// }

	var res [][]string 
	
	combs := combinations.All(apps)
	
	log.Println(combs)

	for _, comb := range combs {
		
		if len(comb) <= 2 {
			continue
		}

		res = append(res, permutations(comb)...)

	}

	log.Println(res)

	return res

}

func loadPairwiseSchemaMappings() (*config.SchemaMappings, error) {

	var SchemaMappingsObj config.SchemaMappings

	// dir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(dir)

	jsonFile, err := os.Open(INPUTFILEPATH)
	if err != nil {
		log.Println(err)
		return &SchemaMappingsObj, CannotOpenPSMFile
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	jsonAsBytes, err1 := ioutil.ReadAll(jsonFile)
	if err1 != nil {
		log.Fatal(err1)
	}

	json.Unmarshal(jsonAsBytes, &SchemaMappingsObj)

	// dbConn := db.GetDBConn(db.STENCIL_DB)

	// log.Println(SchemaMappingsObj)

	// for i, mapping := range SchemaMappingsObj.AllMappings {
	// 	for j, toApp := range mapping.ToApps {
	// 		appID := db.GetAppIDByAppName(dbConn, toApp.Name)
	// 		for k, toAppMapping := range toApp.Mappings {
	// 			for l, toTable := range toAppMapping.ToTables {
	// 				ToTableID, err := db.TableID(dbConn, toTable.Table, appID);
	// 				if  err != nil{
	// 					log.Println("LoadSchemaMappings: Unable to resolve ToTableID for table: ", 
	// 						toTable.Table, toApp.Name, appID)
	// 					log.Fatal(err)
	// 				}
	// 				SchemaMappingsObj.AllMappings[i].ToApps[j].Mappings[k].ToTables[l].TableID 
	// 					= ToTableID
	// 				// fmt.Println(toTable.Table, toApp.Name, appID, ToTableID)
	// 			}
	// 		}
	// 	}
	// }
	// fmt.Println(SchemaMappingsObj.AllMappings[0].ToApps[0].Mappings[0].ToTables)
	// log.Fatal()

	return &SchemaMappingsObj, nil

}
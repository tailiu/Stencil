package schema_mappings

import (
	"stencil/config"
	"stencil/common_funcs"
	"strings"
	"log"
	"os"
	"io/ioutil"
	"encoding/json"
	"math/rand"
	"time"
	combinations "github.com/mxschmitt/golang-combinations"
)

func removeREFAndParatheses(data string) string {

	tmp := strings.Replace(data, "#REF(", "", -1)

	tmp1 := strings.Replace(tmp, ")", "", -1)

	return tmp1
}

// Return the first argument of #REF
func getFirstArgFromREF(ref string) string {

	tmp := strings.Split(ref, ",")
	
	return tmp[0]

}

func getThirdArgFromREFIfExists(ref string) string {

	tmp := strings.Split(ref, ",")

	if len(tmp) == 3 {
		return tmp[2]
	} else {
		return ""
	}

}

// For example, we have 
// "#REF(#FETCH(posts.id,posts.guid,photos.status_message_guid),posts.id,statuses)"
// to get "statuses"
func getThirdArgFromREFContainingFETCHIfExists(ref string) string {

	tmp1 := strings.Replace(ref, "#REF(#FETCH(", "", -1)

	tmp2 := strings.Replace(tmp1, ")", "", -1)
	
	tmp3 := strings.Split(tmp2, ",")

	if len(tmp3) != 5 {
		return ""
	} else {
		return tmp3[4]
	}

}

func containREF(data string) bool {

	if strings.Contains(data, "#REF(") {
		
		return true

	} else {

		return false
	}
}

func GetAllApps(allMappings *config.SchemaMappings) []string {

	var allApps []string

	for _, mappings := range allMappings.AllMappings {
		allApps = append(allApps, mappings.FromApp)
	}

	return allApps

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

// For example: #REF(#ASSIGN(messages.id),messages.id)
// will become #REF(messages.id,messages.id
func RemoveASSIGNAllRightParenthesesIfExists(data string) string {

	if strings.Contains(data, "#ASSIGN(") {

		tmp := strings.Replace(data, "#ASSIGN(", "", -1)

		tmp = strings.Replace(tmp, ")", "", -1)

		// log.Println(tmp)

		return tmp

	} else {

		return data
	}

}

func GetMappedAttributesToUpdateOthers(
	allMappings *config.SchemaMappings,
	fromApp, fromTable, fromAttr,
	toApp, toTable string) ([]string, error) {

	var attributes []string

	// In the case: diaspora posts posts.id mastodon statuses
	// since posts are mapped to statuses in two different conditions: 
	// "posts.type": "StatusMessage" and "posts.type": "Reshare"
	// there are two ids in the result if we don't use uniqueAttrs
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

							// If there does not exist #REF
							if !containREF(fAttr) {

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

	for attr := range uniqueAttrs {

		attributes = append(attributes, attr)

	}

	if len(attributes) == 0 {

		return nil, NoMappedAttrFound
	
	} else {

		return attributes, nil
	}

}

func GetMappedAttributesToBeUpdated(
	allMappings *config.SchemaMappings,
	fromApp, fromTable, fromAttr,
	toApp, toTable string) (map[string]string, error) {

	attributes := make(map[string]string)

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

							// If there exists #REF
							if containREF(fAttr) {

								fAttr = removeREFAndParatheses(fAttr)

								firstArg := getFirstArgFromREF(fAttr)

								if firstArg == fromAttr {
									
									thirdArg := getThirdArgFromREFIfExists(fAttr)

									if thirdArg != "" {
										// There could be cases where duplicate tAttrs are found
										// For example, diaspora posts posts.author_id mastodon statuses
										// duplicate tAttr does not influence results
										if _, ok := attributes[tAttr]; ok {
											log.Println(duplicateToAttrWithThirdArg)
										} else {
											attributes[tAttr] = thirdArg
										}
									} else {
										if _, ok := attributes[tAttr]; ok {
											log.Println(duplicateToAttrWithoutThirdArg)
										} else {
											attributes[tAttr] = ""
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

func GetMappedAttributesToBeUpdatedByFETCH(
	allMappings *config.SchemaMappings,
	fromApp, fromAttr, 
	toApp, toTable string) (map[string]string, error) {
	
	attributes := make(map[string]string)

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
					if containREF(fAttr) {

						if containFETCH(fAttr) {

							firstArgInFETCH := getFirstArgFromFETCH(fAttr)

							if firstArgInFETCH == fromAttr {
								
								thirdArgInREF := getThirdArgFromREFContainingFETCHIfExists(fAttr)

								if thirdArgInREF != "" {
									if _, ok := attributes[tAttr]; ok {
										log.Println(duplicateToAttrWithThirdArg)
									} else {
										attributes[tAttr] = thirdArgInREF
									}
								} else {
									if _, ok := attributes[tAttr]; ok {
										log.Println(duplicateToAttrWithoutThirdArg)
									} else {
										attributes[tAttr] = ""
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
// This function will check all these possiblities
func REFExists(mappings *config.MappedApp, 
	toTable, toAttr string) (bool, error) {

	// log.Println(toTable, toAttr)

	foundMappings := false

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

						foundMappings = true
					}
				} 

				// These are commented to take into account the cases mentioned above
				// else {
				// 	return false, NoMappedAttrFound
				// }
			}
		}
	}

	if foundMappings {
		return false, nil
	} else {
		return false, NoMappedAttrFound
	}

}

func GetFirstArgsInREFByToTableToAttr(mappings *config.MappedApp, 
	toTable, toAttr string) []string {

	// log.Println(toTable, toAttr)
	
	var firstArgs []string

	for _, mapping := range mappings.Mappings {
		
		// toTable
		for _, tTable := range mapping.ToTables {
			if tTable.Table == toTable {
				
				// toAttr
				if mappedAttr, ok := tTable.Mapping[toAttr]; ok {
					// log.Println(mappedAttr)

					// If there exists #REF
					if containREF(mappedAttr) {	

						var firstArg string

						mappedAttr = removeREFAndParatheses(mappedAttr)

						// For example,
						// "status_id":"#REF(#FETCH(posts.id,posts.guid,
						//	photos.status_message_guid),posts.id,statuses)",
						if containFETCH(mappedAttr) {
							
							firstArg = getFirstArgFromFETCH(mappedAttr)

						} else {
							
							firstArg = getFirstArgFromREF(mappedAttr)
						}

						firstArg = RemoveASSIGNAllRightParenthesesIfExists(firstArg)
						
						if !common_funcs.ExistsInSlice(firstArgs, firstArg) {
							firstArgs = append(firstArgs, firstArg)
						}
						
					} 
				} 
			}
		}
	}

	return firstArgs
}

func GetAllMappedAttributesContainingREFInMappings(mappings *config.MappedApp,
	toTable string) []string {

	var  mappedAttrsWithREF []string
	
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
						// if _, ok := mappedAttrsWithREF[k]; !ok {
						// 	mappedAttrsWithREF[k] = true
						// }

						if !common_funcs.ExistsInSlice(mappedAttrsWithREF, k) {
							mappedAttrsWithREF = append(mappedAttrsWithREF, k)
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

	var tmp1, tmp2 []string

	if containREF(data) {

		tmp1 = strings.Split(data, "#REF(")

		tmp2 = strings.Split(tmp1[1], "#FETCH(")

	} else {

		tmp2 = strings.Split(tmp1[1], "#FETCH(")

	}

	return strings.Split(tmp2[1], ",")[0]

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

func isAlreadyChecked(mappingsPathToBeChecked []string, checkedMappingsPaths [][]string) bool {

	for _, checkedMappingsPath := range checkedMappingsPaths {

		if len(checkedMappingsPath) == len(mappingsPathToBeChecked) {

			matched := true
			
			for i, app := range checkedMappingsPath {

				if app != mappingsPathToBeChecked[i] {
					matched = false
					break
				}
			}

			if matched {
				return true
			}
		}
	}

	return false

}

func shuffleSlice(s [][]string) {
	
	rand.Seed(time.Now().UnixNano())
	
	rand.Shuffle(len(s), func(i, j int) { 
		s[i], s[j] = s[j], s[i] 
	})

}
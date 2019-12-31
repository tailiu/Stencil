package schema_mappings

import (
	"stencil/config"
	// "stencil/db"
	"os"
	"log"
	"io/ioutil"
	"encoding/json"
	"strings"
	combinations "github.com/mxschmitt/golang-combinations"
)

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

	pairwiseSchemaMappingFile := "./config/app_settings/pairwise_mappings.json"

	jsonFile, err := os.Open(pairwiseSchemaMappingFile)
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
	// 				SchemaMappingsObj.AllMappings[i].ToApps[j].Mappings[k].ToTables[l].TableID = ToTableID
	// 				// fmt.Println(toTable.Table, toApp.Name, appID, ToTableID)
	// 			}
	// 		}
	// 	}
	// }
	// fmt.Println(SchemaMappingsObj.AllMappings[0].ToApps[0].Mappings[0].ToTables)
	// log.Fatal()

	return &SchemaMappingsObj, nil

}

func containVar(data string) bool {

	if strings.Contains(data, "$") {
		
		return true

	} else {

		return false
	}

}

func getVarName(data string) string {

	tmp := strings.Split(data, "$")

	return tmp[1]

}

func replaceVar(fromTableAttr string, inputs []map[string]string) string {

	varName := getVarName(fromTableAttr)

	for _, input := range inputs {

		for k, v := range input {

			// log.Println("********")
			// log.Println(name)
			// log.Println(val)
			// log.Println("********")
			if k == "name" && v == varName {
				return input["value"]
			}
		}
	}

	return ""
}

func getConditions(toTable *config.ToTable) map[string]string {

	return toTable.Conditions

}

func containFunction(data string) bool {

	if strings.Contains(data, "#") {
		
		return true

	} else {

		return false
	}

}

func procMappingsByRows(toApp *config.MappedApp, isSourceApp bool) map[string]string {

	res := make(map[string]string)

	for _, mapping := range toApp.Mappings {

		for _, toTable := range mapping.ToTables {

			toTableName := toTable.Table

			// When mappings are not accurate, app developers can specify that
			// this should not be used in PSM by setting NotUsedInPSM as true
			// For example, 
			if toTable.NotUsedInPSM {
				continue
			}

			// log.Println(getConditions(&toTable))

			// Conditions are very hard to cope with correctly through PSM
			// In this version, PSM does not process conditions
			if getConditions(&toTable) != nil {
				continue
			}

			if toTable.Conditions == nil {

				for toAttr, fromTableAttr := range toTable.Mapping {
					
					// PSM does not process mappings containing #REF
					// because the function #REF is very complex and must be defined by app developers
					// and should not be got by PSM
					if containREF(fromTableAttr) {
						continue
					}

					// For other functions defined by #, if they are included in the source app,
					// they could be included in the PSM result.
					// If they are included in the intermediate apps, then they cannot be included in
					// the PSM result
					// For example, users.#RANDINT -> users.id, users.id -> accounts.id 
					// 				=> users.#RANDINT -> accounts.id, but
					// 				accounts.id -> users.id, #RANDINT -> users.id, 
					// 				/=> users.#RANDINT -> accounts.id
					// #RANDINT can never be the same as table.attr
					if !isSourceApp && containFunction(fromTableAttr) {
						continue
					}

					// log.Println("toAttr:", toAttr)
					// log.Println("fromTableAttr:", fromTableAttr)

					// Similar to functions, for variables in the mappings, like $follow_action,
					// only when they are included in the source app, will they be included in the results.
					// Further, these variables need to be replaced with real values first since the dst app
					// may not define such kind of inputs
					if containVar(fromTableAttr) {
						if !isSourceApp {
							continue
						} else {
							fromTableAttr = replaceVar(fromTableAttr, toApp.Inputs)
						}
					}

					// Note that toTableName.toAttr could be not unique. For example,
					// Twitter.tweets and Twitter.retweets are both mapped to Mastodon.statuses.
					res[toTableName  + "." + toAttr] = fromTableAttr
				}

			}
		}
	}

	log.Println(res)

	return res

}

// Add mappings by PSM through the mapping path
// For example: 
// through a mapping path: Mastodon -> Twitter -> Gnusocial -> Diaspora,
// we can get mappings from Mastodon to Diaspora.
// This is an old design without considering how to handle conditions
func OldAddMappingsByPSMThroughOnePath(pairwiseMappings *config.SchemaMappings, 
	mappingsPath []string) {

	for i := 0; i < len(mappingsPath) - 1; i++ {

		currApp := mappingsPath[i]

		nextApp := mappingsPath[i + 1]

		for _, mappings := range pairwiseMappings.AllMappings {

			// find the current app
			if mappings.FromApp == currApp {	
				
				for _, toApp := range mappings.ToApps {

					// find the next app
					if toApp.Name == nextApp {

						isSourceApp := true

						if i == 0 {

							// procRes := procMappingsByRows(&toApp, isSourceApp)
							procMappingsByRows(&toApp, isSourceApp)

						} else {

							// procRes := procMappingsByRows(&toApp, isSourceApp)
							procMappingsByRows(&toApp, !isSourceApp)

						}
					}
				}
			}
		}
	}

}

func procMappingsByTables(firstMappings, secondMappings *config.MappedApp) {

	for _, firstMapping := range firstMappings.Mappings {

		for _, firstToTable := range firstMapping.ToTables {

			for _, secondMapping := range secondMappings.Mappings {

				for _, secondFromTable := range secondMapping.FromTables {

					// find matched tables
					if secondFromTable == firstToTable {

						for _, secondToTable := range secondMapping.ToTables {

							conditions := secondToTable.Conditions
							
							// check conditions
							if  

						}

					}

				}

			}

		}

	}

}

func findFromAppToAppMappings(pairwiseMappings *config.SchemaMappings, 
	fromApp, toApp string) (*config.MappedApp, error) {

	for _, mappings := range pairwiseMappings.AllMappings {

		// find the from app
		if mappings.FromApp == fromApp {	
			
			for _, toApp := range mappings.ToApps {

				// find the to app
				if toApp.Name == toApp {

					return &toApp, nil
				}
			}
		}
	}

	return nil, CannotFindPairwiseMappings

}

// This new design adds a condition handler
func addMappingsByPSMThroughOnePath(pairwiseMappings *config.SchemaMappings, 
	mappingsPath []string) {

	for i := 0; i < len(mappingsPath) - 2; i++ {

		currApp := mappingsPath[i]

		nextApp := mappingsPath[i + 1]

		nextNextApp := mappingsPath[i + 2]

		firstMappings, err1 := findFromAppToAppMappings(pairwiseMappings, currApp, nextApp)
		if err1 != nil {
			log.Fatal(err1)
		}

		secondMappings, err2 := findFromAppToAppMappings(pairwiseMappings, nextApp, nextNextApp)
		if err2 != nil {
			log.Fatal(err2)
		}
		
		procMappingsByTables(firstMappings, secondMappings)

	}
	
}


func DeriveMappingsByPSM() (*config.SchemaMappings, error) {

	pairwiseMappings, err := loadPairwiseSchemaMappings()
	if err != nil {
		log.Fatal(err)
	}

	apps := getApplications(pairwiseMappings)

	log.Println(apps)

	mappingsPaths := getMappingsPaths(apps)

	for _, mappingsPath := range mappingsPaths {

		addMappingsByPSMThroughOnePath(pairwiseMappings, mappingsPath)

	}

	// return pairwiseMappings, nil

	return nil, nil

}
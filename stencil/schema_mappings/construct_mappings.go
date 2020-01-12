package schema_mappings

import (
	"stencil/config"
	"log"
	"strings"
)

func createFromAppWhenMissingFromApp(pairwiseMappings *config.SchemaMappings, 
	fromAppName, toAppName string) *config.MappedApp {

	toApp := config.MappedApp {
		Name: toAppName,
	}

	mappings := config.SchemaMapping {
		FromApp: fromAppName,
		ToApps: []config.MappedApp {
			toApp,
		},
	}

	pairwiseMappings.AllMappings = append(pairwiseMappings.AllMappings, mappings)

	mappingsLen := len(pairwiseMappings.AllMappings)

	return &pairwiseMappings.AllMappings[mappingsLen - 1].ToApps[0]
		

}

func createToAppWhenMissingToApp(pairwiseMappings *config.SchemaMappings, 
	fromAppName, toAppName string) (*config.MappedApp, error) {

	for i, mappings := range pairwiseMappings.AllMappings {

		// find the from app
		if mappings.FromApp == fromAppName {	
			
			toApps := mappings.ToApps
			
			toApp := config.MappedApp {
				Name: toAppName,
			}

			toApps = append(toApps, toApp)

			pairwiseMappings.AllMappings[i].ToApps = toApps

			return &pairwiseMappings.AllMappings[i].ToApps[len(toApps) - 1], nil
		}
	}

	return nil, CannotCreateToApp

}

func getFromTablesVariablesFromToTable(toTable config.ToTable) (map[string]bool, map[string]bool) {

	fromTables := make(map[string]bool) 
	variables := make(map[string]bool)

	for _, v := range toTable.Mapping {

		tmp := strings.Split(v, ".")

		if len(tmp) == 2 {
			if _, ok := fromTables[tmp[0]]; !ok {
				fromTables[tmp[0]] = true
			}
		} 

		if len(tmp) == 1 {
			if _, ok := variables[tmp[0]]; !ok {
				variables[tmp[0]] = true
			}
		}

		if len(tmp) > 2 {
			panic("Should never happen here!")
		}

	}

	return fromTables, variables
	
}

func mergeSecondMapToFirstMap(m1, m2 map[string]bool) {

	for k2, v2 := range m2 {
		m1[k2] = v2
	}

}

func getMappingsByFromTable(mappedApp *config.MappedApp, 
	fromTableName string) (*config.Mapping, error) {

	for _, mappings := range mappedApp.Mappings {

		for _, fromTable := range mappings.FromTables {

			if fromTable == fromTableName {
				return &mappings, nil
			}
		}
	}

	return nil, CannotGetMappingsByFromTable

}

func getToTableByName(mappings *config.Mapping, 
	tableName string) (*config.ToTable, error) {

	for i, toTable := range mappings.ToTables {

		if toTable.Table == tableName {
			return &mappings.ToTables[i], nil
		}
	}

	return nil, CannotGetToTableByName

}

func isInFromTables(data string, fromTables []string) bool {

	for _, table := range fromTables {

		if table == data {
			return true
		}
	}

	return false
}

// Compared with createToTableWhenMissingToTable and oldCreateToTableWhenMissingToTable1,
// This function does not work because it just changes the address pointed by existingMappings
// but the mapping in the mappedApp still points to the old address
func oldCreateToTableWhenMissingToTable2(existingMappings *config.Mapping, 
	toTable *config.ToTable, fromTables []string) {

	newToTable := config.ToTable {
		Table: toTable.Table,
		Mapping: make(map[string]string),
	}
	
	// We need to consider whether value contains variables or functions
	// Note that variables do not contain "$" because we have already 
	// it with real values, but they do not contain "."
	// For functions and variables, we simply add them even though they do not 
	// have from tables
	for k, v := range toTable.Mapping {
		
		if containFunction(v) {
			newToTable.Mapping[k] = v
			continue
		}

		tmp := strings.Split(v, ".") 
		
		// This indicates that it is an variable
		// The variable should be written as "$var"
		if len(tmp) == 1 {
			
			newToTable.Mapping[k] = "$" + v

		} else if len(tmp) == 2 {

			if isInFromTables(tmp[0], fromTables) {
				newToTable.Mapping[k] = v
			}

		}
	}
	
	// log.Println(newToTable)

	// Note that existingMappings
	existingMappings.ToTables = append(existingMappings.ToTables, newToTable)

	// return &existingMappings

	// log.Println(existingMappings)

}

func addVariablesIfNotExist(mappedApp *config.MappedApp, totalVariables map[string]bool) {

	for variable, _ := range totalVariables {
		
		alreadyExisted := false
		
		// For simplicity, I use the variable as both the name and the value
		// Also, even though there could be other variables with the same value,
		// I just create new variables.
		for _, input := range mappedApp.Inputs {
			if input["name"] == variable && input["value"] == variable {
				alreadyExisted = true
				break
			}
		}

		if !alreadyExisted {

			newVar := map[string]string {
				"name": variable,
				"value": variable,
			}
			mappedApp.Inputs = append(mappedApp.Inputs, newVar)

		}

	}
}

func createFromTablesWhenMissingFromTables(mappedApp *config.MappedApp, 
	missingFromTables []string, toTable *config.ToTable, existingFromTableGroups [][]string) {

	newMappings := config.Mapping {
		FromTables: missingFromTables,
		ToTables:   []config.ToTable {
			config.ToTable {
				Table: toTable.Table,
				Mapping: make(map[string]string),
			}}}

	for k, v := range toTable.Mapping {

		if containFunction(v) {
			newMappings.ToTables[0].Mapping[k] = v
			continue
		}

		tmp := strings.Split(v, ".") 
		
		// This indicates that it is an variable
		// The variable should be written as "$var"
		if len(tmp) == 1 {
			
			newMappings.ToTables[0].Mapping[k] = "$" + v

		} else if len(tmp) == 2 {

			if isInFromTables(tmp[0], missingFromTables) {
				newMappings.ToTables[0].Mapping[k] = v
			}
		}
	}

	mappedApp.Mappings = append(mappedApp.Mappings, newMappings)

}

func addMappingsIfNotExist(existingToTable, toTable *config.ToTable, fromTables []string) {

	// Similarly, we need to cope with values containing variables or functions in the same ways
	// If the existingToTable already has mappings from a key, we use the existing ones
	// The only uncerntain thing is whether we need to check existing conditions 
	// before adding new mappings
	// For now, the decision is not to consider conditions since we have considered
	// conditions during PSM transitive transformation and if there are many multiple 
	// mappings, e.g., in the path: gnusocial twitter mastodon
	// notice -> tweets/retweets -> statuses
	// notice -> statuses ("notice.reply_to": "#NULL") / statuses ("notice.reply_to": "#NOTNULL")
	// we just add to each.
	for k, v := range toTable.Mapping {

		if _, ok := existingToTable.Mapping[k]; ok {
			continue
		}

		if containFunction(v) {
			existingToTable.Mapping[k] = v
			continue
		}

		tmp := strings.Split(v, ".") 
		
		// This indicates that it is an variable
		// The variable should be written as "$var"
		if len(tmp) == 1 {
			
			existingToTable.Mapping[k] = "$" + v

		} else if len(tmp) == 2 {

			if isInFromTables(tmp[0], fromTables) {
				existingToTable.Mapping[k] = v
			}
		}
	}

}

func oldCreateToTableWhenMissingToTable1(pairwiseMappings *config.SchemaMappings, 
	fromAppName, toAppName string, toTable *config.ToTable, fromTables []string) {

	newToTable := config.ToTable {
		Table: toTable.Table,
		Mapping: make(map[string]string),
	}

	for seq1, mappings := range pairwiseMappings.AllMappings {

		// find the from app
		if mappings.FromApp == fromAppName {
			
			for seq2, toApp := range mappings.ToApps {

				// find the to app
				if toApp.Name == toAppName {

					for seq3, mappings := range toApp.Mappings {

						for _, fromTable := range mappings.FromTables {
				
							if isInFromTables(fromTable, fromTables) {
				
								// We need to consider whether value contains variables or functions
								// Note that variables do not contain "$" because we have already 
								// it with real values, but they do not contain "."
								// For functions and variables, we simply add them 
								// even though they do not have from tables
								for k, v := range toTable.Mapping {
									
									if containFunction(v) {
										newToTable.Mapping[k] = v
										continue
									}
				
									tmp := strings.Split(v, ".") 
									
									// This indicates that it is an variable
									// The variable should be written as "$var"
									if len(tmp) == 1 {
										
										newToTable.Mapping[k] = "$" + v
				
									} else if len(tmp) == 2 {
				
										if isInFromTables(tmp[0], fromTables) {
											newToTable.Mapping[k] = v
										}
				
									}
								}
				
								pairwiseMappings.AllMappings[seq1].
									ToApps[seq2].Mappings[seq3].ToTables = 
										append(
											pairwiseMappings.AllMappings[seq1].
												ToApps[seq2].Mappings[seq3].ToTables,
											newToTable)
								
								return 
							}
							
						}
					}
					
				}
			}
		}
	}
}

func createToTableWhenMissingToTable(mappedApp *config.MappedApp, 
	toTable *config.ToTable, fromTables []string) {

	newToTable := config.ToTable {
		Table: toTable.Table,
		Mapping: make(map[string]string),
	}
	
	for seq1, mappings := range mappedApp.Mappings {

		for _, fromTable := range mappings.FromTables {

			if isInFromTables(fromTable, fromTables) {

				// We need to consider whether value contains variables or functions
				// Note that variables do not contain "$" because we have already 
				// it with real values, but they do not contain "."
				// For functions and variables, we simply add them 
				// even though they do not have from tables
				for k, v := range toTable.Mapping {
					
					if containFunction(v) {
						newToTable.Mapping[k] = v
						continue
					}

					tmp := strings.Split(v, ".") 
					
					// This indicates that it is an variable
					// The variable should be written as "$var"
					if len(tmp) == 1 {
						
						newToTable.Mapping[k] = "$" + v

					} else if len(tmp) == 2 {

						if isInFromTables(tmp[0], fromTables) {
							newToTable.Mapping[k] = v
						}

					}
				}

				mappedApp.Mappings[seq1].ToTables = append(
							mappedApp.Mappings[seq1].ToTables,
							newToTable)

				return 
			}
		}
	}
}

func constructMappingsUsingProcMappings(pairwiseMappings *config.SchemaMappings, 
	procMappings []config.ToTable, srcApp, dstApp string) {

	mappedApp, err := findFromAppToAppMappings(pairwiseMappings, srcApp, dstApp)

	if err != nil {

		if err == CannotFindFromApp {

			mappedApp = createFromAppWhenMissingFromApp(pairwiseMappings, srcApp, dstApp)

		} else {

			mappedApp, err = createToAppWhenMissingToApp(pairwiseMappings, srcApp, dstApp)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// log.Println(mappedApp)

	totalVariables := make(map[string]bool)

	// For each toTable in processed mappings
	for _, toTable := range procMappings {

		// log.Println(toTable)

		// Get from tables and variables in this toTable
		fromTables, variables := getFromTablesVariablesFromToTable(toTable)

		// log.Println(fromTables, variables)

		mergeSecondMapToFirstMap(totalVariables, variables)

		var missingFromTables []string

		var existingFromTableGroups [][]string

		checkedFromTable := make(map[string]*config.Mapping)

		// For each from table in the above toTable
		for fromTable, _ := range fromTables {

			if _, ok := checkedFromTable[fromTable]; ok {
				continue
			}

			// Get existing mappings by the from table
			mappings, err1 := getMappingsByFromTable(mappedApp, fromTable)

			// If there is no from table in the existing mappings
			if err1 != nil {
				log.Println(err1)
				missingFromTables = append(missingFromTables, fromTable)
				checkedFromTable[fromTable] = nil
				continue
			}

			checkedFromTable[fromTable] = mappings

			fromTableGroup := []string {
				fromTable,
			}

			// If there are any from tables already in the exsiting mappings,
			// these from tables need to be grouped based on the structure of
			// the existing mappings
			for _, existingFromTable := range mappings.FromTables {
				
				if existingFromTable != fromTable {

					if _, ok := fromTables[existingFromTable]; ok {

						fromTableGroup = append(fromTableGroup, existingFromTable)

						checkedFromTable[existingFromTable] = mappings
					}
				}
			}

			existingFromTableGroups = append(existingFromTableGroups, fromTableGroup)			
			
		}

		log.Println("existing from table groups")
		log.Println(existingFromTableGroups)

		// The basic rule here is to add mappings based on existing from groups to respect existing structures,
		// For the from tables not in those groups, group those tables together and add them
		// This design decision works in most cases, but may cause app developers 
		// to restructure in the case of one from table in a from table group maps to an entire toTable
		// for example, user, profile -> users, credentials, among the mappings, 
		// profile -> users and user, profile -> credentials
		// so in the result, profile and user are separated because profile already becomes a from group
		// and user, as a missing group, becomes a separate from group
		// profile -> users, credentials
		// user -> users, credentials
		for _, fromTableGroup := range existingFromTableGroups {

			existingMappings := checkedFromTable[fromTableGroup[0]]

			// log.Println(existingMappings)

			toTable1, err2 := getToTableByName(existingMappings, toTable.Table)

			if err2 != nil {

				// log.Print(toTable.Table + ":")
				log.Println(err2)
				createToTableWhenMissingToTable(mappedApp, &toTable, fromTableGroup)
				// log.Println(existingMappings)
				// log.Println(mappedApp)

			} else {

				addMappingsIfNotExist(toTable1, &toTable, fromTableGroup)
				// log.Println(existingMappings)
				// log.Println(mappedApp)
			}

		}

		log.Println("missing from tables:")
		log.Println(missingFromTables)

		if len(missingFromTables) != 0 {
			// all non-existing from tables will be combined and added
			// If 
			createFromTablesWhenMissingFromTables(mappedApp, missingFromTables, 
				&toTable, existingFromTableGroups)
		}

	}

	// all missing variables need to be defined
	addVariablesIfNotExist(mappedApp, totalVariables)

}

func constructMappingsByToTables(pairwiseMappings *config.SchemaMappings, 
	ToTables []config.ToTable, srcApp, dstApp string) *config.MappedApp {

	constructMappingsUsingProcMappings(pairwiseMappings, 
		ToTables, srcApp, dstApp)
	
	mappedApp, err := findFromAppToAppMappings(pairwiseMappings, srcApp, dstApp)

	if err != nil {
		log.Fatal(err)
	}

	return mappedApp

}
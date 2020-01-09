package schema_mappings

import (
	"stencil/config"
	"log"
	"strings"
)

// A list of conditions not considered while processing mappings
// I decided to add this list because sometimes conditions should not be
// considered, for example, in the mapping path: Twitter Mastodon Diaspora
// statuses.reply should not be considered in the mappings from statuses to
// posts/comments while statuses.visibility should be
var conditionsNotConsideredList = []conditionsNotConsidered {

	conditionsNotConsidered {
		fromApp:		"mastodon",
		toApp:			"diaspora",
		fromTables:		[]string {
			"statuses", 
			"status_stats",
		},
		toTable:		"posts",
		condName:		"statuses.reply",
		condVal:		"false",
	},

	conditionsNotConsidered {
		fromApp:		"mastodon",
		toApp:			"diaspora",
		fromTables:		[]string {
			"statuses", 
			"status_stats",
		},
		toTable:		"comments",
		condName:		"statuses.reply",
		condVal:		"true",
	},

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
					// only when they are included in the source app, 
					// will they be included in the results.
					// Further, these variables need to be replaced with real values first 
					// since the dst app may not define such kind of inputs
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

func findFromAppToAppMappings(pairwiseMappings *config.SchemaMappings, 
	fromAppName, toAppName string) (*config.MappedApp, error) {
	
	fromAppExists := false

	for _, mappings := range pairwiseMappings.AllMappings {

		// find the from app
		if mappings.FromApp == fromAppName {	
			
			fromAppExists = true
			
			for _, toApp := range mappings.ToApps {

				// find the to app
				if toApp.Name == toAppName {

					return &toApp, nil
				}
			}
		}
	}

	if !fromAppExists {
		return nil, CannotFindFromApp
	} else {
		return nil, CannotFindToApp
	}
	
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

func replaceVar(variable string, inputs []map[string]string) string {

	varName := getVarName(variable)

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

func areTwoSlicesIdenticalWithoutOrder(s1, s2 []string) bool {

	xMap := make(map[string]int)
    yMap := make(map[string]int)

    for _, xElem := range s1 {
        xMap[xElem]++
    }
    for _, yElem := range s2 {
        yMap[yElem]++
    }

    for xMapKey, xMapVal := range xMap {
        if yMap[xMapKey] != xMapVal {
            return false
        }
	}
	
    return true
}

func isInNotConsideredList(fromApp, toApp, toTable, condName, condVal string,
	fromTables []string) bool {

	// log.Println(conditionName)
	// log.Println(conditionValue)

	for _, conditionInList := range conditionsNotConsideredList {

		if conditionInList.fromApp == fromApp &&
			conditionInList.toApp == toApp &&
			conditionInList.toTable == toTable &&
			conditionInList.condName == condName &&
			conditionInList.condVal == condVal && 
			areTwoSlicesIdenticalWithoutOrder(conditionInList.fromTables, fromTables) {
				return true
		}
	}
	
	return false
}

func satisfyConditions(conditions map[string]string, 
	toTable *config.ToTable, inputs []map[string]string, 
	fromTables []string, toTableName, fromApp, toApp string) bool {

	tableName := toTable.Table

	log.Println(inputs)

	for k, v := range conditions {

		// log.Println(fromApp, toApp, toTableName, k, v, fromTables)
		// log.Println(k, v, isInNotConsideredList(fromApp, toApp, toTableName, k, v, fromTables))

		if isInNotConsideredList(fromApp, toApp, toTableName, k, v, fromTables) {
			continue
		}

		satisfyThisCondition := false

		// If conditions contain functions like #NOTNULL or #NULL,
		// such conditions are used when migrating data and not used in PSM
		if containFunction(v) {
			continue
		}

		for k1, v1 := range toTable.Mapping {

			// #REF is not involved in conditions
			if containREF(v1) {
				continue
			}

			if tableName + "." + k1 == k {
				
				// log.Println(tableName + "." + k1)
				// log.Println(v1)

				// v1 may contain variables like "$reshare"
				if containVar(v1) {
					v1 = replaceVar(v1, inputs)
					// log.Println(v1)
					// log.Println(v)
				}

				if v1 == v {
					satisfyThisCondition = true
					break
				}
				// log.Println(v1, v)
			}
		}

		if !satisfyThisCondition {
			// log.Println("not satisfied:")
			// log.Println(k, v)
			return false
		}

	}

	return true

}

func mergeTwoMappings(firstToTable, secondToTable *config.ToTable,
	firstInputs []map[string]string) config.ToTable {

	mergedToTable := config.ToTable {
		Table: secondToTable.Table,
		Mapping: make(map[string]string),
	}

	firstTableName := firstToTable.Table

	for k1, v1 := range firstToTable.Mapping {

		for k2, v2 := range secondToTable.Mapping {
			
			// log.Println(k2, v2)

			// PSM does not process #REF this is because even though through PSM 
			// mappings in #REF can be got, they generally have to be further processed and
			// formatted in the #REF in the destination app.
			// Functions cannot be matched
			// For functions defined by #, if they are included in the source app,
			// they could be included in the PSM result.
			// If they are included in the intermediate apps, then they cannot be included in
			// the PSM result
			// For example, users.#RANDINT -> users.id, users.id -> accounts.id 
			// 				=> users.#RANDINT -> accounts.id, but
			// 				accounts.id -> users.id, #RANDINT -> users.id, 
			// 				/=> users.#RANDINT -> accounts.id
			// #RANDINT can never be the same as table.attr
			// Similary, variables in the intermediate apps cannot be matched and included in the
			// PSM result
			if containREF(v2) || containFunction(v2) || containVar(v2) {
				continue
			} 

			// PSM does not process #REF 
			if containREF(v1) {
				continue
			}

			// Find a match
			if firstTableName + "." + k1 == v2 {

				// The variable in v1 needs to be replaced with the real value
				// because the variable is only defined in the first app
				if containVar(v1) {
					v1 = replaceVar(v1, firstInputs)
				}

				mergedToTable.Mapping[k2] = v1
 			}

		}
	}

	return mergedToTable

}


func mergeTwoSameToTables(table1, table2 *config.ToTable) config.ToTable {

	mergedToTable := config.ToTable {
		Table: table1.Table,
		Mapping: make(map[string]string),
	}

	m1 := table1.Mapping
	m2 := table2.Mapping

	for k1, v1 := range m1 {

		if v2, ok := m2[k1]; ok {

			// If we find duplicate (k, v), we simply merge them
			// If we find the same key with different values, we cannot
			// be sure which value to include, so we exclude such key
			if v1 == v2 {
				mergedToTable.Mapping[k1] = v1
			}
		
		// If we do not find the key, we need to include this unique key in m1
		} else {
			mergedToTable.Mapping[k1] = v1
		}

	}

	for k2, v2 := range m2 {

		if _, ok := m1[k2]; !ok {
			
			// Since we alreay delt with the commone keys of the two mappings,
			// we only need to add the unique keys in m2 to the result
			mergedToTable.Mapping[k2] = v2
		}

	}

	return mergedToTable

}

// The most complex part in processing mappings is to handle conditions
// We process mappings on the table level because conditions are defined on the table level,
// in other words, either one table can be mapped or not depending on conditions.
// There could be several special cases: 
// 1. Same source table -> different intermediate tables -> same destination table
// 	e.g., in the path: gnusocial twitter mastodon
// 	notice -> tweets/retweets -> statuses
// 2. Different source tables -> same intermediate table -> same destination table
// 	e.g., in the path: twitter gnusocial mastodon
//	tweets/retweets -> notice -> statuses
// 3. Same source table -> same intermediate tables -> same destination table
// 	e.g., in the path: gnusocial mastodon twitter
//	notice -> statuses ("notice.reply_to": "#NULL") / statuses ("notice.reply_to": "#NOTNULL") 
//	-> tweets
// 4. Same source table -> same intermediate tables -> different destination table
//	e.g., in the path: mastodon gnusocial twitter
//  statuses -> notice -> tweets/retweets
// The general rule to cope with those cases is to 
// keep the path with unqiue (fromTable, toTable) pair,
// so different paths in 1, 3 will be merged and in 2, 4 will be kept
func procMappingsByTables(firstMappings, secondMappings *config.MappedApp) []config.ToTable {

	var mergedMappings []config.ToTable

	firstInputs := firstMappings.Inputs

	// Since mergedMappings stores all merged tables, 
	// we need to use a global sequence
	seq := 0

	for _, firstMapping := range firstMappings.Mappings {

		// We initialize mergedTableNameIndex here
		// because we only want to merge the mappings from same tables to the same table
		// For example, in the path: twitter gnusocial mastodon, 
		// if we initialize these outside the for loop,
		// we may also merge tweets -> notice -> statuses and retweets -> notice -> statuses,
		// which should not be merged
		mergedTableNameIndex := make(map[string]int)

		for _, firstToTable := range firstMapping.ToTables {

			// When mappings are not accurate, app developers can specify that
			// those mappings should not be used in PSM by setting NotUsedInPSM as true
			// For example, the mappings from twitter.conversations to mastodon.conversations
			// (twitter.conversations are the conversations for messages 
			// while mastodon.conversations are the conversations for statuses including messages)
			// and the mappings from mastodon.conversations to gnusocial.conversation are inaccurate. 
			// (gnusocial.conversation are the conversations only for posts not messages)
			// Then if these mappings are used in PSM, 
			// we will get twitter.conversations -> gnusocial.conversation,
			// which is incorrect.
			if firstToTable.NotUsedInPSM {
				continue
			}

			for _, secondMapping := range secondMappings.Mappings {

				secondMappingFromTables := secondMapping.FromTables

				for _, secondFromTable := range secondMappingFromTables {

					// find matched tables
					if secondFromTable == firstToTable.Table {
						
						// log.Println(secondFromTable)

						for _, secondToTable := range secondMapping.ToTables {

							if secondToTable.NotUsedInPSM {
								continue
							}
							
							conditions := getConditions(&secondToTable)
							
							// log.Println(secondToTable.Table)
							// log.Println(satisfyConditions(conditions, &firstToTable, firstInputs))

							// check conditions
							if satisfyConditions(conditions, &firstToTable, firstInputs,
								secondMappingFromTables, secondToTable.Table,
								firstMappings.Name, secondMappings.Name) {

								mergedTable := mergeTwoMappings(&firstToTable, 
									&secondToTable, firstInputs)

								if index, ok := mergedTableNameIndex[mergedTable.Table]; ok {

									// For example, in the path: gnusocial mastodon twitter,
									// If there is no merging, there will be two almost the 
									// same toTables of tweets and retweets 
									// because notice map to statuses in two different conditions. 
									// In this case, we need to merge the two toTable results. 
									// {tweets map[] false map[id:notice.id content:notice.content 
									// 	updated_at:notice.modified created_at:notice.created] map[] } 
									// {retweets map[] false map[created_at:notice.created 
									//  updated_at:notice.modified id:notice.id] map[] } 
									// {tweets map[] false map[content:notice.content 
									// created_at:notice.created updated_at:notice.modified 
									// id:notice.id] map[] } 
									// {retweets map[] false map[id:notice.id created_at:notice.created 
									//  updated_at:notice.modified] map[] }
									mergedTable = mergeTwoSameToTables(&mergedMappings[index], 
										&mergedTable)	
									
									// log.Println("Merge two tables results:", mergedTable)
									
									mergedMappings[index] = mergedTable

								} else {

									// Only add to merged mappings when 
									// there are combined mappings returned
									if len(mergedTable.Mapping) != 0 {

										mergedMappings = append(mergedMappings, mergedTable) 

										mergedTableNameIndex[mergedTable.Table] = seq
										seq += 1

									}	
								}	
							}
						}
					}
				}
			}
		}
	}

	return mergedMappings

}
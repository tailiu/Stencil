package app_dependency_handler

import (
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/app_display"
	// "strconv"
	"strings"
)

// func checkResolveReferenceInGetDataInParentNode(displayConfig *config.DisplayConfig,
// 	id, table0, col0, table1, col1, value string) (map[string]interface{}, error) {

// 	log.Println("+++++++++++++++++++")
// 	log.Println(table0)
// 	log.Println(col0)
// 	log.Println(table1)
// 	log.Println(col1)
// 	log.Println("+++++++++++++++++++")
	
// 	table0ID := displayConfig.AppConfig.TableNameIDPairs[table0]
// 	table1ID := displayConfig.AppConfig.TableNameIDPairs[table1]

// 	// First, we need to get the attribute that requires reference resolution
// 	// For example, we have *account.id*, and we want to get *users.account_id*
// 	// We check whether account.id needs to be resolved
// 	if reference_resolution.NeedToResolveReference(displayConfig, table0, col0) {

// 		log.Println("Before checking reference1 resolved or not")

// 		// If account.id should be resolved (in this case, it should not),
// 		// we check whether the reference has been resolved or not
// 		newVal := reference_resolution.ReferenceResolved(displayConfig, table0ID, col0, id)
		
// 		// If the reference has been resolved, then use the new reference to get data
// 		if newVal != "" {

// 			log.Println("reference has been resolve1")

// 			return getOneRowBasedOnDependency(displayConfig, table1, col1, newVal)
		
// 		// Otherwise, we try to resolve the reference
// 		} else {

// 			hint0 := app_display.CreateHint(table0, table0ID, id)
// 			log.Println("Before resolving reference1: ", hint0)

// 			updatedAttrs, _ := reference_resolution.ResolveReference(displayConfig, hint0)

// 			// We check whether the desired attr (col0) has been resolved
// 			foundResolvedAttr := false
// 			for attr, val := range updatedAttrs {
// 				if attr == col0 {
// 					newVal = val
// 					foundResolvedAttr = true
// 					break
// 				}
// 			}

// 			// If we find that col0 has been resolved, then we can use it to get other data
// 			if foundResolvedAttr {

// 				return getOneRowBasedOnDependency(displayConfig, table1, col1, newVal)
			
// 			// Otherwise we cannot use the unresolved reference to get other data in node
// 			} else {

// 				return nil, app_display.CannotResolveReferencesGetDataInNode
// 			}
// 		}

// 	// We check if users.account_id needs be resolved (of course, in this case, it should be)
// 	// However we don't know its id. 
// 	} else if reference_resolution.NeedToResolveReference(displayConfig, table1, col1) {

// 		log.Println("Before checking reference2 resolved or not")

// 		// We assume that users.account_id has already been resolved and get its data
// 		data, err := getOneRowBasedOnDependency(displayConfig, table1, col1, value)
// 		if err != nil {
// 			return nil, app_display.CannotFindRemainingData
// 		}

// 		// Now we have the id of the data, we should check whether it has been resolved before, 
// 		// but actually if we can get one, it should always be the one we want to get because
// 		// otherwise there will be multiple rows corresponding to one member.
// 		// There could be the case where ids are not changed, 
// 		// so even if references are not resolved, 
// 		// we can still get the rows we want, but we need to resolve it further.
// 		newVal := reference_resolution.ReferenceResolved(displayConfig, 
// 			table1ID, col1, fmt.Sprint(data["id"]))

// 		if newVal != "" {

// 			// Theoretically, if it has been resolved, then it should be the value we have 
// 			// given that one member corresponds to one row
// 			if newVal == value {

// 				log.Println("reference has been resolve2")
// 				log.Println(data)

// 				return data, nil

// 			} else {

// 				panic("Should not happen given one member corresponds to one row for now!")
			
// 			}

// 		} else {
			
// 			hint1 := app_display.CreateHint(table1, table1ID, fmt.Sprint(data["id"]))
// 			log.Println("Before resolving reference2: ", hint1)

// 			updatedAttrs, _ := reference_resolution.ResolveReference(displayConfig, hint1)

// 			// We check whether the desired attr (col1) has been resolved 
// 			// (until this point, it should be resolved)
// 			foundResolvedAttr := false
// 			for attr, val := range updatedAttrs {
// 				if attr == col1 {
// 					newVal = val
// 					foundResolvedAttr = true
// 					break
// 				}
// 			}

// 			// If we find that col0 has been resolved, then we can use it to get other data
// 			if foundResolvedAttr {

// 				if newVal == value {

// 					return data, nil
	
// 				} else {
	
// 					panic("Should not happen given one member corresponds to one row for now!")
// 				}
			
// 			// This should not happen
// 			} else {

// 				panic("Should not happen given one member corresponds to one row for now!")
// 			}
// 		}
	
// 	// Normally, there must exist one that needs to be resolved. 
// 	// Howver, the following can happen when there is no mapping
// 	// For example: 
// 	// When migrating from Diaspora to Mastodon:
// 	// there is no mapping to stream_entries.activity_id.
// 	} else {

// 		return nil, app_display.NoMappingAndNoReferenceToResolve
// 	}
// }

func getHintsInParentNode(displayConfig *config.DisplayConfig, 
	hints []*app_display.HintStruct, conditions []string) (*app_display.HintStruct, error) {
	
	var data map[string]interface{}
	var err, err1 error
	var table string

	hintID := -1

	for i, condition := range conditions {

		log.Println(condition)

		tableAttr1 := strings.Split(condition, ":")[0]
		tableAttr2 := strings.Split(condition, ":")[1]

		t1 := strings.Split(tableAttr1, ".")[0]
		a1 := strings.Split(tableAttr1, ".")[1]

		t2 := strings.Split(tableAttr2, ".")[0]
		a2 := strings.Split(tableAttr2, ".")[1]

		// log.Println(t1, a1, t2, a2)

		if i == 0 {

			// There could be mutliple pieces of data in nodes
			// For example:
			// A statuses node contains status, conversation, and status_stats
			for j, hint := range hints {

				if hint.Table == t1 {
					hintID = j
				}

			}

			if hintID == -1 {

				// In this case, since data may be incomplete, 
				// we cannot get the data in the parent node
				return nil, app_display.CannotFindAnyDataInParent
			
			} else {

				// This can happen when the data this data depends on is not migrated,
				// e.g., a post does not have correpsonding conversation in Diaspora, 
				// so when it is migrated to Mastodon, and becomes a status, 
				// it does not have conversation_id,  
				// which is actually necessary for each status in Mastodon.
				if hints[hintID].Data[t1 + "." + a1] == nil {

					return nil, app_display.CannotFindAnyDataInParent

				}

				query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s", 
					t2, a2, fmt.Sprint(hints[hintID].Data[a1]))

				data, err = db.DataCall1(displayConfig.AppConfig.DBConn, query)
				if err != nil {
					log.Fatal(err)
				}
			
				// log.Println(".....first check......")
				// log.Println(data)
				// log.Println("...........")

				if len(data) == 0 {
					return nil, app_display.CannotFindAnyDataInParent
				}

				table = t2
			}

		// This is mainly to solve the case in which
		// conversation cannot directly depend on root
		// conversation depends on statuses, which in turn depends on root. 
		// This is now obsolete because there is no dependency between other nodes with root
		// For now, there is always only one condition.
		} else {

			query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s", 
				t2, a2, fmt.Sprint(data[a1]))

			data, err1 = db.DataCall1(displayConfig.AppConfig.DBConn, query)
			if err1 != nil {
				log.Fatal(err1)
			}

			if len(data) == 0 {
				return nil, app_display.CannotFindAnyDataInParent
			}

			table = t2
		}
	}

	log.Println("...........")
	log.Println(data)
	log.Println("...........")

	return app_display.TransformRowToHint(displayConfig, data, table), nil

}

// func getHintsInParentNode(displayConfig *config.DisplayConfig, 
// 	hints []*app_display.HintStruct, conditions []string) (*app_display.HintStruct, error) {
	
// 	query := fmt.Sprintf("SELECT %s.* FROM ", "t"+strconv.Itoa(len(conditions)))
// 	from := ""
// 	table := ""
// 	hintID := -1

// 	for i, condition := range conditions {

// 		tableAttr1 := strings.Split(condition, ":")[0]
// 		tableAttr2 := strings.Split(condition, ":")[1]

// 		t1 := strings.Split(tableAttr1, ".")[0]
// 		a1 := strings.Split(tableAttr1, ".")[1]

// 		t2 := strings.Split(tableAttr2, ".")[0]
// 		a2 := strings.Split(tableAttr2, ".")[1]

// 		seq1 := "t" + strconv.Itoa(i)
// 		seq2 := "t" + strconv.Itoa(i+1)

// 		if i == 0 {

// 			// There could be mutliple pieces of data in nodes
// 			// For example:
// 			// A statuses node contains status, conversation, and status_stats
// 			for j, hint := range hints {

// 				if hint.Table == t1 {
// 					hintID = j
// 				}

// 			}

// 			// In this case, since data may be incomplete, 
// 			// we cannot get the data in the parent node
// 			if hintID == -1 {

// 				return nil, app_display.CannotFindAnyDataInParent

// 			} else {
				
// 				// For example:
// 				// if a condition is [favourites.status_id:statuses.id], 
// 				// from will be "favourites t0 JOIN statuses t1 ON t0.status_id = t1.id"
// 				from += fmt.Sprintf("%s %s JOIN %s %s ON %s.%s = %s.%s ",
// 					t1, seq1, t2, seq2, seq1, a1, seq2, a2)

// 				checkResolveReferenceInGetDataInParentNode()

// 			}

// 		// This is mainly to solve the case in which
// 		// conversation cannot directly depend on root
// 		// conversation depends on statuses, which in turn depends on root. 
// 		// This is now obsolete because there is no dependency between other nodes with root
// 		// For now, there is always only one condition.
// 		} else {

// 			from += fmt.Sprintf("JOIN %s %s on %s.%s = %s.%s ",
// 				t2, seq2, seq1, a1, seq2, a2)
			
// 			checkResolveReferenceInGetDataInParentNode()

// 		}

// 		//The last condition
// 		if i == len(conditions)-1 {

// 			var depDataKey string
// 			var depDataValue int

// 			for k, v := range hints[hintID].KeyVal {

// 				depDataKey = k
// 				depDataValue = v

// 			}

// 			// Following the above example,
// 			// the whole query will be:
// 			// SELECT t1.* 
// 			// FROM favourites t0 JOIN statuses t1 ON t0.status_id = t1.id
// 			// WHERE t0.status_id = 80
// 			where := fmt.Sprintf("WHERE %s.%s = %d", "t0", depDataKey, depDataValue)
// 			table = t2
// 			query += from + where

// 		}
// 	}
// 	// fmt.Println(query)

// 	data, err := db.DataCall1(displayConfig.AppConfig.DBConn, query)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// fmt.Println(data)
// 	if len(data) == 0 {

// 		return nil, app_display.CannotFindAnyDataInParent

// 	} else {

// 		return app_display.TransformRowToHint(displayConfig, data, table), nil

// 	}
// }

func replaceKey(displayConfig *config.DisplayConfig, tag string, key string) string {

	for _, tag1 := range displayConfig.AppConfig.Tags {

		if tag1.Name == tag {
			// fmt.Println(tag)

			for k, v := range tag1.Keys {

				if k == key {

					member := strings.Split(v, ".")[0]
					
					attr := strings.Split(v, ".")[1]
					
					for k1, table := range tag1.Members {

						if k1 == member {

							return table + "." + attr
						}
					}
				}
			}
		}
	}

	return ""

}

func dataFromParentNodeExists(displayConfig *config.DisplayConfig,
	hints []*app_display.HintStruct, pTag string) (bool, error) {
	
	displayExistenceSetting, _ := hints[0].GetDisplayExistenceSetting(displayConfig, pTag)

	// If display existence setting is not set, 
	// then we have to try to get data in the parent node in any case
	if displayExistenceSetting == "" {

		return true, nil

	} else {

		tag, _ := hints[0].GetTagName(displayConfig)
		tableCol := replaceKey(displayConfig, tag, displayExistenceSetting)
		table := strings.Split(tableCol, ".")[0]

		for _, hint := range hints {

			if hint.Table == table {

				if hint.Data[tableCol] == nil {

					return false, app_display.NotDependsOnAnyData

				} else {

					return true, nil

				}
			}
		}

	}

	// In this case, since data may be incomplete, 
	// we cannot find the existence of the data in a parent node
	// This also implies that it cannot find any data in a parent node
	return false, app_display.CannotFindAnyDataInParent

}

// Note: this function may return multiple hints based on dependencies
func GetdataFromParentNode(displayConfig *config.DisplayConfig,
	hints []*app_display.HintStruct, pTag string) (*app_display.HintStruct, error) {

	// Before getting data from a parent node, 
	// we check the existence of the data based on the cols of a child node
	if exists, err := dataFromParentNodeExists(displayConfig, hints, pTag); !exists {
		return nil, err
	}

	tag, _ := hints[0].GetTagName(displayConfig)
	conditions, _ := displayConfig.AppConfig.GetDependsOnConditions(tag, pTag)
	pTag, _ = hints[0].GetOriginalTagNameFromAliasOfParentTagIfExists(displayConfig, pTag)

	var procConditions []string
	var from, to string

	if len(conditions) == 1 {

		condition := conditions[0]
		from = replaceKey(displayConfig, tag, condition.TagAttr)
		to = replaceKey(displayConfig, pTag, condition.DependsOnAttr)
		procConditions = append(procConditions, from+":"+to)

	} else {

		for i, condition := range conditions {

			if i == 0 {

				from = replaceKey(displayConfig, tag, condition.TagAttr)

				to = replaceKey(displayConfig,
					strings.Split(condition.DependsOnAttr, ".")[0], 
					strings.Split(condition.DependsOnAttr, ".")[1])

			} else if i == len(conditions)-1 {

				from = replaceKey(displayConfig, 
					strings.Split(condition.TagAttr, ".")[0], 
					strings.Split(condition.TagAttr, ".")[1])
				
				to = replaceKey(displayConfig, pTag, condition.DependsOnAttr)

			}

			procConditions = append(procConditions, from+":"+to)

		}

	}

	// fmt.Println(procConditions)
	// fmt.Println(hints)

	return getHintsInParentNode(displayConfig, hints, procConditions)

}

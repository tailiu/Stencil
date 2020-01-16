package SA1_display

import (
	"errors"
	"fmt"
	"log"
	"stencil/db"
	"stencil/reference_resolution"
	"stencil/schema_mappings"
	"strconv"
	"strings"
)

func getOneRowBasedOnDependency(displayConfig *displayConfig,
	table, col, value string) (map[string]interface{}, error) {

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s", table, col, value)
	log.Println(query)

	data, err := db.DataCall1(displayConfig.dstAppConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)
	if len(data) == 0 {

		return nil, CannotFindRemainingData

	} else {

		return data, nil

	}
}

func checkResolveReferenceInGetDataInNode(displayConfig *displayConfig,
	id, table0, col0, table1, col1, value string) (map[string]interface{}, error) {

	// We use table0 and col0 to get table1 and col1
	log.Println("+++++++++++++++++++")
	log.Println("id:", id)
	log.Println(table0)
	log.Println(col0)
	log.Println("value:", value)
	log.Println(table1)
	log.Println(col1)
	log.Println("+++++++++++++++++++")
	
	table0ID := displayConfig.dstAppConfig.tableNameIDPairs[table0]
	table1ID := displayConfig.dstAppConfig.tableNameIDPairs[table1]

	// First, we need to get the attribute that requires reference resolution
	// For example, we have *account.id*, and we want to get *users.account_id*
	// We check whether account.id needs to be resolved
	if reference_resolution.NeedToResolveReference(displayConfig.refResolutionConfig, 
		table0, col0) {

		log.Println("Before checking reference1 resolved or not")

		// If account.id should be resolved (in this case, it should not),
		// we check whether the reference has been resolved or not
		newVal := reference_resolution.ReferenceResolved(displayConfig.refResolutionConfig, 
			table0ID, col0, id)
		
		// If the reference has been resolved, then use the new reference to get data
		if newVal != "" {

			log.Println("reference has been resolved")

			return getOneRowBasedOnDependency(displayConfig, table1, col1, newVal)
		
		// Otherwise, we try to resolve the reference
		} else {

			hint0 := CreateHint(table0, table0ID, id)
			log.Println("Before resolving reference1: ", hint0)

			ID0 := hint0.TransformHintToIdenity(displayConfig)

			updatedAttrs, _ := reference_resolution.ResolveReference(
				displayConfig.refResolutionConfig, ID0)

			// We check whether the desired attr (col0) has been resolved
			foundResolvedAttr := false
			for attr, val := range updatedAttrs {
				if attr == col0 {
					newVal = val
					foundResolvedAttr = true
					break
				}
			}

			// If we find that col0 has been resolved, then we can use it to get other data
			if foundResolvedAttr {

				return getOneRowBasedOnDependency(displayConfig, table1, col1, newVal)
			
			// Otherwise we cannot use the unresolved reference to get other data in node
			} else {

				return nil, CannotResolveReferencesGetDataInNode
			}
		}

	// We check if users.account_id needs be resolved (of course, in this case, it should be)
	// However we don't know its id (this is the differece from the above case!!). 
	// Also if the value is the value of "id", this could not be used directly 
	} else if fromAttrsfirstArg := schema_mappings.GetFirstArgsInREFByToTableToAttr(
		displayConfig.mappingsFromSrcToDst, table1, col1); len(fromAttrsfirstArg) != 0 {

		log.Println("Before checking reference2 resolved or not")
		
		log.Println("From attributes:")
		log.Println(fromAttrsfirstArg)
		
		data := make(map[string]interface{})
		var err error

		fromAttrfirstArgContainID := false

		// Because we only use toTable and toAttr to get the first argument in fromAttrs,
		// there could be multiple results. For example, toTable = status_stats and
		// toAttr = status_id, then fromAttr could be posts.id, comments.id, or messages.id
		for fromAttrfirstArg, _ := range fromAttrsfirstArg {

			log.Println("Check a from attribute:", fromAttrfirstArg)

			// If the first argument of the from attribute contains "id", this indicates 
			// we need to get the original id as the value to get the data
			// (the current id is the newly generated one)
			if doesArgAttributeContainID(fromAttrfirstArg) {

				fromAttrfirstArgContainID = true

				// log.Println(displayConfig.dstAppConfig.tableNameIDPairs)
				// log.Println(table0)
				// log.Println(table0ID)
				
				dataID := reference_resolution.CreateIdentity(
					displayConfig.dstAppConfig.appID,
					table0ID,
					id,
				)
				
				// tableInFirstArg := getTableInArg(fromAttrfirstArg)
				// srcTableID := displayConfig.srcAppConfig.tableNameIDPairs[tableInFirstArg]
	
				prevID := reference_resolution.GetPreviousID(displayConfig.refResolutionConfig, 
					dataID)
				
				log.Println("Previous id:", prevID)

				// since there is only one mapping to this toAttr, as long as we find one, 
				// we can set the value as the prevID
				if prevID == "" {
					continue
				} else {
					// value = prevID
					// break

					data, err = getOneRowBasedOnDependency(displayConfig, table1, col1, prevID)
					if err != nil {
						log.Println("The first argument of the from attribute contains id, but")
						log.Println(err)
						break
					}
				}
			}
		} 

		// If the first argument of the from attribute does not contain "id", 
		// this indicates we can use the current data and the relationship indicated
		// by table0, col0, table1, col1, value to get data
		if !fromAttrfirstArgContainID {

			log.Println(`The from attributes don't contain id`)

			data, err = getOneRowBasedOnDependency(displayConfig, table1, col1, value)
			// This could happen when no data is migrated or there is no mappings.
			// For example, statuses, id, mentions, status_id
			if err != nil {
				return nil, err
			}

		}

		// If the first argument of the from attribute contains id and
		// we cannot get data, there could be two cases:
		// 1. The reference has been resolved, so the data contains the up-to-date value
		// 2. The reference has been resolved, but reference resolution crashes before
		// inserting the reference into the resolution resolved table
		// In both cases, try to get data with the new value 
		// For 1, it will be checked afterwards
		// For 2, do the reference resolution again since it does not matter and in the
		// second time, we can remove the reference and add it to the resolved resolution table
		if fromAttrfirstArgContainID && data == nil {

			log.Println(`The from attributes contain id but we cannot get data,
				so we try to get data with the current id value`)

			data, err = getOneRowBasedOnDependency(displayConfig, table1, col1, value)
			// This could happen when the resolved and displayed data is deleted
			if err != nil {
				return nil, err
			}
			
		}

		// Now we have the id of the data, we should check whether it has been resolved before, 
		// but actually if we can get one, it should always be the one we want to get because
		// otherwise there will be multiple rows corresponding to one member.
		// There could be the case where ids are not changed, 
		// so even if references are not resolved, 
		// we can still get the rows we want, but we need to resolve it further.
		newVal := reference_resolution.ReferenceResolved(displayConfig.refResolutionConfig, 
			table1ID, col1, fmt.Sprint(data["id"]))

		// The reference has been resolved
		if newVal != "" {

			// Theoretically, if it has been resolved, then it should be the value we have 
			// given that one member corresponds to one row
			if newVal == value {

				log.Println("reference has been resolved")
				log.Println(data)

				return data, nil

			} else {

				panic(`The reference has been resolved, but the value is not what we want. 
					Should not happen given one member corresponds to one row for now!`)
			
			}
		} else {
			
			hint1 := CreateHint(table1, table1ID, fmt.Sprint(data["id"]))
			log.Println("Before resolving reference2:", hint1)

			ID1 := hint1.TransformHintToIdenity(displayConfig)

			updatedAttrs, _ := reference_resolution.ResolveReference(
				displayConfig.refResolutionConfig, ID1)
			
			log.Println("Updated attributes:")
			log.Println(updatedAttrs)

			// We check whether the desired attr (col1) has been resolved 
			// (until this point, it should be resolved)
			foundResolvedAttr := false
			for attr, val := range updatedAttrs {
				if attr == col1 {
					newVal = val
					foundResolvedAttr = true
					break
				}
			}

			// If we find that col0 has been resolved, then we can use it to get other data
			if foundResolvedAttr {

				if newVal == value {

					return data, nil
	
				} else {
	
					panic(`Find the resolved attribute, but the value is not what we want. 
						Should not happen given one member corresponds to one row for now!`)
				}
			
			// This should not happen
			} else {
				
				panic(`Does not find resolved attributes. Should not happen 
					given one member corresponds to one row for now!`)
			}
		}
	
	// Normally, there must exist one that needs to be resolved. 
	// However, up to now there is a case breaking the above rule. 
	// When there is no mapping
	// For example: 
	// When migrating from Diaspora to Mastodon, 
	// there is no mapping to stream_entries.activity_id.
	} else {

		return nil, NoReferenceToResolve
	}
}

func getRemainingDataInNode(displayConfig *displayConfig,
	dependencies []map[string]string, 
	members map[string]string, 
	hint *HintStruct) ([]*HintStruct, error) {
	
	var result []*HintStruct

	procDependencies := make(map[string][]string)

	for _, dependency := range dependencies {

		for k, v := range dependency {

			memberSeqInKey := strings.Split(k, ".")[0]
			memberSeqInVal := strings.Split(v, ".")[0]

			memberTableInKey := members[memberSeqInKey]
			memberTableInVal := members[memberSeqInVal]

			newKey := strings.Replace(k, memberSeqInKey, memberTableInKey, 1)
			newVal := strings.Replace(v, memberSeqInVal, memberTableInVal, 1)

			procDependencies[newKey] = append(procDependencies[newKey], newVal)
			procDependencies[newVal] = append(procDependencies[newVal], newKey)

		}
	}

	// fmt.Println(procDependencies)

	tag := hint.Tag

	result = append(result, hint)

	queue := []DataInDependencyNode{DataInDependencyNode{
		Table: hint.Table,
		Data:  hint.Data,
	}}

	for len(queue) != 0 && len(procDependencies) != 0 {
		
		// fmt.Println(queue)
		// fmt.Println(procDependencies)

		dataInDependencyNode := queue[0]

		queue = queue[1:]

		table := dataInDependencyNode.Table

		for col, val := range dataInDependencyNode.Data {

			if deps, ok := procDependencies[table+"."+col]; ok {

				for _, dep := range deps {

					// fmt.Println(dep)

					table1 := strings.Split(dep, ".")[0]
					key1 := strings.Split(dep, ".")[1]

					var data map[string]interface{} 
					var err1 error

					// If resolving reference is required
					if displayConfig.resolveReference {

						// We assume that val is an integer value 
						// otherwise we have to define it in dependency config
						data, err1 = checkResolveReferenceInGetDataInNode(displayConfig, 
							fmt.Sprint(dataInDependencyNode.Data["id"]),
							table, col, table1, key1, fmt.Sprint(val))

					} else {

						data, err1 = getOneRowBasedOnDependency(
							displayConfig, table1, key1, fmt.Sprint(val))
						
					}

					// fmt.Println(data)

					if err1 != nil {
						log.Println(err1)
						// fmt.Println(result)
						continue
					}

					queue = append(queue, DataInDependencyNode{
						Table: table1,
						Data:  data,
					})
					// fmt.Println(queue)

					intPK, err2 := strconv.Atoi(fmt.Sprint(data["id"]))
					if err2 != nil {
						log.Fatal(err2)
					}
					keyVal := map[string]int{
						"id": intPK,
					}

					result = append(result, &HintStruct{
						Table: table1,
						TableID: displayConfig.dstAppConfig.tableNameIDPairs[table1],
						KeyVal: keyVal,
						Data: data,
						Tag: tag,
					})

					deps1 := procDependencies[table1+"."+key1]
					for i, val2 := range deps1 {
						if val2 == table+"."+col {
							deps1 = append(deps1[:i], deps1[i+1:]...)
							break
						}
					}
					if len(deps1) == 0 {

						delete(procDependencies, table1+"."+key1)

					} else {

						procDependencies[table1+"."+key1] = deps1

					}
				}

				delete(procDependencies, table+"."+col)

			}
		}
	}

	// fmt.Println(procDependencies)
	// fmt.Println(result)
	if len(procDependencies) == 0 {

		return result, nil

	} else {

		return result, NodeIncomplete

	}

}

func getOneRowBasedOnHint(displayConfig *displayConfig, 
	hint *HintStruct) (map[string]interface{}, error) {
	
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = %d", hint.Table, hint.KeyVal["id"])
	
	// log.Println(query)
	
	data, err := db.DataCall1(displayConfig.dstAppConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if len(data) == 0 {

		return nil, DataNotExists

	} else {

		return data, nil

	}

}

func getDataInNode(displayConfig *displayConfig, 
	hint *HintStruct) ([]*HintStruct, error) {

	if hint.Data == nil {

		data, err := getOneRowBasedOnHint(displayConfig, hint)
		if err != nil {

			return nil, err

		} else {

			hint.Data = data

		}
	}
	
	for _, tag := range displayConfig.dstAppConfig.dag.Tags {

		for _, member := range tag.Members {

			if hint.Table == member {

				if len(tag.Members) == 1 {

					return []*HintStruct{hint}, nil

				} else {

					// Note: we assume that one dependency represents that one row
					// 		in one table depends on another row in another table
					return getRemainingDataInNode(displayConfig,
						 tag.InnerDependencies, tag.Members, hint)

				}
			}
		}
	}

	return nil, errors.New("Error: the hint does not match any tags")

}

// A recursive function checks whether all the data one data recursively depends on exists
// We only checks whether the table depended on exists, which is sufficient for now
func checkDependsOnExists(displayConfig *displayConfig, 
	allData []*HintStruct, data *HintStruct) bool {
	
	memberID, _ := data.GetMemberID(displayConfig)
	// fmt.Println(memberID)

	dependsOnTables := data.GetDependsOnTables(displayConfig, memberID)
	// fmt.Println(dependsOnTables)
	
	if len(dependsOnTables) == 0 {

		return true

	} else {

		for _, dependsOnTable := range dependsOnTables {

			exists := false

			for _, oneData := range allData {

				if oneData.Table == dependsOnTable {

					if !checkDependsOnExists(displayConfig, allData, oneData) {

						return false

					} else {

						exists = true
						break

					}
				}
			}
			if !exists {

				return false

			}
		}
	}

	return true
	
}

func trimDataBasedOnInnerDependencies(displayConfig *displayConfig,
	 allData []*HintStruct) []*HintStruct {
	
	var trimmedData []*HintStruct

	for _, data := range allData {

		if checkDependsOnExists(displayConfig, allData, data) {

			trimmedData = append(trimmedData, data)
			
		}

	}

	return trimmedData

}

func GetDataInNodeBasedOnDisplaySetting(displayConfig *displayConfig, 
	hint *HintStruct) ([]*HintStruct, error) {
	
	displaySetting, _ := hint.GetTagDisplaySetting(displayConfig)

	// Whether a node is complete or not, get all the data in a node.
	// If the node is complete, err is nil, otherwise, err is "node is not complete".
	if data, err := getDataInNode(displayConfig, hint); err != nil {

		// The setting "default_display_setting" means only display a node 
		// when the node is complete.
		// Therefore, return nil and error message when node is not complete.
		if displaySetting == "default_display_setting" {
			return nil, err

		// The setting "display_based_on_inner_dependencies" means 
		// display as much data in a node as possible based on inner dependencies.
		// Note: if a piece of data in a node depends on some data not existing in the node,
		// it needs to be deleted from the data set and cannot be displayed.
		} else if displaySetting == "display_based_on_inner_dependencies" {

			// fmt.Println(data)
			return trimDataBasedOnInnerDependencies(displayConfig, data), err

		}
	
		// If a node is complete, return all the data in the node regardless of the setting.
	} else {

		return data, nil

	}

	panic("Should never happen")

}

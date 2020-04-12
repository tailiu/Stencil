package SA1_display

import (
	"errors"
	"fmt"
	"log"
	"stencil/db"
	"stencil/common_funcs"
	"stencil/reference_resolution_v2"
	"strconv"
	"strings"
)

func (display *display) getOneRowBasedOnDependency(table, col, value string) (map[string]interface{}, error) {

	var query string
	
	if !display.markAsDelete {
		query = fmt.Sprintf(
			`SELECT * FROM "%s" WHERE %s = '%s'`, 
			table, col, value,
		)
	} else {
		query = fmt.Sprintf(
			`SELECT * FROM "%s" WHERE %s = '%s' and mark_as_delete = false`, 
			table, col, value,
		)
	}
	
	log.Println(query)

	data, err := db.DataCall1(display.dstAppConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)
	if len(data) == 0 {

		return nil, common_funcs.CannotFindRemainingData

	} else {

		return data, nil

	}
}

func (display *display) getRowsBasedOnDependency(table, col, value string) ([]map[string]interface{}, error) {

	var query string

	if !display.markAsDelete {
		query = fmt.Sprintf(
			`SELECT * FROM "%s" WHERE %s = '%s'`, 
			table, col, value,
		)
	} else {
		query = fmt.Sprintf(
			`SELECT * FROM "%s" WHERE %s = '%s' and mark_as_delete = false`, 
			table, col, value,
		)
	}
	
	// log.Println(query)

	data, err := db.DataCall(display.dstAppConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)
	if len(data) == 0 {

		return nil, common_funcs.CannotFindRemainingData

	} else {

		return data, nil

	}
}

func (display *display) checkReferenceIndeedResolved(table, col, 
	tableID, colID, value string) (map[string]interface{}, error) {

	// First we must assume that it has already been resolved. If it has not been resolved,
	// then we cannot get data. Otherwise we just return the obtained data
	// Note that if we first assume that it has not been resolved and get data using 
	// prevID, then we could get wrong results
	data, err := display.getRowsBasedOnDependency(table, col, value)
	
	// This could happen when table1 and col1 have been resolved
	if err == nil {

		log.Println("The reference seems to have already been resolved")
		
		// Now we have not encountered data1 with more than one piece of data
		for _, data1 := range data {

			resolvedVal := display.rr.ReferenceResolved(tableID, colID, fmt.Sprint(data1["id"]))
			
			if resolvedVal == value {
				log.Println("It was indeed resolved")
				return data1, nil 
			
			// This could happen when there was some data happening to satisfy the condition,
			// but that data actually does not have relationships with table0 and col0
			// For example, table0: accounts, col0: id, table1: users, col1: account_id
			// we can get a data with users.account_id pointing to accounts.id, but 
			// the users.account_id of that data is actually old application value and the data
			// may be from other migrations.
			} else {
				log.Println("There happened to be some data satisfying but that data is not what we want")
			}
		}

		log.Println("The reference has not been resolved")

		return nil, DataNotWanted

	} else {
		
		log.Println("The reference has not been resolved")

		return nil, DataNotFound
	} 
}

func (display *display) checkResolveRefWithIDInData(table, col, tableID, colID, id, attrVal string) (string, error) {

	// If favourites.status_id should be resolved (in this case, it should be),
	// we check whether the reference has been resolved or not
	newVal := display.rr.ReferenceResolved(tableID, colID, id)
	
	// If the reference has been resolved, then use the new reference to get data
	// Otherwise, we try to resolve the reference
	if newVal != "" {
		return newVal, nil	
	} else {

		attr0 := reference_resolution_v2.CreateAttribute(display.dstAppConfig.appID, tableID, colID, attrVal, id)
		log.Println("Before resolving reference: ", attr0)

		display.rr.ResolveReference(attr0)
		
		// Here we check again to get updated attributes and values
		// instead of using the returned values from the ResolveReference
		// because ResolveReference only returns the updated values in that
		// function call. Values could be updated by other threads and in this
		// case, ResolveReference does not return the updated attribute and value
		// Therefore, we check all updated attributes again here by calling 
		// GetUpdatedAttributes
		updatedAttrs := display.rr.GetUpdatedAttributes(tableID, id)

		log.Println("Updated attributes and values:")
		log.Println(updatedAttrs)
		
		// We check whether the desired attr (col0) has been resolved
		foundResolvedAttr := false
		for attr, val := range updatedAttrs {
			if attr == colID {
				newVal = val
				foundResolvedAttr = true
				break
			}
		}

		// If we find that col has been resolved, then we can use it to get other data
		// Otherwise we cannot use the unresolved reference to get other data in node
		if foundResolvedAttr {
			return newVal, nil
		} else {
			return "", CannotResolveRefersWithIDInData
		}
	}
}

func (display *display) checkResolveReferenceInGetDataInNode(
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
	
	table0ID := display.dstAppConfig.tableNameIDPairs[table0]
	table1ID := display.dstAppConfig.tableNameIDPairs[table1]

	col0ID := display.dstAppConfig.colNameIDPairs[table0 + ":" + col0]
	col1ID := display.dstAppConfig.colNameIDPairs[table1 + ":" + col1]

	// First, we need to get the attribute that requires reference resolution
	// For example, we have *account.id*, and we want to get *users.account_id*
	// We check whether account.id needs to be resolved
	if display.needToResolveReference(table0, col0) {

		log.Println("Checking reference1 resolved or not")

		newVal, err := display.checkResolveRefWithIDInData(table0, col0, table0ID, col0ID, id, value)

		// If account.id should be resolved (in this case, it should not),
		// we check whether the reference has been resolved or not
		// newVal := display.rr.ReferenceResolved(table0ID, col0ID, id)
		
		// If the reference has been resolved, then use the new reference to get data
		if newVal != "" {
			log.Println("reference1 has been resolved")
			return display.getOneRowBasedOnDependency(table1, col1, newVal)
		} else {
			return nil, err
		}

	// We check if users.account_id needs to be resolved (of course, in this case, it should be)
	// However we don't know its id (this is the differece from the above case!!). 
	// Also if the value is the value of "id", this could not be used directly 
	// We decide to make so much efforts to resolve "backwards" because inner-dependencies like
	// "statuses.id":"statuses.status_id"
	// "statuses.id":"mentions.status_id"
	// "statuses.id":"stream_entries.activity_id"
	// force us to do in this way. Otherwise, we cannot get other data in a node through statuses.id
	} else if display.needToResolveReference(table1, col1) {

		log.Println("Checking reference2 resolved or not")
		
		data1, err1 := display.checkReferenceIndeedResolved(table1, col1, table1ID, col1ID, value)
		if err1 != nil {
			log.Println(err1)
		} else {
			return data1, nil
		}
		
		attr1 := reference_resolution_v2.CreateAttribute(display.dstAppConfig.appID, table0ID, col0ID, id, id)
		log.Println("Resolving reference2: ", attr1)

		display.rr.ResolveReference(attr1)
		
		data2, err2 := display.checkReferenceIndeedResolved(table1, col1, table1ID, col1ID, value)
		if err2 != nil {
			return nil, CannotGetDataAfterResolvingRef2
		} else {
			return data2, nil
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

func (display *display) getRemainingDataInNode(dependencies []map[string]string, 
	members map[string]string, hint *HintStruct) ([]*HintStruct, error) {
	
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

	queue := []common_funcs.DataInDependencyNode{
		common_funcs.DataInDependencyNode{
			Table: hint.Table,
			Data:  hint.Data,
		},
	}

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
					if display.resolveReference {

						// We assume that val is an integer value 
						// otherwise we have to define it in dependency config
						data, err1 = display.checkResolveReferenceInGetDataInNode(
							fmt.Sprint(dataInDependencyNode.Data["id"]),
							table, col, table1, key1, fmt.Sprint(val),
						)

					} else {
						data, err1 = display.getOneRowBasedOnDependency(table1, key1, fmt.Sprint(val))
					}

					// fmt.Println(data)

					if err1 != nil {
						log.Println(err1)
						// fmt.Println(result)
						continue
					}

					queue = append(
						queue, 
						common_funcs.DataInDependencyNode{
							Table: table1,
							Data:  data,
						},
					)
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
						TableID: display.dstAppConfig.tableNameIDPairs[table1],
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

		return result, common_funcs.NodeIncomplete

	}

}

func (display *display) getOneRowBasedOnHint(hint *HintStruct) (map[string]interface{}, error) {
	
	var query string

	if !display.markAsDelete {
		query = fmt.Sprintf(
			`SELECT * FROM "%s" WHERE id = %d`, 
			hint.Table, hint.KeyVal["id"])
	} else {
		query = fmt.Sprintf(
			`SELECT * FROM "%s" WHERE id = %d and mark_as_delete = false`, 
			hint.Table, hint.KeyVal["id"])
	}
	
	// log.Println(query)
	
	data, err := db.DataCall1(display.dstAppConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if len(data) == 0 {

		return nil, common_funcs.DataNotExists

	} else {

		return data, nil

	}

}

func (display *display) getDataInNode(hint *HintStruct) ([]*HintStruct, error) {

	if hint.Data == nil {

		data, err := display.getOneRowBasedOnHint(hint)
		if err != nil {
			return nil, err
		} else {

			hint.Data = data

		}
	}

	// log.Println("My data is:", hint.Data)
	
	for _, tag := range display.dstAppConfig.dag.Tags {

		for _, member := range tag.Members {

			if hint.Table == member {

				if len(tag.Members) == 1 {

					return []*HintStruct{hint}, nil

				} else {

					// Note: we assume that one dependency represents that one row
					// 		in one table depends on another row in another table
					hints, err1 := display.getRemainingDataInNode(
						 tag.InnerDependencies, tag.Members, hint)
					
					// Refresh the cached results which could have changed due to
					// reference resolution 
					display.refreshCachedDataHints(hints)

					return hints, err1
				}
			}
		}
	}

	return nil, errors.New("Error: the hint does not match any tags")

}

// A recursive function checks whether all the data one data recursively depends on exists
// We only checks whether the table depended on exists, which is sufficient for now
func (display *display) checkDependsOnExists(allData []*HintStruct, data *HintStruct) bool {
	
	memberID, _ := data.GetMemberID(display)
	// fmt.Println(memberID)

	dependsOnTables := data.GetDependsOnTables(display, memberID)
	// fmt.Println(dependsOnTables)
	
	if len(dependsOnTables) == 0 {

		return true

	} else {

		for _, dependsOnTable := range dependsOnTables {

			exists := false

			for _, oneData := range allData {

				if oneData.Table == dependsOnTable {

					if !display.checkDependsOnExists(allData, oneData) {
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

func (display *display) trimDataBasedOnInnerDependencies(allData []*HintStruct) []*HintStruct {
	
	var trimmedData []*HintStruct

	for _, data := range allData {
		if display.checkDependsOnExists(allData, data) {
			trimmedData = append(trimmedData, data)
		}
	}

	return trimmedData
}

func (display *display) GetDataInNodeBasedOnDisplaySetting(hint *HintStruct) ([]*HintStruct, error) {
	
	displaySetting, _ := hint.GetTagDisplaySetting(display)

	// Whether a node is complete or not, get all the data in a node.
	// If the node is complete, err is nil, otherwise, err is "node is not complete".
	if data, err := display.getDataInNode(hint); err != nil {

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
			return display.trimDataBasedOnInnerDependencies(data), err

		}
	
		// If a node is complete, return all the data in the node regardless of the setting.
	} else {

		return data, nil

	}

	panic("Should never happen")

}

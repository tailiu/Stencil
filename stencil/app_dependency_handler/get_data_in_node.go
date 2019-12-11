package app_dependency_handler

import (
	"errors"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/app_display"
	"stencil/reference_resolution"
	"strconv"
	"strings"
)

func getOneRowBasedOnDependency(displayConfig *config.DisplayConfig,
	table, col, value string) (map[string]interface{}, error) {

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s", table, col, value)
	// fmt.Println(query)

	data, err := db.DataCall1(displayConfig.AppConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)
	if len(data) == 0 {

		return nil, app_display.CannotFindRemainingData

	} else {

		return data, nil

	}
}

func checkResolveReference(displayConfig *config.DisplayConfig,
	id, table0, col0, table1, col1, value string) (map[string]interface{}, error) {

	log.Println("+++++++++++++++++++")
	log.Println(table0)
	log.Println(col0)
	log.Println(table1)
	log.Println(col1)
	log.Println("+++++++++++++++++++")

	// First, we need to get the attribute that requires reference resolution
	// For example, we have account.id, and we want to get users.account_id
	// We check whether account.id needs be resolved
	if reference_resolution.NeedToResolveReference(displayConfig, table0, col0) {

		// If account.id should be resolved (in this case, it should not),
		// we check whether the reference has been resolved or not
		newVal := reference_resolution.ReferenceResolved(displayConfig, table0, col0, id)
		
		// If the reference has been resolved, then use the new reference to get data
		if newVal != "" {

			value = newVal
		
		// Otherwise, we try to resolve the reference
		} else {

			updatedAttrs := reference_resolution.ResolveReference(displayConfig, &hint)

			// We check whether the desired attr (col0) has been resolved
			foundResolvedAttr := false
			for _, attr := range updatedAttrs {
				if attr == col0 {
					foundResolvedAttr = true
					break
				}
			}

			// If we find that col0 has been resolved, then we can use 
			if foundResolvedAttr {

				return getOneRowBasedOnDependency(displayConfig, table1, col1, value)
			

			} else {
				return 
			}
		}

	// We check if users.account_id needs be resolved (of course, in this case, it should be)
	} else {
		
		if reference_resolution.NeedToResolveReference(displayConfig, table1, col1) {
			
		} else {

			panic("Should have at least one attribute to resolve!")
		}
	}
}

func getRemainingDataInNode(displayConfig *config.DisplayConfig,
	 dependencies []map[string]string, 
	 members map[string]string, 
	 hint *app_display.HintStruct) ([]*app_display.HintStruct, error) {
	
	var result []*app_display.HintStruct

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
					if displayConfig.ResolveReference {

						// We assume that val is an integer value 
						// otherwise we have to define it in dependency config
						checkResolveReference(displayConfig, 
							fmt.Sprint(dataInDependencyNode.Data["id"]),
							table0, col0, table1, col1, val)

					} else {

						data, err1 = getOneRowBasedOnDependency(
							displayConfig, table1, key1, val)
						
					}

					// fmt.Println(data)

					if err1 != nil {
						// fmt.Println(err1)
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

					result = append(result, &app_display.HintStruct{
						Table: table1,
						TableID: displayConfig.AppConfig.TableNameIDPairs[table1],
						KeyVal: keyVal,
						Data: data,
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

		return result, app_display.NodeIncomplete

	}

}

func getOneRowBasedOnHint(displayConfig *config.DisplayConfig, 
	hint *app_display.HintStruct) (map[string]interface{}, error) {
	
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = %d", hint.Table, hint.KeyVal["id"])
	
	// log.Println(query)
	
	data, err := db.DataCall1(displayConfig.AppConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if len(data) == 0 {

		return nil, app_display.DataNotExists

	} else {

		return data, nil

	}

}

func getDataInNode(displayConfig *config.DisplayConfig, 
	hint *app_display.HintStruct) ([]*app_display.HintStruct, error) {

	if hint.Data == nil {

		data, err := getOneRowBasedOnHint(displayConfig, hint)
		if err != nil {

			return nil, err

		} else {

			hint.Data = data

		}
	}
	
	for _, tag := range displayConfig.AppConfig.Tags {

		for _, member := range tag.Members {

			if hint.Table == member {

				if len(tag.Members) == 1 {

					return []*app_display.HintStruct{hint}, nil

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
func checkDependsOnExists(displayConfig *config.DisplayConfig, 
	allData []*app_display.HintStruct, tagName string, data *app_display.HintStruct) bool {
	
	memberID, _ := data.GetMemberID(displayConfig, tagName)
	// fmt.Println(memberID)

	dependsOnTables := displayConfig.AppConfig.GetDependsOnTables(tagName, memberID)
	// fmt.Println(dependsOnTables)
	
	if len(dependsOnTables) == 0 {

		return true

	} else {

		for _, dependsOnTable := range dependsOnTables {

			exists := false

			for _, oneData := range allData {

				if oneData.Table == dependsOnTable {

					if !checkDependsOnExists(displayConfig, allData, tagName, oneData) {

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

func trimDataBasedOnInnerDependencies(displayConfig *config.DisplayConfig,
	 allData []*app_display.HintStruct, tagName string) []*app_display.HintStruct {
	
	var trimmedData []*app_display.HintStruct

	for _, data := range allData {

		if checkDependsOnExists(displayConfig, allData, tagName, data) {

			trimmedData = append(trimmedData, data)
			
		}

	}

	return trimmedData

}

func GetDataInNodeBasedOnDisplaySetting(displayConfig *config.DisplayConfig, 
	hint *app_display.HintStruct) ([]*app_display.HintStruct, error) {
	
	var data []*app_display.HintStruct

	tagName, err := hint.GetTagName(displayConfig)
	if err != nil {
		return nil, err
	}

	displaySetting, _ := displayConfig.AppConfig.GetTagDisplaySetting(tagName)

	// Whether a node is complete or not, get all the data in a node.
	// If the node is complete, err is nil, otherwise, err is "node is not complete".
	if data, err = getDataInNode(displayConfig, hint); err != nil {

		// The setting "default_display_setting" means only display a node when the node is complete.
		// Therefore, return nil and error message when node is not complete.
		if displaySetting == "default_display_setting" {
			return nil, err

		// The setting "display_based_on_inner_dependencies" means 
		// display as much data in a node as possible based on inner dependencies.
		// Note: if a piece of data in a node depends on some data not existing in the node,
		// it needs to be deleted from the data set and cannot be displayed.
		} else if displaySetting == "display_based_on_inner_dependencies" {

			// fmt.Println(data)
			return trimDataBasedOnInnerDependencies(displayConfig, data, tagName), err

		}
	
		// If a node is complete, return all the data in the node regardless of the setting.
	} else {

		return data, nil

	}

	panic("Should never happen")

}

package SA1_display

import (
	"errors"
	"fmt"
	"log"
	"stencil/db"
	"stencil/common_funcs"
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

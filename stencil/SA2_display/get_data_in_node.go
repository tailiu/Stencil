package SA2_display

import (
	"errors"
	"log"
	"stencil/common_funcs"
	"strings"
	"fmt"
)

func (displayConfig *displayConfig) getOneRowBasedOnDependency(val string, dep string) (map[string]interface{}, error) {
	
	table := strings.Split(dep, ".")[0]
	key := strings.Split(dep, ".")[1]

	// log.Println(table)
	// log.Println(key)
	// log.Println(val)
	
	data := GetData1FromPhysicalSchema(
		displayConfig, 
		table + ".*", table, 
		table + "." + key, "=", val,
	)

	if len(data) == 0 {
		return nil, common_funcs.CannotFindRemainingData
	} else {
		return data, nil
	}
}

func (displayConfig *displayConfig) getRemainingDataInNode(dependencies []map[string]string, 
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
	// log.Println(procDependencies)
	
	tag := hint.Tag

	result = append(result, hint)

	queue := []common_funcs.DataInDependencyNode {
		common_funcs.DataInDependencyNode {
			Table: hint.TableName,
			Data:  hint.Data,
		},
	}

	for len(queue) != 0 && len(procDependencies) != 0 {
		// log.Println(queue)
		// log.Println(procDependencies)

		dataInDependencyNode := queue[0]
		queue = queue[1:]

		// table := dataInDependencyNode.Table
		for tableCol, val := range dataInDependencyNode.Data {

			if deps, ok := procDependencies[tableCol]; ok {
				// We assume that this is an integer value 
				// otherwise we have to define it in dependency config
				for _, dep := range deps {

					// log.Println(dep)
					// log.Println(tableCol)
					// log.Println(dataInDependencyNode.Data)
					if val == nil {
						log.Println("Fail to get one data because the value of the relevant column is nil")
						continue
					}
					data1, err1 := displayConfig.getOneRowBasedOnDependency(
						fmt.Sprint(val), 
						dep,
					)
					
					if err1 != nil {
						// log.Println(err1)
						// fmt.Println(result)
						continue
					}

					table1 := strings.Split(dep, ".")[0]
					key1 := strings.Split(dep, ".")[1]
					queue = append(
						queue, 
						common_funcs.DataInDependencyNode{
							Table: table1,
							Data:  data1,
						},
					)

					result = append(result, &HintStruct{
						TableName: table1,
						TableID: displayConfig.dstAppConfig.tableNameIDPairs[table1],
						RowIDs: GetRowIDsFromData(data1),
						Data: data1,
						Tag: tag, 
					})

					deps1 := procDependencies[table1+"."+key1]

					// log.Println("before delete: ", deps1)
					// log.Println("to delete: ", tableCol)
					
					for i, val2 := range deps1 {
						if val2 == tableCol {
							deps1 = append(deps1[:i], deps1[i+1:]...)
							break
						}
					}

					// log.Println("after delete: ", deps1)
					
					if len(deps1) == 0 {
						delete(procDependencies, table1+"."+key1)
					} else {
						procDependencies[table1+"."+key1] = deps1
					}
				}

				delete(procDependencies, tableCol)
			}
		}
	}

	// log.Println(procDependencies)
	// log.Println(result)

	if len(procDependencies) == 0 {
		return result, nil
	} else {
		return result, errors.New("Error: node is not complete")
	}
}

func (displayConfig *displayConfig) getOneRowBasedOnHint(hint *HintStruct) (map[string]interface{}, error) {
	
	restrictions, err := hint.GetRestrictionsInTag(displayConfig)
	if err != nil {
		log.Fatal(err)
	}

	data := GetData1FromPhysicalSchemaByRowID(displayConfig, 
		hint.TableName + ".*", hint.TableName, 
		hint.RowIDs, restrictions,
	)

	if len(data) == 0 {
		return nil, common_funcs.DataNotExists
	} else {
		return data, nil
	}
}

func (displayConfig *displayConfig) getDataInNode(hint *HintStruct) ([]*HintStruct, error) {
	
	// Get and cache hint.Data if it is not there
	if len(hint.Data) == 0 {

		data, err := displayConfig.getOneRowBasedOnHint(hint)
		if err != nil {
			return nil, err
		}

		hint.Data = data

	}

	log.Println("My data is:", hint.Data)

	for _, tag := range displayConfig.dstAppConfig.dag.Tags {

		for _, member := range tag.Members {

			if hint.TableName == member {

				if len(tag.Members) == 1 {

					return []*HintStruct{hint}, nil

				} else {

					// Note: we assume that one dependency represents that one row
					// 		in one table depends on another row in another table
					return displayConfig.getRemainingDataInNode(
						tag.InnerDependencies, 
						tag.Members, hint,
					)
				}
			}
		}
	}

	return nil, errors.New("Error: the hint does not match any tags")
}

// A recursive function checks whether all the data one data recursively depends on exists
// We only checks whether the table depended on exists, which is sufficient for now
func (displayConfig *displayConfig) checkDependsOnExists(allData []*HintStruct, data *HintStruct) bool {
	
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
				
				if oneData.TableName == dependsOnTable {
					
					if !displayConfig.checkDependsOnExists(allData, oneData) {
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
		if displayConfig.checkDependsOnExists(allData, data) {
			trimmedData = append(trimmedData, data)
		}
	}

	return trimmedData
}

func (displayConfig *displayConfig) GetDataInNodeBasedOnDisplaySetting(hint *HintStruct) ([]*HintStruct, error) {
		
	// tagName, err := hint.GetTagName(appConfig)
	// if err != nil {
	// 	return nil, err
	// }
	// displaySetting, _ := appConfig.GetTagDisplaySetting(tagName)

	displaySetting, _ := hint.GetTagDisplaySetting(displayConfig)

	// Whether a node is complete or not, get all the data in a node.
	// If the node is complete, err is nil, otherwise, err is "node is not complete".
	if data, err := displayConfig.getDataInNode(hint); err != nil {

		// log.Println("++++++++++++")
		// log.Println(data)
		// log.Println("++++++++++++")
		
		// The setting "default_display_setting" means only display a node when the node is complete.
		// Therefore, return nil and error message when node is not complete.
		if displaySetting == "default_display_setting" {

			return nil, err
			
		// The setting "display_based_on_inner_dependencies" means display as much data in a node as possible
		// based on inner dependencies.
		// Note: if a piece of data in a node depends on some data not existing in the node,
		// it needs to be deleted from the data set and cannot be displayed.
		} else if displaySetting == "display_based_on_inner_dependencies" {
			return trimDataBasedOnInnerDependencies(displayConfig, data), err
		}

	// If a node is complete, return all the data in the node regardless of the setting.
	} else {
		return data, nil
	}

	panic("Should never happen")

}

package dependency_handler

import (
	"database/sql"
	"errors"
	"log"
	"stencil/config"
	"stencil/display"
	// "strconv"
	"strings"
	"fmt"
)

func getOneRowBasedOnDependency(appConfig *config.AppConfig, stencilDBConn *sql.DB, val string, dep string) (map[string]interface{}, error) {
	table := strings.Split(dep, ".")[0]
	key := strings.Split(dep, ".")[1]
	// log.Println(table)
	// log.Println(key)
	// log.Println(val)
	data := display.GetData1FromPhysicalSchema(stencilDBConn, appConfig.QR, appConfig.AppID, table + ".*", 
		table, table + "." + key, "=", val)

	if len(data) == 0 {
		return nil, errors.New("Error: Cannot Find One Remaining Data in the Node")
	} else {
		return data, nil
	}
}

func getRemainingDataInNode(appConfig *config.AppConfig, stencilDBConn *sql.DB, dependencies []map[string]string, members map[string]string, hint display.HintStruct) ([]display.HintStruct, error) {
	var result []display.HintStruct

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

	result = append(result, hint)

	queue := []DataInDependencyNode{DataInDependencyNode{
		Table: hint.TableName,
		Data:  hint.Data,
	}}
	for len(queue) != 0 && len(procDependencies) != 0 {
		// log.Println(queue)
		// log.Println(procDependencies)

		dataInDependencyNode := queue[0]
		queue = queue[1:]

		// table := dataInDependencyNode.Table
		for tableCol, val := range dataInDependencyNode.Data {
			if deps, ok := procDependencies[tableCol]; ok {
				// We assume that this is an integer value otherwise we have to define it in dependency config
				for _, dep := range deps {
					// log.Println(dep)
					// log.Println(tableCol)
					// log.Println(dataInDependencyNode.Data)
					if val == nil {
						log.Println("Fail to get one data because the value of the relevant column is nil")
						continue
					}
					data1, err1 := getOneRowBasedOnDependency(appConfig, stencilDBConn, fmt.Sprint(val), dep)
					if err1 != nil {
						// log.Println(err1)
						// fmt.Println(result)
						continue
					}

					table1 := strings.Split(dep, ".")[0]
					key1 := strings.Split(dep, ".")[1]
					queue = append(queue, DataInDependencyNode{
						Table: table1,
						Data:  data1,
					})

					rowID := display.GetRowIDsFromData(data1)
					tableID := display.GetTableIDByTableName(stencilDBConn, table1, appConfig.AppID)
					result = append(result, display.HintStruct{
						TableName: table1,
						TableID: tableID,
						RowIDs: rowID,
						Data: data1,
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

func getOneRowBasedOnHint(appConfig *config.AppConfig, stencilDBConn *sql.DB, hint display.HintStruct) (map[string]interface{}, error) {
	restrictions, err := hint.GetRestrictionsInTag(appConfig)
	if err != nil {
		log.Fatal(err)
	}

	data := display.GetData1FromPhysicalSchemaByRowID(stencilDBConn, appConfig.QR, appConfig.AppID, hint.TableName + ".*", hint.TableName, hint.RowIDs, restrictions)

	if len(data) == 0 {
		return nil, errors.New("Error: the Data in a Data Hint Does Not Exist")
	} else {
		return data, nil
	}
}

func getDataInNode(appConfig *config.AppConfig, hint display.HintStruct, stencilDBConn *sql.DB) ([]display.HintStruct, error) {
	// Get and cache hint.Data if it is not there
	if len(hint.Data) == 0 {
		data, err := getOneRowBasedOnHint(appConfig, stencilDBConn, hint)
		if err != nil {
			return nil, err
		}
		hint.Data = data
	}

	for _, tag := range appConfig.Tags {
		for _, member := range tag.Members {
			if hint.TableName == member {
				if len(tag.Members) == 1 {
					return []display.HintStruct{hint}, nil
				} else {
					// Note: we assume that one dependency represents that one row
					// 		in one table depends on another row in another table
					return getRemainingDataInNode(appConfig, stencilDBConn, tag.InnerDependencies, tag.Members, hint)
				}
			}
		}
	}
	return nil, errors.New("Error: the hint does not match any tags")
}

// A recursive function checks whether all the data one data recursively depends on exists
// We only checks whether the table depended on exists, which is sufficient for now
func checkDependsOnExists(appConfig *config.AppConfig, allData []display.HintStruct, tagName string, data display.HintStruct) bool {
	memberID, _ := data.GetMemberID(appConfig, tagName)
	// fmt.Println(memberID)
	dependsOnTables := appConfig.GetDependsOnTables(tagName, memberID)
	// fmt.Println(dependsOnTables)
	if len(dependsOnTables) == 0 {
		return true
	} else {
		for _, dependsOnTable := range dependsOnTables {
			exists := false
			for _, oneData := range allData {
				if oneData.TableName == dependsOnTable {
					if !checkDependsOnExists(appConfig, allData, tagName, oneData) {
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

func trimDataBasedOnInnerDependencies(appConfig *config.AppConfig, allData []display.HintStruct, tagName string) []display.HintStruct {
	var trimmedData []display.HintStruct

	for _, data := range allData {
		if checkDependsOnExists(appConfig, allData, tagName, data) {
			trimmedData = append(trimmedData, data)
		}
	}

	return trimmedData
}

func GetDataInNodeBasedOnDisplaySetting(appConfig *config.AppConfig, hint display.HintStruct, stencilDBConn *sql.DB) ([]display.HintStruct, error) {
	var data []display.HintStruct
	
	tagName, err := hint.GetTagName(appConfig)
	if err != nil {
		return nil, err
	}

	log.Println("Check Data ", hint)

	displaySetting, _ := appConfig.GetTagDisplaySetting(tagName)
	// Whether a node is complete or not, get all the data in a node.
	// If the node is complete, err is nil, otherwise, err is "node is not complete".
	if data, err = getDataInNode(appConfig, hint, stencilDBConn); err != nil {
		// log.Println(("**************"))
		// log.Println(data)
		// log.Println(("**************"))
		// The setting "default_display_setting" means only display a node when the node is complete.
		// Therefore, return nil and error message when node is not complete.
		if displaySetting == "default_display_setting" {
			return nil, err
			// The setting "display_based_on_inner_dependencies" means display as much data in a node as possible
			// based on inner dependencies.
			// Note: if a piece of data in a node depends on some data not existing in the node,
			// it needs to be deleted from the data set and cannot be displayed.
		} else if displaySetting == "display_based_on_inner_dependencies" {
			// fmt.Println(data)
			return trimDataBasedOnInnerDependencies(appConfig, data, tagName), err
		}
		// If a node is complete, return all the data in the node regardless of the setting.
	} else {
		return data, nil
	}

	panic("Should never happen")
}

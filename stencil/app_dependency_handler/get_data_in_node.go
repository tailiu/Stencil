package app_dependency_handler

import (
	"errors"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/app_display"
	"strconv"
	"strings"
)

func getOneRowBasedOnDependency(appConfig *config.AppConfig, val int, dep string) (map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %d", strings.Split(dep, ".")[0], strings.Split(dep, ".")[1], val)
	// fmt.Println(query)

	data, err := db.DataCall1(appConfig.DBConn, query)
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

func getRemainingDataInNode(appConfig *config.AppConfig, dependencies []map[string]string, members map[string]string, hint *app_display.HintStruct) ([]*app_display.HintStruct, error) {
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
				// We assume that this is an integer value otherwise we have to define it in dependency config
				intVal, err := strconv.Atoi(fmt.Sprint(val))
				if err != nil {
					log.Fatal("Error in Getting Data in Node: Converting '%s' to Integer", val)
				}
				for _, dep := range deps {
					data, err1 := getOneRowBasedOnDependency(appConfig, intVal, dep)
					// fmt.Println(data)

					if err1 != nil {
						// fmt.Println(err1)
						// fmt.Println(result)
						continue
					}
					// fmt.Println(dep)

					table1 := strings.Split(dep, ".")[0]
					key1 := strings.Split(dep, ".")[1]
					// fmt.Println(queue)
					queue = append(queue, DataInDependencyNode{
						Table: table1,
						Data:  data,
					})

					intPK, err2 := strconv.Atoi(fmt.Sprint(data["id"]))
					if err2 != nil {
						log.Fatal(err2)
					}
					keyVal := map[string]int{
						"id": intPK,
					}
					result = append(result, &app_display.HintStruct{
						Table: table1,
						TableID: appConfig.TableNameIDPairs[table1],
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

func getOneRowBasedOnHint(appConfig *config.AppConfig, hint *app_display.HintStruct) (map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = %d", hint.Table, hint.KeyVal["id"])
	log.Println(query)
	data, err := db.DataCall1(appConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if len(data) == 0 {
		return nil, app_display.DataNotExists
	} else {
		return data, nil
	}
}

func getDataInNode(appConfig *config.AppConfig, hint *app_display.HintStruct) ([]*app_display.HintStruct, error) {

	if hint.Data == nil {
		data, err := getOneRowBasedOnHint(appConfig, hint)
		if err != nil {
			return nil, err
		} else {
			hint.Data = data
		}
	}
	
	for _, tag := range appConfig.Tags {
		for _, member := range tag.Members {
			if hint.Table == member {
				if len(tag.Members) == 1 {
					return []*app_display.HintStruct{hint}, nil
				} else {
					// Note: we assume that one dependency represents that one row
					// 		in one table depends on another row in another table
					return getRemainingDataInNode(appConfig, tag.InnerDependencies, tag.Members, hint)
				}
			}
		}
	}
	return nil, errors.New("Error: the hint does not match any tags")
}

// A recursive function checks whether all the data one data recursively depends on exists
// We only checks whether the table depended on exists, which is sufficient for now
func checkDependsOnExists(appConfig *config.AppConfig, allData []*app_display.HintStruct, tagName string, data *app_display.HintStruct) bool {
	
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
				if oneData.Table == dependsOnTable {
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

func trimDataBasedOnInnerDependencies(appConfig *config.AppConfig, allData []*app_display.HintStruct, tagName string) []*app_display.HintStruct {
	var trimmedData []*app_display.HintStruct

	for _, data := range allData {
		if checkDependsOnExists(appConfig, allData, tagName, data) {
			trimmedData = append(trimmedData, data)
		}
	}

	return trimmedData
}

func GetDataInNodeBasedOnDisplaySetting(appConfig *config.AppConfig, hint *app_display.HintStruct) ([]*app_display.HintStruct, error) {
	var data []*app_display.HintStruct

	tagName, err := hint.GetTagName(appConfig)
	if err != nil {
		return nil, err
	}

	displaySetting, _ := appConfig.GetTagDisplaySetting(tagName)

	// Whether a node is complete or not, get all the data in a node.
	// If the node is complete, err is nil, otherwise, err is "node is not complete".
	if data, err = getDataInNode(appConfig, hint); err != nil {

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

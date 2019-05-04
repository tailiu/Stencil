package dependency_handler 

import (
	"transaction/config"
	"transaction/display"
	"transaction/db"
	"fmt"
	"strings"
	"log"
	"database/sql"
	"strconv"
	"errors"
)

func getOneRowBasedOnHint(dbConn *sql.DB, app, depDataTable, depDataKey string, depDataValue int) (map[string]string, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %d LIMIT 1;", depDataTable, depDataKey, depDataValue)

	data := db.GetAllColsOfRows(dbConn, query)
	if len(data) == 0 {
		return nil, errors.New("Error: the Data in a Data Hint Does Not Exist")
	} else {
		return data[0], nil
	}
}

func getOneRowBasedOnDependency(dbConn *sql.DB, app string, val int, dep string) (map[string]string, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %d LIMIT 1;", strings.Split(dep, ".")[0], strings.Split(dep, ".")[1], val)
	// fmt.Println(query)
	data := db.GetAllColsOfRows(dbConn, query)
	// fmt.Println(data)
	if len(data) == 0 {
		return nil, errors.New("Error: the Remaining Data in a Node Does Not Exist")
	} else {
		return data[0], nil
	}
}

func getRemainingData(dbConn *sql.DB, dependencies []map[string]string, members map[string]string, hint display.HintStruct, app string) ([]display.HintStruct, error) {
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
	fmt.Println(procDependencies)

	var data map[string]string
	var err error
	for k, v := range hint.KeyVal {
		data, err = getOneRowBasedOnHint(dbConn, app, hint.Table, k, v)
		if err != nil {
			// log.Println(err)
			return nil, err
		}
	}

	result = append(result, hint)

	queue := []DataInDependencyNode{DataInDependencyNode{
		Table:	hint.Table,
		Data:	data,
	}}
	for len(queue) != 0 && len(procDependencies) != 0 {
		// fmt.Println(queue)
		// fmt.Println(procDependencies)

		dataInDependencyNode := queue[0]
		queue = queue[1:]
		
		table := dataInDependencyNode.Table
		for col, val := range dataInDependencyNode.Data {
			if deps, ok := procDependencies[table + "." + col]; ok {
				// We assume that this is an integer value otherwise we have to define it in dependency config
				intVal, err := strconv.Atoi(val)
				if err != nil {
					log.Fatal("Error in Getting Data in Node: Converting '%s' to Integer", val)
				}
				for _, dep := range deps {
					data, err = getOneRowBasedOnDependency(dbConn, app, intVal, dep)
					// fmt.Println(data)
					
					if err != nil {
						// fmt.Println(err)
						fmt.Println(result)
						return nil, err
					}
					// fmt.Println(dep)

					table1 := strings.Split(dep, ".")[0]
					key1 := strings.Split(dep, ".")[1]
					// fmt.Println(queue)
					queue = append(queue, DataInDependencyNode{
						Table:	table1,
						Data:	data,
					})

					pk, err1 := db.GetPrimaryKeyOfTable(dbConn, table1)
					if err1 != nil {
						log.Fatal(err1)
					}
					intPK, err2 := strconv.Atoi(data[pk])
					if err2 != nil {
						log.Fatal(err2)
					}
					keyVal := map[string]int {
						pk:		intPK,
					}
					result = append(result, display.HintStruct{
						Table:		table1,
						KeyVal:		keyVal,
					})

					deps1 := procDependencies[table1 + "." + key1]
					for i, val2 := range deps1 {
						if val2 == table + "." + col {
							deps1 = append(deps1[:i], deps1[i+1:]...)
							break
						}
					}
					if len(deps1) == 0 {
						delete(procDependencies, table1 + "." + key1)
					} else {
						procDependencies[table1 + "." + key1] = deps1
					}
				}
				delete(procDependencies, table + "." + col)
			}
		}
	}

	// fmt.Println(procDependencies)
	// fmt.Println(result)
	if len(procDependencies) == 0 {
		return result, nil
	} else {
		return result, errors.New("Error: node is not complete")
	}
}

func GetDataInNode(appConfig *config.AppConfig, hint display.HintStruct) ([]display.HintStruct, error) {
	for _, tag := range appConfig.Tags {
		for _, member := range tag.Members{
			if hint.Table == member {
				if len(tag.Members) == 1 {
					return []display.HintStruct{hint}, nil
				} else {
					// Note: we assume that one dependency represents that one row 
					// 		in one table depends on another row in another table
					return getRemainingData(appConfig.DBConn, tag.InnerDependencies, tag.Members, hint, appConfig.AppName)
				}
			}
		}
	}
	return []display.HintStruct{}, errors.New("Error: the hint does not match any tags")
}


func GetDataInNodeBasedOnDisplaySetting(appConfig *config.AppConfig, hint display.HintStruct) ([]display.HintStruct, error) {
	var data []display.HintStruct

	tagName, err := hint.GetTagName(appConfig.Tags)
	if err != nil {
		return nil, err
	}
	// fmt.Println(tagName)

	displaySetting, _ := appConfig.GetTagDisplaySetting(tagName)
	if data, err = GetDataInNode(appConfig, hint); err != nil {
		if displaySetting == "default_display_setting" {
			return nil, err
		} else if displaySetting == "display_based_on_inner_dependencies" {
			return data, err
		}
	} else {
		return data, nil
	}

	panic("Should never happen")
}
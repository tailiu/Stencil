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
	"mastodon/auxiliary"
)

const getOneDataFromParentNodeAttemptTimes = 10

type DataInDependencyNode struct {
	Table 	string
	Data	map[string]string
}

func checkRemainingDataExists(dependencies []map[string]string, members map[string]string, hint display.HintStruct, app string, dbConn *sql.DB) ([]display.HintStruct, bool) {
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

	data, err := db.GetOneRowBasedOnHint(dbConn, app, hint.Value, hint.ValueType, hint.Key, hint.Table)
	// fmt.Println(data)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	result = append(result, display.HintStruct{
		Table:		hint.Table,
		Key:		hint.Key,
		Value:		data[hint.Key],
		ValueType: 	"int",
	})

	queue := []DataInDependencyNode{DataInDependencyNode{
		Table:	hint.Table,
		Data:	data,
	}}
	for len(queue) != 0 && len(procDependencies) != 0 {
		dataInDependencyNode, queue := queue[0], queue[1:]
		
		table := dataInDependencyNode.Table
		for col, val := range dataInDependencyNode.Data {
			if deps, ok := procDependencies[table + "." + col]; ok {
				// We assume that this is an integer value otherwise we have to define it in dependency config
				intVal, err := strconv.Atoi(val)
				if err != nil {
					log.Fatal("Dependency Handler: Converting '%s' to Integer Errors", val)
				}
				for _, dep := range deps {
					data, err = db.GetOneRowBasedOnDependency(dbConn, app, intVal, dep)
					if err != nil {
						fmt.Println(err)
						return nil, false
					}
					// fmt.Println(dep)
					// fmt.Println(data["account_id"])

					table1 := strings.Split(dep, ".")[0]
					key1 := strings.Split(dep, ".")[1]
					queue = append(queue, DataInDependencyNode{
						Table:	table1,
						Data:	data,
					})
					result = append(result, display.HintStruct{
						Table:		table1,
						Key:		key1,
						Value:		data[key1],
						ValueType:	"int",
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
	if len(procDependencies) == 0 {
		return result, true
	} else {
		return nil, false
	}
}

func getDependsOn(dependencies []config.Dependency, tag string) ([]config.DependsOn, error) {
	for _, dependency := range dependencies {
		// fmt.Println(dependency)
		if dependency.Tag == tag {
			return dependency.DependsOn, nil
		}
	}
	return nil, errors.New("Cannot Find Any Parent Tags")
}

func GetTagName(innerDependencies []config.Tag, hint display.HintStruct) (string, error) {
	for _, innerDependency := range innerDependencies {
		for _, member := range innerDependency.Members{
			if hint.Table == member {
				return innerDependency.Name, nil
			}
		}
	}
	return "", errors.New("No Corresponding Tag Found!")
}

func replaceKey(innerDependencies []config.Tag, tag string, key string) string {
	for _, dependsOnTag := range innerDependencies {
		if dependsOnTag.Name == tag {
			// fmt.Println(dependsOnTag)
			for k, v := range dependsOnTag.Keys {
				if k == key {
					member := strings.Split(v, ".")[0]
					attr := strings.Split(v, ".")[1]
					for k1, table := range dependsOnTag.Members {
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

func GetOneDataFromParentNode(appConfig config.AppConfig, hint display.HintStruct, app string, dbConn *sql.DB) (display.HintStruct, error){
	hintData := display.HintStruct{}
	data1 := DataInDependencyNode{}

	tag, err := GetTagName(appConfig.Tags, hint)
	if err != nil {
		log.Fatal(err)
	} 
	
	dependsOn, err1 := getDependsOn(appConfig.Dependencies, tag)
	if err1 != nil {
		log.Println(err1)
		return hint, err1
	}

	// fmt.Println(dependsOn)
	for i := 0; i < getOneDataFromParentNodeAttemptTimes; i ++ {
		oneDependensOn := dependsOn[auxiliary.RandomNonnegativeIntWithUpperBound(len(dependsOn))]
		// fmt.Println(oneDependensOn)

		var conditions []string
		if len(oneDependensOn.Conditions) == 1 {
			condition := oneDependensOn.Conditions[0]
			from := replaceKey(appConfig.Tags, tag, condition.TagAttr)
			to := replaceKey(appConfig.Tags, oneDependensOn.Tag, condition.DependsOnAttr)
			conditions = append(conditions, from + ":" + to)
		}

		fmt.Println(conditions)
		fmt.Println(hint)

		data1.Data, data1.Table, err1 = db.GetOneRowInParentNodeRandomly(dbConn, hint.Value, hint.ValueType, hint.Key, hint.Table, conditions)
		if err1 != nil {
			fmt.Println(err1)
		} else {
			fmt.Println(data1)
			break
		}
	}

	hintData, err2 := display.TransformRowToHint(dbConn, data1.Data, data1.Table)
	if err2 != nil {
		log.Fatal(err2)
	} 
	fmt.Println(hintData)
	return hintData, nil
}

func CheckNodeComplete(innerDependencies []config.Tag, hint display.HintStruct, app string, dbConn *sql.DB) bool {
	for _, innerDependency := range innerDependencies {
		for _, member := range innerDependency.Members{
			if hint.Table == member {
				if len(innerDependency.Members) == 1 {
					return true
				} else {
					// Note: we assume that one dependency represents that one row 
					// 		in one table depends on another row in another table
					if result, ok:= checkRemainingDataExists(innerDependency.InnerDependencies, innerDependency.Members, hint, app, dbConn); ok {
						fmt.Println(result)
						return true
					} else {
						return false
					}
				}
			}
		}
	}
	return false
}

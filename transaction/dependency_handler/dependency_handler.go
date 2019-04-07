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

/**************************** Check data in Node complete and return data if it is complete *****************/

func getOneRowBasedOnHint(dbConn *sql.DB, app, depDataTable, depDataKey string, depDataValue int) (map[string]string, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %d LIMIT 1;", depDataTable, depDataKey, depDataValue)

	data := db.GetAllColsOfRows(dbConn, query)
	if len(data) == 0 {
		return nil, errors.New("Check Remaining Data Exists Error: Original Data Not Exists")
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
		return nil, errors.New("Check Remaining Data Exists Error: Data Not Exists")
	} else {
		return data[0], nil
	}
}

func checkRemainingDataExists(dbConn *sql.DB, dependencies []map[string]string, members map[string]string, hint display.HintStruct, app string) ([]display.HintStruct, bool) {
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
	// fmt.Println(procDependencies)

	var data map[string]string
	var err error
	for k, v := range hint.KeyVal {
		data, err = getOneRowBasedOnHint(dbConn, app, hint.Table, k, v)
		if err != nil {
			log.Println(err)
			return nil, false
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
					log.Fatal("Dependency Handler: Converting '%s' to Integer Errors", val)
				}
				for _, dep := range deps {
					data, err = getOneRowBasedOnDependency(dbConn, app, intVal, dep)
					// fmt.Println(data)
					
					if err != nil {
						fmt.Println(err)
						return nil, false
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
		return result, true
	} else {
		return nil, false
	}
}

func CheckNodeComplete(dbConn *sql.DB, innerDependencies []config.Tag, hint display.HintStruct, app string) (bool, []display.HintStruct) {
	for _, innerDependency := range innerDependencies {
		for _, member := range innerDependency.Members{
			if hint.Table == member {
				if len(innerDependency.Members) == 1 {
					var completeData []display.HintStruct
					completeData = append(completeData, hint)
					return true, completeData
				} else {
					// Note: we assume that one dependency represents that one row 
					// 		in one table depends on another row in another table
					if completeData, ok := checkRemainingDataExists(dbConn, innerDependency.InnerDependencies, innerDependency.Members, hint, app); ok {
						fmt.Println(completeData)
						return true, completeData
					} else {
						return false, nil
					}
				}
			}
		}
	}
	return false, nil
}

/**************************** end **************************/

func getOneRowInParentNodeRandomly(dbConn *sql.DB, hint display.HintStruct, conditions []string) (map[string]string, string, error) {
	query := fmt.Sprintf("SELECT %s.* FROM ", "t"+strconv.Itoa(len(conditions)))
	from := ""
	table := ""
	for i, condition := range conditions {
		table1 := strings.Split(condition, ":")[0]
		table2 := strings.Split(condition, ":")[1]
		t1 := strings.Split(table1, ".")[0]
		a1 := strings.Split(table1, ".")[1]
		t2 := strings.Split(table2, ".")[0]
		a2 := strings.Split(table2, ".")[1]
		seq1 := "t" + strconv.Itoa(i)
		seq2 := "t" + strconv.Itoa(i+1)
		if i == 0 {
			from += fmt.Sprintf("%s %s JOIN %s %s ON %s.%s = %s.%s ",
				t1, seq1, t2, seq2, seq1, a1, seq2, a2)
		} else {
			from += fmt.Sprintf("JOIN %s %s on %s.%s = %s.%s ",
				t2, seq2, seq1, a1, seq2, a2)
		}
		if i == len(conditions)-1 {
			var depDataKey string
			var depDataValue int
			for k, v := range hint.KeyVal {
				depDataKey = k
				depDataValue = v
			}
			where := fmt.Sprintf("WHERE %s.%s = %d ORDER BY RANDOM() LIMIT 1;", "t0", depDataKey, depDataValue)
			table = t2
			query += from + where
		}
	}
	// fmt.Println(query)

	data := db.GetAllColsOfRows(dbConn, query)
	if len(data) == 0 {
		return nil, "", errors.New("Error In Get Data: Fail To Get One Data This Data Depends On")
	} else {
		return data[0], table, nil
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

func GetTagName(innerDependencies []config.Tag, hint display.HintStruct) (string, error) {
	for _, innerDependency := range innerDependencies {
		for _, member := range innerDependency.Members {
			if hint.Table == member {
				return innerDependency.Name, nil
			}
		}
	}
	return "", errors.New("No Corresponding Tag Found!")
}

func GetParentTags(appConfig config.AppConfig, data display.HintStruct) ([]string, error) {
	tag, err := GetTagName(appConfig.Tags, data)
	if err != nil {
		return nil, err
	}
	if tag == "root" {
		return nil, nil
	}

	var parentTags []string
	for _, dependency := range appConfig.Dependencies {
		if dependency.Tag == tag {
			for _, dependsOn := range dependency.DependsOn {
				parentTags = append(parentTags, dependsOn.Tag)
			}
		}
	}

	if len(parentTags) == 0 {
		return nil, errors.New("Check Parent Tag Name error: Does Not Find Any Parent Node!")
	} else {
		return parentTags, nil
	}
}

func GetOneDataFromParentNodeRandomly(dbConn *sql.DB, appConfig config.AppConfig, hint display.HintStruct, app string) (display.HintStruct, error){
	hintData := display.HintStruct{}
	data1 := DataInDependencyNode{}

	tag, err := GetTagName(appConfig.Tags, hint)
	if err != nil {
		log.Fatal(err)
	} 
	
	dependsOn, err1 := getDependsOn(appConfig.Dependencies, tag)
	if err1 != nil {
		log.Println(err1)
		return hintData, err1
	}

	// fmt.Println(dependsOn)
	for i := 0; i < getOneDataFromParentNodeAttemptTimes; i ++ {
		oneDependensOn := dependsOn[auxiliary.RandomNonnegativeIntWithUpperBound(len(dependsOn))]
		// fmt.Println(oneDependensOn)

		var conditions []string
		var from, to string
		if len(oneDependensOn.Conditions) == 1 {
			condition := oneDependensOn.Conditions[0]
			from = replaceKey(appConfig.Tags, tag, condition.TagAttr)
			to = replaceKey(appConfig.Tags, oneDependensOn.Tag, condition.DependsOnAttr)
			conditions = append(conditions, from + ":" + to)
		} else {
			for i, condition := range(oneDependensOn.Conditions) {
				if i == 0 {
					from = replaceKey(appConfig.Tags, tag, condition.TagAttr)
					to = replaceKey(appConfig.Tags, strings.Split(condition.DependsOnAttr, ".")[0], strings.Split(condition.DependsOnAttr, ".")[1])
				} else if i == len(oneDependensOn.Conditions) - 1 {
					from = replaceKey(appConfig.Tags, strings.Split(condition.TagAttr, ".")[0], strings.Split(condition.TagAttr, ".")[1])
					to = replaceKey(appConfig.Tags, oneDependensOn.Tag, condition.DependsOnAttr)
				} 
				conditions = append(conditions, from + ":" + to)
			}
		}

		// fmt.Println(conditions)
		// fmt.Println(hint)

		data1.Data, data1.Table, err1 = getOneRowInParentNodeRandomly(dbConn, hint, conditions)
		if err1 != nil {
			fmt.Println(err1)
		} else {
			// fmt.Println(data1)
			break
		}
	}

	hintData, err2 := display.TransformRowToHint(dbConn, data1.Data, data1.Table)
	if err2 != nil {
		log.Fatal(err2)
	} 
	// fmt.Println(hintData)
	return hintData, nil
}

func CheckDisplayCondition() bool {
	return false
}
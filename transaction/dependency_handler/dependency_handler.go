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
)

type DataInDependencyNode struct {
	Table 	string
	Data	map[string]string
}

func GetParent() {

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

	data, err := db.GetOneRowBasedOnHint(dbConn, app, hint)
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
					// fmt.Println()

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

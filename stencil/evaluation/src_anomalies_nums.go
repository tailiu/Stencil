package evaluation

// import (
// 	"log"
// 	"strconv"
// 	"strings"
// 	"fmt"
// 	"stencil/db"
// 	// "reflect"
// 	"time"
// )

// func getTableKeyDeletedAt(evalConfig *EvalConfig, migrationID string) []map[string]interface{} {
// 	conditions := "dst_table != 'n/a'"
// 	query := fmt.Sprintf("select src_table, src_id, deleted_at from evaluation where migration_id = '%s' and %s;", 
// 		migrationID, conditions)
	
// 	data, err := db.DataCall(evalConfig.StencilDBConn, query)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
	
// 	return data
// }

// func getDeletedAtInEvaluation(evalConfig *EvalConfig, migrationID, dependsOnTable string, pKey int64) map[string]interface{} {
// 	query := fmt.Sprintf("select deleted_at from evaluation where migration_id = '%s' and src_table = '%s' and src_id = '%s'",
// 		migrationID, dependsOnTable, strconv.FormatInt(pKey, 10))
	
// 	data, err := db.DataCall1(evalConfig.StencilDBConn, query)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return data
// }

// func getDependsOnTableKeys(evalConfig *EvalConfig, app, table string) []string {
// 	return evalConfig.Dependencies[app][table]
// }

// func violateDependencies(evalConfig *EvalConfig, table string, pKey int, deleted_at time.Time, migrationID string) int {
// 	dependsOnTableKeys := getDependsOnTableKeys(evalConfig, "diaspora", table)
// 	if len(dependsOnTableKeys) == 0 {
// 		return 0
// 	}

// 	// log.Println(table)
// 	// log.Println(dependsOnTableKeys)
// 	log.Println(table)
// 	log.Println(pKey)

// 	violationNum := 0

// 	row := getLogicalRow(evalConfig.DiasporaDBConn, table, pKey)
// 	for _, dependsOnTableKey := range dependsOnTableKeys {
// 		fromAttr := strings.Split(dependsOnTableKey, ":")[0]
// 		if row[fromAttr] == nil {
// 			continue
// 		}
// 		dependsOnTable := strings.Split(strings.Split(dependsOnTableKey, ":")[1], ".")[0]
// 		dependsOnKey := strings.Split(strings.Split(dependsOnTableKey, ":")[1], ".")[1]
// 		query := fmt.Sprintf("select * from %s where %s = %d", dependsOnTable, dependsOnKey, row[fromAttr].(int64))
// 		// log.Println(query)
// 		row1, err := db.DataCall1(evalConfig.DiasporaDBConn, query)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// log.Println(row1)
// 		row2 := getDeletedAtInEvaluation(evalConfig, migrationID, dependsOnTable, row1["id"].(int64))
// 		if row2["deleted_at"] == nil {
// 			// This can happen when migration is not complete
// 			log.Println("dependsOn_deleted_at is nil!!")
// 			continue
// 		}
// 		dependsOn_deleted_at := row2["deleted_at"].(time.Time)
// 		log.Println(dependsOn_deleted_at)
// 		log.Println(deleted_at)
// 		if dependsOn_deleted_at.Before(deleted_at) {
// 			log.Println("Got one")
// 			violationNum += 1
// 		}
// 	}

// 	return violationNum
// }

// func GetAnomaliesNumsInSrc(evalConfig *EvalConfig, migrationID string, side string) {
// 	// var sourceAnomalies map[string]int
// 	data := getTableKeyDeletedAt(evalConfig, migrationID)
// 	// fmt.Println(data)
// 	checkedRow := make(map[string]bool) 
// 	for _, data1 := range data {
// 		// log.Println(reflect.TypeOf(data1["deleted_at"]))
// 		table, pKey := transformTableKeyToNormalType(data1)
// 		key := table + ":" + strconv.Itoa(pKey)
// 		if _, ok := checkedRow[key]; ok {
// 			continue
// 		} else {
// 			checkedRow[key] = true
// 			violateDependencies(evalConfig, table, pKey, data1["deleted_at"].(time.Time), migrationID)
// 		}
// 	}
// }
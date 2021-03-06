package evaluation

import (
	"log"
	"strconv"
	"strings"
	"fmt"
	"stencil/db"
	// "reflect"
	"time"
)

func getTableKeyDeletedAt(evalConfig *EvalConfig, 
	migrationID string) []map[string]interface{} {
	
	conditions := "dst_table != 'n/a'"
	
	query := fmt.Sprintf(
		`select src_table, src_id, deleted_at 
		from evaluation where migration_id = '%s' and %s;`, 
		migrationID, conditions)
	
	data, err := db.DataCall(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	
	return data
}

func getDeletedAtInEvaluation(evalConfig *EvalConfig, 
	migrationID, dependsOnTable string, 
	pKey int64) map[string]interface{} {
	
	query := fmt.Sprintf(
		`select deleted_at from evaluation 
		where migration_id = '%s' and 
		src_table = '%s' and src_id = '%s'`,
		migrationID, dependsOnTable, 
		strconv.FormatInt(pKey, 10))
	
	data, err := db.DataCall1(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func srcViolateDependencies(evalConfig *EvalConfig, 
	table string, pKey int, deleted_at time.Time, 
	migrationID string) (map[string]int, []time.Duration) {

	violateStats := make(map[string]int)
	var interruptionDuration []time.Duration

	dependsOnTableKeys := getDependsOnTableKeys(evalConfig, "diaspora", table)
	if len(dependsOnTableKeys) == 0 {
		return violateStats, interruptionDuration
	}

	// log.Println(table)
	// log.Println(dependsOnTableKeys)
	log.Println(table)
	// log.Println(pKey)

	row := getLogicalRow(evalConfig.DiasporaDBConn, table, pKey)
	for _, dependsOnTableKey := range dependsOnTableKeys {
		
		statsKey := table + "." + dependsOnTableKey
		log.Println(statsKey)

		fromAttr := strings.Split(dependsOnTableKey, ":")[0]
		if row[fromAttr] == nil {
			continue
		}

		dependsOnTable := strings.Split(strings.Split(dependsOnTableKey, ":")[1], ".")[0]
		dependsOnKey := strings.Split(strings.Split(dependsOnTableKey, ":")[1], ".")[1]
		
		var query string
		// Diaspora reshare is very special
		if fromAttr == "root_guid" {
			
			query = fmt.Sprintf("select * from %s where %s = '%s'", 
			dependsOnTable, dependsOnKey, row[fromAttr].(string))
		
		} else {

			query = fmt.Sprintf("select * from %s where %s = %d", 
			dependsOnTable, dependsOnKey, row[fromAttr].(int64))
		
		}

		// log.Println(query)

		row1, err := db.DataCall1(evalConfig.DiasporaDBConn, query)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println(row1)

		// This gets the data from the current user
		row2 := getDeletedAtInEvaluation(evalConfig, 
			migrationID, dependsOnTable, row1["id"].(int64))

		if row2["deleted_at"] == nil {
			// This can happen when the data the checked data depends on 
			// does not belong to the current checked user
			log.Println("dependsOn_deleted_at is nil!!")
			continue
		}

		dependsOn_deleted_at := row2["deleted_at"].(time.Time)
		
		log.Println(dependsOn_deleted_at)
		log.Println(deleted_at)
		
		if dependsOn_deleted_at.Before(deleted_at) {

			interruptionDuration = append(interruptionDuration, 
				deleted_at.Sub(dependsOn_deleted_at))
			
			increaseMapValOneByKey(violateStats, statsKey)

			log.Println("Got one")

		}
	}

	return violateStats, interruptionDuration
}

func getTotalUnmigratedLikes(evalConfig *EvalConfig) int64 {
	
	query := fmt.Sprintf("select count(*) from likes where mark_as_delete = 'false'")
	
	data, err := db.DataCall1(evalConfig.DiasporaDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return data["count"].(int64)
}

func GetAnomaliesNumsInSrc(evalConfig *EvalConfig, 
	migrationID string) (map[string]int, []time.Duration, map[string]int64) {
	
	danglingDataStats := make(map[string]int64)
	violateStats := make(map[string]int)
	
	var interruptionDuration []time.Duration

	// ********** This is added for testing **********
	// unmigratedLikesNumBefore := getTotalUnmigratedLikes(evalConfig)
	// ********** This is added for testing **********

	data := getTableKeyDeletedAt(evalConfig, migrationID)
	// fmt.Println(data)
	
	checkedRow := make(map[string]bool) 
	
	for _, data1 := range data {
		
		// log.Println(reflect.TypeOf(data1["deleted_at"]))
		table, pKey := transformTableKeyToNormalType(data1)
		key := table + ":" + strconv.Itoa(pKey)
		
		if _, ok := checkedRow[key]; ok {
			
			continue
		} else {

			checkedRow[key] = true

			danglingDataStats1 := srcDanglingDataNonSystem(evalConfig, migrationID, table, pKey)
			violateStats1, interruptionDuration1 := srcViolateDependencies(evalConfig, 
				table, pKey, data1["deleted_at"].(time.Time), migrationID)

			IncreaseMapValByMapInt64(danglingDataStats, danglingDataStats1)
			IncreaseMapValByMap(violateStats, violateStats1)
			interruptionDuration = append(interruptionDuration, interruptionDuration1...)
			
			log.Println("+++++++++++++++++++++++++++++++++++++++++++++++")
			log.Println("Source Violation Statistics:", violateStats)
			log.Println("Source Interruption Duration:", interruptionDuration)
			log.Println("Source Dangling Statistics:", danglingDataStats)
			log.Println("+++++++++++++++++++++++++++++++++++++++++++++++")
		}
	}

	// ********** This is added for testing **********
	// var LikesBeforePost int64
	// unmigratedLikesNumAfter := getTotalUnmigratedLikes(evalConfig)
	// if val, ok := violateStats["likes.target_id:posts.id"]; ok {
	// 	LikesBeforePost = unmigratedLikesNumBefore - unmigratedLikesNumAfter - int64(val)
	// } else {
	// 	LikesBeforePost = unmigratedLikesNumBefore - unmigratedLikesNumAfter
	// }
	// violateStats["totalLikes"] = int(unmigratedLikesNumAfter)
	// violateStats["LikesBeforePosts"] = int(LikesBeforePost)
	// ********** This is added for testing **********
	
	return violateStats, interruptionDuration, danglingDataStats
}
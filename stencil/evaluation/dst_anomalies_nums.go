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

func getTableKeyAddedAt(evalConfig *EvalConfig, migrationID string) []map[string]interface{} {
	conditions := "dst_table != 'n/a'"
	query := fmt.Sprintf("select dst_table, dst_id, added_at from evaluation where migration_id = '%s' and %s;", 
		migrationID, conditions)
	
	data, err := db.DataCall(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	
	return data
}

func getAddedAtInEvaluation(evalConfig *EvalConfig, migrationID, dependsOnTable string, pKey int64) map[string]interface{} {
	query := fmt.Sprintf("select added_at from evaluation where migration_id = '%s' and dst_table = '%s' and dst_id = '%s'",
		migrationID, dependsOnTable, strconv.FormatInt(pKey, 10))
	
	data, err := db.DataCall1(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func dstViolateDependencies(evalConfig *EvalConfig, table string, pKey int, added_at time.Time, migrationID string) (map[string]int, map[string]int) {
	log.Println(table)
	log.Println(pKey)

	violateStats := make(map[string]int)
	depNotMigratedStats := make(map[string]int)

	dependsOnTableKeys := getDependsOnTableKeys(evalConfig, "mastodon", table)
	if len(dependsOnTableKeys) == 0 {
		return violateStats, depNotMigratedStats
	}

	row := getLogicalRow(evalConfig.MastodonDBConn, table, pKey)
	log.Println(row)

	checkFavourte := 0
	for _, dependsOnTableKey := range dependsOnTableKeys {
		// log.Println(dependsOnTableKey)
		statsKey := table + "." + dependsOnTableKey
		log.Println(statsKey)

		fromAttr := strings.Split(dependsOnTableKey, ":")[0]
		// This can happen when it does not depends on
		if row[fromAttr] == nil {
			log.Println("fromAttr is nil!!")
			continue
		}
		dependsOnTable := strings.Split(strings.Split(dependsOnTableKey, ":")[1], ".")[0]
		dependsOnKey := strings.Split(strings.Split(dependsOnTableKey, ":")[1], ".")[1]
		query := fmt.Sprintf("select * from %s where %s = %d", dependsOnTable, dependsOnKey, row[fromAttr].(int64))
		// log.Println(query)
		row1, err := db.DataCall1(evalConfig.MastodonDBConn, query)
		if err != nil {
			log.Fatal(err)
		}
		// This can happen when data depended on is not migrated
		if row1["id"] == nil {
			// This works when commneting or favouriting on other's data or on user's own data
			if table == "favourites" { 
				if checkFavourte != 1 {
					checkFavourte += 1
				} else {
					favourteStatsKey := "favourites.status_id:statuses.id:comments.id"
					increaseMapValOneByKey(depNotMigratedStats, favourteStatsKey)
				}
			} else {
				increaseMapValOneByKey(depNotMigratedStats, statsKey)
			}
			log.Println("dependsOn is nil!!")
			continue
		}
		// log.Println(row1)
		row2 := getAddedAtInEvaluation(evalConfig, migrationID, dependsOnTable, row1["id"].(int64))
		// This can happen when migration is not complete
		if row2["added_at"] == nil {
			continue
		}
		dependsOn_added_at := row2["added_at"].(time.Time)
		log.Println(dependsOn_added_at)
		log.Println(added_at)
		if added_at.Before(dependsOn_added_at) {
			increaseMapValOneByKey(violateStats, statsKey)
			log.Println("Got one")
		}
	}

	return violateStats, depNotMigratedStats
}

func GetAnomaliesNumsInDst(evalConfig *EvalConfig, migrationID string) (map[string]int, map[string]int){
	violateStats := make(map[string]int)
	depNotMigratedStats := make(map[string]int)

	// var sourceAnomalies map[string]int
	data := getTableKeyAddedAt(evalConfig, migrationID)
	// fmt.Println(data)
	checkedRow := make(map[string]bool) 
	for _, data1 := range data {
		// log.Println(reflect.TypeOf(data1["deleted_at"]))
		table, pKey := transformTableKeyToNormalTypeInDstApp(data1)
		key := table + ":" + strconv.Itoa(pKey)
		if _, ok := checkedRow[key]; ok {
			continue
		} else {
			checkedRow[key] = true
			violateStats1, depNotMigratedStats1 := dstViolateDependencies(evalConfig, table, pKey, data1["added_at"].(time.Time), migrationID)
			IncreaseMapValByMap(violateStats, violateStats1)
			IncreaseMapValByMap(depNotMigratedStats, depNotMigratedStats1)
			
			log.Println("+++++++++++++++++++++++++++++++++++++++++++++++")
			log.Println("Destination Violation Statistics:", violateStats)
			log.Println("Destination data depended on not migrated statistics:", depNotMigratedStats)
			log.Println("+++++++++++++++++++++++++++++++++++++++++++++++")
		}
	}

	return violateStats, depNotMigratedStats
}
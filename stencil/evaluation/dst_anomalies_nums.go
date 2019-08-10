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

func getDependsOnTableKeys(evalConfig *EvalConfig, app, table string) []string {
	return evalConfig.Dependencies[app][table]
}

func violateDependencies(evalConfig *EvalConfig, table string, pKey int, added_at time.Time, migrationID string) (map[string]int, map[string]int) {
	violateStats := make(map[string]int)
	depNotMigratedStats := make(map[string]int)

	dependsOnTableKeys := getDependsOnTableKeys(evalConfig, "mastodon", table)
	if len(dependsOnTableKeys) == 0 {
		return violateStats, depNotMigratedStats
	}

	log.Println("***********************")
	log.Println(table)
	log.Println(pKey)
	// log.Println(dependsOnTableKeys)

	row := getLogicalRow(evalConfig.OldMastodonDBConn, table, pKey)
	for _, dependsOnTableKey := range dependsOnTableKeys {
		log.Println(dependsOnTableKey)
		
		statsKey := table + "." + dependsOnTableKey

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
		row1, err := db.DataCall1(evalConfig.OldMastodonDBConn, query)
		if err != nil {
			log.Fatal(err)
		}
		// This can happen when data depended on is not migrated
		if row1["id"] == nil {
			log.Println("dependsOn is nil!!")
			if _, ok := depNotMigratedStats[statsKey]; ok {
				depNotMigratedStats[statsKey] += 1
			} else {
				depNotMigratedStats[statsKey] = 1
			}
			continue
		}
		// log.Println(row1)
		row2 := getAddedAtInEvaluation(evalConfig, migrationID, dependsOnTable, row1["id"].(int64))
		// if row2["added_at"] == nil {
		// 	log.Println("dependsOn_added_at is nil!!")
		// 	continue
		// }
		dependsOn_added_at := row2["added_at"].(time.Time)
		log.Println(dependsOn_added_at)
		log.Println(added_at)
		if added_at.Before(dependsOn_added_at) {
			log.Println("Got one")
			if _, ok := violateStats[statsKey]; ok {
				violateStats[statsKey] += 1
			} else {
				violateStats[statsKey] = 1
			}
		}
	}
	log.Println("***********************")

	return violateStats, depNotMigratedStats
}

func GetAnomaliesNumsInDst(evalConfig *EvalConfig, migrationID string, side string) (map[string]int, map[string]int){
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
			violateStats1, depNotMigratedStats1 := violateDependencies(evalConfig, table, pKey, data1["added_at"].(time.Time), migrationID)
			
			for k, v := range violateStats1 {
				if _, ok := violateStats[k]; ok {
					violateStats[k] += v
				} else {
					violateStats[k] = v
				}
			}

			for k, v := range depNotMigratedStats1 {
				if _, ok := depNotMigratedStats[k]; ok {
					depNotMigratedStats[k] += v
				} else {
					depNotMigratedStats[k] = v
				}
			}
		}
	}

	return violateStats, depNotMigratedStats
}
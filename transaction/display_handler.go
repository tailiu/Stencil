package main

import (
	"fmt"
	"log"
	// "transaction/atomicity"
	"transaction/db"
	"transaction/display"
	"transaction/dependency_handler"
	"transaction/config"
	"database/sql"
	"time"
	"encoding/json"
	"strconv"
	// "errors"
)

const StencilDBName = "stencil"
const checkInterval = 200 * time.Millisecond

var displayedData = make(map[string]int)

func procData(rawData []sql.NullString) []display.HintStruct {
	var processedData []display.HintStruct
	for _, oneData := range rawData {
		var oneSetOfHints []display.HintStruct
		json.Unmarshal([]byte(oneData.String), &oneSetOfHints)
		for _, hint := range oneSetOfHints {
			processedData = append(processedData, hint)
		}
	}
	return processedData
}

func getMigratedData(migrationID int, stencilDBConn *sql.DB) []display.HintStruct {
	var displayHints []sql.NullString
	var hintString sql.NullString

	op := fmt.Sprintf("SELECT display_hint FROM txn_log WHERE action_id = %d", migrationID)
	rows, err := stencilDBConn.Query(op)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&hintString); err != nil {
            log.Fatal(err)
		}
		// Only get log entries with display_hint not NULL, so BEGIN, COMMIT, etc. will be ignored 
		if hintString.Valid {
			displayHints = append(displayHints, hintString)
		}
	}

	// display hints contain info to find data in destination application
	// E.g.,[{"Table":"conversations","Key":"account_id","Value":"1517102025","ValueType":"int"},
	// 		{"Table":"account_stats","Key":"account_id","Value":"1918176832","ValueType":"int"}]
	return procData(displayHints)
}

func checkMigrationComplete(migrationID int, dbConn *sql.DB) bool {
	var complete bool
	op := fmt.Sprintf("SELECT id FROM txn_log WHERE action_id = %d and action_type='COMMIT' LIMIT 1", migrationID)
	rows, err := dbConn.Query(op)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		complete = true
	}

	return complete
}

func DisplayThread(dstApp string, migrationID int) {
	var secondRound bool

	stencilDBConn := db.GetDBConn(StencilDBName)
	destAppDBConn := db.GetDBConn(dstApp)
	appConfig, err := config.CreateAppConfig(dstApp)
	if err != nil {
		log.Fatal(err)
	}

	// For now just assume this is an infinite loop
	for migratedData := getMigratedData(migrationID, stencilDBConn); checkMigrationComplete(migrationID, stencilDBConn); migratedData = getMigratedData(migrationID, stencilDBConn) {
		for _, oneMigratedData := range migratedData {
			checkDisplayOneMigratedData(stencilDBConn, destAppDBConn, appConfig, oneMigratedData, migratedData, dstApp, secondRound)
		}
		time.Sleep(checkInterval)
	}

	secondRound = true
}


func checkDisplayOneMigratedData(stencilDBConn *sql.DB, destAppDBConn *sql.DB, appConfig config.AppConfig, oneMigratedData display.HintStruct, migratedData []display.HintStruct, dstApp string, secondRound bool) (bool, error) {
	// fmt.Println(oneMigratedData)
	val, err1 := strconv.Atoi(oneMigratedData.Value)
	if err1 != nil {
		log.Fatal("Check  Display One Data: Converting '%s' to Integer Errors", oneMigratedData.Value)
	}
	displayed, err2 := display.CheckDisplayFlag(stencilDBConn, val, oneMigratedData.Table)
	// fmt.Println(displayed)
	if err2 != nil {
		fmt.Println(err2)
	} else {
		if displayed {
			return true, nil
		} else {
			if !dependency_handler.CheckNodeComplete(appConfig.Tags, oneMigratedData, dstApp, destAppDBConn) {
				return false, nil
			} else {

			}
		}
	}
	return false, nil
}

// func DisplayController(migrationID int) {
// 	for migratedNode := GetMigratedData(migrationID); 
// 		!IsMigrationComplete(migrationID);  
// 		migratedNode = GetMigratedData(migrationID){
// 		if migratedNode {
// 			go CheckDisplay(migratedNode. false)
// 		}
// 	}
// 	// Only Executed After The Migration Is Complete
// 	// Remaning Migration Nodes:
// 	// -> The Migrated Nodes In The Destination Application That Still Have Their Migration Flags Raised
// 	for migratedNode := range GetRemainingMigratedNodes(migrationID){
// 		go CheckDisplay(migratedNode, true)
// 	}
// }

// func CheckDisplay(node *DependencyNode, finalRound bool) bool {
// 	try:
// 		if AlreadyDisplayed(node) {
// 			return true
// 		}
// 		if t.Root == node.GetParent() {
// 			Display(node)
// 			return true
// 		} else {
// 			if CheckDisplay(node.GetParent(), finalRound) {
// 				Display(node)
// 				return true
// 			}
// 		}
// 		if finalRound && node.DisplayFlag {
// 			Display(node)
// 			return true
// 		}
// 		return  false
// 	catch NodeNotFound:
// 		return false
// }

func main() {
	dstApp := "mastodon"
	// DisplayThread(dstApp, 808810123)

	dbConn := db.GetDBConn(dstApp)
	if appConfig, err := config.CreateAppConfig(dstApp); err != nil {
		fmt.Println(err)
	} else {
		// fmt.Println(appConfig)
		// fmt.Println(appConfig.Tags)
		hint := display.HintStruct {
			Table: "accounts",
			Key: "id",
			Value: "62632", 
			ValueType: "int",
		} 
		dependency_handler.CheckNodeComplete(appConfig.Tags, hint, dstApp, dbConn)
	}

	// atomicity.CreateTxnLogTable()

	// dbConn := db.GetDBConn(StencilDBName)
	// data := getMigratedData(1134814368, dbConn)
	// fmt.Println(data)

	// var displayHints []display.HintStruct 
	// json.Unmarshal([]byte(data[2].data.String), &displayHints)

	// fmt.Println(displayHints)
	// fmt.Println(displayHints[0].Table)

	// fmt.Println(checkMigrationComplete(1134814368, dbConn))
}

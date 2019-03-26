package main

import (
	"fmt"
	"log"
	// "transaction/config"
	// "transaction/atomicity"
	"transaction/db"
	"database/sql"
	"time"
	"encoding/json"
)

const StencilDBName = "stencil"
const maxDataPerThread = 20
const checkInterval = 200 * time.Millisecond

type HintStruct struct {
	Table string		`json:"Table"`
	Key string			`json:"Key"`
	Value string		`json:"Value"`
	ValueType string	`json:"ValueType"`
} 

type MigratedData struct {
	log_id int
	data sql.NullString
}

var displayedData = make(map[string]int)

func procData(rawData []sql.NullString) []HintStruct {
	var processedData []HintStruct
	for _, oneData := range rawData {
		var oneSetOfHints []HintStruct
		json.Unmarshal([]byte(oneData.String), &oneSetOfHints)
		for _, hint := range oneSetOfHints {
			processedData = append(processedData, hint)
		}
	}
	return processedData
}

func getMigratedData(migrationID int, dbConn *sql.DB) []HintStruct {
	var displayHints []sql.NullString
	var hintString sql.NullString

	op := fmt.Sprintf("SELECT display_hint FROM txn_log WHERE action_id = %d", migrationID)
	rows, err := dbConn.Query(op)
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

// func distributeCheckDisplayTasks(migratedData []MigratedData, secondRound bool) {
// 	j := 0
// 	for i := 0; i < len(migratedData); i++ {
// 		if (i + 1) % maxDataPerThread == 0 {
// 			go checkDisplay(migratedData[j:i+1], secondRound)
// 			j = i + 1
// 		}
// 	}
// 	if j != len(migratedData) {
// 		go checkDisplay(migratedData[j:len(migratedData)], secondRound)
// 	}
// }

func DisplayThread(migrationID int) {
	var secondRound bool
	dbConn := db.GetDBConn(StencilDBName)

	// For now just assume this is an infinite loop
	for migratedData := getMigratedData(migrationID, dbConn); checkMigrationComplete(migrationID, dbConn); migratedData = getMigratedData(migrationID, dbConn) {
		// distributeCheckDisplayTasks(migratedData, secondRound)
		checkDisplay(migratedData, secondRound)
		time.Sleep(checkInterval)
	}

	secondRound = true
}

// func alreadyDisplayed(data) (bool, error) {

// }

// func checkDisplayOneData(oneData ,secondRound bool) bool {
// 	displayed, err = alreadyDisplayed(oneData)
// 	if err != nil {
// 		fmt.Println(err)
// 		continue
// 	}
// 	if displayed {
// 		return true
// 	} else {

// 	}
// }

func checkDisplay(migratedData []HintStruct, secondRound bool) {
	for _, oneData := range migratedData {
		fmt.Println(oneData)
		// for _, oneData := range oneSetOfData {
		// 	fmt.Println(oneData)
		// 	// checkDisplayOneData(oneData, secondRound)
		// }
	}
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
	// atomicity.CreateTxnLogTable()
	DisplayThread(1134814368)

	// dbConn := db.GetDBConn(StencilDBName)
	// data := getMigratedData(1134814368, dbConn)
	// fmt.Println(data)

	// var displayHints []HintStruct 
	// json.Unmarshal([]byte(data[2].data.String), &displayHints)

	// fmt.Println(displayHints)
	// fmt.Println(displayHints[0].Table)

	// fmt.Println(checkMigrationComplete(1134814368, dbConn))
}

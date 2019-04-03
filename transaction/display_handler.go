package main

import (
	"fmt"
	// "log"
	// "transaction/atomicity"
	// "transaction/db"
	"transaction/display"
	// "transaction/dependency_handler"
	"transaction/config"
	"database/sql"
	"time"
	// "strconv"
	// "errors"
)

const checkInterval = 200 * time.Millisecond

var displayedData = make(map[string]int)

func DisplayThread(app string, migrationID int) {
	stencilDBConn, appDBConn, appConfig, pks := display.Initialize(app)

	// For now just assume this is an infinite loop
	var secondRound bool
	for migratedData := display.GetUndisplayedMigratedData(stencilDBConn, app, migrationID, pks); display.CheckMigrationComplete(stencilDBConn, migrationID); migratedData = display.GetUndisplayedMigratedData(stencilDBConn, app, migrationID, pks) {
		for _, oneMigratedData := range migratedData {
			// fmt.Println(oneMigratedData)
			checkDisplayOneMigratedData(stencilDBConn, appDBConn, appConfig, oneMigratedData, migratedData, app, secondRound)
		}
		time.Sleep(checkInterval)
	}

	secondRound = true
}


func checkDisplayOneMigratedData(stencilDBConn *sql.DB, appDBConn *sql.DB, appConfig config.AppConfig, oneMigratedData display.HintStruct, migratedData []display.HintStruct, app string, secondRound bool) (bool, error) {
	// fmt.Println(oneMigratedData)
	var val int
	for _, v := range oneMigratedData.KeyVal {
		val = v
	}
	displayed, _ := display.GetDisplayFlag(stencilDBConn, app, oneMigratedData.Table, val)
	fmt.Println(displayed)
	// if err2 != nil {
	// 	fmt.Println(err2)
	// } else {
	// 	if displayed {
	// 		return true, nil
	// 	} else {
	// 		if !dependency_handler.CheckNodeComplete(appConfig.Tags, oneMigratedData, app, appDBConn) {
	// 			return false, nil
	// 		} else {
	// 			// dependency_handler.GetOneDataFromParentNode()
	// 		}
	// 	}
	// }
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
	DisplayThread(dstApp, 534782464)

	// dbConn := db.GetDBConn(dstApp)
	// if appConfig, err := config.CreateAppConfig(dstApp); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	// fmt.Println(appConfig)
	// 	// fmt.Println(appConfig.Tags)
	// 	keyVal := map[string]int {
	// 		"id": 62632,
	// 	}
	// 	hint := display.HintStruct {
	// 		Table: "accounts",
	// 		KeyVal: keyVal,
	// 	} 
	// 	dependency_handler.CheckNodeComplete(dbConn, appConfig.Tags, hint, dstApp)
	// }

	// dstApp := "mastodon"
	// dbConn := db.GetDBConn(dstApp)
	// if appConfig, err := config.CreateAppConfig(dstApp); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	// fmt.Println(appConfig)
	// 	// fmt.Println(appConfig.Tags)
	// 	keyVal := map[string]int {
	// 		"id": 23550,
	// 	}
	// 	hint := display.HintStruct {
	// 		Table: "statuses",
	// 		KeyVal: keyVal,
	// 	} 
	// 	// hint := display.HintStruct {
	// 	// 	Table: "conversations",
	// 	// 	Key: "id",
	// 	// 	Value: "211",
	// 	// 	ValueType: "int",
	// 	// }
	// 	dependency_handler.GetOneDataFromParentNode(appConfig, hint, dstApp, dbConn)
	// }

	// atomicity.CreateTxnLogTable()

	// dbConn := db.GetDBConn(StencilDBName)
	// data := getMigratedData("mastodon", 1134814368, dbConn)
	// fmt.Println(data)

	// var displayHints []display.HintStruct 
	// json.Unmarshal([]byte(data[2].data.String), &displayHints)

	// fmt.Println(displayHints)
	// fmt.Println(displayHints[0].Table)

	// fmt.Println(checkMigrationComplete(1134814368, dbConn))

	// display.CreateDisplayFlagsTable(dbConn)
}

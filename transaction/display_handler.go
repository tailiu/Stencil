package main

import (
	"fmt"
	// "log"
	// "transaction/atomicity"
	"transaction/db"
	"transaction/display"
	// "transaction/dependency_handler"
	// "transaction/config"
	// "database/sql"
	"time"
	// "strconv"
	// "errors"
)

const StencilDBName = "stencil"
const checkInterval = 200 * time.Millisecond

var displayedData = make(map[string]int)

func DisplayThread(app string, migrationID int) {
	// var secondRound bool

	stencilDBConn := db.GetDBConn(StencilDBName)
	// destAppDBConn := db.GetDBConn(app)
	// appConfig, err := config.CreateAppConfig(app)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// For now just assume this is an infinite loop
	for migratedData := display.GetMigratedData(stencilDBConn, app, migrationID); display.CheckMigrationComplete(stencilDBConn, migrationID); migratedData = display.GetMigratedData(stencilDBConn, app, migrationID) {
		for _, oneMigratedData := range migratedData {
			fmt.Println(oneMigratedData)
			// checkDisplayOneMigratedData(stencilDBConn, destAppDBConn, appConfig, oneMigratedData, migratedData, app, secondRound)
		}
		time.Sleep(checkInterval)
	}

	// secondRound = true
}


// func checkDisplayOneMigratedData(stencilDBConn *sql.DB, destAppDBConn *sql.DB, appConfig config.AppConfig, oneMigratedData display.HintStruct, migratedData []display.HintStruct, dstApp string, secondRound bool) (bool, error) {
// 	// fmt.Println(oneMigratedData)
// 	val, err1 := strconv.Atoi(oneMigratedData.Value)
// 	if err1 != nil {
// 		log.Fatal("Check  Display One Data: Converting '%s' to Integer Errors", oneMigratedData.Value)
// 	}
// 	displayed, err2 := display.GetDisplayFlag(stencilDBConn, val, oneMigratedData.Table)
// 	// fmt.Println(displayed)
// 	if err2 != nil {
// 		fmt.Println(err2)
// 	} else {
// 		if displayed {
// 			return true, nil
// 		} else {
// 			if !dependency_handler.CheckNodeComplete(appConfig.Tags, oneMigratedData, dstApp, destAppDBConn) {
// 				return false, nil
// 			} else {
// 				// dependency_handler.GetOneDataFromParentNode()
// 			}
// 		}
// 	}
// 	return false, nil
// }

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
	// 	hint := display.HintStruct {
	// 		Table: "accounts",
	// 		Key: "id",
	// 		Value: "62632", 
	// 		ValueType: "int",
	// 	} 
	// 	dependency_handler.CheckNodeComplete(appConfig.Tags, hint, dstApp, dbConn)
	// }

	// dstApp := "mastodon"
	// dbConn := db.GetDBConn(dstApp)
	// if appConfig, err := config.CreateAppConfig(dstApp); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(appConfig)
	// 	fmt.Println(appConfig.Tags)
	// 	hint := display.HintStruct {
	// 		Table: "statuses",
	// 		Key: "id",
	// 		Value: "23550", 
	// 		ValueType: "int",
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

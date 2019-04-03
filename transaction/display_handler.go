package main

import (
	"fmt"
	// "log"
	// "transaction/atomicity"
	// "transaction/db"
	"transaction/display"
	"transaction/dependency_handler"
	"transaction/config"
	"database/sql"
	"time"
	// "strconv"
	"errors"
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
			checkDisplayOneMigratedData(stencilDBConn, appDBConn, appConfig, oneMigratedData, migratedData, app, pks, secondRound)
		}
		time.Sleep(checkInterval)
	}

	secondRound = true
}


func checkDisplayOneMigratedData(stencilDBConn *sql.DB, appDBConn *sql.DB, appConfig config.AppConfig, oneMigratedData display.HintStruct, migratedData []display.HintStruct, app string, pks map[string]string, secondRound bool) (bool, error) {
	// fmt.Println(oneMigratedData)
	var val int
	for _, v := range oneMigratedData.KeyVal {
		val = v
	}
	displayed, err1 := display.GetDisplayFlag(stencilDBConn, app, oneMigratedData.Table, val)

	// fmt.Println(displayed)
	if err1 != nil {
		fmt.Println(err1)
		return false, err1
	} else {
		if displayed {
			return true, nil
		} else {
			complete, completeDataHints := dependency_handler.CheckNodeComplete(appDBConn, appConfig.Tags, oneMigratedData, app)
			if !complete {
				return false, errors.New("Data of a Node is Not Complete")
			} else {
				tags, err2 := dependency_handler.GetParentTags(appConfig, oneMigratedData)
				if err2 != nil {
					fmt.Println(err2)
					return false, err2
				} else {
					// This should not happen in Stencil case, because root node data should
					// be separated stored
					if tags == nil {
						fmt.Println("This Data Already Belongs To Root Node!")
						return true, nil
					} else {
						for _, tag := range tags {
							if tag == "root" {
								err3 := display.Display(stencilDBConn, app, completeDataHints, pks)
								if err3 != nil {
									fmt.Println(err3)
									return false, err3
								} else {
									return true, nil
								}
							}
						}
						oneDataInParentNode, err4 := dependency_handler.GetOneDataFromParentNodeRandomly(appDBConn, appConfig, oneMigratedData, app)
						if err4 != nil {
							fmt.Println(err4)
							return false, err4
						} else {
							result, err5 := checkDisplayOneMigratedData(stencilDBConn, appDBConn, appConfig, oneDataInParentNode, migratedData, app, pks, secondRound)
							if err5 != nil {
								fmt.Println(err5)
								return false, err5
							} else {
								if result {
									err6 := display.Display(stencilDBConn, app, completeDataHints, pks)
									if err6 != nil {
										fmt.Println(err6)
										return false, err6
									} else {
										return true, nil
									}
								} else {
									return false, nil
								}
							}
						}
					}
				}
			}
		}
	}
	return false, nil
}

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

func main() {
	dstApp := "mastodon"
	DisplayThread(dstApp, 534782464)

	// var completeDataHints []display.HintStruct
	// stencilDBConn, _, _, pks := display.Initialize(dstApp)
	// display.Display(stencilDBConn, dstApp, completeDataHints, pks)

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

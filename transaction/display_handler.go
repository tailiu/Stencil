package main

import (
	"fmt"
	"log"
	// "transaction/atomicity"
	// "transaction/db"
	"transaction/display"
	"transaction/dependency_handler"
	"transaction/config"
	"database/sql"
	"time"
	"errors"
)

const checkInterval = 200 * time.Millisecond

var displayedData = make(map[string]int)

func DisplayThread(app string, migrationID int) {
	stencilDBConn, appDBConn, appConfig, pks := display.Initialize(app)
	
	fmt.Println("--------- First Phase --------")
	secondRound := false
	for migratedData := display.GetUndisplayedMigratedData(stencilDBConn, app, migrationID, pks); 
			!display.CheckMigrationComplete(stencilDBConn, migrationID);
			migratedData = display.GetUndisplayedMigratedData(stencilDBConn, app, migrationID, pks) {
		
		for _, oneMigratedData := range migratedData {
			// fmt.Println(oneMigratedData)
			checkDisplayOneMigratedData(stencilDBConn, appDBConn, appConfig, oneMigratedData, app, pks, secondRound)
		}
		time.Sleep(checkInterval)
	}

	fmt.Println("--------- Second Phase ---------")
	secondRound = true
	secondRoundMigratedData := display.GetUndisplayedMigratedData(stencilDBConn, app, migrationID, pks)
	for _, oneSecondRoundMigratedData := range secondRoundMigratedData {
		checkDisplayOneMigratedData(stencilDBConn, appDBConn, appConfig, oneSecondRoundMigratedData, app, pks, secondRound)
	}
}


func checkDisplayOneMigratedData(stencilDBConn *sql.DB, appDBConn *sql.DB, appConfig config.AppConfig, oneMigratedData display.HintStruct, app string, pks map[string]string, secondRound bool) (bool, error) {
	fmt.Println(oneMigratedData)
	var val int
	for _, v := range oneMigratedData.KeyVal {
		val = v
	}
	displayed, err1 := display.GetDisplayFlag(stencilDBConn, app, oneMigratedData.Table, val)

	fmt.Println(displayed)
	if err1 != nil {
		log.Fatal(err1)
		return false, err1
	} else {
		if displayed {
			return true, nil
		} else {
			// This should be different for the second round because, based on config, nodes could be displayed despite incomplete 
			complete, completeDataHints := dependency_handler.CheckNodeComplete(appDBConn, appConfig.Tags, oneMigratedData, app)
			fmt.Println(complete, completeDataHints)
			if !complete {
				return false, errors.New("Data of a Node is Not Complete")
			} else {
				tags, err2 := dependency_handler.GetParentTags(appConfig, oneMigratedData)
				if err2 != nil {
					log.Fatal(err2)
					return false, err2
				} else {
					// This should not happen in Stencil's case, because root node data should
					// be stored separatedly
					if tags == nil {
						log.Fatal("This Data Already Belongs To Root Node!")
						return true, nil
					} else {
						for _, tag := range tags {
							if tag == "root" {
								err3 := display.Display(stencilDBConn, app, completeDataHints, pks)
								if err3 != nil {
									log.Fatal(err3)
									return false, err3
								} else {
									return true, nil
								}
							}
						}
						// This function should also be different for the second round
						// because we may end up with always getting some data in a node that could not be displayed but other data in that  
						// node may have already been displayed
						oneDataInParentNode, err4 := dependency_handler.GetOneDataFromParentNodeRandomly(appDBConn, appConfig, oneMigratedData, app)
						if err4 != nil {
							log.Fatal(err4)
							return false, err4
						} else {
							result, err5 := checkDisplayOneMigratedData(stencilDBConn, appDBConn, appConfig, oneDataInParentNode, app, pks, secondRound)
							if err5 != nil {
								log.Fatal(err5)
								return false, err5
							} else {
								if result {
									err6 := display.Display(stencilDBConn, app, completeDataHints, pks)
									if err6 != nil {
										log.Fatal(err6)
										return false, err6
									} else {
										return true, nil
									}
								} else {
									if secondRound && dependency_handler.CheckDisplayCondition() {
										err6 := display.Display(stencilDBConn, app, completeDataHints, pks)
										if err6 != nil {
											log.Fatal(err6)
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
	DisplayThread(dstApp, 857232446)

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
	// 		"id": 440047296002523137,
	// 	}
	// 	hint := display.HintStruct {
	// 		Table: "users",
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

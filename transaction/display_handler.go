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
	// "errors"
)

const checkInterval = 200 * time.Millisecond

func DisplayThread(app string, migrationID int) {
	fmt.Println("--------- Start of Display Check ---------")

	stencilDBConn, appDBConn, appConfig, pks := display.Initialize(app)
	
	fmt.Println("--------- First Phase --------")
	secondRound := false
	for migratedData := display.GetUndisplayedMigratedData(stencilDBConn, app, migrationID, pks); 
			!display.CheckMigrationComplete(stencilDBConn, migrationID);
			migratedData = display.GetUndisplayedMigratedData(stencilDBConn, app, migrationID, pks) {
		
		for _, oneMigratedData := range migratedData {
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

	fmt.Println("--------- End of Display Check ---------")
}

func checkDisplayOneMigratedData(stencilDBConn *sql.DB, appDBConn *sql.DB, appConfig config.AppConfig, oneMigratedData display.HintStruct, app string, pks map[string]string, secondRound bool) (bool, error) {
	fmt.Println("Check Data ", oneMigratedData)
	var val int
	for _, v := range oneMigratedData.KeyVal {
		val = v
	}
	displayed, err0 := display.GetDisplayFlag(stencilDBConn, app, oneMigratedData.Table, val)

	fmt.Println("Displayed: ", displayed)
	if err0 != nil {
		log.Println(err0)
		return false, err0
	} else {
		// There could be a case where only this data of a node has been displayed and other data in that node
		// is not displayed, but we still return true because  
		if displayed {
			return true, nil
		} else {


			dataInNode, err1 := dependency_handler.GetDataInNodeBasedOnDisplaySetting(&appConfig, oneMigratedData)
			if dataInNode == nil {
				log.Println(err1)
				return false, err1
			} else {
				for _, oneDataInNode := range dataInNode {
					var val int
					for _, v := range oneMigratedData.KeyVal {
						val = v
					}
					displayed, err0 := display.GetDisplayFlag(stencilDBConn, app, oneDataInNode.Table, val)
				}
				tags, err2 := oneMigratedData.GetParentTags(appConfig)
				if err2 != nil {
					log.Println(err2)
					return false, err2
				} else {
					if tags == nil {
						log.Println("This Data's Tag Does not Depend on Any Other Tag!")
						err3 := display.Display(stencilDBConn, app, dataInNode, pks)
						if err3 != nil {
							log.Println(err3)
							return false, err3
						} else {
							return true, nil
						}
					} else {
						// This function should also be different for the second round
						// because we may end up with always getting some data in a node that could not be displayed but other data in that  
						// node may have already been displayed
						oneDataInParentNode, err4 := dependency_handler.GetOneDataFromParentNodeRandomly(appDBConn, appConfig, oneMigratedData, app)
						if err4 != nil {
							log.Println(err4)
							return false, err4
						} else {
							result, err5 := checkDisplayOneMigratedData(stencilDBConn, appDBConn, appConfig, oneDataInParentNode, app, pks, secondRound)
							if err5 != nil {
								log.Println(err5)
								return false, err5
							} else {
								if result {
									err6 := display.Display(stencilDBConn, app, dataInNode, pks)
									if err6 != nil {
										log.Println(err6)
										return false, err6
									} else {
										return true, nil
									}
								} else {
									if secondRound && dependency_handler.CheckDisplayCondition() {
										err6 := display.Display(stencilDBConn, app, dataInNode, pks)
										if err6 != nil {
											log.Println(err6)
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

func CheckDisplay(oneUndisplayedMigratedData dataStruct, finalRound bool) bool {
	dataInNode, err := GetDataInNodeBasedOnDisplaySetting(oneUndisplayedMigratedData)
	if dataInNode == nil {
		return false
	} else {
		for _, oneDataInNode := range dataInNode {

		}
	}
	if AlreadyDisplayed(node) {
		return true
	}
	if t.Root == node.GetParent() {
		Display(node)
		return true
	} else {
		if CheckDisplay(node.GetParent(), finalRound) {
			Display(node)
			return true
		}
	}
	if finalRound && node.DisplayFlag {
		Display(node)
		return true
	}
	return  false
}

func DisplayController(migrationID int) {
	for undisplayedMigratedData := GetUndisplayedMigratedData(migrationID); 
		!CheckMigrationComplete(migrationID);  
		undisplayedMigratedData = GetUndisplayedMigratedData(migrationID){
			for _, oneUndisplayedMigratedData := range undisplayedMigratedData {
				CheckDisplay(oneUndisplayedMigratedData, false)
			}
	}
	// Only Executed After The Migration Is Complete
	// Remaning Migration Nodes:
	// -> The Migrated Nodes In The Destination Application That Still Have Their Migration Flags Raised
	for _, oneUndisplayedMigratedData := range GetUndisplayedMigratedData(migrationID){
		CheckDisplay(oneUndisplayedMigratedData, true)
	}
}

func main() {
	dstApp := "mastodon"
	// DisplayThread(dstApp, 857232446)

	// var dataInNode []display.HintStruct
	// stencilDBConn, _, _, pks := display.Initialize(dstApp)
	// display.Display(stencilDBConn, dstApp, dataInNode, pks)

	// dbConn := db.GetDBConn(dstApp)
	if appConfig, err := config.CreateAppConfig(dstApp); err != nil {
		fmt.Println(err)
	} else {
		// keyVal := map[string]int {
		// 	"id": 14435263,
		// }
		// hint := display.HintStruct {
		// 	Table: "favourites",
		// 	KeyVal: keyVal,
		// } 

		// keyVal := map[string]int {
		// 	"id": 4630,
		// }
		// hint := display.HintStruct {
		// 	Table: "accounts",
		// 	KeyVal: keyVal,
		// } 

		keyVal := map[string]int {
			"id": 28300,
		}
		hint := display.HintStruct {
			Table: "status_stats",
			KeyVal: keyVal,
		} 
		// dependency_handler.CheckNodeComplete(dbConn, appConfig.Tags, hint, dstApp)
		data, err := dependency_handler.GetDataInNodeBasedOnDisplaySetting(&appConfig, hint)
		if err != nil {
			fmt.Println(err)
			fmt.Println(data)
		} else {
			fmt.Println(data)
		}
	}

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

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

func returnResultBasedOnNodeCompleteness(err error) (string, error) {
	if err != nil {
		return "Data In a Node Can be partially Displayed", err
	} else {
		return "Data In a Node Can be completely Displayed", nil
	}
}

func returnDisplayConditionWhenCannotGetDataFromParentNode(displaySetting string, secondRound bool) bool {
	if !secondRound {
		if displaySetting != "parent_node_not_displays_without_check" {
			return true
		} else {
			return false
		}
	} else {
		if displaySetting == "parent_node_not_displays_with_check" or displaySetting == "parent_node_not_displays_without_check" {
			return true
		} else {
			return false
		}
	}
}

func returnDisplayConditionWhenGetPartialDataFromParentNode(displaySetting string) bool {
	if displaySetting != "parent_node_complete_displays" {
		return true
	} else {
		return false
	}
}

func checkDisplayConditions() bool {

}

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

func checkDisplayOneMigratedData(stencilDBConn *sql.DB, appDBConn *sql.DB, appConfig config.AppConfig, oneMigratedData display.HintStruct, app string, pks map[string]string, secondRound bool) (string, error) {
	// var val int
	// for _, v := range oneMigratedData.KeyVal {
	// 	val = v
	// }
	// displayed, err0 := display.GetDisplayFlag(stencilDBConn, app, oneMigratedData.Table, val)

	// fmt.Println("Displayed: ", displayed)
	// if err0 != nil {
	// 	log.Println(err0)
	// 	return false, err0
	// } else {
	// 	// There could be a case where only this data of a node has been displayed and other data in that node
	// 	// is not displayed, but we still return true because  
	// 	if displayed {
	// 		return true, nil
	// 	} else {

			fmt.Println("Check Data ", oneMigratedData)
			dataInNode, err1 := dependency_handler.GetDataInNodeBasedOnDisplaySetting(&appConfig, oneMigratedData)
			if dataInNode == nil {
				log.Println(err1)
				return "No Data In a Node Can be Displayed", err1
			} else {

				var displayedData, notDisplayedData []display.HintStruct
				for _, oneDataInNode := range dataInNode {
					var val int
					for _, v := range oneDataInNode.KeyVal {
						val = v
					}
					displayed, err0 := display.GetDisplayFlag(stencilDBConn, app, oneDataInNode.Table, val)
					if err0 != nil {
						log.Fatal(err0)
					}
					if !displayed {
						notDisplayedData = append(notDisplayedData, oneDataInNode)
					} else {
						displayedData = append(displayedData, oneDataInNode)
					}
				}
				// Note: This will be changed when considering ongoing application services 
				// and the existence of other display threads !!
				if len(displayedData) != 0 {
					err6 := display.Display(stencilDBConn, app, notDisplayedData, pks)
					if err6 != nil {
						log.Fatal(err6)
					} 
					returnResultBasedOnNodeCompleteness(err1)
				} 

				pTags, err2 := oneMigratedData.GetParentTags(&appConfig)
				if err2 != nil {
					log.Fatal(err2)
				} else {
					if pTags == nil {
						log.Println("This Data's Tag Does not Depend on Any Other Tag!")
						err3 := display.Display(stencilDBConn, app, dataInNode, pks)
						if err3 != nil {
							log.Fatal(err3)
						}
						returnResultBasedOnNodeCompleteness(err1)
					} else {
						pTagConditions := make(map[string]bool)
						for _, pTag := range pTags {
							dataInParentNode, err4 := dependency_handler.GetdataFromParentNode(&appConfig, oneMigratedData, pTag)
							// fmt.Println(dataInParentNode, err4)
							displaySetting, err5 := dependency_handler.GetDisplaySettingInDependencies(&appConfig, oneMigratedData, pTag)
							if err5 != nil {
								log.Fatal(err5)
							}
							if err4 != nil {
								switch err4 {
								case errors.New("This Data Does not Depend on Any Data in the Parent Node"):
									pTagConditions[pTag] = true
								case errors.New("Fail To Get Any Data in the Parent Node"):
									pTagConditions[pTag] = returnDisplayConditionWhenCannotGetDataFromParentNode(displaySetting, secondRound)
								}
							} else {
								// For now, there is no case where there is more than one piece of data in a parent node
								if len(dataInParentNode) != 1 {
									log.Fatal("Find more than one piece of data in a parent node!!")
								}
								result, err7 := checkDisplayOneMigratedData(stencilDBConn, appDBConn, appConfig, dataInParentNode[0], app, pks, secondRound)
								if err7 != nil {
									log.Println(err7)
								}
								switch result {
								case "No Data In a Node Can be Displayed":
									pTagConditions[pTag] = returnDisplayConditionWhenCannotGetDataFromParentNode(displaySetting, secondRound)
								case "Data In a Node Can be partially Displayed":
									pTagConditions[pTag] = returnDisplayConditionWhenGetPartialDataFromParentNode(displaySetting)
								case "Data In a Node Can be completely Displayed":
									pTagConditions[pTag] = true
								}
							}
						}
						checkResult := checkDisplayConditions(pTagConditions)
						


						// This function should also be different for the second round
						// because we may end up with always getting some data in a node that could not be displayed but other data in that  
						// node may have already been displayed
						// oneDataInParentNode, err4 := dependency_handler.GetOneDataFromParentNodeRandomly(appDBConn, appConfig, oneMigratedData, app)
						// if err4 != nil {
						// 	log.Println(err4)
						// 	return false, err4
						// } else {
						// 	result, err5 := checkDisplayOneMigratedData(stencilDBConn, appDBConn, appConfig, oneDataInParentNode, app, pks, secondRound)
						// 	if err5 != nil {
						// 		log.Println(err5)
						// 		return false, err5
						// 	} else {
						// 		if result {
						// 			err6 := display.Display(stencilDBConn, app, dataInNode, pks)
						// 			if err6 != nil {
						// 				log.Println(err6)
						// 				return false, err6
						// 			} else {
						// 				return true, nil
						// 			}
						// 		} else {
						// 			if secondRound && dependency_handler.CheckDisplayCondition() {
						// 				err6 := display.Display(stencilDBConn, app, dataInNode, pks)
						// 				if err6 != nil {
						// 					log.Println(err6)
						// 					return false, err6
						// 				} else {
						// 					return true, nil
						// 				}
						// 			} else {
						// 				return false, nil
						// 			}
						// 		}
						// 	}
						// }
					}
				}
			}
	// 	}
	// }
	return "", nil
}

// func CheckDisplay(oneUndisplayedMigratedData dataStruct, finalRound bool) bool {
// 	dataInNode, err := GetDataInNodeBasedOnDisplaySetting(oneUndisplayedMigratedData)
// 	if dataInNode == nil {
// 		return false
// 	} else {
// 		for _, oneDataInNode := range dataInNode {

// 		}
// 	}
// 	if AlreadyDisplayed(node) {
// 		return true
// 	}
// 	if t.Root == node.GetParent() {
// 		Display(node)
// 		return true
// 	} else {
// 		if CheckDisplay(node.GetParent(), finalRound) {
// 			Display(node)
// 			return true
// 		}
// 	}
// 	if finalRound && node.DisplayFlag {
// 		Display(node)
// 		return true
// 	}
// 	return  false
// }

// func DisplayController(migrationID int) {
// 	for undisplayedMigratedData := GetUndisplayedMigratedData(migrationID); 
// 		!CheckMigrationComplete(migrationID);  
// 		undisplayedMigratedData = GetUndisplayedMigratedData(migrationID){
// 			for _, oneUndisplayedMigratedData := range undisplayedMigratedData {
// 				CheckDisplay(oneUndisplayedMigratedData, false)
// 			}
// 	}
// 	// Only Executed After The Migration Is Complete
// 	// Remaning Migration Nodes:
// 	// -> The Migrated Nodes In The Destination Application That Still Have Their Migration Flags Raised
// 	for _, oneUndisplayedMigratedData := range GetUndisplayedMigratedData(migrationID){
// 		CheckDisplay(oneUndisplayedMigratedData, true)
// 	}
// }

func main() {
	// dstApp := "mastodon"
	// DisplayThread(dstApp, 857232446)

	// // var dataInNode []display.HintStruct
	// // stencilDBConn, _, _, pks := display.Initialize(dstApp)
	// // display.Display(stencilDBConn, dstApp, dataInNode, pks)

	// dbConn := db.GetDBConn(dstApp)
	// if appConfig, err := config.CreateAppConfig(dstApp); err != nil {
	// 	fmt.Println(err)
	// } else {

	// 	// keyVal := map[string]int {
	// 	// 	"id": 14435263,
	// 	// }
	// 	// hint := display.HintStruct {
	// 	// 	Table: "favourites",
	// 	// 	KeyVal: keyVal,
	// 	// } 

	// 	// keyVal := map[string]int {
	// 	// 	"id": 4630,
	// 	// }
	// 	// hint := display.HintStruct {
	// 	// 	Table: "accounts",
	// 	// 	KeyVal: keyVal,
	// 	// } 

	// 	keyVal := map[string]int {
	// 		"id": 28300,
	// 	}
	// 	hint := display.HintStruct {
	// 		Table: "status_stats",
	// 		KeyVal: keyVal,
	// 	} 
	// 	fmt.Println(dependency_handler.GetDisplaySettingInDependencies(&appConfig, hint, "S1"))
	// 	// dependency_handler.CheckNodeComplete(dbConn, appConfig.Tags, hint, dstApp)
	// 	data, err := dependency_handler.GetDataInNodeBasedOnDisplaySetting(&appConfig, hint)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		fmt.Println(data)
	// 	} else {
	// 		fmt.Println(data)
	// 	}
	// }

	// dstApp := "mastodon"
	// if appConfig, err := config.CreateAppConfig(dstApp); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	// fmt.Println(appConfig)
	// 	// fmt.Println(appConfig.Tags)
	// 	keyVal := map[string]int {
	// 		"id": 32999907,
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
	// 	fmt.Println(hint.GetParentTags(&appConfig))
	// 	data, err := dependency_handler.GetdataFromParentNode(&appConfig, hint, "S2")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	} else {
	// 		fmt.Println(data)
	// 	}
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

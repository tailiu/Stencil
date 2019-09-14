package main

import (
	"database/sql"
	"errors"
	"log"
	"time"
	"stencil/config"
	"stencil/dependency_handler"
	"stencil/display"
	"fmt"
)

const checkInterval = 200 * time.Millisecond

func DisplayThread(app string, migrationID int, deletionHoldEnable bool) {
	startTime := time.Now()
	log.Println("--------- Start of Display Check ---------")

	stencilDBConn, appConfig, threadID := display.Initialize(app)

	// display.CreateDeletionHoldTable(stencilDBConn)
	log.Println("Thread ID:", threadID)

	log.Println("--------- First Phase --------")
	secondRound := false
	for migratedData := display.GetUndisplayedMigratedData(stencilDBConn, migrationID, appConfig); 
		!display.CheckMigrationComplete(stencilDBConn, migrationID); 
		migratedData = display.GetUndisplayedMigratedData(stencilDBConn, migrationID, appConfig) {
		
		var dhStack [][]int
		for _, oneMigratedData := range migratedData {
			_, dhStack, _ = checkDisplayOneMigratedData(stencilDBConn, appConfig, oneMigratedData, secondRound, deletionHoldEnable, dhStack, threadID)
			if deletionHoldEnable {
				display.RemoveDeletionHold(stencilDBConn, dhStack, threadID)
			}
		}
		time.Sleep(checkInterval)
	}

	log.Println("--------- Second Phase ---------")
	secondRound = true
	secondRoundMigratedData := display.GetUndisplayedMigratedData(stencilDBConn, migrationID, appConfig)
	for _, oneSecondRoundMigratedData := range secondRoundMigratedData {
		var dhStack [][]int
		_, dhStack, _ = checkDisplayOneMigratedData(stencilDBConn, appConfig, oneSecondRoundMigratedData, secondRound, deletionHoldEnable, dhStack, threadID)
		if deletionHoldEnable {
			display.RemoveDeletionHold(stencilDBConn, dhStack, threadID)
		}
	}

	log.Println("--------- End of Display Check ---------")
	endTime := time.Now()
	log.Println("Time used: ", endTime.Sub(startTime))
}

// Two-way display check
func checkDisplayOneMigratedData(stencilDBConn *sql.DB, appConfig *config.AppConfig, oneMigratedData display.HintStruct, secondRound bool, deletionHoldEnable bool, dhStack [][]int, threadID int) (string, [][]int, error) {
	if oneMigratedData.TableName == "" {
		oneMigratedData.TableName = display.GetTableNameByTableID(stencilDBConn, oneMigratedData.TableID)
	}
	dataInNode, err1 := dependency_handler.GetDataInNodeBasedOnDisplaySetting(appConfig, oneMigratedData, stencilDBConn)
	log.Println("-----------")
	log.Println(dataInNode)
	log.Println("-----------")
	if len(dataInNode) == 0 {
		log.Println(err1)
		if secondRound {
			err10 := display.PutIntoDataBag(stencilDBConn, appConfig.AppID, []display.HintStruct{oneMigratedData})
			// Found path conflicts
			if err10 != nil {
				return "No Data In a Node Can be Displayed", dhStack, err10
			} else {
				return "No Data In a Node Can be Displayed", dhStack, err1
			}
		} else {
			return "No Data In a Node Can be Displayed", dhStack, err1
		}
	} else {

		var displayedData, notDisplayedData []display.HintStruct
		for _, dataInNode1 := range dataInNode {
			displayed := display.CheckDisplay(stencilDBConn, appConfig.AppID, dataInNode1)
			if displayed == 1 {
				notDisplayedData = append(notDisplayedData, dataInNode1)
			} else {
				displayedData = append(displayedData, dataInNode1)
			}
		}
		// Note: This will be changed when considering ongoing application services
		// and the existence of other display threads !!
		if len(displayedData) != 0 {
			var err6 error
			err6, dhStack = display.Display(stencilDBConn, appConfig.AppID, notDisplayedData, deletionHoldEnable, dhStack, threadID)
			if err6 != nil {
				return "", dhStack, err6
			}
			return display.ReturnResultBasedOnNodeCompleteness(err1, dhStack)
		}

		pTags, err2 := oneMigratedData.GetParentTags(appConfig)
		if err2 != nil {
			log.Fatal(err2)
		} else {
			if pTags == nil {
				log.Println("This Data's Tag Does not Depend on Any Other Tag!")
				// Need to change this display.Display function
				var err3 error
				err3, dhStack = display.Display(stencilDBConn, appConfig.AppID, dataInNode, deletionHoldEnable, dhStack, threadID)
				// Found path conflicts
				if err3 != nil {
					return "", dhStack, err3
				}
				return display.ReturnResultBasedOnNodeCompleteness(err1, dhStack)
			} else {
				pTagConditions := make(map[string]bool)
				for _, pTag := range pTags {
					log.Println(pTag)
					dataInParentNode, err4 := dependency_handler.GetdataFromParentNode(stencilDBConn, appConfig, dataInNode, pTag)
					log.Println(dataInParentNode, err4)
					displaySetting, err5 := dependency_handler.GetDisplaySettingInDependencies(appConfig, oneMigratedData, pTag)
					if err5 != nil {
						log.Fatal(err5)
					}
					if err4 != nil {
						switch err4.Error() {
						case "This Data Does not Depend on Any Data in the Parent Node":
							pTagConditions[pTag] = true
						case "Fail To Get Any Data in the Parent Node":
							pTagConditions[pTag] = display.ReturnDisplayConditionWhenCannotGetDataFromParentNode(displaySetting, secondRound)
						}
					} else {
						// For now, there is no case where there is more than one piece of data in a parent node
						// if len(dataInParentNode) != 1 {
						// 	log.Fatal("Find more than one piece of data in a parent node!!")
						// }
						var result string
						var err7 error
						result, dhStack, err7 = checkDisplayOneMigratedData(stencilDBConn, appConfig, dataInParentNode, secondRound, deletionHoldEnable, dhStack, threadID)
						if err7 != nil {
							log.Println(err7)
							// If there is path confilct, just return until to the original layer and remove the already added deletion hold
							if err7.Error() == "Path conflict" {
								return "", dhStack, err7
							}
						}
						log.Println(result, err7)
						switch result {
						case "No Data In a Node Can be Displayed":
							pTagConditions[pTag] = display.ReturnDisplayConditionWhenCannotGetDataFromParentNode(displaySetting, secondRound)
						case "Data In a Node Can be partially Displayed":
							pTagConditions[pTag] = display.ReturnDisplayConditionWhenGetPartialDataFromParentNode(displaySetting)
						case "Data In a Node Can be completely Displayed":
							pTagConditions[pTag] = true
						}
					}
				}
				// log.Println(pTagConditions)

				if checkResult := display.CheckCombinedDisplayConditions(appConfig, pTagConditions, oneMigratedData); checkResult {
					var err8 error
					err8, dhStack = display.Display(stencilDBConn, appConfig.AppID, dataInNode, deletionHoldEnable, dhStack, threadID)
					// Found path conflicts
					if err8 != nil {
						return "", dhStack, err8
					}
					return display.ReturnResultBasedOnNodeCompleteness(err1, dhStack)
				} else {
					if secondRound {
						err9 := display.PutIntoDataBag(stencilDBConn, appConfig.AppID, dataInNode)
						// Found path conflicts
						if err9 != nil {
							return "No Data In a Node Can be Displayed", dhStack, err9
						} else {
							return "No Data In a Node Can be Displayed", dhStack, errors.New("Display Setting does not allow the data in the node to be displayed")
						}
					} else {
						return "No Data In a Node Can be Displayed", dhStack, errors.New("Display Setting does not allow the data in the node to be displayed")
					}
				}
			}
		}
	}
	panic("Should never happen here!")
}

func main() {
	threadNum := 1
	dstApp := "mastodon"
	migrationID := 2014441389
	deletionHoldEnable := false

	for i := 0; i < threadNum; i++ {
		go DisplayThread(dstApp, migrationID, deletionHoldEnable)
	}

	for {
		fmt.Scanln()
	}
}

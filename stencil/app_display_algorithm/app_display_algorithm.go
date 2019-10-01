package app_display_algorithm

import (
	"database/sql"
	"errors"
	"log"
	"stencil/config"
	"stencil/app_dependency_handler"
	"stencil/app_display"
	"time"
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
		if displaySetting == "parent_node_not_displays_with_check" || displaySetting == "parent_node_not_displays_without_check" {
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

func checkDisplayConditions(appConfig *config.AppConfig, pTagConditions map[string]bool, oneMigratedData app_display.HintStruct) bool {
	for _, result := range pTagConditions {
		if result {
			return true
		}
	}
	return false
}

func DisplayThread(app string, migrationID int) {
	startTime := time.Now()
	log.Println("--------- Start of Display Check ---------")

	stencilDBConn, appConfig, pks := app_display.Initialize(app, "")

	log.Println("--------- First Phase --------")
	secondRound := false
	for migratedData := app_display.GetUndisplayedMigratedData(stencilDBConn, app, migrationID, pks); 
		!app_display.CheckMigrationComplete(stencilDBConn, migrationID); 
		migratedData = app_display.GetUndisplayedMigratedData(stencilDBConn, app, migrationID, pks) {

		for _, oneMigratedData := range migratedData {
			checkDisplayOneMigratedData(stencilDBConn, appConfig, oneMigratedData, app, pks, secondRound)
		}
		time.Sleep(checkInterval)
	}

	log.Println("--------- Second Phase ---------")
	secondRound = true
	secondRoundMigratedData := app_display.GetUndisplayedMigratedData(stencilDBConn, app, migrationID, pks)
	for _, oneSecondRoundMigratedData := range secondRoundMigratedData {
		checkDisplayOneMigratedData(stencilDBConn, appConfig, oneSecondRoundMigratedData, app, pks, secondRound)
	}

	log.Println("--------- End of Display Check ---------")
	endTime := time.Now()
	log.Println("Time used: ", endTime.Sub(startTime))
}

func checkDisplayOneMigratedData(stencilDBConn *sql.DB, appConfig config.AppConfig, oneMigratedData app_display.HintStruct, app string, pks map[string]string, secondRound bool) (string, error) {

	log.Println("Check Data ", oneMigratedData)
	dataInNode, err1 := app_dependency_handler.GetDataInNodeBasedOnDisplaySetting(&appConfig, oneMigratedData)
	if dataInNode == nil {
		log.Println(err1)
		return "No Data In a Node Can be Displayed", err1
	} else {

		var displayedData, notDisplayedData []app_display.HintStruct
		for _, oneDataInNode := range dataInNode {
			var val int
			for _, v := range oneDataInNode.KeyVal {
				val = v
			}
			displayed, err0 := app_display.GetDisplayFlag(stencilDBConn, app, oneDataInNode.Table, val)
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
		// and the existence of other app_display threads !!
		if len(displayedData) != 0 {
			err6 := app_display.Display(stencilDBConn, app, notDisplayedData, pks)
			if err6 != nil {
				log.Fatal(err6)
			}
			return returnResultBasedOnNodeCompleteness(err1)
		}

		pTags, err2 := oneMigratedData.GetParentTags(&appConfig)
		if err2 != nil {
			log.Fatal(err2)
		} else {
			if pTags == nil {
				log.Println("This Data's Tag Does not Depend on Any Other Tag!")
				err3 := app_display.Display(stencilDBConn, app, dataInNode, pks)
				if err3 != nil {
					log.Fatal(err3)
				}
				return returnResultBasedOnNodeCompleteness(err1)
			} else {
				pTagConditions := make(map[string]bool)
				for _, pTag := range pTags {
					dataInParentNode, err4 := app_dependency_handler.GetdataFromParentNode(&appConfig, dataInNode, pTag)
					log.Println(dataInParentNode, err4)
					displaySetting, err5 := app_dependency_handler.GetDisplaySettingInDependencies(&appConfig, oneMigratedData, pTag)
					if err5 != nil {
						log.Fatal(err5)
					}
					if err4 != nil {
						switch err4.Error() {
						case "This Data Does not Depend on Any Data in the Parent Node":
							pTagConditions[pTag] = true
						case "Fail To Get Any Data in the Parent Node":
							pTagConditions[pTag] = returnDisplayConditionWhenCannotGetDataFromParentNode(displaySetting, secondRound)
						}
					} else {
						// For now, there is no case where there is more than one piece of data in a parent node
						if len(dataInParentNode) != 1 {
							log.Fatal("Find more than one piece of data in a parent node!!")
						}
						result, err7 := checkDisplayOneMigratedData(stencilDBConn, appConfig, dataInParentNode[0], app, pks, secondRound)
						if err7 != nil {
							log.Println(err7)
						}
						log.Println(result, err7)
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
				log.Println(pTagConditions)
				// For now, without checking the combined_display_setting,
				// this check app_display condition func will return true
				// as long as one pTagCondition is true
				if checkResult := checkDisplayConditions(&appConfig, pTagConditions, oneMigratedData); checkResult {
					err8 := app_display.Display(stencilDBConn, app, dataInNode, pks)
					if err8 != nil {
						log.Fatal(err8)
					}
					return returnResultBasedOnNodeCompleteness(err1)
				} else {
					return "No Data In a Node Can be Displayed", errors.New("Display Setting does not allow the data in the node to be displayed")
				}
			}
		}
	}
	panic("Should never happen here!")
}
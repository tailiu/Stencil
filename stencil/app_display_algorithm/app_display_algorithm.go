package app_display_algorithm

import (
	"database/sql"
	"log"
	"stencil/config"
	"stencil/app_dependency_handler"
	"stencil/app_display"
	"time"
)

const checkInterval = 200 * time.Millisecond

func DisplayThread(app string, migrationID int) {
	startTime := time.Now()
	log.Println("--------- Start of Display Check ---------")

	stencilDBConn, appConfig := app_display.Initialize(app)

	log.Println("--------- First Phase --------")
	secondRound := false
	for migratedData := app_display.GetUndisplayedMigratedData(stencilDBConn, appConfig, migrationID); 
		!app_display.CheckMigrationComplete(stencilDBConn, migrationID); 
		migratedData = app_display.GetUndisplayedMigratedData(stencilDBConn, appConfig, migrationID) {

		for _, oneMigratedData := range migratedData {
			checkDisplayOneMigratedData(stencilDBConn, appConfig, oneMigratedData, secondRound)
		}

		time.Sleep(checkInterval)
	}

	log.Println("--------- Second Phase ---------")
	secondRound = true
	secondRoundMigratedData := app_display.GetUndisplayedMigratedData(stencilDBConn, appConfig, migrationID)
	for _, oneSecondRoundMigratedData := range secondRoundMigratedData {
		checkDisplayOneMigratedData(stencilDBConn, appConfig, oneSecondRoundMigratedData, secondRound)
	}

	log.Println("--------- End of Display Check ---------")
	endTime := time.Now()
	log.Println("Time used: ", endTime.Sub(startTime))
}

func checkDisplayOneMigratedData(stencilDBConn *sql.DB, appConfig *config.AppConfig, oneMigratedData *app_display.HintStruct, secondRound bool) error {

	log.Println("Check Data ", *oneMigratedData)
	dataInNode, err1 := app_dependency_handler.GetDataInNodeBasedOnDisplaySetting(appConfig, oneMigratedData)
	if dataInNode == nil {
		log.Println(err1)
		return app_display.NoNodeCanBeDisplayed
	} else {

		var displayedData, notDisplayedData []*app_display.HintStruct
		for _, oneDataInNode := range dataInNode {
			displayed := app_display.CheckDisplay(stencilDBConn, appConfig, oneDataInNode)
			if !displayed {
				notDisplayedData = append(notDisplayedData, oneDataInNode)
			} else {
				displayedData = append(displayedData, oneDataInNode)
			}
		}
		// Note: This will be changed when considering ongoing application services
		// and the existence of other app_display threads !!
		if len(displayedData) != 0 {
			err6 := app_display.Display(stencilDBConn, appConfig, notDisplayedData)
			if err6 != nil {
				log.Fatal(err6)
			}
			return app_display.ReturnResultBasedOnNodeCompleteness(err1)
		}

		pTags, err2 := oneMigratedData.GetParentTags(appConfig)
		if err2 != nil {
			log.Fatal(err2)
		} else {
			if pTags == nil {
				log.Println("This Data's Tag Does not Depend on Any Other Tag!")
				err3 := app_display.Display(stencilDBConn, appConfig, dataInNode)
				if err3 != nil {
					log.Fatal(err3)
				}
				return app_display.ReturnResultBasedOnNodeCompleteness(err1)
			} else {
				pTagConditions := make(map[string]bool)
				for _, pTag := range pTags {
					dataInParentNode, err4 := app_dependency_handler.GetdataFromParentNode(appConfig, dataInNode, pTag)
					log.Println(dataInParentNode, err4)
					displaySetting, err5 := app_dependency_handler.GetDisplaySettingInDependencies(appConfig, oneMigratedData, pTag)
					if err5 != nil {
						log.Fatal(err5)
					}
					if err4 != nil {
						switch err4 {
						case app_display.NotDependsOnAnyData:
							pTagConditions[pTag] = true
						case app_display.CannotFindAnyDataInParent:
							pTagConditions[pTag] = app_display.ReturnDisplayConditionWhenCannotGetDataFromParentNode(displaySetting, secondRound)
						}
					} else {
						// For now, there is no case where there is more than one piece of data in a parent node
						// if len(dataInParentNode) != 1 {
						// 	log.Fatal("Find more than one piece of data in a parent node!!")
						// }
						err7 := checkDisplayOneMigratedData(stencilDBConn, appConfig, dataInParentNode, secondRound)
						if err7 != nil {
							log.Println(err7)
						}
						switch err7 {
						case app_display.NoNodeCanBeDisplayed:
							pTagConditions[pTag] = app_display.ReturnDisplayConditionWhenCannotGetDataFromParentNode(displaySetting, secondRound)
						case app_display.PartiallyDisplayed:
							pTagConditions[pTag] = app_display.ReturnDisplayConditionWhenGetPartialDataFromParentNode(displaySetting)
						case app_display.CompletelyDisplayed:
							pTagConditions[pTag] = true
						}
					}
				}
				log.Println(pTagConditions)
				// For now, without checking the combined_display_setting,
				// this check app_display condition func will return true
				// as long as one pTagCondition is true
				if checkResult := app_display.CheckDisplayConditions(appConfig, pTagConditions, oneMigratedData); checkResult {
					err8 := app_display.Display(stencilDBConn, appConfig, dataInNode)
					if err8 != nil {
						log.Fatal(err8)
					}
					return app_display.ReturnResultBasedOnNodeCompleteness(err1)
				} else {
					return app_display.NoNodeCanBeDisplayed
				}
			}
		}
	}
	panic("Should never happen here!")
}
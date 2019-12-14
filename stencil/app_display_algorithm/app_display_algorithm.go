package app_display_algorithm

import (
	"log"
	"stencil/config"
	"stencil/app_dependency_handler"
	"stencil/app_display"
	"time"
)

const checkInterval = 200 * time.Millisecond

func DisplayThread(displayConfig *config.DisplayConfig) {

	startTime := time.Now()

	log.Println("--------- Start of Display Check ---------")

	log.Println("--------- First Phase --------")

	secondRound := false

	for migratedData := app_display.GetUndisplayedMigratedData(displayConfig); 
		!app_display.CheckMigrationComplete(displayConfig); 
		migratedData = app_display.GetUndisplayedMigratedData(displayConfig) {

		for _, oneMigratedData := range migratedData {

			checkDisplayOneMigratedData(displayConfig, oneMigratedData, secondRound)

		}

		time.Sleep(checkInterval)
	}


	log.Println("--------- Second Phase ---------")
	
	secondRound = true

	secondRoundMigratedData := app_display.GetUndisplayedMigratedData(displayConfig)
	
	for _, oneSecondRoundMigratedData := range secondRoundMigratedData {

		checkDisplayOneMigratedData(displayConfig, oneSecondRoundMigratedData, secondRound)

	}

	log.Println("--------- End of Display Check ---------")
	
	endTime := time.Now()
	log.Println("Time used: ", endTime.Sub(startTime))

}

func checkDisplayOneMigratedData(
	displayConfig *config.DisplayConfig, 
	oneMigratedData *app_display.HintStruct,
	secondRound bool) error {

	log.Println("Check Data:", *oneMigratedData)

	dataInNode, err1 := app_dependency_handler.GetDataInNodeBasedOnDisplaySetting(
		displayConfig, oneMigratedData)

	// Either this data is not in the destination application,
	// e.g., this data is displayed by other threads and deleted by application services, 
	// or the data is not able to be displayed 
	// because of missing some other data it depends on within the node,
	// e.g., this node only has status_stats without status
	if dataInNode == nil {

		log.Println(err1)
		
		if secondRound {

			err9 := app_display.PutIntoDataBag(displayConfig, []*app_display.HintStruct{oneMigratedData})

			log.Fatal(err9)

			return app_display.NoNodeCanBeDisplayed

		} else {

			return app_display.NoNodeCanBeDisplayed

		}
	} else {

		// This is to display data once there is any data already displayed in a node
		var displayedData, notDisplayedData []*app_display.HintStruct

		for _, oneDataInNode := range dataInNode {

			displayed := app_display.CheckDisplay(displayConfig, oneDataInNode)

			if !displayed {

				notDisplayedData = append(notDisplayedData, oneDataInNode)

			} else {

				displayedData = append(displayedData, oneDataInNode)

			}
		}

		// Note: This will be changed when considering ongoing application services
		// and the existence of other app_display threads !!
		if len(displayedData) != 0 {

			err6 := app_display.Display(displayConfig, notDisplayedData)
			if err6 != nil {
				log.Fatal(err6)
			}

			return app_display.ReturnResultBasedOnNodeCompleteness(err1)

		}

		pTags, err2 := oneMigratedData.GetParentTags(displayConfig)
		if err2 != nil {

			log.Fatal(err2)

		} else {

			if pTags == nil {

				log.Println("This Data's Tag Does not Depend on Any Other Tag!")

				err3 := app_display.Display(displayConfig, dataInNode)
				if err3 != nil {
					log.Fatal(err3)
				}
				
				return app_display.ReturnResultBasedOnNodeCompleteness(err1)

			} else {

				pTagConditions := make(map[string]bool)

				for _, pTag := range pTags {

					dataInParentNode, err4 := app_dependency_handler.GetdataFromParentNode(
						displayConfig, dataInNode, pTag)
					
					if err4 != nil {
						log.Println(err4)
					} else {
						log.Println(*dataInParentNode)
					}

					displaySetting, err5 := app_dependency_handler.GetDisplaySettingInDependencies(
						displayConfig, oneMigratedData, pTag)
					if err5 != nil {
						log.Fatal(err5)
					}

					if err4 != nil {

						switch err4 {

						case app_display.NotDependsOnAnyData:
							
							pTagConditions[pTag] = true

						case app_display.CannotFindAnyDataInParent:

							pTagConditions[pTag] = 
								app_display.ReturnDisplayConditionWhenCannotGetDataFromParentNode(
								displaySetting, secondRound)
							
						}
						
					} else {

						// For now, there is no case where 
						// there is more than one piece of data in a parent node
						// if len(dataInParentNode) != 1 {
						// 	log.Fatal("Find more than one piece of data in a parent node!!")
						// }
						err7 := checkDisplayOneMigratedData(
							displayConfig, dataInParentNode, secondRound)

						if err7 != nil {
							log.Println(err7)
						}

						switch err7 {

						case app_display.NoNodeCanBeDisplayed:

							pTagConditions[pTag] = app_display.
								ReturnDisplayConditionWhenCannotGetDataFromParentNode(
									displaySetting, secondRound)

						case app_display.PartiallyDisplayed:

							pTagConditions[pTag] = app_display.
								ReturnDisplayConditionWhenGetPartialDataFromParentNode(displaySetting)

						case app_display.CompletelyDisplayed:

							pTagConditions[pTag] = true

						}
					}
				}
				log.Println("Get parent nodes results:", pTagConditions)

				// For now, without checking the combined_display_setting,
				// this check app_display condition func will return true
				// as long as one pTagCondition is true
				if checkResult := app_display.CheckDisplayConditions(
					displayConfig, pTagConditions, oneMigratedData); checkResult {
					
					err8 := app_display.Display(displayConfig, dataInNode)
					if err8 != nil {
						log.Fatal(err8)
					}

					return app_display.ReturnResultBasedOnNodeCompleteness(err1)

				} else {

					if secondRound {
						
						err10 := app_display.PutIntoDataBag(displayConfig, dataInNode)

						log.Fatal(err10)
					}

					return app_display.NoNodeCanBeDisplayed

				}
			}
		}
	}

	panic("Should never happen here!")

}
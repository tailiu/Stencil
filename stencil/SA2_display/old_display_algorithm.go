package SA2_display

import (
	"database/sql"
	"errors"
	"log"
	"time"
	"stencil/config"
)

const checkInterval = 200 * time.Millisecond

func oldDisplayThread(migrationID int, deletionHoldEnable bool) {

	startTime := time.Now()

	log.Println("--------- Start of Display Check ---------")

	stencilDBConn, appConfig, threadID, userID := Initialize(migrationID)

	// CreateDeletionHoldTable(stencilDBConn)

	log.Println("Thread ID:", threadID)

	log.Println("--------- First Phase --------")

	secondRound := false

	for migratedData := GetUndisplayedMigratedData(stencilDBConn, migrationID, appConfig); 
		!CheckMigrationComplete(stencilDBConn, migrationID); 
		migratedData = GetUndisplayedMigratedData(stencilDBConn, migrationID, appConfig) {
		
		var dhStack [][]int

		for _, oneMigratedData := range migratedData {
			_, dhStack, _ = oldCheckDisplayOneMigratedData(
				stencilDBConn, appConfig, oneMigratedData, 
				secondRound, deletionHoldEnable, dhStack, threadID, userID)

			if deletionHoldEnable {
				RemoveDeletionHold(stencilDBConn, dhStack, threadID)
			}

		}
		time.Sleep(checkInterval)
	}

	log.Println("--------- Second Phase ---------")

	secondRound = true

	secondRoundMigratedData := GetUndisplayedMigratedData(
		stencilDBConn, migrationID, appConfig)

	for _, oneSecondRoundMigratedData := range secondRoundMigratedData {

		var dhStack [][]int
		
		_, dhStack, _ = checkDisplayOneMigratedData(
			stencilDBConn, appConfig, oneSecondRoundMigratedData, 
			secondRound, deletionHoldEnable, 
			dhStack, threadID, userID,
		)

		if deletionHoldEnable {
			RemoveDeletionHold(stencilDBConn, dhStack, threadID)
		}

	}

	log.Println("--------- End of Display Check ---------")

	logDisplayEndTime(stencilDBConn, migrationID)

	endTime := time.Now()

	log.Println("Time used in this display thread: ", endTime.Sub(startTime))


}

// Two-way display check
func oldCheckDisplayOneMigratedData(stencilDBConn *sql.DB, 
	appConfig *config.AppConfig, oneMigratedData HintStruct, 
	secondRound bool, deletionHoldEnable bool, dhStack [][]int, 
	threadID int, userID string) (string, [][]int, error) {

	// CheckAndGetTableNameAndID(stencilDBConn, &oneMigratedData, appConfig.AppID)
	log.Println("Check Data ", oneMigratedData)
	
	log.Println("==================== Check Intra-node dependencies ====================")

	dataInNode, err1 := GetDataInNodeBasedOnDisplaySetting(appConfig, 
		oneMigratedData, stencilDBConn)

	log.Println("Data in Node:")

	for _, oneDataInNode := range dataInNode {
		log.Println(oneDataInNode)
	}

	// Either this data is not in the destination application,
	// e.g., this data is displayed by other threads and deleted by application services, 
	// or the data is not able to be displayed 
	// because of missing some other data it depends on within the node,
	// e.g., this node only has status_stats without status
	if len(dataInNode) == 0 {

		log.Println(err1)

		if secondRound {

			err10 := PutIntoDataBag(stencilDBConn, 
				appConfig.AppID, []HintStruct{oneMigratedData}, userID)

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

		var displayedData, notDisplayedData []HintStruct
		for _, dataInNode1 := range dataInNode {

			displayed := CheckDisplay(stencilDBConn, appConfig.AppID, dataInNode1)

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

			err6, dhStack = Display(stencilDBConn, 
				appConfig.AppID, notDisplayedData, deletionHoldEnable, dhStack, threadID)

			if err6 != nil {

				return "", dhStack, err6

			}

			return ReturnResultBasedOnNodeCompleteness(err1, dhStack)
		}

		pTags, err2 := oneMigratedData.GetParentTags(appConfig)
		if err2 != nil {
			log.Fatal(err2)
		} else {

			if pTags == nil {

				log.Println("This Data's Tag Does not Depend on Any Other Tag!")

				// Need to change this Display function
				var err3 error
				err3, dhStack = Display(stencilDBConn, 
					appConfig.AppID, dataInNode, deletionHoldEnable, dhStack, threadID)

				// Found path conflicts
				if err3 != nil {
					return "", dhStack, err3
				}

				return ReturnResultBasedOnNodeCompleteness(err1, dhStack)

			} else {
				
				pTagConditions := make(map[string]bool)

				for _, pTag := range pTags {
					log.Println(pTag)
					
					dataInParentNode, err4 := GetdataFromParentNode(
						stencilDBConn, appConfig, dataInNode, pTag)

					log.Println(dataInParentNode, err4)

					displaySetting, err5 := GetDisplaySettingInDependencies(
						appConfig, oneMigratedData, pTag)

					if err5 != nil {
						log.Fatal(err5)
					}
					if err4 != nil {

						switch err4.Error() {

						case "This Data Does not Depend on Any Data in the Parent Node":

							pTagConditions[pTag] = true

						case "Fail To Get Any Data in the Parent Node":

							pTagConditions[pTag] = ReturnDisplayConditionWhenCannotGetDataFromParentNode(
								displaySetting, secondRound)
						}
					} else {

						// For now, there is no case where 
						// there is more than one piece of data in a parent node
						// if len(dataInParentNode) != 1 {
						// 	log.Fatal("Find more than one piece of data in a parent node!!")
						// }
						var result string
						var err7 error

						result, dhStack, err7 = checkDisplayOneMigratedData(
							stencilDBConn, appConfig, dataInParentNode, 
							secondRound, deletionHoldEnable, dhStack, threadID, userID)
						
						if err7 != nil {
							log.Println(err7)

							// If there is path confilct, 
							// just return until to the original layer 
							// and remove the already added deletion hold
							if err7.Error() == "Path conflict" {
								return "", dhStack, err7
							}

						}

						log.Println(result, err7)

						switch result {

						case "No Data In a Node Can be Displayed":
							pTagConditions[pTag] = 
							ReturnDisplayConditionWhenCannotGetDataFromParentNode(
								displaySetting, secondRound)

						case "Data In a Node Can be partially Displayed":
							pTagConditions[pTag] = 
							ReturnDisplayConditionWhenGetPartialDataFromParentNode(
								displaySetting)

						case "Data In a Node Can be completely Displayed":
							pTagConditions[pTag] = true
						}
					}
				}
				// log.Println(pTagConditions)

				if checkResult := CheckCombinedDisplayConditions(
					appConfig, pTagConditions, oneMigratedData); 
					checkResult {

					var err8 error

					err8, dhStack = Display(
						stencilDBConn, appConfig.AppID, 
						dataInNode, deletionHoldEnable, dhStack, threadID)

					// Found path conflicts
					if err8 != nil {
						return "", dhStack, err8
					}
					return ReturnResultBasedOnNodeCompleteness(err1, dhStack)

				} else {

					if secondRound {
						
						err9 := PutIntoDataBag(stencilDBConn, appConfig.AppID, dataInNode, userID)
						
						// Found path conflicts

						if err9 != nil {
							return "No Data In a Node Can be Displayed", 
							dhStack, 
							err9

						} else {
							return "No Data In a Node Can be Displayed", 
							dhStack, 
							errors.New("Display Setting does not allow the data in the node to be displayed")
						}

					} else {
						return "No Data In a Node Can be Displayed", 
						dhStack, 
						errors.New("Display Setting does not allow the data in the node to be displayed")
					}
				}
			}
		}
	}

	panic("Should never happen here!")

}
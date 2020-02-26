package SA2_display

import (
	"database/sql"
	"errors"
	"log"
	"time"
	"stencil/config"
)

const checkInterval = 200 * time.Millisecond

func DisplayThread(migrationID int, deletionHoldEnable bool) {

	startTime := time.Now()

	log.Println("--------- Start of Display Check ---------")

	stencilDBConn, appConfig, threadID, userID, dstDAG := Initialize(migrationID)

	defer stencilDBConn.Close()
	defer appConfig.DBConn.Close()

	// CreateDeletionHoldTable(stencilDBConn)

	log.Println("Thread ID:", threadID)

	log.Println("--------- First Phase --------")

	secondRound := false

	for migratedData := GetUndisplayedMigratedData(stencilDBConn, migrationID, appConfig); 
		!CheckMigrationComplete(stencilDBConn, migrationID); 
		migratedData = GetUndisplayedMigratedData(stencilDBConn, migrationID, appConfig) {
		
		var dhStack [][]int

		for _, oneMigratedData := range migratedData {
			_, dhStack, _ = checkDisplayOneMigratedData(
				stencilDBConn, appConfig, oneMigratedData, 
				secondRound, deletionHoldEnable, dhStack, threadID, userID, dstDAG,
			)

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
			dhStack, threadID, userID, dstDAG,
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
func checkDisplayOneMigratedData(stencilDBConn *sql.DB, 
	appConfig *config.AppConfig, oneMigratedData HintStruct, 
	secondRound bool, deletionHoldEnable bool, dhStack [][]int, 
	threadID int, userID string, dstDAG *DAG) (string, [][]int, error) {

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

			err10 := PutIntoDataBag(stencilDBConn, appConfig.AppID,
				[]HintStruct{oneMigratedData}, userID)

			// Found path conflicts
			if err10 != nil {

				return NoDataInNodeCanBeDisplayed, dhStack, err10

			} else {

				return NoDataInNodeCanBeDisplayed, dhStack, err1

			}
		} else {

			return NoDataInNodeCanBeDisplayed, dhStack, err1

		}
	} else {

		var displayedData, notDisplayedData []HintStruct
		for _, dataInNode1 := range dataInNode {

			displayed := CheckDisplay(stencilDBConn, appConfig.AppID, dataInNode1)

			if !displayed {

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
				appConfig.AppID, notDisplayedData, 
				deletionHoldEnable, dhStack, threadID,
			)

			if err6 != nil {

				return "", dhStack, err6

			}

			return ReturnResultBasedOnNodeCompleteness(err1, dhStack)
		}
		
		log.Println("==================== Check Ownership ====================")

		nodeTag := oneMigratedData.GetTagName(appConfig)

		log.Println("Check data with tag:", nodeTag)
		
		// If the tag of this node is the root, the node could be the migrating user's 
		// or other users' root. Regardless of that, this node will be displayed
		// and there is no need to further check data dependencies since root node does not
		// depend on other nodes
		if nodeTag == "root" {

			log.Println("The checked data is a root node")

			err15 := Display(stencilDBConn, 
				appConfig.AppID, dataInNode, 
				deletionHoldEnable, dhStack, threadID,
			)

			if err15 != nil {
				log.Fatal(err15)
			} else {
				log.Println("Display a root node when checking ownership")
			}
			
			return ReturnResultBasedOnNodeCompleteness(err1)
		
		// If the tag of this node is not the root,
		// we need to check the ownership and sharing relationships of this data.
		// The check of sharing conditions for now is not implemented for now.
		} else {
			
			dataOwnershipSpec, err12 := oneMigratedData.GetOwnershipSpec(dstDAG)
			
			// Mastodon conversations have no ownership settings. In this case
			// we cannot check ownership settings
			if err12 != nil {
				
				log.Println(err12)
				log.Println("Skip this ownership check")
			
			} else {
				// log.Println(dataOwnershipSpec)

				dataInOwnerNode, err13 := getOwner(displayConfig, dataInNode, dataOwnershipSpec)

				// The root node could be incomplete
				if err13 != nil {
					log.Println("An error in getting the checked node's owner:")
					log.Println(err13)
				}

				// Display the data not displayed in the root node
				// this root node should be could be the migrating user's root node
				// or other users' root nodes
				if len(dataInOwnerNode) != 0 {

					displayedDataInOwnerNode, notDisplayedDataInOwnerNode := checkDisplayConditionsInNode(
						displayConfig, dataInOwnerNode)
					
					if len(displayedDataInOwnerNode) != 0 {

						err6 := Display(displayConfig, notDisplayedDataInOwnerNode)
						if err6 != nil {
							log.Fatal(err6)
						}

					}

					var displayedDataInOwnerNode, notDisplayedDataInOwnerNode []HintStruct
					for _, dataInOwnerNode1 := range dataInOwnerNode {

						displayed := CheckDisplay(stencilDBConn, appConfig.AppID, dataInOwnerNode1)

						if !displayed {

							notDisplayedDataInOwnerNode = append(notDisplayedDataInOwnerNode, dataInOwnerNode1)

						} else {

							displayedDataInOwnerNode = append(displayedDataInOwnerNode, dataInOwnerNode1)
						}
					}

					// Note: This will be changed when considering ongoing application services
					// and the existence of other display threads !!
					if len(displayedData) != 0 {

						var err6 error

						err6, dhStack = Display(stencilDBConn, 
							appConfig.AppID, notDisplayedDataInOwnerNode, 
							deletionHoldEnable, dhStack, threadID,
						)

						if err6 != nil {

							return "", dhStack, err6

						}
					}

				}

				// If based on the ownership display settings this node is allowed to be displayed,
				// then continue to check dependencies.
				// Otherwise, no data in the node can be displayed.
				if displayResultBasedOnOwnership := CheckOwnershipCondition(
					dataOwnershipSpec.Display_setting, err13); !displayResultBasedOnOwnership {

					log.Println(`Ownership display settings are not satisfied, 
						so this node cannot be displayed`)

					error16	:= chechPutIntoDataBag(stencilDBConn, appConfig.AppID, 
						dataInNode, userID, secondRound)

					return NoDataInNodeCanBeDisplayed, dhStack, error16

				} else {

					log.Println("Ownership display settings are satisfied")

				}
			}
			
		}

		log.Println("==================== Check Inter-node dependencies ====================")
		
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
							secondRound, deletionHoldEnable, dhStack, threadID, userID, dstDAG,
						)
						
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

						case NoDataInNodeCanBeDisplayed:
							pTagConditions[pTag] = 
								ReturnDisplayConditionWhenCannotGetDataFromParentNode(
									displaySetting, secondRound)

						case PartiallyDisplayed:
							pTagConditions[pTag] = 
								ReturnDisplayConditionWhenGetPartialDataFromParentNode(
									displaySetting)

						case CompletelyDisplayed:
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
							return NoDataInNodeCanBeDisplayed, 
							dhStack, 
							err9

						} else {
							return NoDataInNodeCanBeDisplayed, 
							dhStack, 
							errors.New("Display Setting does not allow the data in the node to be displayed")
						}

					} else {
						return NoDataInNodeCanBeDisplayed, 
						dhStack, 
						errors.New("Display Setting does not allow the data in the node to be displayed")
					}
				}
			}
		}
	}

	panic("Should never happen here!")

}
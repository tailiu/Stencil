package SA1_display

import (
	"log"
	"stencil/config"
	"time"
)

const checkInterval = 200 * time.Millisecond

func DisplayThread(displayConfig *config.DisplayConfig) {

	startTime := time.Now()

	log.Println("--------- Start of Display Check ---------")

	log.Println("--------- First Phase --------")

	secondRound := false

	for migratedData := GetUndisplayedMigratedData(displayConfig); 
		!CheckMigrationComplete(displayConfig); 
		migratedData = GetUndisplayedMigratedData(displayConfig) {

		for _, oneMigratedData := range migratedData {

			checkDisplayOneMigratedData(displayConfig, oneMigratedData, secondRound)

		}

		time.Sleep(checkInterval)
	}


	log.Println("--------- Second Phase ---------")
	
	secondRound = true

	secondRoundMigratedData := GetUndisplayedMigratedData(displayConfig)
	
	for _, oneSecondRoundMigratedData := range secondRoundMigratedData {

		checkDisplayOneMigratedData(displayConfig, oneSecondRoundMigratedData, secondRound)

	}

	log.Println("--------- End of Display Check ---------")
	
	endTime := time.Now()
	log.Println("Time used: ", endTime.Sub(startTime))

}

func checkDisplayOneMigratedData(
	displayConfig *config.DisplayConfig, 
	oneMigratedData *HintStruct,
	secondRound bool) error {

	log.Println("Check Data:", *oneMigratedData)

	// Get data in the node based on intra-node data dependencies
	dataInNode, err1 := GetDataInNodeBasedOnDisplaySetting(
		displayConfig, oneMigratedData)

	// If dataInNode is nil, either this data is not in the destination application,
	// e.g., this data is displayed by other threads and deleted by application services, 
	// or the data is not able to be displayed 
	// because of missing some other data it depends on within the node,
	// e.g., this node only has status_stats without status
	if dataInNode == nil {

		log.Println(err1)
		
		if secondRound {

			err9 := PutIntoDataBag(displayConfig, []*HintStruct{oneMigratedData})
			if err9 != nil {
				log.Fatal(err9)
			}

			return NoNodeCanBeDisplayed

		} else {

			return NoNodeCanBeDisplayed

		}
	} else {

		
		displayedData, notDisplayedData := checkDisplayConditionsInNode(displayConfig, dataInNode)

		// This is to display data once there is any data already displayed in a node
		// Note: This will be changed when considering ongoing application services
		// and the existence of other app_display threads !!
		if len(displayedData) != 0 {

			err6 := Display(displayConfig, notDisplayedData)
			if err6 != nil {
				log.Fatal(err6)
			}

			return ReturnResultBasedOnNodeCompleteness(err1)

		}


		// If the tag of this node is the root, it should be other user's root node since
		// the migrating user root node is connected with the migrated data with ownership
		// In this case, we do not need to further check 
		// the inter-node data dependencies, ownership, or sharing relationships
		// of the current root node.
		// As the display thread only displays this migrating user's data, even if there is
		// some data not displayed in the root node in this case, it will not display it, but
		// just return results
		if oneMigratedData.Tag == "root" {

			return ReturnResultBasedOnNodeCompleteness(err1)
		
		// If the tag of this node is not the root,
		// we need to check the ownership and sharing relationship of this data.
		// The check of sharing conditions for now is not implemented.
		// Since we only migrate users' own data and migration threads migrate the root node
		// at the beginning of migrations, checking data ownership
		// should be always true given the current display settings in the ownership.	
		} else {

			dataOwnership, err12 := oneMigratedData.GetOwnership(displayConfig)
			if err12 != nil {
				log.Fatal(err12)
			}

			dataInOwnerNode, err13 := getOwner(displayConfig, dataInNode, dataOwnership)

			// Display the data not displayed in the root node because this root node
			// is the migrating user's root node
			if len(dataInOwnerNode) != 0 {

				displayedDataInOwnerNode, notDisplayedDataInOwnerNode := checkDisplayConditionsInNode(
					displayConfig, dataInOwnerNode)
				
				if len(notDisplayedDataInOwnerNode) != 0 {

					err6 := Display(displayConfig, notDisplayedDataInOwnerNode)
					if err6 != nil {
						log.Fatal(err6)
					}

				}

			}

			displayResultBasedOnOwnership := CheckOwnershipCondition(
				dataOwnership.Display_setting, err13)
			
			if !displayResultBasedOnOwnership {

				return ReturnResultBasedOnOwnershipCheck(err13)

			}

		}
		
		
		// Start to check inter-node data dependencies if this is required, and 
		// inner-node data dependencies, ownership and sharing relationships are satified
		// Basically, overall display results = 
		// 		intra-node check results AND 
		// 		ownership relationship check results AND 
		// 		sharing relationship check results AND
		// 		inter-node check results 
		pTags, err2 := oneMigratedData.GetParentTags(displayConfig)
		if err2 != nil {

			log.Fatal(err2)

		} else {

			// When pTags is nil, it means that the tag of the data being checked
			// does not depend on any other tag. 
			if pTags == nil {

				log.Println("This Data's Tag Does not Depend on Any Other Tag!")

				err3 := Display(displayConfig, dataInNode)
				if err3 != nil {
					log.Fatal(err3)
				}
				
				return ReturnResultBasedOnNodeCompleteness(err1)

			} else {

				pTagConditions := make(map[string]bool)

				for _, pTag := range pTags {

					dataInParentNode, err4 := GetdataFromParentNode(
						displayConfig, dataInNode, pTag)
					
					if err4 != nil {
						log.Println(err4)
					} else {
						log.Println(*dataInParentNode)
					}

					displaySettingInDeps, err5 := oneMigratedData.GetDisplaySettingInDependencies(
						displayConfig, pTag)

					if err5 != nil {
						log.Fatal(err5)
					}

					if err4 != nil {

						switch err4 {

							case NotDependsOnAnyData:

								pTagConditions[pTag] = true

							case CannotFindAnyDataInParent:

								pTagConditions[pTag] = 
									ReturnDisplayConditionWhenCannotGetDataFromParentNode(
										displaySettingInDeps, secondRound)
							
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

							case NoNodeCanBeDisplayed:

								pTagConditions[pTag] = 
									ReturnDisplayConditionWhenCannotGetDataFromParentNode(
										displaySettingInDeps, secondRound)

							case PartiallyDisplayed:

								pTagConditions[pTag] = 
									ReturnDisplayConditionWhenGetPartialDataFromParentNode(
										displaySettingInDeps)

							case CompletelyDisplayed:

								pTagConditions[pTag] = true

						}
					}
				}
				log.Println("Get parent nodes results:", pTagConditions)

				// Check the combined_display_setting from all parent nodes
				// to decide whether to display the current node
				if checkResult := CheckCombinedDisplayConditions(
					displayConfig, pTagConditions, oneMigratedData); checkResult {
					
					err8 := Display(displayConfig, dataInNode)
					if err8 != nil {
						log.Fatal(err8)
					}

					return ReturnResultBasedOnNodeCompleteness(err1)

				} else {

					if secondRound {
						
						err10 := PutIntoDataBag(displayConfig, dataInNode)
						if err10 != nil {
							log.Fatal(err10)
						}
						
					}

					return NoNodeCanBeDisplayed

				}
			}
		}
	}

	panic("Should never happen here!")

}
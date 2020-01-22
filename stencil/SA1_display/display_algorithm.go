package SA1_display

import (
	"log"
	"time"
)

const checkInterval = 200 * time.Millisecond

func DisplayThread(displayConfig *displayConfig) {

	if displayConfig.wg != nil {
		defer displayConfig.wg.Done()
	}

	startTime := time.Now()

	log.Println("--------- Start of Display Check In One Thread ---------")

	log.Println("--------- First Phase --------")

	secondRound := false

	// Since we get all undisplayed data seen so far, during check, there could be cases
	// where some data has already been displayed either by the current display thread or other threads 
	// (this is difficult to know).
	// For that data, we will continue to check it. This does not violate correctness
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

	log.Println("--------- End of Display Check In One Thread ---------")
	
	endTime := time.Now()
	log.Println("Time used in this display thread: ", endTime.Sub(startTime))

}


// Basically, overall display results = 
// 		intra-node check results AND 
// 		ownership relationship check results AND 
// 		sharing relationship check results AND
// 		inter-node check results 
// Display can DISPLAY the migrating user's and other users' data 
// (since data is connected and other users' data also could be checked),
// but can ONLY put the migrating user's data into data bags.
func checkDisplayOneMigratedData(displayConfig *displayConfig, 
	oneMigratedData *HintStruct, secondRound bool) error {
	
	log.Println("Check Data:", *oneMigratedData)

	log.Println("==================== Check Intra-node dependencies ====================")

	// Get data in the node based on intra-node data dependencies
	dataInNode, err1 := GetDataInNodeBasedOnDisplaySetting(
		displayConfig, oneMigratedData)

	log.Println("Data in Node:")

	for _, oneDataInNode := range dataInNode {
		log.Println(*oneDataInNode)
	}

	// If dataInNode is nil, either this data is not in the destination application,
	// e.g., this data is displayed by other threads and deleted by application services, 
	// or the data is not able to be displayed 
	// because of missing some other data it depends on within the node,
	// e.g., this node only has status_stats without status
	if dataInNode == nil {

		log.Println(err1)
		
		return chechPutIntoDataBag(displayConfig, 
			secondRound, []*HintStruct{oneMigratedData})

	} else {

		displayedData, notDisplayedData := checkDisplayConditionsInNode(
			displayConfig, dataInNode)

		// This is to display data once there is any data already displayed in a node
		// Note: This will be changed when considering ongoing application services
		// and the existence of other app_display threads !!
		if len(displayedData) != 0 {

			log.Println("There is already some displayed data in the node")

			err6 := Display(displayConfig, notDisplayedData)
			if err6 != nil {
				log.Fatal(err6)
			}

			return ReturnResultBasedOnNodeCompleteness(err1)

		}

		log.Println("==================== Check Ownership ====================")

		// If the tag of this node is the root, the node could be the migrating user's 
		// or other users' root. Regardless of that, this node will be displayed
		// and there is no need to further check data dependencies since root node does not
		// depend on other nodes
		if oneMigratedData.Tag == "root" {

			log.Println("The checked data is a root node")

			err15 := Display(displayConfig, dataInNode)
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

			dataOwnershipSpec, err12 := oneMigratedData.GetOwnershipSpec(displayConfig)
			if err12 != nil {
				log.Fatal(err12)
			}

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

			}

			// If based on the ownership display settings this node is allowed to be displayed,
			// then continue to check dependencies.
			// Otherwise, no data in the node can be displayed.
			if displayResultBasedOnOwnership := CheckOwnershipCondition(
				dataOwnershipSpec.Display_setting, err13); 
				!displayResultBasedOnOwnership {

				log.Println(`Ownership display settings are not satisfied, 
					so this node cannot be displayed`)

				return chechPutIntoDataBag(displayConfig, 
					secondRound, dataInNode)

			} else {

				log.Println("Ownership display settings are satisfied")

			}
		}
		
		log.Println("==================== Check Inter-node dependencies ====================")

		// After intra-node data dependencies, and ownership and sharing relationships are satified,
		// start to check inter-node data dependencies if this is required.
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
					
					log.Println("Check a Parent Tag:", pTag)

					dataInParentNode, err4 := GetdataFromParentNode(
						displayConfig, dataInNode, pTag)
					
					// There could be cases where the display thread cannot get the data
					// For example, follows require both migrating user's root node (ownership)
					// and also the user followed who might be not in the dest app
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

					return chechPutIntoDataBag(displayConfig, 
						secondRound, dataInNode)

				}
			}
		}
	}

	panic("Should never happen here!")

}
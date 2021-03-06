package SA2_display

import (
	"stencil/common_funcs"
	"log"
	"time"
)

const CHECK_INTERVAL = 200 * time.Millisecond

func (display *display) DisplayThread() {

	defer display.closeDBConns()

	startTime := time.Now()

	log.Println("--------- Start of Display Check In One Thread ---------")

	secondRound := false

	if display.displayInFirstPhase {

		log.Println("--------- First Phase --------")

		for migratedData := display.GetUndisplayedMigratedData(); 
			!display.CheckMigrationComplete(); 
			migratedData = display.GetUndisplayedMigratedData() {

			for _, oneMigratedData := range migratedData {
				display.checkDisplayOneMigratedData(oneMigratedData, secondRound)
			}

			time.Sleep(CHECK_INTERVAL)
		}
	}

	log.Println("--------- Second Phase ---------")

	secondRound = true

	secondRoundMigratedData := display.GetUndisplayedMigratedData()

	for _, oneSecondRoundMigratedData := range secondRoundMigratedData {
		display.checkDisplayOneMigratedData(oneSecondRoundMigratedData, secondRound)
	}

	log.Println("--------- End of Display Check ---------")

	display.logDisplayEndTime()

	endTime := time.Now()

	log.Println("Time used in this display thread: ", endTime.Sub(startTime))

}

// Two-way display check
// func checkDisplayOneMigratedData(stencilDBConn *sql.DB, 
// 	appConfig *config.AppConfig, oneMigratedData HintStruct, 
// 	secondRound bool, deletionHoldEnable bool, dhStack [][]int, 
// 	threadID int, userID string, dstDAG *DAG) (string, [][]int, error) {

func (display *display) checkDisplayOneMigratedData(oneMigratedData *HintStruct, secondRound bool) error {

	log.Println("==================== Check Data ====================")

	// CheckAndGetTableNameAndID(stencilDBConn, &oneMigratedData, appConfig.AppID)
	log.Println(*oneMigratedData)
	
	log.Println("==================== Check Intra-node dependencies ====================")

	dataInNode, err1 := display.GetDataInNodeBasedOnDisplaySetting(oneMigratedData)

	log.Println("Data in Node:")

	for _, oneDataInNode := range dataInNode {
		log.Println(*oneDataInNode)
	}

	// Either this data is not in the destination application,
	// e.g., this data is displayed by other threads and deleted by application services, 
	// or the data is not able to be displayed 
	// because of missing some other data it depends on within the node,
	// e.g., this node only has status_stats without status
	if len(dataInNode) == 0 {

		log.Println(err1)

		return display.chechPutIntoDataBag(secondRound, []*HintStruct{oneMigratedData})

	} else {

		displayedData, notDisplayedData := display.checkDisplayConditionsInNode(dataInNode)

		// Note: This will be changed when considering ongoing application services
		// and the existence of other display threads !!
		if len(displayedData) != 0 {

			log.Println("There is already some displayed data in the node")

			err6 := display.Display(notDisplayedData)
			if err6 != nil {
				log.Println(err6)
			}

			return common_funcs.ReturnResultBasedOnNodeCompleteness(err1)
		}
		
		log.Println("==================== Check Ownership ====================")

		nodeTag := oneMigratedData.Tag

		log.Println("Check data with tag:", nodeTag)
		
		// If the tag of this node is the root, the node could be the migrating user's 
		// or other users' root. Regardless of that, this node will be displayed
		// and there is no need to further check data dependencies since root node does not
		// depend on other nodes
		if nodeTag == "root" {

			log.Println("The checked data is a root node")

			err15 := display.Display(dataInNode)
			if err15 != nil {
				log.Println(err15)
			} else {
				log.Println("Display a root node when checking ownership")
			}
			
			return common_funcs.ReturnResultBasedOnNodeCompleteness(err1)
		
		// If the tag of this node is not the root,
		// we need to check the ownership and sharing relationships of this data.
		// The check of sharing conditions for now is not implemented for now.
		} else {
			
			dataOwnershipSpec, err12 := oneMigratedData.GetOwnershipSpec(display)
			
			// Mastodon conversations have no ownership settings. In this case
			// we cannot check ownership settings
			if err12 != nil {
				
				log.Println(err12)
				log.Println("Skip this ownership check")
			
			} else {
				// log.Println(dataOwnershipSpec)

				dataInOwnerNode, err13 := display.getOwner(dataInNode, dataOwnershipSpec)

				// The root node could be incomplete
				if err13 != nil {
					log.Println("An error in getting the checked node's owner:")
					log.Println(err13)
				}

				// Display the data not displayed in the root node
				// this root node should be could be the migrating user's root node
				// or other users' root nodes
				if len(dataInOwnerNode) != 0 {

					displayedDataInOwnerNode, notDisplayedDataInOwnerNode := display.checkDisplayConditionsInNode(dataInOwnerNode)
					
					if len(displayedDataInOwnerNode) != 0 {

						err6 := display.Display(notDisplayedDataInOwnerNode)
						if err6 != nil {
							log.Fatal(err6)
						}

					}

				}

				// If based on the ownership display settings this node is allowed to be displayed,
				// then continue to check dependencies.
				// Otherwise, no data in the node can be displayed.
				if displayResultBasedOnOwnership := common_funcs.CheckOwnershipCondition(
					dataOwnershipSpec.Display_setting, err13); 
					!displayResultBasedOnOwnership {

					log.Println(`Ownership display settings are not satisfied, 
						so this node cannot be displayed`)

					return display.chechPutIntoDataBag(secondRound, dataInNode)			

				} else {

					log.Println("Ownership display settings are satisfied")

				}
			}
			
		}

		log.Println("==================== Check Inter-node dependencies ====================")
		
		pTags, err2 := oneMigratedData.GetParentTags(display)
		if err2 != nil {
			log.Fatal(err2)
		} else {

			if pTags == nil {

				log.Println("This Data's Tag Does not Depend on Any Other Tag!")

				err3 := display.Display(dataInNode)
				if err3 != nil {
					log.Println(err3)
				}

				return common_funcs.ReturnResultBasedOnNodeCompleteness(err1)

			} else {
				
				pTagConditions := make(map[string]bool)

				for _, pTag := range pTags {

					log.Println("Check a Parent Tag:", pTag)
					
					dataInParentNode, err4 := display.GetdataFromParentNode(dataInNode, pTag)

					if err4 != nil {
						log.Println(err4)
					} else {
						log.Println(*dataInParentNode)
					}

					displaySettingInDeps, err5 := oneMigratedData.GetDisplaySettingInDependencies(
						display, pTag)

					if err5 != nil {
						log.Fatal(err5)
					}

					if err4 != nil {

						switch err4 {

						case common_funcs.NotDependsOnAnyData:

							pTagConditions[pTag] = true

						case common_funcs.CannotFindAnyDataInParent:

							pTagConditions[pTag] = 
								common_funcs.ReturnDisplayConditionWhenCannotGetDataFromParentNode(
								displaySettingInDeps, secondRound,
							)
						}
					} else {

						// For now, there is no case where 
						// there is more than one piece of data in a parent node
						// if len(dataInParentNode) != 1 {
						// 	log.Fatal("Find more than one piece of data in a parent node!!")
						// }
						
						err7 := display.checkDisplayOneMigratedData(dataInParentNode, secondRound)
						
						if err7 != nil {
							log.Println(err7)
						}

						switch err7 {

						case common_funcs.NoDataInNodeCanBeDisplayed:

							pTagConditions[pTag] = 
								common_funcs.ReturnDisplayConditionWhenCannotGetDataFromParentNode(
									displaySettingInDeps, secondRound)

						case common_funcs.PartiallyDisplayed:

							pTagConditions[pTag] = 
								common_funcs.ReturnDisplayConditionWhenGetPartialDataFromParentNode(
									displaySettingInDeps)

						case common_funcs.CompletelyDisplayed:

							pTagConditions[pTag] = true
						}
					}
				}
				// log.Println(pTagConditions)

				if checkResult := display.CheckCombinedDisplayConditions(
					pTagConditions, oneMigratedData); checkResult {

					err8 := display.Display(dataInNode)
					// Found path conflicts
					if err8 != nil {
						log.Println(err8)
					}

					return common_funcs.ReturnResultBasedOnNodeCompleteness(err1)

				} else {

					return display.chechPutIntoDataBag(secondRound, dataInNode)
				}
			}
		}
	}

	panic("Should never happen here!")

}
package SA2_display

import (
	"log"
	"stencil/config"
	"strings"
	"strconv"
)

func ReturnDisplayConditionWhenCannotGetDataFromParentNode(displaySetting string, 
	secondRound bool) bool {
	
	if !secondRound {

		if displaySetting == "parent_node_not_displays_without_check" {

			return true

		} else {

			return false

		}

	} else {

		if displaySetting == "parent_node_not_displays_with_check" ||
			displaySetting == "parent_node_not_displays_without_check" {
			
			return true

		} else {

			return false
		}
	}
}

func ReturnDisplayConditionWhenGetPartialDataFromParentNode(displaySetting string) bool {

	if displaySetting != "parent_node_complete_displays" {
		
		return true

	} else {

		return false
	}
}

func CheckCombinedDisplayConditions(appConfig *config.AppConfig, 
	pTagConditions map[string]bool, oneMigratedData HintStruct) bool {	
	
	if len(pTagConditions) == 1 {

		for _, result := range pTagConditions{
			return result
		}

	}

	combinedSettings, err := oneMigratedData.GetCombinedDisplaySettings(appConfig)
	if err != nil {
		log.Fatal(err)
	}
	
	for key, val := range pTagConditions {
		combinedSettings = strings.Replace(combinedSettings, key, strconv.FormatBool(val), 1)
	}
	strs := strings.Split(combinedSettings, " ")

	var combinedResults bool 
	var operator string

	for i, val := range strs {

		if i == 0 {

			result, err := strconv.ParseBool(val)
			if err != nil {
				log.Fatal(err)
			}
			
			combinedResults = result
		
		} else if i % 2 == 1 {
			
			operator = val
		
		} else {

			result, err := strconv.ParseBool(val)
			if err != nil {
				log.Fatal(err)
			}

			if operator == "or" {

				combinedResults = combinedResults || result

			} else {
				
				combinedResults = combinedResults && result
			}
		}
	}

	return combinedResults

}

func CheckOwnershipCondition(displaySettingInOwnership string, err error) bool {

	// Currently there are only two display settings in the ownership node:
	// 1. Not setting ownership means that by default only display data when the ownership node is complete
	// 2. parent_node_partially_displays means that data can be displayed 
	//		when the ownership node is partially displayed.
	// err is nil, meaning that the ownership node is complete
	if (displaySettingInOwnership == "" && err == nil) || 
		(displaySettingInOwnership == "parent_node_partially_displays" && err == NodeIncomplete) {

		return true

	} else {

		return false
	}

}
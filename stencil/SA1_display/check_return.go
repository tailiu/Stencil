package SA1_display

import (
	"log"
	"strconv"
	"strings"
)

func ReturnResultBasedOnNodeCompleteness(err error) error {

	if err != nil {

		return PartiallyDisplayed
	
	} else {

		return CompletelyDisplayed
	}
}

func ReturnDisplayConditionWhenCannotGetDataFromParentNode(
	displaySetting string, secondRound bool) bool {
		
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

func CheckDisplayConditions(displayConfig *displayConfig, 
	pTagConditions map[string]bool, oneMigratedData *HintStruct) bool {
		
	for _, result := range pTagConditions {

		if result {

			return true
		}
	}

	return false
}

func CheckCombinedDisplayConditions(displayConfig *displayConfig, 
	pTagConditions map[string]bool, oneMigratedData *HintStruct) bool {	
	
	if len(pTagConditions) == 1 {

		for _, result := range pTagConditions{
			return result
		}

	}

	combinedSettings, err := oneMigratedData.GetCombinedDisplaySettings(displayConfig)
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

	if (displaySettingInOwnership == "" && err == nil) || 
		(displaySettingInOwnership == "parent_node_partially_displays" && err == NodeIncomplete) {

		return true

	} else {

		return false
	}

}

func ReturnResultBasedOnOwnershipCheck(err error) error {

	if err == NodeIncomplete {
		return DataNotDisplayedDueToIncompleteOwnerNode
	} else if err == DataNotExists {
		return DataNotDisplayedDueToNoDataInOwnerNode
	}

	panic("Should not happen given the current display settings in ownership!")

}
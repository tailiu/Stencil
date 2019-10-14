package app_display

import (
	"stencil/config"
)

func ReturnResultBasedOnNodeCompleteness(err error) error {
	if err != nil {
		return PartiallyDisplayed
	} else {
		return CompletelyDisplayed
	}
}

func ReturnDisplayConditionWhenCannotGetDataFromParentNode(displaySetting string, secondRound bool) bool {
	if !secondRound {
		if displaySetting == "parent_node_not_displays_without_check" {
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

func ReturnDisplayConditionWhenGetPartialDataFromParentNode(displaySetting string) bool {
	if displaySetting != "parent_node_complete_displays" {
		return true
	} else {
		return false
	}
}

func CheckDisplayConditions(appConfig *config.AppConfig, pTagConditions map[string]bool, oneMigratedData *HintStruct) bool {
	for _, result := range pTagConditions {
		if result {
			return true
		}
	}
	return false
}
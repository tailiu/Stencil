package common_funcs

// import (
// 	"log"
// 	"strconv"
// 	"strings"
// )

/**
 *
 * "parent_node_not_displays_without_check": this condition is satisfied
 * 		in any case no matter of the phase or whether the parent node is displayed or not.
 *
 * "parent_node_not_displays_with_check": in the first phase, if parent nodes are not displayed,
 * 		then this condition is not satisfied. In the second phase, this condition is always satisfied,
 *		no matter whether parent nodes are displayed or not.
 *
 *
 * "parent_node_complete_displays": only when the parent node is complete, this condition is satisfied
 *
 * "parent_node_partially_displays": only when the parent node is partially/completely displayed, 
 *		this condition is satisfied
 *
 */

func ReturnResultBasedOnNodeCompleteness(err error) error {

	if err != nil {

		return PartiallyDisplayed
	
	} else {

		return CompletelyDisplayed
	}
}

func CheckOwnershipCondition(displaySettingInOwnership string, err error) bool {

	// Currently there are only two display settings in the ownership node:
	// 1. Not setting ownership means that by default only display data when the ownership node is complete
	// 2. parent_node_partially_displays means that data can be displayed 
	//		when the ownership node is partially displayed.
	// err is nil, meaning that the ownership node is complete
	if (displaySettingInOwnership == "" && err == nil) || 
		(displaySettingInOwnership == "parent_node_partially_displays" &&
		 err == NodeIncomplete) {

		return true

	} else {

		return false
	}

}

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
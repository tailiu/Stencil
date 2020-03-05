package common_funcs

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
		(displaySettingInOwnership == "parent_node_partially_displays" 
			&& err == NodeIncomplete) {

		return true

	} else {

		return false
	}

}
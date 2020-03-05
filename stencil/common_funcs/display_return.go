package common_funcs

func ReturnResultBasedOnNodeCompleteness(err error) error {

	if err != nil {

		return PartiallyDisplayed
	
	} else {

		return CompletelyDisplayed
	}
}
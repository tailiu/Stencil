package SA2_display

import (
	"errors"
)

// Display algorithm
var NoDataInNodeCanBeDisplayed = 
	errors.New("No Data In a Node Can be Displayed")

var PartiallyDisplayed = 
	errors.New("Data In a Node Can be partially Displayed")

var CompletelyDisplayed = 
	errors.New("Data In a Node Can be completely Displayed")

var PathConflictsWhenPuttingInBags =
	errors.New("Found that there is a path conflict!! When putting data in a databag")

var PathConflictsWhenDisplayingData =
	errors.New("Found that there is a path conflict!! When displaying data")

var CannotFindAnyDataInParent = errors.New("Fail To Get Any Data in the Parent Node")
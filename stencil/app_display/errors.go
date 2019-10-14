package app_display

import (
	"errors"
)
// Display algorithm
var NoNodeCanBeDisplayed = errors.New("No Data In a Node Can be Displayed")
var PartiallyDisplayed = errors.New("Data In a Node Can be partially Displayed")
var CompletelyDisplayed = errors.New("Data In a Node Can be completely Displayed")

// Get data in a node
var CannotFindRemainingData = errors.New("Error: Cannot Find One Remaining Data in the Node")
var NodeIncomplete = errors.New("Error: node is not complete")
var DataNotExists = errors.New("Error: the Data in a Data Hint Does Not Exist")

// Get data in a parent node
var CannotFindAnyDataInParent = errors.New("Fail To Get Any Data in the Parent Node")
var NotDependsOnAnyData = errors.New("This Data Does not Depend on Any Data in the Parent Node")
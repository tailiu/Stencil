package common_funcs

import (
	"errors"
)

// DAG specification
var CannotFindRootMemberAttr = 
	errors.New("Cannot find the root member and attribute")

var CannotFindDependencyDisplaySetting = 
	errors.New("No dependency display setting is found!")

var NoTableFound = errors.New("Error: No Table Found For the Provided Member ID")

// Display algorithm
var NoDataInNodeCanBeDisplayed = 
	errors.New("No Data In a Node Can be Displayed")

var PartiallyDisplayed = 
	errors.New("Data In a Node Can be partially Displayed")

var CompletelyDisplayed = 
	errors.New("Data In a Node Can be completely Displayed")

var NodeIncomplete = 
	errors.New("Error: node is not complete")

// Get data in a node
var CannotFindRemainingData = 
	errors.New("Error: Cannot Find One Remaining Data in the Node")

var DataNotExists = 
	errors.New("Error: the Data in a Data Hint Does Not Exist")

// Get data in a parent node
var CannotFindAnyDataInParent = 
	errors.New("Fail To Get Any Data in the Parent Node")

var NotDependsOnAnyData = 
	errors.New("This Data Does not Depend on Any Data in the Parent Node")

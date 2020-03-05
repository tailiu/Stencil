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
var NoNodeCanBeDisplayed = 
	errors.New("No Data In a Node Can be Displayed")

var PartiallyDisplayed = 
	errors.New("Data In a Node Can be partially Displayed")

var CompletelyDisplayed = 
	errors.New("Data In a Node Can be completely Displayed")

var NodeIncomplete = 
	errors.New("Error: node is not complete")

package SA1_display

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

// Resolve reference when getting data in node
var CannotResolveReferencesGetDataInNode = 
	errors.New("Fail to resolve references when remaining data in node")

var NoMappingAndNoReferenceToResolve =
	errors.New("There is no need to resolve reference and there is no mapping")

// Ownership
var CannotFindDataInOwnership = errors.New("Fail to get any Data by the ownership relationship")

var DataNotDisplayedDueToIncompleteOwnerNode = 
	errors.New("Data is not displayed because the ownership node is not complete")

var DataNotDisplayedDueToNoDataInOwnerNode = 
	errors.New("Data is not displayed because no data can be displayed in the ownership node")
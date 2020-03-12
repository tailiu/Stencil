package SA1_display

import (
	"errors"
)

// Resolve reference when getting data in node or parent node
var CannotResolveReferencesGetDataInParentNode =
	errors.New("Fail to resolve references when getting data in a parent node")

var CannotResolveReferencesGetDataInNode = 
	errors.New("Fail to resolve references when getting remaining data in a node")

var NoReferenceToResolve =
	errors.New("There is no need to resolve reference")

var CannotFindResolvedAttributes = errors.New(`Does not find resolved attributes. 
	Should not happen given one member corresponds to one row for now!`)

var CannotGetPrevID = errors.New(`Cannot get previous ids because of the row has 
	not been inserted into the identity table`)

var CannotGetDataAfterResolvingRef2 = errors.New(`Cannot get data after resolving reference2`)

// Ownership
var CannotFindDataInOwnership = 
	errors.New("Fail to get any Data by the ownership relationship")

// Ownership in old functions
var NotMigratingUserRootNode = errors.New("Not migrating user root node")

var CannotFindRootTable = errors.New("Cannot find root table")


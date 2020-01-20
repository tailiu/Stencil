package schema_mappings

import (
	"errors"
)

var NoMappedAttrFound = errors.New("No mapped attribute found")

var NoREFFound = errors.New("The mapped attribute does not contain #REF")

var MappingsToAppNotFound = errors.New("No mappings to this app are found")

var CannotOpenPSMFile = errors.New("Can't open the pairwise schema mapping json file")

var CannotFindFromApp = errors.New("Cannot find the from app in the mappings")

var CannotFindToApp = errors.New("Cannot find the to app in the mappings")

var CannotCreateToApp = errors.New("Cannot create a to app in the mappings")

var CannotGetMappingsByFromTable = errors.New("Cannot get mappings by the from table")

var CannotGetToTableByName = errors.New("Cannot get to table by the table name")

var duplicateToAttrWithThirdArg = errors.New("Found a duplicate toAttr also with third Arg")

var duplicateToAttrWithoutThirdArg = errors.New("Found a duplicate toAttr also without third Arg")
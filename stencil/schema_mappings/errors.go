package schema_mappings

import (
	"errors"
)

var NoMappedAttrFound = errors.New("No mapped attribute found")

var MappingsToAppNotFound = errors.New("No mappings to this app are found")

var CannotOpenPSMFile = errors.New("Can't open the pairwise schema mapping json file")

var CannotFindPairwiseMappings = errors.New("Cannot find the pairwise mappings")

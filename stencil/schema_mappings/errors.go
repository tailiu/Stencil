package schema_mappings

import (
	"errors"
)

var NoMappedAttrFound = errors.New("No mapped attribute found")

var MappingsToAppNotFound = errors.New("No mappings to this app are found")

var CannotOpenPSMFile = errors.New("can't open pairwise schema mapping json file")

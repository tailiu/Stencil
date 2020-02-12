package reference_resolution

import (
	"errors"
)

var dataToUpdateOtherDataNotFound = errors.New("No data found for updating other data")

var notMigrated = errors.New("Data not migrated maybe due to the lack of schema mappings")

var alreadySolved = errors.New("Reference has probably been implicitly resolved")

var notOneAttributeFound = errors.New("The number of attribute to update other attributes is not one")

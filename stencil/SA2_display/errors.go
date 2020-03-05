package SA2_display

import (
	"errors"
)

// Display algorithm
var PathConflictsWhenPuttingInBags =
	errors.New("Found that there is a path conflict!! When putting data in a databag")

var PathConflictsWhenDisplayingData =
	errors.New("Found that there is a path conflict!! When displaying data")


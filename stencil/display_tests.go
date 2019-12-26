package main

import (
	"stencil/SA1_display"
)

func main() {

	// threadNum := 1
	
	migrationID := 955012936

	// If the destination app database is not in the new server, newDB is false
	newDB := false

	// If the display controller needs to resolve references, resolveReference is true
	resolveReference := true

	SA1_display.CreateDisplayConfig(migrationID, resolveReference, newDB)

	
}
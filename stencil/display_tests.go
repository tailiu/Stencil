package main

import (
	"stencil/SA1_display"
	"stencil/reference_resolution"
	"log"
)

func test1() {
	
	migrationID := 955012936

	// If the destination app database is not in the new server, newDB is false
	newDB := false

	// If the display controller needs to resolve references, resolveReference is true
	resolveReference := true

	SA1_display.CreateDisplayConfig(migrationID, resolveReference, newDB)

}

func test2() {

	prevUserIDs := reference_resolution.GetPrevUserIDs("2", "13")

	// preUserIDs := reference_resolution.GetPreUserIDs("1", "44778")

	log.Println(prevUserIDs)
}

func main() {

	// test1()

	test2()
	
}
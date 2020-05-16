package main

import (
	"data_generators/data_generator"
)

func main() {

	app := "diaspora"
	db := "diaspora_1000_0"

	dataGen := data_generator.Initialize(db)

	switch app {
	case data_generator.DIASPORA:
		dataGen.GenDataDiaspora()
	case data_generator.MASTODON:
		dataGen.GenDataMastodon()
	}
	
}
package main

import (
	"data_generators/data_generator"
)

func main() {

	app := "twitter"
	db := "twitter_1000"

	dataGen := data_generator.Initialize(db)

	switch app {
	case data_generator.DIASPORA:
		dataGen.GenDataDiaspora()
	case data_generator.MASTODON:
		dataGen.GenDataMastodon()
	case data_generator.TWITTER:
		dataGen.GenDataTwitter()
	default:
        panic("Unrecognized application!")
	}
	
}
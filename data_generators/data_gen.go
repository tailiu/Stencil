package main

import (
	"data_generators/data_generator"
)

func main() {

	app := "mastodon"
	db := "mastodon_10000"

	dataGen := data_generator.Initialize(db)

	switch app {
	case data_generator.DIASPORA:
		dataGen.DiasporaGenData()
	case data_generator.MASTODON:
		dataGen.MastodonGenData()
	case data_generator.TWITTER:
		dataGen.TwitterGenData()
	case data_generator.GNUSOCIAL:
		dataGen.GnusocialGenData()
	default:
        panic("Unrecognized application!")
	}
	
}
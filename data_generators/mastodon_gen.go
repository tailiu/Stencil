package main

import (
	"data_generators/data_generator"
)

func main() {

	app := "mastodon_1000"

	dataGen := data_generator.Initialize(app)

	dataGen.genData()

}
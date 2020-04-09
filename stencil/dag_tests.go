package main

import (
	"stencil/common_funcs"
	"log"
)

func main() {
	
	app := "diaspora"

	dag, err := common_funcs.LoadDAG(app)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(dag.IfDependsOn("notification_actors", "notification_id"))

}
package main

import (
	"stencil/common_funcs"
	"log"
)

func test1() {

	app := "diaspora"

	dag, err := common_funcs.LoadDAG(app)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(dag.IfDependsOnBasedOnDag("notification_actors", "notification_id"))

}

func test2() {

	app := "diaspora"

	dag, err := common_funcs.LoadDAG(app)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(dag.GetAllAttrsDepsOnBasedOnDag("photos"))

}

func test3() {

	app := "diaspora"

	dag, err := common_funcs.LoadDAG(app)
	if err != nil {
		log.Fatal(err)
	}

	// attr, attrTable, attrToUpdate, attrToUpdateTable := "id", "people", "author_id", "posts"
	attr, attrTable, attrToUpdate, attrToUpdateTable := "id", "users", "user_id", "aspects"

	log.Println(dag.ReferenceExists(attr, attrTable, attrToUpdate, attrToUpdateTable))

}

func main() {
	
	// test1()

	// test2()

	test3()

}
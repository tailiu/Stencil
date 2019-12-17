package main

import (
	"diaspora/data_generator"
	"log"
)

func main() {
	
	alpha := data_generator.ALPHA

	xm := data_generator.XM

	numOfUsers := 100000

	totalData := 802920 

	scores := data_generator.ParetoScores(alpha, xm, numOfUsers)

	dis := data_generator.AssignDataToUsersByUserScores(scores, totalData)

	log.Println(dis)

}
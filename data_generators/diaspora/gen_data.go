package main

import (
	"diaspora/datagen"
	// "diaspora/db"
	"diaspora/data_generator"
	"log"
)

const APP = "diaspora" 
const USERNUM = 1000

func genUsers(genConfig *data_generator.GenConfig) {
	for i := 0; i < USERNUM; i++ {
		uid, _, _ := datagen.NewUser(genConfig.DBConn)
		log.Println(uid)
	}
}

func genPosts() {

}

func genComments() {

}

func main() {
	genConfig := data_generator.Initialize(APP, USERNUM)
	// log.Println(genConfig.LikeScores)
	// log.Println(genConfig.CommentScores)
	genUsers(genConfig)
}
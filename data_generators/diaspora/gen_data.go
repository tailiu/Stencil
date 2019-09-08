package main

import (
	"diaspora/datagen"
	"diaspora/data_generator"
	// "log"
)

const APP = "diaspora" 
const USERNUM = 1000
const POSTNUM = 8030

func genUsers(genConfig *data_generator.GenConfig) []data_generator.User {
	var users []data_generator.User
	for i := 0; i < USERNUM; i++ {
		var user data_generator.User
		user.User_ID, user.Person_ID, user.Aspects = datagen.NewUser(genConfig.DBConn)
		users = append(users, user)
	}
	return users
}

func genPosts(genConfig *data_generator.GenConfig, users []data_generator.User) {
	assignment := data_generator.AssignPostsToUsersByPopScores(genConfig, USERNUM, POSTNUM)

	sequence := 0
	for _, user := range users {
		for n := 0; n < assignment[sequence]; n++ {
			datagen.NewPost(genConfig.DBConn, user.User_ID, user.Person_ID, user.Aspects)
		}
		sequence += 1
	}
}

func genFollows(genConfig *data_generator.GenConfig, users []data_generator.User) {
	
}

func genComments() {

}

func main() {
	genConfig := data_generator.Initialize(APP, USERNUM)
	// log.Println(genConfig.LikeScores)
	// log.Println(genConfig.CommentScores)
	users := genUsers(genConfig)
	genFollows(genConfig, users)
	genPosts(genConfig, users)
}
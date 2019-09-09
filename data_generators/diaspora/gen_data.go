package main

import (
	"diaspora/datagen"
	"diaspora/data_generator"
	"diaspora/helper"
	"log"
)

const APP = "diaspora" 
const USERNUM = 1000
const FOLLOWNUM = 30375
const POSTNUM = 8030
const COMMENTNUM = 13970

func genUsers(genConfig *data_generator.GenConfig) []data_generator.User {
	var users []data_generator.User
	for i := 0; i < USERNUM; i++ {
		var user data_generator.User
		user.User_ID, user.Person_ID, user.Aspects = datagen.NewUser(genConfig.DBConn)
		users = append(users, user)
	}
	return users
}

func genFollows(genConfig *data_generator.GenConfig, users []data_generator.User) {
	followedAssignment := data_generator.AssignDataToUsersByUserScores(genConfig.UserPopularityScores, FOLLOWNUM)
	
	log.Println(followedAssignment)
	log.Println(data_generator.GetSumOfIntSlice(followedAssignment))

	// alreadyFollowedEnough := make(map[int]bool)
	
	for seq1, user1 := range users {
		ableToBeFollowed := true
		personID1 := user1.Person_ID
		toBeFollowed := followedAssignment[seq1]
		// if _, ok := alreadyFollowedEnough[seq1]; ok {
		// 	continue
		// }
		// if fNum := datagen.GetFollowedNum(genConfig.DBConn, personID1); fNum == toBeFollowed {
		// 	alreadyFollowedEnough[seq1] = true
		// 	continue
		// } else {
		// 	toBeFollowed = toBeFollowed - fNum
		// }
		for n := 0; n < toBeFollowed; n++ {
			haveTried := make(map[int]bool)
			for {
				if len(haveTried) == USERNUM - 1 {
					log.Println("Cannot find more users to follow this user!!")
					log.Println("Total users to follow this user:", followedAssignment[seq1])
					log.Println("Have been followed by:", n + followedAssignment[seq1] - toBeFollowed)
					ableToBeFollowed = false
					break
				}
				seq2 := data_generator.RandomNonnegativeIntWithUpperBound(USERNUM)
				if seq2 == seq1 {
					continue
				}
				if _, ok := haveTried[seq2]; ok {
					continue
				}
				// if _, ok := alreadyFollowedEnough[seq2]; ok {
				// 	haveTried[seq2] = true
				// 	continue
				// }
				personID2 := users[seq2].Person_ID
				// if datagen.GetFollowedNum(genConfig.DBConn, personID2) == followedAssignment[seq2] {
				// 	alreadyFollowedEnough[seq2] = true
				// 	haveTried[seq2] = true
				// 	continue
				// }
				if datagen.CheckFollowed(genConfig.DBConn, personID1, personID2) {
					haveTried[seq2] = true
				} else {
					aspect_idx := helper.RandomNumber(0, len(user1.Aspects) - 1)
					datagen.FollowUser(genConfig.DBConn, personID2, personID1, user1.Aspects[aspect_idx])
					break
				}
			}
			if !ableToBeFollowed {
				break
			}
		}
	}
}

func genPosts(genConfig *data_generator.GenConfig, users []data_generator.User) map[int]float64 {
	postAssignment := data_generator.AssignDataToUsersByUserScores(genConfig.UserPopularityScores, POSTNUM)

	for seq, user := range users {
		for n := 0; n < postAssignment[seq]; n++ {
			datagen.NewPost(genConfig.DBConn, user.User_ID, user.Person_ID, user.Aspects)
		}
	}

	return data_generator.AssignScoresToPosts(datagen.GetAllPostIDs(genConfig.DBConn))
}

func genComments(genConfig *data_generator.GenConfig, postScores map[int]float64) {
	commentAssignment := data_generator.AssignDataToUsersByUserScores(genConfig.UserCommentScores, COMMENTNUM)
	log.Println(commentAssignment)

}

func main() {
	genConfig := data_generator.Initialize(APP, USERNUM)
	// log.Println(genConfig.LikeScores)
	// log.Println(genConfig.CommentScores)
	users := genUsers(genConfig)
	// var users []data_generator.User
	genFollows(genConfig, users)
	// postScores := genPosts(genConfig, users)
	// genComments(genConfig, postScores)
}
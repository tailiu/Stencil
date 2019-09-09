package main

import (
	"diaspora/datagen"
	"diaspora/data_generator"
	"diaspora/helper"
	"log"
)

const APP = "diaspora" 
const USER_NUM = 1000
const FOLLOW_NUM = 30375
const POST_NUM = 8030
const COMMENT_NUM = 13970
const RECIPROCAL_FOLLOW_PERCENTAGE = 0.3

func genUsers(genConfig *data_generator.GenConfig) []data_generator.User {
	var users []data_generator.User
	for i := 0; i < USER_NUM; i++ {
		var user data_generator.User
		user.User_ID, user.Person_ID, user.Aspects = datagen.NewUser(genConfig.DBConn)
		users = append(users, user)
	}
	return users
}

func genFollows(genConfig *data_generator.GenConfig, users []data_generator.User) {
	followedAssignment := data_generator.AssignDataToUsersByUserScores(genConfig.UserPopularityScores, FOLLOW_NUM)
	
	log.Println(followedAssignment)
	log.Println(data_generator.GetSumOfIntSlice(followedAssignment))
	
	for seq1, user1 := range users {
		var toBeFollowedByPersons []int
		ableToBeFollowed := true
		personID1 := user1.Person_ID
		alreadyFollowedByPersons := datagen.GetFollowedUsers(genConfig.DBConn, personID1)
		toBeFollowed := followedAssignment[seq1] - len(alreadyFollowedByPersons)
		toBeFollowedByPersons = append(toBeFollowedByPersons, data_generator.GetSeqsByPersonIDs(users, alreadyFollowedByPersons)...)
		// log.Println("Check user:", seq1)
		// log.Println("Total users to follow this user:", followedAssignment[seq1])
		// log.Println("already followed by", alreadyFollowedByPersons)
		// log.Println("To be followed by:", toBeFollowed)
		if toBeFollowed > followedAssignment[seq1] {
			log.Fatal("cannot happend1!!!!")
		}
		for n := 0; n < toBeFollowed; n++ {
			haveTried := make(map[int]bool)
			for {
				if len(haveTried) == USER_NUM - 1 {
					log.Println("Cannot find more users to follow this user!!")
					log.Println("Total users to follow this user:", followedAssignment[seq1])
					log.Println("Have been followed by:", n + followedAssignment[seq1] - toBeFollowed)
					ableToBeFollowed = false
					log.Fatal("cannot happend2!!!!")
					break
				}
				seq2 := data_generator.RandomNonnegativeIntWithUpperBound(USER_NUM)
				if seq2 == seq1 {
					continue
				}
				if _, ok := haveTried[seq2]; ok {
					continue
				}
				personID2 := users[seq2].Person_ID

				if datagen.CheckFollowed(genConfig.DBConn, personID1, personID2) {
					haveTried[seq2] = true
				} else {
					aspect_idx := helper.RandomNumber(0, len(user1.Aspects) - 1)
					datagen.FollowUser(genConfig.DBConn, personID2, personID1, user1.Aspects[aspect_idx])
					toBeFollowedByPersons = append(toBeFollowedByPersons, seq2)
					break
				}
			}
			if !ableToBeFollowed {
				break
			}
		}

		// Note that this RECIPROCAL_FOLLOW_PERCENTAGE cannot guarantee that
		// this user can follow this percentage of users following the current user
		// because maybe most of those users have already fully followed by other users
		toFollowNum := int(float64(followedAssignment[seq1]) * RECIPROCAL_FOLLOW_PERCENTAGE)
		// log.Println("Total Num", followedAssignment[seq1])
		// log.Println("to Follow Num", toFollowNum)
		currentlyFollowNum := 0
		for _, seq3 := range toBeFollowedByPersons {
			personID3 := users[seq3].Person_ID
			if currentlyFollowNum == toFollowNum {
				break
			}
			if datagen.CheckFollowed(genConfig.DBConn, personID3, personID1) {
				currentlyFollowNum += 1
				continue
			} else {
				if datagen.GetFollowedNum(genConfig.DBConn, personID3) >= followedAssignment[seq3] {
					continue
				} else {
					aspect_idx := helper.RandomNumber(0, len(users[seq3].Aspects) - 1)
					datagen.FollowUser(genConfig.DBConn, personID1, personID3, user1.Aspects[aspect_idx])
					currentlyFollowNum += 1
				}
			}
		}
		if currentlyFollowNum < toFollowNum {
			// log.Println("Fail to follow enough followers!!")
		}
		
	}
}

func genPosts(genConfig *data_generator.GenConfig, users []data_generator.User) map[int]float64 {
	postAssignment := data_generator.AssignDataToUsersByUserScores(genConfig.UserPopularityScores, POST_NUM)

	for seq, user := range users {
		for n := 0; n < postAssignment[seq]; n++ {
			datagen.NewPost(genConfig.DBConn, user.User_ID, user.Person_ID, user.Aspects)
		}
	}

	return data_generator.AssignScoresToPosts(datagen.GetAllPostIDs(genConfig.DBConn))
}

func genComments(genConfig *data_generator.GenConfig, postScores map[int]float64) {
	commentAssignment := data_generator.AssignDataToUsersByUserScores(genConfig.UserCommentScores, COMMENT_NUM)
	log.Println(commentAssignment)

}

func main() {
	genConfig := data_generator.Initialize(APP, USER_NUM)
	// log.Println(genConfig.LikeScores)
	// log.Println(genConfig.CommentScores)
	users := genUsers(genConfig)
	// var users []data_generator.User
	genFollows(genConfig, users)
	// postScores := genPosts(genConfig, users)
	// genComments(genConfig, postScores)
	// log.Println(datagen.GetFollowedDistribution(genConfig.DBConn))
}
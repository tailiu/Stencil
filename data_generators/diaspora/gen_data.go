package main

import (
	"diaspora/datagen"
	"diaspora/data_generator"
	"diaspora/helper"
	"log"
)

const APP = "diaspora" 
const USERNUM = 1000
const POSTNUM = 8030
const FRIENDSHIPNUM = 30375

func genUsers(genConfig *data_generator.GenConfig) []data_generator.User {
	var users []data_generator.User
	for i := 0; i < USERNUM; i++ {
		var user data_generator.User
		user.User_ID, user.Person_ID, user.Aspects = datagen.NewUser(genConfig.DBConn)
		users = append(users, user)
	}
	return users
}

func genFriends(genConfig *data_generator.GenConfig, users []data_generator.User) {
	friendshipAssignment := data_generator.AssignDataToUsersByUserPopScores(genConfig, USERNUM, FRIENDSHIPNUM)
	
	log.Println(friendshipAssignment)
	log.Println(data_generator.GetSumOfIntSlice(friendshipAssignment))

	alreadyMakeEnoughFriends := make(map[int]bool)
	
	for seq1, user1 := range users {
		ableToMakeFriends := true
		personID1 := user1.Person_ID
		friendsTomake := friendshipAssignment[seq1]
		if _, ok := alreadyMakeEnoughFriends[seq1]; ok {
			continue
		}
		if fNum := datagen.GetFriendsNum(genConfig.DBConn, personID1); fNum == friendsTomake {
			alreadyMakeEnoughFriends[seq1] = true
			continue
		} else {
			friendsTomake = friendsTomake - fNum
		}
		for n := 0; n < friendsTomake; n++ {
			haveTried := make(map[int]bool)
			for {
				if len(haveTried) == USERNUM - 1 {
					log.Println("Cannot find more users to make friends!!")
					log.Println("Total friends to make:", friendshipAssignment[seq1])
					log.Println("Have made friends:", n)
					ableToMakeFriends = false
					break
				}
				seq2 := data_generator.RandomNonnegativeIntWithUpperBound(USERNUM)
				if seq2 == seq1 {
					continue
				}
				if _, ok := haveTried[seq2]; ok {
					continue
				}
				if _, ok := alreadyMakeEnoughFriends[seq2]; ok {
					continue
				}
				personID2 := users[seq2].Person_ID
				if datagen.GetFriendsNum(genConfig.DBConn, personID2) == friendshipAssignment[seq2] {
					alreadyMakeEnoughFriends[seq2] = true
					haveTried[seq2] = true
					continue
				}
				exists, _ := datagen.ContactExists(genConfig.DBConn, personID1, personID2)
				if exists {
					haveTried[seq2] = true
				} else {
					aspect_idx := helper.RandomNumber(0, len(user1.Aspects) - 1)
					datagen.FollowUser(genConfig.DBConn, personID1, personID2, user1.Aspects[aspect_idx])
					datagen.FollowUser(genConfig.DBConn, personID2, personID1, users[seq2].Aspects[aspect_idx])
					break
				}
			}
			if !ableToMakeFriends {
				break
			}
		}
	}
}

func genPosts(genConfig *data_generator.GenConfig, users []data_generator.User) {
	postAssignment := data_generator.AssignDataToUsersByUserPopScores(genConfig, USERNUM, POSTNUM)

	for seq, user := range users {
		for n := 0; n < postAssignment[seq]; n++ {
			datagen.NewPost(genConfig.DBConn, user.User_ID, user.Person_ID, user.Aspects)
		}
	}

	// genConfig.PostPopularityScores = data_generator.shuffleSlices(
	// 	data_generator.ParetoScores(data_generator.ALPHA, data_generator.XM, data_generator.GetSumOfIntSlice(postAssignment)))
}

func genComments(genConfig *data_generator.GenConfig) {
	
}

func main() {
	genConfig := data_generator.Initialize(APP, USERNUM)
	// log.Println(genConfig.LikeScores)
	// log.Println(genConfig.CommentScores)
	users := genUsers(genConfig)
	genFriends(genConfig, users)
	genPosts(genConfig, users)
	genComments(genConfig)
}
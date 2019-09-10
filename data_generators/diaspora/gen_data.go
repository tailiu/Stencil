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
const LIKE_NUM = 85680
const RECIPROCAL_FOLLOW_PERCENTAGE = 0.3
const MESSAGE_NUM = 4015

func genUsers(genConfig *data_generator.GenConfig) []data_generator.User {
	var users []data_generator.User
	for i := 0; i < USER_NUM; i++ {
		var user data_generator.User
		user.User_ID, user.Person_ID, user.Aspects = datagen.NewUser(genConfig.DBConn)
		users = append(users, user)
	}
	return users
}

// We use the user popularity score to generate how many followers a user has.
// For some of those followers, the user also follows back.
// Note: RECIPROCAL_FOLLOW_PERCENTAGE cannot guarantee that
// this user can follow this percentage of followers,
// because maybe most of those users have already fully followed by other users.
func genFollows(genConfig *data_generator.GenConfig, users []data_generator.User) {
	followedAssignment := data_generator.AssignDataToUsersByUserScores(genConfig.UserPopularityScores, FOLLOW_NUM)
	
	log.Println("Followed assignment to users:", followedAssignment)
	log.Println("Total followed:", data_generator.GetSumOfIntSlice(followedAssignment))
	
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
		// if currentlyFollowNum < toFollowNum {
			// log.Println("Fail to follow enough followers!!")
		// }

	}
}

// We randomly assign posts to the users proportionally to the popularity of users. 
func genPosts(genConfig *data_generator.GenConfig, users []data_generator.User) map[int]float64 {
	postAssignment := data_generator.AssignDataToUsersByUserScores(genConfig.UserPopularityScores, POST_NUM)
	log.Println("Posts assignments to users:", postAssignment)
	log.Println("Total posts:", data_generator.GetSumOfIntSlice(postAssignment))

	for seq, user := range users {
		for n := 0; n < postAssignment[seq]; n++ {
			datagen.NewPost(genConfig.DBConn, user.User_ID, user.Person_ID, user.Aspects)
		}
	}

	return data_generator.AssignParetoDistributionScoresToData(datagen.GetAllPostIDs(genConfig.DBConn))
}

// Only for test
func prepareTest(genConfig *data_generator.GenConfig) ([]data_generator.User, map[int]float64){
	var users []data_generator.User
	users1 := datagen.GetAllUsersWithAspectsOrderByID(genConfig.DBConn)
	for _, user1 := range users1 {
		var user data_generator.User
		user.User_ID, user.Person_ID, user.Aspects = user1.User_ID, user1.Person_ID, user1.Aspects
		users = append(users, user)
	}
	return users, data_generator.AssignParetoDistributionScoresToData(datagen.GetAllPostIDs(genConfig.DBConn))
}

// We randomly assign comments to posts proportionally to the popularity of posts of friends, 
// including posts by the commenter.
func genComments(genConfig *data_generator.GenConfig, users []data_generator.User, postScores map[int]float64) {
	commentAssignment := data_generator.AssignDataToUsersByUserScores(genConfig.UserCommentScores, COMMENT_NUM)
	log.Println("Comments assignments to users:", commentAssignment)
	log.Println("Total comments:", data_generator.GetSumOfIntSlice(commentAssignment))

	for seq1, user1 := range users {
		// log.Println("Check user:", seq1)
		var posts []*data_generator.Post
		var scores []float64
		commentNum := commentAssignment[seq1]
		// log.Println("Comment number:", commentNum)
		personID := user1.Person_ID
		totalUsers := datagen.GetFollowingUsers(genConfig.DBConn, personID)
		// log.Println(user1)
		// log.Println(totalUsers)
		totalUsers = append(totalUsers, personID)

		for _, user2 := range totalUsers {
			posts1 := datagen.GetPostsForUser(genConfig.DBConn, user2)
			for _, post1 := range posts1 {
				post := new(data_generator.Post)
				post.ID = post1.ID
				post.Author = post1.Author
				post.Score = postScores[post1.ID]
				posts = append(posts, post)
				scores = append(scores, post.Score)
			}
		}
		
		likeNumsOfPosts := data_generator.RandomNumWithProbGenerator(scores, commentNum)
		for seq2, post := range posts {
			for i := 0; i < likeNumsOfPosts[seq2]; i++ {
				datagen.NewComment(genConfig.DBConn, post.ID, personID, post.Author)
			}
		}
	}
}

// We randomly assign likes to posts proportionally to the popularity of posts of friends, 
// including posts by the liker.
// The difference between generating comments and likes is that
// a user make several comments on the same post, but can only like once on that post.
func genLikes(genConfig *data_generator.GenConfig, users []data_generator.User, postScores map[int]float64) {
	likeAssignment := data_generator.AssignDataToUsersByUserScores(genConfig.UserLikeScores, LIKE_NUM)
	log.Println("Likes assignments to users:", likeAssignment)
	log.Println("Total likes:", data_generator.GetSumOfIntSlice(likeAssignment))

	for seq1, user1 := range users {
		// log.Println("Check user:", seq1)
		var posts []*data_generator.Post
		var scores []float64
		likeNum := likeAssignment[seq1]
		// log.Println("Like number:", likeNum)
		personID := user1.Person_ID
		totalUsers := datagen.GetFollowingUsers(genConfig.DBConn, personID)
		// log.Println(user1)
		// log.Println(totalUsers)
		totalUsers = append(totalUsers, personID)

		for _, user2 := range totalUsers {
			posts1 := datagen.GetPostsForUser(genConfig.DBConn, user2)
			for _, post1 := range posts1 {
				post := new(data_generator.Post)
				post.ID = post1.ID
				post.Author = post1.Author
				post.Score = postScores[post1.ID]
				posts = append(posts, post)
				scores = append(scores, post.Score)
			}
		}
		
		likeNumsOfPosts := data_generator.RandomNumWithProbGenerator(scores, likeNum)
		// log.Println(likeNumsOfPosts)
		for seq2, post := range posts {
			if _, ok := likeNumsOfPosts[seq2]; ok {
				datagen.NewLike(genConfig.DBConn, post.ID, personID, post.Author)
			}
		}
	}
}

func genConversationsAndMessages(genConfig *data_generator.GenConfig, users []data_generator.User) {
	messageAssignment := data_generator.AssignDataToUsersByUserScores(genConfig.UserMessageScores, MESSAGE_NUM)
	log.Println("Messages assignments to users:", messageAssignment)
	log.Println("Total messages:", data_generator.GetSumOfIntSlice(messageAssignment))
	conversationNum := 0

	for seq1, user1 := range users {
		personID := user1.Person_ID
		messageNum := messageAssignment[seq1]
		friends := datagen.GetRealFriendsOfUser(genConfig.DBConn, personID)
		friendCI := data_generator.AssignParetoDistributionScoresToDataReturnSlice(len(friends))
		// log.Println(friends)
		// log.Println(friendCI)
		// log.Println(messageNum)
		
		messageNumsOfConversations := data_generator.RandomNumWithProbGenerator(friendCI, messageNum)
		// log.Println(messageNumsOfConversations)

		for seq2, messageNum := range messageNumsOfConversations {
			exists, conv_id := datagen.CheckConversationBetweenTwoUsers(genConfig.DBConn, personID, friends[seq2])
			if exists {
				for i := 0; i < messageNum; i++ {
					datagen.NewMessage(genConfig.DBConn, personID, conv_id)
				}
			} else {
				new_conv, _ := datagen.NewConversation(genConfig.DBConn, personID, friends[seq2])
				conversationNum += 1
				for i := 0; i < messageNum; i++ {
					datagen.NewMessage(genConfig.DBConn, personID, new_conv)
				}
			}
		}
	}
	log.Println("Total conversations:", conversationNum)
}

func main() {
	genConfig := data_generator.Initialize(APP, USER_NUM)
	users, _ := prepareTest(genConfig)
	// users := genUsers(genConfig)
	// postScores := genPosts(genConfig, users)
	// genFollows(genConfig, users)
	// genComments(genConfig, users, postScores)
	// genLikes(genConfig, users, postScores)
	genConversationsAndMessages(genConfig, users)
}
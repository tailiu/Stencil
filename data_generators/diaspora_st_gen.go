package main

import (
	"data_generators/diaspora/datagen"
	"data_generators/diaspora/helper"
	"data_generators/data_generator"
	"time"
	"log"
)

/**
 * This is to generate data for Diaspora with single thread
*/

const APP = "diaspora_1000_st" 
const USER_NUM = 1000
const FOLLOW_NUM = 20660
const POST_NUM = 8029
const COMMENT_NUM = 13971
const LIKE_NUM = 85672
const RECIPROCAL_FOLLOW_PERCENTAGE = 0.3
const MESSAGE_NUM = 40146
const IMAGE_NUM = 3693

// const APP = "diaspora" 
// const USER_NUM = 10000
// const FOLLOW_NUM = 206600
// const POST_NUM = 80292
// const COMMENT_NUM = 139708
// const LIKE_NUM = 856715
// const RECIPROCAL_FOLLOW_PERCENTAGE = 0.3
// const MESSAGE_NUM = 40146
// const IMAGE_NUM = 36934

// const APP = "diaspora_100000" 
// const USER_NUM = 100500
// const FOLLOW_NUM = 323000
// const POST_NUM = 802920
// const COMMENT_NUM = 1397080
// const LIKE_NUM = 8567156
// const RECIPROCAL_FOLLOW_PERCENTAGE = 0.3
// const MESSAGE_NUM = 401460
// const IMAGE_NUM = 369343

// const APP = "diaspora_1000000" 
// const USER_NUM = 1010000
// const FOLLOW_NUM = 3230000
// const POST_NUM = 8029200
// const COMMENT_NUM = 13970800
// const LIKE_NUM = 85671564
// const RECIPROCAL_FOLLOW_PERCENTAGE = 0.3
// const MESSAGE_NUM = 4014600
// const IMAGE_NUM = 3693432


// Function genUsers() tries to create USER_NUM users, but it cannot guarantee
func genUsers(dataGen *data_generator.DataGen) []data_generator.DUser {

	var users []data_generator.DUser

	for i := 0; i < USER_NUM; i++ {

		var user data_generator.DUser

		var err error
		
		user.User_ID, user.Person_ID, user.Aspects, err = datagen.NewUser(dataGen.DBConn)

		if err != nil {

			// log.Println(err)

		} else {
			
			users = append(users, user)
			
		}
	}

	log.Println("Total number of users:", len(users))

	return users
	
}

// We use the user popularity score to generate how many followers a user has.
// For some of those followers, the user also follows back.
// Note: RECIPROCAL_FOLLOW_PERCENTAGE cannot guarantee that
// this user can follow this percentage of followers,
// because maybe most of those users have already fully followed by other users.
func genFollows(dataGen *data_generator.DataGen, users []data_generator.DUser) {

	followedAssignment := data_generator.AssignDataToUsersByUserScores(
		dataGen.UserPopularityScores, FOLLOW_NUM)
	
	log.Println("Followed assignment to users:", followedAssignment)
	log.Println("Total followed:", data_generator.GetSumOfIntSlice(followedAssignment))
	
	for seq1, user1 := range users {

		var toBeFollowedByPersons []int

		ableToBeFollowed := true
		personID1 := user1.Person_ID
		alreadyFollowedByPersons := datagen.GetFollowedUsers(dataGen.DBConn, personID1)

		toBeFollowed := followedAssignment[seq1] - len(alreadyFollowedByPersons)
		toBeFollowedByPersons = append(toBeFollowedByPersons, 
			data_generator.GetSeqsByPersonIDs(users, alreadyFollowedByPersons)...)

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

				if len(haveTried) == len(users) - 1 {

					log.Println("Cannot find more users to follow this user!!")
					log.Println("Total users to follow this user:", followedAssignment[seq1])
					log.Println("Have been followed by:", n + followedAssignment[seq1] - toBeFollowed)
					
					ableToBeFollowed = false
					
					log.Fatal("cannot happend2!!!!")
					break

				}

				seq2 := data_generator.RandomNonnegativeIntWithUpperBound(len(users))
				if seq2 == seq1 {
					continue
				}
				if _, ok := haveTried[seq2]; ok {
					continue
				}
				personID2 := users[seq2].Person_ID

				if datagen.CheckFollowed(dataGen.DBConn, personID1, personID2) {

					haveTried[seq2] = true

				} else {

					aspect_idx := helper.RandomNumber(0, len(user1.Aspects) - 1)
					datagen.FollowUser(dataGen.DBConn, 
						personID2, personID1, user1.Aspects[aspect_idx])
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

			if datagen.CheckFollowed(dataGen.DBConn, personID3, personID1) {
				
				currentlyFollowNum += 1

				continue

			} else {

				if datagen.GetFollowedNum(dataGen.DBConn, personID3) >= followedAssignment[seq3] {
					
					continue

				} else {

					aspect_idx := helper.RandomNumber(0, len(users[seq3].Aspects) - 1)
					datagen.FollowUser(dataGen.DBConn, personID1, personID3, user1.Aspects[aspect_idx])
					
					currentlyFollowNum += 1

				}
			}
		}

		// if currentlyFollowNum < toFollowNum {
			// log.Println("Fail to follow enough followers!!")
		// }

	}
}

// The number of posts of users is proportional to the popularity of users.
// We also randomly assign images to the posts proportionally to the popularity of posts.
// The scores assigned to posts are in pareto distributiuon.
// so it is more likely that popular users will have popular posts because they have more posts
func genPosts(dataGen *data_generator.DataGen, users []data_generator.DUser) map[int]float64 {

	postAssignment := data_generator.AssignDataToUsersByUserScores(dataGen.UserPopularityScores, POST_NUM)
	totalPosts := data_generator.GetSumOfIntSlice(postAssignment)
	
	log.Println("Posts assignments to users:", postAssignment)
	log.Println("Total posts:", totalPosts)

	seqNum := data_generator.MakeRange(0, totalPosts - 1)
	seqScores := data_generator.AssignParetoDistributionScoresToDataReturnSlice(len(seqNum))
	imageNumsOfSeq := data_generator.RandomNumWithProbGenerator(seqScores, IMAGE_NUM)

	postSeq := 0
	imageNums := 0
	postScores := make(map[int]float64)
	
	for userSeq, user := range users {

		for n := 0; n < postAssignment[userSeq]; n++ {

			var postID int
			imageNum := imageNumsOfSeq[postSeq]
			if imageNum == 0 {

				postID = datagen.NewPost(dataGen.DBConn, 
					user.User_ID, user.Person_ID, user.Aspects)

			} else {

				postID = datagen.NewPhotoPost(dataGen.DBConn, 
					user.User_ID, user.Person_ID, user.Aspects, imageNum)

			}

			postScores[postID] = seqScores[postSeq]
			postSeq += 1
			imageNums += imageNum

		}
	}
	
	log.Println("Total images:", imageNums)

	return postScores

}

// Only for test
func prepareTest(dataGen *data_generator.DataGen) ([]data_generator.DUser, map[int]float64){

	var users []data_generator.DUser
	users1 := datagen.GetAllUsersWithAspectsOrderByID(dataGen.DBConn)

	for _, user1 := range users1 {

		var user data_generator.DUser
		user.User_ID, user.Person_ID, user.Aspects = user1.User_ID, user1.Person_ID, user1.Aspects
		users = append(users, user)

	}

	return users, data_generator.AssignParetoDistributionScoresToData(
		datagen.GetAllPostIDs(dataGen.DBConn))

}

// We randomly assign comments to posts proportionally to the popularity of posts of friends, 
// including posts by the commenter.
func genComments(dataGen *data_generator.DataGen, 
	users []data_generator.DUser, postScores map[int]float64) {

	commentAssignment := data_generator.AssignDataToUsersByUserScores(
		dataGen.UserCommentScores, COMMENT_NUM)

	log.Println("Comments assignments to users:", commentAssignment)
	log.Println("Total comments:", data_generator.GetSumOfIntSlice(commentAssignment))

	for seq1, user1 := range users {

		// log.Println("Check user:", seq1)
		
		var posts []*data_generator.Post
		var scores []float64
		commentNum := commentAssignment[seq1]

		// log.Println("Comment number:", commentNum)
		
		personID := user1.Person_ID
		totalUsers := datagen.GetFollowingUsers(dataGen.DBConn, personID)

		// log.Println(user1)
		// log.Println(totalUsers)

		totalUsers = append(totalUsers, personID)

		for _, user2 := range totalUsers {

			posts1 := datagen.GetPostsForUser(dataGen.DBConn, user2)

			for _, post1 := range posts1 {

				post := new(data_generator.Post)
				post.ID = post1.ID
				post.Author = post1.Author
				post.Score = postScores[post1.ID]

				posts = append(posts, post)
				scores = append(scores, post.Score)

			}
		}
		
		commentNumsOfPosts := data_generator.RandomNumWithProbGenerator(scores, commentNum)

		for seq2, post := range posts {

			for i := 0; i < commentNumsOfPosts[seq2]; i++ {

				datagen.NewComment(dataGen.DBConn, post.ID, personID, post.Author)
			}
		}
	}
}

// We randomly assign likes to posts proportionally to the popularity of posts of friends, 
// including posts by the liker.
// The difference between generating comments and likes is that
// a user make several comments on the same post, but can only like once on that post.
func genLikes(dataGen *data_generator.DataGen, 
	users []data_generator.DUser, postScores map[int]float64) {

	likeAssignment := data_generator.AssignDataToUsersByUserScores(
		dataGen.UserLikeScores, LIKE_NUM)

	log.Println("Likes assignments to users:", likeAssignment)
	log.Println("Total likes based on assignments:", data_generator.GetSumOfIntSlice(likeAssignment))
	
	totalLikeNum := 0

	for seq1, user1 := range users {

		// log.Println("Check user:", seq1)
		
		var posts []*data_generator.Post
		var scores []float64
		
		likeNum := likeAssignment[seq1]
		// log.Println("Like number:", likeNum)
		
		personID := user1.Person_ID
		totalUsers := datagen.GetFollowingUsers(dataGen.DBConn, personID)
	
		// log.Println(user1)
		// log.Println(totalUsers)
		
		totalUsers = append(totalUsers, personID)

		for _, user2 := range totalUsers {
			
			posts1 := datagen.GetPostsForUser(dataGen.DBConn, user2)
			
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

				datagen.NewLike(dataGen.DBConn, post.ID, personID, post.Author)
				totalLikeNum += 1
			}
		}
	}

	log.Println("In reality, the num of total likes is:", totalLikeNum)

}

// Pareto-distributed message scores determine the number of messages each user should have.
// Friendships have pareto-distributed closeness indexes. 
// We randomly assign messages to users (or conversations) proportionally 
// to the closeness indexes.
// Two users talk with each other sharing the same conversation.
func genConversationsAndMessages(dataGen *data_generator.DataGen, users []data_generator.DUser) {

	messageAssignment := data_generator.AssignDataToUsersByUserScores(
		dataGen.UserMessageScores, MESSAGE_NUM)

	log.Println("Messages assignments to users:", messageAssignment)
	log.Println("Total messages:", data_generator.GetSumOfIntSlice(messageAssignment))
	
	conversationNum := 0

	for seq1, user1 := range users {

		// oneUserConversationNum := 0
		
		personID := user1.Person_ID
		messageNum := messageAssignment[seq1]
		friends := datagen.GetRealFriendsOfUser(dataGen.DBConn, personID)
		friendCI := data_generator.AssignParetoDistributionScoresToDataReturnSlice(len(friends))
		
		// log.Println(friends)
		// log.Println(friendCI)
		// log.Println(messageNum)
		
		messageNumsOfConversations := data_generator.RandomNumWithProbGenerator(friendCI, messageNum)
		// log.Println(messageNumsOfConversations)

		for seq2, messageNum := range messageNumsOfConversations {

			exists, conv_id := datagen.CheckConversationBetweenTwoUsers(dataGen.DBConn, 
				personID, friends[seq2])
			
			if exists {

				for i := 0; i < messageNum; i++ {

					datagen.NewMessage(dataGen.DBConn, personID, conv_id)

				}

			} else {

				new_conv, _ := datagen.NewConversation(dataGen.DBConn, 
					personID, friends[seq2])
				
				conversationNum += 1
				
				// oneUserConversationNum += 1
				
				for i := 0; i < messageNum; i++ {

					datagen.NewMessage(dataGen.DBConn, personID, new_conv)

				}
			}
		}

		// log.Println(oneUserConversationNum)
	}

	log.Println("Total conversations:", conversationNum)

}

func main() {
	
	startTime := time.Now()

	log.Println("--------- Start of Data Generation ---------")

	// users, postScores := prepareTest(dataGen)

	dataGen := data_generator.Initialize(APP)

	users := genUsers(dataGen)

	data_generator.InitializeWithUserNum(dataGen, len(users))
	
	postScores := genPosts(dataGen, users)
	
	genFollows(dataGen, users)
	
	genComments(dataGen, users, postScores)
	
	genLikes(dataGen, users, postScores)
	
	genConversationsAndMessages(dataGen, users)

	log.Println("--------- End of Data Generation ---------")

	endTime := time.Now()

	log.Println("Time used: ", endTime.Sub(startTime))

}

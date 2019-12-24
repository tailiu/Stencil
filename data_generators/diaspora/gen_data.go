package main

import (
	"diaspora/datagen"
	"diaspora/data_generator"
	"diaspora/helper"
	"time"
	"log"
	"sync"
)

// Note that THREAD_NUM must be larger than 0
const THREAD_NUM = 50

// const APP = "diaspora" 
// const USER_NUM = 10000
// const FOLLOW_NUM = 206600
// const POST_NUM = 80292
// const COMMENT_NUM = 139708
// const LIKE_NUM = 856715
// const RECIPROCAL_FOLLOW_PERCENTAGE = 0.3
// const MESSAGE_NUM = 40146
// const IMAGE_NUM = 36934

const APP = "diaspora_100000" 
const USER_NUM = 100500
const FOLLOW_NUM = 323000
const POST_NUM = 802920
const COMMENT_NUM = 1397080
const LIKE_NUM = 8567156
const RECIPROCAL_FOLLOW_PERCENTAGE = 0.3
const MESSAGE_NUM = 401460
const IMAGE_NUM = 369343

// const APP = "diaspora_1000000" 
// const USER_NUM = 1010000
// const FOLLOW_NUM = 3230000
// const POST_NUM = 8029200
// const COMMENT_NUM = 13970800
// const LIKE_NUM = 85671564
// const RECIPROCAL_FOLLOW_PERCENTAGE = 0.3
// const MESSAGE_NUM = 4014600
// const IMAGE_NUM = 3693432


func genUsers(genConfig *data_generator.GenConfig, num int, wg *sync.WaitGroup, res chan<- []data_generator.User) {

	defer wg.Done()

	var users []data_generator.User

	for i := 0; i < num; i++ {

		var user data_generator.User

		var err error
		
		user.User_ID, user.Person_ID, user.Aspects, err = datagen.NewUser(genConfig.DBConn)

		if err != nil {

			// log.Println(err)

		} else {
			
			users = append(users, user)
			
		}
	}

	// log.Println("Total number of users:", len(users))

	res <- users
	
}

// Function genUsersController() tries to create USER_NUM users, but it cannot guarantee
func genUsersController(genConfig *data_generator.GenConfig) []data_generator.User {
	
	var users []data_generator.User

	channel := make(chan []data_generator.User, THREAD_NUM)

	var wg sync.WaitGroup

	wg.Add(THREAD_NUM)

	for i := 0; i < THREAD_NUM; i++ {
		
		go genUsers(genConfig, USER_NUM / THREAD_NUM, &wg, channel)
		 
	}

	wg.Wait()

	close(channel)

	for res := range channel {

		users = append(users, res...)

	}

	log.Println("Total number of users:", len(users))

	return users

}

func genFollows(genConfig *data_generator.GenConfig, wg *sync.WaitGroup, 
	userSeqStart, userSeqEnd int, followedAssignment []int, users []data_generator.User) {
	
	for seq1 := userSeqStart; seq1 < userSeqEnd; seq1++ {

		user1 := users[seq1]

		var toBeFollowedByPersons []int

		ableToBeFollowed := true
		personID1 := user1.Person_ID
		alreadyFollowedByPersons := datagen.GetFollowedUsers(genConfig.DBConn, personID1)

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

				if datagen.CheckFollowed(genConfig.DBConn, personID1, personID2) {

					haveTried[seq2] = true

				} else {

					aspect_idx := helper.RandomNumber(0, len(user1.Aspects) - 1)
					datagen.FollowUser(genConfig.DBConn, 
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

// We use the user popularity score to generate how many followers a user has.
// For some of those followers, the user also follows back.
// Note: RECIPROCAL_FOLLOW_PERCENTAGE cannot guarantee that
// this user can follow this percentage of followers,
// because maybe most of those users have already fully followed by other users.
func genFollowsController(genConfig *data_generator.GenConfig, users []data_generator.User) {

	followedAssignment := data_generator.AssignDataToUsersByUserScores(
		genConfig.UserPopularityScores, FOLLOW_NUM)
	
	log.Println("Followed assignment to users:", followedAssignment)
	log.Println("Total followed:", data_generator.GetSumOfIntSlice(followedAssignment))
	
	var wg sync.WaitGroup

	wg.Add(THREAD_NUM)

	userSeqStart := 0
	
	userSeqStep := len(users) / THREAD_NUM

	for i := 0; i < THREAD_NUM; i++ {
		
		if i != THREAD_NUM - 1 {

			// Start included, end (start + step) not included
			go genFollows(genConfig, &wg, userSeqStart, userSeqStart + userSeqStep, 
				followedAssignment, users)

		} else {

			// Start included, end (start + step) not included
			go genFollows(genConfig, &wg, userSeqStart, len(users), 
				followedAssignment, users)
		
		}

		userSeqStart += userSeqStep
	}

	wg.Wait()

}

func genPosts(genConfig *data_generator.GenConfig, wg *sync.WaitGroup, 
	res1 chan<- map[int]float64, res2 chan<- int, 
	userSeqStart, userSeqEnd, postSeqStart int,
	users []data_generator.User, postAssignment []int, 
	imageNumsOfSeq map[int]int, seqScores []float64) {

	defer wg.Done()
	
	postScores := make(map[int]float64)
	
	postSeq := postSeqStart

	imageNums := 0

	for i := userSeqStart; i < userSeqEnd; i++ {
		
		user := users[i]
		
		for n := 0; n < postAssignment[i]; n++ {

			var postID int
			
			imageNum := imageNumsOfSeq[postSeq]

			if imageNum == 0 {

				postID = datagen.NewPost(genConfig.DBConn, 
					user.User_ID, user.Person_ID, user.Aspects)

			} else {

				postID = datagen.NewPhotoPost(genConfig.DBConn, 
					user.User_ID, user.Person_ID, user.Aspects, imageNum)

			}

			postScores[postID] = seqScores[postSeq]
			postSeq += 1
			imageNums += imageNum

		}

	}
	
	res1 <- postScores

	res2 <- imageNums

}

// The number of posts of users is proportional to the popularity of users.
// We also randomly assign images to the posts proportionally to the popularity of posts.
// The scores assigned to posts are in pareto distributiuon.
// so it is more likely that popular users will have popular posts because they have more posts
func genPostsController(genConfig *data_generator.GenConfig, users []data_generator.User) map[int]float64 {

	postAssignment := data_generator.AssignDataToUsersByUserScores(genConfig.UserPopularityScores, POST_NUM)
	totalPosts := data_generator.GetSumOfIntSlice(postAssignment)
	
	log.Println("Posts assignments to users:", postAssignment)
	log.Println("Total posts:", totalPosts)

	seqNum := data_generator.MakeRange(0, totalPosts - 1)
	seqScores := data_generator.AssignParetoDistributionScoresToDataReturnSlice(len(seqNum))
	imageNumsOfSeq := data_generator.RandomNumWithProbGenerator(seqScores, IMAGE_NUM)
	
	imageNums := 0
	
	postScores := make(map[int]float64)
	
	channel1 := make(chan map[int]float64, THREAD_NUM)

	channel2 := make(chan int, THREAD_NUM)

	var wg sync.WaitGroup

	wg.Add(THREAD_NUM)

	userSeqStart := 0
	
	userSeqStep := len(users) / THREAD_NUM

	postSeqStart := 0

	for i := 0; i < THREAD_NUM; i++ {
		
		if i != THREAD_NUM - 1 {

			// Start included, end (start + step) not included
			go genPosts(genConfig, &wg, channel1, channel2, 
				userSeqStart, userSeqStart + userSeqStep, 
				postSeqStart, users, postAssignment, imageNumsOfSeq, seqScores)

		} else {

			// Start included, end (start + step) not included
			go genPosts(genConfig, &wg, channel1, channel2, 
				userSeqStart, len(users), 
				postSeqStart, users, postAssignment, imageNumsOfSeq, seqScores)
		
		}

		postSeqStart = data_generator.CalculateNextPostSeqStart(
			postSeqStart, userSeqStart, userSeqStart + userSeqStep, postAssignment)

		userSeqStart += userSeqStep
	}

	wg.Wait()

	close(channel1)

	close(channel2)

	for res1 := range channel1 {

		postScores = data_generator.MergeTwoMaps(postScores, res1) 

	}

	for res2 := range channel2 {

		imageNums += res2

	}

	log.Println("Total images:", imageNums)

	return postScores

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

	return users, data_generator.AssignParetoDistributionScoresToData(
		datagen.GetAllPostIDs(genConfig.DBConn))

}

// We randomly assign comments to posts proportionally to the popularity of posts of friends, 
// including posts by the commenter.
func genComments(genConfig *data_generator.GenConfig, 
	users []data_generator.User, postScores map[int]float64) {

	commentAssignment := data_generator.AssignDataToUsersByUserScores(
		genConfig.UserCommentScores, COMMENT_NUM)

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
		
		commentNumsOfPosts := data_generator.RandomNumWithProbGenerator(scores, commentNum)

		for seq2, post := range posts {

			for i := 0; i < commentNumsOfPosts[seq2]; i++ {

				datagen.NewComment(genConfig.DBConn, post.ID, personID, post.Author)
			}
		}
	}
}

// We randomly assign likes to posts proportionally to the popularity of posts of friends, 
// including posts by the liker.
// The difference between generating comments and likes is that
// a user make several comments on the same post, but can only like once on that post.
func genLikes(genConfig *data_generator.GenConfig, 
	users []data_generator.User, postScores map[int]float64) {

	likeAssignment := data_generator.AssignDataToUsersByUserScores(
		genConfig.UserLikeScores, LIKE_NUM)

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
func genConversationsAndMessages(genConfig *data_generator.GenConfig, users []data_generator.User) {

	messageAssignment := data_generator.AssignDataToUsersByUserScores(
		genConfig.UserMessageScores, MESSAGE_NUM)

	log.Println("Messages assignments to users:", messageAssignment)
	log.Println("Total messages:", data_generator.GetSumOfIntSlice(messageAssignment))
	
	conversationNum := 0

	for seq1, user1 := range users {

		// oneUserConversationNum := 0
		
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

			exists, conv_id := datagen.CheckConversationBetweenTwoUsers(genConfig.DBConn, 
				personID, friends[seq2])
			
			if exists {

				for i := 0; i < messageNum; i++ {

					datagen.NewMessage(genConfig.DBConn, personID, conv_id)

				}

			} else {

				new_conv, _ := datagen.NewConversation(genConfig.DBConn, 
					personID, friends[seq2])
				
				conversationNum += 1
				
				// oneUserConversationNum += 1
				
				for i := 0; i < messageNum; i++ {

					datagen.NewMessage(genConfig.DBConn, personID, new_conv)

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

	users, postScores := prepareTest(genConfig)

	// genConfig := data_generator.Initialize(APP)

	// users := genUsersController(genConfig)

	data_generator.InitializeWithUserNum(genConfig, len(users))

	genPostsController(genConfig, users)

	// genFollowsController(genConfig, users)

	// users := genUsers(genConfig)

	// data_generator.InitializeWithUserNum(genConfig, len(users))
	
	// postScores := genPosts(genConfig, users)
	
	// genFollows(genConfig, users)
	
	// genComments(genConfig, users, postScores)
	
	// genLikes(genConfig, users, postScores)
	
	// genConversationsAndMessages(genConfig, users)

	log.Println("--------- End of Data Generation ---------")

	endTime := time.Now()

	log.Println("Time used: ", endTime.Sub(startTime))

}
package data_generator

import (
	"data_generators/diaspora/datagen"
	"data_generators/diaspora/helper"
	"time"
	"log"
	"sync"
)

/**
 * This is a general algorithm to generate data
 * for the four test apps with multiple threads
**/
func (dataGen *DataGen) genUsers(num int, wg *sync.WaitGroup, res chan<- []DUser) {

	defer wg.Done()

	var users []DUser

	for i := 0; i < num; i++ {

		var user DUser

		var err error
		
		user.User_ID, user.Person_ID, user.Aspects, err = datagen.NewUser(dataGen.DBConn)

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
func (dataGen *DataGen) genUsersController() []DUser {
	
	var users []DUser

	channel := make(chan []DUser, THREAD_NUM)

	var wg sync.WaitGroup

	wg.Add(THREAD_NUM)

	for i := 0; i < THREAD_NUM; i++ {
		go dataGen.genUsers(USER_NUM / THREAD_NUM, &wg, channel)
	}

	wg.Wait()

	close(channel)

	for res := range channel {
		users = append(users, res...)
	}

	log.Println("Total number of users:", len(users))

	return users
}

func (dataGen *DataGen) genFollows(wg *sync.WaitGroup, 
	userSeqStart, userSeqEnd int, followedAssignment []int, users []DUser) {
	
	defer wg.Done()

	for seq1 := userSeqStart; seq1 < userSeqEnd; seq1++ {

		user1 := users[seq1]

		var toBeFollowedByPersons []int

		// ableToBeFollowed := true
		
		personID1 := user1.Person_ID

		alreadyFollowedByPersons := datagen.GetFollowedUsers(dataGen.DBConn, personID1)

		toBeFollowed := followedAssignment[seq1] - len(alreadyFollowedByPersons)

		toBeFollowedByPersons = append(toBeFollowedByPersons, 
			GetSeqsByPersonIDs(users, alreadyFollowedByPersons)...)

		// log.Println("Check user:", seq1)
		// log.Println("Total users to follow this user:", followedAssignment[seq1])
		// log.Println("already followed by", alreadyFollowedByPersons)
		// log.Println("To be followed by:", toBeFollowed)
		
		if toBeFollowed > followedAssignment[seq1] {
			log.Fatal("cannot happend1!!!!")
		}

		haveTried := make(map[int]bool)

		for _, alreadyFollowedByPersonsID := range alreadyFollowedByPersons {
			
			haveTried[alreadyFollowedByPersonsID] = true

		}

		for n := 0; n < toBeFollowed; n++ {

			for {

				if len(haveTried) >= len(users) - 1 {

					log.Println("Cannot find more users to follow this user!!")
					log.Println("Total users to follow this user:", followedAssignment[seq1])
					log.Println("Have been followed by:", n + followedAssignment[seq1] - toBeFollowed)
					
					// ableToBeFollowed = false
					
					log.Fatal("cannot happend2!!!!")
					break

				}

				seq2 := RandomNonnegativeIntWithUpperBound(len(users))
				
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

					// Note: in the multiple-thread data generation, 
					// person1 could be followed by person2 twice
					aspect_idx := helper.RandomNumber(0, len(user1.Aspects) - 1)

					datagen.FollowUser(dataGen.DBConn, 
						personID2, personID1, user1.Aspects[aspect_idx])

					toBeFollowedByPersons = append(toBeFollowedByPersons, seq2)
					
					haveTried[seq2] = true

					break

				}

			}

			// if !ableToBeFollowed {
			// 	break
			// }
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

// We use the user popularity score to generate how many followers a user has.
// For some of those followers, the user also follows back.
// Note: RECIPROCAL_FOLLOW_PERCENTAGE cannot guarantee that
// this user can follow this percentage of followers,
// because maybe most of those users have already fully followed by other users.
// Use the query: "select count(*) from contacts where sharing = true;" 
// to get the total number of follows
// Use the query: "select count(*) from contacts where sharing = true and receiving = true;" 
// to get the a number and divide this number by 2 to get the total pairs of friends
// The exact number could be more than FOLLOW_NUM 
// because a user could be followed by the same other user twice
// due to multiple-thread data generation.
func (dataGen *DataGen) genFollowsController(users []DUser) {

	followedAssignment := AssignDataToUsersByUserScores(
		dataGen.UserPopularityScores, FOLLOW_NUM)
	
	log.Println("Followed assignment to users:", followedAssignment)
	log.Println("Total followed:", GetSumOfIntSlice(followedAssignment))
	
	var wg sync.WaitGroup

	wg.Add(THREAD_NUM)

	userSeqStart := 0
	
	userSeqStep := len(users) / THREAD_NUM

	for i := 0; i < THREAD_NUM; i++ {
		
		if i != THREAD_NUM - 1 {

			// Start included, end (start + step) not included
			go dataGen.genFollows(&wg, userSeqStart, userSeqStart + userSeqStep, 
				followedAssignment, users)

		} else {

			// Start included, end (start + step) not included
			go dataGen.genFollows(&wg, userSeqStart, len(users), 
				followedAssignment, users)
		
		}

		userSeqStart += userSeqStep
	}

	wg.Wait()

}

func (dataGen *DataGen) genPosts(wg *sync.WaitGroup, 
	res1 chan<- map[int]float64, res2 chan<- int, 
	userSeqStart, userSeqEnd, postSeqStart int,
	users []DUser, postAssignment []int, 
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
	
	res1 <- postScores

	res2 <- imageNums

}

// The number of posts of users is proportional to the popularity of users.
// We also randomly assign images to the posts proportionally to the popularity of posts.
// The scores assigned to posts are in pareto distributiuon.
// so it is more likely that popular users will have popular posts because they have more posts
func (dataGen *DataGen) genPostsController(users []DUser) map[int]float64 {

	postAssignment := AssignDataToUsersByUserScores(dataGen.UserPopularityScores, POST_NUM)
	totalPosts := GetSumOfIntSlice(postAssignment)
	
	log.Println("Posts assignments to users:", postAssignment)
	log.Println("Total posts:", totalPosts)

	seqNum := MakeRange(0, totalPosts - 1)
	seqScores := AssignParetoDistributionScoresToDataReturnSlice(len(seqNum))
	imageNumsOfSeq := RandomNumWithProbGenerator(seqScores, IMAGE_NUM)
	
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
			go dataGen.genPosts(&wg, channel1, channel2, 
				userSeqStart, userSeqStart + userSeqStep, 
				postSeqStart, users, postAssignment, imageNumsOfSeq, seqScores)

		} else {

			// Start included, end (start + step) not included
			go dataGen.genPosts(&wg, channel1, channel2, 
				userSeqStart, len(users), 
				postSeqStart, users, postAssignment, imageNumsOfSeq, seqScores)
		
		}

		postSeqStart = CalculateNextPostSeqStart(
			postSeqStart, userSeqStart, userSeqStart + userSeqStep, postAssignment)

		userSeqStart += userSeqStep
	}

	wg.Wait()

	close(channel1)

	close(channel2)

	for res1 := range channel1 {

		postScores = MergeTwoMaps(postScores, res1) 

	}

	for res2 := range channel2 {

		imageNums += res2

	}

	log.Println("Total images:", imageNums)

	return postScores

}

// Only for test
func (dataGen *DataGen) prepareTest() ([]DUser, map[int]float64) {

	var users []DUser
	users1 := datagen.GetAllUsersWithAspectsOrderByID(dataGen.DBConn)

	for _, user1 := range users1 {

		var user DUser
		user.User_ID, user.Person_ID, user.Aspects = user1.User_ID, user1.Person_ID, user1.Aspects
		users = append(users, user)

	}

	return users, AssignParetoDistributionScoresToData(
		datagen.GetAllPostIDs(dataGen.DBConn))

}

func (dataGen *DataGen) genComments(wg *sync.WaitGroup, 
	userSeqStart, userSeqEnd int, commentAssignment []int, 
	users []DUser, postScores map[int]float64) {
	
	defer wg.Done()

	for seq1 := userSeqStart; seq1 < userSeqEnd; seq1++ {

		user1 := users[seq1]

		// log.Println("Check user:", seq1)
		
		var posts []*Post
		var scores []float64
		commentNum := commentAssignment[seq1]

		// log.Println("Comment number:", commentNum)
		
		personID := user1.Person_ID
		
		// Even if a user is followed by the same user (U1) twice, it does not influcence much
		// This can result in the posts of U1 being added twice, so the posts of U1 will be
		// commented twice more than expected.
		totalUsers := datagen.GetFollowingUsers(dataGen.DBConn, personID)

		// log.Println(user1)
		// log.Println(totalUsers)

		totalUsers = append(totalUsers, personID)

		for _, user2 := range totalUsers {

			posts1 := datagen.GetPostsForUser(dataGen.DBConn, user2)

			for _, post1 := range posts1 {

				post := new(Post)
				post.ID = post1.ID
				post.Author = post1.Author
				post.Score = postScores[post1.ID]

				posts = append(posts, post)
				scores = append(scores, post.Score)

			}
		}
		
		commentNumsOfPosts := RandomNumWithProbGenerator(scores, commentNum)

		for seq2, post := range posts {

			for i := 0; i < commentNumsOfPosts[seq2]; i++ {

				datagen.NewComment(dataGen.DBConn, post.ID, personID, post.Author)
			}
		}
	}
}

// We randomly assign comments to posts proportionally to the popularity of posts of friends, 
// including posts by the commenter.
// Mutiple threads do not cause much influence to comments generation
func (dataGen *DataGen) genCommentsController(users []DUser, postScores map[int]float64) {
	
	commentAssignment := AssignDataToUsersByUserScores(
		dataGen.UserCommentScores, COMMENT_NUM)

	log.Println("Comments assignments to users:", commentAssignment)
	log.Println("Total comments:", GetSumOfIntSlice(commentAssignment))

	var wg sync.WaitGroup

	wg.Add(THREAD_NUM)

	userSeqStart := 0
	
	userSeqStep := len(users) / THREAD_NUM

	for i := 0; i < THREAD_NUM; i++ {
		
		if i != THREAD_NUM - 1 {

			// Start included, end (start + step) not included
			go dataGen.genComments(&wg, userSeqStart, userSeqStart + userSeqStep, 
				commentAssignment, users, postScores)

		} else {

			// Start included, end (start + step) not included
			go dataGen.genComments(&wg, userSeqStart, len(users), 
				commentAssignment, users, postScores)
		
		}

		userSeqStart += userSeqStep
	}

	wg.Wait()
		
}

func (dataGen *DataGen) genLikes(wg *sync.WaitGroup, 
	userSeqStart, userSeqEnd int, likeAssignment []int, 
	users []DUser, postScores map[int]float64, res chan<- int) {
		
	defer wg.Done()

	totalLikeNum := 0

	for seq1 := userSeqStart; seq1 < userSeqEnd; seq1++ {

		user1 := users[seq1]

		// log.Println("Check user:", seq1)
		
		var posts []*Post
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

				post := new(Post)
				post.ID = post1.ID
				post.Author = post1.Author
				post.Score = postScores[post1.ID]
				posts = append(posts, post)
				scores = append(scores, post.Score)

			}
		}
		
		likeNumsOfPosts := RandomNumWithProbGenerator(scores, likeNum)

		// log.Println(likeNumsOfPosts)
		
		for seq2, post := range posts {

			if _, ok := likeNumsOfPosts[seq2]; ok {

				datagen.NewLike(dataGen.DBConn, post.ID, personID, post.Author)
				totalLikeNum += 1
			}
		}
	}

	res <- totalLikeNum

}

// We randomly assign likes to posts proportionally to the popularity of posts of friends, 
// including posts by the liker.
// The difference between generating comments and likes is that
// a user make several comments on the same post, but can only like once on that post.
// Mutiple threads do not cause much influence to likes generation
func (dataGen *DataGen) genLikesController(users []DUser, postScores map[int]float64) {

	likeAssignment := AssignDataToUsersByUserScores(
		dataGen.UserLikeScores, LIKE_NUM)

	log.Println("Likes assignments to users:", likeAssignment)
	log.Println("Total likes based on assignments:", GetSumOfIntSlice(likeAssignment))

	channel := make(chan int, THREAD_NUM)

	totalLikeNum := 0

	var wg sync.WaitGroup

	wg.Add(THREAD_NUM)

	userSeqStart := 0
	
	userSeqStep := len(users) / THREAD_NUM

	for i := 0; i < THREAD_NUM; i++ {
		
		if i != THREAD_NUM - 1 {

			// Start included, end (start + step) not included
			go dataGen.genLikes(&wg, userSeqStart, userSeqStart + userSeqStep, 
				likeAssignment, users, postScores, channel)

		} else {

			// Start included, end (start + step) not included
			go dataGen.genLikes(&wg, userSeqStart, len(users), 
				likeAssignment, users, postScores, channel)
		
		}

		userSeqStart += userSeqStep
	}

	wg.Wait()

	close(channel)

	for res := range channel {

		totalLikeNum += res

	}

	log.Println("In reality, the num of total likes is:", totalLikeNum)

}

func (dataGen *DataGen) genConversationsAndMessages(wg *sync.WaitGroup, 
	userSeqStart, userSeqEnd int, messageAssignment []int, 
	users []DUser, res chan<- int) {
		
	defer wg.Done()

	conversationNum := 0

	for seq1 := userSeqStart; seq1 < userSeqEnd; seq1++ {

		user1 := users[seq1]

		// oneUserConversationNum := 0
		
		personID := user1.Person_ID
		messageNum := messageAssignment[seq1]

		// There could be cases in which the user has no friend
		friends := datagen.GetRealFriendsOfUser(dataGen.DBConn, personID)
		friendCloseIndex := AssignParetoDistributionScoresToDataReturnSlice(len(friends))
		
		// log.Println(friends)
		// log.Println(friendCloseIndex)
		// log.Println(messageNum)
		
		messageNumsOfConversations := RandomNumWithProbGenerator(friendCloseIndex, messageNum)
		// log.Println(messageNumsOfConversations)

		for seq2, messageNum := range messageNumsOfConversations {

			exists, conv_id := datagen.CheckConversationBetweenTwoUsers(dataGen.DBConn, 
				personID, friends[seq2])
			
			if exists {

				for i := 0; i < messageNum; i++ {

					datagen.NewMessage(dataGen.DBConn, personID, conv_id)

				}

			} else {
				
				// Given the multiple-thread data generator, there could be 
				// two conversations between two same users
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

	res <- conversationNum

}

// Pareto-distributed message scores determine the number of messages each user should have.
// Friendships have pareto-distributed closeness indexes. 
// We randomly assign messages to users (or conversations) proportionally 
// to the closeness indexes.
// Due to the multiple threads data generation, two users could 
// talk with each other sharing the same conversation.
// Also note that with number of users generated increasing, there could be cases
// in which a user has no friend, so the messages allocated to this user cannot be sent.
// Therefore the actual message number is lower than the calculated total number 
// according to the messageAssignment
func (dataGen *DataGen) genConversationsAndMessagesController(users []DUser) {

	messageAssignment := AssignDataToUsersByUserScores(
		dataGen.UserMessageScores, MESSAGE_NUM)

	log.Println("Messages assignments to users:", messageAssignment)
	log.Println("Total messages:", GetSumOfIntSlice(messageAssignment))

	channel := make(chan int, THREAD_NUM)

	conversationNum := 0

	var wg sync.WaitGroup

	wg.Add(THREAD_NUM)

	userSeqStart := 0
	
	userSeqStep := len(users) / THREAD_NUM

	for i := 0; i < THREAD_NUM; i++ {
		
		if i != THREAD_NUM - 1 {

			// Start included, end (start + step) not included
			go dataGen.genConversationsAndMessages(&wg, userSeqStart, userSeqStart + userSeqStep, 
				messageAssignment, users, channel)

		} else {

			// Start included, end (start + step) not included
			go dataGen.genConversationsAndMessages(&wg, userSeqStart, len(users), 
				messageAssignment, users, channel)
		
		}

		userSeqStart += userSeqStep
	}

	wg.Wait()

	close(channel)

	for res := range channel {

		conversationNum += res

	}

	log.Println("Total conversations:", conversationNum)

}


func (dataGen *DataGen) GenDataDiaspora() {
	
	startTime := time.Now()

	log.Println("--------- Start of Data Generation ---------")

	// users, postScores := dataGen.prepareTest()

	users := dataGen.genUsersController()

	// After getting the exact user number, the data generator needs
	// to initialize UserPopularityScores, UserCommentScores, etc.
	dataGen.InitializeWithUserNum(len(users))

	postScores := dataGen.genPostsController(users)

	dataGen.genFollowsController(users)

	dataGen.genCommentsController(users, postScores)

	dataGen.genLikesController(users, postScores)
	
	dataGen.genConversationsAndMessagesController(users)

	log.Println("--------- End of Data Generation ---------")

	endTime := time.Now()

	log.Println("Time used: ", endTime.Sub(startTime))

}
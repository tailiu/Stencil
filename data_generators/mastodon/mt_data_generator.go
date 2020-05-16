/*
 * Data Generator for Mastodon
 */

package main

import (
	"mastodon/mtDataGen"
	"mastodon/database"
	"mastodon/auxiliary"
    "fmt"
	"log"
	"database/sql"
)

const (
	host     = "10.230.12.86"
	port     = 5432
	user     = "cow"
	password = "123456"
	dbname   = "mastodon"
)

var address = fmt.Sprintf("host=%s port=%d user=%s "+" password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

func getAllAccountIDs(dbConn *sql.DB) *[]int{
	sql1 := fmt.Sprintf("SELECT id FROM accounts ORDER BY id;")
	rows, err := dbConn.Query(sql1)
	if err != nil {
		log.Fatal(err)
	}
	var accountIDs []int
	var accountID int
	for rows.Next() {
		if err := rows.Scan(&accountID); err != nil {
			log.Fatal(err)
		}
		accountIDs = append(accountIDs, accountID)
	}
	return &accountIDs
}

func getAllPostIDs(dbConn *sql.DB) *[]int {
	sql1 := fmt.Sprintf("SELECT id FROM statuses;")
	rows, err := dbConn.Query(sql1)
	if err != nil {
		log.Fatal(err)
	}
	var postIDs []int
	var postID int
	for rows.Next() {
		if err := rows.Scan(&postID); err != nil {
			log.Fatal(err)
		}
		postIDs = append(postIDs, postID)
	}
	// for _, postID := range postIDs {
	// 	fmt.Println(postID)
	// }
	return &postIDs
} 

/************************************** Create Users ***********************************************************/

func createUsers(dbConn *sql.DB, num int, c chan int) {
	for i := 0; i < num; i++ {
		newUser := mtDataGen.User{
			auxiliary.RandStrSeq(4) + "@" + auxiliary.RandStrSeq(5),
			auxiliary.RandStrSeq(10),
			auxiliary.RandStrSeq(25)}
		mtDataGen.Signup(dbConn, &newUser)

		newUserLog := fmt.Sprintf("Create %d users with username='%s', email='%s', password='%s'",
			i + 1, newUser.Username, newUser.Email, newUser.Password)
		fmt.Println(newUserLog)
	}
	c <- num
}

/*************************************************************************************************/

/**************************************** Create Public Posts *********************************************************/

func createPostsThread(dbConn *sql.DB, slicedAccountIDs []int) {
	var haveMedia bool
	var numOfPosts int
	for _, accountID := range slicedAccountIDs {
		// Each user publishes 0 - 350 statuses
		numOfPosts = auxiliary.RandomNonnegativeIntWithUpperBound(350)
		for j := 0; j < numOfPosts; j++ {
			if j % 25 == 0 {
				haveMedia = true
			} else {
				haveMedia = false
			}
			var mentionedAccounts []int
			mtDataGen.PublishStatus(dbConn, accountID, auxiliary.RandStrSeq(50), haveMedia, 0, mentionedAccounts)
			newStatusLog := fmt.Sprintf("User %d creates a post with Media %t",
				accountID, haveMedia)
			fmt.Println(newStatusLog)
		}	
	}
}

func createPublicPostsController(accountIDs *[]int) {
	j := 0
	var dbConn *sql.DB
	for i := 0; i < len(*accountIDs); i++ {
		// There are about 100,000 accounts, so there will be 100000/1000 = 100 threads
		if i != 0 && i % 1000 == 0 {
			dbConn = database.ConnectToDB(address)
			go createPostsThread(dbConn, (*accountIDs)[j:i])
			j = i
		}
	}
}

/*************************************************************************************************/

/***************************************** Follow Friends ********************************************************/

func followFriendsThread(dbConn *sql.DB, slicedAccountIDs []int, accountIDs *[]int) {
	var targetAccountID int
	for _, accountID := range slicedAccountIDs {
		// Each user has 0 - 240 friends
		numOfFriends := auxiliary.RandomNonnegativeIntWithUpperBound(240)
		for j := 0; j < numOfFriends; j++ {
			targetAccountID = (*accountIDs)[auxiliary.RandomNonnegativeIntWithUpperBound(len(*accountIDs))]
			mtDataGen.Follow(dbConn, accountID, targetAccountID)
			newFriendLog := fmt.Sprintf("User %d follows user %d",
				accountID, targetAccountID)
			fmt.Println(newFriendLog)
		}
	}
}

func followFriendsController(dbConn *sql.DB, accountIDs *[]int) {
	accountNum, j := len(*accountIDs), 0
	for i := 0; i < accountNum; i++ {
		// There are about 100,000 accounts, so there will be 100000/500 = 200 threads
		if i != 0 && i % 500 == 0 {
			go followFriendsThread(dbConn, (*accountIDs)[j:i], accountIDs)
			j = i
		}
	}
}

/*************************************************************************************************/

/******************************* Create Message Groups ****************************************************/

func createOneDirectMessage(dbConn *sql.DB, accountID int, accountIDs *[]int) ([]int, int) {
	var haveMedia bool
	if auxiliary.RandomNonnegativeIntWithUpperBound(10) == 0 {
		haveMedia = true
	} else {
		haveMedia = false
	}
	var mentionedAccounts []int
	// Each message group initially has 2 - 11 users 
	numOfMentionedUsers := auxiliary.RandomNonnegativeIntWithUpperBound(10) + 1
	for i := 0; i < numOfMentionedUsers; i++ {
		mentionedAccounts = append(mentionedAccounts, 
			(*accountIDs)[auxiliary.RandomNonnegativeIntWithUpperBound(len(*accountIDs))])
	}
	messageID := mtDataGen.PublishStatus(dbConn, accountID, auxiliary.RandStrSeq(50), haveMedia, 3, mentionedAccounts)
	if messageID != -1 {
		newMessageGroupLog := fmt.Sprintf("User %d creates a message group with %d users",
			accountID, len(mentionedAccounts) + 1)
		fmt.Println(newMessageGroupLog)
	}
	return mentionedAccounts, messageID
}

func createOneReplyUsingMentionedAccount(dbConn *sql.DB, replyToStatusID int, newLayer *[]int, allAccounts []int, visibility int) {
	accountID := allAccounts[auxiliary.RandomNonnegativeIntWithUpperBound(len(allAccounts))]
	for i, targetAccountID := range allAccounts {
		if targetAccountID == accountID {
			allAccounts = append(allAccounts[:i], allAccounts[i+1:]...)
			break
		}
	}
	result := mtDataGen.ReplyToStatus(dbConn, accountID, auxiliary.RandStrSeq(50), replyToStatusID, visibility, allAccounts)
	if result != -1 {
		newReply := fmt.Sprintf("Create a message %d as a reply to a message %d with visibility %d", result, replyToStatusID, visibility)
		fmt.Println(newReply)
		*newLayer = append(*newLayer, result)
	} else {
		FailedToCreateNewReply := fmt.Sprintf("Failed to create a reply to a message %d with visibility %d", replyToStatusID, visibility)
		fmt.Println(FailedToCreateNewReply)
	}
}

func createDirectMessageGroupsThread(dbConn *sql.DB, slicedAccountIDs []int, accountIDs *[]int) {
	var numOfMessageGroups int
	for _, accountID := range slicedAccountIDs {
		// Each user creates 0 - 10 message groups
		numOfMessageGroups = auxiliary.RandomNonnegativeIntWithUpperBound(10)
		for j := 0; j < numOfMessageGroups; j++ {
			var newLayer []int
			var allAccounts []int
			// Each message group has 0 - 15 messages
			numOfLayers := auxiliary.RandomNonnegativeIntWithUpperBound(15)
			for k := 0; k < numOfLayers; k++ {
				if k == 0 {
					// Create one direct message with 1 - 10 mentioned users
					mentionedAccounts, messageID := createOneDirectMessage(dbConn, accountID, accountIDs)
					if messageID == -1 {
						break
					}
					allAccounts = append(allAccounts, mentionedAccounts...)
					allAccounts = append(allAccounts, accountID)
					newLayer = append(newLayer, messageID)
				} else {
					// Copy newLayer to layers and set newLayer to nil
					length := len(newLayer)
					layers := make([]int, length)
					copy(layers, newLayer)
					newLayer = nil
					for _, replyToStatusID := range layers {
						// Each message has 0 - 2 replies
						numOfReplies := auxiliary.RandomNonnegativeIntWithUpperBound(2)
						for i := 0; i < numOfReplies; i++ {
							createOneReplyUsingMentionedAccount(dbConn, replyToStatusID, &newLayer, allAccounts, 3)
						}
					}	
				}
				if len(newLayer) == 0 {
					break
				}
			}
		}	
	}
}

func createDirectMessagesController(dbConn *sql.DB, accountIDs *[]int) {
	j := 0
	// var dbConn *sql.DB
	for i := 0; i < len(*accountIDs); i++ {
		// There are about 100,000 accounts, so there will be 100000/500 = 200 threads
		if i != 0 && i % 500 == 0 {
			// dbConn = database.ConnectToDB(address)
			go createDirectMessageGroupsThread(dbConn, (*accountIDs)[j:i], accountIDs)
			j = i
		}
	}
}

/*************************************************************************************************/

/******************************* Create Replies To Posts ****************************************************/

func createOneReply(dbConn *sql.DB, replyToStatusID int, newLayer *[]int, accountIDs *[]int, visibility int, mentionedAccounts []int) {
	accountID := (*accountIDs)[auxiliary.RandomNonnegativeIntWithUpperBound(len(*accountIDs))]
	result := mtDataGen.ReplyToStatus(dbConn, accountID, auxiliary.RandStrSeq(50), replyToStatusID, 0, mentionedAccounts)
	if result != -1 {
		newReply := fmt.Sprintf("Create a new reply to %d", replyToStatusID)
		fmt.Println(newReply)
		*newLayer = append(*newLayer, result)
	}
}

func createRepliesToPostsThread(dbConn *sql.DB, slicedPostIDs []int, accountIDs *[]int) {
	var mentionedAccounts []int
	for _, postID := range slicedPostIDs {
		var newLayer []int
		// Replay layers are 0 - 10 levels
		numOfLayers := auxiliary.RandomNonnegativeIntWithUpperBound(10)
		for j := 0; j < numOfLayers; j++ {
			if j == 0 {
				// Each post has 0 - 3 replies
				numOfReplies := auxiliary.RandomNonnegativeIntWithUpperBound(3)
				for i := 0; i < numOfReplies; i++ {
					createOneReply(dbConn, postID, &newLayer, accountIDs, 0, mentionedAccounts)
				}
			} else {
				length := len(newLayer)
				layers := make([]int, length)
				copy(layers, newLayer)
				newLayer = nil
				// This is a mistake...
				for _, replyToStatusID := range layers {
					createOneReply(dbConn, replyToStatusID, &newLayer, accountIDs, 0, mentionedAccounts)
				}	
			}
			if len(newLayer) == 0 {
				break
			}
		}
	}
}

func createRepliesToPostsController(dbConn *sql.DB, accountIDs *[]int, postIDs *[]int) {
	postNum, j := len(*postIDs), 0
	// var dbConn *sql.DB
	for i := 0; i < postNum; i++ {
		// There are about 20,000,000 posts, so there will be 20000000/200000 = 100 threads
		if i != 0 && i % 200000 == 0 {
			// dbConn = database.ConnectToDB(address)
			go createRepliesToPostsThread(dbConn, (*postIDs)[j:i], accountIDs)
			j = i
		}
	}
}

/*************************************************************************************************/

func main() {
	dbConn := database.ConnectToDB(address)

	// threadNum := 50
	// c := make(chan int, threadNum)
	// for i := 0; i < threadNum; i++ {
	// 	go createUsers(dbConn, 50000/threadNum, c)
	// }
	// for i := range c {
	// 	fmt.Println(i)
	// }
	// createPublicPostsController(getAllAccountIDs(dbConn))
	// followFriendsController(dbConn, getAllAccountIDs(dbConn))
	createRepliesToPostsController(dbConn, getAllAccountIDs(dbConn), getAllPostIDs(dbConn))
	// createDirectMessagesController(dbConn, getAllAccountIDs(dbConn))

	for {
		fmt.Scanln()
	}

	// var mentionedAccounts []int
	// mentionedAccounts = append(mentionedAccounts, 3243277, 3258353)
	// mtDataGen.PublishStatus(dbConn, 925840864, "five media COOOOL", true, 0, mentionedAccounts)

	// favourite(dbConn, 925840864, 1389362391)

	// unfavourite(dbConn, 925840864, 1389362391)

	// signup(dbConn, "tai@nyu.edu", "zaincow", "cowcow")

	// follow(dbConn, 1217195077, 1042906640)

	// unfollow(dbConn, 1217195077, 1042906640)
	
	// replyToStatus(dbConn, 829522384, "a reply", 2042450516)

	// reblog(dbConn, 735104489, 614615112)
}

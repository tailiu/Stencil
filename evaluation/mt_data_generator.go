/*
 * Data Generator for Mastodon
 */

package main

import (
	"evaluation/functions"
	"evaluation/database"
	"evaluation/auxiliary"
    "fmt"
	"log"
	"database/sql"
)

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

func createUsers(dbConn *sql.DB, num int, c chan int) {
	for i := 0; i < num; i++ {
		newUser := functions.User{
			auxiliary.RandStrSeq(4) + "@" + auxiliary.RandStrSeq(5),
			auxiliary.RandStrSeq(10),
			auxiliary.RandStrSeq(25)}
		functions.Signup(dbConn, &newUser)

		newUserLog := fmt.Sprintf("Create %d users with username='%s', email='%s', password='%s'",
			i + 1, newUser.Username, newUser.Email, newUser.Password)
		fmt.Println(newUserLog)
	}
	c <- num
}

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
			functions.PublishStatus(dbConn, accountID, auxiliary.RandStrSeq(50), haveMedia, 0, mentionedAccounts)
			newStatusLog := fmt.Sprintf("User %d creates a post with Media %t",
				accountID, haveMedia)
			fmt.Println(newStatusLog)
		}	
	}
}

func createPostsController(dbConn *sql.DB, accountIDs *[]int) {
	j := 0
	for i := 0; i < len(*accountIDs); i++ {
		// There are about 100,000 accounts, so there will be 100000/500 = 200 threads
		if i != 0 && i % 1000 == 0 {
			go createPostsThread(dbConn, (*accountIDs)[j:i])
			j = i
		}
	}
}

func followFriendsThread(dbConn *sql.DB, slicedAccountIDs []int, accountIDs *[]int) {
	var targetAccountID int
	for _, accountID := range slicedAccountIDs {
		// Each user has 0 - 240 friends
		numOfFriends := auxiliary.RandomNonnegativeIntWithUpperBound(240)
		for j := 0; j < numOfFriends; j++ {
			targetAccountID = (*accountIDs)[auxiliary.RandomNonnegativeIntWithUpperBound(len(*accountIDs))]
			functions.Follow(dbConn, accountID, targetAccountID)
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

func main() {
	address := "postgresql://root@10.230.12.75:26257/mastodon?sslmode=disable"
	dbConn := database.ConnectToDB(address)

	// threadNum := 1
	// c := make(chan int, threadNum)
	// for i := 0; i < threadNum; i++ {
		// go createUsers(dbConn, 50000/threadNum, c)
	// }
	// for i := range c {
	// 	fmt.Println(i)
	// }

	createPostsController(dbConn, getAllAccountIDs(dbConn))
	// followFriendsController(dbConn, getAllAccountIDs(dbConn))
	
	for {
		fmt.Scanln()
	}

	// var mentionedAccounts []int
	// mentionedAccounts = append(mentionedAccounts, 3243277, 3258353)
	// functions.PublishStatus(dbConn, 925840864, "five media COOOOL", true, 0, mentionedAccounts)

	// favourite(dbConn, 925840864, 1389362391)

	// unfavourite(dbConn, 925840864, 1389362391)

	// signup(dbConn, "tai@nyu.edu", "zaincow", "cowcow")

	// follow(dbConn, 1217195077, 1042906640)

	// unfollow(dbConn, 1217195077, 1042906640)
	
	// replyToStatus(dbConn, 829522384, "a reply", 2042450516)

	// reblog(dbConn, 735104489, 614615112)
}

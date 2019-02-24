/*
 * Data Generator for Mastodon
 */

package main

import (
	"evaluation/functions"
	"evaluation/database"
	"evaluation/auxiliary"
    "fmt"
	// "log"
	"database/sql"
)

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

func main() {
	address := "postgresql://root@10.230.12.75:26257/mastodon?sslmode=disable"
	dbConn := database.ConnectToDB(address)

	c := make(chan int, 2)

	go createUsers(dbConn, 20, c)
	go createUsers(dbConn, 20, c)

	for i := range c {
		fmt.Println(i)
	}

	// go createUsers(dbConn, 5)
	// functions.PublishStatus(dbConn, 925840864, "five media", true)

	// favourite(dbConn, 925840864, 1389362391)

	// unfavourite(dbConn, 925840864, 1389362391)

	// signup(dbConn, "tai@nyu.edu", "zaincow", "cowcow")

	// follow(dbConn, 1217195077, 1042906640)

	// unfollow(dbConn, 1217195077, 1042906640)
	
	// replyToStatus(dbConn, 829522384, "a reply", 2042450516)

	// reblog(dbConn, 735104489, 614615112)
}

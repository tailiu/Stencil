package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"twitter/config"
	"twitter/datagen"
	"twitter/db"
)

func createNewUsers(dbConn *sql.DB, num int) {
	for i := 0; i < num; i++ {
		uid := datagen.CreateNewUser(dbConn)
		fmt.Printf("User: %4d/%4d | id : %s \n", i, num, uid)
	}
}

func createNewConversations(dbConn *sql.DB, num int, userID string) {
	for i := 0; i < num; i++ {
		id := datagen.CreateNewConversation(dbConn, userID)
		fmt.Printf("Conversation: %4d/%4d | id : %s \n", i, num, id)
	}
}

func createNewPosts(dbConn *sql.DB, num int, profileID string) {
	for i := 0; i < num; i++ {
		id := datagen.CreateNewPost(dbConn, profileID)
		fmt.Printf("Post: %4d/%4d | id : %s \n", i, num, id)
	}
}

func createNewPostMedia(dbConn *sql.DB, num int, profileID string) {
	for i := 0; i < num; i++ {
		id := datagen.CreateTweetWithPhoto(dbConn, profileID)
		fmt.Printf("Media: %4d/%4d | id : %s \n", i, num, id)
	}
}

func createNewReshares(dbConn *sql.DB, num int, profileID, postID string) {
	for i := 0; i < num; i++ {
		id := datagen.CreateNewReshare(dbConn, profileID, postID)
		fmt.Printf("Reshare: %4d/%4d | id : %s \n", i, num, id)
	}
}

func createNewComments(dbConn *sql.DB, num int, profileID, postID string) {
	for i := 0; i < num; i++ {
		id := datagen.CreateNewComment(dbConn, profileID, postID)
		fmt.Printf("Comment: %4d/%4d | id : %s \n", i, num, id)
	}
}

func createNewLikes(dbConn *sql.DB, num int, profileID, postID string) {
	for i := 0; i < num; i++ {
		id := datagen.CreateNewLike(dbConn, postID, profileID)
		fmt.Printf("Like: %4d/%4d | id : %s \n", i, num, id)
	}
}

func createNewMessages(dbConn *sql.DB, num int, profileID, conversationID string) {
	for i := 0; i < num; i++ {
		id := datagen.CreateNewMessage(dbConn, profileID, conversationID)
		fmt.Printf("Message: %4d/%4d | id : %s \n", i, num, id)
	}
}

func createNewFollows(dbConn *sql.DB, fromProfileID, toProfileID string) {

	id := datagen.CreateNewFollow(dbConn, fromProfileID, toProfileID)
	fmt.Printf("Subscription | id : %s \n", id)
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Enter cmd args.")
	}

	dbConn := db.GetDBConn(config.DB_NAME)

	arg := os.Args[1]

	personID1, personID2 := "264509681", "1088267574"
	postID, conversationID := "977512205", "460128162"

	switch arg {
	case "users":
		fmt.Println("Creating New Users!")
		createNewUsers(dbConn, 10)
	case "conversations":
		fmt.Println("Creating New Conversations!")
		createNewConversations(dbConn, 10, personID1)
	case "posts":
		fmt.Println("Creating New posts!")
		createNewPosts(dbConn, 10, personID1)
	case "comments":
		fmt.Println("Creating New comments!")
		createNewComments(dbConn, 10, personID2, postID)
	case "reshares":
		fmt.Println("Creating New reshares!")
		createNewReshares(dbConn, 10, personID2, postID)
	case "media":
		fmt.Println("Creating New media!")
		createNewPostMedia(dbConn, 10, personID2)
	case "message":
		fmt.Println("Creating New messages!")
		createNewMessages(dbConn, 10, personID1, conversationID)
		createNewMessages(dbConn, 10, personID2, conversationID)
	case "follow":
		fmt.Println("Creating New subscriptions!")
		createNewFollows(dbConn, personID1, personID2)
	case "likes":
		fmt.Println("Creating New likes!")
		createNewLikes(dbConn, 10, personID2, postID)
	}
}

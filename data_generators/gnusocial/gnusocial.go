package main

import (
	"database/sql"
	"fmt"
	"gnusocial/config"
	"gnusocial/datagen"
	"gnusocial/db"
	"log"
	"os"
	"stencil/helper"
)

func createNewUsers(dbConn *sql.DB, num int) {
	for i := 0; i < num; i++ {
		uid := datagen.CreateNewUser(dbConn)
		fmt.Printf("User: %4d/%4d | id : %s \n", i, num, uid)
	}
}

func createNewConversations(dbConn *sql.DB, num int) {
	for i := 0; i < num; i++ {
		id := datagen.CreateNewConversation(dbConn)
		fmt.Printf("Conversation: %4d/%4d | id : %s \n", i, num, id)
	}
}

func createNewPosts(dbConn *sql.DB, num int, profileID, conversationID string) {
	for i := 0; i < num; i++ {
		id := datagen.CreateNewPost(dbConn, conversationID, profileID)
		fmt.Printf("Post: %4d/%4d | id : %s \n", i, num, id)
	}
}

func createNewPostMedia(dbConn *sql.DB, num int, profileID, conversationID, postID string) {
	for i := 0; i < num; i++ {
		id := datagen.CreateMediaForPost(dbConn, postID, profileID)
		fmt.Printf("Media: %4d/%4d | id : %s \n", i, num, id)
	}
}

func createNewReshares(dbConn *sql.DB, num int, profileID, conversationID, postID string) {
	for i := 0; i < num; i++ {
		id := datagen.CreateNewReshare(dbConn, conversationID, profileID, postID)
		fmt.Printf("Reshare: %4d/%4d | id : %s \n", i, num, id)
	}
}

func createNewComments(dbConn *sql.DB, num int, profileID, conversationID, postID string) {
	for i := 0; i < num; i++ {
		id := datagen.CreateNewComment(dbConn, conversationID, profileID, postID)
		fmt.Printf("Comment: %4d/%4d | id : %s \n", i, num, id)
	}
}

func createNewLikes(dbConn *sql.DB, num int, profileID, postID string) {
	for i := 0; i < num; i++ {
		id := datagen.CreateNewLike(dbConn, postID, profileID)
		fmt.Printf("Like: %4d/%4d | id : %s \n", i, num, id)
	}
}

func createNewMessages(dbConn *sql.DB, num int, fromProfileID, toProfileID string) {
	for i := 0; i < num; i++ {
		id := datagen.CreateNewMessage(dbConn, fromProfileID, toProfileID)
		fmt.Printf("Message: %4d/%4d | id : %s \n", i, num, id)
	}
}

func createNewSubscriptions(dbConn *sql.DB, fromProfileID, toProfileID string) {

	id := datagen.CreateNewSubscription(dbConn, fromProfileID, toProfileID)
	fmt.Printf("Subscription | id : %s \n", id)
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Enter cmd args.")
	}

	helper.Init()
	dbConn := db.GetDBConn(config.DB_NAME)

	arg := os.Args[1]

	personID1, personID2 := "1157670463", "700549107"
	postID, conversationID := "1648058497", "195160395"

	switch arg {
	case "users":
		fmt.Println("Creating New Users!")
		createNewUsers(dbConn, 10)
	case "conversations":
		fmt.Println("Creating New Conversations!")
		createNewConversations(dbConn, 10)
	case "posts":
		fmt.Println("Creating New posts!")
		createNewPosts(dbConn, 10, personID1, conversationID)
	case "comments":
		fmt.Println("Creating New comments!")
		createNewComments(dbConn, 10, personID1, conversationID, postID)
	case "reshares":
		fmt.Println("Creating New reshares!")
		createNewReshares(dbConn, 10, personID1, conversationID, postID)
	case "media":
		fmt.Println("Creating New media!")
		createNewPostMedia(dbConn, 10, personID1, conversationID, postID)
	case "message":
		fmt.Println("Creating New messages!")
		createNewMessages(dbConn, 10, personID1, personID2)
	case "subscription":
		fmt.Println("Creating New subscriptions!")
		createNewSubscriptions(dbConn, personID1, personID2)
	case "likes":
		fmt.Println("Creating New likes!")
		createNewLikes(dbConn, 10, personID1, postID)
	}
}

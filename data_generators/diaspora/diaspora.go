package main

import (
	"database/sql"
	"diaspora/config"
	"diaspora/datagen"
	"diaspora/db"
	"diaspora/helper"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func WaitForAWhile() {
	time.Sleep(10 * time.Minute)
}

func createNewUsers(dbConn *sql.DB, num, thread int) {
	for i := 0; i < num; i++ {
		uid, _, _ := datagen.NewUser(dbConn)
		fmt.Println(fmt.Sprintf("Thread: %3d, User: %4d/%4d | uid : %d", thread, i, num, uid))
	}
}

func createNewPostsForUsers(dbConn *sql.DB, users []*datagen.User, thread_num int) {

	for uidx, user := range users {
		num_of_posts := helper.RandomNumber(0, 500)
		for i := 0; i <= num_of_posts; i++ {
			log.Println(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Posts %3d/%3d", thread_num, uidx, len(users), i, num_of_posts))
			datagen.NewPost(dbConn, user.User_ID, user.Person_ID, user.Aspects)
		}
	}
}

func createNewMentionsForUsers(dbConn *sql.DB, users []*datagen.User) {

}

func makeUsersFriends(dbConn *sql.DB, users []*datagen.User, thread_num int) {

	for uidx, user := range users {
		indices := rand.Perm(len(users))
		num_of_friends := helper.RandomNumber(60, 80)
		for i := 0; i <= num_of_friends; i++ {
			if i >= len(indices) {
				break
			}
			if index := indices[i]; index != uidx {
				log.Println(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Friends %3d/%3d", thread_num, uidx, len(users), i, num_of_friends))
				user2 := users[index]
				aspect_idx := helper.RandomNumber(0, len(user.Aspects)-1)
				datagen.FollowUser(dbConn, user.Person_ID, user2.Person_ID, user.Aspects[aspect_idx])
				if helper.RandomNumber(1, 50)%2 == 0 {
					aspect_idx := helper.RandomNumber(0, len(user2.Aspects)-1)
					datagen.FollowUser(dbConn, user2.Person_ID, user.Person_ID, user2.Aspects[aspect_idx])
				}
			}
		}
	}
}

func makeUsersTalk(dbConn *sql.DB, users []*datagen.User, thread_num int) {
	fmt.Println(fmt.Sprintf("Thread # %d, reporting for duty! ", thread_num))
	num_users := len(users)
	for uidx, user := range users {
		friends_of_user := datagen.GetFriendsOfUser(dbConn, user.User_ID)
		num_frnds := len(friends_of_user)
		for fidx, friend := range friends_of_user {
			if helper.RandomNumber(1, 100)%3 == 0 {
				conversation_id, err := datagen.NewConversation(dbConn, user.Person_ID, friend.Person_ID)
				if err == nil && conversation_id != -1 {
					num_of_msgs := helper.RandomNumber(50, 1000)
					for i := 0; i <= num_of_msgs; {
						fmt.Println(fmt.Sprintf("{THREAD: %3d} [Users %4d/%4d | Frnds %3d/%3d] | Msg # %3d/%3d | Conversation: %d", thread_num, uidx, num_users, fidx, num_frnds, i, num_of_msgs, conversation_id))
						if helper.RandomNumber(1, 100)%2 == 0 {
							_, err := datagen.NewMessage(dbConn, friend.Person_ID, conversation_id)
							if err != nil {
								log.Println(err)
								WaitForAWhile()
							}
							i++
						}
						if helper.RandomNumber(1, 100)%4 != 0 {
							_, err := datagen.NewMessage(dbConn, user.Person_ID, conversation_id)
							if err != nil {
								log.Println(err)
								WaitForAWhile()
							}
							i++
						}
					}
					datagen.UpdateConversation(dbConn, conversation_id)
				}
			}
		}
	}

}

// comments, like and reshare friends posts
func interactWithPosts(dbConn *sql.DB, users []*datagen.User, thread_num int) {

	num_users := len(users)
	for uidx, user := range users {
		friends_of_user := datagen.GetFriendsOfUser(dbConn, user.User_ID)
		posts := datagen.GetPostsForUser(dbConn, user.Person_ID)
		if len(posts) <= 0 || len(friends_of_user) <= 0 {
			continue
		}
		num_frnds, num_posts := len(friends_of_user), helper.RandomNumber(0, len(posts))
		for pidx, post := range posts[0:num_posts] {
			for fidx, friend := range friends_of_user {
				fmt.Println(fmt.Sprintf("{THREAD: %3d} Users %3d/%4d | Frnds %3d/%4d | Posts %4d/%4d ", thread_num, uidx, num_users, fidx, num_frnds, pidx, num_posts))

				if helper.RandomNumber(1, 100)%4 == 0 { // 25%, Friend Likes The Post
					datagen.NewLike(dbConn, post.ID, friend.Person_ID, user.Person_ID)
				}
				if helper.RandomNumber(1, 100)%10 == 0 { // 10%, Friend Reshares The Post
					datagen.NewReshare(dbConn, *post, friend.Person_ID)
				}
				if helper.RandomNumber(1, 100)%5 == 0 { // 20%, Comments On The Post
					loopcount := helper.RandomNumber(1, 10)
					for l := 0; l < loopcount; l++ {
						if helper.RandomNumber(1, 100)%2 == 0 { // Friend Comments
							datagen.NewComment(dbConn, post.ID, friend.Person_ID, user.Person_ID)
						}
						if helper.RandomNumber(1, 100)%2 == 0 { // Owner Comments
							datagen.NewComment(dbConn, post.ID, user.Person_ID, user.Person_ID)
						}
					}
				}
			}
		}
	}
}

func runinteractWithPosts() {
	dbConn := db.GetDBConn(config.APP_NAME)
	users := datagen.GetAllUsersWithAspects(dbConn)
	num_users := len(users)
	inc := 500
	for i, j := 0, inc; i < num_users && j < num_users; i, j = j+1, j+inc {
		thread_num := j / inc
		go interactWithPosts(db.GetDBConn(config.APP_NAME), users[i:j], thread_num)
	}

	for {
		fmt.Scanln()
	}
}

func runMakeUsersTalk() {

	dbConn := db.GetDBConn(config.APP_NAME)
	users := datagen.GetAllUsersWithAspects(dbConn)

	// makeUsersTalk(dbConn, users, 0)

	num_users := len(users)
	inc := 500

	for thread_num, i, j := 0, 0, inc; i < num_users && j < num_users; i, j, thread_num = j+1, j+inc, thread_num+1 {
		go makeUsersTalk(dbConn, users[i:j], thread_num)
	}

	for {
		fmt.Scanln()
	}
}

func runCreateNewUsers() {
	// dbConn := db.GetDBConn(config.APP_NAME)
	for i := 0; i < 25; i++ {
		go createNewUsers(db.GetDBConn(config.APP_NAME), 40, i)
	}
	for {
		fmt.Scanln()
	}
}

func runMakeUsersFriends() {
	dbConn := db.GetDBConn(config.APP_NAME)
	users := datagen.GetAllUsersWithAspects(dbConn)
	num_users := len(users)
	inc := 100
	// makeUsersFriends(dbConn, users, 0)

	for thread_num, i, j := 0, 0, inc; i < num_users && j < num_users; i, j, thread_num = j+1, j+inc, thread_num+1 {
		go makeUsersFriends(dbConn, users[i:j], thread_num)
	}

	for {
		fmt.Scanln()
	}
}

func runCreateNewPosts() {
	dbConn := db.GetDBConn(config.APP_NAME)
	users := datagen.GetAllUsersWithAspects(dbConn)
	num_users := len(users)
	inc := 500
	for i, j := 0, inc; i < num_users && j < num_users; i, j = j+1, j+inc {
		thread_num := j / inc
		go createNewPostsForUsers(dbConn, users[i:j], thread_num)
	}

	for {
		fmt.Scanln()
	}
}

func main() {

	arg := os.Args[1]

	switch arg {
	case "posts":
		fmt.Println("Creating New Posts!")
		runCreateNewPosts()
	case "comments":
		fmt.Println("Interacting With Posts!")
		runinteractWithPosts()
	case "likes":
		fmt.Println("Interacting With Posts!")
		runinteractWithPosts()
	case "reshares":
		fmt.Println("Interacting With Posts!")
		runinteractWithPosts()
	case "messages":
		fmt.Println("Creating New Messages!")
		runMakeUsersTalk()
	case "friends":
		fmt.Println("Making People Friends!")
		runMakeUsersFriends()
	case "newusers":
		fmt.Println("Creating New Users!")
		runCreateNewUsers()
	}
}

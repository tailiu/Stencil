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
		num_of_friends := helper.RandomNumber(0, 300)
		for i := 0; i <= num_of_friends; i++ {
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
	num_users := len(users)
	for uidx, user := range users {
		friends_of_user := datagen.GetFriendsOfUser(dbConn, user.User_ID)
		num_frnds := len(friends_of_user)
		for fidx, friend := range friends_of_user {
			if helper.RandomNumber(1, 100)%3 == 0 {
				conversation_id, conversation_visibilities_id, err := datagen.NewConversation(dbConn, user.Person_ID, friend.Person_ID)
				if err == nil {
					num_of_msgs := helper.RandomNumber(50, 1000)
					for i := 0; i <= num_of_msgs; {
						fmt.Println(fmt.Sprintf("{THREAD: %d} [Users %d/%d | Frnds %d/%d] | UserID: %d, FriendID:%d | Conversation: %d | Msg # %d ", thread_num, uidx, num_users, fidx, num_frnds, user.Person_ID, friend.Person_ID, conversation_id, i))
						if helper.RandomNumber(1, 100)%2 == 0 {
							_, err := datagen.NewMessage(dbConn, friend.Person_ID, conversation_id, conversation_visibilities_id)
							if err != nil {
								log.Println(err)
								WaitForAWhile()
							}
							i++
						}
						if helper.RandomNumber(1, 100)%4 != 0 {
							_, err := datagen.NewMessage(dbConn, user.Person_ID, conversation_id, conversation_visibilities_id)
							if err != nil {
								log.Println(err)
								WaitForAWhile()
							}
							i++
						}
					}
					datagen.UpdateConversation(dbConn, conversation_id, conversation_visibilities_id)
				} else {
					log.Println(err)
					WaitForAWhile()
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
		num_frnds, num_posts := len(friends_of_user), helper.RandomNumber(0, len(posts))
		for pidx, post := range posts[0:num_posts] {
			for fidx, friend := range friends_of_user {
				thereIsAnError := true
				fmt.Println(fmt.Sprintf("{THREAD: %3d} Users %3d/%4d | Frnds %3d/%4d | Posts %4d/%4d ", thread_num, uidx, num_users, fidx, num_frnds, pidx, num_posts))

				if helper.RandomNumber(1, 100)%4 == 0 { // 25%, Friend Likes The Post
					if _, err := datagen.NewLike(dbConn, post.ID, friend.Person_ID, user.Person_ID); err != nil {
						// log.Println(err)
						// WaitForAWhile()
					} else {
						thereIsAnError = false
					}
				}
				if helper.RandomNumber(1, 100)%10 == 0 { // 10%, Friend Reshares The Post
					if _, err := datagen.NewReshare(dbConn, *post, friend.Person_ID); err != nil {
						// log.Println(err)
						// WaitForAWhile()
					} else {
						thereIsAnError = false
					}
				}
				if helper.RandomNumber(1, 100)%5 == 0 { // 20%, Comments On The Post
					loop_count := helper.RandomNumber(1, 10)
					for l := 0; l < loop_count; l++ {
						if helper.RandomNumber(1, 100)%2 == 0 { // Friend Comments
							if _, err := datagen.NewComment(dbConn, post.ID, friend.Person_ID, user.Person_ID); err != nil {
								log.Println(err)
								// WaitForAWhile()
							} else {
								thereIsAnError = false
							}
						}
						if helper.RandomNumber(1, 100)%2 == 0 { // Owner Comments
							if _, err := datagen.NewComment(dbConn, post.ID, user.Person_ID, user.Person_ID); err != nil {
								log.Println(err)
								// WaitForAWhile()
							} else {
								thereIsAnError = false
							}
						}
					}
				}
				if thereIsAnError {
					// log.Println("Waiting...")
					// WaitForAWhile()
				}
			}
		}
	}
}

func runinteractWithPosts() {
	dbConn := db.GetDBConn(config.APP_NAME)
	users := datagen.GetAllUsersWithAspects(dbConn)
	num_users := len(users)
	inc := 100
	for i, j := 0, inc; i < num_users && j < num_users; i, j = j+1, j+inc {
		thread_num := j / inc
		go interactWithPosts(db.GetDBConn(config.APP_NAME), users[i:j], thread_num)
		// fmt.Println(i, j, num_users, thread_num)
	}

	for {
		fmt.Scanln()
	}
}

func runMakeUsersTalk() {
	dbConn := db.GetDBConn(config.APP_NAME)
	users := datagen.GetAllUsersWithAspectsExcept(dbConn, "author_id", "conversations")
	num_users := len(users)
	inc := 100
	for i, j := 0, inc; i < num_users && j < num_users; i, j = j+1, j+inc {
		thread_num := j / inc
		go makeUsersTalk(db.GetDBConn(config.APP_NAME), users[i:j], thread_num)
		// fmt.Println(i, j, num_users, thread_num)
	}

	for {
		fmt.Scanln()
	}
}

func runCreateNewUsers() {
	// dbConn := db.GetDBConn(config.APP_NAME)
	for i := 0; i < 100; i++ {
		go createNewUsers(db.GetDBConn(config.APP_NAME), 500, i)
	}
	for {
		fmt.Scanln()
	}
}

func runMakeUsersFriends() {
	dbConn := db.GetDBConn(config.APP_NAME)
	users := datagen.GetAllUsersWithAspects(dbConn)
	num_users := len(users)
	inc := 1000
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

	// dbConn := db.GetDBConn(config.APP_NAME)
	// runCreateNewUsers()
	runMakeUsersFriends()
	// runCreateNewPosts()
	// runinteractWithPosts()
	// runMakeUsersTalk()
	// users := datagen.GetAllUsersWithAspects(dbConn)

}

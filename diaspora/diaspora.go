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

func createNewUsers(dbConn *sql.DB, num int) {
	for i := 0; i < num; i++ {
		datagen.NewUser(dbConn)
	}
}

func blockAUser() {

}

func createNewPostsForUsers(dbConn *sql.DB, users []*datagen.User) {

	for _, user := range users {
		num_of_posts := helper.RandomNumber(0, 500)
		for i := 0; i <= num_of_posts; i++ {
			datagen.NewPost(dbConn, user.User_ID, user.Person_ID, user.Aspects)
		}
	}
}

func createNewCommentsOnPosts(dbConn *sql.DB, users []*datagen.User) {

	// for _, user := range users {
	// 	friends_of_user := datagen.GetFriendsOfUser(dbConn, user.User_ID)
	// 	num_of_posts := helper.RandomNumber(0, 500)
	// 	for i := 0; i <= num_of_posts; i++ {
	// 		datagen.NewPost(dbConn, user.User_ID, user.Person_ID, user.Aspects)
	// 	}
	// }
}

func createNewMentionsForUsers(dbConn *sql.DB, users []*datagen.User) {

}

func makeUsersFriends(dbConn *sql.DB, users []*datagen.User) {

	for uidx, user := range users {
		helper.Init()
		indices := rand.Perm(len(users))
		num_of_friends := helper.RandomNumber(0, 500)
		for i := 0; i <= num_of_friends; i++ {
			if index := indices[i]; index == uidx {
				continue
			} else {
				user2 := users[index]
				aspect_idx := helper.RandomNumber(0, len(user.Aspects)-1)
				datagen.FollowUser(dbConn, user.User_ID, user2.User_ID, user.Aspects[aspect_idx])
				log.Println(fmt.Sprintf("User: %d added User: %d to Aspect: %d", user.User_ID, user2.User_ID, user.Aspects[aspect_idx]))
				if helper.RandomNumber(1, 50)%2 == 0 {
					aspect_idx := helper.RandomNumber(0, len(user2.Aspects)-1)
					datagen.FollowUser(dbConn, user2.User_ID, user.User_ID, user2.Aspects[aspect_idx])
					log.Println(fmt.Sprintf("User: %d added User: %d to Aspect: %d", user2.User_ID, user.User_ID, user2.Aspects[aspect_idx]))
				}
			}
		}
	}
}

func WaitForAWhile() {
	time.Sleep(10 * time.Minute)
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

func interactWithPosts(dbConn *sql.DB, users []*datagen.User, thread_num int) {
	// comments, like and reshare friends posts

	num_users := len(users)
	for uidx, user := range users {
		friends_of_user := datagen.GetFriendsOfUser(dbConn, user.User_ID)
		posts := datagen.GetPostsForUser(dbConn, user.Person_ID)
		num_frnds, num_posts := len(friends_of_user), helper.RandomNumber(0, len(posts))
		for pidx, post := range posts[0:num_posts] {
			for fidx, friend := range friends_of_user {
				fmt.Println(fmt.Sprintf("{THREAD: %3d} Users %3d/%4d | Frnds %3d/%4d | Posts %4d/%4d ", thread_num, uidx, num_users, fidx, num_frnds, pidx, num_posts))

				if helper.RandomNumber(1, 100)%4 == 0 { // 25%, Friend Likes The Post
					if _, err := datagen.NewLike(dbConn, post.ID, friend.Person_ID, user.Person_ID); err != nil {
						// log.Println(err)
						// WaitForAWhile()
					}
				}
				if helper.RandomNumber(1, 100)%10 == 0 { // 10%, Friend Reshares The Post
					if _, err := datagen.NewReshare(dbConn, *post, friend.Person_ID); err != nil {
						// log.Println(err)
						// WaitForAWhile()
					}
				}
				if helper.RandomNumber(1, 100)%5 == 0 { // 20%, Comments On The Post
					loop_count := helper.RandomNumber(1, 10)
					for l := 0; l < loop_count; l++ {
						if helper.RandomNumber(1, 100)%2 == 0 { // Friend Comments
							if _, err := datagen.NewComment(dbConn, post.ID, friend.Person_ID, user.Person_ID); err != nil {
								log.Println(err)
								WaitForAWhile()
							}
						}
						if helper.RandomNumber(1, 100)%2 == 0 { // Owner Comments
							if _, err := datagen.NewComment(dbConn, post.ID, user.Person_ID, user.Person_ID); err != nil {
								log.Println(err)
								WaitForAWhile()
							}
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
	inc := 500
	for i, j := 0, inc; i < num_users && j < num_users; i, j = j+1, j+inc {
		thread_num := j / inc
		go makeUsersTalk(db.GetDBConn(config.APP_NAME), users[i:j], thread_num)
		// fmt.Println(i, j, num_users, thread_num)
	}

	for {
		fmt.Scanln()
	}
}

func main() {

	// dbConn := db.GetDBConn(config.APP_NAME)
	runinteractWithPosts()
	// runMakeUsersTalk()
	// users := datagen.GetAllUsersWithAspects(dbConn)

}

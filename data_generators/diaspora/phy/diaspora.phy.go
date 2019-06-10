package main

import (
	"diaspora/config"
	"diaspora/datagen/phy"
	"diaspora/helper"
	"fmt"
	"log"
	"math/rand"
	"stencil/qr"
	"time"
)

var QR = qr.NewQR(config.APP_NAME, config.APP_ID)

func WaitForAWhile() {
	time.Sleep(10 * time.Minute)
}

func createNewUsers(num, thread int) {
	for i := 0; i < num; i++ {
		uid, _, _ := phy.NewUser(QR)
		fmt.Println(fmt.Sprintf("Thread: %3d, User: %4d/%4d | uid : %d", thread, i, num, uid))
	}
}

func createNewPostsForUsers(users []*phy.User, thread_num int) {

	for uidx, user := range users {
		num_of_posts := helper.RandomNumber(0, 500)
		for i := 0; i <= num_of_posts; i++ {
			log.Println(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Posts %3d/%3d", thread_num, uidx, len(users), i, num_of_posts))
			phy.NewPost(QR, user.User_ID, user.Person_ID, user.Aspects)
		}
	}
}

func createNewMentionsForUsers(users []*phy.User) {

}

func makeUsersFriends(users []*phy.User, thread_num int) {

	for uidx, user := range users {
		indices := rand.Perm(len(users))
		num_of_friends := helper.RandomNumber(0, 300)
		for i := 0; i <= num_of_friends; i++ {
			if i >= len(indices) {
				break
			}
			if index := indices[i]; index != uidx {
				log.Println(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Friends %3d/%3d", thread_num, uidx, len(users), i, num_of_friends))
				user2 := users[index]
				aspect_idx := helper.RandomNumber(0, len(user.Aspects)-1)
				phy.FollowUser(QR, user.Person_ID, user2.Person_ID, user.Aspects[aspect_idx])
				if helper.RandomNumber(1, 50)%2 == 0 {
					aspect_idx := helper.RandomNumber(0, len(user2.Aspects)-1)
					phy.FollowUser(QR, user2.Person_ID, user.Person_ID, user2.Aspects[aspect_idx])
				}
			}
		}
	}
}

func makeUsersTalk(users []*phy.User, thread_num int) {
	fmt.Println(fmt.Sprintf("Thread # %d, reporting for duty! ", thread_num))
	num_users := len(users)
	for uidx, user := range users {
		friends_of_user := phy.GetFriendsOfUser(QR, user.User_ID)
		num_frnds := len(friends_of_user)
		for fidx, friend := range friends_of_user {
			if helper.RandomNumber(1, 100)%3 == 0 {
				conversation_id, err := phy.NewConversation(QR, user.Person_ID, friend.Person_ID)
				if err == nil && conversation_id != -1 {
					num_of_msgs := helper.RandomNumber(50, 1000)
					for i := 0; i <= num_of_msgs; {
						fmt.Println(fmt.Sprintf("{THREAD: %3d} [Users %4d/%4d | Frnds %3d/%3d] | Msg # %3d/%3d | Conversation: %d", thread_num, uidx, num_users, fidx, num_frnds, i, num_of_msgs, conversation_id))
						if helper.RandomNumber(1, 100)%2 == 0 {
							_, err := phy.NewMessage(QR, friend.Person_ID, conversation_id)
							if err != nil {
								log.Println(err)
								WaitForAWhile()
							}
							i++
						}
						if helper.RandomNumber(1, 100)%4 != 0 {
							_, err := phy.NewMessage(QR, user.Person_ID, conversation_id)
							if err != nil {
								log.Println(err)
								WaitForAWhile()
							}
							i++
						}
					}
					phy.UpdateConversation(QR, conversation_id)
				}
			}
		}
	}

}

// comments, like and reshare friends posts
func interactWithPosts(users []*phy.User, thread_num int) {

	num_users := len(users)
	for uidx, user := range users {
		friends_of_user := phy.GetFriendsOfUser(QR, user.User_ID)
		posts := phy.GetPostsForUser(QR, user.Person_ID)
		if len(posts) <= 0 || len(friends_of_user) <= 0 {
			continue
		}
		num_frnds, num_posts := len(friends_of_user), helper.RandomNumber(0, len(posts))
		for pidx, post := range posts[0:num_posts] {
			for fidx, friend := range friends_of_user {
				fmt.Println(fmt.Sprintf("{THREAD: %3d} Users %3d/%4d | Frnds %3d/%4d | Posts %4d/%4d ", thread_num, uidx, num_users, fidx, num_frnds, pidx, num_posts))

				if helper.RandomNumber(1, 100)%4 == 0 { // 25%, Friend Likes The Post
					phy.NewLike(QR, post.ID, friend.Person_ID, user.Person_ID)
				}
				if helper.RandomNumber(1, 100)%10 == 0 { // 10%, Friend Reshares The Post
					phy.NewReshare(QR, *post, friend.Person_ID)
				}
				if helper.RandomNumber(1, 100)%5 == 0 { // 20%, Comments On The Post
					loopcount := helper.RandomNumber(1, 10)
					for l := 0; l < loopcount; l++ {
						if helper.RandomNumber(1, 100)%2 == 0 { // Friend Comments
							phy.NewComment(QR, post.ID, friend.Person_ID, user.Person_ID)
						}
						if helper.RandomNumber(1, 100)%2 == 0 { // Owner Comments
							phy.NewComment(QR, post.ID, user.Person_ID, user.Person_ID)
						}
					}
				}
			}
		}
	}
}

func runinteractWithPosts() {
	users := phy.GetAllUsersWithAspects(QR)
	num_users := len(users)
	inc := 500
	for i, j := 0, inc; i < num_users && j < num_users; i, j = j+1, j+inc {
		thread_num := j / inc
		go interactWithPosts(users[i:j], thread_num)
	}

	for {
		fmt.Scanln()
	}
}

func runMakeUsersTalk() {

	users := phy.GetAllUsersWithAspects(QR)

	num_users := len(users)
	inc := 500

	for thread_num, i, j := 0, 0, inc; i < num_users && j < num_users; i, j, thread_num = j+1, j+inc, thread_num+1 {
		go makeUsersTalk(users[i:j], thread_num)
	}

	for {
		fmt.Scanln()
	}
}

func runCreateNewUsers() {
	// dbConn := db.GetDBConn(config.APP_NAME)
	for i := 0; i < 100; i++ {
		go createNewUsers(500, i)
	}
	for {
		fmt.Scanln()
	}
}

func runMakeUsersFriends() {

	users := phy.GetAllUsersWithAspects(QR)
	num_users := len(users)
	inc := 500

	for thread_num, i, j := 0, 0, inc; i < num_users && j < num_users; i, j, thread_num = j+1, j+inc, thread_num+1 {
		go makeUsersFriends(users[i:j], thread_num)
	}

	for {
		fmt.Scanln()
	}
}

func runCreateNewPosts() {

	users := phy.GetAllUsersWithAspects(QR)
	num_users := len(users)
	inc := 500
	for i, j := 0, inc; i < num_users && j < num_users; i, j = j+1, j+inc {
		thread_num := j / inc
		go createNewPostsForUsers(users[i:j], thread_num)
	}

	for {
		fmt.Scanln()
	}
}

func main() {

	createNewUsers(1, 0)

	// arg := os.Args[1]

	// switch arg {
	// case "posts":
	// 	fmt.Println("Creating New Posts!")
	// 	runCreateNewPosts()
	// case "comments":
	// 	fmt.Println("Interacting With Posts!")
	// 	runinteractWithPosts()
	// case "messages":
	// 	fmt.Println("Creating New Messages!")
	// 	runMakeUsersTalk()
	// case "friends":
	// 	fmt.Println("Making People Friends!")
	// 	runMakeUsersFriends()
	// case "newusers":
	// 	fmt.Println("Creating New Users!")
	// 	runCreateNewUsers()
	// }
}

package main

import (
	"database/sql"
	"diaspora/config"
	"diaspora/dist"
	"diaspora/datagen/phy"
	"diaspora/helper"
	"fmt"
	"log"
	"math/rand"
	"os"
	"stencil/db"
	"stencil/qr"
	"strings"
	"time"
	"sync"
)

var QR = qr.NewQR(config.APP_NAME, config.APP_ID)

func WaitForAWhile() {
	time.Sleep(10 * time.Minute)
}

func createNewUsers(dbConn *sql.DB, num, thread int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < num; i++ {
		fmt.Println(fmt.Sprintf("Thread: %3d, User: %4d/%4d", thread, i, num))
		phy.NewUser(QR, dbConn)
	}
}

func createNewPostsForUsers(dbConn *sql.DB, users []*phy.User, thread_num int) {

	for uidx, user := range users {
		num_of_posts := helper.RandomNumber(100, 5000)
		for i := 0; i <= num_of_posts; i++ {
			log.Println(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Posts %3d/%3d", thread_num, uidx, len(users), i, num_of_posts))
			phy.NewPost(QR, dbConn, user.User_ID, user.Person_ID, user.Aspects)
		}
	}
}

func runGenerateFriendships() {

	var wg sync.WaitGroup

	numThreads := 1
	dbConn := db.GetDBConn(config.STENCILDB)
	users := phy.GetAllUsersWithAspects(QR, dbConn)
	numUsers := len(users)
	friendlist := dist.AssignFriendsToUsers(numUsers)
	dist.VerifyFriendsDistribution(friendlist)
	// log.Fatal(phy.GetFriendsDistribution(QR, dbConn))
	time.Sleep(time.Duration(10) * time.Second)
	inc := int(numUsers / numThreads)
	// log.Fatal("here")
	for t, i, j := 1, 0, inc; t <= numThreads; t, i, j = t+1, j+1, j+inc {
		wg.Add(1)
		go createFriends(dbConn, users, friendlist, t, i, j, &wg)
	}

	wg.Wait()
	time.Sleep(time.Duration(2) * time.Second)
	dist.VerifyFriendsDistribution(friendlist)
	fmt.Println(phy.GetFriendsDistribution(QR, dbConn))
}

func createFriends(dbConn *sql.DB, users []*phy.User, friendlist map[int][]int, thread_num, start, end int, wg *sync.WaitGroup) {
	defer wg.Done()

	uidx := 0
	// for uid, friends := range friendlist {
	for i := start; i < end; i++ {
		user1 := users[i]
		for fid, friend := range friendlist[i] {
			user2 := users[friend]
			aspect_idx := helper.RandomNumber(0, len(user1.Aspects)-1)
			phy.FollowUser(QR, dbConn, user1.Person_ID, user2.Person_ID, user1.Aspects[aspect_idx])
			phy.FollowUser(QR, dbConn, user2.Person_ID, user1.Person_ID, user2.Aspects[aspect_idx])
			log.Println(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Friends %3d/%3d", thread_num, uidx, end-start, fid, len(friendlist[i])))
		}
		uidx++
	}
}

func makeUsersFriends(dbConn *sql.DB, users []*phy.User, thread_num int) {

	for uidx, user := range users {
		indices := rand.Perm(len(users))
		num_of_friends := helper.RandomNumber(10, len(users))
		for i := 0; i <= num_of_friends && i < len(indices); i++ {
			if index := indices[i]; index != uidx {
				log.Println(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Friends %3d/%3d", thread_num, uidx, len(users), i, num_of_friends))
				user2 := users[index]
				aspect_idx := helper.RandomNumber(0, len(user.Aspects)-1)
				phy.FollowUser(QR, dbConn, user.Person_ID, user2.Person_ID, user.Aspects[aspect_idx])
				if helper.RandomNumber(1, 50)%2 == 0 {
					aspect_idx := helper.RandomNumber(0, len(user2.Aspects)-1)
					phy.FollowUser(QR, dbConn, user2.Person_ID, user.Person_ID, user2.Aspects[aspect_idx])
				}
			}
		}
	}
}

func makeUsersTalk(dbConn *sql.DB, users []*phy.User, thread_num int) {

	num_users := len(users)
	for uidx, user := range users {
		fmt.Println(fmt.Sprintf("{THREAD: %3d} Fetching friends ", thread_num))
		friends_of_user := phy.GetFriendsOfUser(QR, user.Person_ID)
		num_frnds := len(friends_of_user)
		fmt.Println(fmt.Sprintf("{THREAD: %3d} UID/PID: %d/%d | Friends Fetched: %d ", thread_num, user.User_ID, user.Person_ID, len(friends_of_user)))
		for fidx, friend := range friends_of_user {
			if helper.RandomNumber(1, 100)%3 == 0 {
				conversation_id, err := phy.NewConversation(QR, dbConn, fmt.Sprint(user.Person_ID), fmt.Sprint(friend.Person_ID))
				if err == nil && conversation_id != -1 {
					num_of_msgs := helper.RandomNumber(50, 200)
					for i := 0; i <= num_of_msgs; {
						fmt.Println(fmt.Sprintf("{THREAD: %3d} [Users %4d/%4d | Frnds %3d/%3d] | Msg # %3d/%3d | Conversation: %d", thread_num, uidx, num_users, fidx, num_frnds, i, num_of_msgs, conversation_id))
						if helper.RandomNumber(1, 100)%2 == 0 {
							err := phy.NewMessage(QR, dbConn, fmt.Sprint(friend.Person_ID), fmt.Sprint(conversation_id))
							if err != nil {
								log.Println(err)
								// WaitForAWhile()
							}
							i++
						}
						if helper.RandomNumber(1, 100)%4 != 0 {
							err := phy.NewMessage(QR, dbConn, fmt.Sprint(user.Person_ID), fmt.Sprint(conversation_id))
							if err != nil {
								log.Println(err)
								// WaitForAWhile()
							}
							i++
						}
					}
					// phy.UpdateConversation(QR, dbConn, conversation_id)
				}
			}
		}
	}

}

// comments, like and reshare friends posts
func interactWithPosts(dbConn *sql.DB, users []*phy.User, thread_num int, itype string) {

	num_users := len(users)
	for uidx, user := range users {
		fmt.Println(fmt.Sprintf("{THREAD: %3d} Fetching friends ", thread_num))
		friends_of_user := phy.GetFriendsOfUser(QR, user.Person_ID)
		fmt.Println(fmt.Sprintf("{THREAD: %3d} Fetching posts ", thread_num))
		posts := phy.GetPostsForUser(QR, dbConn, user.Person_ID)
		fmt.Println(fmt.Sprintf("{THREAD: %3d} UID/PID: %d/%d | Friends Fetched: %d | Posts Fetched: %d", thread_num, user.User_ID, user.Person_ID, len(friends_of_user), len(posts)))
		if len(posts) <= 0 || len(friends_of_user) <= 0 {
			continue
		}
		num_frnds, num_posts := len(friends_of_user), helper.RandomNumber(0, len(posts))
		for pidx, post := range posts[0:num_posts] {
			for fidx, friend := range friends_of_user {
				fmt.Println(fmt.Sprintf("{THREAD: %3d} Users %3d/%4d | Frnds %3d/%4d | Posts %4d/%4d ", thread_num, uidx, num_users, fidx, num_frnds, pidx, num_posts))

				if strings.EqualFold(itype, "likes") && helper.RandomNumber(1, 100)%4 == 0 { // 25%, Friend Likes The Post
					phy.NewLike(QR, dbConn, post.ID, friend.Person_ID, user.Person_ID)
				}
				if strings.EqualFold(itype, "reshares") && helper.RandomNumber(1, 100)%10 == 0 { // 10%, Friend Reshares The Post
					phy.NewReshare(QR, dbConn, *post, friend.Person_ID)
				}
				if strings.EqualFold(itype, "comments") && helper.RandomNumber(1, 100)%5 == 0 { // 20%, Comments On The Post
					loopcount := helper.RandomNumber(1, 10)
					for l := 0; l < loopcount; l++ {
						if helper.RandomNumber(1, 100)%2 == 0 { // Friend Comments
							phy.NewComment(QR, dbConn, post.ID, friend.Person_ID, user.Person_ID)
						}
						if helper.RandomNumber(1, 100)%2 == 0 { // Owner Comments
							phy.NewComment(QR, dbConn, post.ID, user.Person_ID, user.Person_ID)
						}
					}
				}
			}
		}
	}
}

func runinteractWithPosts(itype string) {
	dbConn := db.GetDBConn(config.STENCILDB)
	users := phy.GetAllUsersWithAspects(QR, dbConn)
	num_users := len(users)
	inc := 50
	for i, j := 0, inc; i < num_users && j < num_users; i, j = j+1, j+inc {
		
		thread_num := j / inc
		go interactWithPosts(dbConn, users[i:j], thread_num, itype)
		// break
	}

	for {
		fmt.Scanln()
	}
}

func runMakeUsersTalk() {

	dbConn := db.GetDBConn(config.STENCILDB)
	users := phy.GetAllUsersWithAspects(QR, dbConn)

	num_users := len(users)
	inc := 50

	for thread_num, i, j := 0, 0, inc; i < num_users && j < num_users; i, j, thread_num = j+1, j+inc, thread_num+1 {
		
		go makeUsersTalk(dbConn, users[i:j], thread_num)
		// break
	}

	for {
		fmt.Scanln()
	}
}

func runCreateNewUsers() {
	var wg sync.WaitGroup

	dbConn := db.GetDBConn(config.STENCILDB)
	thread_num := 1
	for i := 0; i < thread_num; i++ {
		wg.Add(1)
		go createNewUsers(dbConn, 1000, i, &wg)
	}
	wg.Wait()
}

func runMakeUsersFriends() {

	dbConn := db.GetDBConn(config.STENCILDB)
	for count := 1; count <= 1; count++ {
		users := phy.GetAllUsersWithAspects(QR, dbConn)
		num_users := len(users)
		inc := 50 * count
		for thread_num, i, j := 0, 0, inc; i < num_users && j < num_users; i, j, thread_num = j+1, j+inc, thread_num+1 {
			
			go makeUsersFriends(dbConn, users[i:j], thread_num)
			// break
		}
	}
	for {
		fmt.Scanln()
	}
}

func runCreateNewPosts() {
	dbConn := db.GetDBConn(config.STENCILDB)

	users := phy.GetAllUsersWithAspects(QR, dbConn)
	num_users := len(users)
	inc := 10
	for i, j := 0, inc; i < num_users && j < num_users; i, j = j+1, j+inc {
		thread_num := j / inc
		
		go createNewPostsForUsers(dbConn, users[i:j], thread_num)
		// break
	}
	for {
		fmt.Scanln()
	}
}

func main() {

	helper.Init()

	arg := os.Args[1]

	switch arg {
	case "posts":
		fmt.Println("Creating New Posts!")
		// runCreateNewPosts("posts", 8030)
	case "comments":
		fmt.Println("Interacting With Posts!")
		// runCreateNewPosts("comments", 13970)
	case "likes":
		fmt.Println("Interacting With Posts!")
		// runCreateNewPosts("likes", 85680)
	case "reshares":
		fmt.Println("Interacting With Posts!")
		// runCreateNewPosts("reshares", 1550)
	case "messages":
		fmt.Println("Creating New Messages!")
		// runCreateNewPosts("messages", 4015)
		// runMakeUsersTalk()
	case "makefriends":
		fmt.Println("Making People Friends!")
		// runMakeUsersFriends()
	case "genfriends":
		fmt.Println("Generating Friendships!")
		runGenerateFriendships()
	case "newusers":
		fmt.Println("Creating New Users!")
		runCreateNewUsers()
	case "addfriendforpersoninrange":
		person_id := os.Args[2]
		// lower_bound := os.Args[3]
		// upper_bound := os.Args[4]
		// limit := os.Args[5]
		fmt.Println("Add New Friends For:", person_id)
		// AddFriendsForUser(person_id, lower_bound, upper_bound, limit)
	}
}

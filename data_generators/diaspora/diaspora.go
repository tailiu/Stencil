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
	"strconv"
	"strings"
	"sync"
	"time"
)

type Range struct {
	lower_bound, upper_bound int
	pshare                   float32
}

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

func checkFriendsDistribution(friendcount map[int]int) {

	o, i, j, k, l, m, n := 0, 0, 0, 0, 0, 0, 0

	for _, fnum := range friendcount {
		if fnum < 1 {
			o++
		} else if fnum >= 1 && fnum < 100 {
			i++
		} else if fnum >= 100 && fnum < 200 {
			j++
		} else if fnum >= 200 && fnum < 300 {
			k++
		} else if fnum >= 300 && fnum < 400 {
			l++
		} else if fnum >= 400 && fnum < 500 {
			m++
		} else if fnum >= 500 {
			n++
		}
	}
	fmt.Println(fmt.Sprintf("friendcount <100: %3d; 100-200: %3d; 200-300: %3d; 300-400: %3d; 400-500: %3d; 500-600: %3d; | 0: %3d", i, j, k, l, m, n, o))
}

func verifyFriendsDistribution(friendlist map[int][]int) bool {

	o, i, j, k, l, m, n := 0, 0, 0, 0, 0, 0, 0

	for _, friends := range friendlist {
		fnum := len(friends)
		if fnum < 1 {
			o++
		} else if fnum >= 1 && fnum < 100 {
			i++
		} else if fnum >= 100 && fnum < 200 {
			j++
		} else if fnum >= 200 && fnum < 300 {
			k++
		} else if fnum >= 300 && fnum < 400 {
			l++
		} else if fnum >= 400 && fnum < 500 {
			m++
		} else if fnum >= 500 {
			n++
		}
	}
	fmt.Println(fmt.Sprintf("friendlist  <100: %3d; 100-200: %3d; 200-300: %3d; 300-400: %3d; 400-500: %3d; 500-600: %3d; | 0: %3d", i, j, k, l, m, n, o))

	return i > 0 && j > 0 && k > 0 && l > 0 && m > 0 && n > 0
}

func genFriendsCountForUsers(totalUsers int) map[int]int {
	friendcount := make(map[int]int)
	for fRange, unum := range getFriendsBuckets(totalUsers) {
		rTokens := strings.Split(fRange, "-")
		lower_bound, _ := strconv.Atoi(rTokens[0])
		upper_bound, _ := strconv.Atoi(rTokens[1])
		ucount := 0
		for uid := 0; uid < totalUsers; uid++ {
			if _, ok := friendcount[uid]; !ok {
				friendcount[uid] = helper.RandomNumber(lower_bound, upper_bound-1)
				ucount++
			}
			if ucount == unum {
				break
			}
		}
	}
	return friendcount
}

func idExistsInList(id int, list []int) bool {
	for _, pID := range list {
		if pID == id {
			return true
		}
	}
	return false
}

func assignFriendsToUsers(totalUsers int) map[int][]int {
	friendcount := genFriendsCountForUsers(totalUsers)
	friendlist := make(map[int][]int)
	for uid := 0; uid < totalUsers; uid++ {
		for _, fid := range rand.Perm(totalUsers) {
			if uid == fid {
				continue
			}
			// check if friend already has enough friends of their own
			if friendIDs, ok := friendlist[fid]; ok && len(friendIDs) >= friendcount[fid] {
				continue
			}
			if friendIDs, ok := friendlist[uid]; ok && len(friendIDs) >= friendcount[uid] {
				break
			} else {
				if !idExistsInList(fid, friendlist[uid]) && !idExistsInList(uid, friendlist[fid]) {
					friendlist[uid] = append(friendlist[uid], fid)
					friendlist[fid] = append(friendlist[fid], uid)
				}
			}
		}
	}

	// checkFriendsDistribution(friendcount)
	// verifyFriendsDistribution(friendlist)

	return friendlist
}

func getFriendsBuckets(totalUsers int) map[string]int {
	buckets := make(map[string]int)
	for _, fRange := range getFriendRanges() {
		key := fmt.Sprintf("%d-%d", fRange.lower_bound, fRange.upper_bound)
		buckets[key] = int(fRange.pshare * float32(totalUsers))
	}
	return buckets
}

func getFriendRanges() []Range {
	// Buckets: 32% <100, 15% 100-200, 13% 200-300, 8% 300-400, 7% 400-500, 25% 500+
	var ranges []Range
	ranges = append(ranges, Range{1, 100, 0.32})
	ranges = append(ranges, Range{100, 200, 0.15})
	ranges = append(ranges, Range{200, 300, 0.13})
	ranges = append(ranges, Range{300, 400, 0.08})
	ranges = append(ranges, Range{400, 500, 0.07})
	ranges = append(ranges, Range{500, 600, 0.25})
	return ranges
}

func getBoundsForBuckets(bucketsGenerated, bucketsRequired map[string]int) (int, int) {
	lower_bound, upper_bound := -1, -1
	if bucketsGenerated["500+"] < bucketsRequired["500+"] {
		lower_bound, upper_bound = 500, 600
	} else if bucketsGenerated["400-500"] < bucketsRequired["400-500"] {
		lower_bound, upper_bound = 400, 499
	} else if bucketsGenerated["300-400"] < bucketsRequired["300-400"] {
		lower_bound, upper_bound = 300, 399
	} else if bucketsGenerated["200-300"] < bucketsRequired["200-300"] {
		lower_bound, upper_bound = 200, 299
	} else if bucketsGenerated["100-200"] < bucketsRequired["100-200"] {
		lower_bound, upper_bound = 100, 199
	} else if bucketsGenerated["<100"] < bucketsRequired["<100"] {
		lower_bound, upper_bound = 1, 99
	}
	return lower_bound, upper_bound
}

func createFriends(dbConn *sql.DB, users []*datagen.User, friendlist map[int][]int, thread_num, start, end int, wg *sync.WaitGroup) {
	defer wg.Done()

	uidx := 0
	// for uid, friends := range friendlist {
	for i := start; i < end; i++ {
		user1 := users[i]
		for fid, friend := range friendlist[i] {
			user2 := users[friend]
			aspect_idx := helper.RandomNumber(0, len(user1.Aspects)-1)
			datagen.FollowUser(dbConn, user1.Person_ID, user2.Person_ID, user1.Aspects[aspect_idx])
			datagen.FollowUser(dbConn, user2.Person_ID, user1.Person_ID, user2.Aspects[aspect_idx])
			log.Println(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Friends %3d/%3d", thread_num, uidx, end-start, fid, len(friendlist[i])))
		}
		uidx++
	}
}

func generateFriendships(dbConn *sql.DB, users []*datagen.User, allUsers []*datagen.User, thread_num int, wg *sync.WaitGroup) {
	defer wg.Done()

	totalUsers, countAllUsers := len(users), len(allUsers)
	bucketsRequired := getFriendsBuckets(countAllUsers)
	for uidx, user := range users {
		bucketsGenerated := datagen.GetFriendsDistribution(dbConn)
		if lower_bound, upper_bound := getBoundsForBuckets(bucketsGenerated, bucketsRequired); lower_bound > 0 && upper_bound > 0 {
			if existing_friends := datagen.GetTotalNumberOfFriendsOfUser(dbConn, user.Person_ID); existing_friends < lower_bound {
				number_of_friends := helper.RandomNumber(lower_bound, upper_bound)
				indices := rand.Perm(countAllUsers)
				loop_count := number_of_friends
				for count := 0; count < loop_count; count++ {
					if existing_friends = datagen.GetTotalNumberOfFriendsOfUser(dbConn, user.Person_ID); existing_friends < number_of_friends {
						log.Println(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Friends %3d/%3d | Bounds %3d-%3d", thread_num, uidx, totalUsers, existing_friends, number_of_friends, lower_bound, upper_bound))
						user2 := allUsers[indices[count]]
						if user.Person_ID != user2.Person_ID {
							aspect_idx := helper.RandomNumber(0, len(user.Aspects)-1)
							datagen.FollowUser(dbConn, user.Person_ID, user2.Person_ID, user.Aspects[aspect_idx])
							datagen.FollowUser(dbConn, user2.Person_ID, user.Person_ID, user.Aspects[aspect_idx])
						} else {
							loop_count++
						}
					} else {
						break
					}
				}
			}
		} else {
			log.Println("Negative bounds returned. Buckets completed?")
			fmt.Println(bucketsRequired)
			log.Fatal(bucketsGenerated)
			break
		}
	}
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

func AddFriendsForUser(person_id, lower_bound, upper_bound, limit string) {

	lim, _ := strconv.Atoi(limit)
	dbConn := db.GetDBConn(config.APP_NAME)
	user := datagen.GetUserByPersonID(dbConn, person_id)
	usersInRange := datagen.GetUsersWithFriendCountInRange(dbConn, lower_bound, upper_bound)
	for i, friend := range usersInRange {
		aspect_idx := helper.RandomNumber(0, len(user.Aspects)-1)
		datagen.FollowUser(dbConn, user.Person_ID, friend.Person_ID, user.Aspects[aspect_idx])
		datagen.FollowUser(dbConn, friend.Person_ID, user.Person_ID, friend.Aspects[aspect_idx])
		log.Println(fmt.Sprintf("Friends %3d/%3d", i, len(usersInRange)))
		if i == lim {
			break
		}
	}

}

func runGenerateFriendships() {

	var wg sync.WaitGroup

	dbConn := db.GetDBConn(config.APP_NAME)
	users := datagen.GetAllUsersWithAspects(dbConn)
	numUsers := datagen.GetTotalNumberOfUsers(dbConn)
	friendlist := assignFriendsToUsers(numUsers)
	numThreads := 100
	verifyFriendsDistribution(friendlist)
	log.Fatal(datagen.GetFriendsDistribution(dbConn))
	log.Fatal()
	time.Sleep(time.Duration(10) * time.Second)
	inc := int(numUsers / numThreads)
	for t, i, j := 1, 0, inc; t <= numThreads; t, i, j = t+1, j+1, j+inc {
		wg.Add(1)
		go createFriends(dbConn, users, friendlist, t, i, j, &wg)

		// i, j = j+1, j+thread_num
	}

	wg.Wait()
	time.Sleep(time.Duration(2) * time.Second)
	verifyFriendsDistribution(friendlist)
	fmt.Println(datagen.GetFriendsDistribution(dbConn))
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

	helper.Init()

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
		lower_bound := os.Args[3]
		upper_bound := os.Args[4]
		limit := os.Args[5]
		fmt.Println("Add New Friends For:", person_id)
		AddFriendsForUser(person_id, lower_bound, upper_bound, limit)
	}
}

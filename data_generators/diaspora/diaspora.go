package main

import (
	"database/sql"
	"diaspora/config"
	"diaspora/datagen"
	"diaspora/db"
	"diaspora/dist"
	"diaspora/helper"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
	"os/exec"
	"strings"
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

func createNewPostsForUsers(dbConn *sql.DB, users []*datagen.User, userPostCounts map[int]int, start, end, thread_num int, wg *sync.WaitGroup) {
	defer wg.Done()

	for uid := start - 1; uid < end; uid++ {
		for n := 0; n < userPostCounts[uid]; n++ {
			log.Print(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Posts %3d/%3d", thread_num, uid, end, n, userPostCounts[uid]))
			// datagen.NewPost(dbConn, users[uid].User_ID, users[uid].Person_ID, users[uid].Aspects)
		}
	}
}

func createNewCommentsForUsers(dbConn *sql.DB, users []*datagen.User, userCommentCounts map[int]int, start, end, thread_num int, wg *sync.WaitGroup) {
	defer wg.Done()
	// log.Fatal(start, end)
	for uidx := start; uidx < end; uidx++ {

		totalComments := userCommentCounts[uidx]

		if totalComments > 10 {
			selfComments := int(0.1 * float32(totalComments))
			uposts := datagen.GetPostsForUserLimit(dbConn, users[uidx].Person_ID, selfComments)
			for _, upost := range uposts {
				if _, err := datagen.NewComment(dbConn, upost.ID, users[uidx].Person_ID, upost.Author); err == nil {
					totalComments--
					log.Print(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Comments %3d/%3d", thread_num, uidx, end, totalComments, userCommentCounts[uidx]))
				} else {
					log.Print(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Comments %3d/%3d | ERROR: %s", thread_num, uidx, end, totalComments, userCommentCounts[uidx], err))
				}
			}
		}

		friends := datagen.GetFriendsOfUser(dbConn, users[uidx].Person_ID)
		commentsPerFriend := int(math.Ceil(float64(totalComments) / float64(len(friends))))

		for {
			for fidx, friend := range friends {
				fposts := datagen.GetPostsForUserLimit(dbConn, friend.Person_ID, commentsPerFriend)
				for pidx := 0; pidx < len(fposts) && totalComments > 0; pidx++ {
					fpost := fposts[pidx]
					if _, err := datagen.NewComment(dbConn, fpost.ID, users[uidx].Person_ID, fpost.Author); err == nil {
						totalComments--
						log.Print(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Friends %3d/%3d | Comments %3d/%3d", thread_num, uidx, end, fidx, len(friends), totalComments, userCommentCounts[uidx]))
					} else {
						log.Print(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Friends %3d/%3d | Comments %3d/%3d", thread_num, uidx, end, fidx, len(friends), totalComments, userCommentCounts[uidx]))
					}
				}
			}
			if totalComments <= 0 {
				break
			}
		}
	}
}

func createNewLikesForUsers(dbConn *sql.DB, users []*datagen.User, userLikeCounts map[int]int, start, end, thread_num int, wg *sync.WaitGroup) {
	defer wg.Done()
	// log.Fatal(start, end)
	for uidx := start; uidx < end; uidx++ {
		personID := users[uidx].Person_ID
		totalLikes := userLikeCounts[uidx]

		if totalLikes > 10 {
			selfLikes := int(0.1 * float32(totalLikes))
			uposts := datagen.GetPostsForUserLimit(dbConn, personID, selfLikes)
			for _, upost := range uposts {
				if _, err := datagen.NewLike(dbConn, upost.ID, personID, upost.Author); err == nil {
					totalLikes--
					log.Print(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Likes %3d/%3d", thread_num, uidx, end, totalLikes, userLikeCounts[uidx]))
				} else {
					log.Fatal(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Likes %3d/%3d | u%d-f%d | ERROR: %s", thread_num, uidx, end, totalLikes, userLikeCounts[uidx], personID, upost.Author, err))
				}
			}
		}

		friends := datagen.GetFriendsOfUser(dbConn, personID)
		likesPerFriend := int(math.Ceil(float64(totalLikes) / float64(len(friends))))

		for fidx, friend := range friends {
			fposts := datagen.GetPostsForUserLimit(dbConn, friend.Person_ID, likesPerFriend)
			for pidx := 0; pidx < len(fposts) && totalLikes > 0; pidx++ {
				fpost := fposts[pidx]
				if _, err := datagen.NewLike(dbConn, fpost.ID, personID, fpost.Author); err == nil {
					totalLikes--
					log.Print(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Friends %3d/%3d | Likes %3d/%3d | u%d-f%d", thread_num, uidx, end, fidx, len(friends), totalLikes, userLikeCounts[uidx], personID, friend.Person_ID))
				} else {
					log.Fatal(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Friends %3d/%3d | Likes %3d/%3d | u%d-f%d | Error: %s", thread_num, uidx, end, fidx, len(friends), totalLikes, userLikeCounts[uidx], personID, friend.Person_ID, err))
				}
			}
		}
	}
}

func createNewResharesForUsers(dbConn *sql.DB, users []*datagen.User, userResharesCounts map[int]int, start, end, thread_num int, wg *sync.WaitGroup) {
	defer wg.Done()
	// log.Fatal(start, end)
	for uidx := start; uidx < end; uidx++ {
		personID := users[uidx].Person_ID
		totalReshares := userResharesCounts[uidx]

		if totalReshares > 10 {
			selfLikes := int(0.1 * float32(totalReshares))
			uposts := datagen.GetPostsForUserLimit(dbConn, personID, selfLikes)
			for _, upost := range uposts {
				if _, err := datagen.NewReshare(dbConn, upost, personID); err == nil {
					totalReshares--
					log.Print(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Reshares %3d/%3d", thread_num, uidx, end, totalReshares, userResharesCounts[uidx]))
				} else {
					log.Fatal(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Reshares %3d/%3d | u%d-f%d | ERROR: %s", thread_num, uidx, end, totalReshares, userResharesCounts[uidx], personID, upost.Author, err))
				}
			}
		}

		friends := datagen.GetFriendsOfUser(dbConn, personID)
		likesPerFriend := int(math.Ceil(float64(totalReshares) / float64(len(friends))))

		for fidx, friend := range friends {
			fposts := datagen.GetPostsForUserLimit(dbConn, friend.Person_ID, likesPerFriend)
			for pidx := 0; pidx < len(fposts) && totalReshares > 0; pidx++ {
				fpost := fposts[pidx]
				if _, err := datagen.NewReshare(dbConn, fpost, personID); err == nil {
					totalReshares--
					log.Print(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Friends %3d/%3d | Reshares %3d/%3d | u%d-f%d", thread_num, uidx, end, fidx, len(friends), totalReshares, userResharesCounts[uidx], personID, friend.Person_ID))
				} else {
					log.Fatal(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Friends %3d/%3d | Reshares %3d/%3d | u%d-f%d | Error: %s", thread_num, uidx, end, fidx, len(friends), totalReshares, userResharesCounts[uidx], personID, friend.Person_ID, err))
				}
			}
		}
	}
}

func createNewMessagesForUsers(dbConn *sql.DB, users []*datagen.User, userMessageCounts map[int]int, start, end, thread_num int, wg *sync.WaitGroup) {
	defer wg.Done()
	const dunbarsNumber int = 150

	for uidx := start; uidx < end; uidx++ {

		totalMessages := userMessageCounts[uidx]

		if totalMessages <= 0 {
			continue
		}

		friends := datagen.GetFriendsOfUser(dbConn, users[uidx].Person_ID)
		if len(friends) > dunbarsNumber {
			friends = friends[:dunbarsNumber]
		}

		messagesPerFriend := int(math.Ceil(float64(totalMessages) / float64(len(friends))))

		for fidx, friend := range friends {
			conv_id, _ := datagen.NewConversation(dbConn, users[uidx].Person_ID, friend.Person_ID)
			for mID := 0; mID < messagesPerFriend; mID++ {
				// user sends message
				if _, err := datagen.NewMessage(dbConn, users[uidx].Person_ID, conv_id); err == nil {
					totalMessages--
					userMessageCounts[uidx]--
					log.Print(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Friends %3d/%3d | Messages %3d/%3d", thread_num, uidx, end, fidx, len(friends), totalMessages, userMessageCounts[uidx]))
				} else {
					log.Print(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Friends %3d/%3d | Messages %3d/%3d | ERROR: %s", thread_num, uidx, end, fidx, len(friends), totalMessages, userMessageCounts[uidx], err))
				}
				// friend sends message
				if fIDXInList := datagen.FindIndexInUserListByPersonID(users, friend.Person_ID); userMessageCounts[fIDXInList] > 0 {
					if _, err := datagen.NewMessage(dbConn, friend.Person_ID, conv_id); err == nil {
						userMessageCounts[fIDXInList]--
						log.Print(fmt.Sprintf("Thread # %3d | Users: %3d/%3d | Friends %3d/%3d | Messages %3d/%3d | FRIEND REPLY", thread_num, uidx, end, fidx, len(friends), totalMessages, userMessageCounts[uidx]))
					} else {

					}
				}
			}
			if totalMessages <= 0 {
				break
			}
		}
	}
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
	bucketsRequired := dist.GetFriendsBuckets(countAllUsers)
	for uidx, user := range users {
		bucketsGenerated := datagen.GetFriendsDistribution(dbConn)
		if lower_bound, upper_bound := dist.GetBoundsForBuckets(bucketsGenerated, bucketsRequired); lower_bound > 0 && upper_bound > 0 {
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
		friends_of_user := datagen.GetFriendsOfUser(dbConn, user.Person_ID)
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
					datagen.NewReshare(dbConn, post, friend.Person_ID)
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

	numThreads := 100
	dbConn := db.GetDBConn(config.APP_NAME)
	users := datagen.GetAllUsersWithAspects(dbConn)
	numUsers := datagen.GetTotalNumberOfUsers(dbConn)
	friendlist := dist.AssignFriendsToUsers(numUsers)
	dist.VerifyFriendsDistribution(friendlist)
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
	dist.VerifyFriendsDistribution(friendlist)
	fmt.Println(datagen.GetFriendsDistribution(dbConn))
}

func runCreateNewPosts(ptype string, totalPosts int) {

	var wg sync.WaitGroup

	numThreads := 1
	userPostCounts := dist.AssignPostsToUsers(totalPosts)
	// log.Fatal(userPostCounts)
	dbConn := db.GetDBConn(config.APP_NAME)
	users := datagen.GetUsersOrderedByFriendCount(dbConn)
	users = users[:len(userPostCounts)]
	numUsers := len(users)
	inc := int(numUsers / numThreads)
	for t, i, j := 1, 0, inc; t <= numThreads; t, i, j = t+1, j+1, j+inc {
		wg.Add(1)
		switch ptype {
		case "posts":
			{
				go createNewPostsForUsers(dbConn, users, userPostCounts, i, j, t, &wg)
			}
		case "likes":
			{
				go createNewLikesForUsers(dbConn, users, userPostCounts, i, j, t, &wg)
			}
		case "comments":
			{
				go createNewCommentsForUsers(dbConn, users, userPostCounts, i, j, t, &wg)
			}
		case "reshares":
			{
				go createNewResharesForUsers(dbConn, users, userPostCounts, i, j, t, &wg)
			}
		case "messages":
			{
				go createNewMessagesForUsers(dbConn, users, userPostCounts, i, j, t, &wg)
			}
		}
	}
	wg.Wait()
} 

func ParetoNewPosts(alpha, total int){
	dbConn := db.GetDBConn(config.APP_NAME)
	users := datagen.GetAllUsersWithAspects(dbConn)
	{
		totalUsers := len(users)
		popDist := getParetoResultFromPython(2.0, totalUsers)
		commDist := getParetoResultFromPython(2.1, totalUsers)
		likeDist := getParetoResultFromPython(2.2, totalUsers)
		for i := 0; i < totalUsers; i++ {
			users[i].PopularityScore = popDist[i]
			users[i].CommentScore = commDist[i]
			users[i].LikeScore = likeDist[i]
		}
	}
	
	// fmt.Println(dist)
}

func getParetoResultFromPython(alpha float32, total int) []float64 { // alpha = 2, 3?
	var dist []float64
	cmd := exec.Command("python", "../pareto.py", fmt.Sprint(alpha), fmt.Sprint(total))
    if out, err := cmd.CombinedOutput(); err != nil {
		log.Fatal(err); 
	}else{
		nums := strings.Split(string(out), ",")
		for _, num := range nums {
			num = strings.Replace(num, "\n", "", -1)
			if value, err := strconv.ParseFloat(num, 32); err == nil {
				dist = append(dist, value)
			}else{
				log.Fatal("Crashed while converting pareto val to float32:", err)
			}
		}
	}
	return dist
}

func test() {
	dbConn := db.GetDBConn(config.APP_NAME)
	users := datagen.GetAllUsersWithAspects(dbConn)
	for uid := 0; uid < len(users); uid++ {
		datagen.NewPhotoPost(dbConn, users[uid].User_ID, users[uid].Person_ID, users[uid].Aspects, 5)
	}	
}

func main() {

	helper.Init()

	arg := os.Args[1]

	switch arg {
	case "test":
		test()
	case "posts":
		fmt.Println("Creating New Posts!")
		ParetoNewPosts(2, 8030)
	case "comments":
		fmt.Println("Interacting With Posts!")
		runCreateNewPosts("comments", 13970)
	case "likes":
		fmt.Println("Interacting With Posts!")
		runCreateNewPosts("likes", 85680)
	case "reshares":
		fmt.Println("Interacting With Posts!")
		runCreateNewPosts("reshares", 1550)
	case "messages":
		fmt.Println("Creating New Messages!")
		runCreateNewPosts("messages", 4015)
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
		lower_bound := os.Args[3]
		upper_bound := os.Args[4]
		limit := os.Args[5]
		fmt.Println("Add New Friends For:", person_id)
		AddFriendsForUser(person_id, lower_bound, upper_bound, limit)
	}
}

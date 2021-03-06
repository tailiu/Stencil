package datagen

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
	"twitter/db"
	"twitter/helper"
)

type Post struct {
	ID     int
	Author int
}

func CreateNewUser(dbConn *sql.DB) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewUser: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "users")
	nickname := helper.RandomString(helper.RandomNumber(4, 10))
	fullname := fmt.Sprintf("%s %s", helper.RandomString(5), helper.RandomString(5))
	password := "$2a$10$408zooOxx9.C.sNm9Csg0.uY83YZ.1f6qX1m4tn3D8tD03jbPPs62"
	email := fmt.Sprintf("%s@%s.com", nickname, helper.RandomString(helper.RandomNumber(2, 8)))
	bio := helper.RandomText(helper.RandomNumber(20, 200))

	sql := "INSERT INTO \"users\" (id, created_at, updated_at, name, handle, bio) VALUES ($1, $2, $3, $4, $5, $6)"

	txErr := db.RunTxWQnArgs(tx, sql, id, time.Now(), time.Now(), fullname, nickname, bio)
	if txErr != nil {
		log.Fatal("Error while creating new user ('users'): ", txErr)
	}

	sql = "INSERT INTO credentials (id, email, password, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	txErr = db.RunTxWQnArgs(tx, sql, id, email, password, id, time.Now(), time.Now())
	if txErr != nil {
		log.Fatal("Error while creating new user ('credentials'): ", txErr)
	}

	tx.Commit()

	return id
}

func CreateNewPost(dbConn *sql.DB, userID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewPost: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "tweets")
	content := helper.RandomText(helper.RandomNumber(10, 100))

	sql := "INSERT INTO tweets (id, created_at, updated_at, content, user_id) VALUES ($1, $2, $3, $4, $5)"

	txErr := db.RunTxWQnArgs(tx, sql, id, time.Now(), time.Now(), content, userID)
	if txErr != nil {
		log.Fatal("Error while creating new tweets ", txErr)
	}

	tx.Commit()

	return id
}

func CreateNewComment(dbConn *sql.DB, userID, postID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewComment: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "tweets")
	content := helper.RandomText(helper.RandomNumber(10, 100))

	sql := "INSERT INTO tweets (id, created_at, updated_at, content, user_id, reply_to_id) VALUES ($1, $2, $3, $4, $5, $6)"

	txErr := db.RunTxWQnArgs(tx, sql, id, time.Now(), time.Now(), content, userID, postID)
	if txErr != nil {
		log.Fatal("Error while creating new comments ", txErr)
	}

	tweetOwnerID := GetTweetOwner(dbConn, postID)

	if err := CreateNotification(dbConn, tx, userID, tweetOwnerID, postID, "comment"); err != nil {
		log.Fatal("Error while creating new follow > notification: ", err)
	}

	tx.Commit()

	return id
}

func CreateNewReshare(dbConn *sql.DB, userID, postID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewReshare: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "retweets")

	sql := "INSERT INTO retweets (id, user_id, tweet_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"

	txErr := db.RunTxWQnArgs(tx, sql, id, userID, postID, time.Now(), time.Now())
	if txErr != nil {
		log.Fatal("Error while creating new retweets ", txErr)
	}

	tweetOwnerID := GetTweetOwner(dbConn, postID)

	if err := CreateNotification(dbConn, tx, userID, tweetOwnerID, postID, "retweet"); err != nil {
		log.Fatal("Error while creating new follow > notification: ", err)
	}

	tx.Commit()

	return id
}

func CreateTweetWithPhoto(dbConn *sql.DB, userID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateTweetWithPhoto: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "tweets")
	content := helper.RandomText(helper.RandomNumber(10, 100))

	path := "/home/zain/project/resources/"
	photoID := helper.RandomNumber(1, 5)
	photoName := fmt.Sprintf("%d.jpg", photoID)
	photoPath := path + photoName
	// jsonMedia, err := json.Marshal(photoPath)
	// if err != nil {
	// 	fmt.Println("path: ", photoPath)
	// 	log.Fatal("Error while converting path to json: ", err)
	// }
	mediaType := "photo"

	sql := "INSERT INTO tweets (id, created_at, updated_at, content, user_id, tweet_media, media_type) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	txErr := db.RunTxWQnArgs(tx, sql, id, time.Now(), time.Now(), content, userID, photoPath, mediaType)
	if txErr != nil {
		log.Fatal("Error while creating new tweets with photo: ", txErr)
	}

	tx.Commit()

	return id
}

func CreateNewLike(dbConn *sql.DB, postID, userID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewLike: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "likes")

	sql := "INSERT INTO likes (id, tweet_id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"

	txErr := db.RunTxWQnArgs(tx, sql, id, postID, userID, time.Now(), time.Now())
	if txErr != nil {
		log.Fatal("Error while creating new like: ", txErr)
	}

	tweetOwnerID := GetTweetOwner(dbConn, postID)

	if err := CreateNotification(dbConn, tx, userID, tweetOwnerID, postID, "like"); err != nil {
		log.Fatal("Error while creating new follow > notification: ", err)
	}

	tx.Commit()

	return id
}

func conversationParticipantsExists(dbConn *sql.DB, userID, conversationID string) bool {

	q := fmt.Sprintf(
		`SELECT id FROM conversation_participants c1 
		WHERE c1.user_id = %s and c1.conversation_id = %s`,
		userID, conversationID,
	)

	v, err := db.DataCall1(dbConn, q)
	if err != nil {
		log.Fatal(err)
	}

	if len(v) > 0 {
		return true
	} else {
		return false
	}

}

func CreateNewConversation(dbConn *sql.DB, userID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewConversation: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "conversations")
	role := "creator"

	sql := "INSERT INTO conversations (id, created_at, updated_at) VALUES ($1, $2, $3)"
	txErr := db.RunTxWQnArgs(tx, sql, id, time.Now(), time.Now())
	if txErr != nil {
		log.Fatal("Error while creating new conversation: ", txErr)
	}

	cpid := db.GetNewRowIDForTable(dbConn, "conversation_participants")
	sql = "INSERT INTO conversation_participants (id, conversation_id, user_id, created_at, updated_at, role, saw_new_messages, saw_messages_until) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	txErr = db.RunTxWQnArgs(tx, sql, cpid, id, userID, time.Now(), time.Now(), role, true, time.Now())
	if txErr != nil {
		log.Fatal("Error while creating new conversation_participants: ", txErr)
	}

	tx.Commit()

	return id
}

func CreateNewMessage(dbConn *sql.DB, userID, conversationID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewMessage: Tx error:", err)
	}

	if !conversationParticipantsExists(dbConn, userID, conversationID) {
		cpid := db.GetNewRowIDForTable(dbConn, "conversation_participants")
		sql0 := "INSERT INTO conversation_participants (id, conversation_id, user_id, created_at, updated_at, role, saw_new_messages, saw_messages_until) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
		txErr0 := db.RunTxWQnArgs(tx, sql0, cpid, conversationID, userID, time.Now(), time.Now(), "participant", true, time.Now())
		if txErr0 != nil {
			log.Fatal("Error while creating new conversation_participants: ", txErr0)
		}
	}

	id := db.GetNewRowIDForTable(dbConn, "messages")
	content := helper.RandomText(helper.RandomNumber(10, 100))

	sql := "INSERT INTO messages (id, created_at, updated_at, content, conversation_id, user_id) VALUES ($1, $2, $3, $4, $5, $6)"
	txErr := db.RunTxWQnArgs(tx, sql, id, time.Now(), time.Now(), content, conversationID, userID)
	if txErr != nil {
		log.Fatal("Error while creating new message: ", txErr)
	}

	tx.Commit()

	return id
}

func CreateNewMessageWithPhoto(dbConn *sql.DB, userID, conversationID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewMessage: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "messages")
	content := helper.RandomText(helper.RandomNumber(10, 100))

	path := "/home/zain/project/resources/"
	photoID := helper.RandomNumber(1, 5)
	photoName := fmt.Sprintf("%d.jpg", photoID)
	photoPath := path + photoName
	// jsonMedia, err := json.Marshal(photoPath)
	// if err != nil {
	// 	fmt.Println("path: ", photoPath)
	// 	log.Fatal("Error while converting path to json: ", err)
	// }
	mediaType := "photo"

	sql := "INSERT INTO messages (id, created_at, updated_at, content, conversation_id, user_id, message_media, media_type) VALUES ($1, $2, $3, $4, $5, $6)"
	txErr := db.RunTxWQnArgs(tx, sql, id, time.Now(), time.Now(), content, conversationID, userID, photoPath, mediaType)
	if txErr != nil {
		log.Fatal("Error while creating new message: ", txErr)
	}

	tx.Commit()

	return id
}

func GetPostsForUser(dbConn *sql.DB, userID int) []*Post {

	q := fmt.Sprintf(`SELECT id FROM tweets WHERE user_id = %d`, userID)

	v := db.DataCall(dbConn, q)

	var posts []*Post
	for _, v1 := range v {
		pid, err := strconv.Atoi(v1["id"])
		if err != nil {
			log.Fatal(err)
		}
		post := new(Post)
		post.Author = userID
		post.ID = pid
		posts = append(posts, post)
	}

	return posts

}

func GetFollowingUsers(dbConn *sql.DB, fromUserID int) []int {

	q := fmt.Sprintf(
		`SELECT to_user_id FROM user_actions 
		WHERE from_user_id = %d and action_type = 'follows'`, fromUserID,
	)

	v := db.DataCall(dbConn, q)

	var result []int
	for _, v1 := range v {
		res1, err := strconv.Atoi(v1["to_user_id"])
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, res1)
	}

	return result

}

func GetFollowedNum(dbConn *sql.DB, toUserID int) int {

	q := fmt.Sprintf(
		`SELECT COUNT(*) FROM user_actions 
		WHERE to_user_id = %d and action_type = 'follows'`, toUserID,
	)

	v, err := db.DataCall1(dbConn, q)
	if err != nil {
		log.Fatal(err)
	}

	count, err1 := strconv.Atoi(fmt.Sprint(v["count"]))
	if err1 != nil {
		log.Fatal(err1)
	}
	return count

}

func CheckFollowed(dbConn *sql.DB, toUserID, fromUserID int) bool {

	q := fmt.Sprintf(
		`SELECT id FROM user_actions 
		WHERE to_user_id = %d and from_user_id = %d and action_type = 'follows'`,
		toUserID, fromUserID,
	)

	v, err := db.DataCall1(dbConn, q)
	if err != nil {
		log.Fatal(err)
	}

	if len(v) > 0 {
		return true
	} else {
		return false
	}
}

func CheckConversationBetweenTwoUsers(dbConn *sql.DB, userID1, userID2 int) (bool, int) {

	q := fmt.Sprintf(
		`SELECT c1.conversation_id FROM conversation_participants c1 JOIN conversation_participants c2 ON 
		c1.conversation_id = c2.conversation_id WHERE c1.user_id = %d and c2.user_id = %d`,
		userID1, userID2,
	)

	v, err := db.DataCall1(dbConn, q)
	if err != nil {
		log.Fatal(err)
	}

	if len(v) > 0 {
		cIDStr := fmt.Sprint(v["conversation_id"])
		cID, err1 := strconv.Atoi(cIDStr)
		if err1 != nil {
			log.Fatal(err1)
		}
		return true, cID
	} else {
		return false, -1
	}

}

func GetRealFriendsOfUser(dbConn *sql.DB, fromUserID int) []int {

	q := fmt.Sprintf(
		`SELECT DISTINCT a1.to_user_id FROM user_actions a1 JOIN user_actions a2 ON 
		a1.from_user_id = a2.to_user_id AND a1.to_user_id = a2.from_user_id
		WHERE a1.from_user_id = %d`, fromUserID,
	)

	v := db.DataCall(dbConn, q)

	var result []int
	for _, v1 := range v {
		res1, err := strconv.Atoi(v1["to_user_id"])
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, res1)
	}

	return result

}

func GetFollowedUsers(dbConn *sql.DB, toUserID int) []int {

	q := fmt.Sprintf(
		`SELECT from_user_id FROM user_actions 
		WHERE to_user_id = %d and action_type = 'follows'`, toUserID,
	)

	v := db.DataCall(dbConn, q)

	var result []int
	for _, v1 := range v {
		res1, err := strconv.Atoi(v1["from_user_id"])
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, res1)
	}

	return result

}

func CreateNewFollow(dbConn *sql.DB, fromUserID, toUserID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewFollow: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "user_actions")
	actionType := "follows"

	sql := "INSERT INTO user_actions (id, from_user_id, to_user_id, action_type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"

	txErr := db.RunTxWQnArgs(tx, sql, id, fromUserID, toUserID, actionType, time.Now(), time.Now())
	if txErr != nil {
		log.Fatal("Error while creating new follow: ", txErr)
	}

	if err := CreateNotification(dbConn, tx, fromUserID, toUserID, "0", "follow"); err != nil {
		log.Fatal("Error while creating new follow > notification: ", err)
	}

	tx.Commit()

	return id
}

func CreateNotification(dbConn *sql.DB, tx *sql.Tx, fromUser, toUser, postID, notifType string) error {

	id := db.GetNewRowIDForTable(dbConn, "notifications")

	sql := "INSERT INTO notifications (id, notification_type, user_id, from_user, tweet, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	txErr := db.RunTxWQnArgs(tx, sql, id, notifType, toUser, fromUser, postID, time.Now(), time.Now())
	if txErr != nil {
		log.Fatal("Error while creating new notification: ", txErr)
		return txErr
	}
	return nil
}

func GetTweetOwner(dbConn *sql.DB, tweetID string) string {

	sql := "SELECT user_id FROM tweets WHERE id = $1"

	if data, err := db.DataCall1(dbConn, sql, tweetID); err != nil {
		log.Fatal("Can't find tweet owner with tweet id: ", tweetID, " | ", err)
		return ""
	} else {
		return fmt.Sprint(data["user_id"])
	}
}

package datagen

import (
	"database/sql"
	"fmt"
	"gnusocial/db"
	"gnusocial/helper"
	"log"
	"time"
	"strconv"
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

	id := db.GetNewRowIDForTable(dbConn, "profile")
	nickname := helper.RandomString(helper.RandomNumber(4, 10))
	fullname := fmt.Sprintf("%s %s", helper.RandomString(5), helper.RandomString(5))
	password := "$2a$10$408zooOxx9.C.sNm9Csg0.uY83YZ.1f6qX1m4tn3D8tD03jbPPs62"
	email := fmt.Sprintf("%s@%s.com", nickname, helper.RandomString(helper.RandomNumber(2, 8)))
	language := "en"
	profileurl := fmt.Sprintf("www.gnusocial.com/profile/%s", nickname)
	homepage := fmt.Sprintf("www.gnusocial.com/home/%s", nickname)
	bio := helper.RandomText(helper.RandomNumber(20, 200))

	sql := "INSERT INTO \"user\" (id, nickname, password, email, language, created, modified) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id "

	txErr := db.RunTxWQnArgs(tx, sql, id, nickname, password, email, language, time.Now(), time.Now())
	if txErr != nil {
		log.Fatal("Error while creating new user ('user'): ", txErr)
	}

	sql = "INSERT INTO profile (id, nickname, fullname, profileurl, homepage, bio, created, modified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id "
	txErr = db.RunTxWQnArgs(tx, sql, id, nickname, fullname, profileurl, homepage, bio, time.Now(), time.Now())
	if txErr != nil {
		log.Fatal("Error while creating new user ('profile'): ", txErr)
	}

	tx.Commit()

	return id
}

func CreateNewConversation(dbConn *sql.DB) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewConversation: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "conversation")
	uri := fmt.Sprintf("www.gnusocial.com/conversation/%s", id)

	sql := "INSERT INTO conversation (id, uri, created, modified) VALUES ($1, $2, $3, $4)"

	txErr := db.RunTxWQnArgs(tx, sql, id, uri, time.Now(), time.Now())
	if txErr != nil {
		log.Fatal("Error while creating new conversation: ", txErr)
	}

	tx.Commit()

	return id
}

func CreateNewPost(dbConn *sql.DB, conversationID, profileID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewNotice: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "notice")
	url := fmt.Sprintf("www.gnusocial.com/notice/%s", id)
	content := helper.RandomText(helper.RandomNumber(10, 100))

	sql := "INSERT INTO notice (id, profile_id, uri, content, url, created, modified, conversation) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	txErr := db.RunTxWQnArgs(tx, sql, id, profileID, url, content, url, time.Now(), time.Now(), conversationID)
	if txErr != nil {
		log.Fatal("Error while creating new post ", txErr)
	}

	tx.Commit()

	return id
}

func CreateNewComment(dbConn *sql.DB, conversationID, profileID, postID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewNotice: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "notice")
	url := fmt.Sprintf("www.gnusocial.com/notice/%s", id)
	content := helper.RandomText(helper.RandomNumber(10, 100))

	sql := "INSERT INTO notice (id, profile_id, uri, content, url, created, modified, reply_to, conversation) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	txErr := db.RunTxWQnArgs(tx, sql, id, profileID, url, content, url, time.Now(), time.Now(), postID, conversationID)
	if txErr != nil {
		log.Fatal("Error while creating new comment: ", txErr)
	}

	tx.Commit()

	return id
}

func CreateNewReshare(dbConn *sql.DB, conversationID, profileID, postID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewNotice: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "notice")
	url := fmt.Sprintf("www.gnusocial.com/notice/%s", id)
	content := helper.RandomText(helper.RandomNumber(10, 100))

	sql := "INSERT INTO notice (id, profile_id, uri, content, url, created, modified, conversation, repeat_of) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	txErr := db.RunTxWQnArgs(tx, sql, id, profileID, url, content, url, time.Now(), time.Now(), conversationID, postID)
	if txErr != nil {
		log.Fatal("Error while creating new reshare ", txErr)
	}

	tx.Commit()

	return id
}

func CreateMediaForPost(dbConn *sql.DB, postID, profileID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateMediaForPost: Tx error:", err)
	}

	fileID := db.GetNewRowIDForTable(dbConn, "file")
	fileToPostID := db.GetNewRowIDForTable(dbConn, "file_to_post")
	path := "/home/zain/project/resources/"
	photoID := helper.RandomNumber(1, 5)
	photoName := fmt.Sprintf("%d.jpg", photoID)
	photoPath := path + photoName
	urlhash := helper.GenerateHashOfString(photoPath + fileID)
	height := 60
	width := 60

	title := fmt.Sprintf("%s %s", helper.RandomString(5), helper.RandomString(5))

	sql := "INSERT INTO file (id, urlhash, url, title, date, filename, width, height, modified, profile_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	txErr := db.RunTxWQnArgs(tx, sql, fileID, urlhash, photoPath, title, time.Now(), photoName, width, height, time.Now(), profileID)
	if txErr != nil {
		log.Fatal("Error while creating new file: ", txErr)
	}

	sql = "INSERT INTO file_to_post (id, file_id, post_id, modified) VALUES ($1, $2, $3, $4)"
	txErr = db.RunTxWQnArgs(tx, sql, fileToPostID, fileID, postID, time.Now())
	if txErr != nil {
		log.Fatal("Error while creating new file_to_post: ", txErr)
	}

	tx.Commit()

	return fileID
}

func CreateNewLike(dbConn *sql.DB, postID, profileID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewLike: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "fave")
	url := fmt.Sprintf("www.gnusocial.com/fave/%s", id)

	sql := "INSERT INTO fave (id, notice_id, user_id, uri, created, modified) VALUES ($1, $2, $3, $4, $5, $6)"

	txErr := db.RunTxWQnArgs(tx, sql, id, postID, profileID, url, time.Now(), time.Now())
	if txErr != nil {
		log.Fatal("Error while creating new like: ", txErr)
	}

	tx.Commit()

	return id
}

func CreateNewMessage(dbConn *sql.DB, fromProfileID, toProfileID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewMessage: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "message")
	url := fmt.Sprintf("www.gnusocial.com/message/%s", id)
	content := helper.RandomText(helper.RandomNumber(10, 100))

	sql := "INSERT INTO message (id, uri, from_profile, to_profile, content, url, created, modified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	txErr := db.RunTxWQnArgs(tx, sql, id, url, fromProfileID, toProfileID, content, url, time.Now(), time.Now())
	if txErr != nil {
		log.Fatal("Error while creating new message: ", txErr)
	}

	tx.Commit()

	return id
}

func CreateNewSubscription(dbConn *sql.DB, fromProfileID, toProfileID string) string {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("@CreateNewSubscription: Tx error:", err)
	}

	id := db.GetNewRowIDForTable(dbConn, "subscription")
	url := fmt.Sprintf("www.gnusocial.com/subscription/%s", id)

	sql := "INSERT INTO subscription (id, subscriber, subscribed, uri, created, modified) VALUES ($1, $2, $3, $4, $5, $6)"

	txErr := db.RunTxWQnArgs(tx, sql, id, fromProfileID, toProfileID, url, time.Now(), time.Now())
	if txErr != nil {
		log.Fatal("Error while creating new subscription: ", txErr)
	}

	tx.Commit()

	return id
}

func GetFollowedUsers(dbConn *sql.DB, toUserID int) []int {

	q := fmt.Sprintf(`SELECT subscriber FROM subscription WHERE subscribed = %d`, toUserID)

	v := db.DataCall(dbConn, q)
	
	var result []int 
	for _, v1 := range v {
		res1, err := strconv.Atoi(v1["subscriber"])
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, res1)
	} 

	return result
}

func CheckFollowed(dbConn *sql.DB, toUserID, fromUserID int) bool {

	q := fmt.Sprintf(
		`SELECT id FROM subscription WHERE subscribed = %d and subscriber = %d`, 
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

func GetFollowedNum(dbConn *sql.DB, toUserID int) int {

	q := fmt.Sprintf(`SELECT COUNT(*) FROM subscription WHERE subscribed = %d`, toUserID)

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

func GetFollowingUsers(dbConn *sql.DB, fromUserID int) []int {

	q := fmt.Sprintf(`SELECT subscribed FROM subscription WHERE subscriber = %d`, fromUserID)

	v := db.DataCall(dbConn, q)
	
	var result []int 
	for _, v1 := range v {
		res1, err := strconv.Atoi(v1["subscribed"])
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, res1)
	} 

	return result
}

func GetPostsForUser(dbConn *sql.DB, userID int) []*Post {

	q := fmt.Sprintf(`SELECT id FROM notice WHERE profile_id = %d`, userID)

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

func GetConversationIDOfPost(dbConn *sql.DB, pid int) string {

	q := fmt.Sprintf(`SELECT conversation FROM notice WHERE id = %d`, pid)

	v, err := db.DataCall1(dbConn, q)
	if err != nil {
		log.Fatal(err)
	}

	if len(v) == 0 {
		log.Fatal("Fail to find the conversation id of a notice")
	}
	
	return fmt.Sprint(v["conversation"])

}

func GetRealFriendsOfUser(dbConn *sql.DB, fromUserID int) []int {

	q := fmt.Sprintf(
		`SELECT DISTINCT a1.subscribed FROM subscription a1 JOIN subscription a2 ON 
		a1.subscriber = a2.subscribed AND a1.subscribed = a2.subscriber
		WHERE a1.subscriber = %d`, fromUserID,
	)

	v := db.DataCall(dbConn, q)
	
	var result []int 
	for _, v1 := range v {
		res1, err := strconv.Atoi(v1["subscribed"])
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, res1)
	} 

	return result

}
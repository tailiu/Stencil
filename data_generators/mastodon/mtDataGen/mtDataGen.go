/*
 * Basic Functions for Mastodon
 */

package mtDataGen

import (
	"database/sql"
	"mastodon/database"
    "fmt"
	"log"
	"time"
	"strconv"
	"encoding/json"
	"mastodon/auxiliary"
	"errors"
	_ "github.com/lib/pq"
)

type FileMeta struct {
	Width int
	Height int
	Size int
	Aspect int	
}

type Post struct {
	ID     int
	Author int
}

func updateAccountStats(dbConn *sql.DB, accountID int, updateType string, difference int) {
	var statusNum int
	var	sqls []string 
	tx := database.BeginTx(dbConn)
	sql1 := fmt.Sprintf(
		"SELECT %s FROM account_stats WHERE account_stats.account_id = %d LIMIT %d;",
		updateType, accountID, 1)
	rows, err := tx.Query(sql1)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&statusNum); err != nil {
			log.Fatal(err)
		}
	}
	sql2 := fmt.Sprintf(
		"UPDATE account_stats SET %s = %d, last_status_at = now(), updated_at = now() WHERE account_stats.account_id = %d;",
		updateType, statusNum + difference, accountID) // These last_status_at and updated_at are generated randomly

	sqls = append(sqls, sql2)
	database.Execute(tx, sqls)
	tx.Commit()
}

func updateStatusStats(dbConn *sql.DB, accountID int, statusID int, updateType string, difference int) {
	var sqls []string
	count := -1
	tx := database.BeginTx(dbConn)
	sql1 := fmt.Sprintf(
		"SELECT status_stats.%s FROM status_stats WHERE status_stats.status_id = %d LIMIT %d;",
		updateType, statusID, 1)
	rows, err := tx.Query(sql1)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}
	var sql2 string
	t := time.Now().Format(time.RFC3339)
	if count == -1 {
		sql2 = fmt.Sprintf(
			"INSERT INTO status_stats (status_id, %s, created_at, updated_at) VALUES (%d, %d, '%s', '%s');",
			updateType, statusID, 1, t, t)
	} else {
		sql2 = fmt.Sprintf(
			"UPDATE status_stats SET %s = %d, updated_at = '%s' WHERE status_id = %d;",
			updateType, count + difference, t, statusID)
	}
	sqls = append(sqls, sql2)
	database.Execute(tx, sqls)
	tx.Commit()
}

func insertIntoStreamEntries(activityID int, accountID int, hidden bool) string {
	return fmt.Sprintf(
		`INSERT INTO stream_entries (activity_id, activity_type, created_at, updated_at, account_id, hidden) 
		VALUES (%d, '%s', now(), now(), %d, %t);`,
		activityID, "Status", accountID, hidden,
	)
}

/*
 * Visibility 
 * 0: Public Statuses
 * 3: Direct Messages
 */
func NewStatus(dbConn *sql.DB, accountID int, haveMedia bool, visibility int) (int, error) {
	
	var sqls1 []string
	var hidden bool

	conversationID := auxiliary.RandomNonnegativeInt()
	statusID := auxiliary.RandomNonnegativeInt()
	content := auxiliary.RandStrSeq(50)
	uri := "http://localhost:3000/users/admin/statuses/" + strconv.Itoa(statusID)

	if visibility == 0 {
		hidden = false
	} else if visibility == 3 {
		hidden = true
	}

	tx := database.BeginTx(dbConn)
	sql1 := fmt.Sprintf("INSERT INTO conversations (id, created_at, updated_at) VALUES (%d, now(), now());", conversationID)
	sql2 := fmt.Sprintf(
		`INSERT INTO statuses (id, text, created_at, updated_at, language, conversation_id, local, account_id, application_id, uri, visibility) 
		VALUES (%d, '%s', now(), now(), '%s', %d, %t, %d, %d, '%s', %d);`,  
		statusID, content, "en", conversationID, true, accountID, 1, uri, visibility)
	sql3 := insertIntoStreamEntries(statusID, accountID, hidden)
	sqls1 = append(sqls1, sql1, sql2, sql3)
	
	if haveMedia {

		// Photo Params
		photo_path := "/home/zain/project/resources/"
		photo_id := auxiliary.RandomNumber(1, 5)
		photo_name := fmt.Sprintf("%d.jpg", photo_id)
		file_file_name := photo_path + photo_name

		file_content_type := "image/jpeg"
		file_file_size := auxiliary.RandomNonnegativeInt()
		shortCode := auxiliary.RandStrSeq(20)
		file_meta, err := json.Marshal(FileMeta{auxiliary.RandomNonnegativeInt(), auxiliary.RandomNonnegativeInt(), auxiliary.RandomNonnegativeInt(), auxiliary.RandomNonnegativeInt()})
		if err != nil {
			log.Fatal(err)
		}

		sql4 := fmt.Sprintf(
			`INSERT INTO media_attachments (status_id, file_file_name, file_content_type, file_file_size, file_updated_at, created_at, updated_at, shortcode, file_meta, account_id) 
			VALUES (%d, '%s', '%s', %d, now(), now(), now(), '%s', '%s', %d);`,
		statusID, file_file_name, file_content_type, file_file_size, shortCode, file_meta, accountID)
		sqls1 = append(sqls1, sql4)
	}

	// if len(mentionedAccounts) != 0 {
	// 	var sql5 string
	// 	for _, mentionedAccount := range mentionedAccounts {
	// 		sql5 = fmt.Sprintf(
	// 			`INSERT INTO mentions (status_id, created_at, updated_at, account_id) 
	// 			VALUES (%d, now(), now(), %d);`,
	// 			statusID, mentionedAccount)
	// 		sqls1 = append(sqls1, sql5)
	// 	}
	// }

	result := database.Execute(tx, sqls1)
	if result {
		tx.Commit()
		updateAccountStats(dbConn, accountID, "statuses_count", 1)
		return statusID, nil
	} else {
		tx.Rollback()
		return -1, errors.New("Fail to create a new status")
	}
}

func GetPostsForUser(dbConn *sql.DB, accountID int) []*Post {

	var posts []*Post

	sql := `SELECT id, account_id FROM statuses WHERE account_id = $1 and reply = false 
			ORDER by random()`

	for _, row := range database.DataCall(dbConn, sql, accountID) {
		if pid, err := strconv.Atoi(row["id"]); err == nil {
			if uid, err := strconv.Atoi(row["account_id"]); err == nil {
				post := new(Post)
				post.Author = uid
				post.ID = pid
				posts = append(posts, post)
			}
		}
	}

	return posts
}

func Favourite(dbConn *sql.DB, accountID int, statusID int) {
	tx := database.BeginTx(dbConn)

	var sqls []string

	sql1 := fmt.Sprintf(
		"SELECT 1 AS one FROM favourites WHERE favourites.status_id = %d AND favourites.account_id = %d LIMIT %d;",
		statusID, accountID, 1)
	if exists := database.CheckExists(tx, sql1); exists == 1 {
		tx.Commit()
		return 
	}

	t := time.Now().Format(time.RFC3339)
	sql3 := fmt.Sprintf(
		"INSERT INTO favourites (created_at, updated_at, account_id, status_id) VALUES ('%s', '%s', %d, %d);",
		t, t, accountID, statusID)
	sqls = append(sqls, sql3)
	result := database.Execute(tx, sqls)
	if result {
		tx.Commit()
		updateStatusStats(dbConn, accountID, statusID, "favourites_count", 1)
	} else {
		tx.Rollback()
	}
	
}

func Unfavourite(dbConn *sql.DB, accountID int, statusID int) {
	tx := database.BeginTx(dbConn)

	var sqls []string

	sql1 := fmt.Sprintf(
		"SELECT 1 AS one FROM favourites WHERE favourites.status_id = %d AND favourites.account_id = %d LIMIT %d;",
		statusID, accountID, 1)
	if exists := database.CheckExists(tx, sql1); exists == 0 {
		tx.Commit()
		return 
	}

	sql3 := fmt.Sprintf(
		"DELETE FROM favourites WHERE account_id = %d and status_ID = %d;",
		accountID, statusID)
	sqls = append(sqls, sql3)
	result := database.Execute(tx, sqls)
	if result {
		tx.Commit()
		updateStatusStats(dbConn, accountID, statusID, "favourites_count", -1)
	} else {
		tx.Rollback()
	}
	
}

func NewUser(dbConn *sql.DB) (int, error) {

	tx := database.BeginTx(dbConn)

	accountID := auxiliary.RandomNonnegativeInt()
	privateKey := auxiliary.RandStrSeq(32)
	publicKey := auxiliary.RandStrSeq(32)
	email := auxiliary.RandStrSeq(4) + "@" + auxiliary.RandStrSeq(5)
	username := auxiliary.RandStrSeq(10)
	password := auxiliary.RandStrSeq(25)
	
	var sqls []string

	sql1 := fmt.Sprintf(
		`INSERT INTO ACCOUNTS (id, username, private_key, public_key, created_at, updated_at, silenced, suspended, locked, protocol, memorial) 
		VALUES (%d, '%s', '%s', '%s', now(), now(), %t, %t, %t, %d, %t);`,
		accountID, username, privateKey, publicKey, false, false, false, 0, false)
	sql2 := fmt.Sprintf(
		`INSERT INTO USERS (email, created_at, updated_at, encrypted_password, remember_created_at, current_sign_in_ip, last_sign_in_ip, admin, account_id, disabled, moderator) 
		VALUES ('%s', now(), now(), '%s', now(), '%s', '%s', %t, %d, %t, %t);`,
		email, password, "127.0.0.1", "127.0.0.1", false, accountID, false, false)
	sql3 := fmt.Sprintf(
		`INSERT INTO ACCOUNT_STATS (account_id, statuses_count, following_count, followers_count, created_at, updated_at) 
		VALUES (%d, %d, %d, %d, now(), now());`,
		accountID, 0, 0, 0)
	sqls = append(sqls, sql1, sql2, sql3)
	result := database.Execute(tx, sqls)

	if result {
		tx.Commit()
		return accountID, nil
	} else {
		tx.Rollback()
		return -1, errors.New("Fail to create a new user")
	}
}

func CheckFollowed(dbConn *sql.DB, accountID1, accountID2 int) bool {

	query1 := fmt.Sprintf(
		"select id from follows where account_id = %d and target_account_id = %d",
		accountID2, accountID1,
	)

	res1 := database.DataCall(dbConn, query1)
	if len(res1) > 0 {
		return true
	} else {
		return false
	}
}

func GetFollowedUsers(dbConn *sql.DB, accountID int) []int {

	var users []int

	sql := fmt.Sprintf("SELECT account_id FROM follows WHERE target_account_id = %d;", accountID)

	res1 := database.DataCall(dbConn, sql)
	for _, row := range res1 {
		userID, _ := strconv.Atoi(row["account_id"])
		users = append(users, userID)
	}

	return users
}

func GetFollowedNum(dbConn *sql.DB, accountID int) int {

	query1 := fmt.Sprintf("select count(*) from follows where target_account_id = %d", accountID)

	res1 := database.DataCall(dbConn, query1)
	num1, err1 := strconv.Atoi(res1[0]["count"])
	if err1 != nil {
		log.Fatal(err1)
	}
	return num1

}

// Get the users that accountID follows
func GetFollowingUsers(dbConn *sql.DB, accountID int) []int {

	var users []int
	
	query1 := fmt.Sprintf(
		"select target_account_id from follows where account_id = %d", 
		accountID,
	)
	res1 := database.DataCall(dbConn, query1)
	for _, row := range res1 {
		pID, _ := strconv.Atoi(row["target_account_id"])
		users = append(users, pID)
	}

	return users
}

func Follow(dbConn *sql.DB, accountID int, targetAccountID int) {

	if accountID == targetAccountID {
		return 
	}

	tx := database.BeginTx(dbConn)

	sql1 := fmt.Sprintf(
		"SELECT 1 AS one FROM follows WHERE account_id = %d AND target_account_id = %d LIMIT %d;",
		accountID, targetAccountID, 1,
	)
	if exists := database.CheckExists(tx, sql1); exists == 1 {
		tx.Commit()
		return 
	}

	uri := "http://localhost:3000/" + auxiliary.RandStrSeq(32)
	var sqls1 []string
	sql4 := fmt.Sprintf(
		`INSERT INTO FOLLOWS (created_at, updated_at, account_id, target_account_id, uri) 
		VALUES (now(), now(), %d, %d, '%s');`,  
		accountID, targetAccountID, uri,
	)
	sqls1 = append(sqls1, sql4)
	result := database.Execute(tx, sqls1)
	if result {
		tx.Commit()
		updateAccountStats(dbConn, accountID, "following_count", 1)
		updateAccountStats(dbConn, targetAccountID, "followers_count", 1)
	} else {
		tx.Rollback()
	}
	
}

func Unfollow(dbConn *sql.DB, accountID int, targetAccountID int) {
	if accountID == targetAccountID {
		return 
	}
	
	tx := database.BeginTx(dbConn)
	sql1 := fmt.Sprintf(
		"SELECT 1 AS one FROM follows WHERE follows.account_id = %d AND follows.target_account_id = %d LIMIT %d;",accountID, targetAccountID, 1)
	if exists := database.CheckExists(tx, sql1); exists == 0 {
		tx.Commit()
		return 
	}
	var sqls []string
	sql2 := fmt.Sprintf("DELETE FROM follows WHERE follows.account_id = %d AND follows.target_account_id = %d",  
	accountID, targetAccountID)
	sqls = append(sqls, sql2)
	result := database.Execute(tx, sqls)

	if result {
		tx.Commit()
		updateAccountStats(dbConn, accountID, "following_count", -1)
		updateAccountStats(dbConn, targetAccountID, "followers_count", -1)
	} else {
		tx.Rollback()
	}
}

func ReplyToStatus(dbConn *sql.DB, replyToStatusID, accountID int, visibility int) int {
	
	var conversationID int
	var replyToAccountID int
	var sqls1 []string

	statusID := auxiliary.RandomNonnegativeInt()
	uri := "http://localhost:3000/users/admin/statuses/" + strconv.Itoa(statusID)
	hidden := false
	content := auxiliary.RandStrSeq(50)

	tx := database.BeginTx(dbConn)
	sql1 := fmt.Sprintf(`
		SELECT conversation_id, account_id FROM statuses WHERE statuses.id = %d LIMIT %d;`,
		replyToStatusID, 1,
	)
	rows, err := tx.Query(sql1)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&conversationID, &replyToAccountID); err != nil {
			log.Fatal(err)
		}
	}
	sql2 := fmt.Sprintf(
		`INSERT INTO statuses (id, text, created_at, updated_at, language, conversation_id, local, account_id, application_id, uri, in_reply_to_id, reply, in_reply_to_account_id, visibility) 
		VALUES (%d, '%s', now(), now(), '%s', %d, %t, %d, %d, '%s', %d, %t, %d, %d);`,  
		statusID, content, "en", conversationID, true, accountID, 1, uri, replyToStatusID, true, replyToAccountID, visibility)
	sql3 := insertIntoStreamEntries(statusID, accountID, hidden)
	sqls1 = append(sqls1, sql2, sql3)

	result := database.Execute(tx, sqls1)

	if result {
		tx.Commit()
		updateAccountStats(dbConn, accountID, "statuses_count", 1)
		updateStatusStats(dbConn, accountID, statusID, "replies_count", 1)
		return statusID
	} else {
		tx.Rollback()
		return -1
	}
	
}

func Reblog(dbConn *sql.DB, accountID int, reblogStatusID int) int {
	t := time.Now().Format(time.RFC3339)
	conversationID := auxiliary.RandomNonnegativeInt()
	statusID := auxiliary.RandomNonnegativeInt()
	uri := "http://localhost:3000/users/admin/statuses/" + strconv.Itoa(statusID)
	var sqls []string
	hidden := false

	tx := database.BeginTx(dbConn)
	sql1 := fmt.Sprintf(
		"SELECT 1 AS one FROM statuses WHERE statuses.reblog_of_id = %d AND statuses.account_id = %d LIMIT %d;", 
		statusID, accountID, 1)
	if exists := database.CheckExists(tx, sql1); exists == 1 {
		tx.Commit()
		return -1
	}
	sql2 := fmt.Sprintf(
		"INSERT INTO conversations (id, created_at, updated_at) VALUES (%d, '%s', '%s')", 
		conversationID, t, t)
	sql3 := fmt.Sprintf(
		"INSERT INTO statuses (id, created_at, updated_at, reblog_of_id, conversation_id, local, account_id, uri) VALUES (%d, '%s', '%s', %d, %d, %t, %d, '%s');",
		statusID, t, t, reblogStatusID, conversationID, true, accountID, uri)
	sql4 := insertIntoStreamEntries(statusID, accountID, hidden)
	sqls = append(sqls, sql2, sql3, sql4)
	result := database.Execute(tx, sqls)
	if result {
		tx.Commit()
		updateAccountStats(dbConn, accountID, "statuses_count", 1)
		updateStatusStats(dbConn, accountID, statusID, "reblogs_count", 1)
		return statusID
	} else {
		tx.Rollback()
		return -1
	}
	
}

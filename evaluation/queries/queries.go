/*
 * Query Generator for Mastodon
 */

package main

import (
    "database/sql"
    "fmt"
	"log"
	"time"
	"math/rand"
	"strconv"
    _ "github.com/lib/pq"
)

func randStrSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func randomNonnegativeInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2147483647)
}

func beginTx(dbConn *sql.DB) *sql.Tx {
	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	return tx
}

func execute(tx *sql.Tx, queries []string) {
	for _, query := range queries {
		fmt.Println(query)
		if _, err := tx.Exec(query); err != nil {
			log.Fatal(err)
		}	
	}
	tx.Commit()
}

func checkExists(tx *sql.Tx, query string) int {
	var exists int
	rows, err := tx.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&exists); err != nil {
			log.Fatal(err)
		}
	}
	return exists
}

func updateAccountStats(dbConn *sql.DB, accountID int, updateType string, difference int) {
	var statusNum int
	var	sqls []string 
	tx := beginTx(dbConn)
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
	t1 := time.Now().Format(time.RFC3339)
	sql2 := fmt.Sprintf(
		"UPDATE account_stats SET %s = %d, last_status_at = '%s', updated_at = '%s' WHERE account_stats.account_id = %d;",
		updateType, statusNum + difference, t1, t1, accountID) // These last_status_at and updated_at are generated randomly

	sqls = append(sqls, sql2)
	execute(tx, sqls)
	tx.Commit()
}

func updateStatusStats(dbConn *sql.DB, accountID int, statusID int, updateType string, difference int) {
	var sqls []string
	count := -1
	tx := beginTx(dbConn)
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
	execute(tx, sqls)
	tx.Commit()
}

func insertIntoStreamEntries(activityID int, t string, accountID int) string {
	return fmt.Sprintf(
		"INSERT INTO stream_entries (activity_id, activity_type, created_at, updated_at, account_id) VALUES (%d, '%s', '%s', '%s', %d);",
		activityID, "Status", t, t, accountID)
}

func publishStatus(dbConn *sql.DB, accountID int, content string) {
	t := time.Now().Format(time.RFC3339)
	activityID := randomNonnegativeInt()
	conversationID := randomNonnegativeInt()
	statusID := randomNonnegativeInt()
	uri := "http://localhost:3000/users/admin/statuses/" + strconv.Itoa(activityID)
	var sqls1 []string

	tx := beginTx(dbConn)
	sql1 := fmt.Sprintf(
		"INSERT INTO conversations (id, created_at, updated_at) VALUES (%d, '%s', '%s');", 
		conversationID, t, t)
	sql2 := fmt.Sprintf(
		"INSERT INTO statuses (id, text, created_at, updated_at, language, conversation_id, local, account_id, application_id, uri) VALUES (%d, '%s', '%s', '%s', '%s', %d, %t, %d, %d, '%s');",  
		statusID, content, t, t, "en", conversationID, true, accountID, 1, uri)
	sql3 := insertIntoStreamEntries(activityID, t, accountID)
	sqls1 = append(sqls1, sql1, sql2, sql3)
	execute(tx, sqls1)
	tx.Commit()

	updateAccountStats(dbConn, accountID, "statuses_count", 1)
}

func favourite(dbConn *sql.DB, accountID int, statusID int) {
	tx := beginTx(dbConn)

	var sqls []string

	sql1 := fmt.Sprintf(
		"SELECT 1 AS one FROM favourites WHERE favourites.status_id = %d AND favourites.account_id = %d LIMIT %d;",
		statusID, accountID, 1)
	if exists := checkExists(tx, sql1); exists == 1 {
		tx.Commit()
		return 
	}

	t := time.Now().Format(time.RFC3339)
	sql3 := fmt.Sprintf(
		"INSERT INTO favourites (created_at, updated_at, account_id, status_id) VALUES ('%s', '%s', %d, %d);",
		t, t, accountID, statusID)
	sqls = append(sqls, sql3)
	execute(tx, sqls)
	tx.Commit()

	updateStatusStats(dbConn, accountID, statusID, "favourites_count", 1)
}

func unfavourite(dbConn *sql.DB, accountID int, statusID int) {
	tx := beginTx(dbConn)

	var sqls []string

	sql1 := fmt.Sprintf(
		"SELECT 1 AS one FROM favourites WHERE favourites.status_id = %d AND favourites.account_id = %d LIMIT %d;",
		statusID, accountID, 1)
	if exists := checkExists(tx, sql1); exists == 0 {
		tx.Commit()
		return 
	}

	sql3 := fmt.Sprintf(
		"DELETE FROM favourites WHERE account_id = %d and status_ID = %d;",
		accountID, statusID)
	sqls = append(sqls, sql3)
	execute(tx, sqls)
	tx.Commit()

	updateStatusStats(dbConn, accountID, statusID, "favourites_count", -1)
}

func signup(dbConn *sql.DB, email string, username string, password string) {
	tx := beginTx(dbConn)

	t := time.Now().Format(time.RFC3339)
	accountID := randomNonnegativeInt()
	privateKey := randStrSeq(32)
	publicKey := randStrSeq(32)
	var sqls []string

	sql1 := fmt.Sprintf(
		"INSERT INTO ACCOUNTS (id, username, private_key, public_key, created_at, updated_at, silenced, suspended, locked, protocol, memorial) VALUES (%d, '%s', '%s', '%s', '%s', '%s', %t, %t, %t, %d, %t);",
		accountID, username, privateKey, publicKey, t, t, false, false, false, 0, false)
	sql2 := fmt.Sprintf(
		"INSERT INTO USERS (email, created_at, updated_at, encrypted_password, remember_created_at, current_sign_in_ip, last_sign_in_ip, admin, account_id, disabled, moderator) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, %d, %t, %t);",
		email, t, t, password, t, "127.0.0.1", "127.0.0.1", false, accountID, false, false)
	sql3 := fmt.Sprintf(
		"INSERT INTO ACCOUNT_STATS (account_id, statuses_count, following_count, followers_count, created_at, updated_at) VALUES (%d, %d, %d, %d, '%s', '%s');",
		accountID, 0, 0, 0, t, t)
	sqls = append(sqls, sql1, sql2, sql3)
	execute(tx, sqls)
}

func follow(dbConn *sql.DB, accountID int, targetAccountID int) {
	if accountID == targetAccountID {
		return 
	}

	tx := beginTx(dbConn)

	sql1 := fmt.Sprintf("SELECT 1 AS one FROM follows WHERE follows.account_id = %d AND follows.target_account_id = %d LIMIT %d;",
	accountID, targetAccountID, 1)
	if exists := checkExists(tx, sql1); exists == 1 {
		tx.Commit()
		return 
	}

	t := time.Now().Format(time.RFC3339)
	uri := "http://localhost:3000/" + randStrSeq(32)
	var sqls1 []string
	sql4 := fmt.Sprintf("INSERT INTO FOLLOWS (created_at, updated_at, account_id, target_account_id, uri) VALUES ('%s', '%s', %d, %d, '%s');",  
	t, t, accountID, targetAccountID, uri)
	sqls1 = append(sqls1, sql4)
	execute(tx, sqls1)
	tx.Commit()

	updateAccountStats(dbConn, accountID, "following_count", 1)
	updateAccountStats(dbConn, targetAccountID, "followers_count", 1)
}

func unfollow(dbConn *sql.DB, accountID int, targetAccountID int) {
	if accountID == targetAccountID {
		return 
	}
	
	tx := beginTx(dbConn)
	sql1 := fmt.Sprintf(
		"SELECT 1 AS one FROM follows WHERE follows.account_id = %d AND follows.target_account_id = %d LIMIT %d;",accountID, targetAccountID, 1)
	if exists := checkExists(tx, sql1); exists == 0 {
		tx.Commit()
		return 
	}
	var sqls []string
	sql2 := fmt.Sprintf("DELETE FROM follows WHERE follows.account_id = %d AND follows.target_account_id = %d",  
	accountID, targetAccountID)
	sqls = append(sqls, sql2)
	execute(tx, sqls)
	tx.Commit()

	updateAccountStats(dbConn, accountID, "following_count", -1)
	updateAccountStats(dbConn, targetAccountID, "followers_count", -1)
}

func replyToStatus(dbConn *sql.DB, accountID int, content string, replyToStatusID int) {
	t := time.Now().Format(time.RFC3339)
	activityID := randomNonnegativeInt()
	statusID := randomNonnegativeInt()
	uri := "http://localhost:3000/users/admin/statuses/" + strconv.Itoa(activityID)
	var conversationID int
	var replyToAccountID int
	var sqls1 []string

	tx := beginTx(dbConn)
	sql1 := fmt.Sprintf("SELECT conversation_id, account_id FROM statuses WHERE statuses.id = %d LIMIT %d;",
	replyToStatusID, 1)
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
		"INSERT INTO statuses (id, text, created_at, updated_at, language, conversation_id, local, account_id, application_id, uri, in_reply_to_id, reply, in_reply_to_account_id) VALUES (%d, '%s', '%s', '%s', '%s', %d, %t, %d, %d, '%s', %d, %t, %d);",  
		statusID, content, t, t, "en", conversationID, true, accountID, 1, uri, replyToStatusID, true, replyToAccountID)
	sql3 := insertIntoStreamEntries(activityID, t, accountID)
	sqls1 = append(sqls1, sql2, sql3)
	execute(tx, sqls1)
	tx.Commit()

	updateAccountStats(dbConn, accountID, "statuses_count", 1)
	updateStatusStats(dbConn, accountID, statusID, "replies_count", 1)
}

func reblog(dbConn *sql.DB, accountID int, reblogStatusID int) {
	t := time.Now().Format(time.RFC3339)
	activityID := randomNonnegativeInt()
	conversationID := randomNonnegativeInt()
	statusID := randomNonnegativeInt()
	uri := "http://localhost:3000/users/admin/statuses/" + strconv.Itoa(activityID)
	var sqls []string

	tx := beginTx(dbConn)
	sql1 := fmt.Sprintf(
		"SELECT 1 AS one FROM statuses WHERE statuses.reblog_of_id = %d AND statuses.account_id = %d LIMIT %d;", 
		statusID, accountID, 1)
	if exists := checkExists(tx, sql1); exists == 1 {
		tx.Commit()
		return 
	}
	sql2 := fmt.Sprintf(
		"INSERT INTO conversations (id, created_at, updated_at) VALUES (%d, '%s', '%s')", 
		conversationID, t, t)
	sql3 := fmt.Sprintf(
		"INSERT INTO statuses (id, created_at, updated_at, reblog_of_id, conversation_id, local, account_id, uri) VALUES (%d, '%s', '%s', %d, %d, %t, %d, '%s');",
		statusID, t, t, reblogStatusID, conversationID, true, accountID, uri)
	sql4 := insertIntoStreamEntries(activityID, t, accountID)
	sqls = append(sqls, sql2, sql3, sql4)
	execute(tx, sqls)
	tx.Commit()

	updateAccountStats(dbConn, accountID, "statuses_count", 1)
	updateStatusStats(dbConn, accountID, statusID, "reblogs_count", 1)
}

func main() {
    dbConn, err := sql.Open("postgres", "postgresql://root@10.230.12.75:26257/mastodon?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
	}
	
	// publishStatus(dbConn, 925840864, "okkkkkkkkkkkk")

	// favourite(dbConn, 925840864, 1389362391)

	// unfavourite(dbConn, 925840864, 1389362391)

	// signup(dbConn, "tai@nyu.edu", "zaincow", "cowcow")

	// follow(dbConn, 1217195077, 1042906640)

	// unfollow(dbConn, 1217195077, 1042906640)
	
	// replyToStatus(dbConn, 829522384, "a reply", 2042450516)

	reblog(dbConn, 735104489, 614615112)
}

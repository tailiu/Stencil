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

func publishStatus(dbConn *sql.DB, accountID int, content string) {
	t := time.Now().Format(time.RFC3339)
	activityID := randomNonnegativeInt()
	conversationID := randomNonnegativeInt()
	statusID := randomNonnegativeInt()
	uri := "http://localhost:3000/users/admin/statuses/" + strconv.Itoa(activityID)
	var statusNum int
	var sqls1 []string
	var sqls2 []string

	tx := beginTx(dbConn)
	sql1 := fmt.Sprintf(
		"INSERT INTO conversations (id, created_at, updated_at) VALUES (%d, '%s', '%s');", 
		conversationID, t, t)
	sql2 := fmt.Sprintf(
		"INSERT INTO statuses (id, text, created_at, updated_at, language, conversation_id, local, account_id, application_id, uri) VALUES (%d, '%s', '%s', '%s', '%s', %d, %t, %d, %d, '%s');",  
		statusID, content, t, t, "en", conversationID, true, accountID, 1, uri)
	sql3 := fmt.Sprintf(
		"INSERT INTO stream_entries (activity_id, activity_type, created_at, updated_at, account_id) VALUES (%d, '%s', '%s', '%s', %d);",
		activityID, "Status", t, t, accountID)
	sqls1 = append(sqls1, sql1, sql2, sql3)
	execute(tx, sqls1)
	tx.Commit()

	tx = beginTx(dbConn)
	sql4 := fmt.Sprintf(
		"SELECT statuses_count FROM account_stats WHERE account_stats.account_id = %d LIMIT %d;",
		accountID, 1)
	rows, err := tx.Query(sql4)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&statusNum); err != nil {
			log.Fatal(err)
		}
	}
	t1 := time.Now().Format(time.RFC3339)
	sql5 := fmt.Sprintf(
		"UPDATE account_stats SET statuses_count = %d, last_status_at = '%s', updated_at = '%s' WHERE account_stats.account_id = %d;",
		statusNum + 1, t1, t1, accountID) // These last_status_at and updated_at are generated randomly

	sqls2 = append(sqls2, sql5)
	execute(tx, sqls2)
	tx.Commit()
}

func favourite(dbConn *sql.DB, userID int, statusID int) {
	tx := beginTx(dbConn)

	var favouriteExist int
	var favouritesCount int 
	var sqls []string

	sql1 := fmt.Sprintf(
		"SELECT 1 AS one FROM favourites WHERE favourites.status_id = %d AND favourites.account_id = %d LIMIT %d;",
		statusID, userID, 1)
	rows, err := tx.Query(sql1)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&favouriteExist); err != nil {
			log.Fatal(err)
		}
	}
	if favouriteExist == 1 {
		tx.Commit()
		return 
	}

	sql2 := fmt.Sprintf(
		"SELECT status_stats.favourites_count FROM status_stats WHERE status_stats.status_id = %d LIMIT %d;",
		statusID, 1)
	rows, err = tx.Query(sql2)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&favouritesCount); err != nil {
			log.Fatal(err)
		}
	}

	t := time.Now().Format(time.RFC3339)
	sql3 := fmt.Sprintf(
		"INSERT INTO favourites (created_at, updated_at, account_id, status_id) VALUES ('%s', '%s', %d, %d);",
		t, t, userID, statusID)
	sql4 := fmt.Sprintf(
		"INSERT INTO status_stats (status_id, favourites_count, created_at, updated_at) VALUES (%d, %d, '%s', '%s');",
		statusID, favouritesCount + 1, t, t)
	sqls = append(sqls, sql3, sql4)
	execute(tx, sqls)
	tx.Commit()
}

func unfavourite(dbConn *sql.DB, userID int, statusID int) {
	tx := beginTx(dbConn)

	var favouriteExist int
	var favouritesCount int
	var sqls []string

	sql1 := fmt.Sprintf(
		"SELECT 1 AS one FROM favourites WHERE favourites.status_id = %d AND favourites.account_id = %d LIMIT %d;",
		statusID, userID, 1)
	rows, err := tx.Query(sql1)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&favouriteExist); err != nil {
			log.Fatal(err)
		}
	}
	if favouriteExist == 0 {
		tx.Commit()
		return 
	}

	sql2 := fmt.Sprintf(
		"SELECT status_stats.favourites_count FROM status_stats WHERE status_stats.status_id = %d LIMIT %d;",
		statusID, 1)
	rows, err = tx.Query(sql2)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&favouritesCount); err != nil {
			log.Fatal(err)
		}
	}

	sql3 := fmt.Sprintf(
		"DELETE FROM favourites WHERE account_id = %d and status_ID = %d;",
		userID, statusID)
	sql4 := fmt.Sprintf(
		"UPDATE status_stats SET favourites_count = %d WHERE status_ID = %d;",
		favouritesCount - 1, statusID)
	sqls = append(sqls, sql3, sql4)
	execute(tx, sqls)
	tx.Commit()
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

	var followExists int
	sql1 := fmt.Sprintf("SELECT 1 AS one FROM follows WHERE follows.account_id = %d AND follows.target_account_id = %d LIMIT %d;",
	accountID, targetAccountID, 1)
	rows, err := tx.Query(sql1)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&followExists); err != nil {
			log.Fatal(err)
		}
	}
	if followExists == 1 {
		tx.Commit()
		return 
	}

	t := time.Now().Format(time.RFC3339)
	uri := "http://localhost:3000/" + randStrSeq(32)
	var sqls1 []string
	var sqls2 []string
	var counts []int
	var count int

	sql2 := fmt.Sprintf("SELECT following_count FROM ACCOUNT_STATS WHERE account_stats.account_id = %d LIMIT %d;", accountID, 1)
	sql3 := fmt.Sprintf("SELECT followers_count FROM ACCOUNT_STATS WHERE account_stats.account_id = %d LIMIT %d;", targetAccountID, 1)
	sqls1 = append(sqls1, sql2, sql3)

	for _, sql := range sqls1 {
		fmt.Println(sql)
		rows, err = tx.Query(sql)
		if err != nil {
			log.Fatal(err)
		}
		for rows.Next() {
			if err := rows.Scan(&count); err != nil {
				log.Fatal(err)
			}
			counts = append(counts, count)
		}
	}

	sql4 := fmt.Sprintf("INSERT INTO FOLLOWS (created_at, updated_at, account_id, target_account_id, uri) VALUES ('%s', '%s', %d, %d, '%s');",  
	t, t, accountID, targetAccountID, uri)
	sql5 := fmt.Sprintf("UPDATE ACCOUNT_STATS SET following_count = %d, updated_at = '%s' WHERE account_stats.account_id = %d;",
	counts[0] + 1, t, accountID)
	sql6 := fmt.Sprintf("UPDATE ACCOUNT_STATS SET followers_count = %d, updated_at = '%s' WHERE account_stats.account_id = %d;",
	counts[1] + 1, t, targetAccountID)
	sqls2 = append(sqls2, sql4, sql5, sql6)
	execute(tx, sqls2)

	tx.Commit()
}

func unfollow(dbConn *sql.DB, accountID int, targetAccountID int) {
	if accountID == targetAccountID {
		return 
	}

	tx := beginTx(dbConn)

	var followExists int
	sql1 := fmt.Sprintf("SELECT 1 AS one FROM follows WHERE follows.account_id = %d AND follows.target_account_id = %d LIMIT %d;",
	accountID, targetAccountID, 1)
	rows, err := tx.Query(sql1)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&followExists); err != nil {
			log.Fatal(err)
		}
	}
	if followExists == 0 {
		tx.Commit()
		return
	}

	t := time.Now().Format(time.RFC3339)
	var sqls1 []string
	var sqls2 []string
	var counts []int
	var count int

	sql2 := fmt.Sprintf("SELECT following_count FROM ACCOUNT_STATS WHERE account_stats.account_id = %d LIMIT %d;", accountID, 1)
	sql3 := fmt.Sprintf("SELECT followers_count FROM ACCOUNT_STATS WHERE account_stats.account_id = %d LIMIT %d;", targetAccountID, 1)
	sqls1 = append(sqls1, sql2, sql3)

	for _, sql := range sqls1 {
		fmt.Println(sql)
		rows, err = tx.Query(sql)
		if err != nil {
			log.Fatal(err)
		}
		for rows.Next() {
			if err := rows.Scan(&count); err != nil {
				log.Fatal(err)
			}
			counts = append(counts, count)
		}
	}

	sql4 := fmt.Sprintf("DELETE FROM follows WHERE follows.account_id = %d AND follows.target_account_id = %d",  
	accountID, targetAccountID)
	sql5 := fmt.Sprintf("UPDATE ACCOUNT_STATS SET following_count = %d, updated_at = '%s' WHERE account_stats.account_id = %d;",
	counts[0] - 1, t, accountID)
	sql6 := fmt.Sprintf("UPDATE ACCOUNT_STATS SET followers_count = %d, updated_at = '%s' WHERE account_stats.account_id = %d;",
	counts[1] - 1, t, targetAccountID)
	sqls2 = append(sqls2, sql4, sql5, sql6)
	execute(tx, sqls2)

	tx.Commit()
}

func replyToStatus(dbConn *sql.DB, accountID int, content string, replyToStatusID int) {
	t := time.Now().Format(time.RFC3339)
	activityID := randomNonnegativeInt()
	statusID := randomNonnegativeInt()
	uri := "http://localhost:3000/users/admin/statuses/" + strconv.Itoa(activityID)
	var conversationID int
	var replyToAccountID int
	var statusNum int
	replyCount := -1
	var sqls1 []string
	var sqls2 []string
	var sqls3 []string

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
	sql3 := fmt.Sprintf(
		"INSERT INTO stream_entries (activity_id, activity_type, created_at, updated_at, account_id) VALUES (%d, '%s', '%s', '%s', %d);",
		activityID, "Status", t, t, accountID)
	sqls1 = append(sqls1, sql2, sql3)
	execute(tx, sqls1)
	tx.Commit()


	tx = beginTx(dbConn)
	sql4 := fmt.Sprintf(
		"SELECT statuses_count FROM account_stats WHERE account_stats.account_id = %d LIMIT %d;",
		accountID, 1)
	rows, err = tx.Query(sql4)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&statusNum); err != nil {
			log.Fatal(err)
		}
	}
	t1 := time.Now().Format(time.RFC3339)
	fmt.Println(statusNum)
	sql5 := fmt.Sprintf(
		"UPDATE account_stats SET statuses_count = %d, last_status_at = '%s', updated_at = '%s' WHERE account_stats.account_id = %d;",
		statusNum + 1, t1, t1, accountID) // These last_status_at and updated_at are generated randomly
	sqls2 = append(sqls2, sql5)
	execute(tx, sqls2)
	tx.Commit()


	tx = beginTx(dbConn)

	sql6 := fmt.Sprintf(
		"SELECT status_stats.replies_count FROM status_stats WHERE status_stats.status_id = %d LIMIT %d;",
		statusID, 1)
	rows, err = tx.Query(sql6)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&replyCount); err != nil {
			log.Fatal(err)
		}
	}
	var sql7 string
	if replyCount == -1 {
		sql7 = fmt.Sprintf(
			"INSERT INTO status_stats (status_id, replies_count, created_at, updated_at) VALUES (%d, %d, '%s', '%s');",
			statusID, 1, t, t)
	} else {
		sql7 = fmt.Sprintf(
			"UPDATE status_stats SET replies_count = %d, updated_at = '%s' WHERE status_id = %d,;",
			replyCount + 1, t, statusID)
	}
	sqls3 = append(sqls3, sql7)
	execute(tx, sqls3)
	tx.Commit()
}

func main() {
    dbConn, err := sql.Open("postgres", "postgresql://root@10.230.12.75:26257/mastodon?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
	}
	
	// publishStatus(dbConn, 925840864, "oooooooooooooooo")

	// favourite(dbConn, 1, 100584447)

	// unfavourite(dbConn, 1, 100584447)

	// signup(dbConn, "zainzainzainzainzainzain@nyu.edu", "zainzainzainzainzainzaincow", "cowcow")

	// follow(dbConn, 829522384, 1042906640)

	// unfollow(dbConn, 829522384, 1042906640)
	
	// replyToStatus(dbConn, 829522384, "jkjkjkjkjkj", 614615112)
}

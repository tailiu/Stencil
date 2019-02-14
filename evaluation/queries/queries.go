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

func randomNonnegativeInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2147483647)
}

func publishStatus(dbConn *sql.DB, userID int, content string) {
	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now().Format(time.RFC3339)
	activityID := randomNonnegativeInt()
	conversationID := randomNonnegativeInt()
	statusID := randomNonnegativeInt()
	uri := "http://localhost:3000/users/admin/statuses/" + strconv.Itoa(activityID)
	var statusNum int

	var sqls []string
	sql1 := fmt.Sprintf(
		"INSERT INTO conversations (id, created_at, updated_at) VALUES (%d, '%s', '%s');", 
		conversationID, t, t)
	sql2 := fmt.Sprintf(
		"INSERT INTO statuses (id, text, created_at, updated_at, language, conversation_id, local, account_id, application_id, uri) VALUES (%d, '%s', '%s', '%s', '%s', %d, %t, %d, %d, '%s');",  
		statusID, content, t, t, "en", conversationID, true, userID, 1, uri)
	sql3 := fmt.Sprintf(
		"INSERT INTO stream_entries (activity_id, activity_type, created_at, updated_at, account_id) VALUES (%d, '%s', '%s', '%s', %d);",
		activityID, "Status", t, t, userID)

	sql4 := fmt.Sprintf(
		"SELECT COUNT(id) FROM statuses WHERE account_id = %d;", userID)
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
		"UPDATE account_stats SET statuses_count = %d, last_status_at = '%s', updated_at = '%s' WHERE account_stats.id = %d;",
		statusNum + 1, t1, t1, userID) // These last_status_at and updated_at are generated randomly

	sqls = append(sqls, sql1, sql2, sql3, sql5)
	for _, sql := range sqls {
		fmt.Println(sql)
		if _, err := tx.Exec(sql); err != nil {
			log.Fatal(err)
		}	
	}
	tx.Commit()
}

func favourite(dbConn *sql.DB, userID int, statusID int) {
	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal(err)
	}

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
	for _, sql := range sqls {
		fmt.Println(sql)
		if _, err := tx.Exec(sql); err != nil {
			log.Fatal(err)
		}	
	}
	tx.Commit()
}

func unfavourite(dbConn *sql.DB, userID int, statusID int) {
	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal(err)
	}

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
	for _, sql := range sqls {
		fmt.Println(sql)
		if _, err := tx.Exec(sql); err != nil {
			log.Fatal(err)
		}	
	}
	tx.Commit()
}

func signUp(dbConn *sql.DB, email string, username string, password string) {
	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now().Format(time.RFC3339)
	userID := randomNonnegativeInt()
	privateKey := strconv.Itoa(randomNonnegativeInt())
	publicKey := strconv.Itoa(randomNonnegativeInt())
	var sqls []string

	sql1 := fmt.Sprintf(
		"INSERT INTO ACCOUNTS (id, username, private_key, public_key, created_at, updated_at, silenced, suspended, locked, protocol, memorial) VALUES (%d, '%s', '%s', '%s', '%s', '%s', %t, %t, %t, %d, %t);",
		userID, username, privateKey, publicKey, t, t, false, false, false, 0, false)
	sql2 := fmt.Sprintf(
		"INSERT INTO USERS (id, email, created_at, updated_at, encrypted_password, remember_created_at, current_sign_in_ip, last_sign_in_ip, admin, account_id, disabled, moderator) VALUES (%d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, %d, %t, %t);",
		userID, email, t, t, password, t, "127.0.0.1", "127.0.0.1", false, userID, false, false)
	sqls = append(sqls, sql1, sql2)
	for _, sql := range sqls {
		fmt.Println(sql)
		if _, err := tx.Exec(sql); err != nil {
			log.Fatal(err)
		}	
	}
	tx.Commit()
}

func main() {
    dbConn, err := sql.Open("postgres", "postgresql://root@10.230.12.75:26257/mastodon?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
	}
	
	// publishStatus(dbConn, 1, "gogogogogogogo")

	// favourite(dbConn, 1, 100584447)

	// unfavourite(dbConn, 1, 100584447)

	signUp(dbConn, "zain@nyu.edu", "zaincow", "cowcow")
}

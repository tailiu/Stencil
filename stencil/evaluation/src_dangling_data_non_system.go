package evaluation

import (
	"log"
	"fmt"
	"stencil/db"
	"database/sql"
)

func getCounts(dbConn *sql.DB, query string) int64 {

	data, err := db.DataCall1(dbConn, query)

	if err != nil {
		log.Fatal(err)
	}

	return data["count"].(int64)
}

func getDanglingLikesNum(evalConfig *EvalConfig, 
	danglingDataStats map[string]int64, pKey int) {

	key := "likes:posts"

	query := fmt.Sprintf(
		`SELECT count(*) from likes 
		where target_id = %d and mark_as_delete = false`,
		pKey)

	data := getCounts(evalConfig.DiasporaDBConn, query)

	if data != 0 {
		danglingDataStats[key] = data
		log.Println("Got Dangling Data!")
	}	
}

func getDanglingCommentsNum(evalConfig *EvalConfig, 
	danglingDataStats map[string]int64, pKey int) {
	
	key := "comments:posts"
	
	query := fmt.Sprintf(
		`SELECT count(*) from comments 
		where commentable_id = %d and mark_as_delete = false`,
		pKey)
	
	data := getCounts(evalConfig.DiasporaDBConn, query)
	if data != 0 {
		danglingDataStats[key] = data
		log.Println("Got Dangling Data!")
	}	
}

func getDanglingMessagesNum(evalConfig *EvalConfig, 
	danglingDataStats map[string]int64, pKey int) {
	
	key := "messages:conversations"

	query := fmt.Sprintf(
		`SELECT count(*) from messages 
		where conversation_id = %d and mark_as_delete = false`,
		pKey)
	
	data := getCounts(evalConfig.DiasporaDBConn, query)

	if data != 0 {
		danglingDataStats[key] = data
		log.Println("Got Dangling Data!")
	}
}

func srcDanglingDataNonSystem(evalConfig *EvalConfig, 
	migrationID string, table string, pKey int) map[string]int64 {
	
	danglingDataStats := make(map[string]int64)

	log.Println("Dangling data to check:", table, pKey)

	if table == "comments" {
		getDanglingLikesNum(evalConfig, danglingDataStats, pKey)
	} else if table == "posts" {
		getDanglingLikesNum(evalConfig, danglingDataStats, pKey)
		getDanglingCommentsNum(evalConfig, danglingDataStats, pKey)
	} else if table == "conversations" {
		getDanglingMessagesNum(evalConfig, danglingDataStats, pKey)
	}

	return danglingDataStats
}
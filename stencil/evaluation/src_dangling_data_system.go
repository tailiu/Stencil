package evaluation

import (
	"log"
	"fmt"
	"stencil/db"
	"database/sql"
)

func getCountsSystem(dbConn *sql.DB, query string) int64 {
	data, err := db.DataCall(dbConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return data[0]["count"].(int64)
}

func getDanglingLikesNumSystem(evalConfig *EvalConfig, danglingDataStats map[string]int64) {
	key := "likes:posts"
	query := fmt.Sprintf("SELECT count(*) from likes join posts on likes.target_id = posts.id where likes.mark_as_delete = false and posts.mark_as_delete = true;")
	data := getCountsSystem(evalConfig.DiasporaDBConn, query)
	if data != 0 {
		danglingDataStats[key] = data
	}	
}

func getDanglingCommentsNumSystem(evalConfig *EvalConfig, danglingDataStats map[string]int64) {
	key := "comments:posts"
	query := fmt.Sprintf("SELECT count(*) from comments join posts on commentable_id = posts.id where comments.mark_as_delete = false and posts.mark_as_delete = true;")
	data := getCountsSystem(evalConfig.DiasporaDBConn, query)
	if data != 0 {
		danglingDataStats[key] = data
	}	
}

func getDanglingMessagesNumSystem(evalConfig *EvalConfig, danglingDataStats map[string]int64) {
	key := "messages:conversations"
	query := fmt.Sprintf("SELECT count(*) from messages join conversations on messages.conversation_id = conversations.id where messages.mark_as_delete = false and conversations.mark_as_delete = true;")
	data := getCountsSystem(evalConfig.DiasporaDBConn, query)
	if data != 0 {
		danglingDataStats[key] = data
	}
}

func srcDanglingDataSystem(evalConfig *EvalConfig) map[string]int64 {
	danglingDataStats := make(map[string]int64)

	getDanglingLikesNumSystem(evalConfig, danglingDataStats)
	getDanglingCommentsNumSystem(evalConfig, danglingDataStats)
	getDanglingMessagesNumSystem(evalConfig, danglingDataStats)

	return danglingDataStats
}
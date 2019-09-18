package evaluation

import (
	"fmt"
	"stencil/db"
	"log"
)

func getDanglingfavouritesNumSystemWithoutIntegration(evalConfig *EvalConfig, danglingDataStats map[string]int64, userID int64) {
	key := "favourites:statuses"
	query := fmt.Sprintf("SELECT count(*) FROM favourites WHERE status_id NOT IN (SELECT id FROM statuses) and account_id = %d;",
		userID)
	data := getCountsSystem(evalConfig.MastodonDBConn, query)
	if data != 0 {
		danglingDataStats[key] = data
	}	
}

func getDanglingfavouritesNumSystemWithIntegration(evalConfig *EvalConfig, danglingDataStats map[string]int64) {
	key := "favourites:statuses(integration)"
	query := fmt.Sprintf("SELECT count(*) FROM favourites WHERE status_id NOT IN (SELECT id FROM statuses);")
	data := getCountsSystem(evalConfig.MastodonDBConn, query)
	if data != 0 {
		danglingDataStats[key] = data
	}	
}

func getDanglingStatusesNumSystem(evalConfig *EvalConfig, danglingDataStats map[string]int64) {
	key := "statuses:conversations"
	query := fmt.Sprintf("SELECT count(*) FROM statuses WHERE conversation_id NOT IN (SELECT id FROM conversations);")
	data := getCountsSystem(evalConfig.MastodonDBConn, query)
	if data != 0 {
		danglingDataStats[key] = data
	}	
}

func getUserID(evalConfig *EvalConfig, migrationID string) int64 {
	query := fmt.Sprintf("SELECT user_id from migration_registration where migration_id = %s", migrationID)
	data1, err1 := db.DataCall1(evalConfig.StencilDBConn, query)
	if err1 != nil {
		log.Fatal(err1)
	}
	return data1["user_id"].(int64)
}

func dstDanglingDataSystem(evalConfig *EvalConfig, migrationID string) map[string]int64 {
	danglingDataStats := make(map[string]int64)

	userID := getUserID(evalConfig, migrationID)
	getDanglingfavouritesNumSystemWithoutIntegration(evalConfig, danglingDataStats, userID)
	// getDanglingfavouritesNumSystemWithIntegration(evalConfig, danglingDataStats)
	// getDanglingStatusesNumSystem(evalConfig, danglingDataStats)
	
	return danglingDataStats
}
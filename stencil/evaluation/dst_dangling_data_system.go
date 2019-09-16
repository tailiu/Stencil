package evaluation

import (
	"fmt"
)

func getDanglingfavouritesNumSystem(evalConfig *EvalConfig, danglingDataStats map[string]int64) {
	key := "favourites:statuses"
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

func dstDanglingDataSystem(evalConfig *EvalConfig) map[string]int64 {
	danglingDataStats := make(map[string]int64)

	getDanglingfavouritesNumSystem(evalConfig, danglingDataStats)
	getDanglingStatusesNumSystem(evalConfig, danglingDataStats)
	
	return danglingDataStats
}
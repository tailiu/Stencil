package evaluation

import (
	"stencil/db"
	"stencil/config"
	"stencil/qr"
	"database/sql"
	"log"
	"fmt"
	// "strings"
)

func getColsSizeOfDataInStencilModel(evalConfig *EvalConfig, dstAppConfig *config.AppConfig, tableName string, rowIDs []string) map[string]interface{} {
	qs := qr.CreateQS(dstAppConfig.QR)
	qs.FromTable(map[string]string{"table":tableName, "mflag": "0", "mark_as_delete": "false", "bag": "false"})
	qs.SelectColumns(tableName + ".*")
	var strRowIDs string 
	for i, rowID := range rowIDs {
		if i == 0 {
			strRowIDs += rowID
		} else {
			strRowIDs += "," + rowID
		}
	}
	qs.RowIDs(strRowIDs)
	physicalQuery := qs.GenSQLSize()

	// log.Println(physicalQuery)

	result, err := db.DataCall1(evalConfig.StencilDBConn, physicalQuery)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)

	return result
}

func getData1FromPhysicalSchema(stencilDBConn *sql.DB, QR *qr.QR, appID, cols, from, col, op, val string) map[string]interface{}  {	
	qs := qr.CreateQS(QR)
	qs.SelectColumns(cols)
	// Note that we don't care about mflag here
	qs.FromTable(map[string]string{"table":from, "mflag": "0,1"})
	qs.AddWhereWithValue(col, op, val)
	physicalQuery := qs.GenSQL()
	// log.Println(physicalQuery)

	result, err := db.DataCall1(stencilDBConn, physicalQuery)
	if err != nil {
		log.Fatal(err)
	}
	
	return result
}

func filterColsAndResultsBasedOnSchemaMapping(data map[string]interface{}, evalConfig *EvalConfig, dstAppConfig *config.AppConfig, tableName, srcApp, dstApp string) int64 {
	var size int64
	for _, v := range data {
		// if strings.Contains(k, ".mark_as_delete") {
		// 	continue
		if v == nil {
			continue
		} else {
			size += v.(int64)
		}
		// } else if srcApp == "diaspora" && dstApp == "mastodon" {
		// 	if tableName == "status_stats" {
		// 		break
		// 	} else if tableName == "conversations" {
		// 		data1 := getData1FromPhysicalSchema(evalConfig.StencilDBConn, dstAppConfig.QR, dstAppConfig.AppID, "statuses.*", "statuses", 
		// 			"statuses.id", "=", fmt.Sprint(data["conversations.id"]))
		// 		if len(data1) == 0 {
		// 			continue
		// 		} else {
		// 			break
		// 		}
		// 	} else {
		// 		size += v.(int64)
		// 	}
		// } else if srcApp == "diaspora" && dstApp == "twitter" {
		// 	if tableName == "conversation_participants" && k == "conversation_participants.role" {
		// 		continue
		// 	} else if tableName == "user_actions" && k == "user_actions.action_type" {
		// 		continue
		// 	} else {
		// 		size += v.(int64)
		// 	}
		// } else if srcApp == "diaspora" && dstApp == "gnusocial" {
		// 	if tableName == "conversation_participants" && k == "conversation_participants.role" {
		// 		continue
		// 	} else if tableName == "profile" && (k == "profile.nickname" || k == "profile.id") {
		// 		continue
		// 	} else {
		// 		size += v.(int64)
		// 	}
		// } else if srcApp == "mastodon" && dstApp == "diaspora" {
		// 	if strings.Contains(k, ".guid") || strings.Contains(k, ".commentable_type") || strings.Contains(k, ".target_type") {
		// 		continue
		// 	} else if tableName == "notification_actors" && (strings.Contains(k, ".created_at") || strings.Contains(k, ".updated_at")) {
		// 		continue
		// 	} else {
		// 		size += v.(int64)
		// 	}
		// } else if srcApp == "mastodon" && dstApp == "twitter" {
		// 	if tableName == "credentials" && (strings.Contains(k, ".id") || strings.Contains(k, ".user_id")) {
		// 		continue
		// 	} else if strings.Contains(k, ".action_type") {
		// 		continue
		// 	} else {
		// 		size += v.(int64)
		// 	}
		// } else if srcApp == "mastodon" && dstApp == "gnusocial" { 
		// 	if tableName == "profile" && strings.Contains(k, ".id") {
		// 		continue
		// 	} else {
		// 		size += v.(int64)
		// 	}
		// } else if srcApp == "twitter" && dstApp == "mastodon" {
		// 	if tableName == "users" && (strings.Contains(k, ".id") || strings.Contains(k, ".account_id")) {
		// 		continue
		// 	} else if tableName == "media_attachments" {
		// 		break
		// 	} else {
		// 		size += v.(int64)
		// 	}
		// } else if srcApp == "twitter" && dstApp == "diaspora" { 
		// 	if tableName == "profiles" {
		// 		break
		// 	} else if tableName == "posts" && strings.Contains(k, ".type") {
		// 		continue
		// 	} else if strings.Contains(k, ".receiving") || strings.Contains(k, ".sharing"){
		// 		continue
		// 	} else if tableName == "notification_actors" && (strings.Contains(k, ".created_at") || strings.Contains(k, ".updated_at")) {
		// 		continue 
		// 	} else {
		// 		size += v.(int64)
		// 	}
		// } else if srcApp == "twitter" && dstApp == "gnusocial" {
		// 	if tableName == "profile" && (k == "profile.nickname" || k == "profile.id") {
		// 		continue
		// 	} else {
		// 		size += v.(int64)
		// 	}
		// } else {
		// 	size += v.(int64)
		// }
	}
	return size
}

func calculateDisplayedDataSizeInBagEvaluation(evalConfig *EvalConfig, dstAppConfig *config.AppConfig, srcApp, dstApp string, displayedData []DisplayedData) int64 {
	var size int64
	for _, data := range displayedData {
		tableName := GetTableNameByTableID(evalConfig, data.TableID)
		size += filterColsAndResultsBasedOnSchemaMapping(getColsSizeOfDataInStencilModel(evalConfig, dstAppConfig, tableName, data.RowIDs), evalConfig, dstAppConfig, tableName, srcApp, dstApp)
	}
	return size
}

func getDisplayedDataSize(evalConfig *EvalConfig, srcApp, dstApp, migrationID string) int64 {
	dstAppConfig := getAppConfig(evalConfig, dstApp)
	displayedData := getAllDisplayedData(evalConfig, migrationID, dstAppConfig.AppID)
	log.Println(displayedData)
	return calculateDisplayedDataSizeInBagEvaluation(evalConfig, dstAppConfig, srcApp, dstApp, displayedData)
}

// We use dstApp here to get the total migrated node size in the source application
func getTotalMigratedNodeSize(evalConfig *EvalConfig, dstApp string, migrationID string) int64 {
	dstAppID := db.GetAppIDByAppName(evalConfig.StencilDBConn, dstApp)
	query := fmt.Sprintf("select msize from migration_registration where dst_app = %s and migration_id = %s", dstAppID, migrationID)
	result, err := db.DataCall1(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)

	return result["msize"].(int64)
}
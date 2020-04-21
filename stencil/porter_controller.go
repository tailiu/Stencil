package main

import (
	"fmt"
	"log"
	"stencil/apis"
	"stencil/db"
	"strings"
)

func checkIfFuzool(table string) bool {
	fuzoolTables := []string{"blocks", "schema_migrations", "ar_internal_metadata", "pods", "mentions", "o_embed_caches", "user_preferences", "chat_offline_messages", "simple_captcha_data", "comment_signatures", "o_auth_applications", "signature_orders", "o_auth_access_tokens", "account_deletions", "account_migrations", "authorizations", "poll_participations", "services", "open_graph_caches", "participations", "invitation_codes", "polls", "ppid", "references", "reports", "aspect_memberships", "poll_answers", "roles", "chat_contacts", "like_signatures", "poll_participation_signatures", "tag_followings", "tags", "taggings", "chat_fragments", "locations", "aspects", "contacts", "comments", "conversations", "messages", "notifications", "notification_actors"}
	for _, fTable := range fuzoolTables {
		if strings.EqualFold(fTable, table) {
			return true
		}
	}
	return false
}

func main() {
	limit := int64(1000)
	appName, appID := "diaspora", "1"
	appDB, stencilDB := db.GetDBConn(appName), db.GetDBConn("stencil")
	tables := db.GetTablesOfDB(appDB, appName)
	fmt.Println(tables)
	for _, table := range tables {
		if checkIfFuzool(table) {
			continue
		}
		if totalRows, err := db.GetRowCount(appDB, table); err == nil {
			for offset := int64(0); offset <= totalRows; offset += limit {
				log.Printf("Table: %s | Porting rows => %v - %v out of %v\n", table, offset, offset+limit, totalRows)
				apis.Port(appName, appID, table, limit, offset, appDB, stencilDB)
			}
		} else {
			log.Fatal("main.GetRowCount: ", err)
		}
	}
	appDB.Close()
	stencilDB.Close()
}

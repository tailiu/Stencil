package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"stencil/db"
	"stencil/qr"
	"strings"
	"sync"
	"time"
	// escape "github.com/tj/go-pg-escape"
)

func FilterTablesFromList(tables []string, tablesToRemove []string) []string {
	var filteredTables []string

	for _, table := range tables {
		remove := false
		for _, tableToRemove := range tablesToRemove {
			if strings.EqualFold(table, tableToRemove) {
				remove = true
				break
			}
		}
		if !remove {
			filteredTables = append(filteredTables, table)
		}
	}

	return filteredTables
}

func runTx(dbConn *sql.DB, QIs []*qr.QI) bool {
	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal("transaction can't even begin")
	}

	success := true

	for _, qi := range QIs {
		query, args := qi.GenSQL()
		// fmt.Println(query)
		if _, err := tx.Exec(query, args...); err != nil {
			success = false
			fmt.Println("Some error:", err)
			fmt.Println(query, args)
			fmt.Println(qi)
			break
		}
	}

	if success {
		tx.Commit()
	} else {
		tx.Rollback()
		fmt.Println("QIs :=v")
		fmt.Println(QIs)
		log.Fatal()
	}

	return success
}

func _transfer(QR *qr.QR, appDB, stencilDB *sql.DB, table string, limit, offset int64) {
	q := fmt.Sprintf("SELECT * FROM \"%s\" ORDER BY id LIMIT %d OFFSET %d", table, limit, offset)
	if ldata, err := db.DataCall(appDB, q); err != nil {
		fmt.Println(q)
		log.Fatal("Some problem with logical data query:", err)
	} else {
		for _, ldatum := range ldata {
			var cols []string
			var vals []interface{}
			for col, val := range ldatum {
				cols, vals = append(cols, col), append(vals, val)
			}
			qi := qr.CreateQI(table, cols, vals, qr.QTInsert)
			rowid := db.GetNewRowID(stencilDB)
			qis := QR.ResolveInsert(qi, rowid)
			runTx(stencilDB, qis)
		}
	}
}

func transfer(QR *qr.QR, appDB, stencilDB *sql.DB, table string, wg *sync.WaitGroup) {

	log.Println("Populating ", table)
	if totalRows, err := db.GetRowCount(appDB, table); err == nil {
		limit := int64(25000)
		for offset := int64(0); offset < totalRows; offset += limit {
			log.Println(fmt.Sprintf(">> %s: %d - %d of %d | Remaining: %d", table, offset, offset+limit, totalRows, totalRows-offset))
			_transfer(QR, appDB, stencilDB, table, limit, offset)
		}
	} else {
		log.Fatal("Error while fetching total rows", err)
	}
	log.Println("Done:", table)
	if wg != nil {
		wg.Done()
	}
}

func checkIfFuzool(table string) bool {
	fuzoolTables := []string{"blocks", "schema_migrations", "ar_internal_metadata", "pods", "mentions", "o_embed_caches", "user_preferences", "chat_offline_messages", "simple_captcha_data", "comment_signatures", "o_auth_applications", "signature_orders", "o_auth_access_tokens", "account_deletions", "account_migrations", "authorizations", "poll_participations", "services", "open_graph_caches", "participations", "invitation_codes", "polls", "ppid", "references", "reports", "aspect_memberships", "poll_answers", "roles", "chat_contacts", "like_signatures", "poll_participation_signatures", "tag_followings", "tags", "taggings", "chat_fragments", "locations"}
	for _, fTable := range fuzoolTables {
		if strings.EqualFold(fTable, table) {
			return true
		}
	}
	return false
}

func main() {
	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())

	appName, appID := "diaspora", "1"
	// stencilDB := db.GetDBConn(db.STENCIL_DB)
	appDB := db.GetDBConn(appName)
	QR := qr.NewQR(appName, appID)
	tables := db.GetTablesOfDB(appDB, appName)
	// tables = FilterTablesFromList(tables, []string{"messages"})
	// tables := []string{"messages"}
	// log.Fatal(tables)
	// tables = []string{"comments"}
	current_threads := 0
	for _, table := range tables {
		if checkIfFuzool(table) {
			continue
		}
		// transfer(QR, appDB, stencilDB, table, nil)
		wg.Add(1)
		current_threads++
		go transfer(QR, db.GetDBConn(appName), db.GetDBConn(db.STENCIL_DB), table, &wg)
		if current_threads > 2 {
			wg.Wait()
			current_threads = 0
		}
	}
	wg.Wait()
}

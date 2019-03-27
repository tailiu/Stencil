package main

import (
	"mastodon/auxiliary"
	"transaction/atomicity"
	"encoding/json"
	"log"
	"fmt"
	"strconv"
	"transaction/display"
	"transaction/db"
)

const StencilDBName = "stencil"
const logEntryNum = 100
const maxRowsPerHint = 10 
var tables = [...]string{"accounts", "account_stats", "conversations", "favourites", "follows",
						"media_attachments", "mentions", "statuses", "status_stats", "stream_entries", "users"}


type HintStruct struct {
	Table 		string
	Key 		string
	Value 		string
	ValueType 	string
} 

func genSerializedDisplayHint() []byte {
	var displayHints []HintStruct

	dbConn := db.GetDBConn(StencilDBName)

	// Generate 1 - 10 rows for each hint
	rowNum := auxiliary.RandomNonnegativeIntWithUpperBound(maxRowsPerHint) + 1
	for i := 0; i < rowNum; i ++ {
		tableNum := auxiliary.RandomNonnegativeIntWithUpperBound(len(tables))
		table := tables[tableNum]
		id := auxiliary.RandomNonnegativeInt()

		var hint = HintStruct {
			Table: 		table,
			Key: 		"account_id",
			Value: 		strconv.Itoa(id),
			ValueType: 	"int",
		}
		
		display.GenDisplayFlag(dbConn, id, table, false)
		
		displayHints = append(displayHints, hint)
	}

	encodedData, err := json.Marshal(displayHints)
	if err != nil {
		log.Fatal("Encoding errors!")
	}

	return encodedData
}

func genDummyMigrationLogs() {
	log_txn := atomicity.BeginTransaction()

	var undo_action string
	var display_hint []byte
	for i := 0; i < logEntryNum; i++ {
		undo_action = auxiliary.RandStrSeq(20)
		display_hint = genSerializedDisplayHint()
		atomicity.LogChange(undo_action, display_hint, log_txn)
	}

	atomicity.LogOutcome(log_txn, "COMMIT")
}

func main() {
	// genDummyMigrationLogs()


	// display.UpdateDisplayFlag(dbConn, 1, "account_id", true)
	dbConn := db.GetDBConn(StencilDBName)
	// display.CreateDisplayFlagsTable(dbConn)

	display_flag, err := display.CheckDisplayFlag(dbConn, 1233, "jjjj")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(display_flag)
	}
}
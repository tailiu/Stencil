package SA2_display

import (
	"database/sql"
	"log"
	"fmt"
	"stencil/db"
)

func CreateDeletionHoldTable(dbConn *sql.DB) {
	op := `CREATE TABLE deletion_hold (
			row_id int NOT NULL,
			thread_id int NOT NULL,
			hold boolean NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL);
			CREATE INDEX idx_row_id
				ON deletion_hold(row_id);
			CREATE INDEX idx_thread_id
				ON deletion_hold(thread_id);`
	if _, err := dbConn.Exec(op); err != nil {
		log.Fatal(err)
	}
}

func RemoveDeletionHold(stencilDBConn *sql.DB, dhStack [][]int, threadID int) {
	for len(dhStack) > 0 {
		n := len(dhStack) - 1
		hintRowIDGroup := dhStack[n]
		
		var queries []string
		for _, hintRowID := range hintRowIDGroup {
			query := fmt.Sprintf("UPDATE deletion_hold SET hold = %t, updated_at = now() WHERE row_id = %d and thread_id = %d;",
				false, hintRowID, threadID)
			log.Println("**************************************")
			log.Println(query)
			log.Println("**************************************")
			queries = append(queries, query)
		}
		if err := db.TxnExecute(stencilDBConn, queries); err != nil {
			log.Fatal(err)
		}
		dhStack = dhStack[:n]
	}
}

func AddToDeletionHoldStack(dhStack [][]int, dataHints []HintStruct, threadID int) ([]string, [][]int) {
	var	hintRowIDs []int 
	var queries []string

	for _, dataHint := range dataHints {
		for _, rowID := range dataHint.RowIDs {
			hintRowIDs = append(hintRowIDs, rowID)
			query := fmt.Sprintf("INSERT INTO deletion_hold (row_id, thread_id, hold, created_at, updated_at) VALUES (%d, %d, %t, now(), now());",
				rowID, threadID, true)
			log.Println("**************************************")
			log.Println(query)
			log.Println("**************************************")
			queries = append(queries, query)
		}
	}

	dhStack = append(dhStack, hintRowIDs)

	// log.Println("&&&&&&&&&&&")
	// log.Println(dhStack)
	// log.Println("&&&&&&&&&&&")

	return queries, dhStack
}
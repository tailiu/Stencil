package main

import (
	"stencil/reference_resolution"
	"stencil/app_display"
)

func PublishStatus(dbConn *sql.DB, accountID int, content string, haveMedia bool, visibility int, mentionedAccounts []int) int {

	t := time.Now().Format(time.RFC3339)
	conversationID := auxiliary.RandomNonnegativeInt()
	statusID := auxiliary.RandomNonnegativeInt()
	
	uri := "http://localhost:3000/users/admin/statuses/" + strconv.Itoa(statusID)
	var sqls1 []string
	var hidden bool

	if visibility == 0 {
		hidden = false
	} else if visibility == 3 {
		hidden = true
	}

	tx := database.BeginTx(dbConn)
	
	sql1 := fmt.Sprintf(
		"INSERT INTO conversations (id, created_at, updated_at) VALUES (%d, '%s', '%s');", 
		conversationID, t, t)
	
	sql2 := fmt.Sprintf(
		"INSERT INTO statuses (id, text, created_at, updated_at, language, conversation_id, local, account_id, application_id, uri, visibility) VALUES (%d, '%s', '%s', '%s', '%s', %d, %t, %d, %d, '%s', %d);",  
		statusID, content, t, t, "en", conversationID, true, accountID, 1, uri, visibility)
	
	sql3 := insertIntoStreamEntries(statusID, t, accountID, hidden)

	sqls1 = append(sqls1, sql1, sql2, sql3)
	
	if haveMedia {

		file_file_name := auxiliary.RandStrSeq(20) + "jpeg"
		file_content_type := "image/jpeg"
		file_file_size := auxiliary.RandomNonnegativeInt()
		shortCode := auxiliary.RandStrSeq(20)
		file_meta, err := json.Marshal(FileMeta{auxiliary.RandomNonnegativeInt(), auxiliary.RandomNonnegativeInt(), auxiliary.RandomNonnegativeInt(), auxiliary.RandomNonnegativeInt()})
		if err != nil {
			log.Fatal(err)
		}

		sql4 := fmt.Sprintf("INSERT INTO media_attachments (status_id, file_file_name, file_content_type, file_file_size, file_updated_at, created_at, updated_at, shortcode, file_meta, account_id) VALUES (%d, '%s', '%s', %d, '%s', '%s', '%s', '%s', '%s', %d);",
		statusID, file_file_name, file_content_type, file_file_size, t, t, t, shortCode, file_meta, accountID)
		
		sqls1 = append(sqls1, sql4)

	}

	if len(mentionedAccounts) != 0 {

		var sql5 string

		for _, mentionedAccount := range mentionedAccounts {
			sql5 = fmt.Sprintf("INSERT INTO mentions (status_id, created_at, updated_at, account_id) VALUES (%d, '%s', '%s', %d);",
				statusID, t, t, mentionedAccount)
			sqls1 = append(sqls1, sql5)
		}

	}

	result := database.Execute(tx, sqls1)
	
	if result {

		tx.Commit()

		updateAccountStats(dbConn, accountID, "statuses_count", 1)
		
		return statusID
	
	} else {
		
		tx.Rollback()
		
		return -1

	}
}

func preTest1() {

	PublishStatus()

}

/**
 *
 * Diaspora -> Mastodon
 * 
 * Identity:
 * 	a like (id:12) in Diaspora likes (id:8) -> a favourite (id:24) in Mastodon favourite (id:72)
 * 	a status (id:40) in Diaspora statuses (id:37) -> a post (id:90) in Mastodon posts (id:92)
 *
 * Reference:
 *	Diaspora (1), likes (id:8), like (id:12), target_id -> Mastodon (2), statuses (id:37), status (id:40), id
 *
**/
func test1(app string, migrationID int) {

	displayConfig := app_display.CreateDisplayConfig(app, migrationID, true)

	var hint = app_display.HintStruct{
		Table:		"favourites",
		TableID:	"72",
		KeyVal:		map[string]int{"id":24},
	}

	reference_resolution.ResolveReference(displayConfig, &hint)

}

func main() { 

	migrationID := 434969759

	app := "mastodon"

	test1(app, migrationID)
}

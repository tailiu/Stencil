package evaluation

import (
	"fmt"
	"log"
	"stencil/db"
	"strings"
	"strconv"
)

// There could be three cases:
// For example, 
// 1. mTable: comments
// mCols: text,created_at,updated_at,comments.commentable_id,posts.id,statuses,
// commentable_id,comments.author_id,people.id,$visibility0,$reply2
// 2. mTable: photos
// mCols: photos.author_id,people.id,posts.id,posts.guid,photos.status_message_guid),
// posts.id,statuses,text,updated_at,created_at,remote_photo_path
// 3. mTable: posts
// mCols: created_at,updated_at,id,author_id,people,text
// Case 1: comments.commentable_id and posts.id
// Case 2: $visibility0
// Case 3: statuses
// Case 4: photos.status_message_guid)
// Case 5: people
// Case 6: There are two rows (photos and posts) to each Media attachment 
//	(calculate twice media sizes)
func procMigratedCols(evalConfig *EvalConfig, 
	mCols string, mTable string) []string {

	// log.Println(mCols)
	var procMCols []string

	tmp := strings.Split(mCols, ",")

	for _, col := range tmp { 

		// if mTable == "profiles" && 
		// 	(col == "username" || 
		// 	col == "serialized_public_key" ||
		// 	col == "serialized_private_key") {
		// 	continue
		// }

		// if mTable == "users" && 
		// 	(col == "serialized_public_key" ||
		// 	col == "diaspora_handle" ||
		// 	col == "image_url" ||
		// 	col == "bio") {
		// 	continue
		// }

		// if mTable == "people" && 
		// 	(col == "username" ||
		// 	col == "image_url" ||
		// 	col == "bio" ||
		// 	col == "serialized_private_key") {
		// 	continue
		// }

		col = strings.Replace(col, ")", "", -1)

		if strings.Contains(col, "$") {
			continue
		}

		if _, ok := evalConfig.MastodonTableNameIDPairs[col]; ok {
			continue
		}

		if _, ok := evalConfig.DiasporaTableNameIDPairs[col]; ok {
			continue
		}

		if strings.Contains(col, ".") {
			tmp1 := strings.Split(col, ".")
			if tmp1[0] != mTable {
				continue
			} else {
				col = tmp1[1]
			}
		}

		procMCols = append(procMCols, col)
	}

	return procMCols
	
}

func GetMigratedDataSizeV2(evalConfig *EvalConfig, migrationID string) int64 {

	query1 := fmt.Sprintf(
		`SELECT src_table, src_id, src_cols, dst_table, dst_id 
		FROM evaluation WHERE migration_id = '%s' and 
		dst_table != 'n/a'`, migrationID)

	// log.Println(query1)
	
	result, err := db.DataCall(evalConfig.StencilDBConn, query1)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(result)

	var tSize int64

	checkedMedia := make(map[string]bool)

	for _, data1 := range result {
		
		mTable := evalConfig.TableIDNamePairs[fmt.Sprint(data1["src_table"])]

		dTable := evalConfig.TableIDNamePairs[fmt.Sprint(data1["dst_table"])]

		log.Println(mTable)
		log.Println(dTable)

		mCols := procMigratedCols(evalConfig, 
			fmt.Sprint(data1["src_cols"]), mTable)

		log.Println(data1["src_cols"])
		log.Println(mCols)

		mID, err := strconv.Atoi(fmt.Sprint(data1["src_id"]))
		if err != nil {
			log.Fatal(err)
		}

		log.Println(mID)

		whetherCheckMediaSize := true

		if dTable == "media_attachments" {
			key := dTable + ":" + fmt.Sprint(data1["dst_id"])
			if _, ok :=	checkedMedia[key]; ok {
				whetherCheckMediaSize = false
			} else {
				checkedMedia[key] = true
			}
		}

		size := calculateRowSize(
			evalConfig.DiasporaDBConn, 
			mCols,
			mTable,
			mID,
			evalConfig.DiasporaAppID,
			whetherCheckMediaSize,
		)

		log.Println(size)

		tSize += size

	}

	return tSize

}
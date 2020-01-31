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

func GetMigratedDataSizeV2(evalConfig *EvalConfig, 
	migrationID string) int64 {

	query1 := fmt.Sprintf(
		`SELECT src_table, src_id, src_cols, 
		dst_table, dst_id, dst_cols  
		FROM evaluation WHERE migration_id = '%s'`,
		migrationID)

	// log.Println(query1)
	
	result, err := db.DataCall(evalConfig.StencilDBConn, query1)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(result)

	var tSize int64

	// checkedMedia := make(map[string]bool)

	for _, data1 := range result {
		
		mTable := evalConfig.TableIDNamePairs[fmt.Sprint(data1["src_table"])]

		dTable := evalConfig.TableIDNamePairs[fmt.Sprint(data1["dst_table"])]

		log.Println("src table:", mTable)
		log.Println("dst table:", dTable)

		mCols := procMigratedCols(evalConfig, 
			fmt.Sprint(data1["src_cols"]), mTable)

		log.Println("src cols:", data1["src_cols"])
		log.Println("processed src cols:", mCols)

		mID, err := strconv.Atoi(fmt.Sprint(data1["src_id"]))
		if err != nil {
			log.Fatal(err)
		}

		log.Println("dst cols:", fmt.Sprint(data1["dst_cols"]))
		
		// log.Println(mID)
		// log.Println(fmt.Sprint(data1["dst_id"]))

		// whetherCheckMediaSize := true

		// if dTable == "media_attachments" {
		// 	key := dTable + ":" + fmt.Sprint(data1["dst_id"])
		// 	if _, ok :=	checkedMedia[key]; ok {
		// 		whetherCheckMediaSize = false
		// 	} else {
		// 		checkedMedia[key] = true
		// 	}
		// }

		// log.Println(whetherCheckMediaSize)

		size := calculateRowSize(
			evalConfig.DiasporaDBConn, 
			mCols,
			mTable,
			mID,
			evalConfig.DiasporaAppID,
			true,
		)

		log.Println("size:", size)

		tSize += size

	}

	return tSize

}

// The dst cols don't contain id
func procDstCols(dstCols string) []string {

	tmp := strings.Split(dstCols, ",")

	return append(tmp, "id")

}

func GetMigratedDataSizeFromDst(evalConfig *EvalConfig, 
	migrationID string) int64 {

	query1 := fmt.Sprintf(
		`SELECT dst_table, dst_id, dst_cols  
		FROM evaluation WHERE migration_id = '%s'`,
		migrationID)

	// log.Println(query1)
	
	result, err := db.DataCall(evalConfig.StencilDBConn, query1)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(result)

	var tSize int64

	checkedData := make(map[string]bool)

	for _, data1 := range result {

		dstTable := evalConfig.TableIDNamePairs[fmt.Sprint(data1["dst_table"])]

		dstCols := fmt.Sprint(data1["dst_cols"])

		log.Println("dst table:", dstTable)
		log.Println("dst cols:", dstCols)
		
		strDstID := fmt.Sprint(data1["dst_id"])

		dstID, err := strconv.Atoi(strDstID)
		if err != nil {
			log.Fatal(err)
		}

		key := dstCols + ":" + strDstID
		if _, ok := checkedData[key]; ok {
			log.Println("duplicate")
			continue
		} else {
			checkedData[key] = true
		}

		procDstCols := procDstCols(dstCols)
		log.Println("processed dst cols:", procDstCols)

		log.Println("dst id:", dstID)

		size := calculateRowSize(
			evalConfig.MastodonDBConn, 
			procDstCols,
			dstTable,
			dstID,
			evalConfig.MastodonAppID,
			true,
		)

		log.Println("size:", size)

		tSize += size

	}

	return tSize

}

func procSrcCols(srcCols string) []string {

	tmp := strings.Split(srcCols, ",")

	return tmp

}

func GetMigratedDataSizeFromSrc(evalConfig *EvalConfig, 
	migrationID string) int64 {

	query1 := fmt.Sprintf(
		`SELECT src_table, src_id, src_cols  
		FROM evaluation WHERE migration_id = '%s'`,
		migrationID)

	// log.Println(query1)
	
	result, err := db.DataCall(evalConfig.StencilDBConn, query1)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(result)

	var tSize int64

	for _, data1 := range result {

		srcTable := evalConfig.TableIDNamePairs[fmt.Sprint(data1["src_table"])]

		srcCols := fmt.Sprint(data1["src_cols"])

		log.Println("src table:", srcTable)
		log.Println("src cols:", srcCols)
		
		strSrcID := fmt.Sprint(data1["src_id"])

		srcID, err := strconv.Atoi(strSrcID)
		if err != nil {
			log.Fatal(err)
		}

		procSrcCols := procSrcCols(srcCols)
		log.Println("processed src cols:", procSrcCols)

		log.Println("src id:", srcID)

		size := calculateRowSize(
			evalConfig.DiasporaDBConn, 
			procSrcCols,
			srcTable,
			srcID,
			evalConfig.DiasporaAppID,
			true,
		)

		log.Println("size:", size)

		tSize += size

	}

	return tSize

}
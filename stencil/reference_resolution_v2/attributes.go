package reference_resolution_v2

import (
	"database/sql"
	"fmt"
	"log"
	"stencil/db"
	"stencil/common_funcs"
)

/*
 * This new identity file aims to generalize reference resolution to different migrations
 * The functions here DO NOT consider migration id and GetPreviousID DOES NOT restrict
 * from app and from member
 */

func CreateAttribute(app, member, attrName, val string) *Attribute {

	attr := &Attribute{
		app:    	  app,
		member: 	  member,
		attrName:     attrName,
		val:		  val,
	}

	return attr
}

func createAttributeByRefRowToPart(procAttrRow map[string]string) *Attribute {

	attr := CreateAttribute(
		procAttrRow["to_app"],
		procAttrRow["to_member"],
		procAttrRow["to_attr"],
		procAttrRow["to_val"],
	)

	return attr
}

func getRowsFromAttrChangesTableByTo(refResolutionConfig *RefResolutionConfig,
	attr *Attribute) []map[string]interface{} {

	query := fmt.Sprintf(
		`SELECT * FROM attribute_changes WHERE 
		to_app = %s and to_member = %s and to_attr and to_val = '%s'`,
		attr.app, attr.member, attr.attrName, attr.val,
	)

	log.Println(query)

	data, err := db.DataCall(refResolutionConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return data
}

func getRowsFromAttrChangesTableByFrom(refResolutionConfig *RefResolutionConfig,
	attr *Attribute) []map[string]interface{} {

	query := fmt.Sprintf(
		`SELECT * FROM attribute_changes WHERE
		from_app = %s and from_member = %s and from_attr = %s and from_val = '%s'`,
		attr.app, attr.member, attr.attrName, attr.val,
	)

	data, err := db.DataCall(refResolutionConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return data

}

func forwardTraverseAttrChangesTable(refResolutionConfig *RefResolutionConfig,
	attr, orgAttr *Attribute, inRecurrsion bool) []*Attribute {

	var res []*Attribute

	attrRows := getRowsFromAttrChangesTableByFrom(refResolutionConfig, attr)
	// log.Println(IDRows)

	for _, attrRow := range attrRows {

		procAttrRow := common_funcs.TransformInterfaceToString(attrRow)

		nextData := createAttributeByRefRowToPart(procAttrRow)

		res = append(res, forwardTraverseAttrChangesTable(refResolutionConfig, nextData, orgAttr, true)...)

	}

	// If recurrsion has not started yet and we cannot find IDRows, then
	// we should directly return null result and 
	// not execute the code in the if block since this could directly return
	// the provided ID to us which is wrong!
	if len(attrRows) == 0 && inRecurrsion {

		// We don't need to test ID.id != orginalID.id becaseu as long as
		// ID.member != orginalID.member, this means that
		// this is different from the original row.
		// ID.id may be the same as orginalID.id in the scenario in which
		// migration does not change ids.
		// We don't find the cases in which ID.member == orginalID.member but
		// ID.id != orginalID.id, however this may happen..
		// Before changing:
		// if ID.app == refResolutionConfig.AppConfig.AppID &&
		// 	ID.member != orginalID.member && ID.id != orginalID.id {
		// if ID.app == refResolutionConfig.appID &&
		// 	(ID.member != orginalID.member || ID.id != orginalID.id) {
		// We remove the conditions: ID.member != orginalID.member || ID.id != orginalID.id
		// because there could be data referencing the same data
		// For example, when resolving tweets.guid in the mapping Twitter.tweets -> Diaspora.posts, 
		// tweets.guid is referring to the id of its own tweets.
		// In this case, we cannot use the condition to filter its own data. 
		// If there are more than one mappings to different tables, 
		// for example, when resolving Diaspora.notification_actors.notification_id, there are two ids:
		// Twitter.notifications.id -> Diaspora.notifications.id and Diaspora.notification_actors.id,
		// we use the third argument in the #REF to point out which table we are referring to 
		// (excluding Twitter.notifications.id -> Diaspora.notifications.id (same as the original data under check)
		// in that example). 
		if attr.app == refResolutionConfig.appID {
			res = append(res, attr)
		}
	}

	return res
}

func getUpdateToAttrInAttrChangesTableQuery(refResolutionConfig *RefResolutionConfig,
	memberName, attrToBeUpdated, attrValToBeUpdated, newAttrVal string) string {
		
	query := fmt.Sprintf(
		`UPDATE attribute_changes SET to_val = %s 
		WHERE to_app = %s and to_member = %s 
		and to_attr = %s and to_val = %s`,
		newAttrVal, 
		refResolutionConfig.appID, 
		refResolutionConfig.appTableNameIDPairs[memberName],
		refResolutionConfig.appAttrNameIDPairs[memberName + ":" + attrToBeUpdated],
		attrValToBeUpdated,
	)

	return query

}

func GetPreviousIDWithoutFromMember(refResolutionConfig *RefResolutionConfig,
	ID *Identity) string {

	query := fmt.Sprintf(`SELECT from_id FROM identity_table 
		WHERE to_app = %s and to_member = %s and to_id = %s`,
		ID.app, ID.member, ID.id)

	log.Println(query)

	data, err := db.DataCall1(refResolutionConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(data)

	if data["from_id"] == nil {
		return ""
	} else {
		return fmt.Sprint(data["from_id"])
	}

}

func GetPreviousID(refResolutionConfig *RefResolutionConfig,
	ID *Identity, fromMember string) string {

	query := fmt.Sprintf(
		`SELECT from_id FROM identity_table 
		WHERE from_member = %s and to_app = %s 
		and to_member = %s and to_id = %s`,
		fromMember, ID.app, ID.member, ID.id,
	)

	log.Println(query)

	data, err := db.DataCall1(refResolutionConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(data)

	if data["from_id"] == nil {
		return ""
	} else {
		return fmt.Sprint(data["from_id"])
	}

}

func GetPreIDByBackTraversal(refResolutionConfig *RefResolutionConfig,
	ID *Identity, fromMember string) string {

	query := fmt.Sprintf(`
		SELECT from_app, from_member, from_id FROM identity_table 
		WHERE to_app = %s and to_member = %s and to_id = %s`,
		ID.app, ID.member, ID.id,
	)

	log.Println(query)

	data, err := db.DataCall(refResolutionConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	for _, data1 := range data {
		
		from_app := fmt.Sprint(data1["from_app"])
		from_member := fmt.Sprint(data1["from_member"])
		from_id := fmt.Sprint(data1["from_id"])

		if data1["from_member"] == fromMember {
			return from_id
		}

		prevID := CreateIdentity(from_app, from_member, from_id)

		res := GetPreIDByBackTraversal(refResolutionConfig, prevID, fromMember)
		if res != "" {
			return res
		}

	}

	return ""

}

func getRootMembersOfApps(stencilDBConn *sql.DB) map[string]string {

	query := `SELECT app_id, root_member_id from app_root_member`

	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	rootMembers := make(map[string]string)
	
	for _, data1 := range data {
		rootMembers[fmt.Sprint(data1["app_id"])] =
			fmt.Sprint(data1["root_member_id"])
	}

	return rootMembers

}

func GetNextUserID(stencilDBConn *sql.DB, migrationID string) string {

	appRootMembers := getRootMembersOfApps(stencilDBConn)

	query := fmt.Sprintf(
		`SELECT user_id, src_app, dst_app FROM migration_registration
		WHERE migration_id = %s`,
		migrationID,
	)

	log.Println(query)

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	userID := fmt.Sprint(data["user_id"])
	srcApp := fmt.Sprint(data["src_app"])
	dstApp := fmt.Sprint(data["dst_app"])
	
	query1 := fmt.Sprintf(
		`SELECT to_id FROM identity_table 
		WHERE from_app = %s and from_member = %s and from_id = %s 
		and to_app = %s and to_member =%s and migration_id = %s`,
		srcApp, appRootMembers[srcApp], userID,
		dstApp, appRootMembers[dstApp], migrationID,
	)

	log.Println(query1)

	data1, err1 := db.DataCall1(stencilDBConn, query1)
	if err1 != nil {
		log.Fatal(err1)
	}

	if data1["to_id"] == nil {
		log.Fatal("Cannot get user ID in the destination application!")
	}

	return fmt.Sprint(data1["to_id"])

}

func getAppRootMemberID(stencilDBConn *sql.DB, appID string) string {

	query := fmt.Sprintf(`SELECT root_member_id from app_root_member 
		where app_id = %s`, appID)

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return fmt.Sprint(data["root_member_id"])

}

// func getPrevUserIDsByBackTraversal(stencilDBConn *sql.DB,
// 	appID, rootMemberID, userID string) [][]string {

// 	query := fmt.Sprintf(`SELECT * FROM identity_table
// 		WHERE to_app = %s and to_member = %s and to_id = %s`,
// 		appID, rootMemberID, userID)

// 	// log.Println(query)

// 	data, err := db.DataCall(stencilDBConn, query)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// log.Println(data)

// 	if len(data) == 0 {

// 		return nil

// 	} else {

// 		prevAppID := fmt.Sprint(data[0]["from_app"])

// 		prevRootMemberID := getAppRootMemberID(stencilDBConn, prevAppID)

// 		var prevUserID string

// 		for _, data1 := range data {
// 			if fmt.Sprint(data1["from_member"]) == prevRootMemberID {
// 				prevUserID = fmt.Sprint(data1["from_id"])
// 			}
// 		}

// 		preUserData := [][]string{
// 			[]string{
// 				prevAppID, prevUserID,
// 			}}

// 		prevPrevUserData := getPrevUserIDsByBackTraversal(stencilDBConn,
// 			prevAppID, prevRootMemberID, prevUserID)

// 		return append(preUserData, prevPrevUserData...)

// 	}

// }

func getPrevUserIDsByBackTraversal(stencilDBConn *sql.DB,
	appID, rootMemberID, userID string) map[string]string {

	query := fmt.Sprintf(`SELECT * FROM identity_table 
		WHERE to_app = %s and to_member = %s and to_id = %s`,
		appID, rootMemberID, userID)

	// log.Println(query)

	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	if len(data) == 0 {

		return nil

	} else {

		prevAppID := fmt.Sprint(data[0]["from_app"])

		prevRootMemberID := getAppRootMemberID(stencilDBConn, prevAppID)

		var prevUserID string

		for _, data1 := range data {
			if fmt.Sprint(data1["from_member"]) == prevRootMemberID {
				prevUserID = fmt.Sprint(data1["from_id"])
			}
		}

		preUserData := map[string]string{prevAppID: prevUserID}

		prevPrevUserData := getPrevUserIDsByBackTraversal(stencilDBConn,
			prevAppID, prevRootMemberID, prevUserID)

		for k, v := range prevPrevUserData {
			preUserData[k] = v
		}

		return preUserData

	}

}

// This function gets all previous user ids in previous applications
// Note that for now it can only get previous user ids
// when user has a line migration history, for example,
// A -> B -> C and we have the appID and userID in C, it can get A and B
// For the case A -> B and A -> C concurrently (not a line history)
// and we have the appID and userID in C, it can only get A not B
func GetPrevUserIDs(stencilDBConn *sql.DB, appID, userID string) map[string]string {

	rootMemberID := getAppRootMemberID(stencilDBConn, appID)

	return getPrevUserIDsByBackTraversal(stencilDBConn,
		appID, rootMemberID, userID)

}

func CreateAttributeChangesTable(dbConn *sql.DB) {
	
	var queries []string

	query1 := `CREATE TABLE attribute_changes (
			from_app 		INT8	NOT NULL,
			from_member 	INT8	NOT NULL,
			from_attr 		INT8	NOT NULL,
			from_val 		VARCHAR NOT NULL,
			to_app 			INT8	NOT NULL,
			to_member 		INT8	NOT NULL,
			to_attr 		INT8	NOT NULL,
			to_val 			VARCHAR NOT NULL,
			migration_id	INT8	NOT NULL,
			pk				SERIAL PRIMARY KEY,
			FOREIGN KEY (from_member) REFERENCES app_tables (pk),
			FOREIGN KEY (to_member) REFERENCES app_tables (pk),
			FOREIGN KEY (from_app) REFERENCES apps (pk),
			FOREIGN KEY (to_app) REFERENCES apps (pk),
			FOREIGN KEY (from_attr) REFERENCES app_schemas (pk),
			FOREIGN KEY (to_attr) REFERENCES app_schemas (pk))`

	query2 := "CREATE INDEX ON attribute_changes(to_app, to_member, to_attr, to_val)"

	query3 := "CREATE INDEX ON attribute_changes(from_app, from_member, from_attr, from_val)"

	queries = append(queries, query1, query2, query3)

	for _, query := range queries {
		err := db.TxnExecute1(dbConn, query)
		if err != nil {
			log.Fatal(err)
		}
	}

}
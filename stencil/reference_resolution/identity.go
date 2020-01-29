package reference_resolution

import (
	"database/sql"
	"fmt"
	"log"
	"stencil/db"
)

/*
 * This new identity file aims to generalize reference resolution to different migrations
 * The functions here DO NOT consider migration id and GetPreviousID DOES NOT restrict
 * from app and from member
 */

// app, member, id are all integer corresponding to names
func CreateIdentity(app, member, id string) *Identity {

	ID := &Identity{
		app:    app,
		member: member,
		id:     id,
	}

	return ID
}

func getRowsFromIDTableByTo(refResolutionConfig *RefResolutionConfig,
	ID *Identity) []map[string]interface{} {

	query := fmt.Sprintf(`SELECT * FROM identity_table 
		WHERE to_app = %s and to_member = %s and to_id = %s`,
		ID.app, ID.member, ID.id)

	log.Println(query)

	data, err := db.DataCall(refResolutionConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return data

}

func getRowsFromIDTableByFrom(refResolutionConfig *RefResolutionConfig,
	ID *Identity) []map[string]interface{} {

	query := fmt.Sprintf(`SELECT * FROM identity_table 
		WHERE from_app = %s and from_member = %s and from_id = %s`,
		ID.app, ID.member, ID.id)

	data, err := db.DataCall(refResolutionConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return data

}

func forwardTraverseIDTable(refResolutionConfig *RefResolutionConfig,
	ID, orginalID *Identity) []*Identity {

	var res []*Identity

	IDRows := getRowsFromIDTableByFrom(refResolutionConfig, ID)
	// log.Println(IDRows)

	for _, IDRow := range IDRows {

		procIDRow := transformInterfaceToString(IDRow)

		nextData := CreateIdentity(
			procIDRow["to_app"],
			procIDRow["to_member"],
			procIDRow["to_id"])

		res = append(res, forwardTraverseIDTable(refResolutionConfig, nextData, orginalID)...)

	}

	if len(IDRows) == 0 {

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
		if ID.app == refResolutionConfig.appID &&
			(ID.member != orginalID.member || ID.id != orginalID.id) {

			resData := CreateIdentity(ID.app, ID.member, ID.id)

			res = append(res, resData)

		}

	}

	return res
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

	query := fmt.Sprintf(`SELECT from_id FROM identity_table 
		WHERE from_member = %s and to_app = %s 
		and to_member = %s and to_id = %s`,
		fromMember, ID.app, ID.member, ID.id)

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

func getPrevUserIDsByBackTraversal(stencilDBConn *sql.DB,
	appID, rootMemberID, userID string) [][]string {

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

		preUserData := [][]string{
			[]string{
				prevAppID, prevUserID,
			}}

		prevPrevUserData := getPrevUserIDsByBackTraversal(stencilDBConn,
			prevAppID, prevRootMemberID, prevUserID)

		return append(preUserData, prevPrevUserData...)

	}

}

// This function gets all previous user ids in previous applications
// Note that for now it can only get previous user ids
// when user has a line migration history, for example,
// A -> B -> C and we have the appID and userID in C, it can get A and B
// For the case A -> B and A -> C concurrently (not a line history)
// and we have the appID and userID in C, it can only get A not B
func GetPrevUserIDs(appID, userID string) [][]string {

	stencilDBConn := db.GetDBConn(db.STENCIL_DB)

	rootMemberID := getAppRootMemberID(stencilDBConn, appID)

	return getPrevUserIDsByBackTraversal(stencilDBConn,
		appID, rootMemberID, userID)

}

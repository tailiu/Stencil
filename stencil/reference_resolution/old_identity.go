package reference_resolution

import (
	"stencil/db"
	"stencil/config"
	"database/sql"
	"fmt"
	"log"
)

// app, member, id are all integer corresponding to names
func oldCreateIdentity(app, member, id string) *Identity {
	
	ID := &Identity{
		app: 	app,
		member:	member,
		id:		id,
	}

	return ID
}

func oldGetRowsFromIDTableByTo(refResolutionConfig *RefResolutionConfig, 
	ID *Identity) []map[string]interface{} {

	query := fmt.Sprintf(`SELECT * FROM identity_table 
		WHERE to_app = %s and to_member = %s and to_id = %s and migration_id = %d`,
		ID.app, ID.member, ID.id, refResolutionConfig.migrationID)
	
	log.Println(query)

	data, err := db.DataCall(refResolutionConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return data

}


func oldGetRowsFromIDTableByFrom(refResolutionConfig *RefResolutionConfig, 
	ID *Identity) []map[string]interface{} {
	
	query := fmt.Sprintf(`SELECT * FROM identity_table 
		WHERE from_app = %s and from_member = %s and from_id = %s and migration_id = %d`,
		ID.app, ID.member, ID.id, refResolutionConfig.migrationID)

	data, err := db.DataCall(refResolutionConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return data

}

func oldGetPreviousID(refResolutionConfig *RefResolutionConfig, 
	ID *Identity, from_app, fromMember string) string {

	query := fmt.Sprintf(`SELECT from_id FROM identity_table 
		WHERE from_app = %s and from_member = %s and to_app = %s and 
		to_member = %s and to_id = %s and migration_id = %d`,
		from_app, fromMember, ID.app, ID.member, ID.id, 
		refResolutionConfig.migrationID)
	
	// log.Println(query)

	data, err := db.DataCall1(refResolutionConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	if data["from_id"] == nil {
		return ""
	} else {
		return fmt.Sprint(data["from_id"])
	}
	
}

func oldForwardTraverseIDTable(refResolutionConfig *RefResolutionConfig, 
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

		res = append(res, oldForwardTraverseIDTable(refResolutionConfig, nextData, orginalID)...)

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

func oldGetAllAppsConnections() {

	appsConnections := make(map[string]*sql.DB)

	stencilDBConn := db.GetDBConn(config.StencilDBName)

	query1 := fmt.Sprintf(`SELECT pk, app_name from apps`)
	
	data, err := db.DataCall(stencilDBConn, query1)
	if err != nil {
		log.Fatal(err)
	}

	for _, data1 := range data {

		appName := fmt.Sprint(data1["app_name"])

		appID := fmt.Sprint(data1["pk"])

		appsConnections[appID] = db.GetDBConn(appName)
	}

}
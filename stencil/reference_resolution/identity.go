package reference_resolution

import (
	"stencil/db"
	"fmt"
	"log"
)

// app, member, id are all integer corresponding to names
func CreateIdentity(app, member, id string) *Identity {
	
	ID := &Identity{
		app: 	app,
		member:	member,
		id:		id,
	}

	return ID
}

func getRowsFromIDTableByTo(refResolutionConfig *RefResolutionConfig, 
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


func getRowsFromIDTableByFrom(refResolutionConfig *RefResolutionConfig, 
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
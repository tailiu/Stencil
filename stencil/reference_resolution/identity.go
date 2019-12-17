package reference_resolution

import (
	"stencil/db"
	"stencil/config"
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

func getRowsFromIDTableByTo(displayConfig *config.DisplayConfig, 
	ID *Identity) []map[string]interface{} {

	query := fmt.Sprintf(`SELECT * FROM identity_table 
		WHERE to_app = %s and to_member = %s and to_id = %s and migration_id = %d`,
		ID.app, ID.member, ID.id, displayConfig.MigrationID)
	
	log.Println(query)

	data, err := db.DataCall(displayConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return data

}


func getRowsFromIDTableByFrom(displayConfig *config.DisplayConfig, 
	ID *Identity) []map[string]interface{} {
	
	query := fmt.Sprintf(`SELECT * FROM identity_table 
		WHERE from_app = %s and from_member = %s and from_id = %s and migration_id = %d`,
		ID.app, ID.member, ID.id, displayConfig.MigrationID)

	data, err := db.DataCall(displayConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return data

}

func forwardTraverseIDTable(displayConfig *config.DisplayConfig, 
	ID, orginalID *Identity) []*Identity {
	
	var res []*Identity

	IDRows := getRowsFromIDTableByFrom(displayConfig, ID)
	// log.Println(IDRows)

	for _, IDRow := range IDRows {
		
		procIDRow := transformInterfaceToString(IDRow)

		nextData := CreateIdentity(
			procIDRow["to_app"], 
			procIDRow["to_member"], 
			procIDRow["to_id"])

		res = append(res, forwardTraverseIDTable(displayConfig, nextData, orginalID)...)

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
		// if ID.app == displayConfig.AppConfig.AppID && 
		// 	ID.member != orginalID.member && ID.id != orginalID.id {
		if ID.app == displayConfig.AppConfig.AppID && 
			(ID.member != orginalID.member || ID.id != orginalID.id) {
			
			resData := CreateIdentity(ID.app, ID.member, ID.id)

			res = append(res, resData)

		}

	}

	return res
}
package reference_resolution

import (
	"stencil/db"
	"stencil/config"
	"stencil/app_display"
	"strconv"
	"fmt"
	"log"
)

func createIdentity(app, member, id string) *identity {
	
	ID := &identity{
		app: 	app,
		member:	member,
		id:		id,
	}

	return ID
}

func transformHintToIdenity(displayConfig *config.DisplayConfig, hint *app_display.HintStruct) *identity {

	return createIdentity(displayConfig.AppConfig.AppID, hint.TableID, strconv.Itoa(hint.KeyVal["id"]))

}

func getRowsFromIDTableByTo(displayConfig *config.DisplayConfig, ID *identity) []map[string]interface{} {

	query := fmt.Sprintf("SELECT * FROM identity_table WHERE to_app = %s and to_member = %s and to_id = %s and migration_id = %d",
		ID.app, ID.member, ID.id, displayConfig.MigrationID)
	
	data, err := db.DataCall(displayConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return data

}


func getRowsFromIDTableByFrom(displayConfig *config.DisplayConfig, ID *identity) []map[string]interface{} {
	
	query := fmt.Sprintf("SELECT * FROM identity_table WHERE from_app = %s and from_member = %s and from_id = %s and migration_id = %d",
		ID.app, ID.member, ID.id, displayConfig.MigrationID)

	data, err := db.DataCall(displayConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return data

}

func forwardTraverseIDTable(displayConfig *config.DisplayConfig, ID, orginalID *identity, dstAppID string) []*identity {
	
	var res []*identity

	IDRows := getRowsFromIDTableByFrom(displayConfig, ID)
	// log.Println(IDRows)

	for _, IDRow := range IDRows {
		
		procIDRow := transformInterfaceToString(IDRow)
		nextData := &identity {
			app: 	procIDRow["to_app"],
			member:	procIDRow["to_member"],
			id:		procIDRow["to_id"],
		}
		res = append(res, forwardTraverseIDTable(displayConfig, nextData, orginalID, dstAppID)...)

	}

	if len(IDRows) == 0 {

		if ID.app == dstAppID && ID.member != orginalID.member && ID.id != orginalID.id {
			
			resData := &identity {
				app: 	ID.app,
				member:	ID.member,
				id:		ID.id,
			}
			res = append(res, resData)
		}

	}

	return res

}
package reference_resolution

import (
	"stencil/db"
	"database/sql"
	"stencil/config"
	"stencil/app_display"
	"fmt"
	"log"
)

func getRowsFromIDTableByTo(stencilDBConn *sql.DB, appConfig *config.AppConfig, migrationID int, hint *app_display.HintStruct) []map[string]interface{} {

	query := fmt.Sprintf("SELECT * FROM identity_table WHERE to_app = %s and to_member = %s and to_id = %d and migration_id = %d",
		appConfig.AppID, hint.TableID, hint.KeyVal["id"], migrationID)
	
	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return data

}


func getRowsFromIDTableByFrom(stencilDBConn *sql.DB, migrationID int, ID *identity) []map[string]interface{} {
	
	query := fmt.Sprintf("SELECT * FROM identity_table WHERE from_app = %s and from_member = %s and from_id = %s and migration_id = %d",
		ID.app, ID.member, ID.id, migrationID)

	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return data

}

func forwardTraverseIDTable(stencilDBConn *sql.DB, migrationID int, ID, orginalID *identity, dstAppID string) []*identity {
	
	var res []*identity

	IDRows := getRowsFromIDTableByFrom(stencilDBConn, migrationID, ID)
	// log.Println(IDRows)

	for _, IDRow := range IDRows {
		
		procIDRow := transformInterfaceToString(IDRow)
		nextData := &identity {
			app: 	procIDRow["to_app"],
			member:	procIDRow["to_member"],
			id:		procIDRow["to_id"],
		}
		res = append(res, forwardTraverseIDTable(stencilDBConn, migrationID, nextData, orginalID, dstAppID)...)

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
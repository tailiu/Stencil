package reference_resolution

import (
	"stencil/db"
	"database/sql"
	"fmt"
)

func GetRowsFromIDTableByTo(stencilDBConn *sql.DB, appConfig *config.AppConfig, migrationID int) []map[string]interface{} {

	query := fmt.Sprintf("SELECT * FROM identity_table WHERE to_app = %d and to_member = %d and to_id = %d and migration_id = %d",
		appConfig.AppID, member, id, migrationID)
	
	data, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)

	return data

}


func GetRowsFromIDTableByFrom(fromApp, member, id) {
	
}

func ForwardTraverseIDTable() {

	IDRows := GetRowsFromIDTableByFrom(app, member, id)

	for _, IDRow := range IDRows {
		ForwardTraverseIDTable(IDRow.ToApp, IDRow.ToMember, IDRow.ToID, org_member, org_id, listToBeReturned)
	}

	if len(IDRows) == 0 {
		if app == t.DstApp && member != org_member && id != org_id {
			listToBeReturned = append(listToBeReturned, []string{tag, member, id})
		}
	}

}
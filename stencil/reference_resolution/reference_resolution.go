package reference_resolution

import (
	"stencil/app_display"
	"database/sql"
	"stencil/config"
)

// You are on the left/from part
func updateMyDataBasedOnReferences(stencilDBConn *sql.DB, appConfig *config.AppConfig, migrationID int, IDRow map[string]string) {
	
	for _, ref := range getFromReferences(stencilDBConn, migrationID, IDRow) {
		
		proRef := transformInterfaceToString(ref)	
		data := &identity{
			app: 	proRef["app"],
			member:	proRef["from_member"],
			id:		proRef["from_id"],
		}
		forwardTraverseIDTable(stencilDBConn, migrationID, data, data, appConfig.AppID)

	}

}

// You are on the right/to part
func updateOtherDataBasedOnReferences() {

}

func ResolveReferenceByBackTraversal(stencilDBConn *sql.DB, appConfig *config.AppConfig, migrationID int, hint *app_display.HintStruct) {
	
	for _, IDRow := range getRowsFromIDTableByTo(stencilDBConn, appConfig, migrationID, hint) {
		
		proIDRow := transformInterfaceToString(IDRow)

		// You are on the left/from part
		updateMyDataBasedOnReferences(stencilDBConn, appConfig, migrationID, proIDRow)

		// You are on the right/to part
		// updateOtherDataBasedOnReferences()

		// Traverse back
		// ResolveReferenceByBackTraversal()
	}

}
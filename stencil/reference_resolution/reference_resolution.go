package reference_resolution

import (
	"stencil/app_display"
	"database/sql"
	"stencil/config"
)

// You are on the left/from part
func updateMyDataBasedOnReferences() {

}

// You are on the right/to part
func updateOtherDataBasedOnReferences() {

}

func ResolveReferenceByBackTraversal(stencilDBConn *sql.DB, appConfig *config.AppConfig, migrationID int, hint *app_display.HintStruct) {
	
	for _, IDRow := range GetRowsFromIDTableByTo(stencilDBConn, appConfig, migrationID, member, id) {
		// You are on the left/from part
		updateMyDataBasedOnReferences()

		// You are on the right/to part
		updateOtherDataBasedOnReferences()

		// Traverse back
		ResolveReferenceByBackTraversal()
	}

}
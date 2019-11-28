package reference_resolution

import (
	"stencil/app_display"
	"stencil/config"
	"stencil/schema_mapping"
	"database/sql"
	"log"
)

// You are on the left/from part
func updateMyDataBasedOnReferences(stencilDBConn *sql.DB, appConfig *config.AppConfig, migrationID int, IDRow map[string]string) {
	
	for _, ref := range getFromReferences(stencilDBConn, migrationID, IDRow) {
		
		proRef := transformInterfaceToString(ref)
		// log.Println(proRef)

		data := &identity{
			app: 	proRef["app"],
			member:	proRef["to_member"],
			id:		proRef["to_id"],
		}
		refIdentityRows := forwardTraverseIDTable(stencilDBConn, migrationID, data, data, appConfig.AppID)
		// log.Println(refIdentityRows[0])

		if len(refIdentityRows) > 0 {
			
			for _, refIdentityRow := range refIdentityRows {
				schema_mapping.GetMappedAttributeFromSchemaMappings(
					proRef["App"], proRef["to_member"], proRef["to_reference"], appConfig.AppID, refIdentityRow.member)
			}

		} else if proRef["app"] == appConfig.AppID {
			
		}

	}

}

// You are on the right/to part
func updateOtherDataBasedOnReferences() {

}

func ResolveReferenceByBackTraversal(stencilDBConn *sql.DB, appConfig *config.AppConfig, migrationID int, hint *app_display.HintStruct) {
	
	for _, IDRow := range getRowsFromIDTableByTo(stencilDBConn, appConfig, migrationID, hint) {
		
		proIDRow := transformInterfaceToString(IDRow)
		log.Println(proIDRow)

		// You are on the left/from part
		updateMyDataBasedOnReferences(stencilDBConn, appConfig, migrationID, proIDRow)

		// You are on the right/to part
		// updateOtherDataBasedOnReferences()

		// Traverse back
		// ResolveReferenceByBackTraversal()
	}

}
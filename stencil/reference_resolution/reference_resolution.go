package reference_resolution

import (
	"stencil/app_display"
	"stencil/config"
	"stencil/schema_mapping"
	"database/sql"
	"log"
)

// You are on the left/from part
func updateMyDataBasedOnReferences(displayConfig *config.DisplayConfig, IDRow map[string]string) {
	
	for _, ref := range getFromReferences(displayConfig, IDRow) {
		
		proRef := transformInterfaceToString(ref)
		// log.Println(proRef)

		data := &identity{
			app: 	proRef["app"],
			member:	proRef["to_member"],
			id:		proRef["to_id"],
		}
		refIdentityRows := forwardTraverseIDTable(displayConfig, data, data, displayConfig.AppConfig.AppID)
		// log.Println(refIdentityRows[0])

		if len(refIdentityRows) > 0 {
			
			for _, refIdentityRow := range refIdentityRows {

				attr := schema_mapping.GetMappedAttributeFromSchemaMappings(displayConfig, 
					proRef["app"], proRef["to_member"], proRef["to_reference"], 
					appConfig.AppName, refIdentityRow.member)
				
				log.Println(attr)
				// for _, attrToUpdate := range schema_mapping.GetMappedAttributeFromSchemaMappings(
				// 	proRef["app"], proRef["from_member"], proRef["from_reference"], appConfig.AppName, org_member) {
					
					// updateReferences(ref, refIdentityRow.ToMember, refIdentityRow.ToID, attr, org_member, org_id, AttrToUpdate)
					
				// }
			}

		} else if proRef["app"] == appConfig.AppID {
			
		}

	}

}

// You are on the right/to part
func updateOtherDataBasedOnReferences() {

}

func ResolveReferenceByBackTraversal(displayConfig *config.DisplayConfig, hint *app_display.HintStruct) {
	
	for _, IDRow := range getRowsFromIDTableByTo(displayConfig, hint) {
		
		proIDRow := transformInterfaceToString(IDRow)
		log.Println(proIDRow)

		// You are on the left/from part
		updateMyDataBasedOnReferences(displayConfig, proIDRow)

		// You are on the right/to part
		// updateOtherDataBasedOnReferences()

		// Traverse back
		// ResolveReferenceByBackTraversal()
	}

}
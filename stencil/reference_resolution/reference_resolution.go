package reference_resolution

import (
	"stencil/app_display"
	"stencil/config"
	"stencil/schema_mappings"
	"log"
)

// You are on the left/from part
func updateMyDataBasedOnReferences(displayConfig *config.DisplayConfig, IDRow map[string]string) {
	
	for _, ref := range getFromReferences(displayConfig, IDRow) {
		
		proRef := transformInterfaceToString(ref)
		// log.Println(proRef)

		data := createIdentity(proRef["app"], proRef["to_member"], proRef["to_id"])

		refIdentityRows := forwardTraverseIDTable(displayConfig, data, data, displayConfig.AppConfig.AppID)
		log.Println(refIdentityRows[0])

		if len(refIdentityRows) > 0 {
			
			for _, refIdentityRow := range refIdentityRows {

				attr, err := schema_mappings.GetMappedAttributeFromSchemaMappings( 
					displayConfig.AppIDNamePairs[proRef["app"]], 
					displayConfig.TableIDNamePairs[proRef["to_member"]],
					displayConfig.TableIDNamePairs[proRef["to_member"]] + "." + proRef["to_reference"], 
					displayConfig.AppConfig.AppName,
					displayConfig.TableIDNamePairs[refIdentityRow.member])

				if err != nil {
					log.Fatal(err)
				}

				log.Println(attr)
				// for _, attrToUpdate := range schema_mapping.GetMappedAttributeFromSchemaMappings(
				// 	proRef["app"], proRef["from_member"], proRef["from_reference"], appConfig.AppName, org_member) {
					
					// updateReferences(ref, refIdentityRow.ToMember, refIdentityRow.ToID, attr, org_member, org_id, AttrToUpdate)
					
				// }
			}

		} else if proRef["app"] == displayConfig.AppConfig.AppID {
			
		}

	}

}

// You are on the right/to part
func updateOtherDataBasedOnReferences() {

}

func resolveReferenceByBackTraversal(displayConfig *config.DisplayConfig, ID *identity) {
	
	for _, IDRow := range getRowsFromIDTableByTo(displayConfig, ID) {

		proIDRow := transformInterfaceToString(IDRow)
		log.Println(proIDRow)

		// You are on the left/from part
		updateMyDataBasedOnReferences(displayConfig, proIDRow)

		// You are on the right/to part
		// updateOtherDataBasedOnReferences()

		// Traverse back
		// resolveReferenceByBackTraversal()
	}

}

func ResolveReference(displayConfig *config.DisplayConfig, hint *app_display.HintStruct) {
	
	ID := transformHintToIdenity(displayConfig, hint)
	
	resolveReferenceByBackTraversal(displayConfig, ID)

}
package reference_resolution

import (
	"stencil/app_display"
	"stencil/config"
	"stencil/schema_mappings"
	"log"
)

// You are on the left/from part
func updateMyDataBasedOnReferences(displayConfig *config.DisplayConfig, 
	IDRow map[string]string, orgID *identity) {
	
	for _, ref := range getFromReferences(displayConfig, IDRow) {
		
		proRef := transformInterfaceToString(ref)
		// log.Println(proRef)

		data := createIdentity(proRef["app"], proRef["to_member"], proRef["to_id"])

		refIdentityRows := forwardTraverseIDTable(
			displayConfig, data, data, displayConfig.AppConfig.AppID)
		log.Println(refIdentityRows[0])

		if len(refIdentityRows) > 0 {
			
			for _, refIdentityRow := range refIdentityRows {

				// This is a little bit dirty hack. For example, when trying to find
				// the mapped attribute from Diaspora Posts Posts.id to Mastodon Statuses,
				// if we consider arguments in #REF, 
				// there are two results: id and conversation_id which should not be included
				ignoreREF := true

				attrs, err := schema_mappings.GetMappedAttributesFromSchemaMappings( 
					displayConfig.AppIDNamePairs[proRef["app"]], 
					displayConfig.TableIDNamePairs[proRef["to_member"]],
					displayConfig.TableIDNamePairs[proRef["to_member"]] + 
						"." + proRef["to_reference"], 
					displayConfig.AppConfig.AppName,
					displayConfig.TableIDNamePairs[refIdentityRow.member],
					ignoreREF)

				if err != nil {
					log.Fatal(err)
				}

				log.Println(attrs)

				ignoreREF = false 

				attrsToUpdate, err1 := schema_mappings.GetMappedAttributesFromSchemaMappings(
					displayConfig.AppIDNamePairs[proRef["app"]], 
					displayConfig.TableIDNamePairs[proRef["from_member"]], 
					displayConfig.TableIDNamePairs[proRef["from_member"]] +
						"." + proRef["from_reference"], 
					displayConfig.AppConfig.AppName,
					displayConfig.TableIDNamePairs[orgID.member],
					ignoreREF)

				if err1 != nil {
					log.Fatal(err1)
				}

				log.Println(attrsToUpdate)

				for _, attrToUpdate := range attrsToUpdate {
					
					log.Println(attrToUpdate)

					err2 := updateReferences(displayConfig,
						proRef["pk"], 
						displayConfig.TableIDNamePairs[refIdentityRow.member], 
						refIdentityRow.id, 
						attrs[0], 
						displayConfig.TableIDNamePairs[orgID.member], 
						orgID.id, 
						attrToUpdate)
					
					if err2 != nil {
						log.Println(err2)
					}
					
				}
			}

		} else if proRef["app"] == displayConfig.AppConfig.AppID {

			attr := proRef["to_reference"]

			log.Println(attr)

			ignoreREF := false 

			attrsToUpdate, err := schema_mappings.GetMappedAttributesFromSchemaMappings(
				displayConfig.AppIDNamePairs[proRef["app"]],
				displayConfig.TableIDNamePairs[proRef["from_member"]],
				displayConfig.TableIDNamePairs[proRef["from_member"]] +
					"." + proRef["from_reference"], 
				displayConfig.AppConfig.AppName, 
				displayConfig.TableIDNamePairs[orgID.member],
				ignoreREF)
			
			if err != nil {
				log.Fatal(err)
			}

			log.Println(attrsToUpdate)

			for _, attrToUpdate := range attrsToUpdate {

				err1 := updateReferences(displayConfig,
					proRef["pk"],  
					displayConfig.TableIDNamePairs[proRef["to_member"]], 
					displayConfig.TableIDNamePairs[proRef["to_id"]], 
					attr, 
					displayConfig.TableIDNamePairs[orgID.member], 
					orgID.id, 
					attrToUpdate)
				
				if err1 != nil {
					log.Println(err1)
				}

			}
		}

	}

}

// You are on the right/to part
func updateOtherDataBasedOnReferences(displayConfig *config.DisplayConfig, 
	IDRow map[string]string, orgID *identity) {
	
	for _, ref := range getToReferences(displayConfig, IDRow) {
		proRef := transformInterfaceToString(ref)
		// log.Println(proRef)

		data := createIdentity(proRef["app"], proRef["from_member"], proRef["from_id"])

		refIdentityRows := forwardTraverseIDTable(
			displayConfig, data, data, displayConfig.AppConfig.AppID)
		log.Println(refIdentityRows[0])

	}

}

func resolveReferenceByBackTraversal(displayConfig *config.DisplayConfig, 
	ID *identity, orgID *identity) {
	
	for _, IDRow := range getRowsFromIDTableByTo(displayConfig, ID) {

		proIDRow := transformInterfaceToString(IDRow)
		log.Println(proIDRow)

		// You are on the left/from part
		updateMyDataBasedOnReferences(displayConfig, proIDRow, orgID)

		// You are on the right/to part
		updateOtherDataBasedOnReferences(displayConfig, proIDRow, orgID)

		// Traverse back
		preID := createIdentity(IDRow["from_app"], IDRow["from_member"], IDRow["from_id"])

		resolveReferenceByBackTraversal(displayConfig, preID, orgID)
	}

}

func ResolveReference(displayConfig *config.DisplayConfig, hint *app_display.HintStruct) {
	
	ID := transformHintToIdenity(displayConfig, hint)
	
	resolveReferenceByBackTraversal(displayConfig, ID, ID)

}
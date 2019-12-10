package reference_resolution

import (
	"stencil/app_display"
	"stencil/config"
	"stencil/schema_mappings"
	"log"
)

// You are on the left/from part
func updateMyDataBasedOnReferences(displayConfig *config.DisplayConfig, 
	IDRow map[string]string, orgID *identity) []string {
	
	var updatedAttrs []string

	for _, ref := range getFromReferences(displayConfig, IDRow) {
		
		procRef := transformInterfaceToString(ref)
		log.Println("ref_row: ", procRef)

		data := createIdentity(procRef["app"], procRef["to_member"], procRef["to_id"])

		refIdentityRows := forwardTraverseIDTable(displayConfig, data, orgID)
		// log.Println(refIdentityRows[0])

		if len(refIdentityRows) > 0 {
			
			for _, refIdentityRow := range refIdentityRows {

				// This is a little bit dirty hack. For example, when trying to find
				// the mapped attribute from Diaspora Posts Posts.id to Mastodon Statuses,
				// if we consider arguments in #REF, 
				// there are two results: id and conversation_id which should not be included
				ignoreREF := true

				attrs, err := schema_mappings.GetMappedAttributesFromSchemaMappings( 
					displayConfig.AppIDNamePairs[procRef["app"]], 
					displayConfig.TableIDNamePairs[procRef["to_member"]],
					displayConfig.TableIDNamePairs[procRef["to_member"]] + 
						"." + procRef["to_reference"], 
					displayConfig.AppConfig.AppName,
					displayConfig.TableIDNamePairs[refIdentityRow.member],
					ignoreREF)

				if err != nil {
					log.Println(err)
				}

				log.Println(attrs)

				ignoreREF = false 

				attrsToUpdate, err1 := schema_mappings.GetMappedAttributesFromSchemaMappings(
					displayConfig.AppIDNamePairs[procRef["app"]], 
					displayConfig.TableIDNamePairs[procRef["from_member"]], 
					displayConfig.TableIDNamePairs[procRef["from_member"]] +
						"." + procRef["from_reference"], 
					displayConfig.AppConfig.AppName,
					displayConfig.TableIDNamePairs[orgID.member],
					ignoreREF)

				if err1 != nil {
					log.Println(err1)
				}

				log.Println(attrsToUpdate)

				for _, attrToUpdate := range attrsToUpdate {
					
					log.Println(attrToUpdate)

					err2 := updateReferences(
						displayConfig,
						procRef["pk"], 
						displayConfig.TableIDNamePairs[refIdentityRow.member], 
						refIdentityRow.id, 
						attrs[0], 
						displayConfig.TableIDNamePairs[orgID.member], 
						orgID.id, 
						attrToUpdate)
					
					updatedAttrs = append(updatedAttrs, attrToUpdate)

					if err2 != nil {
						log.Println(err2)
					}
					
				}
			}

		} else if procRef["app"] == displayConfig.AppConfig.AppID {

			attr := procRef["to_reference"]

			log.Println(attr)

			ignoreREF := false 

			attrsToUpdate, err := schema_mappings.GetMappedAttributesFromSchemaMappings(
				displayConfig.AppIDNamePairs[procRef["app"]],
				displayConfig.TableIDNamePairs[procRef["from_member"]],
				displayConfig.TableIDNamePairs[procRef["from_member"]] +
					"." + procRef["from_reference"], 
				displayConfig.AppConfig.AppName, 
				displayConfig.TableIDNamePairs[orgID.member],
				ignoreREF)
			
			if err != nil {
				log.Println(err)
			}

			log.Println(attrsToUpdate)

			for _, attrToUpdate := range attrsToUpdate {

				err1 := updateReferences(
					displayConfig,
					procRef["pk"],  
					displayConfig.TableIDNamePairs[procRef["to_member"]], 
					procRef["to_id"], 
					attr, 
					displayConfig.TableIDNamePairs[orgID.member], 
					orgID.id, 
					attrToUpdate)
				
				updatedAttrs = append(updatedAttrs, attrToUpdate)

				if err1 != nil {
					log.Println(err1)
				}

			}
		}

	}

	return updatedAttrs

}

// You are on the right/to part
func updateOtherDataBasedOnReferences(displayConfig *config.DisplayConfig, 
	IDRow map[string]string, orgID *identity) []string {
	
	var updatedAttrs []string

	for _, ref := range getToReferences(displayConfig, IDRow) {

		procRef := transformInterfaceToString(ref)
		// log.Println(procRef)

		data := createIdentity(procRef["app"], procRef["from_member"], procRef["from_id"])

		refIdentityRows := forwardTraverseIDTable(displayConfig, data, orgID)
		// log.Println(refIdentityRows[0])

		if len(refIdentityRows) > 0 {

			for _, refIdentityRow := range refIdentityRows {

				ignoreREF := true

				attrs, err := schema_mappings.GetMappedAttributesFromSchemaMappings(
					displayConfig.AppIDNamePairs[procRef["app"]], 
					displayConfig.TableIDNamePairs[procRef["to_member"]], 
					displayConfig.TableIDNamePairs[procRef["to_member"]] + 
						"." + procRef["to_reference"], 
					displayConfig.AppConfig.AppName,  
					displayConfig.TableIDNamePairs[orgID.member],
					ignoreREF) 
				
				if err != nil {
					log.Println(err)
				}

				log.Println(attrs)

				ignoreREF = false

				attrsToUpdate, err1 := schema_mappings.GetMappedAttributesFromSchemaMappings(
					displayConfig.AppIDNamePairs[procRef["app"]], 
					displayConfig.TableIDNamePairs[procRef["from_member"]], 
					displayConfig.TableIDNamePairs[procRef["from_member"]] +
						"." + procRef["from_reference"], 
					displayConfig.AppConfig.AppName, 
					displayConfig.TableIDNamePairs[refIdentityRow.member],
					ignoreREF)
				
				if err1 != nil {
					log.Println(err1)
				}

				log.Println(attrsToUpdate)
				
				for _, attrToUpdate := range attrsToUpdate {

					err2 := updateReferences(
						displayConfig,
						procRef["pk"],
						displayConfig.TableIDNamePairs[orgID.member], 
						orgID.id, 
						attrs[0], 
						displayConfig.TableIDNamePairs[refIdentityRow.member], 
						refIdentityRow.id, 
						attrToUpdate)
					
					updatedAttrs = append(updatedAttrs, attrToUpdate)

					if err2 != nil {
						log.Println(err2)
					}
				
				}
			}

		} else if procRef["app"] == displayConfig.AppConfig.AppID {

			attr := procRef["to_reference"]

			log.Println(attr)

			ignoreREF := false

			attrsToUpdate, err := schema_mappings.GetMappedAttributesFromSchemaMappings(
				displayConfig.AppIDNamePairs[procRef["app"]], 
				displayConfig.TableIDNamePairs[procRef["from_member"]], 
				displayConfig.TableIDNamePairs[procRef["from_member"]] + 
					"." + procRef["from_reference"], 
				displayConfig.AppConfig.AppName,
				displayConfig.TableIDNamePairs[procRef["to_member"]], 
				ignoreREF)

			if err != nil {
				log.Println(err)
			}

			for _, attrToUpdate := range attrsToUpdate {

				err1 := updateReferences(
					displayConfig,
					procRef["pk"],
					displayConfig.TableIDNamePairs[orgID.member],
					orgID.id,
					attr, 
					displayConfig.TableIDNamePairs[procRef["from_member"]], 
					procRef["from_id"], 
					attrToUpdate)
				
				updatedAttrs = append(updatedAttrs, attrToUpdate)

				if err1 != nil {
					log.Println(err1)
				}

			}
		}

	}

	return updatedAttrs 

}

func resolveReferenceByBackTraversal(displayConfig *config.DisplayConfig, 
	ID *identity, orgID *identity) []string {

	var updatedAttrs []string

	for _, IDRow := range getRowsFromIDTableByTo(displayConfig, ID) {

		procIDRow := transformInterfaceToString(IDRow)
		
		log.Println("id_row: ", procIDRow)

		// You are on the left/from part
		updatedAttrs = append(updatedAttrs, 
			updateMyDataBasedOnReferences(displayConfig, procIDRow, orgID)...)

		// You are on the right/to part
		updatedAttrs = append(updatedAttrs,
			updateOtherDataBasedOnReferences(displayConfig, procIDRow, orgID)...)

		// Traverse back
		preID := createIdentity(
			procIDRow["from_app"], procIDRow["from_member"], procIDRow["from_id"])

		updatedAttrs = append(updatedAttrs,
			resolveReferenceByBackTraversal(displayConfig, preID, orgID)...)
		
	}

	return updatedAttrs

}

func ResolveReference(displayConfig *config.DisplayConfig, hint *app_display.HintStruct) []string {
	
	ID := transformHintToIdenity(displayConfig, hint)
	
	return resolveReferenceByBackTraversal(displayConfig, ID, ID)

}

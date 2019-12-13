package reference_resolution

import (
	"stencil/app_display"
	"stencil/config"
	"stencil/schema_mappings"
	"log"
)

// You are on the left/from part
func updateMyDataBasedOnReferences(displayConfig *config.DisplayConfig, 
	IDRow map[string]string, orgID *identity) map[string]string {
	
	log.Println("You are on the left/from part")

	updatedAttrs := make(map[string]string)

	for _, ref := range getFromReferences(displayConfig, IDRow) {
		
		procRef := transformInterfaceToString(ref)
		log.Println("ref_row: ", procRef)

		data := createIdentity(procRef["app"], procRef["to_member"], procRef["to_id"])

		refIdentityRows := forwardTraverseIDTable(displayConfig, data, orgID)
		// log.Println("refIdentityRows: ", refIdentityRows)

		if len(refIdentityRows) > 0 {
			
			for _, refIdentityRow := range refIdentityRows {

				log.Println("refIdentityRow: ", refIdentityRow)

				// For example, when trying to find
				// the mapped attribute from Diaspora Posts Posts.id to Mastodon Statuses,
				// if we consider arguments in #REF, 
				// there are two results: id and conversation_id which should not be included
				// Basically, the attribute to update other atrributes should not contain #REF
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

				log.Println("attr: ", attrs)

				// There are cases in which no attribute can be found
				// For example: diaspora posts posts.id mastodon media_attachments
				if len(attrs) != 1 {
					
					log.Println(notOneAttributeFound)
					
					continue

				}

				// Basically, the attributes to be updated should always contain #REF
				// Otherwise, the following inputs:
				// "diaspora", "comments", "comments.commentable_id", "mastodon", "statuses", false
				// will return both status_id and id which should not be contained
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
					
					log.Println("attr to be updated:", attrToUpdate)

					updatedVal, err2 := updateReferences(
						displayConfig,
						procRef["pk"], 
						displayConfig.TableIDNamePairs[refIdentityRow.member], 
						refIdentityRow.id, 
						attrs[0], 
						displayConfig.TableIDNamePairs[orgID.member], 
						orgID.id, 
						attrToUpdate)

					if err2 != nil {
						log.Println(err2)
					} else {
						updatedAttrs[attrToUpdate] = updatedVal
					}
					
				}
			}

		} else if procRef["app"] == displayConfig.AppConfig.AppID {

			attr := procRef["to_reference"]

			log.Println("attr: ", attr)

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

			// log.Println(attrsToUpdate)

			for _, attrToUpdate := range attrsToUpdate {

				log.Println("attr to be updated:", attrToUpdate)

				updatedVal, err1 := updateReferences(
					displayConfig,
					procRef["pk"],  
					displayConfig.TableIDNamePairs[procRef["to_member"]], 
					procRef["to_id"], 
					attr, 
					displayConfig.TableIDNamePairs[orgID.member], 
					orgID.id, 
					attrToUpdate)
				
				if err1 != nil {
					log.Println(err1)
				} else {
					updatedAttrs[attrToUpdate] = updatedVal
				}

			}
		}

	}

	return updatedAttrs

}

// You are on the right/to part
func updateOtherDataBasedOnReferences(displayConfig *config.DisplayConfig, 
	IDRow map[string]string, orgID *identity) map[string]string {
	
	log.Println("You are on the right/to part")

	updatedAttrs := make(map[string]string)

	for _, ref := range getToReferences(displayConfig, IDRow) {

		procRef := transformInterfaceToString(ref)
		log.Println(procRef)

		data := createIdentity(procRef["app"], procRef["from_member"], procRef["from_id"])

		refIdentityRows := forwardTraverseIDTable(displayConfig, data, orgID)
		// log.Println(refIdentityRows[0])

		if len(refIdentityRows) > 0 {

			for _, refIdentityRow := range refIdentityRows {

				log.Println("refIdentityRow: ", refIdentityRow)

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

				log.Println("attr: ", attrs)

				if len(attrs) != 1 {
					
					log.Println(notOneAttributeFound)
					
					continue

				}

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

				// log.Println(attrsToUpdate)
				
				for _, attrToUpdate := range attrsToUpdate {

					log.Println("attr to be updated:", attrToUpdate)

					updatedVal, err2 := updateReferences(
						displayConfig,
						procRef["pk"],
						displayConfig.TableIDNamePairs[orgID.member], 
						orgID.id, 
						attrs[0], 
						displayConfig.TableIDNamePairs[refIdentityRow.member], 
						refIdentityRow.id, 
						attrToUpdate)

					if err2 != nil {
						log.Println(err2)
					} else {
						updatedAttrs[refIdentityRow.id + ":" + attrToUpdate] = updatedVal
					}

				}
			}

		} else if procRef["app"] == displayConfig.AppConfig.AppID {

			attr := procRef["to_reference"]

			log.Println("attr: ", attr)

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

				log.Println("attr to be updated:", attrToUpdate)

				updatedVal, err1 := updateReferences(
					displayConfig,
					procRef["pk"],
					displayConfig.TableIDNamePairs[orgID.member],
					orgID.id,
					attr, 
					displayConfig.TableIDNamePairs[procRef["from_member"]], 
					procRef["from_id"], 
					attrToUpdate)

				if err1 != nil {
					log.Println(err1)
				} else {
					updatedAttrs[procRef["from_id"] + ":" + attrToUpdate] = updatedVal
				}

			}
		}

	}

	return updatedAttrs 

}

func resolveReferenceByBackTraversal(displayConfig *config.DisplayConfig, 
	ID *identity, orgID *identity) (map[string]string, map[string]string) {

	myUpdatedAttrs := make(map[string]string)
	
	othersUpdatedAttrs := make(map[string]string)

	for _, IDRow := range getRowsFromIDTableByTo(displayConfig, ID) {

		procIDRow := transformInterfaceToString(IDRow)
		
		log.Println("id_row: ", procIDRow)

		// You are on the left/from part
		currentMyupdatedAttrs := updateMyDataBasedOnReferences(displayConfig, procIDRow, orgID)

		myUpdatedAttrs = combineTwoMaps(myUpdatedAttrs, currentMyupdatedAttrs)

		// You are on the right/to part
		currentOthersUpdatedAttrs := updateOtherDataBasedOnReferences(displayConfig, 
			procIDRow, orgID)
		
		othersUpdatedAttrs = combineTwoMaps(othersUpdatedAttrs, currentOthersUpdatedAttrs)

		// Traverse back
		preID := createIdentity(
			procIDRow["from_app"], procIDRow["from_member"], procIDRow["from_id"])

		nextMyUpdatedAttrs, nextOthersUpdatedAttrs := 
			resolveReferenceByBackTraversal(displayConfig, preID, orgID)
		
		myUpdatedAttrs = combineTwoMaps(myUpdatedAttrs, nextMyUpdatedAttrs)
		othersUpdatedAttrs = combineTwoMaps(othersUpdatedAttrs, nextOthersUpdatedAttrs)

	}

	return myUpdatedAttrs, othersUpdatedAttrs

}

// In terms of the argument *hint*, the first return value is my updated attributes, and 
// the second return value is others' updated attributes.
// Note that my updated attributes will not have collision, but
// others' updated attributes may have some collision, 
// so we use *id:updatedAttr*, which is unique, as the key in the second return value.
func ResolveReference(displayConfig *config.DisplayConfig, 
	hint *app_display.HintStruct) (map[string]string, map[string]string) {
	
	ID := transformHintToIdenity(displayConfig, hint)
	
	return resolveReferenceByBackTraversal(displayConfig, ID, ID)

}

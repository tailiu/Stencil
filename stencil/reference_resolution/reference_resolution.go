package reference_resolution

import (
	"stencil/config"
	"stencil/schema_mappings"
	"log"
)


// Basically, reference resolution traces back to corresponding data 
// in old applications by the identity table, gets any references in an old application
// through the reference table, checks the current application of the reference through
// the identity table again. If the current application is same as the application of the checked
// data, then make some updates.
// Note: migration threads always need to insert rows in the identity table if there does not exist one 
// when encountering #REF or #REF and #FETCH. 

// For example, Diaspora: notifications, notification_actors -> Mastodon: notifications
// Let's focus on the mapping "from_account_id":"#REF(notification_actors.person_id,people.id)" 
// because this is the only place notification_actors is involved. 
// As long as notification_actors are involved in the first argument of #REF, 
// there must exist a row in the identity tabe from Diaspora.notification_actors to 
// Mastodon.notifications if there does not exist (in this case notification_actors is only involved here, 
// so if the #REF does not result in one, there will not be an id row) 
// because otherwise from_account_id cannot be updated. 

// Specifically, when we are trying to resolve references for Mastodon.notifications, 
// we can get a row in the identity table from Diaspora.notification_actors to Mastodon.notifications. 
// Migration threads has already inserted rows (#REF.first_arg -> #REF.second_arg) in the reference table 
// once they encounter #REF. Through the reference (Diaspora.notification_actors.person_id -> Diaspora.people.id),
// display threads find Diaspora.people.id. Then through the id row (Diaspora.people.id -> Mastodon.accounts.id),
// they find Mastodon.accounts.id. Then attr to update other attributes is 'account.id', and the attribute to be updated
// is 'from_account_id' by calling the function GetMappedAttributeFromSchemaMappings(Diaspora, 
// notification_actors, person_id, Mastodon, notification). 

// For #REF(#FETCH($1, $2, $3), $4), migration threads use $2 and $3 to find $1, and then insert a row ($1 -> $4)
// in the reference table. ** Even though the member of $1 is not in the fromTables, there must exist a row in the 
// id table to enable display threads to resolve corresponding reference **.

// For example, Diaspora: photos -> Mastodon: media_attachments
// Let's focus on the mapping "status_id":"#REF(#FETCH(posts.id,posts.guid,photos.status_message_guid),posts.id)".
// Migration threads use photos.status_message_guid and posts.guid to find the post the photo belongs to, 
// insert a row (posts.id -> posts.id) in the reference table, and also insert a row (posts.id -> media_attachments.id)
// in the id table. Then display threads follow the path: media_attachments.id  -> id table -> posts.id -> ref table ->
// posts.id -> id table -> statuses.id 
// attr: statuses.id
// attr to be updated: 
// GetMappedAttributeFromSchemaMappings(Diaspora, posts, id, Mastodon, media_attachments) -> status_id


// You are on the left/from part
func updateMyDataBasedOnReferences(displayConfig *config.DisplayConfig, 
	IDRow map[string]string, orgID *Identity) map[string]string {
	
	log.Println("You are on the left/from part")

	updatedAttrs := make(map[string]string)

	for _, ref := range getFromReferences(displayConfig, IDRow) {
		
		procRef := transformInterfaceToString(ref)
		log.Println("ref_row: ", procRef)

		data := CreateIdentity(procRef["app"], procRef["to_member"], procRef["to_id"])

		refIdentityRows := forwardTraverseIDTable(displayConfig, data, orgID)
		// log.Println("refIdentityRows: ", refIdentityRows)

		if len(refIdentityRows) > 0 {
			
			for _, refIdentityRow := range refIdentityRows {

				log.Println("refIdentityRow: ", refIdentityRow)

				combineTwoMaps(updatedAttrs, updateRefOnLeftUsingRefIDRow(displayConfig, 
					refIdentityRow, procRef, orgID))

			}

		} else if procRef["app"] == displayConfig.AppConfig.AppID {

			combineTwoMaps(updatedAttrs,
				 updateRefOnLeftNotUsingRefIDRow(displayConfig, procRef, orgID))

		}

	}

	return updatedAttrs

}

// You are on the right/to part
func updateOtherDataBasedOnReferences(displayConfig *config.DisplayConfig, 
	IDRow map[string]string, orgID *Identity) map[string]string {
	
	log.Println("You are on the right/to part")

	updatedAttrs := make(map[string]string)

	for _, ref := range getToReferences(displayConfig, IDRow) {

		procRef := transformInterfaceToString(ref)
		log.Println(procRef)

		data := CreateIdentity(procRef["app"], procRef["from_member"], procRef["from_id"])

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
	ID *Identity, orgID *Identity) (map[string]string, map[string]string) {

	myUpdatedAttrs := make(map[string]string)
	
	othersUpdatedAttrs := make(map[string]string)

	for _, IDRow := range getRowsFromIDTableByTo(displayConfig, ID) {

		procIDRow := transformInterfaceToString(IDRow)
		
		log.Println("id_row: ", procIDRow)

		// You are on the left/from part
		currentMyupdatedAttrs := updateMyDataBasedOnReferences(displayConfig, procIDRow, orgID)

		combineTwoMaps(myUpdatedAttrs, currentMyupdatedAttrs)

		// You are on the right/to part
		currentOthersUpdatedAttrs := updateOtherDataBasedOnReferences(displayConfig, 
			procIDRow, orgID)
		
		combineTwoMaps(othersUpdatedAttrs, currentOthersUpdatedAttrs)

		// Traverse back
		preID := CreateIdentity(
			procIDRow["from_app"], procIDRow["from_member"], procIDRow["from_id"])

		nextMyUpdatedAttrs, nextOthersUpdatedAttrs := 
			resolveReferenceByBackTraversal(displayConfig, preID, orgID)
		
		combineTwoMaps(myUpdatedAttrs, nextMyUpdatedAttrs)
		combineTwoMaps(othersUpdatedAttrs, nextOthersUpdatedAttrs)

	}

	return myUpdatedAttrs, othersUpdatedAttrs

}

// In terms of the argument *hint*, the first return value is my updated attributes, and 
// the second return value is others' updated attributes.
// Note that my updated attributes will not have collision, but
// others' updated attributes may have some collision, 
// so we use *id:updatedAttr*, which is unique, as the key in the second return value.
func ResolveReference(displayConfig *config.DisplayConfig, 
	ID *Identity) (map[string]string, map[string]string) {
	
	return resolveReferenceByBackTraversal(displayConfig, ID, ID)

}

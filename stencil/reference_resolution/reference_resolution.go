package reference_resolution

import (
	"log"
	"stencil/config"
	"database/sql"
)

/**
 *
 * Basically, reference resolution traces back to corresponding data 
 * in old applications by the identity table, gets any references in an old application
 * through the reference table, checks the current application of the reference through
 * the identity table again. If the current application is same as the application of the checked
 * data, then make some updates.
 * Note: migration threads always need to insert rows in the identity table if there does not exist one 
 * when encountering #REF or #REF and #FETCH. 
 *
 * For now reference resolution only tries to link data with other data with *new IDs*.
 * We can easily generalize the reference resolution by storing other data changes other than
 * ID changes in some new table.
 *
 */

/**
 *
 * For example, Diaspora: notifications, notification_actors -> Mastodon: notifications
 * Let's focus on the mapping "from_account_id":"#REF(notification_actors.person_id,people.id)" 
 * because this is the only place notification_actors is involved. 
 * As long as notification_actors are involved in the first argument of #REF, 
 * there must exist a row in the identity tabe from Diaspora.notification_actors to 
 * Mastodon.notifications if there does not exist (in this case notification_actors is only involved here, 
 * so if the #REF does not result in one, there will not be an id row) 
 * because otherwise from_account_id cannot be updated. 
 *
 */

/**
 *
 * Specifically, when we are trying to resolve references for Mastodon.notifications, 
 * we can get a row in the identity table from Diaspora.notification_actors to Mastodon.notifications. 
 * Migration threads has already inserted rows (#REF.first_arg -> #REF.second_arg) in the reference table 
 * once they encounter #REF. Through the reference (Diaspora.notification_actors.person_id -> Diaspora.people.id),
 * display threads find Diaspora.people.id. Then through the id row (Diaspora.people.id -> Mastodon.accounts.id),
 * they find Mastodon.accounts.id. Then attr to update other attributes is 'account.id', and the attribute to be updated
 * is 'from_account_id' by calling the function GetMappedAttributeFromSchemaMappings(Diaspora, 
 * notification_actors, person_id, Mastodon, notification). 
 *
 */

/**
 *
 * For #REF(#FETCH($1, $2, $3), $4), migration threads use $2 and $3 to find $1, and then insert a row ($1 -> $4)
 * in the reference table. ** Even though the member of $1 is not in the fromTables, there must exist a row in the 
 * id table to enable display threads to resolve corresponding reference **.
 *
 */

/**
 *
 * For example, Diaspora: photos -> Mastodon: media_attachments
 * Let's focus on the mapping "status_id":"#REF(#FETCH(posts.id,posts.guid,photos.status_message_guid),posts.id)".
 * Migration threads use photos.status_message_guid and posts.guid to find the post the photo belongs to, 
 * insert a row (posts.id -> posts.id) in the reference table, and also insert a row (posts.id -> media_attachments.id)
 * in the id table. Then display threads follow the path: media_attachments.id  -> id table -> posts.id -> ref table ->
 * posts.id -> id table -> statuses.id 
 * attr: statuses.id
 * attr to be updated: 
 * GetMappedAttributeFromSchemaMappings(Diaspora, posts, id, Mastodon, media_attachments) -> status_id
 *
 */
 

// You are on the left/from part
func updateMyDataBasedOnReferences(refResolutionConfig *RefResolutionConfig, 
	IDRow map[string]string, orgID *Identity) map[string]string {
	
	log.Println("You are on the left/from part")

	updatedAttrs := make(map[string]string)
	
	fromRefs := getFromReferences(refResolutionConfig, IDRow)

	log.Println("Get", len(fromRefs), "from reference(s)")

	for _, ref := range fromRefs {
		
		procRef := transformInterfaceToString(ref)
		log.Println("ref_row: ", procRef)

		data := CreateIdentity(procRef["app"], procRef["to_member"], procRef["to_id"])

		refIdentityRows := forwardTraverseIDTable(refResolutionConfig, data, orgID)
		// log.Println("refIdentityRows: ", refIdentityRows)

		log.Println("After traversing forward the ID table:")
		log.Println("Get", len(refIdentityRows), "refIdentity row(s)")

		log.Println("procRef['app']:", procRef["app"], 
			"refResolutionConfig.appID:", refResolutionConfig.appID)

		if len(refIdentityRows) > 0 {
			
			for _, refIdentityRow := range refIdentityRows {

				log.Println("refIdentityRow: ", refIdentityRow)

				oneUpdatedAttr := updateRefOnLeftBasedOnMappingsUsingRefIDRow(
					refResolutionConfig, refIdentityRow, procRef, orgID)

				combineTwoMaps(updatedAttrs, oneUpdatedAttr)

				if len(oneUpdatedAttr) > 0 {
					break
				}

			}

		} else if procRef["app"] == refResolutionConfig.appID {

			oneUpdatedAttr := updateRefOnLeftBasedOnMappingsNotUsingRefIDRow(
				refResolutionConfig, procRef, orgID)

			combineTwoMaps(updatedAttrs, oneUpdatedAttr)

		}

	}

	return updatedAttrs

}

// You are on the right/to part
func updateOtherDataBasedOnReferences(refResolutionConfig *RefResolutionConfig, 
	IDRow map[string]string, orgID *Identity) map[string]string {
	
	log.Println("You are on the right/to part")

	updatedAttrs := make(map[string]string)
	
	toRefs := getToReferences(refResolutionConfig, IDRow)

	log.Println("Get", len(toRefs), "to reference(s)")

	for _, ref := range toRefs {

		procRef := transformInterfaceToString(ref)
		log.Println(procRef)

		data := CreateIdentity(procRef["app"], procRef["from_member"], procRef["from_id"])

		refIdentityRows := forwardTraverseIDTable(refResolutionConfig, data, orgID)
		// log.Println(refIdentityRows[0])

		log.Println("After traversing forward the ID table:")
		log.Println("Get", len(refIdentityRows), "refIdentity row(s)")

		log.Println("procRef['app']:", procRef["app"], 
			"refResolutionConfig.appID:", refResolutionConfig.appID)

		if len(refIdentityRows) > 0 {

			for _, refIdentityRow := range refIdentityRows {

				log.Println("refIdentityRow: ", refIdentityRow)

				oneUpdatedAttr := updateRefOnRightBasedOnMappingsUsingRefIDRow(
					refResolutionConfig, refIdentityRow, procRef, orgID)

				combineTwoMaps(updatedAttrs, oneUpdatedAttr)
				
				if len(oneUpdatedAttr) > 0 {
					break
				}
			}

		} else if procRef["app"] == refResolutionConfig.appID {

			oneUpdatedAttr := updateRefOnRightBasedOnMappingsNotUsingRefIDRow(
				refResolutionConfig, procRef, orgID)

			combineTwoMaps(updatedAttrs, oneUpdatedAttr)

		}
	}

	return updatedAttrs 

}

func resolveReferenceByBackTraversal(refResolutionConfig *RefResolutionConfig, 
	ID *Identity, orgID *Identity) (map[string]string, map[string]string) {

	myUpdatedAttrs := make(map[string]string)
	
	othersUpdatedAttrs := make(map[string]string)
	
	idRows := getRowsFromIDTableByTo(refResolutionConfig, ID)

	log.Println("Get", len(idRows), "id row(s)")

	for _, IDRow := range idRows {

		procIDRow := transformInterfaceToString(IDRow)
		
		log.Println("id_row: ", procIDRow)

		// You are on the left/from part
		currentMyupdatedAttrs := updateMyDataBasedOnReferences(refResolutionConfig, procIDRow, orgID)

		combineTwoMaps(myUpdatedAttrs, currentMyupdatedAttrs)

		// You are on the right/to part
		currentOthersUpdatedAttrs := updateOtherDataBasedOnReferences(refResolutionConfig, 
			procIDRow, orgID)
		
		combineTwoMaps(othersUpdatedAttrs, currentOthersUpdatedAttrs)

		// Traverse back
		preID := CreateIdentity(
			procIDRow["from_app"], procIDRow["from_member"], procIDRow["from_id"])

		nextMyUpdatedAttrs, nextOthersUpdatedAttrs := 
			resolveReferenceByBackTraversal(refResolutionConfig, preID, orgID)
		
		combineTwoMaps(myUpdatedAttrs, nextMyUpdatedAttrs)
		combineTwoMaps(othersUpdatedAttrs, nextOthersUpdatedAttrs)

	}

	return myUpdatedAttrs, othersUpdatedAttrs

}


func InitializeReferenceResolution(migrationID int, 
	appID, appName string, appDBConn, StencilDBConn *sql.DB, 
	appTableNameIDPairs map[string]string,
	appIDNamePairs map[string]string,
	tableIDNamePairs map[string]string,
	allMappings *config.SchemaMappings,
	mappingsFromSrcToDst *config.MappedApp) *RefResolutionConfig {

	var refResolutionConfig RefResolutionConfig

	refResolutionConfig.stencilDBConn = StencilDBConn
	refResolutionConfig.appDBConn = appDBConn
	refResolutionConfig.migrationID = migrationID
	refResolutionConfig.appID = appID
	refResolutionConfig.appName = appName
	refResolutionConfig.allMappings = allMappings
	refResolutionConfig.appTableNameIDPairs = appTableNameIDPairs
	refResolutionConfig.appIDNamePairs = appIDNamePairs
	refResolutionConfig.tableIDNamePairs = tableIDNamePairs
	refResolutionConfig.mappingsFromSrcToDst = mappingsFromSrcToDst
	
	return &refResolutionConfig
}

// In terms of the argument *ID*, the first return value is my updated attributes, and 
// the second return value is others' updated attributes.
// Note that my updated attributes will not have collision because a table does not have
// duplicate attributes, 
// but others' updated attributes may have some collision, 
// so we use *id:updatedAttr*, which is unique, as the key in the second return value.
func ResolveReference(refResolutionConfig *RefResolutionConfig, 
	ID *Identity) (map[string]string, map[string]string) {
	
	return resolveReferenceByBackTraversal(refResolutionConfig, ID, ID)

}

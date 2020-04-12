package reference_resolution_v2

import (
	"database/sql"
	"log"
	"stencil/common_funcs"
	"stencil/config"
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

/**
 *
 * Actually, for reference resolution, the handling of #REF and #REFHARD is the same,
 * because it assumes that when creating a reference row, the id (from_id) of the first arg refers to
 * the id (to_id) of the second arg, which is a general case. However, the migration algo does in a
 * different way, in #REF, it is convenient to directly uses the value of the first arg as the to_id,
 * while in #REFHARD, it actually uses the general way, i.e., using the id of the second arg, to set to_id.
 *
 */

// You are on the left/from part
func (rr *RefResolution) updateMyDataBasedOnReferences(attrRow map[string]string, orgAttr *Attribute) map[string]string {

	log.Println("You are on the left/from part")

	updatedAttrs := make(map[string]string)

	// Using ID in getting from references here is to remove possible duplicate reference rows
	// Otherwise, all the references will be resolved and reference rows will be deleted
	// For example, multiple comments reply to the same post, but only the attribute in one comment should be resolved
	fromRefs := rr.getFromReferencesUsingID(attrRow)

	log.Println("Get", len(fromRefs), "from reference(s)")

	for _, ref := range fromRefs {

		procRef := common_funcs.TransformInterfaceToString(ref)

		rr.LogRefRow(procRef)

		// There is no to_id in the reference table, but to_val is enough here 
		// since the to_val of the referenced data should not be STENCIL_NULL 
		data := CreateAttribute(
			procRef["app"],
			procRef["to_member"],
			procRef["to_attr"],
			procRef["to_val"],
		)

		refAttributeRows := rr.forwardTraverseAttrChangesTable(data, orgAttr, false)

		log.Println("After traversing forward the attribute_changes table:")
		log.Println("Get", len(refAttributeRows), "refAttribute row(s)")

		if len(refAttributeRows) > 0 {

			for _, refAttributeRow := range refAttributeRows {

				rr.logRefAttrRow(refAttributeRow)

				oneUpdatedAttr := make(map[string]string)

				// If we can get refIdentityRows, but the app of the reference table is the same as
				// the destination application, this means that the data is migrated back
				// Since there could be multiple pieces of data migrated back, we only update data
				// based on the table of the reference row
				if procRef["app"] == rr.appID {

					log.Println("The data has been migrated back to the dest app1")

					// Up to now it seems that refIdentityRow.member == procRef["to_member"] is enough,
					// but procRef["from_member"] == orgID.member is added in case
					if refAttributeRow.member == procRef["to_member"] &&
						refAttributeRow.attrName == procRef["to_attr"] &&
						procRef["from_member"] == orgAttr.member &&
						procRef["from_attr"] == orgAttr.attrName {

						oneUpdatedAttr = updateRefOnLeftByRefAttrRow1(
							rr, procRef, orgAttr, refAttributeRow.val)

						updatedAttrs = combineTwoMaps(updatedAttrs, oneUpdatedAttr)

					} else {
						log.Println("The migrated data is not in the desired reference table1")
					}

				} else {

					oneUpdatedAttr = updateRefOnLeftByRefAttrRow(
						rr, refAttributeRow, procRef, orgAttr)

					updatedAttrs = combineTwoMaps(updatedAttrs, oneUpdatedAttr)

				}

				if len(oneUpdatedAttr) > 0 {
					break
				}

			}

			// If we cannot get any refIdentityRows which means the referenced row has not been migrated
			// and the application of that row is exactly the destination app, then we can directly
			// use the id in the reference row to update attributes
			// Checking the from member in the reference row is important because
			// there could be cases where after back traversal there are multiple rows
			// For example, when trying to update notification_actors and after back traversal,
			// we get two reference rows for notifications and notification_actors,
			// then notifications should be igored
		} else if procRef["app"] == rr.appID &&
			procRef["from_member"] == orgAttr.member &&
			procRef["from_attr"] == orgAttr.attrName {

			log.Println("The data has not been migrated and is in the dest app1")

			oneUpdatedAttr := updateRefOnLeftNotUsingRefAttrRow(
				rr, procRef, orgAttr)

			updatedAttrs = combineTwoMaps(updatedAttrs, oneUpdatedAttr)

		}

	}

	return updatedAttrs

}

// You are on the right/to part
func (rr *RefResolution) updateOtherDataBasedOnReferences(attrRow map[string]string, orgAttr *Attribute) map[string]string {

	log.Println("You are on the right/to part")

	updatedAttrs := make(map[string]string)

	toRefs := rr.getToReferences(attrRow)

	log.Println("Get", len(toRefs), "to reference(s)")

	for _, ref := range toRefs {

		procRef := common_funcs.TransformInterfaceToString(ref)

		rr.LogRefRow(procRef)

		data := CreateAttribute(
			procRef["app"],
			procRef["from_member"],
			procRef["from_attr"],
			procRef["from_val"],
			procRef["from_id"],
		)

		refAttributeRows := rr.forwardTraverseAttrChangesTable(data, orgAttr, false)

		log.Println("After traversing forward the attribute_changes table:")
		log.Println("Get", len(refAttributeRows), "refAttribute row(s)")

		if len(refAttributeRows) > 0 {

			for _, refAttributeRow := range refAttributeRows {

				rr.logRefAttrRow(refAttributeRow)

				oneUpdatedAttr := make(map[string]string)

				if procRef["app"] == rr.appID {

					log.Println("The data has been migrated back to the dest app2")

					if refAttributeRow.member == procRef["from_member"] &&
						refAttributeRow.attrName == procRef["from_attr"] &&
						procRef["to_member"] == orgAttr.member &&
						procRef["to_attr"] == orgAttr.attrName {

						oneUpdatedAttr = updateRefOnRightByRefAttrRow1(
							rr, procRef, orgAttr, refAttributeRow.val)

						updatedAttrs = combineTwoMaps(updatedAttrs, oneUpdatedAttr)

					} else {
						log.Println("The migrated data is not in the desired reference table2")
					}

				} else {

					oneUpdatedAttr = updateRefOnRightByRefAttrRow(
						rr, refAttributeRow, procRef, orgAttr)

					updatedAttrs = combineTwoMaps(updatedAttrs, oneUpdatedAttr)

				}

				if len(oneUpdatedAttr) > 0 {
					break
				}
			}

		} else if procRef["app"] == rr.appID &&
			procRef["to_member"] == orgAttr.member &&
			procRef["to_attr"] == orgAttr.attrName {

			log.Println("The data has not been migrated and is in the dest app2")

			oneUpdatedAttr := updateRefOnRightNotUsingRefAttrRow(
				rr, procRef, orgAttr)

			updatedAttrs = combineTwoMaps(updatedAttrs, oneUpdatedAttr)

		}
	}

	return updatedAttrs

}

func (rr *RefResolution) resolveReferenceByBackTraversal(attr, orgAttr *Attribute) (map[string]string, map[string]string) {

	log.Println("Resolve references by back traversal")

	myUpdatedAttrs := make(map[string]string)

	othersUpdatedAttrs := make(map[string]string)

	attrRows := rr.getRowsFromAttrChangesTableByTo(attr)

	log.Println("Get", len(attrRows), "attribute row(s)")

	for _, attrRow := range attrRows {

		procAttrRow := common_funcs.TransformInterfaceToString(attrRow)

		rr.logAttrChangeRow(procAttrRow)

		// You are on the left/from part
		currentMyupdatedAttrs := rr.updateMyDataBasedOnReferences(procAttrRow, orgAttr)

		myUpdatedAttrs = combineTwoMaps(myUpdatedAttrs, currentMyupdatedAttrs)

		// You are on the right/to part
		currentOthersUpdatedAttrs := rr.updateOtherDataBasedOnReferences(
			procAttrRow, orgAttr)

		othersUpdatedAttrs = combineTwoMaps(othersUpdatedAttrs, currentOthersUpdatedAttrs)

		// Traverse back
		preAttr := CreateAttribute(
			procAttrRow["from_app"], procAttrRow["from_member"],
			procAttrRow["from_attr"], procAttrRow["from_val"],
		)

		nextMyUpdatedAttrs, nextOthersUpdatedAttrs := rr.resolveReferenceByBackTraversal(preAttr, orgAttr)

		myUpdatedAttrs = combineTwoMaps(myUpdatedAttrs, nextMyUpdatedAttrs)
		othersUpdatedAttrs = combineTwoMaps(othersUpdatedAttrs, nextOthersUpdatedAttrs)

	}

	return myUpdatedAttrs, othersUpdatedAttrs

}

func InitializeReferenceResolution(migrationID int,
	appID, appName string, appDBConn, StencilDBConn *sql.DB,
	appTableNameIDPairs map[string]string,
	appIDNamePairs map[string]string,
	tableIDNamePairs map[string]string,
	attrIDNamePairs map[string]string,
	appAttrNameIDPairs map[string]string,
	allMappings *config.SchemaMappings,
	mappingsFromSrcToDst *config.MappedApp,
	mappingsFromOtherAppsToDst map[string]*config.MappedApp) *RefResolution {

	var rr RefResolution

	rr.stencilDBConn = StencilDBConn
	rr.appDBConn = appDBConn
	rr.migrationID = migrationID
	rr.appID = appID
	rr.appName = appName
	rr.appTableNameIDPairs = appTableNameIDPairs
	rr.appIDNamePairs = appIDNamePairs
	rr.tableIDNamePairs = tableIDNamePairs
	rr.attrIDNamePairs = attrIDNamePairs
	rr.appAttrNameIDPairs = appAttrNameIDPairs
	rr.allMappings = allMappings
	rr.mappingsFromSrcToDst = mappingsFromSrcToDst
	rr.mappingsFromOtherAppsToDst = mappingsFromOtherAppsToDst

	return &rr
}

// In terms of the argument *ID*, the first return value is my updated attributes, and
// the second return value is others' updated attributes.
// Note that my updated attributes will not have collision because a table does not have
// duplicate attributes,
// but others' updated attributes may have some collision,
// so we use *id:updatedAttr*, which is unique, as the key in the second return value.
func (rr *RefResolution) ResolveReference(attr *Attribute) (map[string]string, map[string]string) {

	return rr.resolveReferenceByBackTraversal(attr, attr)

}

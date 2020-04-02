package reference_resolution_v2

import (
	"stencil/schema_mappings"
	"log"
)

/**
 *
 * Actually, we can use two possible ways to update references: 
 * by dependencies (inter and intra-node deps) and by schema mappings. 
 * Simply updating refs by dependencies is not sufficient because when migrating to a new 
 * app, the deps of the source app and the destination app are different, and inserting and 
 * updating references based on the source app could be both incorrect and insufficient for the destination app,
 * For example, a node (N1) depends on another node N2 in the source app, so we insert N1 -> N2
 * in the reference table if we use deps in the source app. However, N1 deps on N3 in the destination app,
 * and N1 and N2 happen to be migrated to the destination app. Then in this case, the reference will be updated
 * incorrectly. (Although we cannot find a scenario in our test apps, this might happen)
 * so actually we should not do that.
 * Thus we use #REF in the schema mappings in which both the source and destination
 * apps need to be involved to indicate how we change refs in the destination app.
 * Specifying #REF needs to think in two steps:
 * 1. How does the attr1 in a destination app depends on another attr2 in the destination app?
 * 2. Where does attr1 and attr2 come from in the source app?
 *
 */


// It should be noted that since now we are using the attribute_changes table instead of identity tables,
// we are sure about which attributes to update or to be updated,
// and there is no need for us to use get attr or attrToUpdate by mappings.
// But in the display algorithm, we have to check all the attributes to be updated based on mappings.  
func updateRefOnLeftByRefAttrRow(refResolutionConfig *RefResolutionConfig, 
	refAttributeRow *Attribute, procRef map[string]string, orgAttr *Attribute) map[string]string {

	updatedAttr := make(map[string]string)

	attr := refResolutionConfig.attrIDNamePairs(refAttributeRow.attrName)
	attrToUpdate := refResolutionConfig.attrIDNamePairs(orgAttr.attrName)

	log.Println("attr:", attr)
	log.Println("attr to be updated:", attrToUpdate)

	updatedVal, err1 := updateReferences(
		refResolutionConfig,
		procRef["pk"], 
		refResolutionConfig.tableIDNamePairs[refAttributeRow.member], 
		refAttributeRow.val, 
		attr, 
		refResolutionConfig.tableIDNamePairs[orgAttr.member], 
		orgAttr.val, 
		attrToUpdate,
	)

	if err1 != nil {
		log.Println(err1)
	} else {
		updatedAttr[attrToUpdate] = updatedVal
	}

	return updatedAttr

}

func updateRefOnLeftByRefAttrRow1(refResolutionConfig *RefResolutionConfig, 
	procRef map[string]string, orgAttr *Attribute, refAttrRowVal string) map[string]string {

	attr := refResolutionConfig.attrIDNamePairs(procRef["to_attr"])
	attrToUpdate := refResolutionConfig.attrIDNamePairs(procRef["from_attr"])

	updatedAttr := make(map[string]string)

	log.Println("attr:", attr)
	log.Println("attr to be updated:", attrToUpdate)

	updatedVal, err1 := updateReferences(
		refResolutionConfig,
		procRef["pk"],  
		refResolutionConfig.tableIDNamePairs[procRef["to_member"]], 
		refAttrRowVal, 
		attr, 
		refResolutionConfig.tableIDNamePairs[orgAttr.member], 
		orgAttr.val, 
		attrToUpdate,
	)
	
	if err1 != nil {
		log.Println(err1)
	} else {
		updatedAttr[attrToUpdate] = updatedVal
	}

	return updatedAttr
}

func updateRefOnLeftNotUsingRefAttrRow(refResolutionConfig *RefResolutionConfig, 
	procRef map[string]string, orgAttr *Attribute) map[string]string {

	attr := refResolutionConfig.attrIDNamePairs(procRef["to_attr"])
	attrToUpdate := refResolutionConfig.attrIDNamePairs(procRef["from_attr"])

	updatedAttr := make(map[string]string)

	log.Println("attr:", attr)
	log.Println("attr to be updated:", attrToUpdate)

	updatedVal, err1 := updateReferences(
		refResolutionConfig,
		procRef["pk"],  
		refResolutionConfig.tableIDNamePairs[procRef["to_member"]], 
		procRef["to_val"], 
		attr, 
		refResolutionConfig.tableIDNamePairs[orgAttr.member], 
		orgAttr.val, 
		attrToUpdate,
	)
	
	if err1 != nil {
		log.Println(err1)
	} else {
		updatedAttr[attrToUpdate] = updatedVal
	}

	return updatedAttr
}

func updateRefOnRightBasedOnMappingsUsingRefIDRow(refResolutionConfig *RefResolutionConfig, 
	refIdentityRow *Identity, procRef map[string]string, orgID *Identity) map[string]string {

	updatedAttrs := make(map[string]string)

	// ignoreREF := true

	attrs, err := schema_mappings.GetMappedAttributesToUpdateOthers(
		refResolutionConfig.allMappings,
		refResolutionConfig.appIDNamePairs[procRef["app"]], 
		refResolutionConfig.tableIDNamePairs[procRef["to_member"]], 
		refResolutionConfig.tableIDNamePairs[procRef["to_member"]] + 
			"." + procRef["to_reference"], 
		refResolutionConfig.appName,  
		refResolutionConfig.tableIDNamePairs[orgID.member]) 
	
	if err != nil {
		log.Println("Error in Getting attributes to update other attributes:")
		log.Println(err)

		return nil
	}

	log.Println("attr:", attrs)

	if len(attrs) != 1 {
		
		log.Println(notOneAttributeFound)
		
		return nil

	}

	attrsToUpdateNotInFETCH := make(map[string]string)
	attrsToUpdateInFETCH := make(map[string]string)

	var err1, err2 error

	// ignoreREF = false

	attrsToUpdateNotInFETCH, err1 = schema_mappings.GetMappedAttributesToBeUpdated(
		refResolutionConfig.allMappings,
		refResolutionConfig.appIDNamePairs[procRef["app"]], 
		refResolutionConfig.tableIDNamePairs[procRef["from_member"]], 
		refResolutionConfig.tableIDNamePairs[procRef["from_member"]] +
			"." + procRef["from_reference"], 
		refResolutionConfig.appName, 
		refResolutionConfig.tableIDNamePairs[refIdentityRow.member])
	
	if err1 != nil {
		log.Println("Error in Getting mapped attributes from schema mappings:")
		log.Println(err1)
	}

	// log.Println("total attrs to be updated:", attrsToUpdateNotInFETCH)

	attrsToUpdateInFETCH, err2 = schema_mappings.GetMappedAttributesToBeUpdatedByFETCH(
		refResolutionConfig.allMappings,
		refResolutionConfig.appIDNamePairs[procRef["app"]], 
		refResolutionConfig.tableIDNamePairs[procRef["from_member"]] +
			"." + procRef["from_reference"], 
		refResolutionConfig.appName,
		refResolutionConfig.tableIDNamePairs[refIdentityRow.member])

	if err2 != nil {
		log.Println("Error in Getting mapped attributes from schema mappings by #FETCH:")
		log.Println(err2)

	}

	attrsToUpdate := combineTwoMaps(attrsToUpdateNotInFETCH, attrsToUpdateInFETCH)

	log.Println("total attrs to be updated:",attrsToUpdate)
	
	for attrToUpdate, thirdArgInREF := range attrsToUpdate {

		log.Println("one attr to be checked and updated:", attrToUpdate)
		log.Println("Third argument in #REF:", thirdArgInREF)

		// For example,
		// diaspora posts posts.id mastodon conversations
		// attr:  [id]
		// diaspora likes likes.target_id mastodon favourites
		// total attrs to be updated: [status_id]
		// Obviously, if there is no the third argument (statuses in this example) indicating that
		// it is the statuses table not the conversations table should update status_id,
		// then there will be errors
		if thirdArgInREF != "" && thirdArgInREF != 
			refResolutionConfig.tableIDNamePairs[orgID.member] {
			
			log.Println("Third argument in #REF", 
				thirdArgInREF, "is not equal to toTable", 
				refResolutionConfig.tableIDNamePairs[orgID.member])
			
			continue
		}

		updatedVal, err2 := updateReferences(
			refResolutionConfig,
			procRef["pk"],
			refResolutionConfig.tableIDNamePairs[orgID.member], 
			orgID.id, 
			attrs[0], 
			refResolutionConfig.tableIDNamePairs[refIdentityRow.member], 
			refIdentityRow.id, 
			attrToUpdate)

		if err2 != nil {
			
			log.Println(err2)
		
		} else {

			updatedAttrs[refIdentityRow.id + ":" + attrToUpdate] = updatedVal

			// This is an important break because one reference can only
			// update one value
			break
		}

	}

	return updatedAttrs
	
}

func updateRefOnRightByRefAttrRow1(refResolutionConfig *RefResolutionConfig, 
	procRef map[string]string, orgID *Identity, refIdentityRowID string) map[string]string {
	
	attr := procRef["to_reference"]
	attrToUpdate := procRef["from_reference"]

	updatedAttr := make(map[string]string)

	log.Println("attr:", attr)
	log.Println("attr to be updated:", attrToUpdate)

	updatedVal, err1 := updateReferences(
		refResolutionConfig,
		procRef["pk"],
		refResolutionConfig.tableIDNamePairs[orgID.member],
		orgID.id,
		attr, 
		refResolutionConfig.tableIDNamePairs[procRef["from_member"]], 
		refIdentityRowID, 
		attrToUpdate,
	)

	if err1 != nil {
		
		log.Println(err1)
	
	} else {
		
		updatedAttr[refIdentityRowID + ":" + attrToUpdate] = updatedVal
	}

	return updatedAttr 
	
}

func updateRefOnRightBasedOnMappingsNotUsingRefIDRow(refResolutionConfig *RefResolutionConfig, 
	procRef map[string]string, orgID *Identity) map[string]string {
	
	attr := procRef["to_reference"]
	attrToUpdate := procRef["from_reference"]

	updatedAttr := make(map[string]string)

	log.Println("attr:", attr)
	log.Println("attr to be updated:", attrToUpdate)

	updatedVal, err1 := updateReferences(
		refResolutionConfig,
		procRef["pk"],
		refResolutionConfig.tableIDNamePairs[orgID.member],
		orgID.id,
		attr, 
		refResolutionConfig.tableIDNamePairs[procRef["from_member"]], 
		procRef["from_id"], 
		attrToUpdate,
	)

	if err1 != nil {
		
		log.Println(err1)
	
	} else {
		
		updatedAttr[procRef["from_id"] + ":" + attrToUpdate] = updatedVal
	}

	return updatedAttr 
	
}


package reference_resolution

import (
	"stencil/config"
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


func updateRefOnLeftUsingRefIDRow(displayConfig *config.DisplayConfig, 
	refIdentityRow *Identity, procRef map[string]string, orgID *Identity) map[string]string {
	
	updatedAttrs := make(map[string]string)

	// For example, when trying to find
	// the mapped attribute from Diaspora Posts Posts.id to Mastodon Statuses,
	// if we consider arguments in #REF, 
	// there are two results: id and conversation_id which should not be included
	// Basically, the attribute to update other atrributes should not contain #REF
	ignoreREF := true

	attrs, err := schema_mappings.GetMappedAttributesFromSchemaMappings(
		displayConfig.AllMappings,
		displayConfig.AppIDNamePairs[procRef["app"]], 
		displayConfig.TableIDNamePairs[procRef["to_member"]],
		displayConfig.TableIDNamePairs[procRef["to_member"]] + 
			"." + procRef["to_reference"], 
		displayConfig.AppConfig.AppName,
		displayConfig.TableIDNamePairs[refIdentityRow.member],
		ignoreREF)

	if err != nil {
		
		log.Println(err)

		return nil

	}

	log.Println("attr: ", attrs)

	// If #FETCH is ignored, there could be cases in which no attribute can be found.
	// For example: diaspora posts posts.id mastodon media_attachments. This is caused
	// by wrong implementation.
	if len(attrs) != 1 {
		
		log.Println(notOneAttributeFound)
		
		return nil

	}

	var attrsToUpdate, attrsToUpdateInFETCH []string

	var err1, err2 error

	// Basically, the attributes to be updated should always contain #REF
	// Otherwise, the following inputs:
	// "diaspora", "comments", "comments.commentable_id", "mastodon", "statuses", false
	// will return both status_id and id which should not be contained
	ignoreREF = false 

	attrsToUpdate, err1 = schema_mappings.GetMappedAttributesFromSchemaMappings(
		displayConfig.AllMappings,
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

	// log.Println("attrsToUpdate:",attrsToUpdate)

	// #FETCH case is different from the normal cases.
	// For example: diaspora posts posts.id mastodon media_attachments, 
	// in this case, the mappings do not contain the posts table, and 
	// the first argument (posts.id) of #FETCH is needed to be used to resolve photo.status_id
	attrsToUpdateInFETCH, err2 = schema_mappings.GetMappedAttributesFromSchemaMappingsByFETCH(
		displayConfig.AllMappings,
		displayConfig.AppIDNamePairs[procRef["app"]], 
		displayConfig.TableIDNamePairs[procRef["from_member"]] +
			"." + procRef["from_reference"], 
		displayConfig.AppConfig.AppName,
		displayConfig.TableIDNamePairs[orgID.member])

	if err2 != nil {

		log.Println(err2)

	}

	// log.Println("attrsToUpdateInFETCH:", attrsToUpdateInFETCH)

	attrsToUpdate = append(attrsToUpdate, attrsToUpdateInFETCH...)

	log.Println("attrsToUpdate:",attrsToUpdate)

	for _, attrToUpdate := range attrsToUpdate {
		
		log.Println("attr to be updated:", attrToUpdate)

		updatedVal, err3 := updateReferences(
			displayConfig,
			procRef["pk"], 
			displayConfig.TableIDNamePairs[refIdentityRow.member], 
			refIdentityRow.id, 
			attrs[0], 
			displayConfig.TableIDNamePairs[orgID.member], 
			orgID.id, 
			attrToUpdate)

		if err3 != nil {

			log.Println(err3)
		
		} else {

			updatedAttrs[attrToUpdate] = updatedVal
		}
	}

	return updatedAttrs

}

func updateRefOnLeftBasedOnMappingsNotUsingRefIDRow(displayConfig *config.DisplayConfig, 
	procRef map[string]string, orgID *Identity) map[string]string {

	updatedAttrs := make(map[string]string)

	attr := procRef["to_reference"]

	log.Println("attr: ", attr)

	var attrsToUpdate, attrsToUpdateInFETCH []string

	var err1, err2 error

	ignoreREF := false 

	attrsToUpdate, err1 = schema_mappings.GetMappedAttributesFromSchemaMappings(
		displayConfig.AllMappings,
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

	// log.Println("attrsToUpdate:",attrsToUpdate)

	attrsToUpdateInFETCH, err2 = schema_mappings.GetMappedAttributesFromSchemaMappingsByFETCH(
		displayConfig.AllMappings,
		displayConfig.AppIDNamePairs[procRef["app"]], 
		displayConfig.TableIDNamePairs[procRef["from_member"]] +
			"." + procRef["from_reference"], 
		displayConfig.AppConfig.AppName,
		displayConfig.TableIDNamePairs[orgID.member])

	if err2 != nil {

		log.Println(err2)

	}

	// log.Println("attrsToUpdateInFETCH:", attrsToUpdateInFETCH)

	attrsToUpdate = append(attrsToUpdate, attrsToUpdateInFETCH...)

	log.Println("attrsToUpdate:",attrsToUpdate)

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

	return updatedAttrs
}

func updateRefOnRightBasedOnMappingsUsingRefIDRow(displayConfig *config.DisplayConfig, 
	refIdentityRow *Identity, procRef map[string]string, orgID *Identity) map[string]string {

	updatedAttrs := make(map[string]string)

	ignoreREF := true

	attrs, err := schema_mappings.GetMappedAttributesFromSchemaMappings(
		displayConfig.AllMappings,
		displayConfig.AppIDNamePairs[procRef["app"]], 
		displayConfig.TableIDNamePairs[procRef["to_member"]], 
		displayConfig.TableIDNamePairs[procRef["to_member"]] + 
			"." + procRef["to_reference"], 
		displayConfig.AppConfig.AppName,  
		displayConfig.TableIDNamePairs[orgID.member],
		ignoreREF) 
	
	if err != nil {
		
		log.Println(err)

		return nil
	}

	log.Println("attr: ", attrs)

	if len(attrs) != 1 {
		
		log.Println(notOneAttributeFound)
		
		return nil

	}

	var attrsToUpdate, attrsToUpdateInFETCH []string

	var err1, err2 error

	ignoreREF = false

	attrsToUpdate, err1 = schema_mappings.GetMappedAttributesFromSchemaMappings(
		displayConfig.AllMappings,
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

	// log.Println("attrsToUpdate:",attrsToUpdate)

	attrsToUpdateInFETCH, err2 = schema_mappings.GetMappedAttributesFromSchemaMappingsByFETCH(
		displayConfig.AllMappings,
		displayConfig.AppIDNamePairs[procRef["app"]], 
		displayConfig.TableIDNamePairs[procRef["from_member"]] +
			"." + procRef["from_reference"], 
		displayConfig.AppConfig.AppName,
		displayConfig.TableIDNamePairs[refIdentityRow.member])

	if err2 != nil {

		log.Println(err2)

	}

	attrsToUpdate = append(attrsToUpdate, attrsToUpdateInFETCH...)

	log.Println("attrsToUpdate:",attrsToUpdate)
	
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

	return updatedAttrs
	
}

func updateRefOnRightBasedOnMappingsNotUsingRefIDRow(displayConfig *config.DisplayConfig, 
	procRef map[string]string, orgID *Identity) map[string]string {
	
	updatedAttrs := make(map[string]string)

	attr := procRef["to_reference"]

	log.Println("attr: ", attr)

	var attrsToUpdate, attrsToUpdateInFETCH []string

	var err1, err2 error

	ignoreREF := false

	attrsToUpdate, err1 = schema_mappings.GetMappedAttributesFromSchemaMappings(
		displayConfig.AllMappings,
		displayConfig.AppIDNamePairs[procRef["app"]], 
		displayConfig.TableIDNamePairs[procRef["from_member"]], 
		displayConfig.TableIDNamePairs[procRef["from_member"]] + 
			"." + procRef["from_reference"], 
		displayConfig.AppConfig.AppName,
		displayConfig.TableIDNamePairs[procRef["to_member"]], 
		ignoreREF)

	if err1 != nil {
		log.Println(err1)
	}

	// log.Println("attrsToUpdate:",attrsToUpdate)

	attrsToUpdateInFETCH, err2 = schema_mappings.GetMappedAttributesFromSchemaMappingsByFETCH(
		displayConfig.AllMappings,
		displayConfig.AppIDNamePairs[procRef["app"]], 
		displayConfig.TableIDNamePairs[procRef["from_member"]] +
			"." + procRef["from_reference"], 
		displayConfig.AppConfig.AppName,
		displayConfig.TableIDNamePairs[procRef["to_member"]])

	if err2 != nil {

		log.Println(err2)

	}

	attrsToUpdate = append(attrsToUpdate, attrsToUpdateInFETCH...)

	log.Println("attrsToUpdate:",attrsToUpdate)

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

	return updatedAttrs 
	
}


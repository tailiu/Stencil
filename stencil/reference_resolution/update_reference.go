package reference_resolution

import (
	"stencil/config"
	"stencil/schema_mappings"
	"log"
)

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

	log.Println(attrsToUpdate)

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
		displayConfig.TableIDNamePairs[orgID.member],
	)

	if err2 != nil {

		log.Println(err2)

	}

	log.Println(attrsToUpdateInFETCH)

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

func updateRefOnLeftNotUsingRefIDRow(displayConfig *config.DisplayConfig, 
	procRef map[string]string, orgID *Identity) map[string]string {

	updatedAttrs := make(map[string]string)

	attr := procRef["to_reference"]

	log.Println("attr: ", attr)

	ignoreREF := false 

	attrsToUpdate, err := schema_mappings.GetMappedAttributesFromSchemaMappings(
		displayConfig.AllMappings,
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

	return updatedAttrs
}
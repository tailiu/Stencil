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

	// There are cases in which no attribute can be found
	// For example: diaspora posts posts.id mastodon media_attachments
	if len(attrs) != 1 {
		
		log.Println(notOneAttributeFound)
		
		return nil

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

func updateRefOnLeftNotUsingRefIDRow(displayConfig *config.DisplayConfig, 
	procRef map[string]string, orgID *Identity)  map[string]string {

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
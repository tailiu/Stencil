package reference_resolution_v2

import (
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
// When using identity tables, we are not sure about which attributes to update other attributes or 
// which attributes to be updated, so we use fromApp, fromTable, fromAttr, toApp, and toTable to find
// attr and attrToUpdate by virtual of schema mappings. 
// But with the attribute_change table, we are sure about which attributes to update or to be updated.
// It should also be noted that in the display algorithm, 
// we have to check all the attributes to be updated based on mappings.  
func updateRefOnLeftByRefAttrRow(refResolutionConfig *RefResolutionConfig, 
	refAttributeRow *Attribute, procRef map[string]string, orgAttr *Attribute) map[string]string {

	updatedAttr := make(map[string]string)

	attr := refResolutionConfig.attrIDNamePairs[refAttributeRow.attrName]
	attrToUpdate := refResolutionConfig.attrIDNamePairs[orgAttr.attrName]

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

	attr := refResolutionConfig.attrIDNamePairs[procRef["to_attr"]]
	attrToUpdate := refResolutionConfig.attrIDNamePairs[procRef["from_attr"]]

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

	attr := refResolutionConfig.attrIDNamePairs[procRef["to_attr"]]
	attrToUpdate := refResolutionConfig.attrIDNamePairs[procRef["from_attr"]]

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

func updateRefOnRightByRefAttrRow(refResolutionConfig *RefResolutionConfig, 
	refAttributeRow *Attribute, procRef map[string]string, orgAttr *Attribute) map[string]string {

	updatedAttrs := make(map[string]string)

	attr := refResolutionConfig.attrIDNamePairs[orgAttr.attrName]
	attrToUpdate := refResolutionConfig.attrIDNamePairs[refAttributeRow.attrName]

	log.Println("attr:", attr)
	log.Println("attr to be updated:", attrToUpdate)

	updatedVal, err1 := updateReferences(
		refResolutionConfig,
		procRef["pk"],
		refResolutionConfig.tableIDNamePairs[orgAttr.member], 
		orgAttr.val, 
		attr, 
		refResolutionConfig.tableIDNamePairs[refAttributeRow.member], 
		refAttributeRow.val, 
		attrToUpdate,
	)

	if err1 != nil {
		log.Println(err1)
	} else {
		updatedAttrs[refAttributeRow.val + ":" + attrToUpdate] = updatedVal
	}

	return updatedAttrs
}

func updateRefOnRightByRefAttrRow1(refResolutionConfig *RefResolutionConfig, 
	procRef map[string]string, orgAttr *Attribute, refAttrRowVal string) map[string]string {

	attr := refResolutionConfig.attrIDNamePairs[procRef["to_attr"]]
	attrToUpdate := refResolutionConfig.attrIDNamePairs[procRef["from_attr"]]

	updatedAttr := make(map[string]string)

	log.Println("attr:", attr)
	log.Println("attr to be updated:", attrToUpdate)

	updatedVal, err1 := updateReferences(
		refResolutionConfig,
		procRef["pk"],
		refResolutionConfig.tableIDNamePairs[orgAttr.member],
		orgAttr.val,
		attr, 
		refResolutionConfig.tableIDNamePairs[procRef["from_member"]], 
		refAttrRowVal, 
		attrToUpdate,
	)

	if err1 != nil {
		log.Println(err1)
	} else {
		updatedAttr[refAttrRowVal + ":" + attrToUpdate] = updatedVal
	}

	return updatedAttr 
	
}

func updateRefOnRightNotUsingRefAttrRow(refResolutionConfig *RefResolutionConfig, 
	procRef map[string]string, orgAttr *Attribute) map[string]string {

	attr := refResolutionConfig.attrIDNamePairs[procRef["to_attr"]]
	attrToUpdate := refResolutionConfig.attrIDNamePairs[procRef["from_attr"]]

	updatedAttr := make(map[string]string)

	log.Println("attr:", attr)
	log.Println("attr to be updated:", attrToUpdate)

	updatedVal, err1 := updateReferences(
		refResolutionConfig,
		procRef["pk"],
		refResolutionConfig.tableIDNamePairs[orgAttr.member],
		orgAttr.val,
		attr, 
		refResolutionConfig.tableIDNamePairs[procRef["from_member"]], 
		procRef["from_val"], 
		attrToUpdate,
	)

	if err1 != nil {
		log.Println(err1)
	} else {
		updatedAttr[procRef["from_val"] + ":" + attrToUpdate] = updatedVal
	}

	return updatedAttr
}
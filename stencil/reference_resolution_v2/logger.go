package reference_resolution_v2

import (
	"log"
	"fmt"
)

func LogRefRow(refResolutionConfig *RefResolutionConfig, 
	refRow map[string]string, returnLogOnly ...bool) string {
	
	appName := refResolutionConfig.appIDNamePairs[refRow["app"]]

	fromMember := refResolutionConfig.tableIDNamePairs[refRow["from_member"]]
	toMember := refResolutionConfig.tableIDNamePairs[refRow["to_member"]]

	fromAttr := refResolutionConfig.attrIDNamePairs[refRow["from_attr"]]
	toAttr := refResolutionConfig.attrIDNamePairs[refRow["to_attr"]]
	
	output := fmt.Sprint(
		"ref_row - from_member:", fromMember, "|",
		"from_attr:", fromAttr, "|",
		"from_val:", refRow["from_val"], "|",
		"to_member:", toMember, "|",
		"to_attr:", toAttr, "|",
		"to_val:", refRow["to_val"], "|",
		"app:", appName, "|",
		"migration_id:", refRow["migration_id"], "|",
		"pk:", refRow["pk"],
	)

	if len(returnLogOnly) == 0 || !returnLogOnly[0] {
		log.Println(output)
		return ""
	} else {
		return output
	}

}

func logAttrChangeRow(refResolutionConfig *RefResolutionConfig, 
	attrRow map[string]string) {
	
	fromApp := refResolutionConfig.appIDNamePairs[attrRow["from_app"]]
	toApp := refResolutionConfig.appIDNamePairs[attrRow["to_app"]]

	fromMember := refResolutionConfig.tableIDNamePairs[attrRow["from_member"]]
	toMember := refResolutionConfig.tableIDNamePairs[attrRow["to_member"]]

	fromAttr := refResolutionConfig.attrIDNamePairs[attrRow["from_attr"]]
	toAttr := refResolutionConfig.attrIDNamePairs[attrRow["to_attr"]]

	log.Println(
		"id_row - from_app:", fromApp, "|", 
		"from_member:", fromMember, "|",
		"from_attr:", fromAttr, "|",
		"from_val:", attrRow["from_val"], "|",
		"to_app:", toApp, "|",
		"to_member:", toMember, "|",
		"to_attr:", toAttr, "|",
		"to_val:", attrRow["to_val"], "|",
		"migration_id:", attrRow["migration_id"], "|",
		"pk:", attrRow["pk"], 
	)

}

func logRefAttrRow(refResolutionConfig *RefResolutionConfig, 
	attribute *Attribute) {
	
	app := refResolutionConfig.appIDNamePairs[attribute.app]
	member := refResolutionConfig.tableIDNamePairs[attribute.member]
	attrName := refResolutionConfig.attrIDNamePairs[attribute.attrName]

	log.Println(
		"refIdentityRow - app:", app, "|",
		"member:", member, "|",
		"attrName:", attrName, "|",
		"val:", attribute.val,
	)
}
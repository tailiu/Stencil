package reference_resolution_v2

import (
	"log"
	"fmt"
)

func (rr *RefResolution) LogRefRow(refRow map[string]string, returnLogOnly ...bool) string {
	
	appName := rr.appIDNamePairs[refRow["app"]]

	fromMember := rr.tableIDNamePairs[refRow["from_member"]]
	toMember := rr.tableIDNamePairs[refRow["to_member"]]

	fromAttr := rr.attrIDNamePairs[refRow["from_attr"]]
	toAttr := rr.attrIDNamePairs[refRow["to_attr"]]
	
	output := fmt.Sprint(
		"ref_row - from_member:", fromMember, " | ",
		"from_attr:", fromAttr, " | ",
		"from_val:", refRow["from_val"], " | ",
		"to_member:", toMember, " | ",
		"to_attr:", toAttr, " | ",
		"to_val:", refRow["to_val"], " | ",
		"app:", appName, " | ",
		"migration_id:", refRow["migration_id"], " | ",
		"pk:", refRow["pk"],
	)

	if len(returnLogOnly) == 0 || !returnLogOnly[0] {
		log.Println(output)
		return ""
	} else {
		return output
	}

}

func (rr *RefResolution) logAttrChangeRow(attrRow map[string]string) {
	
	fromApp := rr.appIDNamePairs[attrRow["from_app"]]
	toApp := rr.appIDNamePairs[attrRow["to_app"]]

	fromMember := rr.tableIDNamePairs[attrRow["from_member"]]
	toMember := rr.tableIDNamePairs[attrRow["to_member"]]

	fromAttr := rr.attrIDNamePairs[attrRow["from_attr"]]
	toAttr := rr.attrIDNamePairs[attrRow["to_attr"]]

	log.Println(
		"attr_row - from_app:", fromApp, "|", 
		"from_member:", fromMember, "|",
		"from_attr:", fromAttr, "|",
		"from_val:", attrRow["from_val"], "|",
		"from_id:", attrRow["from_id"], "|",
		"to_app:", toApp, "|",
		"to_member:", toMember, "|",
		"to_attr:", toAttr, "|",
		"to_val:", attrRow["to_val"], "|",
		"to_id:", attrRow["to_id"], "|",
		"migration_id:", attrRow["migration_id"], "|",
		"pk:", attrRow["pk"], 
	)

}

func (rr *RefResolution) logRefAttrRow(attribute *Attribute) {
	
	app := rr.appIDNamePairs[attribute.app]
	member := rr.tableIDNamePairs[attribute.member]
	attrName := rr.attrIDNamePairs[attribute.attrName]

	log.Println(
		"refAttrRow - app:", app, "|",
		"member:", member, "|",
		"attrName:", attrName, "|",
		"val:", attribute.val, "|",
		"id:", attribute.id, 
	)
}
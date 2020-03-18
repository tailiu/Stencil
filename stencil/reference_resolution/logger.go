package reference_resolution

import (
	"log"
	"fmt"
)

func LogRefRow(refResolutionConfig *RefResolutionConfig, 
	refRow map[string]string, returnLogOnly ...bool) string {

	fromMember := refResolutionConfig.tableIDNamePairs[refRow["from_member"]]
	toMember := refResolutionConfig.tableIDNamePairs[refRow["to_member"]]

	appName := refResolutionConfig.appIDNamePairs[refRow["app"]]
	
	output := fmt.Sprint(
		"ref_row - from_member:", fromMember, "|",
		"from_reference:", refRow["from_reference"], "|",
		"from_id:", refRow["from_id"], "|",
		"to_member:", toMember, "|",
		"to_reference:", refRow["to_reference"], "|",
		"to_id:", refRow["to_id"], "|",
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

func logIDRow(refResolutionConfig *RefResolutionConfig, 
	IDRow map[string]string) {
	
	fromApp := refResolutionConfig.appIDNamePairs[IDRow["from_app"]]
	toApp := refResolutionConfig.appIDNamePairs[IDRow["to_app"]]

	fromMember := refResolutionConfig.tableIDNamePairs[IDRow["from_member"]]
	toMember := refResolutionConfig.tableIDNamePairs[IDRow["to_member"]]

	log.Println("id_row - from_app:", fromApp, "|", 
		"from_member:", fromMember, "|",
		"from_id:", IDRow["from_id"], "|",
		"to_app:", toApp, "|",
		"to_member:", toMember, "|",
		"to_id:", IDRow["to_id"], "|",
		"migration_id:", IDRow["migration_id"], "|",
		"pk:", IDRow["pk"], 
	)

}

func logRefIDRow(refResolutionConfig *RefResolutionConfig, 
	ID *Identity) {
	
	app := refResolutionConfig.appIDNamePairs[ID.app]

	member := refResolutionConfig.tableIDNamePairs[ID.member]

	log.Println("refIdentityRow - app:", app, "|",
		"member:", member, "|",
		"id:", ID.id,
	)
}
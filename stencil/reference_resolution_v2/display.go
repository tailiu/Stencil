package reference_resolution_v2

import (
	"log"
	"fmt"
	"stencil/db"
	"stencil/schema_mappings"
)

func NeedToResolveReference(refResolutionConfig *RefResolutionConfig, 
	toTable, toAttr string) bool {
	
	for _, mapping := range refResolutionConfig.mappingsFromOtherAppsToDst {

		if exists, err := schema_mappings.REFExists(
			mapping, toTable, toAttr); err != nil {
	
			// This can happen when there is no mapping
			// For example: 
			// When migrating from Diaspora to Mastodon:
			// there is no mapping to stream_entries.activity_id.
			log.Println(err)
	
		} else {
			if exists {
				return true
			}
		}

	}

	return false
	
}

func NeedToResolveReferenceOnlyBasedOnSrc(refResolutionConfig *RefResolutionConfig, 
	toTable, toAttr string) bool {
	
	if exists, err := schema_mappings.REFExists(
		refResolutionConfig.mappingsFromSrcToDst, toTable, toAttr); err != nil {

		// This can happen when there is no mapping
		// For example: 
		// When migrating from Diaspora to Mastodon:
		// there is no mapping to stream_entries.activity_id.
		log.Println(err)

		return false

	} else {
		if exists {
			return true
		} else {
			return false
		}
	}
}

func ReferenceResolved(refResolutionConfig *RefResolutionConfig, 
	member, attr, id string) string {

	query := fmt.Sprintf(
		`select attr_val_after_changes from resolved_references where app = %s 
		and member = %s and attr = '%s' and id = %s`,
		refResolutionConfig.appID, 
		member, attr, id,
	)
	
	log.Println(query)
	
	data, err := db.DataCall1(refResolutionConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)
	if len(data) == 0 {
		return ""
	} else {
		return fmt.Sprint(data["attr_val_after_changes"])
	}
}

func GetUpdatedAttributes(refResolutionConfig *RefResolutionConfig, 
	ID *Identity) map[string]string {
	
	updatedAttrs := make(map[string]string) 

	query := fmt.Sprintf(`select reference, value from resolved_references where app = %s 
		and member = %s and id = %s ORDER BY pk`,
		refResolutionConfig.appID, ID.member, ID.id)
	
	log.Println(query)

	data, err := db.DataCall(refResolutionConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// Even though there could be some duplicate keys here since 
	// we don't consider migration_id here,
	// the ORDER BY in the query will give us the latest resolved values
	for _, data1 := range data {

		updatedAttrs[fmt.Sprint(data1["reference"])] = 
			fmt.Sprint(data1["value"])

	}

	return updatedAttrs
	
}

func getUpdateIDInDisplayFlagsQuery(refResolutionConfig *RefResolutionConfig, 
	table, IDToBeUpdated, id string) string {
	
	query := fmt.Sprintf(
		`UPDATE display_flags SET id = %s, updated_at = now() 
		WHERE app_id = %s and table_id = %s 
		and id = %s and migration_id = %d;`,
		id, refResolutionConfig.appID, 
		refResolutionConfig.appTableNameIDPairs[table],
		IDToBeUpdated, refResolutionConfig.migrationID,
	)

	return query

}
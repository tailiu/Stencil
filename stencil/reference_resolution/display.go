package reference_resolution

import (
	"log"
	"fmt"
	"stencil/db"
	"stencil/schema_mappings"
)

func NeedToResolveReference(refResolutionConfig *RefResolutionConfig, 
	toTable, toAttr string) bool {
	
	if exists, err := schema_mappings.REFExists(
		refResolutionConfig.mappingsFromSrcToDst, toTable, toAttr);
		
		err != nil {

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
	member, reference, id string) string {

	query := fmt.Sprintf(`select value from resolved_references where app = %s 
		and member = %s and reference = '%s' and id = %s`,
		refResolutionConfig.appID, member, reference, id)

	data, err := db.DataCall1(refResolutionConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)
	if len(data) == 0 {

		return ""

	} else {

		return fmt.Sprint(data["value"])

	}
}
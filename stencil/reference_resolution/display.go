package reference_resolution

import (
	"stencil/config"
	"stencil/schema_mappings"
	"log"
)

func NeedToResolveReference(displayConfig *config.DisplayConfig, toTable, toAttr string) bool {
	
	if exists, err := schema_mappings.REFExists(displayConfig, toTable, toAttr); err != nil {

		log.Fatal(err)

		return false

	} else {
		if exists {

			return true

		} else {

			return false
		}
	}
}

func ReferenceResolved(displayConfig *config.DisplayConfig, member, reference string, id int) {

	query := fmt.Sprintf("select value from resolved_references where app = %s and member = %s and reference = %s and id = %s",
		displayConfig.AppConfig.AppID, member, reference, id)



}
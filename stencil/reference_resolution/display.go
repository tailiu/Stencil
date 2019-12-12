package reference_resolution

import (
	"stencil/db"
	"stencil/config"
	"stencil/schema_mappings"
	"log"
	"fmt"
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

func ReferenceResolved(displayConfig *config.DisplayConfig, member, reference, id string) string {

	query := fmt.Sprintf(`select value from resolved_references where app = %s 
		and member = %s and reference = %s and id = %s`,
		displayConfig.AppConfig.AppID, member, reference, id)

	data, err := db.DataCall1(displayConfig.StencilDBConn, query)
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
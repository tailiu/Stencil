package reference_resolution

import (
	"stencil/config"
	"stencil/schema_mappings"
)

func NeedToResolveRef(displayConfig *config.DisplayConfig, toTable, toAttr string) bool {
	
	if schema_mappings.REFExists(displayConfig, toTable, toAttr) {

		return true

	} else {

		return false
	}
}

func RefResolved(displayConfig *config.DisplayConfig) {

	// query := fmt.Sprintf("select pk, table_name from app_tables;")

}
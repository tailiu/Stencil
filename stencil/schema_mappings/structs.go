package schema_mappings

const FILEPATH = "./config/app_settings/PSM_mappings.json"

type conditionsNotConsidered struct {
	fromApp			string
	toApp			string
	fromTables		[]string
	toTable			string
	condName		string
	condVal			string
}
package schema_mappings

const INPUTFILEPATH = "./config/app_settings/PSM_test.json"

const OUTPUTFILEPATH = "./config/app_settings/PSM_mappings.json"

type conditionsNotConsidered struct {
	fromApp			string
	toApp			string
	fromTables		[]string
	toTable			string
	condName		string
	condVal			string
}
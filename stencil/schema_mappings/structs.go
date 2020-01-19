package schema_mappings

const INPUTFILEPATH = "./schema_mappings/PSM_tests/PSM_test.json"

const OUTPUTFILEPATH = "./schema_mappings/PSM_tests/PSM_mappings.json"

type conditionsNotConsidered struct {
	fromApp			string
	toApp			string
	fromTables		[]string
	toTable			string
	condName		string
	condVal			string
}
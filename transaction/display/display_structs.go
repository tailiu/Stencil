package display

// The Key should be the primay key of the Table
type HintStruct struct {
	Table string		`json:"Table"`
	Key string			`json:"Key"`
	Value string		`json:"Value"`
	ValueType string	`json:"ValueType"`
}
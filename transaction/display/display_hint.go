package display

import (
	// "fmt"
	// "errors"
	// "log"
	"database/sql"
	"transaction/db"
)

// The Key should be the primay key of the Table
type HintStruct struct {
	Table string		`json:"Table"`
	Key string			`json:"Key"`
	Value string		`json:"Value"`
	ValueType string	`json:"ValueType"`
}

// NOTE: We assume that primary key is only one integer value!!!
func TransformRowToHint(dbConn *sql.DB, row map[string]string, table string) (HintStruct, error) {
	hintData := HintStruct{}
	pk, err := db.GetPrimaryKeyOfTable(dbConn, table)
	if err != nil {
		return hintData, err
	} else {
		hintData.Table = table
		hintData.Key = pk
		hintData.Value = row[pk]
		hintData.ValueType = "int"
	}
	return hintData, nil
}
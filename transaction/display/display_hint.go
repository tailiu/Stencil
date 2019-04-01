package display

import (
	// "fmt"
	// "errors"
	// "log"
	"database/sql"
	"transaction/db"
)

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
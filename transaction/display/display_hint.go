package display

import (
	"log"
	"database/sql"
	"transaction/db"
	"strconv"
)

// The Key should be the primay key of the Table
// type HintStruct struct {
// 	Table string		
// 	id string			
// 	Value string		
// 	ValueType string	
// }

// The Key should be the primay key of the Table
type HintStruct struct {
	Table string		
	KeyVal map[string]int
}

// NOTE: We assume that primary key is only one integer value!!!
func TransformRowToHint(dbConn *sql.DB, row map[string]string, table string) (HintStruct, error) {
	hint := HintStruct{}
	pk, err := db.GetPrimaryKeyOfTable(dbConn, table)
	if err != nil {
		return hint, err
	} else {
		intPK, err1 := strconv.Atoi(row[pk])
		if err1 != nil {
			log.Fatal(err1)
		}
		keyVal := map[string]int {
			pk:	intPK,
		}
		hint.Table = table
		hint.KeyVal = keyVal
	}
	return hint, nil
}
package migrate_v2

import (
	"fmt"
	"log"
	"stencil/db"
	"strings"
)

func (mmd MappedMemberData) GetQueryArgs() (string, string, []interface{}) {
	var colStr, phStr string
	var valList []interface{}
	counter := 1

	for mappedAttr, mmv := range mmd.Data {
		colStr += fmt.Sprintf("\"%s\",", mappedAttr)
		phStr += fmt.Sprintf("$%v,", counter)
		valList = append(valList, mmv.Value)
		counter++
	}

	colStr = strings.Trim(colStr, ",")
	phStr = strings.Trim(phStr, ",")

	return colStr, phStr, valList
}

func (mmd MappedMemberData) ValidateMappedData() bool {
	for mappedAttr, mmv := range mmd.Data {
		if mmv.IsInput || mmv.IsExpression || mappedAttr == "id" {
			continue
		}
		if mmv.Value != nil {
			return true
		}
	}
	return false
}

func (mmd MappedMemberData) SetMember(table string) {
	mmd.ToMember = table
	if mmd.DBConn == nil {
		log.Fatal("@mmd.SetMember: DBConn not set!")
	} else {
		if tableID, err := db.TableID(mmd.DBConn, mmd.ToMember, mmd.AppID); err == nil {
			mmd.ToMemberID = tableID
		} else {
			fmt.Println(mmd.ToMember, mmd.AppID)
			log.Fatal("@SetMember: ", err)
		}
	}
}

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

func (mmd *MappedMemberData) SetMember(table string) {
	mmd.ToMember = table
	if mmd.DBConn == nil {
		log.Fatal("@mmd.SetMember: DBConn not set!")
	} else {
		if tableID, err := db.TableID(mmd.DBConn, mmd.ToMember, mmd.ToAppID); err == nil {
			mmd.ToMemberID = tableID
		} else {
			fmt.Println(mmd.ToMember, mmd.ToAppID)
			log.Fatal("@SetMember: ", err)
		}
	}
}

func (mmd MappedMemberData) SrcTables() []Member {

	var srcTables []Member

	added := make(map[string]bool)

	for _, mmv := range mmd.Data {
		if mmv.IsExpression || mmv.IsInput {
			continue
		}
		if _, ok := added[mmv.FromMemberID]; !ok {
			srcTables = append(srcTables, Member{ID: mmv.FromMemberID, Name: mmv.FromMember})
		}
		added[mmv.FromMemberID] = true
	}

	return srcTables
}

func (mmd MappedMemberData) ToCols() []string {
	var toCols []string
	for toCol := range mmd.Data {
		toCols = append(toCols, toCol)
	}
	return toCols
}

func (mmd MappedMemberData) FromCols(table string) []string {
	var fromCols []string
	for _, mmv := range mmd.Data {
		if strings.EqualFold(mmv.FromMember, table) {
			fromCols = append(fromCols, mmv.FromAttr)
		}
	}
	return fromCols
}

func (mmd MappedMemberData) GetDataMap() DataMap {

	data := make(DataMap)

	for col, mmv := range mmd.Data {
		data[col] = mmv.Value
	}

	return data
}

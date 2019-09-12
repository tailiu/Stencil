package qr

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"stencil/db"
	"github.com/xwb1989/sqlparser"
	"log"
)

func CreateQI(table string, cols []string, vals []interface{}, qtype string) *QI {
	qi := new(QI)
	qi.TableName = table
	qi.Columns = cols
	qi.Values = vals
	qi.Type = qtype
	return qi
}

func (self *QI) valueOfColumn(col string) (interface{}, error) {

	for i, c := range self.Columns {
		if strings.EqualFold(col, c) {
			return self.Values[i], nil
		}
	}
	return "", errors.New("No column: " + col)
}

func (self *QI) Print() {

	fmt.Println(self)

	// for i, c := range self.Columns {
	// 	fmt.Println("i:", i, "Col:", c, " || Val:", self.Values[i])
	// }
}

// func (self *QI) ResolveInsert(QR *QR, rowID int32) []*QI {

// 	var PQIs []*QI
// 	newRowCols := []string{"rowid", "app_id", "table"}
// 	newRowVals := []interface{}{rowID, QR.AppID, self.TableName}
// 	newRowQI := CreateQI("row_desc", newRowCols, newRowVals, QTInsert)
// 	PQIs = append(PQIs, newRowQI)
// 	phyMap := QR.GetPhyMappingForLogicalTable(self.TableName)

// 	for pt, mapping := range phyMap {
// 		isValid := false
// 		pqiCols := []string{"pk"}
// 		pqiVals := []interface{}{rowID}
// 		for _, colmap := range mapping {
// 			if val, err := self.valueOfColumn(colmap[1]); err == nil {
// 				isValid = true
// 				pqiCols = append(pqiCols, colmap[0])
// 				pqiVals = append(pqiVals, val)
// 			}
// 		}
// 		if isValid {
// 			pqi := CreateQI(pt, pqiCols, pqiVals, QTInsert)
// 			PQIs = append(PQIs, pqi)
// 		}

// 	}
// 	return PQIs
// }



func (self *QR) ResolveInsert(qi *QI, rowID int32) []*QI {

	PQIs := self.ResolveInsertWithoutRowDesc(qi, rowID)
	if len(PQIs) > 0 {
		if tableID, err := db.TableID(self.StencilDB, qi.TableName, self.AppID); err == nil {
			newRowCols := []string{"group_id", "row_id", "app_id", "table_id"}
			newRowVals := []interface{}{rowID, rowID, self.AppID, tableID}
			newRowQI := CreateQI("migration_table", newRowCols, newRowVals, QTInsert)
			PQIs = append(PQIs, newRowQI)
		}else{
			fmt.Println("Cant get tableID for ", qi.TableName)
			log.Fatal(err)
		}
		
	}
	return PQIs
}

func (self *QR) ResolveInsertWithoutRowDesc(qi *QI, rowID int32) []*QI {

	var PQIs []*QI
	
	phyMap := self.GetPhyMappingForLogicalTable(qi.TableName)

	for pt, mapping := range phyMap {
		isValid := false
		pqiCols := []string{"pk"}
		pqiVals := []interface{}{rowID}
		for _, colmap := range mapping {
			if val, err := qi.valueOfColumn(colmap[1]); err == nil {
				isValid = true
				pqiCols = append(pqiCols, colmap[0])
				pqiVals = append(pqiVals, val)
			}
		}
		if isValid {
			pqi := CreateQI(pt, pqiCols, pqiVals, QTInsert)
			PQIs = append(PQIs, pqi)
		}

	}
	return PQIs
}



func (self *QI) GenSQL() (string, []interface{}) {

	switch self.Type {
	case QTSelect:
		fmt.Println("!!! Unimplemented type: Select")
	case QTUpdate:
		fmt.Println("!!! Unimplemented type: Update")
	case QTDelete:
		fmt.Println("!!! Unimplemented type: Delete")
	case QTInsert:
		var cols, vals []string
		for i, col := range self.Columns {
			cols = append(cols, col)
			vals = append(vals, fmt.Sprintf("$%d", i+1))
		}
		q := fmt.Sprintf("INSERT INTO \"%s\" (\"%s\") VALUES (%s) ON CONFLICT DO NOTHING;", self.TableName, strings.Join(cols, "\",\""), strings.Join(vals, ","))
		return q, self.Values
	}
	fmt.Println("!!! Unable to identify query type.", self.Type)
	return "", self.Values
}

func getInsertQueryIngs(sql string) *QI {

	qi := new(QI)

	tokens := sqlparser.NewStringTokenizer(sql)

	vswitch := 0

	var cols []string
	var vals []interface{}

	for {
		ttype, tval := tokens.Scan()

		if ttype == 0 {
			break
		}

		if len(string(tval)) <= 0 {
			continue
		}

		if string(tval) == "values" {
			vswitch = 4
			continue
		}

		switch vswitch {
		case 0:
			vswitch++
		case 1:
			vswitch++
		case 2:
			qi.TableName = string(tval)
			vswitch++
		case 3:
			cols = append(cols, string(tval))
		case 4:
			vals = append(vals, string(tval))
		}

	}

	qi.Columns = cols
	qi.Values = vals

	return qi
}

func getDeleteQueryIngs(sql string) *QI {

	qi := new(QI)
	re := regexp.MustCompile(`(?i)(delete from | where )`)
	phrases := deleteEmpty(re.Split(sql, -1))
	if len(phrases) > 0 {
		qi.TableName = phrases[0]
		if len(phrases) > 1 {
			qi.Conditions = phrases[1]
		}
	}
	return qi
}

func getUpdateQueryIngs(sql string) *QI {

	qi := new(QI)

	re := regexp.MustCompile(`(?i)(update | where | set )`)
	phrases := deleteEmpty(re.Split(sql, -1))

	if len(phrases) > 1 {

		qi.TableName = phrases[0]

		updates := strings.Split(phrases[1], ",")

		var cols []string
		var vals []interface{}

		for _, update := range updates {
			items := strings.Split(update, "=")
			cols = append(cols, strings.Trim(items[0], " ,"))
			vals = append(vals, strings.Trim(items[1], " ,'"))
		}

		qi.Columns = cols
		qi.Values = vals

		if len(phrases) > 2 {

			qi.Conditions = phrases[2]
		}
	}

	return qi
}

func getSelectQueryIngs(sql string) *QI {

	qi := new(QI)
	re := regexp.MustCompile(`(?i)(select | from | where )`)
	phrases := deleteEmpty(re.Split(sql, -1))

	if len(phrases) > 1 {
		qi.TableName = strings.Trim(phrases[1], " ,")
		qi.Columns = strings.Split(phrases[0], ",")
		qi.ColumnsWithTable = make(map[string][]string)
		for _, col := range qi.Columns {
			if strings.Contains(col, ".") {
				coltab := strings.Split(col, ".")
				table := strings.Trim(coltab[0], " ,.")
				column := strings.Trim(coltab[1], " ,.")

				// qi.Columns[i] = column
				qi.ColumnsWithTable[table] = append(qi.ColumnsWithTable[table], column)
			}
		}

		if len(phrases) > 2 {
			qi.Conditions = phrases[2]
		}
	}
	return qi
}
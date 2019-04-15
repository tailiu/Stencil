/*
 * Query Resolver
 */

package qr

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"transaction/db"

	_ "github.com/lib/pq" // postgres driver
	escape "github.com/tj/go-pg-escape"
)

func (self QR) NewRowId() (string, error) {
	sql := "SELECT unique_rowid() AS rowid"
	res := db.DataCall1(self.StencilDB, sql)
	if val, ok := res["rowid"]; ok {
		return val.(string), nil
	}
	return "-1", errors.New("can't get new rowid")
}

func (self QR) GetPhyMappingForLogicalTable(ltable string) map[string][][]string {

	var phyMap = make(map[string][][]string)

	for _, mapping := range append(self.BaseMappings, self.SuppMappings...) {
		if strings.EqualFold(ltable, mapping["logical_table"].(string)) {
			ptab := mapping["physical_table"].(string)
			pcol := mapping["physical_column"].(string)
			lcol := mapping["logical_column"].(string)
			var pair []string
			pair = append(pair, pcol, lcol)
			if _, ok := phyMap[ptab]; ok {
				phyMap[ptab] = append(phyMap[ptab], pair)
			} else {
				phyMap[ptab] = [][]string{pair}
			}
		}
	}
	return phyMap
}

func (self QR) GetBaseMappingForLogicalTable(ltable string) map[string][][]string {

	var phyMap = make(map[string][][]string)

	for _, mapping := range self.BaseMappings {
		if strings.EqualFold(ltable, mapping["logical_table"].(string)) {
			ptab := mapping["physical_table"].(string)
			pcol := mapping["physical_column"].(string)
			lcol := mapping["logical_column"].(string)
			var pair []string
			pair = append(pair, pcol, lcol)
			if _, ok := phyMap[ptab]; ok {
				phyMap[ptab] = append(phyMap[ptab], pair)
			} else {
				phyMap[ptab] = [][]string{pair}
			}
			// fmt.Println(i, pair, mapping)
		}
	}

	return phyMap
}

// func (self QR) ResolveUpdate(sql string) []string {
// 	var PQs []string

// 	qi := getUpdateQueryIngs(sql)
// 	rowIDs := self.GetAffectedRowIDs(qi.TableName, qi.Conditions)

// 	phyMap := self.GetPhyMappingForLogicalTable(qi.TableName)

// 	for pt, mapping := range phyMap {
// 		updates := ""
// 		for _, colmap := range mapping {
// 			if val, err := qi.valueOfColumn(colmap[1]); err == nil {
// 				updates += fmt.Sprintf("%s = %s, ", colmap[0], escape.Literal(val))
// 			}
// 		}
// 		if updates != "" {
// 			updates := strings.Trim(updates, ", ")
// 			pq := fmt.Sprintf("UPDATE %s SET %s ", pt, updates)
// 			if len(rowIDs) > 0 {
// 				pq += fmt.Sprintf("WHERE %s_row_id IN (%s)", pt[0:4], strings.Join(rowIDs[:], ","))
// 			}
// 			PQs = append(PQs, pq)
// 		}
// 	}
// 	return PQs
// }

// func (self QR) ResolveDelete(sql string) []string {
// 	var PQs []string

// 	qi := getDeleteQueryIngs(sql)
// 	fmt.Println(qi)
// 	rowIDs := self.GetAffectedRowIDs(qi.TableName, qi.Conditions)
// 	phyMap := self.GetPhyMappingForLogicalTable(qi.TableName)

// 	for pt, _ := range phyMap {
// 		pq := fmt.Sprintf("DELETE FROM %s ", pt)
// 		if len(rowIDs) > 0 {
// 			pq += fmt.Sprintf("WHERE %s_row_id IN (%s)", pt[0:4], strings.Join(rowIDs[:], ","))
// 		}
// 		PQs = append(PQs, pq)
// 	}
// 	return PQs
// }

// func (self QR) GetAffectedRowIDs(table, conds string) []string {

// 	rowIDs := []string{"1", "2"}
// 	return rowIDs

// 	// var rowIDs []string

// 	sql := fmt.Sprintf("SELECT * from %s WHERE %s", table, conds)

// 	// needs to be changed, indicator args replaced for migration
// 	pqs := self.ResolveSelect(sql, true)

// 	if len(pqs) > 0 {
// 		for _, rowMap := range db.DataCall(self.StencilDB, sql) {
// 			for _, val := range rowMap {
// 				rowIDs = append(rowIDs, val.(string))
// 			}
// 		}
// 	}

// 	return rowIDs
// }

func (self QR) getPhyTabCol(ltabcol string) (string, string) {

	tab := strings.Trim(strings.Split(ltabcol, ".")[0], " ")
	col := strings.Trim(strings.Split(ltabcol, ".")[1], " ")

	phyMap := self.GetPhyMappingForLogicalTable(tab)

	for pt, mapping := range phyMap {
		for _, colmap := range mapping {
			if colmap[1] == col {
				return pt, colmap[0]
			}
		}
	}

	return "", ""
}

func (self QR) PhyUpdateAppIDByRowID(new_app_id, ltab string, rowIDs []string) []string {

	var PQs []string

	if len(rowIDs) <= 0 {
		log.Println("Warning: NO ROWIDS!")
	} else {
		phyMap := self.GetBaseMappingForLogicalTable(strings.ToLower(ltab))
		for pt := range phyMap {
			pq := fmt.Sprintf("UPDATE %s SET app_id = %s WHERE app_id = '%s' AND %s_row_id IN (%s);", pt, escape.Literal(new_app_id), self.AppID, pt[0:4], strings.Join(rowIDs[:], ","))
			PQs = append(PQs, pq)
		}
	}
	return PQs
}

func (self QR) ResolveSelectWithoutJoins(sql string, qi *QI, args ...interface{}) []string {

	var PQs []string
	var phyMap map[string][][]string

	if len(args) <= 0 {
		phyMap = self.GetPhyMappingForLogicalTable(qi.TableName)
	} else {
		phyMap = self.GetBaseMappingForLogicalTable(qi.TableName)
	}
	cols := ""
	joins := ""
	conds := qi.Conditions
	prev := ""
	for pt, mapping := range phyMap {
		joined := false
		for _, colmap := range mapping {
			if contains(qi.Columns, colmap[1]) || contains(qi.Columns, "*") {
				if cols == "" {
					cols = fmt.Sprintf("row_desc.rowid AS base_row_id, ")
				}
				if !joined {
					if prev == "" {
						joins += fmt.Sprintf(" %s ", pt)
					} else {
						joins += fmt.Sprintf(" JOIN %s ON %s.%s = %s.%s ", pt, prev, prev[0:4]+"_row_id", pt, pt[0:4]+"_row_id")
					}
					prev = pt
					joined = true
				}
				col := fmt.Sprintf("%s.%s", pt, colmap[0])
				cols += col + ", "
				if nconds := strings.Replace(conds, qi.TableName+"."+colmap[1], col, -1); conds == nconds {
					conds = strings.Replace(conds, colmap[1], col, -1)
				} else {
					conds = nconds
				}
				// conds = strings.Replace(conds, colmap[1], col, -1)
			} else {
			}
		}
	}
	joins += fmt.Sprintf(" JOIN row_desc ON %s.%s = row_desc.rowid", prev, prev[0:4]+"_row_id")
	pq := fmt.Sprintf("SELECT %s FROM %s WHERE row_desc.app_id = '%s'", strings.Trim(cols, ", "), strings.Trim(joins, ", "), self.AppID)
	if len(conds) > 0 {
		pq += fmt.Sprintf(" AND (%s)", strings.Trim(conds, ", "))
	}
	PQs = append(PQs, pq)
	return PQs
}

func (self QR) ResolveSelectWithJoins(sql string, qi *QI, args ...interface{}) []string {

	// parameter args used as indicator for query resolution during migration

	var PQs []string

	re := regexp.MustCompile(`(?i)(join)`)
	phrases := deleteEmpty(re.Split(qi.TableName, -1))
	phyMaps := make(map[string]map[string][][]string)
	bigjoin := ""
	for _, phrase := range phrases {
		phrase = strings.Trim(phrase, " ")
		re := regexp.MustCompile(`(?i)( on )`)
		tabWOnCond := deleteEmpty(re.Split(phrase, -1))
		if len(args) <= 0 {
			phyMaps[tabWOnCond[0]] = self.GetPhyMappingForLogicalTable(tabWOnCond[0])
		} else {
			phyMaps[tabWOnCond[0]] = self.GetBaseMappingForLogicalTable(tabWOnCond[0])
		}
		// pconds := ""
		joins := ""
		prev := ""
		pcols := ""
		for pt, mapping := range phyMaps[tabWOnCond[0]] {
			// pcols += pt + ".*, "
			if pcols == "" {
				pcols = fmt.Sprintf("row_desc.rowid AS base_row_id, ")
			}
			for _, colmap := range mapping {
				pcols += fmt.Sprintf("%s.%s as %s, ", pt, colmap[0], colmap[1])
			}
			if joins == "" {
				joins = pt
			} else {
				joins += fmt.Sprintf(" JOIN %s ON %s.%s = %s.%s ", pt, prev, prev[0:4]+"_row_id", pt, pt[0:4]+"_row_id")
			}
			prev = pt
		}
		joins += fmt.Sprintf(" JOIN row_desc ON %s.%s = row_desc.rowid", prev, prev[0:4]+"_row_id")
		// ptable := fmt.Sprintf(" (SELECT %s FROM %s WHERE app_id = '%s') %s ", strings.Trim(pcols, " ,"), strings.Trim(joins, " ,"), self.AppID, tabWOnCond[0])
		ptable := fmt.Sprintf(" (SELECT %s FROM %s WHERE row_desc.app_id = '%s') %s ", strings.Trim(pcols, " ,"), strings.Trim(joins, " ,"), self.AppID, tabWOnCond[0])
		// ptable := fmt.Sprintf(" (SELECT %s FROM %s) %s ", strings.Trim(pcols, " ,"), strings.Trim(joins, " ,"), tabWOnCond[0])

		if len(tabWOnCond) > 1 {
			bigjoin += fmt.Sprintf(" JOIN %s ON %s ", ptable, tabWOnCond[1])
		} else {
			bigjoin += ptable
		}
	}
	var bigsql string
	if len(qi.Conditions) > 0 {
		bigsql = fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(qi.Columns, ","), bigjoin, qi.Conditions)
	} else {
		bigsql = fmt.Sprintf("SELECT %s FROM %s", strings.Join(qi.Columns, ","), bigjoin)
	}
	PQs = append(PQs, bigsql)
	return PQs
}

func (self QR) ResolveSelect(sql string, args ...interface{}) []string {

	qi := getSelectQueryIngs(sql)
	if strings.Contains(qi.TableName, " join ") {
		return self.ResolveSelectWithJoins(sql, qi, args)
	} else {
		return self.ResolveSelectWithoutJoins(sql, qi, args)
	}
}

func (self QR) ResolveInsert(qi *QI) []*QI {

	var PQIs []*QI
	if rowID, err := self.NewRowId(); err == nil {
		// newRowSQL := fmt.Sprintf("INSERT INTO row_desc (row_id, app_id, table_name) VALUES ('%s', '%s', '%s')", rowID, self.AppID, qi.TableName)
		newRowCols := []string{"rowid", "app_id", "table_name"}
		newRowVals := []interface{}{rowID, self.AppID, qi.TableName}
		newRowQI := CreateQI("row_desc", newRowCols, newRowVals, QTInsert)
		PQIs = append(PQIs, newRowQI)
		phyMap := self.GetPhyMappingForLogicalTable(qi.TableName)

		for pt, mapping := range phyMap {
			isValid := false
			// pqCols := fmt.Sprintf("INSERT INTO %s ( rowid, app_id, ", pt, pt[0:4])
			pqiCols := []string{"rowid", "app_id"}
			// pqVals := fmt.Sprintf("VALUES ( '%s','%s',", rowID, self.AppID)
			pqiVals := []interface{}{rowID, self.AppID}
			for _, colmap := range mapping {
				if val, err := qi.valueOfColumn(colmap[1]); err == nil {
					isValid = true
					pqiCols = append(pqiCols, colmap[0])
					pqiVals = append(pqiVals, val)
					// pqVals += fmt.Sprintf("E'%s',", val)
				}
			}
			if isValid {
				// pqi := strings.Trim(pqCols, ", ") + ") " + strings.Trim(pqVals, ", ") + ");"
				pqi := CreateQI(pt, pqiCols, pqiVals, QTInsert)
				PQIs = append(PQIs, pqi)
			}

		}
	}
	return PQIs
}

// func (self QR) Resolve(sql string, args ...interface{}) []string {

// 	var PQs []string

// 	sql = strings.ToLower(sql)

// qi := getInsertQueryIngs(sql)

// 	if stmt, err := sqlparser.Parse(sql); err != nil {
// 		fmt.Println("Error parsing:", err)
// 	} else {
// 		switch stmt := stmt.(type) {
// 		case *sqlparser.Select:
// 			PQs = self.ResolveSelect(sql, args)
// 		case *sqlparser.Update:
// 			PQs = self.ResolveUpdate(sql)
// 		case *sqlparser.Delete:
// 			PQs = self.ResolveDelete(sql)
// 		case *sqlparser.Insert:
// 			PQs = self.ResolveInsert(sql)
// 		default:
// 			fmt.Println("!!! Unable to identify query type.", stmt)
// 		}
// 	}
// 	return PQs
// }

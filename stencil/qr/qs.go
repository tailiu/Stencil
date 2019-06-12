package qr

import (
	"fmt"
	"strings"
)

func CreateQS(QR *QR) *QS {
	qs := new(QS)
	qs.seen = make(map[string]bool)
	qs.QR = QR
	return qs
}

func (self *QS) ColSimple(colNames string) { // users.id, users.*, user_id, person_id

	cols := strings.Split(colNames, ",")

	for _, col := range cols {
		if i := strings.Index(col, "."); i < 0 {
			self.Columns = append(self.Columns, col)
		} else {
			if column := col[i+1:]; column == "*" {
				table := col[:i]
				for ptab, mappedcols := range self.QR.GetPhyMappingForLogicalTable(table) {
					for _, pair := range mappedcols {
						pColName := fmt.Sprintf("%s.%s as %s", ptab, pair[0], pair[1])
						self.Columns = append(self.Columns, pColName)
					}
				}
			} else {
				ptab, pcol := self.QR.GetPhyTabCol(col)
				pColName := fmt.Sprintf("%s.%s as %s", ptab, pcol, column)
				self.Columns = append(self.Columns, pColName)
			}
		}
	}
}

func (self *QS) ColAlias(col, alias string) { //"users.id", "user_id"

	if i := strings.Index(col, "."); i < 0 {
		pColName := fmt.Sprintf("%s as %s", col, alias)
		self.Columns = append(self.Columns, pColName)
	} else {
		ptab, pcol := self.QR.GetPhyTabCol(col)
		pColName := fmt.Sprintf("%s.%s as %s", ptab, pcol, alias)
		self.Columns = append(self.Columns, pColName)
	}
}

func (self *QS) ColFunction(funcStmt, col, alias string) {
	if i := strings.Index(col, "."); i < 0 {
		pColName := fmt.Sprintf(funcStmt+" as %s", col, alias)
		self.Columns = append(self.Columns, pColName)
	} else {
		ptab, pcol := self.QR.GetPhyTabCol(col)
		pColName := fmt.Sprintf(funcStmt+" as %s", ptab+"."+pcol, alias)
		// pColName := fmt.Sprintf("%s.%s as %s", ptab, pcol, alias)
		self.Columns = append(self.Columns, pColName)
	}
}

func (self *QS) FromSimple(table string) {

	phyTab := self.QR.GetPhyMappingForLogicalTable(table)

	prev := ""
	for ptab := range phyTab {
		if _, ok := self.seen[ptab]; !ok {
			if prev == "" {
				self.From += fmt.Sprintf(" %s ", ptab)
			} else {
				self.From += fmt.Sprintf(" JOIN %s ON %s.pk = %s.pk ", ptab, prev, ptab)
			}
			prev = ptab
			self.seen[ptab] = true
		}
	}
}

// some problems here
// multiple physical tables can appear, should be legal.
// need to carefully consider the need of qr.seen map. Probably don't need it. Need a work around.
func (self *QS) FromJoin(table, condition string) {
	// condition: previous_table.column=current_table.column

	condition_tokens := strings.Split(condition, "=")
	prev_tab, prev_col := self.QR.GetPhyTabCol(condition_tokens[0])
	curr_tab, curr_col := self.QR.GetPhyTabCol(condition_tokens[1])
	self.From += fmt.Sprintf("JOIN %s ON %s.%s = %s.%s AND %s.pk = %s.pk ", curr_tab, prev_tab, prev_col, curr_tab, curr_col, prev_tab, curr_tab)
	self.seen[curr_tab] = true

	phyTab := self.QR.GetPhyMappingForLogicalTable(table)

	prev := curr_tab

	for ptab := range phyTab {
		// if _, ok := self.seen[ptab]; !ok {
		if !strings.EqualFold(curr_tab, ptab) {
			self.From += fmt.Sprintf("JOIN %s ON %s.pk = %s.pk ", ptab, prev, ptab)
			prev = ptab
		}
	}
}

func (self *QS) FromQuery(qs *QS) {
	self.From = fmt.Sprintf("(%s) tab ", qs.GenSQL())
}

func (self *QS) WhereSimple(col1, op, col2 string) {
	ptab1, pcol1 := self.QR.GetPhyTabCol(col1)
	ptab2, pcol2 := self.QR.GetPhyTabCol(col2)
	self.Where = fmt.Sprintf("%s.%s %s %s.%s", ptab1, pcol1, op, ptab2, pcol2)
}

func (self *QS) WhereSimpleVal(col, op, val string) {
	ptab, pcol := self.QR.GetPhyTabCol(col)
	self.Where = fmt.Sprintf("%s.%s %s '%s'", ptab, pcol, op, val)
}

func (self *QS) WhereOperator(operator, col1, op, col2 string) { // AND, OR

}

func (self *QS) WhereQuery(condition string, qs *QS) { // IN, NOT IN

}

func (self *QS) GroupBy(cols string) {
	self.Group = cols
}

func (self *QS) OrderBy(cols string) {
	self.Order = cols
}

func (self *QS) LimitResult(limit string) {
	self.Limit = limit
}

func (self QS) GenSQL() string {
	sql := fmt.Sprintf("SELECT %s FROM %s", strings.Join(self.Columns, ","), self.From)
	if len(self.Where) > 0 {
		sql += fmt.Sprintf("WHERE %s ", self.Where)
	}
	if len(self.Group) > 0 {
		sql += fmt.Sprintf("GROUP BY %s ", self.Group)
	}
	if len(self.Order) > 0 {
		sql += fmt.Sprintf("ORDER BY %s ", self.Order)
	}
	if len(self.Limit) > 0 {
		sql += fmt.Sprintf("LIMIT %s ", self.Limit)
	}
	// fmt.Println("WHERE", self.Where)
	// fmt.Println("GROUPBY", self.Group)
	// fmt.Println("ORDERBY", self.Order)
	return sql
}

package qr

import (
	"fmt"
	"stencil/helper"
	"strings"
)

func CreateQS(QR *QR) *QS {
	qs := new(QS)
	qs.seen = make(map[string]bool)
	qs.TableAliases = make(map[string]map[string]string)
	qs.QR = QR
	return qs
}

func (self *QS) ColSimple(colNames string) { // users.id, users.*, user_id, person_id

	cols := strings.Split(colNames, ",")

	for _, col := range cols {
		if i := strings.Index(col, "."); i < 0 {
			self.Columns = append(self.Columns, col)
		} else {
			table := col[:i]
			if column := col[i+1:]; column == "*" {
				for ptab, mappedcols := range self.QR.GetPhyMappingForLogicalTable(table) {
					for _, pair := range mappedcols {
						pColName := fmt.Sprintf("%s.%s as \"%s.%s\"", self.getTableAlias(table, ptab), pair[0], table, pair[1])
						self.Columns = append(self.Columns, pColName)
					}
				}
			} else {
				ptab, pcol := self.QR.GetPhyTabCol(col)
				pColName := fmt.Sprintf("%s.%s as \"%s.%s\"", self.getTableAlias(table, ptab), pcol, table, column)
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
		table := col[:i]
		ptab, pcol := self.QR.GetPhyTabCol(col)
		pColName := fmt.Sprintf("%s.%s as %s", self.getTableAlias(table, ptab), pcol, alias)
		self.Columns = append(self.Columns, pColName)
	}
}

func (self *QS) ColFunction(funcStmt, col, alias string) {
	if i := strings.Index(col, "."); i < 0 {
		pColName := fmt.Sprintf(funcStmt+" as %s", col, alias)
		self.Columns = append(self.Columns, pColName)
	} else {
		table := col[:i]
		ptab, pcol := self.QR.GetPhyTabCol(col)
		pColName := fmt.Sprintf(funcStmt+" as %s", self.getTableAlias(table, ptab)+"."+pcol, alias)
		// pColName := fmt.Sprintf("%s.%s as %s", ptab, pcol, alias)
		self.Columns = append(self.Columns, pColName)
	}
}

func (self *QS) ColNull(col string) {
	pColName := fmt.Sprintf("NULL as %s", col)
	self.Columns = append(self.Columns, pColName)
}

// func (self *QS) getTableAlias(ltab, ptab string) string {
// 	if _, ok := self.TableAliases[ltab]; !ok {
// 		self.TableAliases[ltab] = make(map[string]string)
// 	}
// 	if _, ok := self.TableAliases[ltab][ptab]; !ok {
// 		self.TableAliases[ltab][ptab] = helper.RandomString(10)
// 	}
// 	return self.TableAliases[ltab][ptab]
// }

func (self *QS) ColPK(table string) {
	for ptab, _ := range self.QR.GetPhyMappingForLogicalTable(table) {
		alias := self.getTableAlias(table, ptab)
		pColName := fmt.Sprintf("%s.pk as \"pk.%s.%s\"", alias, table, alias)
		if !helper.Contains(self.Columns, pColName) {
			self.Columns = append(self.Columns, pColName)
		}
		// break
	}
}

func (self *QS) getTableAlias(ltab, ptab string) string {
	if _, ok := self.TableAliases[ltab]; !ok {
		self.TableAliases[ltab] = make(map[string]string)
	}
	if _, ok := self.TableAliases[ltab][ptab]; !ok {
		self.TableAliases[ltab][ptab] = helper.RandomString(10)
	}
	return self.TableAliases[ltab][ptab]
}

func (self *QS) FromSimple(table string) {

	phyTab := self.QR.GetPhyMappingForLogicalTable(table)

	prev := ""
	for ptab := range phyTab {
		if _, ok := self.seen[ptab]; !ok {
			if prev == "" {
				self.From += fmt.Sprintf(" %s %s ", ptab, self.getTableAlias(table, ptab))
			} else {
				self.From += fmt.Sprintf(" LEFT JOIN %s %s ON %s.pk = %s.pk ", ptab, self.getTableAlias(table, ptab), self.getTableAlias(table, prev), self.getTableAlias(table, ptab))
			}
			prev = ptab
			self.seen[ptab] = true
		}
	}
}

func (self *QS) FromJoin(table, condition string) {
	// condition: previous_table.column=current_table.column

	condition_tokens := strings.Split(condition, "=")
	prev_tab, prev_col := self.QR.GetPhyTabCol(condition_tokens[0])
	curr_tab, curr_col := self.QR.GetPhyTabCol(condition_tokens[1])
	// self.From += fmt.Sprintf("JOIN %s ON %s.%s = %s.%s AND %s.pk = %s.pk ", curr_tab, prev_tab, prev_col, curr_tab, curr_col, prev_tab, curr_tab)
	self.From += fmt.Sprintf("LEFT JOIN %s %s ON %s.%s::text = %s.%s::text ", curr_tab, self.getTableAlias(table, curr_tab), self.getTableAlias(table, prev_tab), prev_col, self.getTableAlias(table, curr_tab), curr_col)
	self.seen[curr_tab] = true

	phyTab := self.QR.GetPhyMappingForLogicalTable(table)

	prev := curr_tab

	for ptab := range phyTab {
		// if _, ok := self.seen[ptab]; !ok {
		if !strings.EqualFold(curr_tab, ptab) {
			self.From += fmt.Sprintf("LEFT JOIN %s %s ON %s.pk = %s.pk ", ptab, self.getTableAlias(table, ptab), self.getTableAlias(table, prev), self.getTableAlias(table, ptab))
			prev = ptab
		}
	}
}

func (self *QS) FromJoinList(table string, conditions []string) {

	// condition: previous_table.column=current_table.column

	var pconditions []string
	var curr_tab, curr_col string

	for _, condition := range conditions {
		condition_tokens := strings.Split(condition, "=")
		condition_lhs_tokens := strings.Split(condition_tokens[0], ".")
		condition_rhs_tokens := strings.Split(condition_tokens[1], ".")
		lhs_table, rhs_table := condition_lhs_tokens[0], condition_rhs_tokens[0]
		prev_tab, prev_col := self.QR.GetPhyTabCol(condition_tokens[0])
		curr_tab, curr_col = self.QR.GetPhyTabCol(condition_tokens[1])
		pconditions = append(pconditions, fmt.Sprintf(" %s.%s::text = %s.%s::text ", self.getTableAlias(lhs_table, prev_tab), prev_col, self.getTableAlias(rhs_table, curr_tab), curr_col))
	}

	self.From += fmt.Sprintf("LEFT JOIN %s %s ON %s ", curr_tab, self.getTableAlias(table, curr_tab), strings.Join(pconditions, "AND"))
	self.seen[curr_tab] = true

	phyTab := self.QR.GetPhyMappingForLogicalTable(table)

	prev := curr_tab

	for ptab := range phyTab {
		// if _, ok := self.seen[ptab]; !ok {
		if !strings.EqualFold(curr_tab, ptab) {
			self.From += fmt.Sprintf("LEFT JOIN %s %s ON %s.pk = %s.pk ", ptab, self.getTableAlias(table, ptab), self.getTableAlias(table, prev), self.getTableAlias(table, ptab))
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
	table1 := strings.Split(col1, ".")[0]
	table2 := strings.Split(col2, ".")[0]
	self.Where = fmt.Sprintf(" %s.%s %s %s.%s ", self.getTableAlias(table1, ptab1), pcol1, op, self.getTableAlias(table2, ptab2), pcol2)
}

func (self *QS) WhereSimpleVal(col, op, val string) {
	ptab, pcol := self.QR.GetPhyTabCol(col)
	table := strings.Split(col, ".")[0]
	self.Where = fmt.Sprintf(" %s.%s %s '%s' ", self.getTableAlias(table, ptab), pcol, op, val)
}

func (self *QS) WhereSimpleInterface(col, op string, val interface{}) {
	ptab, pcol := self.QR.GetPhyTabCol(col)
	self.vals = append(self.vals, val)
	table := strings.Split(col, ".")[0]
	// self.Where = fmt.Sprintf(" %s.%s %s $%d ", ptab, pcol, op, len(self.vals))
	self.Where = fmt.Sprintf("%s.%s::text %s '%s' ", self.getTableAlias(table, ptab), pcol, op, fmt.Sprint(val))
}

func (self *QS) WhereOperator(operator, col1, op, col2 string) { // AND, OR

}

func (self *QS) WhereOperatorVal(operator, col, op, val string) { // AND, OR
	ptab, pcol := self.QR.GetPhyTabCol(col)
	table := strings.Split(col, ".")[0]
	self.Where += fmt.Sprintf(" %s %s.%s %s '%s' ", operator, self.getTableAlias(table, ptab), pcol, op, val)
}

func (self *QS) WhereOperatorInterface(operator, col, op string, val interface{}) { // AND, OR
	ptab, pcol := self.QR.GetPhyTabCol(col)
	table := strings.Split(col, ".")[0]
	self.vals = append(self.vals, val)
	if len(self.Where) > 0 {
		// self.Where += fmt.Sprintf(" %s %s.%s %s $%d ", operator, self.getTableAlias(table, ptab), pcol, op, len(self.vals))
		self.Where += fmt.Sprintf(" %s %s.%s::text %s '%s' ", operator, self.getTableAlias(table, ptab), pcol, op, fmt.Sprint(val))
	} else {
		// self.Where = fmt.Sprintf("%s.%s %s $%d ", ptab, pcol, op, len(self.vals))
		self.Where = fmt.Sprintf("%s.%s::text %s '%s' ", self.getTableAlias(table, ptab), pcol, op, fmt.Sprint(val))
	}
}

func (self *QS) WhereOperatorBool(operator, col, op, val string) { // AND, OR
	ptab, pcol := self.QR.GetPhyTabCol(col)
	table := strings.Split(col, ".")[0]
	self.Where += fmt.Sprintf(" %s %s.%s %s %s ", operator, self.getTableAlias(table, ptab), pcol, op, val)
}

func (self *QS) WhereQuery(condition string, qs *QS) { // IN, NOT IN

}

func (self *QS) WhereString(operator, condition string) { // AND, OR, NOT

	if len(self.Where) > 0 {
		self.Where += fmt.Sprintf(" %s (%s)", operator, condition)
	} else {
		self.Where = fmt.Sprintf("(%s)", condition)
	}
}

func (self *QS) WhereMFlag(condition, flag, app_id string) { // EXISTS/NOT EXISTS, 0,1,2
	var pkcols []string
	for _, pTables := range self.TableAliases {
		for _, alias := range pTables {
			pkcols = append(pkcols, alias+".pk")
		}
	}
	q := fmt.Sprintf("%s (SELECT 1 FROM row_desc WHERE app_id = %s AND mflag IN (%s) AND rowid IN (%s))", condition, app_id, flag, strings.Join(pkcols, ","))
	self.WhereString("AND", q)
}

func (self *QS) GroupBy(col string) {
	ptab, pcol := self.QR.GetPhyTabCol(col)
	if strings.EqualFold(self.Group, "") {
		self.Group = fmt.Sprintf("%s.%s", ptab, pcol)
	} else {
		self.Group += fmt.Sprintf(", %s.%s", ptab, pcol)
	}
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
	// fmt.Println(sql)
	return sql
}
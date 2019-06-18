package qr

import (
	"fmt"
	"stencil/db"
	"strings"
)

func CreateQU(QR *QR) *QU {
	qu := new(QU)
	qu.QR = QR
	qu.Tables = make(map[string]bool)
	qu.Update = make(map[string]string)
	qu.Where = make(map[string]string)
	// qu.affected_tables = make(map[string]bool)
	return qu
}

func (self *QU) SetTable(table string) {
	self.Tables[table] = true
}

func (self *QU) SetUpdate(col, val string) {
	ptab, pcol := self.QR.GetPhyTabCol(col)
	if _, ok := self.Update[ptab]; ok {
		self.Update[ptab] = fmt.Sprintf(" AND %s = '%s'", pcol, val)
	} else {
		self.Update[ptab] = fmt.Sprintf("%s = '%s' ", pcol, val)
	}
}

func (self *QU) SetUpdate_(col, val string) {
	ptab, pcol := self.QR.GetPhyTabCol(col)
	if _, ok := self.Update[ptab]; ok {
		self.Update[ptab] = fmt.Sprintf(" AND %s = %s", pcol, val)
	} else {
		self.Update[ptab] = fmt.Sprintf("%s = %s ", pcol, val)
	}
}

func (self *QU) SetWhere(col, op, val string) {
	ptab, pcol := self.QR.GetPhyTabCol(col)
	if _, ok := self.Where[ptab]; ok {
		self.Where[ptab] = fmt.Sprintf(" AND %s = '%s'", pcol, val)
	} else {
		self.Where[ptab] = fmt.Sprintf("%s = '%s' ", pcol, val)
	}
	// self.affected_tables[ptab] = true
	self.affected_tables = append(self.affected_tables, ptab)
}

func (self *QU) SetWhere_(col, op, val string) {
	ptab, pcol := self.QR.GetPhyTabCol(col)
	if _, ok := self.Where[ptab]; ok {
		self.Where[ptab] = fmt.Sprintf(" AND %s = %s", pcol, val)
	} else {
		self.Where[ptab] = fmt.Sprintf("%s = %s ", pcol, val)
	}
	// self.affected_tables[ptab] = true
	self.affected_tables = append(self.affected_tables, ptab)
}

func (self *QU) SetWhereOperator(logop, col, op, val string) {

}

func (self *QU) SetWhereOperator_(logop, col, op, val string) {

}

func (self *QU) GetAffectedRows() {

	join, where, prev := "", "", ""
	for _, table := range self.affected_tables {
		if strings.EqualFold(join, "") {
			join = table
		} else {
			join += fmt.Sprintf(" JOIN %s ON %s.pk = %s.pk", table, prev, table)
		}
		if cond, ok := self.Where[table]; ok {
			if strings.EqualFold(where, "") {
				where = fmt.Sprintf("%s.%s", table, cond)
			} else {
				where += fmt.Sprintf(" AND %s.%s", table, cond)
			}
		}
		prev = table
	}

	sql := fmt.Sprintf("SELECT %s.pk FROM %s WHERE %s", strings.Join(self.affected_tables, ".pk,"), join, where)
	res := db.DataCall(self.QR.StencilDB, sql)

	for _, row := range res {
		for _, val := range row {
			pk := fmt.Sprint(val)
			if len(pk) > 0 {
				self.affected_rows = append(self.affected_rows, pk)
			}
		}
	}

	// for table := range self.affected_tables {
	// 	sql := fmt.Sprintf("SELECT pk FROM %s WHERE %s", table, self.Where[table])
	// 	fmt.Println(sql)
	// 	res := db.DataCall(self.QR.StencilDB, sql)
	// 	for _, row := range res {
	// 		self.affected_rows = append(self.affected_rows, fmt.Sprint(row["pk"]))
	// 	}
	// }
}

func (self *QU) GenSQL() []string {
	var sqls []string

	for table, update := range self.Update {
		if len(self.affected_tables) > 0 {
			self.GetAffectedRows()
			where := strings.Join(self.affected_rows, ",")
			sql := fmt.Sprintf("UPDATE %s SET %s WHERE pk IN (%s)", table, update, where)
			sqls = append(sqls, sql)
		} else {
			sql := fmt.Sprintf("UPDATE %s SET %s ", table, update)
			sqls = append(sqls, sql)
		}
	}
	return sqls
}

package qr

import (
	"fmt"
	"stencil/helper"
	"strings"
	"stencil/db"
	"log"
)

func CreateQS(QR *QR) *QS {
	qs := new(QS)
	qs.seen = make(map[string]bool)
	qs.PK = false
	qs.TableAliases = make(map[string]map[string]string)
	qs.QR = QR
	return qs
}

func (self *QS) SelectColumns(columns string){
	for _, col := range strings.Split(columns, ","){
		parts := strings.Split(col, ".")
		self.Columns = append(self.Columns, fmt.Sprintf("\"%s\".%s", parts[0], parts[1]))
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

func (self *QS) GenCombinedTableQuery(args map[string]string) string {

	tableID, err := db.TableID(self.QR.StencilDB, args["table"], self.QR.AppID)
	if err != nil {
		fmt.Println("Can't get Table ID for table ", args["table"])
		log.Fatal(err)
	}
	
	if _, ok := args["alias"]; !ok {
		args["alias"] = args["table"]
	}
	if _, ok := args["mflag"]; !ok {
		args["mflag"] = "0"
	}
	if _, ok := args["mark_as_delete"]; !ok {
		args["mark_as_delete"] = "false"
	}
	if _, ok := args["bag"]; !ok {
		args["bag"] = "false"
	}

	var cols, pkCols []string

	from := fmt.Sprintf("(SELECT array_agg(row_id) AS rowids FROM migration_table WHERE table_id = %s and app_id = %s AND mflag IN (%s) AND mark_as_delete = %s AND bag = %s GROUP BY group_id) MT ", tableID, self.QR.AppID, args["mflag"], args["mark_as_delete"], args["bag"])
	if user_id, ok := args["user_id"]; ok {
		from = fmt.Sprintf("(SELECT array_agg(row_id) AS rowids FROM migration_table WHERE table_id = %s and app_id = %s AND mflag IN (%s) AND mark_as_delete = %s AND bag = %s AND user_id = '%s' GROUP BY group_id) MT ", tableID, self.QR.AppID, args["mflag"], args["mark_as_delete"], args["bag"], user_id)
	}

	phyTab := self.QR.GetPhyMappingForLogicalTable(args["table"])
	phyTabKeys := helper.GetKeysOfPhyTabMap(phyTab)

	for _, ptab := range phyTabKeys {
		for _, pair := range phyTab[ptab] {
			pColName := fmt.Sprintf("%s.%s as \"%s.%s\"", self.getTableAlias(args["alias"], ptab), pair[0], args["alias"], pair[1])
			cols = append(cols, pColName)
			pSizeColName := fmt.Sprintf("pg_column_size(%s.\"%s.%s\") as \"%s.%s\"", args["alias"], args["alias"], pair[1], args["alias"], pair[1])
			self.ColumnsWSize = append(self.ColumnsWSize, pSizeColName)
		}
		pTabAlias := self.getTableAlias(args["alias"], ptab)
		from += fmt.Sprintf(" LEFT JOIN %s %s ON %s.pk = ANY(mt.rowids) ", ptab, pTabAlias, pTabAlias)
		pkCols = append(pkCols, pTabAlias + ".pk")
	}

	cols = append(cols, fmt.Sprintf("uniq(sort(array_remove(array[%s]::int4[], null))) as \"%s.rowids\"", strings.Join(pkCols, ","), args["alias"]))
	cols = append(cols, fmt.Sprintf("array_to_string(uniq(sort(array_remove(array[%s]::int4[], null))),',') as \"%s.rowids_str\"", strings.Join(pkCols, ","), args["alias"]))
	tableQuery := fmt.Sprintf(" (SELECT %s FROM %s) \"%s\" ", strings.Join(cols, ","), from, args["alias"])
	
	if len(tableQuery) > 0 {
		return tableQuery	
	}

	log.Fatal("error resolving query for table: " + args["table"])

	return ""
}

func (self *QS) FromTable(args map[string]string) {
	self.From = self.GenCombinedTableQuery(args)
}

func (self *QS) JoinTable(args map[string]string) {
	tableQuery := self.GenCombinedTableQuery(args)

	if _, ok := args["join"]; !ok {
		args["join"] = "JOIN"
	}

	var tableConditions []string

	for key, val := range args {
		if strings.Contains(key, "condition") {
			conditions := strings.Split(val, "=")
			table1 := strings.Split(conditions[0], ".")[0]
			table2 := strings.Split(conditions[1], ".")[0]
			tableConditions = append(tableConditions, fmt.Sprintf("%s.\"%s\"::varchar = %s.\"%s\"::varchar ", table1, conditions[0], table2, conditions[1]))
		}
	}
	self.From += fmt.Sprintf(" %s %s ON %s ", args["join"], tableQuery, strings.Join(tableConditions, " AND "))
}

func (self *QS) AddWhereWithValue(col, op, val string) {
	tokens := strings.Split(col, ".")
	self.Where = fmt.Sprintf("\"%s\".\"%s\" %s '%s'", tokens[0], col, op, val)
}

func (self *QS) AddWhereWithColumn(col1, op, col2 string) {
	tokens1 := strings.Split(col1, ".")
	tokens2 := strings.Split(col1, ".")
	self.Where = fmt.Sprintf("\"%s\".\"%s\" %s \"%s\".\"%s\"", tokens1[0], col1, op, tokens2[0], col2)
}

func (self *QS) AdditionalWhereWithValue(coop, col, op, val string) {
	tokens := strings.Split(col, ".")
	if len(self.Where) > 0 {
		self.Where += fmt.Sprintf(" %s \"%s\".\"%s\" %s '%s'", coop, tokens[0], col, op, val)
	}else{
		self.Where = fmt.Sprintf(" \"%s\".\"%s\" %s '%s'", tokens[0], col, op, val)
	}
}

func (self *QS) AdditionalWhereWithColumn(coop, col1, op, col2 string) {
	tokens1 := strings.Split(col1, ".")
	tokens2 := strings.Split(col1, ".")
	self.Where += fmt.Sprintf(" %s \"%s\".\"%s\" %s %s.\"%s\"", coop, tokens1[0], col1, op, tokens2[0], col2)
}

func (self *QS) AddWhereAsString(operator, condition string) { // AND, OR, NOT

	if len(self.Where) > 0 {
		self.Where += fmt.Sprintf(" %s (%s)", operator, condition)
	} else {
		self.Where = condition
	}
}

func (self *QS) RowIDs(rowids string) {
	if len(rowids) <= 0 {return}
	for table := range self.TableAliases {
		if len(self.Where) > 0 {
			self.Where += " AND "
		}
		self.Where += fmt.Sprintf("array[%s] <@ %s.\"%s.rowids\"", rowids, table, table)	
	}
}

func (self *QS) ExcludeRowIDs(rowids string) {
	if len(rowids) <= 0 {return}
	for table := range self.TableAliases {
		if len(self.Where) > 0 {
			self.Where += " AND "
		}
		self.Where += fmt.Sprintf("NOT array[%s]::int4[] @> %s.\"%s.rowids\"", rowids, table, table)	
	}
}

func (self *QS) GroupByString(col string) {
	table := strings.Split(col, ".")[0]
	if strings.EqualFold(self.Group, "") {
		self.Group = fmt.Sprintf("%s.\"%s\"", table, col)
	} else {
		self.Group += fmt.Sprintf(", %s.\"%s\"", table, col)
	}
}

func (self *QS) OrderByColumn(col string) {
	table := strings.Split(col, ".")[0]
	if strings.EqualFold(self.Order, "") {
		self.Order = fmt.Sprintf("\"%s\".\"%s\"", table, col)
	} else {
		self.Order += fmt.Sprintf(", \"%s\".\"%s\"", table, col)
	}
}

func (self *QS) OrderByFunction(f string) {
	if strings.EqualFold(self.Order, "") {
		self.Order = f
	} else {
		self.Order += ", "+f
	}
}

func (self *QS) LimitResult(limit string) {
	self.Limit = limit
}

func (self *QS) GenSQL() string {
	var arrayRowIDCols []string
	for table := range self.TableAliases {
		arrayRowIDCols = append(arrayRowIDCols, fmt.Sprintf("\"%s\".\"%s.rowids\"", table, table))
	}
	self.Columns = append(self.Columns, fmt.Sprintf("array_to_string(uniq(sort(array[%s])),',') as rowids", strings.Join(arrayRowIDCols, " || ")))
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
	log.Fatal()
	return sql
}

func (self *QS) GenSQLSize() string {
	
	sql := fmt.Sprintf("SELECT %s FROM %s", strings.Join(self.ColumnsWSize, ","), self.From)
	
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
	return sql
}

func (self *QS) GenSQLWithSize() string {
	var arrayRowIDCols []string
	var tableColSize []string
	for table := range self.TableAliases {
		tableColSize = append(tableColSize, fmt.Sprintf("COALESCE(pg_column_size(%s.*) - pg_column_size(%s.\"%s.rowids\"), 0)", table, table, table))
		arrayRowIDCols = append(arrayRowIDCols, fmt.Sprintf("%s.\"%s.rowids\"", table, table))
	}
	self.Columns = append(self.Columns, fmt.Sprintf("array_to_string(uniq(sort(array[%s])),',') as rowids", strings.Join(arrayRowIDCols, " || ")))
	self.Columns = append(self.Columns, fmt.Sprintf("%s as rowsize", strings.Join(tableColSize, " + ")))
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
	return sql
}
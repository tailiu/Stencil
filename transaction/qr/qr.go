/*
 * Query Resolver
 */

package qr

import (
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"transaction/db"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // postgres driver
	"github.com/xwb1989/sqlparser"
)

type QR struct {
	DB           *sql.DB
	PDBNAME      string
	AppName      string
	AppID        string
	BaseMappings []map[string]string
	SuppMappings []map[string]string
}

type QI struct {
	TableName  string
	Columns    []string
	Values     []string
	Conditions string
}

func NewQR(app_name string) *QR {
	qr := new(QR)
	qr.PDBNAME = "stencil"
	qr.DB = db.GetDBConn(qr.PDBNAME)
	qr.AppName = app_name
	qr.SetAppId()
	qr.GetBaseMappings()
	qr.GetSupplementaryMappings()
	return qr
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func Contains(a []string, x string) bool {
	x = strings.Trim(strings.ToLower(x), ", ")
	for _, n := range a {
		n = strings.Trim(strings.ToLower(n), ", ")
		if x == n {
			return true
		}
	}
	return false
}

func (self QI) ValueOfColumn(col string) (string, error) {

	for i, c := range self.Columns {
		if col == c {
			return self.Values[i], nil
		}
	}
	return "", errors.New("No column: " + col)
}

func (self QI) Print() {

	fmt.Println(self)

	// for i, c := range self.Columns {
	// 	fmt.Println("i:", i, "Col:", c, " || Val:", self.Values[i])
	// }
}

func (self QR) TestQuery() {
	q := `INSERT INTO customer (c_id, c_d_id, c_w_id, c_first, c_middle, c_last, c_street_1, c_street_2, c_city, c_state, c_zip, c_phone, c_since, c_credit, c_credit_lim, c_discount, c_balance, c_ytd_payment, c_payment_cnt, c_delivery_cnt, c_data, mark_delete, col33, col32, col27, col38) 
		  VALUES
		  (1, 1, 1, 'rb.#y2u\\t(_p', 'OE', 'BARBARBAR', ':XP13%zkV.Ni4LW', 'jZ)VoU1\\3:6%pGb&', 'bx,UCn[(jh]0Y<8', '(v', '562711111', '1276005133733447', '2018-10-17 16:09:46+00:00', 'GC', 50000.0, 0.26560458540916, 
		  '-99.0', 10.0, 1, 0, 'I1&Wq\\$>n+ubgy$Y(y?ovCL!@kq92@;R./l@O*s^:b"imuK+[MuLECW$pOP9r+=fZn#PyA:W2=+/^Efqbq#D9i8|^;Dx(\\wfx>a#{\\jw{Nzy.hs,:q)H] Uwr^BR^wLSYo(;jUoFGxa4PTZwO?/"sizF,,H#vFI&Rr)K;SQnI[<d2nU[MrB_dx=4nsU[4jMta&Sup#lU*CMzF=RP#N*@$[$(-.H1iL9vZa+G0DAk(L59pV(y&>6wu>s{\\-uTLlS,yKaI;7U8c9a!P996Rz8&7Q=jOUpf*1SYgGBqXvNC@q7xjf?-ZW)G@HTz"]DB)y@+ZdNcano>V@%:1tXg7^%IU^4m?9:txNe:h"2cP
		  w!<y3"M-#i*7lWDp', true, 'dataforcol33', 'dataforcol32', 'dataforcol27', 'dataforcol42')`
	q = `
		INSERT INTO item (i_id, i_im_id, i_name, i_price, i_data, mark_delete) VALUES
		(1, 4762, '#PyA:W2=+/^Efqbq#D9i', 1.4249650239944, '8|^;Dx(\\wfx>a#{\\jw{Nzy.hs,:q)H] Uwr^BR^w', false)
	`

	q = `UPDATE item SET i_id = 129188`
	_ = `UPDATE item SET i_id = 129188 WHERE mark_delete = false AND i_name = 'zain'`
	_ = `UPDATE item SET i_id = 129188, i_data = 'blahblah', col31 = '1233' WHERE ark_delete = false AND i_name = 'zain'`
	_ = `UPDATE customer SET c_id = 129188, c_data = 'blahblah', col38 = '1233' WHERE mark_delete = false AND c_name = 'zain'`
	q = `DELETE FROM customer WHERE mark_delete = false AND c_name = 'zain'`
	_ = `DELETE FROM customer`

	q = `SELECT c_id, c_data, c_d_id, col38, col27 FROM customer WHERE c_id = '123'`

	_ = `SELECT * FROM customer  WHERE c_id = '1234' AND c_data = 'aw23de'`

	_ = `SELECT * FROM customer WHERE c_id = '1234' `

	_ = `SELECT * FROM customer`

	fmt.Println("*QUERY:", q)
	fmt.Println("******************************************************************************")
	for i, q := range self.Resolve(q) {
		fmt.Println("******************************************************************************")
		fmt.Println(i+1, ":", q)
	}
}

func (self QR) NewRowId() string {
	rowid := fmt.Sprintf("%x\n", uuid.Must(uuid.NewRandom()))
	var i big.Int
	i.SetString(strings.Replace(rowid, "-", "", 4), 16)
	return i.String()
}

func (self *QR) SetAppId() string {
	sql := fmt.Sprintf("SELECT row_id from apps WHERE app_name = '%s'", self.AppName)
	if result, err := db.DataCall1(self.PDBNAME, sql); err == nil {
		self.AppID = result["row_id"]
	} else {
		self.AppID = "-1"
	}
	return self.AppID
}

func (self *QR) GetBaseMappings() {
	sql := fmt.Sprintf(`SELECT
							LOWER(app_schemas.table_name) as logical_table, 
							LOWER(app_schemas.column_name) as logical_column, 
							LOWER(physical_schema.table_name) as physical_table,  
							LOWER(physical_schema.column_name) as physical_column
						FROM 	
							physical_mappings 
							JOIN 	app_schemas ON physical_mappings.logical_attribute = app_schemas.row_id
							JOIN 	physical_schema ON physical_mappings.physical_attribute = physical_schema.row_id
						WHERE 	app_schemas.app_id  = '%s' `, self.AppID)

	self.BaseMappings = db.DataCall(self.PDBNAME, sql)
}

func (self *QR) GetSupplementaryMappings() {
	sql := fmt.Sprintf(`SELECT  LOWER(asm.table_name) as logical_table,
							LOWER(asm.column_name)  as logical_column,
							CONCAT('supplementary_',st.row_id::string) as physical_table,
							LOWER(asm.column_name)  as physical_column
						FROM 	app_schemas asm JOIN
						supplementary_tables st ON 
						st.table_name = asm.table_name AND 
						st.app_id = asm.app_id
						WHERE 	asm.app_id  = '%s' AND
						asm.row_id NOT IN (
							SELECT logical_attribute FROM physical_mappings
						)`, self.AppID)

	self.SuppMappings = db.DataCall(self.PDBNAME, sql)
}

func (self QR) GetPhyMappingForLogicalTable(ltable string) map[string][][]string {

	var phyMap = make(map[string][][]string)

	for _, mapping := range append(self.BaseMappings, self.SuppMappings...) {
		if ltable == mapping["logical_table"] {
			ptab := mapping["physical_table"]
			pcol := mapping["physical_column"]
			lcol := mapping["logical_column"]
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

func (self QR) GetInsertQueryIngs(sql string) *QI {

	qi := new(QI)

	tokens := sqlparser.NewStringTokenizer(sql)

	vswitch := 0

	var cols, vals []string

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
			qi.TableName = strings.ToLower(string(tval))
			vswitch++
		case 3:
			cols = append(cols, strings.ToLower(string(tval)))
		case 4:
			vals = append(vals, string(tval))
		}

	}

	qi.Columns = cols
	qi.Values = vals

	return qi
}

func (self QR) GetDeleteQueryIngs(sql string) *QI {

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

func (self QR) GetUpdateQueryIngs(sql string) *QI {

	qi := new(QI)

	re := regexp.MustCompile(`(?i)(update | where | set )`)
	phrases := deleteEmpty(re.Split(sql, -1))

	if len(phrases) > 1 {

		qi.TableName = phrases[0]

		updates := strings.Split(phrases[1], ",")
		var cols, vals []string

		for _, update := range updates {
			items := strings.Split(update, "=")
			cols = append(cols, strings.Trim(items[0], " ,"))
			vals = append(vals, strings.Trim(items[1], " ,"))
		}

		qi.Columns = cols
		qi.Values = vals

		if len(phrases) > 2 {

			qi.Conditions = phrases[2]
		}
	}

	return qi
}

func (self QR) GetSelectQueryIngs(sql string) *QI {

	qi := new(QI)
	re := regexp.MustCompile(`(?i)(select | from | where )`)
	phrases := deleteEmpty(re.Split(sql, -1))
	if len(phrases) > 1 {
		qi.TableName = phrases[1]
		qi.Columns = strings.Split(phrases[0], ",")
		if len(phrases) > 2 {
			qi.Conditions = phrases[2]
		}
	}
	return qi
}

func (self QR) ResolveUpdate(sql string) []string {
	var PQs []string

	qi := self.GetUpdateQueryIngs(sql)
	rowIDs := self.GetAffectedRowIDs(qi.TableName, qi.Conditions)

	phyMap := self.GetPhyMappingForLogicalTable(qi.TableName)

	for pt, mapping := range phyMap {
		updates := ""
		for _, colmap := range mapping {
			if val, err := qi.ValueOfColumn(colmap[1]); err == nil {
				updates += fmt.Sprintf("%s = %s, ", colmap[0], val)
			}
		}
		if updates != "" {
			updates := strings.Trim(updates, ", ")
			pq := fmt.Sprintf("UPDATE %s SET %s ", pt, updates)
			if len(rowIDs) > 0 {
				var rowColName string
				if pt[0:4] == "base" {
					rowColName = "base_row_id"
				} else {
					rowColName = "supp_row_id"
				}
				pq += fmt.Sprintf("WHERE %s IN (%s)", rowColName, strings.Join(rowIDs[:], ","))
			}
			PQs = append(PQs, pq)
		}
	}
	return PQs
}

func (self QR) ResolveDelete(sql string) []string {
	var PQs []string

	qi := self.GetDeleteQueryIngs(sql)
	fmt.Println(qi)
	rowIDs := self.GetAffectedRowIDs(qi.TableName, qi.Conditions)
	phyMap := self.GetPhyMappingForLogicalTable(qi.TableName)

	for pt, _ := range phyMap {
		pq := fmt.Sprintf("DELETE FROM %s ", pt)
		if len(rowIDs) > 0 {
			var rowColName string
			if pt[0:4] == "base" {
				rowColName = "base_row_id"
			} else {
				rowColName = "supp_row_id"
			}
			pq += fmt.Sprintf("WHERE %s IN (%s)", rowColName, strings.Join(rowIDs[:], ","))
		}
		PQs = append(PQs, pq)
	}
	return PQs
}

func (self QR) GetAffectedRowIDs(table, conds string) []string {

	rowIDs := []string{"1", "2"}
	return rowIDs

	// var rowIDs []string

	sql := fmt.Sprintf("SELECT * from %s WHERE %s", table, conds)

	pqs := self.ResolveSelect(sql, true)

	if len(pqs) > 0 {
		for _, rowMap := range db.DataCall(self.AppName, sql) {
			for _, val := range rowMap {
				rowIDs = append(rowIDs, val)
			}
		}
	}

	return rowIDs
}

func (self QR) ResolveSelect(sql string, args ...interface{}) []string {

	var PQs []string
	qi := self.GetSelectQueryIngs(sql)
	phyMap := self.GetPhyMappingForLogicalTable(qi.TableName)

	cols := ""
	joins := ""
	conds := qi.Conditions
	prev := ""
	for pt, mapping := range phyMap {
		joined := false
		for _, colmap := range mapping {
			if Contains(qi.Columns, colmap[1]) || Contains(qi.Columns, "*") || Contains(qi.Columns, "--") {
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
				conds = strings.Replace(conds, colmap[1], col, 1)
			} else {
			}
		}
	}
	if len(args) > 0 {
		cols = "base_row_id"
	}
	pq := fmt.Sprintf("SELECT %s FROM %s", strings.Trim(cols, ", "), strings.Trim(joins, ", "))
	if len(conds) > 0 {
		pq += fmt.Sprintf(" WHERE %s", strings.Trim(conds, ", "))
	}
	PQs = append(PQs, pq)
	return PQs
}

func (self QR) ResolveInsert(sql string) []string {

	var PQs []string
	rowID := self.NewRowId()
	qi := self.GetInsertQueryIngs(sql)
	phyMap := self.GetPhyMappingForLogicalTable(qi.TableName)

	for pt, mapping := range phyMap {
		isValid := false
		pqCols := fmt.Sprintf("INSERT INTO %s ( app_id, row_id, ", pt)
		pqVals := fmt.Sprintf("VALUES ( %s, '%s',", self.AppID, rowID)
		for _, colmap := range mapping {
			if val, err := qi.ValueOfColumn(colmap[1]); err == nil {
				isValid = true
				pqCols += fmt.Sprintf("\"%s\", ", colmap[0])
				pqVals += fmt.Sprintf("E'%s',", val)
			}
		}
		if isValid {
			pq := strings.Trim(pqCols, ", ") + ") " + strings.Trim(pqVals, ", ") + ");"
			PQs = append(PQs, pq)
		}
	}
	return PQs
}

func (self QR) Resolve(sql string) []string {

	var PQs []string

	if stmt, err := sqlparser.Parse(sql); err != nil {
		fmt.Println("Error parsing:", err)
	} else {
		switch stmt := stmt.(type) {
		case *sqlparser.Select:
			PQs = self.ResolveSelect(sql)
		case *sqlparser.Update:
			PQs = self.ResolveUpdate(sql)
		case *sqlparser.Delete:
			PQs = self.ResolveDelete(sql)
		case *sqlparser.Insert:
			PQs = self.ResolveInsert(sql)
		default:
			fmt.Println("!!! Unable to identify query type.", stmt)
		}
	}
	return PQs
}

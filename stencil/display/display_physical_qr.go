package display

import (
	"stencil/db"
	"database/sql"
	"stencil/qr"
	"log"
	"strings"
	"strconv"
)

func GetData1FromPhysicalSchema(stencilDBConn *sql.DB, QR *qr.QR, appID, cols, from, col, op, val string) map[string]interface{}  {	
	qs := qr.CreateQS(QR)
	qs.SelectColumns(cols)
	// Note that we don't care about mflag here
	qs.FromTable(map[string]string{"table":from, "mflag": "0,1"})
	qs.AddWhereWithValue(col, op, val)
	physicalQuery := qs.GenSQL()
	log.Println(physicalQuery)

	result, err := db.DataCall1(stencilDBConn, physicalQuery)
	if err != nil {
		log.Fatal(err)
	}
	
	return result
}

func GetData1FromPhysicalSchemaByRowID(stencilDBConn *sql.DB, QR *qr.QR, appID, cols, from string, rowIDs []int, restrictions []map[string]string) map[string]interface{} {	
	qs := qr.CreateQS(QR)
	// Note that we don't care about mflag here
	qs.FromTable(map[string]string{"table": from, "mflag": "0,1"})
	qs.SelectColumns(cols)
	// qs.AdditionalWhereWithValue("",statuses.id, =, #numl)
	// qs.AdditionalWhereWithValue("AND", "profiles.bio", "=", "student")
	// qs.AdditionalWhereWithValue("OR", "profiles.bio", "=", "student")
	var strRowIDs string 
	for i, rowID := range rowIDs {
		if i == 0 {
			strRowIDs += strconv.Itoa(rowID)
		} else {
			strRowIDs += "," + strconv.Itoa(rowID)
		}
		
	}
	qs.RowIDs(strRowIDs)
	physicalQuery := qs.GenSQL()

	log.Println(physicalQuery)

	result, err := db.DataCall1(stencilDBConn, physicalQuery)
	if err != nil {
		log.Fatal(err)
	}
	
	return result
}

func GetRowIDsFromData(data map[string]interface{}) []int  {
	var rowIDs []int 

	for key, val := range data {
		if strings.Contains(key, ".rowids_str") && val != nil {
			s := strings.Split(val.(string), ",")
			for _, s1 := range s {
				rowID, err := strconv.Atoi(s1)
				if err != nil {
					log.Fatal(err)
				}
				rowIDs = append(rowIDs, rowID)
			}
			return rowIDs
		}
	}
	
	return rowIDs
}

package migrate

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"transaction/config"
	"transaction/db"
	"transaction/qr"
	"transaction/transaction_log"

	escape "github.com/tj/go-pg-escape"
)

func MoveData(srcApp, tgtApp string, sql config.DataQuery, mappings config.Mapping, uid int) error {

	if appMapping, ok := mappings[tgtApp]; ok {

		if tableMapping, ok := appMapping[strings.ToLower(sql.Table)]; ok {

			srcDB := db.GetDBConn(srcApp)
			tgtDB := db.GetDBConn(tgtApp)

			for {

				row, err := db.DataCall1(srcApp, sql.SQL, uid)
				if err == nil {

					ttx, err := tgtDB.Begin()
					if err != nil {
						log.Println("ERROR! TARGET TRANSACTION CAN'T BEGIN")
						return err
					}

					stx, err := srcDB.Begin()
					if err != nil {
						log.Println("ERROR! SOURCE TRANSACTION CAN'T BEGIN")
						return err
					}

					defer ttx.Rollback()
					defer stx.Rollback()

					ucond := ""

					for col, val := range row {
						if !strings.EqualFold(col, "mark_delete") && val != "" {
							ucond += fmt.Sprintf(" %s = %s AND", col, escape.Literal(val))
						}
						// fmt.Println("col", col, "data", columns[i])
					}

					ucond = strings.TrimSuffix(ucond, "AND")
					usql := fmt.Sprintf("UPDATE %s SET mark_delete = 'true' WHERE %s", sql.Table, ucond)

					if _, err = stx.Exec(usql); err != nil {
						fmt.Println(">>>>>>>>>>> Can't update!", err)
						return err
					} else {
						fmt.Println("Updated!")
					}

					for tgtTable, tgtMap := range tableMapping {

						var cols, vals string
						for scol, tcol := range tgtMap {
							cols += tcol + ","
							vals += escape.Literal(row[scol]) + ","
						}
						cols = strings.TrimSuffix(cols, ",")
						vals = strings.TrimSuffix(vals, ",")
						insql := escape.Escape("INSERT INTO %s (%s) VALUES (%s)", tgtTable, cols, vals)

						if _, err = ttx.Exec(insql); err != nil {
							log.Println("# Can't insert!", err)
							return err
						} else {
							fmt.Println("Inserted!")
						}
					}

					stx.Commit()
					ttx.Commit()
				} else if err != nil {
					log.Println("# No more rows!")
					break
				}

			}
			return nil
		}
		return errors.New("mapping doesn't exist for table:" + sql.Table)
	}
	return errors.New("mapping doesn't exist for app:" + tgtApp)
}

func MigrateData(srcApp, tgtApp string, sql config.DataQuery, mappings config.Mapping, uid int, log_txn *transaction_log.Log_txn) error {

	if appMapping, ok := mappings[tgtApp]; ok {

		if _, ok := appMapping[strings.ToLower(sql.Table)]; ok {

			stencilDB := db.GetDBConn("stencil")
			QR := qr.NewQR(srcApp)
			TgtQR := qr.NewQR(tgtApp)
			sql.SQL = strings.Replace(sql.SQL, "$1", fmt.Sprintf("'%d'", uid), 1)

			// transform a logical request into a physical request
			if psqls := QR.Resolve(sql.SQL); len(psqls) > 0 {
				psql := psqls[0]
				// log.Println(psql)
				for {
					// according to the physical request, find one result 
					if data, err := db.DataCall1("stencil", psql); err == nil {

						// according to the row_id of the result, 
						// form queries to update records with the same row_id in different physical tables 
						updQ := QR.PhyUpdateAppIDByRowID(TgtQR.AppID, sql.Table, []string{data["base_row_id"]})
						
						transaction_log.Log_change(QR.AppID, TgtQR.AppID, sql.Table, data["base_row_id"], log_txn)

						// defer tx.Rollback()

						fmt.Println(updQ)

						tx, err := stencilDB.Begin()
						if err != nil {
							log.Println("ERROR! TARGET TRANSACTION CAN'T BEGIN")
							return err
						}

						// update each physical row
						for _, usql := range updQ {
							log.Println(usql)
							if _, err = tx.Exec(usql); err != nil {
								// log.Println("!! updQ => ", updQ)
								log.Println(">> Can't update!", err)
								return err
							}
							// fmt.Println("Updated:", uid, sql.Table, QR.AppID, "=>", TgtQR.AppID, "|", data["base_row_id"])
						}

						tx.Commit()
					} else if err != nil {
						log.Println("# No more rows!")
						break
					}
				}
			} else {
				log.Println(sql.SQL)
				log.Println("Can't convert to physical query!")
			}
			log.Println("Migration complete!")
			return nil
		}
		return errors.New("mapping doesn't exist for table:" + sql.Table)
	}
	return errors.New("mapping doesn't exist for app:" + tgtApp)
}

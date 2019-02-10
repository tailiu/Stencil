package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"transaction/atomicity"
	"transaction/config"
	"transaction/db"
	"transaction/qr"
	"strconv"

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

func MigrateData(srcApp, tgtApp string, sql config.DataQuery, mappings config.Mapping, uid int, log_txn *atomicity.Log_txn) error {

	if appMapping, ok := mappings[tgtApp]; ok {

		if _, ok := appMapping[strings.ToLower(sql.Table)]; ok {

			srcAppID, err := db.GetAppId(srcApp)

			if err != nil {
				panic("Crashing" + err.Error())
			} else {
				// fmt.Println("SrcApp ID:" + srcAppID)
			}

			tgtAppID, err := db.GetAppId(tgtApp)
			if err != nil {
				panic("Crashing" + err.Error())
			} else {
				// fmt.Println("TgtApp ID:" + tgtAppID)
			}

			QR := qr.NewQR(srcAppID, "stencil")
			TgtQR := qr.NewQR(tgtAppID, "stencil")

			sql.SQL = strings.Replace(sql.SQL, "$1", fmt.Sprintf("'%d'", uid), 1)

			// transform a logical request into a physical request
			if psqls := QR.Resolve(sql.SQL, true); len(psqls) > 0 {
				psql := psqls[0]
				log.Println("IN MIGRATE:", psql)
				for {
					// according to this physical request, find one result
					if data, err := db.DataCall1("stencil", psql); err == nil {
						// log.Println("Data:", data)
						// according to the row_id of the result,
						if len(data["base_row_id"]) > 0 {
							// form queries to update records with the same row_id in different physical tables
							// updQ := QR.PhyUpdateAppIDByRowID(TgtQR.AppID, sql.Table, []string{data["base_row_id"]})

							// before migrating, log the logical query
							atomicity.LogChange(QR.AppID, TgtQR.AppID, sql.Table, data["base_row_id"], log_txn)

							// defer tx.Rollback()

							// migrateOneLogicalRow(updQ, QR)
							migrateOneLogicalRow(QR, TgtQR.AppID, data["base_row_id"])
						}
					} else if err != nil {
						log.Println("# No more rows!")
						break
					}
				}
			} else {
				log.Println("Can't convert to physical query!")
			}
			log.Println("Migration complete!")
			return nil
		}
		return errors.New("mapping doesn't exist for table:" + sql.Table)
	}
	return errors.New("mapping doesn't exist for app:" + tgtApp)
}

func migrateOneLogicalRow(QR *qr.QR, tgt_app_ID string, base_row_id string) error {
	tx, err := QR.DB.Begin()
	if err != nil {
		log.Println("ERROR! TARGET TRANSACTION CAN'T BEGIN")
		return err
	}
	int_tgt_app_ID, err := strconv.Atoi(tgt_app_ID)
	int_base_row_id, err1 := strconv.Atoi(base_row_id)
	if err != nil || err1 != nil {
		log.Println("ERROR! CONVERT STRING APP_ID OR ROW_ID TO INT ERROR")
	}
	msql := fmt.Sprintf("UPDATE row_desc SET app_id = %d WHERE row_id = %d;", int_tgt_app_ID, int_base_row_id)
	log.Println(msql)
	if _, err = tx.Exec(msql); err != nil {
		log.Println(">> Can't update!", err)
		return err
	}
	tx.Commit()

	return nil
}

// // migrate one or several physical rows with the same Row_ID, which corresponds to one logical row
// func migrateOneLogicalRow(updQ []string, QR *qr.QR) error {
// 	tx, err := QR.DB.Begin()
// 	if err != nil {
// 		log.Println("ERROR! TARGET TRANSACTION CAN'T BEGIN")
// 		return err
// 	}

// 	// update each physical row
// 	for _, usql := range updQ {
// 		log.Println(usql)
// 		if _, err = tx.Exec(usql); err != nil {
// 			log.Println(">> Can't update!", err)
// 			return err
// 		}
// 	}

// 	tx.Commit()

// 	return nil
// }

func rollbackOneRow(undo_action sql.NullString) {
	parameters := strings.Fields(undo_action.String)

	QR := qr.NewQR(parameters[0], "stencil")

	// updQ := QR.PhyUpdateAppIDByRowID(parameters[1], parameters[2], []string{parameters[3]})
	// fmt.Println(updQ)

	// migrateOneLogicalRow(updQ, QR)
}

func RollbackMigration(txn_id int) {
	stencilDB := db.GetDBConn(atomicity.StencilDBName)

	getLogRecords := fmt.Sprintf("SELECT action_type, undo_action FROM txn_log WHERE action_id = %d ORDER BY PRIMARY KEY txn_log DESC", txn_id)
	rows, err := stencilDB.Query(getLogRecords)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var action_type string
		var undo_action sql.NullString
		if err := rows.Scan(&action_type, &undo_action); err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("%s %s\n", action_type, undo_action)

		switch action_type {
			case "COMMIT":
				log.Fatal("Can't abort an already completed action.")
			case "ABORT", "ABORTED":
				log.Fatal("Can't abort an already aborted action.")
			case "CHANGE":
				rollbackOneRow(undo_action)
			case "BEGIN_TRANSACTION":
				break
		}
	}

	atomicity.LogOutcome(&atomicity.Log_txn{DBconn: stencilDB, Txn_id: txn_id}, "ABORTED")
}

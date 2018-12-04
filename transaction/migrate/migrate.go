package migrate

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"transaction/config"
	"transaction/db"
	"transaction/qr"
)

func MigrateData(srcApp, tgtApp string, sql config.DataQuery, mappings config.Mapping, uid int) error {

	if appMapping, ok := mappings[tgtApp]; ok {

		if _, ok := appMapping[strings.ToLower(sql.Table)]; ok {

			stencilDB := db.GetDBConn("stencil")
			QR := qr.NewQR(srcApp)
			TgtQR := qr.NewQR(tgtApp)
			sql.SQL = strings.Replace(sql.SQL, "$1", fmt.Sprintf("'%d'", uid), 1)
			if psqls := QR.Resolve(sql.SQL); len(psqls) > 0 {
				psql := psqls[0]
				// log.Println(psql)
				for {
					if data, err := db.DataCall1("stencil", psql); err == nil {
						tx, err := stencilDB.Begin()
						if err != nil {
							log.Println("ERROR! TARGET TRANSACTION CAN'T BEGIN")
							return err
						}
						// fmt.Println("TO UPDATE ROWS:", len(data))
						updQ := QR.PhyUpdateAppIDByRowID(TgtQR.AppID, sql.Table, []string{data["base_row_id"]})
						// defer tx.Rollback()
						// fmt.Println("TO UPDATE TABLES:", len(updQ))
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

package apis

import (
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/migrate"
	"stencil/mthread"
	"stencil/transaction"
)

func StartMigration(uid, srcApp, srcAppID, dstApp, dstAppID, mtype string, isBlade, enableBags bool) {

	if logTxn, err := transaction.BeginTransaction(); err == nil {

		switch mtype {
		case "d":
			{
				mtype = migrate.DELETION
			}
		case "i":
			{
				mtype = migrate.INDEPENDENT
			}
		case "c":
			{
				mtype = migrate.CONSISTENT
			}
		case "b":
			{
				mtype = migrate.BAGS
			}
		case "n":
			{
				mtype = migrate.NAIVE
			}
		}

		if err != nil {
			log.Fatal(err)
		}

		mappings := config.GetSchemaMappingsFor(srcApp, dstApp)

		if mappings == nil {
			log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcApp, dstApp))
		}

		if mthread.ThreadControllerV2(uid, srcApp, srcAppID, dstApp, dstAppID, logTxn, mtype, mappings, 1, isBlade, enableBags) {
			transaction.LogOutcome(logTxn, "COMMIT")
			db.FinishMigration(logTxn.DBconn, logTxn.Txn_id, 0)
			logTxn.DBconn.Close()
		} else {
			transaction.LogOutcome(logTxn, "ABORT")
			log.Fatal("Migration transaction aborted!")
		}

	} else {
		log.Fatal("Can't begin migration transaction: ", err)
	}
}

package apis

import (
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/migrate"
	"stencil/mthread"
	"stencil/transaction"

	"github.com/gookit/color"
)

func StartMigration(uid, srcApp, srcAppID, dstApp, dstAppID, mtype string, isBlade, enableBags bool, args ...bool) {

	enableFTP, debug := false, false

	if len(args) > 0 {
		enableFTP = args[1]
	}

	if len(args) > 1 {
		debug = args[2]
	}

	mtController := migrate.MigrationThreadController{
		UID:             uid,
		MType:           mtype,
		SrcAppInfo:      migrate.App{Name: srcApp, ID: srcAppID},
		DstAppInfo:      migrate.App{Name: dstApp, ID: dstAppID},
		Blade:           isBlade,
		EnableBags:      enableBags,
		FTPFlag:         enableFTP,
		LoggerDebugFlag: debug,
	}

	mtController.Init()
	mtController.Run()
	mtController.Stop()
}

func StartMigrationSA2(uid, srcApp, srcAppID, dstApp, dstAppID, mtype string, enableBags bool) {

	if logTxn, err := transaction.BeginTransaction(); err == nil {
		MaD := "0"
		threads := 1

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
		}

		if len(mtype) <= 0 {
			log.Fatal("can't read migration type")
		}

		mappings := config.GetSchemaMappingsFor(srcApp, dstApp)

		if mappings == nil {
			log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcApp, dstApp))
		}

		if msize, err := mthread.ThreadController(uid, srcApp, srcAppID, dstApp, dstAppID, logTxn, mtype, mappings, threads, MaD, enableBags); err == nil {
			transaction.LogOutcome(logTxn, "COMMIT")
			db.FinishMigration(logTxn.DBconn, logTxn.Txn_id, msize)
			logTxn.DBconn.Close()
		} else {
			transaction.LogOutcome(logTxn, "ABORT")
			log.Println("Transaction aborted:", logTxn.Txn_id)
		}

	} else {
		log.Fatal("Can't begin migration transaction", err)
	}

	color.Success.Println("End of Migration")
}

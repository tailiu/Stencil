package apis

import (
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/migrate_v1"
	"stencil/migrate_v2"
	"stencil/transaction"

	"github.com/gookit/color"
)

func StartMigration(uid, srcApp, srcAppID, dstApp, dstAppID, mtype string, isBlade, enableBags bool, args ...bool) {

	enableFTP, debug, rootAlive := false, false, false

	if len(args) > 0 {
		enableFTP = args[0]
	}

	if len(args) > 1 {
		debug = args[1]
	}

	if len(args) > 2 {
		rootAlive = args[2]
	}

	mtController := migrate_v2.MigrationThreadController{
		UID:             uid,
		MType:           mtype,
		SrcAppInfo:      migrate_v2.App{Name: srcApp, ID: srcAppID},
		DstAppInfo:      migrate_v2.App{Name: dstApp, ID: dstAppID},
		Blade:           isBlade,
		EnableBags:      enableBags,
		FTPFlag:         enableFTP,
		LoggerDebugFlag: debug,
		DeleteRootFlag:  rootAlive,
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
				mtype = migrate_v1.DELETION
			}
		case "i":
			{
				mtype = migrate_v1.INDEPENDENT
			}
		case "c":
			{
				mtype = migrate_v1.CONSISTENT
			}
		}

		if len(mtype) <= 0 {
			log.Fatal("can't read migration type")
		}

		mappings := config.GetSchemaMappingsFor(srcApp, dstApp)

		if mappings == nil {
			log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcApp, dstApp))
		}

		if msize, err := migrate_v1.ThreadController(uid, srcApp, srcAppID, dstApp, dstAppID, logTxn, mtype, mappings, threads, MaD, enableBags); err == nil {
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

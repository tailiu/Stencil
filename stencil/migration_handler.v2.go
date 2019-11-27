/*
 * Logical Migration Handler
 */

package main

import (
	"fmt"
	"log"
	"os"
	"stencil/config"
	"stencil/db"
	"stencil/migrate"
	"stencil/mthread"
	"stencil/transaction"
	"strconv"
)

func main() {
	// evalConfig := evaluation.InitializeEvalConfig()
	if logTxn, err := transaction.BeginTransaction(); err == nil {
		srcApp, srcAppID := os.Args[4], os.Args[5]
		dstApp, dstAppID := os.Args[6], os.Args[7]
		dflag := "f"
		if len(os.Args) > 8 {
			dflag = os.Args[8]
		}
		threads, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		uid := os.Args[2]

		var mtype string

		switch os.Args[3] {
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

		// mWorker := migrate.CreateLMigrationWorker(uid, srcApp, srcAppID, dstApp, dstAppID, logTxn, mtype)
		mappings := config.GetSchemaMappingsFor(srcApp, dstApp)
		if mappings == nil {
			log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcApp, dstApp))
		}
		// var wg sync.WaitGroup
		if dflag == "t" {
			// for threadID := 0; threadID < threads; threadID++ {
			// 	wg.Add(1)
			// 	go func() {
			// 		defer wg.Done()
			// 		app_display_algorithm.DisplayThread(dstApp, logTxn.Txn_id)
			// 	}()
			// }
		}
		if mthread.ThreadControllerV2(uid, srcApp, srcAppID, dstApp, dstAppID, logTxn, mtype, mappings, threads) {
			transaction.LogOutcome(logTxn, "COMMIT")
			// wg.Wait()
			db.FinishMigration(logTxn.DBconn, logTxn.Txn_id, 0)
		} else {
			transaction.LogOutcome(logTxn, "ABORT")
		}
		// evaluation.AnomaliesDanglingData(fmt.Sprint(logTxn.Txn_id), evalConfig)
		// evaluation.MigrationRate(fmt.Sprint(logTxn.Txn_id), evalConfig)
		// evaluation.SystemLevelDanglingData(fmt.Sprint(logTxn.Txn_id), evalConfig)
		// evaluation.GetSize(fmt.Sprint(logTxn.Txn_id), evalConfig)
		// evaluation.GetTime(fmt.Sprint(logTxn.Txn_id), evalConfig)
	} else {
		log.Fatal("Can't begin migration transaction", err)
	}
	// migrate.RollbackMigration(1503622861)
}
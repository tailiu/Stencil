/*
 * Logical Migration Handler
 */

package main

import (
	"log"
	"os"
	"sync"
	"stencil/migrate"
	"stencil/config"
	"stencil/mthread"
	"stencil/transaction"
	"stencil/evaluation"
	"stencil/db"
	"stencil/app_display_algorithm"
	"strconv"
	"fmt"
)

func main() {
	evalConfig := evaluation.InitializeEvalConfig()
	if logTxn, err := transaction.BeginTransaction(); err == nil {
		srcApp, srcAppID := "diaspora", "1"
		dstApp, dstAppID := "mastodon", "2"
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
		var wg sync.WaitGroup
		for threadID := 0; threadID < threads; threadID++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				app_display_algorithm.DisplayThread(dstApp, logTxn.Txn_id)
			}()
		}
		if mthread.LThreadController(uid, srcApp, srcAppID, dstApp, dstAppID, logTxn, mtype, mappings, threads) {
			transaction.LogOutcome(logTxn, "COMMIT")
			wg.Wait()
			db.FinishMigration(logTxn.DBconn, logTxn.Txn_id, 0)
		} else {
			transaction.LogOutcome(logTxn, "ABORT")
		}
		// evaluation.AnomaliesDanglingData(fmt.Sprint(logTxn.Txn_id), evalConfig)
		// evaluation.MigrationRate(fmt.Sprint(logTxn.Txn_id), evalConfig)
		// evaluation.SystemLevelDanglingData(fmt.Sprint(logTxn.Txn_id), evalConfig)
		evaluation.GetSize(fmt.Sprint(logTxn.Txn_id), evalConfig)
		evaluation.GetTime(fmt.Sprint(logTxn.Txn_id), evalConfig)
	} else {
		log.Fatal("Can't begin migration transaction", err)
	}
	// migrate.RollbackMigration(1503622861)
}

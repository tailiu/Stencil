/*
 * Migration Handler
 */

package main

import (
	"log"
	"os"
	"stencil/migrate"
	"stencil/mthread"
	"stencil/transaction"
	"stencil/evaluation"
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

		mWorker := migrate.CreateLMigrationWorker(uid, srcApp, srcAppID, dstApp, dstAppID, logTxn, mtype)

		if mthread.LThreadController(mWorker, threads) {
			transaction.LogOutcome(logTxn, "COMMIT")
		} else {
			transaction.LogOutcome(logTxn, "ABORT")
		}
		evaluation.AnomaliesDanglingData(fmt.Sprint(logTxn.Txn_id), evalConfig)
	} else {
		log.Fatal("Can't begin migration transaction", err)
	}
	// migrate.RollbackMigration(1503622861)
}

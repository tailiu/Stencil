/*
 * Migration Handler
 */

package main

import (
	"log"
	"fmt"
	"time"
	"stencil/migrate"
	"stencil/mthread"
	"stencil/transaction"
	"stencil/evaluation"
	"stencil/db"
)

const (
	srcApp, srcAppID =  "diaspora", "1"
	dstApp, dstAppID =  "mastodon", "2"
	threads = 1
)

func main() {
	uids := db.GetUserListFromAppDB(srcApp, "users", "id")
	evalConfig := evaluation.InitializeEvalConfig()
	for _, uid := range uids {
		if logTxn, err := transaction.BeginTransaction(); err == nil {
			mWorker := migrate.CreateLMigrationWorker(fmt.Sprint(uid), srcApp, srcAppID, dstApp, dstAppID, logTxn, migrate.CONSISTENT)
			if mthread.LThreadController(mWorker, threads, evalConfig) {
				transaction.LogOutcome(logTxn, "COMMIT")
			} else {
				transaction.LogOutcome(logTxn, "ABORT")
			}
		} else {
			log.Fatal("Can't begin migration transaction", err)
		}
		time.Sleep(2 * time.Second)
	}
}

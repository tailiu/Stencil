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
)

func main() {

	if logTxn, err := transaction.BeginTransaction(); err == nil {
		srcApp, srcAppID := "diaspora", "1"
		dstApp, dstAppID := "mastodon", "2"
		uid := os.Args[1]

		mWorker := migrate.CreateMigrationWorker(uid, srcApp, srcAppID, dstApp, dstAppID, logTxn, migrate.INDEPENDENT)

		if mthread.ThreadController(mWorker) {
			transaction.LogOutcome(logTxn, "COMMIT")
		} else {
			transaction.LogOutcome(logTxn, "ABORT")
		}
	} else {
		log.Fatal("Can't begin migration transaction", err)
	}
	// migrate.RollbackMigration(1503622861)
}

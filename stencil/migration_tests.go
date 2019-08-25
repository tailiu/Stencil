package main

import (
	"log"
	"stencil/migrate"
	"stencil/mthread"
	"stencil/transaction"
	"stencil/evaluation"
	"stencil/db"
	"stencil/config"
)

const (
	srcApp, srcAppID =  "diaspora", "1"
	dstApp, dstAppID =  "mastodon", "2"
	threads = 1
)

func main() {
	srcAppConfig, err := config.CreateAppConfig(srcApp, srcAppID)
	if err != nil {
		log.Fatal(err)
	}
	dstAppConfig, err := config.CreateAppConfig(dstApp, dstAppID)
	if err != nil {
		log.Fatal(err)
	}
	evalConfig := evaluation.InitializeEvalConfig()
	uids := db.GetUserListFromAppDB(srcApp, "users", "id")
	for _, uid := range uids {
		if logTxn, err := transaction.BeginTransaction(); err == nil {
			mWorker := migrate.CreateLMigrationWorkerWithAppsConfig(uid, srcApp, srcAppID, dstApp, dstAppID, logTxn, migrate.CONSISTENT, srcAppConfig, dstAppConfig)
			if mthread.LThreadController(mWorker, threads, evalConfig) {
				transaction.LogOutcome(logTxn, "COMMIT")
			} else {
				transaction.LogOutcome(logTxn, "ABORT")
			}
			evaluation.AnomaliesDanglingData(fmt.Sprint(logTxn.Txn_id), evalConfig)
		} else {
			log.Fatal("Can't begin migration transaction", err)
		}
	}
}

/*
* Both Migration Handler
*/

package main

import (
	"log"
	"os"
	"sync"
	"stencil/migrate"
	"stencil/db"
	"stencil/config"
	"stencil/mthread"
	"stencil/transaction"
	"stencil/display_algorithm"
	"stencil/evaluation"
	"strconv"
	"fmt"
)

func main() {
	evalConfig := evaluation.InitializeEvalConfig()
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

	appLogTxn, err := transaction.BeginTransaction()
	if err != nil {
		log.Fatal("Can't begin appLogTxn transaction", err)
	}
	mappings := config.GetSchemaMappingsFor(srcApp, dstApp)
	if mappings == nil {
		log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcApp, dstApp))
	}
	if mthread.LThreadController(uid, srcApp, srcAppID, dstApp, dstAppID, appLogTxn, mtype, mappings, threads) {
		transaction.LogOutcome(appLogTxn, "COMMIT")
	} else {
		transaction.LogOutcome(appLogTxn, "ABORT")
		log.Println("Transaction aborted:", appLogTxn.Txn_id)
	}

	stencilLogTxn, err := transaction.BeginTransaction()
	if err != nil {
		log.Fatal("Can't begin stencilLogTxn transaction", err)
	}
	var wg sync.WaitGroup
	for threadID := 0; threadID < threads; threadID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			display_algorithm.DisplayThread(dstApp, stencilLogTxn.Txn_id, false)
		}()
	}
	if msize, err := mthread.ThreadController(uid, srcApp, srcAppID, dstApp, dstAppID, stencilLogTxn, mtype, mappings, threads, "0"); err==nil {
		transaction.LogOutcome(stencilLogTxn, "COMMIT")
		wg.Wait()
		db.FinishMigration(stencilLogTxn.DBconn, stencilLogTxn.Txn_id, msize)
		evaluation.GetTime(fmt.Sprint(stencilLogTxn.Txn_id), evalConfig)
		evaluation.GetTime(fmt.Sprint(appLogTxn.Txn_id), evalConfig)
		evaluation.GetSize(fmt.Sprint(appLogTxn.Txn_id), evalConfig)
		evaluation.GetDataDowntimeInStencil(fmt.Sprint(stencilLogTxn.Txn_id), evalConfig)
		evaluation.GetDataDowntimeInNaiveMigration(fmt.Sprint(stencilLogTxn.Txn_id), fmt.Sprint(appLogTxn.Txn_id), evalConfig)
	} else {
		transaction.LogOutcome(stencilLogTxn, "ABORT")
		log.Println("Transaction aborted:", stencilLogTxn.Txn_id)
	}
}

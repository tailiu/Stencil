/*
* Both Migration Handler
*/

package main

import (
	"log"
	"os"
	"stencil/migrate"
	"stencil/config"
	"stencil/mthread"
	"stencil/transaction"
	// "stencil/evaluation"
	"strconv"
	"fmt"
)

func main() {
	// evalConfig := evaluation.InitializeEvalConfig()
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
	}
	stencilLogTxn, err := transaction.BeginTransaction()
	if err != nil {
		log.Fatal("Can't begin stencilLogTxn transaction", err)
	}
	if mthread.ThreadController(uid, srcApp, srcAppID, dstApp, dstAppID, stencilLogTxn, mtype, mappings, threads, "0") {
		transaction.LogOutcome(stencilLogTxn, "COMMIT")
	} else {
		transaction.LogOutcome(stencilLogTxn, "ABORT")
	}
	// evaluation.GetDataDowntimeInStencil(fmt.Sprint(stencilLogTxn.Txn_id), evalConfig)
	// evaluation.GetDataDowntimeInNaiveMigration(fmt.Sprint(stencilLogTxn.Txn_id), fmt.Sprint(appLogTxn.Txn_id), evalConfig)
}

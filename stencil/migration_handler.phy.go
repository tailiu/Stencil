/*
 * Physical Migration Handler
 */

package main

import (
	"os"
	"stencil/SA2_migrate"
)

func main() {
	// evalConfig := evaluation.InitializeEvalConfig()

	uid := os.Args[2]

	srcApp, srcAppID := os.Args[4], os.Args[5]
	dstApp, dstAppID := os.Args[6], os.Args[7]

	mtype := os.Args[3]

	enableBags := false

	// apis.StartMigrationSA2(uid, srcApp, srcAppID, dstApp, dstAppID, mtype, enableBags)
	threads := 1

	enableDisplay, displayInFirstPhase, enableBags := true, true, true

	SA2_migrate.Controller(
		uid, srcApp, srcAppID, dstApp, dstAppID,
		mtype, threads, enableDisplay,
		displayInFirstPhase, enableBags,
	)
}

// func main() {
// 	// evalConfig := evaluation.InitializeEvalConfig()

// 	if logTxn, err := transaction.BeginTransaction(); err == nil {
// 		MaD := "0"
// 		if len(os.Args) > 8 {
// 			MaD = os.Args[8]
// 		}
// 		srcApp, srcAppID := os.Args[4], os.Args[5]
// 		dstApp, dstAppID := os.Args[6], os.Args[7]
// 		threads, err := strconv.Atoi(os.Args[1])
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		uid := os.Args[2]

// 		var mtype string

// 		switch os.Args[3] {
// 		case "d":
// 			{
// 				mtype = migrate.DELETION
// 			}
// 		case "i":
// 			{
// 				mtype = migrate.INDEPENDENT
// 			}
// 		case "c":
// 			{
// 				mtype = migrate.CONSISTENT
// 			}
// 		}

// 		if len(mtype) <= 0 {
// 			log.Fatal("can't read migration type")
// 		}
// 		mappings := config.GetSchemaMappingsFor(srcApp, dstApp)
// 		if mappings == nil {
// 			log.Fatal(fmt.Sprintf("Can't find mappings from [%s] to [%s].", srcApp, dstApp))
// 		}
// 		// var wg sync.WaitGroup
// 		// for threadID := 0; threadID < threads; threadID++ {
// 		// 	wg.Add(1)
// 		// 	go func() {
// 		// 		defer wg.Done()
// 		// 		return
// 		// 		display_algorithm.DisplayThread(dstApp, logTxn.Txn_id, false)
// 		// 	}()
// 		// }
// 		if msize, err := mthread.ThreadController(uid, srcApp, srcAppID, dstApp, dstAppID, logTxn, mtype, mappings, threads, MaD); err == nil {
// 			transaction.LogOutcome(logTxn, "COMMIT")
// 			// wg.Wait()
// 			db.FinishMigration(logTxn.DBconn, logTxn.Txn_id, msize)
// 		} else {
// 			transaction.LogOutcome(logTxn, "ABORT")
// 			log.Println("Transaction aborted:", logTxn.Txn_id)
// 		}
// 		// evaluation.GetDataBagOfUser(uid, evalConfig)
// 		// evaluation.GetTime(fmt.Sprint(logTxn.Txn_id), evalConfig)
// 	} else {
// 		log.Fatal("Can't begin migration transaction", err)
// 	}
// 	// migrate.RollbackMigration(1503622861)
// }

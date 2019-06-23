/*
 * Migration Handler
 */

package main

import (
	"fmt"
	"log"
	"os"
	"stencil/config"
	m2 "stencil/migrate"
	migrate "stencil/migrate/phy"
	"stencil/transaction"
	"sync"
)

type ThreadChannel struct {
	Finished  bool
	Thread_id int
}

func main() {

	var wg sync.WaitGroup

	srcApp, srcAppID := "diaspora", "1"
	dstApp, dstAppID := "mastodon", "2"
	threads_num := 1
	// uid := 4716
	uid := os.Args[1]
	commitChannel := make(chan ThreadChannel)
	// startFrom, inc := 4670, 10
	// for uid := startFrom; uid < startFrom+inc; uid += 1 {
	config.LoadSchemaMappings()
	logTxn, err := transaction.BeginTransaction()
	if err == nil {
		for thread_id := 1; thread_id <= threads_num; thread_id++ {
			wg.Add(1)
			go func(thread_id int, commitChannel chan ThreadChannel) {
				defer wg.Done()
				if srcAppConfig, err := config.CreateAppConfig(srcApp, srcAppID); err != nil {
					commitChannel <- ThreadChannel{Finished: false, Thread_id: thread_id}
					log.Fatal(err)
				} else {
					if dstAppConfig, err := config.CreateAppConfig(dstApp, dstAppID); err != nil {
						commitChannel <- ThreadChannel{Finished: false, Thread_id: thread_id}
						log.Fatal(err)
					} else {
						migrate.RemoveUserFromApp(uid, dstAppID, logTxn)
						if rootNode := migrate.GetRoot(srcAppConfig, fmt.Sprint(uid), logTxn); rootNode != nil {
							var wList = new(m2.WaitingList)
							var invalidList = new(m2.InvalidList)

							migrate.MigrateProcess(fmt.Sprint(uid), srcAppConfig, dstAppConfig, rootNode, wList, invalidList, logTxn)

						} else {
							fmt.Println("Root Node can't be fetched!")
						}
						// dstAppConfig.CloseDBConn()
					}
					// srcAppConfig.CloseDBConn()
					commitChannel <- ThreadChannel{Finished: true, Thread_id: thread_id}
				}
			}(thread_id, commitChannel)
		}
		go func() {
			wg.Wait()
			close(commitChannel)
		}()
	} else {
		log.Println("Can't begin migration transaction", err)
		transaction.LogOutcome(logTxn, "ABORT")
		// transaction.CloseDBConn(logTxn)
	}

	txnCommit := true

	for threadResponse := range commitChannel {
		fmt.Println("THREAD FINISHED WORKING", threadResponse)
		if !threadResponse.Finished {
			txnCommit = false
		}
	}

	if txnCommit {
		transaction.LogOutcome(logTxn, "COMMIT")
	} else {
		transaction.LogOutcome(logTxn, "ABORT")
	}

	// }

	// settingsFileName := "mappings"
	// // fromApp := "mastodon"
	// // toApp := "diaspora"
	// if schemaMappings, err := config.ReadSchemaMappingSettings(settingsFileName); err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println(schemaMappings)
	// }

	// initAppLevelMigration(7, "app1", "app5")
	// initStencilMigration(61, "app3", "app4")
	// QR := qr.NewQR("app1")
	// QR.TestQuery()

	// migrate.RollbackMigration(1503622861)
}

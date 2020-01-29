package mthread

import (
	"errors"
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/migrate"
	"stencil/transaction"
	"strings"
	"sync"
	"time"
)

func ThreadController(uid, srcApp, srcAppID, dstApp, dstAppID string, logTxn *transaction.Log_txn, mtype string, mappings *config.MappedApp, threads int, MaD string) (int, error) {
	var wg sync.WaitGroup

	commitChannel := make(chan ThreadChannel)

	if threads != 0 {
		if !db.RegisterMigration(uid, srcAppID, dstAppID, mtype, logTxn.Txn_id, threads, logTxn.DBconn, false) {
			log.Fatal("Unable to register migration!")
		} else {
			log.Println("Migration registered:", mtype)
		}
	} else {
		threads = 1
	}

	for threadID := 0; threadID < threads; threadID++ {
		wg.Add(1)
		go func(thread_id int, commitChannel chan ThreadChannel) {
			defer wg.Done()
			mWorker := migrate.CreateMigrationWorker(uid, srcApp, srcAppID, dstApp, dstAppID, logTxn, mtype, MaD, mappings)
			switch mWorker.MType() {
			case migrate.DELETION:
				{
					for {
						if err := mWorker.DeletionMigration(mWorker.GetRoot(), thread_id); err != nil {
							if !strings.Contains(err.Error(), "deadlock") {
								mWorker.RenewDBConn()
							} else {
								fmt.Print(">>>>>>>>>>>>>>>>>>>>>>> RESTART AFTER DEADLOCK <<<<<<<<<<<<<<<<<<<<<<<<<<<")
							}
							continue
						}
						break
					}
					for {
						if err := mWorker.SecondPhase(thread_id); err != nil {
							if !strings.Contains(err.Error(), "deadlock") {
								mWorker.RenewDBConn()
							} else {
								fmt.Print(">>>>>>>>>>>>>>>>>>>>>>> RESTART AFTER DEADLOCK <<<<<<<<<<<<<<<<<<<<<<<<<<<")
							}
							continue
						}
						break
					}
					for {
						if err := mWorker.MigrateProcessBags(thread_id); err != nil {
							if !strings.Contains(err.Error(), "deadlock") {
								mWorker.RenewDBConn()
							} else {
								fmt.Print(">>>>>>>>>>>>>>>>>>>>>>> RESTART AFTER DEADLOCK <<<<<<<<<<<<<<<<<<<<<<<<<<<")
							}
							continue
						}
						break
					}
				}
			case migrate.CONSISTENT:
				{
					for {
						if err := mWorker.ConsistentMigration(thread_id); err != nil {
							if !strings.Contains(err.Error(), "deadlock") {
								mWorker.RenewDBConn()
							} else {
								fmt.Print(">>>>>>>>>>>>>>>>>>>>>>> RESTART AFTER DEADLOCK <<<<<<<<<<<<<<<<<<<<<<<<<<<")
							}
							continue
						}
						break
					}
				}
			case migrate.INDEPENDENT:
				{
					for {
						if err := mWorker.IndependentMigration(thread_id); err != nil {
							if !strings.Contains(err.Error(), "deadlock") {
								mWorker.RenewDBConn()
							} else {
								fmt.Print(">>>>>>>>>>>>>>>>>>>>>>> RESTART AFTER DEADLOCK <<<<<<<<<<<<<<<<<<<<<<<<<<<")
							}
							continue
						}
						break
					}
				}

			}
			commitChannel <- ThreadChannel{Finished: true, Thread_id: thread_id, size: mWorker.Size}
		}(threadID, commitChannel)
	}

	go func() {
		wg.Wait()
		close(commitChannel)
	}()

	finished := true
	msize := 0
	var finished_threads []string
	for threadResponse := range commitChannel {
		fmt.Println("THREAD FINISHED WORKING", threadResponse, strings.Join(finished_threads, ","))
		if !threadResponse.Finished {
			finished = false
		}
		finished_threads = append(finished_threads, fmt.Sprint(threadResponse.Thread_id))
		msize += threadResponse.size
	}

	if mtype == migrate.DELETION {
		// mWorker.HandleLeftOverWaitingNodes()
	}

	// db.FinishMigration(logTxn.DBconn, logTxn.Txn_id, msize)
	if finished {
		return msize, nil
	} else {
		return msize, errors.New("Some thread crashed?")
	}
}

func LThreadController(uid, srcApp, srcAppID, dstApp, dstAppID string, logTxn *transaction.Log_txn, mtype string, mappings *config.MappedApp, threads int) bool {
	var wg sync.WaitGroup

	commitChannel := make(chan ThreadChannel)

	if threads != 0 {
		if !db.RegisterMigration(uid, srcAppID, dstAppID, mtype, logTxn.Txn_id, threads, logTxn.DBconn, true) {
			log.Fatal("Unable to register migration!")
		} else {
			log.Println("Migration registered:", mtype)
		}
	} else {
		threads = 1
	}

	for threadID := 0; threadID < threads; threadID++ {
		wg.Add(1)
		go func(thread_id int, commitChannel chan ThreadChannel) {
			defer wg.Done()
			mWorker := migrate.CreateLMigrationWorker(uid, srcApp, srcAppID, dstApp, dstAppID, logTxn, mtype, mappings)
			switch mWorker.MType() {
			case migrate.DELETION:
				{
					for {
						if err := mWorker.DeletionMigration(mWorker.GetRoot(), thread_id); err != nil {
							if !strings.Contains(err.Error(), "deadlock") {
								mWorker.RenewDBConn()
							} else {
								fmt.Print(">>>>>>>>>>>>>>>>>>>>>>> RESTART AFTER DEADLOCK <<<<<<<<<<<<<<<<<<<<<<<<<<<")
							}
							continue
						}
						break
					}
					for {
						if err := mWorker.SecondPhase(thread_id); err != nil {
							if !strings.Contains(err.Error(), "deadlock") {
								mWorker.RenewDBConn()
							} else {
								fmt.Print(">>>>>>>>>>>>>>>>>>>>>>> RESTART AFTER DEADLOCK <<<<<<<<<<<<<<<<<<<<<<<<<<<")
							}
							continue
						}
						break
					}
				}
			case migrate.CONSISTENT:
				{
					for {
						if err := mWorker.ConsistentMigration(thread_id); err != nil {
							if !strings.Contains(err.Error(), "deadlock") {
								mWorker.RenewDBConn()
							} else {
								fmt.Print(">>>>>>>>>>>>>>>>>>>>>>> RESTART AFTER DEADLOCK <<<<<<<<<<<<<<<<<<<<<<<<<<<")
							}
							continue
						}
						break
					}
				}
			case migrate.INDEPENDENT:
				{
					for {
						if err := mWorker.IndependentMigration(thread_id); err != nil {
							if !strings.Contains(err.Error(), "deadlock") {
								mWorker.RenewDBConn()
							} else {
								fmt.Print(">>>>>>>>>>>>>>>>>>>>>>> RESTART AFTER DEADLOCK <<<<<<<<<<<<<<<<<<<<<<<<<<<")
							}
							continue
						}
						break
					}
				}

			}
			commitChannel <- ThreadChannel{Finished: true, Thread_id: thread_id}
		}(threadID, commitChannel)
	}

	go func() {
		wg.Wait()
		close(commitChannel)
	}()

	finished := true

	var finished_threads []string
	for threadResponse := range commitChannel {
		fmt.Println("THREAD FINISHED WORKING", threadResponse, strings.Join(finished_threads, ","))
		if !threadResponse.Finished {
			finished = false
		}
		finished_threads = append(finished_threads, fmt.Sprint(threadResponse.Thread_id))
	}

	if mtype == migrate.DELETION {
		// mWorker.HandleLeftOverWaitingNodes()
	}

	// db.FinishMigration(logTxn.DBconn, logTxn.Txn_id, 0)
	return finished
}

func ThreadControllerV2(uid, srcApp, srcAppID, dstApp, dstAppID string, logTxn *transaction.Log_txn, mtype string, mappings *config.MappedApp, threads int) bool {
	var wg sync.WaitGroup

	commitChannel := make(chan ThreadChannel)

	if threads != 0 {
		if !db.RegisterMigration(uid, srcAppID, dstAppID, mtype, logTxn.Txn_id, threads, logTxn.DBconn, true) {
			log.Fatal("Unable to register migration!")
		} else {
			log.Println("Migration registered:", mtype)
		}
	} else {
		threads = 1
	}

	for threadID := 0; threadID < threads; threadID++ {
		wg.Add(1)
		go func(thread_id int, commitChannel chan ThreadChannel) {
			defer wg.Done()
			mWorker := migrate.CreateMigrationWorkerV2(uid, srcApp, srcAppID, dstApp, dstAppID, logTxn, mtype, mappings, threadID)

			switch mWorker.MType() {
			case migrate.DELETION:
				{
					for {
						break
						if err := mWorker.MigrateBags(thread_id); err != nil {
							log.Println("@ThreadControllerV2 > MigrateBags | Crashed with error: ", err)
							time.Sleep(time.Second * 5)
							continue
						}
						break
					}

					for {
						// log.Println("@ThreadControllerV2 > DeletionMigration skipped !!!!!!!!!")
						// break
						if err := mWorker.DeletionMigration(mWorker.GetRoot(), thread_id); err != nil {
							if !strings.Contains(err.Error(), "deadlock") {
								mWorker.RenewDBConn()
							} else {
								fmt.Print(">>>>>>>>>>>>>>>>>>>>>>> RESTART AFTER DEADLOCK <<<<<<<<<<<<<<<<<<<<<<<<<<<")
							}
							continue
						}
						break
					}
				}
			case migrate.CONSISTENT:
				{
					for {
						if err := mWorker.ConsistentMigration(thread_id); err != nil {
							if !strings.Contains(err.Error(), "deadlock") {
								mWorker.RenewDBConn()
							} else {
								fmt.Print(">>>>>>>>>>>>>>>>>>>>>>> RESTART AFTER DEADLOCK <<<<<<<<<<<<<<<<<<<<<<<<<<<")
							}
							continue
						}
						break
					}
				}
			case migrate.INDEPENDENT:
				{
					for {
						if err := mWorker.IndependentMigration(thread_id); err != nil {
							if !strings.Contains(err.Error(), "deadlock") {
								mWorker.RenewDBConn()
							} else {
								fmt.Print(">>>>>>>>>>>>>>>>>>>>>>> RESTART AFTER DEADLOCK <<<<<<<<<<<<<<<<<<<<<<<<<<<")
							}
							continue
						}
						break
					}
				}

			}
			commitChannel <- ThreadChannel{Finished: true, Thread_id: thread_id}
		}(threadID, commitChannel)
	}

	go func() {
		wg.Wait()
		close(commitChannel)
	}()

	finished := true

	var finished_threads []string
	for threadResponse := range commitChannel {
		fmt.Println("THREAD FINISHED WORKING", threadResponse, strings.Join(finished_threads, ","))
		if !threadResponse.Finished {
			finished = false
		}
		finished_threads = append(finished_threads, fmt.Sprint(threadResponse.Thread_id))
	}

	if mtype == migrate.DELETION {
		// mWorker.HandleLeftOverWaitingNodes()
	}

	// db.FinishMigration(logTxn.DBconn, logTxn.Txn_id, 0)
	return finished
}

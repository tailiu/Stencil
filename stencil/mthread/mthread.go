package mthread

import (
	"fmt"
	"log"
	"stencil/migrate"
	"strings"
	"sync"
)

func ThreadController(mWorker migrate.MigrationWorker, threads int) bool {
	var wg sync.WaitGroup

	commitChannel := make(chan ThreadChannel)

	if threads != 0 {
		if !mWorker.RegisterMigration(mWorker.MType(), threads) {
			log.Fatal("Unable to register migration!")
		} else {
			log.Println("Migration registered:", mWorker.MType())
		}
	} else {
		threads = 1
	}

	for threadID := 0; threadID < threads; threadID++ {
		wg.Add(1)
		go func(thread_id int, commitChannel chan ThreadChannel) {
			defer wg.Done()
			switch mWorker.MType() {
			case migrate.DELETION:
				{
					for {
						if err := mWorker.DeletionMigration(mWorker.GetRoot(), thread_id); err != nil {
							mWorker.RenewDBConn()
							continue
						}
						break
					}
					for {
						if err := mWorker.SecondPhase(thread_id); err != nil {
							mWorker.RenewDBConn()
							continue
						}
						break
					}
					for {
						if err := mWorker.MigrateProcessBags(thread_id); err != nil {
							mWorker.RenewDBConn()
							continue
						}
						break
					}
				}
			case migrate.CONSISTENT:
				{
					for {
						if err := mWorker.ConsistentMigration(thread_id); err != nil {
							mWorker.RenewDBConn()
							continue
						}
						break
					}
				}
			case migrate.INDEPENDENT:
				{
					for {
						if err := mWorker.IndependentMigration(thread_id); err != nil {
							mWorker.RenewDBConn()
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

	if mWorker.MType() == migrate.DELETION {
		mWorker.HandleLeftOverWaitingNodes()
	}

	mWorker.FinishMigration(mWorker.MType(), threads)
	return finished
}

func LThreadController(mWorker migrate.LMigrationWorker, threads int) bool {
	var wg sync.WaitGroup

	commitChannel := make(chan ThreadChannel)

	if threads != 0 {
		if !mWorker.RegisterMigration(mWorker.MType(), threads) {
			log.Fatal("Unable to register migration!")
		} else {
			log.Println("Migration registered:", mWorker.MType())
		}
	} else {
		threads = 1
	}

	for threadID := 0; threadID < threads; threadID++ {
		wg.Add(1)
		go func(thread_id int, commitChannel chan ThreadChannel) {
			defer wg.Done()
			switch mWorker.MType() {
			case migrate.DELETION:
				{
					for {
						if err := mWorker.DeletionMigration(mWorker.GetRoot(), thread_id); err != nil {
							mWorker.RenewDBConn()
							continue
						}
						break
					}
					for {
						if err := mWorker.SecondPhase(thread_id); err != nil {
							mWorker.RenewDBConn()
							continue
						}
						break
					}
				}
			case migrate.CONSISTENT:
				{
					for {
						if err := mWorker.ConsistentMigration(thread_id); err != nil {
							mWorker.RenewDBConn()
							continue
						}
						break
					}
				}
			case migrate.INDEPENDENT:
				{
					for {
						if err := mWorker.IndependentMigration(thread_id); err != nil {
							mWorker.RenewDBConn()
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

	if mWorker.MType() == migrate.DELETION {
		mWorker.HandleLeftOverWaitingNodes()
	}
	
	mWorker.FinishMigration(mWorker.MType(), threads)
	return finished
}

package mthread

import (
	"fmt"
	"log"
	"stencil/migrate"
	"sync"
	"time"
)

func ThreadController(mWorker migrate.MigrationWorker) bool {
	var wg sync.WaitGroup

	commitChannel := make(chan ThreadChannel)
	threads := 1

	if mWorker.MType() != migrate.DELETION {
		if !mWorker.RegisterMigration(mWorker.MType()) {
			log.Fatal("Unable to register migration!")
		} else {
			log.Println("Migration registered!")
		}
	}

	for threadID := 0; threadID < threads; threadID++ {
		time.Sleep(time.Millisecond * 300)
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
				}
			case migrate.CONSISTENT:
				{
					for {
						if err := mWorker.ConsistentMigration(mWorker.GetRoot(), thread_id); err != nil {
							mWorker.RenewDBConn()
							continue
						}
						break
					}
				}
			case migrate.INDEPENDENT:
				{
					for {
						if err := mWorker.IndependentMigration(mWorker.GetRoot(), thread_id); err != nil {
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
	if mWorker.MType() == migrate.DELETION {
		if bags, err := mWorker.GetUserBags(); err == nil && len(bags) > 0 {
			for ibag, bag := range bags {
				wg.Add(1)
				go func(thread_id int, commitChannel chan ThreadChannel) {
					defer wg.Done()
					for {
						if err := mWorker.MigrateProcessBags(bag); err != nil {
							mWorker.RenewDBConn()
							continue
						}
						break
					}
					commitChannel <- ThreadChannel{Finished: true, Thread_id: thread_id}
				}(ibag+threads, commitChannel)
			}
		}
	}

	go func() {
		wg.Wait()
		close(commitChannel)
	}()

	finished := true

	for threadResponse := range commitChannel {
		fmt.Println("THREAD FINISHED WORKING", threadResponse)
		if !threadResponse.Finished {
			finished = false
		}
	}

	if mWorker.MType() == migrate.DELETION {
		mWorker.HandleLeftOverWaitingNodes()
	}
	return finished
}

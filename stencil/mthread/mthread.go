package mthread

import (
	"fmt"
	"stencil/migrate"
	"sync"
	"time"
)

func ThreadController(mWorker migrate.MigrationWorker) bool {
	var wg sync.WaitGroup

	commitChannel := make(chan ThreadChannel)
	threads := 1

	for threadID := 1; threadID <= threads; threadID++ {
		time.Sleep(time.Millisecond * 300)
		wg.Add(1)
		go func(thread_id int, commitChannel chan ThreadChannel) {
			defer wg.Done()
			for {
				if err := mWorker.MigrateProcess(mWorker.GetRoot()); err != nil {
					mWorker.RenewDBConn()
					continue
				}
				mWorker.Finish()
				break
			}
			commitChannel <- ThreadChannel{Finished: true, Thread_id: thread_id}
		}(threadID, commitChannel)
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
	return finished
}

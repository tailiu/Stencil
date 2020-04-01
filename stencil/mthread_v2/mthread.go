package mthread_v2

import (
	"fmt"
	"log"
	"os"
	"stencil/config"
	"stencil/db"
	migrate "stencil/migrate_v2"
	"stencil/transaction"
	"strings"
	"time"

	"github.com/gookit/color"
	logg "github.com/withmandala/go-log"
)

// Init : Initializes the thread controller
func (mThread *MigrationThreadController) Init() {

	mThread.Logger = logg.New(os.Stderr)

	mThread.Logger.WithTimestamp()
	mThread.Logger.WithColor()
	mThread.Logger.WithDebug()

	if mThread.enableBags {
		mThread.Logger.Info("Bags Enabled")
	}

	if mThread.isBlade {
		mThread.Logger.Info("Destination App Server: Blade")
	}

	self.commitChannel = make(chan ThreadChannel)
	self.stencilDB = db.GetDBConn(db.STENCIL_DB)

	if txnID, err := db.CreateMigrationTransaction(self.stencilDB); err == nil {
		mThread.txnID = txnID
	} else {
		mThread.Logger.Fatal("@MigrationThreadController.Init > CreateMigrationTransaction | ", err)
	}

	if !db.RegisterMigration(mThread.uid, mThread.SrcAppInfo.ID, mThread.DstAppInfo.ID, mThread.mType, mThread.txnID, mThread.totalThreads, mThread.stencilDB, false) {
		mThread.Logger.Fatal("Unable to register migration!")
	} else {
		mThread.Logger.Info("Migration Registered of Type ", mThread.mType)
	}

	mThread.Logger.Info("Migration thread controller intialized!")
}

// CreateMigrationWorker : Creates a new migration worker
func (mThread MigrationThreadController) CreateMigrationWorker(uid, srcApp, srcAppID, dstApp, dstAppID string, logTxn *transaction.Log_txn, mtype string, mappings *config.MappedApp, threadID int, isBlade ...bool) MigrationWorkerV2 {

	srcAppConfig, err := config.CreateAppConfig(srcApp, srcAppID)
	if err != nil {
		log.Fatal(err)
	}

	dstAppConfig, err := config.CreateAppConfig(dstApp, dstAppID, isBlade...)
	if err != nil {
		log.Fatal(err)
	}

	mWorker := migrate.MigrationWorker{
		uid:          uid,
		SrcAppConfig: srcAppConfig,
		DstAppConfig: dstAppConfig,
		mappings:     mappings,
		wList:        WaitingList{},
		unmappedTags: CreateUnmappedTags(),
		SrcDBConn:    db.GetDBConn(srcApp),
		DstDBConn:    db.GetDBConn(dstApp, isBlade...),
		logTxn:       &transaction.Log_txn{DBconn: logTxn.DBconn, Txn_id: logTxn.Txn_id},
		mtype:        mtype,
		visitedNodes: make(map[string]map[string]bool),
		Logger:       logg.New(os.Stderr)}

	mWorker.Logger.WithTimestamp()
	mWorker.Logger.WithColor()
	mWorker.Logger.WithDebug()

	if err := mWorker.FetchRoot(threadID); err != nil {
		mWorker.Logger.Fatal(err)
	}

	mWorker.FTPClient = GetFTPClient()

	mWorker.Logger.Infof("Worker Created for thread: %v", threadID)

	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	return mWorker
}

// NewMigrationThread : Creates new migration thread
func (mThread *MigrationThreadController) NewMigrationThread() {

	mThread.Logger.Info("Attempting to create a new migration thread")

	mThread.waitGroup.Add(1)

	go func(thread_id int, commitChannel chan ThreadChannel) {
		defer mThread.waitGroup.Done()

		mWorker := mThread.CreateMigrationWorker(uid, srcApp, srcAppID, dstApp, dstAppID, logTxn, mtype, mappings, threadID, isBlade)

		defer mWorker.CloseDBConns()

		switch mWorker.MType() {
		case migrate.BAGS:
			{
				for {
					if err := mWorker.MigrateBags(thread_id, isBlade); err != nil {
						mWorker.Logger.Error("@ThreadControllerV2 > MigrateBags | Crashed with error: ", err)
						time.Sleep(time.Second * 5)
						continue
					}
					break
				}
			}
		case migrate.DELETION:
			{
				for {
					if err := mWorker.DeletionMigration(mWorker.GetRoot(), thread_id); err != nil {
						if !strings.Contains(err.Error(), "deadlock") {
							mWorker.RenewDBConn(isBlade)
						} else {
							fmt.Print(">>>>>>>>>>>>>>>>>>>>>>> RESTART AFTER DEADLOCK <<<<<<<<<<<<<<<<<<<<<<<<<<<")
						}
						continue
					}
					break
				}

				if enableBags {
					for {
						if err := mWorker.MigrateBags(thread_id, isBlade); err != nil {
							mWorker.Logger.Error("@ThreadControllerV2 > MigrateBags | Crashed with error: ", err)
							time.Sleep(time.Second * 5)
							continue
						}
						break
					}
				}
			}
		case migrate.CONSISTENT:
			{
				for {
					if err := mWorker.ConsistentMigration(thread_id); err != nil {
						if !strings.Contains(err.Error(), "deadlock") {
							mWorker.RenewDBConn(isBlade)
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
							mWorker.RenewDBConn(isBlade)
						} else {
							fmt.Print(">>>>>>>>>>>>>>>>>>>>>>> RESTART AFTER DEADLOCK <<<<<<<<<<<<<<<<<<<<<<<<<<<")
						}
						continue
					}
					break
				}
			}
		case migrate.NAIVE:
			{
				for {
					if err := mWorker.NaiveMigration(thread_id); err != nil {
						if !strings.Contains(err.Error(), "deadlock") {
							mWorker.RenewDBConn(isBlade)
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

func (mThread MigrationThread) ThreadController(uid, srcApp, srcAppID, dstApp, dstAppID string, logTxn *transaction.Log_txn, mtype string, mappings *config.MappedApp, threads int, isBlade bool, enableBags bool) bool {

	go func() {
		wg.Wait()
		close(commitChannel)
	}()

	finished := true

	var finished_threads []string
	for threadResponse := range commitChannel {
		color.Light.Println("THREAD FINISHED WORKING", threadResponse, strings.Join(finished_threads, ","))
		if !threadResponse.Finished {
			finished = false
		}
		finished_threads = append(finished_threads, fmt.Sprint(threadResponse.Thread_id))
	}

	return finished
}

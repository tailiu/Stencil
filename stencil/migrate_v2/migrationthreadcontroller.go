package migrate_v2

import (
	"fmt"
	"log"
	config "stencil/config/v2"
	"stencil/db"
	"stencil/helper"
	"stencil/transaction"
	"strings"
	"time"
)

// Init : Initializes the thread controller
func (mThread *MigrationThreadController) Init() {

	mThread.Logger = helper.CreateLogger(mThread.LoggerDebugFlag)

	if mThread.EnableBags {
		mThread.Logger.Info("Bags: Enabled")
	} else {
		mThread.Logger.Info("Bags: Disabled")
	}

	if mThread.FTPFlag {
		mThread.Logger.Info("File Transfer: Enabled")
	} else {
		mThread.Logger.Info("File Transfer: Disabled")
	}

	if mThread.Blade {
		mThread.Logger.Info("Destination App Server: Blade")
	} else {
		mThread.Logger.Info("Destination App Server:  Not Blade")
	}

	if len(mThread.DstAppInfo.Name) <= 0 || len(mThread.SrcAppInfo.Name) <= 0 || len(mThread.DstAppInfo.ID) == 0 || len(mThread.SrcAppInfo.ID) == 0 {
		mThread.Logger.Debug("Src App | ", mThread.SrcAppInfo)
		mThread.Logger.Debug("Dst App | ", mThread.DstAppInfo)
		mThread.Logger.Fatal("App Info(s) not set!")
	} else {
		mThread.Logger.Info("Src App | ", mThread.SrcAppInfo)
		mThread.Logger.Info("Dst App | ", mThread.DstAppInfo)
	}

	switch mThread.MType {
	case "d":
		{
			mThread.MType = DELETION
			mThread.Logger.Info("Migration Type: DELETION")
		}
	case "i":
		{
			mThread.MType = INDEPENDENT
			mThread.Logger.Info("Migration Type: INDEPENDENT")
		}
	case "c":
		{
			mThread.MType = CONSISTENT
			mThread.Logger.Info("Migration Type: CONSISTENT")
		}
	case "b":
		{
			mThread.MType = BAGS
			mThread.Logger.Info("Migration Type: BAGS")
		}
	case "n":
		{
			mThread.MType = NAIVE
			mThread.Logger.Info("Migration Type: NAIVE")
		}
	default:
		{
			mThread.Logger.Fatal("Wrong migration type specified: ", mThread.MType)
		}
	}

	mThread.commitChannel = make(chan ThreadChannel)
	mThread.stencilDB = db.GetDBConn(db.STENCIL_DB)

	if mappings := config.GetSchemaMappingsFor(mThread.SrcAppInfo.Name, mThread.DstAppInfo.Name); mappings == nil {
		mThread.Logger.Fatalf("Can't find mappings from '%s' to '%s'", mThread.SrcAppInfo.Name, mThread.DstAppInfo.Name)
	} else {
		mThread.mappings = *mappings
	}

	if mThread.Threads == 0 {
		mThread.Threads = 1
	}

	if txnID, err := db.CreateMigrationTransaction(mThread.stencilDB); err == nil {
		mThread.txnID = txnID
	} else {
		mThread.Logger.Fatal("Can't create migration transaction | ", err)
	}

	if !db.RegisterMigration(mThread.UID, fmt.Sprint(mThread.SrcAppInfo.ID), fmt.Sprint(mThread.DstAppInfo.ID), mThread.MType, mThread.txnID, mThread.Threads, mThread.stencilDB, false) {
		mThread.Logger.Fatal("Unable to register migration!")
	} else {
		mThread.Logger.Info("Migration Registered of Type: ", mThread.MType)
	}

	mThread.Logger.Info("Migration thread controller intialized!")

	fmt.Print("========================================================================\n\n")
}

// CreateMigrationWorker : Creates and returns a new migration worker
func (mThread *MigrationThreadController) CreateMigrationWorker(threadID int) MigrationWorker {

	mThread.Logger.Info("Creating a new migration worker for thread: ", threadID)

	srcAppConfig, err := config.CreateAppConfig(mThread.SrcAppInfo.Name, fmt.Sprint(mThread.SrcAppInfo.ID))
	if err != nil {
		log.Fatal(err)
	}

	dstAppConfig, err := config.CreateAppConfig(mThread.DstAppInfo.Name, fmt.Sprint(mThread.DstAppInfo.ID), mThread.Blade)
	if err != nil {
		log.Fatal(err)
	}

	mWorker := MigrationWorker{
		uid:          mThread.UID,
		threadID:     threadID,
		SrcAppConfig: srcAppConfig,
		DstAppConfig: dstAppConfig,
		mappings:     mThread.mappings,
		logTxn:       &transaction.Log_txn{DBconn: db.GetDBConn(db.STENCIL_DB), Txn_id: mThread.txnID},
		mtype:        mThread.MType,
		visitedNodes: VisitedNodes{},
		mThread:      mThread,
		FTPFlag:      mThread.FTPFlag,
		Logger:       helper.CreateLogger(mThread.LoggerDebugFlag)}

	mWorker.visitedNodes.Init()

	if err := mWorker.FetchRoot(threadID); err != nil {
		mWorker.Logger.Fatal(err)
	}

	if mWorker.FTPFlag {
		mWorker.FTPClient = GetFTPClient()
	}

	mThread.Logger.Info("Migration worker created!")
	fmt.Print("========================================================================\n\n")

	return mWorker
}

// CreateBagWorker : Creates and returns a new migration worker for data bags
func (mThread *MigrationThreadController) CreateBagWorker(uid, srcAppID, dstAppID string, threadID int) MigrationWorker {

	mThread.Logger.Infof("Creating a new bag worker for thread: %d | uid: %s, srcApp: %s, dstApp: %s \n", threadID, uid, srcAppID, dstAppID)

	srcApp, err := db.GetAppNameByAppID(mThread.stencilDB, srcAppID)
	if err != nil {
		log.Fatal(err)
	}

	dstApp, err := db.GetAppNameByAppID(mThread.stencilDB, dstAppID)
	if err != nil {
		log.Fatal(err)
	}

	srcAppConfig, err := config.CreateAppConfig(srcApp, srcAppID)
	if err != nil {
		log.Fatal(err)
	}

	dstAppConfig, err := config.CreateAppConfig(dstApp, dstAppID, mThread.Blade)
	if err != nil {
		log.Fatal(err)
	}

	var mappings *config.MappedApp

	if srcAppID == dstAppID {
		mappings = config.GetSelfSchemaMappings(mThread.stencilDB, srcAppID, srcApp)
	} else if srcApp == mThread.SrcAppInfo.Name && dstApp == mThread.DstAppInfo.Name {
		mappings = &mThread.mappings
	} else {
		mappings = config.GetSchemaMappingsFor(srcAppConfig.AppName, dstAppConfig.AppName)
		if mappings == nil {
			mThread.Logger.Fatalf("Can't find mappings from '%s' to '%s'", srcAppConfig.AppName, dstAppConfig.AppName)
		}
	}

	mWorker := MigrationWorker{
		uid:           uid,
		SrcAppConfig:  srcAppConfig,
		DstAppConfig:  dstAppConfig,
		mappings:      *mappings,
		logTxn:        &transaction.Log_txn{DBconn: db.GetDBConn(db.STENCIL_DB), Txn_id: mThread.txnID},
		mtype:         BAGS,
		processedBags: ProcessedBags{},
		FTPFlag:       mThread.FTPFlag,
		Logger:        helper.CreateLogger(mThread.LoggerDebugFlag)}

	mWorker.processedBags.Init()

	if mWorker.FTPFlag {
		mWorker.FTPClient = GetFTPClient()
	}

	mThread.Logger.Infof("Bag worker created for thread: %d | uid: %s, srcApp: %s, dstApp: %s \n", threadID, uid, srcAppID, dstAppID)
	fmt.Print("========================================================================\n\n")

	return mWorker
}

// NewMigrationThread : Creates new migration thread
func (mThread *MigrationThreadController) NewMigrationThread() {

	mThread.currentThreads++
	newThreadID := mThread.currentThreads

	mThread.Logger.Infof("Creating a new migration thread of type: %v, ID: %v \n", mThread.MType, newThreadID)

	defer mThread.waitGroup.Done()

	mWorker := mThread.CreateMigrationWorker(newThreadID)
	defer mWorker.CloseDBConns()

	switch mThread.MType {
	case BAGS:
		{
			for {
				if err := mWorker.BagsMigration(newThreadID); err != nil {
					mWorker.Logger.Error("NewMigrationThread : MigrateBags | Crashed with error: ", err)
					time.Sleep(time.Second * 5)
					continue
				}
				break
			}
		}
	case DELETION:
		{
			for {
				if err := mWorker.DeletionMigration(mWorker.Root, newThreadID); err != nil {
					mWorker.Logger.Error("NewMigrationThread : DeletionMigration | Crashed with error: ", err)
					mWorker.RenewDBConn(mThread.Blade)
					time.Sleep(time.Second * 5)
					continue
				}
				break
			}

			if mThread.EnableBags {
				for {
					if err := mWorker.BagsMigration(newThreadID); err != nil {
						mWorker.Logger.Error("NewMigrationThread : DeletionMigration > MigrateBags | Crashed with error: ", err)
						mWorker.RenewDBConn(mThread.Blade)
						time.Sleep(time.Second * 5)
						continue
					}
					break
				}
			}
		}
	case CONSISTENT:
		{
			for {
				if err := mWorker.ConsistentMigration(newThreadID); err != nil {
					mWorker.Logger.Error("NewMigrationThread : ConsistentMigration | Crashed with error: ", err)
					mWorker.RenewDBConn(mThread.Blade)
					time.Sleep(time.Second * 5)
					continue
				}
				break
			}
		}
	case INDEPENDENT:
		{
			for {
				if err := mWorker.IndependentMigration(newThreadID); err != nil {
					mWorker.Logger.Error("NewMigrationThread : IndependentMigration | Crashed with error: ", err)
					mWorker.RenewDBConn(mThread.Blade)
					time.Sleep(time.Second * 5)
					continue
				}
				break
			}
		}
	case NAIVE:
		{
			for {
				if err := mWorker.NaiveMigration(newThreadID); err != nil {
					mWorker.Logger.Error("NewMigrationThread : NaiveMigration | Crashed with error: ", err)
					mWorker.RenewDBConn(mThread.Blade)
					time.Sleep(time.Second * 5)
					continue
				}
				break
			}
		}
	}

	mThread.commitChannel <- ThreadChannel{finished: true, threadID: newThreadID, size: mWorker.Size}
}

// Run : Start migration threads
func (mThread *MigrationThreadController) Run() error {

	for i := 0; i < mThread.Threads; i++ {
		mThread.waitGroup.Add(1)
		go mThread.NewMigrationThread()
	}

	go func() {
		mThread.waitGroup.Wait()
		close(mThread.commitChannel)
	}()

	var finishedThreads []string

	for threadResponse := range mThread.commitChannel {
		finishedThreads = append(finishedThreads, fmt.Sprint(threadResponse.threadID))
		mThread.Logger.Infof("MIGRATION THREAD # %v FINISHED WORKING!", threadResponse.threadID)
		mThread.Logger.Info("Finished Migration Thread IDs | ", strings.Join(finishedThreads, ","))
		mThread.size += threadResponse.size
	}

	return nil
}

// Stop : Close db conns and commit migration transaction
func (mThread *MigrationThreadController) Stop() {
	if err := db.MigrationTransactionLogOutcome(mThread.stencilDB, mThread.txnID, "COMMIT"); err != nil {
		mThread.Logger.Fatal(err)
	}
	if finished := db.FinishMigration(mThread.stencilDB, mThread.txnID, mThread.size); !finished {
		mThread.Logger.Fatal("DB error in FinishMigration")
	}
	mThread.stencilDB.Close()

	mThread.Logger.Info("Migration Finished!")
	fmt.Print("========================================================================\n\n")
}

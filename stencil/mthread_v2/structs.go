package mthread_v2

import (
	"database/sql"
	"sync"

	logg "github.com/withmandala/go-log"
)

type App struct {
	Name string
	ID   string
}

type ThreadChannel struct {
	Finished  bool
	Thread_id int
	size      int
}

// MigrationThreadController : SrcAppInfo and DstAppInfo just contain strings of App ids and names.
// Actual App Configs will be created separately for each migration thread.
type MigrationThreadController struct {
	uid            string
	waitGroup      sync.WaitGroup
	commitChannel  chan ThreadChannel
	enableBags     bool
	isBlade        bool
	totalThreads   int
	currentThreads int
	txnID          int
	stencilDB      *sql.DB
	Logger         *logg.Logger
	SrcAppInfo     App
	DstAppInfo     App
	mType          string
}

type MigrationThread struct {
}

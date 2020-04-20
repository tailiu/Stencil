package migrate_v1

import (
	"database/sql"
	"stencil/config"
	"stencil/transaction"
	"sync"

	logg "github.com/withmandala/go-log"
)

const (
	INDEPENDENT = "0"
	CONSISTENT  = "1"
	DELETION    = "3"
	BAGS        = "4"
	NAIVE       = "5"
)

type Transactions struct {
	SrcTx     *sql.Tx
	DstTx     *sql.Tx
	StencilTx *sql.Tx
}

type DependencyNode struct {
	Tag  config.Tag
	SQL  string
	Data map[string]interface{}
	PKs  []int64
}

type WaitingNode struct {
	ContainedNodes []*DependencyNode
	// LookingFor    []config.Tag
	LookingFor map[string]map[string]interface{}
}

type WaitingList struct {
	Nodes []*WaitingNode
}

type InvalidList struct {
	Nodes []*DependencyNode
}

type UnmappedTags struct {
	Mutex *sync.Mutex
	tags  []string
}

type MigrationWorker struct {
	uid          string
	SrcAppConfig config.AppConfig
	DstAppConfig config.AppConfig
	mappings     *config.MappedApp
	wList        WaitingList
	unmappedTags UnmappedTags
	root         *DependencyNode
	DBConn       *sql.DB
	logTxn       *transaction.Log_txn
	mtype        string
	visitedNodes map[string]map[string]bool
	arg          string
	Size         int
	Logger       *logg.Logger
	// threadID     int
}

type ThreadChannel struct {
	Finished  bool
	Thread_id int
	size      int
}

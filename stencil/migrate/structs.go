package migrate

import (
	"database/sql"
	"stencil/config"
	"stencil/transaction"
	"sync"

	"github.com/jlaffaye/ftp"
)

const (
	INDEPENDENT = "0"
	CONSISTENT  = "1"
	DELETION    = "3"
	BAGS        = "4"
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

type MappedData struct {
	cols        string
	vals        string
	orgCols     string
	orgColsLeft string
	srcTables   map[string]bool
	ivals       []interface{}
	undoAction  *transaction.UndoAction
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
	// threadID     int
}

type LMigrationWorker struct {
	uid          string
	SrcAppConfig config.AppConfig
	DstAppConfig config.AppConfig
	mappings     *config.MappedApp
	wList        WaitingList
	unmappedTags UnmappedTags
	root         *DependencyNode
	SrcDBConn    *sql.DB
	DstDBConn    *sql.DB
	logTxn       *transaction.Log_txn
	mtype        string
	visitedNodes map[string]bool
	FTPClient    *ftp.ServerConn
	// threadID     int
}

type MigrationWorkerV2 struct {
	uid          string
	SrcAppConfig config.AppConfig
	DstAppConfig config.AppConfig
	mappings     *config.MappedApp
	wList        WaitingList
	unmappedTags UnmappedTags
	root         *DependencyNode
	SrcDBConn    *sql.DB
	DstDBConn    *sql.DB
	logTxn       *transaction.Log_txn
	mtype        string
	visitedNodes map[string]bool
	FTPClient    *ftp.ServerConn
	tx           Transactions
	// threadID     int
}

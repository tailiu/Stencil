package migrate_v1

import (
	"database/sql"
	"stencil/config"
	"stencil/transaction"
	"sync"

	"github.com/jlaffaye/ftp"
	logg "github.com/withmandala/go-log"
)

const (
	INDEPENDENT = "0"
	CONSISTENT  = "1"
	DELETION    = "3"
	BAGS        = "4"
	NAIVE       = "5"
)

type ValueWithReference struct {
	value interface{}
	ref   *MappingRef
}

type App struct {
	Name string
	ID   int64
}

type Member struct {
	Name string
	ID   int64
}

type AttrRow struct {
	FromApp    App
	FromMember Member
	FromAttr   int64
	FromVal    string
	ToApp      App
	ToMember   Member
	ToAttr     int64
	ToVal      string
}

type IDRow struct {
	FromAppName  string
	FromAppID    string
	FromMember   string
	FromMemberID string
	FromID       int64
	ToAppID      string
	ToAppName    string
	ToMember     string
	ToMemberID   string
	ToID         int64
}

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

type MappingRef struct {
	appID         string
	fromMember    string
	fromAttr      string
	fromVal       string
	toVal         string
	toMember      string
	toAttr        string
	mergedFromBag bool
}

type MappedData struct {
	cols        string
	vals        string
	orgCols     string
	orgColsLeft string
	srcTables   map[string][]string
	ivals       []interface{}
	undoAction  *transaction.UndoAction
	refs        []MappingRef
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
	Logger       *logg.Logger
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
	visitedNodes map[string]map[string]bool
	FTPClient    *ftp.ServerConn
	tx           Transactions
	Logger       *logg.Logger
	// threadID     int
}

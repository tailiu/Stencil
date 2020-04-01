package migrate_v2

import (
	"database/sql"
	config "stencil/config/v2"
	"stencil/transaction"
	"sync"

	"github.com/jlaffaye/ftp"
	logg "github.com/withmandala/go-log"
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
	mtype        string
	threadID     int
	visitedNodes VisitedNodes
	refCreator   ReferenceCreator
	SrcAppConfig config.AppConfig
	DstAppConfig config.AppConfig
	mappings     config.MappedApp
	Root         *DependencyNode
	logTxn       *transaction.Log_txn
	FTPClient    *ftp.ServerConn
	tx           Transactions
	Logger       *logg.Logger
	Size         int
	mThread      *MigrationThreadController
}

type VisitedNodes struct {
	Nodes map[string]map[string]bool
}

type ReferenceCreator struct {
}

// ThreadChannel : Channel for thread communication
type ThreadChannel struct {
	finished bool
	threadID int
	size     int
}

// MigrationThreadController : SrcAppInfo and DstAppInfo just contain strings of App ids and names.
// Actual App Configs will be created separately for each migration thread.
type MigrationThreadController struct {
	UID            string
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
	MType          string
	mappings       config.MappedApp
	size           int
}

type MigrationThread struct {
}

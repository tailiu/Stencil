package migrate_v2

import (
	"database/sql"
	config "stencil/config/v2"
	"stencil/transaction"
	"sync"

	"github.com/jlaffaye/ftp"
	logg "github.com/withmandala/go-log"
)

// DataMap : type representation of data fetched from db
type DataMap map[string]interface{}

// ValueWithReference : Stores fetched value and reference (if created)
type ValueWithReference struct {
	value interface{}
	ref   *MappingRef
}

// App : Data struct
type App struct {
	Name string
	ID   string
}

type Bag struct {
	AppID    string
	ID       string
	UID      string
	PK       string
	MemberID string
	Member   string
	Node     *DependencyNode
}

// Member : Data struct
type Member struct {
	Name string
	ID   string
}

type AttrRow struct {
	FromAppName  string
	FromAppID    string
	FromMember   string
	FromMemberID string
	FromID       int64
	FromAttr     string
	FromVal      string
	ToAppID      string
	ToAppName    string
	ToMember     string
	ToMemberID   string
	ToID         int64
	ToAttr       string
	ToVal        string
}

type IDRow__defunct struct {
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
	Data DataMap
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
	fromMemberID  string
	fromMember    string
	fromAttr      string
	fromAttrID    string
	fromID        int64
	fromVal       string
	toVal         string
	toMemberID    string
	toMember      string
	toAttr        string
	toAttrID      string
	mergedFromBag bool
}

type MappedData__defunct struct {
	cols        string
	vals        string
	orgCols     string
	orgColsLeft string
	srcTables   map[string][]string
	ivals       []interface{}
	undoAction  *transaction.UndoAction
	refs        []MappingRef
}

type MappedMemberData struct {
	ToID       string
	ToAppID    string
	ToMemberID string
	ToMember   string
	Data       map[string]MappedMemberValue
	DBConn     *sql.DB
}

type MappedMemberValue struct {
	IsInput      bool
	IsMethod     bool
	IsExpression bool
	AppID        string
	FromID       string
	FromMemberID string
	FromMember   string
	FromAttr     string
	FromAttrID   string
	ToID         string
	Value        interface{}
	Ref          *MappingRef
	Logger       *logg.Logger
	DBConn       *sql.DB
}

type MigrationWorker struct {
	uid            string
	mtype          string
	threadID       int
	visitedNodes   VisitedNodes
	processedBags  ProcessedBags
	refCreator     ReferenceCreator
	SrcAppConfig   config.AppConfig
	DstAppConfig   config.AppConfig
	mappings       config.MappedApp
	Root           *DependencyNode
	logTxn         *transaction.Log_txn
	FTPClient      *ftp.ServerConn
	tx             Transactions
	Logger         *logg.Logger
	Size           int
	mThread        *MigrationThreadController
	FTPFlag        bool
	DeleteRootFlag bool
}

type VisitedNodes struct {
	Nodes map[string]map[string]bool
}

type ProcessedBags struct {
	Bags map[string]bool
}

type BagManager struct {
	Bags        map[string]*DBBag
	VisitedRows map[string]bool
	PrevUIDs    map[string]string
	DBConn      *sql.DB
}

type DBBag struct {
	Data          DataMap
	PK            string
	ID            string
	UID           string
	AppID         string
	MemberID      string
	TxnID         string
	AttrsToRemove []string
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
	UID             string
	waitGroup       sync.WaitGroup
	commitChannel   chan ThreadChannel
	EnableBags      bool
	Blade           bool
	Threads         int
	currentThreads  int
	txnID           int
	stencilDB       *sql.DB
	Logger          *logg.Logger
	SrcAppInfo      App
	DstAppInfo      App
	MType           string
	mappings        config.MappedApp
	size            int
	LoggerDebugFlag bool
	FTPFlag         bool
	DeleteRootFlag  bool
}

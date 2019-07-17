package migrate

import (
	"database/sql"
	"stencil/config"
	"stencil/transaction"
	"sync"
)

const (
	INDEPENDENT = "0"
	CONSISTENT  = "1"
	DELETION    = "3"
)

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

type MigrationWorker struct {
	uid          string
	srcAppConfig config.AppConfig
	dstAppConfig config.AppConfig
	mappings     *config.MappedApp
	wList        WaitingList
	unmappedTags UnmappedTags
	root         *DependencyNode
	dbConn       *sql.DB
	logTxn       *transaction.Log_txn
	mtype        string
	// threadID     int
}

package evaluation

import (
	"database/sql"
	// "stencil/db"
	// "time"
)

const (
	INDEPENDENT = "0"
	CONSISTENT = "1"
	DELETION = "3"
	logDir = "./evaluation/logs/"
	logCounterDir = "./evaluation/logs_counter/"
)

var mediaSize = map[string]int64 {
	"1.jpg": 512017,
	"2.jpg": 206993,
	"3.jpg": 102796,
	"4.jpg": 51085,
	"5.jpg": 1033414,
}

var appMediaTables = map[string]string {
	"diaspora": "photos",
	"mastodon": "media_attachments",
	"gnusocial": "file",
}

var mediaTables = map[string]string {
	"31": "photos",
	"76": "media_attachments",
	"124": "file",
}

// These databases are default databases if I don't set them in experiments
var	stencilDB = "stencil_exp"
var	stencilDB1 = "stencil_exp1"
var	stencilDB2 = "stencil_exp2"
	
var	mastodon = "mastodon_exp"
var	mastodon1 = "mastodon_exp1"
var	mastodon2 = "mastodon_exp2"

var diaspora = "diaspora_test"
var diaspora1 = "diaspora_test"

var twitter = "twitter_test"
var twitter1 = "twitter_test"

var gnusocial = "gnusocial_test"
var gnusocial1 = "gnusocial_test"

var dependencies = map[string]map[string][]string {
	"diaspora" : map[string][]string {
		"posts": []string {
			"root_guid:posts.guid"},
		"comments": []string {
			"commentable_id:posts.id"},
		"likes": []string {
			"target_id:posts.id"},
		"messages": []string {
			"conversation_id:conversations.id"}},
	"mastodon" : map[string][]string {
		"statuses": []string {
			"reblog_of_id:statuses.id",
			"conversation_id:conversations.id",
			"in_reply_to_id:statuses.id"},
		"favourites": []string {
			"status_id:statuses.id"}}}

type EvalConfig struct {
	Dependencies map[string]map[string][]string
	StencilDBConn *sql.DB
	StencilDBConn1 *sql.DB
	StencilDBConn2 *sql.DB
	MastodonDBConn *sql.DB
	MastodonDBConn1 *sql.DB
	MastodonDBConn2 *sql.DB
	DiasporaDBConn *sql.DB
	TwitterDBConn *sql.DB
	GnusocialDBConn *sql.DB
	DiasporaDBConn1 *sql.DB
	TwitterDBConn1 *sql.DB
	GnusocialDBConn1 *sql.DB
	TableIDNamePairs map[string]string
	AttrNameIDPairsOfApps map[string]map[string]string
	MastodonTableNameIDPairs map[string]string
	DiasporaTableNameIDPairs map[string]string
	MastodonAppID string
	DiasporaAppID string
	AllAppNameIDs map[string]string
	SrcAnomaliesVsMigrationSizeFile string
	DstAnomaliesVsMigrationSizeFile string
	InterruptionDurationFile string
	MigrationRateFile string
	MigratedDataSizeFile string
	MigrationTimeFile string
	SrcDanglingDataInSystemFile string
	DstDanglingDataInSystemFile string
	DataDowntimeInStencilFile string
	DataDowntimeInNaiveFile string
	DataBags string
	MigratedDataSizeByDstFile string
	MigrationTimeByDstFile string
	MigratedDataSizeBySrcFile string
	MigrationTimeBySrcFile string
	DanglingDataFile string
	DanglingObjectsFile string
	Diaspora1KCounterFile string
	Diaspora10KCounterFile string
	Diaspora100KCounterFile string
	Diaspora1MCounterFile string
	DataDowntimeInPercentageInStencilFile string
	DataDowntimeInPercentageInNaiveFile string
}

type DataBagData struct {
	TableID 	string
	RowIDs 		[]string
}

type DisplayedData struct {
	TableID 	string
	RowIDs 		[]string
}

type Counter struct {
	Edges	int	`json:"edges"`
	Nodes	int	`json:"nodes"`
	UserID	int	`json:"userID"`
}

type SA1SizeStruct struct {
	Size	string `json:"size"`
	UserID	string `json:"userID"`
}

type ScalabilityDataStruct struct {
	DisplayTime 			string `json:"displayTime"`
	Edges					string `json:"edges"`
	EdgesAfterMigration		string `json:"edgesAfterMigration"`
	MigrationTime 			string `json:"migrationTime"`
	Nodes					string `json:"nodes"`
	NodesAfterMigration		string `json:"nodesAfterMigration"`
	PersonID 				string `json:"person_id"`
}
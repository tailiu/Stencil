package evaluation

import (
	"database/sql"
	// "stencil/db"
	// "time"
)

var mediaSize = map[string]int64 {
	"1.jpg": 512017,
	"2.jpg": 206993,
	"3.jpg": 102796,
	"4.jpg": 51085,
	"5.jpg": 1033414,
}

const logDir = "./evaluation/logs/"
const logCounterDir = "./evaluation/logs_counter/"

var	stencilDB = "stencil_exp"
var	stencilDB1 = "stencil_exp1"
var	stencilDB2 = "stencil_exp2"
	
var	mastodon = "mastodon_exp"
var	mastodon1 = "mastodon_exp1"
var	mastodon2 = "mastodon_exp2"

var diaspora = "diaspora_1000000"

const (

	INDEPENDENT = "0"
	CONSISTENT = "1"
	DELETION = "3"
)

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
	TableIDNamePairs map[string]string
	MastodonTableNameIDPairs map[string]string
	DiasporaTableNameIDPairs map[string]string
	MastodonAppID string
	DiasporaAppID string
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
	Diaspora1KCounterFile string
	Diaspora10KCounterFile string
	Diaspora100KCounterFile string
}

type DataBagData struct {
	TableID 	string
	RowIDs 		[]string
}

type DisplayedData struct {
	TableID 	string
	RowIDs 		[]string
}

// Messages will be handled in special ways
// var dependencies = map[string]map[string][]string {
// 	"diaspora" : map[string][]string {
// 		"posts": []string {
// 			"root_guid:posts.guid"},
// 		"comments": []string {
// 			"commentable_id:posts.id"},
// 		"likes": []string {
// 			"target_id:posts.id"},
// 		"messages": []string {
// 			"conversation_id:conversations.id"}},
// 	"mastodon" : map[string][]string {
// 		"statuses": []string {
// 			"reblog_of_id:statuses.id",
// 			"conversation_id:conversations.id"},
// 		"comments": []string {
// 			"conversation_id:conversations.id",
// 			"in_reply_to_id:statuses.id"},
// 		"messages": []string {
// 			"conversation_id:conversations.id",
// 			"in_reply_to_id:messages.id"},
// 		"favourites": []string {
// 			"status_id:statuses.id",
// 			"status_id:comments.id"}}}
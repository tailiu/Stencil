package evaluation

import (
	"database/sql"
	"stencil/db"
	// "time"
)

const (
	stencilDB = "stencil"
	mastodon = "mastodon"
	diaspora = "diaspora"

	INDEPENDENT = "0"
	CONSISTENT = "1"
	DELETION = "3"
)

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
	MastodonDBConn *sql.DB
	DiasporaDBConn *sql.DB
	MastodonAppID string
	DiasporaAppID string
	SrcAnomaliesVsMigrationSizeFile string
	DstAnomaliesVsMigrationSizeFile string
	InterruptionDurationFile string
}

func InitializeEvalConfig() *EvalConfig {
	evalConfig := new(EvalConfig)
	evalConfig.StencilDBConn = db.GetDBConn(stencilDB)
	evalConfig.MastodonDBConn = db.GetDBConn2(mastodon)
	evalConfig.DiasporaDBConn = db.GetDBConn(diaspora)
	evalConfig.MastodonAppID = db.GetAppIDByAppName(evalConfig.StencilDBConn, mastodon)
	evalConfig.DiasporaAppID = db.GetAppIDByAppName(evalConfig.StencilDBConn, diaspora)
	evalConfig.Dependencies = dependencies

	// t := time.Now()
	evalConfig.SrcAnomaliesVsMigrationSizeFile, 
	evalConfig.DstAnomaliesVsMigrationSizeFile, 
	evalConfig.InterruptionDurationFile = 
		"srcAnomaliesVsMigrationSize",// + t.String(), 
		"dstAnomaliesVsMigrationSize",// + t.String(),
		"interruptionDuration"// + t.String()

	return evalConfig
}
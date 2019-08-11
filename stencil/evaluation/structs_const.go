package evaluation

import (
	"database/sql"
	"stencil/db"
)

const (
	stencilDB = "stencil"
	mastodon = "mastodon"
	diaspora = "diaspora"
	old_mastodon = "old_mastodon"
)

// Messages will be handled in special ways
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
			"conversation_id:conversations.id"},
		"comments": []string {
			"conversation_id:conversations.id",
			"in_reply_to_id:statuses.id",
			"in_reply_to_id:comments.id"},
		"messages": []string {
			"conversation_id:conversations.id",
			"in_reply_to_id:messages.id"},
		"favourites": []string {
			"status_id:statuses.id",
			"status_id:comments.id",
			"status_id:messages.id"}}}

type EvalConfig struct {
	Dependencies map[string]map[string][]string
	StencilDBConn *sql.DB
	MastodonDBConn *sql.DB
	OldMastodonDBConn *sql.DB
	DiasporaDBConn *sql.DB
	MastodonAppID string
	DiasporaAppID string
}

// type DstViolateStats struct {
//     Messages.conversation_id:conversations.id  int `json:"messages.conversation_id:conversations.id"`
// }

// type DstDepNotMigratedStats struct {
//     Number int    `json:"number"`
//     Title  string `json:"title"`
// }


func InitializeEvalConfig() *EvalConfig {
	evalConfig := new(EvalConfig)
	evalConfig.StencilDBConn = db.GetDBConn(stencilDB)
	evalConfig.MastodonDBConn = db.GetDBConn(mastodon)
	evalConfig.DiasporaDBConn = db.GetDBConn(diaspora)
	evalConfig.OldMastodonDBConn = db.GetDBConn(old_mastodon)
	evalConfig.MastodonAppID = db.GetAppIDByAppName(evalConfig.StencilDBConn, mastodon)
	evalConfig.DiasporaAppID = db.GetAppIDByAppName(evalConfig.StencilDBConn, diaspora)
	evalConfig.Dependencies = dependencies

	return evalConfig
}
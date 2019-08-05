package evaluation

import (
	"database/sql"
	"stencil/db"
)

const (
	stencilDB = "stencil"
	mastodon = "mastodon"
	diaspora = "diaspora"
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
			"conversation_id:conversations.id"}}}

type EvalConfig struct {
	Dependencies map[string]map[string][]string
	StencilDBConn *sql.DB
	MastodonDBConn *sql.DB
	DiasporaDBConn *sql.DB
	MastodonAppID string
	DiasporaAppID string
}

func InitializeEvalConfig() *EvalConfig {
	evalConfig := new(EvalConfig)
	evalConfig.StencilDBConn = db.GetDBConn(stencilDB)
	evalConfig.MastodonDBConn = db.GetDBConn(mastodon)
	evalConfig.DiasporaDBConn = db.GetDBConn(diaspora)
	evalConfig.MastodonAppID = db.GetAppIDByAppName(evalConfig.StencilDBConn, mastodon)
	evalConfig.DiasporaAppID = db.GetAppIDByAppName(evalConfig.StencilDBConn, diaspora)
	evalConfig.Dependencies = dependencies

	return evalConfig
}
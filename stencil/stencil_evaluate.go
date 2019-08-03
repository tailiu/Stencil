package main

import (
	"stencil/evaluation"
	"stencil/db"
	// "stencil/config"
	// "stencil/common_phy_funcs"
	"database/sql"
	"log"
	"strconv"
)

const (
	stencilDB = "stencil"
	mastodon = "mastodon"
	diaspora = "diaspora"
)

type EvalConfig struct {
	stencilDBConn *sql.DB
	mastodonDBConn *sql.DB
	diasporaDBConn *sql.DB
	mastodonAppID string
	diasporaAppID string
}

func anomaliesVsSize(evalConfig *EvalConfig) {
	filterConditions := "and user_id >= 1003 and user_id < 10000"
	for _, dstMigrationID := range evaluation.GetAllMigrationIDsOfAppWithConds(evalConfig.stencilDBConn, evalConfig.mastodonAppID, filterConditions) {
		migrationID := strconv.FormatInt(dstMigrationID["migration_id"].(int64), 10)
		log.Println(migrationID)
		// totalDataSize := evaluation.GetTotalDataSize(evalConfig.stencilDBConn, evalConfig.diasporaDBConn, migrationID)
		// log.Println(totalDataSize)
		evaluation.GetLeftoverDataSize(evalConfig.stencilDBConn, evalConfig.diasporaDBConn, migrationID)

		// evaluation.GetPartiallyMappedRowTotalDataSize(evalConfig.stencilDBConn, evalConfig.mastodonAppID, dstMigrationID)
		// evaluation.GetPartiallyMappedRowDataSize(evalConfig.stencilDBConn, evalConfig.mastodonAppID, dstMigrationID)
	}
}

func initialize() *EvalConfig {
	evalConfig := new(EvalConfig)
	evalConfig.stencilDBConn = db.GetDBConn(stencilDB)
	evalConfig.mastodonDBConn = db.GetDBConn(mastodon)
	evalConfig.diasporaDBConn = db.GetDBConn(diaspora)
	evalConfig.mastodonAppID = db.GetAppIDByAppName(evalConfig.stencilDBConn, mastodon)
	evalConfig.diasporaAppID = db.GetAppIDByAppName(evalConfig.stencilDBConn, diaspora)

	return evalConfig
}

func main() {
	evalConfig := initialize()
	// mastodonConfig, _:= config.CreateAppConfig(mastodon, mastodonAppID)
	// common_phy_funcs.GetRowFromRowIDandTable(stencilDBConn, &mastodonConfig, "1008062662", "comments")

	anomaliesVsSize(evalConfig)
}
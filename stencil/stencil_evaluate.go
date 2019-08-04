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
	leftoverVsMigratedFile = "leftoverVsMigrated"
)

type EvalConfig struct {
	stencilDBConn *sql.DB
	mastodonDBConn *sql.DB
	diasporaDBConn *sql.DB
	mastodonAppID string
	diasporaAppID string
}

func leftoverVsMigrated(evalConfig *EvalConfig) {
	var data []float64
	// Need to be changed once data is ready to use
	filterConditions := "and migration_id = 734616546"
	
	for _, dstMigrationID := range evaluation.GetAllMigrationIDsOfAppWithConds(evalConfig.stencilDBConn, evalConfig.mastodonAppID, filterConditions) {
		migrationID := strconv.FormatInt(dstMigrationID["migration_id"].(int64), 10)		
		log.Println(migrationID)
		migratedDataSize := evaluation.GetMigratedDataSize(evalConfig.stencilDBConn, evalConfig.diasporaDBConn, evalConfig.diasporaAppID, migrationID)
		log.Println("Migrated data size: %d", migratedDataSize)
		leftoverDataSize := evaluation.GetLeftoverDataSize(evalConfig.stencilDBConn, evalConfig.diasporaDBConn, evalConfig.diasporaAppID, migrationID)
		log.Println("Leftover data size: %d", leftoverDataSize)
		percentageOfLeftoverData := float64(leftoverDataSize) / (float64(migratedDataSize) + float64(leftoverDataSize))
		log.Println("Percentage of leftover data size: %f", percentageOfLeftoverData)
		data = append(data, percentageOfLeftoverData)

		// evaluation.GetPartiallyMappedRowTotalDataSize(evalConfig.stencilDBConn, evalConfig.mastodonAppID, dstMigrationID)
		// evaluation.GetPartiallyMappedRowDataSize(evalConfig.stencilDBConn, evalConfig.mastodonAppID, dstMigrationID)
	}
	evaluation.WriteToLog(leftoverVsMigratedFile, evaluation.ConvertFloat64ToString(data))
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

	leftoverVsMigrated(evalConfig)
}
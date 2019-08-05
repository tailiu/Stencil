package main

import (
	"stencil/evaluation"
	"log"
	"strconv"
)

const (
	leftoverVsMigratedFile = "leftoverVsMigrated"
)

func leftoverVsMigrated(evalConfig *evaluation.EvalConfig) {
	var data []float64
	// Need to be changed once data is ready to use
	filterConditions := "and migration_id = 734616546"

	for _, dstMigrationID := range evaluation.GetAllMigrationIDsOfAppWithConds(evalConfig.StencilDBConn, evalConfig.MastodonAppID, filterConditions) {
		migrationID := strconv.FormatInt(dstMigrationID["migration_id"].(int64), 10)		
		log.Println(migrationID)
		migratedDataSize := evaluation.GetMigratedDataSize(evalConfig.StencilDBConn, evalConfig.DiasporaDBConn, evalConfig.DiasporaAppID, migrationID)
		log.Println("Migrated data size:", migratedDataSize)
		leftoverDataSize := evaluation.GetLeftoverDataSize(evalConfig.StencilDBConn, evalConfig.DiasporaDBConn, evalConfig.DiasporaAppID, migrationID)
		log.Println("Leftover data size:", leftoverDataSize)
		percentageOfLeftoverData := float64(leftoverDataSize) / (float64(migratedDataSize) + float64(leftoverDataSize))
		log.Println("Percentage of leftover data size:", percentageOfLeftoverData)
		data = append(data, percentageOfLeftoverData)

		// evaluation.GetPartiallyMappedRowTotalDataSize(evalConfig.StencilDBConn, evalConfig.MastodonAppID, dstMigrationID)
		// evaluation.GetPartiallyMappedRowDataSize(evalConfig.StencilDBConn, evalConfig.MastodonAppID, dstMigrationID)
	}
	evaluation.WriteToLog(leftoverVsMigratedFile, evaluation.ConvertFloat64ToString(data))
}

func anomaliesVsMigrationSize(evalConfig *evaluation.EvalConfig) {
	filterConditions := "and migration_id = 997596076"

	for _, dstMigrationID := range evaluation.GetAllMigrationIDsOfAppWithConds(evalConfig.StencilDBConn, evalConfig.MastodonAppID, filterConditions) {
		migrationID := strconv.FormatInt(dstMigrationID["migration_id"].(int64), 10)
		log.Println(migrationID)
		evaluation.GetAnomaliesNums(evalConfig, migrationID, "src")
	}

}

func main() {
	// mastodonConfig, _:= config.CreateAppConfig(mastodon, mastodonAppID)
	// common_phy_funcs.GetRowFromRowIDandTable(stencilDBConn, &mastodonConfig, "1008062662", "comments")
	
	evalConfig := evaluation.InitializeEvalConfig()
	// leftoverVsMigrated(evalConfig)
	anomaliesVsMigrationSize(evalConfig)
}
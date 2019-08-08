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
	filterConditions := "and user_id in (3300, 3344, 3503, 3482, 3924, 4134, 4322, 4386, 4503, 5001, 5323, 5370, 5458, 5574, 5602, 6012, 6431, 6853, 7168, 7251, 7488, 7557, 8239, 8263, 8506, 8563, 8664, 8894, 9017, 9051, 9739, 9716, 9831, 9857, 10082, 10286, 10781, 10795, 10979, 11141, 11318, 11321, 11351, 11348, 11455, 11487, 12536, 12724, 12726, 12823, 12789, 12963, 13158, 13879, 14031, 14226, 14351, 14777, 15265, 15401, 15495, 15505, 15517, 15579, 16043, 16127, 16900, 16994, 17209, 17279, 17639, 17680, 17732, 17691, 17809, 18208, 18569, 19079, 19372, 19563, 19613, 19656, 19746, 20265, 20269, 20254, 20448, 21119, 21537, 22535)"

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
	leftoverVsMigrated(evalConfig)
	// anomaliesVsMigrationSize(evalConfig)
}
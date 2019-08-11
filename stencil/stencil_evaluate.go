package main

import (
	"stencil/evaluation"
	"log"
	"strconv"
	"time"
	"fmt"
)

const (
	leftoverVsMigratedFile = "leftoverVsMigrated"
	srcAnomaliesVsMigrationSizeFile = "srcAnomaliesVsMigrationSize"
	dstAnomaliesVsMigrationSizeFile = "dstAnomaliesVsMigrationSize"
	interruptionDurationFile = "interruptionDuration"
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
	evaluation.WriteStrArrToLog(leftoverVsMigratedFile, evaluation.ConvertFloat64ToString(data))
}

func anomaliesVsMigrationSize(evalConfig *evaluation.EvalConfig) {
	filterConditions := "and user_id in (3300, 3344, 3503, 3482, 3924, 4134, 4322, 4386, 5001, 5323, 5370, 5458, 5574, 5602, 6012, 6431, 6853, 7168, 7251, 7488, 7557, 8239, 8263, 8506, 8563, 8664, 8894, 9017, 9051, 9739, 9716, 9831, 9857, 10082, 10286, 10781, 10795, 10979, 11141, 11318, 11321, 11351, 11348, 11455, 11487, 12536, 12724, 12726, 12823, 12789, 12963, 13158, 13879, 14031, 14226, 14351, 14777, 15265, 15401, 15495, 15505, 15517, 15579, 16043, 16127, 16900, 16994, 17209, 17279, 17639, 17680, 17732, 17691, 17809, 18208, 18569, 19079, 19372, 19563, 19613, 19656, 19746, 20265, 20269, 20254, 20448, 21119, 21537, 22535)"

	// filterConditions := "and user_id = 4503"

	totalSrcVoliateStats := make(map[string]int)
	var totalSrcInterruptionDuration []time.Duration
	totalDstViolateStats := make(map[string]int)
	totalDstDepNotMigratedStats := make(map[string]int)
	var totalMigratedDataSize int64

	for _, dstMigrationID := range evaluation.GetAllMigrationIDsOfAppWithConds(evalConfig.StencilDBConn, evalConfig.MastodonAppID, filterConditions) {
		migrationID := strconv.FormatInt(dstMigrationID["migration_id"].(int64), 10)
		log.Println(migrationID)

		srcViolateStats, srcInterruptionDuration := evaluation.GetAnomaliesNumsInSrc(evalConfig, migrationID, "src")

		log.Println("Source Violate Statistics:", srcViolateStats)
		log.Println("Source Interruption statistics:", srcInterruptionDuration)

		evaluation.IncreaseMapValByMap(totalSrcVoliateStats, srcViolateStats)
		totalSrcInterruptionDuration = append(totalSrcInterruptionDuration, srcInterruptionDuration...)

		evaluation.WriteStrArrToLog(interruptionDurationFile, evaluation.ConvertDurationToString(srcInterruptionDuration))
		evaluation.WriteStrToLog(srcAnomaliesVsMigrationSizeFile, evaluation.ConvertMapToJSONString(srcViolateStats))

		dstViolateStats, dstDepNotMigratedStats := evaluation.GetAnomaliesNumsInDst(evalConfig, migrationID, "dst")

		migratedDataSize := evaluation.GetMigratedDataSize(evalConfig.StencilDBConn, evalConfig.DiasporaDBConn, evalConfig.DiasporaAppID, migrationID)

		log.Println("Destination Violate Statistics:", dstViolateStats)
		log.Println("Destination Data depended on not migrated statistics:", dstDepNotMigratedStats)
		log.Println("Migrated data size(Bytes):", migratedDataSize)

		evaluation.WriteStrToLog(dstAnomaliesVsMigrationSizeFile, evaluation.ConvertMapToJSONString(dstViolateStats))
		evaluation.WriteStrToLog(dstAnomaliesVsMigrationSizeFile, evaluation.ConvertMapToJSONString(dstDepNotMigratedStats))
		evaluation.WriteStrToLog(dstAnomaliesVsMigrationSizeFile, fmt.Sprintf("%f", migratedDataSize))

		evaluation.IncreaseMapValByMap(totalDstViolateStats, dstViolateStats)
		evaluation.IncreaseMapValByMap(totalDstDepNotMigratedStats, dstDepNotMigratedStats)
		totalMigratedDataSize += migratedDataSize
	}

	log.Println("Destination Total Violate Statistics:", totalDstViolateStats)
	log.Println("Destination Total Data depended on not migrated statistics:", totalDstDepNotMigratedStats)
	log.Println("Source Total Violate Statistics:", totalSrcVoliateStats)
	log.Println("Source Total Interruption statistics:", totalSrcInterruptionDuration)
	log.Println("Total Migrated data size(Bytes):", totalMigratedDataSize)

	// evaluation.WriteStrArrToLog(interruptionDurationFile, evaluation.ConvertDurationToString(totalSrcInterruptionDuration))

}

func main() {
	// mastodonConfig, _:= config.CreateAppConfig(mastodon, mastodonAppID)
	// common_phy_funcs.GetRowFromRowIDandTable(stencilDBConn, &mastodonConfig, "1008062662", "comments")
	
	evalConfig := evaluation.InitializeEvalConfig()
	// leftoverVsMigrated(evalConfig)
	anomaliesVsMigrationSize(evalConfig)
}
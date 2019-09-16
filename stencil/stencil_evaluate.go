package main

import (
	"stencil/evaluation"
	"log"
	"strconv"
	// "time"
)

const (
	leftoverVsMigratedFile = "leftoverVsMigrated"
	// srcAnomaliesVsMigrationSizeFile = "srcAnomaliesVsMigrationSize"
	// dstAnomaliesVsMigrationSizeFile = "dstAnomaliesVsMigrationSize"
	interruptionDurationFile = "interruptionDuration"
	migrationRateFile = "migrationRate"
)

func leftoverVsMigrated(evalConfig *evaluation.EvalConfig) {
	var data []float64
	// Need to be changed once data is ready to use
	// filterConditions := "and user_id in (3300, 3344, 3503, 3482, 3924, 4134, 4322, 4386, 4503, 5001, 5323, 5370, 5458, 5574, 5602, 6012, 6431, 6853, 7168, 7251, 7488, 7557, 8239, 8263, 8506, 8563, 8664, 8894, 9017, 9051, 9739, 9716, 9831, 9857, 10082, 10286, 10781, 10795, 10979, 11141, 11318, 11321, 11351, 11348, 11455, 11487, 12536, 12724, 12726, 12823, 12789, 12963, 13158, 13879, 14031, 14226, 14351, 14777, 15265, 15401, 15495, 15505, 15517, 15579, 16043, 16127, 16900, 16994, 17209, 17279, 17639, 17680, 17732, 17691, 17809, 18208, 18569, 19079, 19372, 19563, 19613, 19656, 19746, 20265, 20269, 20254, 20448, 21119, 21537, 22535)"
	filterConditions := " LIMIT 100 "

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
	// filterConditions := "and user_id in (4386, 3300, 3344, 3503, 3482, 3924, 4134, 4322, 5001, 5323, 5370, 5458, 5574, 5602, 6012, 6431, 6853, 7168, 7251, 7488, 7557, 8239, 8263, 8506, 8563, 8664, 8894, 9017, 9051, 9739, 9716, 9831, 9857, 10082, 10286, 10781, 10795, 10979, 11141, 11318, 11321, 11351, 11348, 11455, 11487, 12536, 12724, 12726, 12823, 12789, 12963, 13158, 13879, 14031, 14226, 14351, 14777, 15265, 15401, 15495, 15505, 15517, 15579, 16043, 16127, 16900, 16994, 17209, 17279, 17639, 17680, 17732, 17691, 17809, 18208, 18569, 19079, 19372, 19563, 19613, 19656, 19746, 20265, 20269, 20254, 20448, 21119, 21537, 22535)"

	// 1 migration
	// filterConditions := "and start_time between '2019-08-17 17:17:00' and '2019-08-17 17:18:00' and is_logical = 'true' "
	// filterConditions := "and user_id = 1008 and is_logical = 't'"
	// 10 simultaneous logical migrations
	// filterConditions := "and start_time between '2019-08-18 09:42:00' and '2019-08-18 09:43:00' and is_logical = 'true' "
	// 20 simultaneous logical migrations
	// filterConditions := "and start_time between '2019-08-18 16:44:00' and '2019-08-18 16:46:00' and is_logical = 'true' "
	// filterConditions := "and user_id in (1815, 1818, 1006, 1075, 1731, 1032, 1020, 1103, 1041, 1044, 1060, 1819, 1028, 1705, 1107)"
	// 30 simultaneous logical migrations
	// 40 simultaneous logical migrations
	// 50 simultaneous logical migrations
	// filterConditions := "and registration_id > 464 and is_logical = 'true' "
	// 100
	// filterConditions := "and registration_id > 514 and is_logical = 'true' "
	
	filterConditions := "and migration_id = 356255340 "

	// totalSrcDanglingDataStats := make(map[string]int64)
	// totalSrcVoliateStats := make(map[string]int)
	// var totalSrcInterruptionDuration []time.Duration
	// totalDstViolateStats := make(map[string]int)
	// totalDstDepNotMigratedStats := make(map[string]int)
	// var totalMigratedDataSize int64

	for _, dstMigrationID := range evaluation.GetAllMigrationIDsOfAppWithConds(evalConfig.StencilDBConn, evalConfig.MastodonAppID, filterConditions) {
		migrationID := strconv.FormatInt(dstMigrationID["migration_id"].(int64), 10)
		
		evaluation.AnomaliesDanglingData(migrationID, evalConfig)
		
		// migrationID := strconv.FormatInt(dstMigrationID["migration_id"].(int64), 10)
		// log.Println(migrationID)

		// dstViolateStats, dstDepNotMigratedStats := evaluation.GetAnomaliesNumsInDst(evalConfig, migrationID)
		// srcViolateStats, srcInterruptionDuration, srcDanglingDataStats := evaluation.GetAnomaliesNumsInSrc(evalConfig, migrationID)
		
		// log.Println("Source Violate Statistics:", srcViolateStats)
		// log.Println("Source Interruption statistics:", srcInterruptionDuration)
		// log.Println("Source Dangling Statistics:", srcDanglingDataStats)

		// evaluation.IncreaseMapValByMapInt64(totalSrcDanglingDataStats, srcDanglingDataStats)
		// evaluation.IncreaseMapValByMap(totalSrcVoliateStats, srcViolateStats)
		// totalSrcInterruptionDuration = append(totalSrcInterruptionDuration, srcInterruptionDuration...)

		// evaluation.WriteStrArrToLog(interruptionDurationFile, evaluation.ConvertDurationToString(srcInterruptionDuration))
		// evaluation.WriteStrToLog(srcAnomaliesVsMigrationSizeFile, evaluation.ConvertMapToJSONString(srcViolateStats))
		// evaluation.WriteStrToLog(srcAnomaliesVsMigrationSizeFile, evaluation.ConvertMapInt64ToJSONString(srcDanglingDataStats))

		// // migratedDataSize := evaluation.GetMigratedDataSize(evalConfig.StencilDBConn, evalConfig.DiasporaDBConn, evalConfig.DiasporaAppID, migrationID)

		// log.Println("Destination Violate Statistics:", dstViolateStats)
		// log.Println("Destination Data depended on not migrated statistics:", dstDepNotMigratedStats)
		// // log.Println("Migrated data size(Bytes):", migratedDataSize)

		// evaluation.WriteStrToLog(dstAnomaliesVsMigrationSizeFile, evaluation.ConvertMapToJSONString(dstViolateStats))
		// evaluation.WriteStrToLog(dstAnomaliesVsMigrationSizeFile, evaluation.ConvertMapToJSONString(dstDepNotMigratedStats))
		// // evaluation.WriteStrToLog(dstAnomaliesVsMigrationSizeFile, evaluation.ConvertInt64ToString(migratedDataSize))

		// evaluation.IncreaseMapValByMap(totalDstViolateStats, dstViolateStats)
		// evaluation.IncreaseMapValByMap(totalDstDepNotMigratedStats, dstDepNotMigratedStats)
		// // totalMigratedDataSize += migratedDataSize
	}

	// log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	// log.Println("Destination Total Violate Statistics:", totalDstViolateStats)
	// log.Println("Destination Total Data depended on not migrated statistics:", totalDstDepNotMigratedStats)
	// log.Println("Source Total Violate Statistics:", totalSrcVoliateStats)
	// log.Println("Source Total Interruption statistics:", totalSrcInterruptionDuration)
	// log.Println("Source Total Dangling Data statistics:", totalSrcDanglingDataStats)
	// // log.Println("Total Migrated data size(Bytes):", totalMigratedDataSize)
	// log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")

	// // evaluation.WriteStrArrToLog(interruptionDurationFile, evaluation.ConvertDurationToString(totalSrcInterruptionDuration))

}

func migrationRate(evalConfig *evaluation.EvalConfig) {
	// filterConditions := "and migration_id = 1203415167"
	
	// for _, migrationIDType := range evaluation.GetAllMigrationIDsAndTypesOfAppWithConds(evalConfig.StencilDBConn, evalConfig.MastodonAppID, filterConditions) {
	// 	var migrationID int64
	// 	for k, v := range migrationIDType {
	// 		if k == "migration_id" {
	// 			migrationID = v.(int64)
	// 		}
	// 	}
		
	// 	time := evaluation.GetMigrationTime(evalConfig.StencilDBConn, migrationID)
	// 	log.Println("Migration time: ", time)
	// 	migratedDataSize := evaluation.GetMigratedDataSize(evalConfig.StencilDBConn, evalConfig.DiasporaDBConn, evalConfig.DiasporaAppID, strconv.FormatInt(migrationID, 10))
	// 	log.Println("Migrated data size: (KB)", migratedDataSize)

	// 	migrationRate := make(map[string]string)
	// 	migrationRate["time"] = evaluation.ConvertSingleDurationToString(time)
	// 	migrationRate["size"] = strconv.FormatInt(migratedDataSize, 10)
	// 	evaluation.WriteStrToLog(migrationRateFile, evaluation.ConvertMapStringToJSONString(migrationRate))
	// }
}

func main() {
	// mastodonConfig, _:= config.CreateAppConfig(mastodon, mastodonAppID)
	// common_phy_funcs.GetRowFromRowIDandTable(stencilDBConn, &mastodonConfig, "1008062662", "comments")
	
	evalConfig := evaluation.InitializeEvalConfig()
	// leftoverVsMigrated(evalConfig)
	// anomaliesVsMigrationSize(evalConfig)
	// evaluation.MigrationRate("1725984712", evalConfig)
	// evaluation.SystemLevelDanglingData(evalConfig)
	// evaluation.GetDataBagOfUser("1590693271", "diaspora", "mastodon", evalConfig)
	// evaluation.GetDataDownTime("1590693271", evalConfig)
	evaluation.GetSize("1725984712", evalConfig)
	evaluation.GetTime("1725984712", evalConfig)
}
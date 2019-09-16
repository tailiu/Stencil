package evaluation

import (
	"log"
	"strconv"
)

func AnomaliesDanglingData(migrationID string, evalConfig *EvalConfig) {
	// log.Println(migrationID)

	dstViolateStats, dstDepNotMigratedStats := GetAnomaliesNumsInDst(evalConfig, migrationID)
	srcViolateStats, srcInterruptionDuration, srcDanglingDataStats := GetAnomaliesNumsInSrc(evalConfig, migrationID)
	
	log.Println("Source Violate Statistics:", srcViolateStats)
	log.Println("Source Interruption statistics:", srcInterruptionDuration)
	log.Println("Source Dangling Statistics:", srcDanglingDataStats)

	WriteStrArrToLog(evalConfig.InterruptionDurationFile, ConvertDurationToString(srcInterruptionDuration))
	WriteStrToLog(evalConfig.SrcAnomaliesVsMigrationSizeFile, ConvertMapToJSONString(srcViolateStats))
	WriteStrToLog(evalConfig.SrcAnomaliesVsMigrationSizeFile, ConvertMapInt64ToJSONString(srcDanglingDataStats))

	// migratedDataSize := evaluation.GetMigratedDataSize(evalConfig.StencilDBConn, evalConfig.DiasporaDBConn, evalConfig.DiasporaAppID, migrationID)

	log.Println("Destination Violate Statistics:", dstViolateStats)
	log.Println("Destination Data depended on not migrated statistics:", dstDepNotMigratedStats)
	// log.Println("Migrated data size(Bytes):", migratedDataSize)

	WriteStrToLog(evalConfig.DstAnomaliesVsMigrationSizeFile, ConvertMapToJSONString(dstViolateStats))
	WriteStrToLog(evalConfig.DstAnomaliesVsMigrationSizeFile, ConvertMapToJSONString(dstDepNotMigratedStats))
	// evaluation.WriteStrToLog(dstAnomaliesVsMigrationSizeFile, evaluation.ConvertInt64ToString(migratedDataSize))
	// totalMigratedDataSize += migratedDataSize
}

func MigrationRate(migrationID string, evalConfig *EvalConfig) {
	// log.Println(migrationID)
	
	migrationID1, err := strconv.Atoi(migrationID)
	if err != nil {
		log.Fatal(err)
	}

	time := GetMigrationTime(evalConfig.StencilDBConn, migrationID1)
	log.Println("Migration time: ", time)
	migratedDataSize := GetMigratedDataSize(evalConfig.StencilDBConn, evalConfig.DiasporaDBConn, evalConfig.DiasporaAppID, migrationID)
	log.Println("Migrated data size: (KB)", migratedDataSize)

	migrationRate := make(map[string]string)
	migrationRate["time"] = ConvertSingleDurationToString(time)
	migrationRate["size"] = strconv.FormatInt(migratedDataSize, 10)
	
	WriteStrToLog(evalConfig.MigrationRateFile, ConvertMapStringToJSONString(migrationRate))
}

func SystemLevelDanglingData(evalConfig *EvalConfig) {
	srcDanglingDataStats := srcDanglingDataSystem(evalConfig)
	log.Println(srcDanglingDataStats)

	dstDanglingDataStats := dstDanglingDataSystem(evalConfig)
	log.Println(dstDanglingDataStats)

	WriteStrToLog(evalConfig.SrcDanglingDataInSystemFile, ConvertMapInt64ToJSONString(srcDanglingDataStats))
	WriteStrToLog(evalConfig.DstDanglingDataInSystemFile, ConvertMapInt64ToJSONString(dstDanglingDataStats))
}

func GetDataBagOfUser(migrationID, sourceApp, dstApp string, evalConfig *EvalConfig) {
	srcDataBagSize := getDataBagSize(evalConfig, sourceApp, migrationID)
	dstDataBagSize := getDataBagSize(evalConfig, dstApp, migrationID)
	log.Println(srcDataBagSize)
	log.Println(dstDataBagSize)
}

func GetDataDownTime(migrationID, evalConfig *EvalConfig) {
	
}
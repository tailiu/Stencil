package evaluation

import (
	"log"
)

func AnomaliesDanglingData(migrationID string, evalConfig *EvalConfig) {
	log.Println(migrationID)

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
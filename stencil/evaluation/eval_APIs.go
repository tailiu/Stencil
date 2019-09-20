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
	log.Println("Migration time: (s)", time)
	migratedDataSize := GetMigratedDataSize(evalConfig.StencilDBConn, evalConfig.DiasporaDBConn, evalConfig.DiasporaAppID, migrationID)
	log.Println("Migrated data size: (Bytes)", migratedDataSize)

	migrationRate := make(map[string]string)
	migrationRate["time"] = ConvertSingleDurationToString(time)
	migrationRate["size"] = strconv.FormatInt(migratedDataSize, 10)
	
	WriteStrToLog(evalConfig.MigrationRateFile, ConvertMapStringToJSONString(migrationRate))
}

func GetSize(migrationID string, evalConfig *EvalConfig) {
	migratedDataSize := GetMigratedDataSize(evalConfig.StencilDBConn, evalConfig.DiasporaDBConn, evalConfig.DiasporaAppID, migrationID)
	log.Println("Migrated data size: (Bytes)", migratedDataSize)

	migration := make(map[string]string)
	migration["size"] = strconv.FormatInt(migratedDataSize, 10)

	WriteStrToLog(evalConfig.MigratedDataSizeFile, ConvertMapStringToJSONString(migration))
}

func GetTime(migrationID string, evalConfig *EvalConfig) {
	migrationID1, err := strconv.Atoi(migrationID)
	if err != nil {
		log.Fatal(err)
	}
	time := GetMigrationTime(evalConfig.StencilDBConn, migrationID1)
	log.Println("Migration time: (s)", time)

	migration := make(map[string]string)
	migration["time"] = ConvertSingleDurationToString(time)
	WriteStrToLog(evalConfig.MigrationTimeFile, ConvertMapStringToJSONString(migration))
}

func SystemLevelDanglingData(migrationID string, evalConfig *EvalConfig) {
	srcDanglingDataStats := srcDanglingDataSystem(evalConfig)
	log.Println(srcDanglingDataStats)

	dstDanglingDataStats := dstDanglingDataSystem(evalConfig, migrationID)
	log.Println(dstDanglingDataStats)

	WriteStrToLog(evalConfig.SrcDanglingDataInSystemFile, ConvertMapInt64ToJSONString(srcDanglingDataStats))
	WriteStrToLog(evalConfig.DstDanglingDataInSystemFile, ConvertMapInt64ToJSONString(dstDanglingDataStats))
}

// func GetDataBagOfUser(migrationID, srcApp, dstApp string, evalConfig *EvalConfig) {
// 	migratedNodeSize := getTotalMigratedNodeSize(evalConfig, dstApp, migrationID)
// 	log.Println(migratedNodeSize)
// 	displayedDataSize := getDisplayedDataSize(evalConfig, srcApp, dstApp, migrationID)
// 	log.Println(displayedDataSize)

// 	dataBags := make(map[string]int64)
// 	dataBags["migratedNodeSize"] = migratedNodeSize
// 	dataBags["displayedDataSize"] = displayedDataSize
// 	WriteStrToLog(evalConfig.DataBags, ConvertMapInt64ToJSONString(dataBags))
// }

func GetDataBagOfUser(migrationID, sourceApp, dstApp string, evalConfig *EvalConfig) {
	srcDataBagSize := getDataBagSize(evalConfig, sourceApp, migrationID)
	dstDataBagSize := getDataBagSize(evalConfig, dstApp, migrationID)
	log.Println(srcDataBagSize)
	log.Println(dstDataBagSize)

	dataBags := make(map[string]int64)
	dataBags["srcDataBagSize"] = srcDataBagSize
	dataBags["dstDataBagSize"] = dstDataBagSize
	WriteStrToLog(evalConfig.DataBags, ConvertMapInt64ToJSONString(dataBags))
}

func GetDataDowntimeInStencil(migrationID string, evalConfig *EvalConfig) {
	dataDowntimeInStencil := getDataDowntimeInStencil(migrationID, evalConfig)
	WriteStrArrToLog(evalConfig.DataDowntimeInStencilFile, ConvertDurationToString(dataDowntimeInStencil))
}

func GetDataDowntimeInNaiveMigration(stencilMigrationID string, naiveMigrationID string, evalConfig *EvalConfig) {
	dataDowntimeInNaive := getDataDowntimeInNaive(stencilMigrationID, naiveMigrationID, evalConfig)
	WriteStrArrToLog(evalConfig.DataDowntimeInNaiveFile, ConvertDurationToString(dataDowntimeInNaive))
}
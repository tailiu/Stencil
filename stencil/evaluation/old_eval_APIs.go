package evaluation

import (
	"log"
	"strconv"
	"stencil/db"
	"stencil/SA1_migrate"
)

func AnomaliesDanglingData(migrationID string, evalConfig *EvalConfig) {
	// log.Println(migrationID)

	dstViolateStats, dstDepNotMigratedStats := 
		GetAnomaliesNumsInDst(evalConfig, migrationID)
	
	srcViolateStats, srcInterruptionDuration, srcDanglingDataStats := 
		GetAnomaliesNumsInSrc(evalConfig, migrationID)
	
	log.Println("Source Violate Statistics:", srcViolateStats)
	log.Println("Source Interruption statistics:", srcInterruptionDuration)
	log.Println("Source Dangling Statistics:", srcDanglingDataStats)

	WriteStrArrToLog(
		evalConfig.InterruptionDurationFile, 
		ConvertDurationToString(srcInterruptionDuration),
	)
	
	WriteStrToLog(
		evalConfig.SrcAnomaliesVsMigrationSizeFile, 
		ConvertMapToJSONString(srcViolateStats),
	)

	WriteStrToLog(
		evalConfig.SrcAnomaliesVsMigrationSizeFile, 
		ConvertMapInt64ToJSONString(srcDanglingDataStats),
	)

	// migratedDataSize := evaluation.GetMigratedDataSize(evalConfig.StencilDBConn, 
		// evalConfig.DiasporaDBConn, evalConfig.DiasporaAppID, migrationID)

	log.Println("Destination Violate Statistics:", dstViolateStats)	
	log.Println("Destination Data depended on not migrated statistics:", 
		dstDepNotMigratedStats)
	
	// log.Println("Migrated data size(Bytes):", migratedDataSize)

	WriteStrToLog(
		evalConfig.DstAnomaliesVsMigrationSizeFile, 
		ConvertMapToJSONString(dstViolateStats),
	)

	WriteStrToLog(
		evalConfig.DstAnomaliesVsMigrationSizeFile, 
		ConvertMapToJSONString(dstDepNotMigratedStats),
	)
	
	// evaluation.WriteStrToLog(dstAnomaliesVsMigrationSizeFile, 
	// evaluation.ConvertInt64ToString(migratedDataSize))
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

	migratedDataSize := GetMigratedDataSize(
		evalConfig.StencilDBConn, 
		evalConfig.DiasporaDBConn, 
		evalConfig.DiasporaAppID, 
		migrationID,
	)
	
	log.Println("Migrated data size: (Bytes)", migratedDataSize)

	migrationRate := make(map[string]string)
	migrationRate["time"] = ConvertSingleDurationToString(time)
	migrationRate["size"] = strconv.FormatInt(migratedDataSize, 10)
	
	WriteStrToLog(
		evalConfig.MigrationRateFile, 
		ConvertMapStringToJSONString(migrationRate),
	)

}

func GetSize(migrationID string, evalConfig *EvalConfig) {
	
	migratedDataSize := GetMigratedDataSize(
		evalConfig.StencilDBConn, 
		evalConfig.DiasporaDBConn, 
		evalConfig.DiasporaAppID, 
		migrationID,
	)
	
	log.Println("Migrated data size: (Bytes)", migratedDataSize)

	migration := make(map[string]string)
	migration["size"] = strconv.FormatInt(migratedDataSize, 10)

	WriteStrToLog(
		evalConfig.MigratedDataSizeFile, 
		ConvertMapStringToJSONString(migration),
	)

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
	
	WriteStrToLog(
		evalConfig.MigrationTimeFile, 
		ConvertMapStringToJSONString(migration),
	)
}

func SystemLevelDanglingData(migrationID string, evalConfig *EvalConfig) {

	srcDanglingDataStats := srcDanglingDataSystem(evalConfig)
	log.Println(srcDanglingDataStats)

	dstDanglingDataStats := dstDanglingDataSystem(evalConfig, migrationID)
	log.Println(dstDanglingDataStats)

	WriteStrToLog(
		evalConfig.SrcDanglingDataInSystemFile, 
		ConvertMapInt64ToJSONString(srcDanglingDataStats),
	)
	WriteStrToLog(
		evalConfig.DstDanglingDataInSystemFile, 
		ConvertMapInt64ToJSONString(dstDanglingDataStats),
	)
}

// func oldGetDataBagOfUser(migrationID, srcApp, dstApp string, 
// 	evalConfig *EvalConfig) {
	
// 	migratedNodeSize := getTotalMigratedNodeSize(evalConfig, dstApp, migrationID)
// 	log.Println(migratedNodeSize)
	
// 	displayedDataSize := getDisplayedDataSize(evalConfig, srcApp, dstApp, migrationID)
// 	log.Println(displayedDataSize)

// 	dataBags := make(map[string]int64)
// 	dataBags["migratedNodeSize"] = migratedNodeSize
// 	dataBags["displayedDataSize"] = displayedDataSize
	
// 	WriteStrToLog(
// 		evalConfig.DataBags, 
// 		ConvertMapInt64ToJSONString(dataBags),
// 	)

// }

func oldGetDataBagOfUserBasedOnApp(migrationID, sourceApp, dstApp string, 
	evalConfig *EvalConfig) {
	
	srcDataBagSize := getDataBagSize(evalConfig, sourceApp, migrationID)
	dstDataBagSize := getDataBagSize(evalConfig, dstApp, migrationID)
	
	log.Println(srcDataBagSize)
	log.Println(dstDataBagSize)

	dataBags := make(map[string]int64)
	dataBags["srcDataBagSize"] = srcDataBagSize
	dataBags["dstDataBagSize"] = dstDataBagSize
	
	WriteStrToLog(
		evalConfig.DataBags, 
		ConvertMapInt64ToJSONString(dataBags),
	)
}

func GetDataBagOfUser(userID string, evalConfig *EvalConfig) {
	apps := getAllAppsOfDataBag(evalConfig, userID)
	
	var size int64 
	
	for _, app := range apps {
		size += getDataBagSize(evalConfig, app, userID)
	}
	
	dataBags := make(map[string]int64)
	dataBags["dataBagSize"] = size

	WriteStrToLog(
		evalConfig.DataBags, 
		ConvertMapInt64ToJSONString(dataBags),
	)
}

func GetDataDowntimeInStencil(migrationID string, evalConfig *EvalConfig) {
	
	dataDowntimeInStencil := getDataDowntimeInStencil(migrationID, evalConfig)
	
	WriteStrArrToLog(
		evalConfig.DataDowntimeInStencilFile, 
		ConvertDurationToString(dataDowntimeInStencil),
	)

}

func GetDataDowntimeInNaiveMigration(stencilMigrationID string, naiveMigrationID string, 
	evalConfig *EvalConfig) {
	
	dataDowntimeInNaive := getDataDowntimeInNaive(
		stencilMigrationID, 
		naiveMigrationID, 
		evalConfig,
	)
	
	WriteStrArrToLog(
		evalConfig.DataDowntimeInNaiveFile, 
		ConvertDurationToString(dataDowntimeInNaive),
	)
}

func oldExp2GetMigratedDataRate() {
	
	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	migrationIDs := GetAllMigrationIDs(evalConfig)

	for _, migrationID := range migrationIDs {

		size := GetMigratedDataSizeV2(
			evalConfig,
			migrationID,
		)

		migrationIDInt, err := strconv.Atoi(migrationID)
		if err != nil {
			log.Fatal(err)
		}
		
		time := GetMigrationTime(
			evalConfig.StencilDBConn,
			migrationIDInt,
		)

		WriteStrToLog(
			evalConfig.MigratedDataSizeFile, 
			ConvertInt64ToString(size),
		)

		WriteStrToLog(
			evalConfig.MigrationTimeFile,
			ConvertSingleDurationToString(time),
		)

	}

}

func migrateUserUsingSA1AndNaive(evalConfig *EvalConfig, 
	SA1StencilDB, SA1SrcDB, SA1DstDB, userID,
	naiveStencilDB, naiveSrcDB, naiveDstDB string, 
	SA1EnableDisplay, SA1DisplayInFirstPhase, 
	naiveEnableDisplay, naiveDisplayInFirstPhase bool) {

	sizeLog := make(map[string]string)
	timeLog := make(map[string]string)

	db.STENCIL_DB = SA1StencilDB
	db.DIASPORA_DB = SA1SrcDB
	db.MASTODON_DB = SA1DstDB

	uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
		userID, "diaspora", "1", "mastodon", "2", "d", 1

	SA1_migrate.Controller(uid, srcAppName, srcAppID, 
		dstAppName, dstAppID, migrationType, threadNum,
		SA1EnableDisplay, SA1DisplayInFirstPhase,
	)

	db.STENCIL_DB = naiveStencilDB
	db.DIASPORA_DB = naiveSrcDB
	db.MASTODON_DB = naiveDstDB

	migrationType = "n"

	SA1_migrate.Controller(uid, srcAppName, srcAppID, 
		dstAppName, dstAppID, migrationType, threadNum,
		naiveEnableDisplay, naiveDisplayInFirstPhase,
	)

	dMigrationID := 
		getMigrationIDBySrcUserIDMigrationType(evalConfig.StencilDBConn, userID, "d")

	nMigrationID := 
		getMigrationIDBySrcUserIDMigrationType(evalConfig.StencilDBConn, userID, "n")

	dMigrationIDInt, err := strconv.Atoi(dMigrationID)
	if err != nil {
		log.Fatal(err)
	}

	dTime := GetMigrationTime(
		evalConfig.StencilDBConn,
		dMigrationIDInt,
	)

	nMigrationIDInt, err := strconv.Atoi(nMigrationID)
	if err != nil {
		log.Fatal(err)
	}

	nTime := GetMigrationTime(
		evalConfig.StencilDBConn,
		nMigrationIDInt,
	)

	timeLog["deletion_time"] = ConvertSingleDurationToString(dTime)	
	timeLog["naive_time"] = ConvertSingleDurationToString(nTime)
	
	size := GetMigratedDataSizeByDst(
		evalConfig,
		dMigrationID,
	)

	sizeLog["size"] = ConvertInt64ToString(size)

	WriteStrToLog(
		evalConfig.MigratedDataSizeFile, 
		ConvertMapStringToJSONString(sizeLog),
	)

	WriteStrToLog(
		evalConfig.MigrationTimeFile,
		ConvertMapStringToJSONString(timeLog),
	)
}

func RecreateDiaspora1MDB() {

	diaspora = "diaspora_test"

	dbConn := db.GetDBConn(diaspora)

	defer closeDBConn(dbConn)

	templateDB := "diaspora_1000000"

	recreateDBByTemplate(dbConn, "diaspora_1000000_exp", templateDB)

	recreateDBByTemplate(dbConn, "diaspora_1000000_exp1", templateDB)

	recreateDBByTemplate(dbConn, "diaspora_1000000_exp2", templateDB)

}	
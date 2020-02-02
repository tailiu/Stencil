package evaluation

import (
	"stencil/db"
	"stencil/SA1_migrate"
	"log"
	"strconv"
)


func migrateUserUsingSA1AndNaive(evalConfig *EvalConfig, 
	migrationNum int, SA1SrcDB, SA1DstDB, naiveSrcDB, naiveDstDB string,
	SA1EnableDisplay, SA1DisplayInFirstPhase, naiveEnableDisplay, 
	naiveDisplayInFirstPhase bool) {

	userIDs := getAllUserIDsInDiaspora(evalConfig)

	shuffleSlice(userIDs)

	// res := make(map[string]string)

	for i := 0; i < migrationNum; i ++ {

		sizeLog := make(map[string]string)
		timeLog := make(map[string]string)

		db.DIASPORA_DB = SA1SrcDB
		db.MASTODON_DB = SA1DstDB

		uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
			userIDs[i], "diaspora", "1", "mastodon", "2", "d", 1

		SA1_migrate.Controller(uid, srcAppName, srcAppID, 
			dstAppName, dstAppID, migrationType, threadNum,
			SA1EnableDisplay, SA1DisplayInFirstPhase,
		)
		
		db.DIASPORA_DB = naiveSrcDB
		db.MASTODON_DB = naiveDstDB

		migrationType = "n"

		SA1_migrate.Controller(uid, srcAppName, srcAppID, 
			dstAppName, dstAppID, migrationType, threadNum,
			naiveEnableDisplay, naiveDisplayInFirstPhase,
		)

		dMigrationID := 
			getMigrationIDBySrcUserIDMigrationType(evalConfig, userIDs[i], "d")

		nMigrationID := 
			getMigrationIDBySrcUserIDMigrationType(evalConfig, userIDs[i], "n")

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
}
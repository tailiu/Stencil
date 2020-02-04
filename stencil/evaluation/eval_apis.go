package evaluation

import (
	"stencil/db"
	"stencil/SA1_migrate"
	"log"
	"database/sql"
	"strconv"
)


func migrateUserFromDiasporaToMastodon(
	evalConfig *EvalConfig, evalStencilDB, evalReferDB *sql.DB, 
	userID, migrationType, stencilDB, srcDB, dstDB, 
	sizeFile, timeFile string,
	enableDisplay, displayInFirstPhase bool) {

	sizeLog := make(map[string]string)
	timeLog := make(map[string]string)

	db.STENCIL_DB = stencilDB
	db.DIASPORA_DB = srcDB
	db.MASTODON_DB = dstDB

	srcAppName, srcAppID, dstAppName, dstAppID, threadNum := 
		"diaspora", "1", "mastodon", "2", 1

	SA1_migrate.Controller(userID, srcAppName, srcAppID, 
		dstAppName, dstAppID, migrationType, threadNum,
		enableDisplay, displayInFirstPhase,
	)

	migrationID := 
		getMigrationIDBySrcUserIDMigrationType(evalStencilDB, userID, migrationType)

	migrationIDInt, err := strconv.Atoi(migrationID)
	if err != nil {
		log.Fatal(err)
	}

	time := GetMigrationTime(
		evalStencilDB,
		migrationIDInt,
	)

	size := GetMigratedDataSizeBySrc(
		evalConfig,
		evalStencilDB,
		evalReferDB,
		migrationID,
	)

	timeLog["time"] = ConvertSingleDurationToString(time)	
	timeLog["userID"] = userID

	sizeLog["size"] = ConvertInt64ToString(size)
	sizeLog["userID"] = userID

	WriteStrToLog(
		timeFile,
		ConvertMapStringToJSONString(timeLog),
	)

	WriteStrToLog(
		sizeFile, 
		ConvertMapStringToJSONString(sizeLog),
	)
}
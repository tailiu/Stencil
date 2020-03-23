package evaluation

import (
	"stencil/db"
	"stencil/SA1_migrate"
	"stencil/reference_resolution"
	"database/sql"
	"log"
	"strconv"
)


func migrateUserFromDiasporaToMastodon(
	evalConfig *EvalConfig, evalStencilDBName, evalReferDBName string, 
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

	log.Println("************ Calculate the Migration Size and Time ************")

	refreshEvalConfigDBConnections(evalConfig)
	
	evalStencilDB := getDBConnByName(evalConfig, evalStencilDBName)
	evalReferDB := getDBConnByName(evalConfig, evalReferDBName)	

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

func migrateUserFromDiasporaToMastodon1(
	userID, migrationType, stencilDB, srcDB, dstDB string,
	enableDisplay, displayInFirstPhase bool) {

	db.STENCIL_DB = stencilDB
	db.DIASPORA_DB = srcDB
	db.MASTODON_DB = dstDB

	srcAppName, srcAppID, dstAppName, dstAppID, threadNum := 
		"diaspora", "1", "mastodon", "2", 1

	SA1_migrate.Controller(userID, srcAppName, srcAppID, 
		dstAppName, dstAppID, migrationType, threadNum,
		enableDisplay, displayInFirstPhase,
	)
}

func migrateUsersInExp7(evalConfig *EvalConfig, stencilDBConnName string, 
	seqNum int, fromApp, toApp, fromAppID, toAppID string,
	migrationIDs, userIDs []string, 
	enableBagsOption, enableDisplayOption bool) []string {

	var stencilDBConn *sql.DB

	for j, userID := range userIDs {
		
		if seqNum != 0 {

			userNum := len(userIDs)

			stencilDBConn = getDBConnByName(evalConfig, stencilDBConnName)

			userID = reference_resolution.GetNextUserID(
				stencilDBConn, 
				migrationIDs[(seqNum - 1) * userNum + j],
			)
		}

		log.Println("Migrating user ID:", userID)
		log.Println("From app:", fromApp)
		log.Println("To app:", toApp)

		uid, migrationType, threadNum := userID, "d", 1

		enableDisplay, displayInFirstPhase, markAsDelete, useBladeServerAsDst, enableBags := 
			enableDisplayOption, false, false, false, enableBagsOption

		SA1_migrate.Controller(uid, fromApp, fromAppID, 
			toApp, toAppID, migrationType, threadNum,
			enableDisplay, displayInFirstPhase, 
			markAsDelete, useBladeServerAsDst, enableBags,
		)

		refreshEvalConfigDBConnections(evalConfig, false)

		stencilDBConn = getDBConnByName(evalConfig, stencilDBConnName)

		migrationID := getMigrationIDBySrcUserIDMigrationTypeFromToAppID(
			stencilDBConn, userID, 
			fromAppID, toAppID, migrationType,
		)
		
		migrationIDs = append(migrationIDs, migrationID)

	}

	return migrationIDs

}
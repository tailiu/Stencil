package evaluation

import (
	"stencil/SA1_migrate"
	"log"
)

func Exp1() {

	evalConfig := InitializeEvalConfig()

	userIDs := getAllUserIDsInDiaspora(evalConfig)

	for _, userID := range userIDs {
		uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
			userID, "diaspora", "1", "mastodon", "2", "d", 1

		SA1_migrate.Controller(uid, srcAppName, srcAppID, 
			dstAppName, dstAppID, migrationType, threadNum)
	}

	var sizes []int64

	for _, userID := range userIDs {
		migrationID := getMigrationIDBySrcUserID(evalConfig, userID)
		sizes = append(sizes, getDanglingDataSizeOfMigration(evalConfig, migrationID))
	}

	log.Println(userIDs)

	log.Println(sizes)
	
	WriteStrArrToLog("exp1", ConvertInt64ArrToStringArr(sizes))

}
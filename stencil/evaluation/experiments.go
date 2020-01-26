package evaluation

import (
	"stencil/SA1_migrate"
	"stencil/db"
	"log"
)

func preExp(evalConfig *EvalConfig) {

	query1 := "TRUNCATE identity_table, migration_registration, reference_table, resolved_references"

	query2 := "SELECT truncate_tables('cow')"

	if err1 := db.TxnExecute1(evalConfig.StencilDBConn, query1); err1 != nil {
		log.Fatal(err1)
	} else {
		if err2 := db.TxnExecute1(evalConfig.MastodonDBConn, query2); err2 != nil {
			log.Fatal(err2)
		} else {
			return
		}
	}

}

func Exp1() {

	evalConfig := InitializeEvalConfig()

	preExp(evalConfig)

	userIDs := getAllUserIDsInDiaspora(evalConfig)

	for _, userID := range userIDs {
		uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
			userID, "diaspora", "1", "mastodon", "2", "d", 1

		SA1_migrate.Controller(uid, srcAppName, srcAppID, 
			dstAppName, dstAppID, migrationType, threadNum)
	}

	log.Println(userIDs)

	var sizes []int64

	for _, userID := range userIDs {
		migrationID := getMigrationIDBySrcUserID(evalConfig, userID)
		sizes = append(sizes, getDanglingDataSizeOfMigration(evalConfig, migrationID))
	}

	log.Println(sizes)
	
	WriteStrArrToLog("exp1", ConvertInt64ArrToStringArr(sizes))

}
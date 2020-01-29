package evaluation

import (
	"stencil/SA1_migrate"
	"stencil/db"
	"log"
)

func preExp(evalConfig *EvalConfig) {

	query1 := `TRUNCATE identity_table, migration_registration, 
		reference_table, resolved_references, txn_logs, 
		evaluation, data_bags, display_flags`

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

// In this experiment, we migrate 1000 users from Diaspora to Mastodon
// Note that in this exp the migration thread should not migrate data from data bags
// The source database needs to be changed to diaspora_1000
func Exp1() {

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	preExp(evalConfig)

	userIDs := getAllUserIDsInDiaspora(evalConfig)

	shuffleSlice(userIDs)

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
	
	WriteStrArrToLog(
		"exp1", 
		ConvertInt64ArrToStringArr(sizes),
	)

}

func Exp1GetMediaSize() {

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	mediaSize := getAllMediaSize(evalConfig)

	log.Println("Total Media Size:", mediaSize, "bytes")
	
}

// The source database needs to be changed to diaspora_1000000_exp1
// Data will be migrated from diaspora_1000000_exp1
// We can get data size by diaspora_1000000_exp
func Exp2() {

	migrationNum := 100

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	preExp(evalConfig)

	userIDs := getAllUserIDsInDiaspora(evalConfig)

	shuffleSlice(userIDs)

	res := make(map[string]string)

	for i := 0; i < migrationNum; i ++ {

		uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
			userIDs[i], "diaspora", "1", "mastodon", "2", "d", 1

		SA1_migrate.Controller(uid, srcAppName, srcAppID, 
			dstAppName, dstAppID, migrationType, threadNum)
		
		res[userIDs[i]] = "true"

	}

	log.Println(res)
	
	WriteStrToLog(
		"exp2",
		ConvertMapStringToJSONString(res),
	)

}

func Exp3() {
	
	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	migrationID := "1111158775"

	// GetMigratedDataSize(
	// 	evalConfig.StencilDBConn, 
	// 	evalConfig.DiasporaDBConn, 
	// 	"1",
	// 	migrationID,
	// )

	GetMigratedDataSizeV2(
		evalConfig,
		migrationID,
	)

}
package evaluation

import (
	"stencil/SA1_migrate"
	"stencil/apis"
	"stencil/db"
	"log"
	"fmt"
	"strconv"
	"time"
)

func preExp(evalConfig *EvalConfig) {

	query1 := `TRUNCATE identity_table, migration_registration, 
		reference_table, resolved_references, txn_logs, 
		evaluation, data_bags, display_flags, display_registration`

	query2 := "SELECT truncate_tables('cow')"

	query3 := "SELECT truncate_tables('cow')"

	if err1 := db.TxnExecute1(evalConfig.StencilDBConn, query1); err1 != nil {
		log.Fatal(err1)
	} else {
		if err2 := db.TxnExecute1(evalConfig.MastodonDBConn, query2); err2 != nil {
			log.Fatal(err2)
		} else {
			if err3 := db.TxnExecute1(evalConfig.MastodonDBConn1, query3); err3 != nil {
				log.Fatal(err3)
			} else {
				return
			}
		}
	}

}

func PreExp() {

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	preExp(evalConfig)

}

// In this experiment, we migrate 1000 users from Diaspora to Mastodon
// Note that in this exp the migration thread should not migrate data from data bags
// The source database needs to be changed to diaspora_1000_exp
func Exp1() {

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	preExp(evalConfig)

	userIDs := getAllUserIDsInDiaspora(evalConfig)

	shuffleSlice(userIDs)

	for _, userID := range userIDs {

		uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
			userID, "diaspora", "1", "mastodon", "2", "d", 1

		enableDisplay, displayInFirstPhase := true, true

		SA1_migrate.Controller(uid, srcAppName, srcAppID, 
			dstAppName, dstAppID, migrationType, threadNum,
			enableDisplay, displayInFirstPhase,
		)

		migrationID := getMigrationIDBySrcUserID(evalConfig, userID)

		danglingData := make(map[string]int64)

		srcDanglingData, dstDanglingData :=
			getDanglingDataSizeOfMigration(evalConfig, migrationID)

		danglingData["srcDanglingData"] = srcDanglingData
		danglingData["dstDanglingData"] = dstDanglingData

		WriteStrToLog(
			evalConfig.DanglingDataFile,
			ConvertMapInt64ToJSONString(danglingData),
		)

	}

}

func Exp1GetMediaSize() {

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	mediaSize := getAllMediaSize(evalConfig)

	log.Println("Total Media Size:", mediaSize, "bytes")
	
}

// The diaspora database needs to be changed to diaspora_1xxxx_exp and diaspora_1xxxx_exp1
// 1. Data will be migrated in deletion migrations from:
// diaspora_1000000_exp, diaspora_100000_exp, diaspora_10000_exp, diaspora_1000_exp
// to mastodon
// 2. Data will be migrated in naive migrations from:
// diaspora_1000000_exp1, diaspora_100000_exp1, diaspora_10000_exp1, diaspora_1000_exp1
// to mastodon_exp
// Notice that enableDisplay, displayInFirstPhase need to be changed in different exps
func Exp2() {

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	preExp(evalConfig)

	migrationNum := 100

	SA1SrcDB, SA1DstDB := "diaspora_1000000_exp", "mastodon"
	
	NaiveSrcDB, NaiveDstDB := "diaspora_1000000_exp1", "mastodon_exp"

	SA1EnableDisplay, SA1DisplayInFirstPhase := false, false

	NaiveEnableDisplay, NaiveDisplayInFirstPhase := false, false

	migrateUserUsingSA1AndNaive(evalConfig, migrationNum, SA1SrcDB, SA1DstDB,
		NaiveSrcDB, NaiveDstDB, SA1EnableDisplay, SA1DisplayInFirstPhase,
		NaiveEnableDisplay, NaiveDisplayInFirstPhase,
	)

}

// For all the three following get migrated data rate functions,
// the diaspora database needs to be changed to diaspora_1xxxx which has complete data
// We can get data size by the following complete dbs:
// diaspora_1000000, diaspora_100000, diaspora_10000, diaspora_1000
func Exp2GetMigratedDataRate() {
	
	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	migrationData := GetMigrationData(evalConfig)

	for _, migrationData1 := range migrationData {

		sizeLog := make(map[string]string)
		timeLog := make(map[string]string)

		sizeLog["userID"] = fmt.Sprint(migrationData1["user_id"])
		timeLog["userID"] = fmt.Sprint(migrationData1["user_id"])

		migrationID := fmt.Sprint(migrationData1["migration_id"])

		log.Println("Migration ID:", migrationID)

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

		sizeLog["size"] = ConvertInt64ToString(size)
		timeLog["time"] = ConvertSingleDurationToString(time)

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

func Exp2GetMigratedDataRateBySrc() {
	
	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	migrationData := GetMigrationData(evalConfig)

	for _, migrationData1 := range migrationData {

		sizeLog := make(map[string]string)
		timeLog := make(map[string]string)

		sizeLog["userID"] = fmt.Sprint(migrationData1["user_id"])
		timeLog["userID"] = fmt.Sprint(migrationData1["user_id"])

		migrationID := fmt.Sprint(migrationData1["migration_id"])

		log.Println("Migration ID:", migrationID)

		size := GetMigratedDataSizeBySrc(
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

		sizeLog["size"] = ConvertInt64ToString(size)
		timeLog["time"] = ConvertSingleDurationToString(time)

		WriteStrToLog(
			evalConfig.MigratedDataSizeBySrcFile, 
			ConvertMapStringToJSONString(sizeLog),
		)

		WriteStrToLog(
			evalConfig.MigrationTimeBySrcFile,
			ConvertMapStringToJSONString(timeLog),
		)

	}

}

func Exp2GetMigratedDataRateByDst() {
	
	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	migrationData := GetMigrationData(evalConfig)

	for _, migrationData1 := range migrationData {

		sizeLog := make(map[string]string)
		timeLog := make(map[string]string)

		sizeLog["userID"] = fmt.Sprint(migrationData1["user_id"])
		timeLog["userID"] = fmt.Sprint(migrationData1["user_id"])

		migrationID := fmt.Sprint(migrationData1["migration_id"])

		log.Println("Migration ID:", migrationID)

		size := GetMigratedDataSizeByDst(
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

		sizeLog["size"] = ConvertInt64ToString(size)
		timeLog["time"] = ConvertSingleDurationToString(time)

		WriteStrToLog(
			evalConfig.MigratedDataSizeByDstFile, 
			ConvertMapStringToJSONString(sizeLog),
		)

		WriteStrToLog(
			evalConfig.MigrationTimeByDstFile,
			ConvertMapStringToJSONString(timeLog),
		)

	}

}

func Exp3() {

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	preExp(evalConfig)

	migrationNum := 200

	SA1SrcDB, SA1DstDB := "diaspora_1000000_exp", "mastodon"
	
	naiveSrcDB, naiveDstDB := "diaspora_1000000_exp1", "mastodon_exp"

	SA1EnableDisplay, SA1DisplayInFirstPhase := true, true

	naiveEnableDisplay, naiveDisplayInFirstPhase := true, false

	migrateUserUsingSA1AndNaive(evalConfig, migrationNum, SA1SrcDB, SA1DstDB,
		naiveSrcDB, naiveDstDB, SA1EnableDisplay, SA1DisplayInFirstPhase,
		naiveEnableDisplay, naiveDisplayInFirstPhase,
	)

}

func Exp3GetDatadowntime() {

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	migrationData := GetMigrationData(evalConfig)

	var dDowntime, nDowntime []time.Duration

	for _, migrationData1 := range migrationData {

		migrationType := fmt.Sprint(migrationData1["migration_type"])

		migrationID := fmt.Sprint(migrationData1["migration_id"])

		if migrationType == "3" {

			downtime := getDataDowntimeOfMigration(evalConfig, migrationID)

			dDowntime = append(dDowntime, downtime...)
		
		} else if migrationType == "5" {

			downtime := getDataDowntimeOfMigration(evalConfig, migrationID)

			nDowntime = append(nDowntime, downtime...)

		}

	}

	// log.Println(tDowntime)

	WriteStrArrToLog(
		evalConfig.DataDowntimeInStencilFile, 
		ConvertDurationToString(dDowntime),
	)

	WriteStrArrToLog(
		evalConfig.DataDowntimeInNaiveFile, 
		ConvertDurationToString(nDowntime),
	)

}

// The diaspora database needs to be changed to diaspora_1000000_exp
// This is to evaluate the scalability of the migration algorithm with edges
func Exp4() {
	
	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	preExp(evalConfig)

	counterStart := 0
	counterNum := 100
	counterInterval := 10

	userIDWithEdges := getEdgesCounter(evalConfig, 
		counterStart, counterNum, counterInterval)

	// log.Println(userIDWithEdges)

	for i := 0; i < len(userIDWithEdges); i ++ {
		
		userID := userIDWithEdges[i]["person_id"]

		uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
			userID, "diaspora", "1", "mastodon", "2", "d", 1

		enableDisplay, displayInFirstPhase := false, false

		SA1_migrate.Controller(uid, srcAppName, srcAppID, 
			dstAppName, dstAppID, migrationType, threadNum,
			enableDisplay, displayInFirstPhase,
		)
		
		migrationID := getMigrationIDBySrcUserID(evalConfig, userID)
		
		migrationIDInt, err := strconv.Atoi(migrationID)
		if err != nil {
			log.Fatal(err)
		}

		time := GetMigrationTime(
			evalConfig.StencilDBConn,
			migrationIDInt,
		)

		userIDWithEdges[i]["time"] = ConvertSingleDurationToString(time)

		WriteStrToLog(
			"migrationScalabilityEdges",
			ConvertMapStringToJSONString(userIDWithEdges[i]),
		)
	}

	log.Println(userIDWithEdges)

}

// The diaspora database needs to be changed to diaspora_1000000_exp
// This is to evaluate the scalability of the migration algorithm with nodes
func Exp5() {
	
	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	preExp(evalConfig)

	counterStart := 0
	counterNum := 100
	counterInterval := 10

	userIDWithNodes := getNodesCounter(evalConfig, 
		counterStart, counterNum, counterInterval)

	log.Println(userIDWithNodes)
	log.Println(len(userIDWithNodes))

	for i := 0; i < len(userIDWithNodes); i ++ {
		
		userID := userIDWithNodes[i]["person_id"]

		uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
			userID, "diaspora", "1", "mastodon", "2", "d", 1

		enableDisplay, displayInFirstPhase := false, false

		SA1_migrate.Controller(uid, srcAppName, srcAppID, 
			dstAppName, dstAppID, migrationType, threadNum,
			enableDisplay, displayInFirstPhase,
		)
		
		migrationID := getMigrationIDBySrcUserID(evalConfig, userID)
		
		migrationIDInt, err := strconv.Atoi(migrationID)
		if err != nil {
			log.Fatal(err)
		}

		time := GetMigrationTime(
			evalConfig.StencilDBConn,
			migrationIDInt,
		)

		userIDWithNodes[i]["time"] = ConvertSingleDurationToString(time)

		WriteStrToLog(
			"migrationScalabilityNodes",
			ConvertMapStringToJSONString(userIDWithNodes[i]),
		)
	}

	log.Println(userIDWithNodes)

}

func Exp4GetEdgesNodes() {

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	counter := getCounter(evalConfig)

	for _, counter1 := range counter {

		WriteStrToLog(
			"counter",
			ConvertMapStringToJSONString(counter1),
		)

	}

}

func Exp4CountEdgesNodes() {

	appName, appID := "diaspora_1000", "1"

	db.DIASPORA_DB = appName

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	userIDs := getAllUserIDsInDiaspora(evalConfig)

	log.Println("total users:", len(userIDs))

	for _, userID := range userIDs {
		
		res := make(map[string]int)

		apis.StartCounter(appName, appID, userID)

		WriteStrToLog(
			evalConfig.Diaspora1KCounterFile,
			ConvertMapIntToJSONString(res),
			true,
		)
	}

}
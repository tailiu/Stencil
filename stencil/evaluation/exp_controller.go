package evaluation

import (
	"stencil/SA1_migrate"
	"stencil/apis"
	"stencil/db"
	"log"
	"fmt"
	"encoding/json"
	"strconv"
	"time"
)

func preExp(evalConfig *EvalConfig) {

	query1 := `TRUNCATE identity_table, migration_registration, 
		reference_table, resolved_references, txn_logs, 
		evaluation, data_bags, display_flags, display_registration`

	query2 := "SELECT truncate_tables('cow')"

	if err1 := db.TxnExecute1(evalConfig.StencilDBConn, query1); err1 != nil {
		log.Fatal(err1)
	} else {
		if err2 := db.TxnExecute1(evalConfig.StencilDBConn1, query1); err2 != nil {
			log.Fatal(err2)
		} else {
			if err3 := db.TxnExecute1(evalConfig.StencilDBConn2, query1); err3 != nil {
				log.Fatal(err3)
			} else {
				if err4 := db.TxnExecute1(evalConfig.MastodonDBConn, query2); err4 != nil {
					log.Fatal(err4)
				} else {
					if err5 := db.TxnExecute1(evalConfig.MastodonDBConn1, query2); err5 != nil {
						log.Fatal(err5)
					} else {
						if err6 := db.TxnExecute1(evalConfig.MastodonDBConn2, query2); err6 != nil {
							log.Fatal(err6)
						} else {
							return
						}
					}
				}
			}
		}
	}

}

func preExp1(evalConfig *EvalConfig) {

	query1 := `TRUNCATE identity_table, migration_registration, 
		reference_table, resolved_references, txn_logs, 
		evaluation, data_bags, display_flags, display_registration`

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

func PreExp() {

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	preExp(evalConfig)

}

// In this experiment, we migrate 1000 users from Diaspora to Mastodon
// Note that in this exp the migration thread should not migrate data from data bags
// The source database needs to be changed to diaspora_1000_exp
func Exp1(firstUserID ...string) {

	// This is the configuration of the first time test
	// stencilDB = "stencil_cow"
	// mastodon = "mastodon"
	// diaspora = "diaspora_1000_exp"

	stencilDB = "stencil_exp4"
	mastodon = "mastodon_exp4"
	diaspora = "diaspora_1000_exp4"

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	// preExp1(evalConfig)

	// This is the configuration of the first time test
	// db.STENCIL_DB = "stencil_cow"
	// db.DIASPORA_DB = "diaspora_1000_exp"
	// db.MASTODON_DB = "mastodon"

	db.STENCIL_DB = "stencil_exp4"
	db.DIASPORA_DB = "diaspora_1000_exp4"
	db.MASTODON_DB = "mastodon_exp4"

	userIDs := getAllUserIDsInDiaspora(evalConfig)

	shuffleSlice(userIDs)

	if len(firstUserID) != 0 {
		userIDs = moveElementToStartOfSlice(userIDs, firstUserID[0])
	}

	log.Println("Total users:", len(userIDs))

	for _, userID := range userIDs {

		log.Println("User ID:", userID)

		uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
			userID, "diaspora", "1", "mastodon", "2", "d", 1

		enableDisplay, displayInFirstPhase := true, true

		SA1_migrate.Controller(uid, srcAppName, srcAppID, 
			dstAppName, dstAppID, migrationType, threadNum,
			enableDisplay, displayInFirstPhase,
		)

		log.Println("************ Calculate Dangling Data Size ************")

		refreshEvalConfigDBConnections(evalConfig)

		migrationID := getMigrationIDBySrcUserID(evalConfig, userID)

		danglingData := make(map[string]int64)

		srcDanglingData, dstDanglingData :=
			getDanglingDataSizeOfMigration(evalConfig, migrationID)

		danglingData["userID"] = ConvertStringtoInt64(userID)
		danglingData["srcDanglingData"] = srcDanglingData
		danglingData["dstDanglingData"] = dstDanglingData

		WriteStrToLog(
			evalConfig.DanglingDataFile,
			ConvertMapInt64ToJSONString(danglingData),
		)

	}

}

func Exp1GetDanglingDataSize(migrationID string) {

	stencilDB = "stencil_exp4"
	mastodon = "mastodon_exp4"
	diaspora = "diaspora_1000_exp4"

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	srcDanglingData, dstDanglingData :=
		getDanglingDataSizeOfMigration(evalConfig, migrationID)
	
	log.Println("Migration ID:", migrationID)
	log.Println("Src dangling data size:", srcDanglingData)
	log.Println("Dst dangling data size:", dstDanglingData)

}

// In diaspora_1000 database:
// Total Media Size in Diaspora: 793878636 bytes
// All Rows Size in Diaspora: 30840457 bytes
// Total Size in Diaspora: 824719093 bytes
// In mastodon database:
// Total Media Size in Mastodon: 789827483 bytes
// All Rows Size in Mastodon: 16585552 bytes
// Dangling Data Size in Mastodon: 330573 bytes
// Total Size in Mastodon: 806743608 bytes
func Exp1GetTotalMigratedDataSize() {

	diaspora = "diaspora_1000"
	stencilDB = "stencil_cow"
	// Note that mastodon needs to be changed in the config file as well
	mastodon = "mastodon"

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	mediaSizeInDiaspora := getAllMediaSize(evalConfig.DiasporaDBConn, 
		"photos", evalConfig.DiasporaAppID)

	log.Println("Total Media Size in Diaspora:", mediaSizeInDiaspora, "bytes")
	
	rowsSizeInDiaspora := getAllRowsSize(evalConfig.DiasporaDBConn)

	log.Println("All Rows Size in Diaspora:", rowsSizeInDiaspora, "bytes")

	log.Println("Total Size in Diaspora:", mediaSizeInDiaspora + rowsSizeInDiaspora, "bytes")

	mediaSizeInMastodon := getAllMediaSize(evalConfig.MastodonDBConn, 
		"media_attachments", evalConfig.MastodonAppID)
	
	log.Println("Total Media Size in Mastodon:", mediaSizeInMastodon, "bytes")
	
	rowsSizeInMastodon := getAllRowsSize(evalConfig.MastodonDBConn)

	log.Println("All Rows Size in Mastodon:", rowsSizeInMastodon, "bytes")

	danglingDataSizeInMastodon := getDanglingDataSizeOfApp(evalConfig, evalConfig.MastodonAppID)

	log.Println("Dangling Data Size in Mastodon:", danglingDataSizeInMastodon, "bytes")

	log.Println("Total Size in Mastodon:", 
		mediaSizeInMastodon + rowsSizeInMastodon + danglingDataSizeInMastodon,
		"bytes",
	)

}

func Exp1GetDanglingObjects() {

	stencilDB = "stencil_cow"
	mastodon = "mastodon"
	diaspora = "diaspora_1000_exp"

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	migrationIDs := GetAllMigrationIDsOrderByEndTime(evalConfig)

	// log.Println(migrationIDs)

	for _, migrationID := range migrationIDs {

		danglingObjects := make(map[string]int64)

		srcDanglingObjects, dstDanglingObjects :=
			getDanglingObjectsOfMigration(evalConfig, migrationID)

		danglingObjects["srcDanglingObjs"] = srcDanglingObjects
		danglingObjects["dstDanglingObjs"] = dstDanglingObjects

		WriteStrToLog(
			evalConfig.DanglingObjectsFile,
			ConvertMapInt64ToJSONString(danglingObjects),
		)

	}
}

func Exp1GetTotalObjects() {
	
	diaspora = "diaspora_1000"
	stencilDB = "stencil_cow"
	// Note that mastodon needs to be changed in the config file as well
	mastodon = "mastodon"

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	totalRowCounts1 := getTotalRowCountsOfDB(evalConfig.DiasporaDBConn)

	totalPhotoCounts1 := getTotalRowCountsOfTable(evalConfig.DiasporaDBConn, "photos")

	log.Println("Total Objects in Diaspora:", totalRowCounts1 + totalPhotoCounts1)

	totalRowCounts2 := getTotalRowCountsOfDB(evalConfig.MastodonDBConn)

	totalPhotoCounts2 := getTotalRowCountsOfTable(evalConfig.MastodonDBConn, "media_attachments")

	danglingObjs2 := getDanglingObjectsOfApp(evalConfig, "2")

	log.Println("Total Objects in Mastodon:", 
		totalRowCounts2 + totalPhotoCounts2 + danglingObjs2)

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

	diaspora = "diaspora_1000000"

	stencilDB = "stencil_exp"
	stencilDB1 = "stencil_exp1"
	stencilDB2 = "stencil_exp2"

	mastodon = "mastodon_exp"
	mastodon1 = "mastodon_exp1"
	mastodon2 = "mastodon_exp2"

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	// preExp(evalConfig)

	migrationNum := 300

	// startNum := 200 // first time and crash at the 69th user
	// startNum := 300 // second time and crash at the 67th user
	// startNum := 400 // third time and stop at the 14th user
	// startNum := 600 // fouth time and stop at the 1st user
	// startNum := 900 // fifth time and stop at the 10th user
	// startNum := 920 // sixth time and crashes at the 52th user
	// startNum := 1500 // seventh time and crashes at the 11th user

	startNum := 1520

	// ************ SA1 ************

	SA1MigrationType := "d"

	SA1StencilDB, SA1SrcDB, SA1DstDB := 
		"stencil_exp", "diaspora_1000000_exp", "mastodon_exp"

	SA1EnableDisplay, SA1DisplayInFirstPhase := true, true

	SA1SizeFile, SA1TimeFile := "SA1Size", "SA1Time"

	// ************ SA1 without Display ************

	SA1WithoutDisplayMigrationType := "d"

	SA1WithoutDisplayStencilDB, SA1WithoutDisplaySrcDB, SA1WithoutDisplayDstDB := 
		"stencil_exp1", "diaspora_1000000_exp1", "mastodon_exp1"

	SA1WithoutDisplayEnableDisplay, SA1WithoutDisplayDisplayInFirstPhase := false, false

	SA1WithoutDisplaySizeFile, SA1WithoutDisplayTimeFile := "SA1WDSize", "SA1WDTime"

	// ************ Naive Migration ************

	naiveMigrationType := "n"

	naiveStencilDB, naiveSrcDB, naiveDstDB := 
		"stencil_exp2", "diaspora_1000000_exp2", "mastodon_exp2"

	naiveEnableDisplay, naiveDisplayInFirstPhase := false, false

	naiveSizeFile, naiveTimeFile := "naiveSize", "naiveTime"


	userIDs := getAllUserIDsSortByPhotosInDiaspora(evalConfig)

	// log.Println(userIDs)

	for i := startNum; i < migrationNum + startNum; i ++ {

		userID := userIDs[i]["author_id"]

		log.Println("User ID:", userID)

		// ************ SA1 ************

		migrateUserFromDiasporaToMastodon(
			evalConfig, SA1StencilDB, diaspora, 
			userID, SA1MigrationType, 
			SA1StencilDB, SA1SrcDB, SA1DstDB,
			SA1SizeFile, SA1TimeFile,
			SA1EnableDisplay, SA1DisplayInFirstPhase,
		)

		log.Println("User ID:", userID)

		// ************ SA1 without Display ************

		migrateUserFromDiasporaToMastodon(
			evalConfig, SA1WithoutDisplayStencilDB, diaspora, 
			userID, SA1WithoutDisplayMigrationType, 
			SA1WithoutDisplayStencilDB, SA1WithoutDisplaySrcDB, SA1WithoutDisplayDstDB,
			SA1WithoutDisplaySizeFile, SA1WithoutDisplayTimeFile,
			SA1WithoutDisplayEnableDisplay, SA1WithoutDisplayDisplayInFirstPhase,
		)

		log.Println("User ID:", userID)
		
		// ************ Naive Migration ************

		migrateUserFromDiasporaToMastodon(
			evalConfig, naiveStencilDB, diaspora, 
			userID, naiveMigrationType, 
			naiveStencilDB, naiveSrcDB, naiveDstDB,
			naiveSizeFile, naiveTimeFile,
			naiveEnableDisplay, naiveDisplayInFirstPhase,
		)

	}

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
			evalConfig.StencilDBConn,
			evalConfig.DiasporaDBConn,
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

	// preExp(evalConfig)

	// migrationNum := 300

	SA1StencilDB, SA1SrcDB, SA1DstDB := 
		"stencil_cow", "diaspora_1000000_exp", "mastodon"
	
	naiveStencilDB, naiveSrcDB, naiveDstDB := 
		"stencil_exp", "diaspora_1000000_exp1", "mastodon_exp"

	SA1EnableDisplay, SA1DisplayInFirstPhase := true, true

	naiveEnableDisplay, naiveDisplayInFirstPhase := true, false

	userID := "1"
	
	migrateUserUsingSA1AndNaive(evalConfig, 
		SA1StencilDB, SA1SrcDB, SA1DstDB, userID,
		naiveStencilDB, naiveSrcDB, naiveDstDB, 
		SA1EnableDisplay, SA1DisplayInFirstPhase,
		naiveEnableDisplay, naiveDisplayInFirstPhase,
	)

}

func Exp3LoadUserIDsFromLog() {

	SA1MigrationFile := "SA1Size"

	naiveStencilDB, naiveSrcDB, naiveDstDB := 
		"stencil_exp5", "diaspora_1000000_exp5", "mastodon_exp5"

	migrationType := "n"

	naiveEnableDisplay, naiveDisplayInFirstPhase := true, false

	data := ReadStrLinesFromLog(SA1MigrationFile)

	log.Println("Migration number:", len(data))

	log.Println(data)

	for _, data1 := range data {
		
		var sizeData SA1SizeStruct

		err := json.Unmarshal([]byte(data1), &sizeData)
		if err != nil {
			log.Fatal(err)
		}

		userID := sizeData.UserID

		log.Println("UserID is", userID)

		migrateUserFromDiasporaToMastodon1(
			userID, migrationType, 
			naiveStencilDB, naiveSrcDB, naiveDstDB, 
			naiveEnableDisplay, naiveDisplayInFirstPhase,
		)
	}

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

			downtime := getDataDowntimeOfMigration(evalConfig.StencilDBConn,
				migrationID)

			dDowntime = append(dDowntime, downtime...)
		
		} else if migrationType == "5" {

			downtime := getDataDowntimeOfMigration(evalConfig.StencilDBConn,
				migrationID)

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

func Exp3GetDatadowntimeByLoadingUserIDFromLog() {

	SA1StencilDB := "stencil_exp"
	naiveStencilDB := "stencil_exp5"

	logFile := "SA1Size"
	migrationNum := 100

	stencilDB = SA1StencilDB
	stencilDB1 = naiveStencilDB

	evalConfig := InitializeEvalConfig()
	defer closeDBConns(evalConfig)

	data := ReadStrLinesFromLog(logFile)

	var SA1Downtime, naiveDowntime []time.Duration

	for i := 0; i < migrationNum; i++ {

		var data1 SA1SizeStruct 

		err := json.Unmarshal([]byte(data[i]), &data1)
		if err != nil {
			log.Fatal(err)
		}

		userID := data1.UserID

		SA1MigrationID := getMigrationIDBySrcUserID1(evalConfig.StencilDBConn, userID)
		naiveMigrationID := getMigrationIDBySrcUserID1(evalConfig.StencilDBConn1, userID)

		SA1Downtime1 := getDataDowntimeOfMigration(evalConfig.StencilDBConn,
			SA1MigrationID)
		naiveDowntime1 := getDataDowntimeOfMigration(evalConfig.StencilDBConn1,
			naiveMigrationID)

		SA1Downtime = append(SA1Downtime, SA1Downtime1...)
		naiveDowntime = append(naiveDowntime, naiveDowntime1...)

	}

	log.Println(SA1Downtime)

	WriteStrArrToLog(
		evalConfig.DataDowntimeInStencilFile, 
		ConvertDurationToString(SA1Downtime),
	)

	WriteStrArrToLog(
		evalConfig.DataDowntimeInNaiveFile, 
		ConvertDurationToString(naiveDowntime),
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

// This function is for us to get nodes and edges from database to plot 
// the relationship between them
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

func Exp4Count1MDBEdgesNodes() {

	appName, appID := "diaspora", "1"
	
	db.DIASPORA_DB = "diaspora_1000000_counter"
	db.STENCIL_DB = "stencil_counter"

	stencilDB = "stencil_counter"
	diaspora = "diaspora_1000000_counter"

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	userIDs := getAllUserIDsInDiaspora(evalConfig, true)

	file := evalConfig.Diaspora1MCounterFile

	counter := getCounter(evalConfig)

	log.Println("Total user:", len(userIDs))

	// for i := len(userIDs) -  1; i > 10000; i-- {  
	// for _, userID := range userIDs {
	for i := 280000; i < len(userIDs); i += 100 {  
		
		userID := userIDs[i]

		log.Println("Counting userID:", userID)

		if isAlreadyCounted(counter, userID) {
			log.Println("userID", userID, "has already been counted")
			continue
		}

		res := make(map[string]int)

		nodeCount, edgeCount := apis.StartCounter(appName, appID, userID)

		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			log.Fatal(err)
		}

		res["userID"] = userIDInt
		res["nodes"] = nodeCount
		res["edges"] = edgeCount

		WriteStrToLog(
			file,
			ConvertMapIntToJSONString(res),
			true,
		)
	}

}

func Exp4LoadCounterResToTable() {

	stencilDB = "stencil_counter"
	counterFile := "diaspora1MCounter"
	counterTable := "dag_counter"

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	data := ReadStrLinesFromLog(counterFile, true)

	// log.Println(data)

	for _, data1 := range data {
		
		var counter1 Counter

		// log.Println(data1)

		err := json.Unmarshal([]byte(data1), &counter1)
		if err != nil {
			log.Fatal(err)
		}

		// log.Println(counter1.UserID)
		
		insertDataIntoCounterTableIfNotExist(evalConfig,
			counterTable, counter1)

	}
	
}

func Exp4CountEdgesNodes() {

	appName, appID := "diaspora", "1"

	db.DIASPORA_DB = "diaspora_100000"
	db.STENCIL_DB = "stencil_cow"
	diaspora = "diaspora_100000"
	
	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	userIDs := getAllUserIDsInDiaspora(evalConfig)

	// log.Println("total users:", len(userIDs))

	file := evalConfig.Diaspora100KCounterFile

	for _, userID := range userIDs {

		res := make(map[string]int)

		nodeCount, edgeCount := apis.StartCounter(appName, appID, userID)

		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			log.Fatal(err)
		}

		res["userID"] = userIDInt
		res["nodes"] = nodeCount
		res["edges"] = edgeCount

		WriteStrToLog(
			file,
			ConvertMapIntToJSONString(res),
			true,
		)
	}

}

func Exp6() {

	stencilDB = "stencil_exp3"
	mastodon = "mastodon_exp3"
	diaspora = "diaspora_1000000_exp3"

	db.STENCIL_DB = "stencil_exp3"
	db.DIASPORA_DB = "diaspora_1000000_exp3"
	db.MASTODON_DB = "mastodon_exp3"

	// counterStart := 0
	// counterNum := 300
	// counterInterval := 10

	counterStart := 0
	counterNum := 100
	counterInterval := 10

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	preExp1(evalConfig)

	res := getEdgesCounter(evalConfig, 
		counterStart, counterNum, counterInterval)

	log.Println("Total Num:", len(res))
	log.Println(res)

	for i := 0; i < len(res); i ++ {

		res1 := res[i]

		userID := res1["person_id"]

		uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
			userID, "diaspora", "1", "mastodon", "2", "d", 1

		enableDisplay, displayInFirstPhase, markAsDelete := true, false, true

		SA1_migrate.Controller(uid, srcAppName, srcAppID, 
			dstAppName, dstAppID, migrationType, threadNum,
			enableDisplay, displayInFirstPhase, markAsDelete,
		)
		
		log.Println("************ Calculate Migration and Display Time ************")

		refreshEvalConfigDBConnections(evalConfig)

		migrationID := getMigrationIDBySrcUserID(evalConfig, userID)
		
		migrationIDInt, err := strconv.Atoi(migrationID)
		if err != nil {
			log.Fatal(err)
		}

		mTime := GetMigrationTime(
			evalConfig.StencilDBConn,
			migrationIDInt,
		)

		dTime := GetDisplayTime(
			evalConfig.StencilDBConn,
			migrationID,
		)

		res1["migrationTime"] = ConvertSingleDurationToString(mTime)
		res1["displayTime"] = ConvertSingleDurationToString(dTime)

		log.Println("************ Calculate Nodes and Edges after Migration ************")

		migratedUserID := getMigratedUserID(evalConfig, migrationID, dstAppID)

		nodeCount, edgeCount := apis.StartCounter(dstAppName, dstAppID, 
			migratedUserID, true)

		res1["nodesAfterMigration"] = strconv.Itoa(nodeCount)
		res1["edgesAfterMigration"] = strconv.Itoa(edgeCount)

		WriteStrToLog(
			"scalability",
			ConvertMapStringToJSONString(res1),
		)
	}

}
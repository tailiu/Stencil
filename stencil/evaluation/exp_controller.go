package evaluation

import (
	"stencil/SA1_migrate"
	"stencil/apis"
	"stencil/db"
	"database/sql"
	"log"
	"fmt"
	"encoding/json"
	"strconv"
	"time"
)

func preExp(evalConfig *EvalConfig) {

	query1 := `TRUNCATE identity_table, migration_registration, 
		reference_table, resolved_references, txn_logs, id_changes,
		evaluation, data_bags, display_flags, display_registration,
		attribute_changes, reference_table_v2`

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
		reference_table, resolved_references, txn_logs, id_changes,
		evaluation, data_bags, display_flags, display_registration,
		attribute_changes, reference_table_v2`

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

func preExp2(evalConfig *EvalConfig) {

	query1 := `TRUNCATE identity_table, migration_registration, 
		reference_table, resolved_references, txn_logs, id_changes,
		evaluation, data_bags, display_flags, display_registration,
		attribute_changes, reference_table_v2`

	query2 := "SELECT truncate_tables('cow')"

	if err1 := db.TxnExecute1(evalConfig.StencilDBConn, query1); err1 != nil {
		log.Fatal(err1)
	} else {
		if err2 := db.TxnExecute1(evalConfig.StencilDBConn1, query1); err2 != nil {
			log.Fatal(err2)
		} else {
			if err4 := db.TxnExecute1(evalConfig.MastodonDBConn, query2); err4 != nil {
				log.Fatal(err4)
			} else {
				if err5 := db.TxnExecute1(evalConfig.MastodonDBConn1, query2); err5 != nil {
					log.Fatal(err5)
				} else {
					return
				}
			}
		}
	}

}

func PreExp() {

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	preExp(evalConfig)

}

func preExp7(evalConfig *EvalConfig, startApp string) {

	query1 := `TRUNCATE identity_table, migration_registration, 
		reference_table, resolved_references, txn_logs, id_changes,
		evaluation, data_bags, display_flags, display_registration,
		attribute_changes, reference_table_v2`

	query2 := "SELECT truncate_tables('cow')"

	query3 := `TRUNCATE TABLE account_migrations, account_deletions, 
		ar_internal_metadata, blocks, chat_offline_messages, like_signatures, 
		locations, invitation_codes, o_embed_caches, open_graph_caches, 
		o_auth_access_tokens, comments, poll_participation_signatures, 
		polls, ppid, poll_answers, schema_migrations, "references", 
		roles, services, reports, signature_orders, tags, simple_captcha_data, 
		user_preferences, conversation_visibilities, share_visibilities, tag_followings, 
		posts, chat_fragments, chat_contacts, poll_participations, aspect_memberships, 
		taggings, aspects, notification_actors, likes, notifications, contacts, pods, 
		comment_signatures, people, authorizations, messages, aspect_visibilities, 
		conversations, mentions, profiles, o_auth_applications, 
		participations, users, photos CASCADE`

	if err1 := db.TxnExecute1(evalConfig.StencilDBConn, query1); err1 != nil {
		log.Fatal(err1)	
	} else {
		if err2 := db.TxnExecute1(evalConfig.StencilDBConn1, query1); err2 != nil {
			log.Fatal(err2)
		} 
	}

	if startApp != "diaspora" {
		if err3 := db.TxnExecute1(evalConfig.DiasporaDBConn, query3); err3 != nil {
			log.Fatal(err3)
		} else {
			if err4 := db.TxnExecute1(evalConfig.DiasporaDBConn1, query3); err4 != nil {
				log.Fatal(err4)
			}
		}
	}

	if startApp != "mastodon" {
		if err5 := db.TxnExecute1(evalConfig.MastodonDBConn, query2); err5 != nil {
			log.Fatal(err5)
		} else {
			if err6 := db.TxnExecute1(evalConfig.MastodonDBConn1, query2); err6 != nil {
				log.Fatal(err6)
			}
		}
	}

	if startApp != "twitter" {
		if err7 := db.TxnExecute1(evalConfig.TwitterDBConn, query2); err7 != nil {
			log.Fatal(err7)
		} else {
			if err8 := db.TxnExecute1(evalConfig.TwitterDBConn1, query2); err8 != nil {
				log.Fatal(err8)
			} 
		}
	}

	if startApp != "gnusocial" {
		if err9 := db.TxnExecute1(evalConfig.GnusocialDBConn, query2); err9 != nil {
			log.Fatal(err9)
		} else {
			if err10 := db.TxnExecute1(evalConfig.GnusocialDBConn1, query2); err10 != nil {
				log.Fatal(err10)
			} 
		}
	}
}

func preExp9(evalConfig *EvalConfig) {

	query1 := `TRUNCATE identity_table, migration_registration, 
		reference_table, resolved_references, txn_logs, id_changes,
		evaluation, data_bags, display_flags, display_registration,
		attribute_changes, reference_table_v2`

	query2 := "SELECT truncate_tables('cow')"

	if err1 := db.TxnExecute1(evalConfig.StencilDBConn, query1); err1 != nil {
		log.Fatal(err1)	
	} else {
		if err2 := db.TxnExecute1(evalConfig.MastodonDBConn, query2); err2 != nil {
			log.Fatal(err2)
		} else {
			if err3 := db.TxnExecute1(evalConfig.TwitterDBConn, query2); err3 != nil {
				log.Fatal(err3)
			} else {
				if err4 := db.TxnExecute1(evalConfig.GnusocialDBConn, query2); err4 != nil {
					log.Fatal(err4)
				} else {
					if err5 := db.TxnExecute1(evalConfig.StencilDBConn1, query1); err5 != nil {
						log.Fatal(err5)
					} else {
						if err6 := db.TxnExecute1(evalConfig.MastodonDBConn1, query2); err6 != nil {
							log.Fatal(err6)
						} else {
							if err7 := db.TxnExecute1(evalConfig.TwitterDBConn1, query2); err7 != nil {
								log.Fatal(err7)
							} else {
								if err8 := db.TxnExecute1(evalConfig.GnusocialDBConn1, query2); err8 != nil {
									log.Fatal(err8)
								} else {
									return
								}
							}
						}
					}
				}
			}
		}
	}

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

	migrationNum := 1000

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

	for i := 0; i < migrationNum; i++ {
	// for _, userID := range userIDs {

		userID := userIDs[i]

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

	// stencilDB = "stencil_cow"
	// mastodon = "mastodon"
	// diaspora = "diaspora_1000_exp"

	stencilDB = "stencil_exp4"

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

func Exp1GetDanglingObjectsV2() {

	// stencilDB = "stencil_cow"
	// mastodon = "mastodon"
	// diaspora = "diaspora_1000_exp"

	stencilDB = "stencil_exp4"

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	migrationIDs := GetAllMigrationIDsOrderByEndTime(evalConfig)

	// log.Println(migrationIDs)

	for _, migrationID := range migrationIDs {

		danglingObjects := make(map[string]int64)

		srcUserID := getSrcUserIDByMigrationID(evalConfig.StencilDBConn,
			migrationID)

		srcDanglingObjects, dstDanglingObjects :=
			getDanglingObjectsOfMigrationV2(evalConfig, migrationID, srcUserID)

		danglingObjects["srcDanglingObjs"] = srcDanglingObjects
		danglingObjects["dstDanglingObjs"] = dstDanglingObjects

		WriteStrToLog(
			evalConfig.DanglingObjectsFile,
			ConvertMapInt64ToJSONString(danglingObjects),
		)

	}
}

func Exp1GetTotalObjects() {
	
	// diaspora = "diaspora_1000"
	// stencilDB = "stencil_cow"
	// // Note that mastodon needs to be changed in the config file as well
	// mastodon = "mastodon"

	diaspora = "diaspora_1000"
	stencilDB = "stencil_exp4"
	mastodon = "mastodon_exp4"

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	totalRowCounts1 := getTotalRowCountsOfDB(evalConfig.DiasporaDBConn)

	totalPhotoCounts1 := getTotalRowCountsOfTable(evalConfig.DiasporaDBConn, "photos")

	log.Println("Total Objects without considering media in Diaspora:", totalRowCounts1)

	log.Println("Total Objects considering media in Diaspora:", totalRowCounts1 + totalPhotoCounts1)

	totalRowCounts2 := getTotalRowCountsOfDB(evalConfig.MastodonDBConn)

	totalPhotoCounts2 := getTotalRowCountsOfTable(evalConfig.MastodonDBConn, "media_attachments")

	danglingObjs2 := getDanglingObjectsOfApp(evalConfig, "2")

	log.Println("Total Objects in Mastodon without dangling objects or media:", 
		totalRowCounts2)

	log.Println("Total Objects in Mastodon without dangling objects:", 
		totalRowCounts2 + totalPhotoCounts2)

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

	log.Println("===================================")
	log.Println("Starting Exp2: Migration Rates Test")
	log.Println("===================================")

	// diaspora = "diaspora_1000000"

	// stencilDB = "stencil_exp"
	// stencilDB1 = "stencil_exp1"
	// stencilDB2 = "stencil_exp2"

	// mastodon = "mastodon_exp"
	// mastodon1 = "mastodon_exp1"
	// mastodon2 = "mastodon_exp2"

	diaspora = "diaspora_100000"

	stencilDB = "stencil_100k_exp1"
	stencilDB1 = "stencil_100k_exp2"
	stencilDB2 = "stencil_100k_exp3"

	mastodon = "mastodon_100k_exp1"
	mastodon1 = "mastodon_100k_exp2"
	mastodon2 = "mastodon_100k_exp3"

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	preExp(evalConfig)

	edgeCounterRangeStart := 750
	edgeCounterRangeEnd := 1200
	migrationNum := 100

	edgeCounter := getEdgesCounterByRange(
		evalConfig,
		edgeCounterRangeStart, 
		edgeCounterRangeEnd, 
		migrationNum,
	)

	log.Println(edgeCounter)
	log.Println("Migration number:", len(edgeCounter))

	// ************ SA1 ************

	SA1MigrationType := "d"

	// SA1StencilDB, SA1SrcDB, SA1DstDB := 
	// 	"stencil_exp", "diaspora_1000000_exp", "mastodon_exp"

	SA1StencilDB, SA1SrcDB, SA1DstDB := 
		"stencil_100k_exp1", "diaspora_100k_exp1", "mastodon_100k_exp1"

	SA1EnableDisplay, SA1DisplayInFirstPhase := true, true

	SA1SizeFile, SA1TimeFile := "SA1Size", "SA1Time"

	// // ************ SA1 without Display ************

	// SA1WithoutDisplayMigrationType := "d"

	// // SA1WithoutDisplayStencilDB, SA1WithoutDisplaySrcDB, SA1WithoutDisplayDstDB := 
	// // 	"stencil_exp1", "diaspora_1000000_exp1", "mastodon_exp1"

	// SA1WithoutDisplayStencilDB, SA1WithoutDisplaySrcDB, SA1WithoutDisplayDstDB := 
	// 	"stencil_100k_exp2", "diaspora_100k_exp2", "mastodon_100k_exp2"

	// SA1WithoutDisplayEnableDisplay, SA1WithoutDisplayDisplayInFirstPhase := false, false

	// SA1WithoutDisplaySizeFile, SA1WithoutDisplayTimeFile := "SA1WDSize", "SA1WDTime"

	// ************ SA1 without Display ************

	SA1IndependentMigrationType := "i"

	// SA1WithoutDisplayStencilDB, SA1WithoutDisplaySrcDB, SA1WithoutDisplayDstDB := 
	// 	"stencil_exp1", "diaspora_1000000_exp1", "mastodon_exp1"

	SA1IndependentStencilDB, SA1IndependentSrcDB, SA1IndependentDstDB := 
		"stencil_100k_exp2", "diaspora_100k_exp2", "mastodon_100k_exp2"

	SA1IndependentEnableDisplay, SA1IndependentDisplayInFirstPhase := true, true

	SA1IndependentSizeFile, SA1IndependentTimeFile := "SA1IndepSize", "SAIndepTime"

	// ************ Naive Migration ************

	naiveMigrationType := "n"

	// naiveStencilDB, naiveSrcDB, naiveDstDB := 
	// 	"stencil_exp2", "diaspora_1000000_exp2", "mastodon_exp2"

	naiveStencilDB, naiveSrcDB, naiveDstDB := 
		"stencil_100k_exp3", "diaspora_100k_exp3", "mastodon_100k_exp3"

	naiveEnableDisplay, naiveDisplayInFirstPhase := false, false

	naiveSizeFile, naiveTimeFile := "naiveSize", "naiveTime"


	for i := 0; i < len(edgeCounter); i ++ {

		userID := edgeCounter[i]["person_id"]

		log.Println("User ID:", userID)

		// ************ SA1 Deletion ************

		migrateUserFromDiasporaToMastodon(
			evalConfig, SA1StencilDB, diaspora, 
			userID, SA1MigrationType, 
			SA1StencilDB, SA1SrcDB, SA1DstDB,
			SA1SizeFile, SA1TimeFile,
			SA1EnableDisplay, SA1DisplayInFirstPhase,
		)

		log.Println("User ID:", userID)

		// // ************ SA1 without Display ************

		// migrateUserFromDiasporaToMastodon(
		// 	evalConfig, SA1WithoutDisplayStencilDB, diaspora, 
		// 	userID, SA1WithoutDisplayMigrationType, 
		// 	SA1WithoutDisplayStencilDB, SA1WithoutDisplaySrcDB, SA1WithoutDisplayDstDB,
		// 	SA1WithoutDisplaySizeFile, SA1WithoutDisplayTimeFile,
		// 	SA1WithoutDisplayEnableDisplay, SA1WithoutDisplayDisplayInFirstPhase,
		// )

		// log.Println("User ID:", userID)
		
		// ************ SA1 Independent Migratoin ************

		migrateUserFromDiasporaToMastodon(
			evalConfig, SA1IndependentStencilDB, diaspora, 
			userID, SA1IndependentMigrationType, 
			SA1IndependentStencilDB, SA1IndependentSrcDB, SA1IndependentDstDB,
			SA1IndependentSizeFile, SA1IndependentTimeFile,
			SA1IndependentEnableDisplay, SA1IndependentDisplayInFirstPhase,
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

func Exp2GetUserIDsByPhotos() {

	log.Println("===================================")
	log.Println("Starting Exp2: Migration Rates Test")
	log.Println("===================================")

	// diaspora = "diaspora_1000000"

	// stencilDB = "stencil_exp"
	// stencilDB1 = "stencil_exp1"
	// stencilDB2 = "stencil_exp2"

	// mastodon = "mastodon_exp"
	// mastodon1 = "mastodon_exp1"
	// mastodon2 = "mastodon_exp2"

	diaspora = "diaspora_100000"

	stencilDB = "stencil_100k_exp1"
	stencilDB1 = "stencil_100k_exp2"
	stencilDB2 = "stencil_100k_exp3"

	mastodon = "mastodon_100k_exp1"
	mastodon1 = "mastodon_100k_exp2"
	mastodon2 = "mastodon_100k_exp3"

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	preExp(evalConfig)

	startNum := 300 
	migrationNum := 100

	userIDs := getAllUserIDsSortByPhotosInDiaspora(evalConfig)
	log.Println(userIDs)

	// startNum := 200 // first time and crash at the 69th user
	// startNum := 300 // second time and crash at the 67th user
	// startNum := 400 // third time and stop at the 14th user
	// startNum := 600 // fouth time and stop at the 1st user
	// startNum := 900 // fifth time and stop at the 10th user
	// startNum := 920 // sixth time and crashes at the 52th user
	// startNum := 1500 // seventh time and crashes at the 11th user
	// startNum := 1520 // eighth time and stop at the 61th user
	// startNum := 2000 // ninth time and crash at the 4th user

	// startNum := 100 // first time and crash at the 4th user
	// startNum := 105 

	// ************ SA1 ************

	SA1MigrationType := "d"

	// SA1StencilDB, SA1SrcDB, SA1DstDB := 
	// 	"stencil_exp", "diaspora_1000000_exp", "mastodon_exp"

	SA1StencilDB, SA1SrcDB, SA1DstDB := 
		"stencil_100k_exp1", "diaspora_100k_exp1", "mastodon_100k_exp1"

	SA1EnableDisplay, SA1DisplayInFirstPhase := true, true

	SA1SizeFile, SA1TimeFile := "SA1Size", "SA1Time"

	// // ************ SA1 without Display ************

	// SA1WithoutDisplayMigrationType := "d"

	// // SA1WithoutDisplayStencilDB, SA1WithoutDisplaySrcDB, SA1WithoutDisplayDstDB := 
	// // 	"stencil_exp1", "diaspora_1000000_exp1", "mastodon_exp1"

	// SA1WithoutDisplayStencilDB, SA1WithoutDisplaySrcDB, SA1WithoutDisplayDstDB := 
	// 	"stencil_100k_exp2", "diaspora_100k_exp2", "mastodon_100k_exp2"

	// SA1WithoutDisplayEnableDisplay, SA1WithoutDisplayDisplayInFirstPhase := false, false

	// SA1WithoutDisplaySizeFile, SA1WithoutDisplayTimeFile := "SA1WDSize", "SA1WDTime"

	// ************ SA1 without Display ************

	SA1IndependentMigrationType := "i"

	// SA1WithoutDisplayStencilDB, SA1WithoutDisplaySrcDB, SA1WithoutDisplayDstDB := 
	// 	"stencil_exp1", "diaspora_1000000_exp1", "mastodon_exp1"

	SA1IndependentStencilDB, SA1IndependentSrcDB, SA1IndependentDstDB := 
		"stencil_100k_exp2", "diaspora_100k_exp2", "mastodon_100k_exp2"

	SA1IndependentEnableDisplay, SA1IndependentDisplayInFirstPhase := true, true

	SA1IndependentSizeFile, SA1IndependentTimeFile := "SA1IndepSize", "SAIndepTime"

	// ************ Naive Migration ************

	naiveMigrationType := "n"

	// naiveStencilDB, naiveSrcDB, naiveDstDB := 
	// 	"stencil_exp2", "diaspora_1000000_exp2", "mastodon_exp2"

	naiveStencilDB, naiveSrcDB, naiveDstDB := 
		"stencil_100k_exp3", "diaspora_100k_exp3", "mastodon_100k_exp3"

	naiveEnableDisplay, naiveDisplayInFirstPhase := false, false

	naiveSizeFile, naiveTimeFile := "naiveSize", "naiveTime"


	for i := startNum; i < migrationNum + startNum; i ++ {

		userID := userIDs[i]["author_id"]

		log.Println("User ID:", userID)

		// ************ SA1 Deletion ************

		migrateUserFromDiasporaToMastodon(
			evalConfig, SA1StencilDB, diaspora, 
			userID, SA1MigrationType, 
			SA1StencilDB, SA1SrcDB, SA1DstDB,
			SA1SizeFile, SA1TimeFile,
			SA1EnableDisplay, SA1DisplayInFirstPhase,
		)

		log.Println("User ID:", userID)

		// // ************ SA1 without Display ************

		// migrateUserFromDiasporaToMastodon(
		// 	evalConfig, SA1WithoutDisplayStencilDB, diaspora, 
		// 	userID, SA1WithoutDisplayMigrationType, 
		// 	SA1WithoutDisplayStencilDB, SA1WithoutDisplaySrcDB, SA1WithoutDisplayDstDB,
		// 	SA1WithoutDisplaySizeFile, SA1WithoutDisplayTimeFile,
		// 	SA1WithoutDisplayEnableDisplay, SA1WithoutDisplayDisplayInFirstPhase,
		// )

		// log.Println("User ID:", userID)
		
		// ************ SA1 Independent Migratoin ************

		migrateUserFromDiasporaToMastodon(
			evalConfig, SA1IndependentStencilDB, diaspora, 
			userID, SA1IndependentMigrationType, 
			SA1IndependentStencilDB, SA1IndependentSrcDB, SA1IndependentDstDB,
			SA1IndependentSizeFile, SA1IndependentTimeFile,
			SA1IndependentEnableDisplay, SA1IndependentDisplayInFirstPhase,
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

func Exp2LoadUserIDsToAddIndepFromLog() {

	log.Println("========================================")
	log.Println("Starting Independent Migrations for Exp2")
	log.Println("========================================")

	diaspora = "diaspora_1000000"
	stencilDB = "stencil_exp14"

	indepSizeFile := "SA1IndepSize"
	indepTimeFile := "SA1IndepTime"

	SA1MigrationFile := "SA1Size"

	indepStencilDB, indepSrcDB, indepDstDB := 
		stencilDB, "diaspora_1000000_exp14", "mastodon_exp14"

	indepMigrationType := "i"

	indepEnableDisplay, indepDisplayInFirstPhase := false, false

	migrationNum := 100

	data := ReadStrLinesFromLog(SA1MigrationFile)

	evalConfig := InitializeEvalConfig()
	defer closeDBConns(evalConfig)

	for i := 0; i < migrationNum; i++ {
		
		data1 := data[i]

		var sizeData SA1SizeStruct

		err := json.Unmarshal([]byte(data1), &sizeData)
		if err != nil {
			log.Fatal(err)
		}

		userID := sizeData.UserID

		log.Println("UserID:", userID)

		migrateUserFromDiasporaToMastodon(
			evalConfig, indepStencilDB, diaspora, 
			userID, indepMigrationType, 
			indepStencilDB, indepSrcDB, indepDstDB,
			indepSizeFile, indepTimeFile,
			indepEnableDisplay, indepDisplayInFirstPhase,
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

	log.Println("=================================")
	log.Println("Starting Exp3: Data Downtime Test")
	log.Println("=================================")
	
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

func Exp3LoadUserIDsByDagCounter() {

	log.Println("=================================")
	log.Println("Starting Exp3: Data Downtime Test")
	log.Println("=================================")


	SA1StencilDB, SA1SrcDB, SA1DstDB := 
		"stencil_100k_exp8", "diaspora_100k_exp8", "mastodon_100k_exp8"

	SA1MigrationType := "d"

	SA1EnableDisplay, SA1DisplayInFirstPhase := true, true


	naiveStencilDB, naiveSrcDB, naiveDstDB := 
		"stencil_100k_exp9", "diaspora_100k_exp9", "mastodon_100k_exp9"

	naiveMigrationType := "n"

	naiveEnableDisplay, naiveDisplayInFirstPhase := true, false


	edgeCounterRangeStart := 750
	edgeCounterRangeEnd := 1200
	migrationNum := 100

	// Note that the data in the table dag_counter in stencilDB is latest
	stencilDB = SA1StencilDB


	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	edgeCounter := getEdgesCounterByRange(
		evalConfig,
		edgeCounterRangeStart, 
		edgeCounterRangeEnd, 
		migrationNum,
	)

	log.Println(edgeCounter)
	log.Println("Migration number:", len(edgeCounter))

	for _, edgeCounter1 := range edgeCounter {
		
		userID := edgeCounter1["person_id"]

		log.Println("UserID is", userID)

		migrateUserFromDiasporaToMastodon1(
			userID, SA1MigrationType, 
			SA1StencilDB, SA1SrcDB, SA1DstDB, 
			SA1EnableDisplay, SA1DisplayInFirstPhase,
		)

		migrateUserFromDiasporaToMastodon1(
			userID, naiveMigrationType, 
			naiveStencilDB, naiveSrcDB, naiveDstDB, 
			naiveEnableDisplay, naiveDisplayInFirstPhase,
		)

	}
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

func Exp3GetDataDowntimeByLoadingUserIDFromLog() {

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

func Exp3GetDataDowntimeInPercentageByLoadingUserIDFromLog() {

	SA1StencilDB := "stencil_exp"
	naiveStencilDB := "stencil_exp5"

	logFile := "SA1Size"
	migrationNum := 100

	stencilDB = SA1StencilDB
	stencilDB1 = naiveStencilDB

	evalConfig := InitializeEvalConfig()
	defer closeDBConns(evalConfig)

	data := ReadStrLinesFromLog(logFile)

	var SA1PercentageOfDowntime, naivePercentageOfDowntime []float64

	for i := 0; i < migrationNum; i++ {

		var data1 SA1SizeStruct 

		err := json.Unmarshal([]byte(data[i]), &data1)
		if err != nil {
			log.Fatal(err)
		}

		userID := data1.UserID

		SA1MigrationID := getMigrationIDBySrcUserID1(evalConfig.StencilDBConn, userID)
		naiveMigrationID := getMigrationIDBySrcUserID1(evalConfig.StencilDBConn1, userID)

		SA1Downtime := getDataDowntimeOfMigration(evalConfig.StencilDBConn,
			SA1MigrationID)
		naiveDowntime := getDataDowntimeOfMigration(evalConfig.StencilDBConn1,
			naiveMigrationID)

		SA1TotalTime := getTotalTimeOfMigration(evalConfig.StencilDBConn, 
			SA1MigrationID)
		naiveTotalTime := getTotalTimeOfMigration(evalConfig.StencilDBConn1, 
			naiveMigrationID)
		
		SA1PercentageOfDowntime1 := calculateTimeInPercentage(SA1Downtime, SA1TotalTime)
		naivePercentageOfDowntime1 := calculateTimeInPercentage(naiveDowntime, naiveTotalTime)
		
		SA1PercentageOfDowntime = append(SA1PercentageOfDowntime, 
			SA1PercentageOfDowntime1...)
		naivePercentageOfDowntime = append(naivePercentageOfDowntime, 
			naivePercentageOfDowntime1...)

	}

	log.Println(SA1PercentageOfDowntime)

	WriteStrArrToLog(
		evalConfig.DataDowntimeInPercentageInStencilFile, 
		ConvertFloat64ToString(SA1PercentageOfDowntime),
	)

	WriteStrArrToLog(
		evalConfig.DataDowntimeInPercentageInNaiveFile, 
		ConvertFloat64ToString(naivePercentageOfDowntime),
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
	for i := 577000; i < len(userIDs); i += 100 {  
		
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

	// stencilDB = "stencil_counter"
	// counterFile := "diaspora1MCounter"
	// counterTable := "dag_counter"

	stencilDB = "stencil_exp9_0"
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

// diaspora_100000: count edges and nodes every 10 users
// diaspora_10000 and diaspora_10000: count edges and nodes of every user
func Exp4CountEdgesNodes() {

	// db.DIASPORA_DB = "diaspora_100000"
	// diaspora = "diaspora_100000"	
	
	db.DIASPORA_DB = "diaspora_100000"
	diaspora = "diaspora_100000"

	// db.DIASPORA_DB = "diaspora_1000"
	// diaspora = "diaspora_1000"

	appName, appID := "diaspora", "1"

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	// The file name needs to be changed
	file := evalConfig.Diaspora100KCounterFile

	userIDs := getAllUserIDsInDiaspora(evalConfig, true)

	log.Println("total users:", len(userIDs))

	for i := 22975; i < len(userIDs); i += 10 {
	// for _, userID := range userIDs {

		userID := userIDs[i]

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

	log.Println("===============================")
	log.Println("Starting Exp6: Scalability Test")
	log.Println("===============================")

	// stencilDB = "stencil_exp3"
	// mastodon = "mastodon_exp3"
	// diaspora = "diaspora_1000000_exp3"

	// db.STENCIL_DB = "stencil_exp3"
	// db.DIASPORA_DB = "diaspora_1000000_exp3"
	// db.MASTODON_DB = "mastodon_exp3"

	stencilDB = "stencil_100k_exp4"
	mastodon = "mastodon_100k_exp4"

	stencilDB1 = "stencil_100k_exp5"
	mastodon = "mastodon_100k_exp5"

	// counterStart := 0
	// counterNum := 300
	// counterInterval := 10

	counterStart := 0
	counterNum := 100
	counterInterval := 10

	evalConfig := InitializeEvalConfig()

	defer closeDBConns(evalConfig)

	preExp2(evalConfig)

	res := getEdgesCounter(evalConfig, 
		counterStart, counterNum, counterInterval)

	log.Println("Total Num:", len(res))
	log.Println(res)

	for i := 0; i < len(res); i ++ {

		res1 := res[i]

		userID := res1["person_id"]

		db.STENCIL_DB = "stencil_100k_exp4"
		db.DIASPORA_DB = "diaspora_100k_exp4"
		db.MASTODON_DB = "mastodon_100k_exp4"

		uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
			userID, "diaspora", "1", "mastodon", "2", "d", 1

		enableDisplay, displayInFirstPhase, markAsDelete := true, false, true

		SA1_migrate.Controller(uid, srcAppName, srcAppID, 
			dstAppName, dstAppID, migrationType, threadNum,
			enableDisplay, displayInFirstPhase, markAsDelete,
		)

		db.STENCIL_DB = "stencil_100k_exp5"
		db.DIASPORA_DB = "diaspora_100k_exp5"
		db.MASTODON_DB = "mastodon_100k_exp5"
		
		migrationType = "i"

		SA1_migrate.Controller(uid, srcAppName, srcAppID, 
			dstAppName, dstAppID, migrationType, threadNum,
			enableDisplay, displayInFirstPhase, markAsDelete,
		)

		log.Println("************ Calculate Migration and Display Time ************")

		refreshEvalConfigDBConnections(evalConfig)

		mTime, dTime := getMigrationAndDisplayTimeBySrcUserID(evalConfig.StencilDBConn, userID)

		mTime1, dTime1 := getMigrationAndDisplayTimeBySrcUserID(evalConfig.StencilDBConn1, userID)

		res1["deletionMigrationTime"] = ConvertSingleDurationToString(mTime)
		res1["deletionDisplayTime"] = ConvertSingleDurationToString(dTime)

		res1["independentMigrationTime"] = ConvertSingleDurationToString(mTime1)
		res1["independentDisplayTime"] = ConvertSingleDurationToString(dTime1)

		log.Println("************ Calculate Nodes and Edges after Migration ************")

		migrationID := getMigrationIDBySrcUserID1(evalConfig.StencilDBConn, userID)

		migratedUserID := getMigratedUserID(evalConfig.StencilDBConn,
			 migrationID, dstAppID)

		// Note that we need to replace the set db.MASTODON_DB here
		db.STENCIL_DB = "stencil_100k_exp4"
		db.DIASPORA_DB = "diaspora_100k_exp4"
		db.MASTODON_DB = "mastodon_100k_exp4"

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

func Exp6AddIndepLoadFromLog() {

	log.Println("===========================================================")
	log.Println("Starting Exp6: Add Indpendent Migration to Scalability Test")
	log.Println("===========================================================")

	// stencilDB = "stencil_exp3"
	// mastodon = "mastodon_exp3"
	// diaspora = "diaspora_1000000_exp3"

	// db.STENCIL_DB = "stencil_exp3"
	// db.DIASPORA_DB = "diaspora_1000000_exp3"
	// db.MASTODON_DB = "mastodon_exp3"

	stencilDB = "stencil_exp15"
	mastodon = "mastodon_exp15"
	diaspora = "diaspora_1000000_exp15"

	scalabilityFile := "scalability"

	scalabilityWithIndependentFile := "scalabilityWithIndependent"

	data := ReadStrLinesFromLog(scalabilityFile)

	evalConfig := InitializeEvalConfig()
	defer closeDBConns(evalConfig)

	for _, data1 := range data {

		var sData ScalabilityDataStruct

		err := json.Unmarshal([]byte(data1), &sData)
		if err != nil {
			log.Fatal(err)
		}

		userID := sData.PersonID

		uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType, threadNum := 
			userID, "diaspora", "1", "mastodon", "2", "i", 1

		enableDisplay, displayInFirstPhase := false, false

		db.STENCIL_DB = stencilDB
		db.DIASPORA_DB = diaspora
		db.MASTODON_DB = mastodon

		SA1_migrate.Controller(uid, srcAppName, srcAppID, 
			dstAppName, dstAppID, migrationType, threadNum,
			enableDisplay, displayInFirstPhase,
		)

		log.Println("************ Calculate Migration Time ************")

		refreshEvalConfigDBConnections(evalConfig)

		mTime := getMigrationTimeBySrcUserID(evalConfig.StencilDBConn, userID)

		res1 := make(map[string]string) 

		res1["indepMigrationTime"] = ConvertSingleDurationToString(mTime)
		res1["displayTime"] = sData.DisplayTime
		res1["edges"] = sData.Edges
		res1["edgesAfterMigration"] = sData.EdgesAfterMigration
		res1["migrationTime"] = sData.MigrationTime
		res1["nodes"] = sData.Nodes
		res1["nodesAfterMigration"] = sData.NodesAfterMigration
		res1["person_id"] = sData.PersonID

		WriteStrToLog(
			scalabilityWithIndependentFile,
			ConvertMapStringToJSONString(res1),
		)
	}

}

func Exp7ReintegrationDataBags() {

	log.Println("===========================================")
	log.Println("Starting Exp7: Dangling Data Reintegration")
	log.Println("===========================================")

	migrationSeqs := [][]string {
		// []string{"diaspora", "mastodon", "gnusocial", "twitter", "diaspora"}, 
		// []string{"mastodon", "gnusocial", "twitter", "diaspora", "mastodon"},
		[]string{"twitter", "diaspora", "mastodon", "gnusocial", "twitter"},
		// []string{"gnusocial", "twitter", "diaspora", "mastodon", "gnusocial"},
	}

	migrationNum := 1
	
	size := "1k"

	for k, migrationSeq := range migrationSeqs {

		progress := fmt.Sprintf("%d/%d", k + 1, len(migrationSeqs))
		log.Println("***************************")
		log.Println(progress, "- Current migration sequence:", migrationSeq)
		log.Println("Migration number:", migrationNum)
		log.Println("***************************")

		logFile, logFile1 := setDatabasesLogFilesForExp7(migrationSeq[0], size)

		evalConfig := InitializeEvalConfig(false)
		defer closeDBConns(evalConfig)

		preExp7(evalConfig, migrationSeq[0])
		
		userIDs := evalConfig.getRootUserIDsRandomly(migrationSeq[0], migrationNum)

		// // edgeCounterRangeStart := 300
		// edgeCounterRangeStart := 500
		// edgeCounterRangeEnd := 1200
		// getCounterNum := 100

		// edgeCounter := getEdgesCounterByRange(
		// 	evalConfig,
		// 	edgeCounterRangeStart, 
		// 	edgeCounterRangeEnd, 
		// 	getCounterNum,
		// )

		// log.Println(edgeCounter)

		for j := 0; j < migrationNum; j++ {

			// userID := edgeCounter[j]["person_id"]
			userID := userIDs[j]
			// userID = "7088175756532884328"
			userID1 := userID

			log.Println("Next User:", userID)

			preExp7(evalConfig, migrationSeq[0])

			var beforeMigObjsInApp int64
			var beforeMigObjsInApp1 int64

			var afterMigObjsInApp int64
			var afterMigObjsInApp1 int64

			var migObjsInSrc int64
			var migObjsInSrc1 int64

			var beforeLastMigObjsInLastApp int64
			var beforeLastMigObjsInLastApp1 int64

			var afterLastMigObjsInLastApp int64
			var afterLastMigObjsInLastApp1 int64

			var migrationID string
			var migrationID1 string

			for i := 0; i < len(migrationSeq) - 1; i++ {

				fromApp := migrationSeq[i]
				toApp := migrationSeq[i+1]
				
				fromAppID := db.GetAppIDByAppName(evalConfig.StencilDBConn, fromApp)
				toAppID := db.GetAppIDByAppName(evalConfig.StencilDBConn, toApp)

				beforeMigObjsInApp = getTotalObjsNotIncludingMediaOfAppInExp7V2(
					evalConfig, fromApp, true)
					
				beforeMigObjsInApp1 = getTotalObjsNotIncludingMediaOfAppInExp7V2(
					evalConfig, fromApp, false)

				if i == len(migrationSeq) - 2 {
					beforeLastMigObjsInLastApp = getTotalObjsNotIncludingMediaOfAppInExp7V2(
						evalConfig, toApp, true)
					beforeLastMigObjsInLastApp1 = getTotalObjsNotIncludingMediaOfAppInExp7V2(
						evalConfig, toApp, false)
				}

				enableDisplay := true
				// if i == len(migrationSeq) - 2 {
				// 	enableDisplay = false
				// }

				enableBags := true

				db.STENCIL_DB = stencilDB
				db.DIASPORA_DB = diaspora
				db.MASTODON_DB = mastodon
				db.TWITTER_DB = twitter
				db.GNUSOCIAL_DB = gnusocial

				migrationID, userID = migrateUsersInSeqOfApps(
					evalConfig, stencilDB,
					i, fromApp, toApp, fromAppID, toAppID,
					migrationID, userID, 
					enableBags, enableDisplay,
				)

				enableBags = false

				db.STENCIL_DB = stencilDB1
				db.DIASPORA_DB = diaspora1
				db.MASTODON_DB = mastodon1
				db.TWITTER_DB = twitter1
				db.GNUSOCIAL_DB = gnusocial1
				
				migrationID1, userID1 = migrateUsersInSeqOfApps(
					evalConfig, stencilDB1,
					i, fromApp, toApp, fromAppID, toAppID,
					migrationID1, userID1,
					enableBags, enableDisplay,
				)

				afterMigObjsInApp = getTotalObjsNotIncludingMediaOfAppInExp7V2(
					evalConfig, fromApp, true)
				
				afterMigObjsInApp1 = getTotalObjsNotIncludingMediaOfAppInExp7V2(
					evalConfig, fromApp, false)
				
				migObjsInSrc = beforeMigObjsInApp - afterMigObjsInApp
				migObjsInSrc1 = beforeMigObjsInApp1 - afterMigObjsInApp1

				// Exclude others' dangling data
				if i == 0 {
					migObjsInSrc -= evalConfig.getOthersDanglingData(
						stencilDB, userID, fromAppID, migrationID)
					migObjsInSrc1 -= evalConfig.getOthersDanglingData(
						stencilDB1, userID1, fromAppID, migrationID1)
				} 

				// Only calculate how much data has been migrated out, so even though some data
				// is left in the app, such as conversations, it will not be counted
				logSeqMigsRes(logFile, userID, migObjsInSrc)
				logSeqMigsRes(logFile1, userID1, migObjsInSrc1)
				
				if i == len(migrationSeq) - 2 {
					afterLastMigObjsInLastApp = getTotalObjsNotIncludingMediaOfAppInExp7V2(
						evalConfig, toApp, true)
					afterLastMigObjsInLastApp1 = getTotalObjsNotIncludingMediaOfAppInExp7V2(
						evalConfig, toApp, false)

					lastUserID := evalConfig.getNextUserID(evalConfig.StencilDBConn, migrationID)
					lastUserID1 := evalConfig.getNextUserID(evalConfig.StencilDBConn1, migrationID1)

					logSeqMigsRes(
						logFile, lastUserID, 
						afterLastMigObjsInLastApp - beforeLastMigObjsInLastApp,
					)
					logSeqMigsRes(
						logFile1, lastUserID1, 
						afterLastMigObjsInLastApp1 - beforeLastMigObjsInLastApp1,
					)
				}
			}
		}
	}
}

func Exp8SA1() {

	log.Println("===================================")
	log.Println("Starting Exp8 for SA1: Dataset Test")
	log.Println("===================================")

	migrationNum := 100

	// exp8LogFile := "diaspora_1K_dataset"

	// stencilDB = "stencil_exp13"
	// diaspora = "diaspora_1k_exp13"
	// mastodon = "mastodon_exp13"

	// diaspora1 = "diaspora_1000"

	
	exp8LogFile := "diaspora_10K_dataset"

	stencilDB = "stencil_exp12"
	diaspora = "diaspora_10k_exp12"
	mastodon = "mastodon_exp12"

	diaspora1 = "diaspora_10000"


	// exp8LogFile := "diaspora_100K_dataset"

	// stencilDB = "stencil_exp11"
	// diaspora = "diaspora_100k_exp11"
	// mastodon = "mastodon_exp11"

	// diaspora1 = "diaspora_100000"


	// exp8LogFile := "diaspora_1M_dataset"

	// stencilDB = "stencil_exp10"
	// diaspora = "diaspora_1m_exp10"
	// mastodon = "mastodon_exp10"

	// diaspora1 = "diaspora_1000000"

	// allowedDBName := map[string]string {
	// 	"diaspora_1k_exp13", "diaspora_10k_exp12",
	// 	"diaspora_100k_exp11", "diaspora_1m_exp10",
	// }

	evalConfig := InitializeEvalConfig()
	defer closeDBConns(evalConfig)

	preExp1(evalConfig)

	userIDsWithCounters := getUserIDsWithSameNodesAcrossDatasets(
		evalConfig.StencilDBConn, diaspora)

	log.Println(userIDsWithCounters)

	for i := 0; i < migrationNum; i++ {

		data1 := userIDsWithCounters[i]
		
		userID := data1["person_id"]
		nodes := data1["nodes"]
		edges := data1["edges"]

		SA1StencilDB, SA1SrcDB, SA1DstDB := stencilDB, diaspora, mastodon

		SA1MigrationType := "d"

		SA1EnableDisplay, SA1DisplayInFirstPhase := true, true

		migrateUserFromDiasporaToMastodon1(
			userID, SA1MigrationType, 
			SA1StencilDB, SA1SrcDB, SA1DstDB, 
			SA1EnableDisplay, SA1DisplayInFirstPhase,
		)

		log.Println("************ Calculate the Migration Size and Time ************")

		logData := make(map[string]string)

		refreshEvalConfigDBConnections(evalConfig)

		migrationID := getMigrationIDBySrcUserIDMigrationType(evalConfig.StencilDBConn, 
			userID, SA1MigrationType)

		migrationIDInt, err := strconv.Atoi(migrationID)
		if err != nil {
			log.Fatal(err)
		}

		time := GetMigrationTime(
			evalConfig.StencilDBConn,
			migrationIDInt,
		)

		size := GetMigratedDataSizeBySrc(
			evalConfig,
			evalConfig.StencilDBConn,
			evalConfig.DiasporaDBConn1,
			migrationID,
		)

		logData["size"] = ConvertInt64ToString(size)
		logData["time"] = ConvertSingleDurationToString(time)	
		logData["userID"] = userID
		logData["nodes"] = nodes
		logData["edges"] = edges

		WriteStrToLog(
			exp8LogFile,
			ConvertMapStringToJSONString(logData),
		)

	}
}

func Exp8SA2() {

	log.Println("===================================")
	log.Println("Starting Exp8 for SA2: Dataset Test")
	log.Println("===================================")

	seq := 1

	migrationNum := 5

	exp8LogFile := "diaspora_100K_dataset_sa2_" + strconv.Itoa(seq)

	db.STENCIL_DB = "stencil_sa2_100k_exp" + strconv.Itoa(seq)

	var stencilDBConn *sql.DB

	stencilDBConn = db.GetDBConn(db.STENCIL_DB)

	userIDsWithCounters := getUserIDsWithSameNodesAcrossDatasets(stencilDBConn, db.STENCIL_DB)

	log.Println("SEQUENCE NUM:", seq)

	log.Println(userIDsWithCounters)

	for i := migrationNum * seq; i < migrationNum * seq + migrationNum; i++ {

		userID := userIDsWithCounters[i]["person_id"]

		srcApp, srcAppID, dstApp, dstAppID, migrationType, enableBags :=
			"diaspora", "1", "mastodon", "2", "d", true

		apis.StartMigrationSA2(
			userID, srcApp, srcAppID, dstApp, dstAppID, migrationType, enableBags,
		)

		log.Println("************ Calculate the Migration Time ************")

		stencilDBConn.Close()
		stencilDBConn = db.GetDBConn(db.STENCIL_DB)

		logData := make(map[string]string)

		migrationID := getMigrationIDBySrcUserID1(stencilDBConn, userID)

		migrationIDInt, err := strconv.Atoi(migrationID)
		if err != nil {
			log.Fatal(err)
		}

		time := GetMigrationTime(stencilDBConn, migrationIDInt)

		logData["time"] = ConvertSingleDurationToString(time)	
		logData["userID"] = userID

		WriteStrToLog(
			exp8LogFile,
			ConvertMapStringToJSONString(logData),
		)
	}
}

func Exp9PreserveDataLinks() {

	log.Println("============================================")
	log.Println("Starting Exp9: Preserve Data Relationships")
	log.Println("============================================")

	migrationSeq := []string {
		"diaspora", "mastodon", "gnusocial", "twitter", 
		"diaspora", "mastodon", "gnusocial", "twitter",
	}
	log.Println("Migration sequence:", migrationSeq)

	seq := 0
	seqStr := strconv.Itoa(seq)
	// log.Println("Sequence:", seq)

	migrationNum := 100
	log.Println("Migration number:", migrationNum)
	
	// Database setup for migrations enabled databags
	stencilDB = "stencil_exp6_" + seqStr
	diaspora = "diaspora_1m_exp6_" + seqStr
	mastodon = "mastodon_exp6_" + seqStr
	twitter = "twitter_exp6_" + seqStr
	gnusocial = "gnusocial_exp6_" + seqStr
	logFile := "resolveRefsEnabled_" + seqStr

	// Database setup for migrations not enabled databags
	stencilDB1 = "stencil_exp7_" + seqStr
	diaspora1 = "diaspora_1m_exp7_" + seqStr
	mastodon1 = "mastodon_exp7_" + seqStr
	twitter1 = "twitter_exp7_" + seqStr
	gnusocial1 = "gnusocial_exp7_" + seqStr
	logFile1 := "resolveRefsNotEnabled_" + seqStr

	// edgeCounterRangeStart := 300
	edgeCounterRangeStart := 700
	edgeCounterRangeEnd := 1200
	getCounterNum := 100

	evalConfig := InitializeEvalConfig(false)
	defer closeDBConns(evalConfig)

	preExp9(evalConfig)
	
	edgeCounter := getEdgesCounterByRange(
		evalConfig,
		edgeCounterRangeStart, 
		edgeCounterRangeEnd, 
		getCounterNum,
	)

	log.Println(edgeCounter)

	for j := seq * migrationNum; j < (seq + 1) * migrationNum; j++ {

		userID := edgeCounter[j]["person_id"]
		userID1 := userID

		log.Println("Next User:", userID)

		preExp9(evalConfig)

		var beforeMigObjsInApp int64
		var beforeMigObjsInApp1 int64

		var afterMigObjsInApp int64
		var afterMigObjsInApp1 int64

		var migObjsInSrc int64
		var migObjsInSrc1 int64

		var beforeLastMigObjsInLastApp int64
		var beforeLastMigObjsInLastApp1 int64

		var afterLastMigObjsInLastApp int64
		var afterLastMigObjsInLastApp1 int64

		var migrationID string
		var migrationID1 string

		for i := 0; i < len(migrationSeq) - 1; i++ {

			fromApp := migrationSeq[i]
			toApp := migrationSeq[i+1]
			
			fromAppID := db.GetAppIDByAppName(evalConfig.StencilDBConn, fromApp)
			toAppID := db.GetAppIDByAppName(evalConfig.StencilDBConn, toApp)

			beforeMigObjsInApp = getTotalObjsNotIncludingMediaOfAppInExp7V2(
				evalConfig, fromApp, true)
				
			beforeMigObjsInApp1 = getTotalObjsNotIncludingMediaOfAppInExp7V2(
				evalConfig, fromApp, false)

			if i == len(migrationSeq) - 2 {
				beforeLastMigObjsInLastApp = getTotalObjsNotIncludingMediaOfAppInExp7V2(
					evalConfig, toApp, true)
				beforeLastMigObjsInLastApp1 = getTotalObjsNotIncludingMediaOfAppInExp7V2(
					evalConfig, toApp, false)
			}

			enableBags := true
			resolveRefs := true

			db.STENCIL_DB = stencilDB
			db.DIASPORA_DB = diaspora
			db.MASTODON_DB = mastodon
			db.TWITTER_DB = twitter
			db.GNUSOCIAL_DB = gnusocial

			migrationID, userID = migrateUsersInSeqOfApps(
				evalConfig, stencilDB,
				i, fromApp, toApp, fromAppID, toAppID,
				migrationID, userID, 
				enableBags, resolveRefs,
			)

			resolveRefs = false

			db.STENCIL_DB = stencilDB1
			db.DIASPORA_DB = diaspora1
			db.MASTODON_DB = mastodon1
			db.TWITTER_DB = twitter1
			db.GNUSOCIAL_DB = gnusocial1
			
			migrationID1, userID1 = migrateUsersInSeqOfApps(
				evalConfig, stencilDB1,
				i, fromApp, toApp, fromAppID, toAppID,
				migrationID1, userID1,
				enableBags, resolveRefs,
			)

			afterMigObjsInApp = getTotalObjsNotIncludingMediaOfAppInExp7V2(
				evalConfig, fromApp, true)
			
			afterMigObjsInApp1 = getTotalObjsNotIncludingMediaOfAppInExp7V2(
				evalConfig, fromApp, false)
			
			migObjsInSrc = beforeMigObjsInApp - afterMigObjsInApp
			migObjsInSrc1 = beforeMigObjsInApp1 - afterMigObjsInApp1

			// Exclude others' dangling data
			if i == 0 {
				migObjsInSrc -= evalConfig.getOthersDanglingData(
					stencilDB, userID, fromAppID, migrationID)
				migObjsInSrc1 -= evalConfig.getOthersDanglingData(
					stencilDB1, userID1, fromAppID, migrationID1)
			} 

			// Only calculate how much data has been migrated out, so even though some data
			// is left in the app, such as conversations, it will not be counted
			logSeqMigsRes(logFile, userID, migObjsInSrc)
			logSeqMigsRes(logFile1, userID1, migObjsInSrc1)
			
			if i == len(migrationSeq) - 2 {
				afterLastMigObjsInLastApp = getTotalObjsNotIncludingMediaOfAppInExp7V2(
					evalConfig, toApp, true)
				afterLastMigObjsInLastApp1 = getTotalObjsNotIncludingMediaOfAppInExp7V2(
					evalConfig, toApp, false)

				lastUserID := evalConfig.getNextUserID(evalConfig.StencilDBConn, migrationID)
				lastUserID1 := evalConfig.getNextUserID(evalConfig.StencilDBConn1, migrationID1)

				logSeqMigsRes(
					logFile, lastUserID, 
					afterLastMigObjsInLastApp - beforeLastMigObjsInLastApp,
				)
				logSeqMigsRes(
					logFile1, lastUserID1, 
					afterLastMigObjsInLastApp1 - beforeLastMigObjsInLastApp1,
				)
			}
		}
	}
}
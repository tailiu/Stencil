/*
 * Migration Handler
 */

package main

import (
	"fmt"
	"log"
	"os"
	"stencil/config"
	"stencil/migrate"
	"stencil/transaction"
	"sync"
)

/*********************--bgn
 * Functions
***************************/

// func traverseDependencies(tableName string, dependencies []config.Dependency, rDependencies *[]config.Dependency) {

// 	for _, dependency := range dependencies {
// 		if strings.ToLower(tableName) == strings.ToLower(dependency.DependsOn) {
// 			if strings.ToLower(dependency.Tag) != strings.ToLower(dependency.DependsOn) {
// 				traverseDependencies(dependency.Tag, dependencies, rDependencies)
// 			}
// 			*rDependencies = append(*rDependencies, dependency)
// 		}
// 	}
// }

func prepareDataQueries(appconfig config.AppConfig) []config.DataQuery {

	var sqls []config.DataQuery

	// sql := fmt.Sprintf("SELECT %s.* FROM %s WHERE %s = $1 ", settings.UserTable, settings.UserTable, settings.KeyCol, settings.UserTable)
	// sqls = append(sqls, config.DataQuery{SQL: sql, Table: settings.UserTable})

	// for _, dependency := range dependencies {
	// 	if strings.ToLower(settings.UserTable) == strings.ToLower(dependency.DependsOn) {

	// 		var subDeps []config.Dependency
	// 		subDeps = append(subDeps, dependency)
	// 		traverseDependencies(dependency.Tag, dependencies, &subDeps)

	// 		for _, subDep := range subDeps {
	// 			sql = fmt.Sprintf("SELECT %s.* FROM %s ", subDep.Tag, subDep.Tag)
	// 			dep := subDep
	// 			for true {
	// 				sql += fmt.Sprintf(" JOIN %s ON ", dep.DependsOn)
	// 				for i, condition := range dep.Conditions {
	// 					sql += fmt.Sprintf("%s.%s = %s.%s", dep.Tag, condition.TagAttr, dep.DependsOn, condition.DependsOnAttr)
	// 					if i < len(dep.Conditions)-1 {
	// 						sql += " AND "
	// 					}
	// 				}
	// 				if strings.EqualFold(dep.DependsOn, settings.UserTable) {
	// 					sql += fmt.Sprintf(" WHERE %s.%s = $1 ", settings.UserTable, settings.KeyCol, subDep.Tag)
	// 					break
	// 				} else {
	// 					var e error
	// 					dep, e = config.FindDependencyByDependsOn(dep.DependsOn, subDeps)
	// 					if e == nil && !strings.EqualFold(dep.DependsOn, settings.UserTable) {
	// 						// fmt.Println("!!! Couldn't find dependency tag that depends on ", dep.DependsOn)
	// 						break
	// 					}
	// 				}
	// 			}
	// 			sqls = append(sqls, config.DataQuery{SQL: sql, Table: subDep.Tag})
	// 		}
	// 	}
	// }
	return sqls
}

// func prepareAppLevelData(settings config.Settings, dependencies []config.Dependency) []config.DataQuery {

// 	// todo: handle self-dependencies ** properly **. some cases are failing, rn. :S

// 	var sqls []config.DataQuery

// 	sql := fmt.Sprintf("SELECT %s.* FROM %s WHERE %s = $1 AND %s.mark_delete != 'true' ", settings.UserTable, settings.UserTable, settings.KeyCol, settings.UserTable)
// 	sqls = append(sqls, config.DataQuery{SQL: sql, Table: settings.UserTable})

// 	for _, dependency := range dependencies {
// 		if strings.ToLower(settings.UserTable) == strings.ToLower(dependency.DependsOn) {

// 			var subDeps []config.Dependency
// 			subDeps = append(subDeps, dependency)
// 			traverseDependencies(dependency.Tag, dependencies, &subDeps)

// 			for _, subDep := range subDeps {
// 				sql = fmt.Sprintf("SELECT %s.* FROM %s ", subDep.Tag, subDep.Tag)
// 				dep := subDep
// 				for true {
// 					sql += fmt.Sprintf(" JOIN %s ON ", dep.DependsOn)
// 					for i, condition := range dep.Conditions {
// 						sql += fmt.Sprintf("%s.%s = %s.%s", dep.Tag, condition.TagAttr, dep.DependsOn, condition.DependsOnAttr)
// 						if i < len(dep.Conditions)-1 {
// 							sql += " AND "
// 						}
// 					}
// 					if strings.EqualFold(dep.DependsOn, settings.UserTable) {
// 						sql += fmt.Sprintf(" WHERE %s.%s = $1 AND %s.mark_delete != 'true'", settings.UserTable, settings.KeyCol, subDep.Tag)
// 						break
// 					} else {
// 						var e error
// 						dep, e = config.FindDependencyByDependsOn(dep.DependsOn, subDeps)
// 						if e == nil && !strings.EqualFold(dep.DependsOn, settings.UserTable) {
// 							// fmt.Println("!!! Couldn't find dependency tag that depends on ", dep.DependsOn)
// 							break
// 						}
// 					}
// 				}
// 				sqls = append(sqls, config.DataQuery{SQL: sql, Table: subDep.Tag})
// 			}
// 		}
// 	}
// 	return sqls
// }

// func preparePhysicalData(settings config.Settings, dependencies []config.Dependency) []config.DataQuery {

// 	// todo: handle self-dependencies ** properly **. some cases are failing, rn. :S

// 	var sqls []config.DataQuery

// 	sql := fmt.Sprintf("SELECT %s.* FROM %s WHERE %s = $1 ", settings.UserTable, settings.UserTable, settings.KeyCol)
// 	sqls = append(sqls, config.DataQuery{SQL: sql, Table: settings.UserTable})

// 	for _, dependency := range dependencies {
// 		if strings.ToLower(settings.UserTable) == strings.ToLower(dependency.DependsOn) {

// 			var subDeps []config.Dependency
// 			subDeps = append(subDeps, dependency)
// 			traverseDependencies(dependency.Tag, dependencies, &subDeps)

// 			for _, subDep := range subDeps {
// 				sql = fmt.Sprintf("SELECT %s.* FROM %s ", subDep.Tag, subDep.Tag)
// 				dep := subDep
// 				for true {
// 					sql += fmt.Sprintf(" JOIN %s ON ", dep.DependsOn)
// 					for i, condition := range dep.Conditions {
// 						sql += fmt.Sprintf("%s.%s = %s.%s", dep.Tag, condition.TagAttr, dep.DependsOn, condition.DependsOnAttr)
// 						if i < len(dep.Conditions)-1 {
// 							sql += " AND "
// 						}
// 					}
// 					if strings.EqualFold(dep.DependsOn, settings.UserTable) {
// 						sql += fmt.Sprintf(" WHERE %s.%s = $1 ", settings.UserTable, settings.KeyCol)
// 						break
// 					} else {
// 						var e error
// 						dep, e = config.FindDependencyByDependsOn(dep.DependsOn, subDeps)
// 						if e == nil && !strings.EqualFold(dep.DependsOn, settings.UserTable) {
// 							// fmt.Println("!!! Couldn't find dependency tag that depends on ", dep.DependsOn)
// 							break
// 						}
// 					}
// 				}
// 				sqls = append(sqls, config.DataQuery{SQL: sql, Table: subDep.Tag})
// 			}
// 		}
// 	}
// 	return sqls
// }

// func initAppLevelMigration(uid int, srcApp, tgApp string) {

// 	log.Printf("Init Migration for Customer '%d'. '%s' => '%s'\n", uid, srcApp, tgApp)
// 	helper.Linebreak("\n")

// 	dependencies, err := config.ReadDependencies(srcApp)
// 	if err != nil {
// 		log.Fatal("error reading dependencies for:"+srcApp, err)
// 	}

// 	settings, err := config.ReadAppSettings(srcApp, false)
// 	if err != nil {
// 		log.Fatal("error reading settings for:"+srcApp, err)
// 	}

// 	sqls := prepareAppLevelData(settings, dependencies)

// 	helper.Linebreak("=", 80)
// 	for _, sql := range sqls {
// 		helper.Linebreak("±", 50)
// 		fmt.Println("#sql => ", sql.Table, ":", sql.SQL)
// 		helper.Linebreak("±", 50)
// 		migrate.MoveData(srcApp, tgApp, sql, settings.Mappings, uid)
// 	}
// 	helper.Linebreak("=", 80)

// }

// func initStencilMigration(uid int, srcApp, tgApp string) {

// 	log.Printf("Init Stencil Migration for Customer '%d'. '%s' => '%s'\n", uid, srcApp, tgApp)
// 	helper.Linebreak("\n")

// 	log_txn := transaction.BeginTransaction()

// 	dependencies, err := config.ReadDependencies(srcApp)
// 	if err != nil {
// 		log.Fatal("error reading dependencies for:"+srcApp, err)
// 	}
// 	log.Println("Dependencies fetched!")

// 	settings, err := config.ReadAppSettings(srcApp, false)
// 	if err != nil {
// 		log.Fatal("error reading settings for:"+srcApp, err)
// 	}
// 	log.Println("App Settings fetched!")

// 	sqls := preparePhysicalData(settings, dependencies)
// 	log.Println("SQLs prepared!")

// 	helper.Reverse(sqls)
// 	helper.Linebreak("=", 80)

// 	for _, sql := range sqls {
// 		helper.Linebreak("±", 50)
// 		fmt.Println("#sql => ", sql.Table, ":", sql.SQL)
// 		helper.Linebreak("±", 50)
// 		migrate.MigrateData(srcApp, tgApp, sql, settings.Mappings, uid, log_txn)
// 	}
// 	helper.Linebreak("=", 80)

// 	transaction.LogOutcome(log_txn, "COMMIT")

// }

/*********************--end
 * Functions
***************************/

/**************************
 * Main
***************************/

type ThreadChannel struct {
	Finished  bool
	Thread_id int
}

func main() {

	var wg sync.WaitGroup

	srcApp := "diaspora"
	dstApp := "mastodon"
	threads_num := 50
	// uid := 4716
	uid := os.Args[1]
	commitChannel := make(chan ThreadChannel)
	// startFrom, inc := 4670, 10
	// for uid := startFrom; uid < startFrom+inc; uid += 1 {
	config.LoadSchemaMappings()
	logTxn, err := transaction.BeginTransaction()
	if err == nil {
		for thread_id := 1; thread_id <= threads_num; thread_id++ {
			wg.Add(1)
			go func(thread_id int, commitChannel chan ThreadChannel) {
				defer wg.Done()
				if srcAppConfig, err := config.CreateAppConfig(srcApp); err != nil {
					commitChannel <- ThreadChannel{Finished: false, Thread_id: thread_id}
					log.Fatal(err)
				} else {
					if dstAppConfig, err := config.CreateAppConfig(dstApp); err != nil {
						commitChannel <- ThreadChannel{Finished: false, Thread_id: thread_id}
						log.Fatal(err)
					} else {
						migrate.ResetUserExistsInApp()
						if rootNode := migrate.GetRoot(srcAppConfig, fmt.Sprint(uid)); rootNode != nil {
							var wList = new(migrate.WaitingList)
							var invalidList = new(migrate.InvalidList)

							migrate.MigrateProcess(fmt.Sprint(uid), srcAppConfig, dstAppConfig, rootNode, wList, invalidList, logTxn)

						} else {
							fmt.Println("Root Node can't be fetched!")
						}
						dstAppConfig.CloseDBConn()
					}
					srcAppConfig.CloseDBConn()
					commitChannel <- ThreadChannel{Finished: true, Thread_id: thread_id}
				}
			}(thread_id, commitChannel)
		}
		go func() {
			wg.Wait()
			close(commitChannel)
		}()
	} else {
		log.Println("Can't begin migration transaction", err)
		transaction.LogOutcome(logTxn, "ABORT")
		// transaction.CloseDBConn(logTxn)
	}

	txnCommit := true

	for threadResponse := range commitChannel {
		fmt.Println("THREAD FINISHED WORKING", threadResponse)
		if !threadResponse.Finished {
			txnCommit = false
		}
	}

	if txnCommit {
		transaction.LogOutcome(logTxn, "COMMIT")
	} else {
		transaction.LogOutcome(logTxn, "ABORT")
	}

	// }

	// settingsFileName := "mappings"
	// // fromApp := "mastodon"
	// // toApp := "diaspora"
	// if schemaMappings, err := config.ReadSchemaMappingSettings(settingsFileName); err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println(schemaMappings)
	// }

	// initAppLevelMigration(7, "app1", "app5")
	// initStencilMigration(61, "app3", "app4")
	// QR := qr.NewQR("app1")
	// QR.TestQuery()

	// migrate.RollbackMigration(1503622861)
}
/*
 * Migration Handler
 */

package main

import (
	"fmt"
	"log"
	"strings"
	"transaction/config"
	"transaction/db"
	"transaction/helper"
)

/*********************--bgn
 * Functions
***************************/

func traverseDependencies(tableName string, dependencies []config.Dependency, rDependencies *[]config.Dependency) {

	for _, dependency := range dependencies {
		if strings.ToLower(tableName) == strings.ToLower(dependency.DependsOn) {
			if strings.ToLower(dependency.Tag) != strings.ToLower(dependency.DependsOn) {
				traverseDependencies(dependency.Tag, dependencies, rDependencies)
			}
			*rDependencies = append(*rDependencies, dependency)
		}
	}
}

func prepareData(settings config.Settings, dependencies []config.Dependency) []config.DataQuery {

	// todo: handle self-dependencies ** properly **. some cases are failing, rn. :S

	var sqls []config.DataQuery

	sql := fmt.Sprintf("SELECT %s.* FROM %s WHERE %s = $1", settings.UserTable, settings.UserTable, settings.KeyCol)
	sqls = append(sqls, config.DataQuery{SQL: sql, Table: settings.UserTable})

	for _, dependency := range dependencies {
		if strings.ToLower(settings.UserTable) == strings.ToLower(dependency.DependsOn) {

			var subDeps []config.Dependency
			subDeps = append(subDeps, dependency)
			traverseDependencies(dependency.Tag, dependencies, &subDeps)

			for _, subDep := range subDeps {
				sql = fmt.Sprintf("SELECT %s.* FROM %s ", subDep.Tag, subDep.Tag)
				dep := subDep
				for true {
					sql += fmt.Sprintf(" JOIN %s ON ", dep.DependsOn)
					for i, condition := range dep.Conditions {
						sql += fmt.Sprintf("%s.%s = %s.%s", dep.Tag, condition.TagAttr, dep.DependsOn, condition.DependsOnAttr)
						if i < len(dep.Conditions)-1 {
							sql += " AND "
						}
					}
					if strings.EqualFold(dep.DependsOn, settings.UserTable) {
						sql += fmt.Sprintf(" WHERE %s.%s = $1", settings.UserTable, settings.KeyCol)
						break
					} else {
						var e error
						dep, e = config.FindDependencyByDependsOn(dep.DependsOn, subDeps)
						if e == nil && !strings.EqualFold(dep.DependsOn, settings.UserTable) {
							// fmt.Println("!!! Couldn't find dependency tag that depends on ", dep.DependsOn)
							break
						}
					}
				}
				sqls = append(sqls, config.DataQuery{SQL: sql, Table: subDep.Tag})
			}
		}
	}
	return sqls
}

func initMigration(uid int, srcApp, tgApp string) {

	log.Printf("Init Migration for Customer '%d'. '%s' => '%s'\n", uid, srcApp, tgApp)
	helper.Linebreak("\n")

	dependencies, err := config.ReadDependencies(srcApp)
	if err != nil {
		log.Fatal("error reading dependencies for:"+srcApp, err)
	}

	settings, err := config.ReadAppSettings(srcApp)
	if err != nil {
		log.Fatal("error reading settings for:"+srcApp, err)
	}

	sqls := prepareData(settings, dependencies)

	helper.Linebreak("=", 80)
	for _, sql := range sqls {
		// fmt.Println("#sql => ", sql.Table, ":", sql.SQL)
		db.MoveData(srcApp, tgApp, sql, settings.Mappings, uid)
	}
	helper.Linebreak("=", 80)

	// usql := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", settings.UserTable, settings.KeyCol)
	// udata := DataCall(srcApp, usql, uid)[0]

	// fmt.Println("Fetched user data for uid: ", udata["c_id"])

}

/*********************--end
 * Functions
***************************/

/**************************
 * Main
***************************/

func main() {

	initMigration(1, "tpcc", "tpcc2")

}

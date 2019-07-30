/*
 * Migration Handler
 */

package main

import (
	"fmt"
	"log"
	"os"
	"stencil/db"
	"stencil/migrate"
	"stencil/mthread"
	"stencil/transaction"
	"strconv"
)

func main() {

	threads, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var mtype string

	switch os.Args[2] {
	case "d":
		{
			mtype = migrate.DELETION
		}
	case "i":
		{
			mtype = migrate.INDEPENDENT
		}
	case "c":
		{
			mtype = migrate.CONSISTENT
		}
	}

	if len(mtype) <= 0 {
		log.Fatal("can't read migration type")
	}

	limit, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	count := 0

	users, err := db.GetUnmigratedUsers()
	if err != nil {
		log.Fatal("UNMIGRATED USERS NOT FETCHED!", err)
	}

	for _, user := range users {
		uid := fmt.Sprint(user["user_id"])

		if logTxn, err := transaction.BeginTransaction(); err == nil {
			srcApp, srcAppID := "diaspora", "1"
			dstApp, dstAppID := "twitter", "3"

			// uid := os.Args[3]

			mWorker := migrate.CreateMigrationWorker(uid, srcApp, srcAppID, dstApp, dstAppID, logTxn, mtype)

			if mthread.ThreadController(mWorker, threads) {
				transaction.LogOutcome(logTxn, "COMMIT")
			} else {
				transaction.LogOutcome(logTxn, "ABORT")
			}
		} else {
			log.Fatal("Can't begin migration transaction", err)
		}
		count++
		if count >= limit {
			break
		}
	}
	log.Println("Users migrated:", count)
	// migrate.RollbackMigration(1503622861)
}

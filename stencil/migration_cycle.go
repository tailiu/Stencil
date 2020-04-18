package main

import (
	"flag"
	"fmt"
	"log"
	"stencil/SA1_display"
	"stencil/apis"
	"stencil/db"
	"sync"
)

func getUID(appID string) string {
	var query string
	switch appID {
	case "2":
		{
			query = "SELECT to_id as id FROM attribute_changes it where to_app = 2 and to_member = 56 "
		}
	case "4":
		{
			query = "SELECT to_id as id FROM attribute_changes it where to_app = 4 and to_member = 164 "
		}
	case "3":
		{
			query = "SELECT to_id as id FROM attribute_changes it where to_app = 3 and to_member = 119 "
		}
	default:
		{
			log.Fatal("Unknown app ", appID)
		}
	}
	dbConn := db.GetDBConn("stencil", false)
	defer dbConn.Close()

	if result, err := db.DataCall1(dbConn, query); err == nil {
		return fmt.Sprint(result["id"])
	} else {
		fmt.Println("query | ", query)
		log.Fatal(err)
	}
	log.Fatal("Why here?!")
	return ""
}

func main() {

	var wg sync.WaitGroup

	threads := flag.Int("threads", 1, "")
	mtype := flag.String("mtype", "d", "")
	uidInput := flag.String("uid", "54123", "")

	blade := flag.Bool("blade", false, "")
	bags := flag.Bool("bags", false, "")
	ftp := flag.Bool("ftp", false, "")
	display := flag.Bool("display", false, "")
	displayInFirstPhase := flag.Bool("firstphase", false, "")
	markAsDelete := flag.Bool("dmad", false, "")
	debug := flag.Bool("debug", false, "")
	rootAlive := flag.Bool("dontkillroot", false, "")

	flag.Parse()

	// apps := [][]string{{"diaspora", "1"}, {"mastodon", "2"}, {"gnusocial", "4"}, {"twitter", "3"}}
	apps := [][]string{{"diaspora", "1"}, {"mastodon", "2"}}

	totalApps := len(apps)

	for i := 0; i < totalApps; i++ {

		uid := *uidInput

		if i > 0 {
			uid = getUID(apps[i][1])
		}

		srcAppName, srcAppID := apps[i][0], apps[i][1]

		var dstAppID, dstAppName string
		if i == totalApps-1 {
			dstAppName, dstAppID = apps[0][0], apps[0][1]
		} else {
			dstAppName, dstAppID = apps[i+1][0], apps[i+1][1]
		}

		wg.Add(1)

		go apis.StartMigration(uid, srcAppName, srcAppID, dstAppName, dstAppID, *mtype, *blade, *bags, *ftp, *debug, *rootAlive)
		go SA1_display.StartDisplay(uid, srcAppID, dstAppID, *mtype, *threads, &wg, *display, *displayInFirstPhase, *markAsDelete, *blade)

		wg.Wait()

		// print spaces before new migration
		for j := 1; j < 50; j++ {
			fmt.Println()
		}
	}
}

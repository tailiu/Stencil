/*
 * Physical Migration Handler
 */

package main

import (
	"flag"
	"stencil/apis"
)

func main() {
	// evalConfig := evaluation.InitializeEvalConfig()

	srcApp := flag.String("srcApp", "diaspora", "")
	srcAppID := flag.String("srcAppID", "1", "")

	dstApp := flag.String("dstApp", "mastodon", "")
	dstAppID := flag.String("dstAppID", "2", "")

	// threads := flag.Int("threads", 1, "")
	mtype := flag.String("mtype", "d", "")
	uid := flag.String("uid", "", "")

	// blade := flag.Bool("blade", false, "")
	bags := flag.Bool("bags", false, "")

	flag.Parse()

	apis.StartMigrationSA2(*uid, *srcApp, *srcAppID, *dstApp, *dstAppID, *mtype, *bags)
	// threads := 1

	// enableDisplay, displayInFirstPhase, enableBags := true, true, true

	// SA2_migrate.Controller(
	// 	uid, srcApp, srcAppID, dstApp, dstAppID,
	// 	mtype, threads, enableDisplay,
	// 	displayInFirstPhase, enableBags,
	// )
}

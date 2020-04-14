/*
 * Logical Migration Handler
 */

package main

import (
	"flag"
	"stencil/SA1_display"
	"stencil/apis"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	srcApp := flag.String("srcApp", "diaspora", "")
	srcAppID := flag.String("srcAppID", "1", "")

	dstApp := flag.String("dstApp", "mastodon", "")
	dstAppID := flag.String("dstAppID", "2", "")

	threads := flag.Int("threads", 1, "")
	mtype := flag.String("mtype", "d", "")
	uid := flag.String("uid", "", "")

	blade := flag.Bool("blade", false, "")
	bags := flag.Bool("bags", false, "")
	ftp := flag.Bool("ftp", false, "")
	display := flag.Bool("display", false, "")
	displayInFirstPhase := flag.Bool("firstphase", false, "")
	markAsDelete := flag.Bool("dmad", false, "")
	debug := flag.Bool("debug", false, "")
	rootAlive := flag.Bool("dontkillroot", false, "")

	flag.Parse()

	wg.Add(1)

	go apis.StartMigration(*uid, *srcApp, *srcAppID, *dstApp, *dstAppID, *mtype, *blade, *bags, *ftp, *debug, *rootAlive)
	go SA1_display.StartDisplay(*uid, *srcAppID, *dstAppID, *mtype, *threads, &wg, *display, *displayInFirstPhase, *markAsDelete, *blade)

	wg.Wait()
}

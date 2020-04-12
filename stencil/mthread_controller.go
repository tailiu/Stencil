package main

import (
	"flag"
	"stencil/apis"
)

func main() {

	srcApp := flag.String("srcApp", "diaspora", "")
	srcAppID := flag.String("srcAppID", "1", "")

	dstApp := flag.String("dstApp", "mastodon", "")
	dstAppID := flag.String("dstAppID", "2", "")

	// threads := flag.Int("threads", 1, "")
	mtype := flag.String("mtype", "d", "")
	uid := flag.String("uid", "", "")

	// display := flag.Bool("display", false, "")
	blade := flag.Bool("blade", false, "")
	bags := flag.Bool("bags", false, "")
	debug := flag.Bool("debug", false, "")
	ftp := flag.Bool("ftp", false, "")

	flag.Parse()

	apis.StartMigration(*uid, *srcApp, *srcAppID, *dstApp, *dstAppID, *mtype, *blade, *bags, *ftp, *debug)
}

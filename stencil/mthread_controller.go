package main

import "stencil/mthread_v2" mthread

func main(){

	mtController := mthread.MigrationThreadController{
		uid: "129188",
		totalThreads: 1,
		mType: "d",
		SrcAppInfo: mthread.App{Name: "Diaspora", ID: "1"},
		DstAppInfo: mthread.App{Name: "Mastodon", ID: "2"}}

	mtController.Init()
	mtController.NewMigrationThread()

}
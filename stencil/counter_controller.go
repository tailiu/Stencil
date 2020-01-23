package main

import (
	"log"
	"stencil/counter"
)

func main() {
	appName, appID := "diaspora", "1"
	ctr := counter.CreateCounter(appName, appID)
	log.Println("Counter Created for app: ", appName, appID)
	ctr.RunCounter()

}

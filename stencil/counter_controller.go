package main

import (
	"fmt"
	"stencil/apis"
)

func main() {
	// appName, appID := "diaspora", "1"
	// ctr := counter.CreateCounter(appName, appID)
	// log.Println("Counter Created for app: ", appName, appID)
	// counter.RunCounter(&ctr)

	for i := 1; i < 1008102; i += 100 {
		uid := fmt.Sprint(i)
		apis.StartCounter("diaspora", "1", uid)
	}
}

package main

import (
	"fmt"
	"log"
	"os"
	"stencil/apis"
)

func main() {
	// appName, appID := "diaspora", "1"
	// ctr := counter.CreateCounter(appName, appID)
	// log.Println("Counter Created for app: ", appName, appID)
	// counter.RunCounter(&ctr)

	// for i := 1; i < 1008102; i += 100 {
	// 	uid := fmt.Sprint(i)
	// 	apis.StartCounter("diaspora", "1", uid, false)
	// }

	if len(os.Args) < 4 {
		fmt.Println("Not enough argurment! Need: appName, appID, uid, blade")
		log.Fatal("Args: ", os.Args)
	}

	appName, appID, uid, blade := os.Args[1], os.Args[2], os.Args[3], false

	switch os.Args[4] {
	case "t":
		blade = true
	}

	nodes, edges := apis.StartCounter(appName, appID, uid, blade)

	fmt.Println(fmt.Sprintf(">>> Nodes: %d, Edges: %d", nodes, edges))
}

package counter

import (
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
)

func CreateCounter(appName, appID string) Counter {
	AppConfig, err := config.CreateAppConfig(appName, appID)
	if err != nil {
		log.Fatal(err)
	}
	AppConfig.QR.Migration = true
	counter := Counter{
		AppConfig:     AppConfig,
		AppDBConn:     db.GetDBConn(appName),
		StencilDBConn: db.GetDBConn(db.STENCIL_DB),
		visitedNodes:  make(map[string]map[string]bool),
		NodeCount:     0,
		EdgeCount:     0}

	return counter
}

func (self *Counter) RunCounter() error {

	offset := 0

	for {
		if uid, err := db.GetNextUserFromAppDB("diaspora", "people", "id", offset); err == nil {
			if len(uid) < 1 {
				break
			}
			self.TraverseFromUser(uid)
		} else {
			fmt.Println("User offset: ", offset)
			log.Fatal("Crashed while running counter: ", err)
		}
	}
	fmt.Println("Counter Finished!")
	fmt.Println("Offset: ", offset)
	fmt.Println("Nodes: ", self.NodeCount)
	fmt.Println("Edges: ", self.EdgeCount)
	return nil
}

func (self *Counter) TraverseFromUser(uid string) error {

	return nil
}

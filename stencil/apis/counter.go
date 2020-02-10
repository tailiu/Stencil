package apis

import (
	"fmt"
	"log"
	"stencil/counter"
	"stencil/db"
)

func StartCounter(appName, appID, person_id string, isBlade ...bool) (int, int) {
	ctr := counter.CreateCounter(appName, appID, isBlade...)
	ctr.AppDBConn = db.GetDBConn(appName, isBlade...)
	ctr.UID = person_id
	ctr.EdgeCount = 0
	ctr.NodeCount = 0
	ctr.VisitedNodes = make(map[string]map[string]bool)
	fmt.Println("------------------------------------------------------------------------")
	log.Println(fmt.Sprintf("Started counter for user \"%s\" in app [%s:%s]", person_id, appID, appName))

	if personNode, err := ctr.FetchUserNode(person_id); err == nil {
		ctr.Root = personNode
		if err := ctr.Traverse(personNode); err != nil {
			log.Fatal("Error while traversing: ", err)
		}
	} else {
		fmt.Println("Passed Args: ", appName, appID, person_id)
		log.Fatal("User Node Not Created: ", err)
	}

	// if err := db.InsertIntoDAGCounter(ctr.StencilDBConn, person_id, ctr.EdgeCount, ctr.NodeCount); err != nil {
	// 	log.Fatal("Insertion Failed into DAGCOUNTER!", err)
	// }

	ctr.AppConfig.CloseDBConns()
	ctr.AppDBConn.Close()
	ctr.StencilDBConn.Close()

	log.Println(fmt.Sprintf("Finished counter for user \"%s\" in app [%s:%s] | Nodes: %d, Edges: %d", person_id, appID, appName, ctr.NodeCount, ctr.EdgeCount))
	fmt.Println("------------------------------------------------------------------------")
	return ctr.NodeCount, ctr.EdgeCount
}

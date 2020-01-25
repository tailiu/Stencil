package SA1_migrate

import (
	"stencil/apis"
	"stencil/SA1_display"
	"sync"
)

func Controller(uid, srcAppName, srcAppID, 
	dstAppName, dstAppID, migrationType string, threadNum int) {
	
	var wg sync.WaitGroup
	
	// Instead of waiting for all display threads to finish,
	// we only need to wait for one display thread to finish
	wg.Add(1)

	go apis.StartMigration(uid, srcAppName, srcAppID, dstAppName, dstAppID, migrationType)

	go SA1_display.StartDisplay(uid, srcAppID, dstAppID, migrationType, threadNum, &wg)

	wg.Wait()

}
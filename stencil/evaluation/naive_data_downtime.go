package evaluation

import (
	"fmt"
	"time"
	"log"
)

func getDowntimeBasedOnStencilMigration(data DisplayedData, appID, naiveMigrationID string, evalConfig *EvalConfig) time.Duration {
	for _, data1 := range data {
		
	}
} 

func getDataDowntimeInNaive(stencilMigrationID string, naiveMigrationID string, evalConfig *EvalConfig) {
	var dataDowntime []time.Duration 
	appID := evalConfig.MastodonAppID
	displayedData := getAllDisplayedData(evalConfig, stencilMigrationID, appID)
	for _, data := range displayedData {
		dataDowntime = append(dataDowntime, getDowntimeBasedOnStencilMigration(data, appID, naiveMigrationID, evalConfig))
	}
	return dataDowntime
}

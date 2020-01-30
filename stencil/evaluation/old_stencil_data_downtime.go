package evaluation

import (
	"stencil/db"
	"fmt"
	"time"
	"log"
)

func getDowntime(data DisplayedData, appID string, 
	evalConfig *EvalConfig) time.Duration {

	query := fmt.Sprintf(`
		select created_at, updated_at from migration_table 
		where bag = false and mark_as_delete = false 
		and app_id = %s and table_id = %s and row_id = %s`,
		appID, data.TableID, data.RowIDs[0])
	
	data1, err1 := db.DataCall1(evalConfig.StencilDBConn, query)

	log.Println(data1)
	
	if err1 != nil {
		log.Fatal(err1)
	}
	
	return data1["updated_at"].(time.Time).
		Sub(data1["created_at"].(time.Time))
	
}

func getDataDowntimeInStencil(migrationID string, 
	evalConfig *EvalConfig) []time.Duration {
	
	var dataDowntime []time.Duration 
	
	appID := evalConfig.MastodonAppID
	
	displayedData := getAllDisplayedData(evalConfig, migrationID, appID)
	
	for _, data := range displayedData {
		
		dataDowntime = append(dataDowntime, 
			getDowntime(data, appID, evalConfig))
		
	}

	return dataDowntime
}
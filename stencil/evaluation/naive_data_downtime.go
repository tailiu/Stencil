package evaluation

import (
	"stencil/config"
	"stencil/db"
	"stencil/qr"
	"fmt"
	"time"
	"log"
	"strconv"
)

func getDeletedAt(evalConfig *EvalConfig, data map[string]interface{}, naiveMigrationID string, tableName string) time.Time {
	log.Println(tableName + ".id")
	query := fmt.Sprintf("select deleted_at from evaluation where migration_id = '%s' and dst_app = '2' and dst_id = '%s' and dst_table = '%s'", 
		naiveMigrationID, data[tableName + ".id"].(string), tableName)
	log.Println(query)
	result, err := db.DataCall1(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Println(result)
	return result["deleted_at"].(time.Time)
}

func getDowntimeBasedOnStencilMigration(data DisplayedData, naiveMigrationID string, evalConfig *EvalConfig, appConfig *config.AppConfig, naiveMigrationEndTime time.Time) time.Duration {
	tableName := GetTableNameByTableID(evalConfig, data.TableID)
	data1 := getData1FromPhysicalSchemaByRowID(evalConfig, appConfig, tableName, data.RowIDs)
	log.Println(data1)
	log.Println(tableName)
	return naiveMigrationEndTime.Sub(getDeletedAt(evalConfig, data1, naiveMigrationID, tableName))
}

func getData1FromPhysicalSchemaByRowID(evalConfig *EvalConfig, appConfig *config.AppConfig, tableName string, rowIDs []string) map[string]interface{} {	
	qs := qr.CreateQS(appConfig.QR)
	qs.FromTable(map[string]string{"table": tableName, "mflag": "0"})
	qs.SelectColumns(tableName + ".*")
	var strRowIDs string 
	for i, rowID := range rowIDs {
		if i == 0 {
			strRowIDs += rowID
		} else {
			strRowIDs += "," + rowID
		}
	}
	qs.RowIDs(strRowIDs)
	physicalQuery := qs.GenSQL()

	log.Println(physicalQuery)
	
	result, err := db.DataCall1(evalConfig.StencilDBConn, physicalQuery)
	if err != nil {
		log.Fatal(err)
	}
	
	return result
}

func getDataDowntimeInNaive(stencilMigrationID string, naiveMigrationID string, evalConfig *EvalConfig) []time.Duration {
	var dataDowntime []time.Duration 
	appConfig := getAppConfig(evalConfig, "mastodon")
	intNaiveMigrationID, err := strconv.Atoi(naiveMigrationID)
	if err != nil {
		log.Fatal(err)
	}
	naiveMigrationEndTime := getMigrationEndTime(evalConfig.StencilDBConn, intNaiveMigrationID)
	displayedData := getAllDisplayedData(evalConfig, stencilMigrationID, appConfig.AppID)
	for _, data := range displayedData {
		dataDowntime = append(dataDowntime, getDowntimeBasedOnStencilMigration(data, naiveMigrationID, evalConfig, appConfig, naiveMigrationEndTime))
	}
	return dataDowntime
}

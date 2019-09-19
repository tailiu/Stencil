package evaluation

import (
	"stencil/db"
	"stencil/config"
	"stencil/qr"
	"log"
	"fmt"
	"strings"
)

func getColsSizeOfDataInStencilModel(evalConfig *EvalConfig, appConfig *config.AppConfig, tableID string, rowIDs []string) map[string]interface{} {
	qs := qr.CreateQS(appConfig.QR)
	tableName := GetTableNameByTableID(evalConfig, tableID)
	qs.FromTable(map[string]string{"table":tableName, "mflag": "0", "mark_as_delete": "false", "bag": "false"})
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
	physicalQuery := qs.GenSQLSize()

	// log.Println(physicalQuery)

	result, err := db.DataCall1(evalConfig.StencilDBConn, physicalQuery)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)

	return result
}

func filterColsAndResultsBasedOnSchemaMapping(data map[string]interface{}, srcApp, dstApp string) int64 {
	var size int64
	for k, v := range data {
		if strings.Contains(k, ".mark_as_delete") {
			continue
		} else {
			size += v.(int64)
		}
	}
	return size
}

func calculateDisplayedDataSizeInBagEvaluation(evalConfig *EvalConfig, appConfig *config.AppConfig, srcApp, dstApp string, displayedData []DisplayedData) int64 {
	var size int64
	for _, data := range displayedData {
		size += filterColsAndResultsBasedOnSchemaMapping(getColsSizeOfDataInStencilModel(evalConfig, appConfig, data.TableID, data.RowIDs), srcApp, dstApp)
	}
	return size
}

func getDisplayedDataSize(evalConfig *EvalConfig, srcApp, dstApp, migrationID string) int64 {
	dstAppConfig := getAppConfig(evalConfig, dstApp)
	displayedData := getAllDisplayedData(evalConfig, migrationID, dstAppConfig.AppID)
	log.Println(displayedData)
	return calculateDisplayedDataSizeInBagEvaluation(evalConfig, dstAppConfig, srcApp, dstApp, displayedData)
}

// We use dstApp here to get the total migrated node size in the source application
func getTotalMigratedNodeSize(evalConfig *EvalConfig, dstApp string, migrationID string) int64 {
	dstAppID := db.GetAppIDByAppName(evalConfig.StencilDBConn, dstApp)
	query := fmt.Sprintf("select msize from migration_registration where dst_app = %s and migration_id = %s", dstAppID, migrationID)
	result, err := db.DataCall1(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)
	// return result["misze"].(int64)
	return 0
}
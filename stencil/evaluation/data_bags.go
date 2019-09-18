package evaluation

import (
	"stencil/db"
	"stencil/config"
	"stencil/qr"
	"log"
	"fmt"
)

func calculateOneDataSizeInStencilModel(evalConfig *EvalConfig, appConfig *config.AppConfig, tableID string, rowIDs []string) int64 {
	qs := qr.CreateQS(appConfig.QR)
	tableName := GetTableNameByTableID(evalConfig, tableID)
	qs.FromTable(map[string]string{"table":tableName, "mflag": "0", "mark_as_delete": "true", "bag": "true"})
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

	var size int64
	for _, v := range result {
		size += v.(int64)
	}

	return size
}

func calculateDisplayedDataSizeInStencilModel(evalConfig *EvalConfig, appConfig *config.AppConfig, displayedData []DisplayedData) int64 {
	var size int64
	for _, data := range displayedData {
		size += calculateOneDataSizeInStencilModel(evalConfig, appConfig, data.TableID, data.RowIDs)
	}
	return size
}

func getDisplayedDataSize(evalConfig *EvalConfig, app, migrationID string) int64 {
	appConfig := getAppConfig(evalConfig, app)
	displayedData := getAllDisplayedData(evalConfig, migrationID, appConfig.AppID)
	return calculateDisplayedDataSizeInStencilModel(evalConfig, appConfig, displayedData)
}

func getTotalMigratedNodeSize(evalConfig *EvalConfig, app, migrationID string) int64 {
	app_id := db.GetAppIDByAppName(evalConfig.StencilDBConn, app)
	query := fmt.Sprintf("select msize from migration_registration where dst_app = %s and migration_id = %s", app_id, migrationID)
	result, err := db.DataCall1(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)
	return result["misze"].(int64)
}
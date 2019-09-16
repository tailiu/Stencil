package evaluation

import (
	"stencil/db"
	"stencil/config"
	"stencil/qr"
	"log"
	"fmt"
	"strings"
)

func getAppConfig(evalConfig *EvalConfig, app string) *config.AppConfig {
	app_id := db.GetAppIDByAppName(evalConfig.StencilDBConn, app)
	appConfig, err := config.CreateAppConfigDisplay(app, app_id)
	if err != nil {
		log.Fatal(err)
	}
	return &appConfig
}

func getAllDataInDataBag(evalConfig *EvalConfig, migrationID string, appConfig *config.AppConfig) []DataBagData {
	query := fmt.Sprintf("select table_id, array_agg(row_id) as row_ids from migration_table where bag = true and app_id = %s and migration_id = %s group by group_id, table_id;",
		appConfig.AppID, migrationID)
	
	data := db.GetAllColsOfRows(evalConfig.StencilDBConn, query)

	var dataBag []DataBagData
	for _, data1 := range data {
		var rowIDs []string
		s := data1["row_ids"][1:len(data1["row_ids"]) - 1]
		s1 := strings.Split(s, ",")
		for _, rowID := range s1 {
			rowIDs = append(rowIDs, rowID)
		}

		dataBagData := DataBagData{}
		dataBagData.TableID = data1["table_id"]
		dataBagData.RowIDs = rowIDs
		dataBag = append(dataBag, dataBagData)
	}

	log.Println(dataBag)
	return dataBag
}

func GetTableNameByTableID(evalConfig *EvalConfig, tableID string) string {
	query := fmt.Sprintf("select table_name from app_tables where pk = %s", tableID)
	data1, err1 := db.DataCall1(evalConfig.StencilDBConn, query)
	if err1 != nil {
		log.Fatal(err1)
	}
	return data1["table_name"].(string)
}

func calculateOneDataSizeInStencilModel(evalConfig *EvalConfig, appConfig *config.AppConfig, tableID string, rowIDs []string) int {
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
	physicalQuery1 := qs.GenSQL()

	log.Println(physicalQuery)
	log.Println(physicalQuery1)

	result, err := db.DataCall1(evalConfig.StencilDBConn, physicalQuery)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)

	return 0
}

func calculateDataSizeInStencilModel(evalConfig *EvalConfig, appConfig *config.AppConfig, dataBag []DataBagData) int {
	size := 0
	for _, data := range dataBag {
		size += calculateOneDataSizeInStencilModel(evalConfig, appConfig, data.TableID, data.RowIDs)
	}
	return size
}

func getDataBagSize(evalConfig *EvalConfig, app, migrationID string) int {
	appConfig := getAppConfig(evalConfig, app)
	dataBag := getAllDataInDataBag(evalConfig, migrationID, appConfig)
	return calculateDataSizeInStencilModel(evalConfig, appConfig, dataBag)
}
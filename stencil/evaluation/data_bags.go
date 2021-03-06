package evaluation

import (
	"stencil/db"
	"stencil/config"
	"stencil/qr"
	"log"
	"fmt"
	"strings"
)

func getAllDataInDataBag(evalConfig *EvalConfig, userID string, appConfig *config.AppConfig) []DataBagData {
	query := fmt.Sprintf("select table_id, array_agg(row_id) as row_ids from migration_table where bag = true and app_id = %s and user_id = %s group by group_id, table_id;",
		appConfig.AppID, userID)
	
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

	// log.Println(dataBag)
	return dataBag
}

func calculateOneDataSizeInStencilModel(evalConfig *EvalConfig, appConfig *config.AppConfig, tableID string, rowIDs []string) int64 {
	qs := qr.CreateQS(appConfig.QR)
	tableName := GetTableNameByTableID(evalConfig, tableID)
	qs.FromTable(map[string]string{"table":tableName, "mflag": "0,1", "mark_as_delete": "true", "bag": "true"})
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

	// log.Println(result)

	var size int64
	for k, v := range result {
		if result[k] == nil {
			continue
		} else {
			size += v.(int64)
		}
	} 

	return size
}

func calculateDataSizeInStencilModel(evalConfig *EvalConfig, appConfig *config.AppConfig, dataBag []DataBagData) int64 {
	var size int64
	for _, data := range dataBag {
		size += calculateOneDataSizeInStencilModel(evalConfig, appConfig, data.TableID, data.RowIDs)
	}
	return size
}

func getAllAppsOfDataBag(evalConfig *EvalConfig, userID string) []string {
	query := fmt.Sprintf("select distinct a.app_name from migration_table m join apps a on m.app_id = a.pk where user_id = %s and bag = true;",
		userID)
	
	result, err := db.DataCall(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	var apps []string 
	for _, app := range result {
		apps = append(apps, fmt.Sprint(app["app_name"]))
	}
	return apps
}

func getDataBagSize(evalConfig *EvalConfig, app, userID string) int64 {
	appConfig := getAppConfig(evalConfig, app)
	dataBag := getAllDataInDataBag(evalConfig, userID, appConfig)
	return calculateDataSizeInStencilModel(evalConfig, appConfig, dataBag)
}
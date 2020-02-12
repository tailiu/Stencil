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

func getDeletedAt(evalConfig *EvalConfig, 
	data map[string]interface{}, naiveMigrationID string, 
	tableName string) time.Time {
	
	// log.Println(tableName + ".id")
	query := fmt.Sprintf(
		`select deleted_at from evaluation 
		where migration_id = '%s' and dst_app = '2' 
		and dst_id = '%s' and dst_table = '%s'`, 
		naiveMigrationID, 
		fmt.Sprint(data[tableName + ".id"]), 
		tableName,
	)

	log.Println(query)

	result, err := db.DataCall1(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Println(result)
	
	if result["deleted_at"] == nil {
		return time.Time{}
	} else {
		return result["deleted_at"].(time.Time)
	} 
}

func getData1FromPhysicalSchemaByRowID(evalConfig *EvalConfig, 
	appConfig *config.AppConfig, tableName string, 
	rowIDs []string) map[string]interface{} {
	
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

	// log.Println(physicalQuery)
	
	result, err := db.DataCall1(evalConfig.StencilDBConn, physicalQuery)
	if err != nil {
		log.Fatal(err)
	}
	
	return result
}

func getDataDowntimeInNaive(stencilMigrationID string, 
	naiveMigrationID string, 
	evalConfig *EvalConfig) []time.Duration {
	
	var dataDowntime []time.Duration 
	
	appConfig := getAppConfig(evalConfig, "mastodon")
	
	intNaiveMigrationID, err := strconv.Atoi(naiveMigrationID)
	if err != nil {
		log.Fatal(err)
	}
	
	naiveMigrationEndTime := getMigrationEndTime(evalConfig.StencilDBConn, intNaiveMigrationID)
	
	displayedData := getAllDisplayedData(evalConfig, stencilMigrationID, appConfig.AppID)
	
	for _, data := range displayedData {
		
		tableName := GetTableNameByTableID(evalConfig, data.TableID)
		
		data1 := getData1FromPhysicalSchemaByRowID(evalConfig, 
			appConfig, tableName, data.RowIDs)
		
		// log.Println(data1)
		// log.Println(tableName)
		
		deletedAt := getDeletedAt(evalConfig, data1, 
			naiveMigrationID, tableName)
		
		if deletedAt.IsZero() {
			log.Println("GOT ONE ZERO", data1, tableName)
			continue
		} else {
			dataDowntime = append(dataDowntime, 
				naiveMigrationEndTime.Sub(deletedAt))
		}
	}
	return dataDowntime
}

func GetNaiveAndStencilDataDowntimeMigrations(evalConfig *EvalConfig) []map[string]interface{} {

	query := fmt.Sprintf(
		`select m2.migration_id as stencil, m1.migration_id as naive 
		from migration_registration as m1 join migration_registration as m2 
		on m1.user_id = m2.user_id 
		where m1.src_app = 1 and m1.dst_app = 2 and 
		m1.is_logical = true and m2.is_logical = false and 
		m2.src_app = 1 and m2.dst_app = 2;`)

	result1, err1 := db.DataCall(evalConfig.StencilDBConn, query)
	if err1 != nil {
		log.Fatal(err1)
	}

	return result1
}
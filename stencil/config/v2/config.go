package config

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"stencil/db"
	"stencil/qr"
	"strings"
	"time"
)

var SchemaMappingsObj *SchemaMappings

func CreateAppConfig(app, app_id string, isBlade ...bool) (AppConfig, error) {

	var appConfig AppConfig
	var dconfig string

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(dir, "/stencil/") {
		dconfig = "../config/dependencies/" + app + ".json"
	} else {
		dconfig = "./config/dependencies/" + app + ".json"
	}

	jsonFile, err := os.Open(dconfig)
	defer jsonFile.Close()

	if err != nil {
		fmt.Println("Some problem with the file: ")
		log.Fatal(err)
		return appConfig, errors.New("can't open file")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &appConfig)

	appConfig.AppName = app
	appConfig.AppID = app_id
	appConfig.DBConn = db.GetDBConn(app, isBlade...)

	rand.Seed(time.Now().UTC().UnixNano())
	appConfig.Rand = rand.New(rand.NewSource(time.Now().Unix()))
	return appConfig, nil
}

func CreateAppConfigDisplay(
	app, app_id string, stencilDBConn *sql.DB, newDB bool) (AppConfig, error) {

	var appConfig AppConfig
	var dconfig string

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// if strings.Contains(dir, "/stencil/") {
	// 	dconfig = "../config/dependencies/" + app + "_display.json"
	// } else {
	// 	dconfig = "./config/dependencies/" + app + "_display.json"
	// }

	// Use the combined display and migration dependency files
	if strings.Contains(dir, "/stencil/") {
		dconfig = "../config/dependencies/" + app + ".json"
	} else {
		dconfig = "./config/dependencies/" + app + ".json"
	}

	jsonFile, err := os.Open(dconfig)
	defer jsonFile.Close()

	if err != nil {
		fmt.Println("Some problem with the file: ")
		fmt.Println(err)
		return appConfig, errors.New("can't open file")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &appConfig)

	appConfig.AppName = app

	appConfig.AppID = app_id

	if newDB {
		appConfig.DBConn = db.GetDBConn(app)
	} else {
		appConfig.DBConn = db.GetDBConn(app, true)
	}

	if app_id != "" {
		appConfig.QR = qr.NewQR(app, app_id)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	appConfig.Rand = rand.New(rand.NewSource(time.Now().Unix()))

	tableIDNamePairs := getTableIDNamePairsInApp(stencilDBConn, app_id)

	appConfig.TableIDNamePairs = make(map[string]string)

	appConfig.TableNameIDPairs = make(map[string]string)

	for _, tableIDNamePair := range tableIDNamePairs {

		appConfig.TableIDNamePairs[fmt.Sprint(tableIDNamePair["pk"])] = fmt.Sprint(tableIDNamePair["table_name"])

		appConfig.TableNameIDPairs[fmt.Sprint(tableIDNamePair["table_name"])] = fmt.Sprint(tableIDNamePair["pk"])

	}

	return appConfig, nil
}

func getTableIDNamePairsInApp(stencilDBConn *sql.DB, app_id string) []map[string]interface{} {
	query := fmt.Sprintf("select pk, table_name from app_tables where app_id = %s", app_id)

	result, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func LoadSchemaMappings() (*SchemaMappings, error) {
	if SchemaMappingsObj == nil {

		SchemaMappingsObj = new(SchemaMappings)

		schemaMappingFile := build.Default.GOPATH + "/src/stencil/config/app_settings/mappings.json"
		jsonFile, err := os.Open(schemaMappingFile)
		if err != nil {
			fmt.Println(err)
			return SchemaMappingsObj, errors.New("can't open schema mapping json file")
		}
		// defer the closing of our jsonFile so that we can parse it later on

		jsonAsBytes, _ := ioutil.ReadAll(jsonFile)

		json.Unmarshal(jsonAsBytes, SchemaMappingsObj)

		dbConn := db.GetDBConn(db.STENCIL_DB)
		defer dbConn.Close()

		for i, mapping := range SchemaMappingsObj.AllMappings {
			for j, toApp := range mapping.ToApps {
				appID := db.GetAppIDByAppName(dbConn, toApp.Name)
				for k, toAppMapping := range toApp.Mappings {
					for l, toTable := range toAppMapping.ToTables {
						ToTableID, err := db.TableID(dbConn, toTable.Table, appID)
						if err != nil {
							fmt.Println("LoadSchemaMappings: Unable to resolve ToTableID for table: ", toTable.Table, toApp.Name, appID)
							log.Fatal(err)
						}
						SchemaMappingsObj.AllMappings[i].ToApps[j].Mappings[k].ToTables[l].TableID = ToTableID
						// fmt.Println(toTable.Table, toApp.Name, appID, ToTableID)
					}
				}
			}
		}
		// fmt.Println(SchemaMappingsObj.AllMappings[0].ToApps[0].Mappings[0].ToTables)
		// log.Fatal()
		jsonFile.Close()
	}
	return SchemaMappingsObj, nil
}

func GetSelfSchemaMappings(dbConn *sql.DB, appID, appName string) *MappedApp {
	mappedApp := new(MappedApp)
	mappedApp.Name = appName

	if res, err := db.GetTablesForApp(dbConn, appID); err == nil {
		var mappings []Mapping
		for _, row := range res {
			tableID := fmt.Sprint(row["table_id"])
			if tableName, err := db.TableName(dbConn, tableID, appID); err == nil {
				if columnsRes, err := db.GetColumnsFromAppSchema(dbConn, tableID); err == nil {
					var toTable ToTable
					toTable.Table = tableName
					toTable.TableID = tableID
					toTable.Mapping = make(map[string]string)
					for _, columnRow := range columnsRes {
						column := fmt.Sprint(columnRow["column_name"])
						toTable.Mapping[column] = tableName + "." + column
					}
					mappings = append(mappings, Mapping{FromTables: []string{tableName}, ToTables: []ToTable{toTable}})
				} else {
					log.Fatal("@Columns:", err)
				}
			} else {
				log.Fatal("@TableName:", err)
			}
		}
		mappedApp.Mappings = mappings
	}
	return mappedApp
}

func GetSchemaMappingsFor(srcApp, dstApp string) *MappedApp {
	if schemaMappings, err := LoadSchemaMappings(); err == nil {
		for _, schemaMapping := range schemaMappings.AllMappings {
			if strings.EqualFold(srcApp, schemaMapping.FromApp) {
				for _, mappedApp := range schemaMapping.ToApps {
					if strings.EqualFold(dstApp, mappedApp.Name) {
						return &mappedApp
					}
				}
			}
		}
	}
	return nil
}

func (self *AppConfig) ShuffleDependencies(vals []Dependency) []Dependency {
	// r := Init()
	ret := make([]Dependency, len(vals))
	perm := self.Rand.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}

func (self *AppConfig) ShuffleOwnerships(vals []Ownership) []Ownership {
	// r := Init()
	ret := make([]Ownership, len(vals))
	perm := self.Rand.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}

func Reverse(numbers []DataQuery) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}

func remove(s []Tag, i int) []Tag {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

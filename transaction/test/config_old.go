/*
 * Configuration Reader/Exporter
 */

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"transaction/db"

	"github.com/tidwall/gjson"
)

/*********************--bgn
 * Structures
***************************/

type Mapping map[string]map[string]map[string]map[string]string

type DataQuery struct {
	SQL, Table string
}

type Dependencies struct {
	Dependencies []Dependency `json:"dependencies"`
}

type Dependency struct {
	Tag        string      `json:"tag"`
	DependsOn  string      `json:"depends_on"`
	Conditions []Condition `json:"conditions"`
}

type Conditions struct {
	Conditions []Condition `json:"conditions"`
}

type Condition struct {
	TagAttr       string `json:"tag_attr"`
	DependsOnAttr string `json:"depends_on_attr"`
}

type Settings struct {
	UserTable string
	KeyCol    string
	Mappings  Mapping
}

/*********************--end
 * Structures
***************************/

/*********************--bgn
 * Functions
***************************/

func FindDependency(tag, depends_on string, dependencies []Dependency) (Dependency, error) {

	for _, dependency := range dependencies {
		if strings.ToLower(dependency.Tag) == strings.ToLower(tag) &&
			strings.ToLower(dependency.DependsOn) == strings.ToLower(depends_on) {
			return dependency, nil
		}
	}
	return *new(Dependency), errors.New("dependency doesn't exist")
}

func FindDependencyByDependsOn(depends_on string, dependencies []Dependency) (Dependency, error) {

	for _, dependency := range dependencies {
		if strings.ToLower(dependency.Tag) == strings.ToLower(depends_on) {
			return dependency, nil
		}
	}
	return *new(Dependency), errors.New("dependency doesn't exist")
}

func ReadDependencies(app string) ([]Dependency, error) {

	var dependencies Dependencies
	dconfig := "./config/dependencies/" + app + ".json"
	jsonFile, err := os.Open(dconfig)

	if err != nil {
		fmt.Println("Some problem with the file: ")
		fmt.Println(err)
		return dependencies.Dependencies, errors.New("can't open file")
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &dependencies)

	return dependencies.Dependencies, nil
}

func GetSchemaMappingsFromDB(app string) Mapping {

	sql := "select apps1.app_name app1, as1.table_name table1, as1.column_name col1, apps2.app_name app2, as2.table_name table2, as2.column_name col2 from schema_mappings sm join app_schemas as1 on sm.source_attribute = as1.row_id join app_schemas as2 on sm.dest_attribute = as2.row_id join apps apps1 on apps1.row_id = as1.app_id join apps apps2 on apps2.row_id = as2.app_id where apps1.app_name = $1"

	mappings := make(Mapping)

	for _, row := range db.DataCall("stencil", sql, app) {
		app2 := row["app2"]
		table1 := row["table1"]
		table2 := row["table2"]
		col1 := row["col1"]
		col2 := row["col2"]
		if _, ok := mappings[app2]; !ok {
			mappings[app2] = make(map[string]map[string]map[string]string)
		}
		if _, ok := mappings[app2][table1]; !ok {
			mappings[app2][table1] = make(map[string]map[string]string)
		}
		if _, ok := mappings[app2][table1][table2]; !ok {
			mappings[app2][table1][table2] = make(map[string]string)
		}
		mappings[app2][table1][table2][col1] = col2
	}

	return mappings
}

func ReadAppSettings(app string, readMappingsFromJSON bool) (Settings, error) {

	var settings Settings
	appSettingsFile := "./config/app_settings/" + app + ".json"
	jsonAsBytes, err := ioutil.ReadFile(appSettingsFile)
	if err != nil {
		fmt.Println(err)
		return settings, errors.New("can't open file")
	}
	json := string(jsonAsBytes)
	settings.UserTable = gjson.Get(json, "user_table").String()
	settings.KeyCol = gjson.Get(json, "key_column").String()
	if readMappingsFromJSON {
		settings.Mappings = GetSchemaMappingsFromJSON(json)
	} else {
		settings.Mappings = GetSchemaMappingsFromDB(app)
	}
	return settings, nil
}

func GetSchemaMappingsFromJSON(json string) Mapping {

	returnmap := make(Mapping)
	mappings := gjson.Get(json, "mappings")
	mappings.ForEach(func(appName, appJSON gjson.Result) bool {
		returnmap[appName.String()] = make(map[string]map[string]map[string]string)
		appXPath := fmt.Sprintf("mappings.%s", appName.String())
		appMapping := gjson.Get(json, appXPath)
		appMapping.ForEach(func(tableName, tableVal gjson.Result) bool {
			returnmap[appName.String()][tableName.String()] = make(map[string]map[string]string)
			tabXPath := fmt.Sprintf(appXPath+".%s", tableName.String())
			tabMapping := gjson.Get(json, tabXPath)
			tabMapping.ForEach(func(mTabName, mTabVal gjson.Result) bool {
				returnmap[appName.String()][tableName.String()][mTabName.String()] = make(map[string]string)
				mTabXPath := fmt.Sprintf(tabXPath+".%s", mTabName.String())
				mTabMapping := gjson.Get(json, mTabXPath)
				mTabMapping.ForEach(func(colName, colMapping gjson.Result) bool {
					returnmap[appName.String()][tableName.String()][mTabName.String()][colName.String()] = colMapping.String()
					return true
				})
				return true
			})
			return true
		})
		return true
	})
	return returnmap
}

/*********************--end
 * Functions
***************************/

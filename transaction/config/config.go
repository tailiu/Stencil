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

func ReadAppSettings(app string) (Settings, error) {

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

	mappings := gjson.Get(json, "mappings")

	settings.Mappings = make(map[string]map[string]map[string]map[string]string)

	mappings.ForEach(func(appName, appJSON gjson.Result) bool {

		settings.Mappings[appName.String()] = make(map[string]map[string]map[string]string)

		appXPath := fmt.Sprintf("mappings.%s", appName.String())

		appMapping := gjson.Get(json, appXPath)

		appMapping.ForEach(func(tableName, tableVal gjson.Result) bool {

			settings.Mappings[appName.String()][tableName.String()] = make(map[string]map[string]string)

			tabXPath := fmt.Sprintf(appXPath+".%s", tableName.String())

			tabMapping := gjson.Get(json, tabXPath)

			tabMapping.ForEach(func(mTabName, mTabVal gjson.Result) bool {

				settings.Mappings[appName.String()][tableName.String()][mTabName.String()] = make(map[string]string)

				mTabXPath := fmt.Sprintf(tabXPath+".%s", mTabName.String())

				mTabMapping := gjson.Get(json, mTabXPath)

				mTabMapping.ForEach(func(colName, colMapping gjson.Result) bool {

					settings.Mappings[appName.String()][tableName.String()][mTabName.String()][colName.String()] = colMapping.String()
					return true
				})

				return true
			})

			return true
		})

		return true
	})

	return settings, nil
}

/*********************--end
 * Functions
***************************/

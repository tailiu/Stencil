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
)

/****************** Dependencies Functions ***********************/

func FindDependency(tag, depends_on string, dependencies []Dependency) (Dependency, error) {

	for _, dependency := range dependencies {
		if strings.ToLower(dependency.Tag) == strings.ToLower(tag) {
			// && strings.ToLower(dependency.DependsOn) == strings.ToLower(depends_on)

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

func CreateAppConfig(app string) (AppConfig, error) {

	var appConfig AppConfig
	dconfig := "./config/dependencies/" + app + ".json"
	jsonFile, err := os.Open(dconfig)

	if err != nil {
		fmt.Println("Some problem with the file: ")
		fmt.Println(err)
		return appConfig, errors.New("can't open file")
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &appConfig)

	appConfig.AppName = app

	return appConfig, nil
}

/****************** Shema Mappings Functions ***********************/

func ReadSchemaMappingSettings(fileName string) (SchemaMappings, error) {
	var schemaMappings SchemaMappings

	schemaMappingFile := "./config/app_settings/" + fileName + ".json"
	jsonFile, err := os.Open(schemaMappingFile)
	if err != nil {
		fmt.Println(err)
		return schemaMappings, errors.New("can't open schema mapping json file")
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	jsonAsBytes, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(jsonAsBytes, &schemaMappings)

	return schemaMappings, nil
}

/*********************--end
 * Functions
***************************/

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

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

func GetSchemaMappings() (SchemaMappings, error) {
	var schemaMappings SchemaMappings

	schemaMappingFile := "./config/app_settings/mappings.json"
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

func GetSchemaMappingsFor(srcApp, dstApp string) *MappedApp {
	if schemaMappings, err := GetSchemaMappings(); err == nil {
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

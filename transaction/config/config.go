package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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

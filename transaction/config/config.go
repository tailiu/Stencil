package config

import (
	"diaspora/db"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
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
	appConfig.DBConn = db.GetDBConn(app)

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

func ShuffleDependencies(vals []Dependency) []Dependency {
	rand.Seed(time.Now().UnixNano())
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]Dependency, len(vals))
	perm := r.Perm(len(vals))
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

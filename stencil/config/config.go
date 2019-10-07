package config

import (
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

func CreateAppConfig(app, app_id string) (AppConfig, error) {

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

	if err != nil {
		fmt.Println("Some problem with the file: ")
		fmt.Println(err)
		return appConfig, errors.New("can't open file")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()
	json.Unmarshal(byteValue, &appConfig)

	appConfig.AppName = app
	appConfig.AppID = app_id
	appConfig.DBConn = db.GetDBConn(app)

	if app_id != "" {
		appConfig.QR = qr.NewQR(app, app_id)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	appConfig.Rand = rand.New(rand.NewSource(time.Now().Unix()))
	return appConfig, nil
}

func CreateAppConfigDisplay(app, app_id string, newDB bool) (AppConfig, error) {

	var appConfig AppConfig
	var dconfig string

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(dir, "/stencil/") {
		dconfig = "../config/dependencies/" + app + "_display.json"
	} else {
		dconfig = "./config/dependencies/" + app + "_display.json"
	}

	jsonFile, err := os.Open(dconfig)

	if err != nil {
		fmt.Println("Some problem with the file: ")
		fmt.Println(err)
		return appConfig, errors.New("can't open file")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()
	json.Unmarshal(byteValue, &appConfig)

	appConfig.AppName = app
	appConfig.AppID = app_id
	if newDB {
		appConfig.DBConn = db.GetDBConn(app)
	} else {
		appConfig.DBConn = db.GetDBConn2(app)
	}

	if app_id != "" {
		appConfig.QR = qr.NewQR(app, app_id)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	appConfig.Rand = rand.New(rand.NewSource(time.Now().Unix()))
	return appConfig, nil
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

		jsonFile.Close()
	}
	return SchemaMappingsObj, nil
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

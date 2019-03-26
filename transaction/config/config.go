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

/*********************--bgn
 * Structures
***************************/

/****************** Shema Mappings Structs ***********************/

type SchemaMappings struct {
	AllMappings []SchemaMapping		`json:"allMappings"`
}

type SchemaMapping struct {
	FromApp		string				`json:"fromApp"`
	VarsFuncs	VarsFuncsConfig		`json:"varsFuncs"`
	ToApps		[]MappedToApp		`json:"toApps"`
}

type VarsFuncsConfig struct {
	Funcs		[]Func				`json:"funcs"`
	Vars		[]Vars				`json:"vars"`
}

type Func struct {
	Name				string 		`json:"name"`
	MappingToFunc		string 		`json:"mapping"`
}

type Vars struct {
	Name				string 		`json:"name"`
	MappingToVar		string 		`json:"mapping"`
}

type MappedToApp struct {
	Name		string 		`json:"app"`
	Mappings	[]Mapping	`json:"mappings"`
}

type Mapping struct {
	FromTables	string 		`json:"fromTables"`
	ToTables	[]ToTable	`json:"toTables"`
}

type ToTable struct {
	Table		string 				`json:"table"`
	Conditions	map[string]string	`json:"conditions"`
	Mapping		map[string]string	`json:"mapping"`
}

/****************** Dependencies Structs ***********************/

type App struct {
	Tables		[]map[string]string	`json:""`
}

type AppConfig struct {
	Tags         []Tag        `json:"tags"`
	Dependencies []Dependency `json:"dependencies"`
	Ownerships   []Ownership  `json:"ownership"`
}

type Ownership struct {
	Tag        string              `json:"tag"`
	DependsOn  string              `json:"owned_by"`
	Conditions []map[string]string `json:"conditions"`
}

type Tag struct {
	Name              string              `json:"name"`
	Members           map[string]string   `json:"members"`
	Keys              map[string]string   `json:"keys"`
	InnerDependencies []map[string]string `json:"inner_dependencies"`
}

type Dependency struct {
	Tag       string      `json:"tag"`
	DependsOn []DependsOn `json:"depends_on"`
}

type DependsOn struct {
	Tag        string       `json:"tag"`
	Conditions []DCondition `json:"conditions"`
}

type DCondition struct {
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

func ReadAppConfig(app string) (AppConfig, error) {

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

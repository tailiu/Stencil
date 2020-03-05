package common_funcs

import (
	"stencil/config"
	"os"
	"strings"
	"fmt"
	"log"
	"errors"
	"encoding/json"
	"io/ioutil"
)

func ReplaceKey(dag *DAG, tag string, key string) string {

	for _, tag1 := range dag.Tags {

		if tag1.Name == tag {
			// fmt.Println(tag)

			for k, v := range tag1.Keys {

				if k == key {

					member := strings.Split(v, ".")[0]
					
					attr := strings.Split(v, ".")[1]
					
					for k1, table := range tag1.Members {

						if k1 == member {

							return table + "." + attr
						}
					}
				}
			}
		}
	}

	return ""

}

func GetTableByMemberID(dag *DAG, tagName string, checkedMemberID string) (string, error) {

	for _, tag := range dag.Tags {
		if tag.Name == tagName {
			for memberID, memberTable := range tag.Members {
				if memberID == checkedMemberID {
					return memberTable, nil
				}
			}
		}
	}

	return "", NoTableFound
}

func GetDepDisplaySetting(dag *DAG, tag string, pTag string) (string, error) {

	for _, dependency := range dag.Dependencies {

		if dependency.Tag == tag {

			for _, dependsOn := range dependency.DependsOn {

				if dependsOn.As != "" {

					if dependsOn.As == pTag {

						return dependsOn.DisplaySetting, nil

					} else {

						continue

					}
				} else {

					if dependsOn.Tag == pTag {

						return dependsOn.DisplaySetting, nil
					}
				}
			}
		}
	}

	return "", CannotFindDependencyDisplaySetting
}

func GetDependsOnConditionsInDeps(dag *DAG, tagName string, 
	pTagName string) ([]config.DCondition, error) {
	
	for _, dp := range dag.Dependencies {

		if dp.Tag == tagName {
			
			for _, dp1 := range dp.DependsOn {
				
				if dp1.As == pTagName {
					
					return dp1.Conditions, nil
				
				} else if dp1.Tag == pTagName {
					
					return dp1.Conditions, nil
				}
			}
		}
	}

	return nil, errors.New("Error: No Conditions Found")
}

func LoadDAG(app string) (*DAG, error) {
	
	var dag DAG

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
		return nil, errors.New("can't open file")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()
	json.Unmarshal(byteValue, &dag)

	return &dag, nil

}
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

func (dag *DAG) ReplaceKey(tag string, key string) string {

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

func (dag *DAG) GetTableByMemberID(tagName string, checkedMemberID string) (string, error) {

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

func (dag *DAG) GetDepDisplaySetting(tag string, pTag string) (string, error) {

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

func (dag *DAG) GetDependsOnConditionsInDeps(tagName string, pTagName string) ([]config.DCondition, error) {
	
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

func (dag *DAG) GetRootMemberAttr() (string, string, error) {

	for _, tag1 := range dag.Tags {
		
		if tag1.Name == "root" {

			for k, v := range tag1.Keys {

				if k == "root_id" {

					memberNum := strings.Split(v, ".")[0]
					
					attr := strings.Split(v, ".")[1]
					
					for k1, member := range tag1.Members {

						if k1 == memberNum {

							return member, attr, nil
						}
					}
				}
			}
		}
	}
	
	return "", "", CannotFindRootMemberAttr
}

func (dag *DAG) IfDependsOn(table, attr string) bool {

	var tag, member, key string

	// Check inner-node dependencies
	for _, tag1 := range dag.Tags {
		for member1, table1 := range tag1.Members {
			if table1 == table {
				tag = tag1.Name
				member = member1
				for _, innerDependency := range tag1.InnerDependencies {
					for _, dependsOn := range innerDependency {
						if dependsOn == member1 + "." + attr {
							return true
						}
					}
				}
				for k, v := range tag1.Keys {
					if v == member1 + "." + attr {
						key = k
					}
				}
			}
		}
	} 
	
	log.Println(tag, member, key)

	// Check inter-node dependencies
	for _, dep := range dag.Dependencies {
		if dep.Tag == tag {
			for _, dependsOn := range dep.DependsOn {
				// For now we only consider one condition in conditions
				for i, condition := range dependsOn.Conditions {
					if condition.TagAttr == key && i == 0 {
						return true
					}
				}
			}
		}
	}

	// Check ownership
	for _, ownership := range dag.Ownerships {
		if ownership.Tag == tag {
			// For now we only consider one condition in conditions
			for i, condition := range ownership.Conditions {
				if condition.TagAttr == key && i == 0 {
					return true
				}
			}
		}
	}
	
	return false

}
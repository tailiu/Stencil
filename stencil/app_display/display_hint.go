package app_display

import (
	"errors"
	"stencil/config"
	"strconv"
	"log"
	"fmt"
)

// The Key should be the primay key of the Table
type HintStruct struct {
	Table 	string
	TableID string
	KeyVal 	map[string]int
	Data   	map[string]interface{}
}

// NOTE: We assume that primary key is only one integer value!!!
func TransformRowToHint(appConfig *config.AppConfig, row map[string]interface{}, table string) *HintStruct {
	hint := HintStruct{}
	hint.Table = table
	intVal, err := strconv.Atoi(fmt.Sprint(row["id"]))
	if err != nil {
		log.Fatal(err)
	}
	hint.KeyVal = map[string]int{"id": intVal}
	hint.TableID = appConfig.TableNameIDPairs[table]
	hint.Data = row
	return &hint
}

func TransformDisplayFlagDataToHint(appConfig *config.AppConfig, data map[string]string) *HintStruct {
	hint := HintStruct{}
	intVal, err := strconv.Atoi(data["id"])
	if err != nil {
		log.Fatal(err)
	}
	hint.KeyVal = map[string]int{"id": intVal}
	hint.Table = appConfig.TableIDNamePairs[data["table_id"]]
	hint.TableID = data["table_id"]
	return &hint
}

func (hint *HintStruct) GetTagName(appConfig *config.AppConfig) (string, error) {
	for _, tag := range appConfig.Tags {
		for _, member := range tag.Members {
			if hint.Table == member {
				return tag.Name, nil
			}
		}
	}
	return "", errors.New("No Corresponding Tag Found!")
}

func (hint *HintStruct) GetMemberID(appConfig *config.AppConfig, tagName string) (string, error) {
	for _, tag := range appConfig.Tags {
		if tag.Name == tagName {
			for memberID, memberTable := range tag.Members {
				if memberTable == hint.Table {
					return memberID, nil
				}
			}
		}
	}
	return "", errors.New("No Corresponding Tag Found!")
}

func (hint *HintStruct) GetParentTags(appConfig *config.AppConfig) ([]string, error) {
	tag, err := hint.GetTagName(appConfig)
	if err != nil {
		return nil, err
	}

	var parentTags []string
	for _, dependency := range appConfig.Dependencies {
		if dependency.Tag == tag {
			for _, dependsOn := range dependency.DependsOn {
				// Use As as the tag name to avoid adding duplicate tag names
				if dependsOn.As != "" {
					parentTags = append(parentTags, dependsOn.As)
				} else {
					parentTags = append(parentTags, dependsOn.Tag)
				}
			}
		}
	}

	return parentTags, nil
}

func (hint *HintStruct) GetOriginalTagNameFromAliasOfParentTagIfExists(appConfig *config.AppConfig, alias string) (string, error) {
	tag, err := hint.GetTagName(appConfig)
	if err != nil {
		return "", err
	}

	for _, dependency := range appConfig.Dependencies {
		if dependency.Tag == tag {
			for _, dependsOn := range dependency.DependsOn {
				if dependsOn.As == alias {
					return dependsOn.Tag, nil
				}
			}
		}
	}

	return alias, errors.New("No Corresponding Tag for the Provided Alias Found!")
}

func (hint *HintStruct) GetDisplayExistenceSetting(appConfig *config.AppConfig, pTag string) (string, error) {
	tag, err := hint.GetTagName(appConfig)
	if err != nil {
		return "", err
	}

	for _, dependency := range appConfig.Dependencies {
		if dependency.Tag == tag {
			for _, dependsOn := range dependency.DependsOn {
				if dependsOn.As != "" {
					if dependsOn.As == pTag {
						return dependsOn.DisplayExistence, nil
					} else {
						continue
					}
				} else {
					if dependsOn.Tag == pTag {
						return dependsOn.DisplayExistence, nil
					}
				}
			}
		}
	}

	return "", errors.New("Find display existence error!")
}

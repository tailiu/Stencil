package app_display

import (
	"database/sql"
	"errors"
	"log"
	"stencil/config"
	"stencil/db"
	"strconv"
)

// The Key should be the primay key of the Table
type HintStruct struct {
	Table  string
	KeyVal map[string]int
}

// NOTE: We assume that primary key is only one integer value!!!
func TransformRowToHint(dbConn *sql.DB, row map[string]string, table string) (HintStruct, error) {
	hint := HintStruct{}
	pk, err := db.GetPrimaryKeyOfTable(dbConn, table)
	if err != nil {
		return hint, err
	} else {
		intPK, err1 := strconv.Atoi(row[pk])
		if err1 != nil {
			log.Fatal(err1)
		}
		keyVal := map[string]int{
			pk: intPK,
		}
		hint.Table = table
		hint.KeyVal = keyVal
	}
	return hint, nil
}

func (hint HintStruct) GetTagName(appConfig *config.AppConfig) (string, error) {
	for _, tag := range appConfig.Tags {
		for _, member := range tag.Members {
			if hint.Table == member {
				return tag.Name, nil
			}
		}
	}
	return "", errors.New("No Corresponding Tag Found!")
}

func (hint HintStruct) GetMemberID(appConfig *config.AppConfig, tagName string) (string, error) {
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

func (hint HintStruct) GetParentTags(appConfig *config.AppConfig) ([]string, error) {
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

func (hint HintStruct) GetOriginalTagNameFromAliasOfParentTagIfExists(appConfig *config.AppConfig, alias string) (string, error) {
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

func (hint HintStruct) GetDisplayExistenceSetting(appConfig *config.AppConfig, pTag string) (string, error) {
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

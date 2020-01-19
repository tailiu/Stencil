package SA2_display

import (
	"errors"
	"log"
	"stencil/config"
	"stencil/db"
	"strconv"
	"strings"
	"database/sql"
	"fmt"
)

type HintStruct struct {
	TableID		string
	TableName 	string
	RowIDs		[]int
	Data		map[string]interface{}
}

func TransformRowToHint(appConfig *config.AppConfig, data map[string]string) HintStruct {
	var rowIDs []int
	s := data["row_ids"][1:len(data["row_ids"]) - 1]
	s1 := strings.Split(s, ",")
	for _, strRowID := range s1 {
		rowID, err1 := strconv.Atoi(strRowID)
		if err1 != nil {
			log.Fatal(err1)
		} 
		rowIDs = append(rowIDs, rowID)
	}

	hint := HintStruct{}
	hint.TableID = data["table_id"]
	hint.TableName = appConfig.TableIDNamePairs[data["table_id"]]
	hint.RowIDs = rowIDs
	return hint
}

// NOTE: We assume that primary key is only one integer value!!!
func TransformRowToHint1(appConfig *config.AppConfig, data map[string]interface{}) HintStruct {
	hint := HintStruct{}
	hint.Data = data
	hint.RowIDs = GetRowIDsFromData(data)
	for key := range data {
		if strings.Contains(key, ".rowids_str") {
			hint.TableName = strings.Split(key, ".")[0]
			break
		}
	}
	hint.TableID = appConfig.TableNameIDPairs[hint.TableName]
	return hint
}

func (hint HintStruct) GetTagName(appConfig *config.AppConfig) (string, error) {
	for _, tag := range appConfig.Tags {
		for _, member := range tag.Members {
			if hint.TableName == member {
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
				if memberTable == hint.TableName {
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

func (hint HintStruct) GetCombinedDisplaySettings(appConfig *config.AppConfig) (string, error) {
	tag, err := hint.GetTagName(appConfig)
	if err != nil {
		return "", err
	}

	for _, dependency := range appConfig.Dependencies {
		if dependency.Tag == tag {
			if dependency.CombinedDisplaySetting == "" {
				return "", errors.New("No combined display settings found!")
			} else {
				return dependency.CombinedDisplaySetting, nil
			}
		}
	}

	return "", errors.New("No combined display settings found!")
}

func (hint HintStruct) GetRestrictionsInTag(appConfig *config.AppConfig) ([]map[string]string, error) {
	tagName, err := hint.GetTagName(appConfig)
	if err != nil {
		return nil, err
	}

	for _, tag := range appConfig.Tags {
		if tag.Name == tagName {
			return tag.Restrictions, nil
		}
	}

	return nil, errors.New("No matched tag found!")
}

func (hint HintStruct) GetAllRowIDs(stencilDBConn *sql.DB, appID string) []map[string]interface{} {
	query := fmt.Sprintf("select row_id from migration_table where app_id = %s and table_id = %s and group_id in (select group_id from migration_table where row_id = %d and table_id = %s and app_id = %s);",
		appID, hint.TableID, hint.RowIDs[0], hint.TableID, appID)
	
	// log.Println("++++++++")
	// log.Println(query)
	// log.Println("++++++++")

	result, err := db.DataCall(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	
	return result
}
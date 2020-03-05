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
	Tag			string
}

func TransformRowToHint(displayConfig *displayConfig, 
	data map[string]string) HintStruct {

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
	hint.TableName = displayConfig.tableIDNamePairs[data["table_id"]]

	hint.RowIDs = rowIDs

	hint.Tag = GetTagName1(displayConfig, hint.TableName)

	return &hint
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

func (hint HintStruct) GetTagName(displayConfig *displayConfig) (string, error) {
	
	for _, tag := range appConfig.Tags {
		for _, member := range tag.Members {
			if hint.TableName == member {
				return tag.Name, nil
			}
		}
	}
	return "", errors.New("No Corresponding Tag Found!")

}

func GetTagName1(displayConfig *displayConfig, 
	table string) (string, error) {

	for _, tag := range displayConfig.dstAppConfig.dag.Tags {

		for _, member := range tag.Members {

			if table == member {

				return tag.Name, nil

			}
		}

	}

	return "", errors.New("No Corresponding Tag Found!")

}

func (hint *HintStruct) GetTagDisplaySetting(
	displayConfig *displayConfig) (string, error) {
	
	for _, tag := range displayConfig.dstAppConfig.dag.Tags {

		if tag.Name == hint.Tag {

			if tag.Display_setting != "" {

				return tag.Display_setting, nil

			} else {

				return "default_display_setting", nil
			}
		}
	}

	return "", errors.New("Error: No Tag Found For the Provided TagName")

}

func (hint *HintStruct) GetMemberID(displayConfig *displayConfig) (string, error) {
	
	for _, tag := range displayConfig.dstAppConfig.dag.Tags {

		if tag.Name == hint.Tag {

			for memberID, memberTable := range tag.Members {

				if memberTable == hint.Table {

					return memberID, nil
					
				}
			}
		}
	}
	
	return "", errors.New("No Corresponding Tag Found!")

}

func (hint *HintStruct) GetDependsOnTables(displayConfig *displayConfig, 
	memberID string) []string {

	var dependsOnTables []string

	for _, tag := range displayConfig.dstAppConfig.dag.Tags {

		if tag.Name == hint.Tag {

			for _, innerDependency := range tag.InnerDependencies {

				for dependsOnMember, member := range innerDependency {

					if memberID == strings.Split(member, ".")[0] {

						table, _ := GetTableByMemberID(displayConfig.dstAppConfig.dag, 
							hint.Tag, strings.Split(dependsOnMember, ".")[0])

						dependsOnTables = append(dependsOnTables, table)

					}
				}
			}
		}
	}
	return dependsOnTables
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

func (hint HintStruct) GetOriginalTagNameFromAliasOfParentTagIfExists(appConfig *config.AppConfig, 
	alias string) (string, error) {
	
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

func (hint HintStruct) GetDisplayExistenceSetting(appConfig *config.AppConfig, 
	pTag string) (string, error) {
	
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

func (hint *HintStruct) GetRestrictionsInTag(
	displayConfig *displayConfig) ([]map[string]string, error) {

	for _, tag := range displayConfig.dstAppConfig.dag.Tags {
		if tag.Name == hint.Tag {
			return tag.Restrictions, nil
		}
	}

	return nil, errors.New("No matched tag found!")
}

func (hint HintStruct) GetAllRowIDs(stencilDBConn *sql.DB, 
	appID string) []map[string]interface{} {
	
	query := fmt.Sprintf(
		`select row_id from migration_table 
		where app_id = %s and table_id = %s and 
		group_id in 
			(select group_id from migration_table 
			where row_id = %d and table_id = %s and app_id = %s)`,
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

func (hint HintStruct) GetOwnershipSpec(dstDAG *DAG) (*config.Ownership, error) {

	for _, ownership := range dstDAG.Ownerships {

		if ownership.Tag == hint.Tag {

			return &ownership, nil
		}

	}

	return nil, errors.New("Error: No Ownership Tag Found For the Provided Data")

}
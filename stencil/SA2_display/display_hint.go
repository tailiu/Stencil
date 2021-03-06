package SA2_display

import (
	"log"
	"errors"
	"stencil/config"
	"stencil/db"
	"strconv"
	"strings"
	"fmt"
)

type HintStruct struct {
	TableID		string
	TableName 	string
	RowIDs		[]int
	Data		map[string]interface{}
	Tag			string
}

func TransformRowToHint(display *display, data map[string]string) *HintStruct {

	var rowIDs []int
	var err2 error

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
	hint.TableName = display.tableIDNamePairs[data["table_id"]]

	hint.RowIDs = rowIDs

	hint.Tag, err2 = GetTagName(display, hint.TableName)
	if err2 != nil {
		log.Fatal(err2)
	}

	return &hint

}

// NOTE: We assume that primary key is only one integer value!!!
func TransformRowToHint1(display *display, data map[string]interface{}) *HintStruct {
	
	hint := HintStruct{}

	for key := range data {
		if strings.Contains(key, ".rowids_str") {
			hint.TableName = strings.Split(key, ".")[0]
			break
		}
	}
	
	hint.Data = data
	hint.RowIDs = GetRowIDsFromData(data)
	hint.TableID = display.dstAppConfig.tableNameIDPairs[hint.TableName]

	var err1 error
	hint.Tag, err1 = GetTagName(display, hint.TableName)
	if err1 != nil {
		log.Fatal(err1)
	}

	return &hint

}

func GetTagName(display *display, table string) (string, error) {

	for _, tag := range display.dstAppConfig.dag.Tags {

		for _, member := range tag.Members {

			if table == member {

				return tag.Name, nil

			}
		}

	}

	return "", errors.New("No Corresponding Tag Found!")

}

func (hint *HintStruct) GetTagDisplaySetting(display *display) (string, error) {
	
	for _, tag := range display.dstAppConfig.dag.Tags {

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

func (hint *HintStruct) GetMemberID(display *display) (string, error) {
	
	for _, tag := range display.dstAppConfig.dag.Tags {

		if tag.Name == hint.Tag {

			for memberID, memberTable := range tag.Members {

				if memberTable == hint.TableName {

					return memberID, nil
					
				}
			}
		}
	}
	
	return "", errors.New("No Corresponding Tag Found!")

}

func (hint *HintStruct) GetDependsOnTables(display *display, 
	memberID string) []string {

	var dependsOnTables []string

	for _, tag := range display.dstAppConfig.dag.Tags {

		if tag.Name == hint.Tag {

			for _, innerDependency := range tag.InnerDependencies {

				for dependsOnMember, member := range innerDependency {

					if memberID == strings.Split(member, ".")[0] {

						table, _ := display.dstAppConfig.dag.GetTableByMemberID(
							hint.Tag, strings.Split(dependsOnMember, ".")[0])

						dependsOnTables = append(dependsOnTables, table)

					}
				}
			}
		}
	}
	return dependsOnTables
}

func (hint *HintStruct) GetParentTags(display *display) ([]string, error) {
	
	var parentTags []string
	
	for _, dependency := range display.dstAppConfig.dag.Dependencies {

		if dependency.Tag == hint.Tag {

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

func (hint *HintStruct) GetOriginalTagNameFromAliasOfParentTagIfExists(
	display *display, 
	alias string) (string, error) {
	
	for _, dependency := range display.dstAppConfig.dag.Dependencies {

		if dependency.Tag == hint.Tag {

			for _, dependsOn := range dependency.DependsOn {

				if dependsOn.As == alias {
					
					return dependsOn.Tag, nil

				}
			}
		}
	}

	return alias, errors.New("No Corresponding Tag for the Provided Alias Found!")
}

func (hint *HintStruct) GetDisplayExistenceSetting(display *display, 
	pTag string) (string, error) {
	
	for _, dependency := range display.dstAppConfig.dag.Dependencies {

		if dependency.Tag == hint.Tag {
			
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

func (hint *HintStruct) GetCombinedDisplaySettings(
	display *display) (string, error) {
	
	for _, dependency := range display.dstAppConfig.dag.Dependencies {
		if dependency.Tag == hint.Tag {
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
	display *display) ([]map[string]string, error) {

	for _, tag := range display.dstAppConfig.dag.Tags {
		if tag.Name == hint.Tag {
			return tag.Restrictions, nil
		}
	}

	return nil, errors.New("No matched tag found!")
}

func (hint *HintStruct) GetAllRowIDs(display *display) []map[string]interface{} {
	
	appID := display.dstAppConfig.appID

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

	result, err := db.DataCall(display.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}
	
	return result
}

func (hint *HintStruct) GetOwnershipSpec(display *display) (*config.Ownership, error) {

	for _, ownership := range display.dstAppConfig.dag.Ownerships {

		if ownership.Tag == hint.Tag {

			return &ownership, nil
		}

	}

	return nil, errors.New("Error: No Ownership Tag Found For the Provided Data")

}

func (hint *HintStruct) GetDisplaySettingInDependencies(display *display, 
	pTag string) (string, error) {
	
	setting, err := display.dstAppConfig.dag.GetDepDisplaySetting(hint.Tag, pTag)

	if err != nil {
		return "", err
	}

	if setting == "" {
		return "parent_node_complete_displays", nil
	} else {
		return setting, nil
	}
}
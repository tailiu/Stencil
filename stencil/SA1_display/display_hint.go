package SA1_display

import (
	"errors"
	"stencil/config"
	"stencil/reference_resolution"
	"strconv"
	"strings"
	"log"
	"fmt"
)

// The Key should be the primay key of the Table
type HintStruct struct {
	Table 							string
	TableID 						string
	KeyVal 							map[string]int
	Data   							map[string]interface{}
	Tag								string
}

func CreateHint(tableName, tableID, id string) *HintStruct {

	hint := &HintStruct{}

	intVal, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	hint.KeyVal = map[string]int{"id": intVal}

	hint.Table = tableName
	hint.TableID = tableID
	
	return hint

}

// NOTE: We assume that primary key is only one integer value!!!
func TransformRowToHint(displayConfig *displayConfig,
	row map[string]interface{}, table, tag string) *HintStruct {
	
	hint := HintStruct{}

	hint.Table = table
	
	intVal, err := strconv.Atoi(fmt.Sprint(row["id"]))
	if err != nil {
		log.Fatal(err)
	}
	hint.KeyVal = map[string]int{"id": intVal}

	hint.TableID = displayConfig.dstAppConfig.tableNameIDPairs[table]
	
	hint.Data = row

	hint.Tag = tag
	
	return &hint

}

func TransformDisplayFlagDataToHint(displayConfig *displayConfig,
	data map[string]string) *HintStruct {
	
	hint := HintStruct{}

	intVal, err := strconv.Atoi(data["id"])
	if err != nil {
		log.Fatal(err)
	}
	hint.KeyVal = map[string]int{"id": intVal}

	hint.Table = displayConfig.tableIDNamePairs[data["table_id"]]
	hint.TableID = data["table_id"]

	tag, err1 := getTagName(displayConfig, hint.Table)
	if err1 != nil {
		log.Fatal(err1)
	}
	hint.Tag = tag
	
	return &hint

}

// // By default, a member name is a table name 
// func (hint *HintStruct) GetMemberName(appConfig *config.AppConfig) (string, error) {
// 	for _, tag := range appConfig.Tags {

// 		for _, member := range tag.Members {
// 			if hint.Table == member {
// 				return member, nil
// 			}
// 		}

// 	}

// 	return "", errors.New("No Corresponding Member Found!")  
// }

func (hint *HintStruct) TransformHintToIdenity(
	displayConfig *displayConfig) *reference_resolution.Identity {

	return reference_resolution.CreateIdentity(displayConfig.dstAppConfig.appID, 
		hint.TableID, strconv.Itoa(hint.KeyVal["id"]))

}

func getTagName(displayConfig *displayConfig, table string) (string, error) {

	for _, tag := range displayConfig.dstAppConfig.dag.Tags {

		for _, member := range tag.Members {

			if table == member {

				return tag.Name, nil

			}
		}

	}

	return "", errors.New("No Corresponding Tag Found!")
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

func (hint *HintStruct) GetParentTags(displayConfig *displayConfig) ([]string, error) {

	var parentTags []string
	
	for _, dependency := range displayConfig.dstAppConfig.dag.Dependencies {

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
	displayConfig *displayConfig, alias string) (string, error) {

	for _, dependency := range displayConfig.dstAppConfig.dag.Dependencies {

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

func (hint *HintStruct) GetDisplayExistenceSetting(
	displayConfig *displayConfig, pTag string) (string, error) {

	for _, dependency := range displayConfig.dstAppConfig.dag.Dependencies {

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
	displayConfig *displayConfig) (string, error) {
	
	for _, dependency := range displayConfig.dstAppConfig.dag.Dependencies {

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

func (hint *HintStruct) GetDisplaySettingInDependencies(displayConfig *displayConfig, 
	pTag string) (string, error) {

	setting, err := common_funcs.GetDepDisplaySetting(displayConfig.dstAppConfig.dag, hint.Tag, pTag)

	if err != nil {
		return "", err
	}

	if setting == "" {

		return "parent_node_complete_displays", nil

	} else {

		return setting, nil

	}
	
}

func (hint *HintStruct) GetDisplaySettingInOwnership(
	displayConfig *displayConfig) (string, error) {

	for _, ownership := range displayConfig.dstAppConfig.dag.Ownerships {

		if ownership.Tag == hint.Tag {

			return ownership.Display_setting, nil
		}

	}

	return "", errors.New("Error: No Tag Found For the Provided TagName")

}

func (hint *HintStruct) GetOwnershipSpec(
	displayConfig *displayConfig) (*config.Ownership, error) {

	for _, ownership := range displayConfig.dstAppConfig.dag.Ownerships {

		if ownership.Tag == hint.Tag {

			return &ownership, nil
		}

	}

	return nil, errors.New("Error: No Ownership Tag Found For the Provided Data")

}

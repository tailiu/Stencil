package dependency_handler

import (
	"errors"
	"log"
	"fmt"
	"stencil/config"
	"stencil/display"
	"database/sql"
	// "strconv"
	"strings"
)

func getHintsInParentNode(stencilDBConn *sql.DB, appConfig *config.AppConfig, hints []display.HintStruct, conditions []string) (display.HintStruct, error) {	
	// log.Println(".....Second check......")
	// log.Println(display.GetData1FromPhysicalSchemaByRowID(stencilDBConn, appConfig.QR, "statuses.*", "statuses", "734487949"))
	// log.Println(display.GetData1FromPhysicalSchema(stencilDBConn, appConfig.QR, "statuses.*", "statuses", "statuses.conversation_id", "=", "1563"))
	// log.Println("...........")
	
	var data map[string]interface{}
	hintID := -1
	for i, condition := range conditions {
		log.Println(condition)
		tableAttr1 := strings.Split(condition, ":")[0]
		tableAttr2 := strings.Split(condition, ":")[1]
		t1 := strings.Split(tableAttr1, ".")[0]
		a1 := strings.Split(tableAttr1, ".")[1]
		t2 := strings.Split(tableAttr2, ".")[0]
		a2 := strings.Split(tableAttr2, ".")[1]
		// log.Println(t1, a1, t2, a2)
		if i == 0 {
			for j, hint := range hints {
				if hint.TableName == t1 {
					hintID = j
				}
			}
			if hintID == -1 {
				// In this case, since data may be incomplete, we cannot get the data in the parent node
				return display.HintStruct{}, errors.New("Fail To Get Any Data in the Parent Node")
			} else {
				// This can happen when the data this data depends on is not migrated,
				// e.g., a post does not have correpsonding conversation in Diaspora, so when it is migrated to Mastodon,
				// and becomes a status, it does not have conversation_id which is actually necessary for each
				// status in Mastodon.
				if hints[hintID].Data[t1 + "." + a1] == nil {
					return display.HintStruct{}, errors.New("Fail To Get Any Data in the Parent Node")
				}
				data = display.GetData1FromPhysicalSchema(stencilDBConn, appConfig.QR, appConfig.AppID, t2 + ".*", t2, t2 + "." + a2, "=", fmt.Sprint(hints[hintID].Data[t1 + "." + a1]))
				// log.Println(".....first check......")
				// log.Println(data)
				// log.Println("...........")
				if len(data) == 0 {
					return display.HintStruct{}, errors.New("Fail To Get Any Data in the Parent Node")
				}
			}
		} else {
			data = display.GetData1FromPhysicalSchema(stencilDBConn, appConfig.QR, appConfig.AppID, t2 + ".*", t2, t2 + "." + a2, "=", fmt.Sprint(data[t1 + "." + a1]) )
			if len(data) == 0 {
				return display.HintStruct{}, errors.New("Fail To Get Any Data in the Parent Node")
			}
		}
	}
	// log.Println("...........")
	// log.Println(data)
	// log.Println("...........")

	return display.TransformRowToHint1(appConfig, data), nil
}

func replaceKey(appConfig *config.AppConfig, tag string, key string) string {
	for _, tag1 := range appConfig.Tags {
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

func dataFromParentNodeExists(stencilDBConn *sql.DB, appConfig *config.AppConfig, hints []display.HintStruct, pTag string) (bool, error) {
	displayExistenceSetting, _ := hints[0].GetDisplayExistenceSetting(appConfig, pTag)

	// If display existence setting is not set, then we have to try to get data in the parent node in any case
	if displayExistenceSetting == "" {
		return true, nil
	} else {
		tag, _ := hints[0].GetTagName(appConfig)
		tableCol := replaceKey(appConfig, tag, displayExistenceSetting)
		table := strings.Split(tableCol, ".")[0]
		for _, hint := range hints {
			if hint.TableName == table {
				if hint.Data[tableCol] == nil {
					return false, errors.New("This Data Does not Depend on Any Data in the Parent Node")
				} else {
					return true, nil
				}
			}
		}

	}
	// In this case, since data may be incomplete, we cannot find the existence of the data in a parent node
	// This also implies that it cannot find any data in a parent node
	return false, errors.New("Fail To Get Any Data in the Parent Node")
}

// Note: this function may return multiple hints based on dependencies
func GetdataFromParentNode(stencilDBConn *sql.DB, appConfig *config.AppConfig, hints []display.HintStruct, pTag string) (display.HintStruct, error) {

	// Before getting data from a parent node, we check the existence of the data based on the cols of a child node
	if exists, err := dataFromParentNodeExists(stencilDBConn, appConfig, hints, pTag); !exists {
		return display.HintStruct{}, err
	}

	tag, _ := hints[0].GetTagName(appConfig)
	conditions, _ := appConfig.GetDependsOnConditions(tag, pTag)
	pTag, _ = hints[0].GetOriginalTagNameFromAliasOfParentTagIfExists(appConfig, pTag)

	var proConditions []string
	var from, to string

	// log.Println(conditions)

	if len(conditions) == 1 {
		condition := conditions[0]
		from = replaceKey(appConfig, tag, condition.TagAttr)
		to = replaceKey(appConfig, pTag, condition.DependsOnAttr)
		proConditions = append(proConditions, from+":"+to)
	} else {
		for i, condition := range conditions {
			if i == 0 {
				from = replaceKey(appConfig, tag, condition.TagAttr)
				to = replaceKey(appConfig, strings.Split(condition.DependsOnAttr, ".")[0], strings.Split(condition.DependsOnAttr, ".")[1])
			} else if i == len(conditions)-1 {
				from = replaceKey(appConfig, strings.Split(condition.TagAttr, ".")[0], strings.Split(condition.TagAttr, ".")[1])
				to = replaceKey(appConfig, pTag, condition.DependsOnAttr)
			}
			proConditions = append(proConditions, from+":"+to)
		}
	}

	return getHintsInParentNode(stencilDBConn, appConfig, hints, proConditions)
}

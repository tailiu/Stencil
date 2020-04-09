package SA2_display

import (
	"log"
	"fmt"
	"stencil/config"
	"stencil/common_funcs"
	"strings"
)

func getHintInParentNode(displayConfig *displayConfig, 
	hints []*HintStruct, conditions []string, 
	pTag string) (*HintStruct, error) {
	
	// log.Println(".....Second check......")
	// log.Println(GetData1FromPhysicalSchemaByRowID(
		//stencilDBConn, appConfig.QR, "statuses.*", "statuses", "734487949"))
	// log.Println(GetData1FromPhysicalSchema(stencilDBConn, 
		// appConfig.QR, "statuses.*", "statuses", "statuses.conversation_id", "=", "1563"))
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

				// In this case, since data may be incomplete, 
				// we cannot get the data in the parent node
				return nil, common_funcs.CannotFindAnyDataInParent
			
			} else {

				// This can happen when the data this data depends on is not migrated,
				// e.g., a post does not have correpsonding conversation in Diaspora, 
				// so when it is migrated to Mastodon,
				// and becomes a status, it does not have conversation_id 
				// which is actually necessary for each
				// status in Mastodon.
				if hints[hintID].Data[t1 + "." + a1] == nil {
					return nil, common_funcs.CannotFindAnyDataInParent
				}

				data = GetData1FromPhysicalSchema(
					displayConfig,
					t2 + ".*", t2, t2 + "." + a2, "=", 
					fmt.Sprint(hints[hintID].Data[t1 + "." + a1]),
				)
				
				// log.Println(".....first check......")
				// log.Println(data)
				// log.Println("...........")

				if len(data) == 0 {
					return nil, common_funcs.CannotFindAnyDataInParent
				}
			}
		} else {

			data = GetData1FromPhysicalSchema(
				displayConfig,
				t2 + ".*", t2, t2 + "." + a2, "=", 
				fmt.Sprint(data[t1 + "." + a1]),
			)

			if len(data) == 0 {
				return nil, common_funcs.CannotFindAnyDataInParent
			}

		}

	}

	// log.Println("...........")
	// log.Println(data)
	// log.Println("...........")

	return TransformRowToHint1(displayConfig, data), nil

}

func dataFromParentNodeExists(displayConfig *displayConfig,
	hints []*HintStruct, pTag string) (bool, error) {

	displayExistenceSetting, _ := hints[0].GetDisplayExistenceSetting(displayConfig, pTag)

	// If display existence setting is not set, 
	// then we have to try to get data in the parent node in any case
	if displayExistenceSetting == "" {

		return true, nil

	} else {

		tableCol := displayConfig.dstAppConfig.dag.ReplaceKey(hints[0].Tag, displayExistenceSetting)
		table := strings.Split(tableCol, ".")[0]

		for _, hint := range hints {

			if hint.TableName == table {

				if hint.Data[tableCol] == nil {

					return false, common_funcs.NotDependsOnAnyData				
				} else {

					return true, nil
				}
			}
		}

	}

	// In this case, since data may be incomplete, 
	// we cannot find the existence of the data in a parent node
	// This also implies that it cannot find any data in a parent node
	return false, common_funcs.CannotFindAnyDataInParent

}

func getProcConditions(displayConfig *displayConfig, 
	tag, pTag string, conditions []config.DCondition) []string {

	var proConditions []string
	var from, to string

	// log.Println(conditions)

	if len(conditions) == 1 {

		condition := conditions[0]
		from = displayConfig.dstAppConfig.dag.ReplaceKey(tag, condition.TagAttr)
		to = displayConfig.dstAppConfig.dag.ReplaceKey(pTag, condition.DependsOnAttr)
		proConditions = append(proConditions, from+":"+to)

	} else {

		for i, condition := range conditions {

			if i == 0 {

				from = displayConfig.dstAppConfig.dag.ReplaceKey(
					tag, condition.TagAttr)

				to = displayConfig.dstAppConfig.dag.ReplaceKey(
					strings.Split(condition.DependsOnAttr, ".")[0], 
					strings.Split(condition.DependsOnAttr, ".")[1],
				)

			} else if i == len(conditions)-1 {

				from = displayConfig.dstAppConfig.dag.ReplaceKey(
					strings.Split(condition.TagAttr, ".")[0], 
					strings.Split(condition.TagAttr, ".")[1],
				)

				to = displayConfig.dstAppConfig.dag.ReplaceKey(
					pTag, condition.DependsOnAttr,
				)

			}

			proConditions = append(proConditions, from+":"+to)

		}
	}

	return proConditions

}


// Note: this function may return multiple hints based on dependencies
func GetdataFromParentNode(displayConfig *displayConfig,
	hints []*HintStruct, pTag string) (*HintStruct, error) {

	// Before getting data from a parent node, 
	// we check the existence of the data based on the cols of a child node
	if exists, err := dataFromParentNodeExists(displayConfig, hints, pTag); !exists {
		return nil, err
	}

	tag := hints[0].Tag
	pTag, _ = hints[0].GetOriginalTagNameFromAliasOfParentTagIfExists(displayConfig, pTag)

	conditions, _ := displayConfig.dstAppConfig.dag.GetDependsOnConditionsInDeps(tag, pTag)

	procConditions := getProcConditions(displayConfig, tag, pTag, conditions)

	return getHintInParentNode(displayConfig, hints, procConditions, pTag)

}

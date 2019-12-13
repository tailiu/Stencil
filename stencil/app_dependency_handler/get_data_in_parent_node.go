package app_dependency_handler

import (
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/app_display"
	"strconv"
	"strings"
)

func getHintsInParentNode(displayConfig *config.DisplayConfig, 
	hints []*app_display.HintStruct, conditions []string) (*app_display.HintStruct, error) {
	
	query := fmt.Sprintf("SELECT %s.* FROM ", "t"+strconv.Itoa(len(conditions)))
	from := ""
	table := ""
	hintID := -1

	for i, condition := range conditions {

		tableAttr1 := strings.Split(condition, ":")[0]
		tableAttr2 := strings.Split(condition, ":")[1]

		t1 := strings.Split(tableAttr1, ".")[0]
		a1 := strings.Split(tableAttr1, ".")[1]

		t2 := strings.Split(tableAttr2, ".")[0]
		a2 := strings.Split(tableAttr2, ".")[1]

		seq1 := "t" + strconv.Itoa(i)
		seq2 := "t" + strconv.Itoa(i+1)

		if i == 0 {

			// There could be mutliple pieces of data in nodes
			// For example:
			// A statuses node contains status, conversation, and status_stats
			for j, hint := range hints {

				if hint.Table == t1 {
					hintID = j
				}

			}

			// In this case, since data may be incomplete, 
			// we cannot get the data in the parent node
			if hintID == -1 {

				return nil, app_display.CannotFindAnyDataInParent

			} else {
				
				// For example:
				// if a condition is [favourites.status_id:statuses.id], 
				// from will be "favourites t1 JOIN statuses t2 ON t1.status_id = t2.id"
				from += fmt.Sprintf("%s %s JOIN %s %s ON %s.%s = %s.%s ",
					t1, seq1, t2, seq2, seq1, a1, seq2, a2)

			}

		// This is mainly to solve the case in which
		// conversation cannot directly depend on root
		// conversation depends on statuses, which in turn depends on root (obsolete now)
		} else {

			from += fmt.Sprintf("JOIN %s %s on %s.%s = %s.%s ",
				t2, seq2, seq1, a1, seq2, a2)

		}

		//The last condition
		if i == len(conditions)-1 {

			var depDataKey string
			var depDataValue int

			for k, v := range hints[hintID].KeyVal {

				depDataKey = k
				depDataValue = v

			}

			where := fmt.Sprintf("WHERE %s.%s = %d", "t0", depDataKey, depDataValue)
			table = t2
			query += from + where

		}
	}
	// fmt.Println(query)


	data, err := db.DataCall1(displayConfig.AppConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)
	if len(data) == 0 {

		return nil, app_display.CannotFindAnyDataInParent

	} else {

		return app_display.TransformRowToHint(displayConfig, data, table), nil

	}
}

func replaceKey(displayConfig *config.DisplayConfig, tag string, key string) string {

	for _, tag1 := range displayConfig.AppConfig.Tags {

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

func dataFromParentNodeExists(displayConfig *config.DisplayConfig,
	hints []*app_display.HintStruct, pTag string) (bool, error) {
	
	displayExistenceSetting, _ := hints[0].GetDisplayExistenceSetting(displayConfig, pTag)

	// If display existence setting is not set, 
	// then we have to try to get data in the parent node in any case
	if displayExistenceSetting == "" {

		return true, nil

	} else {

		tag, _ := hints[0].GetTagName(displayConfig)
		tableCol := replaceKey(displayConfig, tag, displayExistenceSetting)
		table := strings.Split(tableCol, ".")[0]

		for _, hint := range hints {

			if hint.Table == table {

				if hint.Data[tableCol] == nil {

					return false, app_display.NotDependsOnAnyData

				} else {

					return true, nil

				}
			}
		}

	}

	// In this case, since data may be incomplete, 
	// we cannot find the existence of the data in a parent node
	// This also implies that it cannot find any data in a parent node
	return false, app_display.CannotFindAnyDataInParent

}

// Note: this function may return multiple hints based on dependencies
func GetdataFromParentNode(displayConfig *config.DisplayConfig,
	hints []*app_display.HintStruct, pTag string) (*app_display.HintStruct, error) {

	// Before getting data from a parent node, 
	// we check the existence of the data based on the cols of a child node
	if exists, err := dataFromParentNodeExists(displayConfig, hints, pTag); !exists {
		return nil, err
	}

	tag, _ := hints[0].GetTagName(displayConfig)
	conditions, _ := displayConfig.AppConfig.GetDependsOnConditions(tag, pTag)
	pTag, _ = hints[0].GetOriginalTagNameFromAliasOfParentTagIfExists(displayConfig, pTag)

	var procConditions []string
	var from, to string

	if len(conditions) == 1 {

		condition := conditions[0]
		from = replaceKey(displayConfig, tag, condition.TagAttr)
		to = replaceKey(displayConfig, pTag, condition.DependsOnAttr)
		procConditions = append(procConditions, from+":"+to)

	} else {

		for i, condition := range conditions {

			if i == 0 {

				from = replaceKey(displayConfig, tag, condition.TagAttr)

				to = replaceKey(displayConfig,
					strings.Split(condition.DependsOnAttr, ".")[0], 
					strings.Split(condition.DependsOnAttr, ".")[1])

			} else if i == len(conditions)-1 {

				from = replaceKey(displayConfig, 
					strings.Split(condition.TagAttr, ".")[0], 
					strings.Split(condition.TagAttr, ".")[1])
				
				to = replaceKey(displayConfig, pTag, condition.DependsOnAttr)

			}

			procConditions = append(procConditions, from+":"+to)

		}

	}

	fmt.Println(procConditions)
	fmt.Println(*hints)

	return getHintsInParentNode(displayConfig, hints, procConditions)

}

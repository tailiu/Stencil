package dependency_handler

import (
	"transaction/config"
	"transaction/display"
	"transaction/db"
	"fmt"
	"strings"
	"strconv"
	"errors"
)

func getHintsInParentNode(appConfig *config.AppConfig, hint display.HintStruct, conditions []string) ([]display.HintStruct, error) {
	query := fmt.Sprintf("SELECT %s.* FROM ", "t"+strconv.Itoa(len(conditions)))
	from := ""
	table := ""
	for i, condition := range conditions {
		table1 := strings.Split(condition, ":")[0]
		table2 := strings.Split(condition, ":")[1]
		t1 := strings.Split(table1, ".")[0]
		a1 := strings.Split(table1, ".")[1]
		t2 := strings.Split(table2, ".")[0]
		a2 := strings.Split(table2, ".")[1]
		seq1 := "t" + strconv.Itoa(i)
		seq2 := "t" + strconv.Itoa(i+1)
		if i == 0 {
			from += fmt.Sprintf("%s %s JOIN %s %s ON %s.%s = %s.%s ",
				t1, seq1, t2, seq2, seq1, a1, seq2, a2)
		} else {
			from += fmt.Sprintf("JOIN %s %s on %s.%s = %s.%s ",
				t2, seq2, seq1, a1, seq2, a2)
		}
		if i == len(conditions)-1 {
			var depDataKey string
			var depDataValue int
			for k, v := range hint.KeyVal {
				depDataKey = k
				depDataValue = v
			}
			where := fmt.Sprintf("WHERE %s.%s = %d;", "t0", depDataKey, depDataValue)
			table = t2
			query += from + where
		}
	}
	fmt.Println(query)

	data := db.GetAllColsOfRows(appConfig.DBConn, query)
	
	if len(data) == 0 {
		return nil, errors.New("Error In Get Data: Fail To Get Any Data This Data Depends On")
	} else {
		var result []display.HintStruct
		for _, oneData := range data {
			oneHint, err := display.TransformRowToHint(appConfig.DBConn, oneData, table)
			if err != nil {
				return nil, err
			} else {
				result = append(result, oneHint)
			}
		}
		return result, nil
	}
}

func replaceKey(innerDependencies []config.Tag, tag string, key string) string {
	for _, dependsOnTag := range innerDependencies {
		if dependsOnTag.Name == tag {
			// fmt.Println(dependsOnTag)
			for k, v := range dependsOnTag.Keys {
				if k == key {
					member := strings.Split(v, ".")[0]
					attr := strings.Split(v, ".")[1]
					for k1, table := range dependsOnTag.Members {
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

// Note: this function may return multiple hints based on dependencies
func GetdataFromParentNode(appConfig *config.AppConfig, hint display.HintStruct, pTag string) ([]display.HintStruct, error) {

	tag, _ := hint.GetTagName(appConfig)
	conditions, _ := appConfig.GetDependsOnConditions(tag, pTag)
	pTag, _ = hint.GetOriginalTagNameFromAliasOfParentTagIfExists(appConfig, pTag)

	var proConditions []string
	var from, to string

	if len(conditions) == 1 {
		condition := conditions[0]
		from = replaceKey(appConfig.Tags, tag, condition.TagAttr)
		to = replaceKey(appConfig.Tags, pTag, condition.DependsOnAttr)
		proConditions = append(proConditions, from + ":" + to)
	} else {
		for i, condition := range(conditions) {
			if i == 0 {
				from = replaceKey(appConfig.Tags, tag, condition.TagAttr)
				to = replaceKey(appConfig.Tags, strings.Split(condition.DependsOnAttr, ".")[0], strings.Split(condition.DependsOnAttr, ".")[1])
			} else if i == len(conditions) - 1 {
				from = replaceKey(appConfig.Tags, strings.Split(condition.TagAttr, ".")[0], strings.Split(condition.TagAttr, ".")[1])
				to = replaceKey(appConfig.Tags, pTag, condition.DependsOnAttr)
			} 
			proConditions = append(proConditions, from + ":" + to)
		}
	}

	// fmt.Println(proConditions)
	// fmt.Println(hint)

	return getHintsInParentNode(appConfig, hint, proConditions)
}

func CheckDisplayCondition() bool {
	return false
}
package dependency_handler

import (
	"transaction/config"
	"transaction/display"
	"transaction/db"
	"fmt"
	"strings"
	"strconv"
	"errors"
	"log"
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
		return nil, errors.New("Fail To Get Any Data in the Parent Node")
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

func dataFromParentNodeExists(appConfig *config.AppConfig, hint display.HintStruct, pTag string) bool {
	displayExistenceSetting, _ := hint.GetDisplayExistenceSetting(appConfig, pTag)

	if displayExistenceSetting == "" {
		return true
	} else {
		fmt.Println(displayExistenceSetting)
		tag, _ := hint.GetTagName(appConfig)
		col := strings.Split(replaceKey(appConfig, tag, displayExistenceSetting), ".")[1]
		var dataKey string
		var dataValue int
		for k, v := range hint.KeyVal {
			dataKey = k
			dataValue = v
		}
		query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = %d;", col, hint.Table, dataKey, dataValue)
		fmt.Println(query)
		data := db.GetAllColsOfRows(appConfig.DBConn, query)
		if len(data) == 0 {
			log.Fatal("Data is missing??")
		} else {
			if data[0][col] == "NULL" {
				return false
			} else {
				return true
			}
		}
	}

	panic("Should never happen")
}

// Note: this function may return multiple hints based on dependencies
func GetdataFromParentNode(appConfig *config.AppConfig, hint display.HintStruct, pTag string) ([]display.HintStruct, error) {

	if !dataFromParentNodeExists(appConfig, hint, pTag) {
		return nil, fmt.Errorf("This Data Does not Depend on Any Data in the Parent Node %s", pTag)
	}

	tag, _ := hint.GetTagName(appConfig)
	conditions, _ := appConfig.GetDependsOnConditions(tag, pTag)
	pTag, _ = hint.GetOriginalTagNameFromAliasOfParentTagIfExists(appConfig, pTag)

	var proConditions []string
	var from, to string

	if len(conditions) == 1 {
		condition := conditions[0]
		from = replaceKey(appConfig, tag, condition.TagAttr)
		to = replaceKey(appConfig, pTag, condition.DependsOnAttr)
		proConditions = append(proConditions, from + ":" + to)
	} else {
		for i, condition := range(conditions) {
			if i == 0 {
				from = replaceKey(appConfig, tag, condition.TagAttr)
				to = replaceKey(appConfig, strings.Split(condition.DependsOnAttr, ".")[0], strings.Split(condition.DependsOnAttr, ".")[1])
			} else if i == len(conditions) - 1 {
				from = replaceKey(appConfig, strings.Split(condition.TagAttr, ".")[0], strings.Split(condition.TagAttr, ".")[1])
				to = replaceKey(appConfig, pTag, condition.DependsOnAttr)
			} 
			proConditions = append(proConditions, from + ":" + to)
		}
	}

	// fmt.Println(proConditions)
	// fmt.Println(hint)

	return getHintsInParentNode(appConfig, hint, proConditions)
}

// func CheckDisplayCondition() bool {
// 	return false
// }
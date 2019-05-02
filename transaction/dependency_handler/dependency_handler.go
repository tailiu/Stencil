package dependency_handler

import (
	"transaction/config"
	"transaction/display"
	"transaction/db"
	"fmt"
	"strings"
	"log"
	"database/sql"
	"strconv"
	"errors"
	"mastodon/auxiliary"
)

func getOneRowInParentNodeRandomly(dbConn *sql.DB, hint display.HintStruct, conditions []string) (map[string]string, string, error) {
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
			where := fmt.Sprintf("WHERE %s.%s = %d ORDER BY RANDOM() LIMIT 1;", "t0", depDataKey, depDataValue)
			table = t2
			query += from + where
		}
	}
	// fmt.Println(query)

	data := db.GetAllColsOfRows(dbConn, query)
	if len(data) == 0 {
		return nil, "", errors.New("Error In Get Data: Fail To Get One Data This Data Depends On")
	} else {
		return data[0], table, nil
	}
}

func getDependsOn(dependencies []config.Dependency, tag string) ([]config.DependsOn, error) {
	for _, dependency := range dependencies {
		// fmt.Println(dependency)
		if dependency.Tag == tag {
			return dependency.DependsOn, nil
		}
	}
	return nil, errors.New("Cannot Find Any Parent Tags")
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

func GetParentTags(appConfig config.AppConfig, data display.HintStruct) ([]string, error) {
	tag, err := data.GetTagName(appConfig.Tags)
	if err != nil {
		return nil, err
	}
	if tag == "root" {
		return nil, nil
	}

	var parentTags []string
	for _, dependency := range appConfig.Dependencies {
		if dependency.Tag == tag {
			for _, dependsOn := range dependency.DependsOn {
				parentTags = append(parentTags, dependsOn.Tag)
			}
		}
	}

	if len(parentTags) == 0 {
		return nil, errors.New("Check Parent Tag Name error: Does Not Find Any Parent Node!")
	} else {
		return parentTags, nil
	}
}

func GetOneDataFromParentNodeRandomly(dbConn *sql.DB, appConfig config.AppConfig, hint display.HintStruct, app string) (display.HintStruct, error){
	hintData := display.HintStruct{}
	data1 := DataInDependencyNode{}

	tag, err := hint.GetTagName(appConfig.Tags)
	if err != nil {
		log.Fatal(err)
	} 
	
	dependsOn, err1 := getDependsOn(appConfig.Dependencies, tag)
	if err1 != nil {
		log.Println(err1)
		return hintData, err1
	}

	// fmt.Println("all depends on ", dependsOn)
	for i := 0; i < getOneDataFromParentNodeAttemptTimes; i ++ {
		oneDependensOn := dependsOn[auxiliary.RandomNonnegativeIntWithUpperBound(len(dependsOn))]
		// fmt.Println("depends on ", oneDependensOn)

		var conditions []string
		var from, to string
		if len(oneDependensOn.Conditions) == 1 {
			condition := oneDependensOn.Conditions[0]
			from = replaceKey(appConfig.Tags, tag, condition.TagAttr)
			to = replaceKey(appConfig.Tags, oneDependensOn.Tag, condition.DependsOnAttr)
			conditions = append(conditions, from + ":" + to)
		} else {
			for i, condition := range(oneDependensOn.Conditions) {
				if i == 0 {
					from = replaceKey(appConfig.Tags, tag, condition.TagAttr)
					to = replaceKey(appConfig.Tags, strings.Split(condition.DependsOnAttr, ".")[0], strings.Split(condition.DependsOnAttr, ".")[1])
				} else if i == len(oneDependensOn.Conditions) - 1 {
					from = replaceKey(appConfig.Tags, strings.Split(condition.TagAttr, ".")[0], strings.Split(condition.TagAttr, ".")[1])
					to = replaceKey(appConfig.Tags, oneDependensOn.Tag, condition.DependsOnAttr)
				} 
				conditions = append(conditions, from + ":" + to)
			}
		}

		// fmt.Println(conditions)
		// fmt.Println(hint)

		data1.Data, data1.Table, err1 = getOneRowInParentNodeRandomly(dbConn, hint, conditions)
		if err1 != nil {
			fmt.Println(err1)
		} else {
			// fmt.Println(data1)
			break
		}
	}

	hintData, err2 := display.TransformRowToHint(dbConn, data1.Data, data1.Table)
	if err2 != nil {
		log.Fatal(err2)
	} 
	// fmt.Println(hintData)
	return hintData, nil
}

func CheckDisplayCondition() bool {
	return false
}
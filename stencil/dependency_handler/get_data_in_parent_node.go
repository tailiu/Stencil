package dependency_handler

import (
	"errors"
	"fmt"
	"log"
	"stencil/config"
	"stencil/display"
	"database/sql"
	"strconv"
	"strings"
	// "stencil/qr"
	// "stencil/db"
)

// func GetDataFromPhysicalSchemaByJoin(stencilDBConn *sql.DB, QR *qr.QR, cols, from, col, op, val, limit string) []map[string]string {	
// 	qs := qr.CreateQS(QR)
// 	qs.ColSimple(cols)
// 	qs.FromSimple(from)
// 	qs.WhereSimpleVal(col, op, val)
// 	qs.LimitResult(limit)

// 	physicalQuery := qs.GenSQL()
// 	// log.Println(physicalQuery)

// 	return db.GetAllColsOfRows(stencilDBConn, physicalQuery)
// }

func getHintsInParentNode(stencilDBConn *sql.DB, appConfig *config.AppConfig, hints []display.HintStruct, conditions []string) ([]display.HintStruct, error) {	
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
			for j, hint := range hints {
				if hint.Table == t1 {
					hintID = j
				}
			}
			if hintID == -1 {
				// In this case, since data may be incomplete, we cannot get the data in the parent node
				return nil, errors.New("Fail To Get Any Data in the Parent Node")
			} else {
				from += fmt.Sprintf("%s %s JOIN %s %s ON %s.%s = %s.%s ",
					t1, seq1, t2, seq2, seq1, a1, seq2, a2)
			}
		} else {
			from += fmt.Sprintf("JOIN %s %s on %s.%s = %s.%s ",
				t2, seq2, seq1, a1, seq2, a2)
		}
		if i == len(conditions)-1 {
			var depDataKey string
			var depDataValue int
			for k, v := range hints[hintID].KeyVal {
				depDataKey = k
				depDataValue = v
			}
			where := fmt.Sprintf("WHERE %s.%s = %d;", "t0", depDataKey, depDataValue)
			table = t2
			query += from + where
		}
	}

	// for i, condition := range conditions {
	// 	tableAttr1 := strings.Split(condition, ":")[0]
	// 	tableAttr2 := strings.Split(condition, ":")[1]
	// 	t1 := strings.Split(tableAttr1, ".")[0]
	// 	a1 := strings.Split(tableAttr1, ".")[1]
	// 	t2 := strings.Split(tableAttr2, ".")[0]
	// 	a2 := strings.Split(tableAttr2, ".")[1]
	// 	if i == 0 {
	// 		for j, hint := range hints {
	// 			if hint.Table == t1 {
	// 				hintID = j
	// 			}
	// 		}
	// 		if hintID == -1 {
	// 			// In this case, since data may be incomplete, we cannot get the data in the parent node
	// 			return nil, errors.New("Fail To Get Any Data in the Parent Node")
	// 		} else {
	// 			from += fmt.Sprintf("%s %s JOIN %s %s ON %s.%s = %s.%s ",
	// 				t1, seq1, t2, seq2, seq1, a1, seq2, a2)
	// 		}
	// 	}
	// 	display.GetDataFromPhysicalSchema(stencilDBConn, appConfig.QR, cols, from, col, op, , "1")

	// }

	// log.Println(hints)
	// log.Println(conditions)
	// log.Println(query)

	// Need to be changed
	// condition: [favourites.status_id:statuses.id]
	// SELECT t1.* FROM statuses t0 JOIN conversations t1 ON t0.conversation_id = t1.id WHERE t0.id = 34647260;
	// data := GetDataFromPhysicalSchemaByJoin(stencilDBConn, appConfig.QR, )
	var data []map[string]string

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

func dataFromParentNodeExists(stencilDBConn *sql.DB, appConfig *config.AppConfig, hints []display.HintStruct, pTag string) (bool, error) {
	displayExistenceSetting, _ := hints[0].GetDisplayExistenceSetting(appConfig, pTag)

	// If display existence setting is not set, then we have to try to get data in the parent node in any case
	if displayExistenceSetting == "" {
		return true, nil
	} else {
		// fmt.Println(displayExistenceSetting)
		tag, _ := hints[0].GetTagName(appConfig)
		tableCol := replaceKey(appConfig, tag, displayExistenceSetting)
		table := strings.Split(tableCol, ".")[0]
		col := strings.Split(tableCol, ".")[1]
		for _, hint := range hints {
			hT := hint.Table
			if hT == table {
				var dataKey string
				var dataValue int
				for k, v := range hint.KeyVal {
					dataKey = k
					dataValue = v
				}
				// query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = %d;", col, hint.Table, dataKey, dataValue)
				// Need to be changed
				// SELECT in_reply_to_id FROM statuses WHERE id = 15996494;
				log.Println("-------------------------")
				data := display.GetDataFromPhysicalSchema(stencilDBConn, appConfig.QR, 
					hT + "." + col, hT, hT + "." + dataKey, "=", strconv.Itoa(dataValue), "1")
				log.Println(data)
				log.Println("-------------------------")
				if len(data) == 0 {
					log.Fatal("Data is missing??")
				} else {
					if data[0][col] == "NULL" {
						return false, errors.New("This Data Does not Depend on Any Data in the Parent Node")
					} else {
						return true, nil
					}
				}
			}
		}

	}
	// In this case, since data may be incomplete, we cannot find the existence of the data in a parent node
	// This also implies that it cannot find any data in a parent node
	return false, errors.New("Fail To Get Any Data in the Parent Node")
}

// Note: this function may return multiple hints based on dependencies
func GetdataFromParentNode(stencilDBConn *sql.DB, appConfig *config.AppConfig, hints []display.HintStruct, pTag string) ([]display.HintStruct, error) {

	// Before getting data from a parent node, we check the existence of the data based on the cols of a child node
	if exists, err := dataFromParentNodeExists(stencilDBConn, appConfig, hints, pTag); !exists {
		return nil, err
	}

	tag, _ := hints[0].GetTagName(appConfig)
	conditions, _ := appConfig.GetDependsOnConditions(tag, pTag)
	pTag, _ = hints[0].GetOriginalTagNameFromAliasOfParentTagIfExists(appConfig, pTag)

	var proConditions []string
	var from, to string

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

	// fmt.Println(proConditions)
	// fmt.Println(hints)

	return getHintsInParentNode(stencilDBConn, appConfig, hints, proConditions)
}

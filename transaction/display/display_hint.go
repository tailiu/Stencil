package display

import (
	"log"
	"database/sql"
	"transaction/db"
	"strconv"
	"transaction/config"
	"errors"
)

// The Key should be the primay key of the Table
// type HintStruct struct {
// 	Table string		
// 	id string			
// 	Value string		
// 	ValueType string	
// }

// The Key should be the primay key of the Table
type HintStruct struct {
	Table string		
	KeyVal map[string]int
}

// NOTE: We assume that primary key is only one integer value!!!
func TransformRowToHint(dbConn *sql.DB, row map[string]string, table string) (HintStruct, error) {
	hint := HintStruct{}
	pk, err := db.GetPrimaryKeyOfTable(dbConn, table)
	if err != nil {
		return hint, err
	} else {
		intPK, err1 := strconv.Atoi(row[pk])
		if err1 != nil {
			log.Fatal(err1)
		}
		keyVal := map[string]int {
			pk:	intPK,
		}
		hint.Table = table
		hint.KeyVal = keyVal
	}
	return hint, nil
}

func (hint HintStruct) GetTagName(tags []config.Tag) (string, error) {
	for _, tag := range tags {
		for _, member := range tag.Members {
			if hint.Table == member {
				return tag.Name, nil
			}
		}
	}
	return "", errors.New("No Corresponding Tag Found!")
}

func (hint HintStruct) GetMemberID(appConfig *config.AppConfig, tagName string) (string, error) {
	for _, tag := range appConfig.Tags {
		if tag.Name == tagName {
			for memberID, memberTable := range tag.Members {
				if memberTable == hint.Table {
					return memberID, nil
				}
			}
		}
	}
	return "", errors.New("No Corresponding Tag Found!")
}

func (hint HintStruct) GetParentTags(appConfig config.AppConfig) ([]string, error) {
	tag, err := hint.GetTagName(appConfig.Tags)
	if err != nil {
		return nil, err
	}

	var parentTags []string
	for _, dependency := range appConfig.Dependencies {
		if dependency.Tag == tag {
			for _, dependsOn := range dependency.DependsOn {
				parentTags = append(parentTags, dependsOn.Tag)
			}
		}
	}
	
	return parentTags, nil
}
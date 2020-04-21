package config

import (
	"database/sql"
	"errors"
	"fmt"
	"stencil/db"
	"stencil/helper"
	"strings"
)

func (self *Tag) GetTagMembers() []string {

	var tagMembers []string

	for _, member := range self.Members {
		tagMembers = append(tagMembers, member)
	}

	return tagMembers
}

func (tag Tag) ResolveTagAttr(attr string) (string, error) {

	if _, ok := tag.Keys[attr]; ok {
		keyItems := strings.Split(tag.Keys[attr], ".")
		if _, ok := tag.Members[keyItems[0]]; ok {
			Table := tag.Members[keyItems[0]]
			Col := keyItems[1]
			return fmt.Sprintf("%s.%s", Table, Col), nil
		} else {
			return "", errors.New("Tag Not Resolved, Member Not Found: " + attr)
		}
	}
	return "", errors.New(fmt.Sprintf("Tag Not Resolved, Attr Not Found in Tag Keys: [%s].[%s] ;", tag.Name, attr))
}

func (self Tag) GetInDepMembersInOrder() []string {

	var members []string

	for _, inDep := range self.InnerDependencies {
		for mapFrom, mapTo := range inDep {
			mapFromTokens := strings.Split(mapFrom, ".")
			mapToTokens := strings.Split(mapTo, ".")
			if !helper.Contains(members, self.Members[mapFromTokens[0]]) {
				members = append(members, self.Members[mapFromTokens[0]])
			}
			if !helper.Contains(members, self.Members[mapToTokens[0]]) {
				members = append(members, self.Members[mapToTokens[0]])
			}
		}
	}

	return members
}

func (self Tag) GetInDepPairsInOrder() [][]string {

	var members [][]string

	for _, inDep := range self.InnerDependencies {
		for mapFrom, mapTo := range inDep {
			mapFromTokens := strings.Split(mapFrom, ".")
			mapToTokens := strings.Split(mapTo, ".")

			members = append(members, []string{self.Members[mapFromTokens[0]], self.Members[mapToTokens[0]]})

		}
	}

	return members
}

func (self Tag) CreateInDepMapSA2() map[string]map[string][]string {

	joinMap := make(map[string]map[string][]string)

	for _, inDep := range self.InnerDependencies {
		for mapFrom, mapTo := range inDep {

			mapFromItems := strings.Split(mapFrom, ".")
			mapToItems := strings.Split(mapTo, ".")

			mapFromTable := self.Members[mapFromItems[0]]
			mapFromCol := mapFromItems[1]

			mapToTable := self.Members[mapToItems[0]]
			mapToCol := mapToItems[1]

			if _, ok := joinMap[mapFromTable]; !ok {
				joinMap[mapFromTable] = make(map[string][]string)
			}

			var condition string

			condition = fmt.Sprintf("%s.%s=%s.%s", mapFromTable, mapFromCol, mapToTable, mapToCol)

			joinMap[mapFromTable][mapToTable] = append(joinMap[mapFromTable][mapToTable], condition)
		}
	}

	return joinMap
}

func (self Tag) CreateInDepMap(isBag ...bool) map[string]map[string][]string {

	bag := false
	if len(isBag) > 0 {
		if isBag[0] {
			bag = true
		}
	}

	joinMap := make(map[string]map[string][]string)

	for _, inDep := range self.InnerDependencies {
		for mapFrom, mapTo := range inDep {

			mapFromItems := strings.Split(mapFrom, ".")
			mapToItems := strings.Split(mapTo, ".")

			mapFromTable := self.Members[mapFromItems[0]]
			mapFromCol := mapFromItems[1]

			mapToTable := self.Members[mapToItems[0]]
			mapToCol := mapToItems[1]

			if _, ok := joinMap[mapFromTable]; !ok {
				joinMap[mapFromTable] = make(map[string][]string)
			}

			var condition string

			if bag {
				condition = fmt.Sprintf("%s.\"data\"->>'%s.%s'=%s.\"data\"->>'%s.%s'", mapFromTable, mapFromTable, mapFromCol, mapToTable, mapToTable, mapToCol)
			} else {
				condition = fmt.Sprintf("\"%s\".\"%s\"=\"%s\".\"%s\"", mapFromTable, mapFromCol, mapToTable, mapToCol)
			}

			joinMap[mapFromTable][mapToTable] = append(joinMap[mapFromTable][mapToTable], condition)
		}
	}

	return joinMap
}

func (tag Tag) MemberIDs(dbConn *sql.DB, app_id string) (map[string]string, error) {
	tableIDMap := make(map[string]string)
	for _, table := range tag.Members {
		if tableID, err := db.TableID(dbConn, table, app_id); err == nil {
			tableIDMap[table] = tableID
		} else {
			return nil, err
		}
	}
	return tableIDMap, nil
}

func (tag Tag) ResolveRestrictions() string {
	restrictions := ""
	for _, restriction := range tag.Restrictions {
		if restrictionAttr, err := tag.ResolveTagAttr(restriction["col"]); err == nil {
			restrictions += fmt.Sprintf(" AND %s = '%s' ", restrictionAttr, restriction["val"])
		}

	}
	return restrictions
}

func Contains(list []Tag, tagName string) bool {
	for _, v := range list {
		// fmt.Println(v, "==", str)
		if strings.ToLower(v.Name) == strings.ToLower(tagName) {
			return true
		}
	}
	return false
}

package config

import (
	"errors"
	"fmt"
	"log"

	"stencil/helper"
	"stencil/qr"
	"strings"

	"github.com/drgrib/maps"
)

func (self *AppConfig) CloseDBConn() {

	// db.CloseDBConn(self.AppName)
	self.DBConn.Close()
}

func (self AppConfig) GetTag(tagName string) (Tag, error) {

	for _, tag := range self.Tags {
		if strings.EqualFold(tag.Name, tagName) {
			return tag, nil
		}
	}
	return *new(Tag), nil
}

func (self AppConfig) GetTagMembers(tagName string) ([]string, error) {

	var tagMembers []string

	if tag, err := self.GetTag(tagName); err == nil {
		for _, member := range tag.Members {
			tagMembers = append(tagMembers, member)
		}
	} else {
		return tagMembers, errors.New("Tag Not Found")
	}
	return tagMembers, nil
}

func (self Tag) GetTagMembers() []string {

	var tagMembers []string

	for _, member := range self.Members {
		tagMembers = append(tagMembers, member)
	}

	return tagMembers
}

func (self AppConfig) GetDependency(tagName string) (Dependency, error) {

	for _, dep := range self.Dependencies {
		if strings.EqualFold(dep.Tag, tagName) {
			return dep, nil
		}
	}
	return *new(Dependency), nil
}

func (self AppConfig) CheckDependency(tagName, dependsOnTag string) (DependsOn, error) {

	// if deps, err := self.GetDependency(tagName); err == nil {
	for _, dep := range self.Dependencies {
		if strings.EqualFold(dep.Tag, tagName) {
			for _, dependsOn := range dep.DependsOn {
				if strings.EqualFold(dependsOn.Tag, dependsOnTag) {
					return dependsOn, nil
				}
			}
		}
	}

	return *new(DependsOn), nil
}

func (self AppConfig) GetSubDependencies(tagName string) []Dependency {

	var deps []Dependency

	for _, dep := range self.Dependencies {
		for _, do := range dep.DependsOn {
			if strings.EqualFold(do.Tag, tagName) {
				deps = append(deps, dep)
			}
		}
	}

	return deps
}

func (self AppConfig) GetShuffledOwnerships() []Ownership {

	return self.ShuffleOwnerships(self.Ownerships)
}

func (self AppConfig) GetOwnership(tagName string) (Ownership, error) {

	for _, own := range self.Ownerships {
		if strings.EqualFold(own.Tag, tagName) {
			return own, nil
		}
	}
	return *new(Ownership), nil
}

func (self AppConfig) GetItemsFromKey(tag Tag, key string) (string, string) {
	KeyItems := strings.Split(tag.Keys[key], ".")
	Table, Col := tag.Members[KeyItems[0]], KeyItems[1]
	return Table, Col
}

func (tag Tag) ResolveTagAttr(attr string) (string, error) {

	if _, ok := tag.Keys[attr]; ok {
		keyItems := strings.Split(tag.Keys[attr], ".")
		if _, ok := tag.Members[keyItems[0]]; ok {
			Table := tag.Members[keyItems[0]]
			Col := keyItems[1]
			return fmt.Sprintf("%s.%s", Table, Col), nil
		} else {
			return "", errors.New("Tag Not Resolved, Member Not Found")
		}
	}
	return "", errors.New("Tag Not Resolved, Attr Not Found in Tag Keys")
}

func (self Tag) CreateInDepMap() map[string]map[string][]string {

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

			condition := fmt.Sprintf("%s.%s=%s.%s", mapFromTable, mapFromCol, mapToTable, mapToCol)
			joinMap[mapFromTable][mapToTable] = append(joinMap[mapFromTable][mapToTable], condition)
		}
	}

	return joinMap
}

func (self AppConfig) GetTagsByTables(tables []string) []Tag {
	var tags []Tag
	// no member can appear in more than one tag
	for _, tag := range self.Tags {
		if overlap := helper.IntersectString(maps.GetValuesStringString(tag.Members), tables); len(overlap) > 0 {
			tags = append(tags, tag)
		}
	}
	return tags
}

func (self AppConfig) GetTagsByTablesExcept(tables []string, tagName string) []Tag {
	var tags []Tag
	// no member can appear in more than one tag
	for _, tag := range self.Tags {
		if overlap := helper.IntersectString(maps.GetValuesStringString(tag.Members), tables); len(overlap) > 0 {
			if !strings.EqualFold(tag.Name, tagName) {
				tags = append(tags, tag)
			}
		}
	}
	return tags
}

func (self AppConfig) GetDependsOnTables(tagName string, memberID string) []string {
	var dependsOnTables []string
	for _, tag := range self.Tags {
		if tag.Name == tagName {
			for _, innerDependency := range tag.InnerDependencies {
				for dependsOnMember, member := range innerDependency {
					if memberID == strings.Split(member, ".")[0] {
						table, _ := self.GetTableByMemberID(tagName, strings.Split(dependsOnMember, ".")[0])
						dependsOnTables = append(dependsOnTables, table)
					}
				}
			}
		}
	}
	return dependsOnTables
}

func (self Dependency) GetConditionsForTag(tagName string) []DCondition {

	for _, dependsOn := range self.DependsOn {
		if strings.EqualFold(dependsOn.Tag, tagName) {
			return dependsOn.Conditions
		}
	}
	return nil
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

func (self AppConfig) GetTagDisplaySetting(tagName string) (string, error) {

	for _, tag := range self.Tags {
		if tag.Name == tagName {
			if tag.Display_setting != "" {
				return tag.Display_setting, nil
			} else {
				return "default_display_setting", nil
			}
		}
	}

	return "", errors.New("Error: No Tag Found For the Provided TagName")
}

func (self AppConfig) GetTableByMemberID(tagName string, checkedMemberID string) (string, error) {

	for _, tag := range self.Tags {
		if tag.Name == tagName {
			for memberID, memberTable := range tag.Members {
				if memberID == checkedMemberID {
					return memberTable, nil
				}
			}
		}
	}

	return "", errors.New("Error: No Table Found For the Provided Member ID")
}

func (self *AppConfig) GetDependsOnConditions(tagName string, pTagName string) ([]DCondition, error) {
	for _, dp := range self.Dependencies {
		if dp.Tag == tagName {
			for _, dp1 := range dp.DependsOn {
				if dp1.As == pTagName {
					return dp1.Conditions, nil
				} else if dp1.Tag == pTagName {
					return dp1.Conditions, nil
				}
			}
		}
	}

	return nil, errors.New("Error: No Conditions Found")
}

func (self *AppConfig) GetDepDisplaySetting(tag string, pTag string) (string, error) {

	for _, dependency := range self.Dependencies {
		if dependency.Tag == tag {
			for _, dependsOn := range dependency.DependsOn {
				if dependsOn.As != "" {
					if dependsOn.As == pTag {
						return dependsOn.DisplaySetting, nil
					} else {
						continue
					}
				} else {
					if dependsOn.Tag == pTag {
						return dependsOn.DisplaySetting, nil
					}
				}
			}
		}
	}

	return "", errors.New("No dependency display setting is found!")
}

func (self *AppConfig) GetTagQS(tag Tag) *qr.QS {

	qs := qr.CreateQS(self.QR)
	if len(tag.InnerDependencies) > 0 {
		joinMap := tag.CreateInDepMap()
		seenMap := make(map[string]bool)
		for fromTable, toTablesMap := range joinMap {
			if _, ok := seenMap[fromTable]; !ok {
				qs.FromSimple(fromTable)
				qs.ColSimple(fromTable + ".*")
				qs.ColPK(fromTable)
			}
			for toTable, conditions := range toTablesMap {
				if conditions != nil {
					conditions = append(conditions, joinMap[toTable][fromTable]...)
					if joinMap[toTable][fromTable] != nil {
						joinMap[toTable][fromTable] = nil
					}
					qs.FromJoinList(toTable, conditions)
					qs.ColSimple(toTable + ".*")
					qs.ColPK(toTable)
					seenMap[toTable] = true
				}
			}
			seenMap[fromTable] = true
		}
	} else {
		table := tag.Members["member1"]
		qs = qr.CreateQS(self.QR)
		qs.FromSimple(table)
		qs.ColPK(table)
		qs.ColSimple(table + ".*")
	}
	if len(tag.Restrictions) > 0 {
		restrictions := qr.CreateQS(self.QR)
		restrictions.TableAliases = qs.TableAliases
		for _, restriction := range tag.Restrictions {
			if restrictionAttr, err := tag.ResolveTagAttr(restriction["col"]); err == nil {
				restrictions.WhereOperatorInterface("OR", restrictionAttr, "=", restriction["val"])
			}

		}
		qs.WhereString("AND", restrictions.Where)
	}
	return qs
}

func (self *AppConfig) GetTagQSM(tag Tag) *qr.QS {
	log.Fatal("Here in GetTagQSM!")
	qs := qr.CreateQS(self.QR)
	if len(tag.InnerDependencies) > 0 {
		joinMap := tag.CreateInDepMap()
		seenMap := make(map[string]bool)
		for fromTable, toTablesMap := range joinMap {
			if _, ok := seenMap[fromTable]; !ok {
				qs.LSelect(fromTable, "*")
				qs.LTable(fromTable, fromTable)
			}
			for toTable, conditions := range toTablesMap {
				if conditions != nil {
					conditions = append(conditions, joinMap[toTable][fromTable]...)
					if joinMap[toTable][fromTable] != nil {
						joinMap[toTable][fromTable] = nil
					}
					qs.LSelect(toTable, "*")
					qs.LTable(toTable, toTable)
					qs.LJoinOn(conditions)
					seenMap[toTable] = true
				}
			}
			seenMap[fromTable] = true
		}
	} else {
		table := tag.Members["member1"]
		qs = qr.CreateQS(self.QR)
		qs.LSelect(table, "*")
		qs.LTable(table, table)
	}
	if len(tag.Restrictions) > 0 {
		for _, restriction := range tag.Restrictions {
			if restrictionAttr, err := tag.ResolveTagAttr(restriction["col"]); err == nil {
				qs.WhereString("AND", restrictionAttr+"="+restriction["val"])
			}
		}
	}
	return qs
}

package config

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"stencil/db"
	"stencil/helper"
	"stencil/qr"
	"strings"

	"github.com/drgrib/maps"
)

func (self *AppConfig) CloseDBConn() {

	// db.CloseDBConn(self.AppName)
	self.DBConn.Close()
}

func (self *AppConfig) GetTag(tagName string) (Tag, error) {

	for _, tag := range self.Tags {
		if strings.EqualFold(tag.Name, tagName) {
			return tag, nil
		}
	}
	return *new(Tag), nil
}

func (self *AppConfig) GetTagMembers(tagName string) ([]string, error) {

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

func (self *Tag) GetTagMembers() []string {

	var tagMembers []string

	for _, member := range self.Members {
		tagMembers = append(tagMembers, member)
	}

	return tagMembers
}

func (self *AppConfig) GetDependency(tagName string) (Dependency, error) {

	for _, dep := range self.Dependencies {
		if strings.EqualFold(dep.Tag, tagName) {
			return dep, nil
		}
	}
	return *new(Dependency), nil
}

func (self *AppConfig) CheckDependency(tagName, dependsOnTag string) (DependsOn, error) {

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

	return *new(DependsOn), errors.New("Dependency " + tagName + " : " + dependsOnTag + "doesn't exist!")
}

func (self *AppConfig) GetSubDependencies(tagName string) []Dependency {

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

func (self *AppConfig) GetParentDependencies(tagName string) []Dependency {

	var deps []Dependency

	for _, dep := range self.Dependencies {
		if strings.EqualFold(dep.Tag, tagName) {
			deps = append(deps, dep)
		}
	}

	return deps
}

func (self *AppConfig) GetShuffledOwnerships() []Ownership {

	return self.ShuffleOwnerships(self.Ownerships)
}

func (self *AppConfig) GetOrderedOwnerships() []Ownership {

	var orderedOwnerships []Ownership
	orderOfOwnerships := []string{"photo", "post", "like", "comment", "conversation", "message", "contact", "notification"}

	for _, ownershipName := range orderOfOwnerships {
		for _, ownership := range self.Ownerships {
			if strings.EqualFold(ownershipName, ownership.Tag) {
				orderedOwnerships = append(orderedOwnerships, ownership)
			}
		}
	}
	// log.Fatal(orderedOwnerships)
	return orderedOwnerships
}

func (self *AppConfig) GetOwnership(tagName, owner string) *Ownership {

	for _, own := range self.Ownerships {
		if strings.EqualFold(own.Tag, tagName) && strings.EqualFold(own.OwnedBy, owner) {
			return &own
		}
	}
	return nil
}

func (self *AppConfig) GetItemsFromKey(tag Tag, key string) (string, string) {
	if val, ok := tag.Keys[key]; ok {
		KeyItems := strings.Split(val, ".")
		Table, Col := tag.Members[KeyItems[0]], KeyItems[1]
		return Table, Col
	} else {
		fmt.Println(tag.Keys)
		log.Fatal(fmt.Sprintf("@AppConfig.GetItemsFromKey: Key [%s] not found in Tag [%s] in [%s]", key, tag.Name, self.AppName))
		return "", ""
	}
}

func (self *AppConfig) CheckOwnership(tag string) bool {
	for _, ownership := range self.Ownerships {
		if strings.EqualFold(ownership.Tag, tag) {
			return true
		}
	}
	return false
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

func (self *AppConfig) GetTagsByTables(tables []string) []Tag {
	var tags []Tag
	// no member can appear in more than one tag
	for _, tag := range self.Tags {
		if overlap := helper.IntersectString(maps.GetValuesStringString(tag.Members), tables); len(overlap) > 0 {
			tags = append(tags, tag)
		}
	}
	return tags
}

func (self *AppConfig) GetTagByMember(member string) (*Tag, error) {
	// no member can appear in more than one tag
	for _, tag := range self.Tags {
		for _, tagMember := range tag.Members {
			if strings.EqualFold(member, tagMember) {
				return &tag, nil
			}
		}
	}
	return nil, fmt.Errorf("No tag found for member: %s", member)
}

func (self *AppConfig) GetTagsByTablesExcept(tables []string, tagName string) []Tag {
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

func (self *AppConfig) GetDependsOnTables(tagName string, memberID string) []string {
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

func (self *AppConfig) GetTagDisplaySetting(tagName string) (string, error) {

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

func (self *AppConfig) GetTableByMemberID(tagName string, checkedMemberID string) (string, error) {

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

// For finding display setting
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

func (self *AppConfig) GetTagQS(tag Tag, params map[string]string) *qr.QS {

	qs := qr.CreateQS(self.QR)
	if len(tag.InnerDependencies) > 0 {
		joinMap := tag.CreateInDepMap()
		seenMap := make(map[string]bool)
		for fromTable, toTablesMap := range joinMap {
			if _, ok := seenMap[fromTable]; !ok {
				args := map[string]string{"table": fromTable}
				helper.ConcatMaps(args, params)
				qs.FromTable(args)
				qs.SelectColumns(fromTable + ".*")
			}
			for toTable, conditions := range toTablesMap {
				if conditions != nil {
					joinArgs := map[string]string{"table": toTable, "join": "FULL JOIN"}
					helper.ConcatMaps(joinArgs, params)
					conditions = append(conditions, joinMap[toTable][fromTable]...)
					if joinMap[toTable][fromTable] != nil {
						joinMap[toTable][fromTable] = nil
					}
					for i, condition := range conditions {
						joinArgs[fmt.Sprintf("condition%d", i)] = condition
					}
					qs.JoinTable(joinArgs)

					qs.SelectColumns(toTable + ".*")
					seenMap[toTable] = true
				}
			}
			seenMap[fromTable] = true
		}
	} else {
		table := tag.Members["member1"]
		qs = qr.CreateQS(self.QR)
		args := map[string]string{"table": table}
		helper.ConcatMaps(args, params)
		qs.FromTable(args)
		qs.SelectColumns(table + ".*")
	}
	if len(tag.Restrictions) > 0 {
		restrictions := qr.CreateQS(self.QR)
		for _, restriction := range tag.Restrictions {
			if restrictionAttr, err := tag.ResolveTagAttr(restriction["col"]); err == nil {
				restrictions.AdditionalWhereWithValue("OR", restrictionAttr, "=", restriction["val"])
			}
		}
		if restrictions.Where == "" {
			log.Fatal(tag.Restrictions)
		}
		qs.AddWhereAsString("AND", restrictions.Where)
	}
	// log.Fatal(qs.GenSQL())
	return qs
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

func (self *AppConfig) CloseDBConns() {
	self.DBConn.Close()
	self.QR.StencilDB.Close()
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

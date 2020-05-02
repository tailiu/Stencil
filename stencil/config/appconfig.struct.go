package config

import (
	"errors"
	"fmt"
	"log"
	"stencil/helper"
	"stencil/qr"
	"strings"

	"github.com/drgrib/maps"
	"github.com/gookit/color"
)

func (self *AppConfig) CloseDBConn() {

	// db.CloseDBConn(self.AppName)
	color.Info.Printf("AppConfig.CloseDBConn: Closing DB connection for app: \"%s\"\n", self.AppName)
	self.DBConn.Close()
}

func (self *AppConfig) GetTag(tagName string) (Tag, error) {

	for _, tag := range self.Tags {
		if strings.EqualFold(tag.Name, tagName) {
			return tag, nil
		}
	}
	return *new(Tag), errors.New("Tag doesn't exist for: " + tagName)
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
		joinMap := tag.CreateInDepMapSA2()
		seenMap := make(map[string]bool)
		for _, fromTable := range tag.GetInDepMembersInOrder() {
			toTablesMap := joinMap[fromTable]
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

func (self *AppConfig) __defunct__GetTagQS(tag Tag, params map[string]string) *qr.QS {

	qs := qr.CreateQS(self.QR)
	if len(tag.InnerDependencies) > 0 {
		joinMap := tag.CreateInDepMapSA2()
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

func (self *AppConfig) CloseDBConns() {
	color.Info.Printf("AppConfig.CloseDBConns: Closing DB connection for app: \"%s\"\n", self.AppName)
	self.DBConn.Close()
	color.Info.Printf("AppConfig.CloseDBConns: Closing DB connection for stencil in app: \"%s\"\n", self.AppName)
	self.QR.StencilDB.Close()
}

package config

import (
	"errors"
	"fmt"
	"strings"
	"transaction/helper"

	"github.com/drgrib/maps"
)

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

			condition := fmt.Sprintf("%s.%s = %s.%s", mapToTable, mapToCol, mapFromTable, mapFromCol)
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

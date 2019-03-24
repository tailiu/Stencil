package config

import (
	"errors"
	"fmt"
	"strings"
)

func (self AppConfig) GetTag(tagName string) (Tag, error) {

	for _, tag := range self.Tags {
		if strings.EqualFold(tag.Name, tagName) {
			return tag, nil
		}
	}
	return *new(Tag), nil
}

func (self AppConfig) GetDependency(tagName string) (Dependency, error) {

	for _, dep := range self.Dependencies {
		if strings.EqualFold(dep.Tag, tagName) {
			return dep, nil
		}
	}
	return *new(Dependency), nil
}

func (self AppConfig) CheckDependency(tagName, dependsOn string) bool {

	if deps, err := self.GetDependency(tagName); err == nil {
		for _, dep := range deps.DependsOn {
			if strings.EqualFold(dep.Tag, dependsOn) {
				return true
			}
		}
	}

	return false
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

func (self AppConfig) ResolveTagAttr(tagName, attr string) (string, error) {

	if tag, err := self.GetTag(tagName); err == nil {
		if _, ok := tag.Keys[attr]; ok {
			keyItems := strings.Split(tag.Keys[attr], ".")
			if _, ok := tag.Members[keyItems[0]]; ok {
				Table := tag.Members[keyItems[0]]
				Col := keyItems[1]
				return fmt.Sprintf("%s.%s", Table, Col), nil
			} else {
				return "", errors.New("Tag Not Resolved, Member Not Found")
			}
		} else {
			return "", errors.New("Tag Not Resolved, Attr Not Found in Tag Keys")
		}
	}
	return "", errors.New("Tag Not Resolved")
}

func (self AppConfig) CreateInDepMap(tag Tag) map[string]map[string][]string {

	joinMap := make(map[string]map[string][]string)

	for _, inDep := range tag.InnerDependencies {
		for mapFrom, mapTo := range inDep {

			mapFromItems := strings.Split(mapFrom, ".")
			mapToItems := strings.Split(mapTo, ".")

			mapFromTable := tag.Members[mapFromItems[0]]
			mapFromCol := mapFromItems[1]

			mapToTable := tag.Members[mapToItems[0]]
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

package config

import (
	"errors"
	"strings"
)

func (self Dependency) GetConditionsForTag(tagName string) []DCondition {

	for _, dependsOn := range self.DependsOn {
		if strings.EqualFold(dependsOn.Tag, tagName) {
			return dependsOn.Conditions
		}
	}
	return nil
}

func FindDependency(tag, depends_on string, dependencies []Dependency) (Dependency, error) {

	for _, dependency := range dependencies {
		if strings.ToLower(dependency.Tag) == strings.ToLower(tag) {
			// && strings.ToLower(dependency.DependsOn) == strings.ToLower(depends_on)

			return dependency, nil
		}
	}
	return *new(Dependency), errors.New("dependency doesn't exist")
}

func FindDependencyByDependsOn(depends_on string, dependencies []Dependency) (Dependency, error) {

	for _, dependency := range dependencies {
		if strings.ToLower(dependency.Tag) == strings.ToLower(depends_on) {
			return dependency, nil
		}
	}
	return *new(Dependency), errors.New("dependency doesn't exist")
}

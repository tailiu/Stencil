package config

import (
	"errors"
	"fmt"
	"math/rand"
	"stencil/helper"
	"strings"
)

func (self *MappedApp) GetInput(conditionVal string) (string, error) {
	conditionVal = strings.TrimLeft(conditionVal, "$")
	for i, input := range self.Inputs {
		if strings.EqualFold(input["name"], conditionVal) {
			if strings.EqualFold(input["value"], "#RANDINT") {
				self.Inputs[i]["value"] = fmt.Sprint(rand.Int31n(2147483647))
				return self.Inputs[i]["value"], nil
			}
			return input["value"], nil
		}
	}
	return "", errors.New("Can't find Input with name: " + conditionVal)
}

func (self *MappedApp) SetInput(inputName, valueToSet string) (string, error) {
	inputName = strings.TrimLeft(inputName, "$")
	for i, input := range self.Inputs {
		if strings.EqualFold(input["name"], inputName) {
			if strings.Contains(input["value"], "#") {
				self.Inputs[i]["value"] = valueToSet
				return self.Inputs[i]["value"], nil
			}
			return input["value"], nil
		}
	}
	return "", errors.New("Can't find Input with name: " + inputName)
}

func (self *MappedApp) GetMethod(conditionVal string) (string, error) {
	conditionVal = strings.TrimLeft(conditionVal, "#")
	for _, method := range self.Methods {
		if strings.EqualFold(method["name"], conditionVal) {
			return method["value"], nil
		}
	}
	return "", errors.New("Can't find Method with name: " + conditionVal)
}

func (self *ToTable) FromTables() []string {
	var fromTables []string
	for _, mappedTo := range self.Mapping {
		if strings.Contains(mappedTo, "$") || strings.Contains(mappedTo, "#") {
			continue
		}
		table := strings.Split(mappedTo, ".")[0]
		if !helper.Contains(fromTables, table) {
			fromTables = append(fromTables, table)
		}
	}
	return fromTables
}

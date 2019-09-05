package config

import (
	"errors"
	"strings"
	"fmt"
	"math/rand"
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

func (self *MappedApp) GetMethod(conditionVal string) (string, error) {
	conditionVal = strings.TrimLeft(conditionVal, "#")
	for _, method := range self.Methods {
		if strings.EqualFold(method["name"], conditionVal) {
			return method["value"], nil
		}
	}
	return "", errors.New("Can't find Method with name: " + conditionVal)
}

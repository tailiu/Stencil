package config

import (
	"errors"
	"strings"
)

func (self MappedApp) GetInput(key string) (string, error) {
	key = strings.Trim(key, "$")
	for _, input := range self.Inputs {
		if strings.EqualFold(input["name"], key) {
			return input["value"], nil
		}
	}
	return "", errors.New("Can't find Input with name: " + key)
}

func (self MappedApp) GetMethod(key string) (string, error) {
	key = strings.Trim(key, "$")
	for _, input := range self.Methods {
		if strings.EqualFold(input["name"], key) {
			return input["mapping"], nil
		}
	}
	return "", errors.New("Can't find Method with name: " + key)
}

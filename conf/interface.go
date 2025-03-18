// Package conf v0.1.0
package conf

import (
	"errors"
	"fmt"
)

var (
	parameters map[string]string
	hash       string
)

func Parse(parsers ...Parser) error {
	parameters = nil
	hash = getRequiredString()

	result := make(map[string]string)

	for _, parser := range parsers {
		for key, value := range parser.Parse(GetKeyConverter()) {
			if len(result[key]) != 0 {
				return errors.New(fmt.Sprintf("parameter '%s' was duplicated", key))
			}
			if value == hash {
				return errors.New(fmt.Sprintf("required parameter '%s' is not set", key))
			}
			result[key] = value
		}
	}
	parameters = result
	return nil
}

// Get is function to get parameter by key. If parameter not found, returns empty string.
func Get(key string) string {
	c := GetKeyConverter()
	return parameters[c(key)]
}

// Getc is equivalent of Get but function return Converter instead of raw value.
func Getc(key string) Converter {
	return NewConverter(Get(key))
}

// Geta it's short for GetAll. Function returns all handled parameters.
func Geta() map[string]string {
	return parameters
}

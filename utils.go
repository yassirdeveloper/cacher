package main

import (
	"fmt"
	"strconv"
)

// ParseValue attempts to parse a string into the specified type.
func ParseValue(valueType string, value string) (interface{}, error) {
	switch valueType {
	case TypeInt:
		return strconv.Atoi(value) // Parse string to int
	case TypeFloat:
		return strconv.ParseFloat(value, 64) // Parse string to float64
	case TypeBool:
		return strconv.ParseBool(value) // Parse string to bool
	case TypeString:
		return value, nil // No parsing needed for strings
	default:
		return nil, fmt.Errorf("unsupported type: %s", valueType)
	}
}

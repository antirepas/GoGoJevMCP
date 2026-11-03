package ollama

import (
	"errors"
	"fmt"
	"math"
)

func matchesType(value any, expectedType string) bool {
	switch expectedType {
	case "string":
		_, ok := value.(string)
		return ok

	case "number":
		_, ok := value.(float64)
		return ok

	case "integer":
		num, ok := value.(float64)
		return ok && math.Trunc(num) == num

	case "boolean":
		_, ok := value.(bool)
		return ok

	default:
		return false
	}
}

func ValidateArguments(schema map[string]any, args map[string]any) error {
	reqs, _ := schema["required"].([]string)
	for _, required := range reqs {
		_, exists := args[required]
		if !exists {
			return errors.New("required argument does not exist")
		}
	}

	properties, _ := schema["properties"].(map[string]any)
	fmt.Println(properties)
	for key, val := range args {
		props, ok := properties[key].(map[string]any)
		if !ok {
			return errors.New("schema has different arguments")
		}
		exType, ok := props["type"].(string)
		if !ok {
			return errors.New("invalid schema: missing argument type")
		}
		if !matchesType(val, exType) {
			return errors.New("argument has incorrect type")
		}
	}
	return nil
}

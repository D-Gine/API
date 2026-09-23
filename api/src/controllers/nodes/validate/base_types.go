package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type numberRestrictions struct {
	Min *int `json:"min,omitempty"`
	Max *int `json:"max,omitempty"`
}

func GetAllowedFields(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ""
	}

	var fields []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			// Extract tag name before any options like omitempty
			tagName := strings.Split(jsonTag, ",")[0]
			fields = append(fields, fmt.Sprintf("- %s: %s (optional)", tagName, field.Type.Kind()))
		}
	}
	return strings.Join(fields, "\n  ")
}

func verifyNumber(value int, restrictions numberRestrictions) error {
	if value > *restrictions.Max || value < *restrictions.Min {
		return errors.New("invalid number: value out of bounds")
	}
	return nil
}

func ValidateBaseType(valueJson json.RawMessage, name string, restrictionsJson json.RawMessage) error {
	switch name {
	case "number":
		var restrictions numberRestrictions
		var value int

		if err := json.Unmarshal(valueJson, &value); err != nil {
			return fmt.Errorf(`invalid value for "%s": %w`, name, err)
		}
		if err := json.Unmarshal(restrictionsJson, &restrictions); err != nil {
			return fmt.Errorf(`invalid restrictions for "%s": %w
				Allowed fields:
				%s`, name, err, GetAllowedFields(restrictions))
		}
		return verifyNumber(value, restrictions)
	}

	return nil
}

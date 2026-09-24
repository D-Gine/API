package validate

import (
	"api/src/database"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/lib/pq"
)

type numberRestrictions struct {
	Min *int `json:"min,omitempty"`
	Max *int `json:"max,omitempty"`
}

type enumRestrictions struct {
	Tags *[]string `json:"tags,omitempty"`
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
	if restrictions.Max != nil && value > *restrictions.Max {
		return errors.New("value too big")
	}
	if restrictions.Min != nil && value < *restrictions.Min {
		return errors.New("value too low")
	}
	return nil
}

func verifyEnum(value string, restrictions enumRestrictions) error {
	if restrictions.Tags != nil {
		rows := database.Db.QueryRow(`
			SELECT COUNT(*) FROM link_enums_tags WHERE enum_id = $1 AND tag_id = ANY($2::uuid[])`, value, pq.Array(*restrictions.Tags))
		var count int
		rows.Scan(&count)
		if count != len(*restrictions.Tags) {
			return errors.New("does not respect tag restrictions")
		}
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

		if err := verifyNumber(value, restrictions); err != nil {
			return errors.New("invalid number: " + err.Error())
		}
	case "enum":
		var restrictions enumRestrictions
		var value string

		if err := json.Unmarshal(valueJson, &value); err != nil {
			return fmt.Errorf(`invalid value for "%s": %w`, name, err)
		}
		if err := json.Unmarshal(restrictionsJson, &restrictions); err != nil {
			return fmt.Errorf(`invalid restrictions for "%s": %w
				Allowed fields:
				%s`, name, err, GetAllowedFields(restrictions))
		}
		if err := verifyEnum(value, restrictions); err != nil {
			return errors.New("invalid enum: " + err.Error())
		}
	}

	return nil
}

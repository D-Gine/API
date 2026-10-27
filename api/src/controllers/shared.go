package controllers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type ComponentMap = map[string]json.RawMessage

type BaseRestrictions struct {
	OptionalBool *bool `json:"optional,omitempty"`
}

func (b BaseRestrictions) GetBaseRestrictions() BaseRestrictions {
	return b
}

type NumberRestrictions struct {
	BaseRestrictions
	Min *int `json:"min,omitempty"`
	Max *int `json:"max,omitempty"`
}

type StringRestrictions struct {
	BaseRestrictions
	MinLength       *int    `json:"min_length,omitempty"`
	MaxLength       *int    `json:"max_length,omitempty"`
	RegexExpression *string `json:"regex_expression,omitempty"`
}

type EnumRestrictions struct {
	BaseRestrictions
	Tags []string `json:"tags,omitempty"`
}

type ListRestrictions struct {
	BaseRestrictions
	Type      string `json:"type"`
	MinLength *int   `json:"min_length,omitempty"`
	MaxLength *int   `json:"max_length,omitempty"`
}

type DictionaryRestrictions struct {
	BaseRestrictions
	Interface *map[string]string `json:"interface,omitempty"`
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
	for field := range t.Fields() {
		field := field
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			// Extract tag name before any options like omitempty
			tagName, _, _ := strings.Cut(jsonTag, ",")
			fields = append(fields, fmt.Sprintf("- %s: %s (optional)", tagName, field.Type.Kind()))
		}
	}
	return strings.Join(fields, "\n  ")
}

func ParseRestriction[T interface {
	GetBaseRestrictions() BaseRestrictions
}](
	name string,
	restrictionsJson json.RawMessage,
) (T, error) {

	var restrictions T
	if err := json.Unmarshal(restrictionsJson, &restrictions); err != nil {
		return restrictions, fmt.Errorf(`invalid restrictions for "%s": %w
			Allowed fields:
			%s`, name, err, GetAllowedFields(restrictions))
	}
	return restrictions, nil
}

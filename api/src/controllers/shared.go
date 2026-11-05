package controllers

import (
	"fmt"
	"reflect"
	"strings"
)

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

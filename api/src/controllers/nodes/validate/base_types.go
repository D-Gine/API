package validate

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
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

func InitBaseTypeVerification[Value any, Restrictions interface{ GetBaseRestrictions() BaseRestrictions }](
	ctx *gin.Context,
	nodeList *ComponentMap,
	name string,
	valueJson json.RawMessage,
	restrictionsJson json.RawMessage,
	verifyFunc func(*gin.Context, *ComponentMap, Value, Restrictions) error) error {

	var restrictions Restrictions
	var value *Value

	if err := json.Unmarshal(valueJson, &value); err != nil {
		return fmt.Errorf(`invalid value for "%s": %w`, name, err)
	}
	if err := json.Unmarshal(restrictionsJson, &restrictions); err != nil {
		return fmt.Errorf(`invalid restrictions for "%s": %w
			Allowed fields:
			%s`, name, err, GetAllowedFields(restrictions))
	}
	if value == nil {
		opt := restrictions.GetBaseRestrictions().OptionalBool
		if opt != nil && *opt {
			return nil
		} else {
			return fmt.Errorf("invalid %s: can't be null", name)
		}
	}

	if err := verifyFunc(ctx, nodeList, *value, restrictions); err != nil {
		return fmt.Errorf("invalid %s: %w", name, err)
	}
	return nil
}

func ValidateBaseType(ctx *gin.Context, nodeList *ComponentMap, valueJson json.RawMessage, name string, restrictionsJson json.RawMessage) error {
	switch name {
	case "number":
		return InitBaseTypeVerification(ctx, nodeList, name, valueJson, restrictionsJson, VerifyNumber)
	case "string":
		return InitBaseTypeVerification(ctx, nodeList, name, valueJson, restrictionsJson, VerifyString)
	case "enum":
		return InitBaseTypeVerification(ctx, nodeList, name, valueJson, restrictionsJson, VerifyEnum)
	case "list":
		return InitBaseTypeVerification(ctx, nodeList, name, valueJson, restrictionsJson, VerifyList)
	case "dictionary":
		return InitBaseTypeVerification(ctx, nodeList, name, valueJson, restrictionsJson, VerifyDictionary)
	}

	return nil
}

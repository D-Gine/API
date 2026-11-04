package type_validation

import (
	"api/src/controllers/shared_types"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitBaseTypeVerification[Value any, Restrictions interface {
	GetBaseRestrictions() shared_types.BaseRestrictions
}](
	ctx *gin.Context,
	nodeList *shared_types.ComponentMap,
	name string,
	valueJson json.RawMessage,
	restrictionsJson json.RawMessage,
	verifyFunc func(*gin.Context, *shared_types.ComponentMap, Value, Restrictions) error) error {

	var value *Value

	if err := json.Unmarshal(valueJson, &value); err != nil {
		return fmt.Errorf(`invalid value for "%s": %w`, name, err)
	}
	restrictions, err := shared_types.ParseRestriction[Restrictions](name, restrictionsJson)
	if err != nil {
		return err
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

func ValidateBaseType(ctx *gin.Context, nodeList *shared_types.ComponentMap, valueJson json.RawMessage, name string, restrictionsJson json.RawMessage) error {
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

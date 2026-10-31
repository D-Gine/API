package shared_types

import (
	"api/src/controllers"
	"encoding/json"
	"fmt"
)

type BaseVerifications struct {
	BaseName     string          `json:"base_name"`
	Scripts      []string        `json:"scripts"`
	Restrictions json.RawMessage `json:"restrictions"`
}

type TypeVerifications struct {
	TypeId string `json:"type_id"`
	BaseVerifications
}

type Verification interface {
	GetBase() BaseVerifications
}

func (t TypeVerifications) GetBase() BaseVerifications {
	return t.BaseVerifications
}

// restrictions
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
			%s`, name, err, controllers.GetAllowedFields(restrictions))
	}
	return restrictions, nil
}

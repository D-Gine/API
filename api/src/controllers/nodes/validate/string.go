package validate

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
)

type stringRestrictions struct {
	BaseRestrictions
	MinLength       *int    `json:"min_length,omitempty"`
	MaxLength       *int    `json:"max_length,omitempty"`
	RegexExpression *string `json:"regex_expression,omitempty"`
}

func VerifyString(ctx *gin.Context, componentList *ComponentMap, value string, restrictions stringRestrictions) error {
	if restrictions.MaxLength != nil && len(value) > *restrictions.MaxLength {
		return fmt.Errorf("too long, max lenght is %d, current len is %d", *restrictions.MaxLength, len(value))
	}
	if restrictions.MinLength != nil && len(value) < *restrictions.MinLength {
		return fmt.Errorf("too short, min lenght is %d, current len is %d", *restrictions.MinLength, len(value))
	}
	if restrictions.RegexExpression != nil {
		match, err := regexp.MatchString(*restrictions.RegexExpression, value)
		if !match {
			return fmt.Errorf("regex expression not matching, %w", err)
		}
	}
	return nil
}

package validate

import (
	"api/src/controllers"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
)

func VerifyString(ctx *gin.Context, componentList *controllers.ComponentMap, value string, restrictions controllers.StringRestrictions) error {
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

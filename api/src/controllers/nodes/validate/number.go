package validate

import (
	"api/src/controllers"
	"fmt"

	"github.com/gin-gonic/gin"
)

func VerifyNumber(ctx *gin.Context, componentList *controllers.ComponentMap, value int, restrictions controllers.NumberRestrictions) error {
	if restrictions.Max != nil && value > *restrictions.Max {
		return fmt.Errorf("too hight, max value is %d, current value is %d", *restrictions.Max, value)
	}
	if restrictions.Min != nil && value < *restrictions.Min {
		return fmt.Errorf("too low, min value is %d, current value is %d", *restrictions.Min, value)
	}
	return nil
}

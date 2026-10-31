package type_validation

import (
	"api/src/controllers/shared_types"
	"fmt"

	"github.com/gin-gonic/gin"
)

func VerifyNumber(ctx *gin.Context, componentList *shared_types.ComponentMap, value int, restrictions shared_types.NumberRestrictions) error {
	if restrictions.Max != nil && value > *restrictions.Max {
		return fmt.Errorf("too hight, max value is %d, current value is %d", *restrictions.Max, value)
	}
	if restrictions.Min != nil && value < *restrictions.Min {
		return fmt.Errorf("too low, min value is %d, current value is %d", *restrictions.Min, value)
	}
	return nil
}

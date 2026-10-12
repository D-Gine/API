package validate

import (
	"fmt"
)

type numberRestrictions struct {
	Min *int `json:"min,omitempty"`
	Max *int `json:"max,omitempty"`
}

func VerifyNumber(componentList *ComponentMap, value int, restrictions numberRestrictions) error {
	if restrictions.Max != nil && value > *restrictions.Max {
		return fmt.Errorf("too hight, max value is %d, current value is %d", *restrictions.Max, value)
	}
	if restrictions.Min != nil && value < *restrictions.Min {
		return fmt.Errorf("too low, min value is %d, current value is %d", *restrictions.Min, value)
	}
	return nil
}

package validate

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
)

func VerifyType(componentList *ComponentMap, valueJson json.RawMessage, baseVerification Verification) error {
	verification := baseVerification.GetBase()
	if err := ValidateBaseType(componentList, valueJson, verification.BaseName, verification.Restrictions); err != nil {
		return fmt.Errorf("error validation base type: %w", err)
	}

	for i := 0; i < len(verification.Scripts); i++ {
		if err := runValidateScript(
			valueJson,
			verification.Scripts[i],
			verification.Restrictions,
			componentList,
		); err != nil {
			return err
		}
	}
	return nil
}

func ValidateRowTypes(componentList *ComponentMap, rows *sql.Rows) (err error, errors_map map[string]string, count int) {
	errors_map = make(map[string]string)
	var result ComponentVerification
	for count = 0; rows.Next(); count++ {
		if err = rows.Scan(&result.ComponentId, &result.BaseName, pq.Array(&result.Scripts), &result.Restrictions); err != nil {
			return
		}
		if type_err := VerifyType(componentList, (*componentList)[result.ComponentId], result); type_err != nil {
			errors_map[result.ComponentId] = type_err.Error()
		}
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

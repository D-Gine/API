package validate

import (
	"api/src/controllers"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func VerifyType(ctx *gin.Context, componentList *controllers.ComponentMap, valueJson json.RawMessage, baseVerification Verification) error {
	verification := baseVerification.GetBase()
	if err := ValidateBaseType(ctx, componentList, valueJson, verification.BaseName, verification.Restrictions); err != nil {
		return fmt.Errorf("error validation base type: %w", err)
	}

	for i := 0; i < len(verification.Scripts); i++ {
		if err := runValidateScript(
			ctx,
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

func ValidateRowTypes(ctx *gin.Context, componentList *ComponentMap, rows *sql.Rows) (errors_map map[string]string, count int, err error) {
	errors_map = make(map[string]string)
	var result ComponentVerification
	for count = 0; rows.Next(); count++ {
		if err = rows.Scan(&result.ComponentId, &result.BaseName, pq.Array(&result.Scripts), &result.Restrictions); err != nil {
			return
		}
		if type_err := VerifyType(ctx, componentList, (*componentList)[result.ComponentId], result); type_err != nil {
			errors_map[result.ComponentId] = type_err.Error()
		}
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

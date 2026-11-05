package type_validation

import (
	"api/src/controllers/shared_types"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func VerifyType(ctx *gin.Context, componentList *shared_types.ComponentMap, valueJson json.RawMessage, baseVerification shared_types.Verification) error {
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

func ValidateRowTypes(
	ctx *gin.Context,
	toVerify *shared_types.TypeValueMap,
	components *shared_types.ComponentMap,
	rows *sql.Rows,
) (errors_map map[string]string, count int, err error) {
	errors_map = make(map[string]string)
	var result shared_types.TypeVerifications
	for count = 0; rows.Next(); count++ {
		if err = rows.Scan(&result.TypeId, &result.BaseName, pq.Array(&result.Scripts), &result.Restrictions); err != nil {
			return
		}
		if type_err := VerifyType(ctx, components, (*toVerify)[result.TypeId], result); type_err != nil {
			errors_map[result.TypeId] = type_err.Error()
		}
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func ValidateRowComponents(
	ctx *gin.Context,
	components *shared_types.ComponentMap,
	rows *sql.Rows,
) (map[string]string, int, error) {
	return ValidateRowComponentsWithValues(ctx, components, components, rows)
}

func ValidateRowComponentsWithValues(
	ctx *gin.Context,
	toVerify *shared_types.ComponentMap,
	components *shared_types.ComponentMap,
	rows *sql.Rows,
) (errors_map map[string]string, count int, err error) {

	errors_map = make(map[string]string)
	var result shared_types.ComponentVerification
	for count = 0; rows.Next(); count++ {
		if err = rows.Scan(&result.ComponentId, &result.BaseName, pq.Array(&result.Scripts), &result.Restrictions); err != nil {
			return
		}
		if type_err := VerifyType(ctx, components, (*toVerify)[result.ComponentId], result); type_err != nil {
			errors_map[result.ComponentId] = type_err.Error()
		}
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

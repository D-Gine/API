package type_validation

import (
	"api/src/controllers/shared_types"
	"api/src/internal/scripting"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

func runValidateScript(ctx *gin.Context, value json.RawMessage, script string, restrictions json.RawMessage, componentList *shared_types.ComponentMap) error {
	rawResult, err := scripting.ScriptEngine.Eval(
		ctx.Request.Context(),
		script,
		"validate_type",
		[]any{value, restrictions, *componentList},
		nil,
	)
	if err != nil {
		return fmt.Errorf("script execution error: %w", err)
	}
	stringError, err := scripting.AsOptionalString(rawResult)
	if err != nil {
		return fmt.Errorf("script return type error: %w", err)
	}
	if stringError == nil {
		return nil
	}
	return errors.New(*stringError)
}

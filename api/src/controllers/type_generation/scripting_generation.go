package typegeneration

import (
	"api/src/controllers"
	"api/src/internal/scripting"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ScriptingGeneration(ctx *gin.Context, staticValues controllers.ComponentMap, infos GenerationInfos) (json.RawMessage, int, error) {
	value, err := scripting.ScriptEngine.Eval(
		ctx.Request.Context(),
		*infos.Script,
		"generate_type",
		[]any{infos.Restrictions, staticValues},
		nil,
	)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	rawValue, err := json.Marshal(value)
	if err != nil {
		return nil, http.StatusUnauthorized, fmt.Errorf("unable to parse return value for '%s': %w", infos.TypeName, err)
	}
	return rawValue, http.StatusCreated, nil
}

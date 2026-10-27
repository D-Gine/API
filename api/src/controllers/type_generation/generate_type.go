package typegeneration

import (
	"api/src/controllers"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type GenerationInfos struct {
	Script       *string         `json:"script,omitempty"`
	TypeName     string          `json:"type_name"`
	Restrictions json.RawMessage `json:"restrictions"`
}

func GenerateKnownedType(ctx *gin.Context, staticValues controllers.ComponentMap, infos GenerationInfos) (json.RawMessage, int, error) {
	if infos.Script != nil {
		return ScriptingGeneration(ctx, staticValues, infos)
	}
	return BaseTypeGeneration(ctx, staticValues, infos)
}

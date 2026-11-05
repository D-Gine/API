package typegeneration

import (
	"api/src/controllers/shared_types"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GenerationInfos struct {
	Script       *string         `json:"script,omitempty"`
	TypeName     string          `json:"type_name"`
	Restrictions json.RawMessage `json:"restrictions"`
}

func GenerateKnownedType(ctx *gin.Context, staticValues shared_types.ComponentMap, infos GenerationInfos) (json.RawMessage, int, error) {
	if infos.Script != nil {
		return ScriptingGeneration(ctx, staticValues, infos)
	}
	return BaseTypeGeneration(ctx, staticValues, infos)
}

func GenerateRowTypes(
	ctx *gin.Context,
	components shared_types.ComponentMap,
	rows *sql.Rows,
) (generated shared_types.TypeValueMap, code int, err error) {

	generated = make(shared_types.TypeValueMap)
	for rows.Next() {
		var result GenerationInfos
		var type_id string
		if err = rows.Scan(&type_id, &result.Script, &result.TypeName, &result.Restrictions); err != nil {
			code = http.StatusInternalServerError
			return
		}
		var value json.RawMessage
		value, code, err = GenerateKnownedType(ctx, components, result)
		if err != nil {
			return
		}
		generated[type_id] = value
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

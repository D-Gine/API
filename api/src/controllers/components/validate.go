package components

import (
	"api/src/controllers/shared_types"
	"api/src/controllers/type_validation"
	"api/src/database"
	"encoding/json"
	"errors"
	"maps"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func ValidateComponentsValues(
	ctx *gin.Context,
	ruleset_id string,
	toVerify shared_types.ComponentMap,
	components shared_types.ComponentMap,
) (int, map[string]string, error) {
	row, err := database.Db.Query(`
		SELECT
			sct.id AS component_id,
			gv.base_name AS base_name,
			gv.scripts AS scripts,
			gv.restrictions AS restrictions
		FROM static_components_templates sct
		JOIN get_verifications((
			SELECT ARRAY_AGG(type_id)
			FROM static_components_templates
			WHERE ruleset_id = $1 AND id = sct.id
		)::uuid[]) gv ON 1 = 1
		WHERE
			sct.id = ANY($2)
	`, ruleset_id, slices.Collect(maps.Keys(toVerify)))
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	defer database.CloseRows(row)

	error_map, count, err := type_validation.ValidateRowComponentsWithValues(ctx, &toVerify, &components, row)
	if err != nil {
		return http.StatusInternalServerError, error_map, err
	}
	if count != len(toVerify) {
		return http.StatusUnauthorized, error_map, errors.New("some sent components are invalid, cannot compute")
	}
	code := http.StatusOK
	if len(error_map) != 0 {
		code = http.StatusUnauthorized
	}
	return code, error_map, nil
}

// @BasePath /api/:ruleset_id/component/:id/validate
// Nodes godoc
// @Summary Components - Validate one
// @Schemes
// @Description validate a component with component map
// @Tags components
// @Accept json
// @Produce json
// @Param creds body ValidateValue true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} map[string]string "Empty Object"
// @Router /api/:ruleset_id/component/:id/validate [post]
func ValidateComponent(ctx *gin.Context) {
	ruleset_id := ctx.Param("ruleset_id")
	id := ctx.Param("id")

	var body ValidateValue
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Invalid body"})
		return
	}
	code, error_map, err := ValidateComponentsValues(ctx, ruleset_id, map[string]json.RawMessage{id: body.Value}, body.Components)

	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(code, error_map)
}

// @BasePath /api/:ruleset_id/component/validate
// Nodes godoc
// @Summary Components - Validate multiple
// @Schemes
// @Description validate mmultiple component with component map
// @Tags components
// @Accept json
// @Produce json
// @Param creds body ValidateValues true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} map[string]string "Empty Object"
// @Router /api/:ruleset_id/component/validate [post]
func ValidateComponents(ctx *gin.Context) {
	ruleset_id := ctx.Param("ruleset_id")

	var body ValidateValues
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Invalid body"})
		return
	}
	code, error_map, err := ValidateComponentsValues(ctx, ruleset_id, body.Values, body.Components)

	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(code, error_map)
}

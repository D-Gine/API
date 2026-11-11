package types

import (
	"api/src/database"
	"encoding/json"
	"errors"
	"maps"
	"net/http"
	"slices"

	gdt "github.com/D-Gine/Gestionnaire2Type"
	"github.com/D-Gine/Gestionnaire2Type/validate"
	"github.com/gin-gonic/gin"
)

func ValidateTypesValues(
	ruleset_id string,
	values gdt.TypeValueMap,
	components gdt.ComponentMap,
) (int, map[string]string, error) {

	row, err := database.Db.Query(`
		SELECT
			t.id,
			gv.base_name,
			gv.scripts,
			gv.restrictions
		FROM types t
		JOIN get_verifications($2::uuid[]) gv ON gv.target_type_id = t.id
		WHERE t.ruleset_id = $1
	`, ruleset_id, slices.Collect(maps.Keys(values)))
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	defer database.CloseRows(row)

	error_map, count, err := validate.ValidateRowTypes(database.Db, &values, &components, row)
	if err != nil {
		return http.StatusInternalServerError, error_map, err
	}
	if count != len(values) {
		return http.StatusUnauthorized, error_map, errors.New("some sent components are invalid, cannot compute")
	}
	code := http.StatusOK
	if len(error_map) != 0 {
		code = http.StatusUnauthorized
	}
	return code, error_map, nil
}

// @BasePath /api/:ruleset_id/type/:id/validate
// Nodes godoc
// @Summary Types - validate one type
// @Schemes
// @Description validate a type value with component map
// @Tags types
// @Accept json
// @Produce json
// @Param creds body ValidateValue true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} map[string]string "Empty Object"
// @Router /api/:ruleset_id/type/:id/validate [post]
func ValidateType(ctx *gin.Context) {
	ruleset_id := ctx.Param("ruleset_id")
	id := ctx.Param("id")

	var body ValidateValue
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Invalid body"})
		return
	}
	code, error_map, err := ValidateTypesValues(ruleset_id, map[string]json.RawMessage{id: body.Value}, body.Components)

	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(code, error_map)
}

// @BasePath /api/:ruleset_id/type/validate
// Nodes godoc
// @Summary Types - validate multiple types
// @Schemes
// @Description validate mmultiple types with component map
// @Tags types
// @Accept json
// @Produce json
// @Param creds body ValidateValues true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} map[string]string "Empty Object"
// @Router /api/:ruleset_id/type/validate [post]
func ValidateTypes(ctx *gin.Context) {
	ruleset_id := ctx.Param("ruleset_id")

	var body ValidateValues
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Invalid body"})
		return
	}
	code, error_map, err := ValidateTypesValues(ruleset_id, body.Values, body.Components)

	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(code, error_map)
}

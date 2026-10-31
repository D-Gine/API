package components

import (
	"api/src/controllers/shared_types"
	"api/src/controllers/type_validation"
	"api/src/database"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/:ruleset_id/component/:id/validate
// Nodes godoc
// @Summary validate a component with component map
// @Schemes
// @Description validate a component with component map
// @Tags components
// @Accept json
// @Produce json
// @Param creds body shared_types.ComponentMap true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} map[string]string "Empty Object"
// @Router /api/:ruleset_id/component/:id/validate [post]
func ValidateComponent(ctx *gin.Context) {
	ruleset_id := ctx.Param("ruleset_id")
	id := ctx.Param("id")

	var body shared_types.ComponentMap
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Invalid body"})
		return
	}

	var verif shared_types.ComponentVerification
	var raw json.RawMessage
	if err := database.Db.QueryRow(`
		SELECT jsonb_build_object(
			'component_id', $2,
			'base_name', base_name,
			'scripts', scripts,
			'restrictions', restrictions
		)
		FROM get_verifications(ARRAY[(
			SELECT type_id
			FROM static_components_templates
			WHERE ruleset_id = $1 AND id = $2
		)]::uuid[])
	`, ruleset_id, id).Scan(&raw); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := json.Unmarshal(raw, &verif); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := type_validation.VerifyType(ctx, &body, body[id], verif); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

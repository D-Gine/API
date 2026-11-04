package types

import (
	"api/src/controllers/shared_types"
	typegeneration "api/src/controllers/type_generation"
	"api/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func generateTypesByIds(
	ctx *gin.Context,
	toGenerate []string,
	components shared_types.ComponentMap,
) (shared_types.TypeValueMap, int, error) {
	rows, err := database.Db.Query(`
			SELECT
				to_generate.id,
				tg.script,
				tg.type_name,
				tg.restrictions
			FROM unnest($1::uuid[]) AS to_generate(id)
			JOIN type_generation(to_generate.id) tg ON 1 = 1
	`, pq.Array(toGenerate))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()

	return typegeneration.GenerateRowTypes(ctx, components, rows)
}

// @BasePath /api/:ruleset_id/type/:id/generate
// Nodes godoc
// @Summary generate a type value with component map
// @Schemes
// @Description generate a type value with component map
// @Tags type
// @Accept json
// @Produce json
// @Param creds body shared_types.ComponentMap true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} map[string]object "generated components"
// @Router /api/:ruleset_id/type/:id/generate [post]
func GenerateType(ctx *gin.Context) {
	// ruleset_id := ctx.Param("ruleset_id")
	id := ctx.Param("id")

	var body shared_types.ComponentMap
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Invalid body"})
		return
	}
	generated, code, err := generateTypesByIds(ctx, []string{id}, body)

	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(code, generated)
}

// @BasePath /api/:ruleset_id/type/generate
// Nodes godoc
// @Summary generate multiple type values with component map
// @Schemes
// @Description generate multiple type values with component map
// @Tags type
// @Accept json
// @Produce json
// @Param creds body shared_types.GenerationMultipleBody true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} map[string]object "generated components"
// @Router /api/:ruleset_id/type/generate [post]
func GenerateTypes(ctx *gin.Context) {
	// ruleset_id := ctx.Param("ruleset_id")
	var body shared_types.GenerationMultipleBody
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Invalid body"})
		return
	}
	generated, code, err := generateTypesByIds(ctx, body.Ids, body.Components)

	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(code, generated)
}

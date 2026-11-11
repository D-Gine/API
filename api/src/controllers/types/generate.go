package types

import (
	"api/src/controllers/shared_types"
	"api/src/database"
	"net/http"

	gdt "github.com/D-Gine/Gestionnaire2Type"
	typegeneration "github.com/D-Gine/Gestionnaire2Type/type_generation"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func generateTypesByIds(
	toGenerate []string,
	components gdt.ComponentMap,
) (gdt.TypeValueMap, int, error) {
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
	defer database.CloseRows(rows)

	generated, err := typegeneration.GenerateRowTypes(database.Db, components, rows)
	return generated, shared_types.GenerationErrorCode(err), err
}

// @BasePath /api/:ruleset_id/type/:id/generate
// Nodes godoc
// @Summary Types - generate one type
// @Schemes
// @Description generate a type value with component map
// @Tags types
// @Accept json
// @Produce json
// @Param creds body gdt.ComponentMap true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} map[string]object "generated components"
// @Router /api/:ruleset_id/type/:id/generate [post]
func GenerateType(ctx *gin.Context) {
	// ruleset_id := ctx.Param("ruleset_id")
	id := ctx.Param("id")

	var body gdt.ComponentMap
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Invalid body"})
		return
	}
	generated, code, err := generateTypesByIds([]string{id}, body)

	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(code, generated)
}

// @BasePath /api/:ruleset_id/type/generate
// Nodes godoc
// @Summary Types - generate multiple type
// @Schemes
// @Description generate multiple type values with component map
// @Tags types
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
	generated, code, err := generateTypesByIds(body.Ids, body.Components)

	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(code, generated)
}

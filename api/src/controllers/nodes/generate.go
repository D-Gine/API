package nodes

import (
	"api/src/controllers/shared_types"
	"api/src/database"
	"net/http"

	gdt "github.com/D-Gine/Gestionnaire2Type"
	typegeneration "github.com/D-Gine/Gestionnaire2Type/type_generation"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func generateNodesByIds(
	toGenerate []string,
	components gdt.ComponentMap,
) (gdt.NodeMap, int, error) {
	rows, err := database.Db.Query(`
			SELECT
				n.id,
				tg.script,
				tg.type_name,
				tg.restrictions
			FROM nodes n
			JOIN static_components_templates sct ON n.template_id = sct.id
			JOIN type_generation(sct.type_id) tg ON 1 = 1
			WHERE n.id = ANY($1)
	`, pq.Array(toGenerate))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer database.CloseRows(rows)

	generated, err := typegeneration.GenerateRowTypes(database.Db, components, rows)
	return generated, shared_types.GenerationErrorCode(err), err
}

// @BasePath /api/:ruleset_id/node/:id/generate
// Nodes godoc
// @Summary Nodes - generate one
// @Schemes
// @Description generate a type value with component map
// @Tags nodes
// @Accept json
// @Produce json
// @Param creds body gdt.ComponentMap true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} map[string]object "generated components"
// @Router /api/:ruleset_id/component/:id/generate [post]
func GenerateNode(ctx *gin.Context) {
	// ruleset_id := ctx.Param("ruleset_id")
	id := ctx.Param("id")

	var body gdt.ComponentMap
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Invalid body"})
		return
	}
	generated, code, err := generateNodesByIds([]string{id}, body)

	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(code, generated)
}

// @BasePath /api/:ruleset_id/node/generate
// Nodes godoc
// @Summary Nodes - generate multiple
// @Schemes
// @Description generate multiple type values with component map
// @Tags nodes
// @Accept json
// @Produce json
// @Param creds body shared_types.GenerationMultipleBody true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} map[string]object "generated components"
// @Router /api/:ruleset_id/node/generate [post]
func GenerateNodes(ctx *gin.Context) {
	// ruleset_id := ctx.Param("ruleset_id")
	var body shared_types.GenerationMultipleBody
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Invalid body"})
		return
	}
	generated, code, err := generateNodesByIds(body.Ids, body.Components)

	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(code, generated)
}

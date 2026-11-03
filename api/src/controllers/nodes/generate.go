package nodes

import (
	"api/src/controllers/shared_types"
	typegeneration "api/src/controllers/type_generation"
	"api/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func generateNodesByIds(
	ctx *gin.Context,
	toGenerate []string,
	components shared_types.ComponentMap,
) (shared_types.NodeMap, int, error) {
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
	defer rows.Close()

	generated := make(shared_types.NodeMap)
	for rows.Next() {
		var result typegeneration.GenerationInfos
		var node_id string
		if err = rows.Scan(&node_id, &result.Script, &result.TypeName, &result.Restrictions); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		value, code, err := typegeneration.GenerateKnownedType(ctx, components, result)
		if err != nil {
			return nil, code, err
		}
		generated[node_id] = value
	}
	if err = rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return generated, http.StatusOK, nil
}

// @BasePath /api/:ruleset_id/node/:id/generate
// Nodes godoc
// @Summary generate a type value with component map
// @Schemes
// @Description generate a type value with component map
// @Tags nodes
// @Accept json
// @Produce json
// @Param creds body shared_types.ComponentMap true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} map[string]object "generated components"
// @Router /api/:ruleset_id/component/:id/generate [post]
func GenerateNode(ctx *gin.Context) {
	// ruleset_id := ctx.Param("ruleset_id")
	id := ctx.Param("id")

	var body shared_types.ComponentMap
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Invalid body"})
		return
	}
	generated, code, err := generateNodesByIds(ctx, []string{id}, body)

	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(code, generated)
}

// @BasePath /api/:ruleset_id/node/generate
// Nodes godoc
// @Summary generate multiple type values with component map
// @Schemes
// @Description generate multiple type values with component map
// @Tags node
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
	generated, code, err := generateNodesByIds(ctx, body.Ids, body.Components)

	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(code, generated)
}

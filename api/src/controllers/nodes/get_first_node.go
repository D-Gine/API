package nodes

import (
	"api/src/database"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/nodes
// Node godoc
// @Summary Reads first node for a ruleset
// @Schemes
// @Description Returns the first node of the given ruleset with resolved type
// @Tags nodes
// @Param ruleset_id path string true "id of the ruleset"
// @Produce json
// @Success 200 {object} FirstNodeResponse
// @Router /api/nodes/{ruleset_id}/first_node [get]
func ReadRulesetFirstNode(c *gin.Context) {
	rulesetID := c.Param("ruleset_id")
	if rulesetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ruleset_id"})
		return
	}

	var out FirstNodeResponse
	var typeRaw []byte

	err := database.Db.QueryRow(`
		SELECT
			fn.node_id::text,
			COALESCE(get_base_type_json(sct.type_id), '{}'::jsonb)::text
		FROM public.first_nodes fn
		INNER JOIN public.nodes n
			ON n.id = fn.node_id
			AND n.ruleset_id = fn.ruleset_id
		INNER JOIN public.static_components_templates sct
			ON sct.id = n.template_id
		WHERE fn.ruleset_id = $1
		LIMIT 1
	`, rulesetID).Scan(&out.NodeId, &typeRaw)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No first node found for this ruleset"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	if err := json.Unmarshal(typeRaw, &out.Type); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid type payload from database"})
		return
	}
	if err := resolveTypeDetails(&out.Type, rulesetID, 0); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot resolve type: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

package nodes

import (
	"api/src/database"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/nodes
// Node godoc
// @Summary Gets next node info for a target node
// @Schemes
// @Description Returns next node metadata for a given target node
// @Tags nodes
// @Accept json
// @Produce json
// @Param ruleset_id path string true "id of the ruleset"
// @Param target_node path string true "id of the target node"
// @Param payload body GetNextNodeRequest true "node values and current tree"
// @Success 200 {object} GetNextNodeResponse
// @Router /api/nodes/{ruleset_id}/get_next_node/{target_node} [post]
func GetNextNode(c *gin.Context) {
	rulesetID := c.Param("ruleset_id")
	targetNodeID := c.Param("target_node")
	if rulesetID == "" || targetNodeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing path params"})
		return
	}

	var body GetNextNodeRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	var nextScriptID sql.NullString
	var payload []byte

	err := database.Db.QueryRow(`
		SELECT
			n.next_script_id::text,
			jsonb_build_object(
				'next', jsonb_build_object(
					'node_id', n.id::text,
					'final', true,
					'type', COALESCE(get_base_type_json(sct.type_id), '{}'::jsonb)
				)
			)::text
		FROM public.nodes n
		INNER JOIN public.static_components_templates sct
			ON sct.id = n.template_id
		WHERE n.ruleset_id = $1 AND n.id = $2
		LIMIT 1
	`, rulesetID, targetNodeID).Scan(&nextScriptID, &payload)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Target node not found for this ruleset"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// v1 behavior
	if nextScriptID.Valid && nextScriptID.String != "" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "next script execution is not implemented yet"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", payload)
}

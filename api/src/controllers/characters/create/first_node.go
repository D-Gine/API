/*
** D&GINE Project, 2026
** API
** File description:
** characters/creation/first_node.go
 */

package charactersCreate

import (
	"database/sql"
	"net/http"

	"api/src/database"

	"github.com/gin-gonic/gin"
)

type FirstNodeResponse struct {
	NodeId         string              `json:"node_id"`
	Name           string              `json:"name"`
	Last           bool                `json:"last"`
	ComponentTypes []ComponentTemplate `json:"component_types"`
}

// @BasePath /api/characters/create/firstnode
// Characters godoc
// @Summary Returns the first node of the given ruleset
// @Schemes
// @Description Returns the first node of the given ruleset character creation tree with enriched metadata for enum_tag types
// @Tags characters creation
// @Accept json
// @Produce json
// @Param ruleset_id query string true "Ruleset ID"
// @Success 200 {object} FirstNodeResponse
// @Router /api/characters/create/firstnode [get]
func CharacterFirstNode(c *gin.Context) {
	rulesetId, ok := c.GetQuery("ruleset_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: ruleset_id is required"})
		return
	}

	// Find the first node (components_creations with no parent)
	var nodeId string
	err := database.Db.QueryRow(`
		SELECT DISTINCT cc.id
		FROM games.components_creations cc
		INNER JOIN games.creations_links cl ON cl.node_id = cc.id
		WHERE cc.ruleset_id = $1
		AND cl.parent_id IS NULL
		LIMIT 1
	`, rulesetId).Scan(&nodeId)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Given ruleset has no first node set"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// Get templates with enriched metadata (automatically adds possible_values for enum_tag)
	templates, err := getNodeTemplates(nodeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get node templates: " + err.Error()})
		return
	}

	// Check if this node has children
	hasChild, err := hasChildren(nodeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check children: " + err.Error()})
		return
	}

	// Return response
	response := FirstNodeResponse{
		NodeId:         nodeId,
		Name:           "Character Creation",
		Last:           !hasChild,
		ComponentTypes: templates,
	}

	c.JSON(http.StatusOK, response)
}

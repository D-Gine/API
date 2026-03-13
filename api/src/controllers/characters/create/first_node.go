/*
** D&GINE Project, 2026
** API
** File description:
** auth/register.go
 */

package charactersCreate

import (
	"encoding/json"
	"net/http"

	"api/src/database"

	"github.com/gin-gonic/gin"
)

type ComponentType struct {
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Metadata json.RawMessage `json:"metadata"`
}

type RulesetFirstNode struct {
	NodeId         string          `json:"node_id"`
	ComponentId    string          `json:"component_id"`
	Name           string          `json:"name"`
	Last           bool            `json:"last"`
	ComponentTypes []ComponentType `json:"component_types"`
}

// @BasePath /api/characters/create/firstnode
// Characters godoc
// @Summary Returns the first node of the given ruleset
// @Schemes
// @Description Returns the first node of the given ruleset character creation tree
// @Tags characters creation
// @Accept json
// @Produce json
// @Param creds body FirstNodeArgs true "character related informations that will be later needed for the login process"
// @Success 200 {object} RulesetFirstNode
// @Router /api/characters/create/firstnode [get]
func CharacterFirstNode(c *gin.Context) {
	ruleset_id, ok := c.GetQuery("ruleset_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var nodeId, componentId, name string
	err := database.Db.QueryRow(`
	SELECT
		cc.id, ctemp.id, ctemp.name
	FROM games.components_creations cc
	INNER JOIN games.creations_links cl
		ON cl.node_id = cc.id
	INNER JOIN games.creations_templates ct
		ON ct.creation_id = cc.id
	INNER JOIN games.components_templates ctemp
		ON ctemp.id = ct.template_id
	WHERE cc.ruleset_id = $1
		AND cl.parent_id IS NULL`, ruleset_id).Scan(&nodeId, &componentId, &name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Given ruleset has no first node set"})
		return
	}

	var childExists bool
	_ = database.Db.QueryRow(`
	SELECT EXISTS(
		SELECT 1 FROM games.creations_links
		WHERE parent_id = $1
	)`, nodeId).Scan(&childExists)

	rows, err := database.Db.Query(`
	SELECT
		ctype.name, ctype.type, ctype.metadata
	FROM games.creations_templates ct
	INNER JOIN games.components_templates ctemp
		ON ctemp.id = ct.template_id
	INNER JOIN games.components_types ctype
		ON ctemp.type = ctype.id
	WHERE ct.creation_id = $1`, nodeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}
	defer rows.Close()

	var componentTypes []ComponentType
	for rows.Next() {
		var ct ComponentType
		var metadata []byte
		if err := rows.Scan(&ct.Name, &ct.Type, &metadata); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
			return
		}
		ct.Metadata = json.RawMessage(metadata)
		componentTypes = append(componentTypes, ct)
	}

	c.JSON(http.StatusOK, RulesetFirstNode{
		NodeId:         nodeId,
		ComponentId:    componentId,
		Name:           name,
		Last:           !childExists,
		ComponentTypes: componentTypes,
	})
}

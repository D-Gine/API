/*
** D&GINE Project, 2026
** API
** File description:
** auth/register.go
 */

package charactersCreate

import (
	"database/sql"
	"net/http"

	"api/src/database"

	"github.com/gin-gonic/gin"
)

type FirstNodeArgs struct {
	RulesetId string `json:"ruleset_id"`
}

type RulesetFirstNode struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Metadata string `json:"metadata"`
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
// @Success 201 {object} RulesetFirstNode
// @Router /api/characters/create/firstnode [get]
func CharacterFirstNode(c *gin.Context) {
	ruleset_id, ok := c.GetQuery("ruleset_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var id, name, valType, metadata sql.NullString
	err := database.Db.QueryRow(`
	SELECT
		ctemp.id, ctemp.name, ctype.type, ctype.metadata
	FROM games.components_templates ctemp
	INNER JOIN games.first_nodes fn
		ON ctemp.id=fn.first_node
	INNER JOIN games.component_type ctype
		ON ctemp.type=ctype.id
	WHERE fn.ruleset_id=$1`, ruleset_id).Scan(&id, &name, &valType, &metadata)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Given ruleset has no first node set"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error : " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, RulesetFirstNode{
		Id:       id.String,
		Name:     name.String,
		Type:     valType.String,
		Metadata: metadata.String,
	})
}

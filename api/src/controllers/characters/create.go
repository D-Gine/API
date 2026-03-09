/*
** D&GINE Project, 2026
** API
** File description:
** auth/register.go
 */

package characters

import (
	"database/sql"
	"net/http"

	"api/src/database"

	"github.com/gin-gonic/gin"
)

type CreateArgs struct {
	RulesetId string `json:"ruleset_id"`
}

type RulesetFirstNode struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Metadata string `json:"metadata"`
}

// @BasePath /api/characters
// Characters godoc
// @Summary Creates a character
// @Schemes
// @Description ⚠️ Only accessible to admins ⚠️<br><br>Creates a character with arguments in body<br><br>Will return the id of created character
// @Tags characters
// @Accept json
// @Produce json
// @Param creds body CreateArgs true "character related informations that will be later needed for the login process"
// @Success 201 {object} RulesetFirstNode
// @Router /api/characters [post]
func CreateCharacters(c *gin.Context) {
	var args CreateArgs

	// Body Json parsing
	if err := c.ShouldBindJSON(&args); err != nil {
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
	WHERE fn.ruleset_id=$1`, args.RulesetId).Scan(&id, &name, &valType, &metadata)
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

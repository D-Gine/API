/*
** D&GINE Project, 2026
** API
** File description:
** auth/register.go
 */

package charactersCreate

import (
	"api/src/database"
	"api/src/internal/structs"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateArgs struct {
	RulesetId  string `json:"ruleset_id"`
	Name       string `json:"name"`
	Components []struct {
		Values []struct {
			TemplateId string `json:"template_id"`
			Value      any    `json:"value"`
		} `json:"values"`
	} `json:"components"`
}

// @BasePath /api/characters/create/submit
// Characters godoc
// @Summary Submits a character list of components
// @Schemes
// @Description Submits a character list of components and creates the character in database<br><br>Will return the id of created character
// @Tags characters creation
// @Accept json
// @Produce json
// @Param creds body CreateArgs true "character related informations that will be later needed for the login process"
// @Success 201 {object} RulesetFirstNode
// @Router /api/characters/create/submit [post]
func SubmitCharacter(c *gin.Context) {
	var args CreateArgs

	// Body Json parsing
	if err := c.ShouldBindJSON(&args); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var entity_id sql.NullString
	err := database.Db.QueryRow("INSERT INTO games.entities (ruleset_id, owner_id, name) VALUES ($1, $2, $3) RETURNING id", args.RulesetId, "ebeb6c5b-1f6f-4ec6-a634-4c0b7c17a480", args.Name).Scan(&entity_id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to register ruleset"})
		return
	}
	for i := range len(args.Components) {
		var cmpt_id string
		valuesJSON, err := json.Marshal(args.Components[i].Values)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create component"})
			return
		}
		err = database.Db.QueryRow("INSERT INTO games.components (ruleset_id, name, value) VALUES ($1, $2, $3) RETURNING id",
			args.RulesetId, "component", valuesJSON).Scan(&cmpt_id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create component"})
			return
		}
		_, err = database.Db.Exec("INSERT INTO games.components_entities (entity_id, component_id) VALUES ($1, $2)",
			entity_id, cmpt_id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to link component to entity"})
			return
		}
		for j := range len(args.Components[i].Values) {
			_, err = database.Db.Exec("INSERT INTO games.components_templates_links (component_id, template_id) VALUES ($1, $2)",
				cmpt_id, args.Components[i].Values[j].TemplateId)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to link component to template"})
				return
			}
		}
	}
	c.JSON(http.StatusCreated, structs.PostResponse{Id: entity_id.String})
}

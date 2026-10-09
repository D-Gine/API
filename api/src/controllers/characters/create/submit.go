/*
** D&GINE Project, 2026
** API
** File description:
** characters/create/submit.go
 */

package charactersCreate

import (
	"api/src/database"
	"api/src/internal/structs"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ComponentValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type ComponentData struct {
	CompId string           `json:"comp id"` // creation node id
	Values []ComponentValue `json:"values"`
}

type CreateArgs struct {
	RulesetId  string          `json:"ruleset_id"`
	Name       string          `json:"name"`
	Components []ComponentData `json:"components"`
}

// @BasePath /api/characters/create/submit
// Characters godoc
// @Summary Submits a character list of components
// @Schemes
// @Description Submits a character list of components and creates the character in database<br><br>Will return the id of created character
// @Tags characters creation
// @Accept json
// @Produce json
// @Param creds body CreateArgs true "Character creation data with ruleset_id, name, and components"
// @Success 201 {object} structs.PostResponse
// @Router /api/characters/create/submit [post]
func SubmitCharacter(c *gin.Context) {
	var args CreateArgs

	if err := c.ShouldBindJSON(&args); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// TODO: Get user ID from authentication context instead of hardcoding
	userId := "ebeb6c5b-1f6f-4ec6-a634-4c0b7c17a480"

	var entityId string
	err := database.Db.QueryRow(`
		INSERT INTO games.entities (ruleset_id, owner_id, name)
		VALUES ($1, $2, $3)
		RETURNING id
	`, args.RulesetId, userId, args.Name).Scan(&entityId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entity: " + err.Error()})
		return
	}

	for _, component := range args.Components {
		if len(component.Values) == 0 {
			continue
		}
		valueMap := make(map[string]interface{})
		for _, val := range component.Values {
			valueMap[val.Key] = val.Value
		}
		valuesJSON, err := json.Marshal(valueMap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal component values: " + err.Error()})
			return
		}
		var componentName string
		err = database.Db.QueryRow(`
			SELECT name
			FROM games.components_templates
			WHERE id = $1
		`, component.Values[0].Key).Scan(&componentName)

		if err != nil {
			componentName = "Component"
		}
		var componentId string
		err = database.Db.QueryRow(`
			INSERT INTO games.components (ruleset_id, name, value)
			VALUES ($1, $2, $3)
			RETURNING id
		`, args.RulesetId, componentName, valuesJSON).Scan(&componentId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create component: " + err.Error()})
			return
		}

		_, err = database.Db.Exec(`
			INSERT INTO games.components_entities (entity_id, component_id)
			VALUES ($1, $2)
		`, entityId, componentId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link component to entity: " + err.Error()})
			return
		}
		for _, val := range component.Values {
			_, err = database.Db.Exec(`
				INSERT INTO games.components_templates_links (component_id, template_id)
				VALUES ($1, $2)
			`, componentId, val.Key)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link component to template: " + err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusCreated, structs.PostResponse{Id: entityId})
}

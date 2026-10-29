package characterGet

import (
	"api/src/controllers/characters"
	"api/src/database"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath api/:ruleset_id/character/:id
// Nodes godoc
// @Summary get character infos from the db
// @Schemes
// @Description
// @Tags character
// @Accept json
// @Produce json
// @Param creds body string true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} characterGet.GetResponse
// @Router /api/:ruleset_id/character/:id [get]
func GetOneCharacter(c *gin.Context) {
	character_id := c.Param("id")

	// check if character exists with the user

	var jsonData []byte
	err := database.Db.QueryRow(`
		SELECT jsonb_build_object(
			'id', e.id,
			'name', e.name,
			'description', e.description,
			'components', jsonb_build_object(
				'static_components', (
					SELECT jsonb_object_agg(
						sct.id, jsonb_build_object(
							'name', sct.name,
							'value', sc.value,
							'type', get_base_type_json(sct.type_id)
						)
					)
					FROM static_components sc
					JOIN static_components_templates sct ON sct.id = sc.template_id
					WHERE sc.entity_id = $1
				),
				'dynamic_components', (
					SELECT jsonb_object_agg(
						dct.id, jsonb_build_object(
							'name', dct.name,
							'value', dc.value,
							'type', get_base_type_json(dct.type_id)
						)
					)
					FROM dynamic_components dc
					JOIN dynamic_components_templates dct ON dct.id = dc.template_id
					WHERE dc.entity_id = $1
				)
			)
		)
		FROM entities e
		WHERE e.id = $1
	`, character_id).Scan(&jsonData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var res characters.Entity
	err = json.Unmarshal(jsonData, &res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

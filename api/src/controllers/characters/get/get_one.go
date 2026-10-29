package characterGet

import (
	"api/src/controllers/nodes"
	"api/src/database"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ComponentValue = map[string]struct {
	Value json.RawMessage        `json:"value"`
	Type  nodes.NodeTypeResponse `json:"type"`
}

type GetResponse struct {
	StaticComponents  ComponentValue `json:"static_components"`
	DynamicComponents ComponentValue `json:"dynamic_components"`
}

// @BasePath api/:ruleset_id/character/:id
// Nodes godoc
// @Summary get character infos from the db
// @Schemes
// @Description
// @Tags character
// @Accept json
// @Produce json
// @Param creds body string true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} structs.PostResponse
// @Router /api/:ruleset_id/character/:id [get]
func GetOneCharacter(c *gin.Context) {
	character_id := c.Param("id")

	// check if character exists with the user

	var jsonData []byte
	err := database.Db.QueryRow(`
		SELECT jsonb_build_object(
			'static_components', (
				SELECT jsonb_object_agg(
					sct.id, jsonb_build_object(
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
						'value', dc.value,
						'type', get_base_type_json(dct.type_id)
					)
				)
				FROM dynamic_components dc
				JOIN dynamic_components_templates dct ON dct.id = dc.template_id
				WHERE dc.entity_id = $1
			)
		)
	`, character_id).Scan(&jsonData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var res GetResponse
	err = json.Unmarshal(jsonData, &res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

package enums

import (
	"api/src/controllers"
	"api/src/database"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type EnumInfos struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// @BasePath api/:ruleset_id/enums/get_by_restrictions
// Nodes godoc
// @Summary get enums infos by restrictions
// @Schemes
// @Description get enums id, name and description matching the given restrictions
// @Tags enums
// @Accept json
// @Produce json
// @Param creds body CreateRulesetsArgs true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} structs.PostResponse
// @Router api/:ruleset_id/enums/get_by_restrictions [post]
func GetEnumsByRestrictions(c *gin.Context) {
	ruleset_id := c.Param("ruleset_id")
	var restrictions controllers.EnumRestrictions

	if err := c.ShouldBindBodyWithJSON(&restrictions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body: " + err.Error()})
		return
	}
	var response map[string]EnumInfos
	var rawJson json.RawMessage
	err := database.Db.QueryRow(`
		WITH pack_enums AS (
			SELECT
				l.enum_id as enum_id,
				ARRAY_AGG(l.tag_id) as tags
			FROM link_enums_tags l
			JOIN tags t ON l.tag_id = t.id
			WHERE
				t.ruleset_id = $1 AND
				tag_id = ANY($2::uuid[])
			GROUP BY
				enum_id
		)
		SELECT jsonb_object_agg(
			enums.id, jsonb_build_object(
				'name', enums.name,
				'description', enums.description
			)
		)
		FROM pack_enums p
		JOIN enums ON enums.id = p.enum_id
		WHERE ARRAY_LENGTH(p.tags, 1) = $3
	`, ruleset_id, pq.Array(restrictions.Tags), len(restrictions.Tags)).Scan(&rawJson)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = json.Unmarshal(rawJson, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

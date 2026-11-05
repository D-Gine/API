/*
** D&GINE Project, 2026
** API
** File description:
** characters/validate.go
 */

package characters

import (
	"api/src/controllers/shared_types"
	"api/src/controllers/type_validation"
	"api/src/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateDbCharacter(c *gin.Context, id string) (int, map[string]string, error) {
	var rawJSON []byte
	err := database.Db.QueryRow(`
		SELECT jsonb_object_agg(template_id, value)
		FROM static_components
		WHERE entity_id = $1
	`, id).Scan(&rawJSON)
	if err != nil && err != sql.ErrNoRows {
		return http.StatusInternalServerError, nil, err
	}
	componentList := make(shared_types.ComponentMap)
	if len(rawJSON) > 0 {
		if err := json.Unmarshal(rawJSON, &componentList); err != nil {
			return http.StatusInternalServerError, nil, fmt.Errorf("failed to parse components: %w", err)
		}
	}

	verifyRows, err := database.Db.Query(`
		SELECT * FROM get_entity_verifications($1)
	`, id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	defer database.CloseRows(verifyRows)

	error_map, _, valErr := type_validation.ValidateRowComponents(c, &componentList, verifyRows)
	if valErr != nil {
		return http.StatusBadRequest, nil, valErr
	}

	if err = verifyRows.Err(); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	code := http.StatusOK
	if len(error_map) != 0 {
		code = http.StatusUnauthorized
	}
	return code, error_map, nil
}

// @BasePath api/:ruleset_id/character/:id/validate
// Nodes godoc
// @Summary Characters - Validate entity
// @Schemes
// @Description Validate an entity stored in the database by checking every component
// @Tags characters
// @Accept json
// @Produce json
// @Param creds body string true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} map[string]string "errors"
// @Router /api/:ruleset_id/character/:id/validate [post]
func ValidateCharacter(c *gin.Context) {
	id := c.Param("id")
	code, error_map, err := ValidateDbCharacter(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(code, error_map)
}

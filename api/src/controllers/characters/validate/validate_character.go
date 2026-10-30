package validate

import (
	"api/src/controllers"
	"api/src/controllers/nodes/validate"
	"api/src/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateDbCharacter(c *gin.Context, id string) (int, error, map[string]string) {
	var rawJSON []byte
	err := database.Db.QueryRow(`
		SELECT jsonb_object_agg(template_id, value)
		FROM static_components
		WHERE entity_id = $1
	`, id).Scan(&rawJSON)
	if err != nil && err != sql.ErrNoRows {
		return http.StatusInternalServerError, err, nil
	}
	componentList := make(controllers.ComponentMap)
	if len(rawJSON) > 0 {
		if err := json.Unmarshal(rawJSON, &componentList); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("Failed to parse components: %w", err), nil
		}
	}

	verifyRows, err := database.Db.Query(`
		SELECT * FROM get_entity_verifications($1)
	`, id)
	if err != nil {
		return http.StatusInternalServerError, err, nil
	}
	defer func() { _ = verifyRows.Close() }()

	error_map, _, err := validate.ValidateRowTypes(c, &componentList, verifyRows)
	if err != nil {
		return http.StatusBadRequest, err, nil
	}

	if err = verifyRows.Err(); err != nil {
		return http.StatusInternalServerError, err, nil
	}

	code := http.StatusOK
	if len(error_map) != 0 {
		code = http.StatusUnauthorized
	}
	return code, nil, error_map
}

// @BasePath api/:ruleset_id/character/:id/validate
// Nodes godoc
// @Summary validate an entity
// @Schemes
// @Description Validate an entity stored in the database by checking every component
// @Tags character
// @Accept json
// @Produce json
// @Param creds body string true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} map[string]string "errors"
// @Router /api/:ruleset_id/character/:id/validate [post]
func ValidateCharacter(c *gin.Context) {
	id := c.Param("id")
	code, err, error_map := ValidateDbCharacter(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(code, error_map)
}

package validate

import (
	"api/src/controllers/nodes/validate"
	"api/src/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateDbCharacter(c *gin.Context, id string) (int, error) {
	var rawJSON []byte
	err := database.Db.QueryRow(`
		SELECT jsonb_object_agg(template_id, value)
		FROM static_components
		WHERE entity_id = $1
	`, id).Scan(&rawJSON)
	if err != nil && err != sql.ErrNoRows {
		return http.StatusInternalServerError, err
	}
	componentList := make(validate.ComponentMap)
	if len(rawJSON) > 0 {
		if err := json.Unmarshal(rawJSON, &componentList); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("Failed to parse components: %w", err)
		}
	}

	verifyRows, err := database.Db.Query(`
		SELECT * FROM get_entity_verifications($1)
	`, id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer verifyRows.Close()

	valErr, error_map, _ := validate.ValidateRowTypes(c, &componentList, verifyRows)
	if valErr != nil {
		return http.StatusBadRequest, valErr
	}

	if err = verifyRows.Err(); err != nil {
		return http.StatusInternalServerError, err
	}

	code := http.StatusOK
	if len(error_map) != 0 {
		code = http.StatusUnauthorized
	}
	c.JSON(code, error_map)
	return code, nil
}

// @BasePath api/:ruleset_id/character/:id/validate
// Nodes godoc
// @Summary validate an entity
// @Schemes
// @Description Validate an entity stored in the database by checking every component
// @Tags nodes
// @Accept json
// @Produce json
// @Param creds body CreateRulesetsArgs true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} structs.PostResponse
// @Router api/:ruleset_id/character/:id/validate [post]
func ValidateCharacter(c *gin.Context) {
	id := c.Param("id")
	if code, err := ValidateDbCharacter(c, id); err != nil {
		c.JSON(code, gin.H{"error": err.Error()})
	}
}

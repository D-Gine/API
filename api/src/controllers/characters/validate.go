/*
** D&GINE Project, 2026
** API
** File description:
** characters/validate.go
 */

package characters

import (
	"api/src/database"
	"net/http"

	"github.com/D-Gine/Gestionnaire2Type/validate"
	"github.com/gin-gonic/gin"
)

func ValidateDbCharacter(id string) (int, map[string]string, error) {
	error_map, err := validate.ValidateEntity(database.Db, id)
	if err != nil {
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
	code, error_map, err := ValidateDbCharacter(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(code, error_map)
}

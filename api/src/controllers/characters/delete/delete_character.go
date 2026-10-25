package delete

import (
	"api/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath api/:ruleset_id/character/:id
// Nodes godoc
// @Summary Delete a character in the db
// @Schemes
// @Description Delete an empty character, returning his ID
// @Tags character
// @Accept json
// @Produce json
// @Param creds body DeleteRulesetsArgs true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} structs.PostResponse
// @Router api/:ruleset_id/character/:id [put]
func DeleteCharater(c *gin.Context) {
	character_id := c.Param("id")

	// check if character exists with the user

	_, err := database.Db.Exec(`
		DELETE FROM entities
		WHERE id = $1
		`, character_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

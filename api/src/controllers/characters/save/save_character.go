package save

import (
	"api/src/controllers/nodes/validate"
	"api/src/database"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// // @BasePath api/:ruleset_id/character/:id/save
// // Nodes godoc
// // @Summary save a character to db
// // @Schemes
// // @Description save the static components into the database.
// // @Tags character
// // @Accept json
// // @Produce json
// // @Param creds body CreateRulesetsArgs true "ruleset related informations that will be later needed for the login process"
// // @Success 200 {object} structs.PostResponse
// // @Router api/:ruleset_id/character/:id/save [put]
func SaveCharacter(c *gin.Context) {
	id := c.Param("id")
	var body validate.NodeList

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}
	jsonValues, parseErr := json.Marshal(body.NodeValues)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse node values"})
		return
	}

	_, err := database.Db.Exec(`SELECT save_character($1, $2::JSONB)`, id, jsonValues)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

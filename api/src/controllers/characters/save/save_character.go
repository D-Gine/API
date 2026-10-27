package save

import (
	"api/src/controllers/nodes/validate"
	"api/src/database"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WriteCharacterInDb(character_id string, nodes validate.NodeList) (int, error) {
	jsonValues, parseErr := json.Marshal(nodes.NodeValues)
	if parseErr != nil {
		return http.StatusBadRequest, errors.New("Failed to parse node values")
	}
	if _, err := database.Db.Exec(`SELECT save_character($1, $2::JSONB)`, character_id, jsonValues); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

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
	code, err := WriteCharacterInDb(id, body)
	if err != nil {
		c.JSON(code, gin.H{"error": err.Error()})
	} else {
		c.JSON(code, gin.H{})
	}
}

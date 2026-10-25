package submit

import (
	"api/src/controllers/characters/save"
	"api/src/controllers/characters/validate"
	nodeValidate "api/src/controllers/nodes/validate"
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
func SumbitCharacter(c *gin.Context) {
	id := c.Param("id")
	var body nodeValidate.NodeList

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}
	if code, err := save.WriteCharacterInDb(id, body); err != nil {
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	if code, err := validate.ValidateDbCharacter(c, id); code != http.StatusOK {
		if code != http.StatusUnauthorized {
			c.JSON(code, gin.H{"error": err.Error()})
		}
		return
	}
	// create the static components related to the character
}

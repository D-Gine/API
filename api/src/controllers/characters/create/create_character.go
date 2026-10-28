package charactersCreate

import (
	"api/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateCharacterBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// @BasePath api/:ruleset_id/character/create
// Nodes godoc
// @Summary create a character in the db
// @Schemes
// @Description create an empty character, returning his ID
// @Tags character
// @Accept json
// @Produce json
// @Param creds body CreateRulesetsArgs true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} structs.PostResponse
// @Router api/:ruleset_id/character/create [put]
func CreateCharater(c *gin.Context) {
	// user_id := get_user_id(c)
	ruleset := c.Param("ruleset_id")
	var characte_id string

	var body CreateCharacterBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	err := database.Db.QueryRow(`
		INSERT INTO entities (ruleset_id, name, description)
		VALUES ($1, $2, $3)
		RETURNING id
		`, ruleset, body.Name, body.Description).Scan(&characte_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// insert int the character associative table user -> character id

	c.JSON(http.StatusCreated, gin.H{"characharacte_id": characte_id})
}

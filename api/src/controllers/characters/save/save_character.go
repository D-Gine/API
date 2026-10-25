// package save

// import (
// 	"api/src/database"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// // @BasePath api/:ruleset_id/character/:id/save
// // Nodes godoc
// // @Summary save a character to db
// // @Schemes
// // @Description save the static components into the database.
// // then check if character if valid, if he is, generate the dynamic components
// // @Tags character
// // @Accept json
// // @Produce json
// // @Param creds body CreateRulesetsArgs true "ruleset related informations that will be later needed for the login process"
// // @Success 200 {object} structs.PostResponse
// // @Router api/:ruleset_id/character/:id/save [put]
// func ValidateCharacter(c *gin.Context) {
// 	ruleset_id := c.Param("ruleset_id")
// 	character_id := ""
// 	id := c.Param("id")
// 	var body CreateCharacterBody
// 	if err := c.ShouldBindJSON(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
// 		return
// 	}

// 	database.Db.QueryRow("INSERT")
// }

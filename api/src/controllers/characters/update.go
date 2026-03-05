/*
** D&GINE Project, 2026
** Backend
** File description:
** user/update.go
 */

package characters

import (
	"api/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateArgs struct {
	PlayerId string `json:"player_id"`
	Name     string `json:"name"`
}

// @BasePath /api/characters
// Characters godoc
// @Summary Updates given ID character with arguments in body
// @Schemes
// @Description ⚠️ Only accessible to admins ⚠️<br><br>Updates the given character's account data
// @Tags characters
// @Accept json
// @Param id query string true "id of the character to update"
// @Param data body UpdateArgs true "New arguments to be set in the character given"
// @Success 200
// @Router /api/characters/id [put]
func UpdateCharacters(c *gin.Context) {
	idQuery := c.Request.URL.Query().Get("id")
	if idQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong parameters"})
		return
	}

	var args UpdateArgs

	// Parsing des args
	err := c.ShouldBindJSON(&args)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err = database.Db.Exec(
		`UPDATE games.characters
		SET player_id=$2, name=$3
		WHERE id=$1`,
		idQuery, args.PlayerId, args.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

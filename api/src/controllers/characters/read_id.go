/*
** D&GINE Project, 2026
** Backend
** File description:
** user/read.go
 */

package characters

import (
	"api/src/database"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/characters
// Characters godoc
// @Summary Reads given character data
// @Schemes
// @Description Reads the given character data
// @Tags characters
// @Param id query string true "id of the character to read"
// @Produce json
// @Success 200 {object} RowReadResponse
// @Router /api/characters/id [get]
func ReadCharactersId(c *gin.Context) {
	idQuery := c.Request.URL.Query().Get("id")
	if idQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong parameters"})
		return
	}

	var id, playerId, name sql.NullString
	err := database.Db.QueryRow("SELECT id, player_id, name FROM games.characters WHERE id=$1", idQuery).Scan(&id, &playerId, &name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error : " + err.Error()})
	}

	c.JSON(http.StatusOK, RowReadResponse{
		Id:       id.String,
		PlayerId: playerId.String,
		Name:     name.String,
	})

}

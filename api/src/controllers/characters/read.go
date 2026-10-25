/*
** D&GINE Project, 2026
** API
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
// @Summary Reads all characters data
// @Schemes
// @Description Reads all characters data
// @Tags characters
// @Produce json
// @Success 200 {object} []RowReadResponse
// @Router /api/characters [get]
func ReadCharacters(c *gin.Context) {
	result := []RowReadResponse{}
	rows, err := database.Db.Query("SELECT id, player_id, name FROM games.characters")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var id, playerId, name sql.NullString
		if err := rows.Scan(&id, &playerId, &name); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		result = append(result, RowReadResponse{
			Id:       id.String,
			PlayerId: playerId.String,
			Name:     name.String,
		})
	}

	c.JSON(http.StatusOK, result)
}

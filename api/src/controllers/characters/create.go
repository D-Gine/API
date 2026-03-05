/*
** D&GINE Project, 2026
** Backend
** File description:
** auth/register.go
 */

package characters

import (
	"database/sql"
	"net/http"

	"api/src/database"
	"api/src/internal/structs"

	"github.com/gin-gonic/gin"
)

type CreateArgs struct {
	PlayerId string `json:"player_id"`
	Name     string `json:"name"`
}

// @BasePath /api/characters
// Characters godoc
// @Summary Connection of a new character to a new account
// @Schemes
// @Description ⚠️ Only accessible to admins ⚠️<br><br>Creates a character account with given informations
// @Tags characters
// @Accept json
// @Produce json
// @Param creds body CreateArgs true "character related informations that will be later needed for the login process"
// @Success 200 {object} structs.PostResponse
// @Router /api/characters [post]
func CreateCharacters(c *gin.Context) {
	var args CreateArgs

	// Body Json parsing
	if err := c.ShouldBindJSON(&args); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var id sql.NullString
	err := database.Db.QueryRow("INSERT INTO games.characters (player_id, name) VALUES ($1, $2) RETURNING id", args.PlayerId, args.Name).Scan(&id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, structs.PostResponse{Id: id.String})
}

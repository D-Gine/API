/*
** D&GINE Project, 2026
** Backend
** File description:
** user/delete.go
 */

package characters

import (
	"api/src/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/characters
// Characters godoc
// @Summary Deletes given ID character
// @Schemes
// @Param id query string true "id of the character to do the action on"
// @Description Deletes the given character's data<br> <b>Careful, there is no turn back or check, once called, this route will delete the character no matter what</b> <br> It also deletes all content related to the character
// @Tags characters
// @Success 200
// @Router /api/characters/id [delete]
func DeleteCharacters(c *gin.Context) {
	idQuery := c.Request.URL.Query().Get("id")
	if idQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong parameters"})
		return
	}

	_, err := database.Db.Exec("DELETE FROM games.characters WHERE id=$1", idQuery)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

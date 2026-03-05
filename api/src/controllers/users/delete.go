/*
** D&GINE Project, 2026
** Backend
** File description:
** users/delete.go
 */

package users

import (
	"api/src/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/users
// Users godoc
// @Summary Deletes given user
// @Schemes
// @Param id query string true "id of the user to do the action on"
// @Description ⚠️ Only accessible to admins ⚠️<br><br>Deletes the given user<br> <b>Careful, there is no turn back or check, once called, this route will delete the user no matter what</b> <br> It also deletes all content related to the user
// @Tags users
// @Success 200
// @Router /api/users/id [delete]
func DeleteUsers(c *gin.Context) {
	idQuery := c.Request.URL.Query().Get("id")
	if idQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong parameters"})
		return
	}

	_, err := database.Db.Exec("DELETE FROM accounts.users WHERE id=$1", idQuery)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

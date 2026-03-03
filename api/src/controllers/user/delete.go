/*
** D&GINE Project, 2026
** Backend
** File description:
** user/delete.go
 */

package user

import (
	"api/src/controllers/auth"
	"api/src/database"
	"api/src/internal/domains"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/user
// User godoc
// @Summary Deletes the actual user
// @Schemes
// @Description Deletes the actual user <br> <b>Careful, there is no turn back or check, once called, this route will delete the user no matter what</b> <br> It also deletes all content related to the user <br><br><b>⚠️ The user must be logged in<b>
// @Tags user
// @Success 200
// @Router /api/user [delete]
func DeleteUser(c *gin.Context) {
	// Recuperation du user depuis le token
	user, err := auth.GetUserFromToken(c)
	if err != nil {
		return
	}

	_, err = database.Db.Exec("DELETE FROM accounts.users WHERE id=$1", user.Id)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.SetCookie("token", "", -1, "/", domains.TokenDomain, false, true)
	c.JSON(http.StatusOK, nil)
}

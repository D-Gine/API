/*
** D&GINE Project, 2026
** API
** File description:
** account/delete.go
 */

package account

import (
	"api/src/controllers/auth"
	"api/src/database"
	"api/src/internal/config"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/account
// Account godoc
// @Summary Deletes the user's account
// @Schemes
// @Description <b>⚠️ The user must be logged in ⚠️</b><br><br>Deletes the user's account from database, the user will also be logged out<br> <b>Careful, there is no turn back or check, once called, this route will delete the user no matter what</b> <br> It also deletes all content related to the user
// @Tags account
// @Success 204
// @Router /api/account [delete]
func DeleteAccount(c *gin.Context) {
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
	if err := os.Remove(user.Image); err != nil && !os.IsNotExist(err) {
		fmt.Println("Failed to delete user image:", err)
	}
	c.SetCookie("token", "", -1, "/", config.TokenDomain, false, true)
	c.JSON(http.StatusNoContent, nil)
}

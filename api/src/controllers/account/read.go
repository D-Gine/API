/*
** D&GINE Project, 2026
** Backend
** File description:
** account/read.go
 */

package account

import (
	"api/src/controllers/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserReadResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// @BasePath /api/account
// Account godoc
// @Summary Reads the user's account data
// @Schemes
// @Description <b>⚠️ The user must be logged in ⚠️</b><br><br>Reads the user's account data
// @Tags account
// @Produce json
// @Success 200 {object} UserReadResponse
// @Router /api/account [get]
func ReadAccount(c *gin.Context) {
	// Recuperation du user depuis le token
	user, err := auth.GetUserFromToken(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, UserReadResponse{user.Name, user.Email, user.Role})
}

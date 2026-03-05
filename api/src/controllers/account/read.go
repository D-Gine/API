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
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"image"`
	Role  string `json:"role"`
}

// @BasePath /api/account
// Account godoc
// @Summary Reads the user's account data
// @Schemes
// @Description <b>⚠️ The user must be logged in ⚠️</b><br><br>Reads the user's account data<br><br>The image is a HTTP path to the user's profile picture using the same host:port as the API
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

	c.JSON(http.StatusOK, UserReadResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Image: user.Image,
		Role:  user.Role,
	})
}

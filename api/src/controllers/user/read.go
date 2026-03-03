/*
** D&GINE Project, 2026
** Backend
** File description:
** user/read.go
 */

package user

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

// @BasePath /api/user
// User godoc
// @Summary Reads the actual user data
// @Schemes
// @Description Reads the actual user data <br><br><b>⚠️ The user must be logged in<b>
// @Tags user
// @Produce json
// @Success 200 {object} UserReadResponse
// @Router /api/user [get]
func ReadUser(c *gin.Context) {
	// Recuperation du user depuis le token
	user, err := auth.GetUserFromToken(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, UserReadResponse{user.Name, user.Email, user.Role})
}

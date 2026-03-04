/*
** D&GINE Project, 2026
** Backend
** File description:
** auth/logout.go
 */

package auth

import (
	"api/src/internal/domains"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/auth/logout
// Auth godoc
// @Summary Disconnection of the user
// @Schemes
// @Description Disconnection of the user <br>The token in the Cookies will be deleted after this route has been called
// @Tags auth
// @Success 200
// @Router /api/auth/logout [post]
func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", domains.TokenDomain, false, true)
	c.JSON(http.StatusOK, nil)
}

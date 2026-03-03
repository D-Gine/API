/*
** D&GINE Project, 2026
** Backend
** File description:
** internal/endpoints/about.go
 */

package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type About struct {
	Status string `json:"status"`
}

// About godoc
// @BasePath /about.json
// @Summary Gives informations about api
// @Schemes
// @Description Gives informations about api
// @Tags about.json
// @Accept json
// @Produce json
// @Success 200 {object} About
// @Router /about.json [get]
func GetAbout(c *gin.Context) {
	c.JSON(http.StatusOK, About{Status: "D&GINE API is running"})
}

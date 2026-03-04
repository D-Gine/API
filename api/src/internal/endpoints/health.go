/*
** D&GINE Project, 2026
** Backend
** File description:
** internal/endpoints/endpoints.go
 */

package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/health
// Users godoc
// @Summary Healthcheck for the api
// @Schemes
// @Description Healthcheck used (mainly by docker) to check if the api is up and ready to respond
// @Tags health
// @Produce plain
// @Success 200 {string} OK
// @Router /api/health [get]
func GetHealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

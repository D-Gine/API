/*
** D&GINE Project, 2026
** API
** File description:
** internal/apidata/health.go
 */

package apidata

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/health
// Apidata godoc
// @Summary Healthcheck for the api
// @Schemes
// @Description Healthcheck used (mainly by docker) to check if the api is up and ready to respond
// @Tags apidata
// @Produce plain
// @Success 200 {string} OK
// @Router /api/health [get]
func GetHealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

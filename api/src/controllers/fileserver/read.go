/*
** D&GINE Project, 2026
** API
** File description:
** img.go
 */

package fileserver

import (
	_ "api/src/docs"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// @BasePath /img/*
// Image godoc
// @Summary Reads an image file
// @Schemes
// @Description Reads an image file
// @Tags img
// @Produce image/png
// @Success 200
// @Router /img/:filepath: [get]
func GetImage(c *gin.Context) {
	c.Request.URL.Path = c.Param("filepath")
	http.FileServer(http.Dir(os.Getenv("IMAGES_PATH"))).ServeHTTP(c.Writer, c.Request)
}

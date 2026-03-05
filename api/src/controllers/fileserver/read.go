/*
** D&GINE Project, 2026
** Backend
** File description:
** img.go
 */

package fileserver

import (
	_ "api/src/docs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetFile(c *gin.Context) {
	c.Request.URL.Path = c.Param("filepath")

	println(c.Request.URL.Path)
	http.FileServer(http.Dir("./img/")).ServeHTTP(c.Writer, c.Request)
}

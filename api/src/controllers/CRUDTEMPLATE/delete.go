/*
** D&GINE Project, 2026
** Backend
** File description:
** CRUDTEMPLATE/delete.go
 */

package CRUDTEMPLATE

import (
	"api/src/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/CRUDTEMPLATEs
// CRUDTEMPLATEs godoc
// @Summary Deletes given CRUDTEMPLATE
// @Schemes
// @Param id query string true "id of the CRUDTEMPLATE to delete"
// @Description ⚠️ Only accessible to admins ⚠️<br><br>Deletes the given CRUDTEMPLATE<br> <b>Careful, there is no turn back or check, once called, this route will delete the CRUDTEMPLATE no matter what</b> <br> It also deletes all content related to the CRUDTEMPLATE depending on the parameters of the database
// @Tags CRUDTEMPLATEs
// @Success 200
// @Router /api/CRUDTEMPLATEs/id [delete]
func CRUDTEMPLATEs(c *gin.Context) {
	idQuery := c.Request.URL.Query().Get("id")
	if idQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong parameters"})
		return
	}

	_, err := database.Db.Exec("DELETE FROM CRUDTEMPLATEs WHERE id=$1", idQuery)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

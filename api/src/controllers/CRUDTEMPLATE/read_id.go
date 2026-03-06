/*
** D&GINE Project, 2026
** API
** File description:
** CRUDTEMPLATE/read_id.go
 */

package CRUDTEMPLATE

import (
	"api/src/database"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/CRUDTEMPLATE
// CRUDTEMPLATE godoc
// @Summary Reads given CRUDTEMPLATE data
// @Schemes
// @Description ⚠️ Only accessible to admins ⚠️<br><br>Reads the given CRUDTEMPLATE's account data
// @Tags CRUDTEMPLATEs
// @Param id query string true "id of the CRUDTEMPLATE to read"
// @Produce json
// @Success 200 {object} DbReadResponse
// @Router /api/CRUDTEMPLATEs/id [get]
func ReadCRUDTEMPLATEId(c *gin.Context) {
	idQuery := c.Request.URL.Query().Get("id")
	if idQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong parameters"})
		return
	}

	var id sql.NullString
	err := database.Db.QueryRow("SELECT id FROM CRUDTEMPLATEs WHERE id=$1", idQuery).Scan(&id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error : " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, DbReadResponse{
		Id: id.String,
	})

}

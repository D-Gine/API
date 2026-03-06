/*
** D&GINE Project, 2026
** API
** File description:
** CRUDTEMPLATE/read.go
 */

package CRUDTEMPLATE

import (
	"api/src/database"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/CRUDTEMPLATEs
// CRUDTEMPLATEs godoc
// @Summary Reads all CRUDTEMPLATEs data
// @Schemes
// @Description ⚠️ Only accessible to admins ⚠️<br><br>Reads all CRUDTEMPLATEs data
// @Tags CRUDTEMPLATEs
// @Produce json
// @Success 200 {object} []DbReadResponse
// @Router /api/CRUDTEMPLATEs [get]
func ReadCRUDTEMPLATEs(c *gin.Context) {
	result := []DbReadResponse{}
	rows, err := database.Db.Query("SELECT id FROM CRUDTEMPLATEs")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	}
	defer rows.Close()
	for rows.Next() {
		var id sql.NullString
		rows.Scan(&id)
		result = append(result, DbReadResponse{
			Id: id.String,
		})
	}

	c.JSON(http.StatusOK, result)
}

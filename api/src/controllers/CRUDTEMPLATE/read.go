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
		return
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var id sql.NullString
		if err := rows.Scan(&id); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		result = append(result, DbReadResponse{
			Id: id.String,
		})
	}

	c.JSON(http.StatusOK, result)
}

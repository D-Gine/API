/*
** D&GINE Project, 2026
** API
** File description:
** CRUDTEMPLATE/update.go
 */

package CRUDTEMPLATE

import (
	"api/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateArgs struct {
	Name string `json:"name"`
}

// @BasePath /api/CRUDTEMPLATEs
// CRUDTEMPLATEs godoc
// @Summary Updates given CRUDTEMPLATE
// @Schemes
// @Description ⚠️ Only accessible to admins ⚠️<br><br>Updates given CRUDTEMPLATE with arguments in body
// @Tags CRUDTEMPLATEs
// @Accept json
// @Param id query string true "id of the CRUDTEMPLATE to update"
// @Param data body UpdateArgs true "New arguments to be set in the CRUDTEMPLATE given"
// @Success 200
// @Router /api/CRUDTEMPLATEs/id [put]
func UpdateCRUDTEMPLATEs(c *gin.Context) {
	idQuery := c.Request.URL.Query().Get("id")
	if idQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong parameters"})
		return
	}

	var args UpdateArgs

	// Parsing des args
	err := c.ShouldBindJSON(&args)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err = database.Db.Exec("UPDATE CRUDTEMPLATEs SET name=$2 WHERE id=$1", idQuery, args.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

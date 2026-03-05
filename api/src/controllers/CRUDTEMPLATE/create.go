/*
** D&GINE Project, 2026
** API
** File description:
** CRUDTEMPLATE/register.go
 */

package CRUDTEMPLATE

import (
	"database/sql"
	"net/http"

	"api/src/database"
	"api/src/internal/structs"

	"github.com/gin-gonic/gin"
)

type CreateCRUDTEMPLATEsArgs struct {
	Name string `json:"name"`
}

// @BasePath /api/CRUDTEMPLATEs
// CRUDTEMPLATEs godoc
// @Summary Creates a CRUDTEMPLATE
// @Schemes
// @Description ⚠️ Only accessible to admins ⚠️<br><br>Creates a CRUDTEMPLATE with arguments in body<br><br>Will return the id of created CRUDTEMPLATE
// @Tags CRUDTEMPLATEs
// @Accept json
// @Produce json
// @Param creds body CreateCRUDTEMPLATEsArgs true "CRUDTEMPLATE related informations that will be later needed for the login process"
// @Success 200 {object} structs.PostResponse
// @Router /api/CRUDTEMPLATEs [post]
func CreateCRUDTEMPLATEs(c *gin.Context) {
	var args CreateCRUDTEMPLATEsArgs

	// Body Json parsing
	if err := c.ShouldBindJSON(&args); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Checking db if CRUDTEMPLATE already exists
	rows, err := database.Db.Query("SELECT name FROM CRUDTEMPLATEs WHERE name=$1", args.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()
	if rows.Next() {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "CRUDTEMPLATE already exists"})
		return
	}

	// Inserting new CRUDTEMPLATE into db
	var id sql.NullString
	err = database.Db.QueryRow("INSERT INTO CRUDTEMPLATEs (name) VALUES ($1) RETURNING id", args.Name).Scan(&id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to register CRUDTEMPLATE"})
		return
	}

	c.JSON(http.StatusOK, structs.PostResponse{Id: id.String})
}

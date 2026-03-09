/*
** D&GINE Project, 2026
** API
** File description:
** rulesets/read.go
 */

package rulesets

import (
	"api/src/database"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/rulesets
// rulesets godoc
// @Summary Reads all rulesets data
// @Schemes
// @Description Reads all rulesets data
// @Tags rulesets
// @Produce json
// @Success 200 {object} []DbReadResponse
// @Router /api/rulesets [get]
func ReadRulesets(c *gin.Context) {
	result := []DbReadResponse{}
	rows, err := database.Db.Query("SELECT id, name FROM games.rulesets")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	}
	defer rows.Close()
	for rows.Next() {
		var id, name sql.NullString
		rows.Scan(&id, &name)
		result = append(result, DbReadResponse{
			Id:   id.String,
			Name: name.String,
		})
	}

	c.JSON(http.StatusOK, result)
}

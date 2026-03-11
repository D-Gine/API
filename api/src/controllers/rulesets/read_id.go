/*
** D&GINE Project, 2026
** API
** File description:
** rulesets/read_id.go
 */

package rulesets

import (
	"api/src/database"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/rulesets
// Ruleset godoc
// @Summary Reads given ruleset data
// @Schemes
// @Description Reads the given ruleset's account data
// @Tags rulesets
// @Param id query string true "id of the ruleset to read"
// @Produce json
// @Success 200 {object} DbReadResponse
// @Router /api/rulesets/id [get]
func ReadRulesetId(c *gin.Context) {
	idQuery := c.Param("id")
	var id, name sql.NullString
	err := database.Db.QueryRow("SELECT id, name FROM games.rulesets WHERE id=$1", idQuery).Scan(&id, &name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error : " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, DbReadResponse{
		Id:   id.String,
		Name: name.String,
	})
}
